package iox

import (
	"net/http"
	"strings"
)

type ResponseWriterLogger struct {
	builder strings.Builder
	w       http.ResponseWriter
}

func NewResponseWriterLogger(w http.ResponseWriter) *ResponseWriterLogger {
	return &ResponseWriterLogger{
		w: w,
	}
}

func (r *ResponseWriterLogger) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseWriterLogger) Write(b []byte) (int, error) {
	return r.builder.Write(b)
}

func (r *ResponseWriterLogger) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

func (r *ResponseWriterLogger) String() string {
	return r.builder.String()
}

func (r *ResponseWriterLogger) Flush() {
	r.w.Write([]byte(r.builder.String()))
}
