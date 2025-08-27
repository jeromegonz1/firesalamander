package crawler

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/logger"
	"golang.org/x/net/html"
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
	URL             string
	StatusCode      int
	ContentType     string
	Title           string
	Description     string
	MetaDescription string
	Headers         map[string]string
	Body            string
	Links           []Link
	Images          []Image
	Headings        []string
	CrawledAt       time.Time
	ResponseTime    time.Duration
	Error           error
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
		Timeout:        constants.ClientTimeout,
		RetryAttempts:  3,
		RetryDelay:     constants.DefaultRetryDelay,
		CacheDuration:  constants.RobotsCacheDuration,
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
		robotsCache: NewRobotsCache(constants.RobotsCacheDuration),
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
		log.Error(constants.CrawlerFailed + " page", map[string]interface{}{
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
func (c *Crawler) crawlWorker(ctx context.Context, wg *sync.WaitGroup, queue *CrawlQueue, results chan<- *CrawlResult, baseURL *url.URL, robots *LegacyRobotsTxt) {
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

// parseHTML parse le contenu HTML d'une page et extrait les données SEO et les liens
func (c *Crawler) parseHTML(result *CrawlResult) {
	if result.Body == "" {
		log.Debug("No content to parse", map[string]interface{}{"url": result.URL})
		return
	}

	log.Debug("Parsing HTML", map[string]interface{}{"url": result.URL, "content_size": len(result.Body)})

	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(result.Body))
	if err != nil {
		log.Error("Failed to parse HTML", map[string]interface{}{
			"url": result.URL,
			"error": err.Error(),
		})
		return
	}

	// Initialize parsing results
	var title string
	var metaDescription string
	var links []Link
	var images []Image
	var headings []string

	// Recursive function to traverse DOM
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if title == "" && n.FirstChild != nil {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				name := getAttr(n, "name")
				property := getAttr(n, "property")
				if (name == "description" || property == "og:description") && metaDescription == "" {
					metaDescription = strings.TrimSpace(getAttr(n, "content"))
				}
			case "a":
				if href := getAttr(n, "href"); href != "" {
					if resolvedURL := c.resolveURL(result.URL, href); resolvedURL != "" {
						linkText := strings.TrimSpace(extractText(n))
						linkType := "internal"
						if c.isExternalLink(result.URL, resolvedURL) {
							linkType = "external"
						}
						links = append(links, Link{
							URL:    resolvedURL,
							Text:   linkText,
							Type:   linkType,
							Follow: !strings.Contains(getAttr(n, "rel"), "nofollow"),
						})
					}
				}
			case "img":
				if src := getAttr(n, "src"); src != "" {
					if resolvedURL := c.resolveURL(result.URL, src); resolvedURL != "" {
						images = append(images, Image{
							URL:   resolvedURL,
							Alt:   getAttr(n, "alt"),
							Title: getAttr(n, "title"),
						})
					}
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				if text := extractText(n); text != "" {
					headings = append(headings, text)
				}
			}
		}

		// Traverse children
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	// Start traversal
	traverse(doc)

	// Update result with parsed data
	result.Title = title
	result.MetaDescription = metaDescription
	result.Links = links
	result.Images = images
	result.Headings = headings

	log.Info("HTML parsing completed", map[string]interface{}{
		"url": result.URL,
		"title_length": len(title),
		"meta_length": len(metaDescription),
		"links_found": len(links),
		"images_found": len(images),
		"headings_found": len(headings),
	})

	// Add discovered links to crawl queue for further crawling
	c.addLinksToQueue(links, result.URL)
}

// getAttr gets attribute value from HTML node
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// extractText extracts text content from HTML node
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var text strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		text.WriteString(extractText(child))
	}
	return strings.TrimSpace(text.String())
}

// resolveURL resolves relative URLs to absolute URLs
func (c *Crawler) resolveURL(base, href string) string {
	// Skip invalid URLs
	if href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
		return ""
	}

	// Parse base URL
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	// Parse href
	hrefURL, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Resolve relative to base
	resolved := baseURL.ResolveReference(hrefURL)
	
	// Only return HTTP/HTTPS URLs
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return ""
	}

	return resolved.String()
}

// addLinksToQueue adds discovered links to the crawling queue
func (c *Crawler) addLinksToQueue(links []Link, sourceURL string) {
	for _, link := range links {
		// Only add internal links that should be crawled
		if link.Type == "internal" && c.shouldCrawlURL(link.URL) {
			log.Debug("Adding link to queue", map[string]interface{}{
				"source": sourceURL,
				"target": link.URL,
			})
			// This method should be called by the parallel crawler to add URLs to its queue
			// Implementation will depend on how the parallel crawler manages its queue
		}
	}
}

// shouldCrawlURL determines if a URL should be crawled
func (c *Crawler) shouldCrawlURL(url string) bool {
	// Basic validation
	if url == "" {
		return false
	}

	// Skip non-HTTP URLs
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Skip common non-crawlable file types
	for _, ext := range []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".zip", ".rar", ".tar", ".gz"} {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return false
		}
	}

	// TODO: Add more sophisticated filtering based on robots.txt, patterns, etc.
	return true
}

// isExternalLink determines if a link is external to the current domain
func (c *Crawler) isExternalLink(baseURL, targetURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return true
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		return true
	}

	return base.Host != target.Host
}

// fetchRobotsTxt récupère et parse le fichier robots.txt
func (c *Crawler) fetchRobotsTxt(ctx context.Context, baseURL *url.URL) (*LegacyRobotsTxt, error) {
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
	RobotsTxt *LegacyRobotsTxt
	Sitemaps  []*Sitemap
	Stats     *CrawlStats
}