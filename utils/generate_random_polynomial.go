package utils

import (
	"crypto/rand"
	"github.com/consensys/gnark-crypto/ecc"
	"math/big"
)

// GenerateRandomPolynomial 生成一个次数最高为 (degree - 1) 的多项式的系数列表。
//
// 注意：多项式的系数是按低次到高次的顺序排列：{a_0, a_1, ..., a_{degree-1}}
// degree:   多项式系数列表的长度（即最高次数 + 1）。
// constantTerm: 多项式的常数项系数 a_0。
// 返回值:   一个 []*big.Int 数组，表示多项式的系数。
func GenerateRandomPolynomial(degree int, constantTerm *big.Int) []*big.Int {
	q := ecc.BN254.ScalarField()
	if degree <= 0 {
		return []*big.Int{}
	}
	coefficients := make([]*big.Int, degree)
	coefficients[0] = new(big.Int).Set(constantTerm)
	for i := 1; i < degree; i++ {
		randomCoef, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		coefficients[i] = randomCoef
	}
	return coefficients
}
