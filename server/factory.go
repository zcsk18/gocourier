package server

import "godeliver/driver"

func CreateTcpServer(flag string, crypt driver.CryptionDriver) driver.ServerDriver {
	if flag == "srv" {
		srv := CreateTcpProxyServer()
		srv.Crypt = crypt
		return srv
	}
	srv := CreateTcpClientServer()
	srv.Crypt = crypt
	return srv
}


func CreateUdpServer(flag string, crypt driver.CryptionDriver) driver.ServerDriver {
	if flag == "srv" {
		srv := CreateUdpProxyServer()
		srv.Crypt = crypt
		return srv
	}
	srv := CreateUdplientServer()
	srv.Crypt = crypt
	return srv
}