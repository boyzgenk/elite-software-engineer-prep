package main

import (
	"fmt"
	"sort"
	"unicode"
)

// ===== ARRAYS & STRINGS PRACTICE PROBLEMS SOLUTIONS =====

// Problem 1: Two Sum (LeetCode #1)
// Given array of integers and target, return indices of two numbers that add up to target
func twoSum(nums []int, target int) []int {
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

// Problem 2: Remove Duplicates from Sorted Array (LeetCode #26)
// Remove duplicates in-place and return new length
func removeDuplicates(nums []int) int {
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

// Problem 3: Valid Palindrome (LeetCode #125)
// Check if string is palindrome considering only alphanumeric characters
func isPalindrome(s string) bool {
	left, right := 0, len(s)-1

	for left < right {
		// Skip non-alphanumeric characters from left
		for left < right && !isAlphanumeric(rune(s[left])) {
			left++
		}

		// Skip non-alphanumeric characters from right
		for left < right && !isAlphanumeric(rune(s[right])) {
			right--
		}

		// Compare characters (case insensitive)
		if unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])) {
			return false
		}

		left++
		right--
	}

	return true
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// Problem 4: Best Time to Buy and Sell Stock (LeetCode #121)
// Find maximum profit from buying and selling stock once
func maxProfit(prices []int) int {
	if len(prices) <= 1 {
		return 0
	}

	minPrice := prices[0]
	maxProfitSoFar := 0

	for i := 1; i < len(prices); i++ {
		// Update minimum price seen so far
		if prices[i] < minPrice {
			minPrice = prices[i]
		}

		// Calculate profit if we sell today
		profit := prices[i] - minPrice
		if profit > maxProfitSoFar {
			maxProfitSoFar = profit
		}
	}

	return maxProfitSoFar
}

// Problem 5: Implement strStr() (LeetCode #28)
// Find first occurrence of needle in haystack
func strStr(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}

	if len(needle) > len(haystack) {
		return -1
	}

	// Check each possible starting position
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}

	return -1
}

// Problem 6: Plus One (LeetCode #66)
// Add one to number represented as array of digits
func plusOne(digits []int) []int {
	// Start from the rightmost digit
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0 // Set to 0 and carry over
	}

	// If we reach here, all digits were 9
	// Need to create new array with extra digit
	result := make([]int, len(digits)+1)
	result[0] = 1
	return result
}

// Problem 7: Merge Sorted Array (LeetCode #88)
// Merge nums2 into nums1 in-place
func merge(nums1 []int, m int, nums2 []int, n int) {
	// Start from the end to avoid overwriting
	i := m - 1
	j := n - 1
	k := m + n - 1

	for i >= 0 && j >= 0 {
		if nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
		k--
	}

	// Copy remaining elements from nums2
	for j >= 0 {
		nums1[k] = nums2[j]
		j--
		k--
	}
}

// Problem 8: Remove Element (LeetCode #27)
// Remove all instances of val in-place and return new length
func removeElement(nums []int, val int) int {
	writeIndex := 0

	for _, num := range nums {
		if num != val {
			nums[writeIndex] = num
			writeIndex++
		}
	}

	return writeIndex
}

// Problem 9: Container With Most Water (LeetCode #11)
// Find two lines that form container with most water
func maxArea(height []int) int {
	left, right := 0, len(height)-1
	maxWater := 0

	for left < right {
		// Calculate water with current pointers
		width := right - left
		minHeight := min(height[left], height[right])
		water := width * minHeight

		if water > maxWater {
			maxWater = water
		}

		// Move pointer with smaller height
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}

	return maxWater
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Problem 10: 3Sum (LeetCode #15)
// Find all unique triplets that sum to zero
func threeSum(nums []int) [][]int {
	result := [][]int{}
	if len(nums) < 3 {
		return result
	}

	// Sort array first
	sort.Ints(nums)

	for i := 0; i < len(nums)-2; i++ {
		// Skip duplicates for first number
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		left, right := i+1, len(nums)-1
		target := -nums[i]

		for left < right {
			sum := nums[left] + nums[right]
			if sum == target {
				result = append(result, []int{nums[i], nums[left], nums[right]})

				// Skip duplicates
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}

				left++
				right--
			} else if sum < target {
				left++
			} else {
				right--
			}
		}
	}

	return result
}

// ===== ADDITIONAL STRING PROBLEMS =====

// Problem 11: Valid Anagram (LeetCode #242)
// Check if two strings are anagrams
func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	charCount := make(map[rune]int)

	// Count characters in first string
	for _, char := range s {
		charCount[char]++
	}

	// Subtract characters in second string
	for _, char := range t {
		charCount[char]--
		if charCount[char] == 0 {
			delete(charCount, char)
		}
	}

	return len(charCount) == 0
}

