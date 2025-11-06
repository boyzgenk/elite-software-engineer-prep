// Package dp implements comprehensive Knapsack problem variations
// Essential DP patterns for FAANG L3/L4 interviews with production-quality implementations.
//
// Covers all major knapsack types:
// - 0/1 Knapsack (classic discrete optimization)
// - Unbounded Knapsack (infinite items available)
// - Multi-dimensional Knapsack (multiple constraints)
// - Fractional Knapsack (greedy approach for comparison)
//
// Each implementation includes detailed complexity analysis and optimization techniques
package dp

import (
	"errors"
	"sort"
)

// Item represents an item with weight, value, and optional count
type Item struct {
	Weight int
	Value  int
	Count  int // For bounded knapsack, -1 means unlimited
}

// KnapsackResult contains the solution to a knapsack problem
type KnapsackResult struct {
	MaxValue      int
	SelectedItems []int // Indices of selected items
	ItemCounts    []int // How many of each item was selected
	TotalWeight   int
	Efficiency    float64 // Value per unit weight
}

// ============= 0/1 KNAPSACK (CLASSIC) =============

// Knapsack01 solves the classic 0/1 knapsack problem
// Each item can be taken at most once
// Time Complexity: O(n * W) where n = items, W = capacity
// Space Complexity: O(n * W) for DP table
func Knapsack01(items []Item, capacity int) (*KnapsackResult, error) {
	if capacity < 0 {
		return nil, errors.New("capacity must be non-negative")
	}

	if len(items) == 0 {
		return &KnapsackResult{}, nil
	}

	n := len(items)

	// dp[i][w] = maximum value using items 0..i-1 with weight limit w
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	// Fill DP table
	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			item := items[i-1]

			// Option 1: Don't take current item
			dp[i][w] = dp[i-1][w]

			// Option 2: Take current item (if it fits)
			if item.Weight <= w {
				takeValue := dp[i-1][w-item.Weight] + item.Value
				if takeValue > dp[i][w] {
					dp[i][w] = takeValue
				}
			}
		}
	}

	// Backtrack to find selected items
	selectedItems := []int{}
	itemCounts := make([]int, n)
	totalWeight := 0

	i, w := n, capacity
	for i > 0 && w > 0 {
		// If value came from including item i-1
		if dp[i][w] != dp[i-1][w] {
			selectedItems = append(selectedItems, i-1)
			itemCounts[i-1] = 1
			totalWeight += items[i-1].Weight
			w -= items[i-1].Weight
		}
		i--
	}

	// Reverse selectedItems to get correct order
	for left, right := 0, len(selectedItems)-1; left < right; left, right = left+1, right-1 {
		selectedItems[left], selectedItems[right] = selectedItems[right], selectedItems[left]
	}

	efficiency := 0.0
	if totalWeight > 0 {
		efficiency = float64(dp[n][capacity]) / float64(totalWeight)
	}

	return &KnapsackResult{
		MaxValue:      dp[n][capacity],
		SelectedItems: selectedItems,
		ItemCounts:    itemCounts,
		TotalWeight:   totalWeight,
		Efficiency:    efficiency,
	}, nil
}

// Knapsack01Optimized solves 0/1 knapsack with space optimization
// Uses only two rows instead of full DP table
// Time Complexity: O(n * W)
// Space Complexity: O(W) - space optimized
func Knapsack01Optimized(items []Item, capacity int) (int, error) {
	if capacity < 0 {
		return 0, errors.New("capacity must be non-negative")
	}

	if len(items) == 0 {
		return 0, nil
	}

	// Use only one row, process from right to left to avoid overwriting
	dp := make([]int, capacity+1)

	for _, item := range items {
		// Process from right to left to avoid using updated values
		for w := capacity; w >= item.Weight; w-- {
			takeValue := dp[w-item.Weight] + item.Value
			if takeValue > dp[w] {
				dp[w] = takeValue
			}
		}
	}

	return dp[capacity], nil
}

// ============= UNBOUNDED KNAPSACK =============

