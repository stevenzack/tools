package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/StevenZack/tools/fileToolkit"
)

func main() {
	out, e := fileToolkit.Walk(".", func(path string) bool {
		return strings.HasSuffix(path, ".go")
	})
	if e != nil {
		log.Println(e)
		return
	}
	fmt.Println(strings.Join(out, "\n"))
}
