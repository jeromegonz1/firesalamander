package report

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"firesalamander/internal/audit"
	"firesalamander/internal/crawler"
	"firesalamander/internal/semantic"
)

// ReportEngine handles report generation in multiple formats
type ReportEngine struct {
	htmlTemplate *template.Template
}

// AuditResults contains all data needed for report generation
type AuditResults struct {
	AuditID         string                    `json:"audit_id"`
	SiteURL         string                    `json:"site_url"`
	StartedAt       string                    `json:"started_at"`
	Duration        string                    `json:"duration"`
	TotalPages      int                       `json:"total_pages"`
	CrawlData       crawler.CrawlResult       `json:"crawl_data"`
	TechResults     audit.TechResult          `json:"tech_results"`
	SemanticResults semantic.SemanticResult   `json:"semantic_results"`
}

// TemplateData represents data passed to HTML template
type TemplateData struct {
	AuditID        string           `json:"audit_id"`
	SiteURL        string           `json:"site_url"`
	GeneratedAt    string           `json:"generated_at"`
	TotalPages     int              `json:"total_pages"`
	Duration       string           `json:"duration"`
	CriticalIssues int              `json:"critical_issues"`
	HighIssues     int              `json:"high_issues"`
	MediumIssues   int              `json:"medium_issues"`
	LowIssues      int              `json:"low_issues"`
	OverallScore   float64          `json:"overall_score"`
	Pages          []PageSummary    `json:"pages"`
	Issues         []IssueSummary   `json:"issues"`
	Keywords       []KeywordSummary `json:"keywords"`
	Topics         []TopicSummary   `json:"topics"`
}

// PageSummary represents a page in the report
type PageSummary struct {
	URL              string  `json:"url"`
	Title            string  `json:"title"`
	H1               string  `json:"h1"`
	IssuesCount      int     `json:"issues_count"`
	PerformanceScore float64 `json:"performance_score"`
	Depth            int     `json:"depth"`
}

// IssueSummary represents an SEO issue in the report
type IssueSummary struct {
	ID       string   `json:"id"`
	Severity string   `json:"severity"`
	Message  string   `json:"message"`
	Count    int      `json:"count"`
	Pages    []string `json:"pages"`
}