// KnapsackUnbounded solves unbounded knapsack problem
// Each item can be taken unlimited number of times
// Time Complexity: O(n * W)
// Space Complexity: O(W)
func KnapsackUnbounded(items []Item, capacity int) (*KnapsackResult, error) {
	if capacity < 0 {
		return nil, errors.New("capacity must be non-negative")
	}

	if len(items) == 0 {
		return &KnapsackResult{}, nil
	}

	// dp[w] = maximum value with weight limit w
	dp := make([]int, capacity+1)
	parent := make([]int, capacity+1) // To track which item was used

	// Initialize parent array
	for i := range parent {
		parent[i] = -1
	}

	// Fill DP table
	for w := 1; w <= capacity; w++ {
		for i, item := range items {
			if item.Weight <= w {
				takeValue := dp[w-item.Weight] + item.Value
				if takeValue > dp[w] {
					dp[w] = takeValue
					parent[w] = i
				}
			}
		}
	}

	// Backtrack to find solution
	selectedItems := []int{}
	itemCounts := make([]int, len(items))
	totalWeight := 0

	w := capacity
	for w > 0 && parent[w] != -1 {
		itemIndex := parent[w]
		selectedItems = append(selectedItems, itemIndex)
		itemCounts[itemIndex]++
		totalWeight += items[itemIndex].Weight
		w -= items[itemIndex].Weight
	}

	efficiency := 0.0
	if totalWeight > 0 {
		efficiency = float64(dp[capacity]) / float64(totalWeight)
	}

	return &KnapsackResult{
		MaxValue:      dp[capacity],
		SelectedItems: selectedItems,
		ItemCounts:    itemCounts,
		TotalWeight:   totalWeight,
		Efficiency:    efficiency,
	}, nil
}

// ============= BOUNDED KNAPSACK =============

// KnapsackBounded solves bounded knapsack problem
// Each item has a limited count available
// Time Complexity: O(n * W * max_count)
// Space Complexity: O(W)
func KnapsackBounded(items []Item, capacity int) (*KnapsackResult, error) {
	if capacity < 0 {
		return nil, errors.New("capacity must be non-negative")
	}

	if len(items) == 0 {
		return &KnapsackResult{}, nil
	}

	// Convert bounded to 0/1 by expanding items
	expandedItems := []Item{}
	originalIndex := []int{}

	for i, item := range items {
		count := item.Count
		if count <= 0 {
			count = capacity/item.Weight + 1 // Treat as unbounded
		}

		for j := 0; j < count; j++ {
			expandedItems = append(expandedItems, Item{
				Weight: item.Weight,
				Value:  item.Value,
			})
			originalIndex = append(originalIndex, i)
		}
	}

	// Solve as 0/1 knapsack
	result, err := Knapsack01(expandedItems, capacity)
	if err != nil {
		return nil, err
	}

	// Convert back to original item indices and counts
	itemCounts := make([]int, len(items))
	originalSelected := []int{}

	for _, expandedIdx := range result.SelectedItems {
		origIdx := originalIndex[expandedIdx]
		itemCounts[origIdx]++

		// Add to selected list only once per original item
		found := false
		for _, selected := range originalSelected {
			if selected == origIdx {
				found = true
				break
			}
		}
		if !found {
			originalSelected = append(originalSelected, origIdx)
		}
	}

	return &KnapsackResult{
		MaxValue:      result.MaxValue,
		SelectedItems: originalSelected,
		ItemCounts:    itemCounts,
		TotalWeight:   result.TotalWeight,
		Efficiency:    result.Efficiency,
	}, nil
}

// ============= FRACTIONAL KNAPSACK (GREEDY) =============

// ItemEfficiency pairs item with its value-to-weight ratio
type ItemEfficiency struct {
	Index      int
	Efficiency float64
	Item       Item
}

