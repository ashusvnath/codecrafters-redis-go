package main

import "fmt"

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

type Config map[string]string

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
	sl.size--
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
