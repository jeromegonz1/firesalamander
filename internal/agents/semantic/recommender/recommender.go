package recommender

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"firesalamander/internal/agents"
)

// Ensure SemanticRecommender implements the Agent interface
var _ agents.Agent = (*SemanticRecommender)(nil)

// NewSemanticRecommender creates a new SemanticRecommender instance
func NewSemanticRecommender() *SemanticRecommender {
	return &SemanticRecommender{
		maxRecommendations:  10,
		defaultFocus:       "comprehensive",
		confidenceThreshold: 0.5,
	}
}

// Name returns the agent name
func (sr *SemanticRecommender) Name() string {
	return "semantic-recommender"
}

// Process implements the Agent interface for the SemanticRecommender
func (sr *SemanticRecommender) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()

	// Validate and parse input data
	request, err := sr.parseInput(data)
	if err != nil {
		return nil, fmt.Errorf("invalid input data: %w", err)
	}

	// Generate recommendations
	recommendations := sr.generateRecommendations(*request)
	
	// Calculate semantic score
	semanticScore := sr.calculateSemanticScore(request.Content)
	
	// Create metadata
	metadata := RecommendationMetadata{
		ProcessingTime:       time.Since(startTime).Milliseconds(),
		RecommendationsCount: len(recommendations),
		Algorithm:           "semantic-analysis-v1",
		Confidence:          sr.calculateOverallConfidence(recommendations),
		ContentQuality:      sr.assessContentQuality(request.Content),
	}

	// Count high priority recommendations
	for _, rec := range recommendations {
		if rec.Priority == "high" {
			metadata.HighPriorityCount++
		}
	}

	// Create agent result
	agentResult := &agents.AgentResult{
		AgentName: sr.Name(),
		Status:    "success",
		Data: map[string]interface{}{
			"recommendations": recommendations,
			"semantic_score":  semanticScore,
			"metadata":        metadata,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}

	return agentResult, nil
}

// HealthCheck implements the Agent interface
func (sr *SemanticRecommender) HealthCheck() error {
	if sr.maxRecommendations < 1 {
		return fmt.Errorf("max recommendations must be >= 1")
	}
	
	if sr.confidenceThreshold < 0 || sr.confidenceThreshold > 1 {
		return fmt.Errorf("confidence threshold must be between 0 and 1")
	}
	
	return nil
}

// parseInput validates and parses the input data into a RecommendationRequest
func (sr *SemanticRecommender) parseInput(data interface{}) (*RecommendationRequest, error) {
	if data == nil {
		return nil, fmt.Errorf("input data is nil")
	}

	request, ok := data.(RecommendationRequest)
	if !ok {
		return nil, fmt.Errorf("expected RecommendationRequest, got %T", data)
	}

	// Apply default values if not specified
	if request.Options.MaxRecommendations <= 0 {
		request.Options.MaxRecommendations = sr.maxRecommendations
	}
	if request.Options.Focus == "" {
		request.Options.Focus = sr.defaultFocus
	}

	// Calculate word count and reading time if not provided
	if request.Content.WordCount == 0 {
		request.Content.WordCount = sr.countWords(request.Content.Content)
	}
	if request.Content.ReadingTime == 0 {
		request.Content.ReadingTime = sr.estimateReadingTime(request.Content.WordCount)
	}

	return &request, nil
}

