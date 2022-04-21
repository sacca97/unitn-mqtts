package header

const (
	len = 9
)

type Header struct {
	CipherType uint8
	Cipher     uint8
	Nonce      uint64
}

func Encode(ct uint8, cipher uint8, nonce uint64) []byte {
	b := make([]byte, len)
	b[0] = byte(ct)
	return []byte{}
}
