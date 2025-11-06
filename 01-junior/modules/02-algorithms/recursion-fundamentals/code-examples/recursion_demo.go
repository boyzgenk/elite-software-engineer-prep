package main

import (
	"fmt"
	"strings"
)

// ============= BASIC RECURSION PATTERNS =============

// Factorial - classic example of linear recursion
func Factorial(n int) int {
	// Base case
	if n <= 1 {
		return 1
	}

	// Recursive case: n * factorial(n-1)
	return n * Factorial(n-1)
}

// Fibonacci - example with multiple recursive calls
func Fibonacci(n int) int {
	// Base cases
	if n <= 1 {
		return n
	}

	// Recursive case: fib(n-1) + fib(n-2)
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciMemo - optimized with memoization
func FibonacciMemo(n int) int {
	memo := make(map[int]int)
	return fibHelper(n, memo)
}

func fibHelper(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}

	if val, exists := memo[n]; exists {
		return val
	}

	memo[n] = fibHelper(n-1, memo) + fibHelper(n-2, memo)
	return memo[n]
}

// ============= ARRAY RECURSION =============

// SumArray calculates sum of array elements recursively
func SumArray(arr []int) int {
	// Base case: empty array
	if len(arr) == 0 {
		return 0
	}

	// Recursive case: first element + sum of rest
	return arr[0] + SumArray(arr[1:])
}

// FindMax finds maximum element in array recursively
func FindMax(arr []int) int {
	// Base case: single element
	if len(arr) == 1 {
		return arr[0]
	}

	// Recursive case: max of first and max of rest
	maxRest := FindMax(arr[1:])
	if arr[0] > maxRest {
		return arr[0]
	}
	return maxRest
}

// BinarySearch performs binary search recursively
func BinarySearch(arr []int, target int) int {
	return binarySearchHelper(arr, target, 0, len(arr)-1)
}

func binarySearchHelper(arr []int, target, left, right int) int {
	// Base case: element not found
	if left > right {
		return -1
	}

	mid := left + (right-left)/2

	// Base case: element found
	if arr[mid] == target {
		return mid
	}

	// Recursive cases
	if target < arr[mid] {
		return binarySearchHelper(arr, target, left, mid-1)
	}
	return binarySearchHelper(arr, target, mid+1, right)
}

// ============= STRING RECURSION =============

// ReverseString reverses a string recursively
func ReverseString(s string) string {
	// Base case: empty or single character
	if len(s) <= 1 {
		return s
	}

	// Recursive case: last char + reverse of rest
	return string(s[len(s)-1]) + ReverseString(s[:len(s)-1])
}

// IsPalindrome checks if string is palindrome recursively
func IsPalindrome(s string) bool {
	// Base case: empty or single character
	if len(s) <= 1 {
		return true
	}

	// Check first and last characters
	if s[0] != s[len(s)-1] {
		return false
	}

	// Recursive case: check middle substring
	return IsPalindrome(s[1 : len(s)-1])
}

// CountOccurrences counts occurrences of character in string
func CountOccurrences(s string, char rune) int {
	// Base case: empty string
	if len(s) == 0 {
		return 0
	}

	count := 0
	if rune(s[0]) == char {
		count = 1
	}

	// Recursive case: count in first + count in rest
	return count + CountOccurrences(s[1:], char)
}

// ============= MATHEMATICAL RECURSION =============

// Power calculates base^exp recursively
func Power(base, exp int) int {
	// Base case
	if exp == 0 {
		return 1
	}

	// Recursive case
	return base * Power(base, exp-1)
}

// PowerOptimized - more efficient O(log n) approach
func PowerOptimized(base, exp int) int {
	// Base case
	if exp == 0 {
		return 1
	}

	// If exponent is even: base^exp = (base^(exp/2))^2
	if exp%2 == 0 {
		half := PowerOptimized(base, exp/2)
		return half * half
	}

	// If exponent is odd: base^exp = base * base^(exp-1)
	return base * PowerOptimized(base, exp-1)
}

// GCD calculates Greatest Common Divisor using Euclidean algorithm
func GCD(a, b int) int {
	// Base case
	if b == 0 {
		return a
	}

	// Recursive case: gcd(b, a % b)
	return GCD(b, a%b)
}

// ============= BACKTRACKING PATTERNS =============

// GenerateParentheses generates all valid parentheses combinations
// LeetCode #22
func GenerateParentheses(n int) []string {
	result := []string{}
	generateHelper("", 0, 0, n, &result)
	return result
}

func generateHelper(current string, open, close, max int, result *[]string) {
	// Base case: reached maximum length
	if len(current) == max*2 {
		*result = append(*result, current)
		return
	}

	// Add opening parenthesis if we haven't used all
	if open < max {
		generateHelper(current+"(", open+1, close, max, result)
	}

	// Add closing parenthesis if it doesn't exceed opening
	if close < open {
		generateHelper(current+")", open, close+1, max, result)
	}
}

// Permutations generates all permutations of an array
func Permutations(nums []int) [][]int {
	result := [][]int{}
	permuteHelper(nums, []int{}, &result)
	return result
}

func permuteHelper(nums []int, current []int, result *[][]int) {
	// Base case: no more numbers to permute
	if len(nums) == 0 {
		// Make a copy of current permutation
		perm := make([]int, len(current))
		copy(perm, current)
		*result = append(*result, perm)
		return
	}

	// Try each remaining number
	for i, num := range nums {
		// Choose: add number to current permutation
		newCurrent := append(current, num)

		// Create new array without chosen number
		newNums := append([]int{}, nums[:i]...)
		newNums = append(newNums, nums[i+1:]...)

		// Explore: recurse with remaining numbers
		permuteHelper(newNums, newCurrent, result)

		// Unchoose: backtrack (automatic due to slice creation)
	}
}

