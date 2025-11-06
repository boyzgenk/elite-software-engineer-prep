package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// StudySession represents a single study session
type StudySession struct {
	Date        time.Time `json:"date"`
	Week        int       `json:"week"`
	Day         int       `json:"day"`
	Topic       string    `json:"topic"`
	Duration    int       `json:"duration_minutes"`
	Activities  []string  `json:"activities"`
	Deliverable string    `json:"deliverable"`
	Status      string    `json:"status"` // planned, in_progress, completed, skipped
	Notes       string    `json:"notes"`
	Rating      int       `json:"self_rating"` // 1-10 scale
}

// WeeklyAssessment tracks weekly progress and readiness
type WeeklyAssessment struct {
	Week                int                    `json:"week"`
	Date                time.Time              `json:"date"`
	TechnicalScore      int                    `json:"technical_score"` // 1-100
	DeliverableComplete bool                   `json:"deliverable_complete"`
	InterviewReadiness  int                    `json:"interview_readiness"` // 1-10
	SkillsMastered      []string               `json:"skills_mastered"`
	SkillsImproving     []string               `json:"skills_improving"`
	NextWeekGoals       []string               `json:"next_week_goals"`
	TimeSpent           int                    `json:"time_spent_minutes"`
	Challenges          []string               `json:"challenges"`
	Reflections         string                 `json:"reflections"`
	ExtraMetrics        map[string]interface{} `json:"extra_metrics"`
}

// StudyPlan represents a complete study plan (14-week, 18-week, etc.)
type StudyPlan struct {
	Name             string             `json:"name"`
	Duration         int                `json:"duration_weeks"`
	StartDate        time.Time          `json:"start_date"`
	TargetCompletion time.Time          `json:"target_completion"`
	Sessions         []StudySession     `json:"sessions"`
	Assessments      []WeeklyAssessment `json:"assessments"`
	OverallProgress  float64            `json:"overall_progress"`
	LastUpdated      time.Time          `json:"last_updated"`
}

// StudyTracker manages study plan progress
type StudyTracker struct {
	dataFile string
	plan     StudyPlan
}

// NewStudyTracker creates a new study tracker
func NewStudyTracker(planName string, dataDir string) *StudyTracker {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dataFile := filepath.Join(dataDir, fmt.Sprintf("%s_progress.json", planName))
	tracker := &StudyTracker{
		dataFile: dataFile,
	}

	tracker.loadOrCreate(planName)
	return tracker
}

// loadOrCreate loads existing progress or creates new study plan
func (st *StudyTracker) loadOrCreate(planName string) {
	if data, err := os.ReadFile(st.dataFile); err == nil {
		if err := json.Unmarshal(data, &st.plan); err != nil {
			log.Printf("Failed to parse existing data: %v", err)
		} else {
			return
		}
	}

	// Create new study plan
	st.plan = StudyPlan{
		Name:        planName,
		StartDate:   time.Now(),
		Sessions:    []StudySession{},
		Assessments: []WeeklyAssessment{},
		LastUpdated: time.Now(),
	}

	switch planName {
	case "14-week-elite-intensive":
		st.createIntensiveSchedule()
	case "18-week-elite-comprehensive":
		st.createComprehensiveSchedule()
	case "elite-professional-track":
		st.createProfessionalSchedule()
	}

	st.save()
}

