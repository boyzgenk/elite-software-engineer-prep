package production

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// ======================== CUSTOM METRICS REGISTRATION ========================

// MetricCollector manages custom metrics collection and aggregation
type MetricCollector struct {
	mu              sync.RWMutex
	counters        map[string]*AtomicCounter
	gauges          map[string]*AtomicGauge
	histograms      map[string]*ThreadSafeHistogram
	timers          map[string]*Timer
	businessMetrics map[string]*BusinessMetric
	aggregators     map[string]*MetricAggregator
	exporters       []MetricExporter
	samplingRate    float64
}

// AtomicCounter provides thread-safe counter operations
type AtomicCounter struct {
	value     int64
	name      string
	labels    map[string]string
	createdAt time.Time
}

// AtomicGauge provides thread-safe gauge operations
type AtomicGauge struct {
	value     int64 // Stored as int64 * 1000 for precision
	name      string
	labels    map[string]string
	createdAt time.Time
}

// ThreadSafeHistogram provides concurrent histogram operations
type ThreadSafeHistogram struct {
	mu        sync.Mutex
	buckets   []float64
	counts    []int64
	sum       float64
	count     int64
	name      string
	labels    map[string]string
	quantiles map[float64]float64 // Cached quantiles
}

// Timer measures duration and provides statistics
type Timer struct {
	mu         sync.RWMutex
	name       string
	labels     map[string]string
	durations  []time.Duration
	startTimes sync.Map // map[string]time.Time for ongoing measurements
}

// BusinessMetric tracks business-specific metrics
type BusinessMetric struct {
	mu         sync.RWMutex
	name       string
	metricType BusinessMetricType
	value      float64
	dimensions map[string]string
	timestamp  time.Time
	metadata   map[string]interface{}
}

// BusinessMetricType represents different business metric types
type BusinessMetricType int

const (
	BusinessMetricRevenue BusinessMetricType = iota
	BusinessMetricConversion
	BusinessMetricRetention
	BusinessMetricEngagement
	BusinessMetricPerformance
)

// NewMetricCollector creates a new metric collector
func NewMetricCollector(samplingRate float64) *MetricCollector {
	return &MetricCollector{
		counters:        make(map[string]*AtomicCounter),
		gauges:          make(map[string]*AtomicGauge),
		histograms:      make(map[string]*ThreadSafeHistogram),
		timers:          make(map[string]*Timer),
		businessMetrics: make(map[string]*BusinessMetric),
		aggregators:     make(map[string]*MetricAggregator),
		exporters:       make([]MetricExporter, 0),
		samplingRate:    samplingRate,
	}
}

// RegisterCounter creates and registers a new counter
func (mc *MetricCollector) RegisterCounter(name string, labels map[string]string) *AtomicCounter {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.generateKey(name, labels)
	if counter, exists := mc.counters[key]; exists {
		return counter
	}

	counter := &AtomicCounter{
		name:      name,
		labels:    labels,
		createdAt: time.Now(),
	}

	mc.counters[key] = counter
	return counter
}

// RegisterGauge creates and registers a new gauge
func (mc *MetricCollector) RegisterGauge(name string, labels map[string]string) *AtomicGauge {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.generateKey(name, labels)
	if gauge, exists := mc.gauges[key]; exists {
		return gauge
	}

	gauge := &AtomicGauge{
		name:      name,
		labels:    labels,
		createdAt: time.Now(),
	}

	mc.gauges[key] = gauge
	return gauge
}

// RegisterHistogram creates and registers a new histogram
func (mc *MetricCollector) RegisterHistogram(name string, buckets []float64, labels map[string]string) *ThreadSafeHistogram {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.generateKey(name, labels)
	if histogram, exists := mc.histograms[key]; exists {
		return histogram
	}

	// Sort buckets and ensure +Inf is included
	sortedBuckets := make([]float64, len(buckets))
	copy(sortedBuckets, buckets)
	sort.Float64s(sortedBuckets)
	if sortedBuckets[len(sortedBuckets)-1] != math.Inf(1) {
		sortedBuckets = append(sortedBuckets, math.Inf(1))
	}

	histogram := &ThreadSafeHistogram{
		buckets:   sortedBuckets,
		counts:    make([]int64, len(sortedBuckets)),
		name:      name,
		labels:    labels,
		quantiles: make(map[float64]float64),
	}

	mc.histograms[key] = histogram
	return histogram
}

