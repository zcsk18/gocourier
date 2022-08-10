package driver

import (
	"context"
	"godeliver/conf"
	"godeliver/misc"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type ServerDriver interface {
	Run()
}

type ClientServer struct {
	Listener net.Listener
	Crypt    CryptionDriver
}

type ProxyServer struct {
	Listener net.Listener
	Crypt    CryptionDriver
}

func (this *ClientServer) HandleStoC(ctx context.Context, clt io.Writer, srv io.Reader) {
	buf := make([]byte, conf.BufLen*2)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := srv.Read(buf)
			if err != nil {
				return
			}
			n, err = clt.Write(buf[:n])
			if err != nil {
				return
			}
		}
	}
}

func (this *ProxyServer) HandlePtoC(ctx context.Context, clt io.Writer, proxy io.Reader) {
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
			n, err = clt.Write(buf[:n])
			if err != nil {
				return
			}
		}
	}
}

func (this *ProxyServer) HandleCtoP(ctx context.Context, clt net.Conn, proxy net.Conn) {
	for {
		clt.SetDeadline(time.Now().Add(30*time.Second))
		_, msg, err := misc.Recv(clt)
		if err != nil {
			log.Println(err)
			break
		}

		msg = this.Crypt.Decode(msg)
		_, err = proxy.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func (this *ClientServer) HandleCtoS(ctx context.Context, clt net.Conn, srv net.Conn) {
	buf := make([]byte, 512)
	first := true
	for {
		clt.SetDeadline(time.Now().Add(60*time.Second))

		n, err := clt.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}

		if first {
			first = false
			str := string(buf)
			if strings.Contains(str, "CONNECT") ||
				strings.Contains(str, "GET") ||
				strings.Contains(str, "POST") ||
				strings.Contains(str, "HEAD") {
				misc.AnalysisHttp(str)
			}
		}

		msg := this.Crypt.Encode(buf[:n])
		_, err = misc.Send(srv, "p", msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
}



