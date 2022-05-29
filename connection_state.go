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

/* Creates a new ConnectionState with a cryptographic scheme from the
 * given configuration. If no encryption algorithm is matched fallback to  FAME CP-ABE
 */
func newConnectionState(config *Config) *ConnectionState {
	var c crypto.Cipher

	switch config.Crypto {
	case "chacha20poly1305":
		c = crypto.CipherChaChaPoly([32]byte{})
	case "aesgcm":
		c = crypto.CipherAESGCM([32]byte{})
	case "fame":
		c = crypto.CipherFame(config.Publisher, config.EncryptionKey)
	default:
		c = crypto.CipherFame(config.Publisher, config.EncryptionKey)
	}

	return &ConnectionState{
		publisher:            config.Publisher,
		atomicMessageCounter: 0,
		cipher:               c,
		status:               0,
	}
}

func (cs *ConnectionState) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs)
}