// RegisterTimer creates and registers a new timer
func (mc *MetricCollector) RegisterTimer(name string, labels map[string]string) *Timer {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.generateKey(name, labels)
	if timer, exists := mc.timers[key]; exists {
		return timer
	}

	timer := &Timer{
		name:      name,
		labels:    labels,
		durations: make([]time.Duration, 0),
	}

	mc.timers[key] = timer
	return timer
}

// RegisterBusinessMetric creates and registers a business metric
func (mc *MetricCollector) RegisterBusinessMetric(name string, metricType BusinessMetricType, dimensions map[string]string) *BusinessMetric {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.generateKey(name, dimensions)
	if metric, exists := mc.businessMetrics[key]; exists {
		return metric
	}

	metric := &BusinessMetric{
		name:       name,
		metricType: metricType,
		dimensions: dimensions,
		metadata:   make(map[string]interface{}),
		timestamp:  time.Now(),
	}

	mc.businessMetrics[key] = metric
	return metric
}

func (mc *MetricCollector) generateKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += fmt.Sprintf(":%s=%s", k, v)
	}
	return key
}

// Counter operations
func (ac *AtomicCounter) Inc() {
	atomic.AddInt64(&ac.value, 1)
}

func (ac *AtomicCounter) Add(delta int64) {
	atomic.AddInt64(&ac.value, delta)
}

func (ac *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&ac.value)
}

// Gauge operations
func (ag *AtomicGauge) Set(value float64) {
	atomic.StoreInt64(&ag.value, int64(value*1000))
}

func (ag *AtomicGauge) Inc() {
	atomic.AddInt64(&ag.value, 1000)
}

func (ag *AtomicGauge) Dec() {
	atomic.AddInt64(&ag.value, -1000)
}

func (ag *AtomicGauge) Add(delta float64) {
	atomic.AddInt64(&ag.value, int64(delta*1000))
}

func (ag *AtomicGauge) Get() float64 {
	return float64(atomic.LoadInt64(&ag.value)) / 1000
}

// Histogram operations
func (tsh *ThreadSafeHistogram) Observe(value float64) {
	tsh.mu.Lock()
	defer tsh.mu.Unlock()

	tsh.sum += value
	tsh.count++

	// Update bucket counts
	for i, bucket := range tsh.buckets {
		if value <= bucket {
			tsh.counts[i]++
		}
	}

	// Clear cached quantiles
	tsh.quantiles = make(map[float64]float64)
}

func (tsh *ThreadSafeHistogram) GetQuantile(q float64) float64 {
	tsh.mu.Lock()
	defer tsh.mu.Unlock()

	// Check cache
	if cached, exists := tsh.quantiles[q]; exists {
		return cached
	}

	if tsh.count == 0 {
		return 0
	}

	targetCount := q * float64(tsh.count)
	cumulativeCount := int64(0)

	for i, count := range tsh.counts {
		cumulativeCount += count
		if float64(cumulativeCount) >= targetCount {
			var result float64
			if i == 0 {
				result = tsh.buckets[0]
			} else {
				// Linear interpolation
				prevBucket := float64(0)
				if i > 0 && tsh.buckets[i-1] != math.Inf(-1) {
					prevBucket = tsh.buckets[i-1]
				}

				if tsh.buckets[i] == math.Inf(1) {
					result = prevBucket
				} else {
					result = prevBucket + (tsh.buckets[i]-prevBucket)*0.5
				}
			}

			// Cache result
			tsh.quantiles[q] = result
			return result
		}
	}

	return tsh.buckets[len(tsh.buckets)-2]
}

func (tsh *ThreadSafeHistogram) GetSum() float64 {
	tsh.mu.Lock()
	defer tsh.mu.Unlock()
	return tsh.sum
}

func (tsh *ThreadSafeHistogram) GetCount() int64 {
	tsh.mu.Lock()
	defer tsh.mu.Unlock()
	return tsh.count
}

// Timer operations
func (t *Timer) Start(operationID string) {
	t.startTimes.Store(operationID, time.Now())
}

func (t *Timer) Stop(operationID string) time.Duration {
	if startTime, ok := t.startTimes.LoadAndDelete(operationID); ok {
		duration := time.Since(startTime.(time.Time))

		t.mu.Lock()
		t.durations = append(t.durations, duration)

		// Keep only recent measurements (last 1000)
		if len(t.durations) > 1000 {
			t.durations = t.durations[len(t.durations)-1000:]
		}
		t.mu.Unlock()

		return duration
	}
	return 0
}

func (t *Timer) Record(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.durations = append(t.durations, duration)

	// Keep only recent measurements
	if len(t.durations) > 1000 {
		t.durations = t.durations[len(t.durations)-1000:]
	}
}

