package netToolkit

import (
	"io"
	"net/http"
	"os"
)

func DownloadWithProgress(url, dst string, onPro func(offset, total int64)) error {
	rp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer rp.Body.Close()
	f, e := os.OpenFile(dst, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	p := make([]byte, 10240)
	var offset int64
	for {
		n, e := rp.Body.Read(p)
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		if n <= 0 {
			break
		}
		_, e = f.Write(p[:n])
		if e != nil {
			return e
		}
		offset += int64(n)
		if onPro != nil {
			onPro(offset, rp.ContentLength)
		}
	}
	return nil
}
