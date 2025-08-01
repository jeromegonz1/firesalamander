package crawler

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jeromegonz1/firesalamander/internal/logger"
)

var robotsLog = logger.New("ROBOTS")

// RobotsTxt représente un fichier robots.txt parsé
type RobotsTxt struct {
	Rules         map[string]*RobotRules // Map[UserAgent]Rules
	Sitemaps      []string
	CrawlDelay    map[string]time.Duration // Map[UserAgent]Delay
	Host          string
	DefaultRules  *RobotRules
	mu            sync.RWMutex
}

// RobotRules contient les règles pour un user-agent spécifique
type RobotRules struct {
	UserAgent    string
	Allowed      []string
	Disallowed   []string
	CrawlDelay   time.Duration
}

// RobotsCache gère le cache des fichiers robots.txt
type RobotsCache struct {
	cache    map[string]*CacheEntry
	ttl      time.Duration
	mu       sync.RWMutex
}

// CacheEntry représente une entrée dans le cache
type CacheEntry struct {
	Robots    *RobotsTxt
	ExpiresAt time.Time
}

// NewRobotsCache crée un nouveau cache pour robots.txt
func NewRobotsCache(ttl time.Duration) *RobotsCache {
	robotsLog.Debug("Creating robots cache", map[string]interface{}{
		"ttl": ttl,
	})

	cache := &RobotsCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}

	// Nettoyage périodique du cache
	go cache.cleanup()

	return cache
}

// Get récupère un robots.txt du cache
func (rc *RobotsCache) Get(domain string) (*RobotsTxt, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, exists := rc.cache[domain]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	robotsLog.Debug("Robots cache hit", map[string]interface{}{"domain": domain})
	return entry.Robots, true
}

// Set ajoute un robots.txt au cache
func (rc *RobotsCache) Set(domain string, robots *RobotsTxt) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache[domain] = &CacheEntry{
		Robots:    robots,
		ExpiresAt: time.Now().Add(rc.ttl),
	}

	robotsLog.Debug("Robots cached", map[string]interface{}{
		"domain":     domain,
		"expires_at": time.Now().Add(rc.ttl),
	})
}

// cleanup nettoie périodiquement le cache
func (rc *RobotsCache) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for domain, entry := range rc.cache {
			if now.After(entry.ExpiresAt) {
				delete(rc.cache, domain)
				robotsLog.Debug("Removed expired robots from cache", map[string]interface{}{
					"domain": domain,
				})
			}
		}
		rc.mu.Unlock()
	}
}

