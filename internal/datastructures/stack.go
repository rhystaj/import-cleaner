

package datastructures

type StackNode[E any] struct {
	value        *E
	previousNode *StackNode[E]
}

type Stack[E any] struct {
	topNode *StackNode[E]
	size    int
}

func CreateNewStack[E any]() *Stack[E] {
	return &Stack[E]{
		topNode: nil,
		size:    0,
	}
}

func (s *Stack[E]) Size() int {
	return s.size
}

func (s *Stack[E]) Push(value *E) {
	s.topNode = &StackNode[E]{
		value:        value,
		previousNode: s.topNode,
	}
	s.size++
}

func (s *Stack[E]) Peek() *E {
	if s.topNode == nil {
		return nil
	}

	return s.topNode.value
}

func (s *Stack[E]) Pop() *E {
	if s.topNode == nil {
		return nil
	}

	poppedNode := s.topNode
	s.topNode = s.topNode.previousNode

	s.size--

	return poppedNode.value
}
