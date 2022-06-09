package avl_struct

type avlValType interface {
	string | int | int64 | float32 | float64
}

type Node[T avlValType] struct {
	Val    T
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]
	height int
	Times  int
}

func NewNode[T avlValType](val T) *Node[T] {
	return &Node[T]{
		Val:    val,
		left:   nil,
		right:  nil,
		height: 1,
		parent: nil,
		Times:  1,
	}
}
