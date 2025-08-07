package semantic

import (
	"context"
	"testing"
)

// Test basique de l'extracteur de contenu
func TestContentExtractorBasic(t *testing.T) {
	extractor := NewContentExtractor()

	html := `<html><head><title>Test</title></head><body><p>Hello world</p></body></html>`
	
	content, err := extractor.Extract(html)
	if err != nil {
		t.Fatalf("Erreur extraction: %v", err)
	}

	if content.Title != "Test" {
		t.Errorf("Titre attendu 'Test', reçu '%s'", content.Title)
	}

	if content.WordCount == 0 {
		t.Error("Le nombre de mots devrait être > 0")
	}
}

// Test basique de l'analyseur n-grammes
func TestNGramAnalyzerBasic(t *testing.T) {
	analyzer := NewNGramAnalyzer()
	analyzer.minFrequency = 1 // Réduire le seuil minimum pour les tests

	text := "hello world test analysis hello world"
	ngrams := analyzer.Analyze(text)

	if len(ngrams) == 0 {
		t.Error("Des n-grammes devraient être générés")
	}

	if len(ngrams[1]) == 0 {
		t.Error("Des unigrammes devraient être générés")
	}

	t.Logf("Generated ngrams: 1-gram=%d, 2-gram=%d, 3-gram=%d", 
		len(ngrams[1]), len(ngrams[2]), len(ngrams[3]))
}

// Test basique du scorer SEO
func TestSEOScorerBasic(t *testing.T) {
	scorer := NewSEOScorer()

	content := &ExtractedContent{
		Title:           "Test SEO Page",
		MetaDescription: "Description de test pour vérifier le scoring SEO automatique",
		WordCount:       500,
		Headings:        []string{"Titre 1", "Titre 2"},
	}

	analysis := LocalAnalysis{
		Keywords: []Keyword{
			{Term: "seo", Frequency: 5, Density: 1.0, Relevance: 0.8},
		},
		ReadabilityScore: 75.0,
	}

	score := scorer.Score(content, analysis, nil)

	if score.Overall < 0 || score.Overall > 100 {
		t.Errorf("Score SEO invalide: %f", score.Overall)
	}

	if len(score.Factors) == 0 {
		t.Error("Des facteurs SEO devraient être évalués")
	}
}

// Test du module complet avec l'analyseur de test
func TestCompleteSemanticAnalysis(t *testing.T) {
	analyzer := NewTestSemanticAnalyzer()

	html := `
	<!DOCTYPE html>
	<html lang="fr">
	<head>
		<title>Guide SEO Fire Salamander</title>
		<meta name="description" content="Découvrez notre guide complet pour optimiser votre SEO avec Fire Salamander">
	</head>
	<body>
		<h1>Guide SEO Fire Salamander</h1>
		<p>Le référencement naturel est essentiel pour améliorer la visibilité de votre site web.</p>
		<h2>Techniques d'optimisation</h2>
		<p>Utilisez les meilleures pratiques SEO pour optimiser votre contenu.</p>
	</body>
	</html>
	`

	ctx := context.Background()
	result, err := analyzer.AnalyzePage(ctx, "https://example.com/guide-seo", html)

	if err != nil {
		t.Fatalf("Erreur analyse: %v", err)
	}

	// Vérifications de base
	if result.Title != "Guide SEO Fire Salamander" {
		t.Errorf("Titre incorrect: %s", result.Title)
	}

	if result.WordCount == 0 {
		t.Error("Le nombre de mots devrait être > 0")
	}

	if len(result.LocalAnalysis.Keywords) == 0 {
		t.Error("Des mots-clés devraient être identifiés")
	}

	if result.SEOScore.Overall == 0 {
		t.Error("Le score SEO devrait être > 0")
	}

	t.Logf("Analyse terminée - Titre: %s, Mots: %d, Score SEO: %.1f", 
		result.Title, result.WordCount, result.SEOScore.Overall)
}