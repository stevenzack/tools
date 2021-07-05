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

func ParseUnix(str string) (time.Time, error) {
	i, e := strconv.ParseInt(str, 10, 64)
	if e != nil {
		return ZeroTime, e
	}
	if ValidateUnix(i) {
		return time.Unix(i, 0).In(time.Local), nil
	}
	if len(str) < 10 {
		return time.Unix(i, 0).In(time.Local), nil
	}
	it, _ := strconv.ParseInt(str[:10], 10, 64)
	return time.Unix(it, 0).In(time.Local), nil
}
