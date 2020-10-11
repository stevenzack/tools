package ioToolkit

import "io"

type ReadCloser struct {
	Reader io.Reader
	Closer io.Closer
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func (r *ReadCloser) Close() error {
	return r.Closer.Close()
}
