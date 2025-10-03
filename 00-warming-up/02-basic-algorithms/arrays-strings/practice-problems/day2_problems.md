# 📝 Day 2 Practice Problems: Advanced Array & String Patterns
## String Manipulation and Two-Pointer Mastery

### 🎯 Today's Goals
- Master string manipulation techniques
- Apply two-pointer patterns to complex problems
- Understand sliding window concepts
- Practice array rotation and searching

---

## Problem 1: Implement strStr() (Easy) ⭐⭐
**LeetCode #28**

Find the index of the first occurrence of `needle` in `haystack`, or -1 if not found.

**Example:**
```
Input: haystack = "hello", needle = "ll"
Output: 2

Input: haystack = "aaaaa", needle = "bba"
Output: -1
```

**Your Task:**
1. Implement basic string search
2. Handle edge cases (empty strings)
3. Optimize if possible

**Template:**
```python
def str_str(haystack, needle):
    """
    Find first occurrence of needle in haystack
    Time: O(n*m), Space: O(1)
    """
    # Your code here
    pass

# Test cases
assert str_str("hello", "ll") == 2
assert str_str("aaaaa", "bba") == -1
assert str_str("", "") == 0
assert str_str("a", "a") == 0
```

**Learning Focus:**
- String searching algorithms
- Edge case handling
- Basic pattern matching

---

## Problem 2: Longest Common Prefix (Easy) ⭐⭐
**LeetCode #14**

Find the longest common prefix string amongst an array of strings.

**Example:**
```
Input: strs = ["flower","flow","flight"]
Output: "fl"

Input: strs = ["dog","racecar","car"]
Output: ""
```

**Your Task:**
1. Compare characters vertically across all strings
2. Stop when mismatch is found
3. Handle empty array edge case

**Template:**
```python
def longest_common_prefix(strs):
    """
    Find longest common prefix among strings
    Time: O(S) where S is sum of all characters, Space: O(1)
    """
    # Your code here
    pass

# Test cases
assert longest_common_prefix(["flower","flow","flight"]) == "fl"
assert longest_common_prefix(["dog","racecar","car"]) == ""
assert longest_common_prefix(["interspecies","interstellar","interstate"]) == "inters"
assert longest_common_prefix([]) == ""
```

**Learning Focus:**
- Vertical string comparison
- Early termination optimization
- Multiple string processing

---

## Problem 3: Reverse Words in a String III (Easy) ⭐⭐
**LeetCode #557**

Reverse the order of characters in each word within a sentence while preserving whitespace and word order.

**Example:**
```
Input: s = "Let's take LeetCode contest"
Output: "s'teL ekat edoCteeL tsetnoc"
```

**Your Task:**
1. Split string into words
2. Reverse each word individually
3. Join back with spaces

**Template:**
```python
def reverse_words(s):
    """
    Reverse each word in string while preserving word order
    Time: O(n), Space: O(n)
    """
    # Your code here
    pass

# Test cases
assert reverse_words("Let's take LeetCode contest") == "s'teL ekat edoCteeL tsetnoc"
assert reverse_words("God Ding") == "doG gniD"
assert reverse_words("") == ""
```

**Learning Focus:**
- String splitting and joining
- Individual word processing
- Preserving structure

---

## Problem 4: Move Zeroes (Easy) ⭐⭐⭐
**LeetCode #283**

Move all 0's to the end of array while maintaining the relative order of non-zero elements. Must do this in-place.

**Example:**
```
Input: nums = [0,1,0,3,12]
Output: [1,3,12,0,0]
```

**Your Task:**
1. Use two pointers for in-place modification
2. Maintain relative order of non-zero elements
3. Fill remaining positions with zeros

**Template:**
```python
def move_zeroes(nums):
    """
    Move all zeros to end while maintaining order
    Time: O(n), Space: O(1)
    """
    # Your code here - modify nums in-place
    pass

# Test cases
nums1 = [0,1,0,3,12]
move_zeroes(nums1)
assert nums1 == [1,3,12,0,0]

nums2 = [0]
move_zeroes(nums2)
assert nums2 == [0]
```

**Learning Focus:**
- In-place array modification
- Two-pointer technique variations
- Order preservation

---

## Problem 5: Rotate Array (Medium) ⭐⭐⭐
**LeetCode #189**

Rotate array to the right by `k` steps.

**Example:**
```
Input: nums = [1,2,3,4,5,6,7], k = 3
Output: [5,6,7,1,2,3,4]
```

**Your Task:**
1. Implement using array reversal method
2. Handle k > array length
3. Do it in-place with O(1) extra space

**Template:**
```python
def rotate(nums, k):
    """
    Rotate array right by k steps in-place
    Time: O(n), Space: O(1)
    """
    # Your code here - modify nums in-place
    pass

# Test cases
nums1 = [1,2,3,4,5,6,7]
rotate(nums1, 3)
assert nums1 == [5,6,7,1,2,3,4]

nums2 = [-1,-100,3,99]
rotate(nums2, 2)
assert nums2 == [3,99,-1,-100]
```

