package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"

	"github.com/fentec-project/gofe/abe"
	"golang.org/x/crypto/chacha20poly1305"
)

//Cipher is either an AEAD cipher or a CPABE scheme initialized with keys
type Cipher interface {
	Encrypt(uint64, string, string) ([]byte, error)
	Decrypt(uint64, []byte, []byte) (string, error)
	Keygen() error
}

/*
Implementing a new CP-ABE scheme should be easy.
Create the corresponding struct and implement the methods.

TODO: Check KP-ABE schemes
*/

type CipherFunc interface {
	Cipher(key [32]byte) Cipher
	CipherName() string
}

//FameCipher represent the FAME CPABE scheme
type FameCipher struct {
	*abe.FAME
	publicKey *abe.FAMEPubKey
	attribKey *abe.FAMEAttribKeys
}

//AeadCipher is an AEAD cipher
type AeadCipher struct {
	cipher.AEAD
	nonce func(uint64) []byte
}

func CipherChaChaPoly(key [32]byte) Cipher {
	c, err := chacha20poly1305.New(key[:])
	if err != nil {
		panic(err)
	}
	return AeadCipher{
		c,
		func(u uint64) []byte {
			var nonce [12]byte
			binary.LittleEndian.PutUint64(nonce[4:], u)
			return nonce[:]
		}}
}

func CipherAESGCM(key [32]byte) Cipher {
	c, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}
	return AeadCipher{
		gcm,
		func(n uint64) []byte {
			var nonce [12]byte
			binary.BigEndian.PutUint64(nonce[4:], n)
			return nonce[:]
		},
	}
}

func CipherFame() FameCipher {
	return FameCipher{abe.NewFAME(), &abe.FAMEPubKey{}, &abe.FAMEAttribKeys{}}
}

func (c *FameCipher) FameKeygen(attributes []string) (*abe.FAMEPubKey, *abe.FAMEAttribKeys) {
	pk, sk, _ := c.FAME.GenerateMasterKeys()
	ak, _ := c.FAME.GenerateAttribKeys(attributes, sk)
	return pk, ak
}

func (c FameCipher) Keygen() error {
	attributes := []string{"0", "1", "2", "3", "5"}
	pk, ak := c.FameKeygen(attributes)
	c.publicKey = pk
	c.attribKey = ak
	return nil
}

func (c *FameCipher) setPubKey(pk *abe.FAMEPubKey) {
	c.publicKey = pk
}

func (c *FameCipher) setAttribKey(ak *abe.FAMEAttribKeys) {
	c.attribKey = ak
}

// Encrypts a message MSG with a POLICY using the FAME CP-ABE scheme
func (c FameCipher) Encrypt(u uint64, policy, msg string) ([]byte, error) {
	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		return nil, err
	}
	enc, err := c.FAME.Encrypt(msg, msp, c.publicKey)
	if err != nil {
		return nil, err
	}
	cipertext, err := json.Marshal(enc)
	if err != nil {
		return nil, err
	}
	return cipertext, nil
}

// Decrypts a ciphertext generated using the FAME CP-ABE scheme
func (c FameCipher) Decrypt(u uint64, ad, ciphertext []byte) (string, error) {
	var ct *abe.FAMECipher
	err := json.Unmarshal(ciphertext, &ct)
	if err != nil {
		return "", err
	}
	msg, err := c.FAME.Decrypt(ct, c.attribKey, nil)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (c AeadCipher) Keygen() error {
	return nil
}

func (c AeadCipher) Encrypt(n uint64, ad, plaintext string) ([]byte, error) {
	return c.Seal(nil, c.nonce(n), []byte(plaintext), []byte(ad)), nil
}

func (c AeadCipher) Decrypt(n uint64, ad, ciphertext []byte) (string, error) {
	plaintext, err := c.Open(nil, c.nonce(n), ciphertext, ad)
	return string(plaintext), err
}
