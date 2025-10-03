# 🔢 Arrays & Strings Fundamentals
## Week 1, Days 1-2 | Essential Data Structures for Junior Engineers

### 🎯 Learning Objectives
By the end of this section, you will:
- [ ] Understand array structure and memory allocation
- [ ] Master basic string manipulation techniques in Go
- [ ] Implement dynamic array operations using Go slices
- [ ] Solve array/string problems using two-pointer technique
- [ ] Complete 8-10 LeetCode Easy problems

---

## 📚 Concepts Overview

### 🔸 Arrays Fundamentals

**What is an Array?**
- Collection of elements stored in contiguous memory locations
- Elements accessed by index (0-based indexing)
- Fixed size arrays vs dynamic slices in Go

**Key Characteristics:**
- **Random Access:** O(1) time to access any element
- **Sequential Storage:** Elements stored next to each other in memory
- **Homogeneous:** All elements are the same data type

**Common Operations in Go:**
```go
// Array Creation and Access
arr := []int{1, 2, 3, 4, 5}  // Go slice (dynamic array)
fmt.Println(arr[0])          // Access: O(1)
arr[2] = 10                  // Update: O(1)

// Fixed size array (less common)
var fixedArr [5]int = [5]int{1, 2, 3, 4, 5}
```

### 🔸 Dynamic Arrays (Go Slices)

**Dynamic Array Operations:**
```go
// Dynamic operations using Go slices
arr := []int{}               // Empty slice
arr = append(arr, 1)         // Insert at end: O(1) amortized
arr = append([]int{5}, arr...) // Insert at beginning: O(n)
arr = arr[:len(arr)-1]       // Remove from end: O(1)
arr = arr[1:]                // Remove from beginning: O(n)
length := len(arr)           // Get size: O(1)
```

**Memory Management:**
- Slices automatically resize when capacity is exceeded
- Typically double in size when full (amortized O(1) insertion)
- May have unused space for efficiency
- Use `make([]int, 0, capacity)` to pre-allocate

### 🔸 Strings Fundamentals

**String as Array of Characters:**
```go
s := "hello"
fmt.Println(s[0])            // 'h' (byte) - Access: O(1)
fmt.Println(len(s))          // 5 - Length: O(1)

// For Unicode support, use runes
runes := []rune(s)
fmt.Println(runes[0])        // 'h' (rune/int32)
```

**String Immutability in Go:**
```go
s := "hello"
// s[0] = 'H'                // Error! Strings are immutable
newS := "H" + s[1:]          // Create new string instead

// For mutable string operations, use strings.Builder
var builder strings.Builder
builder.WriteString("H")
builder.WriteString(s[1:])
result := builder.String()   // "Hello"
```

**String Operations:**
```go
import (
    "fmt"
    "strings"
    "unicode"
)

// Common string operations
s := "hello world"
upper := strings.ToUpper(s)        // "HELLO WORLD"
words := strings.Split(s, " ")     // ["hello", "world"]
index := strings.Index(s, "world") // 6 (index of substring)
replaced := strings.Replace(s, "world", "golang", 1) // "hello golang"
```

---

## 🛠️ Implementation Challenges

### Challenge 1: Dynamic Array Implementation
Create a simple dynamic array struct that can:
- Resize automatically when capacity is reached
- Support basic operations: append, insert, delete, access

```go
package main

import (
    "fmt"
    "errors"
)

type DynamicArray struct {
    capacity int
    size     int
    data     []interface{}
}

func NewDynamicArray() *DynamicArray {
    return &DynamicArray{
        capacity: 2,
        size:     0,
        data:     make([]interface{}, 2),
    }
}

func (da *DynamicArray) Get(index int) (interface{}, error) {
    if index < 0 || index >= da.size {
        return nil, errors.New("index out of range")
    }
    return da.data[index], nil
}

func (da *DynamicArray) Append(value interface{}) {
    if da.size >= da.capacity {
        da.resize()
    }
    da.data[da.size] = value
    da.size++
}

func (da *DynamicArray) resize() {
    da.capacity *= 2
    newData := make([]interface{}, da.capacity)
    for i := 0; i < da.size; i++ {
        newData[i] = da.data[i]
    }
    da.data = newData
}

func (da *DynamicArray) Size() int {
    return da.size
}

func (da *DynamicArray) String() string {
    result := []interface{}{}
    for i := 0; i < da.size; i++ {
        result = append(result, da.data[i])
    }
    return fmt.Sprintf("%v", result)
}

// Test your implementation
func main() {
    arr := NewDynamicArray()
    arr.Append(1)
    arr.Append(2)
    arr.Append(3)
    fmt.Println(arr) // [1 2 3]
}
```

### Challenge 2: String Builder (Mutable String)
Since strings are immutable in Go, use strings.Builder for efficient string building:

