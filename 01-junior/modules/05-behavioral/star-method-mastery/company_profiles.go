// Package behavioral implements company-specific behavioral interview preparation
// for FAANG L3/L4 positions with deep cultural alignment and value-based storytelling.
//
// This module provides:
// - Company-specific leadership principles and cultural values
// - Tailored story selection algorithms based on company priorities
// - Interview question patterns and follow-up preparation
// - Real interviewer feedback and evaluation criteria
// - Compensation and negotiation strategies aligned with company culture
//
// Target Companies: Google, Meta, Netflix, Amazon, Apple
// Target Levels: L3/L4 (Software Engineer I/II) with leadership potential
// Success Metrics: Stories that resonate with company values and demonstrate cultural fit
package behavioral

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// CompanyProfile defines company-specific behavioral interview characteristics
type CompanyProfile struct {
	Name                 string                `json:"name"`
	LevelMapping         map[string]string     `json:"level_mapping"` // External level to internal level
	LeadershipPrinciples []string              `json:"leadership_principles"`
	CulturalValues       []string              `json:"cultural_values"`
	InterviewFormat      InterviewFormat       `json:"interview_format"`
	CommonQuestions      []BehavioralQuestion  `json:"common_questions"`
	EvaluationCriteria   []EvaluationCriterion `json:"evaluation_criteria"`
	RedFlags             []string              `json:"red_flags"`       // Things to avoid mentioning
	KeywordWeights       map[string]float64    `json:"keyword_weights"` // Importance of specific keywords
	CompensationInfo     CompensationInfo      `json:"compensation_info"`
}

// InterviewFormat describes company-specific interview structure
type InterviewFormat struct {
	Duration           time.Duration   `json:"duration"`
	NumberOfQuestions  int             `json:"number_of_questions"`
	TimePerQuestion    time.Duration   `json:"time_per_question"`
	FollowUpIntensity  string          `json:"followup_intensity"` // Low, Medium, High
	TechnicalMix       bool            `json:"technical_mix"`      // Whether they mix technical and behavioral
	PanelStyle         bool            `json:"panel_style"`        // Single interviewer vs panel
	RequiredStoryTypes []StoryCategory `json:"required_story_types"`
}

// BehavioralQuestion represents company-specific interview questions
type BehavioralQuestion struct {
	Question          string        `json:"question"`
	Category          StoryCategory `json:"category"`
	Frequency         string        `json:"frequency"`       // How often asked (High, Medium, Low)
	LevelRelevance    []string      `json:"level_relevance"` // Which levels this applies to
	FollowUpQuestions []string      `json:"followup_questions"`
	ExpectedDuration  time.Duration `json:"expected_duration"`
	ScoringCriteria   []string      `json:"scoring_criteria"`
}

// EvaluationCriterion defines how companies evaluate behavioral responses
type EvaluationCriterion struct {
	Name        string   `json:"name"`
	Weight      float64  `json:"weight"` // 0.0 - 1.0
	Description string   `json:"description"`
	Examples    []string `json:"examples"` // What good looks like
}

// CompensationInfo provides company-specific compensation insights
type CompensationInfo struct {
	BaseSalaryRange  []int     `json:"base_salary_range"` // [min, max] in USD
	EquityPercentage []float64 `json:"equity_percentage"` // [min, max] as percentage
	BonusStructure   string    `json:"bonus_structure"`
	NegotiationTips  []string  `json:"negotiation_tips"`
	BehavioralImpact string    `json:"behavioral_impact"` // How behavioral performance affects offers
}

// CompanyMatcher provides intelligent company-specific story selection and preparation
type CompanyMatcher struct {
	profiles    map[string]CompanyProfile
	storyBank   *StoryBank
	preferences MatchingPreferences
}

// MatchingPreferences allows customization of story selection algorithms
type MatchingPreferences struct {
	PreferQuantifiedImpact bool            `json:"prefer_quantified_impact"`
	MinimumTeamSize        int             `json:"minimum_team_size"`
	RequiredTechnologies   []string        `json:"required_technologies"`
	PreferredTimeframe     time.Duration   `json:"preferred_timeframe"`
	AvoidCategories        []StoryCategory `json:"avoid_categories"`
}

