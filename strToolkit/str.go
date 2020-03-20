package strToolkit

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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

func IsDigital(r rune) bool {
	if r < 48 || r > 57 {
		return false
	}
	return true
}
func MD5from(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}

func GetLast(s string) string {
	if len(s) == 0 {
		return ""
	}
	index := len(s) - 1
	return s[index:]
}

func RandomPort() string {
	p := rand.Intn(40000) + 10000
	return strconv.Itoa(p)
}
func JsonArray(i interface{}) string {
	if i == nil {
		return "[]"
	}
	b, e := json.Marshal(i)
	if e != nil {
		return "[]"
	}
	return string(b)
}
func JsonObject(i interface{}) string {
	if i == nil {
		return "{}"
	}
	b, e := json.Marshal(i)
	if e != nil {
		return "{}"
	}
	return string(b)
}
func UnJson(str string, v interface{}) {
	json.Unmarshal([]byte(str), v)
}

func CompareVersion(s1, s2 string) (int, error) {
	is1, e := versionToIntegers(s1)
	if e != nil {
		return 0, e
	}
	is2, e := versionToIntegers(s2)
	if e != nil {
		return 0, e
	}
	for i := 0; i < len(is1) && i < len(is2); i++ {
		if is1[i] > is2[i] {
			return 1, nil
		}
		if is1[i] < is2[i] {
			return -1, nil
		}
	}
	return 0, nil
}
func versionToIntegers(s string) ([]int, error) {
	ss := strings.Split(s, ".")
	var is []int
	for _, v := range ss {
		i, e := strconv.ParseUint(v, 10, 64)
		if e != nil {
			return nil, e
		}
		is = append(is, int(i))
	}
	return is, nil
}

func SubBefore(s string, sep, def string) string {
	for i := 0; i < len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return s[:i]
		}
	}
	return def
}
func SubBeforeLast(s, sep, def string) string {
	for i := len(s) - len(sep); i > -1; i-- {
		if s[i:i+len(sep)] == sep {
			return s[:i]
		}
	}
	return def
}

func SubAfter(s, sep, def string) string {
	for i := 0; i < len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return s[i+len(sep):]
		}
	}
	return def
}

func SubAfterLast(s, sep, def string) string {
	for i := len(s) - len(sep); i > -1; i-- {
		if s[i:i+len(sep)] == sep {
			return s[i+len(sep):]
		}
	}
	return def
}

func TrimStart(s, trim string) string {
	if strings.HasPrefix(s, trim) {
		return s[len(trim):]
	}
	return s
}

func TrimEnd(s, trim string) string {
	if strings.HasSuffix(s, trim) {
		return s[:len(s)-len(trim)]
	}
	return s
}

func TrimBoth(s, trim string) string {
	return TrimStart(TrimEnd(s, trim), trim)
}
