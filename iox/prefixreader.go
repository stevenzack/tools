package iox

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

type (
	PrefixReadCloser struct {
		r   io.ReadCloser
		buf *bytes.Buffer
	}
)

func NewPrefixReaderCloser(r io.ReadCloser) *PrefixReadCloser {
	return &PrefixReadCloser{
		r:   r,
		buf: new(bytes.Buffer),
	}
}
func (p *PrefixReadCloser) ReadPrefix(length int) ([]byte, error) {
	b := make([]byte, length)
	n, e := p.r.Read(b)
	if e != nil {
		return nil, e
	}
	if n < length {
		return nil, errors.New("not enough length for prefix[" + strconv.Itoa(length) + "]")
	}
	_, e = p.buf.Write(b[:n])
	if e != nil {
		return nil, e
	}

	return b[:n], nil
}
func (p *PrefixReadCloser) Read(b []byte) (n int, err error) {
	if p.buf.Len() != 0 {
		return p.buf.Read(b)
	}
	return p.r.Read(b)
}
func (p *PrefixReadCloser) Close() error {
	return p.r.Close()
}