// NewCompanyMatcher creates a company matcher with pre-loaded FAANG profiles
func NewCompanyMatcher(storyBank *StoryBank) *CompanyMatcher {
	cm := &CompanyMatcher{
		profiles:  make(map[string]CompanyProfile),
		storyBank: storyBank,
		preferences: MatchingPreferences{
			PreferQuantifiedImpact: true,
			MinimumTeamSize:        2,
			RequiredTechnologies:   []string{"Go"},
			PreferredTimeframe:     6 * 30 * 24 * time.Hour, // 6 months
		},
	}

	// Load all company profiles
	cm.loadGoogleProfile()
	cm.loadMetaProfile()
	cm.loadNetflixProfile()
	cm.loadAmazonProfile()
	cm.loadAppleProfile()

	return cm
}

// scoredStory holds a story with its calculated score
type scoredStory struct {
	story STARStory
	score float64
}

// GetOptimalStories returns the best stories for a specific company and role
func (cm *CompanyMatcher) GetOptimalStories(companyName, level string, maxStories int) ([]STARStory, error) {
	profile, exists := cm.profiles[strings.ToLower(companyName)]
	if !exists {
		return nil, fmt.Errorf("company profile not found: %s", companyName)
	}

	// Get all available stories
	allStories := cm.storyBank.GetAllStories()

	// Score each story against company profile
	scoredStories := make([]scoredStory, 0, len(allStories))
	for _, story := range allStories {
		score := cm.calculateStoryScore(story, profile, level)
		if score > 0.3 { // Minimum threshold
			scoredStories = append(scoredStories, scoredStory{story, score})
		}
	}

	// Sort by score (highest first)
	sort.Slice(scoredStories, func(i, j int) bool {
		return scoredStories[i].score > scoredStories[j].score
	})

	// Ensure diversity across categories
	selectedStories := cm.ensureCategoryDiversity(scoredStories, profile.InterviewFormat.RequiredStoryTypes, maxStories)

	// Extract just the stories
	result := make([]STARStory, len(selectedStories))
	for i, scored := range selectedStories {
		result[i] = scored.story
	}

	return result, nil
}

// PrepareForCompany provides comprehensive interview preparation for a specific company
func (cm *CompanyMatcher) PrepareForCompany(companyName, level string) (*InterviewPreparation, error) {
	profile, exists := cm.profiles[strings.ToLower(companyName)]
	if !exists {
		return nil, fmt.Errorf("company profile not found: %s", companyName)
	}

	// Get optimal stories
	stories, err := cm.GetOptimalStories(companyName, level, 7)
	if err != nil {
		return nil, err
	}

	// Generate question-specific preparation
	questionPrep := make([]QuestionPreparation, 0, len(profile.CommonQuestions))
	for _, question := range profile.CommonQuestions {
		// Find best story for this question
		bestStory := cm.findBestStoryForQuestion(stories, question)
		prep := QuestionPreparation{
			Question:         question,
			RecommendedStory: bestStory,
			KeyPoints:        cm.generateKeyPoints(bestStory, question, profile),
			PitfallsToAvoid:  cm.generatePitfalls(question, profile),
		}
		questionPrep = append(questionPrep, prep)
	}

	preparation := &InterviewPreparation{
		Company:              profile,
		RecommendedStories:   stories,
		QuestionPreparation:  questionPrep,
		CompanySpecificTips:  cm.generateCompanyTips(profile),
		MockInterviewScript:  cm.generateMockScript(profile, stories),
		CompensationStrategy: cm.generateCompensationStrategy(profile, level),
	}

	return preparation, nil
}

// InterviewPreparation contains complete preparation materials for a company
type InterviewPreparation struct {
	Company              CompanyProfile        `json:"company"`
	RecommendedStories   []STARStory           `json:"recommended_stories"`
	QuestionPreparation  []QuestionPreparation `json:"question_preparation"`
	CompanySpecificTips  []string              `json:"company_specific_tips"`
	MockInterviewScript  MockInterviewScript   `json:"mock_interview_script"`
	CompensationStrategy CompensationStrategy  `json:"compensation_strategy"`
}

// QuestionPreparation provides specific guidance for individual questions
type QuestionPreparation struct {
	Question         BehavioralQuestion `json:"question"`
	RecommendedStory STARStory          `json:"recommended_story"`
	KeyPoints        []string           `json:"key_points"`
	PitfallsToAvoid  []string           `json:"pitfalls_to_avoid"`
}

