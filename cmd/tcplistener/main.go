package main

import (
	"fmt"
	"net"
	"tcpToHttp/internal/request"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string, 1)
//
// 	go func() {
// 		defer f.Close()
// 		defer close(out)
// 		str := ""
// 		for {
// 			data := make([]byte, 8)
// 			n, err := f.Read(data)
// 			if err != nil {
// 				break
// 			}
// 			data = data[:n]
// 			if i := bytes.IndexByte(data, '\n'); i != -1 {
// 				str += string(data[:i])
// 				data = data[i+1:]
// 				out <- str
// 				str = ""
// 			}
//
// 			str += string(data)
//
// 		}
//
// 		if len(str) != 0 {
// 			out <- str
// 		}
// 	}()
//
// 	return out
// }

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
		// for line := range getLinesChannel(conn) {
		// 	fmt.Printf("read :%s\n", line)
		// }
	}

}
