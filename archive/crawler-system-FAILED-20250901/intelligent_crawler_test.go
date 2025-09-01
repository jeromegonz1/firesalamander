package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD RED PHASE - TESTS POUR CRAWLER INTELLIGENT
// Tests qui DOIVENT ÉCHOUER car le code n'existe pas encore
// ========================================

// TestIntelligentCrawler_FixesRaceCondition teste le problème exact trouvé dans les logs
// Cas d'usage : septeo.com - 1 page crawlée, bloque sur "Active jobs: 1"
func TestIntelligentCrawler_FixesRaceCondition(t *testing.T) {
	// ARRANGE : Configuration pour reproduire le problème septeo.com
	cfg := &config.CrawlerConfig{
		MaxPages:             constants.DefaultMaxPages,
		MaxDepth:             constants.DefaultMaxDepth,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false, // Pour les tests
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Créer un IntelligentCrawler (interface à implémenter)
	crawler := NewIntelligentCrawler(cfg)
	if crawler == nil {
		t.Fatal("NewIntelligentCrawler devrait retourner un crawler intelligent, got nil")
	}

	// Test de l'interface ICrawlerEngine
	_, ok := crawler.(ICrawlerEngine)
	if !ok {
		t.Fatal("IntelligentCrawler devrait implémenter ICrawlerEngine")
	}

	// Test avec URL mock qui simule septeo.com
	mockURL := constants.CrawlerTestURLExample // Pas de hardcoding
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := crawler.CrawlWithIntelligence(ctx, mockURL)

	// ASSERT : Vérifications anti-race condition
	if err != nil {
		t.Errorf("CrawlWithIntelligence ne devrait pas échouer, got error: %v", err)
	}
	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}
	if len(result.Pages) == 0 {
		t.Error("Le crawler devrait trouver au moins 1 page")
	}

	// Vérification que la race condition est résolue
	// Le crawler doit finir proprement sans rester bloqué
	if result.Duration > 25*time.Second {
		t.Errorf("Le crawler ne devrait pas prendre plus de 25s, got %v", result.Duration)
	}
}

// TestIntelligentCrawler_HasSmartTermination teste la logique de terminaison intelligente
func TestIntelligentCrawler_HasSmartTermination(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             5, // Limite basse pour tester la terminaison
		MaxDepth:             2,
		TimeoutSeconds:       60,
		InitialWorkers:       2,
		MaxWorkers:           5,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Le SmartTerminationController doit être injecté
	crawler := NewIntelligentCrawler(cfg)
	if crawler == nil {
		t.Fatal("NewIntelligentCrawler devrait retourner un crawler intelligent, got nil")
	}
	
	termController := crawler.GetTerminationController()
	if termController == nil {
		t.Fatal("IntelligentCrawler devrait avoir un TerminationController")
	}

	// Vérifier l'interface ITerminationController
	_, ok := termController.(ITerminationController)
	if !ok {
		t.Fatal("TerminationController devrait implémenter ITerminationController")
	}

	// Test des conditions de terminaison
	conditions := termController.GetTerminationConditions()
	if len(conditions) == 0 {
		t.Error("TerminationController devrait avoir des conditions de terminaison")
	}

	// Vérifier que les conditions incluent les jobs actifs
	hasJobsCondition := false
	for _, condition := range conditions {
		if condition.Type() == "active_jobs_zero" {
			hasJobsCondition = true
			break
		}
	}
	if !hasJobsCondition {
		t.Error("Les conditions doivent inclure 'active_jobs_zero'")
	}
}

// TestIntelligentCrawler_ImplementsInterfaces vérifie que toutes les interfaces sont implémentées
func TestIntelligentCrawler_ImplementsInterfaces(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             constants.DefaultMaxPages,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		InitialWorkers:       constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT
	crawler := NewIntelligentCrawler(cfg)
	if crawler == nil {
		t.Fatal("NewIntelligentCrawler devrait retourner un crawler intelligent, got nil")
	}

	// ASSERT : Test des interfaces obligatoires
	if _, ok := crawler.(ICrawlerEngine); !ok {
		t.Error("IntelligentCrawler doit implémenter ICrawlerEngine")
	}

	// Test que le crawler a un service de découverte d'URLs
	urlDiscovery := crawler.GetURLDiscoveryService()
	if urlDiscovery == nil {
		t.Fatal("IntelligentCrawler devrait avoir un URLDiscoveryService")
	}

	if _, ok := urlDiscovery.(IURLDiscoveryService); !ok {
		t.Error("URLDiscoveryService doit implémenter IURLDiscoveryService")
	}
}

