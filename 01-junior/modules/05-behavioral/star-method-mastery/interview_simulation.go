// Package behavioral implements advanced interview simulation and practice assessment tools
// for comprehensive FAANG L3/L4 behavioral interview preparation.
//
// This module provides:
// - Real-time interview simulation with company-specific scenarios
// - AI-powered answer evaluation and feedback
// - Performance analytics and improvement recommendations
// - Timed practice sessions with realistic pressure
// - Weakest area identification and targeted practice
//
// Key Features:
// - Simulates actual interview conditions with follow-up questions
// - Provides quantified scoring based on STAR methodology
// - Tracks improvement over time with performance metrics
// - Generates personalized practice recommendations
// - Company-specific evaluation criteria and scoring
package behavioral

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// InterviewSimulator provides comprehensive practice and assessment tools
type InterviewSimulator struct {
	companyMatcher *CompanyMatcher
	evaluator      *AnswerEvaluator
	analytics      *PerformanceAnalytics
	sessionHistory []SimulationSession
}

// NewInterviewSimulator creates a fully-featured interview simulation system
func NewInterviewSimulator(storyBank *StoryBank) *InterviewSimulator {
	return &InterviewSimulator{
		companyMatcher: NewCompanyMatcher(storyBank),
		evaluator:      NewAnswerEvaluator(),
		analytics:      NewPerformanceAnalytics(),
		sessionHistory: make([]SimulationSession, 0),
	}
}

// SimulationSession represents a complete practice interview session
type SimulationSession struct {
	ID               string             `json:"id"`
	Company          string             `json:"company"`
	Level            string             `json:"level"`
	Date             time.Time          `json:"date"`
	Duration         time.Duration      `json:"duration"`
	Questions        []QuestionResponse `json:"questions"`
	OverallScore     float64            `json:"overall_score"`
	StrengthAreas    []string           `json:"strength_areas"`
	ImprovementAreas []string           `json:"improvement_areas"`
	CompanyAlignment float64            `json:"company_alignment"`
	RecommendedStory string             `json:"recommended_story"`
	NextSteps        []string           `json:"next_steps"`
}

// QuestionResponse captures a question, answer, and evaluation
type QuestionResponse struct {
	Question        BehavioralQuestion `json:"question"`
	Answer          string             `json:"answer"`
	ResponseTime    time.Duration      `json:"response_time"`
	STARScore       STARScore          `json:"star_score"`
	AlignmentScore  float64            `json:"alignment_score"`
	ImpactScore     float64            `json:"impact_score"`
	ClarityScore    float64            `json:"clarity_score"`
	OverallScore    float64            `json:"overall_score"`
	Feedback        []string           `json:"feedback"`
	FollowUpHandled bool               `json:"followup_handled"`
}

// STARScore provides detailed STAR methodology evaluation
type STARScore struct {
	Situation float64 `json:"situation"` // 0-100
	Task      float64 `json:"task"`      // 0-100
	Action    float64 `json:"action"`    // 0-100
	Result    float64 `json:"result"`    // 0-100
	Overall   float64 `json:"overall"`   // 0-100
}

// AnswerEvaluator provides AI-powered answer assessment
type AnswerEvaluator struct {
	criteriaWeights map[string]float64
	commonKeywords  map[string][]string
}

// NewAnswerEvaluator creates an advanced answer evaluation system
func NewAnswerEvaluator() *AnswerEvaluator {
	return &AnswerEvaluator{
		criteriaWeights: map[string]float64{
			"star_methodology":  0.35,
			"quantified_impact": 0.25,
			"company_alignment": 0.20,
			"clarity":           0.15,
			"authenticity":      0.05,
		},
		commonKeywords: map[string][]string{
			"leadership":      {"led", "managed", "coordinated", "organized", "initiated", "drove"},
			"problem_solving": {"analyzed", "diagnosed", "investigated", "solved", "resolved", "optimized"},
			"collaboration":   {"collaborated", "worked with", "partnered", "coordinated", "aligned"},
			"innovation":      {"created", "developed", "designed", "invented", "pioneered", "improved"},
			"impact":          {"increased", "decreased", "improved", "reduced", "achieved", "delivered"},
		},
	}
}

