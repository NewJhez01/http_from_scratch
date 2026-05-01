package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("listener error")
		return
	}
	lines := getLineChannel(listener)
	for v := range lines {
		fmt.Println("read: " + v)
	}
}

func getLineChannel(l net.Listener) <-chan string {
	currentLine := ""
	buffer := make([]byte, 8)
	data := make(chan string)
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("connection failed")
		return data
	}

	go func(c net.Conn) {
		for {
			n, err := c.Read(buffer)
			if err == io.EOF {
				if n > 0 {
					data <- currentLine + string(buffer[:n])
				}
				close(data)
				break
			}
			if err != nil {
				fmt.Println("something went wrong")
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			if len(parts) > 1 {
				data <- currentLine + parts[0]
				currentLine = strings.Join(parts[1:], "\n")
			} else {
				currentLine += parts[0]
			}
		}
	}(conn)
	return data
}
