package main

import (
	"godeliver/cryption"
	"godeliver/driver"
	"godeliver/tcp"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

func handle(conn net.Conn, driver driver.CryptionDriver) {
	defer conn.Close()
	echo(conn, driver)
}

func echo(conn net.Conn, driver driver.CryptionDriver) {
	for {
		conn.SetDeadline(time.Now().Add(10*time.Second))
		c, msg, err := tcp.Recv(conn)
		if err != nil {
			panic(err)
		}
		msg = driver.Decode(msg)
		//log.Printf("recv[%s] %s %s\n", conn.RemoteAddr().String(), c, msg[len(msg)-5:])
		if c == "n" {
			c = "r"
		} else {
			c = "n"
		}

		v := strconv.Itoa(rand.Intn(10000))
		msg = append(msg, []byte("-> hello: "+v)...)
		start := time.Now()
		msg = driver.Encode(msg)
		tcp.Send(conn, c, msg)
		end := time.Now()
		log.Printf("send[%s] %d %s %d\n", conn.RemoteAddr().String(), len(msg), v, end.Sub(start).Microseconds())
		time.Sleep(1000*time.Millisecond)
	}
}


func srv(driver driver.CryptionDriver, wait *sync.WaitGroup) {
	listen, err := net.Listen("tcp", "127.0.0.1:8338")
	if err != nil {
		panic(err)
	}
	wait.Done()
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go handle(conn, driver)
	}
}

func clt(driver driver.CryptionDriver) {
	conn, err := net.Dial("tcp", "127.0.0.1:8338")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	msg := []byte("hello")
	msg = driver.Encode(msg)
	tcp.Send(conn, "n", msg)
	echo(conn, driver)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	crypt := &cryption.RSA{}
	crypt.SetPublicKey("public.pem")
	crypt.SetPrivateKey("private.pem")

	wait := sync.WaitGroup{}
	wait.Add(1)
	go srv(crypt, &wait)
	wait.Wait()

	clt(crypt)
}
