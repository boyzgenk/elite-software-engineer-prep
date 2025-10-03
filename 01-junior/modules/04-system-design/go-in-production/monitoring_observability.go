package production

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ======================== METRICS COLLECTION (PROMETHEUS-STYLE) ========================

// MetricType represents different types of metrics
type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
	MetricTypeSummary
)

// Metric represents a single metric data point
type Metric struct {
	Name      string            `json:"name"`
	Type      MetricType        `json:"type"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
	Help      string            `json:"help"`
}

// MetricRegistry manages metric collection and exposition
type MetricRegistry struct {
	mu      sync.RWMutex
	metrics map[string]*MetricFamily
	prefix  string
}

// MetricFamily groups related metrics
type MetricFamily struct {
	Name    string     `json:"name"`
	Type    MetricType `json:"type"`
	Help    string     `json:"help"`
	Metrics []*Metric  `json:"metrics"`
	mu      sync.RWMutex
}

// NewMetricRegistry creates a new metric registry
func NewMetricRegistry(prefix string) *MetricRegistry {
	return &MetricRegistry{
		metrics: make(map[string]*MetricFamily),
		prefix:  prefix,
	}
}

// Counter implements a Prometheus-style counter
type Counter struct {
	value  *int64
	family *MetricFamily
	labels map[string]string
}

// NewCounter creates a new counter metric
func (mr *MetricRegistry) NewCounter(name, help string, labels map[string]string) *Counter {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	fullName := mr.prefix + name
	family, exists := mr.metrics[fullName]
	if !exists {
		family = &MetricFamily{
			Name:    fullName,
			Type:    MetricTypeCounter,
			Help:    help,
			Metrics: make([]*Metric, 0),
		}
		mr.metrics[fullName] = family
	}

	counter := &Counter{
		value:  new(int64),
		family: family,
		labels: labels,
	}

	return counter
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	atomic.AddInt64(c.value, 1)
	c.updateMetric()
}

// Add adds the given value to the counter
func (c *Counter) Add(delta float64) {
	if delta < 0 {
		panic("counter values cannot decrease")
	}
	atomic.AddInt64(c.value, int64(delta))
	c.updateMetric()
}

// Get returns the current counter value
func (c *Counter) Get() float64 {
	return float64(atomic.LoadInt64(c.value))
}

func (c *Counter) updateMetric() {
	c.family.mu.Lock()
	defer c.family.mu.Unlock()

	metric := &Metric{
		Name:      c.family.Name,
		Type:      c.family.Type,
		Value:     c.Get(),
		Labels:    c.labels,
		Timestamp: time.Now(),
		Help:      c.family.Help,
	}

	// Update or append metric
	found := false
	for i, m := range c.family.Metrics {
		if labelsEqual(m.Labels, c.labels) {
			c.family.Metrics[i] = metric
			found = true
			break
		}
	}
	if !found {
		c.family.Metrics = append(c.family.Metrics, metric)
	}
}

// Gauge implements a Prometheus-style gauge
type Gauge struct {
	value  *int64 // Using int64 with atomic for thread-safety, storing as float64*1000
	family *MetricFamily
	labels map[string]string
}

// NewGauge creates a new gauge metric
func (mr *MetricRegistry) NewGauge(name, help string, labels map[string]string) *Gauge {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	fullName := mr.prefix + name
	family, exists := mr.metrics[fullName]
	if !exists {
		family = &MetricFamily{
			Name:    fullName,
			Type:    MetricTypeGauge,
			Help:    help,
			Metrics: make([]*Metric, 0),
		}
		mr.metrics[fullName] = family
	}

	gauge := &Gauge{
		value:  new(int64),
		family: family,
		labels: labels,
	}

	return gauge
}

// Set sets the gauge to the given value
func (g *Gauge) Set(value float64) {
	atomic.StoreInt64(g.value, int64(value*1000))
	g.updateMetric()
}

// Inc increments the gauge by 1
func (g *Gauge) Inc() {
	atomic.AddInt64(g.value, 1000)
	g.updateMetric()
}

// Dec decrements the gauge by 1
func (g *Gauge) Dec() {
	atomic.AddInt64(g.value, -1000)
	g.updateMetric()
}

// Add adds the given value to the gauge
func (g *Gauge) Add(delta float64) {
	atomic.AddInt64(g.value, int64(delta*1000))
	g.updateMetric()
}

// Get returns the current gauge value
func (g *Gauge) Get() float64 {
	return float64(atomic.LoadInt64(g.value)) / 1000
}

func (g *Gauge) updateMetric() {
	g.family.mu.Lock()
	defer g.family.mu.Unlock()

	metric := &Metric{
		Name:      g.family.Name,
		Type:      g.family.Type,
		Value:     g.Get(),
		Labels:    g.labels,
		Timestamp: time.Now(),
		Help:      g.family.Help,
	}

	// Update or append metric
	found := false
	for i, m := range g.family.Metrics {
		if labelsEqual(m.Labels, g.labels) {
			g.family.Metrics[i] = metric
			found = true
			break
		}
	}
	if !found {
		g.family.Metrics = append(g.family.Metrics, metric)
	}
}

// Histogram implements a Prometheus-style histogram
type Histogram struct {
	mu      sync.Mutex
	buckets []float64
	counts  []int64
	sum     float64
	count   int64
	family  *MetricFamily
	labels  map[string]string
}

// NewHistogram creates a new histogram metric
func (mr *MetricRegistry) NewHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	fullName := mr.prefix + name
	family, exists := mr.metrics[fullName]
	if !exists {
		family = &MetricFamily{
			Name:    fullName,
			Type:    MetricTypeHistogram,
			Help:    help,
			Metrics: make([]*Metric, 0),
		}
		mr.metrics[fullName] = family
	}

	// Ensure buckets are sorted and include +Inf
	sortedBuckets := make([]float64, len(buckets))
	copy(sortedBuckets, buckets)
	sort.Float64s(sortedBuckets)
	if sortedBuckets[len(sortedBuckets)-1] != math.Inf(1) {
		sortedBuckets = append(sortedBuckets, math.Inf(1))
	}

	histogram := &Histogram{
		buckets: sortedBuckets,
		counts:  make([]int64, len(sortedBuckets)),
		family:  family,
		labels:  labels,
	}

	return histogram
}

// Observe adds a single observation to the histogram
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += value
	h.count++

	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
		}
	}

	h.updateMetric()
}

func (h *Histogram) updateMetric() {
	// Create histogram metrics (bucket counts, sum, count)
	h.family.mu.Lock()
	defer h.family.mu.Unlock()

	// Remove old metrics for this label set
	newMetrics := make([]*Metric, 0)
	for _, metric := range h.family.Metrics {
		if !labelsEqual(metric.Labels, h.labels) {
			newMetrics = append(newMetrics, metric)
		}
	}

	// Add bucket metrics
	for i, bucket := range h.buckets {
		bucketLabels := make(map[string]string)
		for k, v := range h.labels {
			bucketLabels[k] = v
		}
		if bucket == math.Inf(1) {
			bucketLabels["le"] = "+Inf"
		} else {
			bucketLabels["le"] = fmt.Sprintf("%.6f", bucket)
		}

		newMetrics = append(newMetrics, &Metric{
			Name:      h.family.Name + "_bucket",
			Type:      MetricTypeCounter,
			Value:     float64(h.counts[i]),
			Labels:    bucketLabels,
			Timestamp: time.Now(),
			Help:      h.family.Help,
		})
	}

	// Add sum metric
	sumLabels := make(map[string]string)
	for k, v := range h.labels {
		sumLabels[k] = v
	}
	newMetrics = append(newMetrics, &Metric{
		Name:      h.family.Name + "_sum",
		Type:      MetricTypeCounter,
		Value:     h.sum,
		Labels:    sumLabels,
		Timestamp: time.Now(),
		Help:      h.family.Help,
	})

	// Add count metric
	countLabels := make(map[string]string)
	for k, v := range h.labels {
		countLabels[k] = v
	}
	newMetrics = append(newMetrics, &Metric{
		Name:      h.family.Name + "_count",
		Type:      MetricTypeCounter,
		Value:     float64(h.count),
		Labels:    countLabels,
		Timestamp: time.Now(),
		Help:      h.family.Help,
	})

	h.family.Metrics = newMetrics
}

// GetQuantile calculates the approximate quantile from histogram data
func (h *Histogram) GetQuantile(q float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.count == 0 {
		return 0
	}

	targetCount := q * float64(h.count)
	cumulativeCount := int64(0)

	for i, count := range h.counts {
		cumulativeCount += count
		if float64(cumulativeCount) >= targetCount {
			if i == 0 {
				return h.buckets[0]
			}
			// Linear interpolation between buckets
			prevBucket := float64(0)
			if i > 0 {
				if h.buckets[i-1] == math.Inf(-1) {
					prevBucket = 0
				} else {
					prevBucket = h.buckets[i-1]
				}
			}

			if h.buckets[i] == math.Inf(1) {
				return prevBucket
			}

			return prevBucket + (h.buckets[i]-prevBucket)*0.5
		}
	}

	return h.buckets[len(h.buckets)-2] // Return second-to-last bucket (before +Inf)
}

// labelsEqual compares two label maps for equality
func labelsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// GetAllMetrics returns all metrics from the registry
func (mr *MetricRegistry) GetAllMetrics() map[string]*MetricFamily {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	result := make(map[string]*MetricFamily, len(mr.metrics))
	for k, v := range mr.metrics {
		result[k] = v
	}
	return result
}

// ======================== DISTRIBUTED TRACING ========================

// TraceContext represents the context for distributed tracing
type TraceContext struct {
	TraceID  string `json:"trace_id"`
	SpanID   string `json:"span_id"`
	ParentID string `json:"parent_id,omitempty"`
	Flags    byte   `json:"flags"`
}

// Span represents a single span in a distributed trace
type Span struct {
	TraceID       string                 `json:"trace_id"`
	SpanID        string                 `json:"span_id"`
	ParentID      string                 `json:"parent_id,omitempty"`
	OperationName string                 `json:"operation_name"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       *time.Time             `json:"end_time,omitempty"`
	Duration      time.Duration          `json:"duration"`
	Tags          map[string]interface{} `json:"tags"`
	Logs          []SpanLog              `json:"logs"`
	Status        SpanStatus             `json:"status"`
	mu            sync.Mutex
}

