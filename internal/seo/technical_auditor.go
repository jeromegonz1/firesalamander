package seo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"firesalamander/internal/constants"

	"golang.org/x/net/html"
)

// TechnicalAuditor auditeur technique SEO
type TechnicalAuditor struct {
	client            *http.Client
	mobileFriendlyRegex *regexp.Regexp
	validHTMLRegex      *regexp.Regexp
}

// TechnicalAuditResult résultat de l'audit technique
type TechnicalAuditResult struct {
	Security      SecurityAudit      `json:"security"`
	Mobile        MobileAudit        `json:"mobile"`
	Structure     StructureAudit     `json:"structure"`
	Accessibility AccessibilityAudit `json:"accessibility"`
	Indexability  IndexabilityAudit  `json:"indexability"`
	Crawlability  CrawlabilityAudit  `json:"crawlability"`
	
	Issues          []string         `json:"issues"`
	Recommendations []string         `json:"recommendations"`
	CriticalIssues  []string         `json:"critical_issues"`
}

// SecurityAudit audit de sécurité
type SecurityAudit struct {
	HasHTTPS        bool     `json:"has_https"`
	ValidSSL        bool     `json:"valid_ssl"`
	HasHSTS         bool     `json:"has_hsts"`
	HasCSP          bool     `json:"has_csp"`
	MixedContent    bool     `json:"mixed_content"`
	InsecureLinks   []string `json:"insecure_links"`
	SecurityScore   float64  `json:"security_score"`
}

// MobileAudit audit mobile
type MobileAudit struct {
	IsResponsive      bool     `json:"is_responsive"`
	HasViewport       bool     `json:"has_viewport"`
	ViewportContent   string   `json:"viewport_content"`
	TouchFriendly     bool     `json:"touch_friendly"`
	FontReadable      bool     `json:"font_readable"`
	TapTargetsGoodSize bool    `json:"tap_targets_good_size"`
	MobileScore       float64  `json:"mobile_score"`
}

// StructureAudit audit de structure
type StructureAudit struct {
	HasSitemap       bool     `json:"has_sitemap"`
	SitemapURL       string   `json:"sitemap_url"`
	HasRobotsTxt     bool     `json:"has_robots_txt"`
	RobotsTxtContent string   `json:"robots_txt_content"`
	ValidHTML        bool     `json:"valid_html"`
	HTMLErrors       []string `json:"html_errors"`
	W3CValidation    bool     `json:"w3c_validation"`
	StructureScore   float64  `json:"structure_score"`
}

// AccessibilityAudit audit d'accessibilité
type AccessibilityAudit struct {
	HasAltTags        bool     `json:"has_alt_tags"`
	AltTagCoverage    float64  `json:"alt_tag_coverage"`
	HasAriaLabels     bool     `json:"has_aria_labels"`
	ColorContrast     bool     `json:"color_contrast"`
	KeyboardNav       bool     `json:"keyboard_navigation"`
	ScreenReaderReady bool     `json:"screen_reader_ready"`
	A11yIssues        []string `json:"a11y_issues"`
	Score             float64  `json:"score"`
}

// IndexabilityAudit audit d'indexabilité
type IndexabilityAudit struct {
	BlockedByRobots   bool     `json:"blocked_by_robots"`
	HasNoIndex        bool     `json:"has_noindex"`
	HasCanonical      bool     `json:"has_canonical"`
	CanonicalIssues   []string `json:"canonical_issues"`
	DuplicateContent  bool     `json:"duplicate_content"`
	IndexabilityScore float64  `json:"indexability_score"`
}

// CrawlabilityAudit audit de crawlabilité
type CrawlabilityAudit struct {
	InternalLinks     int      `json:"internal_links"`
	BrokenLinks       []string `json:"broken_links"`
	RedirectChains    []string `json:"redirect_chains"`
	OrphanPages       bool     `json:"orphan_pages"`
	DeepPageLevel     int      `json:"deep_page_level"`
	CrawlabilityScore float64  `json:"crawlability_score"`
}

// NewTechnicalAuditor crée un nouvel auditeur technique
func NewTechnicalAuditor() *TechnicalAuditor {
	client := &http.Client{
		Timeout: constants.DefaultRequestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("stopped after 5 redirects")
			}
			return nil
		},
	}

	return &TechnicalAuditor{
		client:              client,
		mobileFriendlyRegex: regexp.MustCompile(`width=device-width|initial-scale=1`),
		validHTMLRegex:      regexp.MustCompile(`<!DOCTYPE html>`),
	}
}

