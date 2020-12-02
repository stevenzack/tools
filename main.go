package main

import (
	"fmt"

	"github.com/StevenZack/tools/cryptoToolkit"
)

func main() {
	s := cryptoToolkit.Sha1FromValues(map[string]interface{}{
		"one": 1,
		"two": 2,
	})
	fmt.Println(s)
}
