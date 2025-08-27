package crawler

import (
	"context"

	"firesalamander/internal/config"
)

// IntelligentAdapter adapte le crawler intelligent pour l'interface ParallelCrawler existante
// 
// Cette structure permet l'intégration transparente du crawler intelligent avec cleanHTML()
// dans le code existant sans breaking changes, respectant le principe Open/Closed de SOLID.
//
// Avantages:
// - Compatible 100% avec interface ParallelCrawler
// - Injection automatique de cleanHTML() pour résoudre caractères HTML invalides  
// - Performance améliorée vs timeouts 90s
// - Facilité de rollback si nécessaire
type IntelligentAdapter struct {
	intelligentCrawler ICrawlerEngine // Crawler intelligent avec cleanHTML intégré
}

// NewIntelligentAdapter crée un nouvel adaptateur utilisant le crawler intelligent avec cleanHTML
func NewIntelligentAdapter(cfg *config.CrawlerConfig) *IntelligentAdapter {
	return &IntelligentAdapter{
		intelligentCrawler: NewIntelligentCrawler(cfg),
	}
}

// CrawlWithContext adapte CrawlWithIntelligence vers l'interface ParallelCrawler existante
//
// Cette méthode maintient la signature compatible avec l'orchestrator existant
// tout en utilisant le crawler intelligent avec cleanHTML() automatiquement appliqué.
//
// Résout le problème des caractères HTML invalides qui causaient des timeouts 90s
// sur des sites comme septeo.com, resalys.com, etc.
func (ia *IntelligentAdapter) CrawlWithContext(ctx context.Context, url string) (*ParallelCrawlResult, error) {
	// Phase 1: Utiliser le crawler intelligent avec cleanHTML intégré
	log.Info("🔍 IntelligentAdapter: Starting CrawlWithIntelligence", map[string]interface{}{"url": url})
	intelligentResult, err := ia.intelligentCrawler.CrawlWithIntelligence(ctx, url)
	if err != nil {
		return &ParallelCrawlResult{
			StartURL: url,
			Error:    err,
		}, err
	}

	// Phase 2: Adapter le résultat vers ParallelCrawlResult pour compatibilité totale
	// Les pages sont déjà nettoyées par cleanHTML() dans CrawlWithIntelligence()
	log.Info("🔍 IntelligentAdapter: Got pages from CrawlWithIntelligence", map[string]interface{}{"pages_count": len(intelligentResult.Pages)})
	adaptedResult := &ParallelCrawlResult{
		StartURL: intelligentResult.StartURL,
		Pages:    intelligentResult.Pages, // ✅ Pages avec cleanHTML appliqué
		Duration: intelligentResult.Duration,
		Metrics:  intelligentResult.Metrics,
		Error:    intelligentResult.Error,
	}

	return adaptedResult, nil
}