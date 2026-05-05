package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader1(r io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	n, err := r.Read(buffer)
	if err != nil {
		fmt.Println("unexpected error  " + err.Error())
		return nil, errors.New("io fail")
	}
	rl, _, err := lineParser1(n, buffer)
	if err != nil {
		fmt.Println("function call failed")
		return nil, errors.New("function fail")
	}

	return &Request{RequestLine: *rl}, nil
}

func lineParser1(n int, b []byte) (*RequestLine, int, error) {
	if !strings.Contains(string(b), "\r\n") {
		return nil, 0, nil
	}
	l := strings.Split(string(b), "\r\n")
	requestLine := strings.Split(l[0], " ")

	if len(requestLine) < 3 {
		return nil, 0, errors.New("incomplete request")
	}
	method := requestLine[0]
	target := requestLine[1]
	version := strings.Split(requestLine[2], "/")[1]

	for _, v := range method {
		if !unicode.IsUpper(v) || !unicode.IsLetter(v) {
			fmt.Println("invalid method")
			return nil, 0, errors.New("Failed to parse method")
		}
	}

	if version != "1.1" {
		fmt.Println("invalid version")
		return nil, 0, errors.New("Failed to parse http version")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, n, nil
}
