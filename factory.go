package main

import (
	"godeliver/driver"
	"godeliver/tcp"
	"godeliver/udp"
)

func CreateServer(proto string, flag string, crypt driver.CryptionDriver) driver.ServerDriver {
	if proto == "tcp" {
		return tcp.CreateServer(flag, crypt)
	}

	return udp.CreateServer(flag, crypt)
}
