package main

import (
	"flag"
	"godeliver/conf"
	"godeliver/cryption"
	"godeliver/udp"
)


var strFlag = flag.String("s", "", "-s [srv|clt]")

func main() {
	flag.Parse()
	conf.SetIni("cfg.ini")

	crypt := &cryption.RSA{}
	crypt.SetPublicKey("public.pem")
	crypt.SetPrivateKey("private.pem")

	if *strFlag == "srv" {
		Srv := udp.CreateProxyServer()
		Srv.Crypt = crypt
		Srv.Run()
	} else {
		Clt := udp.CreateClientServer()
		Clt.Crypt = crypt
		Clt.Run()
	}
}