```go
package main

import (
    "fmt"
    "strings"
)

// Custom StringBuilder wrapper for educational purposes
type StringBuilder struct {
    builder strings.Builder
}

func NewStringBuilder() *StringBuilder {
    return &StringBuilder{}
}

func (sb *StringBuilder) Append(s string) {
    sb.builder.WriteString(s)
}

func (sb *StringBuilder) ToString() string {
    return sb.builder.String()
}

// Usage
func main() {
    sb := NewStringBuilder()
    sb.Append("Hello")
    sb.Append(" ")
    sb.Append("World")
    result := sb.ToString() // "Hello World"
    fmt.Println(result)
}
```

---

## 🎯 Essential Problem-Solving Patterns

### Pattern 1: Two Pointers Technique

**When to Use:**
- Problems involving pairs or comparisons
- Palindrome checking
- Finding pairs that sum to target
- Removing duplicates

**Example: Valid Palindrome**
```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

func isPalindrome(s string) bool {
    // Convert to lowercase and keep only alphanumeric
    var cleaned strings.Builder
    for _, char := range s {
        if unicode.IsLetter(char) || unicode.IsDigit(char) {
            cleaned.WriteRune(unicode.ToLower(char))
        }
    }
    
    cleanedStr := cleaned.String()
    left, right := 0, len(cleanedStr)-1
    
    for left < right {
        if cleanedStr[left] != cleanedStr[right] {
            return false
        }
        left++
        right--
    }
    
    return true
}

// Test
func main() {
    fmt.Println(isPalindrome("A man, a plan, a canal: Panama")) // true
}
```

### Pattern 2: Sliding Window

**When to Use:**
- Finding substrings/subarrays with specific properties
- Maximum/minimum subarray problems
- Fixed or variable window size problems

**Example: Maximum Subarray Sum (Fixed Window)**
```go
package main

import (
    "fmt"
    "math"
)

func maxSumSubarray(arr []int, k int) int {
    if len(arr) < k {
        return 0 // or return error
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
        maxSum = int(math.Max(float64(maxSum), float64(windowSum)))
    }
    
    return maxSum
}

// Test
func main() {
    arr := []int{2, 1, 5, 1, 3, 2}
    k := 3
    fmt.Println(maxSumSubarray(arr, k)) // 9 (subarray [5,1,3])
}
```

### Pattern 3: Fast and Slow Pointers

**When to Use:**
- Cycle detection in arrays
- Finding middle element
- Detecting patterns in sequences

