package main

import (
	"flag"
	"godeliver/conf"
	"godeliver/cryption"
	"godeliver/driver"
	"godeliver/server"
)

var strFlag = flag.String("s", "", "-s [srv|clt]")

func CreateServer(proto string, flag string, crypt driver.CryptionDriver) driver.ServerDriver {
	if proto == "tcp" {
		return server.CreateTcpServer(flag, crypt)
	}

	return server.CreateUdpServer(flag, crypt)
}

func main() {
	flag.Parse()
	conf.SetIni("cfg.ini")
	proto := conf.GetIniValue("common", "protocol")
	crypt := conf.GetIniValue("common", "crypt")

	srv := CreateServer(proto, *strFlag, cryption.CreateCryption(crypt))
	srv.Run()
}