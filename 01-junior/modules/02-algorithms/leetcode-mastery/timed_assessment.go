// Timed Assessment Tool - Junior Software Engineer Interview Prep
// Module 1: CS Fundamentals Assessment System
// Interactive timing and evaluation tool

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Problem represents a coding problem for assessment
type Problem struct {
	ID          string
	Title       string
	Difficulty  string
	TimeLimit   int // in minutes
	Description string
	Examples    []Example
	Constraints []string
	StarterCode string
}

// Example represents test cases
type Example struct {
	Input       string
	Output      string
	Explanation string
}

// AssessmentResult stores the results of a timed assessment
type AssessmentResult struct {
	Problem       Problem
	TimeSpent     time.Duration
	Success       bool
	CodeQuality   int // 1-25 points
	Understanding int // 1-25 points
	Communication int // 1-25 points
	OverallScore  int // total points
	Notes         string
}

// =============================================================================
// PROBLEM BANK - Week 1 Data Structures Problems
// =============================================================================

func GetWeek1Problems() []Problem {
	return []Problem{
		{
			ID:         "two-sum",
			Title:      "Two Sum",
			Difficulty: "Easy",
			TimeLimit:  20,
			Description: `Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.`,
			Examples: []Example{
				{
					Input:       "nums = [2,7,11,15], target = 9",
					Output:      "[0,1]",
					Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
				},
				{
					Input:       "nums = [3,2,4], target = 6",
					Output:      "[1,2]",
					Explanation: "Because nums[1] + nums[2] == 6, we return [1, 2].",
				},
			},
			Constraints: []string{
				"2 <= nums.length <= 10^4",
				"-10^9 <= nums[i] <= 10^9",
				"-10^9 <= target <= 10^9",
				"Only one valid answer exists.",
			},
			StarterCode: `func TwoSum(nums []int, target int) []int {
    // Your implementation here
    return []int{}
}`,
		},
		{
			ID:         "valid-parentheses",
			Title:      "Valid Parentheses",
			Difficulty: "Easy",
			TimeLimit:  15,
			Description: `Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

An input string is valid if:
1. Open brackets must be closed by the same type of brackets.
2. Open brackets must be closed in the correct order.
3. Every close bracket has a corresponding open bracket of the same type.`,
			Examples: []Example{
				{
					Input:       `s = "()"`,
					Output:      "true",
					Explanation: "The parentheses are properly matched.",
				},
				{
					Input:       `s = "()[]{}"`,
					Output:      "true",
					Explanation: "All brackets are properly matched.",
				},
				{
					Input:       `s = "(]"`,
					Output:      "false",
					Explanation: "The brackets are not properly matched.",
				},
			},
			Constraints: []string{
				"1 <= s.length <= 10^4",
				"s consists of parentheses only '()[]{}'.",
			},
			StarterCode: `func ValidParentheses(s string) bool {
    // Your implementation here
    return false
}`,
		},
		{
			ID:         "merge-sorted-arrays",
			Title:      "Merge Sorted Array",
			Difficulty: "Easy",
			TimeLimit:  20,
			Description: `You are given two integer arrays nums1 and nums2, sorted in non-decreasing order, and two integers m and n, representing the number of elements in nums1 and nums2 respectively.

Merge nums1 and nums2 into a single array sorted in non-decreasing order.

The final sorted array should not be returned by the function, but instead be stored inside the array nums1. To accommodate this, nums1 has a length of m + n, where the first m elements denote the elements that should be merged, and the last n elements are set to 0 and should be ignored. nums2 has a length of n.`,
			Examples: []Example{
				{
					Input:       "nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3",
					Output:      "[1,2,2,3,5,6]",
					Explanation: "The arrays we are merging are [1,2,3] and [2,5,6]. The result is [1,2,2,3,5,6].",
				},
			},
			Constraints: []string{
				"nums1.length == m + n",
				"nums2.length == n",
				"0 <= m, n <= 200",
				"1 <= m + n <= 200",
				"-10^9 <= nums1[i], nums2[j] <= 10^9",
			},
			StarterCode: `func MergeSortedArrays(nums1 []int, m int, nums2 []int, n int) {
    // Your implementation here - modify nums1 in place
}`,
		},
	}
}

