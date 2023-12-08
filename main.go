package main

import (
	"fmt"

	"github.com/stevenzack/tools/cmdToolkit"
)

func main() {
	fmt.Println(cmdToolkit.PsProc("tcp"))
}
