package fileToolkit

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
)

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
func WriteFile(f string) (*os.File, error) {
	e := os.MkdirAll(GetDirOfFile(f), 0755)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}
	file, e := os.OpenFile(f, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	return file, e
}
func GetDirOfFile(f string) string {
	for i := len(f) - 1; i > -1; i-- {
		if f[i:i+1] == string(os.PathSeparator) {
			return f[:i]
		}
	}
	return f
}
func FileExists(f string) bool {
	_, e := os.Stat(f)
	return e == nil
}
func GetFileExt(f string) string {
	return path.Ext(f)
}
func GetFileMimeType(f string) string {
	t := mime.TypeByExtension(GetFileExt(f))
	return t
}
func GetHomeDir() (string, error) {
	u, e := user.Current()
	return u.HomeDir, e
}
