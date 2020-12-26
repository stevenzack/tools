package fileToolkit

import (
	"os"
	"path/filepath"
)

func Walk(dir string, check func(path string) bool) ([]string, error) {
	out := []string{}
	e := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if check(path) {
			out = append(out, path)
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}
