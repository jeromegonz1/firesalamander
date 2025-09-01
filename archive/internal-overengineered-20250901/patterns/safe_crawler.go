package patterns

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// 🔥🦎 FIRE SALAMANDER - PATTERN OBLIGATOIRE SafeCrawler
// NOUVEAU PROCESS V2.0 - SAFETY FIRST

// SafeCrawler - PATTERN OBLIGATOIRE pour tout crawler
type SafeCrawler struct {
	visitedURLs sync.Map // Thread-safe map pour éviter les doublons
	maxPages    int
	timeout     time.Duration
	maxRetries  int
	urlCounter  sync.Map // Compteur d'accès par URL pour détecter les boucles
}

// CrawlResult - Résultat sécurisé d'un crawl
type CrawlResult struct {
	URL          string
	Success      bool
	Error        error
	Attempts     int
	Duration     time.Duration
	LoopDetected bool
}

// SafetyMetrics - Métriques de sécurité
type SafetyMetrics struct {
	TotalURLs        int
	UniqueURLs       int
	LoopsDetected    int
	TimeoutsOccurred int
	PagesProcessed   int
	Errors           int
}

// NewSafeCrawler - Constructeur avec valeurs par défaut sécurisées
func NewSafeCrawler() *SafeCrawler {
	return &SafeCrawler{
		maxPages:   20,          // Limite raisonnable
		timeout:    30 * time.Second, // Timeout global
		maxRetries: 3,           // Éviter les boucles de retry
	}
}

// WithSafetyLimits - Configuration sécurisée
func (c *SafeCrawler) WithSafetyLimits(maxPages int, timeout time.Duration) *SafeCrawler {
	// Valeurs maximum autorisées pour éviter les abus
	if maxPages > 100 {
		log.Printf("⚠️ SAFETY: maxPages limité à 100 (demandé: %d)", maxPages)
		maxPages = 100
	}
	if timeout > 5*time.Minute {
		log.Printf("⚠️ SAFETY: timeout limité à 5min (demandé: %v)", timeout)
		timeout = 5 * time.Minute
	}
	
	c.maxPages = maxPages
	c.timeout = timeout
	return c
}

// CrawlWithSafety - PATTERN: Circuit breaker obligatoire
func (c *SafeCrawler) CrawlWithSafety(url string) <-chan CrawlResult {
	results := make(chan CrawlResult, c.maxPages)
	
	go func() {
		defer close(results)
		
		// SAFETY 1: Timeout global OBLIGATOIRE
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
		
		// SAFETY 2: Métriques de monitoring
		metrics := &SafetyMetrics{}
		defer func() {
			log.Printf("🛡️ SAFETY REPORT: %+v", metrics)
		}()
		
		// SAFETY 3: Compteur de pages
		pageCount := 0
		
		// SAFETY 4: Détection de boucle renforcée
		urlAccessCount := make(map[string]int)
		
		// SAFETY 5: Canal d'arrêt d'urgence
		emergency := make(chan bool)
		go c.emergencyBreaker(ctx, emergency, metrics)
		
		// File d'URLs à traiter (simulation)
		urlQueue := []string{url}
		
		for len(urlQueue) > 0 && pageCount < c.maxPages {
			select {
			case <-ctx.Done():
				log.Printf("⏱️ SAFETY: Timeout global atteint (%v)", c.timeout)
				metrics.TimeoutsOccurred++
				results <- CrawlResult{
					URL:     url,
					Success: false,
					Error:   ctx.Err(),
				}
				return
				
			case <-emergency:
				log.Println("🚨 SAFETY: Arrêt d'urgence activé")
				results <- CrawlResult{
					URL:     url,
					Success: false,
					Error:   fmt.Errorf("emergency stop activated"),
				}
				return
				
			default:
				// Traiter la prochaine URL
				currentURL := urlQueue[0]
				urlQueue = urlQueue[1:]
				
				// ANTI-BOUCLE CRITIQUE
				urlAccessCount[currentURL]++
				if urlAccessCount[currentURL] > 2 {
					log.Printf("🔄 SAFETY: Boucle détectée pour %s (accès #%d)", 
						currentURL, urlAccessCount[currentURL])
					metrics.LoopsDetected++
					
					results <- CrawlResult{
						URL:          currentURL,
						Success:      false,
						Error:        fmt.Errorf("infinite loop detected"),
						LoopDetected: true,
						Attempts:     urlAccessCount[currentURL],
					}
					continue
				}
				
				// VÉRIFICATION DÉJÀ VU
				if _, seen := c.visitedURLs.LoadOrStore(currentURL, true); seen {
					log.Printf("🔍 SAFETY: URL déjà visitée ignorée: %s", currentURL)
					continue
				}
				
				// SIMULATION DE CRAWL (remplacer par vraie logique)
				start := time.Now()
				success, err := c.simulateCrawl(ctx, currentURL)
				duration := time.Since(start)
				
				result := CrawlResult{
					URL:      currentURL,
					Success:  success,
					Error:    err,
					Duration: duration,
					Attempts: urlAccessCount[currentURL],
				}
				
				if success {
					metrics.PagesProcessed++
				} else {
					metrics.Errors++
				}
				
				results <- result
				pageCount++
				
				// Mise à jour métriques
				metrics.TotalURLs++
				metrics.UniqueURLs = len(urlAccessCount)
				
				log.Printf("📊 SAFETY Progress: %d/%d pages, %d loops, %d errors", 
					pageCount, c.maxPages, metrics.LoopsDetected, metrics.Errors)
			}
		}
		
		if pageCount >= c.maxPages {
			log.Printf("📊 SAFETY: Limite max pages atteinte (%d)", c.maxPages)
		}
	}()
	
	return results
}

