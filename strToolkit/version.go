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

func ToVersionCodeFloat(version string) (float32, error) {
	ss := strings.Split(version, ".")
	out := 0.0
	for i, s := range ss {
		v, e := strconv.ParseFloat(s, 64)
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
		case 3:
			out += v * 0.01
		case 4:
			out += v * 0.0001
		case 5:
			out += v * 0.000001
		default:
			return 0, errors.New("version too long:" + version)
		}
	}
	return float32(out), nil
}