// generateRecommendations creates comprehensive recommendations based on content analysis
func (sr *SemanticRecommender) generateRecommendations(request RecommendationRequest) []Recommendation {
	var allRecommendations []Recommendation

	// Generate different types of recommendations based on focus
	switch request.Options.Focus {
	case "content":
		allRecommendations = append(allRecommendations, sr.generateContentRecommendations(request.Content)...)
	case "seo":
		allRecommendations = append(allRecommendations, sr.generateSEORecommendations(request.Content)...)
	case "engagement":
		allRecommendations = append(allRecommendations, sr.generateEngagementRecommendations(request.Content)...)
	default: // comprehensive
		allRecommendations = append(allRecommendations, sr.generateContentRecommendations(request.Content)...)
		allRecommendations = append(allRecommendations, sr.generateSEORecommendations(request.Content)...)
		allRecommendations = append(allRecommendations, sr.generateEngagementRecommendations(request.Content)...)
		allRecommendations = append(allRecommendations, sr.generateTechnicalRecommendations(request.Content)...)
	}

	// Filter by confidence threshold
	var qualifiedRecommendations []Recommendation
	for _, rec := range allRecommendations {
		if rec.Confidence >= sr.confidenceThreshold {
			qualifiedRecommendations = append(qualifiedRecommendations, rec)
		}
	}

	// Prioritize and limit recommendations
	return sr.prioritizeRecommendations(qualifiedRecommendations, request.Options.MaxRecommendations)
}

// generateContentRecommendations creates content-focused recommendations
func (sr *SemanticRecommender) generateContentRecommendations(content ContentAnalysis) []Recommendation {
	var recommendations []Recommendation
	recID := 1

	// Title optimization
	if len(content.Title) < 30 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("content_%d", recID),
			Title:       "Expand Title Length",
			Description: "Your title is too short. Expand it to 30-60 characters for better SEO and user engagement.",
			Category:    "content",
			Type:        "title_enhancement",
			Impact:      7.5,
			Confidence:  0.9,
			Priority:    "high",
			Effort:      "low",
			Tags:        []string{"title", "seo", "engagement"},
			Examples:    []string{"Add descriptive keywords", "Include benefits or numbers", "Make it more specific"},
			Implementation: Implementation{
				Steps: []string{
					"Analyze current title keywords",
					"Research competitor titles",
					"Add descriptive modifiers",
					"Include target keywords naturally",
				},
				TimeEstimate: "15-30 minutes",
				Difficulty:   "easy",
			},
		})
		recID++
	}

	// Content length optimization
	if content.WordCount < 300 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("content_%d", recID),
			Title:       "Expand Content Length",
			Description: fmt.Sprintf("Your content has %d words. Expand to at least 300 words for better SEO performance.", content.WordCount),
			Category:    "content",
			Type:        "content_expansion",
			Impact:      8.0,
			Confidence:  0.85,
			Priority:    "high",
			Effort:      "medium",
			Tags:        []string{"content", "seo", "depth"},
			Examples: []string{
				"Add more detailed explanations",
				"Include examples and case studies",
				"Expand on key points with supporting data",
			},
			Implementation: Implementation{
				Steps: []string{
					"Identify thin content sections",
					"Research additional subtopics",
					"Add examples and illustrations",
					"Include data and statistics",
				},
				TimeEstimate: "1-2 hours",
				Difficulty:   "medium",
			},
		})
		recID++
	}

	// Keyword density analysis
	if len(content.Keywords) > 0 {
		keywordDensity := sr.analyzeKeywordDensity(content)
		if keywordDensity < 0.005 { // Less than 0.5% density
			recommendations = append(recommendations, Recommendation{
				ID:          fmt.Sprintf("content_%d", recID),
				Title:       "Improve Keyword Usage",
				Description: "Your target keywords appear infrequently in the content. Increase natural keyword usage for better SEO.",
				Category:    "content",
				Type:        "keyword_integration",
				Impact:      6.5,
				Confidence:  0.8,
				Priority:    "medium",
				Effort:      "low",
				Tags:        []string{"keywords", "seo", "density"},
				Implementation: Implementation{
					Steps: []string{
						"Review current keyword placement",
						"Identify opportunities for natural integration",
						"Rewrite sentences to include keywords",
						"Maintain readability and flow",
					},
					TimeEstimate: "30-60 minutes",
					Difficulty:   "easy",
				},
			})
			recID++
		}
	}

	return recommendations
}

