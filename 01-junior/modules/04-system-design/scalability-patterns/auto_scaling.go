package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// FAANG Interview Focus: Auto-scaling Implementation Patterns
// Key Topics: Metrics-based scaling, thresholds, predictive scaling, Go channels

// ScalingMetric represents different types of metrics for scaling decisions
type ScalingMetric int

const (
	CPUUtilization ScalingMetric = iota
	MemoryUtilization
	RequestRate
	QueueDepth
	ResponseTime
	ErrorRate
	CustomMetric
)

// String returns string representation of scaling metric
func (sm ScalingMetric) String() string {
	switch sm {
	case CPUUtilization:
		return "cpu_utilization"
	case MemoryUtilization:
		return "memory_utilization"
	case RequestRate:
		return "request_rate"
	case QueueDepth:
		return "queue_depth"
	case ResponseTime:
		return "response_time"
	case ErrorRate:
		return "error_rate"
	case CustomMetric:
		return "custom_metric"
	default:
		return "unknown"
	}
}

// MetricValue represents a metric measurement
type MetricValue struct {
	Type      ScalingMetric `json:"type"`
	Value     float64       `json:"value"`
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"` // instance ID or source identifier
}

// ScalingThreshold defines when to scale up or down
type ScalingThreshold struct {
	Metric             ScalingMetric `json:"metric"`
	ScaleUpThreshold   float64       `json:"scale_up_threshold"`
	ScaleDownThreshold float64       `json:"scale_down_threshold"`
	EvaluationPeriod   time.Duration `json:"evaluation_period"`
	Cooldown           time.Duration `json:"cooldown"`
	MinInstances       int           `json:"min_instances"`
	MaxInstances       int           `json:"max_instances"`
}

// ScalingDecision represents a scaling action
type ScalingDecision struct {
	Action          string        `json:"action"` // scale_up, scale_down, no_action
	CurrentSize     int           `json:"current_size"`
	TargetSize      int           `json:"target_size"`
	Reason          string        `json:"reason"`
	Metric          ScalingMetric `json:"metric"`
	MetricValue     float64       `json:"metric_value"`
	Timestamp       time.Time     `json:"timestamp"`
	ConfidenceScore float64       `json:"confidence_score"`
}

// MetricsCollector collects and aggregates metrics for scaling decisions
type MetricsCollector struct {
	metrics     map[ScalingMetric][]MetricValue
	mu          sync.RWMutex
	subscribers []chan<- MetricValue
	subMu       sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics:     make(map[ScalingMetric][]MetricValue),
		subscribers: make([]chan<- MetricValue, 0),
	}
}

// RecordMetric records a new metric value
func (mc *MetricsCollector) RecordMetric(metric MetricValue) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.metrics[metric.Type] == nil {
		mc.metrics[metric.Type] = make([]MetricValue, 0)
	}

	mc.metrics[metric.Type] = append(mc.metrics[metric.Type], metric)

	// Keep only last 1000 metrics per type
	if len(mc.metrics[metric.Type]) > 1000 {
		mc.metrics[metric.Type] = mc.metrics[metric.Type][1:]
	}

	// Notify subscribers
	mc.notifySubscribers(metric)
}

// notifySubscribers sends metric updates to all subscribers
func (mc *MetricsCollector) notifySubscribers(metric MetricValue) {
	mc.subMu.RLock()
	defer mc.subMu.RUnlock()

	for _, subscriber := range mc.subscribers {
		select {
		case subscriber <- metric:
		default:
			// Skip if channel is full (non-blocking)
		}
	}
}

// Subscribe adds a subscriber for metric updates
func (mc *MetricsCollector) Subscribe() <-chan MetricValue {
	mc.subMu.Lock()
	defer mc.subMu.Unlock()

	ch := make(chan MetricValue, 100) // Buffered channel
	mc.subscribers = append(mc.subscribers, ch)
	return ch
}

