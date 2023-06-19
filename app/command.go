package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var hashmap = sync.Map{}

type CacheEntry struct {
	value  string
	expiry int64
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
	if len(args) < 2 || len(args) > 4 {
		writer.WriteErrorResponseString("Invalid number of arguments for Set command")
		return nil, fmt.Errorf("set command accepts only one argument")
	}
	return &SetCommand{args: args, writer: writer}, nil
}

func (p SetCommand) ExecuteCommand() {
	entry := CacheEntry{
		value:  p.args[1],
		expiry: -1,
	}
	if len(p.args) > 2 && strings.ToUpper(p.args[2]) == "PX" {
		e, err := strconv.Atoi(p.args[3])
		if err != nil {
			p.writer.WriteErrorResponseString("Invalid expiry value")
		}
		entry.expiry = time.Now().UnixMilli() + int64(e)
	}

	hashmap.Store(p.args[0], entry)
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
		entry := value.(CacheEntry)
		if entry.expiry > -1 && time.Now().UnixMilli() > entry.expiry {
			hashmap.Delete(p.args[0])
			p.writer.WriteNilResponsString()
			return
		}
		p.writer.WriteBulkResponseString(entry.value)
		return
	}
	p.writer.WriteNilResponsString()
}