// SpanLog represents a log entry within a span
type SpanLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields"`
}

// SpanStatus represents the status of a span
type SpanStatus struct {
	Code    StatusCode `json:"code"`
	Message string     `json:"message"`
}

// StatusCode represents span status codes
type StatusCode int

const (
	StatusCodeOK StatusCode = iota
	StatusCodeCancelled
	StatusCodeUnknown
	StatusCodeInvalidArgument
	StatusCodeDeadlineExceeded
	StatusCodeNotFound
	StatusCodeAlreadyExists
	StatusCodePermissionDenied
	StatusCodeResourceExhausted
	StatusCodeFailedPrecondition
	StatusCodeAborted
	StatusCodeOutOfRange
	StatusCodeUnimplemented
	StatusCodeInternal
	StatusCodeUnavailable
	StatusCodeDataLoss
	StatusCodeUnauthenticated
)

// Tracer manages distributed tracing
type Tracer struct {
	serviceName string
	spans       sync.Map // map[string]*Span
	sampler     Sampler
	exporter    SpanExporter
}

// Sampler determines whether a trace should be sampled
type Sampler interface {
	ShouldSample(traceID string, operationName string) bool
}

// SpanExporter exports spans to external systems
type SpanExporter interface {
	Export(spans []*Span) error
}

