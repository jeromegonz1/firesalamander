package seo

import (
	"testing"
)

// Test de l'analyseur SEO réel avec HTML - Title Analysis
func TestRealSEOAnalyzer_AnalyzeTitle(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()

	tests := []struct {
		name           string
		htmlContent    string
		expectedLength int
		expectedIssues int
		shouldHaveContent bool
	}{
		{
			name:           "Valid title",
			htmlContent:    "<html><head><title>Great SEO Title Example</title></head></html>",
			expectedLength: 23,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
		{
			name:           "Missing title",
			htmlContent:    "<html><head></head></html>",
			expectedLength: 0,
			expectedIssues: 1,
			shouldHaveContent: false,
		},
		{
			name:           "Empty title",
			htmlContent:    "<html><head><title></title></head></html>",
			expectedLength: 0,
			expectedIssues: 1,
			shouldHaveContent: false,
		},
		{
			name:           "Title too short",
			htmlContent:    "<html><head><title>Short</title></head></html>",
			expectedLength: 5,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
		{
			name:           "Title too long",
			htmlContent:    "<html><head><title>This is a very long title that exceeds the recommended sixty character limit for SEO</title></head></html>",
			expectedLength: 84,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeTitle(tt.htmlContent)
			
			// Test length
			if result.Length != tt.expectedLength {
				t.Errorf("Expected length %d, got %d", tt.expectedLength, result.Length)
			}
			
			// Test issues count
			if len(result.Issues) != tt.expectedIssues {
				t.Errorf("Expected %d issues, got %d: %v", tt.expectedIssues, len(result.Issues), result.Issues)
			}

			// Test content presence
			if tt.shouldHaveContent && result.Content == "" {
				t.Error("Expected content but got empty string")
			}
			if !tt.shouldHaveContent && result.Content != "" {
				t.Errorf("Expected empty content but got: %s", result.Content)
			}

			// Test score calculation
			if result.Score < 0 || result.Score > 100 {
				t.Errorf("Score should be between 0 and 100, got %d", result.Score)
			}
		})
	}
}

// Test de l'analyseur SEO - Meta Description
func TestRealSEOAnalyzer_AnalyzeMetaDescription(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()

	tests := []struct {
		name           string
		htmlContent    string
		expectedLength int
		expectedIssues int
		shouldHaveContent bool
	}{
		{
			name:           "Valid meta description",
			htmlContent:    `<html><head><meta name="description" content="Great meta description for SEO purposes, optimized for search engines with perfect length."></head></html>`,
			expectedLength: 90,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
		{
			name:           "Missing meta description", 
			htmlContent:    "<html><head></head></html>",
			expectedLength: 0,
			expectedIssues: 1,
			shouldHaveContent: false,
		},
		{
			name:           "Meta description too short",
			htmlContent:    `<html><head><meta name="description" content="Too short"></head></html>`,
			expectedLength: 9,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
		{
			name:           "Meta description too long",
			htmlContent:    `<html><head><meta name="description" content="This is an extremely long meta description that definitely exceeds the recommended maximum of 160 characters for search engine optimization and user experience purposes."></head></html>`,
			expectedLength: 169,
			expectedIssues: 1,
			shouldHaveContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeMetaDescription(tt.htmlContent)
			
			// Test length
			if result.Length != tt.expectedLength {
				t.Errorf("Expected length %d, got %d", tt.expectedLength, result.Length)
			}
			
			// Test issues count
			if len(result.Issues) != tt.expectedIssues {
				t.Errorf("Expected %d issues, got %d: %v", tt.expectedIssues, len(result.Issues), result.Issues)
			}

			// Test content presence
			if tt.shouldHaveContent && result.Content == "" {
				t.Error("Expected content but got empty string")
			}
			if !tt.shouldHaveContent && result.Content != "" {
				t.Errorf("Expected empty content but got: %s", result.Content)
			}

			// Test score
			if result.Score < 0 || result.Score > 100 {
				t.Errorf("Score should be between 0 and 100, got %d", result.Score)
			}
		})
	}
}

// Test de l'analyseur SEO - Headings Analysis  
func TestRealSEOAnalyzer_AnalyzeHeadings(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()

	tests := []struct {
		name           string
		htmlContent    string
		expectedH1Count int
		expectedIssues int
	}{
		{
			name:           "Valid single H1",
			htmlContent:    "<html><body><h1>Main Title</h1><h2>Subtitle</h2></body></html>",
			expectedH1Count: 1,
			expectedIssues: 0,
		},
		{
			name:           "Missing H1",
			htmlContent:    "<html><body><h2>Only H2</h2><h3>And H3</h3></body></html>",
			expectedH1Count: 0,
			expectedIssues: 1,
		},
		{
			name:           "Multiple H1",
			htmlContent:    "<html><body><h1>First</h1><h1>Second</h1></body></html>",
			expectedH1Count: 2,
			expectedIssues: 1,
		},
		{
			name:           "No headings",
			htmlContent:    "<html><body><p>Just text</p></body></html>",
			expectedH1Count: 0,
			expectedIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeHeadings(tt.htmlContent)
			
			// Test H1 count
			if result.H1Count != tt.expectedH1Count {
				t.Errorf("Expected %d H1, got %d", tt.expectedH1Count, result.H1Count)
			}
			
			// Test issues count
			if len(result.Issues) != tt.expectedIssues {
				t.Errorf("Expected %d issues, got %d: %v", tt.expectedIssues, len(result.Issues), result.Issues)
			}

			// Test score
			if result.Score < 0 || result.Score > 100 {
				t.Errorf("Score should be between 0 and 100, got %d", result.Score)
			}
		})
	}
}

// Test de l'analyseur SEO - Images Analysis
func TestRealSEOAnalyzer_AnalyzeImages(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()

	tests := []struct {
		name              string
		htmlContent       string
		expectedImages    int
		expectedWithAlt   int
		expectedCoverage  float64
		expectedIssues    int
	}{
		{
			name:              "All images with alt",
			htmlContent:       `<html><body><img src="1.jpg" alt="Image 1"><img src="2.jpg" alt="Image 2"></body></html>`,
			expectedImages:    2,
			expectedWithAlt:   2,
			expectedCoverage:  1.0,
			expectedIssues:    0,
		},
		{
			name:              "Some images without alt",
			htmlContent:       `<html><body><img src="1.jpg" alt="Image 1"><img src="2.jpg"><img src="3.jpg" alt="Image 3"></body></html>`,
			expectedImages:    3,
			expectedWithAlt:   2,
			expectedCoverage:  0.6667,
			expectedIssues:    1,
		},
		{
			name:              "No images",
			htmlContent:       "<html><body><p>No images here</p></body></html>",
			expectedImages:    0,
			expectedWithAlt:   0,
			expectedCoverage:  0.0,
			expectedIssues:    0,
		},
		{
			name:              "All images without alt", 
			htmlContent:       `<html><body><img src="1.jpg"><img src="2.jpg"></body></html>`,
			expectedImages:    2,
			expectedWithAlt:   0,
			expectedCoverage:  0.0,
			expectedIssues:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeImages(tt.htmlContent)
			
			// Test images count
			if result.TotalImages != tt.expectedImages {
				t.Errorf("Expected %d images, got %d", tt.expectedImages, result.TotalImages)
			}
			
			// Test images with alt
			if result.ImagesWithAlt != tt.expectedWithAlt {
				t.Errorf("Expected %d images with alt, got %d", tt.expectedWithAlt, result.ImagesWithAlt)
			}

			// Test coverage (with tolerance for float comparison)
			tolerance := 0.01
			if tt.expectedImages > 0 {
				if abs(result.AltTextCoverage - tt.expectedCoverage) > tolerance {
					t.Errorf("Expected coverage %.4f, got %.4f", tt.expectedCoverage, result.AltTextCoverage)
				}
			}
			
			// Test issues count
			if len(result.Issues) != tt.expectedIssues {
				t.Errorf("Expected %d issues, got %d: %v", tt.expectedIssues, len(result.Issues), result.Issues)
			}
		})
	}
}

