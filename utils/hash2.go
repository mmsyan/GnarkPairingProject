package utils

import "github.com/consensys/gnark-crypto/ecc/bn254"

func Hash2(gt bn254.GT) []byte {
	gtBytes := gt.Bytes()
	return gtBytes[:]
}