// PerformanceAnalytics tracks improvement over time
type PerformanceAnalytics struct {
	performanceHistory []SessionMetrics
	weaknessPatterns   map[string]int
	improvementTrends  map[string][]float64
}

// NewPerformanceAnalytics creates a performance tracking system
func NewPerformanceAnalytics() *PerformanceAnalytics {
	return &PerformanceAnalytics{
		performanceHistory: make([]SessionMetrics, 0),
		weaknessPatterns:   make(map[string]int),
		improvementTrends:  make(map[string][]float64),
	}
}

// SessionMetrics captures key performance indicators
type SessionMetrics struct {
	Date                time.Time     `json:"date"`
	OverallScore        float64       `json:"overall_score"`
	STARMethodologyAvg  float64       `json:"star_methodology_avg"`
	CompanyAlignmentAvg float64       `json:"company_alignment_avg"`
	ResponseTimeAvg     time.Duration `json:"response_time_avg"`
	Category            StoryCategory `json:"category"`
	Company             string        `json:"company"`
}

// StartSimulation begins a timed interview practice session
func (is *InterviewSimulator) StartSimulation(company, level string, questionCount int) (*SimulationSession, error) {
	profile, exists := is.companyMatcher.profiles[strings.ToLower(company)]
	if !exists {
		return nil, fmt.Errorf("company profile not found: %s", company)
	}

	// Generate session ID
	sessionID := fmt.Sprintf("%s_%s_%d", company, level, time.Now().Unix())

	session := &SimulationSession{
		ID:        sessionID,
		Company:   company,
		Level:     level,
		Date:      time.Now(),
		Questions: make([]QuestionResponse, 0, questionCount),
	}

	// Select questions based on company profile
	selectedQuestions := is.selectSimulationQuestions(profile, questionCount)

	fmt.Printf("🎯 Starting %s L%s Interview Simulation\n", company, level)
	fmt.Printf("📋 %d questions selected based on %s interview patterns\n", len(selectedQuestions), company)
	fmt.Printf("⏱️  Expected duration: %v\n", profile.InterviewFormat.Duration)
	fmt.Println("🚀 Let's begin! Answer as if you're in a real interview.\n")

	startTime := time.Now()

	for i, question := range selectedQuestions {
		fmt.Printf("Question %d/%d:\n", i+1, len(selectedQuestions))
		fmt.Printf("💼 %s\n", question.Question)
		fmt.Printf("📝 Expected duration: %v\n", question.ExpectedDuration)

		// Simulate getting user input (in real implementation, this would be interactive)
		response := is.simulateUserResponse(question)

		questionStart := time.Now()

		// Evaluate the response
		evaluation := is.evaluator.EvaluateAnswer(response, question, profile)
		responseTime := time.Since(questionStart)

		questionResponse := QuestionResponse{
			Question:        question,
			Answer:          response,
			ResponseTime:    responseTime,
			STARScore:       evaluation.STARScore,
			AlignmentScore:  evaluation.AlignmentScore,
			ImpactScore:     evaluation.ImpactScore,
			ClarityScore:    evaluation.ClarityScore,
			OverallScore:    evaluation.OverallScore,
			Feedback:        evaluation.Feedback,
			FollowUpHandled: evaluation.FollowUpHandled,
		}

		session.Questions = append(session.Questions, questionResponse)

		// Provide immediate feedback
		fmt.Printf("⚡ Quick Feedback:\n")
		fmt.Printf("   Overall Score: %.1f/100\n", evaluation.OverallScore)
		fmt.Printf("   STAR Quality: %.1f/100\n", evaluation.STARScore.Overall)
		fmt.Printf("   Company Fit: %.1f/100\n", evaluation.AlignmentScore)
		if len(evaluation.Feedback) > 0 {
			fmt.Printf("   Key Point: %s\n", evaluation.Feedback[0])
		}
		fmt.Println("")
	}

	session.Duration = time.Since(startTime)
	session = is.finalizeSession(session, profile)

	// Save to history
	is.sessionHistory = append(is.sessionHistory, *session)

	// Update analytics
	is.analytics.RecordSession(*session)

	return session, nil
}

