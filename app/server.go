package main

import (
	"bytes"
	"fmt"

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
	parser := NewRESPParser(conn)
	for {

		list, err := parser.GetRequest()
		if err != nil {
			fmt.Printf("Error encountered while reading request: %v. Closing connection", err)
			return
		}
		if list.Size() <= 0 {
			fmt.Printf("Error encountered while reading request: %v. Closing connection", err)
			return
		}

		command := list.Next()
		switch command {
		case String("info"):
			conn.Write([]byte("+go-redis-test\r\n"))
		case String("echo"):
			message := list.Next().String()
			buf := bytes.NewBufferString("+")
			buf.WriteString(message)
			buf.WriteString("\r\n")
			conn.Write(buf.Bytes())
		case String("ping"):
			conn.Write([]byte("+pong\r\n"))
		}
	}
}
