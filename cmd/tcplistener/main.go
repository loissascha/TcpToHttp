package main

import (
	"fmt"
	"net"
	"tcpToHttp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
	}

}
