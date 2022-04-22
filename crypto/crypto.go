package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/fentec-project/gofe/abe"
)

type CipherFunc interface {
	Cipher(k [32]byte) Cipher
	CipherName() string
}

type aeadCipher struct {
	cipher.AEAD
	nonce func(uint64) []byte
}

func (c Cpabe) GenerateAndSaveMasterKeys() error {
	pk, sk, err := c.PubKeygen()
	if err != nil {
		return err
	}
	bytesPublic, _ := json.Marshal(pk)
	bytesSecret, _ := json.Marshal(sk)

	pub := &pem.Block{
		Type:  "FAME PUBLIC KEY",
		Bytes: bytesPublic,
	}
	sec := &pem.Block{
		Type:  "FAME SECRET KEY",
		Bytes: bytesSecret,
	}
	pemPublic, err := os.Create("public.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := pem.Encode(pemPublic, pub); err != nil {
		panic(err)
	}

	pemSecret, err := os.Create("secret.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := pem.Encode(pemSecret, sec); err != nil {
		panic(err)
	}

	pemPublic.Close()
	pemSecret.Close()

	return nil
}

func (c Cpabe) PubKeygen() (*abe.FAMEPubKey, *abe.FAMESecKey, error) {
	pk, sk, err := c.scheme.GenerateMasterKeys()
	if err != nil {
		return nil, nil, err
	}
	return pk, sk, err
}

func (c Cpabe) PrivKeygen(sk *abe.FAMESecKey, attributes []string) (*abe.FAMEAttribKeys, error) {
	ak, err := c.scheme.GenerateAttribKeys(attributes, sk)
	if err != nil {
		return nil, err
	}
	return ak, nil
}

type PublicKey abe.FAMEPubKey
type AttrKey abe.FAMEAttribKeys
type SecretKey abe.FAMESecKey

func MarshalPubKey(k abe.FAMEPubKey) ([]byte, error) {
	return json.Marshal(k)
}

func MarshalAttribKey(k abe.FAMEAttribKeys) ([]byte, error) {
	return json.Marshal(k)
}

func UnmarshalFamePubKey(data []byte) (*PublicKey, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("nil byte array")
	}
	var k *PublicKey
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}
	return k, err
}

//Decode a FameKey from a byte array
func UnmarshalFameAttrKey(data []byte) (*AttrKey, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("nil byte array")
	}
	var k *AttrKey
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}
	return k, err
}

type Cpabe struct {
	scheme *abe.FAME
	pubKey *abe.FAMEPubKey
	attKey *abe.FAMEAttribKeys
}

func NewCPABE() *Cpabe {
	return &Cpabe{
		scheme: abe.NewFAME(),
	}
}

func (c *Cpabe) GetKeys() {
	fmt.Println("Public Key:")
	fmt.Println(c.pubKey)
	fmt.Println("Attribute Key:")
	fmt.Println(c.attKey)
}

func (c *Cpabe) SetPublicKey(key *abe.FAMEPubKey) {
	c.pubKey = key
}

func (c *Cpabe) SetAttribKey(key *abe.FAMEAttribKeys) {
	c.attKey = key
}

// FAME encryption function wrapper
func (c Cpabe) Encrypt(msg string, msp *abe.MSP) (*abe.FAMECipher, error) {
	return c.scheme.Encrypt(msg, msp, c.pubKey)
}

// FAME decryption function wrapper
func (c Cpabe) Decrypt(ciphertext *abe.FAMECipher) (string, error) {
	return c.scheme.Decrypt(ciphertext, c.attKey, nil)
}

// Encrypt and encode
func (c Cpabe) EncryptEncode(policy string, plaintext string) ([]byte, error) {
	msp, err := abe.BooleanToMSP(policy, false)
	if err != nil {
		return nil, err
	}

	ct, err := c.Encrypt(plaintext, msp)
	if err != nil {
		return nil, err
	}

	ciphertext, err := json.Marshal(ct)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func loadKey(path string) (any, error) {
	//load pk
	return nil, nil
}

func (c Cpabe) DecryptDecode(ciphertext []byte) (string, error) {
	var ct *abe.FAMECipher
	err := json.Unmarshal(ciphertext, &ct)
	if err != nil {
		return "", fmt.Errorf("error decoding ciphertext: %s", err)
	}
	pt, err := c.Decrypt(ct)
	if err != nil {
		return "", err
	}
	return pt, nil
}

func SymKeygen() [32]byte {
	var k [32]byte
	_, err := rand.Read(k[:])
	if err != nil {
		panic(err)
	}
	return k
}