// Audit effectue l'audit technique complet
func (ta *TechnicalAuditor) Audit(ctx context.Context, targetURL string, doc *html.Node, htmlContent string) (*TechnicalAuditResult, error) {
	result := &TechnicalAuditResult{
		Issues:          []string{},
		Recommendations: []string{},
		CriticalIssues:  []string{},
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("URL invalide: %w", err)
	}

	// 1. Audit de sécurité
	result.Security = ta.auditSecurity(ctx, targetURL, parsedURL, htmlContent)

	// 2. Audit mobile
	result.Mobile = ta.auditMobile(doc, htmlContent)

	// 3. Audit de structure
	result.Structure = ta.auditStructure(ctx, parsedURL, doc, htmlContent)

	// 4. Audit d'accessibilité
	result.Accessibility = ta.auditAccessibility(doc, htmlContent)

	// 5. Audit d'indexabilité
	result.Indexability = ta.auditIndexability(doc, htmlContent)

	// 6. Audit de crawlabilité
	result.Crawlability = ta.auditCrawlability(ctx, parsedURL, doc)

	// 7. Consolider les issues et recommandations
	ta.consolidateResults(result)

	return result, nil
}

// auditSecurity effectue l'audit de sécurité
func (ta *TechnicalAuditor) auditSecurity(ctx context.Context, targetURL string, parsedURL *url.URL, htmlContent string) SecurityAudit {
	audit := SecurityAudit{
		InsecureLinks: []string{},
	}

	// Vérifier HTTPS
	audit.HasHTTPS = parsedURL.Scheme == "https"

	// Test SSL/TLS
	if audit.HasHTTPS {
		audit.ValidSSL = ta.testSSLCertificate(ctx, parsedURL.Host)
	}

	// Headers de sécurité
	if resp, err := ta.fetchHeaders(ctx, targetURL); err == nil {
		audit.HasHSTS = resp.Header.Get("Strict-Transport-Security") != ""
		audit.HasCSP = resp.Header.Get("Content-Security-Policy") != ""
	}

	// Détecter le contenu mixte
	if audit.HasHTTPS {
		audit.MixedContent = strings.Contains(htmlContent, constants.HTTPPrefix) ||
			strings.Contains(htmlContent, `src="http:`) ||
			strings.Contains(htmlContent, `href="http:`)

		// Identifier les liens non sécurisés
		httpLinks := regexp.MustCompile(`https?://[^"'\s]+`).FindAllString(htmlContent, -1)
		for _, link := range httpLinks {
			if strings.HasPrefix(link, constants.HTTPPrefix) {
				audit.InsecureLinks = append(audit.InsecureLinks, link)
			}
		}
	}

	// Calculer le score de sécurité
	audit.SecurityScore = ta.calculateSecurityScore(audit)

	return audit
}

// auditMobile effectue l'audit mobile
func (ta *TechnicalAuditor) auditMobile(doc *html.Node, htmlContent string) MobileAudit {
	audit := MobileAudit{}

	// Viewport meta tag
	viewport := ta.findMetaByName(doc, "viewport")
	audit.HasViewport = viewport != ""
	audit.ViewportContent = viewport

	if audit.HasViewport {
		audit.IsResponsive = ta.mobileFriendlyRegex.MatchString(viewport)
	}

	// Vérifier les polices lisibles
	audit.FontReadable = ta.checkFontReadability(htmlContent)

	// Vérifier les cibles tactiles
	audit.TapTargetsGoodSize = ta.checkTapTargets(htmlContent)

	// Touch-friendly
	audit.TouchFriendly = audit.IsResponsive && audit.TapTargetsGoodSize

	// Score mobile
	audit.MobileScore = ta.calculateMobileScore(audit)

	return audit
}

// auditStructure effectue l'audit de structure
func (ta *TechnicalAuditor) auditStructure(ctx context.Context, parsedURL *url.URL, doc *html.Node, htmlContent string) StructureAudit {
	audit := StructureAudit{
		HTMLErrors: []string{},
	}

	// Vérifier sitemap.xml
	sitemapURL := fmt.Sprintf("%s://%s/sitemap.xml", parsedURL.Scheme, parsedURL.Host)
	if ta.urlExists(ctx, sitemapURL) {
		audit.HasSitemap = true
		audit.SitemapURL = sitemapURL
	}

	// Vérifier robots.txt
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	if content := ta.fetchTextContent(ctx, robotsURL); content != "" {
		audit.HasRobotsTxt = true
		audit.RobotsTxtContent = content
	}

	// Validation HTML basique
	audit.ValidHTML = ta.validHTMLRegex.MatchString(htmlContent)
	if !audit.ValidHTML {
		audit.HTMLErrors = append(audit.HTMLErrors, "DOCTYPE HTML5 manquant")
	}

	// Vérifications HTML supplémentaires
	audit.HTMLErrors = append(audit.HTMLErrors, ta.validateHTMLStructure(doc)...)

	// Score de structure
	audit.StructureScore = ta.calculateStructureScore(audit)

	return audit
}

