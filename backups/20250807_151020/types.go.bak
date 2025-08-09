package semantic

import "time"

// AnalysisResult résultat complet de l'analyse sémantique
type AnalysisResult struct {
	URL             string                 `json:"url"`
	Title           string                 `json:"title"`
	MetaDescription string                 `json:"meta_description"`
	Language        string                 `json:"language"`
	ContentType     string                 `json:"content_type"`
	WordCount       int                    `json:"word_count"`
	
	// Analyse locale (n-grammes)
	LocalAnalysis   LocalAnalysis          `json:"local_analysis"`
	
	// Enrichissement IA (sélectif)
	AIEnrichment    *AIAnalysis            `json:"ai_enrichment,omitempty"`
	
	// Scoring SEO
	SEOScore        SEOScore               `json:"seo_score"`
	
	// Métadonnées
	ProcessedAt     time.Time              `json:"processed_at"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	UseAI           bool                   `json:"use_ai"`
	CacheHit        bool                   `json:"cache_hit"`
}

// LocalAnalysis analyse sémantique locale sans IA
type LocalAnalysis struct {
	Keywords        []Keyword              `json:"keywords"`
	NGrams          map[int][]NGram        `json:"ngrams"` // 1-gram, 2-gram, 3-gram
	Topics          []Topic                `json:"topics"`
	Sentiment       SentimentScore         `json:"sentiment"`
	ReadabilityScore float64               `json:"readability_score"`
	Statistics      ContentStatistics      `json:"statistics"`
}

// AIAnalysis enrichissement via GPT-3.5 (sélectif)
type AIAnalysis struct {
	Intent          string                 `json:"intent"`
	MainTopics      []string               `json:"main_topics"`
	TargetAudience  string                 `json:"target_audience"`
	ContentGaps     []string               `json:"content_gaps"`
	Recommendations []string               `json:"recommendations"`
	CompetitiveEdge []string               `json:"competitive_edge"`
	ProcessingCost  float64                `json:"processing_cost"` // en tokens
}

// Keyword mot-clé avec métriques
type Keyword struct {
	Term        string  `json:"term"`
	Frequency   int     `json:"frequency"`
	Density     float64 `json:"density"`
	Relevance   float64 `json:"relevance"`
	Position    int     `json:"position"` // Position première occurrence
	InTitle     bool    `json:"in_title"`
	InMeta      bool    `json:"in_meta"`
	InHeadings  bool    `json:"in_headings"`
}

// NGram n-gramme avec score
type NGram struct {
	Text      string  `json:"text"`
	Count     int     `json:"count"`
	Score     float64 `json:"score"`
	Type      string  `json:"type"` // unigram, bigram, trigram
}

// Topic sujet identifié
type Topic struct {
	Name       string   `json:"name"`
	Confidence float64  `json:"confidence"`
	Keywords   []string `json:"keywords"`
	Coverage   float64  `json:"coverage"` // % du contenu
}

// SentimentScore analyse de sentiment basique
type SentimentScore struct {
	Polarity   float64 `json:"polarity"`   // -1 (négatif) à +1 (positif)
	Subjectivity float64 `json:"subjectivity"` // 0 (objectif) à 1 (subjectif)
	Confidence float64 `json:"confidence"`
}

// ContentStatistics statistiques du contenu
type ContentStatistics struct {
	Sentences      int     `json:"sentences"`
	Paragraphs     int     `json:"paragraphs"`
	AvgWordsPerSentence float64 `json:"avg_words_per_sentence"`
	AvgSentencesPerParagraph float64 `json:"avg_sentences_per_paragraph"`
	UniqueWords    int     `json:"unique_words"`
	LexicalDiversity float64 `json:"lexical_diversity"`
}

// SEOScore scoring SEO détaillé
type SEOScore struct {
	Overall         float64            `json:"overall"`
	Factors         map[string]float64 `json:"factors"`
	Recommendations []string           `json:"recommendations"`
	Issues          []string           `json:"issues"`
}

// Les types ExtractedContent, Link et Image sont définis dans content_extractor.go