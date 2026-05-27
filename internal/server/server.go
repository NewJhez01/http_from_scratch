package server

import (
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
}

type HandlerError struct {
	statusCode string
	message    string
}

type Handler func(io.Writer, request.Request) *HandlerError

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		l,
		atomic.Bool{},
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
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
}
