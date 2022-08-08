package tcp

import (
	"context"
	"godeliver/conf"
	"godeliver/driver"
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
	go this.handleStoC(ctx, conn, srv)
	this.handleCtoS(ctx, conn, srv)
	cancel()
}

func (this *ClientServer) handleStoC(ctx context.Context, conn net.Conn, srv net.Conn) {
	buf := make([]byte, conf.BufLen*2)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			//log.Printf("recv srv[%s] %d\n", srv.RemoteAddr().String(), n)
			n, err = conn.Write(buf[:n])
			if err != nil {
				return
			}
			//log.Printf("send clt[%s] %d\n",  conn.RemoteAddr().String(), n)
		}
	}
}

func (this *ClientServer) handleCtoS(ctx context.Context, conn net.Conn, srv net.Conn) {
	buf := make([]byte, conf.BufLen)
	for {
		conn.SetDeadline(time.Now().Add(60*time.Second))

		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		//log.Printf("recv clt[%s] %d\n", conn.RemoteAddr().String(), n)
		msg := this.Crypt.Encode(buf[:n])

		_, err = srv.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
		//log.Printf("send srv[%s] %d\n", srv.RemoteAddr().String(), n)
	}
}

func (this * ClientServer) handShakeWithSrv() (net.Conn, error) {
	ip := conf.GetIniValue("srv", "ip")
	port := conf.GetIniValue("srv", "port")

	srv, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("proxy connect : %s\n", srv.RemoteAddr().String())

	return srv, err
}
