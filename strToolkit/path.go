package strToolkit

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func GetOS() string {
	return runtime.GOOS
}
func HandleTmpDir(pkgDir string) {
	path, _ := filepath.Abs(pkgDir)
	if GetOS() == "android" {
		e := os.MkdirAll(path+"/tmp", 0755)
		if e != nil {
			fmt.Println("mkdirAll() failed:", e)
			return
		} else {
			os.Setenv("TMPDIR", path+"/tmp/")
		}
	}
}

func Getrpath(path string) string {
	if len(path) == 0 {
		return ""
	}
	sep := string(os.PathSeparator)
	if GetLast(path) == sep {
		return path
	}
	return path + sep
}

func Getunpath(path string) string {
	if len(path) == 0 {
		return ""
	}
	sep := string(os.PathSeparator)
	if GetLast(path) != sep {
		return path
	}
	return path[:len(path)-1]
}

func GetUserHomeDir() string {
	c, e := user.Current()
	if e != nil {
		fmt.Println(e)
		d, _ := os.Getwd()
		return d
	}
	return c.HomeDir
}

func GetDirOfFile(path string) string {
	if path == "" {
		return path
	}
	sep := string(os.PathSeparator)
	for i := len(path) - 1; i > -1; i-- {
		if path[i:i+1] == sep {
			return path[:i+1]
		}
	}
	return "." + string(os.PathSeparator)
}