// ============= TREE RECURSION EXAMPLES =============

// TreeNode for tree problems
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// TreeHeight calculates height of binary tree
func TreeHeight(root *TreeNode) int {
	// Base case: empty tree
	if root == nil {
		return 0
	}

	// Recursive case: 1 + max height of subtrees
	leftHeight := TreeHeight(root.Left)
	rightHeight := TreeHeight(root.Right)

	return 1 + max(leftHeight, rightHeight)
}

// TreeSum calculates sum of all nodes in tree
func TreeSum(root *TreeNode) int {
	// Base case: empty tree
	if root == nil {
		return 0
	}

	// Recursive case: current value + sum of subtrees
	return root.Val + TreeSum(root.Left) + TreeSum(root.Right)
}

// TreePaths finds all root-to-leaf paths
func TreePaths(root *TreeNode) []string {
	if root == nil {
		return []string{}
	}

	result := []string{}
	findPaths(root, "", &result)
	return result
}

func findPaths(node *TreeNode, path string, result *[]string) {
	if node == nil {
		return
	}

	// Add current node to path
	if path == "" {
		path = fmt.Sprintf("%d", node.Val)
	} else {
		path = fmt.Sprintf("%s->%d", path, node.Val)
	}

	// If leaf node, add path to result
	if node.Left == nil && node.Right == nil {
		*result = append(*result, path)
		return
	}

	// Recurse on children
	findPaths(node.Left, path, result)
	findPaths(node.Right, path, result)
}

// ============= DIVIDE AND CONQUER =============

// MergeSort implements merge sort algorithm
func MergeSort(arr []int) []int {
	// Base case: single element or empty
	if len(arr) <= 1 {
		return arr
	}

	// Divide: split array in half
	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	// Conquer: merge sorted halves
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	// Merge while both arrays have elements
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Add remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

// QuickSort implements quicksort algorithm
func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	return quickSortHelper(arr, 0, len(arr)-1)
}

func quickSortHelper(arr []int, low, high int) []int {
	if low < high {
		// Partition and get pivot index
		pivotIndex := partition(arr, low, high)

		// Recursively sort elements before and after partition
		quickSortHelper(arr, low, pivotIndex-1)
		quickSortHelper(arr, pivotIndex+1, high)
	}
	return arr
}

func partition(arr []int, low, high int) int {
	pivot := arr[high] // Choose last element as pivot
	i := low - 1       // Index of smaller element

	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// ============= UTILITY FUNCTIONS =============

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============= DEMONSTRATION =============

func main() {
	fmt.Println("Recursion Examples and Patterns")
	fmt.Println(strings.Repeat("=", 40))

	// Basic recursion
	fmt.Println("=== Basic Recursion ===")
	fmt.Printf("Factorial(5) = %d\n", Factorial(5))
	fmt.Printf("Fibonacci(8) = %d\n", Fibonacci(8))
	fmt.Printf("FibonacciMemo(8) = %d\n", FibonacciMemo(8))

	// Array recursion
	fmt.Println("\n=== Array Recursion ===")
	arr := []int{1, 5, 3, 8, 2}
	fmt.Printf("Array: %v\n", arr)
	fmt.Printf("Sum: %d\n", SumArray(arr))
	fmt.Printf("Max: %d\n", FindMax(arr))

	sortedArr := []int{1, 3, 5, 7, 9, 11}
	fmt.Printf("Binary search for 7 in %v: index %d\n", sortedArr, BinarySearch(sortedArr, 7))

	// String recursion
	fmt.Println("\n=== String Recursion ===")
	str := "hello"
	fmt.Printf("Reverse of '%s': '%s'\n", str, ReverseString(str))
	palindrome := "racecar"
	fmt.Printf("'%s' is palindrome: %t\n", palindrome, IsPalindrome(palindrome))
	fmt.Printf("Count of 'l' in '%s': %d\n", str, CountOccurrences(str, 'l'))

	// Mathematical recursion
	fmt.Println("\n=== Mathematical Recursion ===")
	fmt.Printf("2^10 = %d\n", Power(2, 10))
	fmt.Printf("2^10 (optimized) = %d\n", PowerOptimized(2, 10))
	fmt.Printf("GCD(48, 18) = %d\n", GCD(48, 18))

	// Backtracking
	fmt.Println("\n=== Backtracking ===")
	parens := GenerateParentheses(3)
	fmt.Printf("Valid parentheses for n=3: %v\n", parens)

	nums := []int{1, 2, 3}
	perms := Permutations(nums)
	fmt.Printf("Permutations of %v: %v\n", nums, perms)

	// Divide and conquer
	fmt.Println("\n=== Divide and Conquer ===")
	unsorted := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original: %v\n", unsorted)

	mergeSorted := MergeSort(append([]int{}, unsorted...))
	fmt.Printf("Merge sorted: %v\n", mergeSorted)

	quickSorted := QuickSort(append([]int{}, unsorted...))
	fmt.Printf("Quick sorted: %v\n", quickSorted)

	fmt.Println("\n=== Key Recursion Principles ===")
	fmt.Println("1. Base Case: Condition to stop recursion")
	fmt.Println("2. Recursive Case: Function calls itself with modified input")
	fmt.Println("3. Progress: Each call should get closer to base case")
	fmt.Println("4. Trust: Assume recursive calls work correctly")
	fmt.Println("5. Combine: Use results from recursive calls appropriately")
}
