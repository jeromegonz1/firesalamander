package semantic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// TestSemanticAnalyzer version simplifiée pour les tests
type TestSemanticAnalyzer struct {
	contentExtractor *ContentExtractor
	ngramAnalyzer   *NGramAnalyzer
	seoScorer       *SEOScorer
}

// NewTestSemanticAnalyzer crée un analyseur pour les tests
func NewTestSemanticAnalyzer() *TestSemanticAnalyzer {
	return &TestSemanticAnalyzer{
		contentExtractor: NewContentExtractor(),
		ngramAnalyzer:   NewNGramAnalyzer(),
		seoScorer:       NewSEOScorer(),
	}
}

// AnalyzePage version simplifiée pour les tests
func (tsa *TestSemanticAnalyzer) AnalyzePage(ctx context.Context, url string, htmlContent string) (*AnalysisResult, error) {
	startTime := time.Now()
	
	log.Printf("Test analysis starting for URL: %s", url)

	// Extraction du contenu
	content, err := tsa.contentExtractor.Extract(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("erreur extraction contenu: %w", err)
	}

	// Analyse locale
	localAnalysis := tsa.performLocalAnalysis(content)

	// Préparation du résultat
	result := &AnalysisResult{
		URL:             url,
		Title:           content.Title,
		MetaDescription: content.MetaDescription,
		Language:        content.Language,
		ContentType:     content.Type,
		WordCount:       content.WordCount,
		LocalAnalysis:   localAnalysis,
		ProcessedAt:     time.Now(),
		ProcessingTime:  time.Since(startTime),
		UseAI:           false,
		CacheHit:        false,
	}

	// Scoring SEO
	result.SEOScore = tsa.seoScorer.Score(content, localAnalysis, nil)

	log.Printf("Test analysis completed - Score: %.1f", result.SEOScore.Overall)

	return result, nil
}

// performLocalAnalysis version simplifiée
func (tsa *TestSemanticAnalyzer) performLocalAnalysis(content *ExtractedContent) LocalAnalysis {
	// Analyse des n-grammes
	ngrams := tsa.ngramAnalyzer.Analyze(content.CleanText)

	// Extraction des mots-clés
	keywords := tsa.extractKeywords(content, ngrams)

	// Analyse de sentiment basique
	sentiment := tsa.analyzeSentiment(content.CleanText)

	// Score de lisibilité
	readability := tsa.calculateReadability(content.CleanText)

	// Statistiques du contenu
	stats := tsa.calculateStatistics(content)

	return LocalAnalysis{
		Keywords:         keywords,
		NGrams:          ngrams,
		Topics:          []Topic{}, // Simplifié pour les tests
		Sentiment:       sentiment,
		ReadabilityScore: readability,
		Statistics:      stats,
	}
}

// extractKeywords version simplifiée
func (tsa *TestSemanticAnalyzer) extractKeywords(content *ExtractedContent, ngrams map[int][]NGram) []Keyword {
	var keywords []Keyword

	// Traiter les unigrammes
	for _, ngram := range ngrams[1] {
		if len(ngram.Text) >= 3 && ngram.Count >= 1 {
			keyword := Keyword{
				Term:      ngram.Text,
				Frequency: ngram.Count,
				Density:   float64(ngram.Count) / float64(content.WordCount) * 100,
				Relevance: ngram.Score,
				InTitle:   containsIgnoreCase(content.Title, ngram.Text),
				InMeta:    containsIgnoreCase(content.MetaDescription, ngram.Text),
			}
			keywords = append(keywords, keyword)
		}
	}

	// Limiter aux 20 premiers
	if len(keywords) > 20 {
		keywords = keywords[:20]
	}

	return keywords
}

// analyzeSentiment version simplifiée
func (tsa *TestSemanticAnalyzer) analyzeSentiment(text string) SentimentScore {
	return SentimentScore{
		Polarity:     0.0,
		Subjectivity: 0.5,
		Confidence:   0.5,
	}
}

// calculateReadability version simplifiée
func (tsa *TestSemanticAnalyzer) calculateReadability(text string) float64 {
	words := len(strings.Fields(text))
	sentences := len(strings.Split(text, "."))
	
	if sentences == 0 {
		return 50.0
	}
	
	avgWordsPerSentence := float64(words) / float64(sentences)
	
	// Score simplifié
	readabilityScore := 100 - (avgWordsPerSentence * 2)
	if readabilityScore < 0 {
		readabilityScore = 0
	}
	if readabilityScore > 100 {
		readabilityScore = 100
	}
	
	return readabilityScore
}

// calculateStatistics version simplifiée
func (tsa *TestSemanticAnalyzer) calculateStatistics(content *ExtractedContent) ContentStatistics {
	words := strings.Fields(content.CleanText)
	sentences := strings.Split(content.CleanText, ".")
	
	// Compter les mots uniques
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		uniqueWords[strings.ToLower(word)] = true
	}
	
	lexicalDiversity := 0.0
	if len(words) > 0 {
		lexicalDiversity = float64(len(uniqueWords)) / float64(len(words))
	}
	
	avgWordsPerSentence := 0.0
	if len(sentences) > 0 {
		avgWordsPerSentence = float64(len(words)) / float64(len(sentences))
	}
	
	return ContentStatistics{
		Sentences:           len(sentences),
		Paragraphs:          1, // Simplifié
		AvgWordsPerSentence: avgWordsPerSentence,
		UniqueWords:         len(uniqueWords),
		LexicalDiversity:    lexicalDiversity,
	}
}

