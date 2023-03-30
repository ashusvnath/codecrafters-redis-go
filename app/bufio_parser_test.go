package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufioRESPParser_GetInput(t *testing.T) {
	assert := assert.New(t)
	t.Run("should return all inputs split by \\r\\n", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		brp := NewBufioRESPParser(buf)
		buf.WriteString(":1234\r\n")
		input := brp.GetInput()
		assert.Equal(1, len(input))
		assert.Equal(":1234", input[0])
	})

	t.Run("should return all inputs split by \\r\\n", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		brp := NewBufioRESPParser(buf)
		buf.WriteString(":1234\r\n")
		input := brp.GetInput()
		assert.Equal(1, len(input))
		assert.Equal(":1234", input[0])
	})
}

func TestSscanf(t *testing.T) {
	assert := assert.New(t)
	t.Run("no error on proper scan of number", func(t *testing.T) {
		input := ":1234\r\n"
		var number int
		var typ rune
		n, err := fmt.Sscanf(input, "%c%d\r\n", &typ, &number)
		assert.NoError(err)
		assert.Equal(2, n)
		assert.Equal(rune(':'), typ)
		assert.Equal(1234, number)
	})

	t.Run("error on proper scan of number", func(t *testing.T) {
		input := ":asdf\r\n"
		var number int
		var typ rune
		n, err := fmt.Sscanf(input, "%c%d\r\n", &typ, &number)
		assert.Error(err)
		assert.Equal(1, n)
		assert.Equal(rune(':'), typ)
	})

	t.Run("no error on proper scan of bulk string", func(t *testing.T) {
		input := "$12\r\nASDFASDFASDF\r\nASDF"
		var number int
		var str string
		n, err := fmt.Sscanf(input, "$%d\r\n%s\r\n", &number, &str)
		assert.NoError(err)
		assert.Equal(2, n)
		assert.Equal(12, number)
		assert.Equal("ASDFASDFASDF", str)
	})

	t.Run("no error on proper scan of simple string", func(t *testing.T) {
		input := "+ASDFASDFASDF\r\nASDF"
		var str string
		n, err := fmt.Sscanf(input, "+%s\r\n", &str)
		assert.NoError(err)
		assert.Equal(1, n)
		assert.Equal("ASDFASDFASDF", str)
	})

	t.Run("no error on proper scan of error string", func(t *testing.T) {
		input := "-ASDFASDFASDF\r\n"
		var str string
		n, err := fmt.Sscanf(input, "-%s\r\n", &str)
		assert.NoError(err)
		assert.Equal(1, n)
		assert.Equal("ASDFASDFASDF", str)
	})


	t.Run("no error on proper scan of list with size", func(t *testing.T) {
		input := "*1234\r\n"
		var number int
		n, err := fmt.Sscanf(input, "*%d\r\n", &number)
		assert.NoError(err)
		assert.Equal(1, n)
		assert.Equal(1234, number)
	})
}
