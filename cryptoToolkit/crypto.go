package cryptoToolkit

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

func Md5FromBytes(b []byte) []byte {
	h := md5.New()
	io.Copy(h, bytes.NewReader(b))
	return h.Sum(nil)
}

func Md5Str(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
