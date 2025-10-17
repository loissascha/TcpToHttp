package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo: barbar  \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, found := headers.Get("Host")
	assert.True(t, found)
	foofoo, found := headers.Get("Foofoo")
	assert.True(t, found)
	missingKey, found := headers.Get("MissingKey")
	assert.False(t, found)
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "barbar", foofoo)
	assert.Equal(t, "", missingKey)
	assert.Equal(t, 43, n)
	assert.True(t, done)

	// valid multi header
	headers = NewHeaders()
	data = []byte("Set-Person: SomePerson1  \r\nSet-Person:   AnotherPerson\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	setPerson, found := headers.Get("set-person")
	assert.True(t, found)
	assert.Equal(t, "SomePerson1,AnotherPerson", setPerson)
	assert.Equal(t, 58, n)
	assert.True(t, done)

	// Test: Invalid header name
	headers = NewHeaders()
	data = []byte("Höst: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
