package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	PING = "PING"
	ECHO = "ECHO"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func(conn net.Conn) {
			for {
				reader := NewReader(conn)

				line, err := reader.ReadLine()

				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Println("Error parsing command", err.Error())
					break
				}

				if line[0] != '*' {
					fmt.Println("Invalid RESP command. RESP command must be an array")
				}

				length, _ := strconv.Atoi(string(line[1:]))

				if length < 1 {
					fmt.Println("Invalid number of arguments for a command")
					break
				}

				cmd, err := reader.ReadBulkString()

				fmt.Println("Read command: ", cmd)

				if err != nil {
					fmt.Println("Unable to parse command", err.Error())
					break
				}

				arguments := make([]string, length-1)

				for i := 0; i < length-1; i++ {
					val, err := reader.ReadBulkString()
					if err != nil {
						fmt.Println("Unable to parse command", err.Error())
						break
					}
					arguments[i] = val
				}

				writer := NewWriter(conn)
				var command Command
				switch strings.ToUpper(cmd) {
				case PING:
					command, err = NewPingCommand(arguments, writer)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
				case ECHO:
					command, err = NewEchoCommand(arguments, writer)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
				}

				if command == nil {
					fmt.Println("Unsupported command")
					break
				}
				command.ExecuteCommand()
			}
			conn.Close()
		}(conn)
	}
}
