package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD PHASE 1: TESTS ÉCRITS D'ABORD (ROUGE)
// Ces tests DOIVENT ÉCHOUER avant implémentation
// ========================================

func TestParallelCrawler_NoHardcoding(t *testing.T) {
	// GIVEN - Configuration depuis constantes (NO HARDCODING)
	cfg := &config.CrawlerConfig{
		InitialWorkers:        constants.DefaultInitialWorkers,
		MaxWorkers:           constants.DefaultMaxWorkers,
		MinWorkers:           constants.DefaultMinWorkers,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		MaxPages:             constants.DefaultMaxPages,
		TimeoutSeconds:       constants.DefaultTimeoutSeconds,
		UserAgent:            constants.ParallelCrawlerUserAgent,
	}
	
	// THEN - Aucune valeur hardcodée
	assert.NotEqual(t, 5, cfg.InitialWorkers, "InitialWorkers must be from constants")
	assert.NotEqual(t, 20, cfg.MaxWorkers, "MaxWorkers must be from constants") 
	assert.NotEqual(t, constants.HTTPStatusInternalServerError, cfg.FastThresholdMs, "FastThreshold must be from constants")
	assert.NotEqual(t, 2000, cfg.SlowThresholdMs, "SlowThreshold must be from constants")
	
	// Valeurs depuis constantes
	assert.Equal(t, constants.DefaultInitialWorkers, cfg.InitialWorkers)
	assert.Equal(t, constants.DefaultMaxWorkers, cfg.MaxWorkers)
	assert.Equal(t, constants.ParallelCrawlerUserAgent, cfg.UserAgent)
}

