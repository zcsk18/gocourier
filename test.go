package main

import (
	"fmt"
	"godeliver/misc"
)



func main() {
	var i = 1<<31 - 1
	fmt.Println(i)
	b := misc.IntToBytes(i)
	fmt.Println(len(b))
	fmt.Println(misc.BytesToInt(b))
}
