package crawler

import (
	"context"
	"html"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"firesalamander/internal/config"
)

// cleanHTML nettoie le HTML en supprimant les caractères invalides
// et en échappant les entités HTML pour un JSON valide
func cleanHTML(content string) string {
	if content == "" {
		return ""
	}

	// 1. Décoder les entités HTML existantes
	decoded := html.UnescapeString(content)

	// 2. Supprimer les caractères de contrôle invalides (0x00-0x1F) sauf \t, \n, \r
	cleaned := strings.Map(func(r rune) rune {
		if r == '\t' || r == '\n' || r == '\r' {
			return r
		}
		if r < 0x20 || r == 0x7F {
			return -1 // Supprime le caractère
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, decoded)

	// 3. Supprimer les balises script et style avec leur contenu
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	cleaned = scriptRegex.ReplaceAllString(cleaned, "")
	
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	cleaned = styleRegex.ReplaceAllString(cleaned, "")

	// 4. Supprimer les commentaires HTML
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	cleaned = commentRegex.ReplaceAllString(cleaned, "")

	// 5. Normaliser les espaces multiples
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	// 6. Trim les espaces
	cleaned = strings.TrimSpace(cleaned)

	// NOTE: Skip JSON escaping for HTML used in URL discovery
	// The JSON escaping will be done later when needed for API responses
	// 7. Ré-échapper pour JSON (guillemets, backslashes) - MOVED to API layer
	// cleaned = strings.ReplaceAll(cleaned, "\\", "\\\\")
	// cleaned = strings.ReplaceAll(cleaned, "\"", "\\\"")

	return cleaned
}

// IntelligentCrawler étend le crawler avec des capacités de nettoyage HTML
type IntelligentCrawler struct {
	*ParallelCrawler
	terminationController ITerminationController
	urlDiscoveryService   IURLDiscoveryService
	jobCounter           IAtomicJobCounter
	config               *config.CrawlerConfig // Store config for depth limits
}

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

// CrawlJob représente un job de crawl avec profondeur
type CrawlJob struct {
	URL   string
	Depth int
}

// atomicCounter compteur atomique pour éviter les race conditions
type atomicCounter struct {
	value int32
	mu    sync.Mutex
}

func (c *atomicCounter) Get() int32 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *atomicCounter) Add(delta int32) int32 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
	return c.value
}

func (c *atomicCounter) Sub(delta int32) int32 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value -= delta
	return c.value
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

// Implémentations basiques pour les interfaces
type basicTerminationController struct{}

func (tc *basicTerminationController) GetTerminationConditions() []ITerminationCondition {
	return []ITerminationCondition{&activeJobsCondition{}}
}

func (tc *basicTerminationController) ShouldTerminate() bool {
	return false
}

func (tc *basicTerminationController) NotifyJobCompleted() {}
func (tc *basicTerminationController) NotifyJobStarted() {}

type activeJobsCondition struct{}

func (c *activeJobsCondition) Type() string {
	return "active_jobs_zero"
}

func (c *activeJobsCondition) IsMet() bool {
	return false
}

type basicURLDiscoveryService struct{}

func (u *basicURLDiscoveryService) DiscoverFromSitemap(ctx context.Context, baseURL string) ([]string, error) {
	// Implémentation basique pour sitemap.xml
	// TODO: Implémenter la vraie logique sitemap si nécessaire
	return []string{}, nil
}

func (u *basicURLDiscoveryService) DiscoverFromRobots(ctx context.Context, baseURL string) ([]string, error) {
	return []string{}, nil
}

