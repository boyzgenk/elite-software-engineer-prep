# Module 1 Assessment Tools
## CS Fundamentals - Junior Software Engineer Interview Prep

### 🎯 Assessment Overview
These tools help measure your progress against the Module 1 learning objectives:
- **Problem-solving in 20-30 minutes** with working solutions
- **Clear explanation** of approach and complexity
- **Code quality** and best practices
- **Systematic thinking** and debugging skills

---

## ⏱️ Timed Coding Assessment Tool

### How to Use
1. Choose a problem from the assessment bank
2. Set a timer for the target time (usually 25-30 minutes)
3. Run the Go assessment program
4. Implement your solution
5. Review your performance against the rubric

### Running Assessments

```bash
# Navigate to assessment tools directory
cd /path/to/junior/modules/01-fundamentals/assessment-tools

# Run timed assessment (interactive)
go run timed_assessment.go

# Run specific problem assessment
go run timed_assessment.go --problem "two-sum" --time 25

# Run full battery of assessments
go run assessment_battery.go
```

---

## 📊 Self-Assessment Rubrics

### Problem-Solving Speed (25 points)
- **Excellent (23-25 points):** Solves problem in ≤20 minutes with working solution
- **Good (20-22 points):** Solves problem in 21-25 minutes with working solution  
- **Satisfactory (15-19 points):** Solves problem in 26-30 minutes with working solution
- **Needs Work (<15 points):** Takes >30 minutes or solution doesn't work

### Code Quality (25 points)
- **Excellent (23-25 points):** Clean, readable code with proper naming and structure
- **Good (20-22 points):** Generally clean code with minor style issues
- **Satisfactory (15-19 points):** Working code but inconsistent style/naming
- **Needs Work (<15 points):** Messy, hard-to-read code or poor structure

### Problem Understanding (25 points)
- **Excellent (23-25 points):** Correctly identifies all edge cases and constraints
- **Good (20-22 points):** Understands main problem with most edge cases
- **Satisfactory (15-19 points):** Basic understanding but misses some edge cases
- **Needs Work (<15 points):** Misunderstands problem requirements

### Technical Communication (25 points)
- **Excellent (23-25 points):** Clearly explains approach, complexity, and trade-offs
- **Good (20-22 points):** Good explanation with minor gaps
- **Satisfactory (15-19 points):** Basic explanation but lacks detail
- **Needs Work (<15 points):** Unclear or incorrect explanations

### Overall Readiness Scale
- **90-100 points:** Ready for junior interviews
- **80-89 points:** Almost ready, minor improvements needed
- **70-79 points:** Developing well, need more practice
- **60-69 points:** Basic skills present, significant practice needed
- **<60 points:** Needs fundamental skill building

---

## 🎯 Assessment Problem Bank

### Week 1 Assessment Problems (Data Structures)

#### Easy Problems (Target: 15-20 minutes)
1. **Two Sum** - Find two numbers that add up to target
2. **Valid Parentheses** - Check if parentheses are properly balanced
3. **Merge Two Sorted Lists** - Combine two sorted linked lists
4. **Remove Duplicates from Sorted Array** - In-place duplicate removal
5. **Maximum Subarray** - Find contiguous subarray with largest sum

#### Medium Problems (Target: 25-30 minutes)
1. **Add Two Numbers** - Add numbers represented as linked lists
2. **Longest Substring Without Repeating Characters** - Sliding window technique
3. **Container With Most Water** - Two pointers approach
4. **3Sum** - Find all unique triplets that sum to zero
5. **Binary Tree Level Order Traversal** - Queue-based BFS

### Week 2 Assessment Problems (Advanced Data Structures)

#### Easy Problems (Target: 15-20 minutes)
1. **Implement Stack using Queues** - Stack operations with queue data structure
2. **Valid Binary Search Tree** - Validate BST properties
3. **Symmetric Tree** - Check if binary tree is mirror of itself
4. **Same Tree** - Compare two binary trees for equality
5. **Path Sum** - Check if root-to-leaf path sum equals target

#### Medium Problems (Target: 25-30 minutes)
1. **Binary Tree Zigzag Level Order Traversal** - Alternating left-right traversal
2. **Construct Binary Tree from Preorder and Inorder** - Tree reconstruction
3. **Validate Binary Search Tree** - In-order traversal approach
4. **Kth Smallest Element in BST** - In-order traversal with counting
5. **Lowest Common Ancestor of Binary Tree** - Recursive approach

### Week 3 Assessment Problems (Algorithms)

#### Easy Problems (Target: 15-20 minutes)
1. **Binary Search** - Standard implementation
2. **Search Insert Position** - Binary search variation
3. **First Bad Version** - Binary search on version space
4. **Sqrt(x)** - Binary search for square root
5. **Guess Number Higher or Lower** - Interactive binary search

#### Medium Problems (Target: 25-30 minutes)
1. **Search in Rotated Sorted Array** - Modified binary search
2. **Find Peak Element** - Binary search on peak condition
3. **Search a 2D Matrix** - 2D binary search
4. **Find First and Last Position** - Binary search boundaries
5. **Search in Rotated Sorted Array II** - With duplicates

### Week 4 Assessment Problems (Integration)

