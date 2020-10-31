package main

import (
	"log"

	"github.com/StevenZack/tools/fileToolkit"
	"github.com/gofaith/go-zero/core/logx"
)

func main() {
	line, i, e := fileToolkit.Tailn1("/home/asd/app-test/base/logs/error.log")
	if e != nil {
		log.Println(e)
		return
	}
	logx.Info(i)
	logx.Info(line)
}
