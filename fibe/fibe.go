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
	q        *big.Int
	msk_ti   []*big.Int
	msk_y    *big.Int
	pk_Ti    []*bn254.G2Affine
	pk_Y     bn254.GT
}

type FIBESecretKey struct {
	userAttributes []int
	di             map[int]*bn254.G1Affine
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
	fibe.q = ecc.BN254.ScalarField()
	_, _, fibe.g1, fibe.g2 = bn254.Generators()
	var err error
	for i := 1; i <= fibe.universe; i++ {
		fibe.msk_ti[i], err = rand.Int(rand.Reader, fibe.q)                          // ti <- Zq
		fibe.pk_Ti[i] = (&bn254.G2Affine{}).ScalarMultiplicationBase(fibe.msk_ti[i]) // Ti = g2^ti
	}
	fibe.msk_y, err = rand.Int(rand.Reader, fibe.q) // y <- Zq
	eG1G2, err := bn254.Pair([]bn254.G1Affine{fibe.g1}, []bn254.G2Affine{fibe.g2})
	fibe.pk_Y = *((new(bn254.GT)).Exp(eG1G2, fibe.msk_y)) // Y = e(g1, g2)^y

	if err != nil {
		panic(err)
	}

}

func (fibe *FIBE) KeyGenerate(userAttributes []int) (*FIBESecretKey, error) {
	di := make(map[int]*bn254.G1Affine)
	polynomial := utils.GenerateRandomPolynomial(fibe.distance, fibe.msk_y)
	for _, i := range userAttributes {
		qi := utils.ComputePolynomialValue(polynomial, new(big.Int).SetInt64(int64(i)))
		// 在有限域 F_q 内计算除法：qiDivTi = qi * (msk_ti[i])^{-1} mod q
		tiInverse := new(big.Int).ModInverse(fibe.msk_ti[i], fibe.q)
		if tiInverse == nil {
			return nil, fmt.Errorf("failed to compute modular inverse for msk_ti[%d]", i)
		}
		qiDivTi := new(big.Int).Mul(qi, tiInverse)
		qiDivTi.Mod(qiDivTi, fibe.q)
		// Di = g1^(q(i)/ti)
		//Di := fibe.g1.ScalarMultiplicationBase(qiDivTi)
		Di := (&bn254.G1Affine{}).ScalarMultiplicationBase(qiDivTi)
		di[i] = Di
	}
	return &FIBESecretKey{
		userAttributes: userAttributes,
		di:             di,
	}, nil
}

func (fibe *FIBE) Encrypt(messageAttributes []int, message bn254.GT) (*FIBECiphertext, error) {
	s, err := rand.Int(rand.Reader, fibe.q)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}
	egg_ys := *(new(bn254.GT)).Exp(fibe.pk_Y, s)

	// e' = message * Y^s = message * (e(g1, g2)^y)^s
	ePrime := *(message.Mul(&message, &egg_ys))

	// ei = Ti^s = (g2^ti)^s
	ei := map[int]*bn254.G2Affine{}
	for _, i := range messageAttributes {
		//ei[i] = fibe.pk_Ti[i].ScalarMultiplicationBase(s)
		ei[i] = (&bn254.G2Affine{}).ScalarMultiplication(fibe.pk_Ti[i], s)
	}
	return &FIBECiphertext{
		messageAttributes: messageAttributes,
		ePrime:            ePrime,
		ei:                ei,
	}, nil

}

func (fibe *FIBE) Decrypt(userSecretKey *FIBESecretKey, ciphertext *FIBECiphertext) (bn254.GT, error) {
	s := utils.FindCommonAttributes(userSecretKey.userAttributes, ciphertext.messageAttributes, fibe.distance)
	if s == nil {
		return bn254.GT{}, fmt.Errorf("failed to find enough common attributes")
	}
	denominator := bn254.GT{}
	denominator.SetOne()
	for _, i := range s {
		di := userSecretKey.di[i]
		ei := ciphertext.ei[i]
		// e(Di, Ei)
		eDiEi, err := bn254.Pair([]bn254.G1Affine{*di}, []bn254.G2Affine{*ei})
		if err != nil {
			return bn254.GT{}, fmt.Errorf("failed to decrypt message")
		}
		delta := utils.ComputeLagrangeBasis(i, s, 0)
		eDiEiDelta := (new(bn254.GT)).Exp(eDiEi, delta)
		denominator.Mul(&denominator, eDiEiDelta)
	}

	decryptedMessage := ciphertext.ePrime.Div(&ciphertext.ePrime, &denominator)
	return *decryptedMessage, nil
}
