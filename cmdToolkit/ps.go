package cmdToolkit

import (
	"strconv"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

type Proc struct {
	Pid  int
	Name string
}

func Ps() ([]Proc, error) {
	rp, e := Run("ps", "-e")
	if e != nil {
		return nil, e
	}
	ss := strings.Split(rp, "\n")
	procs := []Proc{}
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		s = strToolkit.TrimStarts(s, " ")
		s = strToolkit.TrimEnds(s, " ")
		pid := strToolkit.SubBefore(s, " ", "")
		proc := Proc{}
		proc.Pid, e = strconv.Atoi(pid)
		if e != nil {
			continue
		}
		proc.Name = strToolkit.SubAfterLast(s, " ", "")
		if proc.Name == "" {
			panic("bad proc name:" + s)
		}
		procs = append(procs, proc)
	}
	return procs, nil
}

func PsProc(name string) (int, error) {
	procs, e := Ps()
	if e != nil {
		return -1, e
	}
	for _, proc := range procs {
		if proc.Name == name {
			return proc.Pid, nil
		}
	}
	return -1, nil
}
