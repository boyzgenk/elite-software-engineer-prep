// Package dp implements fundamental Dynamic Programming patterns
// for FAANG L3/L4 interview preparation with production-quality Go implementations.
//
// This module covers core DP concepts:
// - Fibonacci sequence (1D DP foundation)
// - Climbing stairs (basic state transition)
// - House robber (decision-based DP)
// - Grid path problems (2D DP introduction)
//
// Time/Space Complexity analysis provided for each implementation
package dp

import (
	"errors"
	"fmt"
)

// ============= 1D DYNAMIC PROGRAMMING PATTERNS =============

// FibonacciDP calculates nth Fibonacci number using bottom-up DP
// Classic example of overlapping subproblems and optimal substructure
// Time Complexity: O(n)
// Space Complexity: O(n) for dp array
func FibonacciDP(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n must be non-negative")
	}

	if n <= 1 {
		return n, nil
	}

	// dp[i] = ith Fibonacci number
	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n], nil
}

// FibonacciOptimized calculates Fibonacci with space optimization
// Only keep track of last two values instead of entire array
// Time Complexity: O(n)
// Space Complexity: O(1) - space optimized
func FibonacciOptimized(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n must be non-negative")
	}

	if n <= 1 {
		return n, nil
	}

	prev2 := 0 // dp[i-2]
	prev1 := 1 // dp[i-1]

	for i := 2; i <= n; i++ {
		current := prev1 + prev2
		prev2 = prev1
		prev1 = current
	}

	return prev1, nil
}

// ClimbingStairs calculates number of ways to climb n stairs
// You can climb either 1 or 2 steps at a time
// LeetCode #70 - Classic DP problem
// Time Complexity: O(n)
// Space Complexity: O(1) - space optimized
func ClimbingStairs(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n must be non-negative")
	}

	if n <= 2 {
		return n, nil
	}

	// Base cases: 1 way to climb 1 stair, 2 ways to climb 2 stairs
	first := 1  // ways to climb 1 stair
	second := 2 // ways to climb 2 stairs

	// For each stair i, ways[i] = ways[i-1] + ways[i-2]
	for i := 3; i <= n; i++ {
		current := first + second
		first = second
		second = current
	}

	return second, nil
}

// ClimbingStairsVariant calculates ways to climb with k different step sizes
// Generalized version where you can take steps of size 1, 2, ..., k
// Time Complexity: O(n * k)
// Space Complexity: O(n)
func ClimbingStairsVariant(n int, k int) (int, error) {
	if n < 0 || k <= 0 {
		return 0, errors.New("n and k must be positive")
	}

	if n == 0 {
		return 1, nil
	}

	// dp[i] = number of ways to reach step i
	dp := make([]int, n+1)
	dp[0] = 1 // One way to stay at ground (take no steps)

	for i := 1; i <= n; i++ {
		// For each position i, sum ways from all possible previous positions
		for step := 1; step <= k && step <= i; step++ {
			dp[i] += dp[i-step]
		}
	}

	return dp[n], nil
}

// HouseRobber calculates maximum money that can be robbed
// Cannot rob two adjacent houses
// LeetCode #198 - Classic decision-based DP
// Time Complexity: O(n)
// Space Complexity: O(1) - space optimized
func HouseRobber(houses []int) (int, error) {
	if len(houses) == 0 {
		return 0, nil
	}

	if len(houses) == 1 {
		return houses[0], nil
	}

	// rob = max money if we rob current house
	// notRob = max money if we don't rob current house
	rob := houses[0]
	notRob := 0

	for i := 1; i < len(houses); i++ {
		currentRob := notRob + houses[i]  // Rob current house
		currentNotRob := max(rob, notRob) // Don't rob current house

		rob = currentRob
		notRob = currentNotRob
	}

	return max(rob, notRob), nil
}

