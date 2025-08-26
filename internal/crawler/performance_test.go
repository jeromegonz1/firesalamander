package crawler

import (
	"context"
	"runtime"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD RED PHASE - TESTS DE PERFORMANCE
// Tests anti-régression : sites moyens doivent finir < 60s
// ========================================

// TestPerformance_Under60Seconds teste que le crawler finit en moins de 60 secondes
// Anti-régression : éviter les timeouts infinis comme avec septeo.com
func TestPerformance_Under60Seconds(t *testing.T) {
	// ARRANGE : Configuration optimisée pour performance
	cfg := &config.CrawlerConfig{
		MaxPages:             30,  // Site moyen
		MaxDepth:             3,   // Profondeur raisonnable
		TimeoutSeconds:       60,  // Contrainte forte : 60 secondes max
		InitialWorkers:       8,   // Plus de workers pour la performance
		MaxWorkers:           15,
		MinWorkers:           3,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false, // Désactivé pour performance max
		DelayMs:              10,    // Délai minimal pour performance
		FastThresholdMs:      200,   // Seuil rapide adapté
		SlowThresholdMs:      1000,  // Seuil lent adapté
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: 2,     // Adaptation plus fréquente
	}

	// ACT : Test de performance avec mesure précise
	performanceCrawler := NewPerformanceCrawler(cfg)
	if performanceCrawler == nil {
		t.Fatal("NewPerformanceCrawler devrait retourner un crawler de performance, got nil")
	}

	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Second) // Marge pour cleanup
	defer cancel()

	startTime := time.Now()
	result, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)
	duration := time.Since(startTime)

	// ASSERT : Contrainte de performance CRITIQUE
	if duration > 60*time.Second {
		t.Errorf("ÉCHEC PERFORMANCE : Crawl trop lent %v > 60s (anti-régression septeo.com)", duration)
	}

	if err != nil {
		t.Errorf("CrawlWithPerformanceOptimization ne devrait pas échouer: %v", err)
	}

	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	// Vérifier qu'on a quand même des résultats malgré l'optimisation
	if len(result.Pages) == 0 {
		t.Error("Le crawler optimisé devrait trouver au moins quelques pages")
	}

	// Vérifier les métriques de performance
	if result.PerformanceReport == nil {
		t.Error("Le crawler de performance devrait fournir un rapport")
	}

	t.Logf("✅ Performance OK: %d pages crawlées en %v (cible: <60s)", len(result.Pages), duration)
}

// TestPerformance_MemoryEfficiency teste l'efficacité mémoire
func TestPerformance_MemoryEfficiency(t *testing.T) {
	// ARRANGE : Configuration pour test mémoire
	cfg := &config.CrawlerConfig{
		MaxPages:             50,  // Plus de pages pour tester la mémoire
		MaxDepth:             4,
		TimeoutSeconds:       90,
		InitialWorkers:       6,
		MaxWorkers:           10,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false,
		DelayMs:              20,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	performanceCrawler := NewPerformanceCrawler(cfg)

	// Mesurer la mémoire avant le crawl
	runtime.GC() // Force garbage collection
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// ACT
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)

	// Mesurer la mémoire après le crawl
	runtime.GC() // Force garbage collection
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// ASSERT : Vérifications mémoire
	if err != nil {
		t.Errorf("Le test de mémoire ne devrait pas échouer: %v", err)
	}

	memoryUsed := memAfter.Alloc - memBefore.Alloc
	memoryUsedMB := float64(memoryUsed) / 1024 / 1024

	// Seuil raisonnable : maximum 100MB pour 50 pages
	if memoryUsedMB > 100 {
		t.Errorf("Consommation mémoire excessive: %.2f MB (limite: 100 MB)", memoryUsedMB)
	}

	// Vérifier les métriques de mémoire dans le rapport
	if result != nil && result.PerformanceReport != nil {
		if result.PerformanceReport.MemoryUsage == 0 {
			t.Error("Le rapport devrait inclure l'usage mémoire")
		}
	}

	t.Logf("✅ Mémoire OK: %.2f MB utilisés pour %d pages", memoryUsedMB, len(result.Pages))
}

