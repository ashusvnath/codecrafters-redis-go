package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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

	data := make(map[string]string)
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
			//Won't implement COMMAND
			conn.Write([]byte("$-1\r\n"))
		case "echo":
			message := list.Next().String()
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)))
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "get":
			key := list.Next().String()
			result, ok := data[key]
			var replyString string
			if ok {
				replyString = fmt.Sprintf("$%d\r\n%s\r\n", len(result), result)
			} else {
				replyString = "$-1\r\n" //Null bulk string
			}
			conn.Write([]byte(replyString))
		case "set":
			key := list.Next().String()
			val := list.Next().String()
			if list.Size() == 2 {
				option := strings.ToLower(list.Next().String())
				if option == "px" {
					timeoutMilliSecondString := list.Next().String()
					timeoutMilliSeconds, err := strconv.Atoi(timeoutMilliSecondString)
					if err == nil {
						fmt.Printf("Scheduling deletion of key %s in %d milliseconds", key, timeoutMilliSeconds)
						time.AfterFunc(time.Millisecond*time.Duration(timeoutMilliSeconds), func() {
							fmt.Printf("Deleting key %s", key)
							delete(data, key)
							fmt.Println("...Done")
						})
					} else {
						message := fmt.Sprintf("Error parsing timeout value '%s', %v", timeoutMilliSecondString, err)
						fmt.Println(message)
						conn.Write([]byte("-" + message + "\r\n"))
						continue
					}
				} else {
					message := fmt.Sprintf("Error unknown option '%s'", option)
					fmt.Println(message)
					conn.Write([]byte("-" + message + "\r\n"))
					continue
				}
			}
			data[key] = val
			conn.Write([]byte("+OK\r\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("-Error uknown command %s\r\n", command)))
		}
	}
}
