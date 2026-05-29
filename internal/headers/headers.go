package headers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[key]
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}
	if strings.HasPrefix(string(data), "\r\n") {
		return 2, true, nil
	}

	line := strings.Split(string(data), "\r\n")
	parts := strings.Split(string(line[0]), ":")
	key := parts[0]
	val := strings.Join(parts[1:], ":")

	if strings.HasPrefix(key, " ") || strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid header key")
	}

	if validate(key) == false {
		return 0, false, errors.New("invalid header key")
	}

	trimmedVal := strings.Trim(val, " ")
	if validate(trimmedVal) == false {
		return 0, false, errors.New("invalid format for val")
	}

	keyArr := strings.Split(key, "")
	for i, v := range keyArr {
		keyArr[i] = strings.ToLower(v)
	}

	valArr := strings.Split(trimmedVal, "")
	for i, v := range valArr {
		valArr[i] = strings.ToLower(v)
	}

	header := strings.Join(keyArr, "")
	value := strings.Join(valArr, "")
	if h[header] != "" {
		h[header] = fmt.Sprintf("%s, %s", h[header], value)
		fmt.Println("header ", h[header])
	} else {
		h[header] = value
		fmt.Println("header ", h[header])
	}

	return len(line[0]) + 2, false, nil
}

func validate(v string) bool {
	m, err := regexp.Match("^[a-zA-Z0-9_!/#$%&'*+.^_`|~:-]*$", []byte(v))
	if m == false || err != nil {
		return false
	}
	return true
}
