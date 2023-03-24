package main

import (
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
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connection successful!!")
	for {
		input := make([]byte, 100)
		n, err := conn.Read(input)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			if err.Error() == "EOF"{
				fmt.Println("Closing connection and exiting.")
				os.Exit(0)
			}
			os.Exit(1)
		}

		inputString := string(input[:n])
		fmt.Printf("C: %#v\n", inputString)
		if inputString == "*1\r\n$4\r\nping\r\n" {
			response := "+PONG\r\n"
			conn.Write([]byte(response))
			fmt.Printf("S: %#v\n", response)
		}
	}
}
