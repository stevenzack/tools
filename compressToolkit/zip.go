package compressToolkit

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func CompressFilesTo(dst io.Writer, paths []string, progress func(offset, total int64)) error {
	const bufSize = 32 << 10
	var total, offset int64
	// total
	for _, path := range paths {
		info, e := os.Stat(path)
		if e != nil {
			return e
		}
		if info.IsDir() {
			e = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
				if !info.IsDir() {
					info, e := d.Info()
					if e != nil {
						return e
					}
					total += info.Size()
				}
				return nil
			})
			if e != nil {
				log.Println(e)
				return e
			}
		} else {
			total += info.Size()
		}
	}

	if progress != nil {
		progress(0, total)
	}
	// write
	zw := zip.NewWriter(dst)
	defer zw.Close()
	for _, root := range paths {
		rootInfo, e := os.Stat(root)
		if e != nil {
			return e
		}
		root, e = filepath.Abs(root)
		if e != nil {
			return e
		}

		base := filepath.Dir(root)
		if rootInfo.IsDir() {
			e = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				info, e := d.Info()
				if e != nil {
					return e
				}
				//dir
				if info.IsDir() {
					return nil
				}

				path, e = filepath.Abs(path)
				if e!=nil {
					return e
				}
				name, e := filepath.Rel(base, path)
				if e != nil {
					log.Println(e)
					return e
				}
				name = filepath.ToSlash(name)

				header, e := zip.FileInfoHeader(info)
				if e != nil {
					log.Println(e)
					return e
				}
				header.Name = name

				writer, e := zw.CreateHeader(header)
				if e != nil {
					log.Println(e)
					return e
				}

				//symlink
				if info.Mode()&os.ModeType == os.ModeSymlink {
					link, e := os.Readlink(path)
					if e != nil {
						log.Println(e)
						return e
					}
					_, e = writer.Write([]byte(filepath.ToSlash(link)))
					if e != nil {
						log.Println(e)
						return e
					}

					return nil
				}

				//file
				fi, e := os.OpenFile(path, os.O_RDONLY, 0644)
				if e != nil {
					log.Println(e)
					return e
				}
				defer fi.Close()
				buf := make([]byte, bufSize)
				for {
					n, e := fi.Read(buf)
					if e != nil {
						if e == io.EOF {
							break
						}
						log.Println(e)
						return e
					}
					_, e = writer.Write(buf[:n])
					if e != nil {
						log.Println(e)
						return e
					}
					offset += int64(n)
					if progress != nil {
						progress(offset, total)
					}
				}
				return nil
			})
			if e != nil {
				log.Println(e)
				return e
			}
			continue
		}
		// is file
		header, e := zip.FileInfoHeader(rootInfo)
		if e != nil {
			log.Println(e)
			return e
		}
		writer, e := zw.CreateHeader(header)
		if e != nil {
			log.Println(e)
			return e
		}
		fi, e := os.OpenFile(root, os.O_RDONLY, 0644)
		if e != nil {
			log.Println(e)
			return e
		}
		buf := make([]byte, bufSize)
		for {
			n, e := fi.Read(buf)
			if e != nil {
				if e == io.EOF {
					break
				}
				log.Println(e)
				fi.Close()
				return e
			}
			_, e = writer.Write(buf[:n])
			if e != nil {
				log.Println(e)
				fi.Close()
				return e
			}
			offset += int64(n)
			if progress != nil {
				progress(offset, total)
			}
		}
		fi.Close()
	}
	progress(total, total)
	return nil
}

func CompressFileTo(dst io.Writer, rootPath string, progress func(offset, total int64)) error {
	var e error
	rootPath, e = filepath.Abs(rootPath)
	if e != nil {
		return e
	}
	const bufSize = 32 * 1024
	var total int64
	rootInfo, e := os.Stat(rootPath)
	if e != nil {
		return e
	}
	if !rootInfo.IsDir() {
		total = rootInfo.Size()
		file, e := os.OpenFile(rootPath, os.O_RDONLY, 0644)
		if e != nil {
			return e
		}
		defer file.Close()

		zw := zip.NewWriter(dst)
		defer zw.Close()
		header, e := zip.FileInfoHeader(rootInfo)
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

	infos := []fs.FileInfo{}
	paths := []string{}
	e = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		info, e := d.Info()
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}
		path, e = filepath.Abs(path)
		if e != nil {
			return e
		}
		if path == rootPath {
			return nil
		}
		infos = append(infos, info)
		paths = append(paths, path)
		total += info.Size()
		return nil
	})
	if e != nil {
		return e
	}

	if progress != nil {
		progress(0, total)
	}

	zw := zip.NewWriter(dst)
	defer zw.Close()
	var offset int64
	for i, path := range paths {
		info := infos[i]
		name, e := filepath.Rel(rootPath, path)
		if e != nil {
			return e
		}
		name = filepath.ToSlash(name)
		//dir
		if info.IsDir() {
			continue
		}

		header, e := zip.FileInfoHeader(info)
		if e != nil {
			log.Println(e)
			return e
		}
		header.Name = name
		writer, e := zw.CreateHeader(header)
		if e != nil {
			log.Println(e)
			return e
		}

		//symlink
		if info.Mode()&os.ModeType == os.ModeSymlink {
			link, e := os.Readlink(path)
			if e != nil {
				log.Println(e)
				return e
			}
			_, e = writer.Write([]byte(filepath.ToSlash(link)))
			if e != nil {
				log.Println(e)
				return e
			}

			continue
		}

		// file
		buf := make([]byte, 32<<10)
		fmt.Println(path)
		fi, e := os.Open(path)
		if e != nil {
			log.Println(e)
			return e
		}
		for {
			n, e := fi.Read(buf)
			if e != nil {
				if e == io.EOF {
					break
				}
				fi.Close()
				log.Println(e)
				return e
			}
			_, e = writer.Write(buf[:n])
			if e != nil {
				fi.Close()
				log.Println(e)
				return e
			}
			offset += int64(n)
			if progress != nil {
				progress(offset, total)
			}
		}
		fi.Close()
	}

	return nil
}

//解压
func Decompress(zipFile, dest string) error {
	reader, e := zip.OpenReader(zipFile)
	if e != nil {
		log.Println(e)
		return e
	}
	defer reader.Close()
	for _, file := range reader.File {
		filename := strToolkit.Getrpath(dest) + file.Name
		rc, e := file.Open()
		if e != nil {
			log.Println(e)
			return e
		}
		defer rc.Close()
		if file.FileInfo().IsDir() {
			os.MkdirAll(filename, 0755)
			continue
		}
		e = os.MkdirAll(getDir(filename), 0755)
		if e != nil {
			log.Println(e)
			return e
		}
		w, e := os.Create(filename)
		if e != nil {
			log.Println(e)
			return e
		}
		defer w.Close()
		_, e = io.Copy(w, rc)
		if e != nil {
			log.Println(e)
			return e
		}
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
