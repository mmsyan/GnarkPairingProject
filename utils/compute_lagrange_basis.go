package utils

import (
	"github.com/consensys/gnark-crypto/ecc"
	"math/big"
)

// ComputeLagrangeBasis 计算拉格朗日基函数在 x 处的值：Delta_{i, S}(x) mod q
func ComputeLagrangeBasis(i int, s []int, x int) *big.Int {
	// q: 有限域的阶 (ecc.BN254.ScalarField())
	q := ecc.BN254.ScalarField()

	iElement := big.NewInt(int64(i))
	xElement := big.NewInt(int64(x))
	delta := big.NewInt(1)

	// 临时 big.Int 变量
	temp := new(big.Int)

	// 确保 i 和 x 都在 Zq 范围内，虽然在这个场景下通常是小的索引
	iElement.Mod(iElement, q)
	xElement.Mod(xElement, q)

	for _, j := range s {
		if i != j {
			jElement := big.NewInt(int64(j))
			jElement.Mod(jElement, q)

			// 1. 计算 分子: (x - j) mod q
			// numerator = (x - j) mod q
			numerator := temp.Sub(xElement, jElement)
			numerator.Mod(numerator, q)

			// 2. 计算 分母: (i - j) mod q
			// denominator = (i - j) mod q
			denominator := new(big.Int).Sub(iElement, jElement)
			denominator.Mod(denominator, q)

			// 3. 计算 模逆: (i - j)^-1 mod q
			// invDenominator = (i - j)^-1 mod q
			invDenominator := new(big.Int).ModInverse(denominator, q)
			if invDenominator == nil {
				// 错误处理：模逆不存在 (虽然对于素数阶域，如果 i != j 不太可能发生)
				return nil
			}

			// 4. 计算分数: (x - j) * (i - j)^-1 mod q
			// fraction = numerator * invDenominator mod q
			fraction := temp.Mul(numerator, invDenominator)
			fraction.Mod(fraction, q)

			// 5. 更新 delta: delta = delta * fraction mod q
			delta.Mul(delta, fraction)
			delta.Mod(delta, q)
		}
	}

	return delta
}
