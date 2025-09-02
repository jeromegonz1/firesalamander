package linking

import (
	"context"
	"testing"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
	"firesalamander/internal/agents/crawler"
)

func TestLinkingMapper_Name(t *testing.T) {
	mapper := NewLinkingMapper()
	
	expected := constants.AgentNameLinking
	if mapper.Name() != expected {
		t.Errorf("Expected name %s, got %s", expected, mapper.Name())
	}
}

func TestLinkingMapper_HealthCheck(t *testing.T) {
	mapper := NewLinkingMapper()
	
	err := mapper.HealthCheck()
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestLinkingMapper_Process(t *testing.T) {
	mapper := NewLinkingMapper()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         interface{}
		expectedStatus string
	}{
		{
			name: "valid crawl result",
			input: &crawler.CrawlResult{
				Pages: []crawler.PageData{
					{
						URL:   "https://example.com",
						Title: "Test Page",
						Content: `<a href="https://example.com/page1">Internal Link</a>
								  <a href="https://external.com">External Link</a>`,
					},
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
			result, err := mapper.Process(ctx, tt.input)
			
			if err != nil {
				t.Errorf("Process returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("Process returned nil result")
				return
			}
			
			if result.AgentName != constants.AgentNameLinking {
				t.Errorf("Expected agent name %s, got %s", constants.AgentNameLinking, result.AgentName)
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

func TestLinkingMapper_MapLinks(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name               string
		crawlResult        *crawler.CrawlResult
		expectError        bool
		expectedInternal   int
		expectedExternal   int
	}{
		{
			name: "mixed internal and external links",
			crawlResult: &crawler.CrawlResult{
				Pages: []crawler.PageData{
					{
						URL:   "https://example.com/page1",
						Title: "Page 1",
						Content: `
							<a href="https://example.com/page2">Internal Link</a>
							<a href="https://external.com">External Link</a>
							<a href="/relative">Relative Internal</a>
							<a href="#anchor">Anchor Link</a>
						`,
					},
				},
			},
			expectError:      false,
			expectedInternal: 2, // /page2, /relative (anchor est séparé)
			expectedExternal: 1, // external.com
		},
		{
			name: "no links",
			crawlResult: &crawler.CrawlResult{
				Pages: []crawler.PageData{
					{
						URL:     "https://example.com",
						Title:   "No Links Page",
						Content: "<p>Just text content</p>",
					},
				},
			},
			expectError:      false,
			expectedInternal: 0,
			expectedExternal: 0,
		},
		{
			name:        "nil crawl result",
			crawlResult: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkMap, err := mapper.MapLinks(tt.crawlResult)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("MapLinks returned error: %v", err)
				return
			}
			
			if linkMap == nil {
				t.Error("MapLinks returned nil link map")
				return
			}
			
			if len(linkMap.InternalLinks) != tt.expectedInternal {
				t.Errorf("Expected %d internal links, got %d", tt.expectedInternal, len(linkMap.InternalLinks))
			}
			
			if len(linkMap.ExternalLinks) != tt.expectedExternal {
				t.Errorf("Expected %d external links, got %d", tt.expectedExternal, len(linkMap.ExternalLinks))
			}
			
			// Vérification des statistiques (ne compte que internal + external, pas anchor)
			expectedTotal := len(linkMap.InternalLinks) + len(linkMap.ExternalLinks)
			if linkMap.Statistics.TotalLinks != expectedTotal {
				t.Errorf("Expected %d total links in statistics, got %d", expectedTotal, linkMap.Statistics.TotalLinks)
			}
		})
	}
}

func TestLinkingMapper_AnalyzeLinkStructure(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name                    string
		links                   []agents.Link
		expectedRecommendations int
	}{
		{
			name: "well structured links",
			links: []agents.Link{
				{
					Source:     "https://example.com/page1",
					Target:     "https://example.com/page2",
					AnchorText: "Good descriptive anchor text",
					Type:       "internal",
					IsNoFollow: false,
				},
				{
					Source:     "https://example.com/page1",
					Target:     "https://external.com",
					AnchorText: "External resource",
					Type:       "external",
					IsNoFollow: true,
				},
			},
			expectedRecommendations: 1, // Should be well optimized
		},
		{
			name: "problematic links",
			links: []agents.Link{
				{
					Source:     "https://example.com/orphan",
					Target:     "https://example.com/page2",
					AnchorText: "x",
					Type:       "internal",
					IsNoFollow: false,
				},
				{
					Source:     "https://example.com/page1",
					Target:     "https://external.com",
					AnchorText: "External",
					Type:       "external",
					IsNoFollow: false,
				},
			},
			expectedRecommendations: 3, // Multiple issues
		},
		{
			name:                    "no links",
			links:                   []agents.Link{},
			expectedRecommendations: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis, err := mapper.AnalyzeLinkStructure(tt.links)
			
			if err != nil {
				t.Errorf("AnalyzeLinkStructure returned error: %v", err)
				return
			}
			
			if analysis == nil {
				t.Error("AnalyzeLinkStructure returned nil analysis")
				return
			}
			
			if len(analysis.Recommendations) != tt.expectedRecommendations {
				t.Errorf("Expected %d recommendations, got %d: %v", 
					tt.expectedRecommendations, len(analysis.Recommendations), analysis.Recommendations)
			}
			
			// Vérifications de base
			if analysis.LinkEquity == nil {
				t.Error("LinkEquity should not be nil")
			}
			
			if analysis.OrphanPages == nil {
				t.Error("OrphanPages should not be nil")
			}
			
			if analysis.HighTrafficPages == nil {
				t.Error("HighTrafficPages should not be nil")
			} else {
				// HighTrafficPages peut être vide mais ne doit pas être nil
			}
		})
	}
}

func TestLinkingMapper_ExtractDomain(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "simple domain",
			url:      "https://example.com/path",
			expected: "example.com",
		},
		{
			name:     "subdomain",
			url:      "https://www.example.com/path",
			expected: "www.example.com",
		},
		{
			name:     "with port",
			url:      "https://example.com:8080/path",
			expected: "example.com:8080",
		},
		{
			name:     "invalid URL",
			url:      "not-a-url",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.extractDomain(tt.url)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLinkingMapper_ResolveURL(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name     string
		href     string
		baseURL  string
		expected string
	}{
		{
			name:     "absolute URL",
			href:     "https://external.com/page",
			baseURL:  "https://example.com/current",
			expected: "https://external.com/page",
		},
		{
			name:     "relative path",
			href:     "/about",
			baseURL:  "https://example.com/current",
			expected: "https://example.com/about",
		},
		{
			name:     "relative path with subdirectory",
			href:     "contact",
			baseURL:  "https://example.com/pages/",
			expected: "https://example.com/pages/contact",
		},
		{
			name:     "parent directory",
			href:     "../home",
			baseURL:  "https://example.com/pages/current",
			expected: "https://example.com/home",
		},
		{
			name:     "invalid base URL",
			href:     "/test",
			baseURL:  "not-a-url",
			expected: "/test", // url.Parse retourne un résultat même avec une base invalide
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.resolveURL(tt.href, tt.baseURL)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLinkingMapper_DetermineLinkType(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name       string
		linkURL    string
		baseDomain string
		expected   string
	}{
		{
			name:       "internal link",
			linkURL:    "https://example.com/page",
			baseDomain: "example.com",
			expected:   "internal",
		},
		{
			name:       "external link",
			linkURL:    "https://external.com/page",
			baseDomain: "example.com",
			expected:   "external",
		},
		{
			name:       "anchor link",
			linkURL:    "#section",
			baseDomain: "example.com",
			expected:   "anchor",
		},
		{
			name:       "subdomain external",
			linkURL:    "https://api.example.com/data",
			baseDomain: "example.com",
			expected:   "external",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.determineLinkType(tt.linkURL, tt.baseDomain)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLinkingMapper_CleanHTML(t *testing.T) {
	mapper := NewLinkingMapper()

	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple text",
			html:     "Simple text",
			expected: "Simple text",
		},
		{
			name:     "with HTML tags",
			html:     "<strong>Bold</strong> text",
			expected: "Bold text",
		},
		{
			name:     "multiple spaces",
			html:     "Text   with    spaces",
			expected: "Text with spaces",
		},
		{
			name:     "complex HTML",
			html:     "<span class='highlight'>Important</span> <em>text</em>",
			expected: "Important text",
		},
		{
			name:     "nested tags",
			html:     "<div><p><strong>Nested</strong> content</p></div>",
			expected: "Nested content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.cleanHTML(tt.html)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}