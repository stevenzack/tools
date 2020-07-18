package netToolkit

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	urlkit "net/url"
	"os"
	"strconv"

	"github.com/StevenZack/tools/strToolkit"
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
func PostMultipartFiles(rawUrl string, headers map[string]string, m map[string][]string) (*http.Response, error) {
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	for k, paths := range m {
		for _, path := range paths {
			fi, e := os.OpenFile(path, os.O_RDONLY, 0644)
			if e != nil {
				return nil, e
			}
			defer fi.Close()
			fo, e := w.CreateFormFile(k, strToolkit.SubAfterLast(fi.Name(), string(os.PathSeparator), fi.Name()))
			if e != nil {
				log.Println(e)
				return nil, e
			}
			_, e = io.Copy(fo, fi)
			if e != nil {
				log.Println(e)
				return nil, e
			}
		}
	}
	r, e := http.NewRequest(http.MethodPost, rawUrl, buf)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	r.Header.Set("Content-Type", w.FormDataContentType())
	r.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
	cli := http.Client{}
	return cli.Do(r)
}
func PostMultipartFilesTCP(rawUrl string, headers map[string]string, m map[string][]string) (*http.Response, error) {
	url, e := urlkit.Parse(rawUrl)
	if e != nil {
		return nil, e
	}

	address := url.Host
	if strToolkit.SubAfterLast(url.Host, ":", "") == "" {
		address = address + ":80"
	}
	tcpAddr, e := net.ResolveTCPAddr("tcp", address)
	if e != nil {
		return nil, e
	}
	conn, e := net.DialTCP("tcp", nil, tcpAddr)
	if e != nil {
		return nil, e
	}
	defer conn.Close()
	w := multipart.NewWriter(conn)

	conn.Write([]byte("POST " + url.RequestURI() + " HTTP/1.1\r\n"))
	conn.Write([]byte("Host: " + address + "\r\n"))
	for k, v := range headers {
		conn.Write([]byte(k + ": " + v + "\r\n"))
	}
	var length int
	for k, paths := range m {
		for _, path := range paths {
			info, e := os.Stat(path)
			if e != nil {
				log.Println(e)
				return nil, e
			}
			length += len("--"+w.Boundary()+"\r\nContent-Disposition: form-data; name=\""+k+"\"; filename=\""+info.Name()+"\"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n") + int(info.Size())
		}
	}
	length += len("--" + w.Boundary() + "\r\n")
	conn.Write([]byte("Content-Length: " + strconv.Itoa(length) + "\r\n"))
	conn.Write([]byte("Content-Type: multipart/form-data; boundary=" + w.Boundary() + "\r\n\r\n"))

	for k, paths := range m {
		for _, path := range paths {
			fi, e := os.OpenFile(path, os.O_RDONLY, 0644)
			if e != nil {
				return nil, e
			}
			defer fi.Close()
			info, e := os.Stat(path)
			if e != nil {
				log.Println(e)
				return nil, e
			}
			fo, e := w.CreateFormFile(k, info.Name())
			_, e = io.Copy(fo, fi)
			if e != nil {
				log.Println(e)
				return nil, e
			}
		}
	}
	conn.Write([]byte("\r\n--" + w.Boundary() + "\r\n"))
	return nil, nil
}
