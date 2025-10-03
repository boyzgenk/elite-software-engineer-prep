# 🕸️ Graph Algorithms - FAANG Production Systems Mastery
## Master the Algorithms That Power Google, Meta, and Netflix

### 🎯 Elite Learning Objectives
**Master graph algorithms used in production systems at top tech companies.**

- **Graph Representations:** Adjacency lists, matrices, and edge lists with Go optimizations
- **Traversal Mastery:** DFS, BFS with applications in real-world systems
- **Shortest Path Excellence:** Dijkstra, Bellman-Ford, Floyd-Warshall for routing and optimization
- **Union-Find Mastery:** Disjoint set operations for network connectivity and clustering
- **Minimum Spanning Trees:** Kruskal and Prim's algorithms for network design
- **Topological Sorting:** Dependency resolution and task scheduling

### 🚀 Elite Success Metrics
- [ ] **Implementation Mastery:** Build production-quality Go implementations of all major graph algorithms
- [ ] **Problem Recognition:** Quickly identify which graph algorithm applies to interview problems
- [ ] **Optimization Skills:** Understand time-space tradeoffs and choose optimal approaches
- [ ] **Real-World Application:** Connect algorithms to actual systems (social networks, maps, distributed systems)

---

## 📚 Learning Modules

### 🔗 [Graph Representations](./graph-representations/)
**Efficient Data Structures for Graph Storage**
- **Adjacency Lists:** Memory-efficient, good for sparse graphs
- **Adjacency Matrices:** Fast edge queries, good for dense graphs  
- **Edge Lists:** Simple representation, good for algorithm implementations
- **Go Optimizations:** Using slices, maps, and custom types effectively

```go
// Adjacency List representation
type Graph struct {
    vertices int
    adjList  [][]int
}

// Adjacency Matrix representation  
type MatrixGraph struct {
    vertices int
    matrix   [][]bool
}
```

### 🔍 [Traversal Algorithms](./traversal-algorithms/)
**DFS and BFS - The Foundation of Graph Algorithms**

#### **Depth-First Search (DFS)**
- **Applications:** Cycle detection, pathfinding, topological sorting
- **Time Complexity:** O(V + E)
- **Space Complexity:** O(V) for recursion stack

#### **Breadth-First Search (BFS)**
- **Applications:** Shortest path in unweighted graphs, level-order traversal
- **Time Complexity:** O(V + E)  
- **Space Complexity:** O(V) for queue storage

### 🎯 [Shortest Path](./shortest-path/)
**Finding Optimal Paths in Weighted Graphs**

#### **Dijkstra's Algorithm**
- **Use Case:** Single-source shortest path with non-negative weights
- **Applications:** GPS navigation, network routing, social networks
- **Time Complexity:** O((V + E) log V) with binary heap
- **Space Complexity:** O(V)

#### **Bellman-Ford Algorithm**
- **Use Case:** Single-source shortest path with negative weights
- **Applications:** Currency arbitrage, network routing with costs
- **Time Complexity:** O(VE)
- **Space Complexity:** O(V)

#### **Floyd-Warshall Algorithm**
- **Use Case:** All-pairs shortest paths
- **Applications:** Transitive closure, graph analysis
- **Time Complexity:** O(V³)
- **Space Complexity:** O(V²)

### 🔗 [Union-Find](./union-find/)
**Disjoint Set Data Structure for Connectivity**
- **Operations:** Find, Union, Path Compression, Union by Rank
- **Applications:** Network connectivity, clustering, Kruskal's MST
- **Time Complexity:** O(α(n)) amortized (α is inverse Ackermann)
- **Space Complexity:** O(n)

### 🌲 [Minimum Spanning Tree](./minimum-spanning-tree/)
**Finding Minimum Cost Connected Subgraphs**

#### **Kruskal's Algorithm**
- **Approach:** Edge-based, uses Union-Find
- **Time Complexity:** O(E log E)
- **Space Complexity:** O(V)

