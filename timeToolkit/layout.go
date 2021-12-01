package timeToolkit

import (
	"errors"
	"time"
)

const (
	LAYOUT_DATE        = "2006-01-02"
	LAYOUT_LUNAR       = "2006年01月02日"
	LAYOUT_DATE_TIME   = "2006-01-02 15:04:05"
	LAYOUT_DATE_TIME_Z = "2006-01-02T15:04:05.000Z"
	LAYOUT_MONTH       = "2006-01"
)

func Parse(s string) (time.Time, error) {
	if s == "" {
		return ZeroTime, nil
	}
	if len(s) == len(LAYOUT_DATE_TIME_Z) {
		return time.ParseInLocation(LAYOUT_DATE_TIME_Z, s, time.Local)
	}
	if len(s) == len(LAYOUT_DATE_TIME) {
		return time.ParseInLocation(LAYOUT_DATE_TIME, s, time.Local)
	}
	if len(s) == len(LAYOUT_LUNAR) {
		return time.ParseInLocation(LAYOUT_LUNAR, s, time.Local)
	}
	if len(s) == len(LAYOUT_DATE) {
		return time.ParseInLocation(LAYOUT_DATE, s, time.Local)
	}
	if len(s) == len(LAYOUT_MONTH) {
		return time.ParseInLocation(LAYOUT_MONTH, s, time.Local)
	}
	return ZeroTime, errors.New("时间格式不正确:" + s)
}

func FormatDateTime(t time.Time) string {
	if t == ZeroTime {
		return ""
	}
	return t.Local().Format(LAYOUT_DATE_TIME)
}
func FormatDate(t time.Time) string {
	if t == ZeroTime {
		return ""
	}
	return t.Local().Format(LAYOUT_DATE)
}

func FormatMonth(t time.Time) string {
	if t == ZeroTime {
		return ""
	}
	return t.Local().Format(LAYOUT_MONTH)
}
