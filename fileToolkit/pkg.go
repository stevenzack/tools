package fileToolkit

import (
	"bufio"
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func IsGoPathPkg(pkgPath string) bool {
	if pkgPath == "" {
		return false
	}
	gopath := GetGOPATH()
	return IsDirExists(gopath + "src/" + pkgPath)
}

func GetPkgNameFromPkg(gopkg string) (string, error) {
	dir := GetGOPATH() + "src/" + gopkg
	gofile, e := GetFirstGoFile(dir)
	if e != nil {
		return "", e
	}
	firstLine, e := ReadFirstLine(gofile)
	if e != nil {
		return "", e
	}
	pkg, e := ReadPkgFromLine(firstLine)
	if e != nil {
		return "", e
	}
	return pkg, nil
}

func GetFirstGoFile(dir string) (string, error) {

	files, e := RangeFilesInDir(dir)
	if e != nil {
		return "", e
	}
	if len(files) == 0 {
		return "", errors.New("no .go files in dir:" + dir)
	}
	return files[0], nil
}

func ReadPkgFromLine(l string) (string, error) {
	l = strings.Replace(l, "\t", "", -1)
	strs := strings.Split(l, " ")
	if len(strs) < 2 {
		return "", errors.New("bad package format")
	}
	return strs[1], nil
}

func ReadFirstLine(filePath string) (string, error) {
	f, e := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if e != nil {
		return "", e
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, e := ReadLine(r)
	return line, e
}

func ReadLine(r *bufio.Reader) (string, error) {
	b, _, e := r.ReadLine()
	if e != nil {
		return "", e
	}
	return string(b), nil
}

func GetGOPATH() string {
	return strToolkit.Getrpath(os.Getenv("GOPATH"))
}
func GetPkgFromDir(dir string) (string, error) {
	dir = strToolkit.Getrpath(dir)
	if !IsDirExists(dir) {
		return "", errors.New("dir " + dir + " does not exists")
	}
	if !strings.Contains(dir, GetGOPATH()) {
		return "", errors.New("dir " + dir + " is not a Go Package")
	}
	pkg := dir[len(GetGOPATH()+"src/"):]
	return strToolkit.Getunpath(pkg), nil
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

func MkdirsOfFilePath(fpath string) string {
	dir := strToolkit.GetDirOfFile(fpath)
	os.MkdirAll(dir, 0755)
	return dir
}

// ParseGoPkg formats path like : ./model , ../me_gengo , /home/asd/go/src/base... into Go Package like github.com/StevenZack/gengo
func ParseGoPkg(curDir, pkg string) string {
	sep := string(os.PathSeparator)
	if strings.Contains(pkg, sep) {
		if strings.HasPrefix(pkg, "."+sep) || strings.HasPrefix(pkg, ".."+sep) {
			return strToolkit.Getrpath(curDir) + pkg
		}
		if strings.HasPrefix(pkg, sep) || runtime.GOOS == "windows" && strings.Contains(pkg, ":"+sep) {
			return pkg
		}
		return GetGOPATH() + pkg
	}
	return strToolkit.Getrpath(curDir) + pkg
}
