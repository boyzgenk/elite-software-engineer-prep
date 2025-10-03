// Package main demonstrates Fan-In/Fan-Out pattern
// Essential for distributed processing at Netflix, Uber scale
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents work to be processed
type Job struct {
	ID   int
	Data string
}

// Result represents processed job result
type Result struct {
	JobID       int
	Output      string
	ProcessedAt time.Time
	WorkerID    int
}

// FanOutFanIn distributes jobs to multiple workers and collects results
// This is the foundation of microservices architecture
func FanOutFanIn(jobs []Job, numWorkers int) []Result {
	// Create channels for job distribution and result collection
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	// Start workers (fan-out phase)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, jobChan, resultChan, &wg)
	}

	// Send jobs to workers
	go func() {
		defer close(jobChan)
		for _, job := range jobs {
			jobChan <- job
		}
	}()

	// Close result channel when all workers finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Fan-in: Collect all results
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// worker processes jobs - simulates CPU-intensive work
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Simulate processing time (realistic CPU work)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		result := Result{
			JobID:       job.ID,
			Output:      fmt.Sprintf("Processed: %s", job.Data),
			ProcessedAt: time.Now(),
			WorkerID:    id,
		}

		results <- result
	}
}

// Advanced: Dynamic worker scaling based on queue length
func FanOutFanInDynamic(jobs []Job, maxWorkers int) []Result {
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	// Dynamic worker management
	var wg sync.WaitGroup
	workerCount := 1 // Start with 1 worker

	// Start initial worker
	wg.Add(1)
	go worker(1, jobChan, resultChan, &wg)

	// Job sender with dynamic scaling logic
	go func() {
		defer close(jobChan)

		for i, job := range jobs {
			jobChan <- job

			// Scale up workers if queue is building up
			queueLen := len(jobChan)
			if queueLen > 2 && workerCount < maxWorkers {
				workerCount++
				wg.Add(1)
				go worker(workerCount, jobChan, resultChan, &wg)
				fmt.Printf("🚀 Scaled up to %d workers (queue: %d)\n", workerCount, queueLen)
			}

			// Small delay to demonstrate scaling
			if i < len(jobs)-1 {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Result collector
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func main() {
	fmt.Println("🔥 Fan-In/Fan-Out Pattern Demo")
	fmt.Println("==============================")

	// Create sample jobs
	jobs := make([]Job, 10)
	for i := 0; i < 10; i++ {
		jobs[i] = Job{
			ID:   i + 1,
			Data: fmt.Sprintf("user_event_%d", i+1),
		}
	}

	// Demo 1: Fixed worker pool
	fmt.Println("\n1. Fixed Worker Pool (3 workers)")
	start := time.Now()
	results := FanOutFanIn(jobs, 3)
	duration := time.Since(start)

	fmt.Printf("Processed %d jobs in %v:\n", len(results), duration)
	for _, result := range results {
		fmt.Printf("  Job %d -> Worker %d: %s\n",
			result.JobID, result.WorkerID, result.Output)
	}

	// Demo 2: Dynamic scaling
	fmt.Println("\n2. Dynamic Worker Scaling (max 4 workers)")
	start = time.Now()
	results = FanOutFanInDynamic(jobs, 4)
	duration = time.Since(start)

	fmt.Printf("Processed %d jobs in %v with dynamic scaling\n", len(results), duration)

	fmt.Println("\n🎯 Key Interview Points:")
	fmt.Println("- Demonstrates understanding of distributed processing")
	fmt.Println("- Shows knowledge of buffered channels and worker pools")
	fmt.Println("- Exhibits production patterns used at scale")
	fmt.Println("- Illustrates both static and dynamic scaling strategies")
}
