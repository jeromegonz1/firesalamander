package crawler

import (
	"context"
	"strings"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// 🔴 TDD PHASE RED - TESTS QUI DOIVENT ÉCHOUER
// ========================================

// TestIntelligentAdapter_ImplementsParallelCrawlerInterface teste l'interface
func TestIntelligentAdapter_ImplementsParallelCrawlerInterface(t *testing.T) {
	// ARRANGE : Configuration standard
	cfg := &config.CrawlerConfig{
		MaxPages:       constants.DefaultMaxPages,
		MaxDepth:       constants.DefaultMaxDepth,
		TimeoutSeconds: constants.DefaultTimeoutSeconds,
		InitialWorkers: constants.DefaultInitialWorkers,
		UserAgent:      constants.ParallelCrawlerUserAgent,
	}

	// ACT : Créer l'adaptateur (DOIT ÉCHOUER - n'existe pas encore)
	adapter := NewIntelligentAdapter(cfg)

	// ASSERT : L'adaptateur doit implémenter l'interface ParallelCrawler
	if adapter == nil {
		t.Fatal("NewIntelligentAdapter devrait retourner un adaptateur, got nil")
	}

	// Test que l'interface CrawlWithContext existe
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := adapter.CrawlWithContext(ctx, constants.CrawlerTestURLExample)
	if err != nil && !strings.Contains(err.Error(), "context") {
		t.Errorf("CrawlWithContext devrait fonctionner ou timeout, got: %v", err)
	}
	if result == nil {
		t.Error("CrawlWithContext devrait retourner un résultat non-nil")
	}
}

// TestIntelligentAdapter_UsesCleanHTML teste que cleanHTML est appliqué
func TestIntelligentAdapter_UsesCleanHTML(t *testing.T) {
	// ARRANGE : Configuration de test
	cfg := &config.CrawlerConfig{
		MaxPages:       1,
		MaxDepth:       1, 
		TimeoutSeconds: 10,
		InitialWorkers: 1,
		UserAgent:      constants.ParallelCrawlerUserAgent,
	}

	// ACT : Crawl une page (DOIT ÉCHOUER - adaptateur n'existe pas)
	adapter := NewIntelligentAdapter(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := adapter.CrawlWithContext(ctx, constants.CrawlerTestURLExample)

	// ASSERT : Vérifier que cleanHTML a été appliqué
	if err != nil {
		t.Fatalf("CrawlWithContext ne devrait pas échouer: %v", err)
	}
	if len(result.Pages) == 0 {
		t.Fatal("Aucune page trouvée")
	}

	// Vérifier qu'il n'y a pas de caractères de contrôle invalides
	for url, page := range result.Pages {
		t.Logf("Testing page: %s", url)
		
		// Test critique : vérifier cleanHTML
		for _, r := range page.Body {
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				t.Errorf("Page %s contient un caractère de contrôle invalide: %d", url, r)
			}
		}
		
		// Test que le titre est nettoyé
		for _, r := range page.Title {
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				t.Errorf("Title %s contient un caractère de contrôle invalide: %d", page.Title, r)
			}
		}
		break // Une page suffit pour le test
	}
}

// TestIntelligentAdapter_ReturnsParallelCrawlResult teste le type de retour
func TestIntelligentAdapter_ReturnsParallelCrawlResult(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:       1,
		TimeoutSeconds: 10,
		InitialWorkers: 1,
		UserAgent:      constants.ParallelCrawlerUserAgent,
	}

	// ACT : (DOIT ÉCHOUER - n'existe pas encore)
	adapter := NewIntelligentAdapter(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := adapter.CrawlWithContext(ctx, constants.CrawlerTestURLExample)

	// ASSERT : Type de retour correct
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}
	
	// Vérifier que c'est bien un ParallelCrawlResult
	if result.StartURL == "" {
		t.Error("StartURL ne devrait pas être vide")
	}
	if result.Pages == nil {
		t.Error("Pages ne devrait pas être nil")
	}
	if result.Duration == 0 {
		t.Error("Duration devrait être > 0")
	}
}

// TestIntelligentAdapter_PerformanceBetter teste les performances améliorées
func TestIntelligentAdapter_PerformanceBetter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// ARRANGE : Configuration pour site complexe
	cfg := &config.CrawlerConfig{
		MaxPages:       3,
		MaxDepth:       2,
		TimeoutSeconds: 30,
		InitialWorkers: 2,
		UserAgent:      constants.ParallelCrawlerUserAgent,
	}

	// ACT : Test de performance (DOIT ÉCHOUER - n'existe pas)
	adapter := NewIntelligentAdapter(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	start := time.Now()
	result, err := adapter.CrawlWithContext(ctx, "https://example.com") // Site stable pour tests
	duration := time.Since(start)

	// ASSERT : Performance acceptable (pas de timeout 90s)
	if err != nil {
		t.Fatalf("Crawl ne devrait pas échouer: %v", err)
	}
	if duration > 25*time.Second {
		t.Errorf("Crawl trop lent: %v (max 25s attendu)", duration)
	}
	if len(result.Pages) == 0 {
		t.Error("Aucune page trouvée")
	}

	t.Logf("✅ Performance test: %v pour %d pages", duration, len(result.Pages))
}