// GetAverageMetric calculates average metric over a time period
func (mc *MetricsCollector) GetAverageMetric(metricType ScalingMetric, period time.Duration) float64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics, exists := mc.metrics[metricType]
	if !exists || len(metrics) == 0 {
		return 0.0
	}

	cutoff := time.Now().Add(-period)
	var sum float64
	var count int

	for _, metric := range metrics {
		if metric.Timestamp.After(cutoff) {
			sum += metric.Value
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return sum / float64(count)
}

// GetMetricPercentile calculates percentile over a time period
func (mc *MetricsCollector) GetMetricPercentile(metricType ScalingMetric, period time.Duration, percentile float64) float64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics, exists := mc.metrics[metricType]
	if !exists || len(metrics) == 0 {
		return 0.0
	}

	cutoff := time.Now().Add(-period)
	var values []float64

	for _, metric := range metrics {
		if metric.Timestamp.After(cutoff) {
			values = append(values, metric.Value)
		}
	}

	if len(values) == 0 {
		return 0.0
	}

	// Sort values
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[i] > values[j] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}

	index := int(float64(len(values)) * percentile / 100.0)
	if index >= len(values) {
		index = len(values) - 1
	}

	return values[index]
}

// AutoScaler implements auto-scaling logic
type AutoScaler struct {
	thresholds       map[ScalingMetric]ScalingThreshold
	metricsCollector *MetricsCollector
	currentInstances int32
	lastScaleTime    time.Time
	scalingHistory   []ScalingDecision
	mu               sync.RWMutex
	enabled          bool
	predictiveModel  *PredictiveScaler
}

// NewAutoScaler creates a new auto-scaler
func NewAutoScaler(collector *MetricsCollector, initialInstances int) *AutoScaler {
	return &AutoScaler{
		thresholds:       make(map[ScalingMetric]ScalingThreshold),
		metricsCollector: collector,
		currentInstances: int32(initialInstances),
		lastScaleTime:    time.Now(),
		scalingHistory:   make([]ScalingDecision, 0),
		enabled:          true,
		predictiveModel:  NewPredictiveScaler(),
	}
}

// AddThreshold adds a scaling threshold
func (as *AutoScaler) AddThreshold(threshold ScalingThreshold) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.thresholds[threshold.Metric] = threshold
}

// EvaluateScaling evaluates whether scaling is needed
func (as *AutoScaler) EvaluateScaling() *ScalingDecision {
	as.mu.RLock()
	defer as.mu.RUnlock()

	if !as.enabled {
		return &ScalingDecision{
			Action:      "no_action",
			Reason:      "scaling disabled",
			Timestamp:   time.Now(),
			CurrentSize: int(atomic.LoadInt32(&as.currentInstances)),
			TargetSize:  int(atomic.LoadInt32(&as.currentInstances)),
		}
	}

	currentSize := int(atomic.LoadInt32(&as.currentInstances))

	// Check cooldown period
	for _, threshold := range as.thresholds {
		if time.Since(as.lastScaleTime) < threshold.Cooldown {
			return &ScalingDecision{
				Action:      "no_action",
				Reason:      "in cooldown period",
				Timestamp:   time.Now(),
				CurrentSize: currentSize,
				TargetSize:  currentSize,
			}
		}
	}

	// Evaluate each threshold
	for metricType, threshold := range as.thresholds {
		decision := as.evaluateThreshold(metricType, threshold)
		if decision.Action != "no_action" {
			return decision
		}
	}

	// Check predictive scaling
	if predictiveDecision := as.predictiveModel.PredictScaling(as.metricsCollector, currentSize); predictiveDecision != nil {
		return predictiveDecision
	}

	return &ScalingDecision{
		Action:      "no_action",
		Reason:      "all metrics within thresholds",
		Timestamp:   time.Now(),
		CurrentSize: currentSize,
		TargetSize:  currentSize,
	}
}

