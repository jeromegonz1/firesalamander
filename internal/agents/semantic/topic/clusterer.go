package topic

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"firesalamander/internal/agents"
)

// Ensure TopicClusterer implements the Agent interface
var _ agents.Agent = (*TopicClusterer)(nil)

// NewTopicClusterer creates a new TopicClusterer instance
func NewTopicClusterer() *TopicClusterer {
	return &TopicClusterer{
		minPagesPerCluster:  1,
		maxKeywords:         10,
		similarityThreshold: 0.3,
	}
}

// Name returns the agent name
func (tc *TopicClusterer) Name() string {
	return "topic-clusterer"
}

// Process implements the Agent interface for the TopicClusterer
func (tc *TopicClusterer) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()

	// Validate and parse input data
	request, err := tc.parseInput(data)
	if err != nil {
		return nil, fmt.Errorf("invalid input data: %w", err)
	}

	// Perform the clustering operation
	clusters := tc.clusterPages(request.Pages, request.NumClusters)
	
	// Extract topic themes
	topics := tc.extractTopics(clusters)
	
	// Create clustering metadata
	metadata := ClusteringMetadata{
		TotalPages:     len(request.Pages),
		ClustersCount:  len(clusters),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Algorithm:      "k-means-text",
		Quality:        tc.calculateClusteringQuality(clusters),
	}

	// Create agent result
	agentResult := &agents.AgentResult{
		AgentName: tc.Name(),
		Status:    "success",
		Data: map[string]interface{}{
			"clusters": clusters,
			"topics":   topics,
			"metadata": metadata,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}

	return agentResult, nil
}

// HealthCheck implements the Agent interface
func (tc *TopicClusterer) HealthCheck() error {
	if tc.minPagesPerCluster < 1 {
		return fmt.Errorf("minimum pages per cluster must be >= 1")
	}
	
	if tc.maxKeywords < 1 {
		return fmt.Errorf("maximum keywords must be >= 1")
	}
	
	if tc.similarityThreshold < 0 || tc.similarityThreshold > 1 {
		return fmt.Errorf("similarity threshold must be between 0 and 1")
	}
	
	return nil
}

// parseInput validates and parses the input data into a ClusterRequest
func (tc *TopicClusterer) parseInput(data interface{}) (*ClusterRequest, error) {
	if data == nil {
		return nil, fmt.Errorf("input data is nil")
	}

	request, ok := data.(ClusterRequest)
	if !ok {
		return nil, fmt.Errorf("expected ClusterRequest, got %T", data)
	}

	// Apply default values if not specified
	if request.NumClusters <= 0 {
		request.NumClusters = int(math.Max(1, float64(len(request.Pages))/5)) // Default: 1 cluster per 5 pages
	}

	// Apply options defaults
	if request.Options.MinPagesPerCluster > 0 {
		tc.minPagesPerCluster = request.Options.MinPagesPerCluster
	}
	if request.Options.MaxKeywords > 0 {
		tc.maxKeywords = request.Options.MaxKeywords
	}
	if request.Options.SimilarityThreshold > 0 {
		tc.similarityThreshold = request.Options.SimilarityThreshold
	}

	return &request, nil
}

// clusterPages performs the main clustering logic
func (tc *TopicClusterer) clusterPages(pages []agents.PageInfo, numClusters int) []Cluster {
	if len(pages) == 0 {
		return []Cluster{}
	}
	
	// Adjust cluster count based on available pages
	actualClusters := int(math.Min(float64(numClusters), float64(len(pages))))
	
	clusters := make([]Cluster, actualClusters)
	
	// Initialize clusters
	for i := 0; i < actualClusters; i++ {
		clusters[i] = Cluster{
			ID:          fmt.Sprintf("cluster-%d", i+1),
			Label:       "",
			Pages:       make([]ClusterPageInfo, 0),
			Keywords:    make([]ClusterKeyword, 0),
			Similarity:  0.0,
			PageCount:   0,
			Description: "",
		}
	}
	
	// Simple clustering algorithm: assign pages based on content similarity
	for _, page := range pages {
		bestClusterIdx := tc.findBestCluster(page, clusters)
		
		clusterPage := ClusterPageInfo{
			URL:            page.URL,
			Title:          page.Title,
			Relevance:      tc.calculatePageRelevance(page, clusters[bestClusterIdx]),
			KeywordMatches: tc.countKeywordMatches(page, clusters[bestClusterIdx]),
		}
		
		clusters[bestClusterIdx].Pages = append(clusters[bestClusterIdx].Pages, clusterPage)
		clusters[bestClusterIdx].PageCount++
	}
	
	// Generate cluster labels and keywords
	for i := range clusters {
		if clusters[i].PageCount > 0 {
			clusters[i].Keywords = tc.extractClusterKeywords(tc.getClusterPages(pages, clusters[i]))
			clusters[i].Label = tc.generateClusterLabel(clusters[i].Keywords)
			clusters[i].Similarity = tc.calculateIntraClusterSimilarity(clusters[i])
			clusters[i].Description = tc.generateClusterDescription(clusters[i])
		}
	}
	
	// Filter out empty clusters
	var nonEmptyClusters []Cluster
	for _, cluster := range clusters {
		if cluster.PageCount >= tc.minPagesPerCluster {
			nonEmptyClusters = append(nonEmptyClusters, cluster)
		}
	}
	
	return nonEmptyClusters
}

