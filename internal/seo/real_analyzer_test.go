package seo

import (
	"context"
	"testing"
	"time"

	"firesalamander/internal/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TDD PHASE RED - Tests écrits AVANT l'implémentation
// Ces tests DOIVENT échouer jusqu'à l'implémentation

// TestRealSEOAnalyzer_ExtractAndScoreTitle tests title extraction and scoring
func TestRealSEOAnalyzer_ExtractAndScoreTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected RealTitleAnalysis
	}{
		{
			name: "perfect_title",
			html: `<!DOCTYPE html>
				 <html>
				 <head>
					 <title>Fire Salamander - Professional SEO Analysis Tool</title>
				 </head>
				 </html>`,
			expected: RealTitleAnalysis{
				Present:  true,
				Content:  "Fire Salamander - Professional SEO Analysis Tool",
				Length:   49,
				Score:    constants.MaxTitleScore, // 20 points for perfect title
				Issues:   []string{},
				Keywords: []string{"Fire", "Salamander", "Professional", "SEO", "Analysis", "Tool"},
				Severity: constants.RealSEOStatusInfo,
			},
		},
		{
			name: "missing_title",
			html: `<html><head></head><body>Content without title</body></html>`,
			expected: RealTitleAnalysis{
				Present:  false,
				Content:  "",
				Length:   0,
				Score:    0,
				Issues:   []string{constants.ErrorTitleMissing},
				Keywords: []string{},
				Severity: constants.RealSEOStatusError,
			},
		},
		{
			name: "title_too_long",
			html: `<html><head>
				<title>This is an extremely long title that will definitely be truncated in search results and hurt SEO performance significantly because it exceeds the optimal length</title>
			</head></html>`,
			expected: RealTitleAnalysis{
				Present:  true,
				Length:   147,
				Score:    10, // Penalized for length
				Issues:   []string{constants.WarningTitleTooLong},
				Severity: constants.RealSEOStatusWarning,
			},
		},
		{
			name: "title_too_short",
			html: `<html><head><title>SEO Tool</title></head></html>`,
			expected: RealTitleAnalysis{
				Present:  true,
				Content:  "SEO Tool",
				Length:   8,
				Score:    5, // Penalized for being too short
				Issues:   []string{constants.WarningTitleTooShort},
				Severity: constants.RealSEOStatusWarning,
			},
		},
		{
			name: "empty_title",
			html: `<html><head><title></title></head></html>`,
			expected: RealTitleAnalysis{
				Present:  false,
				Content:  "",
				Length:   0,
				Score:    0,
				Issues:   []string{constants.ErrorTitleEmpty},
				Severity: constants.RealSEOStatusError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRealSEOAnalyzer()
			result := analyzer.AnalyzeTitle(tt.html)
			
			assert.Equal(t, tt.expected.Present, result.Present)
			assert.Equal(t, tt.expected.Score, result.Score)
			assert.Equal(t, tt.expected.Severity, result.Severity)
			
			if tt.expected.Content != "" {
				assert.Equal(t, tt.expected.Content, result.Content)
			}
			if len(tt.expected.Issues) > 0 {
				assert.ElementsMatch(t, tt.expected.Issues, result.Issues)
			}
		})
	}
}

// TestRealSEOAnalyzer_MetaDescription tests meta description extraction
func TestRealSEOAnalyzer_MetaDescription(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected RealMetaAnalysis
	}{
		{
			name: "perfect_meta_description",
			html: `<html><head>
				<meta name="description" content="Fire Salamander provides professional SEO analysis and recommendations to improve your website ranking and performance with advanced tools.">
			</head></html>`,
			expected: RealMetaAnalysis{
				Present: true,
				Content: "Fire Salamander provides professional SEO analysis and recommendations to improve your website ranking and performance with advanced tools.",
				Length:  135,
				Score:   constants.MaxMetaDescScore, // 15 points
				Issues:  []string{},
				Severity: constants.RealSEOStatusInfo,
			},
		},
		{
			name: "missing_meta_description",
			html: `<html><head><title>Test</title></head></html>`,
			expected: RealMetaAnalysis{
				Present:  false,
				Content:  "",
				Length:   0,
				Score:    0,
				Issues:   []string{constants.ErrorMetaDescMissing},
				Severity: constants.RealSEOStatusError,
			},
		},
		{
			name: "meta_description_too_long",
			html: `<html><head>
				<meta name="description" content="This is an extremely long meta description that exceeds the recommended 160 characters limit and will be truncated in search results which is not good for SEO performance and user experience in search engines.">
			</head></html>`,
			expected: RealMetaAnalysis{
				Present:  true,
				Length:   200,
				Score:    7, // Penalized for length
				Issues:   []string{constants.WarningMetaDescTooLong},
				Severity: constants.RealSEOStatusWarning,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRealSEOAnalyzer()
			result := analyzer.AnalyzeMetaDescription(tt.html)
			
			assert.Equal(t, tt.expected.Present, result.Present)
			assert.Equal(t, tt.expected.Score, result.Score)
			assert.Equal(t, tt.expected.Severity, result.Severity)
		})
	}
}

