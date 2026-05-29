package server

import (
	"bytes"
	"fmt"
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
	StatusCode response.StatusCode
	Message    string
}

type Handler func(*response.ResponseWriter, request.Request) *HandlerError

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

func writeError(w *response.ResponseWriter, h HandlerError) {
	err := w.WriteStatusLine(h.StatusCode)
	if err != nil {
		logFatal(h.Message)
	}
	headers := response.GetDefaultHeaders(len(h.Message))
	err = w.WriteHeaders(headers)
	if err != nil {
		logFatal(h.Message)
	}
	err = w.WriteBody(h.Message)
	if err != nil {
		logFatal(h.Message)
	}
}

func logFatal(message string) {
	log.Fatalf("total failure prev: %s", message)
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("failed to parse request")
	}
	b := response.CreateNewWriter(bytes.NewBuffer([]byte{}))
	hErr := s.handler(b, req)
	w := response.CreateNewWriter(conn)
	if hErr != nil {
		writeError(w, *hErr)
	}

	err = w.WriteStatusLine(200)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
	h := response.GetDefaultHeaders(len(req.Body))
	err = w.WriteHeaders(h)
	if err != nil {
		log.Fatalf("unexpected error prev: %s", err.Error())
	}
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, b)
}
