package seo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"firesalamander/internal/constants"

	"golang.org/x/net/html"
)

// SEOAnalyzer analyseur SEO technique avancé
type SEOAnalyzer struct {
	client       *http.Client
	userAgent    string
	timeout      time.Duration
	maxRedirects int
	
	// Composants spécialisés
	tagAnalyzer         *TagAnalyzer
	performanceAnalyzer *PerformanceAnalyzer
	technicalAuditor    *TechnicalAuditor
	recommendationEngine *RecommendationEngine
	
	mu sync.RWMutex
}

// SEOAnalysisResult résultat complet de l'analyse SEO
type SEOAnalysisResult struct {
	URL                string                    `json:"url"`
	Domain             string                    `json:"domain"`
	Protocol           string                    `json:"protocol"`
	StatusCode         int                       `json:"status_code"`
	ResponseTime       time.Duration             `json:"response_time"`
	
	// Analyses techniques
	TagAnalysis        TagAnalysisResult         `json:"tag_analysis"`
	PerformanceMetrics PerformanceMetricsResult  `json:"performance_metrics"`
	TechnicalAudit     TechnicalAuditResult      `json:"technical_audit"`
	
	// Recommandations et scoring
	Recommendations    []SEORecommendation       `json:"recommendations"`
	OverallScore       float64                   `json:"overall_score"`
	CategoryScores     map[string]float64        `json:"category_scores"`
	
	// Métadonnées
	AnalyzedAt         time.Time                 `json:"analyzed_at"`
	AnalysisDuration   time.Duration             `json:"analysis_duration"`
}

// NewSEOAnalyzer crée un nouvel analyseur SEO
func NewSEOAnalyzer() *SEOAnalyzer {
	client := &http.Client{
		Timeout: constants.MaxHTTPTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= constants.MaxRedirects {
				return fmt.Errorf("stopped after %d redirects", constants.MaxRedirects)
			}
			return nil
		},
	}

	analyzer := &SEOAnalyzer{
		client:       client,
		userAgent:    constants.SEOBotUserAgent,
		timeout:      constants.MaxHTTPTimeout,
		maxRedirects: constants.MaxRedirects,
	}

	// Initialiser les composants
	analyzer.tagAnalyzer = NewTagAnalyzer()
	analyzer.performanceAnalyzer = NewPerformanceAnalyzer()
	analyzer.technicalAuditor = NewTechnicalAuditor()
	analyzer.recommendationEngine = NewRecommendationEngine()

	return analyzer
}

