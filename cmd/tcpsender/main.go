package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	serverAddr := "127.0.0.1:42069"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to", serverAddr)

	message := "Hello from Go client!"

	i := 0
	for {
		i++
		_, err = fmt.Fprintf(conn, "%s%v\n", message, i)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}
