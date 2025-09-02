package technical

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
)

// TechnicalAuditor implémente l'agent d'audit technique HTML/CSS/Performance
type TechnicalAuditor struct {
	name string
}

// NewTechnicalAuditor crée une nouvelle instance de TechnicalAuditor
func NewTechnicalAuditor() *TechnicalAuditor {
	return &TechnicalAuditor{
		name: constants.AgentNameTechnical,
	}
}

// Name retourne le nom de l'agent
func (t *TechnicalAuditor) Name() string {
	return t.name
}

// Process traite les données d'entrée et effectue l'audit technique
func (t *TechnicalAuditor) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()
	
	pageData, ok := data.(*agents.PageData)
	if !ok {
		return &agents.AgentResult{
			AgentName: t.name,
			Status:    constants.StatusFailed,
			Errors:    []string{"invalid input data type, expected *PageData"},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	report, err := t.AuditPage(pageData)
	if err != nil {
		return &agents.AgentResult{
			AgentName: t.name,
			Status:    constants.StatusFailed,
			Errors:    []string{err.Error()},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	return &agents.AgentResult{
		AgentName: t.name,
		Status:    constants.StatusCompleted,
		Data: map[string]interface{}{
			"technical_report": report,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}, nil
}

// HealthCheck vérifie la santé de l'agent
func (t *TechnicalAuditor) HealthCheck() error {
	// Test simple d'audit
	testPage := &agents.PageData{
		URL:  "http://test.example.com",
		HTML: "<html><head><title>Test</title></head><body><h1>Test</h1></body></html>",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
	
	_, err := t.AuditPage(testPage)
	return err
}

// AuditPage effectue un audit technique complet d'une page
func (t *TechnicalAuditor) AuditPage(page *agents.PageData) (*agents.TechnicalReport, error) {
	if page == nil {
		return nil, fmt.Errorf("page data cannot be nil")
	}

	// Audit de performance
	performance := t.auditPerformance(page)
	
	// Audit d'accessibilité
	accessibility := t.auditAccessibility(page)
	
	// Audit SEO technique
	seo := t.auditSEO(page)
	
	// Collecte des problèmes techniques
	issues := t.collectIssues(page)

	return &agents.TechnicalReport{
		PageURL:       page.URL,
		Performance:   performance,
		Accessibility: accessibility,
		SEO:           seo,
		Issues:        issues,
	}, nil
}

// ValidateStructure valide la structure HTML
func (t *TechnicalAuditor) ValidateStructure(html string) (*agents.StructureResult, error) {
	if html == "" {
		return &agents.StructureResult{
			Valid:        false,
			Errors:       []agents.StructureError{{Message: "HTML content is empty"}},
			Warnings:     []agents.StructureError{},
			HeadingLevel: 0,
		}, nil
	}

	var errors []agents.StructureError
	var warnings []agents.StructureError
	valid := true

	// Vérification de la structure de base
	if !strings.Contains(html, "<html") {
		errors = append(errors, agents.StructureError{
			Message: "Missing <html> tag",
			Element: "html",
		})
		valid = false
	}

	if !strings.Contains(html, "<head") {
		errors = append(errors, agents.StructureError{
			Message: "Missing <head> section",
			Element: "head",
		})
		valid = false
	}

	if !strings.Contains(html, "<body") {
		errors = append(errors, agents.StructureError{
			Message: "Missing <body> section",
			Element: "body",
		})
		valid = false
	}

	// Vérification des balises non fermées
	openTags := t.findUnclosedTags(html)
	for _, tag := range openTags {
		errors = append(errors, agents.StructureError{
			Message: fmt.Sprintf("Unclosed tag: %s", tag),
			Element: tag,
		})
		valid = false
	}

	// Vérification de la hiérarchie des titres
	headingLevel := t.analyzeHeadingHierarchy(html)
	if headingLevel > 6 {
		warnings = append(warnings, agents.StructureError{
			Message: "Heading hierarchy is too deep (>6 levels)",
			Element: "headings",
		})
	}

	return &agents.StructureResult{
		Valid:        valid,
		Errors:       errors,
		Warnings:     warnings,
		HeadingLevel: headingLevel,
	}, nil
}

// auditPerformance évalue les métriques de performance
func (t *TechnicalAuditor) auditPerformance(page *agents.PageData) agents.PerformanceScore {
	score := 100
	resourceCount := 0
	loadTime := int64(100) // Simulation temps de chargement

	// Compte les ressources externes
	externalResourceRegex := regexp.MustCompile(`(?i)(src|href)=["'][^"']*\.(css|js|png|jpg|jpeg|gif|svg)["']`)
	matches := externalResourceRegex.FindAllString(page.HTML, -1)
	resourceCount = len(matches)

	// Pénalise les ressources excessives
	if resourceCount > 20 {
		score -= 20
	} else if resourceCount > 10 {
		score -= 10
	}

	// Vérifie la compression
	if contentEncoding := page.Headers["Content-Encoding"]; contentEncoding == "" {
		score -= 10
	}

	// Vérifie la mise en cache
	if cacheControl := page.Headers["Cache-Control"]; cacheControl == "" {
		score -= 10
	}

	// Vérifie la taille du HTML
	if len(page.HTML) > 100000 { // Plus de 100KB
		score -= 15
		loadTime += 200
	} else if len(page.HTML) > 50000 { // Plus de 50KB
		score -= 10
		loadTime += 100
	}

	if score < 0 {
		score = 0
	}

	return agents.PerformanceScore{
		Score:     score,
		LoadTime:  loadTime,
		Resources: resourceCount,
	}
}

// auditAccessibility évalue l'accessibilité
func (t *TechnicalAuditor) auditAccessibility(page *agents.PageData) agents.AccessibilityScore {
	score := 100
	var issues []string

	// Vérification des images sans alt
	imgRegex := regexp.MustCompile(`<img[^>]*>`)
	imgMatches := imgRegex.FindAllString(page.HTML, -1)
	
	for _, img := range imgMatches {
		if !strings.Contains(img, "alt=") {
			issues = append(issues, "Image without alt attribute found")
			score -= 10
		} else if strings.Contains(img, `alt=""`) || strings.Contains(img, "alt=''") {
			issues = append(issues, "Image with empty alt attribute found")
			score -= 5
		}
	}

	// Vérification des titres de liens
	linkRegex := regexp.MustCompile(`<a[^>]*>[^<]*</a>`)
	linkMatches := linkRegex.FindAllString(page.HTML, -1)
	
	for _, link := range linkMatches {
		linkText := regexp.MustCompile(`>[^<]*<`).FindString(link)
		if linkText != "" {
			linkText = strings.Trim(linkText, "><")
			if len(linkText) < 2 {
				issues = append(issues, "Link with insufficient descriptive text found")
				score -= 5
			}
		}
	}

	// Vérification des labels de formulaires
	inputRegex := regexp.MustCompile(`<input[^>]*>`)
	inputMatches := inputRegex.FindAllString(page.HTML, -1)
	
	for _, input := range inputMatches {
		inputType := t.extractAttribute(input, "type")
		if inputType != "hidden" && inputType != "submit" && inputType != "button" {
			id := t.extractAttribute(input, "id")
			if id != "" {
				labelPattern := fmt.Sprintf(`<label[^>]*for=["']%s["'][^>]*>`, id)
				labelRegex := regexp.MustCompile(labelPattern)
				if !labelRegex.MatchString(page.HTML) {
					issues = append(issues, "Input field without associated label found")
					score -= 8
				}
			} else {
				issues = append(issues, "Input field without ID for label association found")
				score -= 10
			}
		}
	}

	// Vérification du contraste (simulation basique)
	if strings.Contains(page.HTML, `color:#fff`) && strings.Contains(page.HTML, `background:#fff`) {
		issues = append(issues, "Potential contrast issue detected")
		score -= 15
	}

	if score < 0 {
		score = 0
	}

	return agents.AccessibilityScore{
		Score:  score,
		Issues: issues,
	}
}

// auditSEO évalue les éléments SEO techniques
func (t *TechnicalAuditor) auditSEO(page *agents.PageData) agents.SEOScore {
	score := 100
	var missingElements []string

	// Vérification du titre
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]*)</title>`)
	titleMatch := titleRegex.FindStringSubmatch(page.HTML)
	if len(titleMatch) == 0 {
		missingElements = append(missingElements, "title tag")
		score -= 25
	} else {
		title := strings.TrimSpace(titleMatch[1])
		if len(title) == 0 {
			missingElements = append(missingElements, "title content")
			score -= 20
		} else if len(title) < 10 {
			missingElements = append(missingElements, "title too short")
			score -= 10
		} else if len(title) > 60 {
			missingElements = append(missingElements, "title too long")
			score -= 10
		}
	}

	// Vérification de la meta description
	metaDescRegex := regexp.MustCompile(`<meta[^>]*name=["']description["'][^>]*content=["']([^"']*)["']`)
	metaDescMatch := metaDescRegex.FindStringSubmatch(page.HTML)
	if len(metaDescMatch) == 0 {
		missingElements = append(missingElements, "meta description")
		score -= 20
	} else {
		desc := strings.TrimSpace(metaDescMatch[1])
		if len(desc) == 0 {
			missingElements = append(missingElements, "meta description content")
			score -= 15
		} else if len(desc) < 120 {
			missingElements = append(missingElements, "meta description too short")
			score -= 5
		} else if len(desc) > 160 {
			missingElements = append(missingElements, "meta description too long")
			score -= 5
		}
	}

	// Vérification des balises H1
	h1Regex := regexp.MustCompile(`<h1[^>]*>([^<]*)</h1>`)
	h1Matches := h1Regex.FindAllStringSubmatch(page.HTML, -1)
	if len(h1Matches) == 0 {
		missingElements = append(missingElements, "H1 tag")
		score -= 15
	} else if len(h1Matches) > 1 {
		missingElements = append(missingElements, "multiple H1 tags")
		score -= 10
	}

	// Vérification de l'attribut lang
	if !strings.Contains(page.HTML, `lang=`) {
		missingElements = append(missingElements, "lang attribute")
		score -= 10
	}

	// Vérification de la balise meta viewport
	if !strings.Contains(page.HTML, `name="viewport"`) {
		missingElements = append(missingElements, "viewport meta tag")
		score -= 15
	}

	if score < 0 {
		score = 0
	}

	return agents.SEOScore{
		Score:           score,
		MissingElements: missingElements,
	}
}

// collectIssues collecte tous les problèmes techniques détectés
func (t *TechnicalAuditor) collectIssues(page *agents.PageData) []agents.TechnicalIssue {
	var issues []agents.TechnicalIssue

	// Issues de performance
	if len(page.HTML) > 100000 {
		issues = append(issues, agents.TechnicalIssue{
			Type:        "performance",
			Severity:    "high",
			Description: "HTML content size is too large (>100KB)",
		})
	}

	// Issues de sécurité
	if strings.Contains(page.HTML, "http://") && !strings.Contains(page.URL, "http://localhost") {
		issues = append(issues, agents.TechnicalIssue{
			Type:        "security",
			Severity:    "medium",
			Description: "Mixed content detected (HTTP resources on HTTPS page)",
		})
	}

	// Issues de structure
	structureResult, _ := t.ValidateStructure(page.HTML)
	for _, err := range structureResult.Errors {
		issues = append(issues, agents.TechnicalIssue{
			Type:        "structure",
			Severity:    "high",
			Description: err.Message,
			Element:     err.Element,
		})
	}

	for _, warning := range structureResult.Warnings {
		issues = append(issues, agents.TechnicalIssue{
			Type:        "structure",
			Severity:    "medium",
			Description: warning.Message,
			Element:     warning.Element,
		})
	}

	// Issues d'optimisation
	if !strings.Contains(page.HTML, `<meta name="robots"`) {
		issues = append(issues, agents.TechnicalIssue{
			Type:        "seo",
			Severity:    "low",
			Description: "Missing robots meta tag",
		})
	}

	return issues
}

// findUnclosedTags trouve les balises non fermées (implémentation simplifiée)
func (t *TechnicalAuditor) findUnclosedTags(html string) []string {
	var unclosedTags []string
	
	// Tags qui n'ont pas besoin de fermeture
	selfClosingTags := map[string]bool{
		"img": true, "br": true, "hr": true, "meta": true,
		"link": true, "input": true, "area": true, "base": true,
		"col": true, "embed": true, "source": true, "track": true,
		"wbr": true,
	}

	// Expression régulière simple pour les balises ouvrantes
	openTagRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)[^>]*>`)
	closeTagRegex := regexp.MustCompile(`</([a-zA-Z][a-zA-Z0-9]*)>`)

	openTags := openTagRegex.FindAllStringSubmatch(html, -1)
	closeTags := closeTagRegex.FindAllStringSubmatch(html, -1)

	openCount := make(map[string]int)
	closeCount := make(map[string]int)

	for _, match := range openTags {
		tag := strings.ToLower(match[1])
		if !selfClosingTags[tag] {
			openCount[tag]++
		}
	}

	for _, match := range closeTags {
		tag := strings.ToLower(match[1])
		closeCount[tag]++
	}

	for tag, opens := range openCount {
		closes := closeCount[tag]
		if opens > closes {
			unclosedTags = append(unclosedTags, tag)
		}
	}

	return unclosedTags
}

// analyzeHeadingHierarchy analyse la hiérarchie des titres
func (t *TechnicalAuditor) analyzeHeadingHierarchy(html string) int {
	maxLevel := 0
	
	for i := 1; i <= 6; i++ {
		pattern := fmt.Sprintf(`<h%d[^>]*>`, i)
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(html) {
			if i > maxLevel {
				maxLevel = i
			}
		}
	}
	
	return maxLevel
}

// extractAttribute extrait un attribut d'une balise HTML
func (t *TechnicalAuditor) extractAttribute(tag, attr string) string {
	pattern := fmt.Sprintf(`%s=["']([^"']*)["']`, attr)
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(tag)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}