package header

import (
	"encoding/binary"
	"fmt"
	"log"
)

//empty bytes ? (3 bytes)
const (
	Len = 12
)

const (
	SYMMETRIC uint8 = 1
	CPABE     uint8 = 2
	KPABE     uint8 = 3
)

const (
	AESGCM     uint8 = 0
	CHACHAPOLY uint8 = 1
	FAME       uint8 = 0
	//TODO: Define other ciphers
)

var ErrInvalidHeaderLength = fmt.Errorf("invalid header length")

type Header struct {
	Type   uint8
	Cipher uint8
	Nonce  uint64
}

func Create(algo string) []byte {
	h := make([]byte, 12)
	switch algo {
	case "chacha20poly1305":
		Encode(h, SYMMETRIC, CHACHAPOLY, 0)
	case "aesgcm":
		Encode(h, SYMMETRIC, AESGCM, 0)
	case "fame":
		Encode(h, CPABE, FAME, 0)
	default:
		log.Fatal("unsupported algorithm")
	}
	return h
}

func Encode(b []byte, ct uint8, c uint8, nonce uint64) []byte {
	b = b[:Len]
	//b[0] = ct<<4 | byte(c&0x0f)
	b[0] = ct<<2 | byte(c&0x3)
	binary.BigEndian.PutUint64(b[1:9], nonce)
	return b
}

func (h *Header) Decode(b []byte) error {
	if len(b) != Len {
		return fmt.Errorf("invalid header length")
	}
	//h.Type = uint8((b[0] >> 4) & 0x0f)
	h.Type = uint8((b[0] >> 2) & 0x3)
	//h.Cipher = uint8(b[0] & 0x0f)
	h.Cipher = uint8(b[0] & 0x3)
	h.Nonce = binary.BigEndian.Uint64(b[1:9])

	return nil
}

func (h *Header) String() string {
	return fmt.Sprintf("CipherType: %d, Cipher: %d, Nonce: %d", h.Type, h.Cipher, h.Nonce)
}

func (h *Header) IsValid() bool {
	return h.Type != 0 && h.Nonce == 0
}
