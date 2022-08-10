package server

import (
	"context"
	"crypto/sha1"
	"github.com/xtaci/kcp-go/v5"
	"godeliver/conf"
	"godeliver/driver"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"net"
)

type UdpClientServer struct {
	*driver.ClientServer
}

func CreateUdplientServer() *UdpClientServer {
	port := conf.GetIniValue("clt", "port")
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("start srv on: %s", port)
	srv := &UdpClientServer{}
	srv.ClientServer = &driver.ClientServer{}
	srv.Listener = listen
	return srv
}

func (this *UdpClientServer) Run() {
	for {
		s, err := this.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}

func (this *UdpClientServer) handleClt(clt net.Conn) {
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

	srv.SetNoDelay(1,20,2,1)
	ctx, cancel := context.WithCancel(context.Background())
	go this.HandleStoC(ctx, clt, srv)
	this.HandleCtoS(ctx, clt, srv)
	cancel()
}

func (this *UdpClientServer) handShakeWithSrv() (*kcp.UDPSession, error) {
	salt := conf.GetIniValue("udp", "salt")
	password :=  conf.GetIniValue("udp", "password")

	ip := conf.GetIniValue("srv", "ip")
	port := conf.GetIniValue("srv", "port")

	key := pbkdf2.Key([]byte(password), []byte(salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	srv, err := kcp.DialWithOptions(ip+":"+port, block, 10, 3)
	if err != nil {
		panic(err)
	}
	log.Printf("proxy connect : %s\n", srv.RemoteAddr().String())

	return srv, err
}
