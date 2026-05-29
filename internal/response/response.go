package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"http_from_scratch/internal/headers"
)

type (
	StatusCode  int
	WriterState int
)

type ResponseWriter struct {
	status WriterState
	writer io.Writer
}

const (
	StatusOK            StatusCode  = 200
	StatusBadRequest    StatusCode  = 400
	StatusInternalError StatusCode  = 500
	statusLine          WriterState = 0
	header              WriterState = 1
	body                WriterState = 3
)

func CreateNewWriter(w io.Writer) *ResponseWriter {
	return &ResponseWriter{
		statusLine,
		w,
	}
}

func (w *ResponseWriter) WriteStatusLine(sc StatusCode) error {
	switch sc {
	case 200:
		_, err := w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case 400:
		_, err := w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case 500:
		_, err := w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	}
	w.status = header
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func (w *ResponseWriter) WriteHeaders(h headers.Headers) error {
	if w.status != header {
		return errors.New("requestline missing")
	}
	for k, v := range h {
		_, err := fmt.Fprint(w.writer, k, ": ", v, "\r\n")
		if err != nil {
			return err
		}
	}
	w.status = body
	return nil
}

func (w *ResponseWriter) WriteBody(b string) error {
	if w.status != body {
		return errors.New("headers or request line missing")
	}
	_, err := w.writer.Write([]byte(b))
	if err != nil {
		return errors.New("couldn't write to body")
	}
	return nil
}