#### **Prim's Algorithm**  
- **Approach:** Vertex-based, uses priority queue
- **Time Complexity:** O((V + E) log V)
- **Space Complexity:** O(V)

### 📋 [Topological Sort](./topological-sort/)
**Ordering Vertices in Directed Acyclic Graphs**
- **Applications:** Task scheduling, dependency resolution, compilation order
- **Kahn's Algorithm:** BFS-based approach
- **DFS-based:** Using finish times
- **Time Complexity:** O(V + E)
- **Space Complexity:** O(V)

---

## 🏋️ Elite Practice Framework

### 📝 [Practice Problems](./practice-problems/)
**FAANG-Style Graph Algorithm Questions**

#### **Easy Level:**
- Number of Islands (DFS/BFS)
- Valid Tree (Union-Find)
- Course Schedule (Topological Sort)

#### **Medium Level:**
- Word Ladder (BFS)
- Network Delay Time (Dijkstra)
- Accounts Merge (Union-Find)
- Minimum Spanning Tree (Kruskal/Prim)

#### **Hard Level:**
- Alien Dictionary (Topological Sort)
- Critical Connections (Bridges)
- Cheapest Flights Within K Stops (Modified Dijkstra)

### 🧪 [Tests](./tests/)
**Production-Grade Quality Assurance**
- Unit tests for all algorithm implementations
- Edge case handling (empty graphs, disconnected components)
- Performance benchmarks vs standard library
- Memory usage analysis and optimization

---

## 🎯 Real-World Applications

### 🌐 **Social Networks (Meta/Facebook)**
```go
// Friend recommendation using graph traversal
func recommendFriends(userID int, graph *SocialGraph) []int {
    visited := make(map[int]bool)
    recommendations := make(map[int]int) // userID -> mutual friends count
    
    // BFS to find friends of friends
    queue := []int{userID}
    visited[userID] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for _, friend := range graph.GetFriends(current) {
            if !visited[friend] {
                visited[friend] = true
                queue = append(queue, friend)
                
                // Count mutual friends for recommendation
                for _, friendOfFriend := range graph.GetFriends(friend) {
                    if friendOfFriend != userID && !graph.AreFriends(userID, friendOfFriend) {
                        recommendations[friendOfFriend]++
                    }
                }
            }
        }
    }
    
    return getTopRecommendations(recommendations, 10)
}
```

### 🗺️ **GPS Navigation (Google Maps)**
```go
// Route finding using Dijkstra's algorithm
func findShortestRoute(start, end Location, roadNetwork *WeightedGraph) []Location {
    distances := make(map[Location]float64)
    previous := make(map[Location]Location)
    pq := NewPriorityQueue()
    
    // Initialize distances
    for location := range roadNetwork.GetAllLocations() {
        if location == start {
            distances[location] = 0
            pq.Push(location, 0)
        } else {
            distances[location] = math.Inf(1)
        }
    }
    
    for !pq.IsEmpty() {
        current := pq.Pop()
        
        if current == end {
            return reconstructPath(previous, start, end)
        }
        
        for _, neighbor := range roadNetwork.GetNeighbors(current) {
            alt := distances[current] + roadNetwork.GetWeight(current, neighbor)
            if alt < distances[neighbor] {
                distances[neighbor] = alt
                previous[neighbor] = current
                pq.Push(neighbor, alt)
            }
        }
    }
    
    return nil // No path found
}
```