func GetWeek2Problems() []Problem {
	return []Problem{
		{
			ID:         "binary-tree-inorder",
			Title:      "Binary Tree Inorder Traversal",
			Difficulty: "Easy",
			TimeLimit:  25,
			Description: `Given the root of a binary tree, return the inorder traversal of its nodes' values.

Inorder traversal visits nodes in this order: Left -> Root -> Right`,
			Examples: []Example{
				{
					Input:       "root = [1,null,2,3]",
					Output:      "[1,3,2]",
					Explanation: "Inorder traversal: left child (none), root (1), right subtree (3,2).",
				},
			},
			Constraints: []string{
				"The number of nodes in the tree is in the range [0, 100].",
				"-100 <= Node.val <= 100",
			},
			StarterCode: `type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func InorderTraversal(root *TreeNode) []int {
    // Your implementation here
    return []int{}
}`,
		},
	}
}

func GetWeek3Problems() []Problem {
	return []Problem{
		{
			ID:         "binary-search",
			Title:      "Binary Search",
			Difficulty: "Easy",
			TimeLimit:  15,
			Description: `Given an array of integers nums which is sorted in ascending order, and an integer target, write a function to search target in nums. If target exists, then return its index. Otherwise, return -1.

You must write an algorithm with O(log n) runtime complexity.`,
			Examples: []Example{
				{
					Input:       "nums = [-1,0,3,5,9,12], target = 9",
					Output:      "4",
					Explanation: "9 exists in nums and its index is 4",
				},
				{
					Input:       "nums = [-1,0,3,5,9,12], target = 2",
					Output:      "-1",
					Explanation: "2 does not exist in nums so return -1",
				},
			},
			Constraints: []string{
				"1 <= nums.length <= 10^4",
				"-10^4 < nums[i], target < 10^4",
				"All the integers in nums are unique.",
				"nums is sorted in ascending order.",
			},
			StarterCode: `func BinarySearch(nums []int, target int) int {
    // Your implementation here
    return -1
}`,
		},
	}
}

// =============================================================================
// TIMER AND ASSESSMENT FUNCTIONS
// =============================================================================

// StartTimer creates a countdown timer for the assessment
func StartTimer(duration time.Duration) <-chan time.Time {
	return time.After(duration)
}

// DisplayProblem shows the problem details to the user
func DisplayProblem(problem Problem) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Printf("🎯 TIMED ASSESSMENT: %s (%s)\n", problem.Title, problem.Difficulty)
	fmt.Printf("⏱️  Time Limit: %d minutes\n", problem.TimeLimit)
	fmt.Printf("%s\n", strings.Repeat("=", 70))

	fmt.Printf("\n📋 PROBLEM DESCRIPTION:\n")
	fmt.Printf("%s\n", problem.Description)

	fmt.Printf("\n🔍 EXAMPLES:\n")
	for i, example := range problem.Examples {
		fmt.Printf("Example %d:\n", i+1)
		fmt.Printf("  Input:  %s\n", example.Input)
		fmt.Printf("  Output: %s\n", example.Output)
		if example.Explanation != "" {
			fmt.Printf("  Explanation: %s\n", example.Explanation)
		}
		fmt.Println()
	}

	fmt.Printf("📏 CONSTRAINTS:\n")
	for _, constraint := range problem.Constraints {
		fmt.Printf("  • %s\n", constraint)
	}

	fmt.Printf("\n💻 STARTER CODE:\n")
	fmt.Printf("```go\n%s\n```\n", problem.StarterCode)

	fmt.Printf("\n%s\n", strings.Repeat("-", 70))
	fmt.Printf("🚀 READY TO START? Press ENTER to begin the timer...\n")
	fmt.Printf("%s\n", strings.Repeat("-", 70))
}

