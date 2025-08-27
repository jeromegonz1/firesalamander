package crawler

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"firesalamander/internal/logger"
)

var robotsLog = logger.New("ROBOTS")

// LegacyRobotsTxt représente un fichier robots.txt parsé (legacy version)
type LegacyRobotsTxt struct {
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
	Robots    *LegacyRobotsTxt
	ExpiresAt time.Time
}

// NewRobotsCache crée un nouveau cache pour robots.txt
func NewRobotsCache(ttl time.Duration) *RobotsCache {
	robotsLog.Debug("Creating robots cache", map[string]interface{}{
		"ttl": ttl,
	})

	return &RobotsCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}
}

// Get récupère un robots.txt du cache
func (rc *RobotsCache) Get(domain string) (*LegacyRobotsTxt, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, exists := rc.cache[domain]
	if !exists {
		return nil, false
	}

	// Vérifier l'expiration
	if time.Now().After(entry.ExpiresAt) {
		delete(rc.cache, domain)
		return nil, false
	}

	return entry.Robots, true
}

// Set stocke un robots.txt dans le cache
func (rc *RobotsCache) Set(domain string, robots *LegacyRobotsTxt) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache[domain] = &CacheEntry{
		Robots:    robots,
		ExpiresAt: time.Now().Add(rc.ttl),
	}

	robotsLog.Debug("Cached robots.txt", map[string]interface{}{
		"domain":     domain,
		"expires_at": time.Now().Add(rc.ttl),
	})
}

// ParseRobotsTxt parse le contenu d'un fichier robots.txt
func ParseRobotsTxt(content string) (*LegacyRobotsTxt, error) {
	robots := &LegacyRobotsTxt{
		Rules:      make(map[string]*RobotRules),
		Sitemaps:   []string{},
		CrawlDelay: make(map[string]time.Duration),
	}

	var currentUserAgent string
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Ignorer les commentaires et lignes vides
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Séparer clé et valeur
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "user-agent":
			currentUserAgent = value
			if robots.Rules[currentUserAgent] == nil {
				robots.Rules[currentUserAgent] = &RobotRules{
					UserAgent:  currentUserAgent,
					Allowed:    []string{},
					Disallowed: []string{},
				}
			}

		case "disallow":
			if currentUserAgent != "" {
				robots.Rules[currentUserAgent].Disallowed = append(
					robots.Rules[currentUserAgent].Disallowed, value)
			}

		case "allow":
			if currentUserAgent != "" {
				robots.Rules[currentUserAgent].Allowed = append(
					robots.Rules[currentUserAgent].Allowed, value)
			}

		case "crawl-delay":
			if currentUserAgent != "" {
				if delay, err := time.ParseDuration(value + "s"); err == nil {
					robots.Rules[currentUserAgent].CrawlDelay = delay
					robots.CrawlDelay[currentUserAgent] = delay
				}
			}

		case "sitemap":
			robots.Sitemaps = append(robots.Sitemaps, value)

		case "host":
			robots.Host = value
		}
	}

	// Règles par défaut pour *
	if rules, exists := robots.Rules["*"]; exists {
		robots.DefaultRules = rules
	}

	robotsLog.Debug("Parsed robots.txt", map[string]interface{}{
		"user_agents": len(robots.Rules),
		"sitemaps":    len(robots.Sitemaps),
		"host":        robots.Host,
	})

	return robots, scanner.Err()
}

// IsAllowed vérifie si une URL est autorisée pour un user-agent
func (rt *LegacyRobotsTxt) IsAllowed(userAgent, urlPath string) bool {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	// Normaliser l'URL
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return false
	}
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	// Vérifier les règles spécifiques au user-agent
	if rules, exists := rt.Rules[userAgent]; exists {
		return rt.checkRules(rules, path)
	}

	// Vérifier les règles par défaut (*)
	if rt.DefaultRules != nil {
		return rt.checkRules(rt.DefaultRules, path)
	}

	// Par défaut, autorisé
	return true
}

// checkRules vérifie les règles Allow/Disallow pour un path
func (rt *LegacyRobotsTxt) checkRules(rules *RobotRules, path string) bool {
	// D'abord vérifier les règles Allow (plus spécifiques)
	for _, allowPattern := range rules.Allowed {
		if rt.matchPattern(allowPattern, path) {
			return true
		}
	}

	// Ensuite vérifier les règles Disallow
	for _, disallowPattern := range rules.Disallowed {
		if rt.matchPattern(disallowPattern, path) {
			return false
		}
	}

	// Par défaut, autorisé
	return true
}

// matchPattern vérifie si un path correspond à un pattern robots.txt
func (rt *LegacyRobotsTxt) matchPattern(pattern, path string) bool {
	// Pattern vide = tout interdire
	if pattern == "" {
		return path == "/"
	}

	// Pattern exact
	if pattern == path {
		return true
	}

	// Wildcards simples
	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(path, prefix)
	}

	// Préfixe exact
	return strings.HasPrefix(path, pattern)
}

// GetCrawlDelay retourne le délai de crawl pour un user-agent
func (rt *LegacyRobotsTxt) GetCrawlDelay(userAgent string) time.Duration {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	if delay, exists := rt.CrawlDelay[userAgent]; exists {
		return delay
	}

	if delay, exists := rt.CrawlDelay["*"]; exists {
		return delay
	}

	return 0
}

// GetSitemaps retourne la liste des sitemaps
func (rt *LegacyRobotsTxt) GetSitemaps() []string {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.Sitemaps
}

// String retourne une représentation string du robots.txt
func (rt *LegacyRobotsTxt) String() string {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	var builder strings.Builder
	
	for userAgent, rules := range rt.Rules {
		builder.WriteString(fmt.Sprintf("User-agent: %s\n", userAgent))
		
		for _, disallow := range rules.Disallowed {
			builder.WriteString(fmt.Sprintf("Disallow: %s\n", disallow))
		}
		
		for _, allow := range rules.Allowed {
			builder.WriteString(fmt.Sprintf("Allow: %s\n", allow))
		}
		
		if rules.CrawlDelay > 0 {
			builder.WriteString(fmt.Sprintf("Crawl-delay: %.0f\n", rules.CrawlDelay.Seconds()))
		}
		
		builder.WriteString("\n")
	}
	
	for _, sitemap := range rt.Sitemaps {
		builder.WriteString(fmt.Sprintf("Sitemap: %s\n", sitemap))
	}

	return builder.String()
}