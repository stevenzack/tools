package cmdToolkit

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func Run(program string, args ...string) (string, error) {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	out := new(bytes.Buffer)
	c.Stdout = out
	err := new(bytes.Buffer)
	c.Stderr = err
	e := c.Run()
	if e != nil {
		return "", e
	}

	es := err.String()
	if es != "" {
		return "", errors.New(es)
	}
	return out.String(), nil
}
