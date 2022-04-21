package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/fentec-project/gofe/abe"
	"golang.org/x/crypto/chacha20poly1305"
)

//Cipher is ether an AEAD cipher or a CPABE scheme initialized with keys
type Cipher interface {
	EncryptAead(uint64, []byte, []byte) ([]byte, error)
	DecryptAead(uint64, []byte, []byte) ([]byte, error)
	EncryptAbe(string, string) ([]byte, error)
	DecryptAbe([]byte) (any, error)
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

func cipherChaChaPoly(key [32]byte) Cipher {
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

func cipherAESGCM(key [32]byte) Cipher {
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

func CipherFamePub(pk *abe.FAMEPubKey) Cipher {
	c := abe.NewFAME()
	return FameCipher{c, pk, &abe.FAMEAttribKeys{}}
}

func CipherFameSub(atk *abe.FAMEAttribKeys) Cipher {
	c := abe.NewFAME()
	return FameCipher{c, &abe.FAMEPubKey{}, atk}
}

func (c FameCipher) EncryptAead(n uint64, ad, plaintext []byte) ([]byte, error) {
	return nil, errors.New("FAME not supported")
}

func (c FameCipher) DecryptAead(n uint64, ad, ciphertext []byte) ([]byte, error) {
	return nil, errors.New("FAME not supported")
}

func (c FameCipher) EncryptAbe(policy, msg string) ([]byte, error) {
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

func (c FameCipher) DecryptAbe(ciphertext []byte) (any, error) {
	var ct *abe.FAMECipher
	err := json.Unmarshal(ciphertext, &ct)
	if err != nil {
		return nil, err
	}
	msg, err := c.FAME.Decrypt(ct, c.attribKey, nil)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c AeadCipher) EncryptAead(n uint64, ad, plaintext []byte) ([]byte, error) {
	return c.Seal(nil, c.nonce(n), plaintext, ad), nil
}

func (c AeadCipher) DecryptAead(n uint64, ad, ciphertext []byte) ([]byte, error) {
	return c.Open(nil, c.nonce(n), ciphertext, ad)
}

func (c AeadCipher) EncryptAbe(policy, msg string) ([]byte, error) {
	return nil, errors.New("AeadCipher not supported")
}

func (c AeadCipher) DecryptAbe(ciphertext []byte) (any, error) {
	return nil, errors.New("AeadCipher not supported")
}