### 🔧 **Dependency Resolution (Build Systems)**
```go
// Topological sort for build dependency resolution
func resolveBuildOrder(dependencies map[string][]string) ([]string, error) {
    // Calculate in-degrees
    inDegree := make(map[string]int)
    allNodes := make(map[string]bool)
    
    for node, deps := range dependencies {
        allNodes[node] = true
        if _, exists := inDegree[node]; !exists {
            inDegree[node] = 0
        }
        
        for _, dep := range deps {
            allNodes[dep] = true
            inDegree[node]++
        }
    }
    
    // Kahn's algorithm
    queue := make([]string, 0)
    for node := range allNodes {
        if inDegree[node] == 0 {
            queue = append(queue, node)
        }
    }
    
    result := make([]string, 0, len(allNodes))
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        result = append(result, current)
        
        for node, deps := range dependencies {
            for _, dep := range deps {
                if dep == current {
                    inDegree[node]--
                    if inDegree[node] == 0 {
                        queue = append(queue, node)
                    }
                }
            }
        }
    }
    
    if len(result) != len(allNodes) {
        return nil, errors.New("circular dependency detected")
    }
    
    return result, nil
}
```

---

## 🧠 Algorithm Selection Strategy

### **Choose DFS when:**
- Exploring all possible paths
- Detecting cycles in directed graphs
- Topological sorting
- Finding strongly connected components

### **Choose BFS when:**
- Finding shortest path in unweighted graphs
- Level-order traversal
- Finding all nodes at distance K
- Testing bipartiteness

### **Choose Dijkstra when:**
- Single-source shortest path with non-negative weights
- Network routing protocols
- GPS navigation systems

### **Choose Bellman-Ford when:**
- Shortest path with negative weights
- Detecting negative cycles
- Currency arbitrage detection

### **Choose Union-Find when:**
- Dynamic connectivity queries
- Kruskal's MST algorithm
- Percolation problems
- Clustering algorithms

---

## 📊 Complexity Comparison

| Algorithm | Time Complexity | Space Complexity | Use Case |
|-----------|----------------|------------------|----------|
| DFS | O(V + E) | O(V) | Path finding, cycle detection |
| BFS | O(V + E) | O(V) | Shortest path (unweighted) |
| Dijkstra | O((V + E) log V) | O(V) | Shortest path (positive weights) |
| Bellman-Ford | O(VE) | O(V) | Shortest path (negative weights) |
| Floyd-Warshall | O(V³) | O(V²) | All-pairs shortest paths |
| Union-Find | O(α(n)) | O(n) | Dynamic connectivity |
| Kruskal's MST | O(E log E) | O(V) | Minimum spanning tree |
| Prim's MST | O((V + E) log V) | O(V) | Minimum spanning tree |
| Topological Sort | O(V + E) | O(V) | Dependency ordering |

---

## 🎯 Elite Mastery Checklist

### Level 1: Implementation
- [ ] Implement all graph representations in Go
- [ ] Code DFS and BFS from scratch
- [ ] Implement Dijkstra's algorithm with priority queue
- [ ] Build Union-Find with path compression and union by rank
- [ ] Code both Kruskal's and Prim's MST algorithms
- [ ] Implement topological sort using both DFS and Kahn's algorithm

### Level 2: Problem Solving
- [ ] Solve 20+ graph problems on LeetCode
- [ ] Recognize graph patterns in disguised problems
- [ ] Choose optimal algorithms based on problem constraints
- [ ] Handle edge cases (disconnected graphs, cycles, negative weights)

### Level 3: Optimization
- [ ] Optimize implementations for Go's runtime characteristics
- [ ] Implement concurrent versions where applicable
- [ ] Use appropriate data structures (slices vs maps vs custom types)
- [ ] Profile and benchmark algorithm performance

### Level 4: Real-World Application
- [ ] Design graph-based systems (social networks, recommendation engines)
- [ ] Understand how graph algorithms scale in production
- [ ] Connect algorithms to system design interview questions
- [ ] Explain trade-offs in terms of business requirements

---

**Elite Graph Algorithm Readiness:** You can implement any graph algorithm from scratch in 15-20 minutes, choose the optimal algorithm for any graph problem, and explain how these algorithms power real systems at scale.

**Ready to master the algorithms that connect the world? Let's build the graph algorithm expertise that powers FAANG systems! 🔥**