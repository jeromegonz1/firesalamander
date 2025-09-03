package page_profiler

import (
	"context"
	"testing"

	"firesalamander/internal/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageProfiler_ImplementsAgentInterface(t *testing.T) {
	// [RÔLE: QA] Vérifier implémentation interface
	var _ agents.Agent = (*PageProfiler)(nil)
}

func TestNewPageProfiler(t *testing.T) {
	profiler := NewPageProfiler()
	assert.NotNil(t, profiler)
	assert.Equal(t, "page_profiler", profiler.Name())
}

func TestPageProfiler_Name(t *testing.T) {
	profiler := NewPageProfiler()
	assert.Equal(t, "page_profiler", profiler.Name())
}

func TestPageProfiler_HealthCheck(t *testing.T) {
	profiler := NewPageProfiler()
	err := profiler.HealthCheck()
	assert.NoError(t, err)
}

func TestPageProfiler_Process_ValidHTML(t *testing.T) {
	// [RÔLE: QA] Test avec HTML réel
	profiler := NewPageProfiler()
	ctx := context.Background()
	
	htmlContent := `
<!DOCTYPE html>
<html lang="fr">
<head>
	<title>Test Page</title>
	<meta name="description" content="Une page de test">
	<meta name="keywords" content="test, page, html">
	<meta property="og:title" content="Test Page">
</head>
<body>
	<h1>Titre Principal</h1>
	<h2>Sous-titre</h2>
	<p>Contenu avec <strong>texte important</strong>.</p>
	<img src="image.jpg" alt="Image descriptive" width="300" height="200">
	<a href="https://example.com">Lien externe</a>
	<div itemscope itemtype="http://schema.org/Article">
		<h3 itemprop="headline">Article Schema</h3>
	</div>
</body>
</html>`

	input := PageRequest{
		HTML: htmlContent,
		URL:  "https://example.com/test",
	}

	result, err := profiler.Process(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "page_profiler", result.AgentName)
	assert.Equal(t, "completed", result.Status)
	
	// Vérifier structure de données
	data := result.Data
	assert.Contains(t, data, "meta_tags")
	assert.Contains(t, data, "headings")
	assert.Contains(t, data, "images")
	assert.Contains(t, data, "links")
	assert.Contains(t, data, "schema_markup")
	assert.Contains(t, data, "content_stats")
}

func TestPageProfiler_ExtractMetaTags(t *testing.T) {
	// [RÔLE: QA] Test extraction meta tags
	profiler := NewPageProfiler()
	
	html := `<head>
		<title>Page Title</title>
		<meta name="description" content="Page description">
		<meta name="keywords" content="mot1, mot2">
		<meta property="og:title" content="OG Title">
		<meta name="twitter:card" content="summary">
	</head>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	assert.Equal(t, "Page Title", profile.MetaTags["title"])
	assert.Equal(t, "Page description", profile.MetaTags["description"])
	assert.Equal(t, "mot1, mot2", profile.MetaTags["keywords"])
	assert.Equal(t, "OG Title", profile.MetaTags["og:title"])
	assert.Equal(t, "summary", profile.MetaTags["twitter:card"])
}

func TestPageProfiler_ExtractHeadings(t *testing.T) {
	// [RÔLE: QA] Test hiérarchie H1-H6
	profiler := NewPageProfiler()
	
	html := `<body>
		<h1>Heading 1</h1>
		<h2>Heading 2a</h2>
		<h2>Heading 2b</h2>
		<h3>Heading 3</h3>
		<h4>Heading 4</h4>
		<h5>Heading 5</h5>
		<h6>Heading 6</h6>
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	assert.Len(t, profile.Headings.H1, 1)
	assert.Len(t, profile.Headings.H2, 2)
	assert.Len(t, profile.Headings.H3, 1)
	assert.Equal(t, "Heading 1", profile.Headings.H1[0])
	assert.Contains(t, profile.Headings.H2, "Heading 2a")
	assert.Contains(t, profile.Headings.H2, "Heading 2b")
}

