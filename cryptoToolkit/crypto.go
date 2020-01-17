package cryptoToolkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("encrypt ", string(data), "failed", err)
		return nil
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("encrypt ", string(data), "failed", err)
		return nil
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}
func EncryptStr(data string, passphrase string) []byte {
	return Encrypt([]byte(data), passphrase)
}
func DecryptStr(data string, passphrase string) (string, error) {
	de, e := Decrypt([]byte(data), passphrase)
	return string(de), e
}
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) < 12 {
		return nil, errors.New("data too short")
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
func GenConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

func CalcSecWebSocketAccept(input string) string {
	input += "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	sum := sha1.Sum([]byte(input))
	str := base64.StdEncoding.EncodeToString(sum[:])
	return str
}

func HmacSHA1(key, s string) []byte {
	//hmac ,use sha1
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(s))
	return mac.Sum(nil)
}

func Base64Encode(msg string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	return encoded
}
func Base64Decode(msg string) (string, error) {
	decoded, e := base64.StdEncoding.DecodeString(msg)
	if e != nil {
		return "", e
	}
	return string(decoded), nil
}

func MD5(s string) []byte {
	h := md5.New()
	io.WriteString(h, s)
	return h.Sum(nil)
}