// auditAccessibility effectue l'audit d'accessibilité
func (ta *TechnicalAuditor) auditAccessibility(doc *html.Node, htmlContent string) AccessibilityAudit {
	audit := AccessibilityAudit{
		A11yIssues: []string{},
	}

	// Images avec alt
	images := ta.findAllNodesByAtom(doc, "img")
	imagesWithAlt := 0
	for _, img := range images {
		if ta.getAttr(img, "alt") != "" {
			imagesWithAlt++
		}
	}

	if len(images) > 0 {
		audit.AltTagCoverage = float64(imagesWithAlt) / float64(len(images))
		audit.HasAltTags = audit.AltTagCoverage > 0.8
	}

	// Labels ARIA
	audit.HasAriaLabels = strings.Contains(htmlContent, "aria-label") ||
		strings.Contains(htmlContent, "aria-labelledby")

	// Navigation clavier
	audit.KeyboardNav = strings.Contains(htmlContent, "tabindex") ||
		strings.Contains(htmlContent, "accesskey")

	// Screen reader
	audit.ScreenReaderReady = audit.HasAriaLabels && audit.HasAltTags

	// Contraste de couleurs (détection basique)
	audit.ColorContrast = ta.checkColorContrast(htmlContent)

	// Issues d'accessibilité
	if audit.AltTagCoverage < 1.0 {
		audit.A11yIssues = append(audit.A11yIssues, "Images sans texte alternatif")
	}
	if !audit.HasAriaLabels {
		audit.A11yIssues = append(audit.A11yIssues, "Labels ARIA manquants")
	}
	if !audit.KeyboardNav {
		audit.A11yIssues = append(audit.A11yIssues, "Navigation clavier insuffisante")
	}

	// Score d'accessibilité
	audit.Score = ta.calculateAccessibilityScore(audit)

	return audit
}

// auditIndexability effectue l'audit d'indexabilité
func (ta *TechnicalAuditor) auditIndexability(doc *html.Node, htmlContent string) IndexabilityAudit {
	audit := IndexabilityAudit{
		CanonicalIssues: []string{},
	}

	// Meta robots
	robots := ta.findMetaByName(doc, "robots")
	audit.BlockedByRobots = strings.Contains(strings.ToLower(robots), "noindex") ||
		strings.Contains(strings.ToLower(robots), "nofollow")
	audit.HasNoIndex = strings.Contains(strings.ToLower(robots), "noindex")

	// Canonical URL
	canonical := ta.findLinkByRel(doc, "canonical")
	audit.HasCanonical = canonical != ""

	// Vérifications canonical
	if audit.HasCanonical {
		if !strings.HasPrefix(canonical, "http") {
			audit.CanonicalIssues = append(audit.CanonicalIssues, "URL canonique relative")
		}
	} else {
		audit.CanonicalIssues = append(audit.CanonicalIssues, "URL canonique manquante")
	}

	// Contenu dupliqué (heuristique simple)
	audit.DuplicateContent = ta.detectDuplicateContent(htmlContent)

	// Score d'indexabilité
	audit.IndexabilityScore = ta.calculateIndexabilityScore(audit)

	return audit
}

// auditCrawlability effectue l'audit de crawlabilité
func (ta *TechnicalAuditor) auditCrawlability(ctx context.Context, parsedURL *url.URL, doc *html.Node) CrawlabilityAudit {
	audit := CrawlabilityAudit{
		BrokenLinks:    []string{},
		RedirectChains: []string{},
	}

	// Compter les liens internes
	links := ta.findAllNodesByAtom(doc, "a")
	for _, link := range links {
		href := ta.getAttr(link, "href")
		if href != "" && (strings.HasPrefix(href, "/") || strings.Contains(href, parsedURL.Host)) {
			audit.InternalLinks++
		}
	}

	// Tester quelques liens (limité pour éviter la surcharge)
	testedLinks := 0
	for _, link := range links {
		if testedLinks >= 10 { // Limiter les tests
			break
		}

		href := ta.getAttr(link, "href")
		if href != "" && strings.HasPrefix(href, "http") {
			if !ta.urlExists(ctx, href) {
				audit.BrokenLinks = append(audit.BrokenLinks, href)
			}
			testedLinks++
		}
	}

	// Score de crawlabilité
	audit.CrawlabilityScore = ta.calculateCrawlabilityScore(audit)

	return audit
}