// findBestCluster finds the most appropriate cluster for a page
func (tc *TopicClusterer) findBestCluster(page agents.PageInfo, clusters []Cluster) int {
	if len(clusters) == 0 {
		return 0
	}
	
	bestScore := -1.0
	bestIdx := 0
	
	// For the first few pages, distribute them across different clusters
	totalPages := 0
	for _, cluster := range clusters {
		totalPages += cluster.PageCount
	}
	
	if totalPages < len(clusters) {
		// Find first empty cluster
		for i, cluster := range clusters {
			if cluster.PageCount == 0 {
				return i
			}
		}
	}
	
	for i, cluster := range clusters {
		score := tc.calculatePageClusterSimilarity(page, cluster)
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	
	// If similarity is too low, prefer less populated clusters to encourage distribution
	if bestScore < tc.similarityThreshold {
		minPages := clusters[bestIdx].PageCount
		for i, cluster := range clusters {
			if cluster.PageCount < minPages {
				minPages = cluster.PageCount
				bestIdx = i
			}
		}
	}
	
	return bestIdx
}

// calculatePageClusterSimilarity calculates similarity between a page and cluster
func (tc *TopicClusterer) calculatePageClusterSimilarity(page agents.PageInfo, cluster Cluster) float64 {
	if cluster.PageCount == 0 {
		return 0.1 // Low score for empty clusters to encourage distribution
	}
	
	pageWords := tc.extractWords(page.Title + " " + page.Content)
	
	if len(cluster.Keywords) == 0 {
		// If cluster has no keywords yet, calculate similarity with existing pages
		clusterPages := cluster.Pages
		if len(clusterPages) == 0 {
			return 0.1
		}
		
		// Calculate average similarity with existing pages
		totalSim := 0.0
		for _, clusterPage := range clusterPages {
			clusterWords := tc.extractWords(clusterPage.Title)
			sim := tc.calculateWordSimilarity(pageWords, clusterWords)
			totalSim += sim
		}
		
		return totalSim / float64(len(clusterPages))
	}
	
	// Calculate similarity based on keyword overlap with weighting
	totalWeight := 0.0
	matchedWeight := 0.0
	
	for _, keyword := range cluster.Keywords {
		totalWeight += keyword.Weight
		
		for _, word := range pageWords {
			if strings.EqualFold(word, keyword.Term) {
				matchedWeight += keyword.Weight
				break // Only count each keyword once
			}
		}
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	return matchedWeight / totalWeight
}

// calculateWordSimilarity calculates similarity between two word sets
func (tc *TopicClusterer) calculateWordSimilarity(words1, words2 []string) float64 {
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}
	
	set1 := make(map[string]bool)
	for _, word := range words1 {
		set1[strings.ToLower(word)] = true
	}
	
	matches := 0
	for _, word := range words2 {
		if set1[strings.ToLower(word)] {
			matches++
		}
	}
	
	// Jaccard similarity
	union := len(words1) + len(words2) - matches
	if union == 0 {
		return 0.0
	}
	
	return float64(matches) / float64(union)
}

// extractClusterKeywords extracts keywords from cluster pages
func (tc *TopicClusterer) extractClusterKeywords(pages []agents.PageInfo) []ClusterKeyword {
	wordFreq := make(map[string]int)
	totalWords := 0
	
	// Count word frequencies
	for _, page := range pages {
		words := tc.extractWords(page.Title + " " + page.Content)
		for _, word := range words {
			if len(word) >= 3 { // Filter short words
				wordFreq[strings.ToLower(word)]++
				totalWords++
			}
		}
	}
	
	// Convert to keywords with weights
	var keywords []ClusterKeyword
	for term, freq := range wordFreq {
		if freq > 1 || len(pages) == 1 { // Must appear more than once, unless single page
			weight := float64(freq) / float64(totalWords)
			keywords = append(keywords, ClusterKeyword{
				Term:      term,
				Weight:    weight,
				Frequency: freq,
			})
		}
	}
	
	// Sort by weight descending
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Weight > keywords[j].Weight
	})
	
	// Limit to max keywords
	if len(keywords) > tc.maxKeywords {
		keywords = keywords[:tc.maxKeywords]
	}
	
	return keywords
}

