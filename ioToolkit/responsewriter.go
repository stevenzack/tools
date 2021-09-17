package ioToolkit

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"reflect"
)

type ResponseWriterBuffer struct {
	buffer     *bytes.Buffer
	w          http.ResponseWriter
	statusCode int
}

func NewResponseWriterLogger(w http.ResponseWriter) *ResponseWriterBuffer {
	return &ResponseWriterBuffer{
		buffer: new(bytes.Buffer),
		w:      w,
	}
}

func (r *ResponseWriterBuffer) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseWriterBuffer) Write(b []byte) (int, error) {
	return r.buffer.Write(b)
}

func (r *ResponseWriterBuffer) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func (r *ResponseWriterBuffer) String() string {
	return r.buffer.String()
}

func (r *ResponseWriterBuffer) Bytes() []byte {
	return r.buffer.Bytes()
}

func (r *ResponseWriterBuffer) Flush() {
	r.w.WriteHeader(r.statusCode)
	r.w.Write(r.Bytes())
}

func (r *ResponseWriterBuffer) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	h, ok := r.w.(http.Hijacker)
	if ok {
		return h.Hijack()
	}
	panic(fmt.Sprintf("ResponseWriterBuffer.w:%s is not a http.Hijack", reflect.TypeOf(r.w).String()))
}
