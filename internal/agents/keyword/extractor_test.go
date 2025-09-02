package keyword

import (
	"context"
	"strings"
	"testing"

	"firesalamander/internal/constants"
)

func TestKeywordExtractor_Name(t *testing.T) {
	extractor := NewKeywordExtractor()
	
	expected := constants.AgentNameKeyword
	if extractor.Name() != expected {
		t.Errorf("Expected name %s, got %s", expected, extractor.Name())
	}
}

func TestKeywordExtractor_ExtractKeywords(t *testing.T) {
	extractor := NewKeywordExtractor()

	tests := []struct {
		name           string
		content        string
		expectedCount  int
		shouldContain  string
	}{
		{
			name:          "empty content",
			content:       "",
			expectedCount: 0,
			shouldContain: "",
		},
		{
			name:          "simple content",
			content:       "optimisation SEO pour améliorer le référencement naturel",
			expectedCount: 4,
			shouldContain: "optimisation",
		},
		{
			name:          "HTML content",
			content:       "<h1>Titre SEO</h1><p>Contenu avec <strong>mots-clés</strong> importants</p>",
			expectedCount: 4,
			shouldContain: "titre",
		},
		{
			name:          "stop words filtered",
			content:       "le contenu avec des mots très importants pour le référencement",
			expectedCount: 4,
			shouldContain: "contenu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractor.ExtractKeywords(tt.content)
			
			if err != nil {
				t.Errorf("ExtractKeywords returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("ExtractKeywords returned nil result")
				return
			}
			
			if tt.expectedCount > 0 && len(result.Keywords) == 0 {
				t.Errorf("Expected at least %d keywords, got 0", tt.expectedCount)
				return
			}
			
			if tt.shouldContain != "" {
				found := false
				for _, keyword := range result.Keywords {
					if keyword.Term == tt.shouldContain {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find keyword '%s' in results", tt.shouldContain)
				}
			}
			
			// Vérifier que les résultats sont triés par pertinence
			if len(result.Keywords) > 1 {
				for i := 1; i < len(result.Keywords); i++ {
					if result.Keywords[i-1].Relevance < result.Keywords[i].Relevance {
						t.Error("Keywords should be sorted by relevance (descending)")
						break
					}
				}
			}
		})
	}
}

func TestKeywordExtractor_AnalyzeDensity(t *testing.T) {
	extractor := NewKeywordExtractor()
	
	tests := []struct {
		name             string
		keywords         []string
		content          string
		expectedMetrics  int
		shouldHaveRecommendations bool
	}{
		{
			name:             "empty content",
			keywords:         []string{"test"},
			content:          "",
			expectedMetrics:  1,
			shouldHaveRecommendations: true,
		},
		{
			name:             "normal density",
			keywords:         []string{"SEO", "référencement"},
			content:          "Le SEO et le référencement sont importants pour la visibilité. Le SEO technique aide au référencement naturel.",
			expectedMetrics:  2,
			shouldHaveRecommendations: true,
		},
		{
			name:             "high density",
			keywords:         []string{"test"},
			content:          "test test test test test test test test test test",
			expectedMetrics:  1,
			shouldHaveRecommendations: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := extractor.AnalyzeDensity(tt.keywords, tt.content)
			
			if err != nil {
				t.Errorf("AnalyzeDensity returned error: %v", err)
				return
			}
			
			if report == nil {
				t.Error("AnalyzeDensity returned nil report")
				return
			}
			
			if len(report.KeywordMetrics) != tt.expectedMetrics {
				t.Errorf("Expected %d keyword metrics, got %d", tt.expectedMetrics, len(report.KeywordMetrics))
			}
			
			if tt.shouldHaveRecommendations && len(report.Recommendations) == 0 {
				t.Error("Expected recommendations but got none")
			}
			
			// Vérifier que les densités sont dans une plage raisonnable
			for keyword, density := range report.KeywordMetrics {
				if density < 0 || density > 100 {
					t.Errorf("Invalid density for keyword '%s': %f", keyword, density)
				}
			}
		})
	}
}

func TestKeywordExtractor_Process(t *testing.T) {
	extractor := NewKeywordExtractor()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         interface{}
		expectError   bool
		expectedStatus string
	}{
		{
			name:          "valid string input",
			input:         "contenu de test pour l'extraction de mots-clés",
			expectError:   false,
			expectedStatus: constants.StatusCompleted,
		},
		{
			name:          "invalid input type",
			input:         123,
			expectError:   false,
			expectedStatus: constants.StatusFailed,
		},
		{
			name:          "empty string",
			input:         "",
			expectError:   false,
			expectedStatus: constants.StatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractor.Process(ctx, tt.input)
			
			if err != nil {
				t.Errorf("Process returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("Process returned nil result")
				return
			}
			
			if result.AgentName != constants.AgentNameKeyword {
				t.Errorf("Expected agent name %s, got %s", constants.AgentNameKeyword, result.AgentName)
			}
			
			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
			
			if result.Duration < 0 {
				t.Error("Expected non-negative duration")
			}
			
			if tt.expectedStatus == constants.StatusCompleted && result.Data == nil {
				t.Error("Expected data in successful result")
			}
			
			if tt.expectedStatus == constants.StatusFailed && len(result.Errors) == 0 {
				t.Error("Expected errors in failed result")
			}
		})
	}
}

func TestKeywordExtractor_HealthCheck(t *testing.T) {
	extractor := NewKeywordExtractor()
	
	err := extractor.HealthCheck()
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestKeywordExtractor_CleanContent(t *testing.T) {
	extractor := NewKeywordExtractor()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML tags removal",
			input:    "<h1>Titre</h1><p>Contenu</p>",
			expected: "titre contenu",
		},
		{
			name:     "special characters removal",
			input:    "Mot-clé avec @caractères #spéciaux!",
			expected: "mot clé avec caractères spéciaux",
		},
		{
			name:     "multiple spaces normalization",
			input:    "Texte   avec    beaucoup     d'espaces",
			expected: "texte avec beaucoup d espaces",
		},
		{
			name:     "mixed content",
			input:    "<div>Test &amp; @validation #123</div>",
			expected: "test validation 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.cleanContent(tt.input)
			// Convertir en minuscules pour la comparaison
			result = strings.ToLower(result)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestKeywordExtractor_ExtractWords(t *testing.T) {
	extractor := NewKeywordExtractor()
	
	tests := []struct {
		name     string
		input    string
		minCount int
	}{
		{
			name:     "normal content",
			input:    "contenu normal avec mots valides",
			minCount: 4,
		},
		{
			name:     "short words filtered",
			input:    "a de le la test contenu",
			minCount: 2,
		},
		{
			name:     "stop words filtered",
			input:    "le contenu et les mots pour avec",
			minCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := extractor.extractWords(tt.input)
			
			if len(words) < tt.minCount {
				t.Errorf("Expected at least %d words, got %d", tt.minCount, len(words))
			}
			
			// Vérifier qu'aucun mot stop n'est présent
			for _, word := range words {
				if extractor.stopWords[word] {
					t.Errorf("Stop word '%s' should be filtered out", word)
				}
				
				if len(word) < extractor.minWordLength {
					t.Errorf("Word '%s' is shorter than minimum length %d", word, extractor.minWordLength)
				}
			}
		})
	}
}