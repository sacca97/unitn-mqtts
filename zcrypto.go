package mqtts_unitn

import (
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/fentec-project/gofe/abe"
	"golang.org/x/crypto/chacha20poly1305"
)

type CipherFunc interface {
	Cipher(k [32]byte) Cipher
	CipherName() string
}

// A Cipher is a AEAD cipher that has been initialized with a key.
type Cipher interface {
	// Encrypt encrypts the provided plaintext with a nonce and then appends the
	// ciphertext to out along with an authentication tag over the ciphertext
	// and optional authenticated data.
	Encrypt(out []byte, n uint64, ad, plaintext []byte) []byte

	// Decrypt authenticates the ciphertext and optional authenticated data and
	// then decrypts the provided ciphertext using the provided nonce and
	// appends it to out.
	Decrypt(out []byte, n uint64, ad, ciphertext []byte) ([]byte, error)
}

func cipherChaChaPoly(key []byte) Cipher {
	c, err := chacha20poly1305.New(key)
	if err != nil {
		panic(err)
	}
	return aeadCipher{
		c,
		func(u uint64) []byte {
			var nonce [12]byte
			binary.LittleEndian.PutUint64(nonce[4:], u)
			return nonce[:]
		}}
}

func cipherAESGCM(key []byte) Cipher {
	c, err := chacha20poly1305.New(key[:])
	if err != nil {
		panic(err)
	}
	return aeadCipher{
		c,
		func(u uint64) []byte {
			var nonce [12]byte
			binary.LittleEndian.PutUint64(nonce[4:], u)
			return nonce[:]
		}}
}

type aeadCipher struct {
	cipher.AEAD
	nonce func(uint64) []byte
}

func (c aeadCipher) Encrypt(out []byte, n uint64, ad, plaintext []byte) []byte {
	return c.Seal(out, c.nonce(n), plaintext, ad)
}

func (c aeadCipher) Decrypt(out []byte, n uint64, ad, ciphertext []byte) ([]byte, error) {
	return c.Open(out, c.nonce(n), ciphertext, ad)
}

func (c cpabe) GenerateAndSaveMasterKeys() error {
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

func (c cpabe) PubKeygen() (*abe.FAMEPubKey, *abe.FAMESecKey, error) {
	pk, sk, err := c.scheme.GenerateMasterKeys()
	if err != nil {
		return nil, nil, err
	}
	return pk, sk, err
}

func (c cpabe) PrivKeygen(sk *abe.FAMESecKey, attributes []string) (*abe.FAMEAttribKeys, error) {
	ak, err := c.scheme.GenerateAttribKeys(attributes, sk)
	if err != nil {
		return nil, err
	}
	return ak, nil
}

type PublicKey abe.FAMEPubKey
type AttrKey abe.FAMEAttribKeys
type SecretKey abe.FAMESecKey

func (k PublicKey) Marshal() ([]byte, error) {
	return json.Marshal(k)
}

func (k AttrKey) Marshal() ([]byte, error) {
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

type cpabe struct {
	scheme *abe.FAME
	pubKey *abe.FAMEPubKey
	attKey *abe.FAMEAttribKeys
}

func newCPABE(isPublisher bool, keyPath string) *cpabe {

	s := &cpabe{
		scheme: abe.NewFAME(),
	}
	//one is loaded, the other not
	if isPublisher {
		s.pubKey = &abe.FAMEPubKey{}
		s.attKey = &abe.FAMEAttribKeys{}
	} else {
		s.pubKey = &abe.FAMEPubKey{}
		s.attKey = &abe.FAMEAttribKeys{}
	}
	return s
}

// FAME encryption function wrapper
func (c cpabe) Encrypt(msg string, msp *abe.MSP) (*abe.FAMECipher, error) {
	return c.scheme.Encrypt(msg, msp, c.pubKey)
}

// FAME decryption function wrapper
func (c cpabe) Decrypt(ciphertext *abe.FAMECipher) (string, error) {
	return c.scheme.Decrypt(ciphertext, c.attKey, nil)
}

// Encrypt and encode
func (c cpabe) EncryptEncode(key *abe.FAMEPubKey, policy string, plaintext string) ([]byte, error) {
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

func (c cpabe) DecryptDecode(ciphertext []byte) (*string, error) {
	var ct *abe.FAMECipher
	err := json.Unmarshal(ciphertext, &ct)
	pt, err := c.Decrypt(ct)
	if err != nil {
		return nil, err
	}
	return &pt, nil
}
