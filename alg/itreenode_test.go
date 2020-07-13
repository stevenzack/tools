package alg

import "testing"

func TestPrintTree(t *testing.T) {
	tree := Node("1",
		Node("2",
			Node("3"),
			Node("4"),
		),
		Node("5",
			Node("6"),
			Node("7",
				Node("8")),
		),
	)
	PrintTree(tree)
}
