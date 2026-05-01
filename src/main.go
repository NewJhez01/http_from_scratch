package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("error")
	}
	lines := getLineChannel(file)
	for v := range lines {
		fmt.Println("read: " + v)
	}
}

func getLineChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	buffer := make([]byte, 8)
	currentLine := ""
	go func() {
		for {
			n, err := f.Read(buffer)
			if err == io.EOF {
				if n != 0 {
					lines <- currentLine + string(buffer[:n])
				}
				close(lines)
				f.Close()
				break
			}

			currentLine += string(buffer[:n])
			parts := strings.Split(currentLine, "\n")
			if len(parts) > 1 {
				lines <- parts[0]
				currentLine = strings.Join(parts[1:], "\n")
			} else {
				currentLine = parts[0]
			}
		}
	}()
	return lines
}
