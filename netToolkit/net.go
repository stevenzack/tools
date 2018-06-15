package netToolkit

import (
	"bytes"
	"fmt"
	"github.com/StevenZack/tools/fileToolkit"
	"github.com/StevenZack/tools/strToolkit"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
)

func DoPostMultipart(url string, m map[string]interface{}) (string, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	for k, v := range m {
		if vv, ok := v.(string); ok {
			str, e := w.CreateFormField(k)
			if e != nil {
				continue
			}
			str.Write([]byte(vv))
			continue
		}
		if vv, ok := v.(*os.File); ok {
			fo, e := w.CreateFormFile(k, vv.Name())
			if e != nil {
				continue
			}
			io.Copy(fo, vv)
			vv.Close()
			continue
		}
	}
	w.Close()
	r, e := http.NewRequest("POST", url, buf)
	if e != nil {
		return "", e
	}
	var client http.Client
	rp, e := client.Do(r)
	if e != nil {
		return "", e
	}
	defer rp.Body.Close()
	b, e := ioutil.ReadAll(rp.Body)
	return string(b), e
}
func GetIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var strs []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				strs = append(strs, ip.String())
			case *net.IPAddr:
				// ip := v.IP
				// strs = append(strs, ip.String())
			}
		}
	}
	for _, v := range strs {
		if len(v) > 8 && v[:8] == "192.168." {
			return v
		}
	}
	for _, v := range strs {
		if len(v) > 3 && v[:3] == "10.42." {
			return v
		}
	}
	for _, v := range strs {
		if len(v) > 3 && v[:3] == "10." {
			return v
		}
	}
	for _, v := range strs {
		if len(v) > 4 && v[:4] == "172." {
			return v
		}
	}
	for _, v := range strs {
		if v != "127.0.0.1" && v != "::1" && len(v) < 16 {
			return v
		}
	}
	return "127.0.0.1"
}
func GetMacAddr() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var strs, macs []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				strs = append(strs, ip.String())
				macs = append(macs, i.HardwareAddr.String())
			}
		}
	}
	for k, v := range strs {
		if len(v) > 8 && v[:8] == "192.168." {
			return macs[k]
		}
	}
	for k, v := range strs {
		if len(v) > 3 && v[:3] == "10." {
			return macs[k]
		}
	}
	for k, v := range strs {
		if len(v) > 4 && v[:4] == "172." {
			return macs[k]
		}
	}
	for k, v := range strs {
		if v != "127.0.0.1" && v != "::1" && len(v) < 19 {
			return macs[k]
		}
	}
	return ""
}

func IsMyIP(str string) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.String() == str {
					return true
				}
			}
		}
	}
	return false
}
func DoGet(url string) (string, error) {
	rp, e := http.Get(url)
	if e != nil {
		return "", e
	}
	defer rp.Body.Close()
	b, e := ioutil.ReadAll(rp.Body)
	if e != nil {
		return "", e
	}
	return string(b), nil
}
func DownloadFile(url, fdist string) error {
	rp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer rp.Body.Close()
	f, e := fileToolkit.WriteFile(fdist)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = io.Copy(f, rp.Body)
	return e
}
func DownloadFileToDir(url, dir string) (string, error) {
	rp, e := http.Get(url)
	if e != nil {
		return "", e
	}
	defer rp.Body.Close()
	filename := getDispFileName(rp.Request.URL.EscapedPath(), rp.Header.Get("Content-Disposition"))
	f, e := fileToolkit.WriteFile(strToolkit.Getrpath(dir) + filename)
	if e != nil {
		return "", e
	}
	defer f.Close()
	_, e = io.Copy(f, rp.Body)
	return filename, e
}
func getDispFileName(url, str string) string {
	strs := strings.Split(str, ";")
	if len(strs) < 2 {
		return GetFileNameFromEscURL(url)
	}
	for _, v := range strs {
		mindex := strings.Index(v, "filename=")
		if mindex > -1 {
			return strings.Trim(v[mindex+len("filename="):], " ")
		}
	}
	return GetFileNameFromEscURL(url)
}
func GetFileNameFromEscURL(url string) string {
	for i := len(url) - 1; i > -1; i-- {
		if url[i:i+1] == "/" {
			return url[i+1:]
		}
	}
	return url
}