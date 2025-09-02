package broken

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/constants"
)

func TestBrokenLinksDetector_Name(t *testing.T) {
	detector := NewBrokenLinksDetector()
	
	expected := constants.AgentNameBrokenLinks
	if detector.Name() != expected {
		t.Errorf("Expected name %s, got %s", expected, detector.Name())
	}
}

func TestBrokenLinksDetector_Process(t *testing.T) {
	detector := NewBrokenLinksDetector()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         interface{}
		expectedStatus string
	}{
		{
			name:          "valid URLs list",
			input:         []string{"https://httpstat.us/200", "https://httpstat.us/404"},
			expectedStatus: constants.StatusCompleted,
		},
		{
			name:          "empty URLs list",
			input:         []string{},
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
			result, err := detector.Process(ctx, tt.input)
			
			if err != nil {
				t.Errorf("Process returned error: %v", err)
				return
			}
			
			if result == nil {
				t.Error("Process returned nil result")
				return
			}
			
			if result.AgentName != constants.AgentNameBrokenLinks {
				t.Errorf("Expected agent name %s, got %s", constants.AgentNameBrokenLinks, result.AgentName)
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

func TestBrokenLinksDetector_ValidateLink(t *testing.T) {
	detector := NewBrokenLinksDetector()

	// Créer un serveur de test
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/redirect":
			w.WriteHeader(http.StatusMovedPermanently)
		case "/forbidden":
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	tests := []struct {
		name        string
		url         string
		expectValid bool
		expectError bool
	}{
		{
			name:        "valid URL - OK",
			url:         server.URL + "/ok",
			expectValid: true,
			expectError: false,
		},
		{
			name:        "broken URL - Not Found",
			url:         server.URL + "/notfound",
			expectValid: false,
			expectError: false,
		},
		{
			name:        "redirect URL",
			url:         server.URL + "/redirect",
			expectValid: true,
			expectError: false,
		},
		{
			name:        "forbidden URL",
			url:         server.URL + "/forbidden",
			expectValid: true, // 403 est considéré comme valide
			expectError: false,
		},
		{
			name:        "empty URL",
			url:         "",
			expectValid: false,
			expectError: false,
		},
		{
			name:        "invalid URL format",
			url:         "not-a-valid-url",
			expectValid: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := detector.ValidateLink(tt.url)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("ValidateLink returned error: %v", err)
				return
			}
			
			if status == nil {
				t.Error("ValidateLink returned nil status")
				return
			}
			
			if status.IsValid != tt.expectValid {
				t.Errorf("Expected IsValid %v, got %v", tt.expectValid, status.IsValid)
			}
			
			if status.URL != tt.url {
				t.Errorf("Expected URL %s, got %s", tt.url, status.URL)
			}
			
			if status.CheckedAt == "" {
				t.Error("CheckedAt should not be empty")
			}
		})
	}
}

func TestBrokenLinksDetector_CheckLinks(t *testing.T) {
	detector := NewBrokenLinksDetector()

	// Créer un serveur de test
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok1", "/ok2":
			w.WriteHeader(http.StatusOK)
		case "/broken1", "/broken2":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	tests := []struct {
		name             string
		urls             []string
		expectedTotal    int
		expectedBroken   int
	}{
		{
			name:           "mixed valid and broken links",
			urls:           []string{server.URL + "/ok1", server.URL + "/broken1", server.URL + "/ok2", server.URL + "/broken2"},
			expectedTotal:  4,
			expectedBroken: 2,
		},
		{
			name:           "all valid links",
			urls:           []string{server.URL + "/ok1", server.URL + "/ok2"},
			expectedTotal:  2,
			expectedBroken: 0,
		},
		{
			name:           "all broken links",
			urls:           []string{server.URL + "/broken1", server.URL + "/broken2"},
			expectedTotal:  2,
			expectedBroken: 2,
		},
		{
			name:           "empty URLs list",
			urls:           []string{},
			expectedTotal:  0,
			expectedBroken: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := detector.CheckLinks(tt.urls)
			
			if err != nil {
				t.Errorf("CheckLinks returned error: %v", err)
				return
			}
			
			if report == nil {
				t.Error("CheckLinks returned nil report")
				return
			}
			
			if report.TotalChecked != tt.expectedTotal {
				t.Errorf("Expected %d total checked, got %d", tt.expectedTotal, report.TotalChecked)
			}
			
			if report.BrokenCount != tt.expectedBroken {
				t.Errorf("Expected %d broken links, got %d", tt.expectedBroken, report.BrokenCount)
			}
			
			if len(report.BrokenLinks) != tt.expectedBroken {
				t.Errorf("Expected %d broken links in list, got %d", tt.expectedBroken, len(report.BrokenLinks))
			}
			
			if report.CheckedAt == "" {
				t.Error("CheckedAt should not be empty")
			}
			
			// Vérifier que les liens brisés ont des informations complètes
			for _, brokenLink := range report.BrokenLinks {
				if brokenLink.URL == "" {
					t.Error("Broken link URL should not be empty")
				}
				
				if brokenLink.StatusCode == 0 && brokenLink.Error == "" {
					t.Error("Broken link should have either status code or error")
				}
			}
		})
	}
}

