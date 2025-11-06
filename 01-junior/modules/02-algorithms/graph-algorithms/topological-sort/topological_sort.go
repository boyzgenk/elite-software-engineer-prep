// Topological Sort Algorithms in Go
// Implementation of DFS and Kahn's algorithm for topological sorting

package main

import (
	"fmt"
	"sort"
)

// ========================================================================
// GRAPH STRUCTURES
// ========================================================================

// Edge represents a directed edge
type Edge struct {
	To int
}

// DirectedGraph represents a directed graph using adjacency list
type DirectedGraph struct {
	vertices int
	adjList  [][]Edge
}

// NewDirectedGraph creates a new directed graph
func NewDirectedGraph(vertices int) *DirectedGraph {
	adjList := make([][]Edge, vertices)
	for i := range adjList {
		adjList[i] = make([]Edge, 0)
	}

	return &DirectedGraph{
		vertices: vertices,
		adjList:  adjList,
	}
}

// AddEdge adds a directed edge from u to v
func (g *DirectedGraph) AddEdge(from, to int) {
	if from < 0 || from >= g.vertices || to < 0 || to >= g.vertices {
		return
	}
	g.adjList[from] = append(g.adjList[from], Edge{To: to})
}

// GetNeighbors returns neighbors of a vertex
func (g *DirectedGraph) GetNeighbors(vertex int) []Edge {
	if vertex < 0 || vertex >= g.vertices {
		return nil
	}
	return g.adjList[vertex]
}

// GetVertexCount returns number of vertices
func (g *DirectedGraph) GetVertexCount() int {
	return g.vertices
}

// Print displays the graph
func (g *DirectedGraph) Print() {
	fmt.Println("Directed Graph:")
	for i := 0; i < g.vertices; i++ {
		fmt.Printf("Vertex %d: ", i)
		for j, edge := range g.adjList[i] {
			if j > 0 {
				fmt.Print(" -> ")
			}
			fmt.Printf("%d", edge.To)
		}
		fmt.Println()
	}
}

// ========================================================================
// DFS-BASED TOPOLOGICAL SORT
// ========================================================================

// TopologicalSortDFS performs topological sort using DFS
// Time Complexity: O(V + E)
// Space Complexity: O(V) for recursion stack and visited array
type DFSTopologicalSort struct {
	graph   *DirectedGraph
	visited []bool
	stack   []int
}

// NewDFSTopologicalSort creates a new DFS-based topological sorter
func NewDFSTopologicalSort(graph *DirectedGraph) *DFSTopologicalSort {
	return &DFSTopologicalSort{
		graph:   graph,
		visited: make([]bool, graph.vertices),
		stack:   make([]int, 0),
	}
}

// Sort performs topological sort using DFS
func (dfs *DFSTopologicalSort) Sort() ([]int, error) {
	// Reset state
	dfs.visited = make([]bool, dfs.graph.vertices)
	dfs.stack = make([]int, 0)

	// Check for cycles using DFS
	if dfs.hasCycle() {
		return nil, fmt.Errorf("graph contains a cycle, topological sort not possible")
	}

	// Reset for actual topological sort
	dfs.visited = make([]bool, dfs.graph.vertices)
	dfs.stack = make([]int, 0)

	// Visit all vertices
	for i := 0; i < dfs.graph.vertices; i++ {
		if !dfs.visited[i] {
			dfs.dfsUtil(i)
		}
	}

	// Reverse the stack to get topological order
	result := make([]int, len(dfs.stack))
	for i, v := range dfs.stack {
		result[len(dfs.stack)-1-i] = v
	}

	return result, nil
}

// dfsUtil is the recursive DFS utility function
func (dfs *DFSTopologicalSort) dfsUtil(vertex int) {
	dfs.visited[vertex] = true

	// Visit all neighbors first
	for _, edge := range dfs.graph.GetNeighbors(vertex) {
		if !dfs.visited[edge.To] {
			dfs.dfsUtil(edge.To)
		}
	}

	// Add current vertex to stack after visiting all neighbors
	dfs.stack = append(dfs.stack, vertex)
}

// hasCycle checks if the graph has a cycle using DFS
func (dfs *DFSTopologicalSort) hasCycle() bool {
	visited := make([]bool, dfs.graph.vertices)
	recStack := make([]bool, dfs.graph.vertices)

	for i := 0; i < dfs.graph.vertices; i++ {
		if !visited[i] {
			if dfs.isCyclicUtil(i, visited, recStack) {
				return true
			}
		}
	}
	return false
}

