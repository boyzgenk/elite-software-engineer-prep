// Package behavioral provides a comprehensive leadership story bank with quantified impact metrics
// for FAANG L3/L4 behavioral interviews. Each story leverages technical knowledge from previous modules
// and demonstrates progressive leadership responsibility with measurable business outcomes.
//
// Story Categories:
// - Technical Leadership: Architecture decisions, team coordination, mentoring
// - Problem Solving: Production issues, performance optimization, debugging
// - Innovation: New technologies, process improvements, creative solutions
// - Collaboration: Cross-functional projects, stakeholder management
// - Learning & Growth: Skill development, knowledge sharing, continuous improvement
//
// Target Companies: Google, Meta, Netflix, Amazon, Apple, Microsoft, ByteDance, Uber, Airbnb
// Success Metrics: 7-10 versatile stories covering all common behavioral question patterns
package behavioral

import (
	"fmt"
	"strings"
	"time"
)

// StoryBank contains pre-built STAR stories optimized for different interview scenarios
type StoryBank struct {
	framework *STARFramework
	stories   map[string]STARStory // keyed by story ID for easy retrieval
}

// NewStoryBank creates a story bank with pre-loaded Go-specific leadership stories
func NewStoryBank(targetCompany, targetLevel string) *StoryBank {
	sb := &StoryBank{
		framework: NewSTARFramework(targetCompany, targetLevel),
		stories:   make(map[string]STARStory),
	}

	// Load all pre-built stories
	sb.loadTechnicalLeadershipStories()
	sb.loadProblemSolvingStories()
	sb.loadInnovationStories()
	sb.loadCollaborationStories()
	sb.loadLearningAgilityStories()
	sb.loadFailureRecoveryStories()
	sb.loadInitiativeStories()

	return sb
}

// GetStoryByID retrieves a specific story by its ID
func (sb *StoryBank) GetStoryByID(id string) (STARStory, bool) {
	story, exists := sb.stories[id]
	return story, exists
}

// GetStoriesForCategory returns all stories matching a category
func (sb *StoryBank) GetStoriesForCategory(category StoryCategory) []STARStory {
	return sb.framework.GetStoriesForCategory(category)
}

// GetRecommendedStories returns the best stories for a specific interview context
func (sb *StoryBank) GetRecommendedStories(company string, role string) []STARStory {
	// Get company-specific stories
	companyStories := sb.framework.GetStoriesForCompany(company)

	// Recommend diverse categories for comprehensive coverage
	recommended := make([]STARStory, 0)
	categoriesNeeded := []StoryCategory{
		TechnicalLeadership, ProblemSolving, Innovation,
		Collaboration, LearningAgility, Failure, Initiative,
	}

	for _, category := range categoriesNeeded {
		for _, story := range companyStories {
			if story.Category == category && len(recommended) < 7 {
				recommended = append(recommended, story)
				break
			}
		}
	}

	return recommended
}