// ParseRobotsTxt parse le contenu d'un fichier robots.txt
func ParseRobotsTxt(content string) (*RobotsTxt, error) {
	robotsLog.Debug("Parsing robots.txt")

	robots := &RobotsTxt{
		Rules:      make(map[string]*RobotRules),
		CrawlDelay: make(map[string]time.Duration),
		Sitemaps:   []string{},
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentRules *RobotRules
	var currentUserAgent string

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Ignorer les lignes vides et les commentaires
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parser la ligne
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		directive := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		// Supprimer les commentaires inline
		if idx := strings.Index(value, "#"); idx >= 0 {
			value = strings.TrimSpace(value[:idx])
		}

		switch directive {
		case "user-agent":
			// Nouvelle section user-agent
			userAgent := normalizeUserAgent(value)
			currentUserAgent = userAgent
			
			if _, exists := robots.Rules[userAgent]; !exists {
				currentRules = &RobotRules{
					UserAgent:  userAgent,
					Allowed:    []string{},
					Disallowed: []string{},
				}
				robots.Rules[userAgent] = currentRules
			} else {
				currentRules = robots.Rules[userAgent]
			}

			// Si c'est pour tous les robots
			if userAgent == "*" {
				robots.DefaultRules = currentRules
			}

		case "disallow":
			if currentRules != nil && value != "" {
				// Normaliser le path
				path := normalizePath(value)
				currentRules.Disallowed = append(currentRules.Disallowed, path)
			}

		case "allow":
			if currentRules != nil && value != "" {
				path := normalizePath(value)
				currentRules.Allowed = append(currentRules.Allowed, path)
			}

		case "crawl-delay":
			if currentUserAgent != "" {
				if delay, err := parseCrawlDelay(value); err == nil {
					robots.CrawlDelay[currentUserAgent] = delay
					if currentRules != nil {
						currentRules.CrawlDelay = delay
					}
				}
			}

		case "sitemap":
			if value != "" {
				robots.Sitemaps = append(robots.Sitemaps, value)
			}

		case "host":
			if value != "" {
				robots.Host = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning robots.txt: %w", err)
	}

	robotsLog.Debug("Robots.txt parsed", map[string]interface{}{
		"user_agents": len(robots.Rules),
		"sitemaps":    len(robots.Sitemaps),
	})

	return robots, nil
}

// IsAllowed vérifie si une URL est autorisée pour un user-agent donné
func (r *RobotsTxt) IsAllowed(userAgent, targetURL string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Parser l'URL
	u, err := url.Parse(targetURL)
	if err != nil {
		robotsLog.Error("Failed to parse URL", map[string]interface{}{
			"url":   targetURL,
			"error": err.Error(),
		})
		return false
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	// Chercher les règles spécifiques au user-agent
	userAgent = normalizeUserAgent(userAgent)
	rules := r.findRulesForUserAgent(userAgent)
	
	if rules == nil {
		// Pas de règles, tout est autorisé
		return true
	}

	// Vérifier les règles Allow d'abord (elles ont la priorité)
	for _, allowed := range rules.Allowed {
		if matchesPath(path, allowed) {
			robotsLog.Debug("URL allowed by robots.txt", map[string]interface{}{
				"url":        targetURL,
				"user_agent": userAgent,
				"rule":       allowed,
			})
			return true
		}
	}

	// Vérifier les règles Disallow
	for _, disallowed := range rules.Disallowed {
		if matchesPath(path, disallowed) {
			robotsLog.Debug("URL disallowed by robots.txt", map[string]interface{}{
				"url":        targetURL,
				"user_agent": userAgent,
				"rule":       disallowed,
			})
			return false
		}
	}

	// Par défaut, autorisé
	return true
}

// GetCrawlDelay retourne le crawl-delay pour un user-agent
func (r *RobotsTxt) GetCrawlDelay(userAgent string) time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userAgent = normalizeUserAgent(userAgent)
	
	// Chercher d'abord pour le user-agent spécifique
	if delay, exists := r.CrawlDelay[userAgent]; exists {
		return delay
	}

	// Sinon, utiliser le délai par défaut
	if delay, exists := r.CrawlDelay["*"]; exists {
		return delay
	}

	return 0
}

// GetSitemaps retourne la liste des sitemaps déclarés
func (r *RobotsTxt) GetSitemaps() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return append([]string{}, r.Sitemaps...)
}

// findRulesForUserAgent trouve les règles applicables pour un user-agent
func (r *RobotsTxt) findRulesForUserAgent(userAgent string) *RobotRules {
	// Chercher d'abord une correspondance exacte
	if rules, exists := r.Rules[userAgent]; exists {
		return rules
	}

	// Chercher une correspondance partielle
	for ua, rules := range r.Rules {
		if strings.Contains(strings.ToLower(userAgent), strings.ToLower(ua)) {
			return rules
		}
	}

	// Utiliser les règles par défaut
	return r.DefaultRules
}

// normalizeUserAgent normalise un user-agent
func normalizeUserAgent(ua string) string {
	ua = strings.TrimSpace(ua)
	ua = strings.ToLower(ua)
	
	// Gérer les wildcards
	if ua == "" || ua == "*" {
		return "*"
	}
	
	return ua
}

// normalizePath normalise un path
func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	
	// S'assurer que le path commence par /
	if !strings.HasPrefix(path, "/") && path != "*" {
		path = "/" + path
	}
	
	return path
}

// matchesPath vérifie si un path correspond à un pattern
func matchesPath(path, pattern string) bool {
	// Correspondance exacte
	if path == pattern {
		return true
	}

	// Wildcard complet
	if pattern == "*" || pattern == "/*" {
		return true
	}

	// Pattern se terminant par *
	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, prefix)
	}

	// Pattern se terminant par $
	if strings.HasSuffix(pattern, "$") {
		exactPath := pattern[:len(pattern)-1]
		return path == exactPath
	}

	// Le path doit commencer par le pattern
	return strings.HasPrefix(path, pattern)
}

// parseCrawlDelay parse une valeur de crawl-delay
func parseCrawlDelay(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	
	// Essayer de parser comme un float (secondes)
	var seconds float64
	if _, err := fmt.Sscanf(value, "%f", &seconds); err == nil {
		return time.Duration(seconds * float64(time.Second)), nil
	}

	// Essayer de parser comme une durée Go
	if duration, err := time.ParseDuration(value + "s"); err == nil {
		return duration, nil
	}

	return 0, fmt.Errorf("invalid crawl-delay: %s", value)
}