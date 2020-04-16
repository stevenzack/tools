package compressToolkit

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func CompressFileTo(dst io.Writer, path string, progress func(offset, total int64)) error {
	const bufSize = 32 * 1024
	var total int64
	info, e := os.Stat(path)
	if e != nil {
		return e
	}
	if !info.IsDir() {
		total = info.Size()
		file, e := os.OpenFile(path, os.O_RDONLY, 0644)
		if e != nil {
			return e
		}
		defer file.Close()

		zw := zip.NewWriter(dst)
		defer zw.Close()
		header, e := zip.FileInfoHeader(info)
		if e != nil {
			return e
		}
		writer, e := zw.CreateHeader(header)
		if e != nil {
			return e
		}

		buf := make([]byte, bufSize)
		var offset int64
		for {
			n, e := file.Read(buf)
			if e != nil {
				if e == io.EOF {
					break
				}
				return e
			}
			_, e = writer.Write(buf[:n])
			if e != nil {
				return e
			}
			offset += int64(n)
			if progress != nil {
				progress(offset, total)
			}
		}
		return nil
	}

	e = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		total += info.Size()
		return nil
	})
	if e != nil {
		return e
	}

	zw := zip.NewWriter(dst)
	defer zw.Close()
	var offset int64
	base := strToolkit.SubBeforeLast(strToolkit.Getunpath(path), string(os.PathSeparator), strToolkit.Getrpath(path)) + string(os.PathSeparator)

	e = filepath.Walk(path, func(item string, info os.FileInfo, e error) error {
		if info.IsDir() {
			return nil
		}
		header, e := zip.FileInfoHeader(info)
		if e != nil {
			return e
		}
		header.Name = item[len(base):]
		writer, e := zw.CreateHeader(header)
		if e != nil {
			return e
		}
		buf := make([]byte, 32*1024)
		file, e := os.OpenFile(item, os.O_RDONLY, 0644)
		if e != nil {
			return e
		}
		defer file.Close()
		for {
			n, e := file.Read(buf)
			if e != nil {
				if e == io.EOF {
					break
				}
				return e
			}
			_, e = writer.Write(buf[:n])
			if e != nil {
				return e
			}
			offset += int64(n)
			if progress != nil {
				progress(offset, total)
			}
		}
		return nil
	})
	return nil
}

//解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filename := strToolkit.Getrpath(dest) + file.Name
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		if file.FileInfo().IsDir() {
			os.MkdirAll(filename, 0755)
			continue
		}
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