// RunTimedAssessment manages the timed coding session
func RunTimedAssessment(problem Problem) AssessmentResult {
	DisplayProblem(problem)

	// Wait for user to press Enter
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	// Start timing
	startTime := time.Now()
	timeLimit := time.Duration(problem.TimeLimit) * time.Minute
	timer := StartTimer(timeLimit)

	fmt.Printf("\n⏰ TIMER STARTED! You have %d minutes.\n", problem.TimeLimit)
	fmt.Printf("📝 Implement your solution now...\n")
	fmt.Printf("⌨️  When finished, type 'DONE' and press Enter\n\n")

	// Monitor progress
	done := make(chan bool)
	go func() {
		for {
			scanner.Scan()
			input := strings.TrimSpace(strings.ToUpper(scanner.Text()))
			if input == "DONE" {
				done <- true
				return
			}
		}
	}()

	// Wait for completion or timeout
	select {
	case <-timer:
		fmt.Printf("\n⏰ TIME'S UP! Assessment period has ended.\n")
		return AssessmentResult{
			Problem:   problem,
			TimeSpent: timeLimit,
			Success:   false,
			Notes:     "Timed out",
		}
	case <-done:
		elapsed := time.Since(startTime)
		fmt.Printf("\n✅ COMPLETED! Time taken: %v\n", elapsed.Round(time.Second))

		// Collect self-assessment
		return CollectSelfAssessment(problem, elapsed)
	}
}

// CollectSelfAssessment gathers user feedback on their performance
func CollectSelfAssessment(problem Problem, timeSpent time.Duration) AssessmentResult {
	scanner := bufio.NewScanner(os.Stdin)
	result := AssessmentResult{
		Problem:   problem,
		TimeSpent: timeSpent,
	}

	fmt.Printf("\n📊 SELF-ASSESSMENT TIME!\n")
	fmt.Printf("%s\n", strings.Repeat("-", 50))

	// Check if solution works
	fmt.Printf("❓ Did your solution work correctly on all test cases? (y/n): ")
	scanner.Scan()
	result.Success = strings.ToLower(scanner.Text()) == "y"

	// Code quality assessment
	fmt.Printf("\n📝 CODE QUALITY (1-25 points)\n")
	fmt.Printf("Rate your code quality (clean, readable, well-structured): ")
	scanner.Scan()
	if score, err := strconv.Atoi(scanner.Text()); err == nil && score >= 1 && score <= 25 {
		result.CodeQuality = score
	}

	// Problem understanding
	fmt.Printf("\n🧠 PROBLEM UNDERSTANDING (1-25 points)\n")
	fmt.Printf("Rate your understanding (edge cases, constraints, requirements): ")
	scanner.Scan()
	if score, err := strconv.Atoi(scanner.Text()); err == nil && score >= 1 && score <= 25 {
		result.Understanding = score
	}

	// Communication assessment
	fmt.Printf("\n💬 COMMUNICATION (1-25 points)\n")
	fmt.Printf("Rate your ability to explain the approach and complexity: ")
	scanner.Scan()
	if score, err := strconv.Atoi(scanner.Text()); err == nil && score >= 1 && score <= 25 {
		result.Communication = score
	}

	// Speed assessment (1-25 points based on time)
	speedScore := calculateSpeedScore(problem.TimeLimit, int(timeSpent.Minutes()))

	// Calculate overall score
	result.OverallScore = speedScore + result.CodeQuality + result.Understanding + result.Communication

	// Additional notes
	fmt.Printf("\n📝 Any additional notes or lessons learned? ")
	scanner.Scan()
	result.Notes = scanner.Text()

	// Display results
	DisplayAssessmentResults(result, speedScore)

	return result
}