// TestIntelligentCrawler_AtomicJobCounter teste que le compteur de jobs est atomique
func TestIntelligentCrawler_AtomicJobCounter(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             10,
		TimeoutSeconds:       30,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		InitialWorkers:       3, // Plusieurs workers pour tester l'atomicité
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	crawler := NewIntelligentCrawler(cfg)
	if crawler == nil {
		t.Fatal("NewIntelligentCrawler devrait retourner un crawler intelligent, got nil")
	}

	// ACT : Accès au compteur de jobs atomique
	jobCounter := crawler.GetJobCounter()
	if jobCounter == nil {
		t.Fatal("IntelligentCrawler devrait exposer un JobCounter atomique")
	}

	// Test des opérations atomiques
	initial := jobCounter.Get()
	jobCounter.Add(5)
	after_add := jobCounter.Get()
	jobCounter.Sub(3)
	after_sub := jobCounter.Get()

	// ASSERT
	if after_add != initial+5 {
		t.Errorf("JobCounter.Add() échoué: expected %d, got %d", initial+5, after_add)
	}
	if after_sub != after_add-3 {
		t.Errorf("JobCounter.Sub() échoué: expected %d, got %d", after_add-3, after_sub)
	}
}

// ========================================
// INTERFACES À IMPLÉMENTER (TDD CONTRACT)
// ========================================

// ========================================
// TDD TESTS POUR CRAWLER MULTI-PAGES
// Phase RED : Tests qui DOIVENT échouer
// ========================================

// TestBasicURLDiscoveryService_DiscoverFromHTML_MultipleLinks - Test découverte de liens
func TestBasicURLDiscoveryService_DiscoverFromHTML_MultipleLinks(t *testing.T) {
	service := &basicURLDiscoveryService{}
	baseURL := "https://example.com/page1"
	
	// HTML avec plusieurs types de liens
	html := `
		<!DOCTYPE html>
		<html>
		<body>
			<a href="/about">About Us</a>
			<a href="./contact">Contact</a>
			<a href="../services">Services</a>
			<a href="https://example.com/products">Products</a>
			<a href="https://external.com/link">External</a>
			<a href="#anchor">Anchor</a>
			<a href="javascript:void(0)">JS Link</a>
		</body>
		</html>
	`
	
	// WHEN - Découvrir les URLs
	urls, err := service.DiscoverFromHTML(html, baseURL)
	if err != nil {
		t.Fatalf("DiscoverFromHTML failed: %v", err)
	}
	
	// THEN - DOIT découvrir les liens internes convertis en URLs absolues
	expectedURLs := []string{
		"https://example.com/about",     // Lien relatif /about
		"https://example.com/contact",   // Lien relatif ./contact  
		"https://example.com/services",  // Lien relatif ../services
		"https://example.com/products",  // Lien absolu même domaine
	}
	
	if len(urls) != len(expectedURLs) {
		t.Errorf("❌ TDD RED: Expected %d URLs, got %d. URLs found: %v", 
			len(expectedURLs), len(urls), urls)
	}
	
	// Vérifier que les URLs sont converties correctement
	for _, expectedURL := range expectedURLs {
		found := false
		for _, actualURL := range urls {
			if actualURL == expectedURL {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("❌ TDD RED: Expected URL %s not found in discovered URLs: %v", 
				expectedURL, urls)
		}
	}
	
	// Vérifier que les liens externes sont filtrés
	for _, url := range urls {
		if strings.Contains(url, "external.com") {
			t.Errorf("❌ TDD RED: External URL should be filtered: %s", url)
		}
	}
}

// TestBasicURLDiscoveryService_RelativeURLConversion - Test conversion URLs relatives
func TestBasicURLDiscoveryService_RelativeURLConversion(t *testing.T) {
	service := &basicURLDiscoveryService{}
	
	testCases := []struct {
		name     string
		baseURL  string
		href     string
		expected string
	}{
		{
			name:     "Absolute path",
			baseURL:  "https://example.com/dir/page.html",
			href:     "/contact",
			expected: "https://example.com/contact",
		},
		{
			name:     "Relative path current dir",
			baseURL:  "https://example.com/dir/page.html",
			href:     "./about",
			expected: "https://example.com/dir/about",
		},
		{
			name:     "Relative path parent dir",
			baseURL:  "https://example.com/dir/subdir/page.html",
			href:     "../services",
			expected: "https://example.com/dir/services",
		},
		{
			name:     "Query parameters",
			baseURL:  "https://example.com/",
			href:     "/search?q=test",
			expected: "https://example.com/search?q=test",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			html := `<a href="` + tc.href + `">Link</a>`
			
			urls, err := service.DiscoverFromHTML(html, tc.baseURL)
			if err != nil {
				t.Fatalf("DiscoverFromHTML failed: %v", err)
			}
			
			if len(urls) != 1 {
				t.Fatalf("❌ TDD RED: Expected 1 URL, got %d", len(urls))
			}
			
			if urls[0] != tc.expected {
				t.Errorf("❌ TDD RED: Expected %s, got %s", tc.expected, urls[0])
			}
		})
	}
}