// NewTracer creates a new tracer
func NewTracer(serviceName string, sampler Sampler, exporter SpanExporter) *Tracer {
	return &Tracer{
		serviceName: serviceName,
		sampler:     sampler,
		exporter:    exporter,
	}
}

// StartSpan creates a new span
func (t *Tracer) StartSpan(operationName string, parentCtx *TraceContext) *Span {
	traceID := generateTraceID()
	if parentCtx != nil && parentCtx.TraceID != "" {
		traceID = parentCtx.TraceID
	}

	// Check sampling decision
	if !t.sampler.ShouldSample(traceID, operationName) {
		return &Span{} // Return empty span for non-sampled traces
	}

	span := &Span{
		TraceID:       traceID,
		SpanID:        generateSpanID(),
		OperationName: operationName,
		StartTime:     time.Now(),
		Tags:          make(map[string]interface{}),
		Logs:          make([]SpanLog, 0),
		Status:        SpanStatus{Code: StatusCodeOK},
	}

	if parentCtx != nil {
		span.ParentID = parentCtx.SpanID
	}

	// Add service tags
	span.SetTag("service.name", t.serviceName)
	span.SetTag("service.version", "1.0.0") // This could be configurable

	t.spans.Store(span.SpanID, span)
	return span
}

// SetTag sets a tag on the span
func (s *Span) SetTag(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tags[key] = value
}

// LogFields adds a log entry to the span
func (s *Span) LogFields(fields map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Logs = append(s.Logs, SpanLog{
		Timestamp: time.Now(),
		Fields:    fields,
	})
}

// SetStatus sets the span status
func (s *Span) SetStatus(code StatusCode, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = SpanStatus{
		Code:    code,
		Message: message,
	}
}

// Finish completes the span
func (s *Span) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.EndTime = &now
	s.Duration = now.Sub(s.StartTime)
}

// GetContext returns the trace context for this span
func (s *Span) GetContext() *TraceContext {
	return &TraceContext{
		TraceID:  s.TraceID,
		SpanID:   s.SpanID,
		ParentID: s.ParentID,
		Flags:    0,
	}
}

