package main

import (
	"bufio"
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

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(rd),
	}
}

func (r *Reader) ReadLine() ([]byte, error) {
	b, err := r.reader.ReadSlice('\n')

	if err != nil {
		return nil, err
	}
	// check if it is a valid line
	// TODO: is \n not expected anywhere else?
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("invalid command")
	}
	// excluding \r\n
	return b[:len(b)-2], nil
}

func (r *Reader) ReadBulkString() (string, error) {
	line, err := r.ReadLine()
	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return "", fmt.Errorf("Error parsing argument:%s", err.Error())
	}
	b := make([]byte, length+2)
	_, err = io.ReadFull(r.reader, b)
	if err != nil {
		return "", err
	}
	return string(b[:length]), nil
}

type Writer struct {
	writer *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
	}
}

func (w Writer) WriteResponseString(resp string) {
	response := "+" + resp + "\r\n"
	w.writer.Write([]byte(response))
	w.writer.Flush()
}

func (w Writer) WriteBulkResponseString(resp string) {
	n := len(resp)
	response := "$" + strconv.Itoa(n) + "\r\n" + resp + "\r\n"
	w.writer.Write([]byte(response))
	fmt.Println(response)
	w.writer.Flush()
}

func (w Writer) WriteErrorResponseString(resp string) {
	response := "-" + resp + "\r\n"
	w.writer.Write([]byte(response))
	w.writer.Flush()
}

type Command interface {
	ExecuteCommand()
}

type PingCommand struct {
	writer *Writer
	args   []string
}

func NewPingCommand(args []string, writer *Writer) (*PingCommand, error) {
	if len(args) > 2 {
		writer.WriteErrorResponseString("Invalid number of arguments for ping command")
		return nil, fmt.Errorf("Ping Command accepts only one argument")
	}
	return &PingCommand{args: args, writer: writer}, nil
}

func (p PingCommand) ExecuteCommand() {
	if len(p.args) == 0 {
		p.writer.WriteResponseString("PONG")
		return
	}
	p.writer.WriteBulkResponseString(p.args[0])
}

type EchoCommand struct {
	writer *Writer
	args   []string
}

func NewEchoCommand(args []string, writer *Writer) (*EchoCommand, error) {
	if len(args) != 1 {
		writer.WriteErrorResponseString("Invalid number of arguments for Echo command")
		return nil, fmt.Errorf("Echo Command accepts only one argument")
	}
	return &EchoCommand{args: args, writer: writer}, nil
}

func (p EchoCommand) ExecuteCommand() {
	p.writer.WriteBulkResponseString(p.args[0])
}
