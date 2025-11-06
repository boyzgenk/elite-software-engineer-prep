// Package theory implements mathematical foundations of complexity analysis
// Essential algorithms and concepts for FAANG L3+ interview complexity questions
package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// 🎯 BIG-O NOTATION - Formal Mathematical Definitions
// Understanding asymptotic behavior and mathematical rigor

// ComplexityClass represents different asymptotic complexity classes
type ComplexityClass string

const (
	O1     ComplexityClass = "O(1)"       // Constant
	OLogN  ComplexityClass = "O(log n)"   // Logarithmic
	ON     ComplexityClass = "O(n)"       // Linear
	ONLogN ComplexityClass = "O(n log n)" // Linearithmic
	ON2    ComplexityClass = "O(n²)"      // Quadratic
	ON3    ComplexityClass = "O(n³)"      // Cubic
	O2N    ComplexityClass = "O(2^n)"     // Exponential
	ONFact ComplexityClass = "O(n!)"      // Factorial
)

// AlgorithmAnalysis represents formal complexity analysis
type AlgorithmAnalysis struct {
	Name            string
	TimeComplexity  ComplexityClass
	SpaceComplexity ComplexityClass
	BestCase        ComplexityClass
	AverageCase     ComplexityClass
	WorstCase       ComplexityClass
	Explanation     string
}

// 📐 MASTER THEOREM IMPLEMENTATION
// For analyzing divide-and-conquer recurrences: T(n) = aT(n/b) + f(n)

// MasterTheoremCase represents the three cases of Master Theorem
type MasterTheoremCase int

const (
	Case1 MasterTheoremCase = iota // f(n) = O(n^c) where c < log_b(a)
	Case2                          // f(n) = Θ(n^c * log^k(n)) where c = log_b(a)
	Case3                          // f(n) = Ω(n^c) where c > log_b(a)
)

// MasterTheoremInput represents parameters for Master Theorem analysis
type MasterTheoremInput struct {
	A int     // Number of subproblems
	B int     // Factor by which problem size is reduced
	C float64 // Exponent in f(n) = O(n^c)
	K int     // Logarithmic factor in Case 2
}

// AnalyzeMasterTheorem applies Master Theorem to determine complexity
func AnalyzeMasterTheorem(input MasterTheoremInput) (ComplexityClass, MasterTheoremCase, string) {
	logBA := math.Log(float64(input.A)) / math.Log(float64(input.B))

	if input.C < logBA {
		// Case 1: f(n) is polynomially smaller than n^(log_b(a))
		return ComplexityClass(fmt.Sprintf("O(n^%.2f)", logBA)), Case1,
			fmt.Sprintf("Case 1: f(n) = O(n^%.2f) < O(n^%.2f)", input.C, logBA)
	} else if math.Abs(input.C-logBA) < 0.001 {
		// Case 2: f(n) is asymptotically equal to n^(log_b(a))
		if input.K >= 0 {
			return ComplexityClass(fmt.Sprintf("O(n^%.2f * log^%d(n))", logBA, input.K+1)), Case2,
				fmt.Sprintf("Case 2: f(n) = Θ(n^%.2f * log^%d(n))", input.C, input.K)
		}
		return ComplexityClass(fmt.Sprintf("O(n^%.2f * log(n))", logBA)), Case2,
			fmt.Sprintf("Case 2: f(n) = Θ(n^%.2f)", input.C)
	} else {
		// Case 3: f(n) is polynomially larger than n^(log_b(a))
		return ComplexityClass(fmt.Sprintf("O(n^%.2f)", input.C)), Case3,
			fmt.Sprintf("Case 3: f(n) = Ω(n^%.2f) > O(n^%.2f)", input.C, logBA)
	}
}

// 📊 AMORTIZED ANALYSIS TECHNIQUES
// For algorithms where individual operations vary in cost

// AmortizedAnalysisMethod represents different amortized analysis approaches
type AmortizedAnalysisMethod string

const (
	AggregateMethod  AmortizedAnalysisMethod = "Aggregate"
	AccountingMethod AmortizedAnalysisMethod = "Accounting"
	PotentialMethod  AmortizedAnalysisMethod = "Potential"
)

// AmortizedAnalysis represents amortized cost analysis
type AmortizedAnalysis struct {
	Method      AmortizedAnalysisMethod
	WorstCase   int     // Worst-case cost of single operation
	Amortized   float64 // Amortized cost per operation
	TotalCost   int     // Total cost of n operations
	Operations  int     // Number of operations
	Explanation string
}

// 🔍 ALGORITHM COMPLEXITY ANALYZER
// Practical tool for analyzing common algorithms

var CommonAlgorithms = []AlgorithmAnalysis{
	{
		Name:            "Binary Search",
		TimeComplexity:  OLogN,
		SpaceComplexity: O1,
		BestCase:        O1,
		AverageCase:     OLogN,
		WorstCase:       OLogN,
		Explanation:     "Eliminates half the search space each iteration",
	},
	{
		Name:            "Merge Sort",
		TimeComplexity:  ONLogN,
		SpaceComplexity: ON,
		BestCase:        ONLogN,
		AverageCase:     ONLogN,
		WorstCase:       ONLogN,
		Explanation:     "Divide-and-conquer with consistent O(n log n) performance",
	},
	{
		Name:            "Quick Sort",
		TimeComplexity:  ONLogN,
		SpaceComplexity: OLogN,
		BestCase:        ONLogN,
		AverageCase:     ONLogN,
		WorstCase:       ON2,
		Explanation:     "Worst case occurs with poor pivot selection",
	},
	{
		Name:            "DFS Graph Traversal",
		TimeComplexity:  ComplexityClass("O(V + E)"),
		SpaceComplexity: ComplexityClass("O(V)"),
		BestCase:        ComplexityClass("O(V + E)"),
		AverageCase:     ComplexityClass("O(V + E)"),
		WorstCase:       ComplexityClass("O(V + E)"),
		Explanation:     "Visits each vertex and edge exactly once",
	},
}

