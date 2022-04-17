package mqtts_unitn

import (
	"testing"
)

func TestKeygen(t *testing.T) {
	a := newCPABE(false, "")
	a.GenerateAndSaveMasterKeys()
}

/*
func TestKeyLoader(t *testing.T) {

}

func TestFAME(t *testing.T) {
	a := newCPABE(false, "")
	pk, sk, err := a.PubKeygen()
	if err != nil {
		t.Fatalf("Failed to generate master keys: %v", err)
	}

	msg := "This is a test message"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"
	ciphertext, err := a.EncryptEncode(pk, policy, msg)
	if err != nil {
		t.Fatalf("Encryption failure: %v", err)
	}

	gamma := []string{"0", "1", "2", "3", "5"}
	//keys, err := a.PrivKeygen(sk, gamma)

	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v", err)
	}

	plaintext, err := a.DecryptDecode(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	fmt.Println(&plaintext)

}

func keygen() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func TestAEAD(t *testing.T) {
	key, err := keygen()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	c := cipherChaChaPoly(key)
	msg := "This is a test message a mammt"

	ciphertext := c.Encrypt(nil, 0, nil, []byte(msg))
	plaintext, err := c.Decrypt(nil, 0, nil, ciphertext)

	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	fmt.Println(string(plaintext))
}
*/
