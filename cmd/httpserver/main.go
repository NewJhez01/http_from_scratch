package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"http_from_scratch/internal/request"
	"http_from_scratch/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, defaultHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func defaultHandler(w io.Writer, r request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: "400",
			Message:    "bad request",
		}
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: "500",
			Message:    "unexpected error",
		}
	}
	w.Write([]byte("status: ok"))
	return nil
}