// evaluateThreshold evaluates a single threshold
func (as *AutoScaler) evaluateThreshold(metricType ScalingMetric, threshold ScalingThreshold) *ScalingDecision {
	currentValue := as.metricsCollector.GetAverageMetric(metricType, threshold.EvaluationPeriod)
	currentSize := int(atomic.LoadInt32(&as.currentInstances))

	decision := &ScalingDecision{
		Metric:      metricType,
		MetricValue: currentValue,
		Timestamp:   time.Now(),
		CurrentSize: currentSize,
	}

	// Scale up logic
	if currentValue > threshold.ScaleUpThreshold && currentSize < threshold.MaxInstances {
		targetSize := as.calculateTargetSize(currentSize, currentValue, threshold, true)
		decision.Action = "scale_up"
		decision.TargetSize = targetSize
		decision.Reason = fmt.Sprintf("%s (%.2f) exceeds scale-up threshold (%.2f)",
			metricType.String(), currentValue, threshold.ScaleUpThreshold)
		decision.ConfidenceScore = as.calculateConfidence(currentValue, threshold.ScaleUpThreshold, true)
		return decision
	}

	// Scale down logic
	if currentValue < threshold.ScaleDownThreshold && currentSize > threshold.MinInstances {
		targetSize := as.calculateTargetSize(currentSize, currentValue, threshold, false)
		decision.Action = "scale_down"
		decision.TargetSize = targetSize
		decision.Reason = fmt.Sprintf("%s (%.2f) below scale-down threshold (%.2f)",
			metricType.String(), currentValue, threshold.ScaleDownThreshold)
		decision.ConfidenceScore = as.calculateConfidence(currentValue, threshold.ScaleDownThreshold, false)
		return decision
	}

	decision.Action = "no_action"
	decision.TargetSize = currentSize
	decision.Reason = fmt.Sprintf("%s (%.2f) within thresholds", metricType.String(), currentValue)
	return decision
}

// calculateTargetSize calculates the target number of instances
func (as *AutoScaler) calculateTargetSize(currentSize int, metricValue float64, threshold ScalingThreshold, scaleUp bool) int {
	var targetSize int

	if scaleUp {
		// Scale up: increase by percentage based on how much threshold is exceeded
		overage := (metricValue - threshold.ScaleUpThreshold) / threshold.ScaleUpThreshold
		scalePercent := math.Min(overage, 1.0) // Cap at 100% increase
		increase := int(math.Ceil(float64(currentSize) * scalePercent))
		if increase < 1 {
			increase = 1 // Minimum increase of 1
		}
		targetSize = currentSize + increase

		if targetSize > threshold.MaxInstances {
			targetSize = threshold.MaxInstances
		}
	} else {
		// Scale down: decrease by percentage based on how much below threshold
		underage := (threshold.ScaleDownThreshold - metricValue) / threshold.ScaleDownThreshold
		scalePercent := math.Min(underage, 0.5) // Cap at 50% decrease
		decrease := int(math.Ceil(float64(currentSize) * scalePercent))
		if decrease < 1 {
			decrease = 1 // Minimum decrease of 1
		}
		targetSize = currentSize - decrease

		if targetSize < threshold.MinInstances {
			targetSize = threshold.MinInstances
		}
	}

	return targetSize
}

// calculateConfidence calculates confidence score for scaling decision
func (as *AutoScaler) calculateConfidence(metricValue, threshold float64, scaleUp bool) float64 {
	var ratio float64

	if scaleUp {
		if threshold == 0 {
			return 1.0
		}
		ratio = metricValue / threshold
		// More confidence the higher above threshold
		return math.Min((ratio-1.0)*2.0, 1.0)
	} else {
		if metricValue == 0 {
			return 1.0
		}
		ratio = threshold / metricValue
		// More confidence the lower below threshold
		return math.Min((ratio-1.0)*2.0, 1.0)
	}
}

// ExecuteScaling executes a scaling decision
func (as *AutoScaler) ExecuteScaling(decision *ScalingDecision) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if decision.Action == "no_action" {
		return nil
	}

	// Simulate scaling operation
	log.Printf("Executing scaling decision: %+v", decision)

	// Update instance count
	atomic.StoreInt32(&as.currentInstances, int32(decision.TargetSize))
	as.lastScaleTime = time.Now()

	// Record in history
	as.scalingHistory = append(as.scalingHistory, *decision)

	// Keep only last 100 scaling decisions
	if len(as.scalingHistory) > 100 {
		as.scalingHistory = as.scalingHistory[1:]
	}

	log.Printf("Scaling completed: instances changed from %d to %d",
		decision.CurrentSize, decision.TargetSize)

	return nil
}

