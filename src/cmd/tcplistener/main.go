package main

import (
	"fmt"
	"net"

	"http_from_scratch/src/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("listener error")
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("connection failed")
	}
	lines, err := request.RequestFromReader(conn)
	fmt.Println("Method: " + lines.RequestLine.Method)
	fmt.Println("Target: " + lines.RequestLine.RequestTarget)
	fmt.Println("Version: " + lines.RequestLine.HttpVersion)
}