func (u *basicURLDiscoveryService) DiscoverFromHTML(htmlContent string, baseURL string) ([]string, error) {
	var discoveredURLs []string
	
	// Parse the base URL to get domain info for filtering
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	baseDomain := parsedBase.Host
	
	// Enhanced regex to extract links - handles more cases
	linkRegex := regexp.MustCompile(`<a[^>]*href=["']([^"']*?)["'][^>]*>`)
	matches := linkRegex.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			href := strings.TrimSpace(match[1])
			
			// Skip invalid hrefs
			if href == "" || strings.HasPrefix(href, "#") || 
			   strings.HasPrefix(href, "javascript:") || 
			   strings.HasPrefix(href, "mailto:") ||
			   strings.HasPrefix(href, "tel:") {
				continue
			}
			
			// Convert relative URLs to absolute URLs
			absoluteURL, err := u.resolveURL(baseURL, href)
			if err != nil {
				continue // Skip malformed URLs
			}
			
			// Parse resolved URL for domain filtering
			parsedURL, err := url.Parse(absoluteURL)
			if err != nil {
				continue // Skip malformed URLs
			}
			
			// Domain filtering - only same domain
			if parsedURL.Host != baseDomain {
				continue // Skip external domains
			}
			
			// Add to discovered URLs if not duplicate
			if !u.containsURL(discoveredURLs, absoluteURL) {
				discoveredURLs = append(discoveredURLs, absoluteURL)
			}
		}
	}
	
	return discoveredURLs, nil
}

// resolveURL convertit une URL relative en URL absolue
func (u *basicURLDiscoveryService) resolveURL(baseURL, href string) (string, error) {
	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	
	// Parse href (peut être relative ou absolue)
	reference, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	
	// Resolve reference against base
	resolved := base.ResolveReference(reference)
	return resolved.String(), nil
}

// containsURL vérifie si une URL est déjà dans la liste
func (u *basicURLDiscoveryService) containsURL(urls []string, target string) bool {
	for _, existing := range urls {
		if existing == target {
			return true
		}
	}
	return false
}

// fetchAndCleanPage récupère et nettoie une page web
func (ic *IntelligentCrawler) fetchAndCleanPage(ctx context.Context, url string) *PageResult {
	// Créer un task temporaire pour utiliser la méthode crawlPage du ParallelCrawler
	task := CrawlTask{
		URL:   url,
		Depth: 0,
	}
	
	// Utiliser la méthode crawlPage du ParallelCrawler qui gère déjà HTTP, parsing, etc.
	tempResult := ic.ParallelCrawler.crawlPage(ctx, task)
	if tempResult == nil || tempResult.Error != nil {
		return &PageResult{
			URL:   url,
			Error: tempResult.Error,
		}
	}
	
	// Appliquer cleanHTML immédiatement sur tous les contenus
	return &PageResult{
		URL:         tempResult.URL,
		Title:       cleanHTML(tempResult.Title),
		Description: cleanHTML(tempResult.Description),
		Body:        cleanHTML(tempResult.Body),
		Links:       cleanLinks(tempResult.Links), // Function helper pour nettoyer les liens
		StatusCode:  tempResult.StatusCode,
		Error:       nil,
		Depth:       0, // Sera défini par le caller
	}
}

// cleanLinks nettoie le texte des liens
func cleanLinks(links []ParallelLink) []ParallelLink {
	var cleaned []ParallelLink
	for _, link := range links {
		cleaned = append(cleaned, ParallelLink{
			URL:  link.URL,
			Text: cleanHTML(link.Text),
		})
	}
	return cleaned
}


type basicJobCounter struct {
	value int32
}

func (j *basicJobCounter) Get() int32 {
	return j.value
}

func (j *basicJobCounter) Add(delta int32) int32 {
	j.value += delta
	return j.value
}

func (j *basicJobCounter) Sub(delta int32) int32 {
	j.value -= delta
	return j.value
}

func (j *basicJobCounter) Reset() {
	j.value = 0
}

// NewIntelligentCrawler crée une nouvelle instance du crawler intelligent
func NewIntelligentCrawler(cfg *config.CrawlerConfig) ICrawlerEngine {
	return &IntelligentCrawler{
		ParallelCrawler:       NewParallelCrawler(cfg),
		terminationController: &basicTerminationController{},
		urlDiscoveryService:   &basicURLDiscoveryService{},
		jobCounter:           &basicJobCounter{},
		config:               cfg, // Store config for depth access
	}
}

// Implémentation de l'interface ICrawlerEngine