// emergencyBreaker - Circuit breaker d'urgence
func (c *SafeCrawler) emergencyBreaker(ctx context.Context, emergency chan bool, metrics *SafetyMetrics) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// CONDITIONS D'ARRÊT D'URGENCE
			if metrics.LoopsDetected > 5 {
				log.Printf("🚨 EMERGENCY: Trop de boucles détectées (%d)", metrics.LoopsDetected)
				emergency <- true
				return
			}
			
			if metrics.Errors > 10 {
				log.Printf("🚨 EMERGENCY: Trop d'erreurs (%d)", metrics.Errors)
				emergency <- true
				return
			}
			
			// Détection de comportement anormal (trop lent)
			if metrics.TotalURLs > 0 && float64(metrics.PagesProcessed)/float64(metrics.TotalURLs) < 0.1 {
				log.Printf("🚨 EMERGENCY: Taux de succès trop bas (%d/%d)", 
					metrics.PagesProcessed, metrics.TotalURLs)
				// emergency <- true  // Désactivé pour ne pas être trop strict
			}
		}
	}
}

// simulateCrawl - Simulation de crawl (remplacer par la vraie logique)
func (c *SafeCrawler) simulateCrawl(ctx context.Context, url string) (bool, error) {
	// Simulation de temps de traitement variable
	processingTime := 50 * time.Millisecond
	if url == "https://slow.example.com" {
		processingTime = 2 * time.Second
	}
	
	select {
	case <-time.After(processingTime):
		// Simulation de succès/échec
		if url == "https://error.example.com" {
			return false, fmt.Errorf("simulated error for %s", url)
		}
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// GetMetrics - Exposition des métriques de sécurité
func (c *SafeCrawler) GetMetrics() map[string]interface{} {
	visitedCount := 0
	c.visitedURLs.Range(func(key, value interface{}) bool {
		visitedCount++
		return true
	})
	
	return map[string]interface{}{
		"max_pages_limit":    c.maxPages,
		"timeout_seconds":    int(c.timeout.Seconds()),
		"urls_visited":       visitedCount,
		"safety_pattern":     "SafeCrawler v2.0",
		"circuit_breaker":    "ACTIVE",
		"anti_loop":          "ENABLED",
		"emergency_stop":     "ENABLED",
	}
}

// Reset - Remise à zéro sécurisée
func (c *SafeCrawler) Reset() {
	c.visitedURLs = sync.Map{}
	c.urlCounter = sync.Map{}
	log.Println("🔄 SAFETY: SafeCrawler reset completed")
}