// isCyclicUtil is utility function for cycle detection
func (dfs *DFSTopologicalSort) isCyclicUtil(vertex int, visited, recStack []bool) bool {
	visited[vertex] = true
	recStack[vertex] = true

	for _, edge := range dfs.graph.GetNeighbors(vertex) {
		if !visited[edge.To] {
			if dfs.isCyclicUtil(edge.To, visited, recStack) {
				return true
			}
		} else if recStack[edge.To] {
			return true // Back edge found - cycle detected
		}
	}

	recStack[vertex] = false
	return false
}

// ========================================================================
// KAHN'S ALGORITHM (BFS-BASED)
// ========================================================================

// KahnTopologicalSort implements Kahn's algorithm for topological sorting
// Time Complexity: O(V + E)
// Space Complexity: O(V) for queue and in-degree array
type KahnTopologicalSort struct {
	graph *DirectedGraph
}

// NewKahnTopologicalSort creates a new Kahn's algorithm sorter
func NewKahnTopologicalSort(graph *DirectedGraph) *KahnTopologicalSort {
	return &KahnTopologicalSort{graph: graph}
}

// Sort performs topological sort using Kahn's algorithm
func (kahn *KahnTopologicalSort) Sort() ([]int, error) {
	vertices := kahn.graph.GetVertexCount()

	// Calculate in-degrees for all vertices
	inDegree := make([]int, vertices)
	for u := 0; u < vertices; u++ {
		for _, edge := range kahn.graph.GetNeighbors(u) {
			inDegree[edge.To]++
		}
	}

	// Initialize queue with vertices having in-degree 0
	queue := make([]int, 0)
	for i := 0; i < vertices; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	result := make([]int, 0)

	// Process vertices in topological order
	for len(queue) > 0 {
		// Dequeue vertex with in-degree 0
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// For each neighbor of current vertex
		for _, edge := range kahn.graph.GetNeighbors(current) {
			inDegree[edge.To]--

			// If in-degree becomes 0, add to queue
			if inDegree[edge.To] == 0 {
				queue = append(queue, edge.To)
			}
		}
	}

	// Check if all vertices are included (no cycle)
	if len(result) != vertices {
		return nil, fmt.Errorf("graph contains a cycle, topological sort not possible")
	}

	return result, nil
}

// ========================================================================
// ALL POSSIBLE TOPOLOGICAL ORDERS
// ========================================================================

// AllTopologicalOrders finds all possible topological orderings
type AllTopologicalOrders struct {
	graph    *DirectedGraph
	inDegree []int
	result   [][]int
}

// NewAllTopologicalOrders creates a new all topological orders finder
func NewAllTopologicalOrders(graph *DirectedGraph) *AllTopologicalOrders {
	return &AllTopologicalOrders{
		graph:  graph,
		result: make([][]int, 0),
	}
}

// FindAll finds all possible topological orderings
func (ato *AllTopologicalOrders) FindAll() [][]int {
	vertices := ato.graph.GetVertexCount()

	// Calculate in-degrees
	ato.inDegree = make([]int, vertices)
	for u := 0; u < vertices; u++ {
		for _, edge := range ato.graph.GetNeighbors(u) {
			ato.inDegree[edge.To]++
		}
	}

	ato.result = make([][]int, 0)
	path := make([]int, 0)
	ato.findAllUtil(path)

	return ato.result
}

// findAllUtil is recursive utility for finding all topological orders
func (ato *AllTopologicalOrders) findAllUtil(path []int) {
	// Find all vertices with in-degree 0
	candidates := make([]int, 0)
	for i := 0; i < ato.graph.GetVertexCount(); i++ {
		if ato.inDegree[i] == 0 {
			candidates = append(candidates, i)
		}
	}

	// If no candidates, we have a complete ordering
	if len(candidates) == 0 {
		if len(path) == ato.graph.GetVertexCount() {
			// Make a copy of the path
			order := make([]int, len(path))
			copy(order, path)
			ato.result = append(ato.result, order)
		}
		return
	}

	// Try each candidate
	for _, vertex := range candidates {
		// Mark vertex as used
		ato.inDegree[vertex] = -1

		// Reduce in-degree of neighbors
		for _, edge := range ato.graph.GetNeighbors(vertex) {
			ato.inDegree[edge.To]--
		}

		// Add to current path and recurse
		path = append(path, vertex)
		ato.findAllUtil(path)

		// Backtrack
		path = path[:len(path)-1]

		// Restore in-degrees
		ato.inDegree[vertex] = 0
		for _, edge := range ato.graph.GetNeighbors(vertex) {
			ato.inDegree[edge.To]++
		}
	}
}

// ========================================================================
// LEXICOGRAPHICALLY SMALLEST TOPOLOGICAL ORDER
// ========================================================================

// LexicographicalTopologicalSort finds lexicographically smallest topological order
type LexicographicalTopologicalSort struct {
	graph *DirectedGraph
}