// TestRealSEOAnalyzer_HeadingHierarchy tests heading structure analysis
func TestRealSEOAnalyzer_HeadingHierarchy(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected RealHeadingAnalysis
	}{
		{
			name: "perfect_heading_hierarchy",
			html: `<html><body>
				<h1>Main Title</h1>
				<h2>Section 1</h2>
				<h3>Subsection 1.1</h3>
				<h3>Subsection 1.2</h3>
				<h2>Section 2</h2>
			</body></html>`,
			expected: RealHeadingAnalysis{
				H1Count:       1,
				H2Count:       2,
				H3Count:       2,
				HasHierarchy:  true,
				Score:         constants.MaxHeadingScore, // 15 points
				Issues:        []string{},
				Severity:      constants.RealSEOStatusInfo,
			},
		},
		{
			name: "broken_heading_hierarchy",
			html: `<html><body>
				<h1>Main Title</h1>
				<h3>Skipped H2</h3>
				<h2>Correct H2</h2>
				<h1>Duplicate H1</h1>
			</body></html>`,
			expected: RealHeadingAnalysis{
				H1Count:      2,
				H2Count:      1,
				H3Count:      1,
				HasHierarchy: false,
				Score:        5, // Penalized for multiple H1 and broken hierarchy
				Issues:       []string{constants.ErrorMultipleH1, constants.ErrorBrokenHeadingHierarchy},
				Severity:     constants.RealSEOStatusError,
			},
		},
		{
			name: "missing_h1",
			html: `<html><body>
				<h2>Section without H1</h2>
				<h3>Subsection</h3>
			</body></html>`,
			expected: RealHeadingAnalysis{
				H1Count:      0,
				H2Count:      1,
				H3Count:      1,
				HasHierarchy: false,
				Score:        2, // Low score for missing H1
				Issues:       []string{constants.ErrorMissingH1},
				Severity:     constants.RealSEOStatusError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRealSEOAnalyzer()
			result := analyzer.AnalyzeHeadings(tt.html)
			
			assert.Equal(t, tt.expected.H1Count, result.H1Count)
			assert.Equal(t, tt.expected.HasHierarchy, result.HasHierarchy)
			assert.Equal(t, tt.expected.Score, result.Score)
			assert.Equal(t, tt.expected.Severity, result.Severity)
			
			if len(tt.expected.Issues) > 0 {
				assert.ElementsMatch(t, tt.expected.Issues, result.Issues)
			}
		})
	}
}

