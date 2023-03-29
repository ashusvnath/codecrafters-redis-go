package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)

	t.Run("RESPParser.getType", func(t *testing.T) {
		t.Run("should parse number type and size correctly", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(`:1`))
			typ, size, err := p.getType()
			assert.Equal(typ, Number, "Expected type to be number")
			assert.Equal(-1, size, "Expected size to be 0")
			assert.Nil(err)
		})
		t.Run("should parse bulkString type and size correctly", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(`$4\r\nASDF\r\n`))
			typ, size, err := p.getType()
			assert.Equal(typ, BulkString, "Expected type to be bulkString")
			assert.Equal(4, size)
			assert.Nil(err)
		})

		t.Run("should parse simpleString type and size correctly", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(`+ASDF\r\n`))
			typ, size, err := p.getType()
			assert.Equal(typ, SimpleString, "Expected type to be simpleString")
			assert.Equal(-1, size)
			assert.Nil(err)
		})

		t.Run("should parse error type and size correctly", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(`-ASDF\r\n`))
			typ, size, err := p.getType()
			assert.Equal(typ, RespErr, "Expected type to be error")
			assert.Equal(-1, size)
			assert.Nil(err)
		})

		t.Run("should parse list type and size correctly", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(`*10\r\n$4\r\nping\r\n`))
			typ, size, err := p.getType()
			assert.Equal(typ, RespList, "Expected type to be list")
			assert.Equal(10, size)
			assert.Nil(err)
		})

		t.Run("should handle eof", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader(``))
			typ, size, err := p.getType()
			assert.Equal(typ, EOF, "Expected type to be unknown")
			assert.Equal(-1, size)
			assert.NoError(err)
		})
	})
	t.Run("RESPParser.GetRequest", func(t *testing.T) {
		t.Run("should parse single bulk string request", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader("*1\r\n$5\r\nhello\r\n"))
			parsed, err := p.GetRequest()
			assert.Nil(err)
			assert.Equal(1, parsed.Size())
			assert.Equal(String("hello"), parsed.Next())
		})
		t.Run("should parse single simple string request", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader("*1\r\n+hello\r\n"))
			parsed, err := p.GetRequest()
			assert.Nil(err)
			assert.Equal(1, parsed.Size())
			assert.Equal(String("hello"), parsed.Next())
		})

		t.Run("should parse simple string and integer request", func(t *testing.T) {
			p := NewRESPParser(strings.NewReader("*2\r\n+hello\r\n:1234\r\n"))
			parsed, err := p.GetRequest()
			assert.Nil(err)
			assert.Equal(2, parsed.Size())
			assert.Equal(String("hello"), parsed.Next())
			assert.Equal(Int(1234), parsed.Next())
		})
	})
}