// NewLexicographicalTopologicalSort creates new lexicographical sorter
func NewLexicographicalTopologicalSort(graph *DirectedGraph) *LexicographicalTopologicalSort {
	return &LexicographicalTopologicalSort{graph: graph}
}

// Sort finds lexicographically smallest topological ordering
func (lex *LexicographicalTopologicalSort) Sort() ([]int, error) {
	vertices := lex.graph.GetVertexCount()

	// Calculate in-degrees
	inDegree := make([]int, vertices)
	for u := 0; u < vertices; u++ {
		for _, edge := range lex.graph.GetNeighbors(u) {
			inDegree[edge.To]++
		}
	}

	result := make([]int, 0)

	for len(result) < vertices {
		// Find all vertices with in-degree 0
		candidates := make([]int, 0)
		for i := 0; i < vertices; i++ {
			if inDegree[i] == 0 {
				candidates = append(candidates, i)
			}
		}

		if len(candidates) == 0 {
			return nil, fmt.Errorf("graph contains a cycle")
		}

		// Sort candidates to get lexicographically smallest
		sort.Ints(candidates)

		// Pick the smallest candidate
		chosen := candidates[0]
		result = append(result, chosen)
		inDegree[chosen] = -1 // Mark as processed

		// Update in-degrees of neighbors
		for _, edge := range lex.graph.GetNeighbors(chosen) {
			inDegree[edge.To]--
		}
	}

	return result, nil
}

// ========================================================================
// PRACTICAL APPLICATIONS
// ========================================================================

// TaskScheduler demonstrates practical use of topological sort
type TaskScheduler struct {
	tasks        []string
	dependencies map[string][]string
	graph        *DirectedGraph
	taskToID     map[string]int
	idToTask     map[int]string
}

// NewTaskScheduler creates a new task scheduler
func NewTaskScheduler(tasks []string, dependencies map[string][]string) *TaskScheduler {
	taskToID := make(map[string]int)
	idToTask := make(map[int]string)

	for i, task := range tasks {
		taskToID[task] = i
		idToTask[i] = task
	}

	graph := NewDirectedGraph(len(tasks))

	// Build dependency graph
	for task, deps := range dependencies {
		taskID := taskToID[task]
		for _, dep := range deps {
			depID := taskToID[dep]
			graph.AddEdge(depID, taskID) // dep -> task
		}
	}

	return &TaskScheduler{
		tasks:        tasks,
		dependencies: dependencies,
		graph:        graph,
		taskToID:     taskToID,
		idToTask:     idToTask,
	}
}

// Schedule returns the order in which tasks should be executed
func (ts *TaskScheduler) Schedule() ([]string, error) {
	kahn := NewKahnTopologicalSort(ts.graph)
	order, err := kahn.Sort()
	if err != nil {
		return nil, fmt.Errorf("circular dependency detected: %v", err)
	}

	result := make([]string, len(order))
	for i, id := range order {
		result[i] = ts.idToTask[id]
	}

	return result, nil
}

// ========================================================================
// COMPARISON AND BENCHMARKING
// ========================================================================

type TopologicalSortComparison struct{}

func (tsc *TopologicalSortComparison) CompareAlgorithms(graph *DirectedGraph) {
	fmt.Println("🔍 TOPOLOGICAL SORT ALGORITHM COMPARISON")
	fmt.Println("=========================================")

	// DFS-based approach
	fmt.Println("\n=== DFS-Based Topological Sort ===")
	dfs := NewDFSTopologicalSort(graph)
	dfsResult, dfsErr := dfs.Sort()

	if dfsErr != nil {
		fmt.Printf("❌ DFS Error: %v\n", dfsErr)
	} else {
		fmt.Printf("✅ DFS Result: %v\n", dfsResult)
	}

	// Kahn's algorithm
	fmt.Println("\n=== Kahn's Algorithm (BFS-Based) ===")
	kahn := NewKahnTopologicalSort(graph)
	kahnResult, kahnErr := kahn.Sort()

	if kahnErr != nil {
		fmt.Printf("❌ Kahn Error: %v\n", kahnErr)
	} else {
		fmt.Printf("✅ Kahn Result: %v\n", kahnResult)
	}

	// Lexicographical ordering
	fmt.Println("\n=== Lexicographically Smallest ===")
	lex := NewLexicographicalTopologicalSort(graph)
	lexResult, lexErr := lex.Sort()

	if lexErr != nil {
		fmt.Printf("❌ Lex Error: %v\n", lexErr)
	} else {
		fmt.Printf("✅ Lex Result: %v\n", lexResult)
	}

	// All possible orderings (for small graphs)
	if graph.GetVertexCount() <= 6 {
		fmt.Println("\n=== All Possible Topological Orders ===")
		allOrders := NewAllTopologicalOrders(graph)
		allResults := allOrders.FindAll()
		fmt.Printf("Found %d possible orderings:\n", len(allResults))
		for i, order := range allResults {
			fmt.Printf("  %d: %v\n", i+1, order)
		}
	}
}

