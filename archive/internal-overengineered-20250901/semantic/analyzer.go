package semantic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"firesalamander/internal/config"
)

// SemanticAnalyzer analyseur sémantique complet avec IA
type SemanticAnalyzer struct {
	contentExtractor *ContentExtractor
	ngramAnalyzer   *NGramAnalyzer
	seoScorer       *SEOScorer
	aiEnricher      *AIEnricher
	aiEnabled       bool
	aiModel         string
}

// NewSemanticAnalyzer crée un nouvel analyseur sémantique
func NewSemanticAnalyzer(cfg *config.Config) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		contentExtractor: NewContentExtractor(),
		ngramAnalyzer:   NewNGramAnalyzer(),
		seoScorer:       NewSEOScorer(),
		aiEnricher:      NewAIEnricher(cfg),
		aiEnabled:       cfg.OpenAIAPIKey != "",
		aiModel:         cfg.OpenAIModel,
	}
}

// AnalyzePage analyse une page web complète avec enrichissement IA
func (sa *SemanticAnalyzer) AnalyzePage(ctx context.Context, url string, htmlContent string) (*AnalysisResult, error) {
	startTime := time.Now()
	
	log.Printf("🧠 Starting semantic analysis for URL: %s", url)

	// 1. Extraction du contenu
	content, err := sa.contentExtractor.Extract(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("erreur extraction contenu: %w", err)
	}

	// 2. Analyse locale (gratuite)
	localAnalysis := sa.performLocalAnalysis(content)

	// 3. Déterminer si on utilise l'IA
	useAI := sa.shouldUseAI(localAnalysis)
	
	// 4. Enrichissement IA (si nécessaire)
	var aiEnrichment *AIAnalysis
	if useAI && sa.aiEnabled {
		keywords := sa.extractTopKeywords(localAnalysis.Keywords, 5) // Limiter à 5 pour économiser
		enrichmentResult, err := sa.aiEnricher.EnrichKeywords(ctx, keywords, content.CleanText)
		if err != nil {
			log.Printf("⚠️ AI enrichment failed: %v", err)
			useAI = false
		} else {
			// Convertir EnrichmentResult vers AIAnalysis
			aiEnrichment = sa.convertToAIAnalysis(enrichmentResult)
		}
	}

	// 5. Préparation du résultat
	result := &AnalysisResult{
		URL:             url,
		Title:           content.Title,
		MetaDescription: content.MetaDescription,
		Language:        content.Language,
		ContentType:     content.Type,
		WordCount:       content.WordCount,
		LocalAnalysis:   localAnalysis,
		AIEnrichment:    aiEnrichment,
		ProcessedAt:     time.Now(),
		ProcessingTime:  time.Since(startTime),
		UseAI:           useAI,
		CacheHit:        false, // TODO: implémenter cache
	}

	// 6. Scoring SEO final
	result.SEOScore = sa.seoScorer.Score(content, localAnalysis, aiEnrichment)

	log.Printf("✅ Semantic analysis completed - Score: %.1f (AI: %v)", result.SEOScore.Overall, useAI)

	return result, nil
}

// performLocalAnalysis effectue l'analyse locale (gratuite)
func (sa *SemanticAnalyzer) performLocalAnalysis(content *ExtractedContent) LocalAnalysis {
	// Analyse des n-grammes
	ngrams := sa.ngramAnalyzer.Analyze(content.CleanText)

	// Extraction des mots-clés
	keywords := sa.extractKeywords(content, ngrams)

	// Analyse de sentiment basique
	sentiment := sa.analyzeSentiment(content.CleanText)

	// Score de lisibilité
	readability := sa.calculateReadability(content.CleanText)

	// Extraction des topics
	topics := sa.extractTopics(keywords, ngrams)

	// Statistiques du contenu
	stats := sa.calculateStatistics(content)

	return LocalAnalysis{
		Keywords:         keywords,
		NGrams:          ngrams,
		Topics:          topics,
		Sentiment:       sentiment,
		ReadabilityScore: readability,
		Statistics:      stats,
	}
}

// shouldUseAI détermine si on doit utiliser l'IA basé sur la confiance de l'analyse locale
func (sa *SemanticAnalyzer) shouldUseAI(localAnalysis LocalAnalysis) bool {
	if !sa.aiEnabled {
		return false
	}

	// Calculer un score de confiance basé sur l'analyse locale
	confidenceScore := sa.calculateConfidenceScore(localAnalysis)
	
	// Utiliser l'IA si la confiance est faible (< 0.8)
	return confidenceScore < 0.8
}

