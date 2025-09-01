package semantic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"firesalamander/internal/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSemanticClient(t *testing.T) {
	client := NewSemanticClient("http://localhost:8003")
	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8003", client.BaseURL)
	assert.NotNil(t, client.Client)
}

func TestSemanticClient_HealthCheck(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		
		response := map[string]interface{}{
			"status":  "healthy",
			"service": "semantic-analyzer",
			"version": "sem-v1.0",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewSemanticClient(server.URL)
	err := client.HealthCheck()
	require.NoError(t, err)
}

func TestSemanticClient_Analyze(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/analyze", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		
		// Verify request body
		var request SemanticRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		require.NoError(t, err)
		assert.Equal(t, "test_audit_123", request.AuditID)
		assert.NotEmpty(t, request.CrawlData.Pages)
		
		// Mock response
		response := SemanticResult{
			AuditID:      "test_audit_123",
			ModelVersion: "sem-v1.0",
			Topics: []Topic{
				{
					ID:    "topic_0",
					Label: "Logiciel juridique",
					Terms: []string{"logiciel", "avocat", "cabinet"},
				},
			},
			Suggestions: []Suggestion{
				{
					Keyword:    "logiciel avocat",
					Reason:     "Très pertinent pour le thématique",
					Confidence: 0.85,
					Evidence:   []string{"https://example.com (titre)"},
				},
			},
			Metadata: Metadata{
				SchemaVersion:   "1.0",
				WeightsVersion:  "1.0",
				ExecutionTimeMs: 250,
				Lang:            "fr",
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewSemanticClient(server.URL)
	
	// Test data
	crawlData := crawler.CrawlResult{
		Pages: []crawler.PageData{
			{
				URL:     "https://example.com",
				Lang:    "fr",
				Title:   "Logiciel pour avocats",
				H1:      "Solutions juridiques",
				Content: "Notre logiciel de gestion pour cabinets d'avocats...",
			},
		},
		Metadata: crawler.Metadata{
			TotalPages: 1,
		},
	}

	result, err := client.Analyze("test_audit_123", crawlData)
	require.NoError(t, err)
	assert.Equal(t, "test_audit_123", result.AuditID)
	assert.Equal(t, "sem-v1.0", result.ModelVersion)
	assert.Len(t, result.Topics, 1)
	assert.Len(t, result.Suggestions, 1)
	assert.Equal(t, "logiciel avocat", result.Suggestions[0].Keyword)
}

func TestSemanticClient_GetTopics(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/topics", r.URL.Path)
		
		response := map[string]interface{}{
			"topics": []Topic{
				{
					ID:    "topic_legal",
					Label: "Juridique & Droit",
					Terms: []string{"avocat", "cabinet", "juridique"},
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewSemanticClient(server.URL)
	
	pages := []crawler.PageData{
		{
			URL:     "https://example.com",
			Title:   "Cabinet d'avocats",
			Content: "Services juridiques professionnels",
		},
	}

	topics, err := client.GetTopics(pages)
	require.NoError(t, err)
	assert.Len(t, topics, 1)
	assert.Equal(t, "topic_legal", topics[0].ID)
	assert.Equal(t, "Juridique & Droit", topics[0].Label)
}

func TestSemanticClient_GetKeywords(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/keywords", r.URL.Path)
		
		response := map[string]interface{}{
			"suggestions": []Suggestion{
				{
					Keyword:    "avocat paris",
					Reason:     "Forte pertinence géographique",
					Confidence: 0.92,
					Evidence:   []string{"https://example.com/paris (H2)"},
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewSemanticClient(server.URL)
	
	pages := []crawler.PageData{
		{
			URL:   "https://example.com/paris",
			Title: "Cabinet d'avocats à Paris",
			H2:    []string{"Avocats spécialisés à Paris"},
		},
	}

	suggestions, err := client.GetKeywords(pages, 10)
	require.NoError(t, err)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "avocat paris", suggestions[0].Keyword)
	assert.Equal(t, 0.92, suggestions[0].Confidence)
}

func TestSemanticClient_ErrorHandling(t *testing.T) {
	// Test with non-existent server
	client := NewSemanticClient("http://localhost:9999")
	
	err := client.HealthCheck()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "health check failed")
	
	// Test analyze with connection error
	crawlData := crawler.CrawlResult{Pages: []crawler.PageData{}}
	_, err = client.Analyze("test", crawlData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
}