// MockInterviewScript provides realistic practice scenarios
type MockInterviewScript struct {
	Duration        time.Duration        `json:"duration"`
	Questions       []BehavioralQuestion `json:"questions"`
	ExpectedAnswers []string             `json:"expected_answers"`
	TimingGuidance  []time.Duration      `json:"timing_guidance"`
	EvaluationTips  []string             `json:"evaluation_tips"`
}

// CompensationStrategy provides negotiation guidance
type CompensationStrategy struct {
	ExpectedOffer      CompensationOffer `json:"expected_offer"`
	NegotiationTactics []string          `json:"negotiation_tactics"`
	LeveragePoints     []string          `json:"leverage_points"`
	BehavioralImpact   string            `json:"behavioral_impact"`
}

// CompensationOffer represents expected compensation details
type CompensationOffer struct {
	BaseSalary      int     `json:"base_salary"`
	Equity          float64 `json:"equity"`
	Bonus           int     `json:"bonus"`
	TotalComp       int     `json:"total_comp"`
	VestingSchedule string  `json:"vesting_schedule"`
}

// calculateStoryScore evaluates how well a story fits a company profile
func (cm *CompanyMatcher) calculateStoryScore(story STARStory, profile CompanyProfile, level string) float64 {
	score := 0.0

	// Leadership principles alignment (40% weight)
	principleScore := cm.calculatePrincipleAlignment(story, profile.LeadershipPrinciples)
	score += principleScore * 0.4

	// Cultural values alignment (30% weight)
	cultureScore := cm.calculateCultureAlignment(story, profile.CulturalValues)
	score += cultureScore * 0.3

	// Keyword weight matching (20% weight)
	keywordScore := cm.calculateKeywordAlignment(story, profile.KeywordWeights)
	score += keywordScore * 0.2

	// Story quality factors (10% weight)
	qualityScore := cm.calculateQualityScore(story)
	score += qualityScore * 0.1

	// Level appropriateness modifier
	levelModifier := cm.calculateLevelModifier(story, level)
	score *= levelModifier

	return score
}

// calculatePrincipleAlignment measures story alignment with leadership principles
func (cm *CompanyMatcher) calculatePrincipleAlignment(story STARStory, principles []string) float64 {
	storyText := strings.ToLower(story.Situation + " " + story.Task + " " + story.Action + " " + story.Result)

	matches := 0
	for _, principle := range principles {
		principleKeywords := cm.extractPrincipleKeywords(principle)
		for _, keyword := range principleKeywords {
			if strings.Contains(storyText, strings.ToLower(keyword)) {
				matches++
				break
			}
		}
	}

	return float64(matches) / float64(len(principles))
}

// calculateCultureAlignment measures story alignment with cultural values
func (cm *CompanyMatcher) calculateCultureAlignment(story STARStory, values []string) float64 {
	storyText := strings.ToLower(story.Situation + " " + story.Task + " " + story.Action + " " + story.Result)

	matches := 0
	for _, value := range values {
		if strings.Contains(storyText, strings.ToLower(value)) {
			matches++
		}
	}

	return float64(matches) / float64(len(values))
}