func (t *Timer) GetStats() TimerStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.durations) == 0 {
		return TimerStats{}
	}

	// Copy and sort durations for percentile calculation
	sorted := make([]time.Duration, len(t.durations))
	copy(sorted, t.durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	var sum time.Duration
	for _, d := range sorted {
		sum += d
	}

	return TimerStats{
		Count: int64(len(sorted)),
		Sum:   sum,
		Min:   sorted[0],
		Max:   sorted[len(sorted)-1],
		Mean:  sum / time.Duration(len(sorted)),
		P50:   sorted[len(sorted)*50/100],
		P90:   sorted[len(sorted)*90/100],
		P95:   sorted[len(sorted)*95/100],
		P99:   sorted[len(sorted)*99/100],
	}
}

// TimerStats contains timer statistics
type TimerStats struct {
	Count int64         `json:"count"`
	Sum   time.Duration `json:"sum"`
	Min   time.Duration `json:"min"`
	Max   time.Duration `json:"max"`
	Mean  time.Duration `json:"mean"`
	P50   time.Duration `json:"p50"`
	P90   time.Duration `json:"p90"`
	P95   time.Duration `json:"p95"`
	P99   time.Duration `json:"p99"`
}

// Business metric operations
func (bm *BusinessMetric) Set(value float64) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.value = value
	bm.timestamp = time.Now()
}

func (bm *BusinessMetric) SetWithMetadata(value float64, metadata map[string]interface{}) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.value = value
	bm.metadata = metadata
	bm.timestamp = time.Now()
}

func (bm *BusinessMetric) Get() (float64, map[string]interface{}) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	return bm.value, bm.metadata
}

// ======================== PERFORMANCE MONITORING ========================

// PerformanceMonitor tracks latency, throughput, and errors
type PerformanceMonitor struct {
	collector       *MetricCollector
	requestCounter  *AtomicCounter
	errorCounter    *AtomicCounter
	latencyHist     *ThreadSafeHistogram
	throughputGauge *AtomicGauge

	// Golden signals tracking
	latencyP99  *AtomicGauge
	latencyP95  *AtomicGauge
	latencyMean *AtomicGauge
	errorRate   *AtomicGauge

	// Window-based metrics
	windowSize time.Duration
	windowData *TimeWindowData
}

// TimeWindowData stores metrics within time windows
type TimeWindowData struct {
	mu            sync.RWMutex
	requestCounts []TimedCount
	errorCounts   []TimedCount
	latencies     []TimedLatency
}

// TimedCount represents a count with timestamp
type TimedCount struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

// TimedLatency represents latency measurement with timestamp
type TimedLatency struct {
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(collector *MetricCollector, windowSize time.Duration) *PerformanceMonitor {
	labels := map[string]string{"service": "api"}

	latencyBuckets := []float64{
		0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
	}

	pm := &PerformanceMonitor{
		collector:       collector,
		requestCounter:  collector.RegisterCounter("http_requests_total", labels),
		errorCounter:    collector.RegisterCounter("http_errors_total", labels),
		latencyHist:     collector.RegisterHistogram("http_request_duration_seconds", latencyBuckets, labels),
		throughputGauge: collector.RegisterGauge("http_requests_per_second", labels),
		latencyP99:      collector.RegisterGauge("http_latency_p99", labels),
		latencyP95:      collector.RegisterGauge("http_latency_p95", labels),
		latencyMean:     collector.RegisterGauge("http_latency_mean", labels),
		errorRate:       collector.RegisterGauge("http_error_rate", labels),
		windowSize:      windowSize,
		windowData: &TimeWindowData{
			requestCounts: make([]TimedCount, 0),
			errorCounts:   make([]TimedCount, 0),
			latencies:     make([]TimedLatency, 0),
		},
	}

	// Start background metrics calculation
	go pm.calculateMetrics()

	return pm
}

// RecordRequest records a successful request
func (pm *PerformanceMonitor) RecordRequest(duration time.Duration) {
	pm.requestCounter.Inc()
	pm.latencyHist.Observe(duration.Seconds())

	// Add to window data
	pm.windowData.mu.Lock()
	pm.windowData.requestCounts = append(pm.windowData.requestCounts, TimedCount{
		Timestamp: time.Now(),
		Count:     1,
	})
	pm.windowData.latencies = append(pm.windowData.latencies, TimedLatency{
		Timestamp: time.Now(),
		Duration:  duration,
	})
	pm.windowData.mu.Unlock()
}

// RecordError records an error
func (pm *PerformanceMonitor) RecordError(duration time.Duration) {
	pm.requestCounter.Inc()
	pm.errorCounter.Inc()
	pm.latencyHist.Observe(duration.Seconds())

	// Add to window data
	pm.windowData.mu.Lock()
	pm.windowData.requestCounts = append(pm.windowData.requestCounts, TimedCount{
		Timestamp: time.Now(),
		Count:     1,
	})
	pm.windowData.errorCounts = append(pm.windowData.errorCounts, TimedCount{
		Timestamp: time.Now(),
		Count:     1,
	})
	pm.windowData.latencies = append(pm.windowData.latencies, TimedLatency{
		Timestamp: time.Now(),
		Duration:  duration,
	})
	pm.windowData.mu.Unlock()
}

// calculateMetrics runs in background to compute derived metrics
func (pm *PerformanceMonitor) calculateMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pm.updateGoldenSignals()
		pm.cleanupOldData()
	}
}

