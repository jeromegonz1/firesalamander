package topic

import (
	"firesalamander/internal/agents"
)

// ClusterRequest represents the input data for the TopicClusterer agent
type ClusterRequest struct {
	Pages       []agents.PageInfo `json:"pages"`
	NumClusters int               `json:"num_clusters"`
	Options     ClusterOptions    `json:"options,omitempty"`
}

// ClusterOptions provides configuration for clustering
type ClusterOptions struct {
	MinPagesPerCluster int     `json:"min_pages_per_cluster,omitempty"`
	MaxKeywords        int     `json:"max_keywords,omitempty"`
	SimilarityThreshold float64 `json:"similarity_threshold,omitempty"`
}

// Cluster represents a group of semantically related pages
type Cluster struct {
	ID          string             `json:"id"`
	Label       string             `json:"label"`
	Pages       []ClusterPageInfo  `json:"pages"`
	Keywords    []ClusterKeyword   `json:"keywords"`
	Similarity  float64            `json:"similarity"`
	PageCount   int                `json:"page_count"`
	Description string             `json:"description,omitempty"`
}

// ClusterPageInfo represents a page within a cluster
type ClusterPageInfo struct {
	URL             string  `json:"url"`
	Title           string  `json:"title"`
	Relevance       float64 `json:"relevance"`
	KeywordMatches  int     `json:"keyword_matches"`
}

// ClusterKeyword represents a keyword associated with a cluster
type ClusterKeyword struct {
	Term      string  `json:"term"`
	Weight    float64 `json:"weight"`
	Frequency int     `json:"frequency"`
}

// ClusteringMetadata provides information about the clustering process
type ClusteringMetadata struct {
	TotalPages     int     `json:"total_pages"`
	ClustersCount  int     `json:"clusters_count"`
	ProcessingTime int64   `json:"processing_time_ms"`
	Algorithm      string  `json:"algorithm"`
	Quality        float64 `json:"quality"`
}

// TopicClusterer is the main clustering agent
type TopicClusterer struct {
	minPagesPerCluster  int
	maxKeywords         int
	similarityThreshold float64
}