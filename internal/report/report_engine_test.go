package report

import (
	"os"
	"path/filepath"
	"testing"

	"firesalamander/internal/audit"
	"firesalamander/internal/agents/crawler"
	"firesalamander/internal/agents/semantic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportEngine(t *testing.T) {
	engine := NewReportEngine()
	assert.NotNil(t, engine)
}

func TestGenerateReport(t *testing.T) {
	engine := NewReportEngine()

	// Create test data
	auditResults := AuditResults{
		AuditID:    "test_audit_123",
		SiteURL:    "https://example.com",
		StartedAt:  "2025-09-01T10:00:00Z",
		Duration:   "2m30s",
		TotalPages: 5,
		CrawlData: crawler.CrawlResult{
			Pages: []crawler.PageData{
				{
					URL:     "https://example.com/",
					Title:   "Accueil - Cabinet d'avocats",
					H1:      "Cabinet juridique spécialisé",
					Content: "Notre cabinet propose des services juridiques...",
					Depth:   0,
				},
			},
			Metadata: crawler.Metadata{
				TotalPages:      1,
				MaxDepthReached: 0,
				DurationMs:      2500,
			},
		},
		TechResults: audit.TechResult{
			AuditID: "test_audit_123",
			Findings: []audit.Finding{
				{
					ID:       "missing-meta-description",
					Severity: "medium",
					Message:  "Balise meta description manquante",
					Evidence: []string{"https://example.com/"},
				},
			},
			Scores: audit.Scores{
				Performance:   0.85,
				Accessibility: 0.92,
				BestPractices: 0.88,
				SEO:           0.76,
			},
		},
		SemanticResults: semantic.SemanticResult{
			AuditID:      "test_audit_123",
			ModelVersion: "sem-v1.0",
			Topics: []semantic.Topic{
				{
					ID:    "topic_legal",
					Label: "Services juridiques",
					Terms: []string{"avocat", "cabinet", "juridique"},
				},
			},
			Suggestions: []semantic.Suggestion{
				{
					Keyword:    "cabinet avocat paris",
					Reason:     "Forte pertinence géographique",
					Confidence: 0.92,
					Evidence:   []string{"https://example.com/ (titre)"},
				},
			},
		},
	}

	// Test HTML generation
	htmlReport, err := engine.GenerateHTML(auditResults)
	require.NoError(t, err)
	assert.Contains(t, htmlReport, "Fire Salamander")
	assert.Contains(t, htmlReport, "test_audit_123")
	assert.Contains(t, htmlReport, "cabinet avocat paris")
	assert.Contains(t, htmlReport, "Services juridiques")

	// Test JSON generation
	jsonReport, err := engine.GenerateJSON(auditResults)
	require.NoError(t, err)
	assert.Contains(t, jsonReport, "test_audit_123")
	assert.Contains(t, jsonReport, "cabinet avocat paris")

	// Test CSV generation
	csvReport, err := engine.GenerateCSV(auditResults)
	require.NoError(t, err)
	assert.Contains(t, csvReport, "URL,Title,H1")
	assert.Contains(t, csvReport, "https://example.com/")
}

func TestSaveReport(t *testing.T) {
	engine := NewReportEngine()

	auditResults := AuditResults{
		AuditID: "test_save_123",
		SiteURL: "https://example.com",
	}

	// Create temporary directory
	tempDir := t.TempDir()

	// Test HTML save
	htmlPath, err := engine.SaveReport(auditResults, "html", tempDir)
	require.NoError(t, err)
	assert.True(t, filepath.IsAbs(htmlPath))
	assert.Contains(t, htmlPath, "test_save_123")
	assert.Contains(t, htmlPath, ".html")

	// Verify file exists and has content
	_, err = os.Stat(htmlPath)
	require.NoError(t, err)

	content, err := os.ReadFile(htmlPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Fire Salamander")

	// Test JSON save
	jsonPath, err := engine.SaveReport(auditResults, "json", tempDir)
	require.NoError(t, err)
	assert.Contains(t, jsonPath, ".json")

	// Test CSV save
	csvPath, err := engine.SaveReport(auditResults, "csv", tempDir)
	require.NoError(t, err)
	assert.Contains(t, csvPath, ".csv")
}

func TestTemplateRendering(t *testing.T) {
	engine := NewReportEngine()

	templateData := TemplateData{
		AuditID:        "test_template_123",
		SiteURL:        "https://example.com",
		GeneratedAt:    "2025-09-01 10:00:00",
		TotalPages:     3,
		CriticalIssues: 2,
		HighIssues:     5,
		MediumIssues:   8,
		LowIssues:      12,
		OverallScore:   0.83,
		Pages: []PageSummary{
			{
				URL:           "https://example.com/",
				Title:         "Accueil",
				IssuesCount:   3,
				PerformanceScore: 0.85,
			},
		},
		Issues: []IssueSummary{
			{
				ID:       "missing-title",
				Severity: "critical",
				Message:  "Titre manquant",
				Count:    1,
				Pages:    []string{"https://example.com/contact"},
			},
		},
		Keywords: []KeywordSummary{
			{
				Keyword:    "cabinet avocat",
				Confidence: 0.88,
				Reason:     "Forte pertinence thématique",
			},
		},
	}

	html, err := engine.renderHTMLTemplate(templateData)
	require.NoError(t, err)
	assert.Contains(t, html, "test_template_123")
	assert.Contains(t, html, "cabinet avocat")
	assert.Contains(t, html, "Titre manquant")
}

func TestValidateAuditResults(t *testing.T) {
	engine := NewReportEngine()

	// Valid results
	validResults := AuditResults{
		AuditID: "test_valid_123",
		SiteURL: "https://example.com",
	}
	err := engine.validateAuditResults(validResults)
	assert.NoError(t, err)

	// Invalid results - missing audit ID
	invalidResults := AuditResults{
		SiteURL: "https://example.com",
	}
	err = engine.validateAuditResults(invalidResults)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "audit_id")

	// Invalid results - missing site URL
	invalidResults2 := AuditResults{
		AuditID: "test_123",
	}
	err = engine.validateAuditResults(invalidResults2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "site_url")
}

func TestUnsupportedFormat(t *testing.T) {
	engine := NewReportEngine()

	auditResults := AuditResults{
		AuditID: "test_123",
		SiteURL: "https://example.com",
	}

	// Test unsupported format
	_, err := engine.SaveReport(auditResults, "xml", "/tmp")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}