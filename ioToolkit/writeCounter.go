package ioToolkit

import "time"

type WriteCounter struct {
	Step, Total uint64
	LastSecond  int
	OnProgress  func(i uint64)
}

func (f *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	f.Step += uint64(n)
	second := time.Now().Second()
	if f.OnProgress != nil && second != f.LastSecond {
		f.OnProgress(f.Step)
		f.LastSecond = second
	} else if f.Total > 0 && f.Step >= f.Total {
		f.OnProgress(f.Step)
	}
	return n, nil
}
