// Package main demonstrates Pipeline Pattern
// Used for data transformation chains at Uber, LinkedIn
package main

import (
	"fmt"
	"strings"
	"time"
)

// Pipeline represents a data processing pipeline
type Pipeline struct {
	stages []func(<-chan interface{}) <-chan interface{}
}

// NewPipeline creates a new processing pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage adds a processing stage to the pipeline
func (p *Pipeline) AddStage(stage func(<-chan interface{}) <-chan interface{}) {
	p.stages = append(p.stages, stage)
}

// Process runs data through the entire pipeline
func (p *Pipeline) Process(input <-chan interface{}) <-chan interface{} {
	current := input

	// Chain all stages together
	for _, stage := range p.stages {
		current = stage(current)
	}

	return current
}

// =============================================================================
// PIPELINE STAGES - Each stage is focused and testable
// =============================================================================

// validationStage filters out invalid data
func validationStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			if str, ok := data.(string); ok && len(str) > 0 && !strings.Contains(str, "invalid") {
				output <- str
			} else {
				fmt.Printf("❌ Filtered invalid data: %v\n", data)
			}
		}
	}()

	return output
}

// transformationStage processes and transforms data
func transformationStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			if str, ok := data.(string); ok {
				// Simulate transformation logic
				transformed := strings.ToUpper(str) + "_TRANSFORMED"
				output <- transformed
			}
		}
	}()

	return output
}

// enrichmentStage adds metadata to data
func enrichmentStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			if str, ok := data.(string); ok {
				// Add timestamp and metadata
				enriched := fmt.Sprintf("%s | enriched_at:%d", str, time.Now().Unix())
				output <- enriched
			}
		}
	}()

	return output
}

// persistenceStage simulates saving data to database
func persistenceStage(input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)
		for data := range input {
			// Simulate database write time
			time.Sleep(5 * time.Millisecond)

			persisted := fmt.Sprintf("PERSISTED: %s", data)
			output <- persisted
		}
	}()

	return output
}

// =============================================================================
// ADVANCED: Buffered Pipeline with Error Handling
// =============================================================================

// BufferedPipeline provides better performance with configurable buffer sizes
type BufferedPipeline struct {
	stages     []func(<-chan interface{}) <-chan interface{}
	bufferSize int
}

// NewBufferedPipeline creates pipeline with specified buffer size
func NewBufferedPipeline(bufferSize int) *BufferedPipeline {
	return &BufferedPipeline{bufferSize: bufferSize}
}

// AddStage adds a buffered stage
func (bp *BufferedPipeline) AddStage(stage func(<-chan interface{}) <-chan interface{}) {
	bp.stages = append(bp.stages, stage)
}

// Process with buffered channels for better throughput
func (bp *BufferedPipeline) Process(input <-chan interface{}) <-chan interface{} {
	current := input

	for i, stage := range bp.stages {
		// Create buffered output for each stage
		buffered := make(chan interface{}, bp.bufferSize)

		go func(in <-chan interface{}, out chan interface{}, stageNum int) {
			defer close(out)
			stageOutput := stage(in)

			for data := range stageOutput {
				out <- data
			}
		}(current, buffered, i+1)

		current = buffered
	}

	return current
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

func main() {
	fmt.Println("🔥 Pipeline Pattern Demo")
	fmt.Println("========================")

	// Demo 1: Basic Pipeline
	fmt.Println("\n1. Basic Data Processing Pipeline")

	// Create input data
	input := make(chan interface{}, 10)
	go func() {
		defer close(input)
		testData := []string{
			"user_event_1",
			"invalid_data",
			"user_event_2",
			"user_event_3",
			"",
			"user_event_4",
		}

		for _, data := range testData {
			input <- data
		}
	}()

	// Build pipeline
	pipeline := NewPipeline()
	pipeline.AddStage(validationStage)
	pipeline.AddStage(transformationStage)
	pipeline.AddStage(enrichmentStage)
	pipeline.AddStage(persistenceStage)

	// Process data through pipeline
	output := pipeline.Process(input)

	fmt.Println("\nPipeline Results:")
	for result := range output {
		fmt.Printf("✅ %s\n", result)
	}

	// Demo 2: Buffered Pipeline Performance Test
	fmt.Println("\n2. Buffered Pipeline Performance Comparison")

	// Test data
	largeInput := make(chan interface{}, 100)
	go func() {
		defer close(largeInput)
		for i := 1; i <= 50; i++ {
			largeInput <- fmt.Sprintf("event_%d", i)
		}
	}()

	// Buffered pipeline
	bufferedPipeline := NewBufferedPipeline(10)
	bufferedPipeline.AddStage(transformationStage)
	bufferedPipeline.AddStage(enrichmentStage)

	start := time.Now()
	bufferedOutput := bufferedPipeline.Process(largeInput)

	count := 0
	for range bufferedOutput {
		count++
	}

	duration := time.Since(start)
	fmt.Printf("Processed %d items in %v with buffered pipeline\n", count, duration)

	fmt.Println("\n🎯 Key Interview Points:")
	fmt.Println("- Demonstrates separation of concerns (each stage has one responsibility)")
	fmt.Println("- Shows understanding of channel chaining and data flow")
	fmt.Println("- Exhibits knowledge of buffered channels for performance")
	fmt.Println("- Illustrates composable and testable architecture")
	fmt.Println("- Shows real-world application in data processing systems")
}
