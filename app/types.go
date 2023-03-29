package main

import (
	"fmt"
	"math"
)

type _respType int

const (
	Unknown      _respType = -1
	Number       _respType = 1
	SimpleString _respType = iota
	BulkString   _respType = iota
	RespList     _respType = iota
	RespErr      _respType = iota
	EOF          _respType = iota
)

func (rt _respType) String() string {
	switch rt {
	case Unknown:
		return "Unkown"
	case Number:
		return "Number"
	case SimpleString:
		return "SimpleString"
	case BulkString:
		return "BulkString"
	case RespList:
		return "List"
	case RespErr:
		return "Error"
	case EOF:
		return "EOF"
	default:
		return "undefined"
	}
}

type RESPPrimitive interface {
	Value() (RESPPrimitive, _respType)
	String() string
	Int() int
}

type Int int
type String string
type List struct {
	data *llnode[RESPPrimitive]
	size int
}

func (i Int) Value() (RESPPrimitive, _respType) {
	return i, Number
}

func (i Int) String() string {
	return fmt.Sprintf("%d", int(i))
}

func (i Int) Int() int {
	return int(i)
}

func (s String) Value() (RESPPrimitive, _respType) {
	return s, SimpleString
}

func (s String) String() string {
	return string(s)
}

func (s String) Int() int {
	return math.MinInt
}

func (l *List) Value() (RESPPrimitive, _respType) {
	return l, RespList
}

func (l List) String() string {
	return ""
}

func (l List) Int() int {
	return math.MinInt
}

func (l *List) Size() int {
	return l.size
}

func (l *List) Next() RESPPrimitive {
	if l.size == 0 {
		return nil
	}
	l.size -= 1
	retval := l.data.data
	l.data = l.data.n
	return retval
}
