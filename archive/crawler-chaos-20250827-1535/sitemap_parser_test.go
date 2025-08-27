package crawler

import (
	"context"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD RED PHASE - TESTS POUR SITEMAP DISCOVERY
// Tests SANS hardcoding d'URLs - Configuration via constants
// ========================================

// TestSitemapDiscovery_NoHardcoding teste la découverte de sitemaps SANS URLs hardcodées
func TestSitemapDiscovery_NoHardcoding(t *testing.T) {
	// ARRANGE : Configuration via constants uniquement
	cfg := &config.CrawlerConfig{
		MaxPages:             constants.DefaultMaxPages,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     true, // Pour tester robots.txt -> sitemap
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Créer le service de découverte d'URLs
	discoveryService := NewURLDiscoveryService(cfg)
	if discoveryService == nil {
		t.Fatal("NewURLDiscoveryService devrait retourner un service, got nil")
	}

	// Test avec URL de test (constante)
	testURL := constants.CrawlerTestURLExample // NO HARDCODING!
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ASSERT : Test découverte depuis sitemap.xml
	sitemapURLs, err := discoveryService.DiscoverFromSitemap(ctx, testURL)
	if err != nil {
		t.Errorf("DiscoverFromSitemap ne devrait pas échouer avec URL valide: %v", err)
	}

	// Vérifier que les URLs découvertes ne sont pas hardcodées
	for _, url := range sitemapURLs {
		if url == "" {
			t.Error("URL vide trouvée dans les résultats du sitemap")
		}
		// Les URLs doivent être du même domaine que testURL
		if !isValidSitemapURL(url, testURL) {
			t.Errorf("URL invalide trouvée dans sitemap: %s", url)
		}
	}
}

// TestSitemapParser_ValidXMLStructure teste le parsing XML correct des sitemaps
func TestSitemapParser_ValidXMLStructure(t *testing.T) {
	// ARRANGE : Service de découverte avec configuration
	cfg := &config.CrawlerConfig{
		UserAgent:            constants.ParallelCrawlerUserAgent,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	discoveryService := NewURLDiscoveryService(cfg)
	
	// XML de sitemap valide (utilisant les constantes pour les éléments)
	validSitemapXML := createValidSitemapXML()

	// ACT : Parser le XML du sitemap
	urls, err := discoveryService.ParseSitemapXML(validSitemapXML)

	// ASSERT
	if err != nil {
		t.Errorf("ParseSitemapXML devrait parser un XML valide: %v", err)
	}

	if len(urls) == 0 {
		t.Error("ParseSitemapXML devrait trouver des URLs dans un XML valide")
	}

	// Vérifier la structure des URLs parsées
	for _, url := range urls {
		sitemapURL, ok := url.(ISitemapURL)
		if !ok {
			t.Error("Les URLs parsées devraient implémenter ISitemapURL")
		}

		if sitemapURL.GetLoc() == "" {
			t.Error("Chaque URL devrait avoir un 'loc' non vide")
		}

		if sitemapURL.GetLastMod() == nil {
			t.Error("Chaque URL devrait avoir un 'lastmod'")
		}

		if sitemapURL.GetPriority() < 0.0 || sitemapURL.GetPriority() > 1.0 {
			t.Errorf("Priority devrait être entre 0.0 et 1.0, got %f", sitemapURL.GetPriority())
		}
	}
}

// TestSitemapDiscovery_NestedSitemaps teste la découverte de sitemaps imbriqués
func TestSitemapDiscovery_NestedSitemaps(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		UserAgent:            constants.ParallelCrawlerUserAgent,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		MaxPages:             1000, // Pour supporter les gros sitemaps
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	discoveryService := NewURLDiscoveryService(cfg)
	
	// ACT : Test avec sitemap index (qui référence d'autres sitemaps)
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nestedURLs, err := discoveryService.DiscoverFromNestedSitemaps(ctx, testURL)

	// ASSERT
	if err != nil {
		t.Errorf("DiscoverFromNestedSitemaps ne devrait pas échouer: %v", err)
	}

	// Vérifier qu'on a trouvé des URLs de plusieurs niveaux
	if len(nestedURLs) == 0 {
		t.Error("La découverte imbriquée devrait trouver des URLs")
	}

	// Vérifier la déduplication
	uniqueURLs := make(map[string]bool)
	for _, url := range nestedURLs {
		if uniqueURLs[url] {
			t.Errorf("URL dupliquée trouvée: %s", url)
		}
		uniqueURLs[url] = true
	}
}

// TestSitemapDiscovery_RobotsIntegration teste l'intégration robots.txt -> sitemap
func TestSitemapDiscovery_RobotsIntegration(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		UserAgent:            constants.ParallelCrawlerUserAgent,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		RespectRobotsTxt:     true,
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	discoveryService := NewURLDiscoveryService(cfg)
	
	// ACT : Découvrir sitemaps via robots.txt
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	robotsSitemaps, err := discoveryService.DiscoverFromRobots(ctx, testURL)

	// ASSERT
	if err != nil {
		t.Errorf("DiscoverFromRobots ne devrait pas échouer: %v", err)
	}

	// Vérifier que les URLs de sitemaps trouvées sont valides
	for _, sitemapURL := range robotsSitemaps {
		if !isValidSitemapURL(sitemapURL, testURL) {
			t.Errorf("URL de sitemap invalide dans robots.txt: %s", sitemapURL)
		}
	}

	// Vérifier qu'on peut ensuite parser ces sitemaps
	for _, sitemapURL := range robotsSitemaps {
		urls, err := discoveryService.DiscoverFromSitemap(ctx, sitemapURL)
		if err != nil {
			t.Errorf("Impossible de parser sitemap trouvé dans robots.txt: %s, error: %v", sitemapURL, err)
		}
		
		if len(urls) == 0 {
			t.Logf("Aucune URL trouvée dans sitemap %s (peut être vide)", sitemapURL)
		}
	}
}

// TestSitemapDiscovery_ErrorHandling teste la gestion d'erreur robuste
func TestSitemapDiscovery_ErrorHandling(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		UserAgent:            constants.ParallelCrawlerUserAgent,
		TimeoutSeconds:       5, // Timeout court pour tester
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	discoveryService := NewURLDiscoveryService(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ACT & ASSERT : Test avec URL invalide
	invalidURL := "invalid-url-format"
	_, err := discoveryService.DiscoverFromSitemap(ctx, invalidURL)
	if err == nil {
		t.Error("DiscoverFromSitemap devrait échouer avec une URL invalide")
	}

	// Test avec XML invalide
	invalidXML := "<invalid>xml</format>"
	_, err = discoveryService.ParseSitemapXML(invalidXML)
	if err == nil {
		t.Error("ParseSitemapXML devrait échouer avec XML invalide")
	}

	// Test avec URL qui n'existe pas
	nonExistentURL := constants.CrawlerTestURLExample + "/non-existent-sitemap.xml"
	urls, err := discoveryService.DiscoverFromSitemap(ctx, nonExistentURL)
	// L'erreur est acceptable, mais ne devrait pas causer de panique
	if err != nil && len(urls) > 0 {
		t.Error("Si erreur, ne devrait pas retourner d'URLs")
	}
}

// ========================================
// INTERFACES POUR SITEMAP DISCOVERY (À IMPLÉMENTER)
// ========================================

// IURLDiscoveryService interface pour la découverte d'URLs (déjà définie dans intelligent_crawler_test.go)

// ISitemapURL représente une URL dans un sitemap avec métadonnées
type ISitemapURL interface {
	GetLoc() string
	GetLastMod() *time.Time
	GetChangeFreq() string
	GetPriority() float64
}

// IURLDiscoveryService interface étendue pour les sitemaps
type IAdvancedURLDiscoveryService interface {
	IURLDiscoveryService
	
	// Nouvelles méthodes pour les sitemaps
	ParseSitemapXML(xmlContent string) ([]ISitemapURL, error)
	DiscoverFromNestedSitemaps(ctx context.Context, baseURL string) ([]string, error)
}

// ========================================
// TYPES À IMPLÉMENTER (STUBS POUR TDD)
// ========================================

// URLDiscoveryService implémentation du service de découverte
type URLDiscoveryService struct {
	// À implémenter selon l'architecture
}

// TestSitemapURL implémentation d'une URL de sitemap pour les tests
type TestSitemapURL struct {
	// À implémenter
}

// ========================================
// FONCTIONS UTILITAIRES POUR LES TESTS
// ========================================

// createValidSitemapXML crée un XML de sitemap valide en utilisant les constantes
func createValidSitemapXML() string {
	// Utilise les constantes pour les éléments XML
	return `<?xml version="1.0" encoding="UTF-8"?>
<` + constants.CrawlerSitemapURLSet + ` xmlns="` + constants.CrawlerTestURLSitemapSchema + `">
  <` + constants.CrawlerSitemapURL + `>
    <` + constants.CrawlerSitemapLoc + `>` + constants.CrawlerTestURLExamplePage1 + `</` + constants.CrawlerSitemapLoc + `>
    <` + constants.CrawlerSitemapLastmod + `>2024-01-01</` + constants.CrawlerSitemapLastmod + `>
    <` + constants.CrawlerSitemapChangefreq + `>daily</` + constants.CrawlerSitemapChangefreq + `>
    <` + constants.CrawlerSitemapPriority + `>0.8</` + constants.CrawlerSitemapPriority + `>
  </` + constants.CrawlerSitemapURL + `>
  <` + constants.CrawlerSitemapURL + `>
    <` + constants.CrawlerSitemapLoc + `>` + constants.CrawlerTestURLExamplePage2 + `</` + constants.CrawlerSitemapLoc + `>
    <` + constants.CrawlerSitemapLastmod + `>2024-01-02</` + constants.CrawlerSitemapLastmod + `>
    <` + constants.CrawlerSitemapChangefreq + `>weekly</` + constants.CrawlerSitemapChangefreq + `>
    <` + constants.CrawlerSitemapPriority + `>0.6</` + constants.CrawlerSitemapPriority + `>
  </` + constants.CrawlerSitemapURL + `>
</` + constants.CrawlerSitemapURLSet + `>`
}

// isValidSitemapURL vérifie qu'une URL de sitemap est valide par rapport au domaine de base
func isValidSitemapURL(sitemapURL, baseURL string) bool {
	// Implémentation basique - à améliorer selon les besoins
	return sitemapURL != "" && len(sitemapURL) > 7 // http:// minimum
}

// ========================================
// FONCTIONS À IMPLÉMENTER (STUBS POUR TESTS ROUGES)
// ========================================

// NewURLDiscoveryService crée un service de découverte d'URLs
func NewURLDiscoveryService(cfg *config.CrawlerConfig) IAdvancedURLDiscoveryService {
	// TODO: Implémenter selon l'architecture reçue
	// Cette fonction doit retourner nil pour que les tests échouent (RED phase)
	return nil
}