// AnswerEvaluation contains comprehensive answer assessment
type AnswerEvaluation struct {
	STARScore       STARScore `json:"star_score"`
	AlignmentScore  float64   `json:"alignment_score"`
	ImpactScore     float64   `json:"impact_score"`
	ClarityScore    float64   `json:"clarity_score"`
	OverallScore    float64   `json:"overall_score"`
	Feedback        []string  `json:"feedback"`
	FollowUpHandled bool      `json:"followup_handled"`
	KeywordMatches  []string  `json:"keyword_matches"`
	MissingElements []string  `json:"missing_elements"`
	ImprovementTips []string  `json:"improvement_tips"`
}

// EvaluateAnswer provides comprehensive scoring and feedback
func (ae *AnswerEvaluator) EvaluateAnswer(answer string, question BehavioralQuestion, profile CompanyProfile) AnswerEvaluation {
	answerLower := strings.ToLower(answer)

	// Evaluate STAR methodology
	starScore := ae.evaluateSTARMethodology(answerLower)

	// Evaluate company alignment
	alignmentScore := ae.evaluateCompanyAlignment(answerLower, profile)

	// Evaluate quantified impact
	impactScore := ae.evaluateQuantifiedImpact(answerLower)

	// Evaluate clarity and structure
	clarityScore := ae.evaluateClarityAndStructure(answerLower)

	// Calculate overall score
	overallScore := (starScore.Overall * ae.criteriaWeights["star_methodology"]) +
		(alignmentScore * ae.criteriaWeights["company_alignment"]) +
		(impactScore * ae.criteriaWeights["quantified_impact"]) +
		(clarityScore * ae.criteriaWeights["clarity"])

	// Generate feedback
	feedback := ae.generateFeedback(starScore, alignmentScore, impactScore, clarityScore, answerLower, profile)

	// Identify keyword matches
	keywordMatches := ae.identifyKeywordMatches(answerLower)

	// Identify missing elements
	missingElements := ae.identifyMissingElements(starScore, answerLower)

	// Generate improvement tips
	improvementTips := ae.generateImprovementTips(starScore, alignmentScore, impactScore, clarityScore)

	return AnswerEvaluation{
		STARScore:       starScore,
		AlignmentScore:  alignmentScore,
		ImpactScore:     impactScore,
		ClarityScore:    clarityScore,
		OverallScore:    overallScore,
		Feedback:        feedback,
		FollowUpHandled: ae.hasFollowUpPreparation(answerLower),
		KeywordMatches:  keywordMatches,
		MissingElements: missingElements,
		ImprovementTips: improvementTips,
	}
}

// evaluateSTARMethodology scores STAR structure compliance
func (ae *AnswerEvaluator) evaluateSTARMethodology(answer string) STARScore {
	score := STARScore{}

	// Situation indicators
	situationKeywords := []string{"situation", "context", "background", "at the time", "project", "when", "faced with"}
	situationCount := 0
	for _, keyword := range situationKeywords {
		if strings.Contains(answer, keyword) {
			situationCount++
		}
	}
	score.Situation = minFloat(float64(situationCount*25), 100.0)

	// Task indicators
	taskKeywords := []string{"task", "goal", "objective", "needed to", "responsible for", "assigned", "challenge"}
	taskCount := 0
	for _, keyword := range taskKeywords {
		if strings.Contains(answer, keyword) {
			taskCount++
		}
	}
	score.Task = minFloat(float64(taskCount*25), 100.0)

	// Action indicators (most important)
	actionKeywords := []string{"i ", "i implemented", "i designed", "i led", "i analyzed", "i created", "i solved", "my approach", "i decided"}
	actionCount := 0
	for _, keyword := range actionKeywords {
		if strings.Contains(answer, keyword) {
			actionCount++
		}
	}
	score.Action = minFloat(float64(actionCount*20), 100.0)

	// Result indicators
	resultKeywords := []string{"result", "outcome", "achieved", "delivered", "improved", "reduced", "increased", "saved", "impact"}
	resultCount := 0
	for _, keyword := range resultKeywords {
		if strings.Contains(answer, keyword) {
			resultCount++
		}
	}
	score.Result = minFloat(float64(resultCount*25), 100.0)

	// Overall STAR score
	score.Overall = (score.Situation + score.Task + score.Action + score.Result) / 4

	return score
}