// updateGoldenSignals calculates the four golden signals
func (pm *PerformanceMonitor) updateGoldenSignals() {
	// Update latency percentiles
	pm.latencyP99.Set(pm.latencyHist.GetQuantile(0.99))
	pm.latencyP95.Set(pm.latencyHist.GetQuantile(0.95))
	pm.latencyP99.Set(pm.latencyHist.GetQuantile(0.50)) // Using P99 gauge for mean temporarily

	// Calculate error rate
	totalRequests := pm.requestCounter.Get()
	totalErrors := pm.errorCounter.Get()

	if totalRequests > 0 {
		errorRate := float64(totalErrors) / float64(totalRequests)
		pm.errorRate.Set(errorRate)
	}

	// Calculate throughput (requests per second in current window)
	pm.windowData.mu.RLock()
	cutoff := time.Now().Add(-pm.windowSize)

	var requestsInWindow int64
	for _, count := range pm.windowData.requestCounts {
		if count.Timestamp.After(cutoff) {
			requestsInWindow += count.Count
		}
	}
	pm.windowData.mu.RUnlock()

	throughput := float64(requestsInWindow) / pm.windowSize.Seconds()
	pm.throughputGauge.Set(throughput)
}

// cleanupOldData removes data outside the window
func (pm *PerformanceMonitor) cleanupOldData() {
	pm.windowData.mu.Lock()
	defer pm.windowData.mu.Unlock()

	cutoff := time.Now().Add(-pm.windowSize)

	// Clean request counts
	var newRequestCounts []TimedCount
	for _, count := range pm.windowData.requestCounts {
		if count.Timestamp.After(cutoff) {
			newRequestCounts = append(newRequestCounts, count)
		}
	}
	pm.windowData.requestCounts = newRequestCounts

	// Clean error counts
	var newErrorCounts []TimedCount
	for _, count := range pm.windowData.errorCounts {
		if count.Timestamp.After(cutoff) {
			newErrorCounts = append(newErrorCounts, count)
		}
	}
	pm.windowData.errorCounts = newErrorCounts

	// Clean latencies
	var newLatencies []TimedLatency
	for _, latency := range pm.windowData.latencies {
		if latency.Timestamp.After(cutoff) {
			newLatencies = append(newLatencies, latency)
		}
	}
	pm.windowData.latencies = newLatencies
}

// GetPerformanceReport generates a comprehensive performance report
func (pm *PerformanceMonitor) GetPerformanceReport() PerformanceReport {
	pm.windowData.mu.RLock()
	defer pm.windowData.mu.RUnlock()

	report := PerformanceReport{
		Timestamp:        time.Now(),
		TotalRequests:    pm.requestCounter.Get(),
		TotalErrors:      pm.errorCounter.Get(),
		ErrorRate:        pm.errorRate.Get(),
		Throughput:       pm.throughputGauge.Get(),
		LatencyP50:       pm.latencyHist.GetQuantile(0.50),
		LatencyP95:       pm.latencyHist.GetQuantile(0.95),
		LatencyP99:       pm.latencyHist.GetQuantile(0.99),
		WindowSize:       pm.windowSize,
		RequestsInWindow: int64(len(pm.windowData.requestCounts)),
		ErrorsInWindow:   int64(len(pm.windowData.errorCounts)),
	}

	return report
}

// PerformanceReport contains performance metrics summary
type PerformanceReport struct {
	Timestamp        time.Time     `json:"timestamp"`
	TotalRequests    int64         `json:"total_requests"`
	TotalErrors      int64         `json:"total_errors"`
	ErrorRate        float64       `json:"error_rate"`
	Throughput       float64       `json:"throughput"`
	LatencyP50       float64       `json:"latency_p50"`
	LatencyP95       float64       `json:"latency_p95"`
	LatencyP99       float64       `json:"latency_p99"`
	WindowSize       time.Duration `json:"window_size"`
	RequestsInWindow int64         `json:"requests_in_window"`
	ErrorsInWindow   int64         `json:"errors_in_window"`
}

