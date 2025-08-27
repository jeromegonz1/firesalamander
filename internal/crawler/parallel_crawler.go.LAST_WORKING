package crawler

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	
	"golang.org/x/net/html"
	
	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD PHASE 2: IMPLÉMENTATION MINIMALE (GREEN)
// Uniquement le code nécessaire pour faire passer les tests
// ========================================

// ParallelCrawler implements adaptive parallel web crawling
type ParallelCrawler struct {
	config      *config.CrawlerConfig
	urlQueue    chan CrawlTask
	results     chan *PageResult  
	errors      chan error
	metrics     *CrawlerMetrics
	workerPool  *WorkerPool
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	robotsCache map[string]*ParallelRobotsTxt
	client      *http.Client
	
	// 🔥 FIRE SALAMANDER - CONDITION D'ARRÊT
	activeJobs    int32           // Compteur atomique de jobs actifs
	urlsSeen      map[string]bool // URLs déjà vues
	pagesCrawled  int32           // Pages crawlées
	done          chan struct{}   // Channel pour signaler la fin
	doneOnce      sync.Once       // Assurer une seule fermeture de done
	urlSeenMu     sync.RWMutex    // Mutex pour urlsSeen
}

// CrawlerMetrics tracks real-time performance metrics
type CrawlerMetrics struct {
	mu              sync.RWMutex
	PagesPerSecond  float64       `json:"pages_per_second"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	ErrorRate       float64       `json:"error_rate"`
	CurrentWorkers  int          `json:"current_workers"`
	PagesProcessed  int          `json:"pages_processed"`
	TotalErrors     int          `json:"total_errors"`
	StartTime       time.Time    `json:"start_time"`
	LastUpdate      time.Time    `json:"last_update"`
}

// WorkerPool manages the adaptive worker pool
type WorkerPool struct {
	mu          sync.RWMutex
	workers     int
	maxWorkers  int
	minWorkers  int
	lastAdapt   time.Time
	responseTimes []time.Duration
	errorCount    int
	totalRequests int
}

// CrawlTask represents a URL to crawl with metadata
type CrawlTask struct {
	URL     string
	Depth   int
	Referer string
}

// PageResult represents the result of crawling a single page
type PageResult struct {
	URL          string            `json:"url"`
	StatusCode   int              `json:"status_code"`
	ContentType  string           `json:"content_type"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Headers      map[string]string `json:"headers"`
	Body         string           `json:"body"`
	Links        []ParallelLink   `json:"links"`
	ResponseTime time.Duration    `json:"response_time"`
	Error        error            `json:"error,omitempty"`
	CrawledAt    time.Time        `json:"crawled_at"`
	Depth        int              `json:"depth"`
}

// ParallelCrawlResult represents the complete crawling result
type ParallelCrawlResult struct {
	StartURL  string                  `json:"start_url"`
	Pages     map[string]*PageResult  `json:"pages"`
	Duration  time.Duration          `json:"duration"`
	Metrics   *CrawlerMetrics        `json:"metrics"`
	Error     error                  `json:"error,omitempty"`
}

// ParallelLink represents a discovered link
type ParallelLink struct {
	URL    string `json:"url"`
	Text   string `json:"text"`
	Type   string `json:"type"` // "internal", "external"
}

// ParallelRobotsTxt represents parsed robots.txt rules
type ParallelRobotsTxt struct {
	UserAgentRules map[string][]string
	CrawlDelays    map[string]time.Duration
	Sitemaps       []string
}