// generateSEORecommendations creates SEO-focused recommendations
func (sr *SemanticRecommender) generateSEORecommendations(content ContentAnalysis) []Recommendation {
	var recommendations []Recommendation
	recID := 1

	// Semantic keyword recommendations
	if len(content.Keywords) < 5 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("seo_%d", recID),
			Title:       "Expand Keyword Strategy",
			Description: "Add more related keywords to improve semantic SEO and capture long-tail traffic.",
			Category:    "seo",
			Type:        "keyword_expansion",
			Impact:      7.0,
			Confidence:  0.8,
			Priority:    "high",
			Effort:      "medium",
			Tags:        []string{"keywords", "semantic", "longtail"},
			Examples: []string{
				"Research related terms using keyword tools",
				"Include synonyms and variations",
				"Add location-based keywords if relevant",
			},
			Implementation: Implementation{
				Steps: []string{
					"Use keyword research tools",
					"Analyze competitor keywords",
					"Identify semantic variations",
					"Naturally integrate new keywords",
				},
				TimeEstimate: "1-2 hours",
				Difficulty:   "medium",
			},
		})
		recID++
	}

	// Title-content keyword alignment
	titleKeywords := sr.extractKeywordsFromText(content.Title)
	contentKeywords := sr.extractKeywordsFromText(content.Content)
	alignment := sr.calculateKeywordAlignment(titleKeywords, contentKeywords)
	
	if alignment < 0.3 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("seo_%d", recID),
			Title:       "Align Title with Content Keywords",
			Description: "Your title and content keywords don't align well. Better alignment improves topical relevance.",
			Category:    "seo",
			Type:        "keyword_alignment",
			Impact:      6.0,
			Confidence:  0.75,
			Priority:    "medium",
			Effort:      "low",
			Tags:        []string{"title", "keywords", "alignment"},
			Implementation: Implementation{
				Steps: []string{
					"Identify main content themes",
					"Update title to reflect content focus",
					"Ensure primary keywords appear in both",
				},
				TimeEstimate: "15-30 minutes",
				Difficulty:   "easy",
			},
		})
		recID++
	}

	// Topic depth analysis
	if len(content.Topics) < 2 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("seo_%d", recID),
			Title:       "Expand Topic Coverage",
			Description: "Cover more related topics to create comprehensive, authoritative content that ranks better.",
			Category:    "seo",
			Type:        "topic_expansion",
			Impact:      7.5,
			Confidence:  0.8,
			Priority:    "high",
			Effort:      "medium",
			Tags:        []string{"topics", "authority", "comprehensive"},
			Implementation: Implementation{
				Steps: []string{
					"Research related subtopics",
					"Add sections for each subtopic",
					"Link topics together logically",
					"Maintain content flow",
				},
				TimeEstimate: "2-3 hours",
				Difficulty:   "medium",
			},
		})
		recID++
	}

	return recommendations
}

// generateEngagementRecommendations creates engagement-focused recommendations
func (sr *SemanticRecommender) generateEngagementRecommendations(content ContentAnalysis) []Recommendation {
	var recommendations []Recommendation
	recID := 1

	// Reading time optimization
	if content.ReadingTime > 15 {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("engagement_%d", recID),
			Title:       "Optimize Reading Time",
			Description: fmt.Sprintf("Your content takes %d minutes to read. Consider breaking it into sections or shorter pieces.", content.ReadingTime),
			Category:    "engagement",
			Type:        "readability",
			Impact:      5.5,
			Confidence:  0.7,
			Priority:    "medium",
			Effort:      "medium",
			Tags:        []string{"readability", "time", "structure"},
			Implementation: Implementation{
				Steps: []string{
					"Add subheadings to break up text",
					"Use bullet points and lists",
					"Consider splitting into multiple pages",
					"Add visual elements",
				},
				TimeEstimate: "1-2 hours",
				Difficulty:   "medium",
			},
		})
		recID++
	}

	// Content structure
	if !sr.hasGoodStructure(content.Content) {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("engagement_%d", recID),
			Title:       "Improve Content Structure",
			Description: "Add clear headings, subheadings, and formatting to make your content more scannable.",
			Category:    "engagement",
			Type:        "structure",
			Impact:      6.0,
			Confidence:  0.8,
			Priority:    "medium",
			Effort:      "low",
			Tags:        []string{"structure", "headings", "formatting"},
			Implementation: Implementation{
				Steps: []string{
					"Add H2 and H3 headings",
					"Use bullet points for lists",
					"Bold important phrases",
					"Add white space between sections",
				},
				TimeEstimate: "30-60 minutes",
				Difficulty:   "easy",
			},
		})
		recID++
	}

	return recommendations
}