func (tsc *TopologicalSortComparison) ShowCharacteristics() {
	fmt.Println("\n📊 ALGORITHM CHARACTERISTICS")
	fmt.Println("=============================")

	fmt.Println("\n=== DFS-Based Topological Sort ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • Simple recursive implementation")
	fmt.Println("  • Natural with DFS traversal")
	fmt.Println("  • Memory efficient (uses call stack)")
	fmt.Println("❌ Cons:")
	fmt.Println("  • Requires cycle detection as separate step")
	fmt.Println("  • Stack overflow risk for deep recursion")
	fmt.Println("  • Less intuitive order generation")

	fmt.Println("\n=== Kahn's Algorithm ===")
	fmt.Println("✅ Pros:")
	fmt.Println("  • Intuitive approach (remove vertices with no dependencies)")
	fmt.Println("  • Cycle detection built-in")
	fmt.Println("  • No recursion (iterative)")
	fmt.Println("  • Can track progress during execution")
	fmt.Println("❌ Cons:")
	fmt.Println("  • Requires extra space for in-degree array")
	fmt.Println("  • Need to compute in-degrees initially")

	fmt.Println("\n=== Use Cases ===")
	fmt.Println("• Course Prerequisites: Schedule courses based on dependencies")
	fmt.Println("• Build Systems: Compile files in dependency order")
	fmt.Println("• Package Management: Install packages in correct order")
	fmt.Println("• Task Scheduling: Execute tasks respecting dependencies")
	fmt.Println("• Spreadsheet Recalculation: Update cells in dependency order")
}

// ========================================================================
// DEMONSTRATION AND TESTING
// ========================================================================

func main() {
	fmt.Println("📊 TOPOLOGICAL SORT ALGORITHMS DEMO")
	fmt.Println("====================================")

	// Create a sample DAG representing course prerequisites
	// Courses: 0=Math, 1=Physics, 2=Chemistry, 3=Biology, 4=CompSci, 5=AI
	graph := NewDirectedGraph(6)
	graph.AddEdge(0, 1) // Math -> Physics
	graph.AddEdge(0, 2) // Math -> Chemistry
	graph.AddEdge(1, 4) // Physics -> CompSci
	graph.AddEdge(2, 3) // Chemistry -> Biology
	graph.AddEdge(4, 5) // CompSci -> AI

	fmt.Println("\n=== Sample DAG (Course Prerequisites) ===")
	courseNames := []string{"Math", "Physics", "Chemistry", "Biology", "CompSci", "AI"}
	for i, name := range courseNames {
		fmt.Printf("Vertex %d: %s\n", i, name)
	}
	graph.Print()

	// Compare different algorithms
	comparison := &TopologicalSortComparison{}
	comparison.CompareAlgorithms(graph)

	// Task scheduling example
	fmt.Println("\n=== PRACTICAL EXAMPLE: TASK SCHEDULING ===")
	tasks := []string{"setup", "compile", "test", "package", "deploy"}
	dependencies := map[string][]string{
		"compile": {"setup"},
		"test":    {"compile"},
		"package": {"test"},
		"deploy":  {"package"},
	}

	scheduler := NewTaskScheduler(tasks, dependencies)
	schedule, err := scheduler.Schedule()

	if err != nil {
		fmt.Printf("❌ Scheduling Error: %v\n", err)
	} else {
		fmt.Printf("✅ Task Execution Order: %v\n", schedule)
	}

	// Example with cycle detection
	fmt.Println("\n=== CYCLE DETECTION EXAMPLE ===")
	cyclicGraph := NewDirectedGraph(3)
	cyclicGraph.AddEdge(0, 1)
	cyclicGraph.AddEdge(1, 2)
	cyclicGraph.AddEdge(2, 0) // Creates cycle: 0 -> 1 -> 2 -> 0

	fmt.Println("Cyclic Graph:")
	cyclicGraph.Print()

	dfsSort := NewDFSTopologicalSort(cyclicGraph)
	_, err = dfsSort.Sort()
	if err != nil {
		fmt.Printf("✅ Cycle detected correctly: %v\n", err)
	}

	// Show algorithm characteristics
	comparison.ShowCharacteristics()

	fmt.Println("\n🎯 Key Takeaways:")
	fmt.Println("• Use DFS-based for simple implementations")
	fmt.Println("• Use Kahn's for better cycle detection and progress tracking")
	fmt.Println("• Both have O(V + E) time complexity")
	fmt.Println("• Essential for dependency resolution problems")
}
