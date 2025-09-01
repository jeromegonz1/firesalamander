package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - LOGGING TESTS
// TDD + Zero Hardcoding Policy
// ========================================

// TestLoggerCreation teste la création d'un logger
func TestLoggerCreation(t *testing.T) {
	config := &LogConfig{
		Level:         DEBUG,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

// TestLoggerWithNilConfig teste la création avec config nil
func TestLoggerWithNilConfig(t *testing.T) {
	logger, err := NewLogger(nil)
	if err == nil {
		t.Fatal("Should return error for nil config")
	}
	
	if logger != nil {
		t.Fatal("Logger should be nil for invalid config")
	}
}

// TestLogLevels teste les différents niveaux de log
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         DEBUG,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// Rediriger la sortie vers notre buffer pour les tests
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	// Tester tous les niveaux
	logger.Debug(constants.LogCategorySystem, "Debug message")
	logger.Info(constants.LogCategorySystem, "Info message")
	logger.Warn(constants.LogCategorySystem, "Warning message")
	logger.Error(constants.LogCategorySystem, "Error message", fmt.Errorf("test error"))
	
	output := buf.String()
	
	// Vérifier que tous les messages sont présents
	if !strings.Contains(output, "Debug message") {
		t.Error("Debug message not found in output")
	}
	if !strings.Contains(output, "Info message") {
		t.Error("Info message not found in output")
	}
	if !strings.Contains(output, "Warning message") {
		t.Error("Warning message not found in output")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Error message not found in output")
	}
}

// TestLogLevelFiltering teste le filtrage par niveau
func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         WARN, // Seulement WARN et plus
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	logger.Debug(constants.LogCategorySystem, "Debug message") // Ne devrait pas apparaître
	logger.Info(constants.LogCategorySystem, "Info message")   // Ne devrait pas apparaître
	logger.Warn(constants.LogCategorySystem, "Warning message") // Devrait apparaître
	logger.Error(constants.LogCategorySystem, "Error message", fmt.Errorf("test error")) // Devrait apparaître
	
	output := buf.String()
	
	// DEBUG et INFO ne devraient pas être présents
	if strings.Contains(output, "Debug message") {
		t.Error("Debug message should be filtered out")
	}
	if strings.Contains(output, "Info message") {
		t.Error("Info message should be filtered out")
	}
	
	// WARN et ERROR devraient être présents
	if !strings.Contains(output, "Warning message") {
		t.Error("Warning message should be present")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Error message should be present")
	}
}

// TestJSONFormat teste le format JSON
func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	testData := map[string]interface{}{
		"user_id": "123",
		"action":  "login",
	}
	
	logger.Info(constants.LogCategorySystem, "Test message", testData)
	
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	if len(lines) == 0 {
		t.Fatal("No output generated")
	}
	
	// Vérifier que c'est du JSON valide
	var logEntry LogEntry
	if err := json.Unmarshal([]byte(lines[0]), &logEntry); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
	
	// Vérifier les champs
	if logEntry.Level != constants.LogLevelInfo {
		t.Errorf("Expected level %s, got %s", constants.LogLevelInfo, logEntry.Level)
	}
	if logEntry.Category != constants.LogCategorySystem {
		t.Errorf("Expected category %s, got %s", constants.LogCategorySystem, logEntry.Category)
	}
	if logEntry.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", logEntry.Message)
	}
}

// TestContextualLogging teste le logging avec contexte
func TestContextualLogging(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	// Test avec Request ID
	contextLogger := logger.WithRequestID("req-123")
	contextLogger.Info(constants.LogCategoryAPI, "API request")
	
	output := buf.String()
	
	if !strings.Contains(output, "req-123") {
		t.Error("Request ID should be present in log output")
	}
}

// TestHTTPLogger teste le logger HTTP
func TestHTTPLogger(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	httpLogger := logger.HTTP()
	httpLogger.RequestReceived("GET", "/api/test", "127.0.0.1", "test-agent", "req-123")
	
	output := buf.String()
	
	// Vérifier que les informations HTTP sont présentes
	if !strings.Contains(output, "GET") {
		t.Error("HTTP method should be present")
	}
	if !strings.Contains(output, "/api/test") {
		t.Error("URL should be present")
	}
	if !strings.Contains(output, "127.0.0.1") {
		t.Error("Remote address should be present")
	}
}

