package dstruct

type Node[T any] struct {
	data     *T
	parent   *Node[T]
	children map[string]*Node[T]
}

func (n *Node[T]) AddNode(name string, data *T) {
	n.children[name] = &Node[T]{
		data:     data,
		children: make(map[string]*Node[T]),
		parent:   n,
	}
}

func (n *Node[T]) DeleteNode(name string) {
	delete(n.children, name)
}

func (n *Node[T]) GetNode(name string) *Node[T] {
	return n.children[name]
}

func (n *Node[T]) HasChildren() bool {
	return len(n.children) > 0
}

func (n *Node[T]) Copy() *Node[T] {
	newNode := &Node[T]{
		data:     new(T),
		children: make(map[string]*Node[T]),
	}

	dataCopy := *n.data
	*newNode.data = dataCopy

	for name := range n.children {
		newNode.children[name] = n.children[name].Copy()
		newNode.children[name].parent = newNode
	}
	return newNode
}
