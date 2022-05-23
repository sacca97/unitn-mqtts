package mqtts

import (
	"encoding/json"

	"github.com/sacca97/unitn-mqtts/crypto"
)

type ConnectionState struct {
	publisher            bool
	atomicMessageCounter uint64
	cipher               crypto.Cipher
	status               int8
}

func newConnectionState(publisher bool, cipher string) *ConnectionState {

	//Decide how to select the cipher
	var c crypto.Cipher

	switch cipher {
	case "chacha20poly1305":
		c = crypto.CipherChaChaPoly([32]byte{})
	case "aesgcm":
		c = crypto.CipherAESGCM([32]byte{})
	case "fame":
		c = crypto.CipherFame(publisher)

	default:
		//Back on FAME algo in case of wrong setup
		c = crypto.CipherFame(publisher)
	}

	return &ConnectionState{
		publisher:            publisher,
		atomicMessageCounter: 0,
		cipher:               c,
		status:               0,
	}
}

func (cs *ConnectionState) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs)
}