// calculateSpeedScore assigns points based on completion speed
func calculateSpeedScore(timeLimit, actualMinutes int) int {
	if actualMinutes <= int(float64(timeLimit)*0.75) { // Within 75% of time limit
		return 25 // Excellent
	} else if actualMinutes <= timeLimit { // Within time limit
		return 20 // Good
	} else if actualMinutes <= int(float64(timeLimit)*1.2) { // 20% over
		return 15 // Satisfactory
	} else {
		return 10 // Needs work
	}
}

// DisplayAssessmentResults shows the final assessment results
func DisplayAssessmentResults(result AssessmentResult, speedScore int) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("📊 ASSESSMENT RESULTS: %s\n", result.Problem.Title)
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	fmt.Printf("⏱️  Time Spent: %v (Limit: %d minutes)\n",
		result.TimeSpent.Round(time.Second), result.Problem.TimeLimit)
	fmt.Printf("✅ Solution Works: %v\n", result.Success)

	fmt.Printf("\n📋 SCORING BREAKDOWN:\n")
	fmt.Printf("  Speed Assessment:     %d/25 points\n", speedScore)
	fmt.Printf("  Code Quality:         %d/25 points\n", result.CodeQuality)
	fmt.Printf("  Problem Understanding: %d/25 points\n", result.Understanding)
	fmt.Printf("  Technical Communication: %d/25 points\n", result.Communication)
	fmt.Printf("  %s\n", strings.Repeat("-", 40))
	fmt.Printf("  TOTAL SCORE:          %d/100 points\n", result.OverallScore)

	// Provide feedback based on score
	fmt.Printf("\n🎯 PERFORMANCE ASSESSMENT:\n")
	switch {
	case result.OverallScore >= 90:
		fmt.Printf("🌟 EXCELLENT! You're ready for junior interviews!\n")
	case result.OverallScore >= 80:
		fmt.Printf("👍 GOOD! Almost ready, minor improvements needed.\n")
	case result.OverallScore >= 70:
		fmt.Printf("📈 DEVELOPING! Good progress, keep practicing.\n")
	case result.OverallScore >= 60:
		fmt.Printf("📚 LEARNING! Basic skills present, need more practice.\n")
	default:
		fmt.Printf("🎯 BUILDING! Focus on fundamentals, you'll get there!\n")
	}

	if result.Notes != "" {
		fmt.Printf("\n📝 Notes: %s\n", result.Notes)
	}

	// Improvement suggestions
	fmt.Printf("\n💡 IMPROVEMENT SUGGESTIONS:\n")
	if speedScore < 20 {
		fmt.Printf("  • Practice similar problems to improve speed\n")
	}
	if result.CodeQuality < 20 {
		fmt.Printf("  • Focus on clean code practices and proper naming\n")
	}
	if result.Understanding < 20 {
		fmt.Printf("  • Spend more time analyzing problem constraints and edge cases\n")
	}
	if result.Communication < 20 {
		fmt.Printf("  • Practice explaining your approach and complexity analysis\n")
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
}

// =============================================================================
// MAIN ASSESSMENT INTERFACE
// =============================================================================

func ShowMainMenu() {
	fmt.Printf("\n🎯 JUNIOR DEVELOPER TIMED ASSESSMENT TOOL\n")
	fmt.Printf("Module 1: CS Fundamentals Assessment\n")
	fmt.Printf("%s\n", strings.Repeat("=", 50))
	fmt.Printf("1. Week 1 Assessment (Data Structures)\n")
	fmt.Printf("2. Week 2 Assessment (Advanced Data Structures)\n")
	fmt.Printf("3. Week 3 Assessment (Algorithms)\n")
	fmt.Printf("4. Custom Problem Assessment\n")
	fmt.Printf("5. Full Assessment Battery\n")
	fmt.Printf("6. View Assessment Guidelines\n")
	fmt.Printf("0. Exit\n")
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	fmt.Printf("Choose an option (0-6): ")
}