// TestDebugURLDiscovery - Test debug pour voir ce qui se passe
func TestDebugURLDiscovery(t *testing.T) {
	service := &basicURLDiscoveryService{}
	
	html := `
		<!DOCTYPE html>
		<html>
		<body>
			<h1>Welcome</h1>
			<a href="/about">About</a>
			<a href="/products">Products</a>
			<a href="/contact">Contact</a>
		</body>
		</html>
	`
	
	baseURL := "http://127.0.0.1:8080"
	
	urls, err := service.DiscoverFromHTML(html, baseURL)
	if err != nil {
		t.Fatalf("DiscoverFromHTML failed: %v", err)
	}
	
	t.Logf("DEBUG: Base URL: %s", baseURL)
	t.Logf("DEBUG: Discovered URLs: %v", urls)
	t.Logf("DEBUG: Count: %d", len(urls))
	
	// This should find 3 URLs
	if len(urls) != 3 {
		t.Errorf("Expected 3 URLs, got %d", len(urls))
	}
}

// TestIntelligentCrawler_CrawlWithIntelligence_MultiplePages - Test crawl multi-pages
func TestIntelligentCrawler_CrawlWithIntelligence_MultiplePages(t *testing.T) {
	// Mock HTTP server avec plusieurs pages liées
	pageContent := map[string]string{
		"/": `
			<!DOCTYPE html>
			<html>
			<head><title>Home Page</title></head>
			<body>
				<h1>Welcome</h1>
				<a href="/about">About</a>
				<a href="/products">Products</a>
				<a href="/contact">Contact</a>
			</body>
			</html>
		`,
		"/about": `
			<!DOCTYPE html>
			<html>
			<head><title>About Page</title></head>
			<body>
				<h1>About Us</h1>
				<a href="/">Home</a>
				<a href="/team">Our Team</a>
			</body>
			</html>
		`,
		"/products": `
			<!DOCTYPE html>
			<html>
			<head><title>Products Page</title></head>
			<body>
				<h1>Our Products</h1>
				<a href="/">Home</a>
				<a href="/product/item1">Product 1</a>
			</body>
			</html>
		`,
		"/contact": `
			<!DOCTYPE html>
			<html>
			<head><title>Contact Page</title></head>
			<body>
				<h1>Contact Us</h1>
				<a href="/">Home</a>
			</body>
			</html>
		`,
	}
	
	// Créer un serveur de test
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, exists := pageContent[r.URL.Path]
		if !exists {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(content))
	}))
	defer server.Close()
	
	// Créer le crawler intelligent avec configuration de test
	cfg := &config.CrawlerConfig{
		MaxPages:         10,
		MaxDepth:         3,
		TimeoutSeconds:   30,
		UserAgent:        "FireSalamander-Test/1.0",
		MaxWorkers:       2,
		InitialWorkers:   2,
		MinWorkers:       1,
		FastThresholdMs:  constants.DefaultFastThresholdMs,
		SlowThresholdMs:  constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}
	
	crawler := NewIntelligentCrawler(cfg)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// WHEN - Crawl du site complet
	t.Logf("DEBUG: Starting crawl of: %s", server.URL)
	
	result, err := crawler.CrawlWithIntelligence(ctx, server.URL)
	if err != nil {
		t.Fatalf("CrawlWithIntelligence failed: %v", err)
	}
	
	// DEBUG: Log what we actually got
	crawledURLs := getPageURLs(result.Pages)
	t.Logf("DEBUG: Crawl completed. Found %d pages:", len(result.Pages))
	for i, pageURL := range crawledURLs {
		pageResult := result.Pages[pageURL]
		t.Logf("DEBUG:   %d. %s (depth: %d, status: %d)", 
			i+1, pageURL, pageResult.Depth, pageResult.StatusCode)
		
		// Test URL discovery on the first page
		if i == 0 && pageResult.Body != "" {
			// Log the actual body content to see what we're working with
			t.Logf("DEBUG:      Page body (first 500 chars): %s", 
				pageResult.Body[:min(500, len(pageResult.Body))])
			
			service := &basicURLDiscoveryService{}
			discoveredURLs, discErr := service.DiscoverFromHTML(pageResult.Body, pageURL)
			if discErr == nil {
				t.Logf("DEBUG:      Discovered URLs from this page: %v", discoveredURLs)
			} else {
				t.Logf("DEBUG:      URL discovery error: %v", discErr)
			}
		}
	}
	
	// THEN - DOIT avoir crawlé plusieurs pages (pas juste 1)
	if len(result.Pages) <= 1 {
		t.Errorf("❌ TDD RED: Expected multiple pages crawled, got %d pages: %v", 
			len(result.Pages), getPageURLs(result.Pages))
	}
	
	// DOIT avoir crawlé au moins les pages principales
	expectedPages := []string{
		server.URL + "/",
		server.URL + "/about", 
		server.URL + "/products",
		server.URL + "/contact",
	}
	
	// crawledURLs already declared above, just use it
	for _, expectedURL := range expectedPages {
		found := false
		for _, crawledURL := range crawledURLs {
			if crawledURL == expectedURL {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("❌ TDD RED: Expected page %s not crawled. Crawled: %v", 
				expectedURL, crawledURLs)
		}
	}
	
	// Vérifier que chaque page a un contenu (sauf les 404s)
	for url, pageResult := range result.Pages {
		if pageResult.StatusCode == 200 {
			if pageResult.Title == "" {
				t.Errorf("❌ TDD RED: Page %s has empty title (status: %d)", url, pageResult.StatusCode)
			}
			if pageResult.Body == "" {
				t.Errorf("❌ TDD RED: Page %s has empty body (status: %d)", url, pageResult.StatusCode)
			}
		}
		// 404 pages are OK to have empty title/body
	}
}

