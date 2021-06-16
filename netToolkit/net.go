package netToolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/StevenZack/tools/strToolkit"
)

func DoPostMultipartWithHeaders(url string, m map[string]interface{}, headers map[string]string) (string, error) {
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
			st, e := vv.Stat()
			if e != nil {
				continue
			}
			fo, e := w.CreateFormFile(k, st.Name())
			if e != nil {
				continue
			}
			io.Copy(fo, vv)
			vv.Close()
			continue
		}
		if vv, ok := v.([]*os.File); ok {
			for _, vvv := range vv {
				st, e := vvv.Stat()
				if e != nil {
					continue
				}
				fo, e := w.CreateFormFile(k, st.Name())
				if e != nil {
					continue
				}
				io.Copy(fo, vvv)
				vvv.Close()
			}
			continue
		}
	}
	w.Close()
	r, e := http.NewRequest("POST", url, buf)
	if e != nil {
		return "", e
	}
	r.Header.Set("Content-Type", "multipart/form-data; boundary="+w.Boundary())
	if headers != nil {
		for k, v := range headers {
			r.Header.Set(k, v)
		}
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
func DoPostMultipart(url string, m map[string]interface{}) (string, error) {
	return DoPostMultipartWithHeaders(url, m, nil)
}
func DoPostJson(url string, i interface{}) ([]byte, error) {
	b, e := json.Marshal(i)
	if e != nil {
		return nil, e
	}
	r := bytes.NewReader(b)
	client := http.Client{}
	rp, e := client.Post(url, "application/json", r)
	if e != nil {
		if strings.Contains(e.Error(), "refuse") {
			return nil, fmt.Errorf("服务器连接失败")
		}
		if strings.Contains(e.Error(), "network is unreachable") {
			return nil, fmt.Errorf("无网络连接")
		}
		return nil, e
	}
	defer rp.Body.Close()
	s, e := ioutil.ReadAll(rp.Body)
	return s, e
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
	client := http.Client{
		Timeout: time.Second * 5,
	}
	rp, e := client.Get(url)
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
	os.MkdirAll(strToolkit.SubBeforeLast(fdist, string(os.PathSeparator), fdist), 0755)
	f, e := os.OpenFile(fdist, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = io.Copy(f, rp.Body)
	return e
}

// DownloadFileToDir return filename
func DownloadFileToDir(url, dir string) (string, error) {
	rp, e := http.Get(url)
	if e != nil {
		return "", e
	}
	defer rp.Body.Close()
	filename := GetDispFileName(rp)
	f, e := os.OpenFile(strToolkit.Getrpath(dir)+filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return "", e
	}
	defer f.Close()
	_, e = io.Copy(f, rp.Body)
	return filename, e
}
func GetDispFileName(rp *http.Response) string {
	url := rp.Request.URL.EscapedPath()
	str := rp.Header.Get("Content-Disposition")
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

func GetIPs(ipv6 bool) []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	strs := []string{}
	maxAddr := ""
	maxValue := 0
	maxIndex := 0
	ipv6s := []string{}
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
				if strings.HasSuffix(ip.String(), "::1") || ip.String() == "127.0.0.1" {
					continue
				}
				if strings.Contains(ip.String(), ":") {
					if ipv6 {
						ipv6s = append(ipv6s, "["+ip.String()+"]")
					}
					continue
				}
				strs = append(strs, ip.String())
				value, e := strconv.Atoi(strToolkit.SubBefore(ip.String(), ".", "0"))
				if e != nil {
					continue
				}
				if value > maxValue {
					maxValue = value
					maxAddr = ip.String()
					maxIndex = len(strs) - 1
				}
			case *net.IPAddr:
				// ip := v.IP
				// strs = append(strs, ip.String())
			}
		}
	}
	if len(strs) == 0 {
		return nil
	}
	strs = append([]string{maxAddr}, append(strs[:maxIndex], strs[maxIndex+1:]...)...)
	return append(strs, ipv6s...)
}

func DoJSONRequest(url string, i interface{}) (string, error) {
	b, e := json.Marshal(i)
	if e != nil {
		return "", e
	}
	body := bytes.NewReader(b)
	r, e := http.NewRequest("POST", url, body)
	if e != nil {
		return "", e
	}
	c := http.Client{}
	r.Header.Set("Content-Type", "application/json")
	rp, e := c.Do(r)
	if e != nil {
		return "", e
	}
	defer rp.Body.Close()
	back, e := ioutil.ReadAll(rp.Body)
	if e != nil {
		return "", e
	}
	return string(back), nil
}

func DownloadFileUnknownSize(url, dst string, onProgress func(rcv uint64) bool) error {
	rp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer rp.Body.Close()

	if rp.StatusCode != 200 {
		b, e := ioutil.ReadAll(rp.Body)
		if e != nil {
			return e
		}
		return errors.New(rp.Status + ":" + string(b))
	}

	fo, e := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	defer fo.Close()
	b := make([]byte, 10240)
	lastSecond := time.Now().Second()
	var readn uint64
	for {
		n, e := rp.Body.Read(b)
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		_, e = fo.Write(b[:n])
		if e != nil {
			return e
		}
		readn += uint64(n)
		if lastSecond == time.Now().Second() {
			continue
		}
		if onProgress != nil {
			if onProgress(readn) {
				return nil
			}
		}
	}
	if onProgress != nil {
		onProgress(readn)
	}
	return nil
}

func DownloadFileWithProgress(url string, header map[string]string, dst string, onProgress func(rcv, total uint64) bool) error {
	r, e := http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		log.Println(e)
		return e
	}
	if header != nil {
		for k, v := range header {
			r.Header.Add(k, v)
		}
	}
	rp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer rp.Body.Close()

	if rp.Status != "200 OK" {
		b, e := ioutil.ReadAll(rp.Body)
		if e != nil {
			return e
		}
		return errors.New(string(b))
	}

	length, e := strconv.ParseUint(rp.Header.Get("Content-Length"), 10, 64)
	if e != nil {
		return e
	}

	var fo *os.File
	st, e := os.Stat(dst)
	if e == nil && st.IsDir() {
		fo, e = os.OpenFile(strToolkit.Getrpath(dst)+GetDispFileName(rp), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if e != nil {
			return e
		}
	} else {
		fo, e = os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if e != nil {
			return e
		}
	}
	defer fo.Close()

	b := make([]byte, 10240)
	lastSecond := time.Now().Second()
	var readn uint64
	for {
		n, e := rp.Body.Read(b)
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		_, e = fo.Write(b[:n])
		if e != nil {
			return e
		}
		readn += uint64(n)
		if lastSecond == time.Now().Second() {
			continue
		}

		if onProgress != nil {
			if onProgress(readn, length) {
				return nil
			}
		}
	}

	if onProgress != nil {
		if onProgress(length, length) {
			return nil
		}
	}
	return nil
}

func RandomPort() int {
	return rand.Intn(10000) + 10000
}
