package main

import (
	"bufio"
	"io"
)

type ParserState int

type BufioRESPParser struct {
	scanner          bufio.Scanner
	requests         *Queue[List]
	connectionClosed bool
}

func (brp *BufioRESPParser) process(in []byte, isEof bool) (advance int, token []byte, err error) {
	if isEof {
		brp.connectionClosed = true
	}
	return
}

func NewBufioRESPParser(in io.Reader) *BufioRESPParser {
	brp := &BufioRESPParser{
		scanner:          *bufio.NewScanner(in),
		requests:         NewQueue[List](),
		connectionClosed: false,
	}
	brp.scanner.Split(brp.process)
	return brp
}

func (brp *BufioRESPParser) GetInput() []string {
	return []string{}
}
