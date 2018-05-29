package tools

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func SplitHans(s string) string {
	rs := []rune(s)
	result := ""
	for k, v := range rs {
		str := string(v)
		if IsChines(str) {
			result += str + " "
			continue
		}
		if !IsEnglish(v) {
			result += str
		}
		left, right := "", ""
		if k != 0 && IsChines(string(rs[k-1])) {
			left = " "
		}
		if k != len(rs)-1 && IsChines(string(rs[k+1])) {
			right = " "
		}
		result += left + str + right
	}
	return result
}
func IsChines(s string) bool {
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	return hzRegexp.MatchString(s)
}
func IsEnglish(r rune) bool {
	if r >= 65 && r <= 90 || r >= 97 && r <= 122 {
		return true
	}
	return false
}
func NewToken() string {
	ct := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
func GetTimeStrNow() string {
	return time.Now().Format("2006-01-01 15:03:02")
}
func NewNumToken() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(int64(r.Intn(60000)), 10)
}
func EncodingGBK(src string) string {
	if runtime.GOOS != "windows" {
		return src
	}
	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	return string(data)
}
func GetOS() string {
	return runtime.GOOS
}
func HandleTmpDir(pkgDir string) {
	path, _ := filepath.Abs(pkgDir)
	if GetOS() == "android" {
		e := os.MkdirAll(path+"/tmp", 0755)
		if e != nil {
			fmt.Println("mkdirAll() failed:", e)
			return
		} else {
			os.Setenv("TMPDIR", path+"/tmp/")
		}
	}
}