// calculateKeywordAlignment measures story alignment with weighted keywords
func (cm *CompanyMatcher) calculateKeywordAlignment(story STARStory, keywordWeights map[string]float64) float64 {
	storyText := strings.ToLower(story.Situation + " " + story.Task + " " + story.Action + " " + story.Result)

	totalWeight := 0.0
	alignedWeight := 0.0

	for keyword, weight := range keywordWeights {
		totalWeight += weight
		if strings.Contains(storyText, strings.ToLower(keyword)) {
			alignedWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0
	}
	return alignedWeight / totalWeight
}

// calculateQualityScore evaluates intrinsic story quality
func (cm *CompanyMatcher) calculateQualityScore(story STARStory) float64 {
	score := 0.0

	// Quantified metrics presence (high value)
	if len(story.Metrics) > 0 {
		score += 0.3
		if len(story.Metrics) >= 3 {
			score += 0.2 // Bonus for multiple metrics
		}
	}

	// Team size appropriateness
	if story.TeamSize >= cm.preferences.MinimumTeamSize {
		score += 0.2
	}

	// Technology alignment
	hasRequiredTech := false
	for _, tech := range cm.preferences.RequiredTechnologies {
		for _, storyTech := range story.Technologies {
			if strings.Contains(strings.ToLower(storyTech), strings.ToLower(tech)) {
				hasRequiredTech = true
				break
			}
		}
		if hasRequiredTech {
			break
		}
	}
	if hasRequiredTech {
		score += 0.2
	}

	// Story duration appropriateness
	if story.Duration <= cm.preferences.PreferredTimeframe {
		score += 0.3
	}

	return score
}

// calculateLevelModifier adjusts score based on level appropriateness
func (cm *CompanyMatcher) calculateLevelModifier(story STARStory, level string) float64 {
	// Base modifier
	modifier := 1.0

	// Adjust based on team size and responsibility level
	switch level {
	case "L3", "E3", "IC3":
		// Junior levels prefer smaller teams and more individual contribution
		if story.TeamSize > 10 {
			modifier *= 0.8
		}
	case "L4", "E4", "IC4":
		// Mid levels prefer moderate team sizes and leadership opportunities
		if story.TeamSize >= 5 && story.TeamSize <= 15 {
			modifier *= 1.2
		}
	case "L5", "E5", "IC5":
		// Senior levels prefer larger teams and significant impact
		if story.TeamSize >= 10 {
			modifier *= 1.3
		}
	}

	return modifier
}

// ensureCategoryDiversity ensures selected stories cover required categories
func (cm *CompanyMatcher) ensureCategoryDiversity(scoredStories []scoredStory, requiredCategories []StoryCategory, maxStories int) []scoredStory {
	selected := make([]scoredStory, 0, maxStories)
	usedCategories := make(map[StoryCategory]bool)

	// First pass: ensure all required categories are covered
	for _, category := range requiredCategories {
		for _, scored := range scoredStories {
			if scored.story.Category == category && !usedCategories[category] && len(selected) < maxStories {
				selected = append(selected, scored)
				usedCategories[category] = true
				break
			}
		}
	}

	// Second pass: fill remaining slots with highest-scoring stories
	for _, scored := range scoredStories {
		if len(selected) >= maxStories {
			break
		}

		// Skip if we already have this exact story
		alreadySelected := false
		for _, sel := range selected {
			if sel.story.Title == scored.story.Title {
				alreadySelected = true
				break
			}
		}

		if !alreadySelected {
			selected = append(selected, scored)
		}
	}

	return selected
}

// extractPrincipleKeywords converts leadership principles to searchable keywords
func (cm *CompanyMatcher) extractPrincipleKeywords(principle string) []string {
	principleMap := map[string][]string{
		"customer obsession":                 {"customer", "user", "client", "service", "satisfaction"},
		"ownership":                          {"responsibility", "accountability", "ownership", "commitment", "follow-through"},
		"invent and simplify":                {"innovation", "creativity", "simplification", "efficiency", "optimization"},
		"are right, a lot":                   {"decision", "judgment", "analysis", "critical thinking", "problem-solving"},
		"learn and be curious":               {"learning", "growth", "curiosity", "development", "exploration"},
		"hire and develop the best":          {"mentoring", "coaching", "talent", "development", "hiring"},
		"insist on the highest standards":    {"quality", "excellence", "standards", "precision", "attention to detail"},
		"think big":                          {"vision", "strategy", "ambitious", "scale", "long-term"},
		"bias for action":                    {"speed", "execution", "initiative", "urgency", "results"},
		"frugality":                          {"efficiency", "cost-conscious", "resource optimization", "lean"},
		"earn trust":                         {"integrity", "transparency", "reliability", "trust", "honesty"},
		"dive deep":                          {"analysis", "investigation", "thorough", "deep-dive", "root cause"},
		"have backbone; disagree and commit": {"feedback", "disagreement", "conviction", "commitment"},
		"deliver results":                    {"results", "outcomes", "delivery", "achievement", "impact"},
		// Google values
		"focus on the user":                                   {"user experience", "user-centric", "customer focus"},
		"it's best to do one thing really well":               {"focus", "specialization", "excellence"},
		"fast is better than slow":                            {"speed", "velocity", "quick", "rapid"},
		"democracy on the web works":                          {"collaboration", "open", "transparency"},
		"you don't need to be at your desk to need an answer": {"remote", "accessibility", "availability"},
		"you can make money without doing evil":               {"ethics", "integrity", "values"},
		"there's always more information out there":           {"data-driven", "research", "information"},
		"the need for information crosses all borders":        {"global", "universal", "accessibility"},
		"you can be serious without a suit":                   {"culture", "informal", "authenticity"},
		"great just isn't good enough":                        {"excellence", "continuous improvement", "perfection"},
	}

	principle = strings.ToLower(principle)
	if keywords, exists := principleMap[principle]; exists {
		return keywords
	}

	// Fallback: extract words from the principle itself
	return strings.Fields(principle)
}

// findBestStoryForQuestion selects the most appropriate story for a specific question
func (cm *CompanyMatcher) findBestStoryForQuestion(stories []STARStory, question BehavioralQuestion) STARStory {
	bestScore := 0.0
	var bestStory STARStory

	for _, story := range stories {
		score := 0.0

		// Category match (highest weight)
		if story.Category == question.Category {
			score += 1.0
		}

		// Duration appropriateness
		if story.Duration <= question.ExpectedDuration*2 { // Allow some flexibility
			score += 0.3
		}

		// Keyword alignment with question
		questionText := strings.ToLower(question.Question)
		storyText := strings.ToLower(story.Title + " " + story.Action + " " + story.Result)

		questionWords := strings.Fields(questionText)
		for _, word := range questionWords {
			if len(word) > 3 && strings.Contains(storyText, word) { // Skip short words
				score += 0.1
			}
		}

		if score > bestScore {
			bestScore = score
			bestStory = story
		}
	}

	return bestStory
}

// generateKeyPoints creates talking points for a story-question combination
func (cm *CompanyMatcher) generateKeyPoints(story STARStory, question BehavioralQuestion, profile CompanyProfile) []string {
	points := []string{
		fmt.Sprintf("Emphasize the %s aspect of your story", question.Category.String()),
		"Highlight quantified business impact and metrics",
		"Connect your actions to the company's leadership principles",
	}

	// Add company-specific points
	switch strings.ToLower(profile.Name) {
	case "google":
		points = append(points, "Demonstrate user focus and technical innovation")
	case "meta":
		points = append(points, "Show how you moved fast and built social value")
	case "netflix":
		points = append(points, "Emphasize high performance and taking ownership")
	case "amazon":
		points = append(points, "Connect to specific leadership principles")
	case "apple":
		points = append(points, "Highlight attention to detail and user experience")
	}

	// Add metric-specific points
	if len(story.Metrics) > 0 {
		for _, metric := range story.Metrics {
			points = append(points, fmt.Sprintf("Mention the %s improvement of %.1f%s",
				metric.Description, metric.Value, metric.Unit))
		}
	}

	return points
}

// generatePitfalls identifies common mistakes to avoid for a question
func (cm *CompanyMatcher) generatePitfalls(question BehavioralQuestion, profile CompanyProfile) []string {
	pitfalls := []string{
		"Don't speak negatively about team members or previous companies",
		"Avoid taking credit for others' work",
		"Don't provide vague answers without specific examples",
		"Avoid focusing only on technical details without business impact",
	}

	// Add company-specific pitfalls
	pitfalls = append(pitfalls, profile.RedFlags...)

	// Add question-specific pitfalls
	switch question.Category {
	case TechnicalLeadership:
		pitfalls = append(pitfalls, "Don't demonstrate micromanagement or lack of delegation")
	case ProblemSolving:
		pitfalls = append(pitfalls, "Don't skip explaining your thought process and methodology")
	case Failure:
		pitfalls = append(pitfalls, "Don't blame others or avoid taking responsibility")
	case Collaboration:
		pitfalls = append(pitfalls, "Don't minimize others' contributions or show poor communication")
	}

	return pitfalls
}

// generateCompanyTips provides company-specific interview guidance
func (cm *CompanyMatcher) generateCompanyTips(profile CompanyProfile) []string {
	baseTips := []string{
		fmt.Sprintf("Study %s's leadership principles and weave them into your answers", profile.Name),
		"Prepare specific examples that demonstrate cultural alignment",
		"Practice timing - stick to the expected duration per question",
		"Research your interviewers on LinkedIn if possible",
	}

	// Add format-specific tips
	if profile.InterviewFormat.PanelStyle {
		baseTips = append(baseTips, "Address all panel members, not just the question asker")
	}

	if profile.InterviewFormat.TechnicalMix {
		baseTips = append(baseTips, "Be prepared to discuss technical details of your behavioral stories")
	}

	switch profile.InterviewFormat.FollowUpIntensity {
	case "High":
		baseTips = append(baseTips, "Expect deep follow-up questions - prepare details for every aspect of your stories")
	case "Medium":
		baseTips = append(baseTips, "Be ready for 2-3 follow-up questions per story")
	case "Low":
		baseTips = append(baseTips, "Keep answers concise - limited follow-up expected")
	}

	return baseTips
}

// generateMockScript creates a realistic practice interview
func (cm *CompanyMatcher) generateMockScript(profile CompanyProfile, stories []STARStory) MockInterviewScript {
	script := MockInterviewScript{
		Duration:       profile.InterviewFormat.Duration,
		Questions:      make([]BehavioralQuestion, 0, profile.InterviewFormat.NumberOfQuestions),
		TimingGuidance: make([]time.Duration, 0),
	}

	// Select questions based on frequency and level
	highFreqQuestions := make([]BehavioralQuestion, 0)
	mediumFreqQuestions := make([]BehavioralQuestion, 0)

	for _, q := range profile.CommonQuestions {
		switch q.Frequency {
		case "High":
			highFreqQuestions = append(highFreqQuestions, q)
		case "Medium":
			mediumFreqQuestions = append(mediumFreqQuestions, q)
		}
	}

	// Prioritize high-frequency questions
	questionsToInclude := min(len(highFreqQuestions), profile.InterviewFormat.NumberOfQuestions)
	script.Questions = append(script.Questions, highFreqQuestions[:questionsToInclude]...)

	// Fill remaining slots with medium frequency
	remaining := profile.InterviewFormat.NumberOfQuestions - len(script.Questions)
	if remaining > 0 && len(mediumFreqQuestions) > 0 {
		toAdd := min(remaining, len(mediumFreqQuestions))
		script.Questions = append(script.Questions, mediumFreqQuestions[:toAdd]...)
	}

	// Generate timing guidance
	for _, q := range script.Questions {
		script.TimingGuidance = append(script.TimingGuidance, q.ExpectedDuration)
	}

	// Generate evaluation tips
	script.EvaluationTips = []string{
		"Focus on specific, quantified examples",
		"Demonstrate leadership and initiative",
		"Show alignment with company values",
		"Maintain appropriate pacing and energy",
	}

	return script
}

// generateCompensationStrategy creates negotiation guidance
func (cm *CompanyMatcher) generateCompensationStrategy(profile CompanyProfile, level string) CompensationStrategy {
	// Calculate expected offer based on level and company
	expectedOffer := cm.calculateExpectedOffer(profile.CompensationInfo, level)

	strategy := CompensationStrategy{
		ExpectedOffer: expectedOffer,
		NegotiationTactics: []string{
			"Demonstrate value through specific accomplishments",
			"Research market rates for your level and location",
			"Consider total compensation, not just base salary",
			"Be prepared to discuss competing offers if available",
		},
		LeveragePoints: []string{
			"Strong behavioral interview performance",
			"Relevant Go and distributed systems experience",
			"Leadership potential demonstrated through stories",
			"Cultural fit with company values",
		},
		BehavioralImpact: profile.CompensationInfo.BehavioralImpact,
	}

	// Add company-specific tactics
	strategy.NegotiationTactics = append(strategy.NegotiationTactics, profile.CompensationInfo.NegotiationTips...)

	return strategy
}

// calculateExpectedOffer estimates compensation based on level and company
func (cm *CompanyMatcher) calculateExpectedOffer(info CompensationInfo, level string) CompensationOffer {
	// Simplified calculation - in practice this would be much more sophisticated
	baseSalary := (info.BaseSalaryRange[0] + info.BaseSalaryRange[1]) / 2
	equity := (info.EquityPercentage[0] + info.EquityPercentage[1]) / 2

	// Adjust based on level
	switch level {
	case "L3", "E3", "IC3":
		baseSalary = int(float64(baseSalary) * 0.9)
		equity *= 0.8
	case "L5", "E5", "IC5":
		baseSalary = int(float64(baseSalary) * 1.2)
		equity *= 1.3
	}

	return CompensationOffer{
		BaseSalary:      baseSalary,
		Equity:          equity,
		Bonus:           baseSalary / 10,                // Rough estimate
		TotalComp:       int(float64(baseSalary) * 1.4), // Including equity and bonus
		VestingSchedule: "4 years with 1-year cliff",
	}
}

// Company profile loading methods

// loadGoogleProfile loads Google's behavioral interview profile
func (cm *CompanyMatcher) loadGoogleProfile() {
	profile := CompanyProfile{
		Name: "Google",
		LevelMapping: map[string]string{
			"L3": "Software Engineer I",
			"L4": "Software Engineer II",
			"L5": "Senior Software Engineer",
		},
		LeadershipPrinciples: []string{
			"focus on the user",
			"it's best to do one thing really well",
			"fast is better than slow",
			"democracy on the web works",
			"you can make money without doing evil",
			"there's always more information out there",
			"great just isn't good enough",
		},
		CulturalValues: []string{
			"innovation", "user-centric", "data-driven", "collaboration",
			"technical excellence", "scalability", "long-term thinking",
		},
		InterviewFormat: InterviewFormat{
			Duration:           45 * time.Minute,
			NumberOfQuestions:  4,
			TimePerQuestion:    10 * time.Minute,
			FollowUpIntensity:  "High",
			TechnicalMix:       true,
			PanelStyle:         false,
			RequiredStoryTypes: []StoryCategory{TechnicalLeadership, ProblemSolving, Innovation, Collaboration},
		},
		CommonQuestions: []BehavioralQuestion{
			{
				Question:          "Tell me about a time you had to lead a technical project",
				Category:          TechnicalLeadership,
				Frequency:         "High",
				LevelRelevance:    []string{"L4", "L5"},
				ExpectedDuration:  8 * time.Minute,
				FollowUpQuestions: []string{"How did you handle disagreements?", "What would you do differently?"},
			},
			{
				Question:         "Describe a complex problem you solved",
				Category:         ProblemSolving,
				Frequency:        "High",
				LevelRelevance:   []string{"L3", "L4", "L5"},
				ExpectedDuration: 7 * time.Minute,
			},
		},
		RedFlags: []string{
			"work-life balance", "easy projects", "avoiding challenges",
			"not data-driven", "poor user focus",
		},
		KeywordWeights: map[string]float64{
			"scale":       1.0,
			"user":        0.9,
			"innovation":  0.8,
			"data":        0.7,
			"performance": 0.8,
		},
		CompensationInfo: CompensationInfo{
			BaseSalaryRange:  []int{150000, 250000},
			EquityPercentage: []float64{0.1, 0.5},
			BonusStructure:   "15% target bonus",
			BehavioralImpact: "Strong behavioral performance can increase offer by 10-20%",
		},
	}
	cm.profiles["google"] = profile
}

// loadMetaProfile loads Meta's behavioral interview profile
func (cm *CompanyMatcher) loadMetaProfile() {
	profile := CompanyProfile{
		Name: "Meta",
		LevelMapping: map[string]string{
			"E3": "Software Engineer I",
			"E4": "Software Engineer II",
			"E5": "Senior Software Engineer",
		},
		LeadershipPrinciples: []string{
			"move fast", "be bold", "focus on impact", "be open",
			"build social value", "meta, metamates, me",
		},
		CulturalValues: []string{
			"impact", "connection", "bold thinking", "moving fast",
			"openness", "social good", "building community",
		},
		InterviewFormat: InterviewFormat{
			Duration:           45 * time.Minute,
			NumberOfQuestions:  3,
			TimePerQuestion:    12 * time.Minute,
			FollowUpIntensity:  "Medium",
			TechnicalMix:       false,
			PanelStyle:         false,
			RequiredStoryTypes: []StoryCategory{TechnicalLeadership, Innovation, Collaboration},
		},
		RedFlags: []string{
			"slow decision making", "perfectionism over progress",
			"individual over team", "closed mindset",
		},
		KeywordWeights: map[string]float64{
			"impact":     1.0,
			"fast":       0.9,
			"community":  0.8,
			"connection": 0.7,
		},
		CompensationInfo: CompensationInfo{
			BaseSalaryRange:  []int{160000, 270000},
			EquityPercentage: []float64{0.2, 0.8},
			BehavioralImpact: "Meta heavily weights behavioral fit - can make or break offer",
		},
	}
	cm.profiles["meta"] = profile
}

// loadNetflixProfile loads Netflix's behavioral interview profile
func (cm *CompanyMatcher) loadNetflixProfile() {
	profile := CompanyProfile{
		Name: "Netflix",
		CulturalValues: []string{
			"high performance", "freedom", "responsibility", "candor",
			"innovation", "curiosity", "courage", "passion",
		},
		InterviewFormat: InterviewFormat{
			Duration:           60 * time.Minute,
			NumberOfQuestions:  5,
			TimePerQuestion:    10 * time.Minute,
			FollowUpIntensity:  "High",
			TechnicalMix:       true,
			PanelStyle:         true,
			RequiredStoryTypes: []StoryCategory{TechnicalLeadership, ProblemSolving, Initiative, Failure},
		},
		RedFlags: []string{
			"low standards", "avoiding feedback", "comfort zone",
			"bureaucracy", "playing it safe",
		},
		KeywordWeights: map[string]float64{
			"performance": 1.0,
			"ownership":   0.9,
			"innovation":  0.8,
			"feedback":    0.7,
		},
		CompensationInfo: CompensationInfo{
			BaseSalaryRange:  []int{180000, 300000},
			EquityPercentage: []float64{0.0, 0.2}, // Netflix prefers cash
			BehavioralImpact: "Netflix cultural fit is critical - poor fit = no offer",
		},
	}
	cm.profiles["netflix"] = profile
}

// loadAmazonProfile loads Amazon's behavioral interview profile
func (cm *CompanyMatcher) loadAmazonProfile() {
	profile := CompanyProfile{
		Name: "Amazon",
		LeadershipPrinciples: []string{
			"customer obsession", "ownership", "invent and simplify",
			"are right, a lot", "learn and be curious", "hire and develop the best",
			"insist on the highest standards", "think big", "bias for action",
			"frugality", "earn trust", "dive deep", "have backbone; disagree and commit",
			"deliver results",
		},
		InterviewFormat: InterviewFormat{
			Duration:           60 * time.Minute,
			NumberOfQuestions:  6,
			TimePerQuestion:    8 * time.Minute,
			FollowUpIntensity:  "High",
			TechnicalMix:       false,
			PanelStyle:         false,
			RequiredStoryTypes: []StoryCategory{CustomerObsession, Ownership, TechnicalLeadership, ProblemSolving},
		},
		KeywordWeights: map[string]float64{
			"customer":  1.0,
			"ownership": 1.0,
			"results":   0.9,
			"standards": 0.8,
		},
		CompensationInfo: CompensationInfo{
			BaseSalaryRange:  []int{130000, 180000}, // Amazon caps base salary
			EquityPercentage: []float64{0.5, 1.5},   // Heavy on equity
			BehavioralImpact: "Leadership principles are 50% of evaluation",
		},
	}
	cm.profiles["amazon"] = profile
}

// loadAppleProfile loads Apple's behavioral interview profile
func (cm *CompanyMatcher) loadAppleProfile() {
	profile := CompanyProfile{
		Name: "Apple",
		CulturalValues: []string{
			"innovation", "quality", "simplicity", "user experience",
			"attention to detail", "excellence", "privacy", "accessibility",
		},
		InterviewFormat: InterviewFormat{
			Duration:           45 * time.Minute,
			NumberOfQuestions:  4,
			TimePerQuestion:    10 * time.Minute,
			FollowUpIntensity:  "Medium",
			TechnicalMix:       true,
			PanelStyle:         false,
			RequiredStoryTypes: []StoryCategory{Innovation, TechnicalLeadership, ProblemSolving},
		},
		RedFlags: []string{
			"rushed products", "compromising quality", "ignoring user experience",
			"over-engineering", "privacy concerns",
		},
		KeywordWeights: map[string]float64{
			"quality":    1.0,
			"user":       0.9,
			"innovation": 0.8,
			"detail":     0.7,
		},
		CompensationInfo: CompensationInfo{
			BaseSalaryRange:  []int{150000, 240000},
			EquityPercentage: []float64{0.1, 0.4},
			BehavioralImpact: "Apple values cultural fit highly for team cohesion",
		},
	}
	cm.profiles["apple"] = profile
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
