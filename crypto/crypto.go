package crypto

import (
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fentec-project/gofe/abe"
)

func Marshal(k any) ([]byte, error) {
	return json.Marshal(k)
}

func UnmarshalFamePrivKey(data []byte) (*abe.FAMESecKey, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("nil byte array")
	}
	var k *abe.FAMESecKey
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	}
	return k, err
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

func loadKey(path string) []byte {
	k, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	key, _ := pem.Decode(k)
	//I have to check stuff here obviously
	return key.Bytes
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

func GenerateFAMEKeys() error {

	fame := abe.NewFAME()
	pk, sk, err := fame.GenerateMasterKeys()
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

func GenerateAttribKeys(sk *abe.FAMESecKey, attributes []string) error {
	fame := abe.NewFAME()
	ak, err := fame.GenerateAttribKeys([]string{"1", "2", "3", "4", "5"}, sk)
	if err != nil {
		return err
	}
	bytesAttr, _ := json.Marshal(ak)
	attr := &pem.Block{
		Type:  "FAME ATTRIBUTES KEY",
		Bytes: bytesAttr,
	}
	pemAttr, err := os.Create("attributes.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := pem.Encode(pemAttr, attr); err != nil {
		panic(err)
	}
	pemAttr.Close()

	return nil
}
