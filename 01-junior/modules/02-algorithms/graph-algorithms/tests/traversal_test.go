package graph

import (
	"reflect"
	"testing"
)

// TestNewGraph tests graph creation
func TestNewGraph(t *testing.T) {
	tests := []struct {
		name     string
		vertices int
		expected int
	}{
		{"Valid graph", 5, 5},
		{"Single vertex", 1, 1},
		{"Zero vertices", 0, 0},
		{"Negative vertices", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph(tt.vertices)
			if g.GetVertices() != tt.expected {
				t.Errorf("NewGraph(%d) = %d vertices, want %d", tt.vertices, g.GetVertices(), tt.expected)
			}
		})
	}
}

// TestAddEdge tests edge addition
func TestAddEdge(t *testing.T) {
	g := NewGraph(4)

	tests := []struct {
		name    string
		src     int
		dest    int
		wantErr bool
	}{
		{"Valid edge", 0, 1, false},
		{"Self loop", 2, 2, false},
		{"Source out of bounds", 5, 1, true},
		{"Destination out of bounds", 1, 5, true},
		{"Negative source", -1, 1, true},
		{"Negative destination", 1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := g.AddEdge(tt.src, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEdge(%d, %d) error = %v, wantErr %v", tt.src, tt.dest, err, tt.wantErr)
			}
		})
	}
}

// TestDFSRecursive tests recursive DFS implementation
func TestDFSRecursive(t *testing.T) {
	// Create test graph: 0->1->2->3, 0->2
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	tests := []struct {
		name     string
		start    int
		expected []int
		wantErr  bool
	}{
		{"DFS from 0", 0, []int{0, 1, 2, 3}, false},
		{"DFS from 2", 2, []int{2, 3}, false},
		{"DFS from 3", 3, []int{3}, false},
		{"Invalid start vertex", 5, nil, true},
		{"Negative start vertex", -1, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := g.DFSRecursive(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("DFSRecursive(%d) error = %v, wantErr %v", tt.start, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DFSRecursive(%d) = %v, want %v", tt.start, result, tt.expected)
			}
		})
	}
}

// TestDFSIterative tests iterative DFS implementation
func TestDFSIterative(t *testing.T) {
	// Create test graph: 0->1->2->3, 0->2
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	tests := []struct {
		name     string
		start    int
		expected []int
		wantErr  bool
	}{
		{"DFS from 0", 0, []int{0, 2, 3, 1}, false}, // Different order due to stack
		{"DFS from 2", 2, []int{2, 3}, false},
		{"DFS from 3", 3, []int{3}, false},
		{"Invalid start vertex", 5, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := g.DFSIterative(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("DFSIterative(%d) error = %v, wantErr %v", tt.start, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DFSIterative(%d) = %v, want %v", tt.start, result, tt.expected)
			}
		})
	}
}

// TestBFS tests breadth-first search
func TestBFS(t *testing.T) {
	// Create test graph: 0->1->3, 0->2->3
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)

	tests := []struct {
		name     string
		start    int
		expected []int
		wantErr  bool
	}{
		{"BFS from 0", 0, []int{0, 1, 2, 3}, false},
		{"BFS from 1", 1, []int{1, 3}, false},
		{"BFS from 3", 3, []int{3}, false},
		{"Invalid start vertex", 5, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := g.BFS(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("BFS(%d) error = %v, wantErr %v", tt.start, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BFS(%d) = %v, want %v", tt.start, result, tt.expected)
			}
		})
	}
}

// TestBFSShortestPath tests shortest path finding
func TestBFSShortestPath(t *testing.T) {
	// Create test graph: 0->1->3, 0->2->3
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)

	tests := []struct {
		name         string
		start        int
		end          int
		expectedPath []int
		expectedDist int
		wantErr      bool
	}{
		{"Path 0 to 3", 0, 3, []int{0, 1, 3}, 2, false},
		{"Path 0 to 0", 0, 0, []int{0}, 0, false},
		{"Path 1 to 3", 1, 3, []int{1, 3}, 1, false},
		{"No path 3 to 0", 3, 0, nil, -1, false},
		{"Invalid start", 5, 1, nil, -1, true},
		{"Invalid end", 1, 5, nil, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, dist, err := g.BFSShortestPath(tt.start, tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("BFSShortestPath(%d, %d) error = %v, wantErr %v", tt.start, tt.end, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(path, tt.expectedPath) {
					t.Errorf("BFSShortestPath(%d, %d) path = %v, want %v", tt.start, tt.end, path, tt.expectedPath)
				}
				if dist != tt.expectedDist {
					t.Errorf("BFSShortestPath(%d, %d) distance = %d, want %d", tt.start, tt.end, dist, tt.expectedDist)
				}
			}
		})
	}
}