// evaluateCompanyAlignment scores alignment with company values
func (ae *AnswerEvaluator) evaluateCompanyAlignment(answer string, profile CompanyProfile) float64 {
	score := 0.0
	totalPossible := 0.0

	// Check leadership principles alignment
	for _, principle := range profile.LeadershipPrinciples {
		totalPossible += 1.0
		principleWords := strings.Fields(strings.ToLower(principle))
		for _, word := range principleWords {
			if len(word) > 3 && strings.Contains(answer, word) {
				score += 1.0
				break
			}
		}
	}

	// Check cultural values alignment
	for _, value := range profile.CulturalValues {
		totalPossible += 1.0
		if strings.Contains(answer, strings.ToLower(value)) {
			score += 1.0
		}
	}

	// Check keyword weights
	for keyword, weight := range profile.KeywordWeights {
		totalPossible += weight
		if strings.Contains(answer, strings.ToLower(keyword)) {
			score += weight
		}
	}

	if totalPossible == 0 {
		return 50.0 // Default neutral score
	}

	return minFloat((score/totalPossible)*100.0, 100.0)
}

// evaluateQuantifiedImpact scores use of metrics and quantified results
func (ae *AnswerEvaluator) evaluateQuantifiedImpact(answer string) float64 {
	score := 0.0

	// Look for numbers and percentages
	numberKeywords := []string{"%", "percent", "times", "million", "thousand", "days", "hours", "users", "customers"}
	for _, keyword := range numberKeywords {
		if strings.Contains(answer, keyword) {
			score += 20.0
		}
	}

	// Look for improvement verbs with quantification
	improvementKeywords := []string{"increased by", "decreased by", "improved by", "reduced by", "saved", "generated"}
	for _, keyword := range improvementKeywords {
		if strings.Contains(answer, keyword) {
			score += 25.0
		}
	}

	// Look for business impact indicators
	businessKeywords := []string{"revenue", "cost", "efficiency", "performance", "productivity", "quality"}
	for _, keyword := range businessKeywords {
		if strings.Contains(answer, keyword) {
			score += 15.0
		}
	}

	return minFloat(score, 100.0)
}

// evaluateClarityAndStructure scores answer organization and clarity
func (ae *AnswerEvaluator) evaluateClarityAndStructure(answer string) float64 {
	score := 0.0

	// Check for logical connectors
	connectors := []string{"first", "then", "next", "finally", "because", "therefore", "as a result", "consequently"}
	for _, connector := range connectors {
		if strings.Contains(answer, connector) {
			score += 12.5
		}
	}

	// Check answer length (not too short, not too long)
	words := strings.Fields(answer)
	wordCount := len(words)
	if wordCount >= 50 && wordCount <= 200 {
		score += 25.0
	} else if wordCount >= 30 && wordCount <= 300 {
		score += 15.0
	}

	// Check for specific examples
	exampleKeywords := []string{"for example", "specifically", "in particular", "such as"}
	for _, keyword := range exampleKeywords {
		if strings.Contains(answer, keyword) {
			score += 12.5
		}
	}

	return minFloat(score, 100.0)
}

// generateFeedback creates actionable improvement suggestions
func (ae *AnswerEvaluator) generateFeedback(starScore STARScore, alignment, impact, clarity float64, answer string, profile CompanyProfile) []string {
	feedback := make([]string, 0)

	// STAR methodology feedback
	if starScore.Overall < 70 {
		if starScore.Situation < 50 {
			feedback = append(feedback, "Add more context about the situation and circumstances")
		}
		if starScore.Task < 50 {
			feedback = append(feedback, "Clarify your specific role and what you were tasked to accomplish")
		}
		if starScore.Action < 50 {
			feedback = append(feedback, "Focus more on YOUR specific actions and decisions")
		}
		if starScore.Result < 50 {
			feedback = append(feedback, "Include quantified results and business impact")
		}
	}

	// Company alignment feedback
	if alignment < 60 {
		feedback = append(feedback, fmt.Sprintf("Better align your story with %s's values and principles", profile.Name))
	}

	// Impact feedback
	if impact < 50 {
		feedback = append(feedback, "Include specific metrics and quantified business impact")
	}

	// Clarity feedback
	if clarity < 60 {
		feedback = append(feedback, "Improve structure with clearer transitions between STAR elements")
	}

	// Positive reinforcement
	if starScore.Overall >= 80 {
		feedback = append(feedback, "Excellent STAR methodology - well structured response!")
	}
	if alignment >= 80 {
		feedback = append(feedback, fmt.Sprintf("Strong alignment with %s culture and values", profile.Name))
	}

	return feedback
}

