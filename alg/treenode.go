package alg

type SNode struct {
	Value    string
	Children []*SNode
}

func Node(v string, children ...*SNode) *SNode {
	return &SNode{
		Value:    v,
		Children: children,
	}
}

func (s *SNode) ChildNum() int {
	return len(s.Children)
}

func (s *SNode) Child(i int) ITreeNode {
	return s.Children[i]
}

func (s *SNode) String() string {
	return s.Value
}
