package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParser_Next(t *testing.T) {
	assert := assert.New(t)
	t.Run("should parse RESP number", func(t *testing.T) {
		input := strings.NewReader(":1234\r\n")
		sp := NewSimpleParser(input)
		v, err := sp.Next()
		assert.NoError(err)
		assert.Equal(Number, v.typ, fmt.Sprintf("Expected Number got %s", v.typ))
		assert.Equal(1234, v.val)
	})

	t.Run("should parse RESP simple string", func(t *testing.T) {
		input := strings.NewReader("+Just a simple string\r\n")
		sp := NewSimpleParser(input)
		v, err := sp.Next()
		assert.NoError(err)
		assert.Equal(SimpleString, v.typ, fmt.Sprintf("Expected SimpleString got %s", v.typ))
		assert.Equal("Just a simple string", v.val)
	})

	t.Run("should parse RESP error", func(t *testing.T) {
		input := strings.NewReader("-Just a simple error\r\n")
		sp := NewSimpleParser(input)
		v, err := sp.Next()
		assert.NoError(err)
		assert.Equal(RespErr, v.typ, fmt.Sprintf("Expected RespErr got %s", v.typ))
		assert.Equal("Just a simple error", v.val)
	})

	t.Run("should parse RESP bulk string", func(t *testing.T) {
		input := strings.NewReader("$12\r\none two four\r\n")
		sp := NewSimpleParser(input)
		v, err := sp.Next()
		assert.NoError(err)
		assert.Equal(BulkString, v.typ, fmt.Sprintf("Expected BulkString got %s", v.typ))
		assert.Equal("one two four", v.val)
	})

	t.Run("should parse RESP list", func(t *testing.T) {
		input := strings.NewReader("*4\r\n$12\r\none two four\r\n+five six\r\n-seven eight\r\n:910\r\n")
		sp := NewSimpleParser(input)
		v, err := sp.Next()
		assert.NoError(err)
		assert.Equal(RespList, v.typ, fmt.Sprintf("Expected List got %s", v.typ))
		l, ok := v.val.(*SList)
		assert.True(ok)
		assert.Equal(4, l.Size())
	})
}