// loadTechnicalLeadershipStories adds technical leadership stories to the bank
func (sb *StoryBank) loadTechnicalLeadershipStories() {
	// Story 1: Go Microservices Architecture Migration
	story1 := STARStory{
		Title:     "Go Microservices Architecture Migration",
		Category:  TechnicalLeadership,
		Situation: "Our monolithic e-commerce platform was struggling with scalability bottlenecks, experiencing 4-second average response times during peak traffic of 75K concurrent users. The development team of 28 engineers was frustrated with deployment dependencies - any hotfix required coordinating with 6 different feature teams and a full platform restart affecting all 150+ features. Customer complaints about checkout timeouts increased 300% during Black Friday traffic.",
		Task:      "As Senior Backend Engineer, I was tasked with leading the architectural transformation to microservices while maintaining zero customer-facing downtime and ensuring 28 engineers could continue feature development velocity. Executive stakeholders required <500ms response times and 99.95% availability for the upcoming holiday season in 4 months.",
		Action:    "I designed a comprehensive Go-based microservices architecture leveraging patterns from our concurrency and system design knowledge. First, I implemented a strangler fig pattern to gradually extract 12 core business domains. I architected a service mesh using Go's channel patterns for inter-service communication, implemented circuit breakers with automatic recovery for fault tolerance, and created a centralized configuration management system with hot-reload capabilities. I personally mentored 12 engineers on Go best practices including memory optimization and goroutine patterns, established GitOps CI/CD pipelines for independent service deployments, and implemented distributed tracing with custom Go middleware for end-to-end observability.",
		Result:    "Over 4 months, we successfully migrated to 12 independent Go microservices with zero customer impact. Response times improved from 4s to 320ms average (92% improvement). Deployment frequency increased from weekly coordinated releases to 15+ independent deployments daily across services. Feature development velocity increased 75% as teams became autonomous. System availability improved from 99.2% to 99.97%. The architecture handled Black Friday traffic of 150K concurrent users with 280ms average response time. Engineering satisfaction scores improved from 6.2/10 to 8.7/10.",
		Metrics: []ImpactMetric{
			{
				Type:        LatencyReduction,
				Description: "Average response time reduction",
				Value:       92,
				Unit:        "%",
				Baseline:    4000,
				Timeframe:   "4 months post-migration",
			},
			{
				Type:        DeploymentFrequency,
				Description: "Daily deployment frequency increase",
				Value:       2000,
				Unit:        "%",
				Timeframe:   "post-migration",
			},
			{
				Type:        DeveloperProductivity,
				Description: "Feature development velocity increase",
				Value:       75,
				Unit:        "%",
				Timeframe:   "6 months post-migration",
			},
		},
		Technologies: []string{"Go", "microservices", "Docker", "Kubernetes", "gRPC", "circuit breakers", "distributed tracing"},
		Duration:     4 * 30 * 24 * time.Hour, // 4 months
		TeamSize:     12,
		Keywords:     []string{"leadership", "architecture", "mentoring", "performance", "scalability"},
	}

	// Story 2: Cross-Team Go Performance Optimization
	story2 := STARStory{
		Title:     "Cross-Team Go Performance Optimization Initiative",
		Category:  TechnicalLeadership,
		Situation: "Our Go-based data processing pipeline was becoming a bottleneck across 4 engineering teams. The system processed 2TB of user analytics data daily but was taking 14 hours to complete, causing downstream ML model training delays and impacting real-time recommendation accuracy. Memory usage spiked to 32GB during processing, and the overnight batch jobs were failing 15% of the time due to OOM errors.",
		Task:      "As Technical Lead, I was asked to coordinate across 4 teams (Backend, Data Science, Infrastructure, DevOps) to optimize the entire pipeline while maintaining data accuracy and ensuring all teams could continue their feature work. The goal was sub-8 hour processing time with <16GB memory footprint and 99%+ reliability.",
		Action:    "I conducted a comprehensive performance audit using Go's built-in profiling tools and established weekly cross-team optimization sync meetings. I identified that the main bottlenecks were inefficient memory allocation patterns and blocking I/O operations. I implemented several key optimizations: introduced sync.Pool for object reuse reducing garbage collection pressure by 60%, redesigned the data processing pipeline using worker pools with buffered channels for parallel processing, implemented streaming processing patterns to handle large datasets in chunks, and created a custom memory-mapped file reader for efficient data access. I also established performance benchmarking as part of the CI process and mentored team leads on Go memory management best practices.",
		Result:    "Processing time decreased from 14 hours to 5.5 hours (61% improvement). Memory usage dropped from 32GB to 12GB peak (62% reduction). Pipeline reliability improved from 85% to 99.2%. The optimizations enabled real-time recommendation updates, increasing user engagement by 18%. All 4 teams adopted the performance patterns I documented, leading to a 40% improvement in overall application performance across the organization.",
		Metrics: []ImpactMetric{
			{
				Type:        PerformanceImprovement,
				Description: "Data processing time reduction",
				Value:       61,
				Unit:        "%",
				Baseline:    14,
				Timeframe:   "post-optimization",
			},
			{
				Type:        MemoryOptimization,
				Description: "Peak memory usage reduction",
				Value:       62,
				Unit:        "%",
				Baseline:    32,
				Timeframe:   "sustained improvement",
			},
			{
				Type:        UserGrowth,
				Description: "User engagement increase",
				Value:       18,
				Unit:        "%",
				Timeframe:   "3 months post-optimization",
			},
		},
		Technologies: []string{"Go", "sync.Pool", "goroutines", "channels", "pprof", "memory mapping", "worker pools"},
		Duration:     6 * 7 * 24 * time.Hour, // 6 weeks
		TeamSize:     16,                     // across 4 teams
		Keywords:     []string{"leadership", "performance", "cross-team", "mentoring", "optimization"},
	}

	// Add stories to bank and framework
	sb.stories["microservices_migration"] = story1
	sb.stories["performance_optimization"] = story2
	sb.framework.AddStory(story1)
	sb.framework.AddStory(story2)
}