// GetScalingHistory returns recent scaling history
func (as *AutoScaler) GetScalingHistory() []ScalingDecision {
	as.mu.RLock()
	defer as.mu.RUnlock()

	// Return copy to avoid race conditions
	history := make([]ScalingDecision, len(as.scalingHistory))
	copy(history, as.scalingHistory)
	return history
}

// PredictiveScaler implements predictive scaling using simple trend analysis
type PredictiveScaler struct {
	lookbackPeriod   time.Duration
	predictionWindow time.Duration
	enabled          bool
}

// NewPredictiveScaler creates a new predictive scaler
func NewPredictiveScaler() *PredictiveScaler {
	return &PredictiveScaler{
		lookbackPeriod:   30 * time.Minute,
		predictionWindow: 10 * time.Minute,
		enabled:          true,
	}
}

// PredictScaling predicts future scaling needs based on trends
func (ps *PredictiveScaler) PredictScaling(collector *MetricsCollector, currentInstances int) *ScalingDecision {
	if !ps.enabled {
		return nil
	}

	// Analyze CPU utilization trend
	cpuTrend := ps.calculateTrend(collector, CPUUtilization)

	// Predict future value
	currentCPU := collector.GetAverageMetric(CPUUtilization, 5*time.Minute)
	predictedCPU := currentCPU + (cpuTrend * ps.predictionWindow.Minutes())

	// Make scaling decision based on prediction
	if predictedCPU > 80.0 && cpuTrend > 1.0 { // Predicted high CPU with positive trend
		return &ScalingDecision{
			Action:          "scale_up",
			CurrentSize:     currentInstances,
			TargetSize:      currentInstances + 1,
			Reason:          fmt.Sprintf("predictive: CPU trending up (%.2f/min), predicted %.2f%%", cpuTrend, predictedCPU),
			Metric:          CPUUtilization,
			MetricValue:     currentCPU,
			Timestamp:       time.Now(),
			ConfidenceScore: math.Min(cpuTrend/5.0, 0.8), // Cap confidence for predictive
		}
	}

	if predictedCPU < 20.0 && cpuTrend < -1.0 { // Predicted low CPU with negative trend
		return &ScalingDecision{
			Action:          "scale_down",
			CurrentSize:     currentInstances,
			TargetSize:      currentInstances - 1,
			Reason:          fmt.Sprintf("predictive: CPU trending down (%.2f/min), predicted %.2f%%", cpuTrend, predictedCPU),
			Metric:          CPUUtilization,
			MetricValue:     currentCPU,
			Timestamp:       time.Now(),
			ConfidenceScore: math.Min(math.Abs(cpuTrend)/5.0, 0.8),
		}
	}

	return nil
}

// calculateTrend calculates the trend (slope) of a metric over time
func (ps *PredictiveScaler) calculateTrend(collector *MetricsCollector, metricType ScalingMetric) float64 {
	collector.mu.RLock()
	defer collector.mu.RUnlock()

	metrics, exists := collector.metrics[metricType]
	if !exists || len(metrics) < 2 {
		return 0.0
	}

	cutoff := time.Now().Add(-ps.lookbackPeriod)
	var dataPoints []struct {
		x float64 // time in minutes from start
		y float64 // metric value
	}

	var startTime time.Time
	for i, metric := range metrics {
		if metric.Timestamp.After(cutoff) {
			if i == 0 || startTime.IsZero() {
				startTime = metric.Timestamp
			}

			x := metric.Timestamp.Sub(startTime).Minutes()
			dataPoints = append(dataPoints, struct {
				x float64
				y float64
			}{x: x, y: metric.Value})
		}
	}

	if len(dataPoints) < 2 {
		return 0.0
	}

	// Calculate linear regression slope (simple trend)
	n := float64(len(dataPoints))
	var sumX, sumY, sumXY, sumX2 float64

	for _, point := range dataPoints {
		sumX += point.x
		sumY += point.y
		sumXY += point.x * point.y
		sumX2 += point.x * point.x
	}

	// Slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0.0
	}

	slope := (n*sumXY - sumX*sumY) / denominator
	return slope
}

