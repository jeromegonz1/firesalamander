package recommender

import (
	"context"
	"testing"

	"firesalamander/internal/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that SemanticRecommender implements the Agent interface
func TestSemanticRecommenderImplementsAgent(t *testing.T) {
	recommender := NewSemanticRecommender()
	
	// Verify that SemanticRecommender implements Agent interface
	var agent agents.Agent = recommender
	assert.NotNil(t, agent)
}

// Test Agent.Name() implementation
func TestSemanticRecommenderName(t *testing.T) {
	recommender := NewSemanticRecommender()
	
	name := recommender.Name()
	assert.Equal(t, "semantic-recommender", name)
	assert.NotEmpty(t, name)
}

// Test Agent.HealthCheck() implementation
func TestSemanticRecommenderHealthCheck(t *testing.T) {
	recommender := NewSemanticRecommender()
	
	err := recommender.HealthCheck()
	assert.NoError(t, err)
}

// Test Agent.Process() implementation with valid input
func TestSemanticRecommenderProcess(t *testing.T) {
	recommender := NewSemanticRecommender()

	// Create test input with content analysis data
	input := RecommendationRequest{
		Content: ContentAnalysis{
			URL:     "https://example.com/test",
			Title:   "Test Article About AI",
			Content: "This article discusses artificial intelligence, machine learning algorithms, and their applications in modern technology.",
			Keywords: []string{"artificial", "intelligence", "machine", "learning", "algorithms", "technology"},
			Topics:   []string{"AI", "Technology"},
		},
		Context: AnalysisContext{
			Domain:       "technology",
			ContentType:  "blog",
			TargetGoals:  []string{"seo", "engagement"},
			Competitors:  []string{"competitor1.com", "competitor2.com"},
		},
		Options: RecommendationOptions{
			MaxRecommendations: 5,
			Focus:             "content",
			Priority:          "high",
		},
	}

	ctx := context.Background()
	result, err := recommender.Process(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "semantic-recommender", result.AgentName)
	assert.Equal(t, "success", result.Status)
	assert.Greater(t, result.Duration, int64(0))

	// Verify the recommendation data is in the result
	assert.Contains(t, result.Data, "recommendations")
	assert.Contains(t, result.Data, "semantic_score")
	assert.Contains(t, result.Data, "metadata")

	// Verify recommendation data structure
	recommendations, ok := result.Data["recommendations"].([]Recommendation)
	assert.True(t, ok, "recommendations should be of type []Recommendation")
	assert.Greater(t, len(recommendations), 0, "should have at least one recommendation")

	// Verify each recommendation has required fields
	for _, rec := range recommendations {
		assert.NotEmpty(t, rec.ID, "recommendation should have an ID")
		assert.NotEmpty(t, rec.Title, "recommendation should have a title")
		assert.NotEmpty(t, rec.Description, "recommendation should have a description")
		assert.NotEmpty(t, rec.Category, "recommendation should have a category")
		assert.Greater(t, rec.Impact, 0.0, "recommendation should have positive impact")
		assert.Greater(t, rec.Confidence, 0.0, "recommendation should have positive confidence")
	}
}

// Test Agent.Process() with invalid input
func TestSemanticRecommenderProcessInvalidInput(t *testing.T) {
	recommender := NewSemanticRecommender()
	ctx := context.Background()

	// Test with nil input
	result, err := recommender.Process(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test with wrong input type
	result, err = recommender.Process(ctx, "invalid input")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test Agent.Process() with minimal content
func TestSemanticRecommenderProcessMinimalContent(t *testing.T) {
	recommender := NewSemanticRecommender()

	input := RecommendationRequest{
		Content: ContentAnalysis{
			URL:     "https://example.com/minimal",
			Title:   "Short Title",
			Content: "Short content",
			Keywords: []string{},
			Topics:   []string{},
		},
		Context: AnalysisContext{
			Domain:      "general",
			ContentType: "page",
		},
	}

	ctx := context.Background()
	result, err := recommender.Process(ctx, input)

	require.NoError(t, err) // Should handle minimal content gracefully
	assert.NotNil(t, result)
	assert.Equal(t, "semantic-recommender", result.AgentName)
	
	// Should return recommendations even for minimal content
	recommendations, ok := result.Data["recommendations"].([]Recommendation)
	assert.True(t, ok)
	assert.Greater(t, len(recommendations), 0, "should provide basic recommendations")
}

// Test content enhancement recommendations
func TestContentEnhancementRecommendations(t *testing.T) {
	recommender := NewSemanticRecommender()

	content := ContentAnalysis{
		URL:     "https://example.com/article",
		Title:   "Short",  // Too short
		Content: "Brief.", // Too brief
		Keywords: []string{"test"},
		Topics:   []string{"testing"},
	}

	recommendations := recommender.generateContentRecommendations(content)

	assert.Greater(t, len(recommendations), 0, "should generate content recommendations")
	
	// Should recommend title expansion
	hasTitleRec := false
	hasContentRec := false
	
	for _, rec := range recommendations {
		if rec.Category == "content" {
			if rec.Type == "title_enhancement" {
				hasTitleRec = true
			}
			if rec.Type == "content_expansion" {
				hasContentRec = true
			}
		}
	}
	
	assert.True(t, hasTitleRec, "should recommend title enhancement")
	assert.True(t, hasContentRec, "should recommend content expansion")
}

// Test SEO recommendations
func TestSEORecommendations(t *testing.T) {
	recommender := NewSemanticRecommender()

	content := ContentAnalysis{
		URL:     "https://example.com/article",
		Title:   "Article About Technology Without Keywords", 
		Content: "This is content that doesn't use the main keywords effectively.",
		Keywords: []string{"artificial", "intelligence", "machine", "learning"},
		Topics:   []string{"AI"},
	}

	recommendations := recommender.generateSEORecommendations(content)

	assert.Greater(t, len(recommendations), 0, "should generate SEO recommendations")
	
	// Should recommend keyword expansion (since we only have 4 keywords, less than 5)
	hasKeywordRec := false
	
	for _, rec := range recommendations {
		if rec.Category == "seo" && rec.Type == "keyword_expansion" {
			hasKeywordRec = true
			break
		}
	}
	
	assert.True(t, hasKeywordRec, "should recommend keyword expansion")
}

// Test semantic scoring
func TestSemanticScoring(t *testing.T) {
	recommender := NewSemanticRecommender()

	content := ContentAnalysis{
		URL:     "https://example.com/good-article",
		Title:   "Comprehensive Guide to Machine Learning and Artificial Intelligence",
		Content: "This comprehensive guide covers machine learning algorithms, artificial intelligence applications, and deep learning techniques. We explore neural networks, supervised learning, unsupervised learning, and reinforcement learning.",
		Keywords: []string{"machine", "learning", "artificial", "intelligence", "algorithms"},
		Topics:   []string{"AI", "Technology"},
	}

	scoreResult := recommender.calculateSemanticScore(content)

	assert.Greater(t, scoreResult.OverallScore, 0.0, "semantic score should be positive")
	assert.LessOrEqual(t, scoreResult.OverallScore, 10.0, "semantic score should not exceed maximum")
	assert.NotEmpty(t, scoreResult.CategoryScores, "should have category scores")
	assert.NotEmpty(t, scoreResult.KeyInsights, "should have key insights")
}

// Test recommendation prioritization
func TestRecommendationPrioritization(t *testing.T) {
	recommender := NewSemanticRecommender()

	recommendations := []Recommendation{
		{
			ID:         "rec1",
			Title:      "Low Impact Rec",
			Category:   "content",
			Impact:     2.0,
			Confidence: 0.5,
			Priority:   "low",
		},
		{
			ID:         "rec2", 
			Title:      "High Impact Rec",
			Category:   "seo",
			Impact:     8.0,
			Confidence: 0.9,
			Priority:   "high",
		},
		{
			ID:         "rec3",
			Title:      "Medium Impact Rec", 
			Category:   "engagement",
			Impact:     5.0,
			Confidence: 0.7,
			Priority:   "medium",
		},
	}

	prioritized := recommender.prioritizeRecommendations(recommendations, 2)

	assert.Len(t, prioritized, 2, "should return requested number of recommendations")
	assert.Equal(t, "rec2", prioritized[0].ID, "should prioritize high impact recommendation first")
	assert.Greater(t, prioritized[0].Impact, prioritized[1].Impact, "first recommendation should have higher impact")
}

// Benchmark SemanticRecommender.Process() performance
func BenchmarkSemanticRecommenderProcess(b *testing.B) {
	recommender := NewSemanticRecommender()
	
	input := RecommendationRequest{
		Content: ContentAnalysis{
			URL:     "https://example.com/benchmark",
			Title:   "Benchmark Article About Technology and AI",
			Content: "This article covers various technology topics including artificial intelligence, machine learning, software development, and digital transformation strategies.",
			Keywords: []string{"technology", "artificial", "intelligence", "machine", "learning", "software"},
			Topics:   []string{"Technology", "AI", "Software"},
		},
		Context: AnalysisContext{
			Domain:      "technology",
			ContentType: "article",
			TargetGoals: []string{"seo", "engagement"},
		},
		Options: RecommendationOptions{
			MaxRecommendations: 10,
			Focus:             "comprehensive",
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := recommender.Process(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}