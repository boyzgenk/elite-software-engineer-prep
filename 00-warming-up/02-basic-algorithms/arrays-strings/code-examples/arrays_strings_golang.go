package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// ===== DYNAMIC ARRAY IMPLEMENTATION =====

// DynamicArray represents a resizable array
type DynamicArray struct {
	data     []interface{}
	size     int
	capacity int
}

// NewDynamicArray creates a new dynamic array with initial capacity
func NewDynamicArray() *DynamicArray {
	return &DynamicArray{
		data:     make([]interface{}, 2),
		size:     0,
		capacity: 2,
	}
}

// Get returns element at index - O(1)
func (da *DynamicArray) Get(index int) (interface{}, error) {
	if index < 0 || index >= da.size {
		return nil, errors.New("index out of range")
	}
	return da.data[index], nil
}

// Set updates element at index - O(1)
func (da *DynamicArray) Set(index int, value interface{}) error {
	if index < 0 || index >= da.size {
		return errors.New("index out of range")
	}
	da.data[index] = value
	return nil
}

// Append adds element to end - O(1) amortized
func (da *DynamicArray) Append(value interface{}) {
	if da.size >= da.capacity {
		da.resize()
	}
	da.data[da.size] = value
	da.size++
}

// Insert adds element at specific index - O(n)
func (da *DynamicArray) Insert(index int, value interface{}) error {
	if index < 0 || index > da.size {
		return errors.New("index out of range")
	}

	if da.size >= da.capacity {
		da.resize()
	}

	// Shift elements to the right
	for i := da.size; i > index; i-- {
		da.data[i] = da.data[i-1]
	}

	da.data[index] = value
	da.size++
	return nil
}

// Delete removes element at index - O(n)
func (da *DynamicArray) Delete(index int) error {
	if index < 0 || index >= da.size {
		return errors.New("index out of range")
	}

	// Shift elements to the left
	for i := index; i < da.size-1; i++ {
		da.data[i] = da.data[i+1]
	}

	da.size--
	return nil
}

// resize doubles the capacity - O(n)
func (da *DynamicArray) resize() {
	fmt.Printf("Resizing array from capacity %d to %d\n", da.capacity, da.capacity*2)

	da.capacity *= 2
	newData := make([]interface{}, da.capacity)

	for i := 0; i < da.size; i++ {
		newData[i] = da.data[i]
	}

	da.data = newData
}

// Size returns current size
func (da *DynamicArray) Size() int {
	return da.size
}

// Capacity returns current capacity
func (da *DynamicArray) Capacity() int {
	return da.capacity
}

// IsEmpty checks if array is empty
func (da *DynamicArray) IsEmpty() bool {
	return da.size == 0
}

// String returns string representation
func (da *DynamicArray) String() string {
	result := "["
	for i := 0; i < da.size; i++ {
		result += fmt.Sprintf("%v", da.data[i])
		if i < da.size-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

// ===== STRING BUILDER IMPLEMENTATION =====

// StringBuilder provides efficient string concatenation
type StringBuilder struct {
	buffer []string
}

// NewStringBuilder creates a new string builder
func NewStringBuilder() *StringBuilder {
	return &StringBuilder{
		buffer: make([]string, 0),
	}
}

// Append adds string to builder - O(1)
func (sb *StringBuilder) Append(s string) {
	sb.buffer = append(sb.buffer, s)
}

// AppendChar adds single character
func (sb *StringBuilder) AppendChar(c rune) {
	sb.buffer = append(sb.buffer, string(c))
}

// ToString returns final string - O(n)
func (sb *StringBuilder) ToString() string {
	return strings.Join(sb.buffer, "")
}

// Length returns total length of all strings
func (sb *StringBuilder) Length() int {
	total := 0
	for _, s := range sb.buffer {
		total += len(s)
	}
	return total
}

// Clear resets the builder
func (sb *StringBuilder) Clear() {
	sb.buffer = sb.buffer[:0]
}

// ===== TWO POINTERS PATTERNS =====

// IsPalindrome checks if string is palindrome using two pointers
func IsPalindrome(s string) bool {
	// Convert to lowercase and keep only alphanumeric
	var cleaned strings.Builder
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			cleaned.WriteRune(unicode.ToLower(char))
		}
	}

	cleanStr := cleaned.String()
	left, right := 0, len(cleanStr)-1

	for left < right {
		if cleanStr[left] != cleanStr[right] {
			return false
		}
		left++
		right--
	}

	return true
}

// TwoSum finds indices of two numbers that sum to target
func TwoSum(nums []int, target int) []int {
	numMap := make(map[int]int) // value -> index

	for i, num := range nums {
		complement := target - num
		if index, exists := numMap[complement]; exists {
			return []int{index, i}
		}
		numMap[num] = i
	}

	return []int{} // No solution found
}

// RemoveDuplicates removes duplicates from sorted array in-place
func RemoveDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}

	writeIndex := 1

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			nums[writeIndex] = nums[i]
			writeIndex++
		}
	}

	return writeIndex
}

// ===== SLIDING WINDOW PATTERNS =====

// MaxSumSubarray finds maximum sum of subarray of size k
func MaxSumSubarray(arr []int, k int) (int, error) {
	if len(arr) < k {
		return 0, errors.New("array size is less than k")
	}

	// Calculate sum of first window
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}

	maxSum := windowSum

	// Slide the window
	for i := k; i < len(arr); i++ {
		windowSum = windowSum - arr[i-k] + arr[i]
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum, nil
}