// consolidateResults consolide les résultats
func (ta *TechnicalAuditor) consolidateResults(result *TechnicalAuditResult) {
	// Issues critiques
	if !result.Security.HasHTTPS {
		result.CriticalIssues = append(result.CriticalIssues, "Site non sécurisé (HTTP)")
	}
	if result.Security.MixedContent {
		result.CriticalIssues = append(result.CriticalIssues, "Contenu mixte détecté")
	}
	if !result.Mobile.HasViewport {
		result.CriticalIssues = append(result.CriticalIssues, "Meta viewport manquante")
	}
	if result.Indexability.HasNoIndex {
		result.CriticalIssues = append(result.CriticalIssues, "Page marquée noindex")
	}

	// Issues générales
	if !result.Structure.HasSitemap {
		result.Issues = append(result.Issues, "Sitemap.xml manquant")
	}
	if !result.Structure.HasRobotsTxt {
		result.Issues = append(result.Issues, "Robots.txt manquant")
	}
	if result.Accessibility.Score < 0.7 {
		result.Issues = append(result.Issues, "Score d'accessibilité faible")
	}

	// Recommandations
	if !result.Security.HasHTTPS {
		result.Recommendations = append(result.Recommendations, "Migrer vers HTTPS")
	}
	if !result.Security.HasHSTS {
		result.Recommendations = append(result.Recommendations, "Ajouter le header HSTS")
	}
	if !result.Mobile.IsResponsive {
		result.Recommendations = append(result.Recommendations, "Rendre le site responsive")
	}
	if len(result.Crawlability.BrokenLinks) > 0 {
		result.Recommendations = append(result.Recommendations, "Corriger les liens brisés")
	}
}

// Fonctions utilitaires et de calcul

func (ta *TechnicalAuditor) testSSLCertificate(ctx context.Context, host string) bool {
	// Test simple de connexion SSL
	req, err := http.NewRequestWithContext(ctx, "HEAD", constants.HTTPSPrefix+host, nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: constants.FastRequestTimeout}
	_, err = client.Do(req)
	return err == nil
}

func (ta *TechnicalAuditor) fetchHeaders(ctx context.Context, targetURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return nil, err
	}
	return ta.client.Do(req)
}

func (ta *TechnicalAuditor) urlExists(ctx context.Context, targetURL string) bool {
	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return false
	}

	resp, err := ta.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400
}

func (ta *TechnicalAuditor) fetchTextContent(ctx context.Context, targetURL string) string {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return ""
	}

	resp, err := ta.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	buf := make([]byte, 4096) // 4KB max
	n, _ := resp.Body.Read(buf)
	return string(buf[:n])
}

func (ta *TechnicalAuditor) checkFontReadability(htmlContent string) bool {
	// Vérifier les tailles de police
	return !strings.Contains(htmlContent, "font-size:0") &&
		!strings.Contains(htmlContent, "font-size: 0")
}

func (ta *TechnicalAuditor) checkTapTargets(htmlContent string) bool {
	// Heuristique simple pour les cibles tactiles
	return strings.Contains(htmlContent, "min-height") ||
		strings.Contains(htmlContent, "padding")
}

func (ta *TechnicalAuditor) checkColorContrast(htmlContent string) bool {
	// Détection basique - chercher les styles de couleur
	hasColors := strings.Contains(htmlContent, "color:") ||
		strings.Contains(htmlContent, "background-color:")
	
	// Supposer un contraste correct si pas de styles inline détectés
	return !hasColors || len(htmlContent) > 1000
}

func (ta *TechnicalAuditor) validateHTMLStructure(doc *html.Node) []string {
	var errors []string

	// Vérifier la structure de base
	if ta.findNodeByAtom(doc, "html") == nil {
		errors = append(errors, "Balise HTML manquante")
	}
	if ta.findNodeByAtom(doc, "head") == nil {
		errors = append(errors, "Balise HEAD manquante")
	}
	if ta.findNodeByAtom(doc, "body") == nil {
		errors = append(errors, "Balise BODY manquante")
	}

	return errors
}