// KnapsackFractional solves fractional knapsack using greedy approach
// Items can be taken in fractions (not a DP problem, but useful for comparison)
// Time Complexity: O(n log n) due to sorting
// Space Complexity: O(n)
func KnapsackFractional(items []Item, capacity int) (*KnapsackResult, error) {
	if capacity < 0 {
		return nil, errors.New("capacity must be non-negative")
	}

	if len(items) == 0 {
		return &KnapsackResult{}, nil
	}

	// Create efficiency array and sort by value-to-weight ratio
	efficiencies := make([]ItemEfficiency, len(items))
	for i, item := range items {
		if item.Weight <= 0 {
			return nil, errors.New("item weight must be positive")
		}
		efficiencies[i] = ItemEfficiency{
			Index:      i,
			Efficiency: float64(item.Value) / float64(item.Weight),
			Item:       item,
		}
	}

	// Sort by efficiency (descending)
	sort.Slice(efficiencies, func(i, j int) bool {
		return efficiencies[i].Efficiency > efficiencies[j].Efficiency
	})

	maxValue := 0.0
	totalWeight := 0
	selectedItems := []int{}
	fractions := make([]float64, len(items))

	remainingCapacity := capacity

	for _, eff := range efficiencies {
		if remainingCapacity == 0 {
			break
		}

		item := eff.Item
		if item.Weight <= remainingCapacity {
			// Take whole item
			maxValue += float64(item.Value)
			totalWeight += item.Weight
			remainingCapacity -= item.Weight
			selectedItems = append(selectedItems, eff.Index)
			fractions[eff.Index] = 1.0
		} else {
			// Take fraction of item
			fraction := float64(remainingCapacity) / float64(item.Weight)
			maxValue += fraction * float64(item.Value)
			totalWeight += remainingCapacity
			selectedItems = append(selectedItems, eff.Index)
			fractions[eff.Index] = fraction
			remainingCapacity = 0
		}
	}

	efficiency := 0.0
	if totalWeight > 0 {
		efficiency = maxValue / float64(totalWeight)
	}

	return &KnapsackResult{
		MaxValue:      int(maxValue),
		SelectedItems: selectedItems,
		TotalWeight:   totalWeight,
		Efficiency:    efficiency,
	}, nil
}

// ============= MULTI-DIMENSIONAL KNAPSACK =============

// MultiConstraint represents multiple constraints (weight, volume, etc.)
type MultiConstraint struct {
	Limits []int // limits[i] = limit for constraint i
}

// MultiItem represents item with multiple constraint values
type MultiItem struct {
	Values []int // values[i] = consumption for constraint i
	Profit int   // Profit/value of the item
}

// KnapsackMultiDimensional solves knapsack with multiple constraints
// Time Complexity: O(n * ∏(limits)) - product of all constraint limits
// Space Complexity: O(∏(limits))
func KnapsackMultiDimensional(items []MultiItem, constraints MultiConstraint) (int, error) {
	if len(constraints.Limits) == 0 {
		return 0, errors.New("at least one constraint required")
	}

	// Validate all items have same number of constraint values
	numConstraints := len(constraints.Limits)
	for _, item := range items {
		if len(item.Values) != numConstraints {
			return 0, errors.New("item constraint count mismatch")
		}
		for _, value := range item.Values {
			if value < 0 {
				return 0, errors.New("item constraint values must be non-negative")
			}
		}
	}

	// For simplicity, implement 2D case (most common in interviews)
	if numConstraints == 2 {
		return knapsack2D(items, constraints.Limits[0], constraints.Limits[1])
	}

	// For higher dimensions, would need recursive approach or more complex indexing
	return 0, errors.New("multi-dimensional knapsack with >2 constraints not implemented")
}