// SystemMetricsProvider provides system metrics for auto-scaling
type SystemMetricsProvider struct {
	collector *MetricsCollector
	stopCh    chan struct{}
	running   bool
	mu        sync.Mutex
}

// NewSystemMetricsProvider creates a system metrics provider
func NewSystemMetricsProvider(collector *MetricsCollector) *SystemMetricsProvider {
	return &SystemMetricsProvider{
		collector: collector,
		stopCh:    make(chan struct{}),
	}
}

// Start begins collecting system metrics
func (smp *SystemMetricsProvider) Start(ctx context.Context) error {
	smp.mu.Lock()
	defer smp.mu.Unlock()

	if smp.running {
		return fmt.Errorf("metrics provider already running")
	}

	smp.running = true

	go smp.collectMetrics(ctx)
	return nil
}

// Stop stops collecting system metrics
func (smp *SystemMetricsProvider) Stop() {
	smp.mu.Lock()
	defer smp.mu.Unlock()

	if !smp.running {
		return
	}

	close(smp.stopCh)
	smp.running = false
}

// collectMetrics collects system metrics periodically
func (smp *SystemMetricsProvider) collectMetrics(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-smp.stopCh:
			return
		case <-ticker.C:
			smp.collectCurrentMetrics()
		}
	}
}

// collectCurrentMetrics collects current system metrics
func (smp *SystemMetricsProvider) collectCurrentMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// CPU utilization (simplified - in production, use proper CPU monitoring)
	cpuUsage := smp.simulateCPUUsage()

	// Memory utilization
	memUsage := float64(memStats.Alloc) / float64(memStats.Sys) * 100

	// Record metrics
	timestamp := time.Now()

	smp.collector.RecordMetric(MetricValue{
		Type:      CPUUtilization,
		Value:     cpuUsage,
		Timestamp: timestamp,
		Source:    "system",
	})

	smp.collector.RecordMetric(MetricValue{
		Type:      MemoryUtilization,
		Value:     memUsage,
		Timestamp: timestamp,
		Source:    "system",
	})
}

// simulateCPUUsage simulates CPU usage (in production, use actual CPU monitoring)
func (smp *SystemMetricsProvider) simulateCPUUsage() float64 {
	// Simulate varying CPU usage between 10-90%
	base := 30.0
	variation := 40.0 * math.Sin(float64(time.Now().Unix())/100.0)
	noise := float64(time.Now().Nanosecond()%1000) / 100.0

	cpu := base + variation + noise
	if cpu < 0 {
		cpu = 10.0
	}
	if cpu > 100 {
		cpu = 90.0
	}

	return cpu
}

// QueueBasedScaler scales based on queue depth
type QueueBasedScaler struct {
	queue            chan interface{}
	targetQueueDepth int
	mu               sync.RWMutex
}

// NewQueueBasedScaler creates a queue-based scaler
func NewQueueBasedScaler(queueSize, targetDepth int) *QueueBasedScaler {
	return &QueueBasedScaler{
		queue:            make(chan interface{}, queueSize),
		targetQueueDepth: targetDepth,
	}
}

// AddJob adds a job to the queue
func (qbs *QueueBasedScaler) AddJob(job interface{}) error {
	select {
	case qbs.queue <- job:
		return nil
	default:
		return fmt.Errorf("queue full")
	}
}

// GetQueueDepth returns current queue depth
func (qbs *QueueBasedScaler) GetQueueDepth() int {
	return len(qbs.queue)
}

