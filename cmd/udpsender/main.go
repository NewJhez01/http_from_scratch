package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("error for resolving")
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("error for estab connection")
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read")
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("failed to write")
		}
	}
}