func (ta *TechnicalAuditor) detectDuplicateContent(htmlContent string) bool {
	// Heuristique simple pour le contenu dupliqué
	words := strings.Fields(htmlContent)
	if len(words) < 50 {
		return false
	}

	wordCount := make(map[string]int)
	for _, word := range words {
		if len(word) > 3 {
			wordCount[strings.ToLower(word)]++
		}
	}

	// Si plus de 40% des mots sont répétés plus de 3 fois
	repeatedWords := 0
	for _, count := range wordCount {
		if count > 3 {
			repeatedWords += count
		}
	}

	return float64(repeatedWords)/float64(len(words)) > 0.4
}

// Fonctions de calcul de score

func (ta *TechnicalAuditor) calculateSecurityScore(audit SecurityAudit) float64 {
	score := 0.0
	
	if audit.HasHTTPS {
		score += 0.4
	}
	if audit.ValidSSL {
		score += 0.2
	}
	if audit.HasHSTS {
		score += 0.2
	}
	if audit.HasCSP {
		score += 0.1
	}
	if !audit.MixedContent {
		score += 0.1
	}

	return score
}

func (ta *TechnicalAuditor) calculateMobileScore(audit MobileAudit) float64 {
	score := 0.0

	if audit.HasViewport {
		score += 0.3
	}
	if audit.IsResponsive {
		score += 0.3
	}
	if audit.TouchFriendly {
		score += 0.2
	}
	if audit.FontReadable {
		score += 0.1
	}
	if audit.TapTargetsGoodSize {
		score += 0.1
	}

	return score
}

func (ta *TechnicalAuditor) calculateStructureScore(audit StructureAudit) float64 {
	score := 0.0

	if audit.HasSitemap {
		score += 0.3
	}
	if audit.HasRobotsTxt {
		score += 0.2
	}
	if audit.ValidHTML {
		score += 0.3
	}
	if len(audit.HTMLErrors) == 0 {
		score += 0.2
	}

	return score
}

func (ta *TechnicalAuditor) calculateAccessibilityScore(audit AccessibilityAudit) float64 {
	score := 0.0

	score += audit.AltTagCoverage * 0.3

	if audit.HasAriaLabels {
		score += 0.2
	}
	if audit.ColorContrast {
		score += 0.2
	}
	if audit.KeyboardNav {
		score += 0.15
	}
	if audit.ScreenReaderReady {
		score += 0.15
	}

	return score
}

func (ta *TechnicalAuditor) calculateIndexabilityScore(audit IndexabilityAudit) float64 {
	score := 1.0

	if audit.BlockedByRobots {
		score -= 0.5
	}
	if audit.HasNoIndex {
		score -= 0.3
	}
	if !audit.HasCanonical {
		score -= 0.1
	}
	if audit.DuplicateContent {
		score -= 0.1
	}

	return max(0.0, score)
}

func (ta *TechnicalAuditor) calculateCrawlabilityScore(audit CrawlabilityAudit) float64 {
	score := 0.5 // Score de base

	if audit.InternalLinks >= 3 {
		score += 0.3
	} else if audit.InternalLinks > 0 {
		score += 0.15
	}

	if len(audit.BrokenLinks) == 0 {
		score += 0.2
	} else {
		score -= float64(len(audit.BrokenLinks)) * 0.05
	}

	return max(0.0, min(1.0, score))
}

// Fonctions utilitaires pour parcourir le DOM

func (ta *TechnicalAuditor) findNodeByAtom(doc *html.Node, tagName string) *html.Node {
	var result *html.Node
	ta.walkHTML(doc, func(n *html.Node) {
		if result == nil && n.Type == html.ElementNode && n.Data == tagName {
			result = n
		}
	})
	return result
}

func (ta *TechnicalAuditor) findAllNodesByAtom(doc *html.Node, tagName string) []*html.Node {
	var results []*html.Node
	ta.walkHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tagName {
			results = append(results, n)
		}
	})
	return results
}

func (ta *TechnicalAuditor) findMetaByName(doc *html.Node, name string) string {
	var content string
	ta.walkHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			nameAttr := ta.getAttr(n, "name")
			if nameAttr == name {
				content = ta.getAttr(n, "content")
			}
		}
	})
	return content
}

func (ta *TechnicalAuditor) findLinkByRel(doc *html.Node, rel string) string {
	var href string
	ta.walkHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			relAttr := ta.getAttr(n, "rel")
			if relAttr == rel {
				href = ta.getAttr(n, "href")
			}
		}
	})
	return href
}

func (ta *TechnicalAuditor) walkHTML(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ta.walkHTML(c, fn)
	}
}

func (ta *TechnicalAuditor) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// Fonctions utilitaires mathématiques
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}