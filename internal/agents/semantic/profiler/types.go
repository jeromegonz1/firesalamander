package profiler

import (
	"time"
	"firesalamander/internal/agents"
)

// ProfileRequest represents the input data for the PageProfiler agent
type ProfileRequest struct {
	Page    agents.PageInfo `json:"page"`
	Options ProfileOptions  `json:"options,omitempty"`
}

// ProfileOptions configures the profiling analysis
type ProfileOptions struct {
	IncludeSemanticAnalysis bool   `json:"include_semantic_analysis,omitempty"`
	IncludeContentQuality   bool   `json:"include_content_quality,omitempty"`
	IncludeTechnicalSEO     bool   `json:"include_technical_seo,omitempty"`
	DetailLevel            string `json:"detail_level,omitempty"` // "basic", "standard", "comprehensive"
	AnalysisDepth          int    `json:"analysis_depth,omitempty"`
}

// PageProfile represents the complete semantic profile of a page
type PageProfile struct {
	URL             string           `json:"url"`
	Title           string           `json:"title"`
	Language        string           `json:"language"`
	LastAnalyzed    time.Time        `json:"last_analyzed"`
	SemanticProfile SemanticProfile  `json:"semantic_profile"`
	ContentProfile  ContentProfile   `json:"content_profile"`
	TechnicalProfile TechnicalProfile `json:"technical_profile"`
	QualityScore    QualityScore     `json:"quality_score"`
	Insights        []ProfileInsight `json:"insights"`
	Recommendations []string         `json:"recommendations"`
}

// SemanticProfile contains semantic analysis results
type SemanticProfile struct {
	Keywords     []SemanticKeyword `json:"keywords"`
	Topics       []TopicInfo       `json:"topics"`
	Entities     []EntityInfo      `json:"entities"`
	Concepts     []ConceptInfo     `json:"concepts"`
	Sentiment    SentimentInfo     `json:"sentiment"`
	Complexity   float64           `json:"complexity"`
	Readability  ReadabilityInfo   `json:"readability"`
	SemanticDensity float64        `json:"semantic_density"`
}

// ContentProfile contains content structure and quality analysis
type ContentProfile struct {
	WordCount        int           `json:"word_count"`
	CharacterCount   int           `json:"character_count"`
	ParagraphCount   int           `json:"paragraph_count"`
	SentenceCount    int           `json:"sentence_count"`
	ReadingTime      int           `json:"reading_time_minutes"`
	ReadabilityScore float64       `json:"readability_score"`
	HeadingStructure []HeadingInfo `json:"heading_structure"`
	ContentType      string        `json:"content_type"`
	ContentDepth     string        `json:"content_depth"` // "shallow", "moderate", "deep"
}

// TechnicalProfile contains technical SEO analysis
type TechnicalProfile struct {
	URLStructure       URLAnalysis      `json:"url_structure"`
	TitleTag           TagAnalysis      `json:"title_tag"`
	MetaDescription    TagAnalysis      `json:"meta_description,omitempty"`
	HeadingTags        HeadingAnalysis  `json:"heading_tags"`
	CanonicalTag       TagAnalysis      `json:"canonical_tag,omitempty"`
	IndexabilityStatus string           `json:"indexability_status"`
	IsIndexable        bool             `json:"is_indexable"`
	StructuredData     StructuredInfo   `json:"structured_data,omitempty"`
}

// SemanticKeyword represents a keyword with semantic properties
type SemanticKeyword struct {
	Term           string  `json:"term"`
	Weight         float64 `json:"weight"`
	Frequency      int     `json:"frequency"`
	Position       string  `json:"position"` // "title", "heading", "body", "meta"
	Relevance      float64 `json:"relevance"`
	SemanticType   string  `json:"semantic_type"` // "primary", "secondary", "supporting"
	Context        string  `json:"context,omitempty"`
}

// TopicInfo represents a classified topic
type TopicInfo struct {
	Name        string            `json:"name"`
	Relevance   float64           `json:"relevance"`
	Confidence  float64           `json:"confidence"`
	Category    string            `json:"category"`
	Keywords    []string          `json:"keywords"`
	Subtopics   []string          `json:"subtopics,omitempty"`
	Context     map[string]string `json:"context,omitempty"`
}