// NewParallelCrawler creates a new adaptive parallel crawler
func NewParallelCrawler(cfg *config.CrawlerConfig) *ParallelCrawler {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Apply defaults from constants if not set
	if cfg.InitialWorkers == 0 {
		cfg.InitialWorkers = constants.DefaultInitialWorkers
	}
	if cfg.MaxWorkers == 0 {
		cfg.MaxWorkers = constants.DefaultMaxWorkers
	}
	if cfg.MinWorkers == 0 {
		cfg.MinWorkers = constants.DefaultMinWorkers
	}
	if cfg.FastThresholdMs == 0 {
		cfg.FastThresholdMs = constants.DefaultFastThresholdMs
	}
	if cfg.SlowThresholdMs == 0 {
		cfg.SlowThresholdMs = constants.DefaultSlowThresholdMs
	}
	if cfg.ErrorThresholdPercent == 0 {
		cfg.ErrorThresholdPercent = constants.DefaultErrorThresholdPercent
	}
	if cfg.AdaptIntervalSeconds == 0 {
		cfg.AdaptIntervalSeconds = constants.DefaultAdaptIntervalSeconds
	}
	if cfg.MaxPages == 0 {
		cfg.MaxPages = constants.DefaultMaxPages
	}
	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = constants.DefaultTimeoutSeconds
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = constants.ParallelCrawlerUserAgent
	}

	// HTTP client with timeout from config (NO HARDCODING)
	client := &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Limit redirects to prevent infinite loops
			if len(via) >= constants.DefaultMaxRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	
	return &ParallelCrawler{
		config:      cfg,
		urlQueue:    make(chan CrawlTask, cfg.MaxPages*2),
		results:     make(chan *PageResult, cfg.MaxPages),
		errors:      make(chan error, cfg.MaxPages),
		metrics: &CrawlerMetrics{
			CurrentWorkers: cfg.InitialWorkers,
			StartTime:      time.Now(),
			LastUpdate:     time.Now(),
		},
		workerPool: &WorkerPool{
			workers:    cfg.InitialWorkers,
			maxWorkers: cfg.MaxWorkers,
			minWorkers: cfg.MinWorkers,
			lastAdapt:  time.Now(),
			responseTimes: make([]time.Duration, 0, constants.DefaultMetricsHistorySize),
		},
		ctx:         ctx,
		cancel:      cancel,
		robotsCache: make(map[string]*ParallelRobotsTxt),
		client:      client,
	}
}

