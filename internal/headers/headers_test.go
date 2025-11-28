package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParsing(t *testing.T) {
	// Test: Valid single header

	h := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	h = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	h = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, "curl/7.81.0", h["User-Agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	h = NewHeaders()
	data = []byte("\r\n a bunch of other stuff")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Empty(t, h)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	h = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// "Valid single header with extra whitespace"

	h = map[string]string{"host": "localhost:42069"}
	data = []byte("auth:sfs4392\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, "sfs4392", h["auth"])
	assert.Equal(t, 14, n)
	assert.False(t, done)

}
