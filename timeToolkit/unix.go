package timeToolkit

import (
	"strconv"
	"time"
)

func ValidateUnix(i int64) bool {
	if i > 1000000000 && i < 1900000000 {
		return true
	}
	return false
}

func NowUnixTime() int64 {
	return UnixTimeOf(time.Now())
}
func UnixTimeOf(t time.Time) int64 {
	i := t.Unix()
	if ValidateUnix(i) {
		return i
	}
	s := strconv.FormatInt(i, 64)
	if len(s) < 10 {
		return i
	}
	it, _ := strconv.ParseInt(s[:10], 10, 64)
	return it
}
