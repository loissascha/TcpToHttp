package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr := "127.0.0.1:42069"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to", serverAddr)

	message := "Hello from Go client!\n"

	_, err = conn.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}
