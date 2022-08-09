package tcp

import "godeliver/driver"

func CreateServer(flag string, crypt driver.CryptionDriver) driver.ServerDriver {
	if flag == "srv" {
		srv := CreateProxyServer()
		srv.Crypt = crypt
		return srv
	}
	srv := CreateClientServer()
	srv.Crypt = crypt
	return srv
}