// createIntensiveSchedule creates the 14-week intensive study schedule
func (st *StudyTracker) createIntensiveSchedule() {
	st.plan.Duration = 14
	st.plan.TargetCompletion = st.plan.StartDate.AddDate(0, 0, 14*7)

	// Phase 1: Go Internals & Runtime Mastery (Weeks 1-5)
	phase1Topics := []string{
		"Memory Management & GC",
		"Goroutine Scheduler (GMP)",
		"Interface Internals & Performance",
		"Channel Implementation & Patterns",
		"Runtime Architecture Deep Dive",
	}

	// Phase 2: Advanced Algorithms & Data Structures (Weeks 6-10)
	phase2Topics := []string{
		"Advanced Tree Structures (AVL, RB)",
		"Graph Algorithms & Applications",
		"Dynamic Programming Mastery",
		"Advanced String Algorithms",
		"Competitive Programming Techniques",
	}

	// Phase 3: Concurrency & System Design (Weeks 11-13)
	phase3Topics := []string{
		"Advanced Concurrency Patterns",
		"Distributed Systems & Go",
		"System Design Case Studies",
	}

	// Phase 4: FAANG Excellence (Week 14)
	phase4Topics := []string{
		"FAANG Interview Mastery",
	}

	allTopics := append(append(append(phase1Topics, phase2Topics...), phase3Topics...), phase4Topics...)

	for week := 1; week <= 14; week++ {
		topic := allTopics[week-1]

		// Create 5 study sessions per week (weekdays)
		for day := 1; day <= 5; day++ {
			session := StudySession{
				Date:     st.plan.StartDate.AddDate(0, 0, (week-1)*7+day-1),
				Week:     week,
				Day:      day,
				Topic:    topic,
				Duration: 240, // 4 hours per day
				Status:   "planned",
			}

			switch day {
			case 1:
				session.Activities = []string{"Theory deep dive", "Concept learning", "Code examples"}
				session.Deliverable = fmt.Sprintf("%s - Theory mastery", topic)
			case 2:
				session.Activities = []string{"Hands-on implementation", "Problem solving", "Tool building"}
				session.Deliverable = fmt.Sprintf("%s - Implementation project", topic)
			case 3:
				session.Activities = []string{"Advanced techniques", "Performance optimization", "Testing"}
				session.Deliverable = fmt.Sprintf("%s - Optimization demo", topic)
			case 4:
				session.Activities = []string{"Production patterns", "Best practices", "Integration"}
				session.Deliverable = fmt.Sprintf("%s - Production code", topic)
			case 5:
				session.Activities = []string{"Capstone project", "Assessment", "Review"}
				session.Deliverable = fmt.Sprintf("%s - Complete mastery demo", topic)
				session.Duration = 300 // 5 hours on Friday
			}

			st.plan.Sessions = append(st.plan.Sessions, session)
		}

		// Create weekend projects (6 hours each day)
		for day := 6; day <= 7; day++ {
			session := StudySession{
				Date:        st.plan.StartDate.AddDate(0, 0, (week-1)*7+day-1),
				Week:        week,
				Day:         day,
				Topic:       topic,
				Duration:    360, // 6 hours
				Activities:  []string{"Extended project work", "Mock interviews", "Assessment"},
				Deliverable: fmt.Sprintf("%s - Weekend project", topic),
				Status:      "planned",
			}
			st.plan.Sessions = append(st.plan.Sessions, session)
		}
	}
}

// createComprehensiveSchedule creates the 18-week comprehensive study schedule
func (st *StudyTracker) createComprehensiveSchedule() {
	st.plan.Duration = 18
	st.plan.TargetCompletion = st.plan.StartDate.AddDate(0, 0, 18*7)

	// Similar to intensive but with more time per topic
	// Implementation would follow similar pattern with different pacing
}

// createProfessionalSchedule creates the professional track schedule
func (st *StudyTracker) createProfessionalSchedule() {
	st.plan.Duration = 22
	st.plan.TargetCompletion = st.plan.StartDate.AddDate(0, 0, 22*7)

	// Designed for working professionals with fewer hours per day
	// Implementation would follow similar pattern with evening/weekend focus
}

// LogSession records progress for a study session
func (st *StudyTracker) LogSession(week, day int, status string, notes string, rating int) error {
	for i := range st.plan.Sessions {
		session := &st.plan.Sessions[i]
		if session.Week == week && session.Day == day {
			session.Status = status
			session.Notes = notes
			session.Rating = rating
			st.plan.LastUpdated = time.Now()
			return st.save()
		}
	}
	return fmt.Errorf("session not found: week %d, day %d", week, day)
}

// AddWeeklyAssessment records weekly assessment results
func (st *StudyTracker) AddWeeklyAssessment(assessment WeeklyAssessment) error {
	assessment.Date = time.Now()

	// Update or add assessment
	found := false
	for i := range st.plan.Assessments {
		if st.plan.Assessments[i].Week == assessment.Week {
			st.plan.Assessments[i] = assessment
			found = true
			break
		}
	}

	if !found {
		st.plan.Assessments = append(st.plan.Assessments, assessment)
		sort.Slice(st.plan.Assessments, func(i, j int) bool {
			return st.plan.Assessments[i].Week < st.plan.Assessments[j].Week
		})
	}

	st.updateOverallProgress()
	st.plan.LastUpdated = time.Now()
	return st.save()
}

// updateOverallProgress calculates overall progress percentage
func (st *StudyTracker) updateOverallProgress() {
	if len(st.plan.Sessions) == 0 {
		st.plan.OverallProgress = 0
		return
	}

	completed := 0
	for _, session := range st.plan.Sessions {
		if session.Status == "completed" {
			completed++
		}
	}

	st.plan.OverallProgress = float64(completed) / float64(len(st.plan.Sessions)) * 100
}