// AnalyzePage effectue une analyse SEO complète d'une page
func (seo *SEOAnalyzer) AnalyzePage(ctx context.Context, targetURL string) (*SEOAnalysisResult, error) {
	startTime := time.Now()
	
	log.Printf("Début analyse SEO pour: %s", targetURL)

	// Validation de l'URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("URL invalide: %w", err)
	}

	result := &SEOAnalysisResult{
		URL:            targetURL,
		Domain:         parsedURL.Host,
		Protocol:       parsedURL.Scheme,
		AnalyzedAt:     time.Now(),
		CategoryScores: make(map[string]float64),
	}

	// 1. Récupération et analyse HTTP de base
	if err := seo.analyzeHTTPResponse(ctx, targetURL, result); err != nil {
		return nil, fmt.Errorf("erreur analyse HTTP: %w", err)
	}

	// 2. Récupération du contenu HTML
	htmlContent, err := seo.fetchHTML(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération HTML: %w", err)
	}

	// 3. Parsing HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("erreur parsing HTML: %w", err)
	}

	// 4. Analyses parallèles pour optimiser les performances
	var wg sync.WaitGroup
	var analysisErrors []error
	var mu sync.Mutex

	// Analyse des balises SEO
	wg.Add(1)
	go func() {
		defer wg.Done()
		tagResult, err := seo.tagAnalyzer.Analyze(doc, htmlContent)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("erreur analyse balises: %w", err))
			mu.Unlock()
			return
		}
		seo.mu.Lock()
		result.TagAnalysis = *tagResult
		seo.mu.Unlock()
	}()

	// Analyse des performances
	wg.Add(1)
	go func() {
		defer wg.Done()
		perfResult, err := seo.performanceAnalyzer.Analyze(ctx, targetURL, htmlContent)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("erreur analyse performance: %w", err))
			mu.Unlock()
			return
		}
		seo.mu.Lock()
		result.PerformanceMetrics = *perfResult
		seo.mu.Unlock()
	}()

	// Audit technique
	wg.Add(1)
	go func() {
		defer wg.Done()
		auditResult, err := seo.technicalAuditor.Audit(ctx, targetURL, doc, htmlContent)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("erreur audit technique: %w", err))
			mu.Unlock()
			return
		}
		seo.mu.Lock()
		result.TechnicalAudit = *auditResult
		seo.mu.Unlock()
	}()

	// Attendre toutes les analyses
	wg.Wait()

	// Vérifier les erreurs d'analyse
	if len(analysisErrors) > 0 {
		log.Printf("Erreurs partielles lors de l'analyse: %v", analysisErrors)
		// Continuer avec les analyses réussies
	}

	// 5. Génération des recommandations
	seo.mu.RLock()
	recommendations := seo.recommendationEngine.GenerateRecommendations(result)
	categoryScores := seo.calculateCategoryScores(result)
	overallScore := seo.calculateOverallScore(categoryScores)
	seo.mu.RUnlock()

	result.Recommendations = recommendations
	result.CategoryScores = categoryScores
	result.OverallScore = overallScore
	result.AnalysisDuration = time.Since(startTime)

	log.Printf("Analyse SEO terminée - Score: %.1f, Durée: %v, Recommandations: %d", 
		result.OverallScore, result.AnalysisDuration, len(result.Recommendations))

	return result, nil
}

// analyzeHTTPResponse analyse la réponse HTTP de base
func (seo *SEOAnalyzer) analyzeHTTPResponse(ctx context.Context, targetURL string, result *SEOAnalysisResult) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", seo.userAgent)

	start := time.Now()
	resp, err := seo.client.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.ResponseTime = responseTime

	return nil
}

// fetchHTML récupère le contenu HTML de la page
func (seo *SEOAnalyzer) fetchHTML(ctx context.Context, targetURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", seo.userAgent)

	resp, err := seo.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// Lire le contenu avec une limite de taille
	buf := make([]byte, 2*1024*1024) // 2MB max
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	return string(buf[:n]), nil
}

// calculateCategoryScores calcule les scores par catégorie
func (seo *SEOAnalyzer) calculateCategoryScores(result *SEOAnalysisResult) map[string]float64 {
	scores := make(map[string]float64)

	// Score des balises (30%)
	scores["tags"] = seo.calculateTagScore(&result.TagAnalysis)

	// Score de performance (25%)
	scores["performance"] = seo.calculatePerformanceScore(&result.PerformanceMetrics)

	// Score technique (25%)
	scores["technical"] = seo.calculateTechnicalScore(&result.TechnicalAudit)

	// Score de base (20%) - status, redirects, etc.
	scores["basics"] = seo.calculateBasicsScore(result)

	return scores
}

// calculateOverallScore calcule le score global
func (seo *SEOAnalyzer) calculateOverallScore(categoryScores map[string]float64) float64 {
	weights := map[string]float64{
		"tags":        0.30,
		"performance": 0.25,
		"technical":   0.25,
		"basics":      0.20,
	}

	totalScore := 0.0
	for category, score := range categoryScores {
		if weight, exists := weights[category]; exists {
			totalScore += score * weight
		}
	}

	return totalScore
}