// identifyKeywordMatches finds positive keyword alignments
func (ae *AnswerEvaluator) identifyKeywordMatches(answer string) []string {
	matches := make([]string, 0)

	for category, keywords := range ae.commonKeywords {
		for _, keyword := range keywords {
			if strings.Contains(answer, keyword) {
				matches = append(matches, fmt.Sprintf("%s: %s", category, keyword))
			}
		}
	}

	return matches
}

// identifyMissingElements finds gaps in the STAR response
func (ae *AnswerEvaluator) identifyMissingElements(starScore STARScore, answer string) []string {
	missing := make([]string, 0)

	if starScore.Situation < 50 {
		missing = append(missing, "Situation context and background")
	}
	if starScore.Task < 50 {
		missing = append(missing, "Specific task or challenge description")
	}
	if starScore.Action < 50 {
		missing = append(missing, "Detailed personal actions taken")
	}
	if starScore.Result < 50 {
		missing = append(missing, "Quantified results and impact")
	}

	// Check for common missing elements
	if !strings.Contains(answer, "team") && !strings.Contains(answer, "collaborate") {
		missing = append(missing, "Team collaboration or stakeholder management")
	}

	if !strings.Contains(answer, "learn") && !strings.Contains(answer, "challenge") {
		missing = append(missing, "Learning or growth from the experience")
	}

	return missing
}

// generateImprovementTips provides specific actionable advice
func (ae *AnswerEvaluator) generateImprovementTips(starScore STARScore, alignment, impact, clarity float64) []string {
	tips := make([]string, 0)

	if starScore.Overall < 70 {
		tips = append(tips, "Practice the STAR method: Situation, Task, Action, Result")
		tips = append(tips, "Use a timer to practice 2-3 minute responses")
	}

	if impact < 60 {
		tips = append(tips, "Quantify everything: percentages, dollar amounts, time saved")
		tips = append(tips, "Focus on business outcomes, not just technical achievements")
	}

	if alignment < 60 {
		tips = append(tips, "Research company values and weave them into your stories")
		tips = append(tips, "Use company-specific language and priorities")
	}

	if clarity < 60 {
		tips = append(tips, "Practice storytelling with clear beginning, middle, and end")
		tips = append(tips, "Use transition phrases to guide the interviewer through your story")
	}

	return tips
}

// hasFollowUpPreparation checks if answer anticipates follow-up questions
func (ae *AnswerEvaluator) hasFollowUpPreparation(answer string) bool {
	followUpIndicators := []string{
		"would do differently", "lessons learned", "if i had to do it again",
		"what i learned", "alternative approach", "other options",
	}

	for _, indicator := range followUpIndicators {
		if strings.Contains(answer, indicator) {
			return true
		}
	}

	return false
}

// selectSimulationQuestions chooses appropriate questions for practice
func (is *InterviewSimulator) selectSimulationQuestions(profile CompanyProfile, count int) []BehavioralQuestion {
	questions := make([]BehavioralQuestion, 0)

	// Prioritize high-frequency questions
	highFreq := make([]BehavioralQuestion, 0)
	mediumFreq := make([]BehavioralQuestion, 0)

	for _, q := range profile.CommonQuestions {
		switch q.Frequency {
		case "High":
			highFreq = append(highFreq, q)
		case "Medium":
			mediumFreq = append(mediumFreq, q)
		}
	}

	// Select questions ensuring category diversity
	usedCategories := make(map[StoryCategory]bool)

	// First, include high-frequency questions
	for _, q := range highFreq {
		if len(questions) >= count {
			break
		}
		if !usedCategories[q.Category] {
			questions = append(questions, q)
			usedCategories[q.Category] = true
		}
	}

	// Fill remaining with medium frequency
	for _, q := range mediumFreq {
		if len(questions) >= count {
			break
		}
		questions = append(questions, q)
	}

	// If we need more questions, add remaining high-frequency ones
	for _, q := range highFreq {
		if len(questions) >= count {
			break
		}
		alreadyIncluded := false
		for _, existing := range questions {
			if existing.Question == q.Question {
				alreadyIncluded = true
				break
			}
		}
		if !alreadyIncluded {
			questions = append(questions, q)
		}
	}

	return questions
}