// calculateConfidenceScore calcule un score de confiance de l'analyse locale
func (sa *SemanticAnalyzer) calculateConfidenceScore(localAnalysis LocalAnalysis) float64 {
	confidence := 0.0
	factors := 0

	// Facteur 1: Nombre de mots-clés pertinents
	if len(localAnalysis.Keywords) > 0 {
		keywordConfidence := float64(len(localAnalysis.Keywords)) / 20.0 // Max 20 keywords
		if keywordConfidence > 1.0 {
			keywordConfidence = 1.0
		}
		confidence += keywordConfidence
		factors++
	}

	// Facteur 2: Diversité lexicale
	if localAnalysis.Statistics.LexicalDiversity > 0 {
		confidence += localAnalysis.Statistics.LexicalDiversity
		factors++
	}

	// Facteur 3: Longueur du contenu
	if localAnalysis.Statistics.Sentences > 0 {
		lengthConfidence := float64(localAnalysis.Statistics.Sentences) / 50.0 // Max 50 phrases
		if lengthConfidence > 1.0 {
			lengthConfidence = 1.0
		}
		confidence += lengthConfidence
		factors++
	}

	if factors == 0 {
		return 0.0
	}

	return confidence / float64(factors)
}

// extractTopKeywords extrait les N meilleurs mots-clés
func (sa *SemanticAnalyzer) extractTopKeywords(keywords []Keyword, limit int) []string {
	result := make([]string, 0, limit)
	
	for i, kw := range keywords {
		if i >= limit {
			break
		}
		result = append(result, kw.Term)
	}
	
	return result
}

// extractKeywords extrait les mots-clés du contenu
func (sa *SemanticAnalyzer) extractKeywords(content *ExtractedContent, ngrams map[int][]NGram) []Keyword {
	var keywords []Keyword

	// Traiter les unigrammes (mots simples)
	for _, ngram := range ngrams[1] {
		if len(ngram.Text) >= 3 && ngram.Count >= 2 {
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

	// Traiter les bigrammes (expressions de 2 mots)
	for _, ngram := range ngrams[2] {
		if ngram.Count >= 2 {
			keyword := Keyword{
				Term:      ngram.Text,
				Frequency: ngram.Count,
				Density:   float64(ngram.Count) / float64(content.WordCount/2) * 100, // Ajuster pour bigrammes
				Relevance: ngram.Score,
				InTitle:   containsIgnoreCase(content.Title, ngram.Text),
				InMeta:    containsIgnoreCase(content.MetaDescription, ngram.Text),
			}
			keywords = append(keywords, keyword)
		}
	}

	// Trier par pertinence et limiter
	if len(keywords) > 30 {
		keywords = keywords[:30]
	}

	return keywords
}

// extractTopics extrait les sujets principaux
func (sa *SemanticAnalyzer) extractTopics(keywords []Keyword, ngrams map[int][]NGram) []Topic {
	topicMap := make(map[string]*Topic)

	// Regrouper les mots-clés par similarité sémantique
	for _, kw := range keywords {
		topicKey := sa.getTopicKey(kw.Term)
		
		if topic, exists := topicMap[topicKey]; exists {
			topic.Keywords = append(topic.Keywords, kw.Term)
			topic.Confidence += kw.Relevance
			topic.Coverage++
		} else {
			topicMap[topicKey] = &Topic{
				Name:       topicKey,
				Keywords:   []string{kw.Term},
				Confidence: kw.Relevance,
				Coverage:   1,
			}
		}
	}

	// Convertir en slice et normaliser les scores
	var topics []Topic
	for _, topic := range topicMap {
		topic.Confidence /= topic.Coverage // Moyenne
		topics = append(topics, *topic)
	}

	// Limiter aux 10 premiers topics
	if len(topics) > 10 {
		topics = topics[:10]
	}

	return topics
}

// getTopicKey détermine la clé de sujet pour un terme (simplifié)
func (sa *SemanticAnalyzer) getTopicKey(term string) string {
	// Logique simplifiée - peut être améliorée avec des algorithmes plus sophistiqués
	lowerTerm := strings.ToLower(term)
	
	// Catégorisation basique par domaines
	if strings.Contains(lowerTerm, "seo") || strings.Contains(lowerTerm, "référencement") {
		return "SEO"
	}
	if strings.Contains(lowerTerm, "web") || strings.Contains(lowerTerm, "site") {
		return "Web Development"
	}
	if strings.Contains(lowerTerm, "marketing") || strings.Contains(lowerTerm, "publicité") {
		return "Marketing"
	}
	if strings.Contains(lowerTerm, "contenu") || strings.Contains(lowerTerm, "article") {
		return "Content"
	}
	
	// Par défaut, utiliser le terme lui-même
	return strings.Title(term)
}

// analyzeSentiment effectue une analyse de sentiment basique
func (sa *SemanticAnalyzer) analyzeSentiment(text string) SentimentScore {
	// Analyse de sentiment simplifiée basée sur des mots-clés
	positiveWords := []string{"bon", "excellent", "super", "parfait", "génial", "formidable"}
	negativeWords := []string{"mauvais", "terrible", "horrible", "nul", "décevant", "problème"}
	
	lowerText := strings.ToLower(text)
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range positiveWords {
		positiveCount += strings.Count(lowerText, word)
	}
	
	for _, word := range negativeWords {
		negativeCount += strings.Count(lowerText, word)
	}
	
	totalWords := len(strings.Fields(text))
	if totalWords == 0 {
		return SentimentScore{Polarity: 0.0, Subjectivity: 0.0, Confidence: 0.0}
	}
	
	polarity := (float64(positiveCount) - float64(negativeCount)) / float64(totalWords) * 100
	subjectivity := float64(positiveCount+negativeCount) / float64(totalWords)
	confidence := 0.5 // Confiance modérée pour l'analyse locale
	
	return SentimentScore{
		Polarity:     polarity,
		Subjectivity: subjectivity,
		Confidence:   confidence,
	}
}

// calculateReadability calcule un score de lisibilité
func (sa *SemanticAnalyzer) calculateReadability(text string) float64 {
	words := strings.Fields(text)
	sentences := strings.Split(text, ".")
	
	if len(sentences) == 0 || len(words) == 0 {
		return 50.0
	}
	
	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))
	
	// Algorithme simplifié inspiré de Flesch
	readabilityScore := 206.835 - (1.015 * avgWordsPerSentence)
	
	// Normaliser entre 0 et 100
	if readabilityScore < 0 {
		readabilityScore = 0
	}
	if readabilityScore > 100 {
		readabilityScore = 100
	}
	
	return readabilityScore
}

