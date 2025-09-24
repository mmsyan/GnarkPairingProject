package IBE

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"math/big"
)

type IBEParams struct {
	MasterKey   *big.Int
	PublicKeyG  bn254.G1Affine
	PublicKeyGx bn254.G1Affine
}

func SetUp() (*IBEParams, error) {
	x, err := rand.Int(rand.Reader, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to set up")
	}
	_, _, g1, _ := bn254.Generators()
	var g1x bn254.G1Affine
	g1x.ScalarMultiplication(&g1x, x)
	return &IBEParams{
		MasterKey:   x,
		PublicKeyG:  g1,
		PublicKeyGx: g1,
	}, nil
}
