// Package behavioral implements comprehensive STAR method framework for FAANG L3/L4 behavioral interviews
//
// This module provides production-grade behavioral interview preparation with:
// - STAR (Situation, Task, Action, Result) method implementation
// - Go-specific technical examples and scenarios
// - Quantified impact metrics and business value demonstration
// - Leadership story development with progressive responsibility
// - Cross-functional collaboration and mentoring examples
//
// Target Companies: Google, Meta, Netflix, Amazon, Apple, Microsoft, ByteDance, Uber, Airbnb
// Target Levels: L3/L4 (Software Engineer I/II) with leadership potential
// Success Metrics: Compelling stories with quantified business impact and technical depth
package behavioral

import (
	"fmt"
	"strings"
	"time"
)

// STARStory represents a structured behavioral interview story using STAR methodology
type STARStory struct {
	Title        string         // Brief story title for easy reference
	Category     StoryCategory  // Type of story (leadership, problem-solving, etc.)
	Situation    string         // Context and background (30-40% of story)
	Task         string         // Your responsibility and challenges (20-25%)
	Action       string         // Specific actions taken (30-35%)
	Result       string         // Quantified outcomes and impact (15-20%)
	Metrics      []ImpactMetric // Quantified business/technical impact
	Technologies []string       // Go technologies and tools used
	Duration     time.Duration  // Timeline of the project/situation
	TeamSize     int            // Number of people involved
	Keywords     []string       // Keywords for company/role matching
}

// StoryCategory defines different types of behavioral stories
type StoryCategory int

const (
	// Core Leadership Categories
	TechnicalLeadership StoryCategory = iota
	ProblemSolving
	Innovation
	Collaboration
	Mentoring

	// Challenging Situations
	Failure
	Conflict
	DeadlinePressure
	Ambiguity

	// Growth & Learning
	LearningAgility
	Initiative
	ContinuousImprovement

	// Company-Specific
	CustomerObsession // Amazon
	Ownership         // Amazon
	Googleyness       // Google
	JediEngineering   // Meta
	NetflixCulture    // Netflix
)

// ImpactMetric represents quantified business or technical impact
type ImpactMetric struct {
	Type        MetricType `json:"type"`
	Description string     `json:"description"`
	Value       float64    `json:"value"`
	Unit        string     `json:"unit"`
	Baseline    float64    `json:"baseline,omitempty"`
	Timeframe   string     `json:"timeframe"`
}

// MetricType defines categories of impact metrics
type MetricType int

const (
	// Performance Metrics
	PerformanceImprovement MetricType = iota
	LatencyReduction
	ThroughputIncrease
	MemoryOptimization
	CPUOptimization

	// Business Metrics
	CostSavings
	RevenueImpact
	UserGrowth
	ErrorReduction

	// Engineering Metrics
	CodeQualityImprovement
	DeploymentFrequency
	TimeToMarket
	TechnicalDebtReduction
	TestCoverage

	// Team Metrics
	DeveloperProductivity
	TeamVelocity
	KnowledgeTransfer
	OnboardingTime
)

// STARFramework provides methods for developing and presenting behavioral stories
type STARFramework struct {
	stories           []STARStory
	targetCompany     string
	targetLevel       string
	technicalKeywords []string
}

// NewSTARFramework creates a new STAR framework instance
func NewSTARFramework(targetCompany, targetLevel string) *STARFramework {
	return &STARFramework{
		stories:           make([]STARStory, 0),
		targetCompany:     targetCompany,
		targetLevel:       targetLevel,
		technicalKeywords: getGoTechnicalKeywords(),
	}
}

// AddStory adds a new STAR story to the framework
func (sf *STARFramework) AddStory(story STARStory) error {
	// Validate story completeness
	if err := sf.validateStory(story); err != nil {
		return fmt.Errorf("story validation failed: %v", err)
	}

	// Auto-enhance story with Go-specific technical details
	story = sf.enhanceWithTechnicalDetails(story)

	sf.stories = append(sf.stories, story)
	return nil
}

// GetStoriesForCategory returns stories matching a specific category
func (sf *STARFramework) GetStoriesForCategory(category StoryCategory) []STARStory {
	var matchingStories []STARStory
	for _, story := range sf.stories {
		if story.Category == category {
			matchingStories = append(matchingStories, story)
		}
	}
	return matchingStories
}

// GetStoriesForCompany returns stories optimized for a specific company
func (sf *STARFramework) GetStoriesForCompany(company string) []STARStory {
	companyKeywords := getCompanyKeywords(company)
	var matchingStories []STARStory

	for _, story := range sf.stories {
		score := sf.calculateCompanyMatchScore(story, companyKeywords)
		if score > 0.6 { // 60% match threshold
			matchingStories = append(matchingStories, story)
		}
	}

	return matchingStories
}

