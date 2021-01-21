package main

import (
	"fmt"
	"log"

	"github.com/StevenZack/tools/strToolkit"
)

func main() {
	e := strToolkit.RangeLines(`one
	two
	three
	`, func(line string) bool {
		fmt.Println(line)
		return false
	})
	if e != nil {
		log.Println(e)
		return
	}

}
