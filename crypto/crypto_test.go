package crypto

import (
	"encoding/pem"
	"io/ioutil"
	"testing"
)

func TestMasterKeygen(t *testing.T) {
	GenerateFAMEKeys()
}

func TestKeyLoader(t *testing.T) {
	pk, err := ioutil.ReadFile("public.key")
	if err != nil {
		panic(err)
	}
	pub, _ := pem.Decode(pk)
	if pub.Type != "FAME PUBLIC KEY" {
		t.Fatalf("Expected FAME PUBLIC KEY, got %s", pub.Type)
	}
	sk, err := ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}
	sec, _ := pem.Decode(sk)
	if sec.Type != "FAME SECRET KEY" {
		t.Fatalf("Expected FAME SECRET KEY, got %s", sec.Type)
	}

	secret, err := UnmarshalFamePrivKey(sec.Bytes)
	if err != nil {
		t.Fatalf("Failed to unmarshal secret key: %s", err)
	}
	GenerateAttribKeys(secret, []string{"1", "2", "3", "4", "5"})
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %s", err)
	}

}

func TestAttrKeygen(t *testing.T) {

}
