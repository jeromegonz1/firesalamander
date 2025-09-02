package technical

import (
	"context"
	"strings"
	"testing"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
)

func TestTechnicalAuditor_Name(t *testing.T) {
	auditor := NewTechnicalAuditor()
	
	expected := constants.AgentNameTechnical
	if auditor.Name() != expected {
		t.Errorf("Expected name %s, got %s", expected, auditor.Name())
	}
}

func TestTechnicalAuditor_HealthCheck(t *testing.T) {
	auditor := NewTechnicalAuditor()
	
	err := auditor.HealthCheck()
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestTechnicalAuditor_Process(t *testing.T) {
	auditor := NewTechnicalAuditor()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         interface{}
		expectedStatus string
	}{
		{
			name: "valid page data",
			input: &agents.PageData{
				URL:  "https://example.com",
				HTML: "<html><head><title>Test Page</title></head><body><h1>Test</h1></body></html>",
				Headers: map[string]string{
					"Content-Type": "text/html",
				},
			},
			expectedStatus: constants.StatusCompleted,
		},
		{
			name:          "invalid input type",
			input:         "invalid",
			expectedStatus: constants.StatusFailed,
		},
		{
			name:          "nil input",
			input:         nil,
			expectedStatus: constants.StatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := auditor.Process(ctx, tt.input)
			
			if err != nil {
				t.Errorf("Process returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("Process returned nil result")
				return
			}
			
			if result.AgentName != constants.AgentNameTechnical {
				t.Errorf("Expected agent name %s, got %s", constants.AgentNameTechnical, result.AgentName)
			}
			
			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
			
			if result.Duration < 0 {
				t.Error("Expected non-negative duration")
			}
		})
	}
}

func TestTechnicalAuditor_AuditPage(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name        string
		pageData    *agents.PageData
		expectError bool
	}{
		{
			name: "complete valid page",
			pageData: &agents.PageData{
				URL: "https://example.com",
				HTML: `
					<html lang="en">
						<head>
							<title>Test SEO Title</title>
							<meta name="description" content="This is a test meta description that is long enough to be considered valid for SEO purposes">
							<meta name="viewport" content="width=device-width, initial-scale=1">
						</head>
						<body>
							<h1>Main Title</h1>
							<img src="test.jpg" alt="Test image">
							<a href="/test">Test Link</a>
						</body>
					</html>
				`,
				Headers: map[string]string{
					"Content-Type":     "text/html",
					"Cache-Control":    "max-age=3600",
					"Content-Encoding": "gzip",
				},
			},
			expectError: false,
		},
		{
			name: "page with SEO issues",
			pageData: &agents.PageData{
				URL: "https://example.com",
				HTML: `
					<html>
						<head>
							<title>Bad</title>
						</head>
						<body>
							<h1>Title 1</h1>
							<h1>Title 2</h1>
							<img src="test.jpg">
						</body>
					</html>
				`,
				Headers: map[string]string{
					"Content-Type": "text/html",
				},
			},
			expectError: false,
		},
		{
			name:        "nil page data",
			pageData:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := auditor.AuditPage(tt.pageData)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("AuditPage returned error: %v", err)
				return
			}
			
			if report == nil {
				t.Error("AuditPage returned nil report")
				return
			}
			
			// Vérifications de base du rapport
			if report.PageURL != tt.pageData.URL {
				t.Errorf("Expected page URL %s, got %s", tt.pageData.URL, report.PageURL)
			}
			
			// Vérification des scores (doivent être entre 0 et 100)
			if report.Performance.Score < 0 || report.Performance.Score > 100 {
				t.Errorf("Performance score out of range: %d", report.Performance.Score)
			}
			
			if report.Accessibility.Score < 0 || report.Accessibility.Score > 100 {
				t.Errorf("Accessibility score out of range: %d", report.Accessibility.Score)
			}
			
			if report.SEO.Score < 0 || report.SEO.Score > 100 {
				t.Errorf("SEO score out of range: %d", report.SEO.Score)
			}
		})
	}
}

func TestTechnicalAuditor_ValidateStructure(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name           string
		html           string
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid HTML structure",
			html: "<html><head><title>Test</title></head><body><h1>Test</h1></body></html>",
			expectedValid: true,
			expectedErrors: 0,
		},
		{
			name: "missing head section",
			html: "<html><body><h1>Test</h1></body></html>",
			expectedValid: false,
			expectedErrors: 1,
		},
		{
			name: "unclosed tag",
			html: "<html><head><title>Test</title></head><body><h1>Test<p>Unclosed</body></html>",
			expectedValid: false,
			expectedErrors: 2, // h1 et p non fermés
		},
		{
			name: "empty HTML",
			html: "",
			expectedValid: false,
			expectedErrors: 1,
		},
		{
			name: "complete missing structure",
			html: "<div>Just content</div>",
			expectedValid: false,
			expectedErrors: 3, // missing html, head, body
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := auditor.ValidateStructure(tt.html)
			
			if err != nil {
				t.Errorf("ValidateStructure returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("ValidateStructure returned nil result")
				return
			}
			
			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid %v, got %v", tt.expectedValid, result.Valid)
			}
			
			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tt.expectedErrors, len(result.Errors))
			}
			
			// Vérification des niveaux de titres
			if result.HeadingLevel < 0 || result.HeadingLevel > 6 {
				t.Errorf("Heading level out of range: %d", result.HeadingLevel)
			}
		})
	}
}