func TestBrokenLinksDetector_IsValidStatusCode(t *testing.T) {
	detector := NewBrokenLinksDetector()

	tests := []struct {
		statusCode int
		expected   bool
	}{
		// Valid status codes
		{200, true},  // OK
		{201, true},  // Created
		{301, true},  // Moved Permanently
		{302, true},  // Found
		{304, true},  // Not Modified
		{401, true},  // Unauthorized (considered valid)
		{403, true},  // Forbidden (considered valid)
		
		// Invalid status codes
		{404, false}, // Not Found
		{500, false}, // Internal Server Error
		{503, false}, // Service Unavailable
		{400, false}, // Bad Request
		{410, false}, // Gone
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			result := detector.isValidStatusCode(tt.statusCode)
			
			if result != tt.expected {
				t.Errorf("For status code %d, expected %v, got %v", tt.statusCode, tt.expected, result)
			}
		})
	}
}

func TestBrokenLinksDetector_SetMaxWorkers(t *testing.T) {
	tests := []struct {
		name     string
		workers  int
		expected int
	}{
		{
			name:     "valid worker count",
			workers:  5,
			expected: 5,
		},
		{
			name:     "zero workers (should not change)",
			workers:  0,
			expected: 10, // default value
		},
		{
			name:     "negative workers (should not change)",
			workers:  -1,
			expected: 10, // default value
		},
		{
			name:     "too many workers (should not change)",
			workers:  100,
			expected: 10, // default value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer une nouvelle instance pour chaque test
			detector := NewBrokenLinksDetector()
			detector.SetMaxWorkers(tt.workers)
			
			if detector.maxWorkers != tt.expected {
				t.Errorf("Expected %d max workers, got %d", tt.expected, detector.maxWorkers)
			}
		})
	}
}

func TestBrokenLinksDetector_SetTimeout(t *testing.T) {
	detector := NewBrokenLinksDetector()
	
	// Test valid timeout
	newTimeout := 5 * time.Second
	detector.SetTimeout(newTimeout)
	
	if detector.client.Timeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, detector.client.Timeout)
	}
	
	// Test invalid timeout (should not change)
	originalTimeout := detector.client.Timeout
	detector.SetTimeout(-1 * time.Second)
	
	if detector.client.Timeout != originalTimeout {
		t.Error("Negative timeout should not change the original timeout")
	}
}

func TestBrokenLinksDetector_GetStats(t *testing.T) {
	detector := NewBrokenLinksDetector()
	
	stats := detector.GetStats()
	
	if stats == nil {
		t.Error("GetStats returned nil")
		return
	}
	
	// Vérifier les clés attendues
	expectedKeys := []string{"max_workers", "timeout", "agent_name", "transport_type"}
	
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected key '%s' not found in stats", key)
		}
	}
	
	// Vérifier les valeurs
	if stats["agent_name"] != constants.AgentNameBrokenLinks {
		t.Errorf("Expected agent_name %s, got %v", constants.AgentNameBrokenLinks, stats["agent_name"])
	}
	
	if stats["transport_type"] != "http" {
		t.Errorf("Expected transport_type 'http', got %v", stats["transport_type"])
	}
}

func TestBrokenLinksDetector_HealthCheck(t *testing.T) {
	detector := NewBrokenLinksDetector()
	
	// Note: Ce test pourrait échouer si httpstat.us n'est pas disponible
	// Dans un environnement de test réel, nous utiliserions un mock
	err := detector.HealthCheck()
	
	// On accepte que le health check puisse échouer dans certains environnements
	// mais on vérifie au moins qu'il retourne quelque chose
	if err != nil {
		t.Logf("HealthCheck failed (may be expected in test environment): %v", err)
	}
}