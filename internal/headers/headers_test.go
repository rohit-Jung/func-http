package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldLineParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	fmt.Println("n, done, error", n, done, err, headers)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "", headers["MissingKey"])
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: valid multiple headers one with extraspace
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nNurse:        registered     \r\n\r\n")
	n, done, err = headers.Parse(data)
	fmt.Println("n, done, error", n, done, err, headers)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "registered", headers["nurse"])
	assert.Equal(t, "", headers["MissingKey"])
	assert.Equal(t, 56, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: multiple headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nNurse:        registered     \r\nNurse:   with doctor \r\nNurse: and other nurse\r\n\r\n")
	n, done, err = headers.Parse(data)
	fmt.Println("n, done, error", n, done, err, headers)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "registered,with doctor,and other nurse", headers["nurse"])
	assert.Equal(t, 103, n)
	assert.True(t, done)
}
