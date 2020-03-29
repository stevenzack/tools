package strToolkit

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func NewToken() string {
	ct := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(ct, 10))
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

func FormatDuration(d time.Duration) string {
	str := ""
	if d < 0 {
		str += "-"
	}
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
