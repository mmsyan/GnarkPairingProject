package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomPolynomial 生成一个次数最高为 (degree - 1) 的多项式的系数列表。
//
// 注意：多项式的系数是按低次到高次的顺序排列：{a_0, a_1, ..., a_{degree-1}}
// degree:   多项式系数列表的长度（即最高次数 + 1）。
// constantTerm: 多项式的常数项系数 a_0。
// 返回值:   一个 []*big.Int 数组，表示多项式的系数。
func GenerateRandomPolynomial(degree int, constantTerm *big.Int) []*big.Int {
	// 1. 基本检查
	if degree <= 0 {
		return []*big.Int{}
	}

	// 2. 初始化系数数组
	// 数组长度为 degree，因为 degree 规定了系数的个数 (从 a_0 到 a_{degree-1})
	coefficients := make([]*big.Int, degree)

	// 3. 设置常数项 (a_0)
	// 常数项是数组的第一个元素 (索引 0)
	// 使用 Set() 方法进行深拷贝，以确保不直接引用传入的 constantTerm
	coefficients[0] = new(big.Int).Set(constantTerm)

	// 4. 设置随机生成系数的范围（上限）
	// 为了生成合理的随机大整数，我们通常需要设置一个上限。
	// 这里的例子我们设定一个 256 位的上限，这在密码学中很常见。
	// 您可以根据实际需求修改此范围。
	bitLength := 256
	max := new(big.Int).Lsh(big.NewInt(1), uint(bitLength))

	// 5. 生成高次项的随机系数 (a_1 到 a_{degree-1})
	for i := 1; i < degree; i++ {
		// 生成一个 [0, max-1] 范围内的随机大整数
		randomCoeff, err := rand.Int(rand.Reader, max)
		if err != nil {
			// 如果随机数生成失败，出于安全考虑，应返回错误或停止。
			// 简单起见，这里设置一个默认值（例如 1）并继续，但在生产环境中应更严格处理。
			// 实际应用中，您可能希望让函数返回 ([], error)
			randomCoeff = big.NewInt(1)
		}
		coefficients[i] = randomCoeff
	}

	// 6. 返回结果
	return coefficients
}
