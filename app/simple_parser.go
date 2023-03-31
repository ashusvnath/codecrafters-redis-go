package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type SimpleParser struct {
	data io.Reader
}



func NewSimpleParser(conn io.Reader) *SimpleParser {
	sp := &SimpleParser{conn}
	return sp
}

func (sp *SimpleParser) Next() (Value, error) {
	return next(bufio.NewReader(sp.data))
}

func next(bufReader *bufio.Reader) (result Value, err error) {
	var byt byte
	byt, err = bufReader.ReadByte()
	switch byt {
	case ':':
		result.typ = Number
		result.val, err = readInt(bufReader)
		return
	case '$':
		return readBulkString(bufReader)
	case '*':
		return readList(bufReader)
	case '+':
		return readString(bufReader, SimpleString)
	case '-':
		return readString(bufReader, RespErr)
	}
	return
}

func readInt(reader *bufio.Reader) (val int, err error) {
	var line string
	line, err = reader.ReadString('\r')
	if err != nil {
		return
	}
	_, err = reader.ReadByte()
	if err != nil {
		return
	}
	val, err = strconv.Atoi(line[:len(line)-1])
	return
}

func readList(reader *bufio.Reader) (result Value, err error) {
	result.typ = RespList
	var siz int
	siz, err = readInt(reader)
	if err != nil {
		return
	}
	l := &SList{}
	for siz > 0 {
		var val Value
		val, err = next(reader)
		l.Append(val)
		siz--
	}
	result.val = l
	return
}

func readString(reader *bufio.Reader, typ _respType) (result Value, err error) {
	result.typ = typ
	var line string
	line, err = reader.ReadString('\r')
	if err != nil {
		return
	}
	_, err = reader.ReadByte()
	if err != nil {
		return
	}
	result.val = line[:len(line)-1]
	return
}

func readBulkString(bufReader *bufio.Reader) (result Value, err error) {
	result.typ = BulkString
	var siz int
	siz, err = readInt(bufReader)
	if err != nil {
		return
	}
	bytesToRead := make([]byte, siz+2)
	var n int
	n, err = io.ReadFull(bufReader, bytesToRead)
	if n < siz+2 || err != nil {
		err = fmt.Errorf("could not read bulk string of length %d: %v", siz, err)
		return
	}
	result.val = string(bytesToRead[:len(bytesToRead)-2])
	return
}