// simulateUserResponse generates realistic practice responses
func (is *InterviewSimulator) simulateUserResponse(question BehavioralQuestion) string {
	// In a real implementation, this would capture user input
	// For demo purposes, we'll generate sample responses based on question category

	responses := map[StoryCategory]string{
		TechnicalLeadership: "In my previous role as a backend engineer, I was tasked with leading the migration of our monolithic Go application to microservices. The situation was that our system was experiencing performance issues with 10,000+ concurrent users. My task was to design and lead the migration while maintaining zero downtime. I took the approach of implementing a gradual migration strategy, starting with the user authentication service. I designed the new architecture, set up proper monitoring, and led a team of 4 engineers through the implementation. As a result, we achieved 40% better response times, 99.9% uptime, and reduced deployment time from 2 hours to 15 minutes.",

		ProblemSolving: "We had a critical production issue where our Go application was consuming excessive memory, causing frequent crashes. The situation was affecting 50,000+ users during peak hours. My task was to identify and resolve the memory leak within 24 hours. I systematically analyzed the application using pprof, identified that our caching layer wasn't properly releasing connections, and implemented a proper connection pooling strategy. The result was a 70% reduction in memory usage and zero crashes over the following month.",

		Innovation: "I noticed our team was spending 3-4 hours daily on manual deployment processes. The situation was slowing down our delivery velocity significantly. My task was to find a way to streamline this process. I researched and implemented a CI/CD pipeline using GitHub Actions and Docker, created comprehensive testing automation, and trained the team on the new process. As a result, we reduced deployment time from 4 hours to 10 minutes and increased our deployment frequency from weekly to daily.",

		Collaboration: "During a critical project integration, I needed to coordinate between our Go backend team and the frontend React team who had different perspectives on the API design. The situation was creating delays and tension between teams. My task was to facilitate alignment and deliver the integration on time. I organized daily sync meetings, created detailed API documentation, set up a shared testing environment, and ensured both teams had clear communication channels. The result was successful on-time delivery and improved cross-team collaboration processes.",

		Failure: "I made a critical error during a database migration that caused 2 hours of downtime for our Go application. The situation affected thousands of users and put our SLA at risk. My task was to take ownership, restore service, and prevent future occurrences. I immediately rolled back the changes, coordinated with the incident response team, and conducted a thorough post-mortem analysis. I then implemented automated migration testing and created detailed rollback procedures. As a result, we prevented similar incidents and improved our deployment safety by 300%.",
	}

	if response, exists := responses[question.Category]; exists {
		return response
	}

	return "I successfully handled a challenging situation in my previous role working with Go applications, where I took specific actions that led to measurable business impact and team success."
}

