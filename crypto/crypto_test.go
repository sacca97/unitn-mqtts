package crypto

import (
	"fmt"
	"testing"
)

func TestKeygen(t *testing.T) {
	a := NewCPABE()
	a.GenerateAndSaveMasterKeys()
}

func TestKeyLoader(t *testing.T) {

}

func TestFAME_OLD(t *testing.T) {
	a := NewCPABE()
	pk, sk, err := a.PubKeygen()
	a.SetPublicKey(pk)
	if err != nil {
		t.Fatalf("Failed to generate master keys: %v", err)
	}

	msg := "This is a test message"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"
	ciphertext, err := a.EncryptEncode(policy, msg)
	fmt.Println(len(msg))
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}

	gamma := []string{"0", "1", "2", "3", "5"}
	attr, err := a.PrivKeygen(sk, gamma)
	a.SetAttribKey(attr)

	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v", err)
	}

	plaintext, err := a.DecryptDecode(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	fmt.Println(plaintext)

}