// ======================== RESOURCE UTILIZATION MONITORING ========================

// ResourceMonitor tracks system resource utilization
type ResourceMonitor struct {
	collector       *MetricCollector
	cpuUsage        *AtomicGauge
	memoryUsage     *AtomicGauge
	memoryAllocated *AtomicGauge
	goroutineCount  *AtomicGauge
	gcPauseTime     *ThreadSafeHistogram

	// Custom resource metrics
	diskUsage       *AtomicGauge
	networkIn       *AtomicCounter
	networkOut      *AtomicCounter
	fileDescriptors *AtomicGauge

	updateInterval time.Duration
	stopCh         chan struct{}
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor(collector *MetricCollector, updateInterval time.Duration) *ResourceMonitor {
	labels := map[string]string{"instance": "local"}

	rm := &ResourceMonitor{
		collector:       collector,
		cpuUsage:        collector.RegisterGauge("cpu_usage_percent", labels),
		memoryUsage:     collector.RegisterGauge("memory_usage_bytes", labels),
		memoryAllocated: collector.RegisterGauge("memory_allocated_bytes", labels),
		goroutineCount:  collector.RegisterGauge("goroutines_count", labels),
		gcPauseTime:     collector.RegisterHistogram("gc_pause_seconds", []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5}, labels),
		diskUsage:       collector.RegisterGauge("disk_usage_bytes", labels),
		networkIn:       collector.RegisterCounter("network_bytes_in_total", labels),
		networkOut:      collector.RegisterCounter("network_bytes_out_total", labels),
		fileDescriptors: collector.RegisterGauge("file_descriptors_count", labels),
		updateInterval:  updateInterval,
		stopCh:          make(chan struct{}),
	}

	go rm.monitor()
	return rm
}

// monitor runs the resource monitoring loop
func (rm *ResourceMonitor) monitor() {
	ticker := time.NewTicker(rm.updateInterval)
	defer ticker.Stop()

	var lastMemStats runtime.MemStats
	runtime.ReadMemStats(&lastMemStats)

	for {
		select {
		case <-ticker.C:
			rm.updateMetrics(&lastMemStats)
		case <-rm.stopCh:
			return
		}
	}
}

// updateMetrics updates all resource metrics
func (rm *ResourceMonitor) updateMetrics(lastMemStats *runtime.MemStats) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Memory metrics
	rm.memoryUsage.Set(float64(memStats.Sys))
	rm.memoryAllocated.Set(float64(memStats.Alloc))

	// Goroutine count
	rm.goroutineCount.Set(float64(runtime.NumGoroutine()))

	// GC pause time (if there was a GC since last check)
	if memStats.NumGC > lastMemStats.NumGC {
		// Calculate average pause time for recent GCs
		var totalPause time.Duration
		var gcCount uint32

		for i := uint32(0); i < memStats.NumGC-lastMemStats.NumGC && i < 256; i++ {
			idx := (memStats.NumGC - 1 - i) % 256
			totalPause += time.Duration(memStats.PauseNs[idx])
			gcCount++
		}

		if gcCount > 0 {
			avgPause := totalPause / time.Duration(gcCount)
			rm.gcPauseTime.Observe(avgPause.Seconds())
		}
	}

	// Update last stats
	*lastMemStats = memStats
}

// Stop stops the resource monitor
func (rm *ResourceMonitor) Stop() {
	close(rm.stopCh)
}

// GetResourceReport generates a resource utilization report
func (rm *ResourceMonitor) GetResourceReport() ResourceReport {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return ResourceReport{
		Timestamp:       time.Now(),
		CPUUsage:        rm.cpuUsage.Get(),
		MemoryUsage:     rm.memoryUsage.Get(),
		MemoryAllocated: rm.memoryAllocated.Get(),
		GoroutineCount:  int(rm.goroutineCount.Get()),
		GCPauseP95:      rm.gcPauseTime.GetQuantile(0.95),
		GCPauseP99:      rm.gcPauseTime.GetQuantile(0.99),
		NetworkBytesIn:  rm.networkIn.Get(),
		NetworkBytesOut: rm.networkOut.Get(),
		FileDescriptors: int(rm.fileDescriptors.Get()),

		// Additional runtime stats
		HeapObjects: memStats.HeapObjects,
		StackInUse:  memStats.StackInuse,
		GCCycles:    memStats.NumGC,
		NextGC:      memStats.NextGC,
	}
}