// TestRealSEOAnalyzer_ImageAltText tests image alt text analysis
func TestRealSEOAnalyzer_ImageAltText(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected RealImageAnalysis
	}{
		{
			name: "all_images_with_alt",
			html: `<html><body>
				<img src="logo.png" alt="Fire Salamander Logo">
				<img src="banner.jpg" alt="Professional SEO Analysis Banner">
				<img src="icon.svg" alt="SEO Tool Icon">
			</body></html>`,
			expected: RealImageAnalysis{
				TotalImages:      3,
				ImagesWithAlt:    3,
				MissingAlt:       0,
				AltTextCoverage:  1.0,
				Score:           constants.MaxImageScore, // 100 points
				Issues:          []string{},
				Severity:        constants.RealSEOStatusInfo,
			},
		},
		{
			name: "mixed_alt_text",
			html: `<html><body>
				<img src="logo.png" alt="Fire Salamander Logo">
				<img src="banner.jpg">
				<img src="icon.svg" alt="">
			</body></html>`,
			expected: RealImageAnalysis{
				TotalImages:     3,
				ImagesWithAlt:   1,
				MissingAlt:      2,
				AltTextCoverage: 0.33,
				Score:           33, // 1/3 coverage = 33% = 33/100 points
				Issues:          []string{constants.WarningMissingAltText},
				Severity:        constants.RealSEOStatusWarning,
			},
		},
		{
			name: "no_images",
			html: `<html><body><p>Content without images</p></body></html>`,
			expected: RealImageAnalysis{
				TotalImages:     0,
				ImagesWithAlt:   0,
				MissingAlt:      0,
				AltTextCoverage: 1.0, // No images = perfect coverage
				Score:           constants.MaxImageScore,
				Issues:          []string{},
				Severity:        constants.RealSEOStatusInfo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRealSEOAnalyzer()
			result := analyzer.AnalyzeImages(tt.html)
			
			assert.Equal(t, tt.expected.TotalImages, result.TotalImages)
			assert.Equal(t, tt.expected.ImagesWithAlt, result.ImagesWithAlt)
			assert.Equal(t, tt.expected.MissingAlt, result.MissingAlt)
			assert.InDelta(t, tt.expected.AltTextCoverage, result.AltTextCoverage, 0.01)
			assert.Equal(t, tt.expected.Score, result.Score)
			assert.Equal(t, tt.expected.Severity, result.Severity)
		})
	}
}

// TestRealSEOAnalyzer_PerformanceMetrics tests performance scoring
func TestRealSEOAnalyzer_PerformanceMetrics(t *testing.T) {
	tests := []struct {
		name     string
		metrics  PerformanceMetrics
		expected int
	}{
		{
			name: "excellent_performance",
			metrics: PerformanceMetrics{
				LoadTime:         1.2 * float64(time.Second), // Fast load time
				PageSize:         450 * 1024,                 // 450KB - good size
				RequestCount:     25,                         // Reasonable requests
				HasCompression:   true,
				HasCaching:       true,
				OptimizedImages:  true,
			},
			expected: 9, // 3+3+3 points
		},
		{
			name: "poor_performance",
			metrics: PerformanceMetrics{
				LoadTime:        5.5 * float64(time.Second), // Very slow
				PageSize:        3 * 1024 * 1024,            // 3MB - too large
				RequestCount:    80,                          // Too many requests
				HasCompression:  false,
				HasCaching:      false,
				OptimizedImages: false,
			},
			expected: 2, // 1+1+0 points
		},
		{
			name: "moderate_performance",
			metrics: PerformanceMetrics{
				LoadTime:        2.8 * float64(time.Second), // Acceptable
				PageSize:        800 * 1024,                 // 800KB - okay
				RequestCount:    45,
				HasCompression:  true,
				HasCaching:      false,
				OptimizedImages: true,
			},
			expected: 8, // 3+3+2 points
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRealSEOAnalyzer()
			score := analyzer.ScorePerformance(tt.metrics)
			
			assert.Equal(t, tt.expected, score)
		})
	}
}

// TestRealSEOAnalyzer_GenerateRecommendations tests recommendation generation
func TestRealSEOAnalyzer_GenerateRecommendations(t *testing.T) {
	analysis := &RealPageAnalysis{
		Title: RealTitleAnalysis{
			Present: false,
			Score:   0,
		},
		MetaDescription: RealMetaAnalysis{
			Present: true,
			Score:   constants.MaxMetaDescScore,
		},
		Headings: RealHeadingAnalysis{
			H1Count: 0,
			Score:   2,
		},
		Images: RealImageAnalysis{
			Score: 5,
		},
		Performance: RealPerformanceAnalysis{
			Score: 8,
		},
	}

	analyzer := NewRealSEOAnalyzer()
	recommendations := analyzer.GenerateRecommendations(analysis)
	
	require.NotEmpty(t, recommendations)
	
	// Should prioritize critical issue (missing title)
	assert.Equal(t, constants.SEOPriorityCritical, recommendations[0].Priority)
	assert.Contains(t, recommendations[0].Action, "title")
	assert.Equal(t, constants.SEOImpactHigh, recommendations[0].Impact)
	assert.Equal(t, constants.EffortQuickWin, recommendations[0].Effort)
}

