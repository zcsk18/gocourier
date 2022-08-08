package cryption

type Empty struct {
}

func (this *Empty) Encode(plainText []byte) []byte {
	return plainText
}

func (this *Empty) Decode(cipherText []byte) []byte {
	return cipherText
}