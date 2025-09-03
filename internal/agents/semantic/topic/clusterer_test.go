package topic

import (
	"context"
	"testing"

	"firesalamander/internal/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that TopicClusterer implements the Agent interface
func TestTopicClustererImplementsAgent(t *testing.T) {
	clusterer := NewTopicClusterer()
	
	// Verify that TopicClusterer implements Agent interface
	var agent agents.Agent = clusterer
	assert.NotNil(t, agent)
}

// Test Agent.Name() implementation
func TestTopicClustererName(t *testing.T) {
	clusterer := NewTopicClusterer()
	
	name := clusterer.Name()
	assert.Equal(t, "topic-clusterer", name)
	assert.NotEmpty(t, name)
}

// Test Agent.HealthCheck() implementation
func TestTopicClustererHealthCheck(t *testing.T) {
	clusterer := NewTopicClusterer()
	
	err := clusterer.HealthCheck()
	assert.NoError(t, err)
}

// Test Agent.Process() implementation with valid input
func TestTopicClustererProcess(t *testing.T) {
	clusterer := NewTopicClusterer()

	// Create test input with crawl data
	input := ClusterRequest{
		Pages: []agents.PageInfo{
			{
				URL:     "https://example.com/tech",
				Title:   "Technology Blog Post",
				Content: "This article discusses machine learning, artificial intelligence, and deep learning technologies in modern software development.",
			},
			{
				URL:     "https://example.com/cooking",
				Title:   "Recipe for Pasta",
				Content: "How to make delicious pasta with tomatoes, garlic, and olive oil. Cooking tips for beginners.",
			},
			{
				URL:     "https://example.com/ai",
				Title:   "AI Revolution",
				Content: "Artificial intelligence and machine learning are transforming industries with neural networks and automation.",
			},
		},
		NumClusters: 2,
	}

	ctx := context.Background()
	result, err := clusterer.Process(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "topic-clusterer", result.AgentName)
	assert.Equal(t, "success", result.Status)
	assert.Greater(t, result.Duration, int64(0))

	// Verify the clustering data is in the result
	assert.Contains(t, result.Data, "clusters")
	assert.Contains(t, result.Data, "topics")
	assert.Contains(t, result.Data, "metadata")

	// Verify cluster data structure
	clusters, ok := result.Data["clusters"].([]Cluster)
	assert.True(t, ok, "clusters should be of type []Cluster")
	assert.Greater(t, len(clusters), 0, "should have at least one cluster")

	// Verify each cluster has required fields
	for _, cluster := range clusters {
		assert.NotEmpty(t, cluster.ID, "cluster should have an ID")
		assert.NotEmpty(t, cluster.Label, "cluster should have a label")
		assert.Greater(t, len(cluster.Pages), 0, "cluster should have pages")
		assert.Greater(t, len(cluster.Keywords), 0, "cluster should have keywords")
	}
}

// Test Agent.Process() with invalid input
func TestTopicClustererProcessInvalidInput(t *testing.T) {
	clusterer := NewTopicClusterer()
	ctx := context.Background()

	// Test with nil input
	result, err := clusterer.Process(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test with wrong input type
	result, err = clusterer.Process(ctx, "invalid input")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test Agent.Process() with empty pages
func TestTopicClustererProcessEmptyPages(t *testing.T) {
	clusterer := NewTopicClusterer()

	input := ClusterRequest{
		Pages:       []agents.PageInfo{},
		NumClusters: 2,
	}

	ctx := context.Background()
	result, err := clusterer.Process(ctx, input)

	require.NoError(t, err) // Should handle empty gracefully
	assert.NotNil(t, result)
	assert.Equal(t, "topic-clusterer", result.AgentName)
	
	// Should return empty clusters
	clusters, ok := result.Data["clusters"].([]Cluster)
	assert.True(t, ok)
	assert.Equal(t, 0, len(clusters))
}

// Test semantic clustering functionality
func TestSemanticClustering(t *testing.T) {
	clusterer := NewTopicClusterer()

	// Create test pages with distinct topics
	pages := []agents.PageInfo{
		{
			URL:     "https://example.com/sports1",
			Title:   "Football Championship",
			Content: "Football soccer championship league tournament players team sport match game",
		},
		{
			URL:     "https://example.com/sports2", 
			Title:   "Basketball News",
			Content: "Basketball players team sport game tournament league championship match",
		},
		{
			URL:     "https://example.com/tech1",
			Title:   "Programming Guide",
			Content: "Programming software development coding computer technology algorithm function",
		},
		{
			URL:     "https://example.com/tech2",
			Title:   "Software Engineering",
			Content: "Software engineering development programming computer technology coding algorithm",
		},
	}

	clusters := clusterer.clusterPages(pages, 2)

	assert.Len(t, clusters, 2, "should create exactly 2 clusters")

	// Verify clusters have distinct topics
	clusterTopics := make(map[string]bool)
	for _, cluster := range clusters {
		assert.NotEmpty(t, cluster.Label)
		assert.Greater(t, len(cluster.Pages), 0)
		assert.Greater(t, len(cluster.Keywords), 0)
		clusterTopics[cluster.Label] = true
	}

	// Should have different cluster labels
	assert.Len(t, clusterTopics, 2, "clusters should have distinct labels")
}

// Test keyword extraction for clusters
func TestClusterKeywordExtraction(t *testing.T) {
	clusterer := NewTopicClusterer()

	pages := []agents.PageInfo{
		{
			Content: "machine learning artificial intelligence neural networks deep learning",
		},
		{
			Content: "machine learning algorithms data science artificial intelligence",
		},
	}

	keywords := clusterer.extractClusterKeywords(pages)

	assert.Greater(t, len(keywords), 0, "should extract keywords")
	
	// Should include common technical terms
	keywordSet := make(map[string]bool)
	for _, kw := range keywords {
		keywordSet[kw.Term] = true
	}

	// Check for expected keywords
	expectedKeywords := []string{"machine", "learning", "artificial", "intelligence"}
	foundCount := 0
	for _, expected := range expectedKeywords {
		if keywordSet[expected] {
			foundCount++
		}
	}
	
	assert.Greater(t, foundCount, 1, "should find some expected keywords")
}

// Benchmark TopicClusterer.Process() performance
func BenchmarkTopicClustererProcess(b *testing.B) {
	clusterer := NewTopicClusterer()
	
	// Create benchmark data
	pages := make([]agents.PageInfo, 50)
	for i := 0; i < 50; i++ {
		pages[i] = agents.PageInfo{
			URL:     "https://example.com/page" + string(rune(i)),
			Title:   "Test Page " + string(rune(i)),
			Content: "This is test content for clustering with various topics and keywords technology business sports",
		}
	}

	input := ClusterRequest{
		Pages:       pages,
		NumClusters: 5,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := clusterer.Process(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}