// ResourceReport contains resource utilization metrics
type ResourceReport struct {
	Timestamp       time.Time `json:"timestamp"`
	CPUUsage        float64   `json:"cpu_usage_percent"`
	MemoryUsage     float64   `json:"memory_usage_bytes"`
	MemoryAllocated float64   `json:"memory_allocated_bytes"`
	GoroutineCount  int       `json:"goroutine_count"`
	GCPauseP95      float64   `json:"gc_pause_p95_seconds"`
	GCPauseP99      float64   `json:"gc_pause_p99_seconds"`
	NetworkBytesIn  int64     `json:"network_bytes_in_total"`
	NetworkBytesOut int64     `json:"network_bytes_out_total"`
	FileDescriptors int       `json:"file_descriptors_count"`
	HeapObjects     uint64    `json:"heap_objects"`
	StackInUse      uint64    `json:"stack_in_use"`
	GCCycles        uint32    `json:"gc_cycles"`
	NextGC          uint64    `json:"next_gc"`
}

// ======================== REAL-TIME METRICS AGGREGATION ========================

// MetricAggregator performs real-time aggregation of metrics
type MetricAggregator struct {
	name            string
	aggregationType AggregationType
	windowSize      time.Duration
	values          []TimestampedValue
	result          float64
	mu              sync.RWMutex
}

// AggregationType represents different aggregation methods
type AggregationType int

const (
	AggregationSum AggregationType = iota
	AggregationAvg
	AggregationMin
	AggregationMax
	AggregationCount
	AggregationRate
	AggregationPercentile
)

// TimestampedValue represents a value with timestamp
type TimestampedValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// NewMetricAggregator creates a new metric aggregator
func NewMetricAggregator(name string, aggType AggregationType, windowSize time.Duration) *MetricAggregator {
	ma := &MetricAggregator{
		name:            name,
		aggregationType: aggType,
		windowSize:      windowSize,
		values:          make([]TimestampedValue, 0),
	}

	go ma.aggregate()
	return ma
}

// AddValue adds a new value to the aggregator
func (ma *MetricAggregator) AddValue(value float64) {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	ma.values = append(ma.values, TimestampedValue{
		Timestamp: time.Now(),
		Value:     value,
	})
}

// aggregate runs the aggregation loop
func (ma *MetricAggregator) aggregate() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ma.mu.Lock()

		// Remove old values
		cutoff := time.Now().Add(-ma.windowSize)
		var recentValues []TimestampedValue

		for _, v := range ma.values {
			if v.Timestamp.After(cutoff) {
				recentValues = append(recentValues, v)
			}
		}

		ma.values = recentValues

		// Calculate aggregation
		if len(ma.values) > 0 {
			ma.result = ma.calculate(ma.values)
		}

		ma.mu.Unlock()
	}
}

// calculate performs the actual aggregation calculation
func (ma *MetricAggregator) calculate(values []TimestampedValue) float64 {
	if len(values) == 0 {
		return 0
	}

	switch ma.aggregationType {
	case AggregationSum:
		var sum float64
		for _, v := range values {
			sum += v.Value
		}
		return sum

	case AggregationAvg:
		var sum float64
		for _, v := range values {
			sum += v.Value
		}
		return sum / float64(len(values))

	case AggregationMin:
		min := values[0].Value
		for _, v := range values[1:] {
			if v.Value < min {
				min = v.Value
			}
		}
		return min

	case AggregationMax:
		max := values[0].Value
		for _, v := range values[1:] {
			if v.Value > max {
				max = v.Value
			}
		}
		return max

	case AggregationCount:
		return float64(len(values))

	case AggregationRate:
		if len(values) < 2 {
			return 0
		}

		oldest := values[0]
		newest := values[len(values)-1]
		timeDiff := newest.Timestamp.Sub(oldest.Timestamp).Seconds()

		if timeDiff > 0 {
			return (newest.Value - oldest.Value) / timeDiff
		}
		return 0

	case AggregationPercentile:
		// Calculate 95th percentile
		sortedValues := make([]float64, len(values))
		for i, v := range values {
			sortedValues[i] = v.Value
		}
		sort.Float64s(sortedValues)

		index := int(float64(len(sortedValues)) * 0.95)
		if index >= len(sortedValues) {
			index = len(sortedValues) - 1
		}
		return sortedValues[index]

	default:
		return 0
	}
}

// GetResult returns the current aggregated result
func (ma *MetricAggregator) GetResult() float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()
	return ma.result
}