func ShowGuidelines() {
	fmt.Printf("\n📋 ASSESSMENT GUIDELINES\n")
	fmt.Printf("%s\n", strings.Repeat("=", 50))
	fmt.Printf("🎯 PURPOSE: Measure readiness for junior developer interviews\n")
	fmt.Printf("⏱️  TARGET: Solve easy problems in 15-20 minutes\n")
	fmt.Printf("📝 FOCUS: Working solution + clear explanation\n")
	fmt.Printf("🏆 SUCCESS: 80+ points consistently\n\n")

	fmt.Printf("📊 SCORING RUBRIC (100 points total):\n")
	fmt.Printf("  • Speed (25 pts): Completion within time limits\n")
	fmt.Printf("  • Code Quality (25 pts): Clean, readable implementation\n")
	fmt.Printf("  • Understanding (25 pts): Handle edge cases correctly\n")
	fmt.Printf("  • Communication (25 pts): Explain approach clearly\n\n")

	fmt.Printf("💡 TIPS FOR SUCCESS:\n")
	fmt.Printf("  • Read problem carefully, identify edge cases\n")
	fmt.Printf("  • Plan approach before coding\n")
	fmt.Printf("  • Write clean, well-named variables\n")
	fmt.Printf("  • Test with provided examples\n")
	fmt.Printf("  • Explain complexity (time and space)\n")
	fmt.Printf("  • Practice regularly to build speed\n")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		ShowMainMenu()
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Printf("\n📚 Week 1 Problems (Data Structures):\n")
			problems := GetWeek1Problems()
			for i, problem := range problems {
				fmt.Printf("%d. %s (%s - %d min)\n", i+1, problem.Title, problem.Difficulty, problem.TimeLimit)
			}
			fmt.Printf("Choose problem (1-%d): ", len(problems))
			scanner.Scan()
			if idx, err := strconv.Atoi(scanner.Text()); err == nil && idx > 0 && idx <= len(problems) {
				RunTimedAssessment(problems[idx-1])
			}

		case "2":
			problems := GetWeek2Problems()
			if len(problems) > 0 {
				fmt.Printf("\n🌳 Week 2 Problems (Advanced Data Structures):\n")
				for i, problem := range problems {
					fmt.Printf("%d. %s (%s - %d min)\n", i+1, problem.Title, problem.Difficulty, problem.TimeLimit)
				}
				fmt.Printf("Choose problem (1-%d): ", len(problems))
				scanner.Scan()
				if idx, err := strconv.Atoi(scanner.Text()); err == nil && idx > 0 && idx <= len(problems) {
					RunTimedAssessment(problems[idx-1])
				}
			} else {
				fmt.Printf("Week 2 problems coming soon!\n")
			}

		case "3":
			problems := GetWeek3Problems()
			if len(problems) > 0 {
				fmt.Printf("\n🔍 Week 3 Problems (Algorithms):\n")
				for i, problem := range problems {
					fmt.Printf("%d. %s (%s - %d min)\n", i+1, problem.Title, problem.Difficulty, problem.TimeLimit)
				}
				fmt.Printf("Choose problem (1-%d): ", len(problems))
				scanner.Scan()
				if idx, err := strconv.Atoi(scanner.Text()); err == nil && idx > 0 && idx <= len(problems) {
					RunTimedAssessment(problems[idx-1])
				}
			} else {
				fmt.Printf("Week 3 problems coming soon!\n")
			}

		case "4":
			fmt.Printf("Custom problem assessment coming soon!\n")

		case "5":
			fmt.Printf("Full assessment battery coming soon!\n")

		case "6":
			ShowGuidelines()

		case "0":
			fmt.Printf("\n👋 Good luck with your interview preparation!\n")
			return

		default:
			fmt.Printf("Invalid choice. Please try again.\n")
		}

		fmt.Printf("\nPress Enter to continue...")
		scanner.Scan()
	}
}
