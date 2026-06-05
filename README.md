# go-http

HTTP/1.1 server from scratch in Go. One test dependency only.

Built to understand how raw TCP bytes become HTTP requests and responses without `net/http`.

## What it does

- **TCP listener**: Binds to a port, accepts connections, handles them concurrently with goroutines
- **Request parsing**: Reads raw bytes from the connection and parses the HTTP/1.1 request line, headers, and body
- **Response generation**: Constructs HTTP/1.1 responses with proper status line, headers, and body serialization
- **Graceful shutdown**: Handles `SIGTERM` to close the listener and finish active connections

## Project Structure

.
├── cmd/server/ # Entry point: starts the TCP listener
├── internal/
│ ├── server/ # TCP listener, connection loop, goroutine spawning
│ ├── http/ # HTTP protocol constants and shared types
│ ├── request/ # HTTP request parsing from raw bytes
│ └── response/ # HTTP response construction and serialization
├── test/ # Integration tests and test fixtures
├── go.mod
└── go.sum
plain

## Run

```bash
go run ./cmd/server

Test
bash

go test ./...

Dependencies

    testify — assertions and test utilities (dev only)

What I learned

    HTTP/1.1 is a text protocol on top of TCP, but parsing it correctly requires handling partial reads, malformed input, and connection edge cases
    Separating request parsing, response generation, and server logic into distinct packages makes the protocol easier to reason about and test
    Testing a raw server means dialing real TCP connections and asserting on byte-level responses
    Concurrent connection handling requires careful goroutine management and clean socket closure
```
