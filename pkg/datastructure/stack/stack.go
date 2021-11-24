package stack

import (
	"fmt"
	"github.com/wolif/gosaber/pkg/datastructure/linkedlist"
)

type stack struct {
	list *linkedlist.List
}

func New() *stack {
	return &stack{list: new(linkedlist.List)}
}

func (s *stack) Push(value interface{}) {
	s.list.Push(value)
}

func (s *stack) Pop() (interface{}, error) {
	value, err := s.list.Pop()
	if err != nil {
		return nil, fmt.Errorf("there is no more value")
	}
	return value, nil
}

func (s *stack) IsEmpty() bool {
	return s.list.Size() == 0
}

func (s *stack) Size() int {
	return s.list.Size()
}
