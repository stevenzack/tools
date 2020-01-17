# Algorithm Toolkit

## TreeNode

![btree](res/binarytree.png)

DSL方式初始化
```go
package main

import (
	"github.com/StevenZack/tools/alg"
)

func main() {
	n := alg.Node(3,
		alg.Node(5,
			alg.Node(6),
			alg.Node(2,
				alg.Node(7),
				alg.Node(4),
			),
		),
		alg.Node(1,
			alg.Node(0),
			alg.Node(8),
		),
	)
	n.Print()
}

```

Slice方式初始化
```go
package main

import (
	"github.com/StevenZack/tools/alg"
)

func main() {
	n := alg.MakeTreeBySlice([]int{3,5,1,6,2,0,8,-1,-1,7,4})
	n.Print()
}

```

遍历
```go
package main

import (
	"fmt"

	"github.com/StevenZack/tools/alg"
)

func main() {
	n := alg.MakeTreeBySlice([]int{3, 5, 1, 6, 2, 0, 8, -1, -1, 7, 4})
	n.RangeBfs(func(n *alg.TreeNode, depth int) {//广度优先
		fmt.Println(n.Val, depth)
	})
	n.RangeDfs(func(n *alg.TreeNode, depth int, isRight bool) {//深度优先
		fmt.Println(n.Val, depth, isRight)
	})
}

```