#### Easy Problems (Target: 15-20 minutes)
1. **Merge Sorted Array** - In-place merge
2. **Intersection of Two Arrays** - Set intersection
3. **Plus One** - Array manipulation
4. **Move Zeroes** - Two pointers technique
5. **Remove Element** - In-place removal

#### Medium Problems (Target: 25-30 minutes)
1. **Sort Colors** - Dutch national flag problem
2. **Top K Frequent Elements** - Heap or quickselect
3. **Kth Largest Element in Array** - Quickselect algorithm
4. **Find All Anagrams in String** - Sliding window with hash map
5. **Subarray Sum Equals K** - Prefix sum with hash map

---

## 🔧 Assessment Configuration

### Difficulty Progression
- **Week 1:** Focus on easy problems, build confidence
- **Week 2:** Mix of easy/medium, introduce complexity
- **Week 3:** Primarily medium problems, algorithm focus  
- **Week 4:** Integration problems, interview simulation

### Time Targets by Experience Level
- **Complete Beginner:** Easy (25-30 min), Medium (35-40 min)
- **Some Experience:** Easy (20-25 min), Medium (30-35 min)
- **Target Performance:** Easy (15-20 min), Medium (25-30 min)

### Assessment Frequency
- **Daily Practice:** 1-2 easy problems with timing
- **Weekly Assessment:** 1 medium problem under time pressure
- **Milestone Assessment:** Battery of 5 problems (mix of easy/medium)
- **Final Assessment:** Mock interview simulation with 2-3 problems

---

## 📈 Progress Tracking

### Performance Metrics to Track
1. **Average completion time** by problem difficulty
2. **Success rate** (working solution within time limit)
3. **Code quality score** (based on rubric)
4. **Explanation clarity** (self-assessed or peer-reviewed)
5. **Edge case coverage** (percentage of test cases passed)

### Weekly Progress Report Template
```
Week X Assessment Report
========================
Date: [Date]
Problems Attempted: [Number]
Problems Completed Successfully: [Number]
Average Time (Easy): [Minutes]
Average Time (Medium): [Minutes]

Strengths:
- [Strength 1]
- [Strength 2]

Areas for Improvement:
- [Area 1]
- [Area 2]

Action Plan for Next Week:
- [Action 1]
- [Action 2]

Overall Readiness Score: [X/100]
```

### Improvement Tracking
- **Speed:** Track time reduction over repeated problem types  
- **Accuracy:** Track percentage of first-attempt solutions that work
- **Consistency:** Track standard deviation in performance
- **Growth:** Track ability to handle increasingly difficult problems

---

## 🎮 Mock Interview Simulations

### Structure of Mock Interviews
1. **Introduction (2 minutes):** Brief personal introduction
2. **Problem Presentation (3 minutes):** Read and understand problem
3. **Solution Development (20-25 minutes):** Code and explain approach
4. **Testing and Debugging (5 minutes):** Verify solution works
5. **Follow-up Questions (5 minutes):** Extensions and optimizations

### Mock Interview Problems by Week

#### Week 1 Mock Interview
- **Problem:** Two Sum or Valid Parentheses
- **Focus:** Basic problem-solving and code clarity
- **Follow-up:** "What if the array was sorted?" or "How would you handle multiple test cases?"

#### Week 2 Mock Interview  
- **Problem:** Binary Tree Level Order Traversal
- **Focus:** Tree traversal and queue usage
- **Follow-up:** "How would you print it in reverse level order?" or "What about zigzag traversal?"

#### Week 3 Mock Interview
- **Problem:** Search in Rotated Sorted Array
- **Focus:** Binary search variations and edge cases
- **Follow-up:** "What if there are duplicates?" or "How would you find the rotation point?"

#### Week 4 Mock Interview
- **Problem:** Sort Colors (Dutch National Flag)
- **Focus:** Advanced problem-solving and optimization
- **Follow-up:** "What if there were 4 colors?" or "How would you verify your solution?"

### Self-Evaluation Questions Post-Interview
1. Did I solve the problem within the time limit?
2. Was my code clean and readable?
3. Did I handle edge cases appropriately?
4. Could I explain my approach clearly?
5. How did I handle questions and follow-ups?
6. What would I do differently next time?

---

## 📋 Assessment Checklist

### Pre-Assessment Setup
- [ ] Environment ready (Go installed, editor configured)
- [ ] Timer available and working
- [ ] Quiet space with no distractions
- [ ] Problem statement clearly understood
- [ ] Test cases identified

### During Assessment
- [ ] Start timer immediately after reading problem
- [ ] Write out approach before coding
- [ ] Test with provided examples
- [ ] Check edge cases
- [ ] Refactor for clarity if time allows

### Post-Assessment Review
- [ ] Record completion time
- [ ] Evaluate code quality using rubric
- [ ] Test with additional edge cases
- [ ] Document lessons learned
- [ ] Plan improvements for next assessment

### Weekly Assessment Goals
- **Week 1:** Complete 3+ easy problems successfully
- **Week 2:** Complete 2+ easy and 1+ medium problem
- **Week 3:** Complete 2+ medium problems consistently  
- **Week 4:** Pass mock interview with good performance

**Success Criteria:** Consistently solve easy problems in <20 minutes and medium problems in <30 minutes with clean, working code.