// ProcessJobs processes jobs from the queue
func (qbs *QueueBasedScaler) ProcessJobs(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-qbs.queue:
			// Simulate job processing
			log.Printf("Worker %d processing job: %v", workerID, job)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ShouldScale determines if scaling is needed based on queue depth
func (qbs *QueueBasedScaler) ShouldScale() (bool, string) {
	depth := qbs.GetQueueDepth()

	if depth > qbs.targetQueueDepth*2 {
		return true, fmt.Sprintf("queue depth %d exceeds scale-up threshold %d", depth, qbs.targetQueueDepth*2)
	}

	if depth < qbs.targetQueueDepth/2 {
		return true, fmt.Sprintf("queue depth %d below scale-down threshold %d", depth, qbs.targetQueueDepth/2)
	}

	return false, "queue depth within target range"
}

// AutoScalingDemo demonstrates comprehensive auto-scaling patterns
func AutoScalingDemo() {
	log.Println("=== Auto-scaling Implementation Demo ===")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create metrics collector
	collector := NewMetricsCollector()

	// Create auto-scaler with initial 3 instances
	scaler := NewAutoScaler(collector, 3)

	// Configure thresholds
	cpuThreshold := ScalingThreshold{
		Metric:             CPUUtilization,
		ScaleUpThreshold:   70.0,
		ScaleDownThreshold: 30.0,
		EvaluationPeriod:   1 * time.Minute,
		Cooldown:           2 * time.Minute,
		MinInstances:       1,
		MaxInstances:       10,
	}
	scaler.AddThreshold(cpuThreshold)

	memoryThreshold := ScalingThreshold{
		Metric:             MemoryUtilization,
		ScaleUpThreshold:   80.0,
		ScaleDownThreshold: 40.0,
		EvaluationPeriod:   1 * time.Minute,
		Cooldown:           2 * time.Minute,
		MinInstances:       1,
		MaxInstances:       10,
	}
	scaler.AddThreshold(memoryThreshold)

	// Start system metrics provider
	metricsProvider := NewSystemMetricsProvider(collector)
	metricsProvider.Start(ctx)
	defer metricsProvider.Stop()

	// Subscribe to metric updates
	metricUpdates := collector.Subscribe()

	// Monitor metrics and make scaling decisions
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				decision := scaler.EvaluateScaling()
				log.Printf("Scaling decision: %+v", decision)

				if decision.Action != "no_action" {
					if err := scaler.ExecuteScaling(decision); err != nil {
						log.Printf("Error executing scaling: %v", err)
					}
				}
			}
		}
	}()

	// Demo queue-based scaling
	log.Println("\n--- Queue-based Scaling ---")
	queueScaler := NewQueueBasedScaler(100, 10)

	// Start workers
	numWorkers := 2
	for i := 0; i < numWorkers; i++ {
		go queueScaler.ProcessJobs(ctx, i)
	}

	// Add jobs to create queue pressure
	go func() {
		for i := 0; i < 50; i++ {
			if err := queueScaler.AddJob(fmt.Sprintf("job-%d", i)); err != nil {
				log.Printf("Failed to add job: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Monitor queue depth
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				shouldScale, reason := queueScaler.ShouldScale()
				log.Printf("Queue depth: %d, Should scale: %t, Reason: %s",
					queueScaler.GetQueueDepth(), shouldScale, reason)
			}
		}
	}()

	// Monitor metric updates
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case metric := <-metricUpdates:
				if metric.Type == CPUUtilization || metric.Type == MemoryUtilization {
					log.Printf("Metric update: %s = %.2f%% at %s",
						metric.Type.String(), metric.Value, metric.Timestamp.Format(time.RFC3339))
				}
			}
		}
	}()

	// Simulate load spikes
	go func() {
		time.Sleep(30 * time.Second)
		log.Println("Simulating high CPU load...")

		for i := 0; i < 10; i++ {
			collector.RecordMetric(MetricValue{
				Type:      CPUUtilization,
				Value:     85.0 + float64(i),
				Timestamp: time.Now(),
				Source:    "load_test",
			})
			time.Sleep(5 * time.Second)
		}

		log.Println("High load simulation complete")
	}()

	// Wait for demo completion
	<-ctx.Done()

	// Print final statistics
	log.Println("\n--- Final Statistics ---")
	log.Printf("Final instance count: %d", atomic.LoadInt32(&scaler.currentInstances))

	history := scaler.GetScalingHistory()
	log.Printf("Total scaling operations: %d", len(history))

	for _, decision := range history {
		log.Printf("  %s: %s -> %s (%s)",
			decision.Timestamp.Format("15:04:05"),
			decision.Action,
			fmt.Sprintf("%d->%d instances", decision.CurrentSize, decision.TargetSize),
			decision.Reason)
	}
}

