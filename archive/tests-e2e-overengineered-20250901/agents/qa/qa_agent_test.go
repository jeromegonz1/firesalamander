package qa

import (
	"os"
	"testing"
	"time"
)

func TestQAAgentCreation(t *testing.T) {
	config := DefaultQAConfig()
	agent := NewQAAgent(config)

	if agent == nil {
		t.Fatal("QA Agent should not be nil")
	}

	if agent.config.MinCoverage != 80.0 {
		t.Errorf("Expected min coverage 80.0, got %f", agent.config.MinCoverage)
	}

	if !agent.config.EnableVet {
		t.Error("Go vet should be enabled by default")
	}
}

func TestDefaultQAConfig(t *testing.T) {
	config := DefaultQAConfig()

	if config.MinCoverage != 80.0 {
		t.Errorf("Expected min coverage 80.0, got %f", config.MinCoverage)
	}

	if config.OutputFormat != "json" {
		t.Errorf("Expected output format json, got %s", config.OutputFormat)
	}

	if !config.EnableVet || !config.EnableLint {
		t.Error("Basic tools should be enabled by default")
	}
}

func TestCalculateOverallScore(t *testing.T) {
	agent := NewQAAgent(DefaultQAConfig())
	
	// Mock des stats de test
	agent.stats = &QAStats{
		Timestamp: time.Now(),
		Coverage: CoverageStats{
			TotalCoverage: 85.0,
			Threshold:     80.0,
			Passed:        true,
		},
		TestResults: TestResults{
			TotalTests:  10,
			PassedTests: 10,
			FailedTests: 0,
		},
		VetIssues:        []VetIssue{},
		LintIssues:       []LintIssue{},
		SecurityIssues:   []SecurityIssue{},
		ComplexityIssues: []ComplexityIssue{},
	}

	agent.calculateOverallScore()

	if agent.stats.OverallScore != 100.0 {
		t.Errorf("Expected perfect score 100.0, got %f", agent.stats.OverallScore)
	}

	if agent.stats.Status != "excellent" {
		t.Errorf("Expected status excellent, got %s", agent.stats.Status)
	}
}

func TestParseCoverageResults(t *testing.T) {
	agent := NewQAAgent(DefaultQAConfig())
	agent.stats = &QAStats{Coverage: CoverageStats{}}

	mockOutput := `firesalamander/internal/agents/crawler/crawler.go	func New	100.0%
firesalamander/internal/agents/crawler/fetcher.go	func Fetch	85.5%
total:							(statements)		90.2%`

	agent.parseCoverageResults(mockOutput)

	if agent.stats.Coverage.TotalCoverage != 90.2 {
		t.Errorf("Expected total coverage 90.2, got %f", agent.stats.Coverage.TotalCoverage)
	}

	if len(agent.stats.Coverage.PackagesCov) == 0 {
		t.Error("Package coverage should be parsed")
	}
}

func TestParseVetResults(t *testing.T) {
	agent := NewQAAgent(DefaultQAConfig())
	agent.stats = &QAStats{}

	mockOutput := `./crawler/crawler.go:123:45: unreachable code
./main.go:67:12: missing return statement`

	agent.parseVetResults(mockOutput)

	if len(agent.stats.VetIssues) != 2 {
		t.Errorf("Expected 2 vet issues, got %d", len(agent.stats.VetIssues))
	}

	firstIssue := agent.stats.VetIssues[0]
	if firstIssue.File != "./crawler/crawler.go" {
		t.Errorf("Expected file ./crawler/crawler.go, got %s", firstIssue.File)
	}
	if firstIssue.Line != 123 {
		t.Errorf("Expected line 123, got %d", firstIssue.Line)
	}
	if firstIssue.Column != 45 {
		t.Errorf("Expected column 45, got %d", firstIssue.Column)
	}
}

func TestGenerateReport(t *testing.T) {
	agent := NewQAAgent(DefaultQAConfig())
	agent.config.ReportPath = "/tmp/test_qa_reports"
	
	agent.stats = &QAStats{
		Timestamp:    time.Now(),
		OverallScore: 85.5,
		Status:       "good",
		Coverage: CoverageStats{
			TotalCoverage: 85.0,
			Passed:        true,
		},
	}

	err := agent.generateReport()
	if err != nil {
		t.Fatalf("Failed to generate report: %v", err)
	}

	// Vérifier que le fichier existe
	reportFile := "/tmp/test_qa_reports/qa_report.json"
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		t.Error("Report file should exist")
	}

	// Nettoyage
	os.RemoveAll("/tmp/test_qa_reports")
}

func TestScoreCalculationWithIssues(t *testing.T) {
	agent := NewQAAgent(DefaultQAConfig())
	
	// Simuler des problèmes
	agent.stats = &QAStats{
		Coverage: CoverageStats{
			TotalCoverage: 70.0, // Sous le seuil de 80%
			Threshold:     80.0,
		},
		TestResults: TestResults{
			TotalTests:  10,
			PassedTests: 8,
			FailedTests: 2, // 20% d'échec
		},
		VetIssues: []VetIssue{
			{Message: "Issue 1"},
			{Message: "Issue 2"},
		},
		LintIssues: []LintIssue{
			{Message: "Lint issue 1", Severity: "warning"},
		},
		SecurityIssues: []SecurityIssue{
			{Severity: "HIGH", Message: "Security issue"},
		},
		ComplexityIssues: []ComplexityIssue{
			{Complexity: 15, Threshold: 10},
		},
	}

	agent.calculateOverallScore()

	// Score devrait être réduit à cause des problèmes
	expectedReductions := 3.0 + 4.0 + 4.0 + 0.5 + 5.0 + 0.2 // Coverage + Tests + Vet + Lint + Security + Complexity
	expectedScore := 100.0 - expectedReductions
	
	t.Logf("Expected score after reductions: %f", expectedScore)

	if agent.stats.OverallScore < 80.0 {
		t.Logf("Score correctly reduced to %f due to issues", agent.stats.OverallScore)
	}

	if agent.stats.Status == "excellent" {
		t.Error("Status should not be excellent with many issues")
	}
}