package qa

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"firesalamander/internal/logger"
)

var log = logger.New("QA-AGENT")

// QAAgent gère la qualité du code Go
type QAAgent struct {
	config *QAConfig
	stats  *QAStats
}

// QAConfig contient la configuration du QA Agent
type QAConfig struct {
	MinCoverage      float64 `json:"min_coverage"`
	EnableVet        bool    `json:"enable_vet"`
	EnableLint       bool    `json:"enable_lint"`
	EnableSecurity   bool    `json:"enable_security"`
	EnableComplexity bool    `json:"enable_complexity"`
	OutputFormat     string  `json:"output_format"` // json, text, html
	ReportPath       string  `json:"report_path"`
}

// QAStats contient les statistiques d'analyse
type QAStats struct {
	Timestamp        time.Time      `json:"timestamp"`
	Coverage         CoverageStats  `json:"coverage"`
	VetIssues        []VetIssue     `json:"vet_issues"`
	LintIssues       []LintIssue    `json:"lint_issues"`
	SecurityIssues   []SecurityIssue `json:"security_issues"`
	ComplexityIssues []ComplexityIssue `json:"complexity_issues"`
	TestResults      TestResults    `json:"test_results"`
	OverallScore     float64        `json:"overall_score"`
	Status           string         `json:"status"` // pass, fail, warning
}

// CoverageStats représente les statistiques de coverage
type CoverageStats struct {
	TotalCoverage float64            `json:"total_coverage"`
	PackagesCov   map[string]float64 `json:"packages_coverage"`
	Threshold     float64            `json:"threshold"`
	Passed        bool               `json:"passed"`
}

// VetIssue représente un problème détecté par go vet
type VetIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Category string `json:"category"`
}

// LintIssue représente un problème de style/lint
type LintIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
}