// CrawlWithContext crawls a website starting from the given URL with context
func (pc *ParallelCrawler) CrawlWithContext(ctx context.Context, startURL string) (*ParallelCrawlResult, error) {
	// Apply timeout from config (NO HARDCODING)
	crawlCtx, cancel := context.WithTimeout(ctx, time.Duration(pc.config.TimeoutSeconds)*time.Second)
	defer cancel()
	
	result := &ParallelCrawlResult{
		StartURL: startURL,
		Pages:    make(map[string]*PageResult),
		Metrics:  pc.metrics,
	}
	
	start := time.Now()
	defer func() {
		result.Duration = time.Since(start)
	}()
	
	// Parse start URL
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return result, fmt.Errorf("invalid start URL: %w", err)
	}
	
	// Check robots.txt if configured
	if pc.config.RespectRobotsTxt {
		robotsTxt, err := pc.fetchRobotsTxt(crawlCtx, parsedURL.Scheme+"://"+parsedURL.Host)
		if err == nil {
			pc.robotsCache[parsedURL.Host] = robotsTxt
		}
	}
	
	// 🔥 FIRE SALAMANDER - INITIALISATION CONDITION D'ARRÊT
	atomic.StoreInt32(&pc.activeJobs, 0)
	atomic.StoreInt32(&pc.pagesCrawled, 0)
	pc.urlSeenMu.Lock()
	pc.urlsSeen = make(map[string]bool)
	pc.urlSeenMu.Unlock()
	pc.done = make(chan struct{})
	
	// Start workers
	stdlog.Printf("🔍 DEBUG CRAWLER: Starting %d workers for URL: %s", pc.config.InitialWorkers, startURL)
	var wg sync.WaitGroup
	for i := 0; i < pc.config.InitialWorkers; i++ {
		wg.Add(1)
		go pc.workerWithCounter(crawlCtx, &wg, parsedURL)
	}
	
	// Start adaptation routine
	go pc.adaptWorkers(crawlCtx, &wg, parsedURL)
	
	// 🔥 FIRE SALAMANDER - MONITEUR DE COMPLETION
	go func() {
		ticker := time.NewTicker(time.Duration(constants.DefaultMonitoringIntervalMs) * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				jobs := atomic.LoadInt32(&pc.activeJobs)
				crawled := atomic.LoadInt32(&pc.pagesCrawled)
				stdlog.Printf("📊 Active jobs: %d, Pages crawled: %d/%d", 
					jobs, crawled, pc.config.MaxPages)
				
				// CONDITIONS D'ARRÊT
				if jobs == 0 || crawled >= int32(pc.config.MaxPages) {
					stdlog.Println("✅ Crawl complete - closing channels")
					pc.doneOnce.Do(func() { close(pc.done) })
					return
				}
				
			case <-crawlCtx.Done():
				stdlog.Println("⚠️ Context cancelled - stopping crawler")
				pc.doneOnce.Do(func() { close(pc.done) })
				return
			}
		}
	}()
	
	// Add initial URL to queue
	stdlog.Printf("🔍 DEBUG CRAWLER: Adding initial URL to queue: %s", startURL)
	atomic.AddInt32(&pc.activeJobs, 1) // Premier job actif
	pc.urlSeenMu.Lock()
	pc.urlsSeen[startURL] = true
	pc.urlSeenMu.Unlock()
	select {
	case pc.urlQueue <- CrawlTask{URL: startURL, Depth: 0}:
		stdlog.Printf("🔍 DEBUG CRAWLER: Initial URL added successfully")
	case <-crawlCtx.Done():
		pc.cancel()
		return result, crawlCtx.Err()
	}
	
	// Collect results
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Ignore panic from closing nil channels
			}
		}()
		wg.Wait()
		if pc.results != nil {
			close(pc.results)
		}
		if pc.errors != nil {
			close(pc.errors)
		}
	}()
	
	// 🔥 FIRE SALAMANDER - BOUCLE DE TRAITEMENT AVEC CONDITION D'ARRÊT
	stdlog.Printf("📊 Starting result processing loop with termination logic")
	for {
		select {
		case pageResult, ok := <-pc.results:
			if !ok {
				pc.results = nil
				continue
			}
			if pageResult != nil {
				stdlog.Printf("📊 Received result for: %s (Total: %d pages)", pageResult.URL, len(result.Pages)+1)
				result.Pages[pageResult.URL] = pageResult
				pc.updateMetrics(pageResult)
				
				// NOTE: Les nouveaux liens sont ajoutés par workerWithCounter
			}
			
		case err, ok := <-pc.errors:
			if !ok {
				pc.errors = nil
				continue
			}
			if err != nil {
				pc.updateErrorMetrics()
			}
			
		case <-pc.done:
			// 🔥 FIRE SALAMANDER - TERMINAISON PROPRE
			stdlog.Printf("✅ Crawl completed naturally - %d pages found", len(result.Pages))
			if pc.cancel != nil {
				pc.cancel()
			}
			return result, nil
			
		case <-crawlCtx.Done():
			stdlog.Printf("⚠️ Context timeout - %d pages found", len(result.Pages))
			pc.cancel()
			if result.Error == nil {
				result.Error = crawlCtx.Err()
			}
			return result, nil
		}
		
		// Exit conditions (backup)
		if pc.results == nil && pc.errors == nil {
			stdlog.Printf("ℹ️ All channels closed - breaking loop")
			break
		}
		if len(result.Pages) >= pc.config.MaxPages {
			stdlog.Printf("✅ Max pages reached: %d/%d", len(result.Pages), pc.config.MaxPages)
			pc.cancel()
			break
		}
	}
	
	return result, result.Error
}

