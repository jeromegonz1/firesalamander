package crawler

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"firesalamander/internal/logger"
)

var cacheLog = logger.New("CACHE")

// PageCache implémente un cache LRU pour les pages crawlées
type PageCache struct {
	capacity  int
	ttl       time.Duration
	cache     map[string]*cacheItem
	evictList *list.List
	mu        sync.RWMutex
}

// cacheItem représente un élément dans le cache
type cacheItem struct {
	key       string
	value     *CrawlResult
	expiresAt time.Time
	element   *list.Element
}

// NewPageCache crée un nouveau cache de pages
func NewPageCache(ttl time.Duration) *PageCache {
	cacheLog.Debug("Creating page cache", map[string]interface{}{
		"ttl":      ttl,
		"capacity": 1000, // Capacité par défaut
	})

	cache := &PageCache{
		capacity:  1000,
		ttl:       ttl,
		cache:     make(map[string]*cacheItem),
		evictList: list.New(),
	}

	// Nettoyage périodique des entrées expirées
	go cache.cleanup()

	return cache
}

// Get récupère une page du cache
func (pc *PageCache) Get(key string) (*CrawlResult, bool) {
	pc.mu.RLock()
	item, exists := pc.cache[key]
	pc.mu.RUnlock()

	if !exists {
		return nil, false
	}

	// Vérifier l'expiration
	if time.Now().After(item.expiresAt) {
		pc.mu.Lock()
		pc.removeElement(item.element)
		pc.mu.Unlock()
		return nil, false
	}

	// Mettre à jour la position LRU
	pc.mu.Lock()
	pc.evictList.MoveToFront(item.element)
	pc.mu.Unlock()

	cacheLog.Debug("Cache hit", map[string]interface{}{"key": key})
	return item.value, true
}

// Set ajoute une page au cache
func (pc *PageCache) Set(key string, value *CrawlResult) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Si l'élément existe déjà, le mettre à jour
	if item, exists := pc.cache[key]; exists {
		item.value = value
		item.expiresAt = time.Now().Add(pc.ttl)
		pc.evictList.MoveToFront(item.element)
		return
	}

	// Si le cache est plein, évincer le plus ancien
	if pc.evictList.Len() >= pc.capacity {
		pc.removeOldest()
	}

	// Ajouter le nouvel élément
	item := &cacheItem{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(pc.ttl),
	}
	element := pc.evictList.PushFront(item)
	item.element = element
	pc.cache[key] = item

	cacheLog.Debug("Page cached", map[string]interface{}{
		"key":        key,
		"expires_at": item.expiresAt,
	})
}

// Remove supprime une page du cache
func (pc *PageCache) Remove(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if item, exists := pc.cache[key]; exists {
		pc.removeElement(item.element)
	}
}

// Clear vide complètement le cache
func (pc *PageCache) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache = make(map[string]*cacheItem)
	pc.evictList.Init()
	
	cacheLog.Info("Cache cleared")
}

// Size retourne le nombre d'éléments dans le cache
func (pc *PageCache) Size() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.evictList.Len()
}

// removeOldest supprime l'élément le moins récemment utilisé
func (pc *PageCache) removeOldest() {
	element := pc.evictList.Back()
	if element != nil {
		pc.removeElement(element)
	}
}

// removeElement supprime un élément du cache
func (pc *PageCache) removeElement(element *list.Element) {
	item := element.Value.(*cacheItem)
	delete(pc.cache, item.key)
	pc.evictList.Remove(element)
}

// cleanup nettoie périodiquement les entrées expirées
func (pc *PageCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		pc.mu.Lock()
		now := time.Now()
		
		// Créer une liste des éléments à supprimer
		var toRemove []*list.Element
		
		for element := pc.evictList.Back(); element != nil; element = element.Prev() {
			item := element.Value.(*cacheItem)
			if now.After(item.expiresAt) {
				toRemove = append(toRemove, element)
			}
		}

		// Supprimer les éléments expirés
		for _, element := range toRemove {
			pc.removeElement(element)
		}

		if len(toRemove) > 0 {
			cacheLog.Debug("Cleaned expired entries", map[string]interface{}{
				"removed": len(toRemove),
			})
		}
		
		pc.mu.Unlock()
	}
}

