package stack

import "testing"

func TestStack(t *testing.T) {
	stack := New()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	t.Log(stack.Pop())
	t.Log(stack.IsEmpty())
	t.Log(stack.Pop())
	t.Log(stack.Size())
	t.Log(stack.Pop())
	t.Log(stack.Pop())
	t.Log(stack.Size())
}
