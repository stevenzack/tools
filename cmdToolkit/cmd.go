package cmdToolkit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run(program string, args ...string) (string, error) {
	c := exec.Command(program, args...)
	fo := new(strings.Builder)
	fe := new(strings.Builder)
	c.Stderr = fe
	c.Stdout = fo
	e := c.Run()
	if c.ProcessState.ExitCode() != 0 {
		return fo.String(), fmt.Errorf("run command [%s %s] failed: %w, %s", program, strings.Join(args, " "), e, fo.String()+fe.String())
	}
	return fo.String() + fe.String(), nil
}

func RunAttach(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
