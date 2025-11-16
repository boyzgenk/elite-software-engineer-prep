package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ProgressData represents the complete tracking data structure
type ProgressData struct {
	StudentID      string                `json:"student_id"`
	StartDate      time.Time             `json:"start_date"`
	LastUpdated    time.Time             `json:"last_updated"`
	WeeklyProgress []WeeklyMetrics       `json:"weekly_progress"`
	LeetCodeStats  LeetCodeProgress      `json:"leetcode_stats"`
	GoMastery      GoMasteryProgress     `json:"go_mastery"`
	SystemDesign   SystemDesignProgress  `json:"system_design"`
	BehavioralPrep BehavioralProgress    `json:"behavioral_prep"`
	MockInterviews []MockInterviewResult `json:"mock_interviews"`
	OverallScore   float64               `json:"overall_score"`
	ReadinessLevel string                `json:"readiness_level"`
	TargetDate     time.Time             `json:"target_interview_date"`
}

// WeeklyMetrics tracks weekly performance
type WeeklyMetrics struct {
	WeekNumber        int       `json:"week_number"`
	StartDate         time.Time `json:"start_date"`
	StudyHours        float64   `json:"study_hours"`
	ProblemsAttempted int       `json:"problems_attempted"`
	ProblemsSolved    int       `json:"problems_solved"`
	GoScore           float64   `json:"go_proficiency_score"`
	AlgorithmScore    float64   `json:"algorithm_score"`
	SystemScore       float64   `json:"system_design_score"`
	BehavioralScore   float64   `json:"behavioral_score"`
	OverallProgress   float64   `json:"overall_progress"`
	Notes             string    `json:"notes"`
}

// LeetCodeProgress tracks algorithm problem solving
type LeetCodeProgress struct {
	TotalSolved    int                    `json:"total_solved"`
	EasySolved     int                    `json:"easy_solved"`
	MediumSolved   int                    `json:"medium_solved"`
	HardSolved     int                    `json:"hard_solved"`
	SuccessRate    float64                `json:"success_rate"`
	AverageTime    float64                `json:"average_time_minutes"`
	PatternMastery map[string]PatternStat `json:"pattern_mastery"`
	RecentProblems []ProblemAttempt       `json:"recent_problems"`
	WeakPatterns   []string               `json:"weak_patterns"`
	StrongPatterns []string               `json:"strong_patterns"`
}

// PatternStat tracks mastery of algorithm patterns
type PatternStat struct {
	Pattern       string    `json:"pattern"`
	Attempted     int       `json:"attempted"`
	Solved        int       `json:"solved"`
	AverageTime   float64   `json:"average_time"`
	MasteryLevel  string    `json:"mastery_level"` // Beginner, Intermediate, Advanced, Expert
	LastPracticed time.Time `json:"last_practiced"`
}

// ProblemAttempt records individual problem attempts
type ProblemAttempt struct {
	Date       time.Time `json:"date"`
	Problem    string    `json:"problem"`
	Difficulty string    `json:"difficulty"`
	Pattern    string    `json:"pattern"`
	TimeSpent  float64   `json:"time_spent_minutes"`
	Solved     bool      `json:"solved"`
	Attempts   int       `json:"attempts"`
	Notes      string    `json:"notes"`
}

// GoMasteryProgress tracks Go language proficiency
type GoMasteryProgress struct {
	RuntimeArchitecture  float64   `json:"runtime_architecture"`
	MemoryManagement     float64   `json:"memory_management"`
	ConcurrencyPatterns  float64   `json:"concurrency_patterns"`
	InterfaceDesign      float64   `json:"interface_design"`
	PerformanceProfiling float64   `json:"performance_profiling"`
	ProductionPatterns   float64   `json:"production_patterns"`
	OverallScore         float64   `json:"overall_score"`
	LastAssessment       time.Time `json:"last_assessment"`
	CompletedExercises   []string  `json:"completed_exercises"`
	CertificationReady   bool      `json:"certification_ready"`
}

// SystemDesignProgress tracks system design capabilities
type SystemDesignProgress struct {
	LoadBalancing    float64          `json:"load_balancing"`
	Caching          float64          `json:"caching"`
	DatabaseDesign   float64          `json:"database_design"`
	Microservices    float64          `json:"microservices"`
	MessageQueues    float64          `json:"message_queues"`
	CDNStorage       float64          `json:"cdn_storage"`
	OverallScore     float64          `json:"overall_score"`
	CompletedDesigns []DesignExercise `json:"completed_designs"`
	ReadyForL3       bool             `json:"ready_for_l3"`
	ReadyForL4       bool             `json:"ready_for_l4"`
}