**Example: Find Array Cycle (Floyd's Algorithm concept)**
```go
package main

import "fmt"

func hasCycle(arr []int) bool {
    // Array contains cycle if following indices leads to revisiting an index
    // arr[i] represents next index to visit
    if len(arr) == 0 {
        return false
    }
    
    slow, fast := 0, 0
    
    for {
        // Move slow pointer one step
        if slow < 0 || slow >= len(arr) {
            return false
        }
        slow = arr[slow]
        
        // Move fast pointer two steps
        if fast < 0 || fast >= len(arr) || arr[fast] < 0 || arr[fast] >= len(arr) {
            return false
        }
        fast = arr[fast]
        
        if fast < 0 || fast >= len(arr) {
            return false
        }
        fast = arr[fast]
        
        // If they meet, there's a cycle
        if slow == fast {
            return true
        }
    }
}

func main() {
    // Example: arr[i] points to next index
    arr := []int{1, 2, 3, 1} // 0->1->2->3->1 (cycle)
    fmt.Println(hasCycle(arr)) // true
}
```

---

## 💻 Practice Problems

### Easy Level (Days 1-2)
**Must-Solve Problems:**
1. **Two Sum** (LeetCode #1)
   - Use hash map to find pair that sums to target
   - Time: O(n), Space: O(n)

2. **Remove Duplicates from Sorted Array** (LeetCode #26)
   - Two pointers to modify array in-place
   - Time: O(n), Space: O(1)

3. **Valid Palindrome** (LeetCode #125)
   - Two pointers with character filtering
   - Time: O(n), Space: O(1)

4. **Best Time to Buy and Sell Stock** (LeetCode #121)
   - Track minimum price and maximum profit
   - Time: O(n), Space: O(1)

### Implementation Practice
5. **Implement strStr()** (LeetCode #28)
   - Find substring in string
   - Learn basic string matching

6. **Plus One** (LeetCode #66)
   - Array manipulation with carry
   - Handle edge cases

7. **Merge Sorted Array** (LeetCode #88)
   - Two pointers from the end
   - In-place merging

8. **Remove Element** (LeetCode #27)
   - Two pointers for in-place removal
   - Array modification

### Challenge Problems (Optional)
9. **Container With Most Water** (LeetCode #11)
   - Two pointers optimization
   - Greedy approach

10. **3Sum** (LeetCode #15)
    - Two pointers with sorting
    - Avoiding duplicates

---

## 🔍 Common Patterns & Techniques

### 1. Array Manipulation Techniques in Go
```go
package main

import "fmt"

// Reverse slice in-place
func reverseSlice(arr []int) {
    left, right := 0, len(arr)-1
    for left < right {
        arr[left], arr[right] = arr[right], arr[left]
        left++
        right--
    }
}

// Rotate slice to the right by k steps
func rotateSlice(arr []int, k int) {
    n := len(arr)
    k = k % n // Handle k > n
    
    // Reverse entire slice
    reverseSlice(arr)
    // Reverse first k elements
    reverseSlice(arr[:k])
    // Reverse remaining elements
    reverseSlice(arr[k:])
}

func main() {
    arr := []int{1, 2, 3, 4, 5}
    rotateSlice(arr, 2)
    fmt.Println(arr) // [4 5 1 2 3]
}
```

### 2. String Processing Techniques in Go
```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

// Count character frequency
func charFrequency(s string) map[rune]int {
    freq := make(map[rune]int)
    for _, char := range s {
        freq[char]++
    }
    return freq
}

// Check if two strings are anagrams
func areAnagrams(s1, s2 string) bool {
    if len(s1) != len(s2) {
        return false
    }
    
    // Convert to slices and sort
    runes1 := []rune(strings.ToLower(s1))
    runes2 := []rune(strings.ToLower(s2))
    
    sort.Slice(runes1, func(i, j int) bool { return runes1[i] < runes1[j] })
    sort.Slice(runes2, func(i, j int) bool { return runes2[i] < runes2[j] })
    
    return string(runes1) == string(runes2)
}

// Or using frequency count
func areAnagramsFreq(s1, s2 string) bool {
    if len(s1) != len(s2) {
        return false
    }
    
    freq1 := charFrequency(strings.ToLower(s1))
    freq2 := charFrequency(strings.ToLower(s2))
    
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

func main() {
    fmt.Println(areAnagrams("listen", "silent")) // true
    fmt.Println(areAnagramsFreq("listen", "silent")) // true
}
```

---

## 📊 Time & Space Complexity Analysis

### Array Operations Complexity
| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Access by index | O(1) | O(1) |
| Insert at end | O(1) amortized | O(1) |
| Insert at beginning | O(n) | O(1) |
| Delete from end | O(1) | O(1) |
| Delete from beginning | O(n) | O(1) |
| Search (unsorted) | O(n) | O(1) |
| Search (sorted) | O(log n) | O(1) |

### String Operations Complexity
| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Access character | O(1) | O(1) |
| Concatenation | O(n + m) | O(n + m) |
| Substring | O(k) | O(k) |
| Find substring | O(n * m) | O(1) |
| Split string | O(n) | O(n) |

---

## 🎓 Study Tips & Best Practices

### 1. Memory Management Awareness in Go
- Understand the difference between arrays and slices
- Know when slices resize and the performance implications
- Use `make()` with capacity when size is known in advance
- Be aware of slice sharing underlying arrays

### 2. Off-by-One Errors Prevention
```go
// Common mistakes and fixes
arr := []int{1, 2, 3, 4, 5}

// Wrong: Index out of bounds
// for i := 0; i <= len(arr); i++ {  // Don't do this!
//     fmt.Println(arr[i])
// }

// Correct: Proper range
for i := 0; i < len(arr); i++ {
    fmt.Println(arr[i])
}

// Better: Use range when you don't need index
for _, value := range arr {
    fmt.Println(value)
}

// When you need both index and value
for i, value := range arr {
    fmt.Printf("arr[%d] = %d\n", i, value)
}
```

### 3. String Building Efficiency in Go
```go
import "strings"

// Inefficient: Creates new string each time
func badStringBuilding(words []string) string {
    result := ""
    for _, word := range words {
        result += word // O(n²) overall!
    }
    return result
}

// Efficient: Use strings.Builder
func goodStringBuilding(words []string) string {
    var builder strings.Builder
    for _, word := range words {
        builder.WriteString(word)
    }
    return builder.String() // O(n) overall
}

// Alternative: Use strings.Join when appropriate
func bestStringBuilding(words []string) string {
    return strings.Join(words, "") // O(n) and very efficient
}
```

---

## 🚀 Next Steps

### Day 2 Goals
- [ ] Complete all 8 practice problems in Go
- [ ] Implement dynamic array from scratch using Go
- [ ] Master two-pointers technique with Go slices
- [ ] Understand time/space complexity for all operations

### Preparation for Day 3 (Linked Lists)
- Review pointer concepts in Go (pointers and references)
- Understand memory allocation with `new()` and `make()`
- Practice drawing memory diagrams for Go structs

### Self-Assessment Questions
1. Can you explain why array access is O(1)?
2. What's the difference between Go arrays and slices?
3. When would you use two pointers vs sliding window?
4. How do you handle string immutability efficiently in Go?
5. When should you use `strings.Builder` vs string concatenation?

### Go-Specific Considerations
- Understand slice capacity vs length
- Know when slices share underlying arrays
- Be familiar with Go's garbage collector impact on large slices
- Practice using Go's built-in string functions effectively

**Ready to master the foundation of all data structures with Go? Let's build it right! 🔢**