// FormatSTARPresentation formats a story for interview presentation
func (sf *STARFramework) FormatSTARPresentation(story STARStory, timeLimit time.Duration) string {
	var presentation strings.Builder

	// Calculate time allocation based on STAR percentages
	situationTime := int(timeLimit.Seconds() * 0.35)
	taskTime := int(timeLimit.Seconds() * 0.22)
	actionTime := int(timeLimit.Seconds() * 0.33)
	resultTime := int(timeLimit.Seconds() * 0.10)

	presentation.WriteString(fmt.Sprintf("=== %s ===\n\n", story.Title))

	// Situation (35% of time)
	presentation.WriteString(fmt.Sprintf("SITUATION (%ds):\n", situationTime))
	presentation.WriteString(sf.formatSection(story.Situation, situationTime))
	presentation.WriteString("\n\n")

	// Task (22% of time)
	presentation.WriteString(fmt.Sprintf("TASK (%ds):\n", taskTime))
	presentation.WriteString(sf.formatSection(story.Task, taskTime))
	presentation.WriteString("\n\n")

	// Action (33% of time)
	presentation.WriteString(fmt.Sprintf("ACTION (%ds):\n", actionTime))
	presentation.WriteString(sf.formatSection(story.Action, actionTime))
	presentation.WriteString("\n\n")

	// Result (10% of time)
	presentation.WriteString(fmt.Sprintf("RESULT (%ds):\n", resultTime))
	presentation.WriteString(sf.formatSection(story.Result, resultTime))

	// Add quantified metrics
	if len(story.Metrics) > 0 {
		presentation.WriteString("\n\nQUANTIFIED IMPACT:\n")
		for _, metric := range story.Metrics {
			presentation.WriteString(fmt.Sprintf("• %s: %.1f%s %s\n",
				metric.Description, metric.Value, metric.Unit, metric.Timeframe))
		}
	}

	return presentation.String()
}

// GenerateFollowUpQuestions generates potential follow-up questions for a story
func (sf *STARFramework) GenerateFollowUpQuestions(story STARStory) []string {
	questions := []string{
		"What would you do differently if you faced this situation again?",
		"How did you measure the success of your solution?",
		"What was the biggest challenge you faced during implementation?",
		"How did you communicate the technical complexity to non-technical stakeholders?",
		"What did you learn from this experience that you apply to your current work?",
	}

	// Add category-specific questions
	categoryQuestions := sf.getCategorySpecificQuestions(story.Category)
	questions = append(questions, categoryQuestions...)

	// Add Go-specific technical questions
	if len(story.Technologies) > 0 {
		questions = append(questions,
			"What Go-specific features or patterns did you leverage in your solution?",
			"How did you ensure your Go code was production-ready and maintainable?",
			"What performance considerations did you account for in your Go implementation?",
		)
	}

	return questions
}

// validateStory ensures story meets STAR framework requirements
func (sf *STARFramework) validateStory(story STARStory) error {
	if story.Title == "" {
		return fmt.Errorf("story title is required")
	}

	if story.Situation == "" {
		return fmt.Errorf("situation section is required")
	}

	if story.Task == "" {
		return fmt.Errorf("task section is required")
	}

	if story.Action == "" {
		return fmt.Errorf("action section is required")
	}

	if story.Result == "" {
		return fmt.Errorf("result section is required")
	}

	// Validate quantified metrics exist
	if len(story.Metrics) == 0 {
		return fmt.Errorf("at least one quantified metric is required")
	}

	return nil
}

// enhanceWithTechnicalDetails adds Go-specific technical context
func (sf *STARFramework) enhanceWithTechnicalDetails(story STARStory) STARStory {
	// Add Go technical keywords if not present
	for _, keyword := range sf.technicalKeywords {
		if strings.Contains(strings.ToLower(story.Action), strings.ToLower(keyword)) {
			story.Technologies = appendUnique(story.Technologies, keyword)
		}
	}

	return story
}

// calculateCompanyMatchScore calculates how well a story matches company values
func (sf *STARFramework) calculateCompanyMatchScore(story STARStory, companyKeywords []string) float64 {
	storyText := strings.ToLower(story.Situation + " " + story.Task + " " + story.Action + " " + story.Result)

	matches := 0
	for _, keyword := range companyKeywords {
		if strings.Contains(storyText, strings.ToLower(keyword)) {
			matches++
		}
	}

	return float64(matches) / float64(len(companyKeywords))
}

