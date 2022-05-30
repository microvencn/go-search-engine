package avl_struct

type Avl[T avlValType] struct {
	head *Node[T]
}

func getTreeHeight[T avlValType](node *Node[T]) int {
	if node == nil {
		return 0
	} else {
		return node.height
	}
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (avl *Avl[T]) InsertNode(node *Node[T], val T) *Node[T] {
	// 递归到此处若 node 为 nil，证明新的节点是上一次递归节点的孩子
	// 所以返回一个新的节点指针
	if node == nil {
		return NewNode(val)
	}

	if val < node.val {
		node.left = avl.InsertNode(node.left, val)
	} else if val > node.val {
		node.right = avl.InsertNode(node.right, val)
	} else {
		return node
	}

	// 更新以 node 为根的树的高度，为其左右子树的最高高度 + 1（其自身）
	node.height = 1 + max(getTreeHeight(node.left), getTreeHeight(node.right))

	balance := getTreeHeight(node.left) - getTreeHeight(node.right)

	// LL
	if balance > 1 && val < node.left.val {
		return avl.rotateRight(node)
	}

	// RR
	if balance < -1 && val > node.right.val {
		return avl.rotateLeft(node)
	}

	// LR
	if balance > 1 && val > node.left.val {
		return avl.rotateLR(node)
	}

	// RL
	if balance < -1 && val < node.right.val {
		return avl.rotateRL(node)
	}

	return node
}

func (avl *Avl[T]) Insert(val T) {
	// 空树
	if avl.head == nil {
		avl.head = NewNode[T](val)
		return
	}
	avl.InsertNode(avl.head, val)
}

func Init[T avlValType]() Avl[T] {
	avl := Avl[T]{
		head: nil,
	}
	return avl
}

// 左旋 (RR)
// right 是 node.right，left 是 node.left，不必再考虑其三之间的相对位置关系，看成独立节点
// 1. right 进入 node 的位置（更改 parent 以及 parent.left 和 parent.right）
// 2. right.left 成为 node.right，由二叉查找树的特性可知，node 始终小于 right.left，此步永远成立
// 3. node 成为 right.left
func (avl *Avl[T]) rotateLeft(node *Node[T]) *Node[T] {
	right := node.right
	// 1
	right.parent = node.parent
	if node.parent != nil {
		if node.parent.left == node {
			node.parent.left = right
		} else {
			node.parent.right = right
		}
	}
	if avl.head == node {
		avl.head = right
	}
	// 2
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	// 3
	node.parent = right
	right.left = node

	// 更新 height
	node.height = 1 + max(getTreeHeight(node.left), getTreeHeight(node.right))
	right.height = 1 + max(getTreeHeight(right.left), getTreeHeight(right.right))

	return right
}

// 右旋 (LL)
// 1. left 取代 node 的位置
// 2. left 的 right 变成 node 的 left
// 3. left 的 right 为 node
func (avl *Avl[T]) rotateRight(node *Node[T]) *Node[T] {
	left := node.left
	// 1
	left.parent = node.parent
	if node.parent != nil {
		if node == node.parent.left {
			node.parent.left = left
		} else {
			node.parent.right = left
		}
	}
	if avl.head == node {
		avl.head = left
	}
	// 2
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	// 3
	node.parent = left
	left.right = node

	// 更新 height
	node.height = 1 + max(getTreeHeight(node.left), getTreeHeight(node.right))
	left.height = 1 + max(getTreeHeight(left.left), getTreeHeight(left.right))

	// 返回子树新的根结点
	return left
}

func (avl *Avl[T]) rotateLR(node *Node[T]) *Node[T] {
	avl.rotateLeft(node.left)
	return avl.rotateRight(node)
}

func (avl *Avl[T]) rotateRL(node *Node[T]) *Node[T] {
	avl.rotateRight(node.right)
	return avl.rotateLeft(node)
}

func (avl *Avl[T]) Inorder() []T {
	s := make([]T, 0)
	node := avl.head
	avl.inorder(node, &s)
	return s
}

func (avl *Avl[T]) inorder(node *Node[T], s *[]T) {
	if node != nil {
		avl.inorder(node.left, s)
		*s = append(*s, node.val)
		avl.inorder(node.right, s)
	}
}
