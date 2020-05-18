package sysToolkit

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

// GetRAMSize returns total RAM size in kB
func GetRAMSize() (int64, error) {
	cmd := exec.Command(`system_profiler`, `SPHardwareDataType`)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	e := cmd.Run()
	if e != nil {
		log.Println(e)
		return 0, e
	}
	reader := bufio.NewReader(buf)
	for {
		line, e := reader.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			log.Println(e)
			return 0, e
		}
		if strings.Contains(line, "Memory: ") {
			line = strToolkit.SubAfter(line, "Memory: ", line)
			line = strToolkit.SubBefore(line, "B", line)
			size, e := strconv.ParseInt(strToolkit.SubBefore(line, " ", line), 10, 64)
			if e != nil {
				log.Println(e)
				return 0, e
			}
			b := strToolkit.SubAfterLast(line, " ", line)
			switch b {
			case "G":
				return size << 20, nil
			case "M":
				return size << 10, nil
			case "K":
				return size, nil
			default:
				return 0, errors.New("bad bit unit:" + b)
			}
		}
	}
	return 0, errors.New("memory size not found")
}