// loadProblemSolvingStories adds complex problem-solving stories
func (sb *StoryBank) loadProblemSolvingStories() {
	// Story 1: Production Memory Leak Investigation
	story1 := STARStory{
		Title:     "Critical Production Memory Leak Resolution",
		Category:  ProblemSolving,
		Situation: "Our Go-based API gateway serving 50M requests/day started experiencing severe memory leaks in production. Memory usage grew from normal 2GB to 16GB+ over 48 hours, causing automatic restarts every 6 hours and impacting 200K+ active users. The issue appeared suddenly after a routine deployment and was causing cascading failures across our microservices ecosystem. Customer support reported 500% increase in timeout errors.",
		Task:      "As on-call senior engineer, I needed to identify and fix the memory leak while maintaining service availability. The challenge was that the issue only manifested in production under full load and wasn't reproducible in staging environments with simulated traffic.",
		Action:    "I immediately implemented emergency mitigation by reducing restart thresholds and scaling horizontally to handle the load. Then I began systematic investigation using Go's profiling tools. I deployed custom instrumentation to collect heap profiles every 15 minutes and analyzed the memory allocation patterns using pprof. The investigation revealed that a new connection pooling implementation wasn't properly releasing HTTP connections due to a subtle goroutine leak in error handling paths. I traced the issue to missing context cancellation in our circuit breaker implementation. I developed a fix that properly closed connections and canceled goroutines in all error scenarios, then validated the solution using chaos engineering techniques to simulate various failure modes.",
		Result:    "Memory usage returned to normal 2GB baseline within 2 hours of deployment. Service restarts eliminated completely. Customer timeout errors dropped by 95% within 4 hours. I documented the debugging methodology and created automated monitoring alerts to detect similar patterns. The investigation process I developed became the standard operating procedure for performance issues, reducing average resolution time from 8 hours to 2 hours for the entire team.",
		Metrics: []ImpactMetric{
			{
				Type:        MemoryOptimization,
				Description: "Memory usage reduction to baseline",
				Value:       87.5,
				Unit:        "%",
				Baseline:    16,
				Timeframe:   "2 hours post-fix",
			},
			{
				Type:        ErrorReduction,
				Description: "Customer timeout errors reduction",
				Value:       95,
				Unit:        "%",
				Timeframe:   "4 hours post-fix",
			},
			{
				Type:        TimeToMarket,
				Description: "Performance issue resolution time improvement",
				Value:       75,
				Unit:        "%",
				Timeframe:   "team-wide adoption",
			},
		},
		Technologies: []string{"Go", "pprof", "goroutines", "HTTP connection pools", "circuit breakers", "context cancellation"},
		Duration:     18 * time.Hour, // 18 hours investigation + fix
		TeamSize:     3,              // on-call team
		Keywords:     []string{"problem-solving", "debugging", "production", "performance", "methodology"},
	}

	// Story 2: Distributed System Race Condition
	story2 := STARStory{
		Title:     "Complex Race Condition in Distributed Payment System",
		Category:  ProblemSolving,
		Situation: "Our Go-based payment processing system was experiencing intermittent duplicate charges affecting 0.3% of transactions (approximately 1,500 customers daily). The issue was particularly problematic because it involved real money and required manual refunds. The bug was extremely difficult to reproduce - it only occurred under high concurrency with specific timing conditions across 3 different microservices.",
		Task:      "I was assigned to lead the investigation and resolution of this critical issue. The challenge was identifying a race condition across distributed services where the bug only appeared under production load patterns and timing.",
		Action:    "I approached this systematically by first implementing comprehensive distributed tracing to understand the exact sequence of operations across services. Using Go's race detector and custom instrumentation, I identified that the issue occurred when payment validation, authorization, and settlement services processed the same request concurrently due to retry logic. I created a deterministic reproduction environment using controlled goroutine scheduling and discovered the race condition was in our idempotency key implementation - two services could simultaneously check for duplicate transactions and both proceed with processing. I designed a solution using distributed locks with Redis and Go's context package for timeout handling, ensuring atomic payment operations across all services.",
		Result:    "Duplicate payment incidents dropped from 1,500/day to zero within 24 hours of deployment. Customer support tickets related to billing issues decreased by 80%. The distributed locking pattern I implemented became the standard for all financial operations across the platform. I created a race condition detection framework that identified 3 additional potential issues before they reached production, preventing an estimated $50K+ in duplicate charges.",
		Metrics: []ImpactMetric{
			{
				Type:        ErrorReduction,
				Description: "Duplicate payment elimination",
				Value:       100,
				Unit:        "%",
				Baseline:    1500,
				Timeframe:   "24 hours post-deployment",
			},
			{
				Type:        CostSavings,
				Description: "Prevented duplicate charges",
				Value:       50000,
				Unit:        "$",
				Timeframe:   "6 months projected",
			},
			{
				Type:        CodeQualityImprovement,
				Description: "Additional race conditions prevented",
				Value:       3,
				Unit:        "issues",
				Timeframe:   "proactive detection",
			},
		},
		Technologies: []string{"Go", "race detector", "distributed tracing", "Redis", "context", "distributed locks"},
		Duration:     2 * 7 * 24 * time.Hour, // 2 weeks
		TeamSize:     5,
		Keywords:     []string{"problem-solving", "race conditions", "distributed systems", "payments", "investigation"},
	}

	sb.stories["memory_leak_debug"] = story1
	sb.stories["race_condition_fix"] = story2
	sb.framework.AddStory(story1)
	sb.framework.AddStory(story2)
}