// TestIntelligentCrawler_DomainFiltering - Test filtrage par domaine
func TestIntelligentCrawler_DomainFiltering(t *testing.T) {
	// HTML avec liens externes
	pageWithExternalLinks := `
		<!DOCTYPE html>
		<html>
		<body>
			<a href="/internal">Internal Link</a>
			<a href="https://example.com/page">Same Domain</a>
			<a href="https://external.com/page">External Domain</a>
			<a href="https://sub.example.com/page">Subdomain</a>
		</body>
		</html>
	`
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(pageWithExternalLinks))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		MaxPages:         10,
		MaxDepth:         2,
		TimeoutSeconds:   10,
		MaxWorkers:       2,
		InitialWorkers:   2,
		MinWorkers:       1,
		UserAgent:        "FireSalamander-Test/1.0",
		FastThresholdMs:  constants.DefaultFastThresholdMs,
		SlowThresholdMs:  constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}
	
	crawler := NewIntelligentCrawler(cfg)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// WHEN - Crawl avec filtrage domaine
	result, err := crawler.CrawlWithIntelligence(ctx, server.URL)
	if err != nil {
		t.Fatalf("CrawlWithIntelligence failed: %v", err)
	}
	
	// THEN - DOIT seulement crawler le même domaine
	crawledURLs := getPageURLs(result.Pages)
	for _, url := range crawledURLs {
		if strings.Contains(url, "external.com") {
			t.Errorf("❌ TDD RED: External domain should not be crawled: %s", url)
		}
	}
}