// SecurityIssue représente un problème de sécurité
type SecurityIssue struct {
	File        string `json:"file"`
	Line        int    `json:"line"`
	Rule        string `json:"rule"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Confidence  string `json:"confidence"`
}

// ComplexityIssue représente un problème de complexité
type ComplexityIssue struct {
	Function   string  `json:"function"`
	File       string  `json:"file"`
	Line       int     `json:"line"`
	Complexity int     `json:"complexity"`
	Threshold  int     `json:"threshold"`
}

// TestResults contient les résultats des tests
type TestResults struct {
	TotalTests  int           `json:"total_tests"`
	PassedTests int           `json:"passed_tests"`
	FailedTests int           `json:"failed_tests"`
	Duration    time.Duration `json:"duration"`
	Packages    []PackageTest `json:"packages"`
}

// PackageTest représente les résultats d'un package
type PackageTest struct {
	Package string        `json:"package"`
	Status  string        `json:"status"`
	Time    time.Duration `json:"time"`
	Tests   []TestCase    `json:"tests"`
}

// TestCase représente un test unitaire
type TestCase struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Output   string        `json:"output,omitempty"`
}

// DefaultQAConfig retourne une configuration par défaut
func DefaultQAConfig() *QAConfig {
	return &QAConfig{
		MinCoverage:      80.0,
		EnableVet:        true,
		EnableLint:       true,
		EnableSecurity:   true,
		EnableComplexity: true,
		OutputFormat:     "json",
		ReportPath:       "tests/reports/qa",
	}
}

// NewQAAgent crée une nouvelle instance du QA Agent
func NewQAAgent(config *QAConfig) *QAAgent {
	if config == nil {
		config = DefaultQAConfig()
	}

	log.Info("QA Agent initialized", map[string]interface{}{
		"min_coverage": config.MinCoverage,
		"tools":        getEnabledTools(config),
	})

	return &QAAgent{
		config: config,
		stats:  &QAStats{Timestamp: time.Now()},
	}
}

// RunFullAnalysis exécute une analyse complète de qualité
func (qa *QAAgent) RunFullAnalysis() (*QAStats, error) {
	log.Info("🔍 Starting full QA analysis")
	
	qa.stats = &QAStats{
		Timestamp: time.Now(),
		Coverage:  CoverageStats{Threshold: qa.config.MinCoverage},
	}

	// 1. Tests unitaires
	log.Debug("Running unit tests")
	if err := qa.runUnitTests(); err != nil {
		log.Error("Unit tests failed", map[string]interface{}{"error": err.Error()})
	}

	// 2. Coverage analysis
	log.Debug("Analyzing test coverage")
	if err := qa.analyzeCoverage(); err != nil {
		log.Error("Coverage analysis failed", map[string]interface{}{"error": err.Error()})
	}

	// 3. Go vet
	if qa.config.EnableVet {
		log.Debug("Running go vet")
		if err := qa.runGoVet(); err != nil {
			log.Warn("Go vet issues found", map[string]interface{}{"error": err.Error()})
		}
	}

	// 4. Linting
	if qa.config.EnableLint {
		log.Debug("Running golangci-lint")
		if err := qa.runLinting(); err != nil {
			log.Warn("Linting issues found", map[string]interface{}{"error": err.Error()})
		}
	}

	// 5. Security analysis
	if qa.config.EnableSecurity {
		log.Debug("Running security analysis")
		if err := qa.runSecurityAnalysis(); err != nil {
			log.Warn("Security analysis issues", map[string]interface{}{"error": err.Error()})
		}
	}

	// 6. Complexity analysis
	if qa.config.EnableComplexity {
		log.Debug("Analyzing code complexity")
		if err := qa.analyzeComplexity(); err != nil {
			log.Warn("Complexity analysis issues", map[string]interface{}{"error": err.Error()})
		}
	}

	// Calculer le score global
	qa.calculateOverallScore()

	// Générer le rapport
	if err := qa.generateReport(); err != nil {
		log.Error("Failed to generate report", map[string]interface{}{"error": err.Error()})
	}

	log.Info("QA analysis completed", map[string]interface{}{
		"overall_score": qa.stats.OverallScore,
		"status":        qa.stats.Status,
		"coverage":      qa.stats.Coverage.TotalCoverage,
	})

	return qa.stats, nil
}

// runUnitTests exécute les tests unitaires
func (qa *QAAgent) runUnitTests() error {
	cmd := exec.Command("go", "test", "-v", "-json", "./...")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run tests: %w", err)
	}

	// Parser la sortie JSON
	qa.parseTestResults(output)
	return nil
}

// analyzeCoverage analyse la couverture de code
func (qa *QAAgent) analyzeCoverage() error {
	// Générer le profil de coverage
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate coverage: %w", err)
	}

	// Analyser le coverage par package
	cmd = exec.Command("go", "tool", "cover", "-func=coverage.out")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to analyze coverage: %w", err)
	}

	qa.parseCoverageResults(string(output))
	
	// Nettoyage
	os.Remove("coverage.out")
	
	return nil
}

// runGoVet exécute go vet
func (qa *QAAgent) runGoVet() error {
	cmd := exec.Command("go", "vet", "./...")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		qa.parseVetResults(string(output))
		return fmt.Errorf("go vet found issues")
	}
	
	return nil
}

// runLinting exécute golangci-lint
func (qa *QAAgent) runLinting() error {
	// Vérifier si golangci-lint est installé
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		log.Warn("golangci-lint not found, skipping lint analysis")
		return nil
	}

	cmd := exec.Command("golangci-lint", "run", "--out-format", "json")
	output, err := cmd.Output()
	
	if err != nil {
		qa.parseLintResults(output)
		return fmt.Errorf("linting issues found")
	}
	
	return nil
}

// runSecurityAnalysis exécute l'analyse de sécurité
func (qa *QAAgent) runSecurityAnalysis() error {
	// Utiliser gosec si disponible
	if _, err := exec.LookPath("gosec"); err != nil {
		log.Warn("gosec not found, skipping security analysis")
		return nil
	}

	cmd := exec.Command("gosec", "-fmt", "json", "./...")
	output, err := cmd.Output()
	
	if err != nil {
		qa.parseSecurityResults(output)
		return fmt.Errorf("security issues found")
	}
	
	return nil
}

// analyzeComplexity analyse la complexité cyclomatique
func (qa *QAAgent) analyzeComplexity() error {
	// Utiliser gocyclo si disponible
	if _, err := exec.LookPath("gocyclo"); err != nil {
		log.Warn("gocyclo not found, skipping complexity analysis")
		return nil
	}

	cmd := exec.Command("gocyclo", "-over", "10", ".")
	output, err := cmd.Output()
	
	if err == nil && len(output) > 0 {
		qa.parseComplexityResults(string(output))
		return fmt.Errorf("high complexity functions found")
	}
	
	return nil
}

// Parser methods
func (qa *QAAgent) parseTestResults(output []byte) {
	// Parser la sortie JSON des tests Go
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		var result map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
			continue
		}
		
		action, ok := result["Action"].(string)
		if !ok {
			continue
		}
		
		switch action {
		case "pass", "fail":
			qa.stats.TestResults.TotalTests++
			if action == "pass" {
				qa.stats.TestResults.PassedTests++
			} else {
				qa.stats.TestResults.FailedTests++
			}
		}
	}
}

func (qa *QAAgent) parseCoverageResults(output string) {
	lines := strings.Split(output, "\n")
	qa.stats.Coverage.PackagesCov = make(map[string]float64)
	
	for _, line := range lines {
		if strings.Contains(line, "total:") {
			// Extraire le pourcentage total
			re := regexp.MustCompile(`(\d+\.\d+)%`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
					qa.stats.Coverage.TotalCoverage = coverage
				}
			}
		} else if strings.Contains(line, ".go") && strings.Contains(line, "%") {
			// Coverage par fichier/package
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				pkg := parts[0]
				covStr := strings.TrimSuffix(parts[len(parts)-1], "%")
				if coverage, err := strconv.ParseFloat(covStr, 64); err == nil {
					qa.stats.Coverage.PackagesCov[pkg] = coverage
				}
			}
		}
	}
	
	qa.stats.Coverage.Passed = qa.stats.Coverage.TotalCoverage >= qa.config.MinCoverage
}

func (qa *QAAgent) parseVetResults(output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Parser format: ./file.go:line:col: message
		parts := strings.SplitN(line, ":", 4)
		if len(parts) >= 4 {
			lineNum, _ := strconv.Atoi(parts[1])
			colNum, _ := strconv.Atoi(parts[2])
			
			issue := VetIssue{
				File:     parts[0],
				Line:     lineNum,
				Column:   colNum,
				Message:  strings.TrimSpace(parts[3]),
				Category: "vet",
			}
			qa.stats.VetIssues = append(qa.stats.VetIssues, issue)
		}
	}
}

func (qa *QAAgent) parseLintResults(output []byte) {
	var result struct {
		Issues []struct {
			FromLinter  string `json:"FromLinter"`
			Text        string `json:"Text"`
			Pos         struct {
				Filename string `json:"Filename"`
				Line     int    `json:"Line"`
				Column   int    `json:"Column"`
			} `json:"Pos"`
			Severity string `json:"Severity"`
		} `json:"Issues"`
	}
	
	if err := json.Unmarshal(output, &result); err != nil {
		return
	}
	
	for _, issue := range result.Issues {
		lintIssue := LintIssue{
			File:     issue.Pos.Filename,
			Line:     issue.Pos.Line,
			Column:   issue.Pos.Column,
			Message:  issue.Text,
			Rule:     issue.FromLinter,
			Severity: issue.Severity,
		}
		qa.stats.LintIssues = append(qa.stats.LintIssues, lintIssue)
	}
}

func (qa *QAAgent) parseSecurityResults(output []byte) {
	var result struct {
		Issues []struct {
			Rule        string `json:"rule"`
			Details     string `json:"details"`
			File        string `json:"file"`
			Line        string `json:"line"`
			Confidence  string `json:"confidence"`
			Severity    string `json:"severity"`
		} `json:"Issues"`
	}
	
	if err := json.Unmarshal(output, &result); err != nil {
		return
	}
	
	for _, issue := range result.Issues {
		lineNum, _ := strconv.Atoi(issue.Line)
		secIssue := SecurityIssue{
			File:       issue.File,
			Line:       lineNum,
			Rule:       issue.Rule,
			Severity:   issue.Severity,
			Message:    issue.Details,
			Confidence: issue.Confidence,
		}
		qa.stats.SecurityIssues = append(qa.stats.SecurityIssues, secIssue)
	}
}

func (qa *QAAgent) parseComplexityResults(output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Format: complexity function file:line
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			complexity, _ := strconv.Atoi(parts[0])
			function := parts[1]
			fileLine := parts[2]
			
			fileLineParts := strings.Split(fileLine, ":")
			if len(fileLineParts) >= 2 {
				file := fileLineParts[0]
				lineNum, _ := strconv.Atoi(fileLineParts[1])
				
				issue := ComplexityIssue{
					Function:   function,
					File:       file,
					Line:       lineNum,
					Complexity: complexity,
					Threshold:  10,
				}
				qa.stats.ComplexityIssues = append(qa.stats.ComplexityIssues, issue)
			}
		}
	}
}

// calculateOverallScore calcule le score global de qualité
func (qa *QAAgent) calculateOverallScore() {
	score := 100.0
	
	// Coverage (30% du score)
	if qa.stats.Coverage.TotalCoverage < qa.config.MinCoverage {
		coveragePenalty := (qa.config.MinCoverage - qa.stats.Coverage.TotalCoverage) * 0.3
		score -= coveragePenalty
	}
	
	// Tests (20% du score)
	if qa.stats.TestResults.TotalTests > 0 {
		failureRate := float64(qa.stats.TestResults.FailedTests) / float64(qa.stats.TestResults.TotalTests)
		score -= failureRate * 20
	}
	
	// Vet issues (20% du score)
	score -= float64(len(qa.stats.VetIssues)) * 2
	
	// Lint issues (15% du score)
	score -= float64(len(qa.stats.LintIssues)) * 0.5
	
	// Security issues (10% du score)
	highSecIssues := 0
	for _, issue := range qa.stats.SecurityIssues {
		if issue.Severity == "HIGH" {
			highSecIssues++
		}
	}
	score -= float64(highSecIssues) * 5
	
	// Complexity issues (5% du score)
	score -= float64(len(qa.stats.ComplexityIssues)) * 0.2
	
	if score < 0 {
		score = 0
	}
	
	qa.stats.OverallScore = score
	
	// Déterminer le statut
	switch {
	case score >= 90:
		qa.stats.Status = "excellent"
	case score >= 80:
		qa.stats.Status = "good"
	case score >= 70:
		qa.stats.Status = "acceptable"
	case score >= 60:
		qa.stats.Status = "needs_improvement"
	default:
		qa.stats.Status = "poor"
	}
}

// generateReport génère le rapport de qualité
func (qa *QAAgent) generateReport() error {
	// Créer le répertoire de rapport
	if err := os.MkdirAll(qa.config.ReportPath, 0755); err != nil {
		return err
	}
	
	// Générer le rapport JSON
	reportFile := filepath.Join(qa.config.ReportPath, "qa_report.json")
	data, err := json.MarshalIndent(qa.stats, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(reportFile, data, 0644)
}

// GetStats retourne les statistiques actuelles
func (qa *QAAgent) GetStats() *QAStats {
	return qa.stats
}

// Helper functions
func getEnabledTools(config *QAConfig) []string {
	var tools []string
	if config.EnableVet {
		tools = append(tools, "go vet")
	}
	if config.EnableLint {
		tools = append(tools, "golangci-lint")
	}
	if config.EnableSecurity {
		tools = append(tools, "gosec")
	}
	if config.EnableComplexity {
		tools = append(tools, "gocyclo")
	}
	return tools
}