// formatSection formats a story section based on time allocation
func (sf *STARFramework) formatSection(content string, timeSeconds int) string {
	// Estimate words per second (average speaking pace: 150-160 words per minute)
	wordsPerSecond := 2.5
	targetWords := int(float64(timeSeconds) * wordsPerSecond)

	words := strings.Fields(content)
	if len(words) <= targetWords {
		return content
	}

	// Truncate to target word count with ellipsis
	return strings.Join(words[:targetWords], " ") + "..."
}

// getCategorySpecificQuestions returns questions specific to story category
func (sf *STARFramework) getCategorySpecificQuestions(category StoryCategory) []string {
	switch category {
	case TechnicalLeadership:
		return []string{
			"How did you influence technical decisions without formal authority?",
			"What strategies did you use to align your team on technical direction?",
			"How did you handle disagreements about technical approaches?",
		}
	case ProblemSolving:
		return []string{
			"What was your debugging methodology for this complex issue?",
			"How did you prioritize which root causes to investigate first?",
			"What tools or techniques helped you identify the solution?",
		}
	case Innovation:
		return []string{
			"How did you validate your innovative approach before implementation?",
			"What resistance did you face and how did you overcome it?",
			"How do you stay current with emerging technologies and patterns?",
		}
	case Mentoring:
		return []string{
			"How did you adapt your mentoring style to different learning preferences?",
			"What metrics did you use to measure mentoring success?",
			"How did you balance mentoring with your own deliverables?",
		}
	default:
		return []string{}
	}
}

// getGoTechnicalKeywords returns Go-specific technical keywords for story enhancement
func getGoTechnicalKeywords() []string {
	return []string{
		// Core Go Concepts
		"goroutines", "channels", "concurrency", "parallelism",
		"interfaces", "embedding", "composition",
		"garbage collection", "memory management", "escape analysis",

		// Go Tools & Ecosystem
		"go modules", "gofmt", "golint", "go test", "go build",
		"pprof", "race detector", "benchmarking",

		// Patterns & Best Practices
		"dependency injection", "factory pattern", "singleton",
		"middleware", "context", "error handling",

		// Performance & Optimization
		"memory pools", "sync.Pool", "atomic operations",
		"lock-free programming", "cache optimization",

		// Production & Deployment
		"microservices", "containerization", "Kubernetes",
		"monitoring", "observability", "distributed tracing",
	}
}

// getCompanyKeywords returns keywords that align with specific company values
func getCompanyKeywords(company string) []string {
	switch strings.ToLower(company) {
	case "google":
		return []string{
			"scale", "innovation", "user focus", "technical excellence",
			"collaboration", "data-driven", "long-term thinking",
		}
	case "meta", "facebook":
		return []string{
			"move fast", "be bold", "focus on impact", "be open",
			"build social value", "user-centric", "data-informed",
		}
	case "netflix":
		return []string{
			"high performance", "freedom", "responsibility", "candor",
			"innovation", "efficiency", "curiosity", "courage",
		}
	case "amazon":
		return []string{
			"customer obsession", "ownership", "invent and simplify",
			"learn and be curious", "hire and develop", "insist on highest standards",
		}
	case "apple":
		return []string{
			"innovation", "quality", "simplicity", "user experience",
			"attention to detail", "excellence", "privacy", "accessibility",
		}
	default:
		return []string{
			"innovation", "collaboration", "quality", "user focus",
			"technical excellence", "continuous learning",
		}
	}
}

// appendUnique adds item to slice if not already present
func appendUnique(slice []string, item string) []string {
	for _, existing := range slice {
		if existing == item {
			return slice
		}
	}
	return append(slice, item)
}

// CategoryString returns string representation of story category
func (sc StoryCategory) String() string {
	categories := map[StoryCategory]string{
		TechnicalLeadership:   "Technical Leadership",
		ProblemSolving:        "Problem Solving",
		Innovation:            "Innovation",
		Collaboration:         "Collaboration",
		Mentoring:             "Mentoring",
		Failure:               "Failure & Learning",
		Conflict:              "Conflict Resolution",
		DeadlinePressure:      "Working Under Pressure",
		Ambiguity:             "Dealing with Ambiguity",
		LearningAgility:       "Learning Agility",
		Initiative:            "Taking Initiative",
		ContinuousImprovement: "Continuous Improvement",
		CustomerObsession:     "Customer Obsession",
		Ownership:             "Ownership",
		Googleyness:           "Googleyness",
		JediEngineering:       "Jedi Engineering",
		NetflixCulture:        "Netflix Culture",
	}

	if category, exists := categories[sc]; exists {
		return category
	}
	return "Unknown Category"
}