// generateTraceID generates a unique trace ID
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateSpanID generates a unique span ID
func generateSpanID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ProbabilisticSampler implements probabilistic sampling
type ProbabilisticSampler struct {
	rate float64
}

// NewProbabilisticSampler creates a new probabilistic sampler
func NewProbabilisticSampler(rate float64) *ProbabilisticSampler {
	return &ProbabilisticSampler{rate: rate}
}

// ShouldSample determines if a trace should be sampled based on probability
func (ps *ProbabilisticSampler) ShouldSample(traceID string, operationName string) bool {
	// Use trace ID to ensure consistent sampling decisions across services
	if len(traceID) < 16 {
		return false
	}

	// Convert last 8 hex chars to int64
	lastBytes := traceID[len(traceID)-8:]
	val, err := strconv.ParseInt(lastBytes, 16, 64)
	if err != nil {
		return false
	}

	threshold := int64(ps.rate * float64(math.MaxInt64))
	return val < threshold
}

// ======================== STRUCTURED LOGGING WITH CORRELATION IDs ========================

// Logger provides structured logging with correlation IDs
type Logger struct {
	serviceName string
	version     string
	output      io.Writer
	level       LogLevel
	mu          sync.Mutex
}

// LogLevel represents log levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp     time.Time              `json:"timestamp"`
	Level         string                 `json:"level"`
	Message       string                 `json:"message"`
	Service       string                 `json:"service"`
	Version       string                 `json:"version"`
	TraceID       string                 `json:"trace_id,omitempty"`
	SpanID        string                 `json:"span_id,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	RequestID     string                 `json:"request_id,omitempty"`
	Fields        map[string]interface{} `json:"fields,omitempty"`
	Error         *ErrorInfo             `json:"error,omitempty"`
	Performance   *PerformanceInfo       `json:"performance,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace,omitempty"`
	Code       string `json:"code,omitempty"`
}

// PerformanceInfo contains performance metrics
type PerformanceInfo struct {
	Duration   time.Duration `json:"duration"`
	MemoryUsed int64         `json:"memory_used"`
	CPUUsed    float64       `json:"cpu_used"`
}

// NewLogger creates a new structured logger
func NewLogger(serviceName, version string, output io.Writer, level LogLevel) *Logger {
	return &Logger{
		serviceName: serviceName,
		version:     version,
		output:      output,
		level:       level,
	}
}

// WithContext creates a logger with trace context
func (l *Logger) WithContext(ctx *TraceContext) *ContextLogger {
	return &ContextLogger{
		logger:  l,
		traceID: ctx.TraceID,
		spanID:  ctx.SpanID,
	}
}

// ContextLogger is a logger with embedded context
type ContextLogger struct {
	logger        *Logger
	traceID       string
	spanID        string
	correlationID string
	userID        string
	requestID     string
	fields        map[string]interface{}
}

// WithCorrelationID adds a correlation ID to the logger
func (cl *ContextLogger) WithCorrelationID(correlationID string) *ContextLogger {
	newLogger := *cl
	newLogger.correlationID = correlationID
	return &newLogger
}

// WithUserID adds a user ID to the logger
func (cl *ContextLogger) WithUserID(userID string) *ContextLogger {
	newLogger := *cl
	newLogger.userID = userID
	return &newLogger
}

// WithRequestID adds a request ID to the logger
func (cl *ContextLogger) WithRequestID(requestID string) *ContextLogger {
	newLogger := *cl
	newLogger.requestID = requestID
	return &newLogger
}

// WithFields adds structured fields to the logger
func (cl *ContextLogger) WithFields(fields map[string]interface{}) *ContextLogger {
	newLogger := *cl
	newLogger.fields = make(map[string]interface{})
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return &newLogger
}

// Info logs an info message
func (cl *ContextLogger) Info(message string) {
	cl.log(LogLevelInfo, message, nil, nil)
}

// Error logs an error message
func (cl *ContextLogger) Error(message string, err error) {
	var errorInfo *ErrorInfo
	if err != nil {
		errorInfo = &ErrorInfo{
			Type:    fmt.Sprintf("%T", err),
			Message: err.Error(),
		}
	}
	cl.log(LogLevelError, message, errorInfo, nil)
}

// Debug logs a debug message
func (cl *ContextLogger) Debug(message string) {
	cl.log(LogLevelDebug, message, nil, nil)
}

// Warn logs a warning message
func (cl *ContextLogger) Warn(message string) {
	cl.log(LogLevelWarn, message, nil, nil)
}

// Fatal logs a fatal message
func (cl *ContextLogger) Fatal(message string, err error) {
	var errorInfo *ErrorInfo
	if err != nil {
		errorInfo = &ErrorInfo{
			Type:    fmt.Sprintf("%T", err),
			Message: err.Error(),
		}
	}
	cl.log(LogLevelFatal, message, errorInfo, nil)
}

// LogPerformance logs performance metrics
func (cl *ContextLogger) LogPerformance(message string, duration time.Duration) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	perfInfo := &PerformanceInfo{
		Duration:   duration,
		MemoryUsed: int64(m.Alloc),
		CPUUsed:    0, // Would need more sophisticated CPU monitoring
	}

	cl.log(LogLevelInfo, message, nil, perfInfo)
}

func (cl *ContextLogger) log(level LogLevel, message string, errorInfo *ErrorInfo, perfInfo *PerformanceInfo) {
	if level < cl.logger.level {
		return
	}

	entry := LogEntry{
		Timestamp:     time.Now().UTC(),
		Level:         levelToString(level),
		Message:       message,
		Service:       cl.logger.serviceName,
		Version:       cl.logger.version,
		TraceID:       cl.traceID,
		SpanID:        cl.spanID,
		CorrelationID: cl.correlationID,
		UserID:        cl.userID,
		RequestID:     cl.requestID,
		Fields:        cl.fields,
		Error:         errorInfo,
		Performance:   perfInfo,
	}

	cl.logger.mu.Lock()
	defer cl.logger.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging
		fmt.Fprintf(cl.logger.output, "%s [%s] %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Level,
			entry.Message)
		return
	}

	cl.logger.output.Write(data)
	cl.logger.output.Write([]byte("\n"))
}

func levelToString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ======================== SLI/SLO MONITORING ========================

// SLI represents a Service Level Indicator
type SLI struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"`
	Labels      map[string]string `json:"labels"`
	Threshold   float64           `json:"threshold"`
	mu          sync.RWMutex
	values      []SLIValue
}