// TestAPILogger teste le logger API
func TestAPILogger(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	apiLogger := logger.API()
	
	// Test requête reçue
	requestData := map[string]interface{}{
		"url": "https://example.com",
	}
	apiLogger.RequestReceived("/api/analyze", "POST", "req-456", requestData)
	
	// Test réponse envoyée
	apiLogger.ResponseSent("/api/analyze", "POST", 200, 150, "req-456")
	
	output := buf.String()
	
	// Vérifier le contenu
	if !strings.Contains(output, "/api/analyze") {
		t.Error("API endpoint should be present")
	}
	if !strings.Contains(output, "POST") {
		t.Error("HTTP method should be present")
	}
	if !strings.Contains(output, "req-456") {
		t.Error("Request ID should be present")
	}
}

// TestPerformanceLogger teste le logger de performance
func TestPerformanceLogger(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         DEBUG,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	perfLogger := logger.Performance()
	
	// Test début d'opération
	perfLogger.OperationStarted("test_operation", map[string]interface{}{
		"param": "value",
	})
	
	// Test fin d'opération
	perfLogger.OperationCompleted("test_operation", 500, map[string]interface{}{
		"result": "success",
	})
	
	output := buf.String()
	
	// Vérifier le contenu
	if !strings.Contains(output, "test_operation") {
		t.Error("Operation name should be present")
	}
	if !strings.Contains(output, "started") {
		t.Error("Operation started status should be present")
	}
	if !strings.Contains(output, "completed") {
		t.Error("Operation completed status should be present")
	}
}

// TestAuditLogger teste le logger d'audit
func TestAuditLogger(t *testing.T) {
	var buf bytes.Buffer
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = &buf
	
	auditLogger := logger.Audit()
	
	// Test action démarrée
	beforeState := map[string]interface{}{
		"status": "inactive",
	}
	auditLogger.ActionStarted("user_login", "user:123", "user-123", "session-456", beforeState)
	
	// Test action terminée
	afterState := map[string]interface{}{
		"status": "active",
	}
	auditLogger.ActionCompleted("user_login", "user:123", "user-123", "session-456", true, afterState)
	
	output := buf.String()
	
	// Vérifier le contenu
	if !strings.Contains(output, "user_login") {
		t.Error("Action name should be present")
	}
	if !strings.Contains(output, "user-123") {
		t.Error("User ID should be present")
	}
	if !strings.Contains(output, "session-456") {
		t.Error("Session ID should be present")
	}
}

// TestLoggerClose teste la fermeture du logger
func TestLoggerClose(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir := "/tmp/fire-salamander-log-test"
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)
	
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: false,
		EnableFile:    true,
		LogDir:        tmpDir,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// Écrire quelques logs
	logger.Info(constants.LogCategorySystem, "Test message")
	
	// Fermer le logger
	err = logger.Close()
	if err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}
	
	// Vérifier que les fichiers de log ont été créés
	logFiles := []string{
		constants.LogFileAccess,
		constants.LogFileError,
		constants.LogFileDebug,
		constants.LogFileAudit,
		constants.LogFilePerformance,
		constants.LogFileSystem,
	}
	
	for _, fileName := range logFiles {
		filePath := fmt.Sprintf("%s/%s", tmpDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Log file %s was not created", fileName)
		}
	}
}

// BenchmarkJSONLogging benchmark pour le logging JSON
func BenchmarkJSONLogging(b *testing.B) {
	config := &LogConfig{
		Level:         INFO,
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	
	// Rediriger vers discard pour éviter la sortie pendant les benchmarks
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = io.Discard
	
	testData := map[string]interface{}{
		"user_id":    "123",
		"action":     "test",
		"timestamp":  time.Now(),
		"ip_address": "127.0.0.1",
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(constants.LogCategorySystem, "Benchmark message", testData)
		}
	})
}

// BenchmarkTextLogging benchmark pour le logging texte
func BenchmarkTextLogging(b *testing.B) {
	config := &LogConfig{
		Level:         INFO,
		Format:        "text",
		EnableConsole: true,
		EnableFile:    false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	
	fsLogger := logger.(*FireSalamanderLogger)
	fsLogger.writers["console"] = io.Discard
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(constants.LogCategorySystem, "Benchmark message")
		}
	})
}