package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

type respParser struct {
	scnr      scanner.Scanner
	scanError error
}

func NewRESPParser(input io.Reader) *respParser {
	scnr := scanner.Scanner{}
	scnr.Init(input)
	parser := &respParser{scnr, nil}
	scnr.Whitespace = 1 << 11
	scnr.Error = func(s *scanner.Scanner, msg string) {
		parser.scanError = errors.New(msg)
	}
	return parser
}

func (p *respParser) getType() (r _respType, size int, parseError error) {
	p.scnr.Mode = scanner.ScanChars
	p.scanError = nil
repeat:
	chr := p.scnr.Scan()
	switch chr {
	case ':':
		r = Number
	case '+':
		r = SimpleString
	case '$':
		r = BulkString
	case '*':
		r = RespList
	case '-':
		r = RespErr
	case '\n':
		goto repeat
	case '\r':
		goto repeat
	case scanner.EOF:
		r = EOF
	default:
		pos := p.scnr.Pos()
		parseError = fmt.Errorf("found %c during type recognition at %v", chr, pos)
	}

	if p.scanError != nil {
		parseError = p.scanError
		p.scanError = nil
	}
	if parseError != nil {
		return
	}

	if r == Number || r == SimpleString || r == RespErr {
		size = -1
		return
	}

	nextInputChar := p.scnr.Peek()
	if nextInputChar < '0' || nextInputChar > '9' {
		size = -1
		return
	}

	p.scnr.Mode = scanner.ScanInts
	p.scnr.Scan()
	size, parseError = strconv.Atoi(p.scnr.TokenText())
	if p.scanError != nil {
		parseError = p.scanError
		p.scanError = nil
	}
	return
}

func (p *respParser) readInteger() (Int, error) {
	p.scnr.Mode = scanner.ScanInts
	p.scnr.Scan()
	if p.scanError != nil {
		err := p.scanError
		p.scanError = nil
		return 0, err
	}
	value, err := strconv.Atoi(p.scnr.TokenText())

	if err != nil {
		return 0, err
	}
	return Int(value), nil
}

func (p *respParser) readBulkString(size int) (String, error) {
	p.scnr.Mode = scanner.ScanChars
	value := []rune{}
	for ; size > 0; size-- {
		chr := p.scnr.Scan()
		if p.scanError != nil {
			err := p.scanError
			p.scanError = nil
			return "", err
		}
		value = append(value, chr)
	}
	return String(value), nil
}

func (p *respParser) readSimpleString() (String, error) {
	p.scnr.Mode = scanner.ScanChars
	value := []rune{}
	for chr := p.scnr.Next(); chr!= rune('\r') && chr != rune('\n') && chr != scanner.EOF; chr = p.scnr.Next(){
		if p.scanError != nil {
			err := p.scanError
			p.scanError = nil
			return "", err
		}
		value = append(value, chr)
	}
	return String(value), nil
}

func (p *respParser) readList(size int) (l *List, err error) {
	var list *llnode[RESPPrimitive] = nil
	l = &List{nil, size}
loop:
	for ; size > 0; size-- {
		typ, innnerSize, _err := p.getType()
		if _err != nil {
			l.data = list
			err = _err
			break loop
		}
		switch typ {
		case Number:
			v, _err := p.readInteger()
			if _err != nil {
				l.data = list
				err = _err
				break loop
			}
			l.data = l.data.Append(v)

		case SimpleString:
			v, _err := p.readSimpleString()
			if _err != nil {
				l.data = list
				err = _err
				break loop
			}
			l.data = l.data.Append(v)

		case BulkString:
			v, _err := p.readBulkString(innnerSize)
			if _err != nil {
				l.data = list
				err = _err
				break loop
			}
			l.data = l.data.Append(v)
		case RespErr:
			v, _err := p.readSimpleString()
			if _err != nil {
				l.data = list
				err = _err
				break loop
			}
			l.data = l.data.Append(v)
		case RespList:
			v, _err := p.readList(innnerSize)
			if _err != nil {
				l.data = list
				err = _err
				break loop
			}
			l.data = l.data.Append(v)
		case Unknown:
			err = errors.New("unknown type during RESP deserialization")
			break loop
		}
	}

	return
}

func (p *respParser) GetRequest() (*List, error) {
	typ, size, err := p.getType()
	if err != nil {
		return nil, err
	}
	if typ != RespList {
		return nil, fmt.Errorf("expected a list got %v", typ)
	}
	return p.readList(size)
}
