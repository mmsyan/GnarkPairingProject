package fibe

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"reflect"
	"testing"
)

func TestIBECorrect(t *testing.T) {
	fibeInstance := NewFIBE(10, 3)
	fibeInstance.SetUp()

	secretKey, err := fibeInstance.KeyGenerate([]int{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}

	message := bn254.GT{}
	message.SetRandom()
	fmt.Print(message)
	fmt.Println()

	ciphertext, err := fibeInstance.Encrypt([]int{2, 3, 4, 5}, message)
	if err != nil {
		t.Error(err)
	}

	decryptedMessage, err := fibeInstance.Decrypt(secretKey, ciphertext)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(decryptedMessage)

	if !reflect.DeepEqual(message, decryptedMessage) {
		t.Error("decrypted message does not match original message")
	}
}
