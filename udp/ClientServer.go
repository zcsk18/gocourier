package udp

import (
	"context"
	"crypto/sha1"
	"github.com/xtaci/kcp-go/v5"
	"godeliver/conf"
	"godeliver/driver"
	"godeliver/misc"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
	"net"
	"time"
)

type ClientServer struct {
	Listener net.Listener
	Crypt    driver.CryptionDriver
}

func CreateClientServer() *ClientServer {
	port := conf.GetIniValue("clt", "port")
	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("start srv on: %s", port)
	srv := &ClientServer{}
	srv.Listener = listen
	return srv
}

func (this *ClientServer) Run() {
	for {
		s, err := this.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go this.handleClt(s)
	}
}

func (this * ClientServer) handleClt(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("clt disconnect %s\n", conn.RemoteAddr().String())
	}()


	sess, err := this.handShakeWithSrv()
	defer func() {
		sess.Close()
		log.Printf("srv disconnect %s\n", sess.RemoteAddr().String())
	}()

	if err != nil {
		log.Printf("err: %s\n", err)
		return
	}

	sess.SetNoDelay(1,20,2,1)
	ctx, cancel := context.WithCancel(context.Background())
	go this.handleStoC(ctx, conn, sess)
	this.handleCtoS(ctx, conn, sess)
	cancel()
}

func (this *ClientServer) handleStoC(ctx context.Context, conn net.Conn, sess *kcp.UDPSession) {
	buf := make([]byte, conf.BufLen*2)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := io.ReadAtLeast(sess, buf, 1)
			if err != nil {
				return
			}
			//log.Printf("recv srv[%s] %d\n", sess.RemoteAddr().String(), n)
			n, err = conn.Write(buf[:n])
			if err != nil {
				return
			}
			//log.Printf("send clt[%s] %d\n",  conn.RemoteAddr().String(), n)
		}
	}
}

func (this *ClientServer) handleCtoS(ctx context.Context, conn net.Conn, sess *kcp.UDPSession) {
	buf := make([]byte, conf.BufLen)
	for {
		conn.SetDeadline(time.Now().Add(60*time.Second))

		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}

		msg := this.Crypt.Encode(buf[:n])
		_, err = misc.Send(sess, "p", msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func (this * ClientServer) handShakeWithSrv() (*kcp.UDPSession, error) {
	salt := conf.GetIniValue("udp", "salt")
	password :=  conf.GetIniValue("udp", "password")

	ip := conf.GetIniValue("srv", "ip")
	port := conf.GetIniValue("srv", "port")

	key := pbkdf2.Key([]byte(password), []byte(salt), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	sess, err := kcp.DialWithOptions(ip+":"+port, block, 10, 3)
	if err != nil {
		panic(err)
	}
	log.Printf("proxy connect : %s\n", sess.RemoteAddr().String())

	return sess, err
}
