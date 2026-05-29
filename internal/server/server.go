package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"http_from_scratch/internal/request"
	"http_from_scratch/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode string
	Message    string
}

type Handler func(io.Writer, request.Request) *HandlerError

func writeError(w io.Writer, h HandlerError) {
	fmt.Fprint(w, h.StatusCode, "\r\n", h.Message)
}

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		l,
		atomic.Bool{},
		h,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Println("server failed prev: " + err.Error())
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("failed to parse request")
	}
	b := bytes.NewBuffer([]byte{})
	hErr := s.handler(b, req)
	if hErr != nil {
		writeError(conn, *hErr)
	}

	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
	h := response.GetDefaultHeaders(len(req.Body))
	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, b)
}
