package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"http_from_scratch/internal/headers"
)

const (
	bufferSize = 8
)

type state int

const (
	requestLine state = 0
	header      state = 1
	body        state = 2
	done        state = 3
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	status      state
	Body        []byte
}

func CreateNewRequest() *Request {
	return &Request{}
}

func RequestFromReader(r io.Reader) (Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	req := Request{status: requestLine}
	for req.status != done {
		if req.status == 2 &&
			(req.Headers.Get("content-length") == "" ||
				req.Headers.Get("content-length") == "0") {
			req.status = 3
			continue
		}
		n, err := r.Read(buffer[readToIndex:])
		if err != nil {
			fmt.Println("read error ", err.Error())
		}
		if err == io.EOF && n == 0 && readToIndex == 0 {
			req.status = 3
		}
		readToIndex += n
		if readToIndex == cap(buffer) {
			tmp := make([]byte, len(buffer)*2)
			copy(tmp, buffer)
			buffer = tmp
		}
		p, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return req, fmt.Errorf("failed to parse request prev: %s", err.Error())
		}
		tmp := make([]byte, len(buffer)-p)
		copy(tmp, buffer[p:])
		buffer = tmp
		readToIndex -= p
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.status {
	case requestLine:
		rql, bytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, errors.New("error")
		}
		if bytesParsed == 0 {
			return 0, nil
		}
		r.RequestLine = *rql
		r.status = header
		r.Headers = headers.NewHeaders()
		return bytesParsed, nil
	case header:
		n := 0
		for {
			consumed, done, err := r.Headers.Parse(data[n:])
			if err != nil {
				return 0, errors.New("failed to parse headers")
			}
			n += consumed
			if done == true {
				r.status = body
				return n, nil
			}

			if consumed == 0 {
				return n, nil
			}
		}
	case body:
		r.Body = append(r.Body, data...)
		return len(data), nil
	case done:
		contentHeader := r.Headers["content-length"]
		contentLen := 0
		if contentHeader != "" {
			n, err := strconv.Atoi(r.Headers["content-length"])
			if err != nil {
				fmt.Println("error ", err.Error())
			}
			contentLen = n
		}

		if contentLen != len(r.Body) {
			return 0, errors.New("content length doesn't match the header def")
		}
		return 0, nil
	}

	return 0, nil
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