// HouseRobberArray tracks the actual houses robbed (for interview follow-up)
// Returns maximum money and the indices of houses robbed
// Time Complexity: O(n)
// Space Complexity: O(n) for tracking decisions
func HouseRobberArray(houses []int) (int, []int, error) {
	if len(houses) == 0 {
		return 0, []int{}, nil
	}

	if len(houses) == 1 {
		return houses[0], []int{0}, nil
	}

	n := len(houses)

	// dp[i] = maximum money that can be robbed from houses 0 to i
	dp := make([]int, n)
	dp[0] = houses[0]
	dp[1] = max(houses[0], houses[1])

	for i := 2; i < n; i++ {
		dp[i] = max(dp[i-1], dp[i-2]+houses[i])
	}

	// Backtrack to find which houses were robbed
	maxMoney := dp[n-1]
	robbedHouses := []int{}

	i := n - 1
	for i >= 0 {
		if i == 0 {
			if dp[i] > 0 {
				robbedHouses = append([]int{i}, robbedHouses...)
			}
			break
		} else if i == 1 {
			if dp[i] == houses[i] {
				robbedHouses = append([]int{i}, robbedHouses...)
			} else {
				robbedHouses = append([]int{0}, robbedHouses...)
			}
			break
		} else {
			if dp[i] == dp[i-1] {
				// Didn't rob house i
				i--
			} else {
				// Robbed house i
				robbedHouses = append([]int{i}, robbedHouses...)
				i -= 2
			}
		}
	}

	return maxMoney, robbedHouses, nil
}

// HouseRobberCircular solves the circular house robber problem
// Houses are arranged in a circle, so first and last are adjacent
// LeetCode #213 - Extension of basic house robber
// Time Complexity: O(n)
// Space Complexity: O(1)
func HouseRobberCircular(houses []int) (int, error) {
	if len(houses) == 0 {
		return 0, nil
	}

	if len(houses) == 1 {
		return houses[0], nil
	}

	if len(houses) == 2 {
		return max(houses[0], houses[1]), nil
	}

	// Case 1: Rob houses from 0 to n-2 (can rob first house)
	robFirst, err := HouseRobber(houses[:len(houses)-1])
	if err != nil {
		return 0, err
	}

	// Case 2: Rob houses from 1 to n-1 (can rob last house)
	robLast, err := HouseRobber(houses[1:])
	if err != nil {
		return 0, err
	}

	return max(robFirst, robLast), nil
}

// CostClimbingStairs calculates minimum cost to reach top of stairs
// You can pay cost[i] to climb 1 or 2 steps from step i
// LeetCode #746 - Cost-based climbing stairs
// Time Complexity: O(n)
// Space Complexity: O(1) - space optimized
func CostClimbingStairs(cost []int) (int, error) {
	if len(cost) < 2 {
		return 0, errors.New("need at least 2 steps")
	}

	n := len(cost)

	// Can start from step 0 or step 1
	first := cost[0]  // min cost to reach step 0
	second := cost[1] // min cost to reach step 1

	// For each step, choose minimum cost path
	for i := 2; i < n; i++ {
		current := cost[i] + min(first, second)
		first = second
		second = current
	}

	// Can reach top from either last or second-to-last step
	return min(first, second), nil
}

// ============= 2D DYNAMIC PROGRAMMING INTRODUCTION =============

// UniquePaths calculates number of unique paths in m×n grid
// Can only move right or down
// LeetCode #62 - Classic 2D DP problem
// Time Complexity: O(m * n)
// Space Complexity: O(m * n)
func UniquePaths(m, n int) (int, error) {
	if m <= 0 || n <= 0 {
		return 0, errors.New("m and n must be positive")
	}

	// dp[i][j] = number of paths to reach cell (i, j)
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	// Base cases: first row and first column
	for i := 0; i < m; i++ {
		dp[i][0] = 1 // Only one way to reach cells in first column
	}
	for j := 0; j < n; j++ {
		dp[0][j] = 1 // Only one way to reach cells in first row
	}

	// Fill DP table
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = dp[i-1][j] + dp[i][j-1] // From above + from left
		}
	}

	return dp[m-1][n-1], nil
}

// UniquePathsOptimized calculates unique paths with space optimization
// Uses only one row instead of full 2D array
// Time Complexity: O(m * n)
// Space Complexity: O(n) - space optimized
func UniquePathsOptimized(m, n int) (int, error) {
	if m <= 0 || n <= 0 {
		return 0, errors.New("m and n must be positive")
	}

	// Use only one row, update in place
	dp := make([]int, n)
	for j := 0; j < n; j++ {
		dp[j] = 1 // First row is all 1s
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[j] = dp[j] + dp[j-1] // dp[j] (from above) + dp[j-1] (from left)
		}
	}

	return dp[n-1], nil
}