func TestTechnicalAuditor_AuditPerformance(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name        string
		pageData    *agents.PageData
		expectScore int // Score minimum attendu
	}{
		{
			name: "optimized page",
			pageData: &agents.PageData{
				URL:  "https://example.com",
				HTML: "<html><head><title>Test</title></head><body>Small content</body></html>",
				Headers: map[string]string{
					"Content-Encoding": "gzip",
					"Cache-Control":    "max-age=3600",
				},
			},
			expectScore: 90,
		},
		{
			name: "heavy page with many resources",
			pageData: &agents.PageData{
				URL: "https://example.com",
				HTML: strings.Repeat(`<img src="image.jpg"><script src="script.js"></script><link href="style.css">`, 30) +
					strings.Repeat("<p>Heavy content</p>", 1000),
				Headers: map[string]string{},
			},
			expectScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			performance := auditor.auditPerformance(tt.pageData)
			
			if performance.Score < tt.expectScore {
				t.Errorf("Expected score >= %d, got %d", tt.expectScore, performance.Score)
			}
			
			if performance.Resources < 0 {
				t.Error("Resource count should be non-negative")
			}
			
			if performance.LoadTime <= 0 {
				t.Error("Load time should be positive")
			}
		})
	}
}

func TestTechnicalAuditor_AuditAccessibility(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name        string
		html        string
		expectScore int // Score minimum attendu
	}{
		{
			name: "accessible content",
			html: `
				<img src="test.jpg" alt="Test image">
				<a href="/test">Descriptive link text</a>
				<input type="text" id="name">
				<label for="name">Name:</label>
			`,
			expectScore: 90,
		},
		{
			name: "inaccessible content",
			html: `
				<img src="test.jpg">
				<img src="test2.jpg" alt="">
				<a href="/test">x</a>
				<input type="text">
			`,
			expectScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageData := &agents.PageData{
				URL:  "https://example.com",
				HTML: tt.html,
			}
			
			accessibility := auditor.auditAccessibility(pageData)
			
			if accessibility.Score < tt.expectScore {
				t.Errorf("Expected score >= %d, got %d", tt.expectScore, accessibility.Score)
			}
			
			if accessibility.Score < 100 && len(accessibility.Issues) == 0 {
				t.Error("Expected issues when score is not perfect")
			}
		})
	}
}

func TestTechnicalAuditor_AuditSEO(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name        string
		html        string
		expectScore int // Score minimum attendu
	}{
		{
			name: "SEO optimized",
			html: `
				<html lang="en">
					<head>
						<title>Perfect SEO Title Length</title>
						<meta name="description" content="This is a perfect meta description that provides valuable information about the page content and is within the optimal length range">
						<meta name="viewport" content="width=device-width, initial-scale=1">
					</head>
					<body>
						<h1>Single Main Heading</h1>
					</body>
				</html>
			`,
			expectScore: 90,
		},
		{
			name: "SEO issues",
			html: `
				<html>
					<head>
						<title>Bad</title>
					</head>
					<body>
						<h1>First</h1>
						<h1>Second</h1>
					</body>
				</html>
			`,
			expectScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageData := &agents.PageData{
				URL:  "https://example.com",
				HTML: tt.html,
			}
			
			seo := auditor.auditSEO(pageData)
			
			if seo.Score < tt.expectScore {
				t.Errorf("Expected score >= %d, got %d", tt.expectScore, seo.Score)
			}
			
			if seo.Score < 100 && len(seo.MissingElements) == 0 {
				t.Error("Expected missing elements when score is not perfect")
			}
		})
	}
}

func TestTechnicalAuditor_ExtractAttribute(t *testing.T) {
	auditor := NewTechnicalAuditor()

	tests := []struct {
		name      string
		tag       string
		attribute string
		expected  string
	}{
		{
			name:      "extract id",
			tag:       `<input type="text" id="username" name="user">`,
			attribute: "id",
			expected:  "username",
		},
		{
			name:      "extract type",
			tag:       `<input type="password" id="pwd">`,
			attribute: "type",
			expected:  "password",
		},
		{
			name:      "missing attribute",
			tag:       `<input type="text">`,
			attribute: "id",
			expected:  "",
		},
		{
			name:      "single quotes",
			tag:       `<input type='email' id='mail'>`,
			attribute: "type",
			expected:  "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := auditor.extractAttribute(tt.tag, tt.attribute)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}