// knapsack2D solves 2-constraint knapsack problem
func knapsack2D(items []MultiItem, limit1, limit2 int) (int, error) {
	// dp[w1][w2] = maximum profit with constraints w1, w2
	dp := make([][]int, limit1+1)
	for i := range dp {
		dp[i] = make([]int, limit2+1)
	}

	for _, item := range items {
		weight1, weight2 := item.Values[0], item.Values[1]
		profit := item.Profit

		// Process from back to avoid using updated values (0/1 knapsack)
		for w1 := limit1; w1 >= weight1; w1-- {
			for w2 := limit2; w2 >= weight2; w2-- {
				takeValue := dp[w1-weight1][w2-weight2] + profit
				if takeValue > dp[w1][w2] {
					dp[w1][w2] = takeValue
				}
			}
		}
	}

	return dp[limit1][limit2], nil
}

// ============= KNAPSACK VARIANTS FOR INTERVIEW PRACTICE =============

// CoinChange calculates minimum coins needed to make amount (unbounded knapsack variant)
// LeetCode #322 - Classic DP problem using knapsack principles
func CoinChange(coins []int, amount int) (int, error) {
	if amount < 0 {
		return -1, errors.New("amount must be non-negative")
	}

	if amount == 0 {
		return 0, nil
	}

	// dp[i] = minimum coins needed to make amount i
	dp := make([]int, amount+1)

	// Initialize with "impossible" value
	for i := 1; i <= amount; i++ {
		dp[i] = amount + 1 // Max possible coins + 1
	}

	for amt := 1; amt <= amount; amt++ {
		for _, coin := range coins {
			if coin <= amt {
				dp[amt] = minInt(dp[amt], dp[amt-coin]+1)
			}
		}
	}

	if dp[amount] > amount {
		return -1, nil // Impossible to make the amount
	}

	return dp[amount], nil
}

// CoinChangeWays counts number of ways to make amount (unbounded knapsack variant)
// LeetCode #518 - Coin Change II
func CoinChangeWays(coins []int, amount int) (int, error) {
	if amount < 0 {
		return 0, errors.New("amount must be non-negative")
	}

	// dp[i] = number of ways to make amount i
	dp := make([]int, amount+1)
	dp[0] = 1 // One way to make amount 0 (use no coins)

	// Process each coin type
	for _, coin := range coins {
		// Update dp array for current coin
		for amt := coin; amt <= amount; amt++ {
			dp[amt] += dp[amt-coin]
		}
	}

	return dp[amount], nil
}

// PartitionEqualSubsetSum checks if array can be partitioned into two equal sum subsets
// LeetCode #416 - Subset sum variant of knapsack
func PartitionEqualSubsetSum(nums []int) (bool, error) {
	if len(nums) < 2 {
		return false, nil
	}

	totalSum := 0
	for _, num := range nums {
		if num < 0 {
			return false, errors.New("numbers must be non-negative")
		}
		totalSum += num
	}

	// If total sum is odd, can't partition equally
	if totalSum%2 != 0 {
		return false, nil
	}

	target := totalSum / 2

	// dp[i] = true if sum i can be achieved
	dp := make([]bool, target+1)
	dp[0] = true // Sum 0 is always achievable (empty subset)

	for _, num := range nums {
		// Process from right to left (0/1 knapsack pattern)
		for j := target; j >= num; j-- {
			dp[j] = dp[j] || dp[j-num]
		}
	}

	return dp[target], nil
}

// ============= UTILITY FUNCTIONS =============

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============= KNAPSACK COMPARISON UTILITIES =============

// CompareKnapsackMethods compares different knapsack solving approaches
func CompareKnapsackMethods(items []Item, capacity int) (*KnapsackComparison, error) {
	result01, err := Knapsack01(items, capacity)
	if err != nil {
		return nil, err
	}

	resultUnbounded, err := KnapsackUnbounded(items, capacity)
	if err != nil {
		return nil, err
	}

	resultFractional, err := KnapsackFractional(items, capacity)
	if err != nil {
		return nil, err
	}

	return &KnapsackComparison{
		ZeroOne:    result01,
		Unbounded:  resultUnbounded,
		Fractional: resultFractional,
	}, nil
}

// KnapsackComparison holds results from different knapsack approaches
type KnapsackComparison struct {
	ZeroOne    *KnapsackResult
	Unbounded  *KnapsackResult
	Fractional *KnapsackResult
}
