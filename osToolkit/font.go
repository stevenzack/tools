package osToolkit

import "runtime"

var (
	darwinChineseFonts = map[string]struct{}{
		"STHeiti Light.ttc":  struct{}{},
		"STHeiti Medium.ttc": struct{}{},
	}
	windowsChineseFonts = map[string]struct{}{
		"msyh.ttc":    struct{}{},
		"simhei.ttf":  struct{}{},
		"simli.ttf":   struct{}{},
		"simyou.ttf":  struct{}{},
		"simkai.ttf":  struct{}{},
		"simfang.ttf": struct{}{},
	}
)

func GetSystemChineseFont() string {
	switch runtime.GOOS {
	case "darwin":
		return "/System/Library/Fonts/STHeiti Light.ttc"
	default:
		return "C:\\Windows\\Fonts\\msyh.ttc"
	}
}