// calculateStatistics calcule les statistiques du contenu
func (sa *SemanticAnalyzer) calculateStatistics(content *ExtractedContent) ContentStatistics {
	words := strings.Fields(content.CleanText)
	sentences := strings.Split(content.CleanText, ".")
	paragraphs := strings.Split(content.CleanText, "\n\n")
	
	// Compter les mots uniques
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		cleanWord := strings.ToLower(strings.Trim(word, ".,!?;:"))
		if len(cleanWord) > 2 { // Ignorer les mots très courts
			uniqueWords[cleanWord] = true
		}
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
		Paragraphs:          len(paragraphs),
		AvgWordsPerSentence: avgWordsPerSentence,
		UniqueWords:         len(uniqueWords),
		LexicalDiversity:    lexicalDiversity,
	}
}

// GetStats retourne les statistiques de l'analyseur
func (sa *SemanticAnalyzer) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"analyzer_type": "semantic",
		"ai_enabled":    sa.aiEnabled,
		"ai_model":      sa.aiModel,
	}
	
	// Ajouter les stats du cache IA
	if sa.aiEnricher != nil {
		aiStats := sa.aiEnricher.GetCacheStats()
		stats["ai_cache"] = aiStats
	}
	
	return stats
}

// convertToAIAnalysis convertit EnrichmentResult vers AIAnalysis
func (sa *SemanticAnalyzer) convertToAIAnalysis(enrichment *EnrichmentResult) *AIAnalysis {
	var mainTopics []string
	var recommendations []string
	
	// Extraire les topics principaux des mots-clés enrichis
	for _, kw := range enrichment.Keywords {
		if kw.Intent != "" {
			mainTopics = append(mainTopics, kw.Intent)
		}
		recommendations = append(recommendations, kw.ContentSuggestions...)
	}
	
	// Ajouter les recommandations SEO
	recommendations = append(recommendations, enrichment.SEORecommendations...)
	
	return &AIAnalysis{
		Intent:          enrichment.SearchIntent,
		MainTopics:      mainTopics,
		TargetAudience:  "General", // Simplifié pour l'instant
		ContentGaps:     enrichment.ContentQuestions,
		Recommendations: recommendations,
		CompetitiveEdge: []string{fmt.Sprintf("Competitiveness Score: %d/100", enrichment.CompetitivenessScore)},
		ProcessingCost:  0.001, // Estimation ~1/1000 dollar par requête
	}
}

// containsIgnoreCase vérifie si une chaîne contient une sous-chaîne (insensible à la casse)
func containsIgnoreCase(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}