// Helper function for float comparison
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// Test d'intégration - Analyse complète d'une page
func TestRealSEOAnalyzer_FullPageAnalysis(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()
	
	// HTML complet avec divers éléments SEO
	htmlContent := `
<!DOCTYPE html>
<html lang="fr">
<head>
    <title>Page de Test SEO - Titre Optimisé</title>
    <meta name="description" content="Description métaoptimisée pour le SEO avec une longueur parfaite et des mots-clés pertinents pour les moteurs de recherche.">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="canonical" href="https://example.com/test">
</head>
<body>
    <h1>Titre Principal H1</h1>
    <h2>Sous-titre H2</h2>
    <h3>Section H3</h3>
    
    <img src="image1.jpg" alt="Description de l'image 1">
    <img src="image2.jpg" alt="Description de l'image 2">
    <img src="image3.jpg">
    
    <p>Contenu de la page avec du texte optimisé pour le référencement naturel.</p>
</body>
</html>`

	// Test title
	titleResult := analyzer.AnalyzeTitle(htmlContent)
	if titleResult.Length == 0 {
		t.Error("Should find title in complete HTML")
	}
	if len(titleResult.Issues) > 0 {
		t.Errorf("Valid title should have no issues, got: %v", titleResult.Issues)
	}

	// Test meta description
	metaResult := analyzer.AnalyzeMetaDescription(htmlContent)
	if metaResult.Length == 0 {
		t.Error("Should find meta description in complete HTML")
	}

	// Test headings
	headingsResult := analyzer.AnalyzeHeadings(htmlContent)
	if headingsResult.H1Count != 1 {
		t.Errorf("Expected 1 H1, got %d", headingsResult.H1Count)
	}

	// Test images
	imagesResult := analyzer.AnalyzeImages(htmlContent)
	if imagesResult.TotalImages != 3 {
		t.Errorf("Expected 3 images, got %d", imagesResult.TotalImages)
	}
	if imagesResult.ImagesWithAlt != 2 {
		t.Errorf("Expected 2 images with alt, got %d", imagesResult.ImagesWithAlt)
	}
}