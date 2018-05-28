package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func ReadAll(r io.Reader) string {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return ""
	}
	return string(b)
}
func GetAllFilesFromFolder(path string) []string {
	rpath := path
	if len(path) > 0 && path[len(path)-1:] != string(os.PathSeparator) {
		rpath += string(os.PathSeparator)
	}
	dir, e := ioutil.ReadDir(path)
	if e != nil {
		return nil
	}
	files := []string{}
	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, GetAllFilesFromFolder(rpath+fi.Name())...)
		} else {
			files = append(files, rpath+fi.Name())
		}
	}
	return files
}
func GetCurrentExecPath() string {
	f, e := exec.LookPath(os.Args[0])
	if e != nil {
		fmt.Println("exec.LookPath() failed:", e)
		return ""
	}
	realPath, e := filepath.Abs(f)
	if e != nil {
		fmt.Println("filepath.Abs(f) failed", e)
		return ""
	}
	return realPath
}