// TestPerformance_ConcurrencyEfficiency teste l'efficacité de la concurrence
func TestPerformance_ConcurrencyEfficiency(t *testing.T) {
	// ARRANGE : Différentes configurations de concurrence
	testConfigs := []struct {
		name           string
		workers        int
		expectedPagesPerSecond float64
	}{
		{"LowConcurrency", 2, 0.5},
		{"MediumConcurrency", 5, 1.0},
		{"HighConcurrency", 10, 1.5},
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.CrawlerConfig{
				MaxPages:             20,
				MaxDepth:             2,
				TimeoutSeconds:       60,
				InitialWorkers:       tc.workers,
				MaxWorkers:           tc.workers + 2,
				MinWorkers:           1,
				UserAgent:            constants.ParallelCrawlerUserAgent,
				RespectRobotsTxt:     false,
				DelayMs:              5, // Délai très court pour tester concurrence
				FastThresholdMs:      constants.DefaultFastThresholdMs,
				SlowThresholdMs:      constants.DefaultSlowThresholdMs,
				ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
				AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
			}

			performanceCrawler := NewPerformanceCrawler(cfg)

			// ACT
			testURL := constants.CrawlerTestURLExample
			ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()

			startTime := time.Now()
			result, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)
			duration := time.Since(startTime)

			// ASSERT
			if err != nil {
				t.Errorf("Test concurrence %s échoué: %v", tc.name, err)
				return
			}

			if len(result.Pages) > 0 {
				pagesPerSecond := float64(len(result.Pages)) / duration.Seconds()
				
				// Pas de seuil strict car dépend du réseau, mais log pour analyse
				t.Logf("%s: %.2f pages/sec (%d pages en %v)",
					tc.name, pagesPerSecond, len(result.Pages), duration)
				
				// Vérifier que plus de workers = plus d'efficacité (tendance générale)
				if result.PerformanceReport != nil {
					avgWorkers := result.PerformanceReport.AverageWorkersUsed
					if avgWorkers == 0 {
						t.Error("Le rapport devrait indiquer le nombre moyen de workers")
					}
				}
			}
		})
	}
}

// TestPerformance_AdaptiveWorkerPool teste l'adaptation automatique des workers
func TestPerformance_AdaptiveWorkerPool(t *testing.T) {
	// ARRANGE : Configuration avec adaptation aggressive
	cfg := &config.CrawlerConfig{
		MaxPages:             40,
		MaxDepth:             3,
		TimeoutSeconds:       120,
		InitialWorkers:       3,
		MaxWorkers:           12, // Large plage pour adaptation
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false,
		DelayMs:              30,
		FastThresholdMs:      300, // Seuils pour déclencher adaptation
		SlowThresholdMs:      1500,
		ErrorThresholdPercent: 15,
		AdaptIntervalSeconds: 3, // Adaptation fréquente
	}

	performanceCrawler := NewPerformanceCrawler(cfg)

	// ACT
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)

	// ASSERT : Vérifier l'adaptation
	if err != nil {
		t.Errorf("Test adaptation workers échoué: %v", err)
	}

	if result == nil {
		t.Fatal("Le résultat ne devrait pas être nil")
	}

	// Vérifier que l'adaptation a eu lieu
	perfReport := result.PerformanceReport
	if perfReport == nil {
		t.Fatal("Le rapport de performance ne devrait pas être nil")
	}

	if perfReport.WorkerAdaptations == 0 {
		t.Log("Aucune adaptation de workers détectée - peut être normal selon les conditions")
	}

	if perfReport.MaxWorkersReached == 0 {
		t.Log("Maximum de workers jamais atteint - peut indiquer un sous-dimensionnement")
	}

	// Vérifier la cohérence des métriques
	if perfReport.AverageWorkersUsed > float64(cfg.MaxWorkers) {
		t.Errorf("Nombre moyen de workers incohérent: %.1f > max %d",
			perfReport.AverageWorkersUsed, cfg.MaxWorkers)
	}
}

// TestPerformance_NoBottlenecks teste l'absence de goulots d'étranglement
func TestPerformance_NoBottlenecks(t *testing.T) {
	// ARRANGE : Configuration pour détecter les bottlenecks
	cfg := &config.CrawlerConfig{
		MaxPages:             25,
		MaxDepth:             3,
		TimeoutSeconds:       90,
		InitialWorkers:       8,
		MaxWorkers:           12,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false,
		DelayMs:              1, // Délai minimal pour stresser le système
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	performanceCrawler := NewPerformanceCrawler(cfg)

	// ACT
	testURL := constants.CrawlerTestURLExample
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)

	// ASSERT : Détecter les bottlenecks
	if err != nil {
		t.Errorf("Test bottlenecks échoué: %v", err)
	}

	if result != nil && result.PerformanceReport != nil {
		report := result.PerformanceReport

		// Vérifier l'efficacité de la queue
		if report.QueueUtilization < 0.1 {
			t.Log("Queue sous-utilisée - peut indiquer un bottleneck dans la production de tâches")
		}

		if report.QueueUtilization > 0.9 {
			t.Log("Queue surchargée - peut indiquer un bottleneck dans le traitement")
		}

		// Vérifier l'équilibrage des workers
		if report.WorkerIdleTime > 50 {
			t.Logf("Workers inactifs %d%% du temps - peut indiquer un déséquilibre", 
				report.WorkerIdleTime)
		}

		// Vérifier les temps d'attente
		if report.AverageWaitTime > 1*time.Second {
			t.Errorf("Temps d'attente moyen trop élevé: %v", report.AverageWaitTime)
		}
	}
}

