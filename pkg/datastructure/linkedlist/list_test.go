package linkedlist

import "testing"

func TestList(t *testing.T) {
	list := new(List)
	list.Push(1)
	list.Push(2)
	t.Log(list.Pop())
	t.Log(list.Pop())
	t.Log(list)

	list.Unshift(2)
	list.Unshift(1)
	t.Log(list.Shift())
	t.Log(list.Shift())
	t.Log(list)

	list.Push(1)
	list.Push(2)
	list.Push(3)
	t.Log(list.Shift())
	t.Log(list.Pop())
	t.Log(list.Shift())
	t.Log(list)

	t.Log(list.Pop())
	t.Log(list.Shift())
}
