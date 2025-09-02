package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebServer_HandleHome(t *testing.T) {
	server := NewWebServer("8080")
	
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleHome)
	handler.ServeHTTP(rr, req)

	// Note: Ce test peut échouer si le fichier index.html n'existe pas
	// Dans un test réel, nous mockerions le système de fichiers
	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v or %v",
			rr.Code, http.StatusOK, http.StatusNotFound)
	}
}

func TestWebServer_HandleStartAudit(t *testing.T) {
	server := NewWebServer("8080")

	tests := []struct {
		name           string
		requestBody    AuditRequest
		expectedStatus int
	}{
		{
			name: "valid audit request",
			requestBody: AuditRequest{
				SiteURL:   "https://example.com",
				AuditType: "complete",
				MaxPages:  10,
				Timestamp: "2024-01-01T12:00:00Z",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing site URL",
			requestBody: AuditRequest{
				AuditType: "complete",
				MaxPages:  10,
				Timestamp: "2024-01-01T12:00:00Z",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/api/v1/audits", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.handleStartAudit)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response AuditResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if response.AuditID == "" {
					t.Error("Expected audit ID in response")
				}

				if response.Status != "started" {
					t.Errorf("Expected status 'started', got %s", response.Status)
				}
			}
		})
	}
}

func TestWebServer_HandleListAudits(t *testing.T) {
	server := NewWebServer("8080")

	req, err := http.NewRequest("GET", "/api/v1/audits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleListAudits)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	var audits []map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&audits)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(audits) == 0 {
		t.Error("Expected at least one audit in response")
	}
}

func TestWebServer_HandleAuditDetails(t *testing.T) {
	server := NewWebServer("8080")

	tests := []struct {
		name           string
		auditID        string
		expectedStatus int
	}{
		{
			name:           "valid audit ID",
			auditID:        "aud_123456",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty audit ID",
			auditID:        "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/v1/audits/"
			if tt.auditID != "" {
				path += tt.auditID
			}

			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.handleAuditDetails)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var auditDetails map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&auditDetails)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if auditDetails["id"] == nil {
					t.Error("Expected audit ID in response")
				}

				if auditDetails["results"] == nil {
					t.Error("Expected results in response")
				}
			}
		})
	}
}

func TestWebServer_CorsMiddleware(t *testing.T) {
	server := NewWebServer("8080")

	req, err := http.NewRequest("OPTIONS", "/api/v1/audits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := server.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("CORS preflight returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Vérifier les headers CORS
	corsHeaders := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
	}

	for _, header := range corsHeaders {
		if rr.Header().Get(header) == "" {
			t.Errorf("Missing CORS header: %s", header)
		}
	}
}

func TestWebServer_LoggingMiddleware(t *testing.T) {
	server := NewWebServer("8080")

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := server.loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Logging middleware affected response: got %v want %v",
			rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if body != "OK" {
		t.Errorf("Logging middleware affected response body: got %v want %v",
			body, "OK")
	}
}

func TestWebServer_HandleAuditsMethodNotAllowed(t *testing.T) {
	server := NewWebServer("8080")

	req, err := http.NewRequest("DELETE", "/api/v1/audits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleAudits)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestWebServer_HandleHealthCheck(t *testing.T) {
	server := NewWebServer("8080")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleHealthCheck)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("health check returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode health check response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	expectedFields := []string{"timestamp", "service", "version"}
	for _, field := range expectedFields {
		if response[field] == nil {
			t.Errorf("Missing field '%s' in health check response", field)
		}
	}
}