// UniquePathsWithObstacles handles grid with obstacles
// LeetCode #63 - Unique paths with obstacles
// Time Complexity: O(m * n)
// Space Complexity: O(m * n)
func UniquePathsWithObstacles(obstacleGrid [][]int) (int, error) {
	if len(obstacleGrid) == 0 || len(obstacleGrid[0]) == 0 {
		return 0, errors.New("grid cannot be empty")
	}

	m, n := len(obstacleGrid), len(obstacleGrid[0])

	// If start or end is blocked, no path possible
	if obstacleGrid[0][0] == 1 || obstacleGrid[m-1][n-1] == 1 {
		return 0, nil
	}

	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	// Initialize first cell
	dp[0][0] = 1

	// Fill first row
	for j := 1; j < n; j++ {
		if obstacleGrid[0][j] == 1 {
			dp[0][j] = 0
		} else {
			dp[0][j] = dp[0][j-1]
		}
	}

	// Fill first column
	for i := 1; i < m; i++ {
		if obstacleGrid[i][0] == 1 {
			dp[i][0] = 0
		} else {
			dp[i][0] = dp[i-1][0]
		}
	}

	// Fill rest of the table
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if obstacleGrid[i][j] == 1 {
				dp[i][j] = 0
			} else {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}

	return dp[m-1][n-1], nil
}

// MinPathSum finds minimum path sum from top-left to bottom-right
// LeetCode #64 - Minimum path sum in grid
// Time Complexity: O(m * n)
// Space Complexity: O(1) - modifies input grid
func MinPathSum(grid [][]int) (int, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0, errors.New("grid cannot be empty")
	}

	m, n := len(grid), len(grid[0])

	// Fill first row
	for j := 1; j < n; j++ {
		grid[0][j] += grid[0][j-1]
	}

	// Fill first column
	for i := 1; i < m; i++ {
		grid[i][0] += grid[i-1][0]
	}

	// Fill rest of the grid
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			grid[i][j] += min(grid[i-1][j], grid[i][j-1])
		}
	}

	return grid[m-1][n-1], nil
}

// ============= UTILITY FUNCTIONS =============

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ============= DP PATTERN ANALYSIS UTILITIES =============

// DPPatternType represents different DP pattern categories
type DPPatternType int

const (
	LinearDP DPPatternType = iota
	DecisionDP
	GridDP
	SequenceDP
)

// DPProblem represents metadata about a DP problem
type DPProblem struct {
	Name            string
	PatternType     DPPatternType
	TimeComplexity  string
	SpaceComplexity string
	StateTransition string
	BaseCase        string
	Description     string
}

// GetDPPatterns returns analysis of fundamental DP patterns covered
func GetDPPatterns() []DPProblem {
	return []DPProblem{
		{
			Name:            "Fibonacci",
			PatternType:     LinearDP,
			TimeComplexity:  "O(n)",
			SpaceComplexity: "O(1) optimized",
			StateTransition: "dp[i] = dp[i-1] + dp[i-2]",
			BaseCase:        "dp[0] = 0, dp[1] = 1",
			Description:     "Classic overlapping subproblems example",
		},
		{
			Name:            "Climbing Stairs",
			PatternType:     LinearDP,
			TimeComplexity:  "O(n)",
			SpaceComplexity: "O(1) optimized",
			StateTransition: "dp[i] = dp[i-1] + dp[i-2]",
			BaseCase:        "dp[1] = 1, dp[2] = 2",
			Description:     "Count number of ways to reach target",
		},
		{
			Name:            "House Robber",
			PatternType:     DecisionDP,
			TimeComplexity:  "O(n)",
			SpaceComplexity: "O(1) optimized",
			StateTransition: "dp[i] = max(dp[i-1], dp[i-2] + nums[i])",
			BaseCase:        "dp[0] = nums[0]",
			Description:     "Optimal decision with constraints",
		},
		{
			Name:            "Unique Paths",
			PatternType:     GridDP,
			TimeComplexity:  "O(m*n)",
			SpaceComplexity: "O(n) optimized",
			StateTransition: "dp[i][j] = dp[i-1][j] + dp[i][j-1]",
			BaseCase:        "dp[0][j] = 1, dp[i][0] = 1",
			Description:     "2D grid traversal counting",
		},
	}
}

// PrintDPPatterns prints analysis of DP patterns (useful for learning)
func PrintDPPatterns() {
	patterns := GetDPPatterns()

	fmt.Println("Dynamic Programming Pattern Analysis")
	fmt.Println("=====================================")

	for _, pattern := range patterns {
		fmt.Printf("\nPattern: %s\n", pattern.Name)
		fmt.Printf("Type: %v\n", pattern.PatternType)
		fmt.Printf("Time: %s\n", pattern.TimeComplexity)
		fmt.Printf("Space: %s\n", pattern.SpaceComplexity)
		fmt.Printf("Transition: %s\n", pattern.StateTransition)
		fmt.Printf("Base Case: %s\n", pattern.BaseCase)
		fmt.Printf("Description: %s\n", pattern.Description)
		fmt.Println("---")
	}
}
