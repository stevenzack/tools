package fileToolkit

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/stevenzack/tools/strToolkit"
)

func Tailn1(path string) (string, int64, error) {
	st, e := os.Stat(path)
	if e != nil {
		return "", 0, e
	}
	if st.IsDir() {
		return "", 0, errors.New(path + " is not file")
	}

	file, e := os.OpenFile(path, os.O_RDONLY, 0644)
	if e != nil {
		return "", 0, e
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var index int64
	var lastLine string
	for {
		line, e := reader.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				return lastLine, index, nil
			}
			return "", 0, e
		}
		index++
		lastLine = strToolkit.TrimEnd(line, "\n")
	}
}
