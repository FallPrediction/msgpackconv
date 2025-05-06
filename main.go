package main

import (
	"fmt"
	"msgpackconv/msgpack"
)

func main() {
	msg := msgpack.FromJSON([]byte(`{"str": "a"}`))
	fmt.Printf("% x\n", msg)
	json := msgpack.ToJSON(msg)
	fmt.Println(string(json))
}
