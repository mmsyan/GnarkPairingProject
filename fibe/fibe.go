package fibe

import (
	"crypto/rand"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"math/big"
)

type FIBE struct {
	universe int
	distance int
	g1       bn254.G1Affine
	g2       bn254.G2Affine
	msk_ti   []*big.Int
	msk_y    *big.Int
	pk_Ti    []*bn254.G1Affine
	pk_Y     bn254.GT
}

func NewFIBE(universe int, distance int) *FIBE {
	// 使用 &FIBE{} 语法创建一个结构体实例并返回其指针。
	return &FIBE{
		universe: universe,
		distance: distance,
		msk_ti:   make([]*big.Int, universe+1),
		pk_Ti:    make([]*bn254.G1Affine, universe+1),
	}
}

func (fibe *FIBE) SetUp() {
	_, _, fibe.g1, fibe.g2 = bn254.Generators()
	var err error
	for i := 0; i < fibe.universe; i++ {
		fibe.msk_ti[i], err = rand.Int(rand.Reader, ecc.BN254.ScalarField())
		fibe.pk_Ti[i] = (fibe.pk_Ti[i]).ScalarMultiplication(&fibe.g1, fibe.msk_ti[i])
	}
	fibe.msk_y, err = rand.Int(rand.Reader, ecc.BN254.ScalarField())
	fibe.pk_Y, err = bn254.Pair([]bn254.G1Affine{fibe.g1}, []bn254.G2Affine{fibe.g2})
	if err != nil {
		panic(err)
	}

}

func (fibe *FIBE) KeyGenerate() {}

func (fibe *FIBE) Encrypt() {}

func (fibe *FIBE) Decrypt() {}
