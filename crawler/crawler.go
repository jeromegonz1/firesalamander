package crawler

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/jeromegonz1/firesalamander/internal/logger"
)

var log = logger.New("CRAWLER")

// Crawler représente le crawler principal
type Crawler struct {
	config      *Config
	fetcher     *Fetcher
	robotsCache *RobotsCache
	sitemapParser *SitemapParser
	cache       *PageCache
	rateLimiter *RateLimiter
	stats       *CrawlStats
	mu          sync.RWMutex
}

// Config contient la configuration du crawler
type Config struct {
	UserAgent       string
	Workers         int
	RateLimit       string // Format: "10/s" ou "60/m"
	MaxDepth        int
	MaxPages        int
	Timeout         time.Duration
	RetryAttempts   int
	RetryDelay      time.Duration
	CacheDuration   time.Duration
	RespectRobots   bool
	FollowSitemaps  bool
	EnableCache     bool
}

// CrawlResult représente le résultat d'un crawl
type CrawlResult struct {
	URL          string
	StatusCode   int
	ContentType  string
	Title        string
	Description  string
	Headers      map[string]string
	Body         string
	Links        []Link
	Images       []Image
	CrawledAt    time.Time
	ResponseTime time.Duration
	Error        error
}

// Link représente un lien trouvé
type Link struct {
	URL    string
	Text   string
	Type   string // internal, external
	Follow bool
}

// Image représente une image trouvée
type Image struct {
	URL    string
	Alt    string
	Title  string
	Width  int
	Height int
}

// CrawlStats contient les statistiques du crawl
type CrawlStats struct {
	mu               sync.RWMutex
	StartTime        time.Time
	EndTime          time.Time
	TotalPages       int
	SuccessfulPages  int
	FailedPages      int
	TotalBytes       int64
	AverageRespTime  time.Duration
	RobotsBlocked    int
	CacheHits        int
	CacheMisses      int
}

// DefaultConfig retourne une configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		UserAgent:      "FireSalamander/1.0 (SEPTEO) SEO Analyzer",
		Workers:        2,
		RateLimit:      "10/s",
		MaxDepth:       3,
		MaxPages:       100,
		Timeout:        30 * time.Second,
		RetryAttempts:  3,
		RetryDelay:     1 * time.Second,
		CacheDuration:  7 * 24 * time.Hour,
		RespectRobots:  true,
		FollowSitemaps: true,
		EnableCache:    true,
	}
}

// New crée une nouvelle instance du crawler
func New(config *Config) (*Crawler, error) {
	log.Debug("Creating new crawler", map[string]interface{}{
		"workers":     config.Workers,
		"rate_limit":  config.RateLimit,
		"user_agent":  config.UserAgent,
	})

	if config == nil {
		config = DefaultConfig()
	}

	// Créer le rate limiter
	rateLimiter, err := NewRateLimiter(config.RateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter: %w", err)
	}

	crawler := &Crawler{
		config:      config,
		fetcher:     NewFetcher(config),
		robotsCache: NewRobotsCache(7 * 24 * time.Hour), // Cache 7 jours
		sitemapParser: NewSitemapParser(),
		cache:       NewPageCache(config.CacheDuration),
		rateLimiter: rateLimiter,
		stats:       &CrawlStats{StartTime: time.Now()},
	}

	log.Info("Crawler created successfully", map[string]interface{}{
		"config": config,
	})

	return crawler, nil
}

