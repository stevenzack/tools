package osToolkit

import "runtime"

var windowsChineseFonts = map[string]struct{}{
	"msyh.ttc":    struct{}{},
	"simhei.ttf":  struct{}{},
	"simli.ttf":   struct{}{},
	"simyou.ttf":  struct{}{},
	"simkai.ttf":  struct{}{},
	"simfang.ttf": struct{}{},
}

func GetSystemChineseFont() string {
	switch runtime.GOOS {
	default:
		return "C:\\Windows\\Fonts\\msyh.ttc"
	}
}