// finalizeSession completes session scoring and recommendations
func (is *InterviewSimulator) finalizeSession(session *SimulationSession, profile CompanyProfile) *SimulationSession {
	totalScore := 0.0
	strengthAreas := make(map[string]int)
	improvementAreas := make(map[string]int)

	// Calculate overall metrics
	for _, response := range session.Questions {
		totalScore += response.OverallScore

		// Track strengths (scores >= 75)
		if response.STARScore.Overall >= 75 {
			strengthAreas["STAR Methodology"]++
		}
		if response.AlignmentScore >= 75 {
			strengthAreas["Company Alignment"]++
		}
		if response.ImpactScore >= 75 {
			strengthAreas["Quantified Impact"]++
		}
		if response.ClarityScore >= 75 {
			strengthAreas["Communication Clarity"]++
		}

		// Track improvement areas (scores < 60)
		if response.STARScore.Overall < 60 {
			improvementAreas["STAR Methodology"]++
		}
		if response.AlignmentScore < 60 {
			improvementAreas["Company Alignment"]++
		}
		if response.ImpactScore < 60 {
			improvementAreas["Quantified Impact"]++
		}
		if response.ClarityScore < 60 {
			improvementAreas["Communication Clarity"]++
		}
	}

	session.OverallScore = totalScore / float64(len(session.Questions))

	// Convert maps to sorted slices
	session.StrengthAreas = is.mapToSortedSlice(strengthAreas)
	session.ImprovementAreas = is.mapToSortedSlice(improvementAreas)

	// Calculate company alignment
	alignmentTotal := 0.0
	for _, response := range session.Questions {
		alignmentTotal += response.AlignmentScore
	}
	session.CompanyAlignment = alignmentTotal / float64(len(session.Questions))

	// Generate next steps
	session.NextSteps = is.generateNextSteps(session, profile)

	return session
}

// mapToSortedSlice converts a count map to sorted slice
func (is *InterviewSimulator) mapToSortedSlice(m map[string]int) []string {
	type kv struct {
		key   string
		value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].value > ss[j].value
	})

	result := make([]string, 0, len(ss))
	for _, kv := range ss {
		result = append(result, kv.key)
	}

	return result
}

// generateNextSteps creates personalized improvement recommendations
func (is *InterviewSimulator) generateNextSteps(session *SimulationSession, profile CompanyProfile) []string {
	steps := make([]string, 0)

	// Based on overall score
	if session.OverallScore >= 85 {
		steps = append(steps, "Excellent performance! Focus on advanced scenarios and leadership examples")
	} else if session.OverallScore >= 70 {
		steps = append(steps, "Strong foundation. Practice specific company scenarios and edge cases")
	} else if session.OverallScore >= 55 {
		steps = append(steps, "Good progress. Focus on STAR methodology and quantifying impact")
	} else {
		steps = append(steps, "Practice fundamental STAR structure with simple, clear examples")
	}

	// Based on company alignment
	if session.CompanyAlignment < 60 {
		steps = append(steps, fmt.Sprintf("Study %s's leadership principles and cultural values more deeply", session.Company))
	}

	// Based on improvement areas
	if len(session.ImprovementAreas) > 0 {
		for i, area := range session.ImprovementAreas {
			if i >= 2 { // Limit to top 2 improvement areas
				break
			}
			switch area {
			case "STAR Methodology":
				steps = append(steps, "Practice STAR structure with detailed situation setup and specific actions")
			case "Quantified Impact":
				steps = append(steps, "Prepare stories with specific metrics, percentages, and business outcomes")
			case "Company Alignment":
				steps = append(steps, "Research company values and incorporate relevant principles into stories")
			case "Communication Clarity":
				steps = append(steps, "Practice storytelling with clear transitions and logical flow")
			}
		}
	}

	// Specific company recommendations
	steps = append(steps, fmt.Sprintf("Schedule mock interviews focusing on %s-specific scenarios", session.Company))

	return steps
}

// RecordSession adds session data to performance analytics
func (pa *PerformanceAnalytics) RecordSession(session SimulationSession) {
	// Calculate averages
	starAvg := 0.0
	alignmentAvg := 0.0
	responseTimeAvg := time.Duration(0)

	for _, response := range session.Questions {
		starAvg += response.STARScore.Overall
		alignmentAvg += response.AlignmentScore
		responseTimeAvg += response.ResponseTime
	}

	if len(session.Questions) > 0 {
		starAvg /= float64(len(session.Questions))
		alignmentAvg /= float64(len(session.Questions))
		responseTimeAvg /= time.Duration(len(session.Questions))
	}

	metrics := SessionMetrics{
		Date:                session.Date,
		OverallScore:        session.OverallScore,
		STARMethodologyAvg:  starAvg,
		CompanyAlignmentAvg: alignmentAvg,
		ResponseTimeAvg:     responseTimeAvg,
		Company:             session.Company,
	}

	pa.performanceHistory = append(pa.performanceHistory, metrics)

	// Track weakness patterns
	for _, area := range session.ImprovementAreas {
		pa.weaknessPatterns[area]++
	}

	// Track improvement trends
	if len(pa.performanceHistory) > 1 {
		for metric, scores := range pa.improvementTrends {
			switch metric {
			case "overall":
				scores = append(scores, session.OverallScore)
			case "star":
				scores = append(scores, starAvg)
			case "alignment":
				scores = append(scores, alignmentAvg)
			}
			pa.improvementTrends[metric] = scores
		}
	} else {
		// Initialize trends
		pa.improvementTrends["overall"] = []float64{session.OverallScore}
		pa.improvementTrends["star"] = []float64{starAvg}
		pa.improvementTrends["alignment"] = []float64{alignmentAvg}
	}
}