// ======================== TIME-SERIES DATA HANDLING ========================

// TimeSeriesDB manages time-series metric data
type TimeSeriesDB struct {
	mu        sync.RWMutex
	series    map[string]*TimeSeries
	retention time.Duration
}

// TimeSeries represents a single time series
type TimeSeries struct {
	Name     string            `json:"name"`
	Labels   map[string]string `json:"labels"`
	Points   []TimeSeriesPoint `json:"points"`
	metadata map[string]interface{}
}

// TimeSeriesPoint represents a single data point
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// NewTimeSeriesDB creates a new time-series database
func NewTimeSeriesDB(retention time.Duration) *TimeSeriesDB {
	tsdb := &TimeSeriesDB{
		series:    make(map[string]*TimeSeries),
		retention: retention,
	}

	go tsdb.cleanup()
	return tsdb
}

// AddPoint adds a new data point to a time series
func (tsdb *TimeSeriesDB) AddPoint(seriesName string, labels map[string]string, value float64) {
	tsdb.mu.Lock()
	defer tsdb.mu.Unlock()

	key := tsdb.generateKey(seriesName, labels)

	series, exists := tsdb.series[key]
	if !exists {
		series = &TimeSeries{
			Name:     seriesName,
			Labels:   labels,
			Points:   make([]TimeSeriesPoint, 0),
			metadata: make(map[string]interface{}),
		}
		tsdb.series[key] = series
	}

	series.Points = append(series.Points, TimeSeriesPoint{
		Timestamp: time.Now(),
		Value:     value,
	})
}

// QueryRange queries time series data within a time range
func (tsdb *TimeSeriesDB) QueryRange(seriesName string, labels map[string]string, start, end time.Time) []TimeSeriesPoint {
	tsdb.mu.RLock()
	defer tsdb.mu.RUnlock()

	key := tsdb.generateKey(seriesName, labels)
	series, exists := tsdb.series[key]

	if !exists {
		return nil
	}

	var result []TimeSeriesPoint
	for _, point := range series.Points {
		if point.Timestamp.After(start) && point.Timestamp.Before(end) {
			result = append(result, point)
		}
	}

	return result
}

// QueryLatest returns the latest N points for a series
func (tsdb *TimeSeriesDB) QueryLatest(seriesName string, labels map[string]string, count int) []TimeSeriesPoint {
	tsdb.mu.RLock()
	defer tsdb.mu.RUnlock()

	key := tsdb.generateKey(seriesName, labels)
	series, exists := tsdb.series[key]

	if !exists {
		return nil
	}

	points := series.Points
	if len(points) <= count {
		return points
	}

	return points[len(points)-count:]
}

// cleanup removes old data points beyond retention period
func (tsdb *TimeSeriesDB) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		tsdb.mu.Lock()

		cutoff := time.Now().Add(-tsdb.retention)

		for key, series := range tsdb.series {
			var recentPoints []TimeSeriesPoint

			for _, point := range series.Points {
				if point.Timestamp.After(cutoff) {
					recentPoints = append(recentPoints, point)
				}
			}

			if len(recentPoints) == 0 {
				delete(tsdb.series, key)
			} else {
				series.Points = recentPoints
			}
		}

		tsdb.mu.Unlock()
	}
}

func (tsdb *TimeSeriesDB) generateKey(seriesName string, labels map[string]string) string {
	key := seriesName
	for k, v := range labels {
		key += fmt.Sprintf(":%s=%s", k, v)
	}
	return key
}

// ======================== METRIC EXPORT INTERFACES ========================

// MetricExporter interface for exporting metrics to external systems
type MetricExporter interface {
	Export(metrics map[string]interface{}) error
	GetName() string
}

// PrometheusExporter exports metrics in Prometheus format
type PrometheusExporter struct {
	endpoint string
}

// NewPrometheusExporter creates a new Prometheus exporter
func NewPrometheusExporter(endpoint string) *PrometheusExporter {
	return &PrometheusExporter{endpoint: endpoint}
}

func (pe *PrometheusExporter) Export(metrics map[string]interface{}) error {
	// Implementation would format metrics in Prometheus format and send to endpoint
	// For now, just log the export
	fmt.Printf("Exporting %d metrics to Prometheus at %s\n", len(metrics), pe.endpoint)
	return nil
}

func (pe *PrometheusExporter) GetName() string {
	return "prometheus"
}

// StatsDExporter exports metrics to StatsD
type StatsDExporter struct {
	address string
}

// NewStatsDExporter creates a new StatsD exporter
func NewStatsDExporter(address string) *StatsDExporter {
	return &StatsDExporter{address: address}
}

