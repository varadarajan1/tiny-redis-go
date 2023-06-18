package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

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

func (w Writer) WriteNilResponsString() {
	response := "$-1" + "\r\n"
	w.writer.Write([]byte(response))
	w.writer.Flush()
}

func (w Writer) WriteErrorResponseString(resp string) {
	response := "-" + resp + "\r\n"
	w.writer.Write([]byte(response))
	w.writer.Flush()
}