// DesignExercise tracks completed system design exercises
type DesignExercise struct {
	Date       time.Time `json:"date"`
	Problem    string    `json:"problem"`
	Duration   float64   `json:"duration_minutes"`
	Score      float64   `json:"score"`
	Feedback   string    `json:"feedback"`
	GoSpecific bool      `json:"go_specific_implementation"`
}

// BehavioralProgress tracks interview story preparation
type BehavioralProgress struct {
	StoriesCompleted    int       `json:"stories_completed"`
	StoriesPolished     int       `json:"stories_polished"`
	TechnicalLeadership int       `json:"technical_leadership_stories"`
	ProblemSolving      int       `json:"problem_solving_stories"`
	Innovation          int       `json:"innovation_stories"`
	Collaboration       int       `json:"collaboration_stories"`
	Mentoring           int       `json:"mentoring_stories"`
	FailureLearning     int       `json:"failure_learning_stories"`
	PracticeHours       float64   `json:"practice_hours"`
	MockSessions        int       `json:"mock_sessions"`
	ConfidenceLevel     float64   `json:"confidence_level"`
	CompanySpecificPrep []string  `json:"company_specific_prep"`
	LastPractice        time.Time `json:"last_practice"`
}

// MockInterviewResult records mock interview performance
type MockInterviewResult struct {
	Date            time.Time `json:"date"`
	Interviewer     string    `json:"interviewer"`
	CompanyStyle    string    `json:"company_style"`
	TechnicalScore  float64   `json:"technical_score"`
	BehavioralScore float64   `json:"behavioral_score"`
	OverallScore    float64   `json:"overall_score"`
	Feedback        string    `json:"feedback"`
	AreasToImprove  []string  `json:"areas_to_improve"`
	Strengths       []string  `json:"strengths"`
	ReadyToApply    bool      `json:"ready_to_apply"`
}

// ProgressTracker manages progress tracking operations
type ProgressTracker struct {
	dataFile string
	data     *ProgressData
}

// NewProgressTracker creates a new tracker instance
func NewProgressTracker(studentID string) *ProgressTracker {
	pt := &ProgressTracker{
		dataFile: fmt.Sprintf("progress_%s.json", studentID),
	}

	// Load existing data or create new
	if err := pt.loadData(); err != nil {
		pt.data = &ProgressData{
			StudentID:      studentID,
			StartDate:      time.Now(),
			LastUpdated:    time.Now(),
			WeeklyProgress: make([]WeeklyMetrics, 0),
			LeetCodeStats: LeetCodeProgress{
				PatternMastery: make(map[string]PatternStat),
				RecentProblems: make([]ProblemAttempt, 0),
			},
			MockInterviews: make([]MockInterviewResult, 0),
		}
	}

	return pt
}

// loadData loads existing progress data
func (pt *ProgressTracker) loadData() error {
	if _, err := os.Stat(pt.dataFile); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(pt.dataFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &pt.data)
}

// saveData saves progress data to file
func (pt *ProgressTracker) saveData() error {
	pt.data.LastUpdated = time.Now()

	data, err := json.MarshalIndent(pt.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pt.dataFile, data, 0644)
}

// RecordWeeklyProgress adds weekly metrics
func (pt *ProgressTracker) RecordWeeklyProgress(metrics WeeklyMetrics) error {
	metrics.StartDate = time.Now()

	// Calculate overall progress
	metrics.OverallProgress = (metrics.GoScore + metrics.AlgorithmScore +
		metrics.SystemScore + metrics.BehavioralScore) / 4

	pt.data.WeeklyProgress = append(pt.data.WeeklyProgress, metrics)

	// Update overall score
	pt.updateOverallScore()

	return pt.saveData()
}