// generateTechnicalRecommendations creates technical optimization recommendations
func (sr *SemanticRecommender) generateTechnicalRecommendations(content ContentAnalysis) []Recommendation {
	var recommendations []Recommendation
	recID := 1

	// URL structure
	if sr.hasSuboptimalURL(content.URL) {
		recommendations = append(recommendations, Recommendation{
			ID:          fmt.Sprintf("technical_%d", recID),
			Title:       "Optimize URL Structure",
			Description: "Your URL could be more SEO-friendly with keywords and better structure.",
			Category:    "technical",
			Type:        "url_optimization",
			Impact:      4.0,
			Confidence:  0.6,
			Priority:    "low",
			Effort:      "high",
			Tags:        []string{"url", "structure", "technical"},
			Implementation: Implementation{
				Steps: []string{
					"Include target keywords in URL",
					"Use hyphens instead of underscores",
					"Keep URL length reasonable",
					"Set up proper redirects",
				},
				TimeEstimate: "1-2 hours",
				Difficulty:   "high",
				Prerequisites: []string{"Access to site configuration"},
			},
		})
		recID++
	}

	return recommendations
}

// calculateSemanticScore calculates overall semantic quality score
func (sr *SemanticRecommender) calculateSemanticScore(content ContentAnalysis) SemanticScore {
	var scores = make(map[string]float64)
	
	// Content length score (0-10)
	lengthScore := math.Min(10, float64(content.WordCount)/150)
	scores["content_length"] = lengthScore
	
	// Keyword density score
	keywordScore := math.Min(10, float64(len(content.Keywords)))
	scores["keyword_diversity"] = keywordScore
	
	// Topic coverage score
	topicScore := math.Min(10, float64(len(content.Topics))*2)
	scores["topic_coverage"] = topicScore
	
	// Title quality score
	titleScore := sr.scoreTitleQuality(content.Title)
	scores["title_quality"] = titleScore
	
	// Calculate overall score
	totalScore := 0.0
	for _, score := range scores {
		totalScore += score
	}
	overallScore := totalScore / float64(len(scores))
	
	// Determine strengths and weaknesses
	var strengths, weaknesses, opportunities []string
	
	for category, score := range scores {
		if score >= 7 {
			strengths = append(strengths, fmt.Sprintf("Strong %s", strings.ReplaceAll(category, "_", " ")))
		} else if score <= 4 {
			weaknesses = append(weaknesses, fmt.Sprintf("Improve %s", strings.ReplaceAll(category, "_", " ")))
			opportunities = append(opportunities, fmt.Sprintf("Enhance %s for better performance", strings.ReplaceAll(category, "_", " ")))
		}
	}
	
	// Key insights
	insights := []string{
		fmt.Sprintf("Content length: %d words", content.WordCount),
		fmt.Sprintf("Keyword coverage: %d terms", len(content.Keywords)),
		fmt.Sprintf("Topic breadth: %d areas", len(content.Topics)),
	}
	
	return SemanticScore{
		OverallScore:   overallScore,
		CategoryScores: scores,
		Strengths:      strengths,
		Weaknesses:     weaknesses,
		Opportunities:  opportunities,
		KeyInsights:    insights,
	}
}

// Helper methods

func (sr *SemanticRecommender) countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func (sr *SemanticRecommender) estimateReadingTime(wordCount int) int {
	// Average reading speed: 200-250 words per minute
	return int(math.Ceil(float64(wordCount) / 225.0))
}

