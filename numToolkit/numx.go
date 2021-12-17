package numToolkit

import (
	"fmt"
	"strings"
)

func Percent(use, total int64) int {
	if use > total {
		return 100
	}
	return int(use * 100 / total)
}

func PercentOverflow(use, total int64) int {
	return int(use * 100 / total)
}

func FormatMBeans(f float64) string {
	if f < 0 {
		return "0"
	}
	s := fmt.Sprintf("%.2f", f)
	return strings.TrimSuffix(s, ".00")
}

func Unsigned(i int) int {
	if i < 0 {
		return 0
	}
	return i
}
