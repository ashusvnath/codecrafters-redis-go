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
	REDIS_RDB_OPCODE_AUX      = 0xFA
)

func RDB_Read(filepath string) (string, error) {
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
	reader.Discard(9) //REDIS0010
	for {
		opcode, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		log.Printf("OPCODE %x", opcode)

		switch opcode {
		case REDIS_RDB_OPCODE_AUX:
			log.Printf("Reading AUX info")
			l, _ := RDB_readLength(reader)
			auxKey, _ := RDB_readString(l, reader) //discard aux key
			l, _ = RDB_readLength(reader)
			auxValue, _ := RDB_readString(l, reader) //discard aux value
			log.Printf("Aux info: %v:%v", auxKey, auxValue)

		case REDIS_RDB_OPCODE_SELECTDB:
			log.Printf("Reading DB info")
			err := RDB_readMeta(reader)
			if err != nil {
				return "", err
			}

		case REDIS_RDB_OPCODE_KEY:
			log.Println("Found value 0. Reading string encoded key")
			l, err := RDB_readLength(reader)
			if err != nil {
				return "", nil
			}
			key, err := RDB_readString(l, reader)
			if err != nil {
				return "", err
			}
			return key, nil

		default:
			// Do nothing
		}
	}
}

func RDB_readMeta(reader *bufio.Reader) error {
	dbnum, err := reader.Peek(2)
	if err != nil {
		return err
	}
	if dbnum[1] == 0xFB {
		log.Printf("Found dbnum = %v", dbnum[0])
		reader.Discard(2)
	}
	RDB_readLength(reader)
	RDB_readLength(reader)
	return nil
}

func RDB_readLength(reader *bufio.Reader) (int, error) {
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
		switch int(b & 0x3F) {
		case 0:
			return 1, nil
		case 1:
			return 2, nil
		case 2:
			return 4, nil
		default:
			return -1, errors.New("unknown sized interger as string")
		}
	}
	return 0, nil
}

func RDB_readString(l int, reader *bufio.Reader) (string, error) {
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

func RDB_readIntAsString(l int, reader *bufio.Reader) (string, error) {
	return "", nil
}

func main1() {
	dir, _ := os.Getwd()
	log.Printf("working directory is %s", dir)
	key, err := RDB_Read("dump.rdb")
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	log.Printf("Key: %v", key)
}
