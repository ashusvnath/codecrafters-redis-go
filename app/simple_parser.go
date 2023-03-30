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

type SList struct {
	data *llnode[Value]
	size int
}

func (sl *SList) Append(v Value) {
	sl.size++
	sl.data = sl.data.Append(v)
}

func (sl *SList) Next() Value {
	result := sl.data.data
	sl.data = sl.data.n
	return result
}

func (sl *SList) Size() int {
	return sl.size
}

func (sl *SList) String() string {
	return fmt.Sprintf("List(%d)[%s]", sl.size, sl.data.String())
}

type Value struct {
	typ _respType
	val interface{}
}

func (v Value) String() string {
	if v.typ == SimpleString || v.typ == BulkString || v.typ == RespErr {
		str, _ := v.val.(string)
		return str
	}
	if v.typ == RespList {
		return fmt.Sprintf("%s", v.val)
	}
	return ""
}

func (v Value) Int() int {
	if v.typ == Number {
		i, _ := v.val.(int)
		return i
	}
	return 0
}

func (v Value) List() *SList {
	if v.typ == RespList {
		l, _ := v.val.(*SList)
		return l
	}
	return nil
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
		err = fmt.Errorf("could not read bulk string of length %d: %v", n, err)
		return
	}
	result.val = string(bytesToRead[:len(bytesToRead)-2])
	return
}