// CrawlSite effectue le crawl d'un site entier
func (c *Crawler) CrawlSite(ctx context.Context, startURL string) (*CrawlReport, error) {
	log.Info("Starting site crawl", map[string]interface{}{
		"url":         startURL,
		"max_pages":   c.config.MaxPages,
		"max_depth":   c.config.MaxDepth,
	})

	// Valider l'URL
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Initialiser le rapport
	report := &CrawlReport{
		StartURL:  startURL,
		Domain:    parsedURL.Host,
		StartTime: time.Now(),
		Pages:     make(map[string]*CrawlResult),
		Stats:     c.stats,
	}

	// Vérifier robots.txt si activé
	if c.config.RespectRobots {
		log.Debug("Checking robots.txt", map[string]interface{}{"domain": parsedURL.Host})
		robots, err := c.fetchRobotsTxt(ctx, parsedURL)
		if err != nil {
			log.Warn("Failed to fetch robots.txt", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			report.RobotsTxt = robots
		}
	}

	// Découvrir les sitemaps si activé
	if c.config.FollowSitemaps {
		log.Debug("Discovering sitemaps")
		sitemaps, err := c.discoverSitemaps(ctx, parsedURL)
		if err != nil {
			log.Warn("Failed to discover sitemaps", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			report.Sitemaps = sitemaps
		}
	}

	// Créer la queue de crawl
	queue := NewCrawlQueue(c.config.MaxPages)
	queue.Add(startURL, 0)

	// Ajouter les URLs des sitemaps
	for _, sitemap := range report.Sitemaps {
		for _, url := range sitemap.URLs {
			queue.Add(url.Loc, 0)
		}
	}

	// Démarrer les workers
	var wg sync.WaitGroup
	results := make(chan *CrawlResult, c.config.Workers*2)
	
	// Worker pool
	for i := 0; i < c.config.Workers; i++ {
		wg.Add(1)
		go c.crawlWorker(ctx, &wg, queue, results, parsedURL, report.RobotsTxt)
	}

	// Collecteur de résultats
	go func() {
		wg.Wait()
		close(results)
	}()

	// Traiter les résultats
	for result := range results {
		c.mu.Lock()
		report.Pages[result.URL] = result
		c.updateStats(result)
		c.mu.Unlock()

		// Ajouter les nouveaux liens à la queue
		if result.Error == nil {
			for _, link := range result.Links {
				if link.Type == "internal" && link.Follow {
					queue.Add(link.URL, 1) // TODO: calculer la profondeur réelle
				}
			}
		}
	}

	// Finaliser le rapport
	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)
	c.stats.EndTime = report.EndTime

	log.Info("Site crawl completed", map[string]interface{}{
		"duration":     report.Duration,
		"pages_crawled": len(report.Pages),
		"successful":   c.stats.SuccessfulPages,
		"failed":       c.stats.FailedPages,
	})

	return report, nil
}

// CrawlPage effectue le crawl d'une page unique
func (c *Crawler) CrawlPage(ctx context.Context, pageURL string) (*CrawlResult, error) {
	log.Debug("Crawling page", map[string]interface{}{"url": pageURL})

	// Vérifier le cache
	if c.config.EnableCache {
		if cached, found := c.cache.Get(pageURL); found {
			log.Debug("Cache hit", map[string]interface{}{"url": pageURL})
			c.stats.mu.Lock()
			c.stats.CacheHits++
			c.stats.mu.Unlock()
			return cached, nil
		}
		c.stats.mu.Lock()
		c.stats.CacheMisses++
		c.stats.mu.Unlock()
	}

	// Rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Fetch la page
	result, err := c.fetcher.Fetch(ctx, pageURL)
	if err != nil {
		log.Error("Failed to fetch page", map[string]interface{}{
			"url":   pageURL,
			"error": err.Error(),
		})
		return result, err
	}

	// Parser le contenu HTML
	if result.ContentType == "text/html" {
		c.parseHTML(result)
	}

	// Mettre en cache si activé
	if c.config.EnableCache && result.Error == nil {
		c.cache.Set(pageURL, result)
	}

	return result, nil
}

// crawlWorker est un worker qui traite les URLs de la queue
func (c *Crawler) crawlWorker(ctx context.Context, wg *sync.WaitGroup, queue *CrawlQueue, results chan<- *CrawlResult, baseURL *url.URL, robots *RobotsTxt) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, ok := queue.Next()
			if !ok {
				return
			}

			// Vérifier robots.txt
			if c.config.RespectRobots && robots != nil {
				if !robots.IsAllowed(c.config.UserAgent, item.URL) {
					log.Debug("Blocked by robots.txt", map[string]interface{}{"url": item.URL})
					c.stats.mu.Lock()
					c.stats.RobotsBlocked++
					c.stats.mu.Unlock()
					results <- &CrawlResult{
						URL:   item.URL,
						Error: fmt.Errorf("blocked by robots.txt"),
					}
					continue
				}
			}

			// Crawler la page
			result, err := c.CrawlPage(ctx, item.URL)
			if err != nil {
				result = &CrawlResult{
					URL:   item.URL,
					Error: err,
				}
			}

			results <- result
		}
	}
}

