# 🧮 Complexity Analysis - FAANG Interview Mastery
## Mathematical Rigor Meets Production Reality

### 🎯 Elite Learning Objectives
**Master algorithmic complexity analysis at the level expected by Google, Meta, Netflix, Amazon, and Apple engineers.**

- **Big-O Mastery:** Formal mathematical analysis of time and space complexity
- **Amortized Analysis:** Understanding average performance over sequences of operations
- **Best/Average/Worst Case:** When each analysis applies and how to calculate them
- **Space-Time Tradeoffs:** Strategic decisions between memory and computational efficiency
- **Go-Specific Analysis:** How Go's runtime, GC, and language features affect complexity

### 🚀 Elite Success Metrics
- [ ] **Mathematical Rigor:** Derive complexity bounds using formal mathematical methods
- [ ] **Practical Application:** Analyze any algorithm's complexity in real-time during interviews
- [ ] **Optimization Skills:** Identify bottlenecks and propose algorithmic improvements
- [ ] **Go Expertise:** Understand how Go's runtime characteristics affect theoretical complexity

---

## 📚 Learning Modules

### 📐 [Theory](./theory/)
**Mathematical Foundations of Complexity Analysis**
- Big-O, Big-Ω, Big-Θ notation and formal definitions
- Limit-based analysis and asymptotic behavior
- Master theorem for divide-and-conquer recurrences
- Amortized analysis techniques (aggregate, accounting, potential method)

### 📊 [Cheat Sheets](./cheat-sheets/)
**Quick Reference for Common Patterns**
- Data structure operation complexities
- Sorting algorithm comparison table
- Graph algorithm complexity reference
- Go-specific complexity considerations

### 🔍 [Case Studies](./case-studies/)
**Real-World Optimization Examples**
- Before/after algorithm improvements with complexity analysis
- Production system optimizations from FAANG companies
- Space-time tradeoff decision case studies
- Go-specific optimization examples

### 💻 [Practice Problems](./practice-problems/)
**FAANG-Style Complexity Analysis Questions**
- Algorithm complexity derivation problems
- Optimization identification challenges
- Trade-off analysis scenarios
- Interview-style explanation practice

### 🔧 [Go-Specific](./go-specific/)
**Go Runtime and Language Complexity Considerations**
- Slice append/copy complexity and memory allocation patterns
- Map operation complexity and hash collision behavior
- Goroutine creation/switching overhead analysis
- Garbage collector impact on algorithmic performance

---

## 🎯 Core Complexity Classes

### O(1) - Constant Time
**Examples:** Array access, hash table lookup (average), stack push/pop
```go
// Hash table lookup - O(1) average, O(n) worst case
value, exists := hashMap[key]

// Array access - O(1) always
element := array[index]
```

### O(log n) - Logarithmic Time
**Examples:** Binary search, balanced tree operations, heap operations
```go
// Binary search - O(log n)
func binarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}
```

### O(n) - Linear Time
**Examples:** Array traversal, linked list operations, simple graph traversal
```go
// Linear search - O(n)
func linearSearch(arr []int, target int) int {
    for i, val := range arr {
        if val == target {
            return i
        }
    }
    return -1
}
```

### O(n log n) - Linearithmic Time
**Examples:** Efficient sorting (merge sort, heap sort), divide-and-conquer algorithms
```go
// Merge sort - O(n log n) time, O(n) space
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    
    return merge(left, right) // O(n) operation
}
```

### O(n²) - Quadratic Time
**Examples:** Bubble sort, selection sort, nested loop algorithms
```go
// Bubble sort - O(n²) time, O(1) space
func bubbleSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
}
```

### O(2ⁿ) - Exponential Time
**Examples:** Recursive fibonacci, subset generation, brute force solutions
```go
// Naive fibonacci - O(2^n) time, O(n) space (recursion stack)
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// Optimized fibonacci - O(n) time, O(1) space
func fibonacciOptimized(n int) int {
    if n <= 1 {
        return n
    }
    
    prev, curr := 0, 1
    for i := 2; i <= n; i++ {
        prev, curr = curr, prev+curr
    }
    return curr
}
```

---

## 🧠 Advanced Analysis Techniques

### Amortized Analysis
**When individual operations have different costs, but average cost is lower**

**Example: Dynamic Array (Go slice) Append Operation**
```go
// Slice append - O(1) amortized, O(n) worst case
var slice []int
for i := 0; i < n; i++ {
    slice = append(slice, i) // Usually O(1), occasionally O(n) when resizing
}
// Total cost: O(n), Average per operation: O(1)
```

### Space Complexity Analysis
**Memory usage including auxiliary space**

