// Package main demonstrates the complete FAANG L3/L4 behavioral interview preparation system
// showcasing company-specific optimization, story selection, and real-time interview simulation.
//
// This comprehensive demonstration includes:
// - Story bank with 7 quantified leadership examples
// - Company-specific preparation for Google, Meta, Netflix, Amazon, Apple
// - Real-time interview simulation with AI-powered evaluation
// - Performance analytics and improvement tracking
// - Compensation strategy and negotiation guidance
//
// Run this demo to experience the complete behavioral interview mastery system.
package main

import (
	"fmt"
	"strings"

	"behavioral/star-method-mastery"
)

func main() {
	fmt.Println("🚀 FAANG L3/L4 Behavioral Interview Mastery System")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Initialize the complete system
	storyBank := behavioral.NewStoryBank("Google", "L4")
	companyMatcher := behavioral.NewCompanyMatcher(storyBank)
	simulator := behavioral.NewInterviewSimulator(storyBank)

	// Phase 1: Demonstrate Story Bank Capabilities
	fmt.Println("📖 Phase 1: Leadership Story Bank")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("📊 Total Stories Available: %d\n", len(storyBank.GetAllStories()))

	// Show story distribution by category
	categories := []behavioral.StoryCategory{
		behavioral.TechnicalLeadership,
		behavioral.ProblemSolving,
		behavioral.Innovation,
		behavioral.Collaboration,
		behavioral.Failure,
	}
	for _, category := range categories {
		stories := storyBank.GetStoriesForCategory(category)
		fmt.Printf("   %s: %d stories\n", category.String(), len(stories))
	}

	// Highlight a key story
	allStories := storyBank.GetAllStories()
	var bestStory behavioral.STARStory
	if len(allStories) > 0 {
		// Find first technical leadership story
		for _, story := range allStories {
			if story.Category == behavioral.TechnicalLeadership {
				bestStory = story
				break
			}
		}
	}
	if bestStory.Title != "" {
		fmt.Printf("\n🏆 Featured Story: %s\n", bestStory.Title)
		fmt.Printf("   Impact: %d team members\n", bestStory.TeamSize)
		if len(bestStory.Metrics) > 0 {
			fmt.Printf("   Key Metric: %s improved by %.1f%s\n",
				bestStory.Metrics[0].Description,
				bestStory.Metrics[0].Value,
				bestStory.Metrics[0].Unit)
		}
	}
	fmt.Println()

	// Phase 2: Company-Specific Optimization Demo
	fmt.Println("🎯 Phase 2: Company-Specific Preparation")
	fmt.Println(strings.Repeat("-", 40))

	companies := []string{"Google", "Meta", "Netflix", "Amazon", "Apple"}
	for _, company := range companies {
		fmt.Printf("📋 %s L4 Interview Preparation:\n", company)

		// Get optimal stories for this company
		optimalStories, err := companyMatcher.GetOptimalStories(company, "L4", 3)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
			continue
		}

		fmt.Printf("   Recommended Stories: %d selected\n", len(optimalStories))
		for i, story := range optimalStories {
			fmt.Printf("   %d. %s (%s)\n", i+1, story.Title, story.Category.String())
		}

		// Show company-specific preparation
		prep, err := companyMatcher.PrepareForCompany(company, "L4")
		if err == nil {
			fmt.Printf("   Company Alignment Score: %.1f%%\n",
				calculateAverageAlignment(optimalStories, prep.Company))
			fmt.Printf("   Questions to Practice: %d\n", len(prep.QuestionPreparation))
			fmt.Printf("   Expected Comp Range: $%d - $%d\n",
				prep.CompensationStrategy.ExpectedOffer.BaseSalary-20000,
				prep.CompensationStrategy.ExpectedOffer.BaseSalary+20000)
		}
		fmt.Println()
	}

	// Phase 3: Interactive Interview Simulation
	fmt.Println("🎪 Phase 3: Live Interview Simulation")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("🔴 LIVE: Google L4 Software Engineer Interview")
	fmt.Println()

	// Run a realistic simulation
	session, err := simulator.StartSimulation("Google", "L4", 3)
	if err != nil {
		fmt.Printf("Simulation error: %v\n", err)
		return
	}

	// Display session results
	fmt.Printf("✅ Simulation Complete!\n")
	fmt.Printf("   Session ID: %s\n", session.ID)
	fmt.Printf("   Duration: %v\n", session.Duration)
	fmt.Printf("   Overall Score: %.1f/100\n", session.OverallScore)
	fmt.Printf("   Company Alignment: %.1f%%\n", session.CompanyAlignment)

	if len(session.StrengthAreas) > 0 {
		fmt.Printf("   Strengths: %s\n", session.StrengthAreas[0])
	}
	if len(session.ImprovementAreas) > 0 {
		fmt.Printf("   Focus Area: %s\n", session.ImprovementAreas[0])
	}

	fmt.Println()
	fmt.Println("📊 Detailed Question Performance:")
	for i, response := range session.Questions {
		fmt.Printf("   Q%d: %.1f/100 (STAR: %.1f, Alignment: %.1f)\n",
			i+1, response.OverallScore, response.STARScore.Overall, response.AlignmentScore)
		if len(response.Feedback) > 0 {
			fmt.Printf("       💡 %s\n", response.Feedback[0])
		}
	}
	fmt.Println()

	// Phase 4: Performance Analytics
	fmt.Println("📈 Phase 4: Performance Analytics & Improvement")
	fmt.Println(strings.Repeat("-", 40))

	// Generate improvement report (simulated for demo purposes)
	fmt.Printf("📊 Ready to start your interview preparation journey!\n")
	fmt.Printf("   Complete practice sessions to see detailed analytics\n")
	fmt.Printf("   Track your progress across STAR methodology components\n")
	fmt.Printf("   Get personalized recommendations for improvement\n")
	fmt.Println()

	// Phase 5: System Capabilities Summary
	fmt.Println("⚡ Phase 5: System Capabilities Summary")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("🎯 Story Bank: %d quantified leadership stories\n", len(storyBank.GetAllStories()))
	fmt.Printf("🏢 Companies: %d FAANG profiles (Google, Meta, Netflix, Amazon, Apple)\n", len(companies))
	fmt.Printf("📝 Question Bank: 15+ company-specific behavioral questions\n")
	fmt.Printf("🔍 AI Evaluation: STAR methodology + company alignment scoring\n")
	fmt.Printf("📊 Analytics: Performance tracking + improvement recommendations\n")
	fmt.Printf("💰 Compensation: Level-specific salary bands + negotiation strategies\n")
	fmt.Println()

	// Ready for interviews message
	fmt.Println("🎉 CONGRATULATIONS!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("You now have access to the most comprehensive FAANG L3/L4")
	fmt.Println("behavioral interview preparation system available.")
	fmt.Println()
	fmt.Println("📋 Your interview readiness includes:")
	fmt.Println("   ✅ 7 quantified leadership stories with business impact")
	fmt.Println("   ✅ Company-specific story optimization for all FAANG")
	fmt.Println("   ✅ Real-time interview simulation with AI feedback")
	fmt.Println("   ✅ Performance analytics and improvement tracking")
	fmt.Println("   ✅ Compensation negotiation strategies")
	fmt.Println()
	fmt.Println("🚀 You're ready to ace your FAANG behavioral interviews!")
	fmt.Println("   Good luck, and remember: preparation meets opportunity!")
}

// Helper function to calculate average alignment score
func calculateAverageAlignment(stories []behavioral.STARStory, profile behavioral.CompanyProfile) float64 {
	if len(stories) == 0 {
		return 0.0
	}

	totalAlignment := 0.0
	for _, story := range stories {
		// Simple alignment calculation based on keyword matching
		storyText := fmt.Sprintf("%s %s %s %s", story.Situation, story.Task, story.Action, story.Result)

		alignmentScore := 0.0
		totalWeight := 0.0

		// Check against company values
		for _, value := range profile.CulturalValues {
			totalWeight += 1.0
			if strings.Contains(strings.ToLower(storyText), strings.ToLower(value)) {
				alignmentScore += 1.0
			}
		}

		// Check against keyword weights
		for keyword, weight := range profile.KeywordWeights {
			totalWeight += weight
			if strings.Contains(strings.ToLower(storyText), strings.ToLower(keyword)) {
				alignmentScore += weight
			}
		}

		if totalWeight > 0 {
			totalAlignment += (alignmentScore / totalWeight) * 100.0
		}
	}

	return totalAlignment / float64(len(stories))
}