// 🔥 FIRE SALAMANDER - WORKER AVEC COMPTEUR DE JOBS
func (pc *ParallelCrawler) workerWithCounter(ctx context.Context, wg *sync.WaitGroup, baseURL *url.URL) {
	defer wg.Done()
	stdlog.Printf("🕷️ Worker started with job counter")
	
	for {
		select {
		case task, ok := <-pc.urlQueue:
			if !ok {
				stdlog.Printf("🛑 Worker: Queue closed, exiting")
				return
			}
			
			stdlog.Printf("🕷️ Worker got job: %s (depth: %d)", task.URL, task.Depth)
			
			// Check robots.txt compliance
			if pc.config.RespectRobotsTxt {
				if !pc.isAllowedByRobots(task.URL, baseURL.Host) {
					// Décrémenter car on ne traite pas ce job
					atomic.AddInt32(&pc.activeJobs, -1)
					stdlog.Printf("🚫 Blocked by robots.txt: %s", task.URL)
					continue
				}
			}
			
			// Crawl the page
			result := pc.crawlPage(ctx, task)
			
			// Incrémenter le compteur de pages crawlées
			crawled := atomic.AddInt32(&pc.pagesCrawled, 1)
			stdlog.Printf("✅ Page crawled (%d): %s", crawled, result.URL)
			
			// Send result
			select {
			case pc.results <- result:
				stdlog.Printf("📤 Result sent for: %s", result.URL)
			case <-ctx.Done():
				return
			case <-pc.done:
				return
			}
			
			// Ajouter nouvelles URLs découvertes
			if result.Error == nil && result.Depth < pc.config.MaxDepth && crawled < int32(pc.config.MaxPages) {
				pc.addNewLinksWithCounter(result.Links, task.URL, task.Depth+1, baseURL)
			}
			
			// Décrémenter le job actif (IMPORTANT!)
			remainingJobs := atomic.AddInt32(&pc.activeJobs, -1)
			stdlog.Printf("📊 Job completed. Remaining active jobs: %d", remainingJobs)
			
			// Vérifier la condition d'arrêt immédiatement (évite les race conditions)
			if remainingJobs == 0 {
				stdlog.Println("🏁 No more active jobs - signaling completion")
				pc.doneOnce.Do(func() { close(pc.done) })
				return
			}
			
		case <-pc.done:
			stdlog.Printf("🛑 Worker stopping - crawl complete")
			return
			
		case <-ctx.Done():
			stdlog.Printf("🛑 Worker stopping - context cancelled")
			return
		}
	}
}

// worker is a goroutine that processes URLs from the queue (OLD VERSION)
func (pc *ParallelCrawler) worker(ctx context.Context, wg *sync.WaitGroup, baseURL *url.URL) {
	defer wg.Done()
	stdlog.Printf("🔍 DEBUG WORKER: Worker started")
	
	for {
		select {
		case task, ok := <-pc.urlQueue:
			if !ok {
				stdlog.Printf("🔍 DEBUG WORKER: Queue closed, exiting")
				return
			}
			stdlog.Printf("🔍 DEBUG WORKER: Got task: %s", task.URL)
			
			// Check robots.txt compliance
			if pc.config.RespectRobotsTxt {
				if !pc.isAllowedByRobots(task.URL, baseURL.Host) {
					continue
				}
			}
			
			// Crawl the page
			stdlog.Printf("🔍 DEBUG WORKER: Crawling page: %s", task.URL)
			result := pc.crawlPage(ctx, task)
			if result != nil {
				stdlog.Printf("🔍 DEBUG WORKER: Page crawled, sending result for: %s", result.URL)
				select {
				case pc.results <- result:
					stdlog.Printf("🔍 DEBUG WORKER: Result sent successfully")
				case <-ctx.Done():
					return
				}
			}
			
		case <-ctx.Done():
			return
		}
	}
}

// crawlPage crawls a single page and returns the result
func (pc *ParallelCrawler) crawlPage(ctx context.Context, task CrawlTask) *PageResult {
	start := time.Now()
	result := &PageResult{
		URL:       task.URL,
		Depth:     task.Depth,
		CrawledAt: start,
		Headers:   make(map[string]string),
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		result.Error = err
		return result
	}
	
	// Set user agent from config (NO HARDCODING)
	req.Header.Set("User-Agent", pc.config.UserAgent)
	if task.Referer != "" {
		req.Header.Set("Referer", task.Referer)
	}
	
	// Make request
	resp, err := pc.client.Do(req)
	if err != nil {
		result.Error = err
		result.ResponseTime = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	
	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.ResponseTime = time.Since(start)
	
	// Copy headers
	for name, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[name] = values[0]
		}
	}
	
	// Read body (limit size to prevent memory issues)
	bodyBuilder := &strings.Builder{}
	maxBodySize := int64(constants.DefaultMaxBodySize1MB) // 1MB limit from constants
	_, err = io.Copy(bodyBuilder, &io.LimitedReader{
		R: resp.Body,
		N: maxBodySize,
	})
	if err != nil {
		result.Error = err
		return result
	}
	result.Body = bodyBuilder.String()
	
	// Parse HTML if content is HTML
	if strings.Contains(result.ContentType, "text/html") {
		pc.parseHTML(result)
	}
	
	return result
}

