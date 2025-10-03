# 📝 Day 1 Practice Problems: Arrays & Strings Basics
## Easy Level Problems for Foundation Building

### 🎯 Today's Goals
- Solve 4-5 problems to build confidence
- Focus on understanding rather than speed
- Practice explaining your approach out loud
- Master basic array and string operations

---

## Problem 1: Two Sum (Easy) ⭐⭐
**LeetCode #1**

Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.

**Example:**
```
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: nums[0] + nums[1] = 2 + 7 = 9
```

**Constraints:**
- Each input has exactly one solution
- Cannot use the same element twice

**Your Task:**
1. First, solve with brute force O(n²) approach
2. Then optimize using hash map O(n) approach
3. Explain why hash map approach is better

**Template:**
```python
def two_sum(nums, target):
    """
    Find two indices that sum to target
    Time: O(n), Space: O(n)
    """
    # Your code here
    pass

# Test cases
assert two_sum([2,7,11,15], 9) == [0,1]
assert two_sum([3,2,4], 6) == [1,2]
assert two_sum([3,3], 6) == [0,1]
```

**Learning Focus:**
- Hash map for O(1) lookups
- Trade space for time complexity
- Handling edge cases

---

## Problem 2: Valid Palindrome (Easy) ⭐⭐
**LeetCode #125**

Check if a string is a palindrome, considering only alphanumeric characters and ignoring case.

**Example:**
```
Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome
```

**Your Task:**
1. Use two pointers approach
2. Handle non-alphanumeric characters properly
3. Make it case-insensitive

**Template:**
```python
def is_palindrome(s):
    """
    Check if string is palindrome (alphanumeric only, case insensitive)
    Time: O(n), Space: O(1)
    """
    # Your code here
    pass

# Test cases
assert is_palindrome("A man, a plan, a canal: Panama") == True
assert is_palindrome("race a car") == False
assert is_palindrome("") == True
```

**Learning Focus:**
- Two pointers technique
- String character filtering
- Case handling

---

## Problem 3: Remove Duplicates from Sorted Array (Easy) ⭐⭐
**LeetCode #26**

Remove duplicates from a sorted array in-place. Return the new length.

**Example:**
```
Input: nums = [1,1,2]
Output: 2, nums = [1,2,_]
Explanation: First 2 elements should be [1,2]
```

**Your Task:**
1. Modify array in-place (don't create new array)
2. Return the new length
3. Use two pointers approach

**Template:**
```python
def remove_duplicates(nums):
    """
    Remove duplicates in-place from sorted array
    Time: O(n), Space: O(1)
    """
    # Your code here
    pass

# Test cases
nums1 = [1,1,2]
length1 = remove_duplicates(nums1)
assert length1 == 2
assert nums1[:length1] == [1,2]

nums2 = [0,0,1,1,1,2,2,3,3,4]
length2 = remove_duplicates(nums2)
assert length2 == 5
assert nums2[:length2] == [0,1,2,3,4]
```

**Learning Focus:**
- In-place array modification
- Two pointers with different speeds
- Understanding array memory layout

---

## Problem 4: Best Time to Buy and Sell Stock (Easy) ⭐⭐⭐
**LeetCode #121**

Find the maximum profit from buying and selling stock once.

**Example:**
```
Input: prices = [7,1,5,3,6,4]
Output: 5
Explanation: Buy on day 2 (price=1), sell on day 5 (price=6), profit = 6-1 = 5
```

**Your Task:**
1. Track minimum price seen so far
2. Calculate maximum profit at each step
3. Return maximum profit found

**Template:**
```python
def max_profit(prices):
    """
    Find maximum profit from single buy/sell
    Time: O(n), Space: O(1)
    """
    # Your code here
    pass

# Test cases
assert max_profit([7,1,5,3,6,4]) == 5
assert max_profit([7,6,4,3,1]) == 0
assert max_profit([1,2,3,4,5]) == 4
```

**Learning Focus:**
- Single pass optimization
- Tracking state variables
- Greedy algorithm thinking

---

## Problem 5: Plus One (Easy) ⭐
**LeetCode #66**

Add one to a number represented as array of digits.

**Example:**
```
Input: digits = [1,2,3]
Output: [1,2,4]
Explanation: 123 + 1 = 124
```

**Edge Case:**
```
Input: digits = [9,9,9]
Output: [1,0,0,0]
Explanation: 999 + 1 = 1000
```

**Your Task:**
1. Handle carry propagation
2. Handle the edge case where all digits are 9
3. Work from right to left

**Template:**
```python
def plus_one(digits):
    """
    Add one to number represented as array
    Time: O(n), Space: O(1) or O(n) for edge case
    """
    # Your code here
    pass

# Test cases
assert plus_one([1,2,3]) == [1,2,4]
assert plus_one([4,3,2,1]) == [4,3,2,2]
assert plus_one([9]) == [1,0]
assert plus_one([9,9,9]) == [1,0,0,0]
```

**Learning Focus:**
- Carry handling in numeric operations
- Array resizing edge cases
- Working backwards through arrays

---

## 🎯 Problem-Solving Framework

For each problem, follow this structure:

### 1. Understand (5 minutes)
- Read problem statement carefully
- Identify inputs and outputs
- List constraints and edge cases
- Draw examples on paper

### 2. Plan (10 minutes)
- Think of brute force approach first
- Consider optimization opportunities
- Choose appropriate data structures
- Estimate time/space complexity

### 3. Code (15 minutes)
- Start with function signature
- Write main logic first
- Add edge case handling
- Use meaningful variable names

### 4. Test (5 minutes)
- Run provided test cases
- Think of additional edge cases
- Trace through your algorithm
- Fix any bugs found

### 5. Optimize (5 minutes)
- Review for potential improvements
- Consider different approaches
- Analyze time/space complexity
- Can you do better?

---

## 📊 Self-Assessment Checklist

After solving each problem, check:

**Understanding:**
- [ ] Can I explain the problem in my own words?
- [ ] Do I understand all the constraints?
- [ ] Have I identified all edge cases?

**Implementation:**
- [ ] Does my solution handle all test cases?
- [ ] Is my code clean and readable?
- [ ] Have I chosen appropriate variable names?

**Complexity:**
- [ ] What's the time complexity? Can I explain why?
- [ ] What's the space complexity? Can I explain why?
- [ ] Is this the most optimal solution?

**Communication:**
- [ ] Can I explain my approach to someone else?
- [ ] Can I walk through the algorithm step by step?
- [ ] Can I justify my design choices?

---

## 🚀 Next Steps

**After completing Day 1:**
- Review your solutions and understand why they work
- Time yourself solving similar problems
- Practice explaining your solutions out loud
- Move to Day 2 problems with string manipulation focus

**If you're struggling:**
- Focus on understanding rather than speed
- Review the code examples in the main README
- Practice drawing the algorithm on paper first
- Don't hesitate to look at hints, but try to implement yourself

**If problems feel too easy:**
- Challenge yourself to find multiple solutions
- Focus on optimizing time/space complexity
- Try to solve without looking at templates
- Help others understand the concepts

Remember: Building strong fundamentals takes time. Focus on understanding over speed! 🌱
