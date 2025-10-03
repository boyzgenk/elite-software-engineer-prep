package main

import (
	"fmt"
	"strings"
)

// ============= STACK IMPLEMENTATION =============

// Stack represents a Last-In-First-Out (LIFO) data structure
type Stack struct {
	items []int
}

// Push adds an element to the top of the stack
func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element from the stack
func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false // Stack is empty
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

// Peek returns the top element without removing it
func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return len(s.items)
}

// ============= QUEUE IMPLEMENTATION =============

// Queue represents a First-In-First-Out (FIFO) data structure
type Queue struct {
	items []int
}

// Enqueue adds an element to the rear of the queue
func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front element from the queue
func (q *Queue) Dequeue() (int, bool) {
	if len(q.items) == 0 {
		return 0, false // Queue is empty
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Front returns the front element without removing it
func (q *Queue) Front() (int, bool) {
	if len(q.items) == 0 {
		return 0, false
	}
	return q.items[0], true
}

// IsEmpty checks if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of elements in the queue
func (q *Queue) Size() int {
	return len(q.items)
}

// ============= STRING STACK FOR PROBLEMS =============

// StringStack for string operations (like parentheses matching)
type StringStack struct {
	items []string
}

func (s *StringStack) Push(item string) {
	s.items = append(s.items, item)
}

func (s *StringStack) Pop() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

func (s *StringStack) Peek() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	return s.items[len(s.items)-1], true
}

func (s *StringStack) IsEmpty() bool {
	return len(s.items) == 0
}

// ============= PRACTICAL APPLICATIONS =============

// ValidParentheses checks if parentheses are balanced
// LeetCode #20
func ValidParentheses(s string) bool {
	stack := &StringStack{}
	mapping := map[string]string{
		")": "(",
		"}": "{",
		"]": "[",
	}

	for _, char := range s {
		charStr := string(char)

		// If it's a closing bracket
		if opener, exists := mapping[charStr]; exists {
			// Check if stack is empty or top doesn't match
			if stack.IsEmpty() {
				return false
			}

			top, _ := stack.Pop()
			if top != opener {
				return false
			}
		} else {
			// It's an opening bracket, push to stack
			stack.Push(charStr)
		}
	}

	// Stack should be empty for valid parentheses
	return stack.IsEmpty()
}

// EvaluateRPN evaluates Reverse Polish Notation
// LeetCode #150
func EvaluateRPN(tokens []string) int {
	stack := &Stack{}

	for _, token := range tokens {
		if token == "+" || token == "-" || token == "*" || token == "/" {
			// Pop two operands
			b, _ := stack.Pop()
			a, _ := stack.Pop()

			var result int
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				result = a / b
			}

			stack.Push(result)
		} else {
			// It's a number, convert and push
			var num int
			fmt.Sscanf(token, "%d", &num)
			stack.Push(num)
		}
	}

	result, _ := stack.Pop()
	return result
}

// DailyTemperatures finds next warmer temperature
// LeetCode #739
func DailyTemperatures(temperatures []int) []int {
	result := make([]int, len(temperatures))
	stack := &Stack{} // Stack to store indices

	for i, temp := range temperatures {
		// While stack is not empty and current temp > temp at stack top
		for !stack.IsEmpty() {
			topIdx, _ := stack.Peek()
			if temp > temperatures[topIdx] {
				idx, _ := stack.Pop()
				result[idx] = i - idx
			} else {
				break
			}
		}
		stack.Push(i)
	}

	return result
}

// ============= QUEUE APPLICATIONS =============

// MovingAverage calculates moving average of numbers
type MovingAverage struct {
	queue *Queue
	size  int
	sum   int
}

// Constructor initializes MovingAverage
func NewMovingAverage(size int) *MovingAverage {
	return &MovingAverage{
		queue: &Queue{},
		size:  size,
		sum:   0,
	}
}

// Next calculates next moving average
func (ma *MovingAverage) Next(val int) float64 {
	ma.queue.Enqueue(val)
	ma.sum += val

	// If queue size exceeds window size, remove oldest
	if ma.queue.Size() > ma.size {
		oldest, _ := ma.queue.Dequeue()
		ma.sum -= oldest
	}

	return float64(ma.sum) / float64(ma.queue.Size())
}

// ============= MONOTONIC STACK PROBLEMS =============

// NextGreaterElement finds next greater element for each element
// LeetCode #496
func NextGreaterElement(nums1 []int, nums2 []int) []int {
	// Build next greater map for nums2
	nextGreater := make(map[int]int)
	stack := &Stack{}

	for _, num := range nums2 {
		// While stack not empty and current > top
		for !stack.IsEmpty() {
			top, _ := stack.Peek()
			if num > top {
				popped, _ := stack.Pop()
				nextGreater[popped] = num
			} else {
				break
			}
		}
		stack.Push(num)
	}

	// Build result for nums1
	result := make([]int, len(nums1))
	for i, num := range nums1 {
		if val, exists := nextGreater[num]; exists {
			result[i] = val
		} else {
			result[i] = -1
		}
	}

	return result
}

