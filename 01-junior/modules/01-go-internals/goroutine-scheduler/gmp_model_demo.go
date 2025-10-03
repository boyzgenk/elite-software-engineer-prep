// Package scheduler demonstrates Go's G-M-P scheduler model
// This is production-level understanding needed for FAANG L3/L4 interviews

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// GMP Model Demonstration
// G = Goroutine (lightweight thread)
// M = OS Thread (machine)
// P = Processor (logical processor, GOMAXPROCS)

func main() {
	fmt.Println("🔥 Go Scheduler (G-M-P Model) Deep Dive")
	fmt.Println("=====================================")

	// Show initial scheduler state
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())

	// Demonstrate work-stealing scheduler
	demonstrateWorkStealing()

	// Show goroutine lifecycle
	demonstrateGoroutineLifecycle()

	// Demonstrate preemptive scheduling
	demonstratePreemption()

	// Show scheduler statistics
	showSchedulerStats()
}

// Work-stealing demonstration
// The Go scheduler uses work-stealing to balance load across P's
func demonstrateWorkStealing() {
	fmt.Println("\n🚀 Work-Stealing Scheduler Demo")
	fmt.Println("-------------------------------")

	const numWorkers = 8
	const numTasks = 1000

	var wg sync.WaitGroup
	taskChannel := make(chan int, numTasks)

	// Create tasks
	for i := 0; i < numTasks; i++ {
		taskChannel <- i
	}
	close(taskChannel)

	// Create workers (more goroutines than P's)
	for worker := 0; worker < numWorkers; worker++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			tasksProcessed := 0

			for task := range taskChannel {
				// Simulate work
				_ = task * task
				tasksProcessed++
			}

			fmt.Printf("Worker %d processed %d tasks\n", workerID, tasksProcessed)
		}(worker)
	}

	wg.Wait()
	fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}

// Goroutine lifecycle: Created -> Runnable -> Running -> Blocked -> Dead
func demonstrateGoroutineLifecycle() {
	fmt.Println("\n🎯 Goroutine Lifecycle Demo")
	fmt.Println("---------------------------")

	// Channel for coordination
	ch := make(chan string, 1)

	// Create goroutine (G created, enters runnable queue)
	go func() {
		fmt.Println("Goroutine: Created and Running")

		// Simulate blocking operation (G moves to waiting state)
		msg := <-ch
		fmt.Printf("Goroutine: Received message: %s\n", msg)

		// Goroutine will terminate (G moves to dead state)
		fmt.Println("Goroutine: Terminating")
	}()

	// Let goroutine start and block
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Goroutines while blocked: %d\n", runtime.NumGoroutine())

	// Unblock the goroutine
	ch <- "Hello from main!"

	// Wait for goroutine to finish
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Goroutines after completion: %d\n", runtime.NumGoroutine())
}

// Preemptive scheduling demonstration
// Go 1.14+ has non-cooperative preemption
func demonstratePreemption() {
	fmt.Println("\n⚡ Preemptive Scheduling Demo")
	fmt.Println("-----------------------------")

	done := make(chan bool)

	// CPU-intensive goroutine (would block in older Go versions)
	go func() {
		count := 0
		start := time.Now()

		// Tight loop without function calls (tests preemption)
		for time.Since(start) < 200*time.Millisecond {
			count++
		}

		fmt.Printf("CPU-intensive goroutine completed %d iterations\n", count)
		done <- true
	}()

	// This goroutine should still run due to preemption
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("Preempted goroutine tick %d\n", i+1)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	<-done
	time.Sleep(100 * time.Millisecond) // Let second goroutine finish
}

// Show detailed scheduler statistics
func showSchedulerStats() {
	fmt.Println("\n📊 Scheduler Statistics")
	fmt.Println("-----------------------")

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("OS threads: %d\n", runtime.GOMAXPROCS(0))

	// Force garbage collection to see scheduler interaction
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("GC cycles: %d\n", m.NumGC)
	fmt.Printf("Heap objects: %d\n", m.HeapObjects)
}

/*
🔥 FAANG Interview Key Points:

1. G-M-P Model:
   - G: Goroutines (logical threads, 2KB initial stack)
   - M: OS threads (expensive, limited by GOMAXPROCS)
   - P: Processors (scheduling contexts, holds runnable G queue)

2. Work Stealing:
   - Each P has local runnable queue
   - When P's queue is empty, it steals from other P's
   - Global runnable queue for overflow

3. Scheduler States:
   - Runnable: Ready to execute
   - Running: Currently executing on M
   - Blocked: Waiting for I/O, channel, etc.

4. Preemption (Go 1.14+):
   - Non-cooperative preemption using signals
   - Prevents CPU-intensive goroutines from starving others
   - Important for fairness in scheduling

5. Performance Implications:
   - Context switching is cheaper than OS threads
   - Stack growth is handled automatically
   - Scheduler overhead is minimal

This understanding demonstrates production-level Go expertise
required for senior backend positions at FAANG companies.
*/
