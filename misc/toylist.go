package main

import "fmt"

// ListNode is
type ListNode struct {
	Data int
	Next *ListNode
}

var head *ListNode

func insertFront(data int) {
	temp := head
	head = &ListNode{Data: data, Next: temp}
}

func insertBack(data int) {
	if head == nil {
		head = &ListNode{Data: data, Next: nil}
	}
	n := head
	for ; n.Next != nil; n = n.Next {
	}
	n.Next = &ListNode{Data: data, Next: nil}
}

func insertBefore(data int, n ListNode, p ListNode) {

}

func printList() {
	for n := head; n != nil; n = n.Next {
		fmt.Print(n.Data, " -> ")
	}
	fmt.Println()
}

func main() {
	insertFront(1)
	insertFront(2)
	printList()
	insertBack(3)
	insertBack(4)
	insertBack(5)
	printList()
}
