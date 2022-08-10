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
	"time"
)

type UdpProxyServer struct {
	*driver.ProxyServer
}

func CreateUdpProxyServer() *UdpProxyServer {
	salt := conf.GetIniValue("udp", "salt")
	password :=  conf.GetIniValue("udp", "password")
	port := conf.GetIniValue("srv", "port")

	key := pbkdf2.Key([]byte(password), []byte(salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	listener, err := kcp.ListenWithOptions(":"+port, block, 10, 3)
	if err != nil {
		panic(err)
	}

	log.Printf("start proxy on: %s", port)
	srv := &UdpProxyServer{}
	srv.ProxyServer = &driver.ProxyServer{}
	srv.Listener = listener
	return srv
}

func (this *UdpProxyServer) Run() {
	for {
		s, err := this.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}

func (this *UdpProxyServer) handleClt(session net.Conn) {
	defer func() {
		session.Close()
		log.Printf("clt disconnect %s\n", session.RemoteAddr().String())
	}()

	session.SetDeadline(time.Now().Add(30*time.Second))

	log.Printf("new clt: %s", session.RemoteAddr().String())

	proxy, err := net.Dial("tcp", conf.GetIniValue("proxy", "ip")+":"+conf.GetIniValue("proxy", "port"))
	if err != nil {
		log.Println("err :", err)
		return
	}
	defer func() {
		proxy.Close()
		log.Printf("proxy disconnect %s\n", proxy.RemoteAddr().String())
	}() // 关闭连接


	ctx, cancel := context.WithCancel(context.Background())
	go this.HandlePtoC(ctx, session, proxy)
	this.HandleCtoP(ctx, session, proxy)
	cancel()
}