// loadInnovationStories adds innovation and technology adoption stories
func (sb *StoryBank) loadInnovationStories() {
	story1 := STARStory{
		Title:     "Go-Based Real-Time Analytics Platform Innovation",
		Category:  Innovation,
		Situation: "Our company was using a third-party analytics service costing $180K annually with 24-hour data latency, preventing real-time decision making. Marketing campaigns couldn't be optimized in real-time, and customer support couldn't access current user behavior data. The existing solution also had strict data volume limits that we were approaching.",
		Task:      "I proposed building an in-house real-time analytics platform using Go to replace the expensive third-party solution. I needed to prove the technical feasibility, design the architecture, and convince leadership that our team could deliver enterprise-grade analytics capabilities.",
		Action:    "I developed a proof-of-concept using Go's concurrency features to process streaming data in real-time. The system used goroutines for parallel event processing, channels for data pipeline coordination, and a custom time-series database optimized for analytics queries. I implemented advanced features like automatic data aggregation, configurable retention policies, and a GraphQL API for flexible querying. I presented benchmarks showing our solution could handle 10x our current data volume while providing sub-second query responses. To gain buy-in, I demonstrated real-time dashboard updates during a live marketing campaign.",
		Result:    "Leadership approved the project, and we delivered the platform in 3 months. The new system eliminated the $180K annual cost, reduced data latency from 24 hours to <5 seconds, and supported 15x our previous data volume. Marketing team reported 35% improvement in campaign performance due to real-time optimization capabilities. The platform has processed over 2 billion events with 99.8% uptime and became a competitive differentiator when selling to enterprise customers.",
		Metrics: []ImpactMetric{
			{
				Type:        CostSavings,
				Description: "Annual third-party service cost elimination",
				Value:       180000,
				Unit:        "$",
				Timeframe:   "annual recurring",
			},
			{
				Type:        LatencyReduction,
				Description: "Data latency improvement",
				Value:       99.97,
				Unit:        "%",
				Baseline:    86400, // 24 hours in seconds
				Timeframe:   "sustained improvement",
			},
			{
				Type:        PerformanceImprovement,
				Description: "Marketing campaign performance increase",
				Value:       35,
				Unit:        "%",
				Timeframe:   "post-implementation",
			},
		},
		Technologies: []string{"Go", "goroutines", "channels", "time-series database", "GraphQL", "real-time processing"},
		Duration:     3 * 30 * 24 * time.Hour, // 3 months
		TeamSize:     4,
		Keywords:     []string{"innovation", "cost-savings", "real-time", "analytics", "leadership-buy-in"},
	}

	sb.stories["analytics_innovation"] = story1
	sb.framework.AddStory(story1)
}

