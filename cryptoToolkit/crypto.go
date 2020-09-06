package cryptoToolkit

import (
	"crypto/md5"
	"io"
)

func MD5(s string) []byte {
	h := md5.New()
	io.WriteString(h, s)
	return h.Sum(nil)
}
