package ibe

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/mmsyan/GnarkPairingProject/utils"
	"math/big"
)

type IBEParams struct {
	MasterKey   *big.Int
	PublicKeyG  bn254.G1Affine
	PublicKeyGx bn254.G1Affine
	DST         []byte
}

type IBECiphertext struct {
	C1 bn254.G1Affine
	C2 []byte
}

func SetUp() (*IBEParams, error) {
	// x <- Zq
	x, err := rand.Int(rand.Reader, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to set up")
	}

	// g <- G1
	// g^x in G1
	_, _, g, _ := bn254.Generators()
	var gx bn254.G1Affine
	gx.ScalarMultiplication(&g, x)
	return &IBEParams{
		MasterKey:   x,
		PublicKeyG:  g,
		PublicKeyGx: gx,
		DST:         []byte("ibe Encryption"),
	}, nil
}

func KeyGenerate(ibeParams IBEParams, id string) (bn254.G2Affine, error) {
	// qid = hashToCurve(id) in G2
	qid, err := bn254.HashToG2([]byte(id), ibeParams.DST)
	if err != nil {
		return bn254.G2Affine{}, fmt.Errorf("failed to generate key")
	}

	// sk = qid^x in G2
	var sk bn254.G2Affine
	sk.ScalarMultiplication(&qid, ibeParams.MasterKey)
	return sk, nil
}

func Encrypt(ibeParams IBEParams, id string, message []byte) (*IBECiphertext, error) {
	// qid = hashToCurve(id) in G2
	qid, err := bn254.HashToG2([]byte(id), ibeParams.DST)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}

	// r <- Zq
	r, err := rand.Int(rand.Reader, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}

	//gid = e(g^x, qid)^r
	eGxQid, err := bn254.Pair([]bn254.G1Affine{ibeParams.PublicKeyGx}, []bn254.G2Affine{qid})
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}
	gid := *(new(bn254.GT).Exp(eGxQid, r))
	fmt.Printf("gid: %v\n", gid)

	//// qidR = qid^r in G2
	//var qidR bn254.G2Affine
	//qidR.ScalarMultiplication(&qid, r)
	//// gid = e(g^x, qid^r)
	//gid, err := bn254.Pair([]bn254.G1Affine{ibeParams.PublicKeyGx}, []bn254.G2Affine{qidR})
	//if err != nil {
	//	return nil, fmt.Errorf("failed to encrypt message")
	//}
	//fmt.Printf("encrypt gid: %v\n", gid)

	// c1 = g^r
	var c1 bn254.G1Affine
	c1.ScalarMultiplication(&ibeParams.PublicKeyG, r)

	gidBytes := utils.Hash2(gid)
	c2 := utils.Xor(message, gidBytes)

	return &IBECiphertext{
		C1: c1,
		C2: c2,
	}, nil
}

func Decrypt(ciphertext IBECiphertext, secretKey bn254.G2Affine) ([]byte, error) {
	// gid = e(c1, sk) = e(g^r, qid^x)
	gid, err := bn254.Pair([]bn254.G1Affine{ciphertext.C1}, []bn254.G2Affine{secretKey})
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message")
	}
	fmt.Printf("decrypt gid: %v\n", gid)

	gidBytes := utils.Hash2(gid)
	return utils.Xor(ciphertext.C2, gidBytes), nil
}
