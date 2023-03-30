package main

import (
	"bytes"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	fmt.Println("Connection successful!!")
	parser := NewSimpleParser(conn)
	for {
		val, err := parser.Next()
		fmt.Printf("Received command %s\n", val)
		if err != nil {
			errorString := fmt.Sprintf("ERROR encountered while reading request: %v. Closing connection", err)
			fmt.Printf(errorString + "\n")
			errorString = "-" + errorString + "\r\n"
			conn.Write([]byte(errorString))
			return
		}
		list, ok := val.val.(*SList)
		if !ok || list.Size() == 0 {
			errorString := "ERROR 0 argument list"
			fmt.Printf(errorString + "\n")
			errorString = "-" + errorString + "\r\n"
			conn.Write([]byte(errorString))
			return
		}

		command := strings.ToLower(list.Next().String())
		switch command {
		case "command":
			// subcommand := strings.ToLower(list.Next().String())
			// if subcommand == "docs" {
			// 	conn.Write([]byte("*3\r\n$4\r\nPING\r\n$4\r\nECHO\r\n$12\r\nCOMMAND DOCS\r\n"))
			// }
			// default:
			conn.Write([]byte("$-1\r\n"))
		case "echo":
			message := list.Next().String()
			buf := bytes.NewBufferString("+")
			buf.WriteString(message)
			buf.WriteString("\r\n")
			conn.Write(buf.Bytes())
		case "ping":
			conn.Write([]byte("$4\r\nPONG\r\n"))
		}
	}
}
