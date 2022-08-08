package driver

type CryptionDriver interface {
	Encode(plainText []byte) []byte
	Decode(cipherText []byte) []byte
}