// generateClusterLabel generates a descriptive label for a cluster
func (tc *TopicClusterer) generateClusterLabel(keywords []ClusterKeyword) string {
	if len(keywords) == 0 {
		return "Untitled Cluster"
	}
	
	// Use top keywords to create label
	topKeywords := make([]string, 0, 3)
	for i, kw := range keywords {
		if i >= 3 {
			break
		}
		topKeywords = append(topKeywords, strings.Title(kw.Term))
	}
	
	return strings.Join(topKeywords, " & ")
}

// extractWords extracts individual words from text
func (tc *TopicClusterer) extractWords(text string) []string {
	// Simple word extraction
	words := strings.Fields(strings.ToLower(text))
	var cleanWords []string
	
	for _, word := range words {
		// Remove common punctuation
		word = strings.Trim(word, ".,!?;:()[]{}\"'-")
		if len(word) > 0 && !tc.isStopWord(word) {
			cleanWords = append(cleanWords, word)
		}
	}
	
	return cleanWords
}

// isStopWord checks if a word is a common stop word
func (tc *TopicClusterer) isStopWord(word string) bool {
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
		"this": true, "these": true, "they": true, "we": true, "you": true, "your": true,
		"our": true, "their": true, "his": true, "her": true, "him": true, "them": true,
	}
	
	return stopWords[word]
}

// getClusterPages returns the original page info for pages in a cluster
func (tc *TopicClusterer) getClusterPages(allPages []agents.PageInfo, cluster Cluster) []agents.PageInfo {
	var clusterPages []agents.PageInfo
	
	pageMap := make(map[string]agents.PageInfo)
	for _, page := range allPages {
		pageMap[page.URL] = page
	}
	
	for _, clusterPage := range cluster.Pages {
		if page, exists := pageMap[clusterPage.URL]; exists {
			clusterPages = append(clusterPages, page)
		}
	}
	
	return clusterPages
}

// calculatePageRelevance calculates how relevant a page is to a cluster
func (tc *TopicClusterer) calculatePageRelevance(page agents.PageInfo, cluster Cluster) float64 {
	// Simple relevance based on keyword overlap
	return tc.calculatePageClusterSimilarity(page, cluster)
}

// countKeywordMatches counts how many cluster keywords match the page
func (tc *TopicClusterer) countKeywordMatches(page agents.PageInfo, cluster Cluster) int {
	pageWords := tc.extractWords(page.Title + " " + page.Content)
	matches := 0
	
	for _, word := range pageWords {
		for _, keyword := range cluster.Keywords {
			if strings.EqualFold(word, keyword.Term) {
				matches++
			}
		}
	}
	
	return matches
}

// calculateIntraClusterSimilarity calculates similarity within cluster
func (tc *TopicClusterer) calculateIntraClusterSimilarity(cluster Cluster) float64 {
	if cluster.PageCount <= 1 {
		return 1.0
	}
	
	// Simple similarity metric based on keyword coverage
	if len(cluster.Keywords) == 0 {
		return 0.0
	}
	
	totalMatches := 0
	for _, page := range cluster.Pages {
		totalMatches += page.KeywordMatches
	}
	
	avgMatches := float64(totalMatches) / float64(cluster.PageCount)
	maxPossibleMatches := float64(len(cluster.Keywords))
	
	if maxPossibleMatches == 0 {
		return 0.0
	}
	
	return avgMatches / maxPossibleMatches
}

// generateClusterDescription generates a description for the cluster
func (tc *TopicClusterer) generateClusterDescription(cluster Cluster) string {
	if cluster.PageCount == 0 {
		return "Empty cluster"
	}
	
	return fmt.Sprintf("Cluster of %d pages about %s", cluster.PageCount, cluster.Label)
}

// extractTopics extracts high-level topics from clusters
func (tc *TopicClusterer) extractTopics(clusters []Cluster) []string {
	var topics []string
	
	for _, cluster := range clusters {
		if cluster.PageCount > 0 {
			topics = append(topics, cluster.Label)
		}
	}
	
	return topics
}

// calculateClusteringQuality calculates overall quality of clustering
func (tc *TopicClusterer) calculateClusteringQuality(clusters []Cluster) float64 {
	if len(clusters) == 0 {
		return 0.0
	}
	
	totalSimilarity := 0.0
	totalPages := 0
	
	for _, cluster := range clusters {
		totalSimilarity += cluster.Similarity * float64(cluster.PageCount)
		totalPages += cluster.PageCount
	}
	
	if totalPages == 0 {
		return 0.0
	}
	
	return totalSimilarity / float64(totalPages)
}