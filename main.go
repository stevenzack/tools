package main

import (
	"fmt"
	"log"

	"github.com/StevenZack/tools/fileToolkit"
)

func main() {
	line, i, e := fileToolkit.Tailn1("main.go")
	if e != nil {
		log.Println(e)
		return
	}
	fmt.Println(i)
	fmt.Println(line)
}