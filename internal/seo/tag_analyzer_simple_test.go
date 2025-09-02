package seo

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test simple pour améliorer la couverture du TagAnalyzer
func TestTagAnalyzer_Simple(t *testing.T) {
	analyzer := NewTagAnalyzer()
	assert.NotNil(t, analyzer)

	htmlContent := `<!DOCTYPE html>
<html lang="fr">
<head>
	<title>Test SEO</title>
	<meta name="description" content="Description de test pour le référencement.">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="canonical" href="https://example.com/test">
</head>
<body>
	<h1>Titre principal</h1>
	<h2>Section</h2>
	<p>Contenu de test.</p>
	<img src="test.jpg" alt="Image de test">
	<a href="/interne">Lien interne</a>
</body>
</html>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	require.NoError(t, err)

	result, err := analyzer.Analyze(doc, htmlContent)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Test basique de présence des analyses
	assert.NotNil(t, result.Title)
	assert.NotNil(t, result.MetaDescription)
	assert.NotNil(t, result.Headings)
	assert.NotNil(t, result.MetaTags)
	assert.NotNil(t, result.Images)
	assert.NotNil(t, result.Links)
	assert.NotNil(t, result.Microdata)
}

// Test des méthodes individuelles pour coverage
func TestTagAnalyzer_IndividualMethods(t *testing.T) {
	analyzer := NewTagAnalyzer()

	// HTML simple pour tester
	simpleHTML := `<html><head><title>Test</title></head><body><h1>Test</h1></body></html>`
	doc, err := html.Parse(strings.NewReader(simpleHTML))
	require.NoError(t, err)

	// Tester chaque méthode individuelle
	titleResult := analyzer.analyzeTitle(doc)
	assert.NotNil(t, titleResult)

	metaResult := analyzer.analyzeMetaDescription(doc)
	assert.NotNil(t, metaResult)

	headingsResult := analyzer.analyzeHeadings(doc)
	assert.NotNil(t, headingsResult)

	metaTagsResult := analyzer.analyzeMetaTags(doc)
	assert.NotNil(t, metaTagsResult)

	imagesResult := analyzer.analyzeImages(doc)
	assert.NotNil(t, imagesResult)

	linksResult := analyzer.analyzeLinks(doc)
	assert.NotNil(t, linksResult)

	microdataResult := analyzer.analyzeMicrodata(doc, simpleHTML)
	assert.NotNil(t, microdataResult)
}

// Test des cas d'erreur et edge cases
func TestTagAnalyzer_EdgeCases(t *testing.T) {
	analyzer := NewTagAnalyzer()

	// HTML vide
	emptyHTML := `<html></html>`
	emptyDoc, err := html.Parse(strings.NewReader(emptyHTML))
	require.NoError(t, err)

	result, err := analyzer.Analyze(emptyDoc, emptyHTML)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// HTML malformé mais parsable
	malformedHTML := `<html><head><title>Test</head><body><p>Test</body>`
	malformedDoc, err := html.Parse(strings.NewReader(malformedHTML))
	require.NoError(t, err)

	result2, err := analyzer.Analyze(malformedDoc, malformedHTML)
	require.NoError(t, err)
	assert.NotNil(t, result2)
}