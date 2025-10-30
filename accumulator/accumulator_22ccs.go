package accumulator

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"math/big"
)

type Accumulator struct {
	capacity *big.Int
	g1       bn254.G1Affine
	g2       bn254.G2Affine
	// 使用 string (即 big.Int.String()) 作为 map 的键比 *big.Int 更安全
	g1_si map[string]*bn254.G1Affine
	g2_si map[string]*bn254.G2Affine
}

type AccumulatorCommitment struct {
	commitment bn254.G1Affine
}

type AccumulatorProof struct {
	proof bn254.G1Affine
}

// NewAccumulator where capacity is the maximum capacity of the accumulator.
func NewAccumulator(capacity *big.Int) *Accumulator {
	return &Accumulator{
		capacity: new(big.Int).Set(capacity),       // 初始化长度为 0
		g1_si:    make(map[string]*bn254.G1Affine), // 初始化 map
		g2_si:    make(map[string]*bn254.G2Affine), // 初始化 map
	}
}

func (a *Accumulator) SetUp() {
	// 1. 生成秘密值 s
	s, err := rand.Int(rand.Reader, ecc.BN254.ScalarField())
	if err != nil {
		panic(fmt.Sprintf("failed to generate random secret: %v", err))
	}

	// 2. 获取 G1 和 G2 的生成元
	_, _, a.g1, a.g2 = bn254.Generators()

	one := big.NewInt(1)
	// si 用于存储 s 的幂 (s^0, s^1, s^2, ...)
	si := new(big.Int).Set(one) // 初始化为 s^0 = 1

	// 3. 循环计算并存储 g1^(s^i) 和 g2^(s^i)
	// 循环条件: i <= a.capacity
	i := big.NewInt(0)
	limit := new(big.Int).Add(a.capacity, one)

	for i.Cmp(limit) < 0 {
		key := i.String()

		// 计算 g1^(s^i) 并存入 map
		g1si := new(bn254.G1Affine)
		g1si.ScalarMultiplication(&a.g1, si)
		a.g1_si[key] = g1si

		// 计算 g2^(s^i) 并存入 map
		g2si := new(bn254.G2Affine)
		g2si.ScalarMultiplication(&a.g2, si)
		a.g2_si[key] = g2si

		// 更新 si 以用于下一次迭代: si = si * s
		si.Mul(si, s)
		si.Mod(si, ecc.BN254.ScalarField())

		// 迭代器递增: i = i + 1
		i.Add(i, one)
	}
}

// Commit for x_i in elements, return g1^[(x_1+s)(x_2+s)...(x_|X|+s)]
func (a *Accumulator) Commit(elements []*big.Int) (*AccumulatorCommitment, error) {
	// 检查元素数量是否超过累加器容量
	numElements := big.NewInt(int64(len(elements)))
	if numElements.Cmp(a.capacity) > 0 {
		return nil, fmt.Errorf("number of elements (%d) exceeds capacity (%s)", len(elements), a.capacity.String())
	}

	// 1. 计算多项式 P(s) = product(x_i + s) 的系数
	// 从 P_0(s) = 1 开始, 其系数为 [1]
	coeffs := []*big.Int{big.NewInt(1)}
	fieldOrder := ecc.BN254.ScalarField()

	for _, x := range elements {
		// newCoeffs 用于存储 P_old(s) * (x+s) 的系数
		newCoeffs := make([]*big.Int, len(coeffs)+1)
		for i := range newCoeffs {
			newCoeffs[i] = new(big.Int)
		}

		// newCoeffs[0] = x * coeffs[0]
		newCoeffs[0].Mul(x, coeffs[0])
		newCoeffs[0].Mod(newCoeffs[0], fieldOrder)

		// 中间项: newCoeffs[j] = x*coeffs[j] + coeffs[j-1]
		for j := 1; j < len(coeffs); j++ {
			term1 := new(big.Int).Mul(x, coeffs[j])
			newCoeffs[j].Add(term1, coeffs[j-1])
			newCoeffs[j].Mod(newCoeffs[j], fieldOrder)
		}

		// 最后一项: newCoeffs[len(coeffs)] = coeffs[len(coeffs)-1]
		newCoeffs[len(coeffs)].Set(coeffs[len(coeffs)-1])

		coeffs = newCoeffs
	}

	// 2. 执行多标量乘法 (MSM): C = Σ c_i * g1^(s^i)
	// 使用雅可比坐标进行高效的加法运算
	var resultJac bn254.G1Jac
	var termAff bn254.G1Affine
	var termJac bn254.G1Jac

	for i, c := range coeffs {
		key := big.NewInt(int64(i)).String()
		precomputedPoint, ok := a.g1_si[key]
		if !ok {
			return nil, fmt.Errorf("setup not run for degree %d, accumulator not ready", i)
		}

		// 计算 c_i * g1^(s^i)
		termAff.ScalarMultiplication(precomputedPoint, c)

		// 将结果累加
		termJac.FromAffine(&termAff)
		resultJac.AddAssign(&termJac)
	}

	var finalCommitment bn254.G1Affine
	finalCommitment.FromJacobian(&resultJac)

	return &AccumulatorCommitment{commitment: finalCommitment}, nil
}

// Add g1^[(x_1+s)(x_2+s)...(x_|X|+s)(I+s)]
func (a *Accumulator) Add(Ax *AccumulatorCommitment, X []*big.Int, I *big.Int) (*AccumulatorCommitment, error) {
	// 1. 检查 I 是否已在集合 X 中
	for _, elem := range X {
		if elem.Cmp(I) == 0 {
			return nil, fmt.Errorf("element %s is already in set X", I.String())
		}
	}

	// 2. 检查增加新元素后是否会超出容量
	newSize := big.NewInt(int64(len(X) + 1))
	if newSize.Cmp(a.capacity) > 0 {
		return nil, fmt.Errorf("adding element would exceed accumulator capacity")
	}

	// 3. 创建新集合 X' = X U {I}
	newElements := make([]*big.Int, len(X)+1)
	copy(newElements, X)
	newElements[len(X)] = I

	// 4. 对新集合重新计算承诺值
	return a.Commit(newElements)
}

func (a *Accumulator) MemProve(elements []*big.Int, element *big.Int) (*big.Int, error) {
	return nil, nil
}

// {0, 1} ← Acc.NonMemVerifypp (AX , y, π y )
func (a *Accumulator) MemVerify(elements []*big.Int, element *big.Int, proof *AccumulatorProof) bool {
	return false
}
