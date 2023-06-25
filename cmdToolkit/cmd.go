package cmdToolkit

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func Run(program string, args ...string) (string, error) {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	out := new(bytes.Buffer)
	c.Stdout = out
	err, e := c.StderrPipe()
	if e != nil {
		return "", e
	}
	scanner := bufio.NewScanner(err)

	e = c.Start()
	if e != nil {
		return "", e
	}
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if text != "" {
			return "", errors.New(text)
		}
	}

	return out.String(), nil
}

func RunAttach(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