// SLIValue represents a point-in-time SLI measurement
type SLIValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Good      bool      `json:"good"`
}

// SLO represents a Service Level Objective
type SLO struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Target       float64       `json:"target"` // e.g., 99.9 for 99.9%
	Period       time.Duration `json:"period"` // e.g., 30 days
	SLIs         []*SLI        `json:"slis"`
	mu           sync.RWMutex
	measurements []SLOMeasurement
}

// SLOMeasurement represents an SLO measurement over time
type SLOMeasurement struct {
	Timestamp       time.Time `json:"timestamp"`
	ActualSLI       float64   `json:"actual_sli"`
	TargetSLI       float64   `json:"target_sli"`
	ErrorBudget     float64   `json:"error_budget"`
	ErrorBudgetUsed float64   `json:"error_budget_used"`
}

// SLOMonitor manages SLI/SLO monitoring
type SLOMonitor struct {
	slos         map[string]*SLO
	alertManager *AlertManager
	mu           sync.RWMutex
}

// NewSLOMonitor creates a new SLO monitor
func NewSLOMonitor(alertManager *AlertManager) *SLOMonitor {
	return &SLOMonitor{
		slos:         make(map[string]*SLO),
		alertManager: alertManager,
	}
}

// RegisterSLO registers a new SLO for monitoring
func (sm *SLOMonitor) RegisterSLO(slo *SLO) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.slos[slo.Name] = slo
}

// RecordSLI records an SLI measurement
func (sm *SLOMonitor) RecordSLI(sloName, sliName string, value float64, good bool) {
	sm.mu.RLock()
	slo, exists := sm.slos[sloName]
	sm.mu.RUnlock()

	if !exists {
		return
	}

	slo.mu.Lock()
	defer slo.mu.Unlock()

	// Find the SLI and record the value
	for _, sli := range slo.SLIs {
		if sli.Name == sliName {
			sli.mu.Lock()
			sli.values = append(sli.values, SLIValue{
				Timestamp: time.Now(),
				Value:     value,
				Good:      good,
			})

			// Keep only recent values (within the SLO period)
			cutoff := time.Now().Add(-slo.Period)
			var recentValues []SLIValue
			for _, v := range sli.values {
				if v.Timestamp.After(cutoff) {
					recentValues = append(recentValues, v)
				}
			}
			sli.values = recentValues
			sli.mu.Unlock()
			break
		}
	}

	// Calculate current SLO performance
	sm.calculateSLO(slo)
}

