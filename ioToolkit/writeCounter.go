package ioToolkit

import "time"

type WriteCounter struct {
	Total      uint64
	LastSecond int
	OnProgress func(i uint64)
}

func (f *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	f.Total += uint64(n)
	second := time.Now().Second()
	if f.OnProgress != nil && second != f.LastSecond {
		f.OnProgress(f.Total)
		f.LastSecond = second
	}
	return n, nil
}
