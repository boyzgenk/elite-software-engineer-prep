# Mock Interview Scripts & Evaluation Rubrics
## Module 1: CS Fundamentals Assessment

### 🎯 Purpose
These scripts simulate real junior developer technical interviews, focusing on:
- Problem-solving approach and communication
- Code implementation quality
- Time management under pressure
- Handling follow-up questions and optimizations

---

## 📋 Mock Interview Format (30 minutes total)

### Standard Structure
1. **Introduction (2 minutes)** - Brief personal introduction
2. **Problem Presentation (3 minutes)** - Understand requirements
3. **Solution Development (20 minutes)** - Code and explain
4. **Testing & Review (3 minutes)** - Verify solution works
5. **Follow-up Questions (2 minutes)** - Extensions and optimizations

---

## 🎪 Week 1 Mock Interview Scripts

### Script A: Arrays & Two Pointers (Easy Level)
**Problem:** Valid Palindrome (LeetCode #125)

**Interviewer Script:**
```
"Hi! Let's start with a string problem. 

Given a string, determine if it's a valid palindrome, considering only alphanumeric characters and ignoring cases.

For example:
- Input: 'A man, a plan, a canal: Panama'
- Output: true

- Input: 'race a car'  
- Output: false

Do you understand the problem? Any questions about the requirements?"

[Wait for clarification questions - good sign!]

"Great! How would you approach this problem?"

[Listen for approach - should mention cleaning string and two pointers]

"Sounds good! Go ahead and implement it. I'll be here if you have questions."

[During coding, note:]
- Do they clean the string properly?
- Do they use two pointers correctly?
- Are they talking through their approach?
- Do they handle edge cases?

[After coding:]
"Can you trace through your solution with the first example?"

[Follow-up questions:]
"What's the time and space complexity?"
"How would you handle Unicode characters?"
"Could we solve this without creating a new cleaned string?"
```

**Expected Solution Pattern:**
```go
func isPalindrome(s string) bool {
    // Two pointers approach with in-place checking
    left, right := 0, len(s)-1
    
    for left < right {
        // Skip non-alphanumeric from left
        for left < right && !isAlphaNumeric(s[left]) {
            left++
        }
        
        // Skip non-alphanumeric from right  
        for left < right && !isAlphaNumeric(s[right]) {
            right--
        }
        
        // Compare characters (case-insensitive)
        if toLowerCase(s[left]) != toLowerCase(s[right]) {
            return false
        }
        
        left++
        right--
    }
    
    return true
}
```

**Evaluation Rubric:**
- **Problem Understanding (25%):** Asks clarifying questions, understands requirements
- **Approach (25%):** Identifies two-pointer technique, explains logic clearly
- **Implementation (25%):** Clean code, handles edge cases, proper syntax
- **Communication (25%):** Talks through solution, explains complexity

---

### Script B: Array Manipulation (Easy-Medium Level)
**Problem:** Remove Duplicates from Sorted Array (LeetCode #26)

**Interviewer Script:**
```
"Let's work on an array problem.

Given a sorted array nums, remove duplicates in-place such that each element appears only once and returns the new length. You must do this by modifying the input array in-place with O(1) extra memory.

For example:
- Input: [1,1,2]  
- Output: 2, and nums = [1,2,_]

- Input: [0,0,1,1,1,2,2,3,3,4]
- Output: 5, and nums = [0,1,2,3,4,_,_,_,_,_]

The order of elements should be maintained. What questions do you have?"

[Look for understanding of in-place requirement]

"How would you approach this?"

[Listen for two-pointer or slow/fast pointer approach]

"Great! Let's see your implementation."

[During coding, observe:]
- Do they use two pointers effectively?
- Do they handle the in-place requirement correctly?
- Are they careful with array bounds?

[After implementation:]
"Walk me through how this works with the first example."

[Follow-up questions:]
"What if the array wasn't sorted?"
"How would you extend this to remove all duplicates (allow at most k occurrences)?"
```

**Expected Solution Pattern:**
```go
func removeDuplicates(nums []int) int {
    if len(nums) <= 1 {
        return len(nums)
    }
    
    writeIndex := 1 // Position to write next unique element
    
    for readIndex := 1; readIndex < len(nums); readIndex++ {
        if nums[readIndex] != nums[readIndex-1] {
            nums[writeIndex] = nums[readIndex]
            writeIndex++
        }
    }
    
    return writeIndex
}
```

---

## 🌳 Week 2 Mock Interview Scripts

### Script C: Binary Tree Traversal (Medium Level)
**Problem:** Binary Tree Inorder Traversal (LeetCode #94)

**Interviewer Script:**
```
"Let's work with binary trees.

Given the root of a binary tree, return the inorder traversal of its nodes' values.

For example, given tree:
    1
     \
      2
     /
    3

Output: [1,3,2]

Can you implement both recursive and iterative solutions?"

[Key points to observe:]
- Do they understand inorder traversal (left, root, right)?
- Can they implement recursive solution easily?
- Do they attempt iterative solution with stack?

[Follow-up questions:]
"What's the space complexity of each approach?"
"How would you do level-order traversal?"
"What if we wanted to do this without recursion or explicit stack?"
```

---

## 🔍 Week 3 Mock Interview Scripts  

### Script D: Binary Search Application (Medium Level)
**Problem:** Search in Rotated Sorted Array (LeetCode #33)

**Interviewer Script:**
```
"Here's an interesting search problem.

You're given a rotated sorted array and a target value. Find the index of the target, or return -1 if not found.

For example:
- Input: nums = [4,5,6,7,0,1,2], target = 0
- Output: 4

- Input: nums = [4,5,6,7,0,1,2], target = 3  
- Output: -1

The array was originally sorted, then rotated at some pivot. How would you solve this efficiently?"

[Look for binary search approach]
[Key insight: at least one half is always properly sorted]

[Follow-up questions:]
"What if there were duplicates?"
"How would you find the rotation point?"
"Can you modify this for finding the minimum element?"
```

---

## 📊 Comprehensive Evaluation Rubrics

### Technical Skills Assessment (70 points)

#### Problem Understanding (15 points)
- **Excellent (13-15):** Asks clarifying questions, identifies edge cases, restates problem correctly
- **Good (10-12):** Understands main requirements with minor gaps
- **Satisfactory (7-9):** Basic understanding but misses some nuances
- **Needs Work (0-6):** Misunderstands core requirements

#### Algorithm Design (20 points)  
- **Excellent (18-20):** Optimal approach, explains time/space complexity correctly
- **Good (14-17):** Good approach with minor inefficiencies
- **Satisfactory (10-13):** Working approach but not optimal
- **Needs Work (0-9):** Poor or incorrect approach

#### Code Implementation (20 points)
- **Excellent (18-20):** Clean, bug-free code with proper edge case handling
- **Good (14-17):** Working code with minor style/efficiency issues
- **Satisfactory (10-13):** Code works but lacks polish or has minor bugs
- **Needs Work (0-9):** Significant bugs or incomplete implementation

#### Testing & Debugging (15 points)
- **Excellent (13-15):** Tests with examples, identifies and fixes bugs quickly
- **Good (10-12):** Tests code but misses some edge cases
- **Satisfactory (7-9):** Basic testing, slow to identify issues
- **Needs Work (0-6):** No testing or unable to debug

### Communication Skills Assessment (30 points)

#### Verbal Communication (15 points)
- **Excellent (13-15):** Clear explanation of approach, thinks out loud effectively
- **Good (10-12):** Generally clear communication with minor gaps
- **Satisfactory (7-9):** Basic communication but lacks detail
- **Needs Work (0-6):** Poor communication, hard to follow thought process

#### Problem-Solving Process (15 points)
- **Excellent (13-15):** Systematic approach, considers alternatives, handles roadblocks well
- **Good (10-12):** Good process with minor disorganization
- **Satisfactory (7-9):** Some structure but inconsistent
- **Needs Work (0-6):** Chaotic or no clear process

### Overall Readiness Levels

#### 90-100 points: INTERVIEW READY
- Strong technical skills across all areas
- Excellent communication and problem-solving process
- Ready for junior developer technical interviews
- **Recommendation:** Start applying, focus on behavioral prep

#### 80-89 points: ALMOST READY  
- Good technical foundation with minor gaps
- Communication needs slight improvement
- **Recommendation:** 1-2 weeks additional practice, focus on weak areas

#### 70-79 points: DEVELOPING WELL
- Basic technical skills present but inconsistent
- Communication adequate but needs development
- **Recommendation:** 2-4 weeks additional practice, focus on consistency

#### 60-69 points: NEEDS FOCUSED PRACTICE
- Technical skills emerging but significant gaps
- Communication needs substantial improvement
- **Recommendation:** 1-2 months additional practice, consider mentoring

#### Below 60 points: FOUNDATION BUILDING
- Fundamental concepts need reinforcement
- **Recommendation:** Focus on basics, consider structured course/bootcamp

---

## 🎯 Conducting Effective Mock Interviews

### For Self-Practice
1. **Set up recording:** Record yourself solving problems
2. **Use timer:** Stick to 20-25 minute coding sessions
3. **Talk aloud:** Practice explaining your thought process
4. **Review performance:** Use rubric to self-assess

### For Peer Practice
1. **Rotate roles:** Take turns being interviewer/candidate
2. **Follow scripts:** Use provided scripts as guidelines
3. **Give feedback:** Use rubrics for constructive feedback
4. **Focus on growth:** Identify specific improvement areas

### For Mentors/Instructors
1. **Create comfortable environment:** Reduce anxiety while maintaining realism
2. **Ask probing questions:** Help candidates think deeper
3. **Provide actionable feedback:** Specific suggestions for improvement
4. **Track progress:** Use rubrics consistently to measure growth

---

**🎪 Ready to ace your technical interviews? Practice these mock interviews regularly and track your progress with the rubrics!**