// BenchmarkCrawler_Performance benchmark de performance pour régression
func BenchmarkCrawler_Performance(b *testing.B) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             10, // Taille fixe pour benchmark reproductible
		MaxDepth:             2,
		TimeoutSeconds:       30,
		InitialWorkers:       4,
		MaxWorkers:           6,
		MinWorkers:           2,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		RespectRobotsTxt:     false,
		DelayMs:              5,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	testURL := constants.CrawlerTestURLExample

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		performanceCrawler := NewPerformanceCrawler(cfg)
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		
		_, err := performanceCrawler.CrawlWithPerformanceOptimization(ctx, testURL)
		if err != nil {
			b.Errorf("Benchmark échoué à l'itération %d: %v", i, err)
		}
		
		cancel()
	}
}

// ========================================
// INTERFACES POUR PERFORMANCE CRAWLER (À IMPLÉMENTER)
// ========================================

// IPerformanceCrawler interface pour crawler optimisé performance
type IPerformanceCrawler interface {
	CrawlWithPerformanceOptimization(ctx context.Context, startURL string) (*PerformanceCrawlResult, error)
	SetPerformanceTargets(targets *PerformanceTargets)
	GetPerformanceMetrics() *RealTimePerformanceMetrics
}

// ========================================
// TYPES POUR PERFORMANCE CRAWLER (À IMPLÉMENTER)
// ========================================

// PerformanceCrawler implémentation optimisée pour performance
type PerformanceCrawler struct {
	// À implémenter selon l'architecture
}

// PerformanceCrawlResult résultat avec métriques de performance détaillées
type PerformanceCrawlResult struct {
	// Résultats basiques
	StartURL  string                  `json:"start_url"`
	Pages     map[string]*PageResult  `json:"pages"`
	Duration  time.Duration          `json:"duration"`
	
	// Rapport de performance détaillé
	PerformanceReport *PerformanceReport `json:"performance_report"`
}

// PerformanceReport rapport détaillé de performance
type PerformanceReport struct {
	// Métriques workers
	AverageWorkersUsed  float64 `json:"average_workers_used"`
	MaxWorkersReached   int     `json:"max_workers_reached"`
	WorkerAdaptations   int     `json:"worker_adaptations"`
	WorkerIdleTime      int     `json:"worker_idle_time_percent"`
	
	// Métriques queue
	QueueUtilization    float64       `json:"queue_utilization"`
	AverageWaitTime     time.Duration `json:"average_wait_time"`
	
	// Métriques système
	MemoryUsage         uint64        `json:"memory_usage_bytes"`
	CPUUsage            float64       `json:"cpu_usage_percent"`
	NetworkUtilization  float64       `json:"network_utilization"`
	
	// Métriques temporelles
	AverageResponseTime time.Duration `json:"average_response_time"`
	P95ResponseTime     time.Duration `json:"p95_response_time"`
	P99ResponseTime     time.Duration `json:"p99_response_time"`
}

// PerformanceTargets cibles de performance à respecter
type PerformanceTargets struct {
	MaxDuration         time.Duration `json:"max_duration"`
	MinPagesPerSecond   float64       `json:"min_pages_per_second"`
	MaxMemoryMB         int           `json:"max_memory_mb"`
	MaxCPUPercent       float64       `json:"max_cpu_percent"`
}

// RealTimePerformanceMetrics métriques en temps réel
type RealTimePerformanceMetrics struct {
	CurrentPagesPerSecond float64       `json:"current_pages_per_second"`
	CurrentWorkers        int           `json:"current_workers"`
	CurrentMemoryMB       float64       `json:"current_memory_mb"`
	EstimatedCompletion   time.Duration `json:"estimated_completion"`
}

// ========================================
// FONCTIONS À IMPLÉMENTER (STUBS POUR TESTS ROUGES)
// ========================================

// NewPerformanceCrawler crée un crawler optimisé pour la performance
func NewPerformanceCrawler(cfg *config.CrawlerConfig) IPerformanceCrawler {
	// TODO: Implémenter selon l'architecture reçue
	// Cette fonction doit retourner nil pour que les tests échouent (RED phase)
	return nil
}