// calculateTagScore calcule le score des balises
func (seo *SEOAnalyzer) calculateTagScore(tagAnalysis *TagAnalysisResult) float64 {
	score := 0.0

	// Title tag (25%)
	if tagAnalysis.Title.Present {
		score += 0.15
		if tagAnalysis.Title.OptimalLength {
			score += 0.10
		}
	}

	// Meta description (20%)
	if tagAnalysis.MetaDescription.Present {
		score += 0.10
		if tagAnalysis.MetaDescription.OptimalLength {
			score += 0.10
		}
	}

	// Headings structure (25%)
	if tagAnalysis.Headings.H1Count == 1 {
		score += 0.15
	}
	if tagAnalysis.Headings.HasHierarchy {
		score += 0.10
	}

	// Meta tags (20%)
	if tagAnalysis.MetaTags.HasRobots {
		score += 0.05
	}
	if tagAnalysis.MetaTags.HasCanonical {
		score += 0.05
	}
	if tagAnalysis.MetaTags.HasOGTags {
		score += 0.05
	}
	if tagAnalysis.MetaTags.HasTwitterCard {
		score += 0.05
	}

	// Images (10%)
	if tagAnalysis.Images.AltTextCoverage > 0.8 {
		score += 0.10
	} else if tagAnalysis.Images.AltTextCoverage > 0.5 {
		score += 0.05
	}

	return score
}

// calculatePerformanceScore calcule le score de performance
func (seo *SEOAnalyzer) calculatePerformanceScore(perfMetrics *PerformanceMetricsResult) float64 {
	score := 0.0

	// Temps de chargement (40%)
	if perfMetrics.LoadTime < constants.FastLoadTime {
		score += 0.4
	} else if perfMetrics.LoadTime < constants.AcceptableLoadTime {
		score += 0.3
	} else if perfMetrics.LoadTime < constants.SlowLoadTime {
		score += 0.2
	}

	// Taille de la page (30%)
	if perfMetrics.PageSize < 500*1024 { // 500KB
		score += 0.3
	} else if perfMetrics.PageSize < 1024*1024 { // 1MB
		score += 0.2
	} else if perfMetrics.PageSize < 2*1024*1024 { // 2MB
		score += 0.1
	}

	// Optimisations (30%)
	if perfMetrics.HasCompression {
		score += 0.1
	}
	if perfMetrics.HasCaching {
		score += 0.1
	}
	if perfMetrics.OptimizedImages {
		score += 0.1
	}

	return score
}

// calculateTechnicalScore calcule le score technique
func (seo *SEOAnalyzer) calculateTechnicalScore(techAudit *TechnicalAuditResult) float64 {
	score := 0.0

	// HTTPS (20%)
	if techAudit.Security.HasHTTPS {
		score += 0.2
	}

	// Mobile-friendly (20%)
	if techAudit.Mobile.IsResponsive {
		score += 0.1
	}
	if techAudit.Mobile.HasViewport {
		score += 0.1
	}

	// Structure (25%)
	if techAudit.Structure.HasSitemap {
		score += 0.1
	}
	if techAudit.Structure.HasRobotsTxt {
		score += 0.05
	}
	if techAudit.Structure.ValidHTML {
		score += 0.1
	}

	// Accessibilité (20%)
	score += techAudit.Accessibility.Score * 0.2

	// Indexabilité (15%)
	if !techAudit.Indexability.BlockedByRobots {
		score += 0.1
	}
	if !techAudit.Indexability.HasNoIndex {
		score += 0.05
	}

	return score
}

// calculateBasicsScore calcule le score de base
func (seo *SEOAnalyzer) calculateBasicsScore(result *SEOAnalysisResult) float64 {
	score := 0.0

	// Status code (40%)
	if result.StatusCode == 200 {
		score += 0.4
	} else if result.StatusCode >= 300 && result.StatusCode < 400 {
		score += 0.2 // Redirections
	}

	// Response time (30%)
	if result.ResponseTime < constants.FastResponseTime {
		score += 0.3
	} else if result.ResponseTime < constants.AcceptableResponseTime {
		score += 0.2
	} else if result.ResponseTime < constants.SlowResponseTime {
		score += 0.1
	}

	// Protocol (20%)
	if result.Protocol == "https" {
		score += 0.2
	} else if result.Protocol == "http" {
		score += 0.05
	}

	// Domain validity (10%)
	if result.Domain != "" && !strings.Contains(result.Domain, constants.ServerDefaultHost) {
		score += 0.1
	}

	return score
}