package alg

import "testing"

func TestPrintTree(t *testing.T) {
	PrintTree(NewSNode("1",
		NewSNode("2",
			NewSNode("4",
				NewSNode("5"),
				NewSNode("6"),
			),
		),
		NewSNode("3",
			NewSNode("7",
				NewSNode("8"),
			),
		),
	))
}
