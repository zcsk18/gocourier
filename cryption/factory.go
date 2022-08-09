package cryption

import (
	"godeliver/conf"
	"godeliver/driver"
)

func CreateCryption (sType string) (driver.CryptionDriver){
	switch sType {
	case "rsa":
		crypt := &RSA{}
		pub := conf.GetIniValue("key", "pub")
		if pub != "" {
			crypt.SetPublicKey(pub)
		}

		pri := conf.GetIniValue("key", "pri")
		if pri != "" {
			crypt.SetPrivateKey(pri)
		}
		return crypt
	}

	return &Empty{}
}