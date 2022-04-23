package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	h := Header{
		Type:   1,
		Cipher: 2,
		Nonce:  0,
	}
	benc := Encode(make([]byte, 9), h.Type, h.Cipher, h.Nonce)
	bdec := Header{}
	bdec.Decode(benc)
	assert.Equal(t, h, bdec)
}
