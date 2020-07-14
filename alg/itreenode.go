package alg

import (
	"fmt"
)

type ITreeNode interface {
	Child(i int) ITreeNode
	ChildNum() int
	String() string
}

const (
	DIGRAME_SPLIT      = "┴"
	DIAGRAM_LEFT       = "┘"
	DIAGRAM_RIGHT      = "└"
	DIAGRAM_LINE       = "─"
	DIAGRAM_LEFT_DOWN  = "┌"
	DIAGRAM_RIGHT_DOWN = "┐"
)

func PrintTree(root ITreeNode) {
	q := []ITreeNode{root, nil}
	ss := [][]string{[]string{}}
	for level := 0; len(q) > 0; {
		p := q[0]
		q = q[1:]
		if p == nil {
			continue
		}
		ss[level] = append(ss[level], p.String())

		for i := 0; i < p.ChildNum(); i++ {
			n := p.Child(i)
			if n != nil {
				q = append(q, n)
			}
		}

		if len(q) > 0 && q[0] == nil {
			fmt.Println("")
			q = q[1:]
			q = append(q, nil)
			level++
			ss = append(ss, []string{})
		}
	}

}