// parseHTML extracts title, description and links from HTML content
func (pc *ParallelCrawler) parseHTML(result *PageResult) {
	doc, err := html.Parse(strings.NewReader(result.Body))
	if err != nil {
		return
	}
	
	pc.extractFromHTML(doc, result)
}

// extractFromHTML recursively extracts information from HTML nodes
func (pc *ParallelCrawler) extractFromHTML(n *html.Node, result *PageResult) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				result.Title = strings.TrimSpace(n.FirstChild.Data)
			}
		case "meta":
			name := pc.getAttr(n, "name")
			content := pc.getAttr(n, "content")
			if name == "description" && content != "" {
				result.Description = content
			}
		case "a":
			href := pc.getAttr(n, "href")
			if href != "" {
				text := pc.getTextContent(n)
				linkType := "internal"
				if strings.Contains(href, "://") && !strings.Contains(href, result.URL) {
					linkType = "external"
				}
				result.Links = append(result.Links, ParallelLink{
					URL:  href,
					Text: text,
					Type: linkType,
				})
			}
		}
	}
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		pc.extractFromHTML(c, result)
	}
}

// getAttr gets an attribute value from an HTML node
func (pc *ParallelCrawler) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getTextContent extracts text content from an HTML node
func (pc *ParallelCrawler) getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	
	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(pc.getTextContent(c))
	}
	return strings.TrimSpace(text.String())
}

// adaptWorkers dynamically adjusts the number of workers based on performance
func (pc *ParallelCrawler) adaptWorkers(ctx context.Context, wg *sync.WaitGroup, baseURL *url.URL) {
	ticker := time.NewTicker(time.Duration(pc.config.AdaptIntervalSeconds) * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			pc.performAdaptation(wg, baseURL)
		case <-ctx.Done():
			return
		}
	}
}

// performAdaptation analyzes metrics and adjusts worker count
func (pc *ParallelCrawler) performAdaptation(wg *sync.WaitGroup, baseURL *url.URL) {
	pc.workerPool.mu.Lock()
	defer pc.workerPool.mu.Unlock()
	
	if len(pc.workerPool.responseTimes) < constants.DefaultMinMetricsForAdaptation {
		return // Not enough data for adaptation
	}
	
	// Calculate average response time
	var totalTime time.Duration
	for _, rt := range pc.workerPool.responseTimes {
		totalTime += rt
	}
	avgTime := totalTime / time.Duration(len(pc.workerPool.responseTimes))
	
	// Calculate error rate
	errorRate := float64(pc.workerPool.errorCount) / float64(pc.workerPool.totalRequests) * 100
	
	oldWorkers := pc.workerPool.workers
	action := constants.WorkerMaintainAction
	reason := "metrics_stable"
	
	// Adaptation logic based on thresholds from config (NO HARDCODING)
	fastThreshold := time.Duration(pc.config.FastThresholdMs) * time.Millisecond
	slowThreshold := time.Duration(pc.config.SlowThresholdMs) * time.Millisecond
	
	if errorRate > float64(pc.config.ErrorThresholdPercent) {
		// High error rate - reduce workers
		if pc.workerPool.workers > pc.workerPool.minWorkers {
			pc.workerPool.workers--
			action = constants.WorkerDecreaseAction
			reason = fmt.Sprintf("high_error_rate_%.1f", errorRate)
		}
	} else if avgTime < fastThreshold {
		// Fast site - increase workers
		if pc.workerPool.workers < pc.workerPool.maxWorkers {
			pc.workerPool.workers++
			action = constants.WorkerIncreaseAction
			reason = fmt.Sprintf("fast_response_%.0fms", avgTime.Seconds()*1000)
			
			// Add new worker
			wg.Add(1)
			go pc.worker(pc.ctx, wg, baseURL)
		}
	} else if avgTime > slowThreshold {
		// Slow site - reduce workers
		if pc.workerPool.workers > pc.workerPool.minWorkers {
			pc.workerPool.workers--
			action = constants.WorkerDecreaseAction
			reason = fmt.Sprintf("slow_response_%.0fms", avgTime.Seconds()*1000)
		}
	}
	
	// Update metrics
	pc.metrics.mu.Lock()
	pc.metrics.CurrentWorkers = pc.workerPool.workers
	pc.metrics.AvgResponseTime = avgTime
	pc.metrics.ErrorRate = errorRate
	pc.metrics.LastUpdate = time.Now()
	pc.metrics.mu.Unlock()
	
	// Log adaptation (using message from constants)
	if action != constants.WorkerMaintainAction {
		fmt.Printf(constants.AdaptingWorkersMsg+"\n", action, oldWorkers, pc.workerPool.workers, reason)
	}
	
	// Debug info for tests (removed after testing)
	
	// Reset metrics for next interval
	pc.workerPool.responseTimes = pc.workerPool.responseTimes[:0]
	pc.workerPool.errorCount = 0
	pc.workerPool.totalRequests = 0
}