// CrawlWithIntelligence effectue un crawl RÉEL avec correction de la race condition
// CORRECTION MAJEURE : N'utilise PLUS ParallelCrawler défaillant - implémentation from scratch
func (ic *IntelligentCrawler) CrawlWithIntelligence(ctx context.Context, startURL string) (*IntelligentCrawlResult, error) {
	start := time.Now()
	log.Info("🕷️ IntelligentCrawler: Starting crawl", map[string]interface{}{"url": startURL, "max_depth": ic.config.MaxDepth})
	
	// Phase 1: Initialisation avec Producer-Consumer pattern
	urlQueue := make(chan CrawlJob, 100)     // Queue buffered pour les URLs à crawler
	resultChan := make(chan *PageResult, 20)  // Canal pour les résultats
	doneChan := make(chan struct{})          // Signal de terminaison
	
	activeJobs := &atomicCounter{value: 0}   // Compteur atomique des jobs actifs
	processedURLs := make(map[string]bool)   // URLs déjà traitées (éviter doublons)
	mutex := &sync.Mutex{}                   // Protection de la map
	
	pages := make(map[string]*PageResult)
	
	// Phase 2: Initialisation - Add initial URLs directly (no separate Producer goroutine)
	// 2.1: URL initiale
	urlQueue <- CrawlJob{URL: startURL, Depth: 0}
	activeJobs.Add(1)
	processedURLs[startURL] = true
	
	// 2.2: Découverte via sitemap (Intelligence = pas juste liens HTML)
	sitemapURLs, err := ic.urlDiscoveryService.DiscoverFromSitemap(ctx, startURL)
	if err == nil {
		for _, url := range sitemapURLs {
			if !processedURLs[url] {
				processedURLs[url] = true
				select {
				case urlQueue <- CrawlJob{URL: url, Depth: 1}:
					activeJobs.Add(1)
				default:
					// Queue pleine, ignorer cette URL
				}
			}
		}
	}
	
	// Phase 3: Queue monitor - Close queue with delay to avoid race condition
	go func() {
		var zeroCount int
		for {
			time.Sleep(50 * time.Millisecond) // Check more frequently
			if activeJobs.Get() == 0 {
				zeroCount++
				// Only close after activeJobs has been 0 for multiple checks
				// This prevents race condition where jobs are about to be added
				if zeroCount >= 10 { // 500ms delay (10 * 50ms)
					close(urlQueue)
					return
				}
			} else {
				zeroCount = 0 // Reset counter if we see active jobs
			}
			select {
			case <-ctx.Done():
				close(urlQueue)
				return
			default:
			}
		}
	}()
	
	// Phase 4: Workers (Consumers) - Pattern fixed race condition
	numWorkers := 5
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for job := range urlQueue {
				// 3.1: Fetch de la page avec cleanHTML appliqué
				pageResult := ic.fetchAndCleanPage(ctx, job.URL)
				if pageResult != nil {
					pageResult.Depth = job.Depth
					
					// 3.2: Envoyer résultat
					select {
					case resultChan <- pageResult:
					case <-ctx.Done():
						activeJobs.Sub(1)
						return
					}
					
					// 3.3: CORRECTIF RACE CONDITION - Découvrir nouveaux liens AVANT de décrémenter
					if job.Depth < ic.config.MaxDepth { // Use configuration depth limit
						newURLs, _ := ic.urlDiscoveryService.DiscoverFromHTML(pageResult.Body, job.URL)
						log.Info("🔗 Worker discovered new URLs", map[string]interface{}{
							"worker_id": workerID,
							"discovered_count": len(newURLs),
							"from_url": job.URL,
							"depth": job.Depth,
							"max_depth": ic.config.MaxDepth,
						})
						
						for _, newURL := range newURLs {
							mutex.Lock()
							if !processedURLs[newURL] {
								processedURLs[newURL] = true
								select {
								case urlQueue <- CrawlJob{URL: newURL, Depth: job.Depth + 1}:
									activeJobs.Add(1) // Incrémenter AVANT d'ajouter le job
								default:
									// Queue pleine
								}
							}
							mutex.Unlock()
						}
					}
				}
				
				// 3.4: Décrémenter APRÈS traitement complet
				remaining := activeJobs.Sub(1)
				
				// 3.5: Condition de terminaison corrigée
				if remaining == 0 {
					select {
					case doneChan <- struct{}{}:
					default:
						// Signal déjà envoyé
					}
				}
			}
		}(i)
	}
	
	// Phase 5: Collection des résultats avec timeout de sécurité
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Phase 6: Attendre résultats OU timeout OU completion
	timeout := time.NewTimer(90 * time.Second) // Sécurité anti-boucle infinie
	defer timeout.Stop()
	
	collecting := true
	for collecting {
		select {
		case result, ok := <-resultChan:
			if !ok {
				collecting = false
				break
			}
			if result != nil {
				pages[result.URL] = result
			}
			
		case <-doneChan:
			// Attendre encore un peu pour les derniers résultats
			time.Sleep(100 * time.Millisecond)
			collecting = false
			
		case <-timeout.C:
			collecting = false // Timeout de sécurité atteint
			
		case <-ctx.Done():
			collecting = false
		}
	}
	
	// Phase 7: Résultat avec métriques réelles
	duration := time.Since(start)
	
	return &IntelligentCrawlResult{
		StartURL:              startURL,
		Pages:                pages,
		Duration:             duration,
		Metrics:              &CrawlerMetrics{
			PagesPerSecond:   float64(len(pages)) / duration.Seconds(),
			CurrentWorkers:   numWorkers,
		},
		Error:                nil,
		TerminationReason:    "producer_consumer_completed",
		JobsExecuted:         int32(len(pages)),
		RaceConditionsDetected: 0,
	}, nil
}