// calculateSLO calculates current SLO performance and triggers alerts if needed
func (sm *SLOMonitor) calculateSLO(slo *SLO) {
	if len(slo.SLIs) == 0 {
		return
	}

	totalGood := 0
	totalMeasurements := 0

	for _, sli := range slo.SLIs {
		sli.mu.RLock()
		for _, value := range sli.values {
			if value.Good {
				totalGood++
			}
			totalMeasurements++
		}
		sli.mu.RUnlock()
	}

	if totalMeasurements == 0 {
		return
	}

	actualSLI := float64(totalGood) / float64(totalMeasurements) * 100
	errorBudget := 100 - slo.Target
	errorBudgetUsed := (slo.Target - actualSLI) / errorBudget * 100

	measurement := SLOMeasurement{
		Timestamp:       time.Now(),
		ActualSLI:       actualSLI,
		TargetSLI:       slo.Target,
		ErrorBudget:     errorBudget,
		ErrorBudgetUsed: errorBudgetUsed,
	}

	slo.measurements = append(slo.measurements, measurement)

	// Trigger alert if SLO is violated
	if actualSLI < slo.Target {
		alert := Alert{
			Name:        fmt.Sprintf("SLO Violation: %s", slo.Name),
			Description: fmt.Sprintf("SLO %s is below target (%.2f%% < %.2f%%)", slo.Name, actualSLI, slo.Target),
			Severity:    SeverityHigh,
			Labels: map[string]string{
				"slo":        slo.Name,
				"actual_sli": fmt.Sprintf("%.2f", actualSLI),
				"target_sli": fmt.Sprintf("%.2f", slo.Target),
			},
			Timestamp: time.Now(),
		}
		sm.alertManager.TriggerAlert(alert)
	}
}

// GetSLOStatus returns current SLO status
func (sm *SLOMonitor) GetSLOStatus(sloName string) (*SLOMeasurement, error) {
	sm.mu.RLock()
	slo, exists := sm.slos[sloName]
	sm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("SLO %s not found", sloName)
	}

	slo.mu.RLock()
	defer slo.mu.RUnlock()

	if len(slo.measurements) == 0 {
		return nil, fmt.Errorf("no measurements available for SLO %s", sloName)
	}

	// Return the latest measurement
	return &slo.measurements[len(slo.measurements)-1], nil
}

// ======================== ALERTING AND ESCALATION PATTERNS ========================

// Alert represents an alert
type Alert struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Severity    Severity          `json:"severity"`
	Status      AlertStatus       `json:"status"`
	Labels      map[string]string `json:"labels"`
	Timestamp   time.Time         `json:"timestamp"`
	ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
}

// Severity represents alert severity levels
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// AlertStatus represents alert status
type AlertStatus int

