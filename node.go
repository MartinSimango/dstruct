package dstruct

// Struct will be represented using a Tree
type Node[T any] struct {
	data     *T
	children map[string]*Node[T]
}

func (n *Node[T]) AddNode(name string, data *T) {
	n.children[name] = &Node[T]{
		data:     data,
		children: make(map[string]*Node[T]),
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
	dataCopy := *n.data
	newNode := &Node[T]{
		data:     new(T),
		children: make(map[string]*Node[T]),
	}
	*newNode.data = dataCopy

	for name := range n.children {
		newNode.children[name] = n.children[name].Copy()
	}
	return newNode
}
