package audit

import (
	"testing"

	"firesalamander/internal/crawler"
	"github.com/stretchr/testify/assert"
)

func TestNewTechnicalAnalyzer(t *testing.T) {
	rules := TechRules{
		Title: TitleRules{
			MinLength:       15,
			MaxLength:       60,
			MissingSeverity: "critical",
		},
	}

	analyzer := NewTechnicalAnalyzer(rules)
	assert.NotNil(t, analyzer)
	assert.Equal(t, "critical", analyzer.Rules.Title.MissingSeverity)
}

func TestValidateTitle(t *testing.T) {
	rules := TechRules{
		Title: TitleRules{
			MinLength:        15,
			MaxLength:        60,
			MissingSeverity:  "critical",
			TooShortSeverity: "high",
			TooLongSeverity:  "medium",
		},
	}

	analyzer := NewTechnicalAnalyzer(rules)

	tests := []struct {
		title    string
		expected []Finding
		name     string
	}{
		{
			title:    "",
			expected: []Finding{{ID: "missing-title", Severity: "critical", Message: "Titre manquant"}},
			name:     "missing title",
		},
		{
			title:    "Short",
			expected: []Finding{{ID: "title-too-short", Severity: "high", Message: "Titre trop court (5 caractères, minimum 15)"}},
			name:     "title too short",
		},
		{
			title:    "This is a very long title that exceeds the maximum length allowed",
			expected: []Finding{{ID: "title-too-long", Severity: "medium", Message: "Titre trop long (73 caractères, maximum 60)"}},
			name:     "title too long",
		},
		{
			title:    "Perfect title length",
			expected: []Finding{},
			name:     "valid title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.ValidateTitle(tt.title, "https://example.com")
			assert.Equal(t, len(tt.expected), len(result))
			if len(tt.expected) > 0 {
				assert.Equal(t, tt.expected[0].ID, result[0].ID)
				assert.Equal(t, tt.expected[0].Severity, result[0].Severity)
			}
		})
	}
}

func TestValidateHeadings(t *testing.T) {
	rules := TechRules{
		Headings: HeadingRules{
			H1: H1Rules{
				Required:         true,
				MissingSeverity:  "critical",
				MultipleSeverity: "medium",
			},
		},
	}

	analyzer := NewTechnicalAnalyzer(rules)

	tests := []struct {
		h1Count  int
		expected []Finding
		name     string
	}{
		{
			h1Count:  0,
			expected: []Finding{{ID: "missing-h1", Severity: "critical"}},
			name:     "missing H1",
		},
		{
			h1Count:  2,
			expected: []Finding{{ID: "multiple-h1", Severity: "medium"}},
			name:     "multiple H1",
		},
		{
			h1Count:  1,
			expected: []Finding{},
			name:     "valid H1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.ValidateHeadings(tt.h1Count, 2, "https://example.com")
			assert.Equal(t, len(tt.expected), len(result))
			if len(tt.expected) > 0 {
				assert.Equal(t, tt.expected[0].ID, result[0].ID)
				assert.Equal(t, tt.expected[0].Severity, result[0].Severity)
			}
		})
	}
}

func TestAnalyzeMesh(t *testing.T) {
	pages := []crawler.PageData{
		{URL: "https://example.com/", OutgoingLinks: []string{"/about", "/contact"}},
		{URL: "https://example.com/about", OutgoingLinks: []string{"/"}},
		{URL: "https://example.com/contact", OutgoingLinks: []string{"/"}},
		{URL: "https://example.com/orphan", OutgoingLinks: []string{}}, // Orphan page
	}

	analyzer := NewTechnicalAnalyzer(TechRules{})
	mesh := analyzer.AnalyzeMesh(pages)

	assert.Contains(t, mesh.Orphans, "https://example.com/orphan")
	assert.Equal(t, 0, mesh.DepthStats.Min)
	assert.GreaterOrEqual(t, mesh.DepthStats.Max, 0)
}

func TestComputeScores(t *testing.T) {
	lighthouseResults := []LighthouseResult{
		{Performance: 0.8, Accessibility: 0.9, BestPractices: 0.85, SEO: 0.95},
		{Performance: 0.7, Accessibility: 0.8, BestPractices: 0.9, SEO: 0.9},
	}

	analyzer := NewTechnicalAnalyzer(TechRules{})
	scores := analyzer.ComputeScores(lighthouseResults)

	assert.InDelta(t, 0.75, scores.Performance, 0.001)
	assert.InDelta(t, 0.85, scores.Accessibility, 0.001)
	assert.InDelta(t, 0.875, scores.BestPractices, 0.001)
	assert.InDelta(t, 0.925, scores.SEO, 0.001)
}