**Learning Focus:**
- Array reversal technique
- In-place operations
- Modular arithmetic

---

## Problem 6: Valid Anagram (Easy) ⭐⭐
**LeetCode #242**

Determine if two strings are anagrams (contain same characters with same frequency).

**Example:**
```
Input: s = "anagram", t = "nagaram"
Output: true

Input: s = "rat", t = "car"
Output: false
```

**Your Task:**
1. First solve using sorting
2. Then optimize using character frequency counting
3. Handle Unicode characters

**Template:**
```python
def is_anagram(s, t):
    """
    Check if two strings are anagrams
    Time: O(n), Space: O(1) - assuming limited character set
    """
    # Your code here
    pass

# Test cases
assert is_anagram("anagram", "nagaram") == True
assert is_anagram("rat", "car") == False
assert is_anagram("", "") == True
```

**Learning Focus:**
- Character frequency analysis
- Hash map usage
- Multiple solution approaches

---

## Bonus Problem: Container With Most Water (Medium) ⭐⭐⭐⭐
**LeetCode #11**

Find two lines that together with x-axis forms a container that holds the most water.

**Example:**
```
Input: height = [1,8,6,2,5,4,8,3,7]
Output: 49
Explanation: Lines at index 1 and 8 form container with area = 8 * 7 = 49
```

**Your Task:**
1. Use two pointers approach
2. Understand why we move the pointer with smaller height
3. Prove this gives optimal solution

**Template:**
```python
def max_area(height):
    """
    Find maximum area that can be formed by two lines
    Time: O(n), Space: O(1)
    """
    # Your code here
    pass

# Test cases
assert max_area([1,8,6,2,5,4,8,3,7]) == 49
assert max_area([1,1]) == 1
assert max_area([4,3,2,1,4]) == 16
```

**Learning Focus:**
- Greedy optimization
- Two-pointer optimization strategy
- Mathematical reasoning

---

## 🧠 Advanced Patterns Review

### Pattern 1: Two-Pointer Variations

**Same Direction (Fast & Slow):**
```python
def remove_element(nums, val):
    write_pos = 0
    for read_pos in range(len(nums)):
        if nums[read_pos] != val:
            nums[write_pos] = nums[read_pos]
            write_pos += 1
    return write_pos
```

**Opposite Direction (Left & Right):**
```python
def reverse_string(s):
    left, right = 0, len(s) - 1
    while left < right:
        s[left], s[right] = s[right], s[left]
        left += 1
        right -= 1
```

### Pattern 2: String Processing Techniques

**Character Frequency Counting:**
```python
def char_frequency(s):
    freq = {}
    for char in s:
        freq[char] = freq.get(char, 0) + 1
    return freq
```

**String Building (Efficient):**
```python
def build_string_efficiently(parts):
    return ''.join(parts)  # O(n) vs O(n²) for concatenation
```

### Pattern 3: In-Place Array Modifications

**Key Principles:**
1. Use extra variables to track positions
2. Work backwards when needed to avoid overwriting
3. Consider swapping elements instead of shifting

---

## 🎯 Problem-Solving Strategies

### For String Problems:
1. **Consider character frequency** - many problems involve counting
2. **Think about string immutability** - may need to use arrays/lists
3. **Use built-in methods wisely** - split(), join(), isalnum(), etc.
4. **Handle edge cases** - empty strings, single characters

### For Array Problems:
1. **Two pointers are powerful** - especially for in-place modifications
2. **Consider sorting first** - sometimes makes problem much easier
3. **Think about invariants** - what properties are maintained?
4. **Draw examples** - visualize what's happening to the array

---

## 📊 Daily Progress Tracker

**Problem Completion:**
- [ ] Problem 1: Implement strStr() 
- [ ] Problem 2: Longest Common Prefix
- [ ] Problem 3: Reverse Words in String III
- [ ] Problem 4: Move Zeroes
- [ ] Problem 5: Rotate Array
- [ ] Problem 6: Valid Anagram
- [ ] Bonus: Container With Most Water

**Skills Practiced:**
- [ ] String searching and pattern matching
- [ ] Character frequency analysis
- [ ] In-place array modifications
- [ ] Two-pointer technique variations
- [ ] String manipulation and building

**Time Management:**
- Target: 20-25 minutes per problem
- Include time for testing and optimization
- Practice explaining your solution

---

## 🚀 Tomorrow's Preview: Linked Lists

**Get Ready For:**
- Pointer manipulation and memory references
- Node-based data structures
- Insertion, deletion, and traversal operations
- Classic linked list problems

**Preparation:**
- Review how references/pointers work in your language
- Practice drawing linked list diagrams
- Understand the difference between arrays and linked lists

Keep building that foundation! Today's patterns will serve you well in more complex problems. 💪