const (
	AlertStatusFiring AlertStatus = iota
	AlertStatusResolved
	AlertStatusSuppressed
)

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	Name        string            `json:"name"`
	Query       string            `json:"query"`
	Condition   string            `json:"condition"` // e.g., "> 0.95", "< 100"
	Duration    time.Duration     `json:"duration"`  // How long condition must be true
	Severity    Severity          `json:"severity"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// AlertManager manages alerting and escalation
type AlertManager struct {
	rules           map[string]*AlertRule
	activeAlerts    map[string]*Alert
	escalationRules []*EscalationRule
	notifiers       map[string]Notifier
	mu              sync.RWMutex
}

// EscalationRule defines escalation behavior
type EscalationRule struct {
	Severity     Severity      `json:"severity"`
	Delay        time.Duration `json:"delay"`
	MaxAttempts  int           `json:"max_attempts"`
	NotifierType string        `json:"notifier_type"`
}

// Notifier interface for different notification channels
type Notifier interface {
	Send(alert Alert) error
	GetType() string
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules:        make(map[string]*AlertRule),
		activeAlerts: make(map[string]*Alert),
		notifiers:    make(map[string]Notifier),
	}
}

// RegisterNotifier registers a notification channel
func (am *AlertManager) RegisterNotifier(notifier Notifier) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.notifiers[notifier.GetType()] = notifier
}

// AddEscalationRule adds an escalation rule
func (am *AlertManager) AddEscalationRule(rule *EscalationRule) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.escalationRules = append(am.escalationRules, rule)
}

// TriggerAlert triggers a new alert or updates an existing one
func (am *AlertManager) TriggerAlert(alert Alert) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Generate ID if not provided
	if alert.ID == "" {
		alert.ID = generateAlertID(alert.Name, alert.Labels)
	}

	// Check if alert already exists
	existingAlert, exists := am.activeAlerts[alert.ID]
	if exists && existingAlert.Status == AlertStatusFiring {
		// Update existing alert
		existingAlert.Timestamp = alert.Timestamp
		return
	}

	// New alert
	alert.Status = AlertStatusFiring
	am.activeAlerts[alert.ID] = &alert

	// Start escalation process
	go am.escalateAlert(&alert)
}

// ResolveAlert resolves an active alert
func (am *AlertManager) ResolveAlert(alertID string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, exists := am.activeAlerts[alertID]
	if exists {
		now := time.Now()
		alert.Status = AlertStatusResolved
		alert.ResolvedAt = &now
	}
}

// escalateAlert handles alert escalation
func (am *AlertManager) escalateAlert(alert *Alert) {
	// Find applicable escalation rules
	var applicableRules []*EscalationRule
	am.mu.RLock()
	for _, rule := range am.escalationRules {
		if rule.Severity <= alert.Severity {
			applicableRules = append(applicableRules, rule)
		}
	}
	am.mu.RUnlock()

	// Sort rules by delay
	sort.Slice(applicableRules, func(i, j int) bool {
		return applicableRules[i].Delay < applicableRules[j].Delay
	})

	// Execute escalation steps
	for _, rule := range applicableRules {
		// Wait for delay
		time.Sleep(rule.Delay)

		// Check if alert is still active
		am.mu.RLock()
		currentAlert, exists := am.activeAlerts[alert.ID]
		am.mu.RUnlock()

		if !exists || currentAlert.Status != AlertStatusFiring {
			// Alert was resolved, stop escalation
			return
		}

		// Send notification
		am.mu.RLock()
		notifier, exists := am.notifiers[rule.NotifierType]
		am.mu.RUnlock()

		if exists {
			for attempt := 0; attempt < rule.MaxAttempts; attempt++ {
				if err := notifier.Send(*alert); err != nil {
					log.Printf("Failed to send alert notification (attempt %d/%d): %v",
						attempt+1, rule.MaxAttempts, err)
					if attempt < rule.MaxAttempts-1 {
						time.Sleep(time.Second * time.Duration(attempt+1))
					}
				} else {
					break
				}
			}
		}
	}
}

// generateAlertID generates a unique alert ID based on name and labels
func generateAlertID(name string, labels map[string]string) string {
	var labelStr strings.Builder
	labelStr.WriteString(name)

	// Create sorted label string for consistent IDs
	var keys []string
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		labelStr.WriteString(k)
		labelStr.WriteString("=")
		labelStr.WriteString(labels[k])
	}

	// Simple hash of the label string
	return fmt.Sprintf("alert_%x", labelStr.String())[:16]
}

// EmailNotifier sends alerts via email
type EmailNotifier struct{}

func (en *EmailNotifier) Send(alert Alert) error {
	// Implementation would integrate with email service
	log.Printf("EMAIL ALERT: [%s] %s - %s",
		severityToString(alert.Severity),
		alert.Name,
		alert.Description)
	return nil
}

func (en *EmailNotifier) GetType() string {
	return "email"
}

// SlackNotifier sends alerts to Slack
type SlackNotifier struct{}

func (sn *SlackNotifier) Send(alert Alert) error {
	// Implementation would integrate with Slack API
	log.Printf("SLACK ALERT: [%s] %s - %s",
		severityToString(alert.Severity),
		alert.Name,
		alert.Description)
	return nil
}

func (sn *SlackNotifier) GetType() string {
	return "slack"
}

// PagerDutyNotifier sends alerts to PagerDuty
type PagerDutyNotifier struct{}

func (pn *PagerDutyNotifier) Send(alert Alert) error {
	// Implementation would integrate with PagerDuty API
	log.Printf("PAGERDUTY ALERT: [%s] %s - %s",
		severityToString(alert.Severity),
		alert.Name,
		alert.Description)
	return nil
}

func (pn *PagerDutyNotifier) GetType() string {
	return "pagerduty"
}

func severityToString(severity Severity) string {
	switch severity {
	case SeverityInfo:
		return "INFO"
	case SeverityLow:
		return "LOW"
	case SeverityMedium:
		return "MEDIUM"
	case SeverityHigh:
		return "HIGH"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ======================== DASHBOARD AND VISUALIZATION ========================

// Dashboard represents a monitoring dashboard
type Dashboard struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Panels      []*Panel  `json:"panels"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Panel represents a dashboard panel
type Panel struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Type       PanelType              `json:"type"`
	Query      string                 `json:"query"`
	Options    map[string]interface{} `json:"options"`
	Position   Position               `json:"position"`
	Thresholds []Threshold            `json:"thresholds"`
}

// PanelType represents different panel types
type PanelType int

const (
	PanelTypeGraph PanelType = iota
	PanelTypeStat
	PanelTypeTable
	PanelTypeHeatmap
	PanelTypeAlert
)

// Position represents panel position on dashboard
type Position struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Threshold represents visualization thresholds
type Threshold struct {
	Value float64 `json:"value"`
	Color string  `json:"color"`
	Label string  `json:"label"`
}