// KeywordSummary represents a keyword suggestion in the report
type KeywordSummary struct {
	Keyword    string  `json:"keyword"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
	Evidence   string  `json:"evidence"`
}

// TopicSummary represents a semantic topic in the report
type TopicSummary struct {
	Label string   `json:"label"`
	Terms []string `json:"terms"`
}

// NewReportEngine creates a new report engine
func NewReportEngine() *ReportEngine {
	// Add custom template functions
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"printf": fmt.Sprintf,
	}
	
	tmpl := template.New("report").Funcs(funcMap)
	
	return &ReportEngine{
		htmlTemplate: template.Must(tmpl.Parse(htmlTemplateContent)),
	}
}

// GenerateHTML generates an HTML report
func (re *ReportEngine) GenerateHTML(results AuditResults) (string, error) {
	if err := re.validateAuditResults(results); err != nil {
		return "", err
	}

	templateData := re.prepareTemplateData(results)
	return re.renderHTMLTemplate(templateData)
}

// GenerateJSON generates a JSON report
func (re *ReportEngine) GenerateJSON(results AuditResults) (string, error) {
	if err := re.validateAuditResults(results); err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}

// GenerateCSV generates a CSV report with page data
func (re *ReportEngine) GenerateCSV(results AuditResults) (string, error) {
	if err := re.validateAuditResults(results); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	headers := []string{"URL", "Title", "H1", "Depth", "Issues", "Performance", "Accessibility", "SEO"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Data rows
	for _, page := range results.CrawlData.Pages {
		issuesCount := re.countPageIssues(page.URL, results.TechResults.Findings)
		
		row := []string{
			page.URL,
			page.Title,
			page.H1,
			strconv.Itoa(page.Depth),
			strconv.Itoa(issuesCount),
			"N/A", // Performance score per page not available in current structure
			"N/A", // Accessibility score per page not available
			"N/A", // SEO score per page not available
		}
		
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}

// SaveReport saves a report to disk in the specified format
func (re *ReportEngine) SaveReport(results AuditResults, format string, outputDir string) (string, error) {
	if err := re.validateAuditResults(results); err != nil {
		return "", err
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("fire-salamander-audit-%s-%s.%s", results.AuditID, timestamp, format)
	filepath := filepath.Join(outputDir, filename)

	// Generate content based on format
	var content string
	var err error

	switch strings.ToLower(format) {
	case "html":
		content, err = re.GenerateHTML(results)
	case "json":
		content, err = re.GenerateJSON(results)
	case "csv":
		content, err = re.GenerateCSV(results)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate %s report: %w", format, err)
	}

	// Write to file
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write report file: %w", err)
	}

	return filepath, nil
}

// prepareTemplateData converts audit results to template data
func (re *ReportEngine) prepareTemplateData(results AuditResults) TemplateData {
	now := time.Now().Format("2006-01-02 15:04:05")
	
	// Count issues by severity
	criticalCount, highCount, mediumCount, lowCount := re.countIssuesBySeverity(results.TechResults.Findings)
	
	// Calculate overall score (average of Lighthouse scores)
	scores := results.TechResults.Scores
	overallScore := (scores.Performance + scores.Accessibility + scores.BestPractices + scores.SEO) / 4.0

	// Prepare page summaries
	pages := make([]PageSummary, len(results.CrawlData.Pages))
	for i, page := range results.CrawlData.Pages {
		issuesCount := re.countPageIssues(page.URL, results.TechResults.Findings)
		pages[i] = PageSummary{
			URL:              page.URL,
			Title:            page.Title,
			H1:               page.H1,
			IssuesCount:      issuesCount,
			PerformanceScore: scores.Performance, // Global score for now
			Depth:            page.Depth,
		}
	}

	// Prepare issue summaries
	issues := re.groupIssuesByType(results.TechResults.Findings)

	// Prepare keyword summaries
	keywords := make([]KeywordSummary, len(results.SemanticResults.Suggestions))
	for i, suggestion := range results.SemanticResults.Suggestions {
		evidence := ""
		if len(suggestion.Evidence) > 0 {
			evidence = suggestion.Evidence[0] // First evidence
		}
		
		keywords[i] = KeywordSummary{
			Keyword:    suggestion.Keyword,
			Confidence: suggestion.Confidence,
			Reason:     suggestion.Reason,
			Evidence:   evidence,
		}
	}

	// Prepare topic summaries
	topics := make([]TopicSummary, len(results.SemanticResults.Topics))
	for i, topic := range results.SemanticResults.Topics {
		topics[i] = TopicSummary{
			Label: topic.Label,
			Terms: topic.Terms,
		}
	}

	return TemplateData{
		AuditID:        results.AuditID,
		SiteURL:        results.SiteURL,
		GeneratedAt:    now,
		TotalPages:     results.TotalPages,
		Duration:       results.Duration,
		CriticalIssues: criticalCount,
		HighIssues:     highCount,
		MediumIssues:   mediumCount,
		LowIssues:      lowCount,
		OverallScore:   overallScore,
		Pages:          pages,
		Issues:         issues,
		Keywords:       keywords,
		Topics:         topics,
	}
}

// renderHTMLTemplate renders the HTML template
func (re *ReportEngine) renderHTMLTemplate(data TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := re.htmlTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// validateAuditResults validates that audit results are complete
func (re *ReportEngine) validateAuditResults(results AuditResults) error {
	if results.AuditID == "" {
		return fmt.Errorf("audit_id is required")
	}
	if results.SiteURL == "" {
		return fmt.Errorf("site_url is required")
	}
	return nil
}

// countIssuesBySeverity counts issues by severity level
func (re *ReportEngine) countIssuesBySeverity(findings []audit.Finding) (critical, high, medium, low int) {
	for _, finding := range findings {
		switch finding.Severity {
		case "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		case "low":
			low++
		}
	}
	return
}

// countPageIssues counts issues for a specific page
func (re *ReportEngine) countPageIssues(pageURL string, findings []audit.Finding) int {
	count := 0
	for _, finding := range findings {
		for _, evidence := range finding.Evidence {
			if evidence == pageURL {
				count++
				break
			}
		}
	}
	return count
}

// groupIssuesByType groups issues by their ID/type
func (re *ReportEngine) groupIssuesByType(findings []audit.Finding) []IssueSummary {
	issueMap := make(map[string]IssueSummary)

	for _, finding := range findings {
		if existing, exists := issueMap[finding.ID]; exists {
			existing.Count++
			existing.Pages = append(existing.Pages, finding.Evidence...)
			issueMap[finding.ID] = existing
		} else {
			issueMap[finding.ID] = IssueSummary{
				ID:       finding.ID,
				Severity: finding.Severity,
				Message:  finding.Message,
				Count:    1,
				Pages:    append([]string{}, finding.Evidence...),
			}
		}
	}

	// Convert map to slice
	issues := make([]IssueSummary, 0, len(issueMap))
	for _, issue := range issueMap {
		issues = append(issues, issue)
	}

	return issues
}

// HTML template content
const htmlTemplateContent = `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Rapport d'audit SEO - {{.AuditID}}</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background-color: #f5f5f5;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; 
            padding: 30px; 
            border-radius: 10px;
            margin-bottom: 30px;
            text-align: center;
        }
        .header h1 { margin: 0; font-size: 2.5em; }
        .header p { margin: 10px 0 0 0; opacity: 0.9; }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .summary-card {
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            text-align: center;
        }
        .summary-card h3 { margin: 0 0 10px 0; color: #333; }
        .summary-card .value { font-size: 2em; font-weight: bold; }
        .score { color: #28a745; }
        .critical { color: #dc3545; }
        .high { color: #fd7e14; }
        .medium { color: #ffc107; }
        .low { color: #6c757d; }
        .section {
            background: white;
            margin-bottom: 30px;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .section-header {
            background: #343a40;
            color: white;
            padding: 20px;
            font-size: 1.3em;
            font-weight: bold;
        }
        .section-content { padding: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: bold; }
        .keyword { 
            display: inline-block; 
            background: #e3f2fd; 
            color: #1976d2; 
            padding: 5px 10px; 
            border-radius: 15px; 
            margin: 2px;
            font-size: 0.9em;
        }
        .confidence { font-weight: bold; }
        .footer { 
            text-align: center; 
            margin-top: 40px; 
            color: #6c757d; 
            font-size: 0.9em;
        }
        .logo { max-height: 40px; margin-bottom: 10px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🦎 Fire Salamander</h1>
        <p>Rapport d'audit SEO - {{.SiteURL}}</p>
        <p>Généré le {{.GeneratedAt}} | Audit ID: {{.AuditID}}</p>
    </div>

    <div class="summary">
        <div class="summary-card">
            <h3>Score Global</h3>
            <div class="value score">{{printf "%.0f" (mul .OverallScore 100)}}%</div>
        </div>
        <div class="summary-card">
            <h3>Pages Analysées</h3>
            <div class="value">{{.TotalPages}}</div>
        </div>
        <div class="summary-card">
            <h3>Problèmes Critiques</h3>
            <div class="value critical">{{.CriticalIssues}}</div>
        </div>
        <div class="summary-card">
            <h3>Durée d'analyse</h3>
            <div class="value">{{.Duration}}</div>
        </div>
    </div>

    <div class="section">
        <div class="section-header">📊 Répartition des Problèmes</div>
        <div class="section-content">
            <div class="summary">
                <div class="summary-card">
                    <h3>Critique</h3>
                    <div class="value critical">{{.CriticalIssues}}</div>
                </div>
                <div class="summary-card">
                    <h3>Élevé</h3>
                    <div class="value high">{{.HighIssues}}</div>
                </div>
                <div class="summary-card">
                    <h3>Moyen</h3>
                    <div class="value medium">{{.MediumIssues}}</div>
                </div>
                <div class="summary-card">
                    <h3>Faible</h3>
                    <div class="value low">{{.LowIssues}}</div>
                </div>
            </div>
        </div>
    </div>

    <div class="section">
        <div class="section-header">🔍 Problèmes Détectés</div>
        <div class="section-content">
            {{if .Issues}}
            <table>
                <thead>
                    <tr>
                        <th>Problème</th>
                        <th>Sévérité</th>
                        <th>Occurrences</th>
                        <th>Pages Affectées</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Issues}}
                    <tr>
                        <td>{{.Message}}</td>
                        <td><span class="{{.Severity}}">{{.Severity}}</span></td>
                        <td>{{.Count}}</td>
                        <td>{{range $i, $page := .Pages}}{{if $i}}, {{end}}{{$page}}{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <p>Aucun problème détecté ! 🎉</p>
            {{end}}
        </div>
    </div>

    <div class="section">
        <div class="section-header">🎯 Suggestions de Mots-clés</div>
        <div class="section-content">
            {{if .Keywords}}
            <table>
                <thead>
                    <tr>
                        <th>Mot-clé</th>
                        <th>Confiance</th>
                        <th>Raison</th>
                        <th>Preuve</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Keywords}}
                    <tr>
                        <td><span class="keyword">{{.Keyword}}</span></td>
                        <td><span class="confidence">{{printf "%.0f" (mul .Confidence 100)}}%</span></td>
                        <td>{{.Reason}}</td>
                        <td>{{.Evidence}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <p>Aucune suggestion de mot-clé générée.</p>
            {{end}}
        </div>
    </div>

    <div class="section">
        <div class="section-header">🏷️ Thématiques Identifiées</div>
        <div class="section-content">
            {{if .Topics}}
            {{range .Topics}}
            <div style="margin-bottom: 20px;">
                <h4>{{.Label}}</h4>
                <div>
                    {{range .Terms}}<span class="keyword">{{.}}</span>{{end}}
                </div>
            </div>
            {{end}}
            {{else}}
            <p>Aucune thématique identifiée.</p>
            {{end}}
        </div>
    </div>

    <div class="section">
        <div class="section-header">📄 Pages Analysées</div>
        <div class="section-content">
            <table>
                <thead>
                    <tr>
                        <th>URL</th>
                        <th>Titre</th>
                        <th>H1</th>
                        <th>Profondeur</th>
                        <th>Problèmes</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Pages}}
                    <tr>
                        <td><a href="{{.URL}}" target="_blank">{{.URL}}</a></td>
                        <td>{{.Title}}</td>
                        <td>{{.H1}}</td>
                        <td>{{.Depth}}</td>
                        <td>{{.IssuesCount}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
    </div>

    <div class="footer">
        <p>Rapport généré par <strong>Fire Salamander v1.0</strong> - SEPTEO</p>
        <p>🦎 Audit SEO nouvelle génération</p>
    </div>
</body>
</html>`