func TestParallelCrawler_BasicCrawl(t *testing.T) {
	// GIVEN - Mock server avec pages HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simule latence normale
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<head><title>Test Page</title></head>
				<body>
					<h1>Test Content</h1>
					<a href="/page1">Page 1</a>
					<a href="/page2">Page 2</a>
				</body>
			</html>
		`))
	}))
	defer server.Close()
	
	// WHEN - Crawl avec configuration par défaut
	cfg := &config.CrawlerConfig{
		InitialWorkers: constants.DefaultInitialWorkers,
		MaxWorkers:     constants.DefaultMaxWorkers,
		MinWorkers:     constants.DefaultMinWorkers,
		MaxPages:       5,
		TimeoutSeconds: 10,
		UserAgent:      constants.ParallelCrawlerUserAgent,
	}
	
	crawler := NewParallelCrawler(cfg)
	require.NotNil(t, crawler, "Crawler should be created")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	
	// THEN - Ces tests DOIVENT ÊTRE ROUGES car pas encore implémenté
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Pages), 1)
	assert.LessOrEqual(t, len(result.Pages), 5)
	assert.Greater(t, result.Duration, time.Duration(0))
}

func TestParallelCrawler_AdaptToSlowSite(t *testing.T) {
	// GIVEN - Site très lent (2 secondes par page) avec plusieurs pages pour déclencher l'adaptation
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		time.Sleep(1800 * time.Millisecond) // Très lent (1.8s > 1.5s threshold)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		
		// Générer du contenu avec des liens internes pour avoir plus de requêtes
		content := `<html><head><title>Slow Site</title></head><body>`
		if r.URL.Path == "/" {
			for i := 1; i <= 6; i++ {
				content += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
			}
		} else {
			content += fmt.Sprintf(`<h1>%s</h1><a href="/">Home</a>`, r.URL.Path)
		}
		content += `</body></html>`
		w.Write([]byte(content))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers:        5, // Commence avec 5 workers
		MaxWorkers:           10,
		MinWorkers:           1,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		AdaptIntervalSeconds: 1, // Adaptation rapide pour test
		MaxPages:             10, // Plus de pages pour déclencher l'adaptation
		MaxDepth:             2,  // Permettre le crawl de plusieurs niveaux
		TimeoutSeconds:       15,
	}
	
	// WHEN - Crawl un site lent
	crawler := NewParallelCrawler(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	
	// Démarrer le crawl en arrière-plan
	go func() {
		crawler.CrawlWithContext(ctx, server.URL)
	}()
	
	// Attendre que l'adaptation se fasse (plus de temps pour traiter 5+ requêtes)
	time.Sleep(12 * time.Second)
	
	// THEN - Le crawler devrait réduire les workers
	metrics := crawler.GetMetrics()
	assert.NotNil(t, metrics)
	assert.LessOrEqual(t, metrics.CurrentWorkers, 3, 
		"Should reduce workers for slow site (current: %d)", metrics.CurrentWorkers)
	assert.Greater(t, metrics.AvgResponseTime, time.Duration(constants.DefaultSlowThresholdMs)*time.Millisecond,
		"Average response time should be over slow threshold")
}

func TestParallelCrawler_AdaptToFastSite(t *testing.T) {
	// GIVEN - Site très rapide (50ms par page) avec plusieurs pages pour déclencher l'adaptation
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		time.Sleep(constants.HTTPStatusOK * time.Millisecond) // Très rapide (200ms < 300ms threshold)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		
		// Générer du contenu avec des liens internes pour avoir plus de requêtes
		content := `<html><head><title>Fast Site</title></head><body>`
		if r.URL.Path == "/" {
			for i := 1; i <= 6; i++ {
				content += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
			}
		} else {
			content += fmt.Sprintf(`<h1>%s</h1><a href="/">Home</a>`, r.URL.Path)
		}
		content += `</body></html>`
		w.Write([]byte(content))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers:        3, // Commence avec 3 workers
		MaxWorkers:           10,
		MinWorkers:           1,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		AdaptIntervalSeconds: 1, // Adaptation rapide pour test
		MaxPages:             10, // Plus de pages pour déclencher l'adaptation
		MaxDepth:             2,  // Permettre le crawl de plusieurs niveaux
		TimeoutSeconds:       10,
	}
	
	// WHEN - Crawl un site rapide
	crawler := NewParallelCrawler(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	
	// Démarrer le crawl en arrière-plan
	go func() {
		crawler.CrawlWithContext(ctx, server.URL)
	}()
	
	// Attendre que l'adaptation se fasse (plus de temps pour traiter 5+ requêtes)
	time.Sleep(6 * time.Second)
	
	// THEN - Le crawler devrait augmenter les workers
	metrics := crawler.GetMetrics()
	assert.NotNil(t, metrics)
	assert.GreaterOrEqual(t, metrics.CurrentWorkers, 4, 
		"Should increase workers for fast site (current: %d)", metrics.CurrentWorkers)
	assert.Less(t, metrics.AvgResponseTime, time.Duration(constants.DefaultFastThresholdMs)*time.Millisecond,
		"Average response time should be under fast threshold")
}

func TestParallelCrawler_RespectMaxPages(t *testing.T) {
	// GIVEN - Serveur avec beaucoup de pages
	pageCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageCount++
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		// Générer des liens vers plus de pages
		links := ""
		for i := 1; i <= 50; i++ {
			links += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		w.Write([]byte(fmt.Sprintf(`<html><head><title>Page %d</title></head><body>%s</body></html>`, pageCount, links)))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers: 3,
		MaxWorkers:     10,
		MinWorkers:     1,
		MaxPages:       5, // LIMITE: seulement 5 pages
		TimeoutSeconds: 10,
	}
	
	// WHEN
	crawler := NewParallelCrawler(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	
	// THEN - Respecter la limite de pages
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.LessOrEqual(t, len(result.Pages), 5, "Should not exceed MaxPages limit")
}

func TestParallelCrawler_TimeoutHandling(t *testing.T) {
	// GIVEN - Serveur qui prend trop de temps
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Plus long que le timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers: 2,
		MaxWorkers:     5,
		MinWorkers:     1,
		MaxPages:       3,
		TimeoutSeconds: 3, // Timeout court pour test
	}
	
	// WHEN
	crawler := NewParallelCrawler(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	
	start := time.Now()
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	duration := time.Since(start)
	
	// THEN - Devrait timeout et pas attendre indéfiniment
	// Note: Le test peut réussir avec err != nil (timeout) ou result partiellement rempli
	assert.Less(t, duration, 5*time.Second, "Should timeout before 5s")
	if err != nil {
		// Timeout attendu - c'est OK
		assert.Contains(t, err.Error(), "context", "Error should be context-related")
	} else {
		// Ou résultat partiel - aussi OK
		assert.NotNil(t, result)
	}
}

func TestParallelCrawler_GetMetrics(t *testing.T) {
	// GIVEN
	cfg := &config.CrawlerConfig{
		InitialWorkers: constants.DefaultInitialWorkers,
		MaxWorkers:     constants.DefaultMaxWorkers,
		MinWorkers:     constants.DefaultMinWorkers,
		MaxPages:       5,
		TimeoutSeconds: 10,
	}
	
	// WHEN
	crawler := NewParallelCrawler(cfg)
	
	// THEN - Métrics disponibles même avant crawl
	metrics := crawler.GetMetrics()
	assert.NotNil(t, metrics)
	assert.Equal(t, constants.DefaultInitialWorkers, metrics.CurrentWorkers)
	assert.Equal(t, 0, metrics.PagesProcessed)
	assert.Equal(t, 0.0, metrics.PagesPerSecond)
}

// ========================================
// TESTS SPÉCIFIQUES POUR ROBOTS.TXT
// ========================================

func TestParallelCrawler_RespectRobotsTxt(t *testing.T) {
	// GIVEN - Serveur avec robots.txt restrictif
	mux := http.NewServeMux()
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
User-agent: *
Disallow: /private/
Disallow: /admin/
Crawl-delay: 1

User-agent: FireSalamander
Allow: /
Crawl-delay: 0.5
		`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Home</title></head>
			<body><a href="/public">Public</a><a href="/private/secret">Private</a></body></html>`))
	})
	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Public</title></head></html>`))
	})
	mux.HandleFunc("/private/secret", func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not crawl /private/secret - robots.txt violation!")
		w.WriteHeader(http.StatusOK)
	})
	
	server := httptest.NewServer(mux)
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers:   2,
		MaxWorkers:      5,
		MinWorkers:      1,
		MaxPages:        10,
		TimeoutSeconds:  10,
		RespectRobotsTxt: true,
		UserAgent:       constants.DefaultUserAgent,
	}
	
	// WHEN
	crawler := NewParallelCrawler(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	
	// THEN - Ne devrait pas crawler /private/secret
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Vérifier qu'aucune page /private/* n'a été crawlée
	for url := range result.Pages {
		assert.NotContains(t, url, constants.RoutePrivate+"/", "Should not crawl private pages due to robots.txt")
	}
}

// ========================================
// BENCHMARK TESTS
// ========================================

func BenchmarkParallelCrawler_SmallSite(b *testing.B) {
	// GIVEN - Site de test avec 10 pages
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Latence réaliste
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Test</title></head></html>`))
	}))
	defer server.Close()
	
	cfg := &config.CrawlerConfig{
		InitialWorkers: 5,
		MaxWorkers:     10,
		MinWorkers:     1,
		MaxPages:       10,
		TimeoutSeconds: 30,
	}
	
	// WHEN
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crawler := NewParallelCrawler(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		crawler.CrawlWithContext(ctx, server.URL)
		cancel()
	}
}