// updateMetrics updates crawler metrics with new page result
func (pc *ParallelCrawler) updateMetrics(result *PageResult) {
	pc.metrics.mu.Lock()
	defer pc.metrics.mu.Unlock()
	
	pc.metrics.PagesProcessed++
	
	// Update worker pool metrics
	pc.workerPool.mu.Lock()
	pc.workerPool.responseTimes = append(pc.workerPool.responseTimes, result.ResponseTime)
	pc.workerPool.totalRequests++
	if result.Error != nil {
		pc.workerPool.errorCount++
		pc.metrics.TotalErrors++
	}
	// Debug removed
	pc.workerPool.mu.Unlock()
	
	// Calculate pages per second
	elapsed := time.Since(pc.metrics.StartTime)
	if elapsed > 0 {
		pc.metrics.PagesPerSecond = float64(pc.metrics.PagesProcessed) / elapsed.Seconds()
	}
	
	pc.metrics.LastUpdate = time.Now()
}

// updateErrorMetrics updates error-related metrics
func (pc *ParallelCrawler) updateErrorMetrics() {
	pc.metrics.mu.Lock()
	defer pc.metrics.mu.Unlock()
	
	pc.metrics.TotalErrors++
	
	// Update error rate
	if pc.metrics.PagesProcessed > 0 {
		pc.metrics.ErrorRate = float64(pc.metrics.TotalErrors) / float64(pc.metrics.PagesProcessed) * 100
	}
}

// GetMetrics returns current crawler metrics (thread-safe)
func (pc *ParallelCrawler) GetMetrics() *CrawlerMetrics {
	pc.metrics.mu.RLock()
	defer pc.metrics.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	return &CrawlerMetrics{
		PagesPerSecond:  pc.metrics.PagesPerSecond,
		AvgResponseTime: pc.metrics.AvgResponseTime,
		ErrorRate:       pc.metrics.ErrorRate,
		CurrentWorkers:  pc.metrics.CurrentWorkers,
		PagesProcessed:  pc.metrics.PagesProcessed,
		TotalErrors:     pc.metrics.TotalErrors,
		StartTime:       pc.metrics.StartTime,
		LastUpdate:      pc.metrics.LastUpdate,
	}
}

