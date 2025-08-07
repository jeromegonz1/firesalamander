package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"firesalamander/internal/config"
)

// AIEnricher utilise OpenAI pour enrichir l'analyse sémantique
type AIEnricher struct {
	config     *config.AIConfig
	httpClient *http.Client
	cache      map[string]*EnrichmentResult
	cacheTTL   time.Duration
}

// EnrichmentResult résultat de l'enrichissement IA
type EnrichmentResult struct {
	Keywords           []EnrichedKeyword    `json:"keywords"`
	ContentQuestions   []string            `json:"content_questions"`
	SEORecommendations []string            `json:"seo_recommendations"`
	SearchIntent       string              `json:"search_intent"`
	CompetitivenessScore int               `json:"competitiveness_score"`
	CachedAt           time.Time           `json:"cached_at"`
}

// EnrichedKeyword mot-clé enrichi par l'IA
type EnrichedKeyword struct {
	Text                string   `json:"text"`
	Intent             string   `json:"intent"`
	Difficulty         int      `json:"difficulty"`
	SearchVolume       string   `json:"search_volume"`
	RelatedKeywords    []string `json:"related_keywords"`
	ContentSuggestions []string `json:"content_suggestions"`
}

// OpenAIRequest structure pour les requêtes OpenAI
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message structure pour les messages OpenAI
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse structure pour les réponses OpenAI
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// NewAIEnricher crée un nouvel enrichisseur IA
func NewAIEnricher(cfg *config.AIConfig) *AIEnricher {
	return &AIEnricher{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		cache:    make(map[string]*EnrichmentResult),
		cacheTTL: time.Duration(cfg.CacheTTL) * time.Second,
	}
}

// EnrichKeywords enrichit une liste de mots-clés avec l'IA
func (ai *AIEnricher) EnrichKeywords(ctx context.Context, keywords []string, content string) (*EnrichmentResult, error) {
	if !ai.config.Enabled {
		return ai.mockEnrichment(keywords), nil
	}

	// Vérifier le cache
	cacheKey := ai.getCacheKey(keywords, content)
	if cached, exists := ai.cache[cacheKey]; exists {
		if time.Since(cached.CachedAt) < ai.cacheTTL {
			return cached, nil
		}
		// Supprimer l'entrée expirée
		delete(ai.cache, cacheKey)
	}

	// Mode mock pour les tests
	if ai.config.MockMode {
		result := ai.mockEnrichment(keywords)
		ai.cache[cacheKey] = result
		return result, nil
	}

	// Appel réel à OpenAI
	result, err := ai.callOpenAI(ctx, keywords, content)
	if err != nil {
		// En cas d'erreur, retourner un mock
		return ai.mockEnrichment(keywords), fmt.Errorf("OpenAI enrichment failed, using mock: %w", err)
	}

	// Mettre en cache
	result.CachedAt = time.Now()
	ai.cache[cacheKey] = result

	return result, nil
}

// callOpenAI fait l'appel réel à l'API OpenAI
func (ai *AIEnricher) callOpenAI(ctx context.Context, keywords []string, content string) (*EnrichmentResult, error) {
	if ai.config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	prompt := ai.buildPrompt(keywords, content)
	
	request := OpenAIRequest{
		Model:       ai.config.Model,
		MaxTokens:   ai.config.MaxTokens,
		Temperature: 0.3, // Faible température pour des résultats consistants
		Messages: []Message{
			{
				Role:    "system",
				Content: "Tu es un expert SEO spécialisé dans l'analyse sémantique et l'optimisation de contenu. Réponds uniquement en JSON valide.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.config.APIKey)

	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status %d", resp.StatusCode)
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parser la réponse JSON
	var result EnrichmentResult
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("error parsing OpenAI response: %w", err)
	}

	return &result, nil
}

// buildPrompt construit le prompt pour OpenAI
func (ai *AIEnricher) buildPrompt(keywords []string, content string) string {
	keywordsList := strings.Join(keywords, ", ")
	contentPreview := content
	if len(content) > 500 {
		contentPreview = content[:500] + "..."
	}

	return fmt.Sprintf(`Analyse ces mots-clés SEO et ce contenu pour fournir un enrichissement sémantique avancé.

MOTS-CLÉS À ANALYSER: %s

CONTENU (extrait): %s

Réponds uniquement avec un JSON valide suivant cette structure exacte:
{
  "keywords": [
    {
      "text": "mot-clé",
      "intent": "informational|navigational|transactional|commercial",
      "difficulty": 1-100,
      "search_volume": "low|medium|high|very_high",
      "related_keywords": ["mot1", "mot2", "mot3"],
      "content_suggestions": ["suggestion1", "suggestion2"]
    }
  ],
  "content_questions": [
    "Question que se posent les utilisateurs",
    "Autre question pertinente"
  ],
  "seo_recommendations": [
    "Recommandation SEO spécifique",
    "Action d'optimisation"
  ],
  "search_intent": "intent principal du contenu",
  "competitiveness_score": 1-100
}

Concentre-toi sur la pertinence SEO et les opportunités d'optimisation concrètes.`, keywordsList, contentPreview)
}

// mockEnrichment retourne des données fictives pour les tests
func (ai *AIEnricher) mockEnrichment(keywords []string) *EnrichmentResult {
	enrichedKeywords := make([]EnrichedKeyword, len(keywords))
	for i, kw := range keywords {
		enrichedKeywords[i] = EnrichedKeyword{
			Text:       kw,
			Intent:     "informational",
			Difficulty: 30 + (i * 10), // Difficulté croissante
			SearchVolume: "medium",
			RelatedKeywords: []string{
				kw + " guide",
				kw + " tips",
				"best " + kw,
			},
			ContentSuggestions: []string{
				"Créer un guide complet sur " + kw,
				"Ajouter des exemples pratiques",
			},
		}
	}

	return &EnrichmentResult{
		Keywords: enrichedKeywords,
		ContentQuestions: []string{
			"Comment optimiser " + strings.Join(keywords, " et ") + " ?",
			"Quelles sont les meilleures pratiques ?",
			"Quels outils utiliser ?",
		},
		SEORecommendations: []string{
			"Optimiser les balises title et meta description",
			"Améliorer la structure des titres (H1-H6)",
			"Ajouter du contenu de qualité sur les sujets connexes",
			"Créer des liens internes pertinents",
		},
		SearchIntent:        "informational",
		CompetitivenessScore: 65,
		CachedAt:           time.Now(),
	}
}

// getCacheKey génère une clé de cache pour les mots-clés et contenu
func (ai *AIEnricher) getCacheKey(keywords []string, content string) string {
	keywordStr := strings.Join(keywords, "|")
	contentHash := fmt.Sprintf("%x", len(content)) // Simple hash basé sur la longueur
	return fmt.Sprintf("%s_%s", keywordStr, contentHash)
}

// GetCacheStats retourne les statistiques de cache
func (ai *AIEnricher) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_entries": len(ai.cache),
		"enabled":       ai.config.Enabled,
		"mock_mode":     ai.config.MockMode,
		"model":         ai.config.Model,
		"cache_ttl":     ai.cacheTTL.String(),
	}
}

// ClearExpiredCache nettoie les entrées de cache expirées
func (ai *AIEnricher) ClearExpiredCache() int {
	cleaned := 0
	now := time.Now()
	
	for key, entry := range ai.cache {
		if now.Sub(entry.CachedAt) > ai.cacheTTL {
			delete(ai.cache, key)
			cleaned++
		}
	}
	
	return cleaned
}