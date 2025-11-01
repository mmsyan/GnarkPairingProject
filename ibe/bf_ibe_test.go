package ibe

import "testing"

func TestIBECorrect(t *testing.T) {
	ibeParams, err := SetUp()
	if err != nil {
		t.Fatal(err)
	}

	sk, err := KeyGenerate(*ibeParams, "ChenBerry")
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, err := Encrypt(*ibeParams, "ChenBerry", []byte("Hello World"))
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(*ciphertext, sk)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != "Hello World" {
		t.Fatalf("decrypted wrong, %s", string(decrypted))
	}
}

func TestIBEError(t *testing.T) {
	ibeParams, err := SetUp()
	if err != nil {
		t.Fatal(err)
	}

	sk, err := KeyGenerate(*ibeParams, "ChenBerry")
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, err := Encrypt(*ibeParams, "ErrorChenBerry", []byte("Hello World"))
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(*ciphertext, sk)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) == "Hello World" {
		t.Fatal("decrypted wrong")
	}
}
