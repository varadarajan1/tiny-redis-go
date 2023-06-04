package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

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
