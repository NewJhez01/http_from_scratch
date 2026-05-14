package headers

import (
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
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithExtraWhitespace(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:           localhost:42069    \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)
}

func TestValid2HeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	headers["Existing"] = "value"

	data := []byte("Host: localhost:42069\r\nContent-Type: text/html\r\n\r\n")

	// Parse first header
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Parse second header (advance past consumed bytes)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "text/html", headers["Content-Type"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Existing header still there
	assert.Equal(t, "value", headers["Existing"])
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
