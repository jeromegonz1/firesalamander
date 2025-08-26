package crawler

import (
	"context"
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

// ICrawlerEngine définit l'interface du moteur de crawl intelligent
type ICrawlerEngine interface {
	CrawlWithIntelligence(ctx context.Context, startURL string) (*IntelligentCrawlResult, error)
	GetTerminationController() ITerminationController
	GetURLDiscoveryService() IURLDiscoveryService
	GetJobCounter() IAtomicJobCounter
}

// ITerminationController gère la terminaison propre du crawler
type ITerminationController interface {
	GetTerminationConditions() []ITerminationCondition
	ShouldTerminate() bool
	NotifyJobCompleted()
	NotifyJobStarted()
}

// ITerminationCondition représente une condition de terminaison
type ITerminationCondition interface {
	Type() string
	IsMet() bool
}

// IURLDiscoveryService découvre des URLs via différentes méthodes
type IURLDiscoveryService interface {
	DiscoverFromSitemap(ctx context.Context, baseURL string) ([]string, error)
	DiscoverFromRobots(ctx context.Context, baseURL string) ([]string, error)
	DiscoverFromHTML(html string, baseURL string) ([]string, error)
}

// IAtomicJobCounter gère le comptage atomique des jobs
type IAtomicJobCounter interface {
	Get() int32
	Add(delta int32) int32
	Sub(delta int32) int32
	Reset()
}

// ========================================
// TYPES POUR LES TESTS (À IMPLÉMENTER)
// ========================================

// IntelligentCrawler est le crawler intelligent à implémenter
type IntelligentCrawler struct {
	// À implémenter selon l'architecture
}

// IntelligentCrawlResult résultat du crawl intelligent
type IntelligentCrawlResult struct {
	StartURL  string                  `json:"start_url"`
	Pages     map[string]*PageResult  `json:"pages"`
	Duration  time.Duration          `json:"duration"`
	Metrics   *CrawlerMetrics        `json:"metrics"`
	Error     error                  `json:"error,omitempty"`
	
	// Nouvelles métriques intelligentes
	TerminationReason string    `json:"termination_reason"`
	JobsExecuted      int32     `json:"jobs_executed"`
	RaceConditionsDetected int `json:"race_conditions_detected"`
}

// ========================================
// FONCTIONS À IMPLÉMENTER (STUBS POUR TESTS)
// ========================================

// NewIntelligentCrawler crée un nouveau crawler intelligent (À IMPLÉMENTER)
func NewIntelligentCrawler(cfg *config.CrawlerConfig) ICrawlerEngine {
	// TODO: Implémenter selon l'architecture reçue
	// Cette fonction doit retourner nil pour que les tests échouent (RED phase)
	return nil
}