package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/fentec-project/gofe/abe"
	"github.com/stretchr/testify/assert"
)

func TestFame(t *testing.T) {
	auth := abe.NewFAME()
	pk, sk, err := auth.GenerateMasterKeys()
	if err != nil {
		t.Fatalf("Failed to generate master keys: %v", err)
	}
	attributes := []string{"0", "1", "2", "3", "5"}
	ak, err := auth.GenerateAttribKeys(attributes, sk)
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v", err)
	}
	pub := CipherFamePub(pk)
	sub := CipherFameSub(ak)

	msg := "This is a test message"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"

	ciphertext, err := pub.EncryptAbe(policy, msg)
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := sub.DecryptAbe(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	assert.Equal(t, msg, plaintext)

}

func TestChaCha20(t *testing.T) {
	var secretKey [32]byte
	_, err := rand.Read(secretKey[:])
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	c := cipherChaChaPoly(secretKey)
	msg := "This is a test message"
	ciphertext, err := c.EncryptAead(0, nil, []byte(msg))
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := c.DecryptAead(0, nil, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	assert.Equal(t, msg, string(plaintext))
}

func TestAESGCM(t *testing.T) {
	var secretKey [32]byte
	_, err := rand.Read(secretKey[:])
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	c := cipherAESGCM(secretKey)
	msg := "This is a test message"
	ciphertext, err := c.EncryptAead(0, nil, []byte(msg))
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := c.DecryptAead(0, nil, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	assert.Equal(t, msg, string(plaintext))
}