func TestPageProfiler_ExtractImages(t *testing.T) {
	// [RÔLE: QA] Test analyse images
	profiler := NewPageProfiler()
	
	html := `<body>
		<img src="image1.jpg" alt="Description 1" width="300" height="200">
		<img src="image2.png" alt="" width="150">
		<img src="image3.gif">
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	assert.Len(t, profile.Images, 3)
	
	img1 := profile.Images[0]
	assert.Equal(t, "image1.jpg", img1.Src)
	assert.Equal(t, "Description 1", img1.Alt)
	assert.Equal(t, 300, img1.Width)
	assert.Equal(t, 200, img1.Height)
	assert.Equal(t, "jpg", img1.Format)
	
	// Image sans alt
	img2 := profile.Images[1]
	assert.Equal(t, "", img2.Alt)
}

func TestPageProfiler_ExtractLinks(t *testing.T) {
	// [RÔLE: QA] Test extraction liens
	profiler := NewPageProfiler()
	
	html := `<body>
		<a href="https://external.com">Lien externe</a>
		<a href="/internal">Lien interne</a>
		<a href="mailto:test@example.com">Email</a>
		<a href="#anchor">Ancre</a>
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	assert.Len(t, profile.Links, 4)
	
	// Vérifier classification des liens
	hasExternal := false
	hasInternal := false
	hasEmail := false
	hasAnchor := false
	
	for _, link := range profile.Links {
		switch link.Type {
		case "external":
			hasExternal = true
		case "internal":
			hasInternal = true
		case "email":
			hasEmail = true
		case "anchor":
			hasAnchor = true
		}
	}
	
	assert.True(t, hasExternal)
	assert.True(t, hasInternal)
	assert.True(t, hasEmail)
	assert.True(t, hasAnchor)
}

func TestPageProfiler_ExtractSchemaMarkup(t *testing.T) {
	// [RÔLE: QA] Test Schema.org markup
	profiler := NewPageProfiler()
	
	html := `<body>
		<div itemscope itemtype="http://schema.org/Article">
			<h1 itemprop="headline">Article Title</h1>
			<span itemprop="author">John Doe</span>
		</div>
		<script type="application/ld+json">
		{
			"@context": "http://schema.org",
			"@type": "Organization",
			"name": "Example Org"
		}
		</script>
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	assert.Len(t, profile.SchemaMarkup.Microdata, 1)
	assert.Len(t, profile.SchemaMarkup.JsonLD, 1)
	
	microdata := profile.SchemaMarkup.Microdata[0]
	assert.Equal(t, "http://schema.org/Article", microdata.Type)
	assert.Contains(t, microdata.Properties, "headline")
	assert.Contains(t, microdata.Properties, "author")
}

func TestPageProfiler_ContentStats(t *testing.T) {
	// [RÔLE: QA] Test statistiques contenu
	profiler := NewPageProfiler()
	
	html := `<body>
		<h1>Titre</h1>
		<p>Premier paragraphe avec <strong>texte important</strong>.</p>
		<p>Deuxième paragraphe plus long avec plusieurs mots et phrases.</p>
		<ul>
			<li>Premier item</li>
			<li>Deuxième item</li>
		</ul>
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	stats := profile.ContentStats
	assert.Greater(t, stats.WordCount, 10)
	assert.Greater(t, stats.CharacterCount, 50)
	assert.Greater(t, stats.ParagraphCount, 1)
	assert.Equal(t, 1, stats.ListCount)
	assert.Greater(t, stats.TextDensity, 0.0)
}

func TestPageProfiler_Process_InvalidInput(t *testing.T) {
	// [RÔLE: QA] Test gestion erreurs
	profiler := NewPageProfiler()
	ctx := context.Background()
	
	// Input invalide
	_, err := profiler.Process(ctx, "invalid input")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input type")
}

func TestPageProfiler_Process_EmptyHTML(t *testing.T) {
	// [RÔLE: QA] Test HTML vide
	profiler := NewPageProfiler()
	ctx := context.Background()
	
	input := PageRequest{
		HTML: "",
		URL:  "https://example.com",
	}
	
	result, err := profiler.Process(ctx, input)
	require.NoError(t, err)
	assert.NotNil(t, result)
	
	// Doit retourner structure vide mais valide
	data := result.Data
	assert.Contains(t, data, "content_stats")
	stats := data["content_stats"].(ContentStats)
	assert.Equal(t, 0, stats.WordCount)
}

func TestPageProfiler_CoreWebVitalsHints(t *testing.T) {
	// [RÔLE: QA] Test indices Core Web Vitals
	profiler := NewPageProfiler()
	
	html := `<head>
		<style>
			.large { font-size: 48px; }
		</style>
		<script async>console.log("async script");</script>
	</head>
	<body>
		<img src="large.jpg" width="2000" height="1500">
		<div style="width:100%;height:500px;background:red;"></div>
	</body>`
	
	profile := profiler.analyzePage(html, "https://example.com")
	
	vitals := profile.CoreWebVitalsHints
	assert.Contains(t, vitals.LargestContentfulPaint, "large image detected")
	assert.Contains(t, vitals.CumulativeLayoutShift, "inline styles detected")
	assert.Contains(t, vitals.FirstInputDelay, "async scripts found")
}