// MetricTypeString returns string representation of metric type
func (mt MetricType) String() string {
	metrics := map[MetricType]string{
		PerformanceImprovement: "Performance Improvement",
		LatencyReduction:       "Latency Reduction",
		ThroughputIncrease:     "Throughput Increase",
		MemoryOptimization:     "Memory Optimization",
		CPUOptimization:        "CPU Optimization",
		CostSavings:            "Cost Savings",
		RevenueImpact:          "Revenue Impact",
		UserGrowth:             "User Growth",
		ErrorReduction:         "Error Reduction",
		CodeQualityImprovement: "Code Quality Improvement",
		DeploymentFrequency:    "Deployment Frequency",
		TimeToMarket:           "Time to Market",
		TechnicalDebtReduction: "Technical Debt Reduction",
		TestCoverage:           "Test Coverage",
		DeveloperProductivity:  "Developer Productivity",
		TeamVelocity:           "Team Velocity",
		KnowledgeTransfer:      "Knowledge Transfer",
		OnboardingTime:         "Onboarding Time",
	}

	if metric, exists := metrics[mt]; exists {
		return metric
	}
	return "Unknown Metric"
}

// Demo function to showcase STAR framework usage
func DemoSTARFramework() {
	fmt.Println("=== STAR Method Framework Demo ===")
	fmt.Println("Elite FAANG L3/L4 Behavioral Interview Preparation")
	fmt.Println()

	// Create framework instance
	framework := NewSTARFramework("Google", "L4")

	// Example technical leadership story
	story := STARStory{
		Title:     "Go Microservices Architecture Migration",
		Category:  TechnicalLeadership,
		Situation: "Our monolithic e-commerce platform was struggling with scalability issues, experiencing 3-second average response times during peak traffic of 50K concurrent users. The engineering team of 25 developers was frustrated with deployment bottlenecks - any code change required full application restart affecting all features.",
		Task:      "As Senior Backend Engineer, I was tasked with leading the architectural migration to microservices while maintaining zero downtime and ensuring the team could continue feature development. The goal was to reduce response times to under 500ms and enable independent service deployments.",
		Action:    "I designed a Go-based microservices architecture using patterns from our Module 03 (concurrency) and Module 04 (system design). First, I implemented a strangler fig pattern to gradually extract services. I built a service mesh with Go using our channel patterns for inter-service communication, implemented circuit breakers for fault tolerance, and created a centralized configuration management system. I mentored 8 engineers on Go best practices, established CI/CD pipelines for each microservice, and implemented comprehensive monitoring with distributed tracing.",
		Result:    "Over 6 months, we successfully migrated to 12 independent Go microservices. Response times improved from 3s to 380ms average (87% improvement). Deployment frequency increased from weekly to multiple times daily per service. The team reported 60% faster feature development cycles, and system availability improved from 99.5% to 99.95%. Zero customer-facing downtime during the migration.",
		Metrics: []ImpactMetric{
			{
				Type:        LatencyReduction,
				Description: "Average response time reduction",
				Value:       87,
				Unit:        "%",
				Baseline:    3000,
				Timeframe:   "6 months",
			},
			{
				Type:        DeploymentFrequency,
				Description: "Deployment frequency increase",
				Value:       500,
				Unit:        "%",
				Timeframe:   "post-migration",
			},
		},
		Technologies: []string{"Go", "microservices", "Docker", "Kubernetes", "gRPC", "circuit breakers"},
		Duration:     6 * 30 * 24 * time.Hour, // 6 months
		TeamSize:     8,
		Keywords:     []string{"leadership", "architecture", "mentoring", "performance"},
	}

	// Add story to framework
	if err := framework.AddStory(story); err != nil {
		fmt.Printf("Error adding story: %v\n", err)
		return
	}

	// Format for interview presentation (2-minute limit)
	presentation := framework.FormatSTARPresentation(story, 2*time.Minute)
	fmt.Println(presentation)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("POTENTIAL FOLLOW-UP QUESTIONS:")
	questions := framework.GenerateFollowUpQuestions(story)
	for i, question := range questions {
		fmt.Printf("%d. %s\n", i+1, question)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("FRAMEWORK ANALYTICS:")
	fmt.Printf("• Total Stories: %d\n", len(framework.stories))
	fmt.Printf("• Target Company: %s\n", framework.targetCompany)
	fmt.Printf("• Target Level: %s\n", framework.targetLevel)
	fmt.Printf("• Technical Keywords: %d\n", len(framework.technicalKeywords))

	// Show company match analysis
	googleStories := framework.GetStoriesForCompany("Google")
	fmt.Printf("• Stories matching Google values: %d\n", len(googleStories))
}
