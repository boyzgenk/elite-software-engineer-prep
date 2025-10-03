# Hash Tables Practice Problems - Day 4

## Learning Objectives
By the end of this session, you should be able to:
- Apply hash tables to solve frequency counting problems
- Use hash maps for fast lookups and duplicate detection
- Implement two-pointer techniques with hash table support
- Handle edge cases in hash table problems

## Problem Templates

### Problem 1: Two Sum (Easy)
**LeetCode #1**

**Problem Statement:**
Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to target.

**Examples:**
```
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: nums[0] + nums[1] = 2 + 7 = 9

Input: nums = [3,2,4], target = 6
Output: [1,2]

Input: nums = [3,3], target = 6
Output: [0,1]
```

**Golang Template:**
```go
func twoSum(nums []int, target int) []int {
    // TODO: Implement using hash map
    // Time: O(n), Space: O(n)
    
}
```

**Approach:**
1. Create a hash map to store value -> index mapping
2. For each number, check if (target - current) exists in map
3. If found, return both indices
4. Otherwise, add current number to map

---

### Problem 2: Valid Anagram (Easy)
**LeetCode #242**

**Problem Statement:**
Given two strings `s` and `t`, return true if `t` is an anagram of `s`, and false otherwise.

**Examples:**
```
Input: s = "anagram", t = "nagaram"
Output: true

Input: s = "rat", t = "car"
Output: false
```

**Golang Template:**
```go
func isAnagram(s string, t string) bool {
    // TODO: Implement using hash map for character counting
    // Time: O(n), Space: O(1) - at most 26 characters
    
}
```

**Approach:**
1. Check if lengths are equal (early return if not)
2. Count frequency of each character in first string
3. Subtract frequency for each character in second string
4. Check if all counts are zero

---

### Problem 3: Group Anagrams (Medium)
**LeetCode #49**

**Problem Statement:**
Given an array of strings `strs`, group the anagrams together. You can return the answer in any order.

**Examples:**
```
Input: strs = ["eat","tea","tan","ate","nat","bat"]
Output: [["bat"],["nat","tan"],["ate","eat","tea"]]

Input: strs = [""]
Output: [[""]]

Input: strs = ["a"]
Output: [["a"]]
```

**Golang Template:**
```go
func groupAnagrams(strs []string) [][]string {
    // TODO: Implement using hash map with sorted string as key
    // Time: O(n * k log k), Space: O(n * k)
    // where n = number of strings, k = max length of string
    
}
```

**Approach:**
1. Create hash map: sorted_string -> list of original strings
2. For each string, sort it to get the key
3. Add original string to the list for that key
4. Return all values from the hash map

---

### Problem 4: Contains Duplicate (Easy)
**LeetCode #217**

**Problem Statement:**
Given an integer array `nums`, return true if any value appears at least twice in the array, and return false if every element is distinct.

**Examples:**
```
Input: nums = [1,2,3,1]
Output: true

Input: nums = [1,2,3,4]
Output: false

Input: nums = [1,1,1,3,3,4,3,2,4,2]
Output: true
```

**Golang Template:**
```go
func containsDuplicate(nums []int) bool {
    // TODO: Implement using hash set
    // Time: O(n), Space: O(n)
    
}
```

**Approach:**
1. Create a hash set (map[int]bool in Go)
2. For each number, check if it exists in set
3. If found, return true (duplicate found)
4. Otherwise, add to set
5. If loop completes, return false

---

### Problem 5: First Unique Character (Easy)
**LeetCode #387**

**Problem Statement:**
Given a string `s`, find the first non-repeating character in it and return its index. If it does not exist, return -1.

**Examples:**
```
Input: s = "leetcode"
Output: 0

Input: s = "loveleetcode"
Output: 2

Input: s = "aabb"
Output: -1
```

**Golang Template:**
```go
func firstUniqChar(s string) int {
    // TODO: Implement using hash map for frequency counting
    // Time: O(n), Space: O(1) - at most 26 characters
    
}
```

**Approach:**
1. First pass: count frequency of each character
2. Second pass: find first character with frequency 1
3. Return its index, or -1 if not found

---

## Advanced Problems

### Problem 6: Subarray Sum Equals K (Medium)
**LeetCode #560**

**Problem Statement:**
Given an array of integers `nums` and an integer `k`, return the total number of continuous subarrays whose sum equals to `k`.

**Examples:**
```
Input: nums = [1,1,1], k = 2
Output: 2

Input: nums = [1,2,3], k = 3
Output: 2
```

**Golang Template:**
```go
func subarraySum(nums []int, k int) int {
    // TODO: Implement using prefix sum + hash map
    // Time: O(n), Space: O(n)
    
}
```

**Approach:**
1. Use prefix sum technique with hash map
2. Map: prefix_sum -> count of occurrences
3. For each position, check if (prefix_sum - k) exists
4. If found, add its count to result
5. Update map with current prefix sum

---

### Problem 7: Longest Consecutive Sequence (Medium)
**LeetCode #128**

**Problem Statement:**
Given an unsorted array of integers `nums`, return the length of the longest consecutive elements sequence.

**Examples:**
```
Input: nums = [100,4,200,1,3,2]
Output: 4
Explanation: The longest consecutive sequence is [1, 2, 3, 4]. Therefore its length is 4.

Input: nums = [0,3,7,2,5,8,4,6,0,1]
Output: 9
```

**Golang Template:**
```go
func longestConsecutive(nums []int) int {
    // TODO: Implement using hash set
    // Time: O(n), Space: O(n)
    
}
```

**Approach:**
1. Put all numbers in hash set for O(1) lookup
2. For each number, check if it's the start of a sequence (num-1 not in set)
3. If it's a start, count consecutive numbers
4. Track maximum length found

---

## Daily Study Plan

### Morning Session (45 minutes)
1. **Review Theory (10 min)**: Hash table operations, collision handling
2. **Solve Problems 1-3 (30 min)**: Focus on basic hash map patterns
3. **Code Review (5 min)**: Check solutions, optimize if needed

### Evening Session (45 minutes)
1. **Solve Problems 4-5 (25 min)**: Practice frequency counting
2. **Attempt Advanced Problems 6-7 (15 min)**: Challenge yourself
3. **Reflection (5 min)**: Note patterns and key insights

## Success Checklist
- [ ] Can implement Two Sum in under 5 minutes
- [ ] Understand when to use hash map vs hash set
- [ ] Can explain time/space complexity trade-offs
- [ ] Comfortable with frequency counting patterns
- [ ] Can handle edge cases (empty arrays, duplicates)
- [ ] Understand prefix sum + hash map technique

## Key Patterns to Remember
1. **Fast Lookup**: Use hash map when you need O(1) lookup
2. **Frequency Counting**: Count occurrences of elements
3. **Complement Search**: Look for target - current_value
4. **Grouping**: Use hash map to group related items
5. **Prefix Sum**: Combine with hash map for subarray problems

## Time Complexity Quick Reference
- Hash map operations (get, put, delete): O(1) average, O(n) worst case
- Building frequency map: O(n)
- Two Sum pattern: O(n) time, O(n) space
- Anagram checking: O(n) time, O(1) space (fixed alphabet)