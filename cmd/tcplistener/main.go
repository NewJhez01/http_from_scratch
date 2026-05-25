package main

import (
	"fmt"
	"net"

	"http_from_scratch/internal/request"
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
	req, err := request.RequestFromReader(conn)
	fmt.Println("Request:")
	fmt.Println("Method: " + req.RequestLine.Method)
	fmt.Println("Target: " + req.RequestLine.RequestTarget)
	fmt.Println("Version: " + req.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for k, v := range req.Headers {
		fmt.Println(k + ": " + v)
	}
	fmt.Println("Body:")
	fmt.Println(string(req.Body))
}
