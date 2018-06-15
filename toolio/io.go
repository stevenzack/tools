package toolio

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func ReadAll(r io.Reader) string {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return ""
	}
	return string(b)
}
func RunAttachedCmd(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	e := c.Run()
	return e
}