// EntityInfo represents named entities extracted from content
type EntityInfo struct {
	Text       string  `json:"text"`
	Type       string  `json:"type"` // "PERSON", "ORG", "LOCATION", "PRODUCT", etc.
	Confidence float64 `json:"confidence"`
	Position   int     `json:"position"`
	Context    string  `json:"context,omitempty"`
}

// ConceptInfo represents abstract concepts identified in content
type ConceptInfo struct {
	Name       string   `json:"name"`
	Weight     float64  `json:"weight"`
	Related    []string `json:"related,omitempty"`
	Confidence float64  `json:"confidence"`
}

// SentimentInfo contains sentiment analysis results
type SentimentInfo struct {
	Overall    string  `json:"overall"` // "positive", "negative", "neutral"
	Score      float64 `json:"score"`   // -1 to 1
	Confidence float64 `json:"confidence"`
	Aspects    map[string]float64 `json:"aspects,omitempty"`
}

// ReadabilityInfo contains readability metrics
type ReadabilityInfo struct {
	FleschKincaid   float64 `json:"flesch_kincaid"`
	FleschReading   float64 `json:"flesch_reading"`
	GradeLevel      float64 `json:"grade_level"`
	AvgWordsPerSentence float64 `json:"avg_words_per_sentence"`
	AvgSyllablesPerWord float64 `json:"avg_syllables_per_word"`
}

// HeadingInfo represents heading structure information
type HeadingInfo struct {
	Level    int    `json:"level"`
	Text     string `json:"text"`
	Position int    `json:"position"`
	Keywords []string `json:"keywords,omitempty"`
}

// URLAnalysis contains URL structure analysis
type URLAnalysis struct {
	Length    int     `json:"length"`
	Score     float64 `json:"score"`
	Analysis  string  `json:"analysis"`
	HasKeywords bool  `json:"has_keywords"`
	Structure string  `json:"structure"` // "good", "fair", "poor"
}

// TagAnalysis contains HTML tag analysis
type TagAnalysis struct {
	Present  bool    `json:"present"`
	Content  string  `json:"content,omitempty"`
	Length   int     `json:"length,omitempty"`
	Score    float64 `json:"score"`
	Issues   []string `json:"issues,omitempty"`
}

// HeadingAnalysis contains heading tags analysis
type HeadingAnalysis struct {
	H1Count   int      `json:"h1_count"`
	H2Count   int      `json:"h2_count"`
	H3Count   int      `json:"h3_count"`
	Structure string   `json:"structure"` // "good", "fair", "poor"
	Issues    []string `json:"issues,omitempty"`
}

// StructuredInfo contains structured data analysis
type StructuredInfo struct {
	Present bool     `json:"present"`
	Types   []string `json:"types,omitempty"`
	Valid   bool     `json:"valid,omitempty"`
	Issues  []string `json:"issues,omitempty"`
}

// QualityScore represents overall quality assessment
type QualityScore struct {
	Overall    float64            `json:"overall"`
	Content    float64            `json:"content"`
	SEO        float64            `json:"seo"`
	Structure  float64            `json:"structure"`
	Semantic   float64            `json:"semantic"`
	Categories map[string]float64 `json:"categories"`
	Strengths  []string           `json:"strengths"`
	Weaknesses []string           `json:"weaknesses"`
}

// ProfileInsight represents actionable insights from the analysis
type ProfileInsight struct {
	Type        string  `json:"type"`        // "opportunity", "issue", "strength"
	Category    string  `json:"category"`    // "content", "seo", "semantic", "technical"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`      // "high", "medium", "low"
	Effort      string  `json:"effort"`      // "high", "medium", "low"
	Priority    int     `json:"priority"`    // 1-10
	ActionItems []string `json:"action_items,omitempty"`
}

// ProfileMetadata provides information about the profiling process
type ProfileMetadata struct {
	ProcessingTime   int64  `json:"processing_time_ms"`
	AnalysisVersion  string `json:"analysis_version"`
	DetailLevel      string `json:"detail_level"`
	ComponentsAnalyzed []string `json:"components_analyzed"`
	Confidence       float64 `json:"confidence"`
	DataQuality      string  `json:"data_quality"` // "excellent", "good", "fair", "poor"
}

// PageProfiler is the main profiling agent
type PageProfiler struct {
	analysisDepth    int
	includeInsights  bool
	qualityThreshold float64
}