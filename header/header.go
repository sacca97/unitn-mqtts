package header

import (
	"encoding/binary"
	"fmt"
)

const (
	Len = 9
)

const (
	Symmetric uint8 = 0
	Cpabe     uint8 = 1
	Kpabe     uint8 = 2
)

const (
	AESGCM uint8 = 0
	CHACHA uint8 = 1
	FAME   uint8 = 0
	//TODO: Define other ciphers
)

var ErrInvalidHeaderLength = fmt.Errorf("invalid header length")

type Header struct {
	Type   uint8
	Cipher uint8
	Nonce  uint64
}

func Encode(b []byte, ct uint8, c uint8, nonce uint64) []byte {
	b = b[:Len]
	b[0] = ct<<4 | byte(c&0x0f)
	binary.BigEndian.PutUint64(b[1:], nonce)
	return b
}

func (h *Header) Decode(b []byte) error {
	if len(b) != Len {
		return fmt.Errorf("invalid header length")
	}
	h.Type = uint8((b[0] >> 4) & 0x0f)
	h.Cipher = uint8(b[0] & 0x0f)
	h.Nonce = binary.BigEndian.Uint64(b[1:])

	return nil
}

func (h *Header) String() string {
	return fmt.Sprintf("CipherType: %d, Cipher: %d, Nonce: %d", h.Type, h.Cipher, h.Nonce)
}

func (h *Header) IsValid() bool {
	return h.Type != 0 && h.Nonce == 0
}
