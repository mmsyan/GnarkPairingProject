package fibe

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/mmsyan/GnarkPairingProject/utils"
	"math/big"
)

type FIBE struct {
	universe int
	distance int
	g1       bn254.G1Affine
	g2       bn254.G2Affine
	msk_ti   []*big.Int
	msk_y    *big.Int
	pk_Ti    []*bn254.G2Affine
	pk_Y     bn254.GT
}

type FIBESecretKey struct {
	di map[int]*bn254.G1Affine
}

type FIBECiphertext struct {
	messageAttributes []int
	ePrime            bn254.GT
	ei                map[int]*bn254.G2Affine
}

func NewFIBE(universe int, distance int) *FIBE {
	// 使用 &FIBE{} 语法创建一个结构体实例并返回其指针。
	return &FIBE{
		universe: universe,
		distance: distance,
		msk_ti:   make([]*big.Int, universe+1),
		pk_Ti:    make([]*bn254.G2Affine, universe+1),
	}
}

func (fibe *FIBE) SetUp() {
	_, _, fibe.g1, fibe.g2 = bn254.Generators()
	var err error
	for i := 0; i < fibe.universe; i++ {
		fibe.msk_ti[i], err = rand.Int(rand.Reader, ecc.BN254.ScalarField())
		fibe.pk_Ti[i] = (fibe.pk_Ti[i]).ScalarMultiplication(&fibe.g2, fibe.msk_ti[i])
	}
	fibe.msk_y, err = rand.Int(rand.Reader, ecc.BN254.ScalarField())
	fibe.pk_Y, err = bn254.Pair([]bn254.G1Affine{fibe.g1}, []bn254.G2Affine{fibe.g2})
	if err != nil {
		panic(err)
	}

}

func (fibe *FIBE) KeyGenerate(userAttributes []int) (*FIBESecretKey, error) {
	di := make(map[int]*bn254.G1Affine)
	q := utils.GenerateRandomPolynomial(fibe.distance, fibe.msk_y)
	for _, i := range userAttributes {
		qi := utils.ComputePolynomialValue(q, new(big.Int).SetInt64(int64(i)))
		qiDivt := qi.Div(&qi, fibe.msk_ti[i])
		Di := fibe.g1.ScalarMultiplicationBase(qiDivt)
		di[i] = Di
	}
	return &FIBESecretKey{
		di: di,
	}, nil
}

func (fibe *FIBE) Encrypt(messageAttributes []int, message bn254.GT) (*FIBECiphertext, error) {
	s, err := rand.Int(rand.Reader, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}
	y, err := bn254.Pair([]bn254.G1Affine{fibe.g1}, []bn254.G2Affine{fibe.g2})
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}
	ys := (new(bn254.GT)).Exp(y, s)
	ePrime := message.Mul(&message, ys)

	ei := map[int]*bn254.G2Affine{}
	for _, i := range messageAttributes {
		ei[i] = fibe.pk_Ti[i].ScalarMultiplicationBase(s)
	}

	return &FIBECiphertext{
		messageAttributes: messageAttributes,
		ePrime:            *ePrime,
		ei:                ei,
	}, nil

}

func (fibe *FIBE) Decrypt() {}