// DashboardManager manages dashboards and visualizations
type DashboardManager struct {
	dashboards map[string]*Dashboard
	mu         sync.RWMutex
}

// NewDashboardManager creates a new dashboard manager
func NewDashboardManager() *DashboardManager {
	return &DashboardManager{
		dashboards: make(map[string]*Dashboard),
	}
}

// CreateDashboard creates a new dashboard
func (dm *DashboardManager) CreateDashboard(title, description string, tags []string) *Dashboard {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dashboard := &Dashboard{
		ID:          generateDashboardID(),
		Title:       title,
		Description: description,
		Panels:      make([]*Panel, 0),
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dm.dashboards[dashboard.ID] = dashboard
	return dashboard
}

// AddPanel adds a panel to a dashboard
func (dm *DashboardManager) AddPanel(dashboardID string, panel *Panel) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dashboard, exists := dm.dashboards[dashboardID]
	if !exists {
		return fmt.Errorf("dashboard %s not found", dashboardID)
	}

	if panel.ID == "" {
		panel.ID = generatePanelID()
	}

	dashboard.Panels = append(dashboard.Panels, panel)
	dashboard.UpdatedAt = time.Now()

	return nil
}

// GetDashboard retrieves a dashboard by ID
func (dm *DashboardManager) GetDashboard(dashboardID string) (*Dashboard, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	dashboard, exists := dm.dashboards[dashboardID]
	if !exists {
		return nil, fmt.Errorf("dashboard %s not found", dashboardID)
	}

	return dashboard, nil
}

// ListDashboards returns all dashboards
func (dm *DashboardManager) ListDashboards() []*Dashboard {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	dashboards := make([]*Dashboard, 0, len(dm.dashboards))
	for _, dashboard := range dm.dashboards {
		dashboards = append(dashboards, dashboard)
	}

	return dashboards
}

func generateDashboardID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "dash_" + hex.EncodeToString(bytes)
}

func generatePanelID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "panel_" + hex.EncodeToString(bytes)
}

// ======================== HTTP HANDLERS FOR OBSERVABILITY ========================

// ObservabilityHandler provides HTTP endpoints for observability data
type ObservabilityHandler struct {
	registry         *MetricRegistry
	tracer           *Tracer
	logger           *Logger
	sloMonitor       *SLOMonitor
	alertManager     *AlertManager
	dashboardManager *DashboardManager
}

// NewObservabilityHandler creates HTTP handlers for observability
func NewObservabilityHandler() *ObservabilityHandler {
	registry := NewMetricRegistry("myapp_")
	sampler := NewProbabilisticSampler(0.1)
	exporter := &ConsoleSpanExporter{}
	tracer := NewTracer("myapp", sampler, exporter)
	logger := NewLogger("myapp", "1.0.0", log.Writer(), LogLevelInfo)
	alertManager := NewAlertManager()
	sloMonitor := NewSLOMonitor(alertManager)
	dashboardManager := NewDashboardManager()

	return &ObservabilityHandler{
		registry:         registry,
		tracer:           tracer,
		logger:           logger,
		sloMonitor:       sloMonitor,
		alertManager:     alertManager,
		dashboardManager: dashboardManager,
	}
}

// HandleMetrics serves Prometheus-compatible metrics
func (oh *ObservabilityHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	metrics := oh.registry.GetAllMetrics()
	json.NewEncoder(w).Encode(metrics)
}

// HandleSLOStatus serves SLO status information
func (oh *ObservabilityHandler) HandleSLOStatus(w http.ResponseWriter, r *http.Request) {
	sloName := r.URL.Query().Get("slo")
	if sloName == "" {
		http.Error(w, "slo parameter required", http.StatusBadRequest)
		return
	}

	status, err := oh.sloMonitor.GetSLOStatus(sloName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// HandleDashboards serves dashboard data
func (oh *ObservabilityHandler) HandleDashboards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dashboards := oh.dashboardManager.ListDashboards()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dashboards)
	case http.MethodPost:
		var req struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dashboard := oh.dashboardManager.CreateDashboard(req.Title, req.Description, req.Tags)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dashboard)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ConsoleSpanExporter exports spans to console for demonstration
type ConsoleSpanExporter struct{}

func (cse *ConsoleSpanExporter) Export(spans []*Span) error {
	for _, span := range spans {
		log.Printf("TRACE: %s [%s] %s -> %s (%.2fms)",
			span.TraceID,
			span.SpanID,
			span.OperationName,
			span.Status.Message,
			float64(span.Duration.Nanoseconds())/1000000)
	}
	return nil
}
