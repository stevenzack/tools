package cryptoToolkit

import (
	"testing"
)

func TestEncryptAES(t *testing.T) {
	key := []byte("asd")
	data := "qiwhd"
	enc, e := EncryptAES([]byte(data), key)
	if e != nil {
		t.Error(e)
		return
	}
	dec, e := DecryptAES(enc, key)
	if e != nil {
		t.Error(e)
		return
	}
	str := string(dec)
	if str != data {
		t.Error("str is not data , but ", str)
		return
	}

	enc, e = EncryptAES(nil, nil)
	if e != nil {
		t.Error(e)
		return
	}
	dec, e = DecryptAES(nil, nil)
	if e != nil {
		t.Error(e)
		return
	}
	if dec != nil {
		t.Error("dec is not nil , but ", dec)
		return
	}
}
