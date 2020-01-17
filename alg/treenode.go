package alg

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func (t *TreeNode) Assign(n **TreeNode) *TreeNode {
	*n = t
	return t
}

func (t *TreeNode) RangeDfs(fn func(node *TreeNode, depth int, isRight bool)) {
	t.rangeDfsRecursively(fn, 0, false)
}

func (t *TreeNode) rangeDfsRecursively(fn func(*TreeNode, int, bool), depth int, isRight bool) {
	fn(t, depth, isRight)
	if t.Left != nil {
		t.Left.rangeDfsRecursively(fn, depth+1, false)
	}
	if t.Right != nil {
		t.Right.rangeDfsRecursively(fn, depth+1, true)
	}
}

func (t *TreeNode) RangeBfs(fn func(node *TreeNode, depth int)) {
	queue1, queue2 := &Queue{}, &Queue{}
	queue1.Push(t)
	depth := 0
	current, next := queue1, queue2
	for queue1.Length() > 0 || queue2.Length() > 0 {
		node := current.Pop().(*TreeNode)
		fn(node, depth)
		if node.Left != nil {
			next.Push(node.Left)
		}
		if node.Right != nil {
			next.Push(node.Right)
		}
		if current.Length() == 0 {
			depth++
			current, next = next, current
		}
	}
}

func (t *TreeNode) Print() {
	maxDepth := 0
	t.RangeDfs(func(node *TreeNode, depth int, isRight bool) {
		if depth > maxDepth {
			maxDepth = depth
		}
	})
	currentDepth := 0
	t.RangeBfs(func(n *TreeNode, depth int) {
		if depth != currentDepth {
			currentDepth = depth
			fmt.Println("")
		}
		fmt.Print(n.Val, ",")
	})
}

func MakeTreeSequentially(list []int) *TreeNode {
	if len(list) == 0 || list[0] < 0 {
		return nil
	}

	root := &TreeNode{Val: list[0]}
	tree := [][]*TreeNode{[]*TreeNode{root}}

	level := 1
	lastSize := 1
	width := 2
	tree = append(tree, make([]*TreeNode, width))
	for i, item := range list {
		if i == 0 {
			continue
		}
		levelIndex := i - lastSize
		var node *TreeNode
		if item > -1 {
			node = &TreeNode{Val: item}
		}
		tree[level][levelIndex] = node
		//next level
		fmt.Print(item, ",")
		if levelIndex == width-1 {
			lastSize += width
			level++
			width = width * 2
			tree = append(tree, make([]*TreeNode, width))
			fmt.Println("")
		}
	}

	// connect
	for i, line := range tree {
		if i == 0 {
			continue
		}
		for j, node := range line {
			if node == nil {
				continue
			}
			father := tree[i-1][j/2]
			if father != nil {
				if j%2 == 0 {
					father.Left = node
				} else {
					father.Right = node
				}
			}
		}
	}
	return root
}

func Node(v int, children ...*TreeNode) *TreeNode {
	n := &TreeNode{Val: v}
	if len(children) > 0 {
		n.Left = children[0]
	}
	if len(children) > 1 {
		n.Right = children[1]
	}
	return n
}
