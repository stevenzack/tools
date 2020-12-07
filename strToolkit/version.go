package strToolkit

import (
	"errors"
	"strconv"
	"strings"
)

func ToVersionCode(version string) (int, error) {
	ss := strings.Split(version, ".")
	out := 0
	for i, s := range ss {
		v, e := strconv.Atoi(s)
		if e != nil {
			return 0, e
		}
		switch i {
		case 0:
			out += v * 10000
		case 1:
			out += v * 100
		case 2:
			out += v
		default:
			return 0, errors.New("version too long:" + version)
		}
	}
	return out, nil
}
