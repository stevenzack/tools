package ioToolkit

type WriteCounter struct {
	Total      uint64
	OnProgress func(i uint64)
}

func (f *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	f.Total += uint64(n)
	if f.OnProgress != nil {
		f.OnProgress(f.Total)
	}
	return n, nil
}