// GetWeeklyReport generates a comprehensive weekly report
func (st *StudyTracker) GetWeeklyReport(week int) map[string]interface{} {
	report := map[string]interface{}{
		"week":               week,
		"sessions_planned":   0,
		"sessions_completed": 0,
		"total_hours":        0,
		"average_rating":     0.0,
		"deliverables":       []string{},
		"challenges":         []string{},
	}

	totalRating := 0
	ratingCount := 0
	var deliverables []string

	for _, session := range st.plan.Sessions {
		if session.Week == week {
			report["sessions_planned"] = report["sessions_planned"].(int) + 1

			if session.Status == "completed" {
				report["sessions_completed"] = report["sessions_completed"].(int) + 1
				report["total_hours"] = report["total_hours"].(int) + session.Duration/60

				if session.Rating > 0 {
					totalRating += session.Rating
					ratingCount++
				}

				if session.Deliverable != "" {
					deliverables = append(deliverables, session.Deliverable)
				}
			}
		}
	}

	if ratingCount > 0 {
		report["average_rating"] = float64(totalRating) / float64(ratingCount)
	}
	report["deliverables"] = deliverables

	// Add assessment data if available
	for _, assessment := range st.plan.Assessments {
		if assessment.Week == week {
			report["technical_score"] = assessment.TechnicalScore
			report["interview_readiness"] = assessment.InterviewReadiness
			report["skills_mastered"] = assessment.SkillsMastered
			report["challenges"] = assessment.Challenges
			break
		}
	}

	return report
}

// GetOverallStats returns comprehensive study statistics
func (st *StudyTracker) GetOverallStats() map[string]interface{} {
	stats := map[string]interface{}{
		"plan_name":           st.plan.Name,
		"start_date":          st.plan.StartDate.Format("2006-01-02"),
		"target_completion":   st.plan.TargetCompletion.Format("2006-01-02"),
		"overall_progress":    st.plan.OverallProgress,
		"weeks_completed":     0,
		"total_hours":         0,
		"average_rating":      0.0,
		"technical_scores":    []int{},
		"interview_readiness": []int{},
	}

	totalHours := 0
	totalRating := 0
	ratingCount := 0
	weeksCompleted := 0

	// Calculate session statistics
	sessionsByWeek := make(map[int]int)
	completedByWeek := make(map[int]int)

	for _, session := range st.plan.Sessions {
		sessionsByWeek[session.Week]++

		if session.Status == "completed" {
			completedByWeek[session.Week]++
			totalHours += session.Duration / 60

			if session.Rating > 0 {
				totalRating += session.Rating
				ratingCount++
			}
		}
	}

	// Count completed weeks (80% completion threshold)
	for week, totalSessions := range sessionsByWeek {
		completed := completedByWeek[week]
		if float64(completed)/float64(totalSessions) >= 0.8 {
			weeksCompleted++
		}
	}

	stats["weeks_completed"] = weeksCompleted
	stats["total_hours"] = totalHours

	if ratingCount > 0 {
		stats["average_rating"] = float64(totalRating) / float64(ratingCount)
	}

	// Assessment statistics
	var techScores, interviewScores []int
	for _, assessment := range st.plan.Assessments {
		techScores = append(techScores, assessment.TechnicalScore)
		interviewScores = append(interviewScores, assessment.InterviewReadiness)
	}

	stats["technical_scores"] = techScores
	stats["interview_readiness"] = interviewScores

	return stats
}

// save persists the study plan to disk
func (st *StudyTracker) save() error {
	data, err := json.MarshalIndent(st.plan, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal study plan: %v", err)
	}

	return os.WriteFile(st.dataFile, data, 0644)
}

