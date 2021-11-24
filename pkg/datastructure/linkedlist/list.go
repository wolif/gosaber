package linkedlist

import "fmt"

type List struct {
	header *Node
	tail   *Node
	size   int
}

func (l *List) Size() int {
	return l.size
}

func (l *List) Unshift(value interface{}) {
	node := &Node{value: value}
	if l.size == 0 {
		l.tail = node
	} else {
		node.next = l.header
	}
	l.header = node
	l.size++
}

func (l *List) Shift() (interface{}, error) {
	if l.size == 0 {
		return nil, fmt.Errorf("there is no more node")
	}

	n := l.header
	l.size--
	if l.size == 0 {
		l.header = nil
		l.tail = nil
	} else {
		l.header = n.next
	}
	return n.value, nil
}

func (l *List) Push(value interface{}) {
	node := &Node{value: value}
	if l.size == 0 {
		l.header = node
	} else {
		l.tail.next = node
	}
	l.tail = node
	l.size++
}

func (l *List) Pop() (interface{}, error) {
	if l.size == 0 {
		return nil, fmt.Errorf("there is no nore node")
	}

	n := l.tail
	l.size--
	if l.size == 0 {
		l.header = nil
		l.tail = nil
	} else {
		current := 1
		curNode := l.header
		for {
			if current == l.size {
				break
			}
			curNode = curNode.next
			current++
		}
		curNode.next = nil
		l.tail = curNode
	}
	return n.value, nil
}
