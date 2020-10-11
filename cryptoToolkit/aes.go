package cryptoToolkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
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
	c := Md5FromBytes(key)
	block, err := aes.NewCipher(c)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, c)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}