// parseHTML parse le contenu HTML d'une page
func (c *Crawler) parseHTML(result *CrawlResult) {
	// TODO: Implémenter le parsing HTML
	// - Extraire le title
	// - Extraire la meta description
	// - Extraire les liens
	// - Extraire les images
	log.Debug("Parsing HTML", map[string]interface{}{"url": result.URL})
}

// fetchRobotsTxt récupère et parse le fichier robots.txt
func (c *Crawler) fetchRobotsTxt(ctx context.Context, baseURL *url.URL) (*RobotsTxt, error) {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", baseURL.Scheme, baseURL.Host)
	
	// Vérifier le cache
	if robots, found := c.robotsCache.Get(baseURL.Host); found {
		return robots, nil
	}

	// Fetcher robots.txt
	result, err := c.fetcher.Fetch(ctx, robotsURL)
	if err != nil {
		return nil, err
	}

	// Parser robots.txt
	robots, err := ParseRobotsTxt(result.Body)
	if err != nil {
		return nil, err
	}

	// Mettre en cache
	c.robotsCache.Set(baseURL.Host, robots)
	
	return robots, nil
}

// discoverSitemaps découvre les sitemaps d'un site
func (c *Crawler) discoverSitemaps(ctx context.Context, baseURL *url.URL) ([]*Sitemap, error) {
	var sitemaps []*Sitemap

	// Vérifier sitemap.xml standard
	sitemapURL := fmt.Sprintf("%s://%s/sitemap.xml", baseURL.Scheme, baseURL.Host)
	sitemap, err := c.fetchSitemap(ctx, sitemapURL)
	if err == nil {
		sitemaps = append(sitemaps, sitemap)
	}

	// TODO: Vérifier robots.txt pour d'autres sitemaps

	return sitemaps, nil
}

// fetchSitemap récupère et parse un sitemap
func (c *Crawler) fetchSitemap(ctx context.Context, sitemapURL string) (*Sitemap, error) {
	result, err := c.fetcher.Fetch(ctx, sitemapURL)
	if err != nil {
		return nil, err
	}

	return c.sitemapParser.Parse(result.Body)
}

// updateStats met à jour les statistiques
func (c *Crawler) updateStats(result *CrawlResult) {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()

	c.stats.TotalPages++
	if result.Error == nil {
		c.stats.SuccessfulPages++
		c.stats.TotalBytes += int64(len(result.Body))
		
		// Calculer la moyenne du temps de réponse
		if c.stats.AverageRespTime == 0 {
			c.stats.AverageRespTime = result.ResponseTime
		} else {
			c.stats.AverageRespTime = (c.stats.AverageRespTime + result.ResponseTime) / 2
		}
	} else {
		c.stats.FailedPages++
	}
}

// GetStats retourne les statistiques du crawl
func (c *Crawler) GetStats() CrawlStats {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()
	return *c.stats
}

// CrawlReport représente le rapport complet d'un crawl
type CrawlReport struct {
	StartURL  string
	Domain    string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Pages     map[string]*CrawlResult
	RobotsTxt *RobotsTxt
	Sitemaps  []*Sitemap
	Stats     *CrawlStats
}