// loadCollaborationStories adds cross-functional collaboration stories
func (sb *StoryBank) loadCollaborationStories() {
	story1 := STARStory{
		Title:     "Cross-Department API Integration Project",
		Category:  Collaboration,
		Situation: "Our Sales, Marketing, and Customer Success teams were using 5 different tools with no data integration, causing duplicate work and inconsistent customer experiences. Sales reps spent 3 hours daily copying data between systems, Marketing couldn't track lead conversion accurately, and Customer Success had incomplete customer interaction history. Leadership mandated a unified customer data platform but the teams had conflicting requirements and different technical comfort levels.",
		Task:      "I was selected to lead the technical integration while collaborating with non-technical stakeholders across 3 departments to design and implement a unified API platform that would satisfy all team requirements without disrupting their existing workflows.",
		Action:    "I organized weekly cross-departmental workshops to understand each team's specific needs and pain points. I designed a Go-based API gateway that could integrate with all 5 existing tools while providing a unified interface. To address varying technical knowledge levels, I created interactive API documentation with real-time examples and implemented webhook notifications so teams could receive data updates in their preferred tools. I worked closely with each department to design custom dashboards and automated workflows. For non-technical users, I built a web interface that translated business logic into API calls. I also implemented robust data validation and rollback capabilities to prevent data corruption across systems.",
		Result:    "Successfully integrated all 5 systems within 6 weeks. Sales team productivity increased by 40% as data entry time dropped from 3 hours to 20 minutes daily. Marketing attribution accuracy improved from 60% to 94%, leading to 25% better budget allocation. Customer Success response time improved by 50% due to complete interaction history. All three department heads endorsed the solution and requested similar integrations for other tools. The API platform now handles 2M+ requests daily with 99.95% uptime.",
		Metrics: []ImpactMetric{
			{
				Type:        DeveloperProductivity,
				Description: "Sales team daily efficiency improvement",
				Value:       40,
				Unit:        "%",
				Timeframe:   "sustained improvement",
			},
			{
				Type:        CodeQualityImprovement,
				Description: "Marketing attribution accuracy improvement",
				Value:       34,
				Unit:        "percentage points",
				Baseline:    60,
				Timeframe:   "post-integration",
			},
			{
				Type:        PerformanceImprovement,
				Description: "Customer Success response time improvement",
				Value:       50,
				Unit:        "%",
				Timeframe:   "post-integration",
			},
		},
		Technologies: []string{"Go", "API gateway", "webhooks", "data integration", "web interface", "real-time dashboards"},
		Duration:     6 * 7 * 24 * time.Hour, // 6 weeks
		TeamSize:     8,                      // across 3 departments
		Keywords:     []string{"collaboration", "cross-functional", "stakeholder-management", "integration", "communication"},
	}

	sb.stories["api_integration_collaboration"] = story1
	sb.framework.AddStory(story1)
}

// loadLearningAgilityStories adds learning and growth stories
func (sb *StoryBank) loadLearningAgilityStories() {
	story1 := STARStory{
		Title:     "Rapid Kubernetes Learning for Production Migration",
		Category:  LearningAgility,
		Situation: "Our company decided to migrate all applications to Kubernetes to improve scalability and reduce infrastructure costs by 40%. However, our team had zero Kubernetes experience, and the migration deadline was set for 3 months to coincide with our largest traffic event (Black Friday). Failure to migrate in time would result in infrastructure costs doubling due to contracted scaling requirements.",
		Task:      "As the team's most senior Go developer, I was asked to lead the Kubernetes adoption and train 8 other engineers while ensuring our Go applications were properly containerized and orchestrated for the Black Friday traffic.",
		Action:    "I created an intensive self-learning plan covering Kubernetes fundamentals, Go container optimization, and production operations. I spent 2-3 hours daily studying official documentation, completing hands-on labs, and building test clusters. I then designed a knowledge transfer program for the team, creating weekly workshops with practical exercises using our actual Go applications. I implemented proper Go application patterns for containerization including graceful shutdown handling, health checks, and configuration management through environment variables. I also established monitoring and logging practices specific to Go applications in Kubernetes environments.",
		Result:    "Successfully completed the migration 2 weeks ahead of schedule. All 12 Go services ran smoothly during Black Friday with 150% higher traffic than previous year. Infrastructure costs decreased by 35% as planned. The entire team became proficient in Kubernetes operations, and I was promoted to Senior Staff Engineer based on the technical leadership and learning velocity demonstrated. The migration patterns I developed became the company standard and have been used for 20+ additional service migrations.",
		Metrics: []ImpactMetric{
			{
				Type:        TimeToMarket,
				Description: "Migration completed ahead of schedule",
				Value:       14,
				Unit:        "days early",
				Timeframe:   "project completion",
			},
			{
				Type:        CostSavings,
				Description: "Infrastructure cost reduction",
				Value:       35,
				Unit:        "%",
				Timeframe:   "annual savings",
			},
			{
				Type:        KnowledgeTransfer,
				Description: "Engineers trained on Kubernetes",
				Value:       8,
				Unit:        "team members",
				Timeframe:   "3 months",
			},
		},
		Technologies: []string{"Kubernetes", "Go", "Docker", "containerization", "monitoring", "graceful shutdown"},
		Duration:     3 * 30 * 24 * time.Hour, // 3 months
		TeamSize:     9,
		Keywords:     []string{"learning-agility", "knowledge-transfer", "self-directed-learning", "migration", "training"},
	}

	sb.stories["kubernetes_learning"] = story1
	sb.framework.AddStory(story1)
}

