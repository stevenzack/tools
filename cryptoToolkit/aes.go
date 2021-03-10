package cryptoToolkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
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
func EncryptAES(src []byte, key []byte) ([]byte, error) {
	c := Md5FromBytes(key)
	block, err := aes.NewCipher(c)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, c)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// 解密
func DecryptAES(src []byte, key []byte) ([]byte, error) {
	if len(src) < len(key) {
		return nil, errors.New("invalid input src length")
	}
	c := Md5FromBytes(key)
	block, err := aes.NewCipher(c)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, c)
	blockMode.CryptBlocks(src, src)
	return unpadding(src)
}