// 🔥 FIRE SALAMANDER - AJOUTER LIENS AVEC COMPTEUR
func (pc *ParallelCrawler) addNewLinksWithCounter(links []ParallelLink, referer string, depth int, baseURL *url.URL) {
	for _, link := range links {
		// Parse et valider l'URL
		parsedLink, err := url.Parse(link.URL)
		if err != nil {
			continue
		}
		
		// Résoudre les URLs relatives
		if !parsedLink.IsAbs() {
			parsedLink = baseURL.ResolveReference(parsedLink)
		}
		
		// Vérifier que c'est sur le même domaine
		if parsedLink.Host != baseURL.Host {
			continue
		}
		
		normalizedURL := parsedLink.String()
		
		// Vérifier si déjà vue
		pc.urlSeenMu.Lock()
		if pc.urlsSeen[normalizedURL] {
			pc.urlSeenMu.Unlock()
			continue
		}
		pc.urlsSeen[normalizedURL] = true
		pc.urlSeenMu.Unlock()
		
		// Ajouter un job actif
		atomic.AddInt32(&pc.activeJobs, 1)
		
		// Ajouter à la queue
		task := CrawlTask{
			URL:     normalizedURL,
			Depth:   depth,
			Referer: referer,
		}
		
		select {
		case pc.urlQueue <- task:
			stdlog.Printf("➕ Added new URL to queue: %s (depth: %d)", normalizedURL, depth)
		default:
			// Queue pleine, annuler le job
			atomic.AddInt32(&pc.activeJobs, -1)
			stdlog.Printf("⚠️ Queue full, skipping: %s", normalizedURL)
		}
	}
}

// addLinksToQueue adds discovered links to the crawl queue (OLD VERSION)
func (pc *ParallelCrawler) addLinksToQueue(links []ParallelLink, referer string, depth int, baseURL *url.URL) {
	for _, link := range links {
		if link.Type == "internal" {
			// Convert relative URLs to absolute
			linkURL, err := url.Parse(link.URL)
			if err != nil {
				continue
			}
			
			if !linkURL.IsAbs() {
				linkURL = baseURL.ResolveReference(linkURL)
			}
			
			// Only add internal links from same domain
			if linkURL.Host == baseURL.Host {
				select {
				case pc.urlQueue <- CrawlTask{
					URL:     linkURL.String(),
					Depth:   depth,
					Referer: referer,
				}:
				default:
					// Queue full, skip this URL
				}
			}
		}
	}
}

// fetchRobotsTxt fetches and parses robots.txt for a domain
func (pc *ParallelCrawler) fetchRobotsTxt(ctx context.Context, baseURL string) (*ParallelRobotsTxt, error) {
	robotsURL := baseURL + "/robots.txt"
	
	req, err := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", pc.config.UserAgent)
	
	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("robots.txt not found: %d", resp.StatusCode)
	}
	
	bodyBuilder := &strings.Builder{}
	_, err = io.Copy(bodyBuilder, resp.Body)
	if err != nil {
		return nil, err
	}
	
	return pc.parseRobotsTxt(bodyBuilder.String()), nil
}

// parseRobotsTxt parses robots.txt content
func (pc *ParallelCrawler) parseRobotsTxt(content string) *ParallelRobotsTxt {
	robots := &ParallelRobotsTxt{
		UserAgentRules: make(map[string][]string),
		CrawlDelays:    make(map[string]time.Duration),
		Sitemaps:       make([]string, 0),
	}
	
	lines := strings.Split(content, "\n")
	currentAgent := "*"
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])
		
		switch key {
		case "user-agent":
			currentAgent = value
		case "disallow":
			robots.UserAgentRules[currentAgent] = append(robots.UserAgentRules[currentAgent], value)
		case "crawl-delay":
			if delay, err := time.ParseDuration(value + "s"); err == nil {
				robots.CrawlDelays[currentAgent] = delay
			}
		case "sitemap":
			robots.Sitemaps = append(robots.Sitemaps, value)
		}
	}
	
	return robots
}

// isAllowedByRobots checks if a URL is allowed by robots.txt
func (pc *ParallelCrawler) isAllowedByRobots(targetURL, domain string) bool {
	robots, exists := pc.robotsCache[domain]
	if !exists {
		return true // No robots.txt found, assume allowed
	}
	
	// Parse URL to get path
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	
	// Check rules for specific user agent first, then wildcard
	userAgents := []string{pc.config.UserAgent, "*"}
	
	for _, agent := range userAgents {
		if rules, exists := robots.UserAgentRules[agent]; exists {
			for _, rule := range rules {
				if rule != "" && strings.HasPrefix(parsedURL.Path, rule) {
					return false // Disallowed
				}
			}
		}
	}
	
	return true // Allowed
}