func (se *StatsDExporter) Export(metrics map[string]interface{}) error {
	fmt.Printf("Exporting %d metrics to StatsD at %s\n", len(metrics), se.address)
	return nil
}

func (se *StatsDExporter) GetName() string {
	return "statsd"
}

// ======================== HTTP HANDLERS FOR METRICS COLLECTION ========================

// MetricsHandler provides HTTP endpoints for metrics
type MetricsHandler struct {
	collector       *MetricCollector
	perfMonitor     *PerformanceMonitor
	resourceMonitor *ResourceMonitor
	timeSeriesDB    *TimeSeriesDB
}

// NewMetricsHandler creates HTTP handlers for metrics collection
func NewMetricsHandler() *MetricsHandler {
	collector := NewMetricCollector(1.0)
	perfMonitor := NewPerformanceMonitor(collector, 5*time.Minute)
	resourceMonitor := NewResourceMonitor(collector, 10*time.Second)
	timeSeriesDB := NewTimeSeriesDB(24 * time.Hour)

	return &MetricsHandler{
		collector:       collector,
		perfMonitor:     perfMonitor,
		resourceMonitor: resourceMonitor,
		timeSeriesDB:    timeSeriesDB,
	}
}

// HandleMetricsEndpoint serves all metrics
func (mh *MetricsHandler) HandleMetricsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"performance": mh.perfMonitor.GetPerformanceReport(),
		"resources":   mh.resourceMonitor.GetResourceReport(),
		"timestamp":   time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// HandleBusinessMetrics serves business metrics
func (mh *MetricsHandler) HandleBusinessMetrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mh.collector.mu.RLock()
		businessMetrics := make(map[string]interface{})
		for key, metric := range mh.collector.businessMetrics {
			value, metadata := metric.Get()
			businessMetrics[key] = map[string]interface{}{
				"value":     value,
				"metadata":  metadata,
				"timestamp": metric.timestamp,
			}
		}
		mh.collector.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(businessMetrics)

	case http.MethodPost:
		var req struct {
			Name       string                 `json:"name"`
			Type       string                 `json:"type"`
			Value      float64                `json:"value"`
			Dimensions map[string]string      `json:"dimensions"`
			Metadata   map[string]interface{} `json:"metadata"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var metricType BusinessMetricType
		switch req.Type {
		case "revenue":
			metricType = BusinessMetricRevenue
		case "conversion":
			metricType = BusinessMetricConversion
		case "retention":
			metricType = BusinessMetricRetention
		case "engagement":
			metricType = BusinessMetricEngagement
		default:
			metricType = BusinessMetricPerformance
		}

		metric := mh.collector.RegisterBusinessMetric(req.Name, metricType, req.Dimensions)
		metric.SetWithMetadata(req.Value, req.Metadata)

		// Also add to time series
		mh.timeSeriesDB.AddPoint(req.Name, req.Dimensions, req.Value)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "metric recorded"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleTimeSeriesQuery serves time series queries
func (mh *MetricsHandler) HandleTimeSeriesQuery(w http.ResponseWriter, r *http.Request) {
	seriesName := r.URL.Query().Get("series")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if seriesName == "" {
		http.Error(w, "series parameter required", http.StatusBadRequest)
		return
	}

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "invalid start time format", http.StatusBadRequest)
			return
		}
	} else {
		start = time.Now().Add(-1 * time.Hour)
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "invalid end time format", http.StatusBadRequest)
			return
		}
	} else {
		end = time.Now()
	}

	// Parse labels from query parameters
	labels := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "series" && key != "start" && key != "end" && len(values) > 0 {
			labels[key] = values[0]
		}
	}

	points := mh.timeSeriesDB.QueryRange(seriesName, labels, start, end)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"series": seriesName,
		"labels": labels,
		"points": points,
		"start":  start,
		"end":    end,
	})
}

// RecordHTTPRequest is middleware to record HTTP request metrics
func (mh *MetricsHandler) RecordHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapper := &ResponseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		// Record metrics
		if wrapper.statusCode >= 400 {
			mh.perfMonitor.RecordError(duration)
		} else {
			mh.perfMonitor.RecordRequest(duration)
		}

		// Add to time series
		mh.timeSeriesDB.AddPoint("http_requests", map[string]string{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": strconv.Itoa(wrapper.statusCode),
		}, 1)

		mh.timeSeriesDB.AddPoint("http_duration", map[string]string{
			"method": r.Method,
			"path":   r.URL.Path,
		}, duration.Seconds())
	})
}

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
