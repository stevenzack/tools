package cryptoToolkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"
)

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) ([]byte, error) {
	n := len(src)
	if n == 0 {
		return nil, errors.New("bad src to unpadding")
	}
	unPadNum := int(src[n-1])
	if n-unPadNum < 0 {
		return nil, errors.New("bad src to unpadding")
	}
	return src[:n-unPadNum], nil
}

// 加密
func EncryptAES(text []byte, key []byte) ([]byte, error) {
	if len(text) == 0 {
		return nil, nil
	}
	h := md5.New()
	io.Copy(h, bytes.NewReader(key))
	// generate a new aes cipher using our 32 byte long key
	c, e := aes.NewCipher(h.Sum(nil))
	// if there are any errors, handle them
	if e != nil {
		return nil, e
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, e := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if e != nil {
		return nil, e
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, e = io.ReadFull(rand.Reader, nonce); e != nil {
		return nil, e
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, []byte(text), nil), nil
}

// 解密
func DecryptAES(text []byte, key []byte) ([]byte, error) {
	if len(text) == 0 {
		return nil, nil
	}
	h := md5.New()
	io.Copy(h, bytes.NewReader([]byte(key)))
	c, e := aes.NewCipher(h.Sum(nil))
	if e != nil {
		return nil, e
	}
	gcm, e := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if e != nil {
		return nil, e
	}

	if len(text) < gcm.NonceSize() {
		return nil, nil
	}

	t, e := gcm.Open(nil, text[:gcm.NonceSize()], text[gcm.NonceSize():], nil)
	if e != nil {
		return nil, e
	}
	return t, nil
}
