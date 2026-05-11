package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	status      int
}

// init = 0 done = 1
func RequestFromReader(r io.Reader) (Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	req := Request{status: 0}
	for req.status == 0 {
		n, err := r.Read(buffer[readToIndex:])
		if err == io.EOF && n == 0 {
			req.status = 1
			break
		}
		readToIndex += n
		if readToIndex == cap(buffer) {
			tmp := make([]byte, len(buffer)*2)
			copy(tmp, buffer)
			buffer = tmp
		}
		p, err := req.parse(buffer)
		if err != nil {
			fmt.Println("error when parsing")
		}
		tmp := make([]byte, len(buffer)-p)
		copy(tmp, buffer[p:])
		buffer = tmp
		readToIndex -= p
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.status == 1 {
		return 0, errors.New("already done")
	}
	rql, bytesParsed, err := parseRequestLine(data)
	if err != nil {
		return 0, errors.New("error")
	}
	if bytesParsed == 0 {
		return 0, nil
	}
	r.RequestLine = *rql
	r.status = 1
	return bytesParsed, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
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
	}, len(l[0]) + 2, nil
}
