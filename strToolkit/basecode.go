package strToolkit

import (
	"errors"
	"strconv"
	"strings"
)

const (
	totalRange = 62
	charRange  = 26
)

func BaseEncode(i int64) string {
	if i == 0 {
		return "0"
	}
	buf := new(strings.Builder)
	for {
		if i == 0 {
			break
		}
		b := i % totalRange
		i = i / totalRange
		// b to rune
		if b < 0 {
			panic("b is less than 0")
		}
		if b <= 9 {
			buf.WriteString(strconv.FormatInt(b, 10))
			continue
		}
		if b <= 9+charRange {
			buf.WriteRune(rune('a' + (b - 10)))
			continue
		}
		if b <= 9+charRange*2 {
			buf.WriteRune(rune('A' + (b - 10 - 26)))
			continue
		}
		panic("invalid b:" + strconv.FormatInt(b, 10))
	}
	s := []rune(buf.String())
	out := new(strings.Builder)
	for i := len(s) - 1; i > -1; i-- {
		out.WriteRune(s[i])
	}
	return out.String()
}

func BaseDecode(s string) (int64, error) {
	rs := []rune(s)
	var num, sep int64 = 0, 1
	for i := len(rs) - 1; i > -1; i-- {
		r := rs[i]
		var l int64
		if r >= '0' && r <= '9' {
			l = int64(r - '0')
		} else if r >= 'a' && r <= 'z' {
			l = int64(r-'a') + 10
		} else if r >= 'A' && r <= 'Z' {
			l = int64(r-'A') + 36
		} else {
			return 0, errors.New("invalid rune '" + string(r) + "' for BaseCodec")
		}

		num += l * sep
		sep *= totalRange
	}
	return num, nil
}
