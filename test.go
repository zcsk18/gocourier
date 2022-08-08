package main

import "godeliver/cryption"

func main() {
	crypt := &cryption.RSA{}
	crypt.GenerateKey()
}
