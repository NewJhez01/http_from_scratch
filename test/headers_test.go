package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithExtraWhitespace(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:           localhost:42069    \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)
}

func TestValid2HeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	headers["existing"] = "value"

	data := []byte("Host: localhost:42069\r\nContent-Type: text/html\r\n\r\n")

	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	data = data[n:]
	n, done, err = headers.Parse(data)
	assert.Equal(t, "text/html", headers["content-type"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	assert.Equal(t, "value", headers["existing"])
}

func TestValid2HeadersWithSameHeader(t *testing.T) {
	headers := NewHeaders()

	data := []byte("Set-Person: value1\r\nSet-Person: value2\r\nSet-Person: value3\r\n\r\n")
	n := 0
	for {
		num, d, err := headers.Parse(data[n:])
		if err != nil {
			t.Errorf("unexpected error")
		}
		n += num
		if d == true {
			break
		}
	}

	assert.Equal(t, len(data), n)
	assert.Equal(t, "value1, value2, value3", headers["set-person"])
}

func TestValidDone(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestInvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestInvalidChars(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	assert.Equal(t, fmt.Errorf("invalid header key"), err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
