package semantic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"firesalamander/internal/agents/crawler"
)

// SemanticClient handles communication with Python semantic analysis agent
type SemanticClient struct {
	BaseURL string
	Client  *http.Client
}

// SemanticRequest represents the request to semantic analyzer
type SemanticRequest struct {
	AuditID   string                `json:"audit_id"`
	CrawlData crawler.CrawlResult `json:"crawl_data"`
}

// SemanticResult represents the response from semantic analyzer
type SemanticResult struct {
	AuditID      string   `json:"audit_id"`
	ModelVersion string   `json:"model_version"`
	Topics       []Topic  `json:"topics"`
	Suggestions  []Suggestion `json:"suggestions"`
	Metadata     Metadata `json:"metadata"`
}

// Topic represents a semantic topic cluster
type Topic struct {
	ID    string   `json:"id"`
	Label string   `json:"label"`
	Terms []string `json:"terms"`
}

// Suggestion represents a keyword suggestion
type Suggestion struct {
	Keyword    string   `json:"keyword"`
	Reason     string   `json:"reason"`
	Confidence float64  `json:"confidence"`
	Evidence   []string `json:"evidence"`
}

// Metadata contains analysis metadata
type Metadata struct {
	SchemaVersion    string `json:"schema_version"`
	WeightsVersion   string `json:"weights_version"`
	ExecutionTimeMs  int    `json:"execution_time_ms"`
	Lang             string `json:"lang"`
	Error            string `json:"error,omitempty"`
}

// NewSemanticClient creates a new client for semantic analysis agent
func NewSemanticClient(baseURL string) *SemanticClient {
	return &SemanticClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Analyze performs semantic analysis on crawl data
func (sc *SemanticClient) Analyze(auditID string, crawlData crawler.CrawlResult) (*SemanticResult, error) {
	request := SemanticRequest{
		AuditID:   auditID,
		CrawlData: crawlData,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := sc.Client.Post(
		sc.BaseURL+"/analyze",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("semantic analysis failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result SemanticResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetTopics extracts only topics from crawl data
func (sc *SemanticClient) GetTopics(pages []crawler.PageData) ([]Topic, error) {
	requestData := map[string]interface{}{
		"pages": pages,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := sc.Client.Post(
		sc.BaseURL+"/topics",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("topic extraction failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Topics []Topic `json:"topics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Topics, nil
}

// GetKeywords generates keyword suggestions from crawl data
func (sc *SemanticClient) GetKeywords(pages []crawler.PageData, limit int) ([]Suggestion, error) {
	requestData := map[string]interface{}{
		"pages": pages,
		"limit": limit,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := sc.Client.Post(
		sc.BaseURL+"/keywords",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("keyword generation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Suggestions []Suggestion `json:"suggestions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Suggestions, nil
}

// HealthCheck verifies the semantic agent is running
func (sc *SemanticClient) HealthCheck() error {
	resp, err := sc.Client.Get(sc.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("semantic agent unhealthy: status %d", resp.StatusCode)
	}

	return nil
}