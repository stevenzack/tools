package strToolkit

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/StevenZack/tools/numToolkit"
)

func SplitHans(s string) string {
	rs := []rune(s)
	result := ""
	for k, v := range rs {
		str := string(v)
		if IsChinese(v) {
			result += str + " "
			continue
		}
		if !IsEnglish(v) {
			result += str
		}
		left, right := "", ""
		if k != 0 && IsChinese(rs[k-1]) {
			left = " "
		}
		if k != len(rs)-1 && IsChinese(rs[k+1]) {
			right = " "
		}
		result += left + str + right
	}
	return result
}
func IsChinese(r rune) bool {
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	return hzRegexp.MatchString(string(r))
}
func HasChinese(s string) bool {
	for _, r := range s {
		if IsChinese(r) {
			return true
		}
	}
	return false
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
	p := numToolkit.Randn(40000) + 10000
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
	for i := 0; i <= len(s)-len(sep); i++ {
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
	for i := 0; i <= len(s)-len(sep); i++ {
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

func TrimStarts(s string, trim string) string {
	for {
		if strings.HasPrefix(s, trim) {
			s = s[len(trim):]
			continue
		}
		break
	}
	return s
}

func TrimEnds(s string, trim string) string {
	for {
		if strings.HasSuffix(s, trim) {
			s = s[:len(s)-len(trim)]
			continue
		}
		break
	}
	return s
}

func TrimStart(s string, trim string) string {
	if strings.HasPrefix(s, trim) {
		return s[len(trim):]
	}
	return s
}

func TrimEnd(s string, trim string) string {
	if strings.HasSuffix(s, trim) {
		return s[:len(s)-len(trim)]
	}
	return s
}

func TrimBoth(s string, trims string) string {
	return TrimStart(TrimEnd(s, trims), trims)
}

func SubBetween(s string, start, end rune) (string, error) {
	var buf *bytes.Buffer
	for _, r := range s {
		if r == start && buf == nil {
			buf = bytes.NewBufferString("")
			continue
		}
		if buf != nil {
			if r == end {
				return buf.String(), nil
			}
			buf.WriteRune(r)
		}
	}
	return "", errors.New("no end " + string(end))
}

func Ellipsis(s string, width int) string {
	if len(s) > width {
		return s[:width] + ".."
	}
	return s
}

func EllipsisRune(s string, width int) string {
	rs := []rune(s)
	if len(rs) > width {
		return string(rs[:width]) + ".."
	}
	return s
}

func RangeLines(s string, fn func(line string) bool) error {
	r := bufio.NewReader(strings.NewReader(s))
	for {
		line, e := r.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		line = TrimEnd(line, "\n")
		if fn(line) {
			break
		}
	}
	return nil
}
