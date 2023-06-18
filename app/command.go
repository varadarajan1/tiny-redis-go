package main

import (
	"fmt"
	"sync"
)

var hashmap = sync.Map{}

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
		return nil, fmt.Errorf("ping command accepts only one argument")
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
		return nil, fmt.Errorf("echo command accepts only one argument")
	}
	return &EchoCommand{args: args, writer: writer}, nil
}

func (p EchoCommand) ExecuteCommand() {
	p.writer.WriteBulkResponseString(p.args[0])
}

type SetCommand struct {
	writer *Writer
	args   []string
}

func NewSetCommand(args []string, writer *Writer) (*SetCommand, error) {
	if len(args) != 2 {
		writer.WriteErrorResponseString("Invalid number of arguments for Set command")
		return nil, fmt.Errorf("set command accepts only one argument")
	}
	return &SetCommand{args: args, writer: writer}, nil
}

func (p SetCommand) ExecuteCommand() {
	hashmap.Store(p.args[0], p.args[1])
	p.writer.WriteResponseString("OK")
}

type GetCommand struct {
	writer *Writer
	args   []string
}

func NewGetCommand(args []string, writer *Writer) (*GetCommand, error) {
	if len(args) != 1 {
		writer.WriteErrorResponseString("Invalid number of arguments for Get command")
		return nil, fmt.Errorf("get command accepts only one argument")
	}
	return &GetCommand{args: args, writer: writer}, nil
}

func (p GetCommand) ExecuteCommand() {
	if value, ok := hashmap.Load(p.args[0]); ok {
		p.writer.WriteBulkResponseString((value).(string))
		return
	}
	p.writer.WriteNilResponsString()
}
