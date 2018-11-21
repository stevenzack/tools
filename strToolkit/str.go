package strToolkit

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
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
func MD5from(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}

func GetTimeStrNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func NewNumToken() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(int64(r.Intn(60000)), 10)
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
func FmtDuration(d time.Duration) string {
	str := ""
	y := d / (time.Hour * 24 * 365)
	if y != 0 {
		str += fmt.Sprintf("%d年", y)
		d -= y * time.Hour * 24 * 365
	}
	m := d / (time.Hour * 24 * 30)
	if m != 0 {
		str += fmt.Sprintf("%d个月", m)
		d -= m * time.Hour * 24 * 30
	}
	day := d / (time.Hour * 24)
	if day != 0 {
		str += fmt.Sprintf("%d天", day)
		d -= day * time.Hour * 24
	}
	hour := d / time.Hour
	if hour != 0 {
		str += fmt.Sprintf("%d个小时", hour)
		d -= hour * time.Hour
	}
	min := d / time.Minute
	if min != 0 {
		str += fmt.Sprintf("%d分钟", min)
		d -= min * time.Minute
	}
	if y == 0 && m == 0 && day == 0 && hour == 0 && min == 0 {
		second := d / time.Second
		if second != 0 {
			str += fmt.Sprintf("%d秒", second)
			d -= second * time.Second
		}
	}
	return str
}

func Getrpath(path string) string {
	p, _ := filepath.Abs(path)
	return p + string(os.PathSeparator)
}
func Getunpath(path string) string {
	p, _ := filepath.Abs(path)
	return p
}
func EndsWith(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	if s[len(s)-len(suffix):] == suffix {
		return true
	}
	return false
}
func StartsWith(s, preffix string) bool {
	if len(preffix) > len(s) {
		return false
	}
	if s[:len(s)-len(preffix)] == preffix {
		return true
	}
	return false
}
func GetUserHomeDir() string {
	c, e := user.Current()
	if e != nil {
		fmt.Println(e)
		d, _ := os.Getwd()
		return d
	}
	return c.HomeDir
}
func RandomPort() string {
	p := rand.Intn(40000) + 10000
	return strconv.Itoa(p)
}
func JsonArray(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		return "[]"
	}
	return string(b)
}
func JsonObject(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		return "{}"
	}
	return string(b)
}
func UnJson(str string, v interface{}) {
	json.Unmarshal([]byte(str), v)
}
