package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFame(t *testing.T) {
	pub := CipherFame()
	sub := CipherFame() //cipherFame()
	pk, sk, err := pub.FAME.GenerateMasterKeys()
	if err != nil {
		t.Fatalf("Failed to generate master keys: %v", err)
	}
	attributes := []string{"0", "1", "2", "3", "5"}
	ak, err := pub.FAME.GenerateAttribKeys(attributes, sk)
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v", err)
	}
	pub.setPubKey(pk)
	sub.setAttribKey(ak)

	msg := "This is a test message"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"

	ciphertext, err := pub.Encrypt(0, policy, msg)
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := sub.Decrypt(0, nil, ciphertext)
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
	c := CipherChaChaPoly(secretKey)
	msg := "This is a test message"
	ciphertext, err := c.Encrypt(0, "", msg)
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := c.Decrypt(0, nil, ciphertext)
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
	c := CipherAESGCM(secretKey)
	msg := "This is a test message"
	ciphertext, err := c.Encrypt(1234, "simoladoveilsugo", msg)
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}
	plaintext, err := c.Decrypt(1234, []byte("simoladoveilsugo"), ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	assert.Equal(t, msg, plaintext)
}
