// Linked List Structure
package main

import "fmt"

type Node struct {
	data int
	next *Node
}

type LinkedList struct {
	head  *Node
	tail  *Node
	count int
}

func List() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) AddElem(n *Node) {
	if l.head == nil {
		l.head = n
		l.tail = n
		l.count++
	}
	l.tail.next = n
	l.tail = n
	l.count++
}

func (l *LinkedList) AddNode(n int) *Node {
	return &Node{data: n}
}

func (l *LinkedList) PrintList() {
	for l.head != nil {
		fmt.Printf("%v ->", l.head.data)
		l.head = l.head.next
	}
	fmt.Println()
}

func main() {
	list := List()
	node1 := list.AddNode(10)
	node2 := list.AddNode(20)
	node3 := list.AddNode(30)
	list.AddElem(node1)
	list.AddElem(node2)
	list.AddElem(node3)
	list.PrintList()
}