// LongestSubstringNoRepeats finds longest substring without repeating characters
func LongestSubstringNoRepeats(s string) int {
	charMap := make(map[rune]int) // character -> last seen index
	maxLength := 0
	start := 0

	for i, char := range s {
		if lastIndex, exists := charMap[char]; exists && lastIndex >= start {
			start = lastIndex + 1
		}

		charMap[char] = i
		currentLength := i - start + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}

// ===== ARRAY MANIPULATION TECHNIQUES =====

// ReverseArray reverses array in-place
func ReverseArray(arr []int) {
	left, right := 0, len(arr)-1

	for left < right {
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}
}

// RotateArray rotates array to the right by k steps
func RotateArray(arr []int, k int) {
	if len(arr) == 0 {
		return
	}

	k = k % len(arr) // Handle k > len(arr)

	// Reverse entire array
	ReverseArray(arr)

	// Reverse first k elements
	reverseRange(arr, 0, k-1)

	// Reverse remaining elements
	reverseRange(arr, k, len(arr)-1)
}

func reverseRange(arr []int, start, end int) {
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

// ===== STRING PROCESSING TECHNIQUES =====

// CharFrequency counts character frequency in string
func CharFrequency(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, char := range s {
		freq[char]++
	}
	return freq
}

// AreAnagrams checks if two strings are anagrams
func AreAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Convert to lowercase and sort
	runes1 := []rune(strings.ToLower(s1))
	runes2 := []rune(strings.ToLower(s2))

	sort.Slice(runes1, func(i, j int) bool {
		return runes1[i] < runes1[j]
	})
	sort.Slice(runes2, func(i, j int) bool {
		return runes2[i] < runes2[j]
	})

	return string(runes1) == string(runes2)
}

// AreAnagramsFreq checks anagrams using frequency counting
func AreAnagramsFreq(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	freq1 := CharFrequency(strings.ToLower(s1))
	freq2 := CharFrequency(strings.ToLower(s2))

	if len(freq1) != len(freq2) {
		return false
	}

	for char, count := range freq1 {
		if freq2[char] != count {
			return false
		}
	}

	return true
}

// ===== DEMONSTRATION FUNCTION =====

func main() {
	fmt.Println("🔢 Arrays & Strings Fundamentals in Golang")
	fmt.Println("==========================================")

	// Dynamic Array demonstration
	fmt.Println("Dynamic Array Demo:")
	da := NewDynamicArray()

	fmt.Printf("Initial: %s (size: %d, capacity: %d)\n", da.String(), da.Size(), da.Capacity())

	da.Append(1)
	da.Append(2)
	da.Append(3) // This should trigger resize

	fmt.Printf("After appending 1,2,3: %s (size: %d, capacity: %d)\n",
		da.String(), da.Size(), da.Capacity())

	da.Insert(1, 99)
	fmt.Printf("After inserting 99 at index 1: %s\n", da.String())

	da.Delete(2)
	fmt.Printf("After deleting index 2: %s\n", da.String())

	// String Builder demonstration
	fmt.Println("\nString Builder Demo:")
	sb := NewStringBuilder()
	sb.Append("Hello")
	sb.Append(" ")
	sb.Append("World")
	sb.Append("!")

	fmt.Printf("Built string: '%s' (length: %d)\n", sb.ToString(), sb.Length())

	// Pattern demonstrations
	fmt.Println("\nPattern Demonstrations:")

	// Palindrome check
	testStr := "A man, a plan, a canal: Panama"
	fmt.Printf("Is '%s' a palindrome? %t\n", testStr, IsPalindrome(testStr))

	// Two Sum
	nums := []int{2, 7, 11, 15}
	target := 9
	result := TwoSum(nums, target)
	fmt.Printf("Two Sum: nums=%v, target=%d, result=%v\n", nums, target, result)

	// Remove duplicates
	duplicates := []int{1, 1, 2, 2, 2, 3, 4, 4, 5}
	uniqueLength := RemoveDuplicates(duplicates)
	fmt.Printf("Remove duplicates: %v -> length=%d, array=%v\n",
		[]int{1, 1, 2, 2, 2, 3, 4, 4, 5}, uniqueLength, duplicates[:uniqueLength])

	// Max subarray sum
	arr := []int{2, 1, 5, 1, 3, 2}
	k := 3
	maxSum, _ := MaxSumSubarray(arr, k)
	fmt.Printf("Max sum subarray of size %d in %v: %d\n", k, arr, maxSum)

	// Longest substring without repeating chars
	longStr := "abcabcbb"
	maxLen := LongestSubstringNoRepeats(longStr)
	fmt.Printf("Longest substring without repeating chars in '%s': %d\n", longStr, maxLen)

	// Array rotation
	rotateArr := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("Original array: %v\n", rotateArr)
	RotateArray(rotateArr, 3)
	fmt.Printf("After rotating right by 3: %v\n", rotateArr)

	// Anagram check
	word1, word2 := "listen", "silent"
	fmt.Printf("Are '%s' and '%s' anagrams? %t\n", word1, word2, AreAnagrams(word1, word2))
}