// Problem 12: Group Anagrams (LeetCode #49)
// Group strings that are anagrams of each other
func groupAnagrams(strs []string) [][]string {
	groups := make(map[string][]string)

	for _, str := range strs {
		// Sort characters to create key
		chars := []rune(str)
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})
		key := string(chars)

		groups[key] = append(groups[key], str)
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// Problem 13: Longest Common Prefix (LeetCode #14)
// Find longest common prefix among array of strings
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	// Find shortest string length
	minLen := len(strs[0])
	for _, str := range strs {
		if len(str) < minLen {
			minLen = len(str)
		}
	}

	// Check each character position
	for i := 0; i < minLen; i++ {
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if strs[j][i] != char {
				return strs[0][:i]
			}
		}
	}

	return strs[0][:minLen]
}

// Problem 14: String to Integer (atoi) (LeetCode #8)
// Convert string to integer with various edge cases
func myAtoi(s string) int {
	const (
		INT_MAX = 2147483647
		INT_MIN = -2147483648
	)

	i := 0
	n := len(s)

	// Skip leading whitespace
	for i < n && s[i] == ' ' {
		i++
	}

	if i == n {
		return 0
	}

	// Handle sign
	sign := 1
	if s[i] == '+' || s[i] == '-' {
		if s[i] == '-' {
			sign = -1
		}
		i++
	}

	// Convert digits
	result := 0
	for i < n && isDigit(s[i]) {
		digit := int(s[i] - '0')

		// Check for overflow
		if result > INT_MAX/10 || (result == INT_MAX/10 && digit > 7) {
			if sign == 1 {
				return INT_MAX
			} else {
				return INT_MIN
			}
		}

		result = result*10 + digit
		i++
	}

	return sign * result
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// ===== DEMONSTRATION FUNCTION =====

func main() {
	fmt.Println("🔢 Arrays & Strings - Practice Problems Solutions")
	fmt.Println("=================================================")

	// Test Two Sum
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Printf("Two Sum: nums=%v, target=%d, result=%v\n", nums, target, result)

	// Test Remove Duplicates
	duplicates := []int{1, 1, 2, 2, 2, 3, 4, 4, 5}
	originalDups := make([]int, len(duplicates))
	copy(originalDups, duplicates)
	uniqueLength := removeDuplicates(duplicates)
	fmt.Printf("Remove Duplicates: %v -> length=%d, unique=%v\n",
		originalDups, uniqueLength, duplicates[:uniqueLength])

	// Test Valid Palindrome
	testStr := "A man, a plan, a canal: Panama"
	fmt.Printf("Is '%s' palindrome? %t\n", testStr, isPalindrome(testStr))

	// Test Max Profit
	prices := []int{7, 1, 5, 3, 6, 4}
	fmt.Printf("Max profit from %v: %d\n", prices, maxProfit(prices))

	// Test String Search
	haystack := "hello world"
	needle := "world"
	fmt.Printf("Find '%s' in '%s': index %d\n", needle, haystack, strStr(haystack, needle))

	// Test Plus One
	digits := []int{9, 9, 9}
	fmt.Printf("Plus one: %v -> %v\n", digits, plusOne(digits))

	// Test Container With Most Water
	heights := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Printf("Max water area with heights %v: %d\n", heights, maxArea(heights))

	// Test Three Sum
	nums3 := []int{-1, 0, 1, 2, -1, -4}
	fmt.Printf("Three sum triplets in %v: %v\n", nums3, threeSum(nums3))

	// Test Anagrams
	word1, word2 := "listen", "silent"
	fmt.Printf("Are '%s' and '%s' anagrams? %t\n", word1, word2, isAnagram(word1, word2))

	// Test Group Anagrams
	words := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	groups := groupAnagrams(words)
	fmt.Printf("Grouped anagrams: %v\n", groups)

	// Test Longest Common Prefix
	strs := []string{"flower", "flow", "flight"}
	fmt.Printf("Longest common prefix of %v: '%s'\n", strs, longestCommonPrefix(strs))

	// Test String to Integer
	testAtoi := "   -42"
	fmt.Printf("String to integer '%s': %d\n", testAtoi, myAtoi(testAtoi))
}
