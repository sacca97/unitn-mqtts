package crypto

import (
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fentec-project/gofe/abe"
)

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

func loadKey(path string) (string, []byte) {
	k, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	key, _ := pem.Decode(k)
	if key == nil {
		panic("invalid key structure")
	}
	return key.Type, key.Bytes
}

func SymKeygen() [32]byte {
	var k [32]byte
	_, err := rand.Read(k[:])
	if err != nil {
		panic(err)
	}
	return k
}

func Encode(name string, k any) *pem.Block {
	out, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	return &pem.Block{
		Type:  name,
		Bytes: out,
	}
}

func GenerateFAMEKeys() {

	fame := abe.NewFAME()
	pk, sk, err := fame.GenerateMasterKeys()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	if err := pem.Encode(pemPublic, pub); err != nil {
		log.Fatal(err)
	}

	pemSecret, err := os.Create("secret.key")
	if err != nil {
		log.Fatal(err)
	}
	if err := pem.Encode(pemSecret, sec); err != nil {
		log.Fatal(err)
	}

	pemPublic.Close()
	pemSecret.Close()
}

func GenerateAttribKeys(sk *abe.FAMESecKey, attributes []string, fileName string) {
	fame := abe.NewFAME()
	ak, err := fame.GenerateAttribKeys(attributes, sk)
	if err != nil {
		log.Fatal(err)
	}
	bytesAttr, _ := json.Marshal(ak)
	attr := &pem.Block{
		Type:  "FAME ATTRIBUTES KEY",
		Bytes: bytesAttr,
	}
	if fileName == "" {
		fileName = "attributes.key"
	}
	pemAttr, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if err := pem.Encode(pemAttr, attr); err != nil {
		log.Fatal(err)
	}
	pemAttr.Close()

}
