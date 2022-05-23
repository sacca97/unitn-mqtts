package header

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	h := Header{
		Type:   1,
		Cipher: 2,
		Nonce:  0,
	}
	benc := Encode(make([]byte, 12), 1, 2, 0)
	fmt.Println(benc)
	bdec := Header{}
	bdec.Decode(benc)
	assert.Equal(t, h, bdec)
}