// PerformanceBenchmark tests auto-scaling performance
func PerformanceBenchmark() {
	log.Println("\n=== Auto-scaling Performance Benchmark ===")

	collector := NewMetricsCollector()
	scaler := NewAutoScaler(collector, 5)

	// Add threshold
	threshold := ScalingThreshold{
		Metric:             CPUUtilization,
		ScaleUpThreshold:   70.0,
		ScaleDownThreshold: 30.0,
		EvaluationPeriod:   30 * time.Second,
		Cooldown:           1 * time.Minute,
		MinInstances:       1,
		MaxInstances:       100,
	}
	scaler.AddThreshold(threshold)

	// Benchmark scaling decision speed
	numEvaluations := 10000
	start := time.Now()

	for i := 0; i < numEvaluations; i++ {
		scaler.EvaluateScaling()
	}

	duration := time.Since(start)
	evaluationsPerSecond := float64(numEvaluations) / duration.Seconds()

	log.Printf("Scaling evaluations: %d in %v (%.2f evaluations/sec)",
		numEvaluations, duration, evaluationsPerSecond)

	// Benchmark metric collection
	numMetrics := 100000
	start = time.Now()

	for i := 0; i < numMetrics; i++ {
		collector.RecordMetric(MetricValue{
			Type:      CPUUtilization,
			Value:     50.0 + float64(i%50),
			Timestamp: time.Now(),
			Source:    "benchmark",
		})
	}

	duration = time.Since(start)
	metricsPerSecond := float64(numMetrics) / duration.Seconds()

	log.Printf("Metric collection: %d metrics in %v (%.2f metrics/sec)",
		numMetrics, duration, metricsPerSecond)
}

// FAANG Interview Discussion Points:

// 1. Auto-scaling Strategies:
//    - Reactive scaling: Based on current metrics
//    - Predictive scaling: Based on historical patterns
//    - Scheduled scaling: Based on known patterns
//    - Proactive scaling: Pre-emptive based on events

// 2. Scaling Metrics:
//    - CPU utilization: Most common, easy to measure
//    - Memory utilization: Important for memory-bound apps
//    - Request rate: Good for request-driven services
//    - Queue depth: Critical for async processing
//    - Custom business metrics: Domain-specific indicators

// 3. Scaling Challenges:
//    - Cold start time: Time for new instances to be ready
//    - Oscillation: Rapid scale up/down cycles
//    - Overshoot: Scaling too aggressively
//    - Undershoot: Not scaling enough
//    - Cost optimization vs performance

// 4. Cooldown Strategies:
//    - Prevent rapid oscillation
//    - Different cooldown for scale up vs scale down
//    - Adaptive cooldown based on confidence
//    - Emergency override for critical situations

// 5. Predictive Scaling:
//    - Time series analysis
//    - Machine learning models
//    - Seasonal pattern recognition
//    - Event-driven predictions

// Real-world Applications:
// - AWS Auto Scaling: CloudWatch metrics + scaling policies
// - Kubernetes HPA: CPU/memory/custom metrics scaling
// - Google Cloud Autoscaler: Multiple metric types
// - Azure VMSS: Predictive and reactive scaling

// Interview Questions to Discuss:
// 1. How do you prevent auto-scaling oscillation?
// 2. What metrics are most important for different application types?
// 3. How do you handle scaling during traffic spikes?
// 4. What are the trade-offs between reactive and predictive scaling?
// 5. How do you implement cost-effective auto-scaling policies?
