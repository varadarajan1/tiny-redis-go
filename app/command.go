package main

import "fmt"

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