**Example: Merge Sort Space Analysis**
```go
// Merge sort space complexity: O(n) auxiliary space
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    // O(n) auxiliary space for temporary arrays
    mid := len(arr) / 2
    left := make([]int, mid)
    right := make([]int, len(arr)-mid)
    
    copy(left, arr[:mid])
    copy(right, arr[mid:])
    
    return merge(mergeSort(left), mergeSort(right))
}
```

### Recursive Complexity Analysis
**Master Theorem and Recursion Tree Method**

**Master Theorem:** T(n) = aT(n/b) + f(n)
- Case 1: f(n) = O(n^(log_b(a-ε))) → T(n) = Θ(n^log_b(a))
- Case 2: f(n) = Θ(n^log_b(a)) → T(n) = Θ(n^log_b(a) * log n)
- Case 3: f(n) = Ω(n^(log_b(a+ε))) → T(n) = Θ(f(n))

**Example: Binary Search**
T(n) = T(n/2) + O(1)
- a=1, b=2, f(n)=O(1)
- log_2(1) = 0, f(n) = Θ(n^0) = Θ(1)
- Case 2: T(n) = Θ(log n)

---

## 🔧 Go-Specific Complexity Considerations

### Slice Operations
```go
// Slice append complexity
var s []int
s = append(s, 1)  // O(1) amortized, O(n) worst case

// Slice copy complexity
copy(dst, src)    // O(min(len(dst), len(src)))

// Slice reslicing complexity
sub := s[i:j]     // O(1) - no copy, shares underlying array
```

### Map Operations
```go
// Map operations complexity
m := make(map[string]int)
m["key"] = value     // O(1) average, O(n) worst case
value, exists := m["key"]  // O(1) average, O(n) worst case
delete(m, "key")     // O(1) average, O(n) worst case
```

### Goroutine Overhead
```go
// Goroutine creation overhead: ~2KB stack, ~µs creation time
go func() {
    // Goroutine work
}()

// Channel operations
ch := make(chan int, 100)
ch <- value    // O(1) for buffered channel with space
value := <-ch  // O(1) for buffered channel with data
```

### Garbage Collector Impact
```go
// GC-friendly code: minimize allocations
// Bad: creates many small allocations
func badConcat(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s  // O(n²) due to string immutability
    }
    return result
}

// Good: single allocation
func goodConcat(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)  // O(n) amortized
    }
    return builder.String()
}
```

---

## 🎯 Elite Practice Framework

### Problem-Solving Approach
1. **Identify Input Size:** What grows as n increases?
2. **Count Operations:** How many operations per input element?
3. **Find Dominant Term:** Which term grows fastest asymptotically?
4. **Consider Space:** Auxiliary space vs. total space usage
5. **Optimize:** Can we reduce time or space complexity?

### Common Patterns Recognition
- **Two nested loops over same input:** Usually O(n²)
- **Divide input in half each iteration:** Usually O(log n)  
- **Process each element once:** Usually O(n)
- **Sort then process:** Usually O(n log n)
- **Recursive calls on subproblems:** Apply Master Theorem

### Interview Communication Framework
1. **State assumptions:** "Assuming n is the array length..."
2. **Explain reasoning:** "We visit each element once, so..."
3. **Give tight bounds:** "This is Θ(n log n), not just O(n log n)"
4. **Discuss optimizations:** "We could reduce space to O(1) by..."
5. **Consider practical factors:** "In Go, this might trigger GC..."

---

## 📈 Mastery Progression

### Level 1: Recognition
- [ ] Identify O(1), O(n), O(n²), O(log n) complexities by inspection
- [ ] Understand basic time-space tradeoffs
- [ ] Calculate simple loop and recursion complexities

### Level 2: Analysis
- [ ] Derive complexity bounds using mathematical methods
- [ ] Apply Master Theorem to divide-and-conquer algorithms
- [ ] Perform amortized analysis for data structures

### Level 3: Optimization
- [ ] Identify algorithmic bottlenecks in existing code
- [ ] Propose improvements with complexity analysis
- [ ] Consider Go-specific performance characteristics

### Level 4: Elite Mastery
- [ ] Explain complex algorithms' complexity bounds during interviews
- [ ] Design algorithms with optimal complexity for given constraints
- [ ] Balance theoretical complexity with practical Go performance

---

**Elite Complexity Analysis Readiness:** You can analyze any algorithm's time and space complexity in real-time, explain your reasoning clearly, and propose optimizations that consider both theoretical bounds and Go's runtime characteristics.

**Ready to master the mathematical foundation of elite algorithmic thinking? Let's dive deep into complexity analysis! 🔥**