// CLI interface for the study tracker
func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Usage: study_tracker init <plan-name>")
			fmt.Println("Plans: 14-week-elite-intensive, 18-week-elite-comprehensive, elite-professional-track")
			return
		}

		planName := os.Args[2]
		tracker := NewStudyTracker(planName, "./study_data")
		fmt.Printf("Initialized study plan: %s\n", planName)
		fmt.Printf("Start date: %s\n", tracker.plan.StartDate.Format("2006-01-02"))
		fmt.Printf("Target completion: %s\n", tracker.plan.TargetCompletion.Format("2006-01-02"))
		fmt.Printf("Total sessions: %d\n", len(tracker.plan.Sessions))

	case "log":
		if len(os.Args) < 7 {
			fmt.Println("Usage: study_tracker log <plan-name> <week> <day> <status> <rating> <notes>")
			fmt.Println("Status: planned, in_progress, completed, skipped")
			fmt.Println("Rating: 1-10")
			return
		}

		planName := os.Args[2]
		week := parseInt(os.Args[3])
		day := parseInt(os.Args[4])
		status := os.Args[5]
		rating := parseInt(os.Args[6])
		notes := ""
		if len(os.Args) > 7 {
			notes = os.Args[7]
		}

		tracker := NewStudyTracker(planName, "./study_data")
		if err := tracker.LogSession(week, day, status, notes, rating); err != nil {
			log.Fatalf("Failed to log session: %v", err)
		}

		fmt.Printf("Logged session: Week %d, Day %d - %s (Rating: %d)\n", week, day, status, rating)

	case "assess":
		if len(os.Args) < 4 {
			fmt.Println("Usage: study_tracker assess <plan-name> <week>")
			return
		}

		planName := os.Args[2]
		week := parseInt(os.Args[3])

		tracker := NewStudyTracker(planName, "./study_data")

		// Interactive assessment input
		assessment := WeeklyAssessment{Week: week}
		fmt.Printf("Weekly Assessment - Week %d\n", week)
		fmt.Print("Technical Score (1-100): ")
		fmt.Scanln(&assessment.TechnicalScore)
		fmt.Print("Interview Readiness (1-10): ")
		fmt.Scanln(&assessment.InterviewReadiness)
		fmt.Print("Time Spent (minutes): ")
		fmt.Scanln(&assessment.TimeSpent)

		if err := tracker.AddWeeklyAssessment(assessment); err != nil {
			log.Fatalf("Failed to add assessment: %v", err)
		}

		fmt.Printf("Assessment recorded for Week %d\n", week)

	case "report":
		if len(os.Args) < 4 {
			fmt.Println("Usage: study_tracker report <plan-name> <week|overall>")
			return
		}

		planName := os.Args[2]
		reportType := os.Args[3]

		tracker := NewStudyTracker(planName, "./study_data")

		if reportType == "overall" {
			stats := tracker.GetOverallStats()
			printOverallReport(stats)
		} else {
			week := parseInt(reportType)
			report := tracker.GetWeeklyReport(week)
			printWeeklyReport(report)
		}

	default:
		printUsage()
	}
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func printUsage() {
	fmt.Println("Elite Study Tracker - FAANG Interview Preparation")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  init <plan-name>                 Initialize a new study plan")
	fmt.Println("  log <plan> <week> <day> <status> <rating> [notes]  Log session progress")
	fmt.Println("  assess <plan> <week>             Record weekly assessment")
	fmt.Println("  report <plan> <week|overall>     Generate progress report")
	fmt.Println("")
	fmt.Println("Study Plans:")
	fmt.Println("  14-week-elite-intensive         25+ hours/week intensive track")
	fmt.Println("  18-week-elite-comprehensive     20+ hours/week comprehensive track")
	fmt.Println("  elite-professional-track        15-20 hours/week working professional track")
}

func printOverallReport(stats map[string]interface{}) {
	fmt.Printf("=== Overall Study Progress Report ===\n")
	fmt.Printf("Plan: %s\n", stats["plan_name"])
	fmt.Printf("Start Date: %s\n", stats["start_date"])
	fmt.Printf("Target Completion: %s\n", stats["target_completion"])
	fmt.Printf("Overall Progress: %.1f%%\n", stats["overall_progress"])
	fmt.Printf("Weeks Completed: %d\n", stats["weeks_completed"])
	fmt.Printf("Total Hours: %d\n", stats["total_hours"])
	fmt.Printf("Average Rating: %.1f/10\n", stats["average_rating"])

	if techScores, ok := stats["technical_scores"].([]int); ok && len(techScores) > 0 {
		fmt.Printf("Technical Scores: %v\n", techScores)
	}

	if interviewScores, ok := stats["interview_readiness"].([]int); ok && len(interviewScores) > 0 {
		fmt.Printf("Interview Readiness: %v\n", interviewScores)
	}
}

func printWeeklyReport(report map[string]interface{}) {
	fmt.Printf("=== Weekly Report - Week %d ===\n", report["week"])
	fmt.Printf("Sessions Completed: %d/%d\n", report["sessions_completed"], report["sessions_planned"])
	fmt.Printf("Total Hours: %d\n", report["total_hours"])
	fmt.Printf("Average Rating: %.1f/10\n", report["average_rating"])

	if deliverables, ok := report["deliverables"].([]string); ok && len(deliverables) > 0 {
		fmt.Printf("Deliverables Completed:\n")
		for _, d := range deliverables {
			fmt.Printf("  - %s\n", d)
		}
	}

	if score, ok := report["technical_score"]; ok {
		fmt.Printf("Technical Score: %d/100\n", score)
	}

	if readiness, ok := report["interview_readiness"]; ok {
		fmt.Printf("Interview Readiness: %d/10\n", readiness)
	}
}
