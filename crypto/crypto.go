package crypto

import (
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/fentec-project/gofe/abe"
)

func Marshal(k any) ([]byte, error) {
	return json.Marshal(k)
}

func UnmarshalFamePubKey(data []byte) (*abe.FAMEPubKey, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("nil byte array")
	}
	var k *abe.FAMEPubKey
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}
	return k, err
}

//Decode a FameKey from a byte array
func UnmarshalFameAttrKey(data []byte) (*abe.FAMEAttribKeys, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("nil byte array")
	}
	var k *abe.FAMEAttribKeys
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}
	return k, err
}

func loadKey(path string) (any, error) {
	//load pk
	return nil, nil
}

func SymKeygen() [32]byte {
	var k [32]byte
	_, err := rand.Read(k[:])
	if err != nil {
		panic(err)
	}
	return k
}

func NewPem(name string, k []byte) *pem.Block {
	out := &pem.Block{
		Type:  name,
		Bytes: k,
	}

	return out
}

func Encode(name string, k any) *pem.Block {
	out, err := Marshal(k)
	if err != nil {
		panic(err)
	}
	return NewPem(name, out)
}

/*
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
}*/