func (sr *SemanticRecommender) analyzeKeywordDensity(content ContentAnalysis) float64 {
	if content.WordCount == 0 || len(content.Keywords) == 0 {
		return 0
	}
	
	contentLower := strings.ToLower(content.Content)
	keywordCount := 0
	
	for _, keyword := range content.Keywords {
		keywordCount += strings.Count(contentLower, strings.ToLower(keyword))
	}
	
	return float64(keywordCount) / float64(content.WordCount)
}

func (sr *SemanticRecommender) extractKeywordsFromText(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var keywords []string
	
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:()")
		if len(word) > 3 && !sr.isStopWord(word) {
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

func (sr *SemanticRecommender) calculateKeywordAlignment(titleKeywords, contentKeywords []string) float64 {
	if len(titleKeywords) == 0 {
		return 0
	}
	
	contentSet := make(map[string]bool)
	for _, keyword := range contentKeywords {
		contentSet[keyword] = true
	}
	
	matches := 0
	for _, titleKeyword := range titleKeywords {
		if contentSet[titleKeyword] {
			matches++
		}
	}
	
	return float64(matches) / float64(len(titleKeywords))
}

func (sr *SemanticRecommender) hasGoodStructure(content string) bool {
	// Simple check for headings or structured content
	return strings.Contains(content, "\n\n") || strings.Contains(content, "#") || strings.Contains(content, "•")
}

func (sr *SemanticRecommender) hasSuboptimalURL(url string) bool {
	// Simple checks for URL optimization
	return strings.Contains(url, "?") || strings.Contains(url, "_") || len(url) > 100
}

func (sr *SemanticRecommender) scoreTitleQuality(title string) float64 {
	score := 5.0 // Base score
	
	// Length bonus/penalty
	if len(title) >= 30 && len(title) <= 60 {
		score += 2
	} else if len(title) < 20 {
		score -= 2
	}
	
	// Word count
	wordCount := len(strings.Fields(title))
	if wordCount >= 6 && wordCount <= 12 {
		score += 1
	}
	
	return math.Max(0, math.Min(10, score))
}

func (sr *SemanticRecommender) isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "or": true, "but": true, "in": true,
		"on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "an": true, "as": true, "is": true,
		"was": true, "are": true, "been": true, "be": true, "have": true,
		"has": true, "had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true,
	}
	return stopWords[word]
}

func (sr *SemanticRecommender) prioritizeRecommendations(recommendations []Recommendation, maxCount int) []Recommendation {
	// Sort by impact score descending, then by confidence descending
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Impact == recommendations[j].Impact {
			return recommendations[i].Confidence > recommendations[j].Confidence
		}
		return recommendations[i].Impact > recommendations[j].Impact
	})
	
	// Limit to max count
	if len(recommendations) > maxCount {
		recommendations = recommendations[:maxCount]
	}
	
	return recommendations
}

func (sr *SemanticRecommender) calculateOverallConfidence(recommendations []Recommendation) float64 {
	if len(recommendations) == 0 {
		return 0
	}
	
	totalConfidence := 0.0
	for _, rec := range recommendations {
		totalConfidence += rec.Confidence
	}
	
	return totalConfidence / float64(len(recommendations))
}

func (sr *SemanticRecommender) assessContentQuality(content ContentAnalysis) string {
	score := 0
	
	// Word count assessment
	if content.WordCount >= 1000 {
		score += 3
	} else if content.WordCount >= 500 {
		score += 2
	} else if content.WordCount >= 300 {
		score += 1
	}
	
	// Keyword diversity
	if len(content.Keywords) >= 10 {
		score += 2
	} else if len(content.Keywords) >= 5 {
		score += 1
	}
	
	// Topic coverage
	if len(content.Topics) >= 3 {
		score += 2
	} else if len(content.Topics) >= 2 {
		score += 1
	}
	
	// Title quality
	if len(content.Title) >= 30 && len(content.Title) <= 60 {
		score += 1
	}
	
	switch {
	case score >= 7:
		return "excellent"
	case score >= 5:
		return "good"
	case score >= 3:
		return "fair"
	default:
		return "poor"
	}
}