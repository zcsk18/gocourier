package tcp

import (
	"godeliver/misc"
	"io"
	"net"
)

func Recv(conn net.Conn) (string, []byte, error) {
	c := make([]byte, 1)
	_,err := conn.Read(c)
	if err != nil {
		return "", nil, err
	}

	l := make([]byte, 4)
	_,err = conn.Read(l)
	if err != nil {
		return "", nil, err
	}
	n := misc.BytesToInt(l)
	msg := make([]byte, n)


	amount, err := io.ReadAtLeast(conn, msg, n)
	if err != nil{
		return "", nil, err
	}
	if amount != n {
		return "", nil, err
	}

	return string(c), msg, nil
}

func Send(conn net.Conn, c string, msg []byte) (int, error) {
	if len(c) != 1 {
		panic("contral len not eq 1")
	}

	n := len(msg)
	_, err := conn.Write([]byte(c))
	if err != nil {
		return 0, err
	}
	_, err =conn.Write(misc.IntToBytes(n))
	if err != nil {
		return 0, err
	}
	_, err =conn.Write(msg)
	if err != nil {
		return 0, err
	}

	return n, nil
}