// ============= DEQUE (Double-ended Queue) =============

// Deque supports insertion and deletion at both ends
type Deque struct {
	items []int
}

// PushFront adds element to the front
func (d *Deque) PushFront(item int) {
	d.items = append([]int{item}, d.items...)
}

// PushBack adds element to the back
func (d *Deque) PushBack(item int) {
	d.items = append(d.items, item)
}

// PopFront removes element from the front
func (d *Deque) PopFront() (int, bool) {
	if len(d.items) == 0 {
		return 0, false
	}

	item := d.items[0]
	d.items = d.items[1:]
	return item, true
}

// PopBack removes element from the back
func (d *Deque) PopBack() (int, bool) {
	if len(d.items) == 0 {
		return 0, false
	}

	index := len(d.items) - 1
	item := d.items[index]
	d.items = d.items[:index]
	return item, true
}

func (d *Deque) IsEmpty() bool {
	return len(d.items) == 0
}

func (d *Deque) Size() int {
	return len(d.items)
}

// ============= SLIDING WINDOW MAXIMUM =============

// MaxSlidingWindow finds maximum in each sliding window
// LeetCode #239
func MaxSlidingWindow(nums []int, k int) []int {
	deque := &Deque{}
	result := []int{}

	for i, num := range nums {
		// Remove elements outside window
		for !deque.IsEmpty() {
			front, _ := deque.items[0], true
			if front <= i-k {
				deque.PopFront()
			} else {
				break
			}
		}

		// Remove smaller elements from back
		for !deque.IsEmpty() {
			back := deque.items[len(deque.items)-1]
			if nums[back] < num {
				deque.PopBack()
			} else {
				break
			}
		}

		deque.PushBack(i)

		// Add to result if window is complete
		if i >= k-1 {
			front, _ := deque.items[0], true
			result = append(result, nums[front])
		}
	}

	return result
}

// ============= DEMONSTRATION FUNCTIONS =============

func demonstrateStack() {
	fmt.Println("=== Stack Operations ===")
	stack := &Stack{}

	// Push elements
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	fmt.Printf("After pushing 1, 2, 3: Size = %d\n", stack.Size())

	// Peek
	if top, ok := stack.Peek(); ok {
		fmt.Printf("Top element: %d\n", top)
	}

	// Pop elements
	for !stack.IsEmpty() {
		if item, ok := stack.Pop(); ok {
			fmt.Printf("Popped: %d\n", item)
		}
	}

	fmt.Printf("Stack is empty: %t\n", stack.IsEmpty())
}

func demonstrateQueue() {
	fmt.Println("\n=== Queue Operations ===")
	queue := &Queue{}

	// Enqueue elements
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	fmt.Printf("After enqueuing 1, 2, 3: Size = %d\n", queue.Size())

	// Front
	if front, ok := queue.Front(); ok {
		fmt.Printf("Front element: %d\n", front)
	}

	// Dequeue elements
	for !queue.IsEmpty() {
		if item, ok := queue.Dequeue(); ok {
			fmt.Printf("Dequeued: %d\n", item)
		}
	}

	fmt.Printf("Queue is empty: %t\n", queue.IsEmpty())
}

func demonstrateApplications() {
	fmt.Println("\n=== Practical Applications ===")

	// Valid Parentheses
	testStrings := []string{"()", "()[]{}", "(]", "([)]", "{[]}"}
	for _, s := range testStrings {
		fmt.Printf("'%s' is valid: %t\n", s, ValidParentheses(s))
	}

	// RPN Evaluation
	rpnTokens := []string{"2", "1", "+", "3", "*"}
	fmt.Printf("RPN %v = %d\n", rpnTokens, EvaluateRPN(rpnTokens))

	// Daily Temperatures
	temps := []int{73, 74, 75, 71, 69, 72, 76, 73}
	fmt.Printf("Daily temperatures %v\n", temps)
	fmt.Printf("Days to warmer: %v\n", DailyTemperatures(temps))

	// Moving Average
	ma := NewMovingAverage(3)
	values := []int{1, 10, 3, 5}
	fmt.Print("Moving averages: ")
	for _, val := range values {
		fmt.Printf("%.1f ", ma.Next(val))
	}
	fmt.Println()
}

// ============= MAIN FUNCTION =============

func main() {
	fmt.Println("Stacks and Queues Demonstration")
	fmt.Println(strings.Repeat("=", 40))

	demonstrateStack()
	demonstrateQueue()
	demonstrateApplications()

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("• Stack (LIFO): Use for expression evaluation, backtracking, undo operations")
	fmt.Println("• Queue (FIFO): Use for BFS, task scheduling, buffering")
	fmt.Println("• Monotonic Stack: Use for next greater/smaller element problems")
	fmt.Println("• Deque: Use for sliding window maximum/minimum problems")
}
