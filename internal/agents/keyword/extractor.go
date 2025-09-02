package keyword

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
)

// KeywordExtractor implémente l'agent d'extraction de mots-clés SEO
type KeywordExtractor struct {
	name           string
	minWordLength  int
	maxKeywords    int
	stopWords      map[string]bool
	langDetector   *LanguageDetector
}

// NewKeywordExtractor crée une nouvelle instance de KeywordExtractor
func NewKeywordExtractor() *KeywordExtractor {
	stopWords := map[string]bool{
		// Français
		"le": true, "la": true, "les": true, "de": true, "du": true, "des": true,
		"et": true, "ou": true, "un": true, "une": true, "ce": true, "cette": true,
		"pour": true, "avec": true, "sur": true, "dans": true, "par": true,
		// Anglais
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "of": true, "with": true, "by": true,
	}

	return &KeywordExtractor{
		name:          constants.AgentNameKeyword,
		minWordLength: constants.KeywordMinLength,
		maxKeywords:   constants.KeywordMaxCount,
		stopWords:     stopWords,
		langDetector:  NewLanguageDetector(),
	}
}

// Name retourne le nom de l'agent
func (k *KeywordExtractor) Name() string {
	return k.name
}

// Process traite les données d'entrée et extrait les mots-clés
func (k *KeywordExtractor) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()
	
	content, ok := data.(string)
	if !ok {
		return &agents.AgentResult{
			AgentName: k.name,
			Status:    constants.StatusFailed,
			Errors:    []string{"invalid input data type, expected string"},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	keywords, err := k.ExtractKeywords(content)
	if err != nil {
		return &agents.AgentResult{
			AgentName: k.name,
			Status:    constants.StatusFailed,
			Errors:    []string{err.Error()},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	return &agents.AgentResult{
		AgentName: k.name,
		Status:    constants.StatusCompleted,
		Data: map[string]interface{}{
			"keywords": keywords,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}, nil
}

// HealthCheck vérifie la santé de l'agent
func (k *KeywordExtractor) HealthCheck() error {
	// Test simple d'extraction
	_, err := k.ExtractKeywords("test content for health check")
	return err
}

// ExtractKeywords extrait et analyse les mots-clés du contenu
func (k *KeywordExtractor) ExtractKeywords(content string) (*agents.KeywordResult, error) {
	if content == "" {
		return &agents.KeywordResult{
			Keywords:   []agents.KeywordItem{},
			TotalCount: 0,
			Language:   "unknown",
		}, nil
	}

	// Détection de la langue
	language := k.langDetector.DetectLanguage(content)

	// Nettoyage du contenu
	cleanContent := k.cleanContent(content)

	// Extraction des mots
	words := k.extractWords(cleanContent)

	// Comptage et filtrage
	wordCounts := k.countWords(words)

	// Génération des KeywordItems
	keywords := k.generateKeywordItems(wordCounts, len(words))

	// Tri par pertinence
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Relevance > keywords[j].Relevance
	})

	// Limitation du nombre de résultats
	if len(keywords) > k.maxKeywords {
		keywords = keywords[:k.maxKeywords]
	}

	return &agents.KeywordResult{
		Keywords:   keywords,
		TotalCount: len(keywords),
		Language:   language,
	}, nil
}

// AnalyzeDensity analyse la densité des mots-clés dans le contenu
func (k *KeywordExtractor) AnalyzeDensity(keywords []string, content string) (*agents.DensityReport, error) {
	keywordMetrics := make(map[string]float64)
	
	if content == "" {
		// Même pour un contenu vide, on doit initialiser les métriques des mots-clés
		for _, keyword := range keywords {
			keywordMetrics[keyword] = 0.0
		}
		return &agents.DensityReport{
			TotalWords:      0,
			KeywordMetrics:  keywordMetrics,
			Recommendations: []string{"Content is empty"},
		}, nil
	}

	cleanContent := k.cleanContent(content)
	words := k.extractWords(cleanContent)
	totalWords := len(words)

	wordCounts := k.countWords(words)

	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		count := wordCounts[lowerKeyword]
		density := 0.0
		if totalWords > 0 {
			density = float64(count) / float64(totalWords) * 100.0
		}
		keywordMetrics[keyword] = density
	}

	recommendations := k.generateDensityRecommendations(keywordMetrics)

	return &agents.DensityReport{
		TotalWords:      totalWords,
		KeywordMetrics:  keywordMetrics,
		Recommendations: recommendations,
	}, nil
}

// cleanContent nettoie le contenu HTML et supprime les caractères indésirables
func (k *KeywordExtractor) cleanContent(content string) string {
	// Suppression des entités HTML communes
	clean := strings.ReplaceAll(content, "&amp;", " ")
	clean = strings.ReplaceAll(clean, "&lt;", " ")
	clean = strings.ReplaceAll(clean, "&gt;", " ")
	clean = strings.ReplaceAll(clean, "&nbsp;", " ")
	
	// Suppression des balises HTML
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	clean = htmlTagRegex.ReplaceAllString(clean, " ")

	// Suppression des caractères spéciaux, garde seulement lettres, chiffres et espaces
	specialCharsRegex := regexp.MustCompile(`[^a-zA-ZÀ-ÿ0-9\s]`)
	clean = specialCharsRegex.ReplaceAllString(clean, " ")

	// Normalisation des espaces
	spaceRegex := regexp.MustCompile(`\s+`)
	clean = spaceRegex.ReplaceAllString(clean, " ")

	return strings.TrimSpace(clean)
}

// extractWords extrait les mots individuels du contenu
func (k *KeywordExtractor) extractWords(content string) []string {
	words := strings.Fields(strings.ToLower(content))
	var validWords []string

	for _, word := range words {
		if len(word) >= k.minWordLength && !k.stopWords[word] {
			validWords = append(validWords, word)
		}
	}

	return validWords
}

// countWords compte l'occurrence de chaque mot
func (k *KeywordExtractor) countWords(words []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}

// generateKeywordItems génère les éléments de mots-clés avec métriques
func (k *KeywordExtractor) generateKeywordItems(wordCounts map[string]int, totalWords int) []agents.KeywordItem {
	var keywords []agents.KeywordItem

	for word, count := range wordCounts {
		density := float64(count) / float64(totalWords) * 100.0
		relevance := k.calculateRelevance(word, count, density)

		keywords = append(keywords, agents.KeywordItem{
			Term:      word,
			Count:     count,
			Density:   density,
			Relevance: relevance,
		})
	}

	return keywords
}

// calculateRelevance calcule la pertinence d'un mot-clé
func (k *KeywordExtractor) calculateRelevance(word string, count int, density float64) float64 {
	// Score basé sur la longueur du mot (mots plus longs = plus pertinents)
	lengthScore := float64(len(word)) / 10.0

	// Score basé sur la fréquence (mais pas trop élevée)
	frequencyScore := float64(count)
	if density > 5.0 { // Pénalise la sur-optimisation
		frequencyScore = frequencyScore * 0.5
	}

	// Score basé sur la densité optimale (1-3% idéal)
	densityScore := 1.0
	if density >= 1.0 && density <= 3.0 {
		densityScore = 2.0
	} else if density > 3.0 {
		densityScore = 0.5
	}

	return (lengthScore + frequencyScore + densityScore) / 3.0
}

// generateDensityRecommendations génère des recommandations basées sur la densité
func (k *KeywordExtractor) generateDensityRecommendations(metrics map[string]float64) []string {
	var recommendations []string

	for keyword, density := range metrics {
		if density < 0.5 {
			recommendations = append(recommendations, 
				"Mot-clé '"+keyword+"' sous-utilisé ("+formatFloat(density)+"%). Considérez augmenter sa présence.")
		} else if density > 5.0 {
			recommendations = append(recommendations, 
				"Mot-clé '"+keyword+"' sur-optimisé ("+formatFloat(density)+"%). Risque de pénalité SEO.")
		} else if density >= 1.0 && density <= 3.0 {
			recommendations = append(recommendations, 
				"Mot-clé '"+keyword+"' bien optimisé ("+formatFloat(density)+"%).")
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Densité des mots-clés dans la plage acceptable.")
	}

	return recommendations
}

// formatFloat formate un float64 avec 2 décimales
func formatFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(sprintf("%.2f", f), "0"), ".")
}

// sprintf simple wrapper pour éviter l'import fmt
func sprintf(format string, a ...interface{}) string {
	// Implementation simple pour le formatage
	result := format
	for _, arg := range a {
		switch v := arg.(type) {
		case float64:
			// Conversion simple de float64 vers string avec 2 décimales
			intPart := int(v)
			fracPart := int((v - float64(intPart)) * 100)
			result = strings.Replace(result, "%.2f", intToString(intPart)+"."+intToString(fracPart), 1)
		}
	}
	return result
}

// intToString convertit un int en string
func intToString(i int) string {
	if i == 0 {
		return "0"
	}
	
	var result strings.Builder
	if i < 0 {
		result.WriteString("-")
		i = -i
	}
	
	var digits []byte
	for i > 0 {
		digits = append([]byte{'0' + byte(i%10)}, digits...)
		i /= 10
	}
	
	result.Write(digits)
	return result.String()
}