package bgw05

import (
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"math/big"
)

type MasterPublicKey struct {
	Bound            int
	G1ExpAlphaPowers []bn254.G1Affine
	G2ExpAlphaPowers []bn254.G2Affine
	G1ExpGamma       fr.Element
}

type MasterSecretKey struct {
	Gamma fr.Element
}

type UserIdentity struct {
	Identity int
}

type UserSecretKey struct {
	SecretKey bn254.G2Affine
}

type Message struct {
	Message bn254.GT
}

type Ciphertext struct {
	C1 bn254.G1Affine
	C2 bn254.G1Affine
	C3 bn254.GT
}

func Setup(n int) (*MasterPublicKey, *MasterSecretKey, error) {
	if n < 1 {
		return nil, nil, fmt.Errorf("invalid parameters")
	}
	alpha, err := new(fr.Element).SetRandom()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to setup master public key")
	}
	gamma, err := new(fr.Element).SetRandom()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to setup master public key")
	}

	g1ExpGamma := new(bn254.G1Affine).ScalarMultiplicationBase(gamma.BigInt(new(big.Int)))
}

func Extract(i *UserIdentity, msk *MasterSecretKey) (*UserSecretKey, error) {

}

func Encrypt(s []*UserIdentity, mpk *MasterPublicKey, m *Message) (*Ciphertext, error) {}

func Decrypt(i *UserIdentity, sk *UserSecretKey, s []*UserIdentity, mpk *MasterPublicKey, c *Ciphertext) (*Message, error) {

}