// AddLeetCodeProblem records a problem attempt
func (pt *ProgressTracker) AddLeetCodeProblem(attempt ProblemAttempt) error {
	attempt.Date = time.Now()

	// Add to recent problems (keep last 50)
	pt.data.LeetCodeStats.RecentProblems = append(pt.data.LeetCodeStats.RecentProblems, attempt)
	if len(pt.data.LeetCodeStats.RecentProblems) > 50 {
		pt.data.LeetCodeStats.RecentProblems = pt.data.LeetCodeStats.RecentProblems[1:]
	}

	// Update overall stats
	pt.data.LeetCodeStats.TotalSolved++
	switch attempt.Difficulty {
	case "Easy":
		pt.data.LeetCodeStats.EasySolved++
	case "Medium":
		pt.data.LeetCodeStats.MediumSolved++
	case "Hard":
		pt.data.LeetCodeStats.HardSolved++
	}

	// Update pattern mastery
	if pattern, exists := pt.data.LeetCodeStats.PatternMastery[attempt.Pattern]; exists {
		pattern.Attempted++
		if attempt.Solved {
			pattern.Solved++
		}
		pattern.AverageTime = (pattern.AverageTime*float64(pattern.Attempted-1) + attempt.TimeSpent) / float64(pattern.Attempted)
		pattern.LastPracticed = time.Now()
		pt.data.LeetCodeStats.PatternMastery[attempt.Pattern] = pattern
	} else {
		pt.data.LeetCodeStats.PatternMastery[attempt.Pattern] = PatternStat{
			Pattern:       attempt.Pattern,
			Attempted:     1,
			Solved:        boolToInt(attempt.Solved),
			AverageTime:   attempt.TimeSpent,
			MasteryLevel:  "Beginner",
			LastPracticed: time.Now(),
		}
	}

	// Recalculate success rate and average time
	pt.recalculateLeetCodeStats()

	return pt.saveData()
}

// AddMockInterviewResult records mock interview performance
func (pt *ProgressTracker) AddMockInterviewResult(result MockInterviewResult) error {
	result.Date = time.Now()
	pt.data.MockInterviews = append(pt.data.MockInterviews, result)

	// Update overall readiness assessment
	pt.updateOverallScore()

	return pt.saveData()
}

// UpdateGoMastery updates Go proficiency scores
func (pt *ProgressTracker) UpdateGoMastery(mastery GoMasteryProgress) error {
	mastery.LastAssessment = time.Now()
	mastery.OverallScore = (mastery.RuntimeArchitecture + mastery.MemoryManagement +
		mastery.ConcurrencyPatterns + mastery.InterfaceDesign +
		mastery.PerformanceProfiling + mastery.ProductionPatterns) / 6

	mastery.CertificationReady = mastery.OverallScore >= 85.0

	pt.data.GoMastery = mastery
	pt.updateOverallScore()

	return pt.saveData()
}

// UpdateSystemDesignProgress updates system design capabilities
func (pt *ProgressTracker) UpdateSystemDesignProgress(design SystemDesignProgress) error {
	design.OverallScore = (design.LoadBalancing + design.Caching + design.DatabaseDesign +
		design.Microservices + design.MessageQueues + design.CDNStorage) / 6

	design.ReadyForL3 = design.OverallScore >= 75.0
	design.ReadyForL4 = design.OverallScore >= 85.0

	pt.data.SystemDesign = design
	pt.updateOverallScore()

	return pt.saveData()
}

// UpdateBehavioralProgress updates interview preparation progress
func (pt *ProgressTracker) UpdateBehavioralProgress(behavioral BehavioralProgress) error {
	behavioral.LastPractice = time.Now()
	pt.data.BehavioralPrep = behavioral
	pt.updateOverallScore()

	return pt.saveData()
}

