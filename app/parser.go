package main

import (
	"io"
	"text/scanner"
)

type _respType int

const (
	unknown      _respType = -1
	number       _respType = 1
	simpleString _respType = iota
	bulkString   _respType = iota
	list         _respType = iota
)

type respParser struct {
	scnr scanner.Scanner
}

type Request struct {
	args []string
}

func NewRESPParser(input io.Reader) *respParser {
	scnr := scanner.Scanner{}
	scnr.Init(input)
	return &respParser{scnr}
}

func (p *respParser) getType() _respType  {
	p.scnr.Mode = scanner.ScanChars
	switch p.scnr.Next() {
	case ':':
		return number
	case '+':
		return simpleString
	case '$':
		return bulkString
	case '*':
		return list
	default:
		return unknown
	}
}

func (p *respParser) Next(r io.Reader) *Request {
	p.getType()
	return nil
}
