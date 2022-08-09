package main

import (
	"flag"
	"godeliver/conf"
	"godeliver/cryption"
)

var strFlag = flag.String("s", "", "-s [srv|clt]")

func main() {
	flag.Parse()
	conf.SetIni("cfg.ini")
	proto := conf.GetIniValue("common", "protocol")
	crypt := conf.GetIniValue("common", "crypt")

	srv := CreateServer(proto, *strFlag, cryption.CreateCryption(crypt))
	srv.Run()
}