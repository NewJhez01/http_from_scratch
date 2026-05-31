package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"http_from_scratch/internal/request"
	"http_from_scratch/internal/response"
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

func defaultHandler(w *response.ResponseWriter, r request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "bad request",
		}
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusInternalError,
			Message:    "unexpected error",
		}
	}
	err := w.WriteStatusLine(response.StatusOK)
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalError,
			Message:    "unexpected error",
		}
	}
	body := "request successfull well done the http server works\n"
	headers := response.GetDefaultHeaders(len([]byte(body)))
	err = w.WriteHeaders(headers)
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "bad request",
		}
	}
	err = w.WriteBody(body)
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "bad request",
		}
	}
	return nil
}