// TestHasCycle tests cycle detection
func TestHasCycle(t *testing.T) {
	tests := []struct {
		name     string
		edges    [][2]int
		vertices int
		expected bool
	}{
		{
			name:     "No cycle",
			edges:    [][2]int{{0, 1}, {1, 2}, {2, 3}},
			vertices: 4,
			expected: false,
		},
		{
			name:     "Has cycle",
			edges:    [][2]int{{0, 1}, {1, 2}, {2, 0}},
			vertices: 3,
			expected: true,
		},
		{
			name:     "Self loop",
			edges:    [][2]int{{0, 0}},
			vertices: 1,
			expected: true,
		},
		{
			name:     "Empty graph",
			edges:    [][2]int{},
			vertices: 3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph(tt.vertices)
			for _, edge := range tt.edges {
				g.AddEdge(edge[0], edge[1])
			}

			result := g.HasCycle()
			if result != tt.expected {
				t.Errorf("HasCycle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsConnected tests connectivity checking
func TestIsConnected(t *testing.T) {
	tests := []struct {
		name     string
		edges    [][2]int
		vertices int
		expected bool
	}{
		{
			name:     "Connected graph",
			edges:    [][2]int{{0, 1}, {1, 2}, {2, 3}},
			vertices: 4,
			expected: true,
		},
		{
			name:     "Disconnected graph",
			edges:    [][2]int{{0, 1}, {2, 3}},
			vertices: 4,
			expected: false,
		},
		{
			name:     "Single vertex",
			edges:    [][2]int{},
			vertices: 1,
			expected: true,
		},
		{
			name:     "Empty graph",
			edges:    [][2]int{},
			vertices: 0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph(tt.vertices)
			for _, edge := range tt.edges {
				g.AddEdge(edge[0], edge[1])
			}

			result := g.IsConnected()
			if result != tt.expected {
				t.Errorf("IsConnected() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAddUndirectedEdge tests undirected edge addition
func TestAddUndirectedEdge(t *testing.T) {
	g := NewGraph(3)

	err := g.AddUndirectedEdge(0, 1)
	if err != nil {
		t.Errorf("AddUndirectedEdge(0, 1) error = %v, want nil", err)
	}

	// Check if both directions exist
	adj0, _ := g.GetAdjacent(0)
	adj1, _ := g.GetAdjacent(1)

	if !contains(adj0, 1) {
		t.Errorf("Edge 0->1 not found in adjacency list")
	}
	if !contains(adj1, 0) {
		t.Errorf("Edge 1->0 not found in adjacency list")
	}
}

// TestGetEdgeCount tests edge counting
func TestGetEdgeCount(t *testing.T) {
	g := NewGraph(4)

	if g.GetEdgeCount() != 0 {
		t.Errorf("Empty graph edge count = %d, want 0", g.GetEdgeCount())
	}

	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)

	if g.GetEdgeCount() != 3 {
		t.Errorf("Graph edge count = %d, want 3", g.GetEdgeCount())
	}
}

// TestClone tests graph cloning
func TestClone(t *testing.T) {
	original := NewGraph(3)
	original.AddEdge(0, 1)
	original.AddEdge(1, 2)

	cloned := original.Clone()

	// Test structure equality
	if cloned.GetVertices() != original.GetVertices() {
		t.Errorf("Cloned vertices = %d, want %d", cloned.GetVertices(), original.GetVertices())
	}

	if cloned.GetEdgeCount() != original.GetEdgeCount() {
		t.Errorf("Cloned edge count = %d, want %d", cloned.GetEdgeCount(), original.GetEdgeCount())
	}

	// Test independence (modifying clone shouldn't affect original)
	cloned.AddEdge(2, 0)

	if cloned.GetEdgeCount() == original.GetEdgeCount() {
		t.Errorf("Clone is not independent of original")
	}
}

// Helper function to check if slice contains element
func contains(slice []int, element int) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Benchmark tests for performance analysis
func BenchmarkDFSRecursive(b *testing.B) {
	g := createLargeGraph(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.DFSRecursive(0)
	}
}

func BenchmarkDFSIterative(b *testing.B) {
	g := createLargeGraph(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.DFSIterative(0)
	}
}

func BenchmarkBFS(b *testing.B) {
	g := createLargeGraph(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.BFS(0)
	}
}

func BenchmarkBFSShortestPath(b *testing.B) {
	g := createLargeGraph(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.BFSShortestPath(0, 999)
	}
}

// Helper function to create a large graph for benchmarking
func createLargeGraph(size int) *Graph {
	g := NewGraph(size)

	// Create a connected graph with some cycles
	for i := 0; i < size-1; i++ {
		g.AddEdge(i, i+1)
		if i%10 == 0 && i > 0 {
			g.AddEdge(i, i-5) // Add some back edges for cycles
		}
	}

	return g
}
