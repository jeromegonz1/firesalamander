package crawler

import (
	"context"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD RED PHASE - TESTS NIVEAU SCREAMING FROG
// Tests pour crawler complexe comme septeo.com avec plusieurs pages
// ========================================

// TestCrawler_ScreamingFrogLevel teste un crawl complet niveau professionnel
// Doit trouver plusieurs pages, pas s'arrêter à 1 comme le problème actuel
func TestCrawler_ScreamingFrogLevel(t *testing.T) {
	// ARRANGE : Configuration pour crawl professionnel
	cfg := &config.CrawlerConfig{
		MaxPages:             50,  // Plus que le minimum
		MaxDepth:             4,   // Plus profond que basique
		TimeoutSeconds:       120, // Temps suffisant pour site complexe
		InitialWorkers:       5,   // Plusieurs workers
		MaxWorkers:           10,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     true,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Créer un crawler de niveau professionnel
	professionalCrawler := NewProfessionalCrawler(cfg)
	if professionalCrawler == nil {
		t.Fatal("NewProfessionalCrawler devrait retourner un crawler professionnel, got nil")
	}

	// Test avec une URL qui simule un site complexe
	complexSiteURL := constants.CrawlerTestURLExample // NO HARDCODING
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := professionalCrawler.CrawlProfessionally(ctx, complexSiteURL)

	// ASSERT : Vérifications niveau Screaming Frog
	if err != nil {
		t.Errorf("CrawlProfessionally ne devrait pas échouer: %v", err)
	}
	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	// CRITÈRE PRINCIPAL : Plus de 1 page (fix du problème septeo.com)
	if len(result.Pages) <= 1 {
		t.Errorf("Un site complexe devrait avoir plus de 1 page, got %d", len(result.Pages))
	}

	// Vérifier les fonctionnalités professionnelles
	if result.SEOMetrics == nil {
		t.Error("Un crawler professionnel devrait fournir des métriques SEO")
	}

	if result.LinkAnalysis == nil {
		t.Error("Un crawler professionnel devrait analyser les liens")
	}

	if result.TechnicalIssues == nil {
		t.Error("Un crawler professionnel devrait détecter les problèmes techniques")
	}

	// Vérifier la qualité des données récoltées
	for url, page := range result.Pages {
		if page.Title == "" {
			t.Errorf("Page %s devrait avoir un titre", url)
		}
		if page.StatusCode == 0 {
			t.Errorf("Page %s devrait avoir un code de statut", url)
		}
		if page.ResponseTime == 0 {
			t.Errorf("Page %s devrait avoir un temps de réponse mesuré", url)
		}
	}
}

// TestProfessionalCrawler_HandlesComplexSite teste la gestion de sites complexes
func TestProfessionalCrawler_HandlesComplexSite(t *testing.T) {
	// ARRANGE : Configuration robuste pour sites complexes
	cfg := &config.CrawlerConfig{
		MaxPages:             100,
		MaxDepth:             5,
		TimeoutSeconds:       300, // 5 minutes pour site complexe
		InitialWorkers:       8,
		MaxWorkers:           15,
		MinWorkers:           3,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     true,
		DelayMs:              50, // Politesse pour site réel
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	professionalCrawler := NewProfessionalCrawler(cfg)
	
	// Simuler un site avec plusieurs types de contenu
	complexSiteURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// ACT
	result, err := professionalCrawler.CrawlProfessionally(ctx, complexSiteURL)

	// ASSERT : Capacités professionnelles
	if err != nil {
		t.Errorf("Le crawler professionnel devrait gérer les sites complexes: %v", err)
	}

	// Vérifier la diversité des types de contenu trouvés
	contentTypes := make(map[string]int)
	for _, page := range result.Pages {
		contentTypes[page.ContentType]++
	}

	if len(contentTypes) == 0 {
		t.Error("Le crawler devrait identifier différents types de contenu")
	}

	// Vérifier la détection de liens internes/externes
	if result.LinkAnalysis.InternalLinks == 0 && result.LinkAnalysis.ExternalLinks == 0 {
		t.Error("Le crawler devrait analyser les liens internes et externes")
	}

	// Vérifier l'analyse technique
	if len(result.TechnicalIssues.Issues) == 0 && len(result.Pages) > 5 {
		t.Log("Aucun problème technique détecté - site très propre ou détection insuffisante")
	}
}

// TestProfessionalCrawler_SEOAnalysis teste l'analyse SEO complète
func TestProfessionalCrawler_SEOAnalysis(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             30,
		MaxDepth:             3,
		TimeoutSeconds:       120,
		InitialWorkers:       4,
		MaxWorkers:           8,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     true,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	professionalCrawler := NewProfessionalCrawler(cfg)
	
	// ACT
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := professionalCrawler.CrawlProfessionally(ctx, testURL)

	// ASSERT : Métriques SEO professionnelles
	if err != nil {
		t.Errorf("L'analyse SEO professionnelle ne devrait pas échouer: %v", err)
	}

	seoMetrics := result.SEOMetrics
	if seoMetrics == nil {
		t.Fatal("Les métriques SEO ne devraient pas être nil")
	}

	// Vérifications SEO obligatoires
	if seoMetrics.TitleIssues == nil {
		t.Error("L'analyse devrait vérifier les problèmes de titres")
	}

	if seoMetrics.MetaDescriptionIssues == nil {
		t.Error("L'analyse devrait vérifier les méta descriptions")
	}

	if seoMetrics.HeadingStructure == nil {
		t.Error("L'analyse devrait vérifier la structure des headings")
	}

	if seoMetrics.ImageAltTextIssues == nil {
		t.Error("L'analyse devrait vérifier les alt text des images")
	}

	// Vérifications de scores SEO
	if seoMetrics.OverallScore < 0 || seoMetrics.OverallScore > 100 {
		t.Errorf("Le score SEO global devrait être entre 0 et 100, got %f", seoMetrics.OverallScore)
	}
}

// TestProfessionalCrawler_PerformanceMonitoring teste le monitoring des performances
func TestProfessionalCrawler_PerformanceMonitoring(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             25,
		MaxDepth:             3,
		TimeoutSeconds:       90,
		InitialWorkers:       6,
		MaxWorkers:           12,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false, // Pour tests de performance
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	professionalCrawler := NewProfessionalCrawler(cfg)
	
	// ACT
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	startTime := time.Now()
	result, err := professionalCrawler.CrawlProfessionally(ctx, testURL)
	totalDuration := time.Since(startTime)

	// ASSERT : Métriques de performance
	if err != nil {
		t.Errorf("Le monitoring de performance ne devrait pas échouer: %v", err)
	}

	perfMetrics := result.PerformanceMetrics
	if perfMetrics == nil {
		t.Fatal("Les métriques de performance ne devraient pas être nil")
	}

	// Vérifier les métriques de temps
	if perfMetrics.AverageResponseTime == 0 {
		t.Error("Le temps de réponse moyen devrait être mesuré")
	}

	if perfMetrics.FastestPage == "" {
		t.Error("La page la plus rapide devrait être identifiée")
	}

	if perfMetrics.SlowestPage == "" {
		t.Error("La page la plus lente devrait être identifiée")
	}

	// Vérifier l'efficacité du crawler
	if len(result.Pages) > 0 {
		pagesPerSecond := float64(len(result.Pages)) / totalDuration.Seconds()
		if pagesPerSecond < 0.1 { // Au moins 1 page toutes les 10 secondes
			t.Errorf("Efficacité du crawler trop faible: %.2f pages/sec", pagesPerSecond)
		}
	}
}

// TestProfessionalCrawler_RobustErrorHandling teste la gestion d'erreur robuste
func TestProfessionalCrawler_RobustErrorHandling(t *testing.T) {
	// ARRANGE : Configuration pour tester la robustesse
	cfg := &config.CrawlerConfig{
		MaxPages:             20,
		MaxDepth:             2,
		TimeoutSeconds:       60,
		InitialWorkers:       3,
		MaxWorkers:           6,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     true,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: 20, // Tolérance d'erreur plus élevée pour tests
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	professionalCrawler := NewProfessionalCrawler(cfg)
	
	// Test avec différents scénarios d'erreur
	testScenarios := []string{
		constants.CrawlerTestURLExample,        // URL valide
		"invalid-url-format",                   // URL invalide
		constants.CrawlerTestURLExample + "/404", // Page inexistante
	}

	for _, testURL := range testScenarios {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		
		// ACT
		result, err := professionalCrawler.CrawlProfessionally(ctx, testURL)
		
		// ASSERT : Le crawler ne devrait pas paniquer
		if testURL == "invalid-url-format" {
			if err == nil {
				t.Error("Le crawler devrait échouer avec une URL invalide")
			}
		} else {
			// Pour les autres cas, le crawler devrait être robuste
			if result != nil && result.ErrorSummary != nil {
				errorRate := float64(result.ErrorSummary.TotalErrors) / float64(result.ErrorSummary.TotalRequests) * 100
				if errorRate > 50 {
					t.Errorf("Taux d'erreur trop élevé pour %s: %.1f%%", testURL, errorRate)
				}
			}
		}
		
		cancel()
	}
}

// ========================================
// INTERFACES POUR CRAWLER PROFESSIONNEL (À IMPLÉMENTER)
// ========================================

// IProfessionalCrawler interface pour crawler de niveau professionnel
type IProfessionalCrawler interface {
	CrawlProfessionally(ctx context.Context, startURL string) (*ProfessionalCrawlResult, error)
	SetSEOAnalysisEnabled(enabled bool)
	SetPerformanceMonitoringEnabled(enabled bool)
	GetCrawlerCapabilities() []string
}

// ========================================
// TYPES POUR CRAWLER PROFESSIONNEL (À IMPLÉMENTER)
// ========================================

// ProfessionalCrawler implémentation professionnelle
type ProfessionalCrawler struct {
	// À implémenter selon l'architecture
}

// ProfessionalCrawlResult résultat de crawl professionnel
type ProfessionalCrawlResult struct {
	// Résultats basiques
	StartURL  string                  `json:"start_url"`
	Pages     map[string]*PageResult  `json:"pages"`
	Duration  time.Duration          `json:"duration"`
	
	// Analyses professionnelles
	SEOMetrics         *SEOMetrics         `json:"seo_metrics"`
	LinkAnalysis       *LinkAnalysis       `json:"link_analysis"`
	TechnicalIssues    *TechnicalIssues    `json:"technical_issues"`
	PerformanceMetrics *PerformanceMetrics `json:"performance_metrics"`
	ErrorSummary       *ErrorSummary       `json:"error_summary"`
}

// SEOMetrics métriques SEO professionnelles
type SEOMetrics struct {
	OverallScore           float64              `json:"overall_score"`
	TitleIssues           []string             `json:"title_issues"`
	MetaDescriptionIssues []string             `json:"meta_description_issues"`
	HeadingStructure      map[string]int       `json:"heading_structure"`
	ImageAltTextIssues    []string             `json:"image_alt_text_issues"`
}

// LinkAnalysis analyse des liens
type LinkAnalysis struct {
	InternalLinks int      `json:"internal_links"`
	ExternalLinks int      `json:"external_links"`
	BrokenLinks   []string `json:"broken_links"`
}

// TechnicalIssues problèmes techniques détectés
type TechnicalIssues struct {
	Issues []TechnicalIssue `json:"issues"`
}

// TechnicalIssue problème technique spécifique
type TechnicalIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	URL         string `json:"url"`
}

// PerformanceMetrics métriques de performance
type PerformanceMetrics struct {
	AverageResponseTime time.Duration `json:"average_response_time"`
	FastestPage         string        `json:"fastest_page"`
	SlowestPage         string        `json:"slowest_page"`
	TotalDataTransfer   int64         `json:"total_data_transfer"`
}

// ErrorSummary résumé des erreurs
type ErrorSummary struct {
	TotalRequests int `json:"total_requests"`
	TotalErrors   int `json:"total_errors"`
	ErrorsByType  map[string]int `json:"errors_by_type"`
}

// ========================================
// FONCTIONS À IMPLÉMENTER (STUBS POUR TESTS ROUGES)
// ========================================

// NewProfessionalCrawler crée un crawler de niveau professionnel
func NewProfessionalCrawler(cfg *config.CrawlerConfig) IProfessionalCrawler {
	// TODO: Implémenter selon l'architecture reçue
	// Cette fonction doit retourner nil pour que les tests échouent (RED phase)
	return nil
}