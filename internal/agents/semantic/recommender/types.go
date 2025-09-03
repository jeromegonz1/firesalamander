package recommender

// RecommendationRequest represents the input data for the SemanticRecommender agent
type RecommendationRequest struct {
	Content ContentAnalysis   `json:"content"`
	Context AnalysisContext   `json:"context"`
	Options RecommendationOptions `json:"options,omitempty"`
}

// ContentAnalysis contains the analyzed content data
type ContentAnalysis struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Keywords    []string `json:"keywords"`
	Topics      []string `json:"topics"`
	WordCount   int      `json:"word_count"`
	ReadingTime int      `json:"reading_time_minutes"`
	Language    string   `json:"language,omitempty"`
}

// AnalysisContext provides contextual information for recommendations
type AnalysisContext struct {
	Domain       string   `json:"domain"`
	ContentType  string   `json:"content_type"` // "blog", "article", "product", "landing"
	TargetGoals  []string `json:"target_goals"` // "seo", "engagement", "conversion"
	Audience     string   `json:"audience,omitempty"`
	Competitors  []string `json:"competitors,omitempty"`
	Industry     string   `json:"industry,omitempty"`
}

// RecommendationOptions configures the recommendation generation
type RecommendationOptions struct {
	MaxRecommendations int    `json:"max_recommendations,omitempty"`
	Focus             string `json:"focus,omitempty"` // "content", "seo", "engagement", "comprehensive"
	Priority          string `json:"priority,omitempty"` // "high", "medium", "low", "all"
	IncludeExamples   bool   `json:"include_examples,omitempty"`
}

// Recommendation represents a semantic optimization recommendation
type Recommendation struct {
	ID            string                 `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"` // "content", "seo", "engagement", "technical"
	Type          string                 `json:"type"`     // specific recommendation type
	Impact        float64                `json:"impact"`   // expected impact score (1-10)
	Confidence    float64                `json:"confidence"` // confidence level (0-1)
	Priority      string                 `json:"priority"` // "high", "medium", "low"
	Effort        string                 `json:"effort"`   // "low", "medium", "high"
	Tags          []string               `json:"tags,omitempty"`
	Examples      []string               `json:"examples,omitempty"`
	Resources     []Resource             `json:"resources,omitempty"`
	Implementation Implementation        `json:"implementation,omitempty"`
	Metrics       RecommendationMetrics  `json:"metrics,omitempty"`
}

// Resource provides additional information or tools for implementation
type Resource struct {
	Type        string `json:"type"`        // "tool", "guide", "example", "reference"
	Title       string `json:"title"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

// Implementation provides step-by-step guidance
type Implementation struct {
	Steps         []string `json:"steps"`
	TimeEstimate  string   `json:"time_estimate"`
	Difficulty    string   `json:"difficulty"`
	Prerequisites []string `json:"prerequisites,omitempty"`
}

// RecommendationMetrics tracks the success metrics for recommendations
type RecommendationMetrics struct {
	ExpectedGains map[string]float64 `json:"expected_gains,omitempty"` // "traffic", "engagement", "conversions"
	Baseline      map[string]float64 `json:"baseline,omitempty"`
	TrackingKPIs  []string           `json:"tracking_kpis,omitempty"`
}

// SemanticScore represents the overall semantic quality score
type SemanticScore struct {
	OverallScore    float64            `json:"overall_score"`    // 0-10 scale
	CategoryScores  map[string]float64 `json:"category_scores"`  // scores by category
	Strengths       []string           `json:"strengths"`
	Weaknesses      []string           `json:"weaknesses"`
	Opportunities   []string           `json:"opportunities"`
	KeyInsights     []string           `json:"key_insights"`
}

// RecommendationMetadata provides information about the recommendation process
type RecommendationMetadata struct {
	ProcessingTime     int64  `json:"processing_time_ms"`
	RecommendationsCount int    `json:"recommendations_count"`
	HighPriorityCount  int    `json:"high_priority_count"`
	Algorithm          string `json:"algorithm"`
	Confidence         float64 `json:"confidence"`
	ContentQuality     string `json:"content_quality"` // "excellent", "good", "fair", "poor"
}

// SemanticRecommender is the main recommendation agent
type SemanticRecommender struct {
	maxRecommendations int
	defaultFocus      string
	confidenceThreshold float64
}