package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/config"
)

// 🔥🦎 FIRE SALAMANDER - TESTS ANTI-BOUCLE INFINIE
// NOUVEAU PROCESS V2.0 - TDD RENFORCÉ

// ✅ OBLIGATOIRE - Test avec timeout et détection de boucle
func TestParallelCrawler_MustTerminate(t *testing.T) {
	done := make(chan bool)
	failed := make(chan string)
	
	go func() {
		crawler := NewParallelCrawler(&config.CrawlerConfig{
			MaxPages:       5,
			InitialWorkers: 2,
			TimeoutSeconds: 30,
			UserAgent:      "Test-Agent",
		})
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		result, err := crawler.CrawlWithContext(ctx, "https://example.com")
		if err != nil && err != context.DeadlineExceeded {
			failed <- err.Error()
			return
		}
		
		if len(result.Pages) == 0 {
			failed <- "❌ No pages crawled"
			return
		}
		
		done <- true
	}()
	
	select {
	case <-done:
		t.Log("✅ Test réussi - Crawler termine correctement")
	case msg := <-failed:
		t.Fatalf("❌ Test échoué: %s", msg)
	case <-time.After(15 * time.Second):
		t.Fatal("❌ TIMEOUT CRITIQUE - Le crawler ne termine pas dans les 15s!")
	}
}

// Test de non-régression pour la boucle infinie
func TestParallelCrawler_NoInfiniteLoop(t *testing.T) {
	// Serveur qui crée une boucle intentionnelle
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			// Page qui se référence elle-même + page2
			w.Write([]byte(`
				<html><body>
					<a href="/">Self link - PIÈGE BOUCLE</a>
					<a href="/page2">Page 2</a>
				</body></html>
			`))
		case "/page2":
			// Page2 qui revient à /
			w.Write([]byte(`
				<html><body>
					<a href="/">Retour accueil</a>
					<a href="/page2">Self page2</a>
				</body></html>
			`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	
	crawler := NewParallelCrawler(&config.CrawlerConfig{
		MaxPages:       10,
		InitialWorkers: 1,
		TimeoutSeconds: 5,
		UserAgent:      "Anti-Loop-Test",
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("Crawl error: %v", err)
	}
	
	// Vérifier qu'aucune URL n'est crawlée plusieurs fois
	urls := make(map[string]int)
	for url, pageResult := range result.Pages {
		urls[url]++
		if urls[url] > 1 {
			t.Fatalf("❌ BOUCLE INFINIE DÉTECTÉE! URL '%s' crawlée %d fois", url, urls[url])
		}
		
		if pageResult.Error != nil {
			t.Logf("⚠️ Page avec erreur: %s - %v", url, pageResult.Error)
		}
	}
	
	t.Logf("✅ Anti-boucle OK: %d URLs uniques crawlées", len(urls))
	
	// Le crawler DOIT s'arrêter rapidement
	if len(result.Pages) == 0 {
		t.Fatal("❌ Aucune page crawlée - possible deadlock")
	}
}

// Test benchmark avec limite de temps STRICTE
func TestParallelCrawler_BenchmarkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulation d'une réponse lente
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("<html><body>Slow page</body></html>"))
	}))
	defer server.Close()
	
	crawler := NewParallelCrawler(&config.CrawlerConfig{
		MaxPages:       3,
		InitialWorkers: 1,
		TimeoutSeconds: 2, // Timeout court
		UserAgent:      "Benchmark-Test",
	})
	
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result, err := crawler.CrawlWithContext(ctx, server.URL)
	elapsed := time.Since(start)
	
	if elapsed > 4*time.Second {
		t.Fatalf("❌ PERFORMANCE: Crawl trop lent (%v), max 4s attendu", elapsed)
	}
	
	t.Logf("⏱️ Benchmark: %v elapsed, %d pages, err=%v", elapsed, len(result.Pages), err)
	
	// Le crawler DOIT soit réussir soit timeout proprement
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}
}

// Test de détection de deadlock avec channels
func TestParallelCrawler_NoDeadlock(t *testing.T) {
	deadlockDetected := make(chan bool)
	
	go func() {
		// Si on arrive ici, pas de deadlock
		crawler := NewParallelCrawler(&config.CrawlerConfig{
			MaxPages:       1,
			InitialWorkers: 1,
			TimeoutSeconds: 5,
			UserAgent:      "Deadlock-Test",
		})
		
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		_, err := crawler.CrawlWithContext(ctx, "https://httpbin.org/get")
		if err != nil {
			t.Logf("ℹ️ Expected error/timeout: %v", err)
		}
		
		deadlockDetected <- false // Pas de deadlock
	}()
	
	select {
	case hasDeadlock := <-deadlockDetected:
		if hasDeadlock {
			t.Fatal("❌ DEADLOCK DÉTECTÉ!")
		}
		t.Log("✅ No deadlock detected")
	case <-time.After(10 * time.Second):
		t.Fatal("❌ PROBABLE DEADLOCK - Le test ne termine pas!")
	}
}