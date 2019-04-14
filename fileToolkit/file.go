package fileToolkit

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func GetExt(p string) string {
	return path.Ext(p)
}
func GetMimeType(f string) string {
	info, e := os.Stat(f)
	if e != nil {
		fmt.Println(e)
		return e.Error()
	}
	if info.IsDir() {
		return "dir"
	}
	switch strings.ToLower(GetExt(info.Name())) {
	case ".mp4", ".webm", ".mkv", ".3gp", ".flv", ".avi", ".mov", ".rmvb", ".wmv", ".m4v":
		return "video"
	case ".mp3", ".wav", ".amr", ".aac", ".wma", ".midi", ".flac":
		return "audio"
	case ".jpg", ".jpeg", ".webp", ".png", ".gif", ".apng", ".bmp", ".tif", ".svg", ".cdr":
		return "image"
	case ".txt", ".md", ".html", ".css", ".js", ".go", ".java", ".py", ".sh":
		return "text"
	default:
		return "file"
	}
}

// recursively
func GetAllFilesFromFolder(path string) []string {
	path = strToolkit.Getrpath(path)
	dir, e := ioutil.ReadDir(path)
	if e != nil {
		return nil
	}
	files := []string{}
	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, GetAllFilesFromFolder(path+fi.Name())...)
		} else {
			files = append(files, path+fi.Name())
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
	info, e := os.Stat(f)
	if e == nil {
		if info.IsDir() {
			return f
		}
	}
	f = strToolkit.Getunpath(f)
	for i := len(f) - 1; i > -1; i-- {
		if f[i:i+1] == string(os.PathSeparator) {
			return f[:i]
		}
	}
	return f
}
func IsFileExists(f string) bool {
	info, e := os.Stat(f)
	if e != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}
func GetFileExt(f string) string {
	return path.Ext(f)
}
func GetFileMimeType(f string) string {
	t := mime.TypeByExtension(GetFileExt(f))
	return t
}
func GetHomeDir() string {
	u, e := user.Current()
	if e != nil {
		return GetCurrentExecPath()
	}
	return u.HomeDir
}

func GetIconURLByFileType(fpath string) string {
	f, e := os.Open(fpath)
	if e != nil {
		fmt.Println(`os.open error :`, e)
		return ""
	}
	server := "https://jywjl.github.io/images/icons/"
	info, e := f.Stat()
	if e != nil {
		return server + "file.png"
	}
	if info.IsDir() {
		return server + "folder.png"
	}
	nameS := strings.Split(f.Name(), ".")
	ext := nameS[len(nameS)-1]
	mimeTypes := strings.Split(mime.TypeByExtension("."+ext), "/")
	switch mimeTypes[0] {
	case "audio":
		return server + "audio.png"
	case "image":
		return "file://" + fpath
	case "video":
		return server + "video.png"
	default:
		return server + "file.png"
	}
}
func FormatFileSize(size int64) string {
	gb := size / 1024 / 1024 / 1024
	if gb != 0 {
		return fmt.Sprint(gb) + "G"
	}
	mb := size / 1024 / 1024
	if mb != 0 {
		return fmt.Sprint(mb) + "M"
	}
	kb := size / 1024
	if kb != 0 {
		return fmt.Sprint(kb) + "K"
	}
	return fmt.Sprint(size) + "B"
}
func IsFile(path string) (bool, error) {
	info, e := os.Stat(path)
	if e != nil {
		return false, e
	}
	return !info.IsDir(), nil
}
func IsDir(dir string) (bool, error) {
	info, e := os.Stat(dir)
	if e != nil {
		return false, e
	}
	return info.IsDir(), nil
}

func GetGOPATH() string {
	return strToolkit.Getrpath(os.Getenv("GOPATH"))
}
func GetCurrentPkgPath() (string, error) {
	wd, e := os.Getwd()
	if e != nil {
		return "", e
	}
	srcPath := GetGOPATH() + "src/"
	wd = strToolkit.Getrpath(wd)
	if !strings.Contains(wd, srcPath) {
		return "", errors.New("not a Go package")
	}
	pkgPath := wd[len(srcPath):]
	return pkgPath, nil
}

func IsDirExists(dir string) bool {
	info, e := os.Stat(dir)
	if e != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}
func IsPathExists(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

func CopyFile(dst, src string) error {
	fi, e := os.OpenFile(src, os.O_RDONLY, 0644)
	if e != nil {
		return e
	}
	defer fi.Close()
	fo, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}
	defer fo.Close()
	_, e = io.Copy(fo, fi)
	return e
}

func ReadFileAll(path string) (string, error) {
	f, e := os.OpenFile(path, os.O_RDONLY, 0644)
	if e != nil {
		return "", e
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return "", e
	}
	return string(b), nil
}