// loadFailureRecoveryStories adds failure and learning stories
func (sb *StoryBank) loadFailureRecoveryStories() {
	story1 := STARStory{
		Title:     "Database Migration Failure and Recovery",
		Category:  Failure,
		Situation: "I was leading a critical database migration from PostgreSQL to a sharded MongoDB setup to handle our growing 100M+ user dataset. The migration was planned for a 4-hour maintenance window during low traffic. However, 2 hours into the migration, we discovered that our Go application's MongoDB driver configuration was incompatible with the production sharding setup, causing 100% of write operations to fail. 500K+ users were unable to access core features.",
		Task:      "I needed to quickly decide whether to continue debugging the MongoDB issues or rollback to PostgreSQL while minimizing customer impact and data loss. The rollback process itself would take 2 hours, extending the outage to 6+ hours total.",
		Action:    "I immediately called for a rollback while simultaneously working on the MongoDB fix with a separate team. During the 2-hour rollback, I identified that the issue was in our Go connection string configuration - we were missing proper read preference settings for the sharded cluster. I coordinated with our DevOps team to script the proper configuration and worked with the database team to optimize the rollback process. Once back on PostgreSQL, I spent the next week thoroughly testing the MongoDB configuration in a production-mirror environment and identified 3 additional configuration issues that would have caused problems.",
		Result:    "The outage lasted 4.5 hours instead of the potentially 8+ hours if we had continued debugging live. We successfully completed the migration 2 weeks later with zero issues and 30-second downtime for traffic switching. The experience taught me to always have fully tested rollback procedures and to mirror production configurations exactly in staging. I implemented automated configuration validation that has prevented 4 similar issues in subsequent migrations. The migration ultimately delivered the planned 60% read performance improvement.",
		Metrics: []ImpactMetric{
			{
				Type:        TimeToMarket,
				Description: "Outage duration minimization",
				Value:       50,
				Unit:        "% of potential impact",
				Timeframe:   "incident resolution",
			},
			{
				Type:        PerformanceImprovement,
				Description: "Database read performance improvement",
				Value:       60,
				Unit:        "%",
				Timeframe:   "post-successful migration",
			},
			{
				Type:        ErrorReduction,
				Description: "Future migration issues prevented",
				Value:       4,
				Unit:        "incidents",
				Timeframe:   "subsequent migrations",
			},
		},
		Technologies: []string{"Go", "PostgreSQL", "MongoDB", "database migration", "sharding", "configuration management"},
		Duration:     3 * 7 * 24 * time.Hour, // 3 weeks total (including retry)
		TeamSize:     6,
		Keywords:     []string{"failure", "learning", "decision-making", "rollback", "risk-management"},
	}

	sb.stories["database_migration_failure"] = story1
	sb.framework.AddStory(story1)
}

