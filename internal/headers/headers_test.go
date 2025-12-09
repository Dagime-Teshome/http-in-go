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
	value, _ := h.Get("Host")
	assert.Equal(t, "localhost:42069", value)
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
	value, _ = h.Get("User-Agent")
	assert.Equal(t, "curl/7.81.0", value)
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
	// Test for invalid key values
	h = make(Headers)
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test for case insensitivity
	h = make(Headers)
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	value, _ = h.Get("Host")
	assert.Equal(t, "localhost:42069", value)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header
	h = map[string]string{"host": "initialValue"}
	data = []byte("Host: anotherValue\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	value, _ = h.Get("Host")
	assert.Equal(t, "initialValue,anotherValue", value)
	assert.Equal(t, 20, n)
	assert.False(t, done)

}
