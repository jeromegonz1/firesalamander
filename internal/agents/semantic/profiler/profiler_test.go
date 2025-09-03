package profiler

import (
	"context"
	"testing"

	"firesalamander/internal/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that PageProfiler implements the Agent interface
func TestPageProfilerImplementsAgent(t *testing.T) {
	profiler := NewPageProfiler()
	
	// Verify that PageProfiler implements Agent interface
	var agent agents.Agent = profiler
	assert.NotNil(t, agent)
}

// Test Agent.Name() implementation
func TestPageProfilerName(t *testing.T) {
	profiler := NewPageProfiler()
	
	name := profiler.Name()
	assert.Equal(t, "page-profiler", name)
	assert.NotEmpty(t, name)
}

// Test Agent.HealthCheck() implementation
func TestPageProfilerHealthCheck(t *testing.T) {
	profiler := NewPageProfiler()
	
	err := profiler.HealthCheck()
	assert.NoError(t, err)
}

// Test Agent.Process() implementation with valid input
func TestPageProfilerProcess(t *testing.T) {
	profiler := NewPageProfiler()

	// Create test input with page data
	input := ProfileRequest{
		Page: agents.PageInfo{
			URL:     "https://example.com/tech-article",
			Lang:    "en",
			Title:   "Understanding Machine Learning: A Comprehensive Guide",
			H1:      "Machine Learning Fundamentals",
			H2:      []string{"Introduction to ML", "Types of Learning", "Common Algorithms"},
			H3:      []string{"Supervised Learning", "Unsupervised Learning", "Neural Networks", "Decision Trees"},
			Content: "Machine learning is a subset of artificial intelligence that enables computers to learn and improve from experience. This comprehensive guide covers supervised learning, unsupervised learning, neural networks, and decision trees. We'll explore algorithms, data preprocessing, model evaluation, and real-world applications in technology, healthcare, and business.",
			Canonical: "https://example.com/tech-article",
			MetaIndex: true,
		},
		Options: ProfileOptions{
			IncludeSemanticAnalysis: true,
			IncludeContentQuality:   true,
			IncludeTechnicalSEO:     true,
			DetailLevel:            "comprehensive",
		},
	}

	ctx := context.Background()
	result, err := profiler.Process(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "page-profiler", result.AgentName)
	assert.Equal(t, "success", result.Status)
	assert.Greater(t, result.Duration, int64(0))

	// Verify the profile data is in the result
	assert.Contains(t, result.Data, "profile")
	assert.Contains(t, result.Data, "metadata")

	// Verify profile data structure
	profile, ok := result.Data["profile"].(PageProfile)
	assert.True(t, ok, "profile should be of type PageProfile")
	
	// Verify profile components
	assert.NotEmpty(t, profile.URL, "profile should have URL")
	assert.NotEmpty(t, profile.Title, "profile should have title")
	assert.Greater(t, len(profile.SemanticProfile.Keywords), 0, "should have semantic keywords")
	assert.Greater(t, len(profile.SemanticProfile.Topics), 0, "should have topics")
	assert.Greater(t, profile.ContentProfile.WordCount, 0, "should have word count")
	assert.Greater(t, profile.QualityScore.Overall, 0.0, "should have overall quality score")
}

// Test Agent.Process() with invalid input
func TestPageProfilerProcessInvalidInput(t *testing.T) {
	profiler := NewPageProfiler()
	ctx := context.Background()

	// Test with nil input
	result, err := profiler.Process(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test with wrong input type
	result, err = profiler.Process(ctx, "invalid input")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test Agent.Process() with minimal page data
func TestPageProfilerProcessMinimalPage(t *testing.T) {
	profiler := NewPageProfiler()

	input := ProfileRequest{
		Page: agents.PageInfo{
			URL:     "https://example.com/minimal",
			Title:   "Short",
			Content: "Brief content.",
		},
		Options: ProfileOptions{
			DetailLevel: "basic",
		},
	}

	ctx := context.Background()
	result, err := profiler.Process(ctx, input)

	require.NoError(t, err) // Should handle minimal data gracefully
	assert.NotNil(t, result)
	assert.Equal(t, "page-profiler", result.AgentName)
	
	// Should return basic profile
	profile, ok := result.Data["profile"].(PageProfile)
	assert.True(t, ok)
	assert.Equal(t, "https://example.com/minimal", profile.URL)
}

// Test semantic analysis functionality
func TestSemanticAnalysis(t *testing.T) {
	profiler := NewPageProfiler()

	page := agents.PageInfo{
		Title:   "Advanced Machine Learning Techniques and Applications",
		Content: "Machine learning algorithms including neural networks, deep learning, natural language processing, computer vision, and reinforcement learning are transforming industries. These AI technologies enable pattern recognition, predictive analytics, and automated decision-making systems.",
	}

	semantic := profiler.analyzeSemantics(page)

	assert.Greater(t, len(semantic.Keywords), 0, "should extract keywords")
	assert.Greater(t, len(semantic.Topics), 0, "should identify topics")
	assert.Greater(t, len(semantic.Entities), 0, "should identify entities")
	assert.Greater(t, semantic.Complexity, 0.0, "should calculate complexity")

	// Check for expected ML-related keywords
	keywordSet := make(map[string]bool)
	for _, kw := range semantic.Keywords {
		keywordSet[kw.Term] = true
	}
	
	expectedKeywords := []string{"machine", "learning", "neural", "networks"}
	foundCount := 0
	for _, expected := range expectedKeywords {
		if keywordSet[expected] {
			foundCount++
		}
	}
	assert.Greater(t, foundCount, 1, "should find ML-related keywords")
}

// Test content analysis functionality
func TestContentAnalysis(t *testing.T) {
	profiler := NewPageProfiler()

	page := agents.PageInfo{
		Title:   "Comprehensive Guide to Web Development",
		H1:      "Web Development Guide",
		H2:      []string{"Frontend Development", "Backend Development", "Full Stack Development"},
		H3:      []string{"HTML & CSS", "JavaScript", "React", "Node.js", "Databases"},
		Content: "Web development encompasses frontend and backend technologies. Frontend development focuses on user interfaces using HTML, CSS, and JavaScript frameworks like React. Backend development involves server-side programming, databases, and APIs. Full stack developers work with both frontend and backend technologies to create complete web applications.",
	}

	content := profiler.analyzeContent(page)

	assert.Greater(t, content.WordCount, 0, "should count words")
	assert.Greater(t, content.ReadingTime, 0, "should estimate reading time")
	assert.Greater(t, len(content.HeadingStructure), 0, "should analyze heading structure")
	assert.Greater(t, content.ReadabilityScore, 0.0, "should calculate readability")

	// Check heading structure
	assert.Equal(t, "Web Development Guide", content.HeadingStructure[0].Text)
	assert.Equal(t, 1, content.HeadingStructure[0].Level)
}

// Test quality assessment functionality
func TestQualityAssessment(t *testing.T) {
	profiler := NewPageProfiler()

	profile := PageProfile{
		URL:   "https://example.com/quality-test",
		Title: "High Quality Article About Technology Trends",
		ContentProfile: ContentProfile{
			WordCount:    800,
			ReadingTime:  4,
			ReadabilityScore: 7.5,
			HeadingStructure: []HeadingInfo{
				{Level: 1, Text: "Technology Trends"},
				{Level: 2, Text: "Artificial Intelligence"},
				{Level: 2, Text: "Cloud Computing"},
			},
		},
		SemanticProfile: SemanticProfile{
			Keywords: []SemanticKeyword{
				{Term: "technology", Weight: 0.1, Frequency: 5},
				{Term: "artificial", Weight: 0.08, Frequency: 4},
				{Term: "intelligence", Weight: 0.08, Frequency: 4},
			},
			Topics: []TopicInfo{
				{Name: "Technology", Relevance: 0.9},
				{Name: "AI", Relevance: 0.8},
			},
		},
	}

	quality := profiler.assessQuality(profile)

	assert.Greater(t, quality.Overall, 0.0, "should have overall quality score")
	assert.Greater(t, quality.Content, 0.0, "should have content quality score")
	assert.Greater(t, quality.SEO, 0.0, "should have SEO quality score")
	assert.Greater(t, quality.Structure, 0.0, "should have structure quality score")
	assert.NotEmpty(t, quality.Strengths, "should identify strengths")
}

// Test technical SEO analysis
func TestTechnicalSEOAnalysis(t *testing.T) {
	profiler := NewPageProfiler()

	page := agents.PageInfo{
		URL:       "https://example.com/seo-test",
		Title:     "SEO Optimized Article About Digital Marketing",
		Canonical: "https://example.com/seo-test",
		MetaIndex: true,
	}

	seo := profiler.analyzeTechnicalSEO(page)

	assert.NotEmpty(t, seo.URLStructure.Analysis, "should analyze URL structure")
	assert.Greater(t, seo.URLStructure.Score, 0.0, "should score URL structure")
	assert.Greater(t, seo.TitleTag.Score, 0.0, "should score title tag")
	assert.NotEmpty(t, seo.IndexabilityStatus, "should determine indexability")
	assert.True(t, seo.IsIndexable, "should be indexable")
}

// Test entity extraction
func TestEntityExtraction(t *testing.T) {
	profiler := NewPageProfiler()

	text := "Apple Inc. released the iPhone 14 in September 2022. CEO Tim Cook presented the new features at the Cupertino headquarters in California."

	entities := profiler.extractEntities(text)

	assert.Greater(t, len(entities), 0, "should extract entities")

	// Check for expected entity types
	entityTypes := make(map[string]bool)
	for _, entity := range entities {
		entityTypes[entity.Type] = true
	}

	// Should identify some entities (exact types may vary)
	assert.True(t, len(entityTypes) > 0, "should identify different entity types")
}

// Test topic classification
func TestTopicClassification(t *testing.T) {
	profiler := NewPageProfiler()

	keywords := []SemanticKeyword{
		{Term: "machine", Weight: 0.1},
		{Term: "learning", Weight: 0.09},
		{Term: "artificial", Weight: 0.08},
		{Term: "intelligence", Weight: 0.08},
		{Term: "neural", Weight: 0.07},
		{Term: "networks", Weight: 0.06},
	}

	content := "Artificial intelligence and machine learning technologies"

	topics := profiler.classifyTopics(keywords, content)

	assert.Greater(t, len(topics), 0, "should classify topics")

	// Should identify AI/ML related topics
	topicNames := make(map[string]bool)
	for _, topic := range topics {
		topicNames[topic.Name] = true
	}

	assert.True(t, len(topicNames) > 0, "should identify relevant topics")
}

// Benchmark PageProfiler.Process() performance
func BenchmarkPageProfilerProcess(b *testing.B) {
	profiler := NewPageProfiler()
	
	input := ProfileRequest{
		Page: agents.PageInfo{
			URL:     "https://example.com/benchmark",
			Title:   "Comprehensive Technology Article for Benchmarking",
			H1:      "Technology Benchmarking Guide",
			H2:      []string{"Performance Metrics", "Analysis Methods", "Best Practices"},
			Content: "This comprehensive guide covers technology benchmarking methodologies, performance metrics analysis, and industry best practices. We explore various measurement techniques, data collection methods, statistical analysis approaches, and reporting frameworks used in technology performance evaluation.",
		},
		Options: ProfileOptions{
			IncludeSemanticAnalysis: true,
			IncludeContentQuality:   true,
			IncludeTechnicalSEO:     true,
			DetailLevel:            "comprehensive",
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := profiler.Process(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}