// GenerateProgressReport creates a comprehensive progress report
func (pt *ProgressTracker) GenerateProgressReport() string {
	report := fmt.Sprintf("🔥 ELITE FOUNDATION PROGRESS REPORT\n")
	report += fmt.Sprintf("Student: %s\n", pt.data.StudentID)
	report += fmt.Sprintf("Start Date: %s\n", pt.data.StartDate.Format("2006-01-02"))
	report += fmt.Sprintf("Last Updated: %s\n", pt.data.LastUpdated.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("Study Duration: %.0f days\n\n", time.Since(pt.data.StartDate).Hours()/24)

	// Overall Score
	report += fmt.Sprintf("🎯 OVERALL READINESS: %.1f%% (%s)\n", pt.data.OverallScore, pt.data.ReadinessLevel)
	report += "===========================================\n\n"

	// Go Mastery
	report += fmt.Sprintf("🚀 GO MASTERY: %.1f%%\n", pt.data.GoMastery.OverallScore)
	report += fmt.Sprintf("  Runtime Architecture: %.1f%%\n", pt.data.GoMastery.RuntimeArchitecture)
	report += fmt.Sprintf("  Memory Management: %.1f%%\n", pt.data.GoMastery.MemoryManagement)
	report += fmt.Sprintf("  Concurrency Patterns: %.1f%%\n", pt.data.GoMastery.ConcurrencyPatterns)
	report += fmt.Sprintf("  Certification Ready: %v\n\n", pt.data.GoMastery.CertificationReady)

	// LeetCode Progress
	report += fmt.Sprintf("💻 ALGORITHM MASTERY\n")
	report += fmt.Sprintf("  Total Solved: %d (Easy: %d, Medium: %d, Hard: %d)\n",
		pt.data.LeetCodeStats.TotalSolved, pt.data.LeetCodeStats.EasySolved,
		pt.data.LeetCodeStats.MediumSolved, pt.data.LeetCodeStats.HardSolved)
	report += fmt.Sprintf("  Success Rate: %.1f%%\n", pt.data.LeetCodeStats.SuccessRate)
	report += fmt.Sprintf("  Average Time: %.1f minutes\n\n", pt.data.LeetCodeStats.AverageTime)

	// System Design
	report += fmt.Sprintf("🏗️ SYSTEM DESIGN: %.1f%%\n", pt.data.SystemDesign.OverallScore)
	report += fmt.Sprintf("  L3 Ready: %v | L4 Ready: %v\n", pt.data.SystemDesign.ReadyForL3, pt.data.SystemDesign.ReadyForL4)
	report += fmt.Sprintf("  Completed Designs: %d\n\n", len(pt.data.SystemDesign.CompletedDesigns))

	// Behavioral Prep
	report += fmt.Sprintf("🎭 BEHAVIORAL PREP\n")
	report += fmt.Sprintf("  Stories Completed: %d/10\n", pt.data.BehavioralPrep.StoriesCompleted)
	report += fmt.Sprintf("  Practice Hours: %.1f\n", pt.data.BehavioralPrep.PracticeHours)
	report += fmt.Sprintf("  Confidence Level: %.1f/10\n\n", pt.data.BehavioralPrep.ConfidenceLevel)

	// Mock Interviews
	if len(pt.data.MockInterviews) > 0 {
		report += fmt.Sprintf("🎯 MOCK INTERVIEW PERFORMANCE\n")
		report += fmt.Sprintf("  Sessions Completed: %d\n", len(pt.data.MockInterviews))

		// Calculate averages
		var techSum, behavSum, overallSum float64
		for _, mock := range pt.data.MockInterviews {
			techSum += mock.TechnicalScore
			behavSum += mock.BehavioralScore
			overallSum += mock.OverallScore
		}

		count := float64(len(pt.data.MockInterviews))
		report += fmt.Sprintf("  Average Technical: %.1f%%\n", techSum/count)
		report += fmt.Sprintf("  Average Behavioral: %.1f%%\n", behavSum/count)
		report += fmt.Sprintf("  Average Overall: %.1f%%\n\n", overallSum/count)
	}

	// Weekly Progress Trend
	if len(pt.data.WeeklyProgress) > 0 {
		report += fmt.Sprintf("📈 WEEKLY PROGRESS TREND\n")

		recent := pt.data.WeeklyProgress[len(pt.data.WeeklyProgress)-1]
		report += fmt.Sprintf("  Week %d Progress: %.1f%%\n", recent.WeekNumber, recent.OverallProgress)
		report += fmt.Sprintf("  Study Hours: %.1f/35\n", recent.StudyHours)
		report += fmt.Sprintf("  Problems Solved: %d/35\n\n", recent.ProblemsSolved)
	}

	// Recommendations
	report += "🎯 RECOMMENDATIONS:\n"
	if pt.data.OverallScore >= 85 {
		report += "✅ INTERVIEW READY: Schedule your FAANG interviews!\n"
	} else if pt.data.OverallScore >= 75 {
		report += "⚡ ALMOST READY: Focus on weak areas for 2-3 weeks\n"
	} else if pt.data.OverallScore >= 65 {
		report += "📚 DEVELOPING: Continue structured practice for 4-6 weeks\n"
	} else {
		report += "🏗️ FOUNDATION: Focus on fundamentals for 6-8 weeks\n"
	}

	return report
}

// Helper functions
func (pt *ProgressTracker) recalculateLeetCodeStats() {
	solved := 0
	totalTime := 0.0

	for _, problem := range pt.data.LeetCodeStats.RecentProblems {
		if problem.Solved {
			solved++
		}
		totalTime += problem.TimeSpent
	}

	total := len(pt.data.LeetCodeStats.RecentProblems)
	if total > 0 {
		pt.data.LeetCodeStats.SuccessRate = float64(solved) / float64(total) * 100
		pt.data.LeetCodeStats.AverageTime = totalTime / float64(total)
	}
}

func (pt *ProgressTracker) updateOverallScore() {
	goScore := pt.data.GoMastery.OverallScore
	algoScore := pt.data.LeetCodeStats.SuccessRate
	systemScore := pt.data.SystemDesign.OverallScore
	behavioralScore := pt.data.BehavioralPrep.ConfidenceLevel * 10 // Convert to percentage

	pt.data.OverallScore = (goScore + algoScore + systemScore + behavioralScore) / 4

	// Update readiness level
	if pt.data.OverallScore >= 85 {
		pt.data.ReadinessLevel = "🔥 ELITE - Interview Ready"
	} else if pt.data.OverallScore >= 75 {
		pt.data.ReadinessLevel = "⚡ STRONG - Almost Ready"
	} else if pt.data.OverallScore >= 65 {
		pt.data.ReadinessLevel = "📚 DEVELOPING - In Progress"
	} else {
		pt.data.ReadinessLevel = "🏗️ FOUNDATION - Building Skills"
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// CLI interface
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Elite Foundation Progress Tracker")
		fmt.Println("Usage:")
		fmt.Println("  progress-tracker [student-id] [command] [options...]")
		fmt.Println("\nCommands:")
		fmt.Println("  report                    - Generate progress report")
		fmt.Println("  weekly [week] [hours]     - Record weekly progress")
		fmt.Println("  leetcode [problem] [time] [solved] [pattern] - Add problem attempt")
		fmt.Println("  mock [company] [tech-score] [behavioral-score] - Add mock interview")
		fmt.Println("  go-update [scores...]     - Update Go mastery scores")
		return
	}

	studentID := os.Args[1]
	tracker := NewProgressTracker(studentID)

	if len(os.Args) < 3 {
		fmt.Println(tracker.GenerateProgressReport())
		return
	}

	command := os.Args[2]

	switch command {
	case "report":
		fmt.Println(tracker.GenerateProgressReport())

	case "weekly":
		if len(os.Args) < 5 {
			fmt.Println("Usage: weekly [week-number] [study-hours]")
			return
		}

		week, _ := strconv.Atoi(os.Args[3])
		hours, _ := strconv.ParseFloat(os.Args[4], 64)

		metrics := WeeklyMetrics{
			WeekNumber: week,
			StudyHours: hours,
		}

		if err := tracker.RecordWeeklyProgress(metrics); err != nil {
			fmt.Printf("Error recording weekly progress: %v\n", err)
		} else {
			fmt.Printf("✅ Week %d progress recorded\n", week)
		}

	case "leetcode":
		if len(os.Args) < 7 {
			fmt.Println("Usage: leetcode [problem-name] [time-minutes] [solved:true/false] [pattern]")
			return
		}

		problem := os.Args[3]
		timeSpent, _ := strconv.ParseFloat(os.Args[4], 64)
		solved := strings.ToLower(os.Args[5]) == "true"
		pattern := os.Args[6]

		attempt := ProblemAttempt{
			Problem:   problem,
			TimeSpent: timeSpent,
			Solved:    solved,
			Pattern:   pattern,
		}

		if err := tracker.AddLeetCodeProblem(attempt); err != nil {
			fmt.Printf("Error adding problem: %v\n", err)
		} else {
			fmt.Printf("✅ Problem '%s' recorded\n", problem)
		}

	case "mock":
		if len(os.Args) < 6 {
			fmt.Println("Usage: mock [company] [technical-score] [behavioral-score]")
			return
		}

		company := os.Args[3]
		techScore, _ := strconv.ParseFloat(os.Args[4], 64)
		behavScore, _ := strconv.ParseFloat(os.Args[5], 64)

		result := MockInterviewResult{
			CompanyStyle:    company,
			TechnicalScore:  techScore,
			BehavioralScore: behavScore,
			OverallScore:    (techScore + behavScore) / 2,
		}

		if err := tracker.AddMockInterviewResult(result); err != nil {
			fmt.Printf("Error adding mock interview: %v\n", err)
		} else {
			fmt.Printf("✅ Mock interview with %s style recorded\n", company)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