// Stats retourne des statistiques sur le cache
func (pc *PageCache) Stats() map[string]interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	stats := map[string]interface{}{
		"size":     pc.evictList.Len(),
		"capacity": pc.capacity,
		"ttl":      pc.ttl,
	}

	// Compter les entrées expirées
	expired := 0
	now := time.Now()
	for _, item := range pc.cache {
		if now.After(item.expiresAt) {
			expired++
		}
	}
	stats["expired"] = expired

	return stats
}

// RateLimiter implémente un limiteur de débit
type RateLimiter struct {
	rate     int           // Requêtes par unité de temps
	unit     time.Duration // Unité de temps (seconde, minute)
	tokens   chan struct{}
	ticker   *time.Ticker
	stopCh   chan struct{}
}

// NewRateLimiter crée un nouveau limiteur de débit
func NewRateLimiter(rateLimit string) (*RateLimiter, error) {
	// Parser le format "10/s" ou "60/m"
	var rate int
	var unit string
	if _, err := fmt.Sscanf(rateLimit, "%d/%s", &rate, &unit); err != nil {
		return nil, fmt.Errorf("invalid rate limit format: %s", rateLimit)
	}

	var duration time.Duration
	switch unit {
	case "s", "sec", "second":
		duration = time.Second
	case "m", "min", "minute":
		duration = time.Minute
	case "h", "hour":
		duration = time.Hour
	default:
		return nil, fmt.Errorf("unknown time unit: %s", unit)
	}

	cacheLog.Debug("Creating rate limiter", map[string]interface{}{
		"rate": rate,
		"unit": duration,
	})

	rl := &RateLimiter{
		rate:   rate,
		unit:   duration,
		tokens: make(chan struct{}, rate),
		stopCh: make(chan struct{}),
	}

	// Remplir initialement les tokens
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}

	// Démarrer le ticker pour renouveler les tokens
	rl.ticker = time.NewTicker(duration / time.Duration(rate))
	go rl.refillTokens()

	return rl, nil
}

// Wait attend qu'un token soit disponible
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.stopCh:
		return fmt.Errorf("rate limiter stopped")
	}
}

// refillTokens remplit périodiquement les tokens
func (rl *RateLimiter) refillTokens() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// Token ajouté
			default:
				// Canal plein, ignorer
			}
		case <-rl.stopCh:
			return
		}
	}
}

// Stop arrête le rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
	rl.ticker.Stop()
}

// CrawlQueue gère la file d'attente des URLs à crawler
type CrawlQueue struct {
	queue    *list.List
	seen     map[string]bool
	maxSize  int
	mu       sync.Mutex
}

// CrawlQueueItem représente un élément dans la queue
type CrawlQueueItem struct {
	URL   string
	Depth int
}

// NewCrawlQueue crée une nouvelle queue de crawl
func NewCrawlQueue(maxSize int) *CrawlQueue {
	return &CrawlQueue{
		queue:   list.New(),
		seen:    make(map[string]bool),
		maxSize: maxSize,
	}
}

// Add ajoute une URL à la queue
func (cq *CrawlQueue) Add(url string, depth int) bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	// Vérifier si l'URL a déjà été vue
	if cq.seen[url] {
		return false
	}

	// Vérifier la taille maximale
	if cq.queue.Len() >= cq.maxSize {
		return false
	}

	// Ajouter à la queue
	cq.queue.PushBack(&CrawlQueueItem{
		URL:   url,
		Depth: depth,
	})
	cq.seen[url] = true

	return true
}

// Next récupère le prochain élément de la queue
func (cq *CrawlQueue) Next() (*CrawlQueueItem, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	element := cq.queue.Front()
	if element == nil {
		return nil, false
	}

	cq.queue.Remove(element)
	item := element.Value.(*CrawlQueueItem)
	
	return item, true
}

// Size retourne la taille de la queue
func (cq *CrawlQueue) Size() int {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.queue.Len()
}

// HasSeen vérifie si une URL a déjà été vue
func (cq *CrawlQueue) HasSeen(url string) bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.seen[url]
}