// loadInitiativeStories adds taking initiative stories
func (sb *StoryBank) loadInitiativeStories() {
	story1 := STARStory{
		Title:     "Proactive Security Vulnerability Assessment",
		Category:  Initiative,
		Situation: "While working on a routine feature update, I noticed that our Go applications were using outdated dependency versions and had no systematic vulnerability scanning. Given the increasing cybersecurity threats in our industry, I realized this could pose significant risks to our 2M+ user platform, but security audits weren't part of our current roadmap or budget.",
		Task:      "Without being asked, I decided to conduct a comprehensive security assessment of our Go codebase and implement automated vulnerability monitoring to protect our users and company from potential security breaches.",
		Action:    "I researched Go security best practices and implemented `govulncheck` into our CI/CD pipeline to automatically scan for known vulnerabilities. I audited all 50+ dependencies across our 15 Go services and created a prioritized remediation plan based on CVSS scores and exploit probability. I also implemented security-focused code review guidelines, added automated dependency update monitoring, and created a security dashboard for tracking our security posture. To ensure buy-in, I prepared a presentation for leadership highlighting the potential risks and showing how other companies in our space had been breached through similar vulnerabilities.",
		Result:    "Identified and fixed 12 high-severity vulnerabilities before they could be exploited. Leadership was impressed with the initiative and allocated budget for a formal security program with me as the technical lead. The automated scanning has since prevented 25+ vulnerabilities from reaching production. Our security posture improvement helped us pass SOC 2 compliance 3 months ahead of schedule, enabling $2M+ in enterprise contracts that required compliance certification. The security framework I built has been adopted by 3 other engineering teams.",
		Metrics: []ImpactMetric{
			{
				Type:        ErrorReduction,
				Description: "High-severity vulnerabilities fixed",
				Value:       12,
				Unit:        "vulnerabilities",
				Timeframe:   "initial assessment",
			},
			{
				Type:        RevenueImpact,
				Description: "Enterprise contracts enabled by compliance",
				Value:       2000000,
				Unit:        "$",
				Timeframe:   "6 months post-implementation",
			},
			{
				Type:        TimeToMarket,
				Description: "SOC 2 compliance acceleration",
				Value:       3,
				Unit:        "months ahead of schedule",
				Timeframe:   "compliance achievement",
			},
		},
		Technologies: []string{"Go", "govulncheck", "security scanning", "dependency management", "CI/CD", "vulnerability assessment"},
		Duration:     4 * 7 * 24 * time.Hour, // 4 weeks
		TeamSize:     1,                      // solo initiative
		Keywords:     []string{"initiative", "security", "proactive", "leadership-presentation", "compliance"},
	}

	sb.stories["security_initiative"] = story1
	sb.framework.AddStory(story1)
}

// GetAllStories returns all stories in the bank
func (sb *StoryBank) GetAllStories() []STARStory {
	stories := make([]STARStory, 0, len(sb.stories))
	for _, story := range sb.stories {
		stories = append(stories, story)
	}
	return stories
}

// GenerateStoryBank creates a complete story bank for demonstration
func GenerateStoryBank() *StoryBank {
	return NewStoryBank("Google", "L4")
}

// DemoStoryBank demonstrates the story bank functionality
func DemoStoryBank() {
	fmt.Println("=== FAANG BEHAVIORAL STORY BANK ===")
	fmt.Println("Comprehensive Leadership Stories for L3/L4 Interviews")
	fmt.Println()

	// Create story bank
	storyBank := GenerateStoryBank()

	// Show all available stories
	allStories := storyBank.GetAllStories()
	fmt.Printf("Total Stories Available: %d\n\n", len(allStories))

	// Show stories by category
	categories := []StoryCategory{
		TechnicalLeadership, ProblemSolving, Innovation,
		Collaboration, LearningAgility, Failure, Initiative,
	}

	for _, category := range categories {
		categoryStories := storyBank.GetStoriesForCategory(category)
		fmt.Printf("📚 %s: %d stories\n", category.String(), len(categoryStories))
		for _, story := range categoryStories {
			fmt.Printf("   • %s\n", story.Title)
		}
		fmt.Println()
	}

	// Show recommended stories for different companies
	companies := []string{"Google", "Meta", "Netflix", "Amazon", "Apple"}
	fmt.Println("🎯 RECOMMENDED STORIES BY COMPANY:")
	for _, company := range companies {
		recommended := storyBank.GetRecommendedStories(company, "Senior Software Engineer")
		fmt.Printf("\n%s (L4 Senior Engineer): %d recommended stories\n", company, len(recommended))
		for i, story := range recommended {
			fmt.Printf("   %d. %s [%s]\n", i+1, story.Title, story.Category.String())
		}
	}

	// Show detailed example story
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("EXAMPLE STORY PRESENTATION:")

	exampleStory, exists := storyBank.GetStoryByID("microservices_migration")
	if exists {
		presentation := storyBank.framework.FormatSTARPresentation(exampleStory, 2*time.Minute)
		fmt.Println(presentation)
	}
}
