package main

import (
	"fmt"
)

// ListNode represents a node in the linked list
type ListNode struct {
	Val  int
	Next *ListNode
}

// LinkedList represents the linked list structure with helper methods
type LinkedList struct {
	Head *ListNode
	Size int
}

// ===== SINGLY LINKED LIST IMPLEMENTATION =====

// NewLinkedList creates a new empty linked list
func NewLinkedList() *LinkedList {
	return &LinkedList{
		Head: nil,
		Size: 0,
	}
}

// Insert adds a new node at the beginning - O(1)
func (ll *LinkedList) Insert(val int) {
	newNode := &ListNode{
		Val:  val,
		Next: ll.Head,
	}
	ll.Head = newNode
	ll.Size++
	fmt.Printf("Inserted %d at head. Size: %d\n", val, ll.Size)
}

// Append adds a new node at the end - O(n)
func (ll *LinkedList) Append(val int) {
	newNode := &ListNode{
		Val:  val,
		Next: nil,
	}

	if ll.Head == nil {
		ll.Head = newNode
		ll.Size++
		fmt.Printf("Appended %d to empty list. Size: %d\n", val, ll.Size)
		return
	}

	current := ll.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	ll.Size++
	fmt.Printf("Appended %d to end. Size: %d\n", val, ll.Size)
}

// Delete removes the first occurrence of val - O(n)
func (ll *LinkedList) Delete(val int) bool {
	if ll.Head == nil {
		return false
	}

	// If head node is to be deleted
	if ll.Head.Val == val {
		ll.Head = ll.Head.Next
		ll.Size--
		fmt.Printf("Deleted %d from head. Size: %d\n", val, ll.Size)
		return true
	}

	current := ll.Head
	for current.Next != nil {
		if current.Next.Val == val {
			current.Next = current.Next.Next
			ll.Size--
			fmt.Printf("Deleted %d from list. Size: %d\n", val, ll.Size)
			return true
		}
		current = current.Next
	}

	fmt.Printf("Value %d not found in list\n", val)
	return false
}

// Display prints all elements in the list - O(n)
func (ll *LinkedList) Display() {
	if ll.Head == nil {
		fmt.Println("List is empty")
		return
	}

	current := ll.Head
	fmt.Print("List: ")
	for current != nil {
		fmt.Print(current.Val)
		if current.Next != nil {
			fmt.Print(" -> ")
		}
		current = current.Next
	}
	fmt.Printf(" (Size: %d)\n", ll.Size)
}

// Reverse reverses the linked list iteratively - O(n)
func (ll *LinkedList) Reverse() {
	var prev *ListNode = nil
	current := ll.Head

	for current != nil {
		nextTemp := current.Next // Store next node
		current.Next = prev      // Reverse the link
		prev = current           // Move prev forward
		current = nextTemp       // Move current forward
	}

	ll.Head = prev
	fmt.Println("List reversed")
}

// FindMiddle finds the middle node using two pointers - O(n)
func (ll *LinkedList) FindMiddle() *ListNode {
	if ll.Head == nil {
		return nil
	}

	slow := ll.Head
	fast := ll.Head

	// Fast moves 2 steps, slow moves 1 step
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	fmt.Printf("Middle node value: %d\n", slow.Val)
	return slow
}

// HasCycle detects if there's a cycle using Floyd's algorithm - O(n)
func (ll *LinkedList) HasCycle() bool {
	if ll.Head == nil || ll.Head.Next == nil {
		return false
	}

	slow := ll.Head
	fast := ll.Head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			fmt.Println("Cycle detected!")
			return true
		}
	}

	fmt.Println("No cycle detected")
	return false
}

func main() {
	fmt.Println("=== Singly Linked List Demo ===")
	ll := NewLinkedList()

	// Test insertions
	ll.Insert(1)
	ll.Append(2)
	ll.Append(3)
	ll.Insert(0)
	ll.Display()

	// Test deletions
	ll.Delete(2)
	ll.Display()

	// Test find middle
	ll.FindMiddle()

	// Test reverse
	ll.Reverse()
	ll.Display()

	// Test cycle detection
	ll.HasCycle()
}