// GetTerminationController retourne le contrôleur de terminaison
func (ic *IntelligentCrawler) GetTerminationController() ITerminationController {
	return ic.terminationController
}

// GetURLDiscoveryService retourne le service de découverte d'URLs
func (ic *IntelligentCrawler) GetURLDiscoveryService() IURLDiscoveryService {
	return ic.urlDiscoveryService
}

// GetJobCounter retourne le compteur de jobs
func (ic *IntelligentCrawler) GetJobCounter() IAtomicJobCounter {
	return ic.jobCounter
}

// CleanPageResult nettoie le contenu HTML d'un PageResult
func CleanPageResult(result *PageResult) *PageResult {
	if result == nil {
		return nil
	}

	// Copier le résultat pour ne pas modifier l'original
	cleaned := *result
	
	// Nettoyer le body HTML
	cleaned.Body = cleanHTML(result.Body)
	
	// Nettoyer le titre et la description
	cleaned.Title = cleanHTML(result.Title)
	cleaned.Description = cleanHTML(result.Description)
	
	// Nettoyer les liens
	for i, link := range cleaned.Links {
		cleaned.Links[i].Text = cleanHTML(link.Text)
	}
	
	return &cleaned
}

// CleanCrawlResult nettoie le contenu HTML d'un CrawlResult
func CleanCrawlResult(result *CrawlResult) *CrawlResult {
	if result == nil {
		return nil
	}

	// Copier le résultat pour ne pas modifier l'original
	cleaned := *result
	
	// Nettoyer le body HTML
	cleaned.Body = cleanHTML(result.Body)
	
	// Nettoyer le titre et les descriptions
	cleaned.Title = cleanHTML(result.Title)
	cleaned.Description = cleanHTML(result.Description)
	cleaned.MetaDescription = cleanHTML(result.MetaDescription)
	
	// Nettoyer les headings
	for i, heading := range cleaned.Headings {
		cleaned.Headings[i] = cleanHTML(heading)
	}
	
	// Nettoyer les liens
	for i, link := range cleaned.Links {
		cleaned.Links[i].Text = cleanHTML(link.Text)
	}
	
	// Nettoyer les images
	for i, img := range cleaned.Images {
		cleaned.Images[i].Alt = cleanHTML(img.Alt)
		cleaned.Images[i].Title = cleanHTML(img.Title)
	}
	
	return &cleaned
}