// 📈 COMPLEXITY GROWTH VISUALIZATION
// Understanding how different complexity classes scale

// ComplexityGrowth calculates operations for different complexity classes
func ComplexityGrowth(n int) map[ComplexityClass]int64 {
	growth := make(map[ComplexityClass]int64)

	growth[O1] = 1
	growth[OLogN] = int64(math.Log2(float64(n)))
	growth[ON] = int64(n)
	growth[ONLogN] = int64(n) * int64(math.Log2(float64(n)))
	growth[ON2] = int64(n * n)
	growth[ON3] = int64(n * n * n)

	// Cap exponential growth to prevent overflow
	if n <= 30 {
		growth[O2N] = int64(math.Pow(2, float64(n)))
	} else {
		growth[O2N] = math.MaxInt64
	}

	if n <= 12 {
		growth[ONFact] = factorial(n)
	} else {
		growth[ONFact] = math.MaxInt64
	}

	return growth
}

// factorial calculates n! for small values
func factorial(n int) int64 {
	if n <= 1 {
		return 1
	}
	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

// 🎯 PRACTICAL COMPLEXITY ANALYSIS TOOLS

// BenchmarkFunction measures actual runtime complexity
func BenchmarkFunction(f func(int), sizes []int) []time.Duration {
	results := make([]time.Duration, len(sizes))

	for i, size := range sizes {
		start := time.Now()
		f(size)
		results[i] = time.Since(start)
	}

	return results
}

// CompareComplexities shows growth differences between complexity classes
func CompareComplexities(maxN int) {
	fmt.Println("📊 COMPLEXITY GROWTH COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("%8s | %8s | %8s | %10s | %10s | %10s | %12s\n",
		"n", "O(1)", "O(log n)", "O(n)", "O(n log n)", "O(n²)", "O(2^n)")
	fmt.Println(strings.Repeat("-", 80))

	testSizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, n := range testSizes {
		if n > maxN {
			break
		}

		growth := ComplexityGrowth(n)
		fmt.Printf("%8d | %8d | %8d | %10d | %10d | %10d | %12d\n",
			n, growth[O1], growth[OLogN], growth[ON],
			growth[ONLogN], growth[ON2], growth[O2N])
	}
}

// 🔬 COMPLEXITY ANALYSIS EXAMPLES

// Example: Analyzing nested loops
func AnalyzeNestedLoops() {
	fmt.Println("\n🔍 NESTED LOOP COMPLEXITY ANALYSIS")
	fmt.Println("Analysis of common nested loop patterns:")

	examples := []struct {
		pattern     string
		complexity  ComplexityClass
		explanation string
	}{
		{
			pattern:     "for i := 0; i < n; i++ { for j := 0; j < n; j++ { ... } }",
			complexity:  ON2,
			explanation: "Two nested loops: n × n = n²",
		},
		{
			pattern:     "for i := 0; i < n; i++ { for j := i; j < n; j++ { ... } }",
			complexity:  ON2,
			explanation: "Triangle pattern: n + (n-1) + ... + 1 = n(n+1)/2 = O(n²)",
		},
		{
			pattern:     "for i := 0; i < n; i++ { for j := 0; j < m; j++ { ... } }",
			complexity:  ComplexityClass("O(nm)"),
			explanation: "Different variables: n × m operations",
		},
		{
			pattern:     "for i := 1; i < n; i *= 2 { for j := 0; j < n; j++ { ... } }",
			complexity:  ONLogN,
			explanation: "Outer loop: log n iterations, inner: n operations",
		},
	}

	for _, ex := range examples {
		fmt.Printf("\nPattern: %s\n", ex.pattern)
		fmt.Printf("Complexity: %s\n", ex.complexity)
		fmt.Printf("Explanation: %s\n", ex.explanation)
	}
}

// Main demonstration
func main() {
	fmt.Println("🎯 COMPLEXITY ANALYSIS THEORY - FAANG Interview Mastery")
	fmt.Println(strings.Repeat("=", 80))

	// Master Theorem Example
	fmt.Println("\n📐 MASTER THEOREM ANALYSIS")
	input := MasterTheoremInput{A: 2, B: 2, C: 1.0, K: 0}
	complexity, caseType, explanation := AnalyzeMasterTheorem(input)
	fmt.Printf("T(n) = %dT(n/%d) + O(n^%.1f)\n", input.A, input.B, input.C)
	fmt.Printf("Result: %s (Case %d: %s)\n", complexity, int(caseType)+1, explanation)

	// Complexity Growth Comparison
	fmt.Println("\n📈 COMPLEXITY GROWTH (n ≤ 64)")
	CompareComplexities(64)

	// Common Algorithm Analysis
	fmt.Println("\n🔍 COMMON ALGORITHM COMPLEXITIES")
	for _, algo := range CommonAlgorithms {
		fmt.Printf("\n%s:\n", algo.Name)
		fmt.Printf("  Time: %s, Space: %s\n", algo.TimeComplexity, algo.SpaceComplexity)
		fmt.Printf("  Best/Avg/Worst: %s/%s/%s\n", algo.BestCase, algo.AverageCase, algo.WorstCase)
		fmt.Printf("  %s\n", algo.Explanation)
	}

	// Nested Loop Analysis
	AnalyzeNestedLoops()

	fmt.Println("\n✅ Theory module complete - Master mathematical complexity analysis!")
}
