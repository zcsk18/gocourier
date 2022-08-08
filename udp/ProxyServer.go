package udp

import (
	"context"
	"crypto/sha1"
	"github.com/xtaci/kcp-go/v5"
	"godeliver/conf"
	"godeliver/driver"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"net"
	"time"
)

type ProxyServer struct {
	Listener *kcp.Listener
	Crypt    driver.CryptionDriver
}

func CreateProxyServer() *ProxyServer {
	salt := conf.GetIniValue("common", "salt")
	password :=  conf.GetIniValue("common", "password")
	port := conf.GetIniValue("srv", "port")

	key := pbkdf2.Key([]byte(password), []byte(salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	listener, err := kcp.ListenWithOptions(":"+port, block, 10, 3)
	if err != nil {
		panic(err)
	}

	log.Printf("start proxy on: %s", port)
	srv := &ProxyServer{}
	srv.Listener = listener
	return srv
}

func (this *ProxyServer) Run() {
	for {
		s, err := this.Listener.AcceptKCP()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}

func (this *ProxyServer) handleClt(session *kcp.UDPSession) {
	defer func() {
		session.Close()
		log.Printf("clt disconnect %s\n", session.RemoteAddr().String())
	}()

	session.SetNoDelay(1,20,2,1)
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
	go this.handlePtoC(ctx, session, proxy)
	this.handleCtoP(ctx, session, proxy)
	cancel()
}

func (this *ProxyServer) handleCtoP(ctx context.Context, session *kcp.UDPSession, proxy net.Conn) {
	buf := make([]byte, conf.BufLen*2)
	for {
		session.SetDeadline(time.Now().Add(30*time.Second))
		n, err := io.ReadAtLeast(session, buf, 1)
		if err != nil {
			log.Println(err)
			break
		}
		//log.Printf("recv clt[%s] %d\n", session.RemoteAddr().String(), n)
		msg := this.Crypt.Decode(buf[:n])
		n, err = proxy.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
		//log.Printf("send proxy[%s] %d\n", proxy.RemoteAddr().String(),  n)
	}
}

func (this *ProxyServer) handlePtoC(ctx context.Context, session *kcp.UDPSession, proxy net.Conn) {
	buf := make([]byte, conf.BufLen)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := proxy.Read(buf)
			if err != nil {
				return
			}
			//log.Printf("recv proxy[%s] %d\n", proxy.RemoteAddr().String(), n)
			n, err = session.Write(buf[:n])
			if err != nil {
				return
			}
			//log.Printf("send clt[%s] %d\n",session.RemoteAddr().String(), n)
		}
	}
}