// GetImprovementReport generates comprehensive progress analysis
func (pa *PerformanceAnalytics) GetImprovementReport() *ImprovementReport {
	if len(pa.performanceHistory) == 0 {
		return &ImprovementReport{
			Message: "Complete your first practice session to see improvement analytics",
		}
	}

	report := &ImprovementReport{
		SessionCount:         len(pa.performanceHistory),
		CurrentScore:         pa.performanceHistory[len(pa.performanceHistory)-1].OverallScore,
		PersistentWeaknesses: pa.getPersistentWeaknesses(),
		ImprovementTrends:    pa.calculateTrends(),
		Recommendations:      pa.generateRecommendations(),
	}

	if len(pa.performanceHistory) > 1 {
		report.ScoreImprovement = report.CurrentScore - pa.performanceHistory[0].OverallScore
	}

	return report
}

// ImprovementReport provides comprehensive progress analysis
type ImprovementReport struct {
	SessionCount         int               `json:"session_count"`
	CurrentScore         float64           `json:"current_score"`
	ScoreImprovement     float64           `json:"score_improvement"`
	PersistentWeaknesses []string          `json:"persistent_weaknesses"`
	ImprovementTrends    map[string]string `json:"improvement_trends"`
	Recommendations      []string          `json:"recommendations"`
	Message              string            `json:"message,omitempty"`
}

// getPersistentWeaknesses identifies consistently problematic areas
func (pa *PerformanceAnalytics) getPersistentWeaknesses() []string {
	persistent := make([]string, 0)

	for weakness, count := range pa.weaknessPatterns {
		if count >= 3 { // Appears in 3+ sessions
			persistent = append(persistent, weakness)
		}
	}

	return persistent
}

// calculateTrends analyzes improvement trajectory
func (pa *PerformanceAnalytics) calculateTrends() map[string]string {
	trends := make(map[string]string)

	for metric, scores := range pa.improvementTrends {
		if len(scores) < 2 {
			trends[metric] = "insufficient_data"
			continue
		}

		recent := scores[len(scores)-3:]
		if len(recent) < 2 {
			recent = scores
		}

		improvement := recent[len(recent)-1] - recent[0]

		switch {
		case improvement > 5:
			trends[metric] = "improving"
		case improvement < -5:
			trends[metric] = "declining"
		default:
			trends[metric] = "stable"
		}
	}

	return trends
}

// generateRecommendations creates personalized improvement advice
func (pa *PerformanceAnalytics) generateRecommendations() []string {
	recommendations := make([]string, 0)

	if len(pa.performanceHistory) < 3 {
		recommendations = append(recommendations, "Complete at least 3 practice sessions for meaningful analytics")
	}

	// Most recent session analysis
	if len(pa.performanceHistory) > 0 {
		recent := pa.performanceHistory[len(pa.performanceHistory)-1]

		if recent.STARMethodologyAvg < 60 {
			recommendations = append(recommendations, "Focus on STAR methodology structure in your next practice session")
		}

		if recent.CompanyAlignmentAvg < 60 {
			recommendations = append(recommendations, "Research and practice company-specific value alignment")
		}

		if recent.ResponseTimeAvg > 10*time.Minute {
			recommendations = append(recommendations, "Practice with timer to improve response conciseness")
		}
	}

	// Persistent weakness recommendations
	for _, weakness := range pa.getPersistentWeaknesses() {
		recommendations = append(recommendations, fmt.Sprintf("Address persistent weakness: %s", weakness))
	}

	return recommendations
}

// minFloat helper function for float64
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
