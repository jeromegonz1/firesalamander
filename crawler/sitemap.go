package crawler

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/jeromegonz1/firesalamander/internal/logger"
)

var sitemapLog = logger.New("SITEMAP")

// Sitemap représente un sitemap XML
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapIndex représente un index de sitemaps
type SitemapIndex struct {
	XMLName  xml.Name       `xml:"sitemapindex"`
	Xmlns    string         `xml:"xmlns,attr"`
	Sitemaps []SitemapEntry `xml:"sitemap"`
}

// SitemapEntry représente une entrée dans un sitemap index
type SitemapEntry struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod,omitempty"`
}

// SitemapURL représente une URL dans un sitemap
type SitemapURL struct {
	Loc        string  `xml:"loc"`
	Lastmod    string  `xml:"lastmod,omitempty"`
	Changefreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// SitemapParser gère le parsing des sitemaps
type SitemapParser struct {
	maxURLs int
}

// NewSitemapParser crée un nouveau parser de sitemap
func NewSitemapParser() *SitemapParser {
	return &SitemapParser{
		maxURLs: 50000, // Limite Google
	}
}

// Parse parse un sitemap XML
func (sp *SitemapParser) Parse(content string) (*Sitemap, error) {
	sitemapLog.Debug("Parsing sitemap")

	// Nettoyer le contenu
	content = strings.TrimSpace(content)
	
	// Vérifier si c'est un sitemap index
	if strings.Contains(content, "<sitemapindex") {
		return sp.parseSitemapIndex(content)
	}

	// Parser comme un sitemap normal
	var sitemap Sitemap
	if err := xml.Unmarshal([]byte(content), &sitemap); err != nil {
		sitemapLog.Error("Failed to parse sitemap", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to parse sitemap: %w", err)
	}

	// Valider et nettoyer les URLs
	validURLs := []SitemapURL{}
	for _, url := range sitemap.URLs {
		if url.Loc != "" {
			// Normaliser l'URL
			url.Loc = strings.TrimSpace(url.Loc)
			
			// Valider la priorité
			if url.Priority < 0 {
				url.Priority = 0
			} else if url.Priority > 1 {
				url.Priority = 1
			}

			// Valider changefreq
			if url.Changefreq != "" {
				url.Changefreq = sp.normalizeChangefreq(url.Changefreq)
			}

			validURLs = append(validURLs, url)
			
			// Limiter le nombre d'URLs
			if len(validURLs) >= sp.maxURLs {
				sitemapLog.Warn("Sitemap URL limit reached", map[string]interface{}{
					"limit": sp.maxURLs,
				})
				break
			}
		}
	}

	sitemap.URLs = validURLs

	sitemapLog.Info("Sitemap parsed successfully", map[string]interface{}{
		"urls_count": len(sitemap.URLs),
	})

	return &sitemap, nil
}

// parseSitemapIndex parse un sitemap index
func (sp *SitemapParser) parseSitemapIndex(content string) (*Sitemap, error) {
	sitemapLog.Debug("Parsing sitemap index")

	var index SitemapIndex
	if err := xml.Unmarshal([]byte(content), &index); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap index: %w", err)
	}

	// Convertir en sitemap simple avec les URLs des sous-sitemaps
	sitemap := &Sitemap{
		URLs: []SitemapURL{},
	}

	for _, entry := range index.Sitemaps {
		if entry.Loc != "" {
			sitemap.URLs = append(sitemap.URLs, SitemapURL{
				Loc:     entry.Loc,
				Lastmod: entry.Lastmod,
			})
		}
	}

	sitemapLog.Info("Sitemap index parsed", map[string]interface{}{
		"sitemaps_count": len(sitemap.URLs),
	})

	return sitemap, nil
}

// normalizeChangefreq normalise la valeur de changefreq
func (sp *SitemapParser) normalizeChangefreq(freq string) string {
	freq = strings.ToLower(strings.TrimSpace(freq))
	
	validFreqs := map[string]bool{
		"always":  true,
		"hourly":  true,
		"daily":   true,
		"weekly":  true,
		"monthly": true,
		"yearly":  true,
		"never":   true,
	}

	if validFreqs[freq] {
		return freq
	}

	return "weekly" // Valeur par défaut
}

// GetPriority retourne la priorité d'une URL (avec valeur par défaut)
func (u *SitemapURL) GetPriority() float64 {
	if u.Priority > 0 {
		return u.Priority
	}
	return 0.5 // Valeur par défaut selon le standard
}

// GetLastModified retourne la date de dernière modification parsée
func (u *SitemapURL) GetLastModified() (time.Time, error) {
	if u.Lastmod == "" {
		return time.Time{}, fmt.Errorf("no lastmod specified")
	}

	// Essayer différents formats de date
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, u.Lastmod); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse lastmod: %s", u.Lastmod)
}

// GetChangeFrequency retourne la fréquence de changement
func (u *SitemapURL) GetChangeFrequency() string {
	if u.Changefreq != "" {
		return u.Changefreq
	}
	return "weekly"
}

// FilterByPriority filtre les URLs par priorité minimale
func (s *Sitemap) FilterByPriority(minPriority float64) []SitemapURL {
	filtered := []SitemapURL{}
	
	for _, url := range s.URLs {
		if url.GetPriority() >= minPriority {
			filtered = append(filtered, url)
		}
	}

	return filtered
}

// FilterByAge filtre les URLs par âge maximum
func (s *Sitemap) FilterByAge(maxAge time.Duration) []SitemapURL {
	filtered := []SitemapURL{}
	cutoff := time.Now().Add(-maxAge)

	for _, url := range s.URLs {
		lastMod, err := url.GetLastModified()
		if err != nil {
			// Si pas de date, inclure l'URL
			filtered = append(filtered, url)
			continue
		}

		if lastMod.After(cutoff) {
			filtered = append(filtered, url)
		}
	}

	return filtered
}

// GetURLsByChangeFreq retourne les URLs groupées par fréquence de changement
func (s *Sitemap) GetURLsByChangeFreq() map[string][]SitemapURL {
	grouped := make(map[string][]SitemapURL)

	for _, url := range s.URLs {
		freq := url.GetChangeFrequency()
		grouped[freq] = append(grouped[freq], url)
	}

	return grouped
}

// Stats retourne des statistiques sur le sitemap
func (s *Sitemap) Stats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_urls": len(s.URLs),
	}

	// Compter par priorité
	priorityCount := make(map[string]int)
	for _, url := range s.URLs {
		priority := fmt.Sprintf("%.1f", url.GetPriority())
		priorityCount[priority]++
	}
	stats["priority_distribution"] = priorityCount

	// Compter par fréquence
	freqCount := make(map[string]int)
	for _, url := range s.URLs {
		freq := url.GetChangeFrequency()
		freqCount[freq]++
	}
	stats["changefreq_distribution"] = freqCount

	// URLs avec lastmod
	withLastmod := 0
	for _, url := range s.URLs {
		if url.Lastmod != "" {
			withLastmod++
		}
	}
	stats["urls_with_lastmod"] = withLastmod

	return stats
}