// TestIntelligentCrawler_DepthLimiting - Test limitation de profondeur configurable
func TestIntelligentCrawler_DepthLimiting(t *testing.T) {
	// Configuration avec profondeur limitée à 1
	cfg := &config.CrawlerConfig{
		MaxPages:         10,
		MaxDepth:         1, // ❌ HARDCODÉ actuellement à 2 dans le code
		TimeoutSeconds:   10,
		MaxWorkers:       2,
		InitialWorkers:   2,
		MinWorkers:       1,
		UserAgent:        "FireSalamander-Test/1.0",
		FastThresholdMs:  constants.DefaultFastThresholdMs,
		SlowThresholdMs:  constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}
	
	pageContent := map[string]string{
		"/": `<a href="/level1">Level 1</a>`,
		"/level1": `<a href="/level2">Level 2</a>`,
		"/level2": `<h1>Level 2 Page</h1>`,
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, exists := pageContent[r.URL.Path]
		if !exists {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	}))
	defer server.Close()
	
	crawler := NewIntelligentCrawler(cfg)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// WHEN - Crawl avec profondeur limitée
	result, err := crawler.CrawlWithIntelligence(ctx, server.URL)
	if err != nil {
		t.Fatalf("CrawlWithIntelligence failed: %v", err)
	}
	
	// THEN - NE DOIT PAS crawler au-delà de la profondeur configurée
	crawledURLs := getPageURLs(result.Pages)
	for _, url := range crawledURLs {
		if strings.Contains(url, "/level2") {
			t.Errorf("❌ TDD RED: Level 2 should not be crawled with MaxDepth=1. Crawled: %v", 
				crawledURLs)
		}
	}
}

// TestIntelligentCrawler_SitemapDiscovery - Test découverte via sitemap.xml
func TestIntelligentCrawler_SitemapDiscovery(t *testing.T) {
	sitemapContent := `<?xml version="1.0" encoding="UTF-8"?>
		<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			<url><loc>https://example.com/</loc></url>
			<url><loc>https://example.com/about</loc></url>
			<url><loc>https://example.com/products</loc></url>
		</urlset>`
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sitemap.xml" {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(sitemapContent))
			return
		}
		// Page normale
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><head><title>Test Page</title></head><body><h1>Content</h1></body></html>"))
	}))
	defer server.Close()
	
	service := &basicURLDiscoveryService{}
	
	// WHEN - Découverte via sitemap
	urls, err := service.DiscoverFromSitemap(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("DiscoverFromSitemap failed: %v", err)
	}
	
	// THEN - DOIT découvrir les URLs du sitemap
	expectedURLs := []string{
		server.URL + "/",
		server.URL + "/about",
		server.URL + "/products",
	}
	
	if len(urls) != len(expectedURLs) {
		t.Errorf("❌ TDD RED: Expected %d URLs from sitemap, got %d: %v", 
			len(expectedURLs), len(urls), urls)
	}
}

// ========================================
// TESTS DE RÉGRESSION POUR PROBLÈME ACTUEL
// ========================================

// TestIntelligentCrawler_CurrentBehavior_OnlyOnePage - Documente le comportement actuel
func TestIntelligentCrawler_CurrentBehavior_OnlyOnePage(t *testing.T) {
	// Ce test documente le problème actuel: seulement 1 page analysée
	
	pageContent := map[string]string{
		"/": `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>Home</h1>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
			</body>
			</html>
		`,
		"/about": `<html><body><h1>About</h1></body></html>`,
		"/contact": `<html><body><h1>Contact</h1></body></html>`,
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, exists := pageContent[r.URL.Path]
		if !exists {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		MaxPages:         10,
		MaxDepth:         3,
		TimeoutSeconds:   30,
		MaxWorkers:       2,
		InitialWorkers:   2,
		MinWorkers:       1,
		UserAgent:        "FireSalamander-Test/1.0",
		FastThresholdMs:  constants.DefaultFastThresholdMs,
		SlowThresholdMs:  constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}
	
	crawler := NewIntelligentCrawler(cfg)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// WHEN - Crawl actuel
	result, err := crawler.CrawlWithIntelligence(ctx, server.URL)
	if err != nil {
		t.Fatalf("CrawlWithIntelligence failed: %v", err)
	}
	
	// DOCUMENTER LE PROBLÈME ACTUEL
	t.Logf("📋 CURRENT BEHAVIOR: Pages crawled = %d", len(result.Pages))
	t.Logf("📋 CURRENT BEHAVIOR: Pages = %v", getPageURLs(result.Pages))
	
	// ASSERTION TEMPORAIRE - Ce test passera maintenant mais échouera après correction
	if len(result.Pages) == 1 {
		t.Logf("✅ DOCUMENTED: Current crawler only analyzes 1 page (needs fixing)")
	} else {
		t.Errorf("🎉 FIXED: Crawler now analyzes %d pages! Update this test", len(result.Pages))
	}
}

// Helper functions pour les tests

// getPageURLs extrait les URLs des pages crawlées
func getPageURLs(pages map[string]*PageResult) []string {
	var urls []string
	for url := range pages {
		urls = append(urls, url)
	}
	return urls
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ========================================
// TYPES POUR LES TESTS - Utilisation des interfaces du fichier d'implémentation
// ========================================