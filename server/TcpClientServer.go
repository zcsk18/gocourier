package server

import (
	"context"
	"godeliver/conf"
	"godeliver/driver"
	"log"
	"net"
)

type TcpClientServer struct {
	*driver.ClientServer
}

func CreateTcpClientServer() *TcpClientServer {
	port := conf.GetIniValue("clt", "port")
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("start srv on: %s", port)
	srv := &TcpClientServer{}
	srv.ClientServer = &driver.ClientServer{}
	srv.Listener = listen
	return srv
}

func (this *TcpClientServer) Run() {
	for {
		s, err := this.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}


func (this *TcpClientServer) handleClt(clt net.Conn) {
	defer func() {
		clt.Close()
		log.Printf("clt disconnect %s\n", clt.RemoteAddr().String())
	}()


	srv, err := this.handShakeWithSrv()
	defer func() {
		srv.Close()
		log.Printf("srv disconnect %s\n", srv.RemoteAddr().String())
	}()

	if err != nil {
		log.Printf("err: %s\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go this.HandleStoC(ctx, clt, srv)
	this.HandleCtoS(ctx, clt, srv)
	cancel()
}

func (this *TcpClientServer) handShakeWithSrv() (net.Conn, error) {
	ip := conf.GetIniValue("srv", "ip")
	port := conf.GetIniValue("srv", "port")

	srv, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("proxy connect : %s\n", srv.RemoteAddr().String())

	return srv, err
}