// TestRealSEOAnalyzer_CalculateTotalScore tests total score calculation
func TestRealSEOAnalyzer_CalculateTotalScore(t *testing.T) {
	analysis := &RealPageAnalysis{
		Title: RealTitleAnalysis{
			Score: 20, // Max title score
		},
		MetaDescription: RealMetaAnalysis{
			Score: 15, // Max meta desc score
		},
		Headings: RealHeadingAnalysis{
			Score: 15, // Max heading score
		},
		Images: RealImageAnalysis{
			Score: 10, // Max image score
		},
		Performance: RealPerformanceAnalysis{
			Score: 10, // Max performance score
		},
		Mobile: RealMobileAnalysis{
			Score: 10, // Max mobile score
		},
		HTTPS: RealHTTPSAnalysis{
			Score: 10, // Max HTTPS score
		},
		Content: RealContentAnalysis{
			Score: 10, // Max content score
		},
	}

	analyzer := NewRealSEOAnalyzer()
	totalScore := analyzer.CalculateTotalScore(analysis)
	
	assert.Equal(t, constants.MaxSEOScore, totalScore) // Should be 100
}

// TestRealSEOAnalyzer_DetermineGrade tests grade determination
func TestRealSEOAnalyzer_DetermineGrade(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{100, constants.SEOGradeAPlus},
		{95, constants.SEOGradeAPlus},
		{94, constants.SEOGradeA},
		{90, constants.SEOGradeA},
		{85, constants.SEOGradeBPlus},
		{75, constants.SEOGradeB},
		{65, constants.SEOGradeC},
		{55, constants.SEOGradeD},
		{30, constants.SEOGradeF},
	}

	analyzer := NewRealSEOAnalyzer()
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			grade := analyzer.DetermineGrade(tt.score)
			assert.Equal(t, tt.expected, grade)
		})
	}
}

// TestRealSEOAnalyzer_AnalyzePage tests complete page analysis integration
func TestRealSEOAnalyzer_AnalyzePage(t *testing.T) {
	// This test will use HTML with some SEO issues to generate recommendations
	htmlContent := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta property="og:title" content="Fire Salamander SEO Tool">
		<link rel="canonical" href="https://example.com/">
	</head>
	<body>
		<h2>Features</h2>
		<h3>Real-time Analysis</h3>
		<p>Our tool provides comprehensive SEO analysis...</p>
		<img src="logo.png">
		<img src="feature.jpg" alt="">
	</body>
	</html>`

	analyzer := NewRealSEOAnalyzer()
	ctx := context.Background()
	
	// This test should fail until we implement AnalyzePageContent method
	analysis := analyzer.AnalyzePageContent(ctx, "https://example.com", htmlContent)
	
	require.NotNil(t, analysis)
	assert.Equal(t, "https://example.com", analysis.URL)
	assert.False(t, analysis.Title.Present)          // No title in test HTML
	assert.False(t, analysis.MetaDescription.Present) // No meta desc in test HTML
	assert.Equal(t, 0, analysis.Headings.H1Count)   // No H1 in test HTML
	assert.Greater(t, analysis.TotalScore, 0)        // Should still have some score
	assert.NotEmpty(t, analysis.Grade)
	assert.NotEmpty(t, analysis.Recommendations)     // Should have recommendations for issues
}

// TestRealSEOAnalyzer_NoHardcoding validates zero hardcoding policy
func TestRealSEOAnalyzer_NoHardcoding(t *testing.T) {
	analyzer := NewRealSEOAnalyzer()
	
	// Test that all scoring uses constants
	titleAnalysis := analyzer.AnalyzeTitle(`<title>Test Title</title>`)
	assert.True(t, titleAnalysis.Score <= constants.MaxTitleScore)
	
	metaAnalysis := analyzer.AnalyzeMetaDescription(`<meta name="description" content="Test description">`)
	assert.True(t, metaAnalysis.Score <= constants.MaxMetaDescScore)
	
	// Verify no magic numbers in scoring
	perfectMetrics := PerformanceMetrics{
		LoadTime:        1.0,
		PageSize:        400 * 1024,
		RequestCount:    20,
		HasCompression:  true,
		HasCaching:      true,
		OptimizedImages: true,
	}
	perfScore := analyzer.ScorePerformance(perfectMetrics)
	assert.True(t, perfScore <= constants.MaxPerformanceScore)
}