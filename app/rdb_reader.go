package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	REDIS_RDB_OPCODE_KEY      = 0x00
	REDIS_RDB_OPCODE_SELECTDB = 0xFE
)

func Read(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	key, err := parseRDB(file)
	return key, err
}

func parseRDB(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	for {
		opcode, err := reader.ReadByte()
		if err != nil {
			return "", err
		}

		switch opcode {
		case REDIS_RDB_OPCODE_SELECTDB:
			err := readDBMeta(reader)
			if err != nil {
				return "", err
			}
		case REDIS_RDB_OPCODE_KEY:
			log.Println("Found value 0. Reading string encoded key")
			l, err := readLength(reader)
			if err != nil {
				return "", nil
			}
			key, err := readKey(l, reader)
			if err != nil {
				return "", err
			}
			return key, nil
		default:
			// Do nothing
		}
	}
}

func readDBMeta(reader *bufio.Reader) error {
	dbnum, err := reader.Peek(2)
	if err != nil {
		return err
	}
	if dbnum[1] == 0xFB {
		log.Printf("Found dbnum = %v", dbnum[0])
		reader.Discard(2)
	}
	return nil
}

func readLength(reader *bufio.Reader) (int, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return -1, err
	}
	msBits := (b >> 6) & 0x03
	switch msBits {

	case 0:
		return int(b & 0x3F), nil

	case 1:
		b2, err := reader.ReadByte()
		if err != nil {
			return -1, err
		}
		b = b & 0x3F
		var ulen uint16
		err = binary.Read(bytes.NewReader([]byte{b, b2}), binary.LittleEndian, &ulen)
		if err != nil {
			return -1, err
		}
		return int(ulen), nil
	case 2:
		numAsBytes := make([]byte, 4)
		n, err := reader.Read(numAsBytes)
		if err != nil {
			return -1, err
		}
		if n != 4 {
			return -1, errors.New("could not read 4 byte length")
		}
		var ulen uint32
		err = binary.Read(bytes.NewReader(numAsBytes), binary.LittleEndian, &ulen)
		if err != nil {
			return -1, err
		}
		return int(ulen), nil

	case 3:
		//Hoping this wont occur
	}
	return 0, nil
}

func readKey(l int, reader *bufio.Reader) (string, error) {
	k := make([]byte, l)
	n, err := reader.Read(k)
	if err != nil {
		return "", err
	}
	if n != l {
		return "", fmt.Errorf("could not read %v bytes", l)
	}
	return string(k), nil
}
