package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	consumed := 0
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}
	if strings.HasPrefix(string(data), "\r\n") {
		return 2, true, nil
	}

	line := strings.Split(string(data), "\r\n")
	consumed += 2
	parts := strings.Split(string(line[0]), ":")
	consumed += 1
	key := parts[0]
	val := strings.Join(parts[1:], ":")

	if strings.HasPrefix(key, " ") || strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid header key")
	}
	trimmedVal := strings.Trim(val, " ")

	h[key] = trimmedVal

	return len(key) + len(val) + consumed, false, nil
}
