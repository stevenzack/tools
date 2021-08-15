package cmdToolkit

import (
	"os"
	"runtime"
)

func CurrentLanguage() string {
	switch runtime.GOOS {
	case "windows":
		out, _ := Run("powershell", "Get-Culture | select -exp Name")
		return out
	default:
		return os.Getenv("LANG")
	}
}
