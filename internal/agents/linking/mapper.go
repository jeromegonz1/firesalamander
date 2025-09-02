package linking

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
	"firesalamander/internal/crawler"
)

// LinkingMapper implémente l'agent de cartographie et analyse des liens
type LinkingMapper struct {
	name string
}

// NewLinkingMapper crée une nouvelle instance de LinkingMapper
func NewLinkingMapper() *LinkingMapper {
	return &LinkingMapper{
		name: constants.AgentNameLinking,
	}
}

// Name retourne le nom de l'agent
func (l *LinkingMapper) Name() string {
	return l.name
}

// Process traite les données d'entrée et effectue la cartographie des liens
func (l *LinkingMapper) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()
	
	crawlResult, ok := data.(*crawler.CrawlResult)
	if !ok {
		return &agents.AgentResult{
			AgentName: l.name,
			Status:    constants.StatusFailed,
			Errors:    []string{"invalid input data type, expected *crawler.CrawlResult"},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	linkMap, err := l.MapLinks(crawlResult)
	if err != nil {
		return &agents.AgentResult{
			AgentName: l.name,
			Status:    constants.StatusFailed,
			Errors:    []string{err.Error()},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	return &agents.AgentResult{
		AgentName: l.name,
		Status:    constants.StatusCompleted,
		Data: map[string]interface{}{
			"link_map": linkMap,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}, nil
}

// HealthCheck vérifie la santé de l'agent
func (l *LinkingMapper) HealthCheck() error {
	// Test simple de cartographie
	testCrawlResult := &crawler.CrawlResult{
		Pages: []crawler.PageData{
			{
				URL:   "https://example.com",
				Title: "Test Page",
				Content: `<a href="https://example.com/page1">Internal Link</a>
						  <a href="https://external.com">External Link</a>`,
			},
		},
	}
	
	_, err := l.MapLinks(testCrawlResult)
	return err
}

// MapLinks crée une cartographie complète des liens du site
func (l *LinkingMapper) MapLinks(crawlResult *crawler.CrawlResult) (*agents.LinkMap, error) {
	if crawlResult == nil {
		return nil, fmt.Errorf("crawl result cannot be nil")
	}

	var internalLinks []agents.Link
	var externalLinks []agents.Link
	
	// Déterminer le domaine de base à partir des pages crawlées
	baseDomain := ""
	if len(crawlResult.Pages) > 0 {
		baseDomain = l.extractDomain(crawlResult.Pages[0].URL)
	}

	// Analyser chaque page
	for _, page := range crawlResult.Pages {
		pageLinks := l.extractLinksFromPage(page, baseDomain)
		
		for _, link := range pageLinks {
			if link.Type == "internal" {
				internalLinks = append(internalLinks, link)
			} else if link.Type == "external" {
				externalLinks = append(externalLinks, link)
			}
		}
	}

	// Générer les statistiques
	statistics := l.generateStatistics(internalLinks, externalLinks, len(crawlResult.Pages))

	return &agents.LinkMap{
		InternalLinks: internalLinks,
		ExternalLinks: externalLinks,
		Statistics:    statistics,
	}, nil
}

// AnalyzeLinkStructure analyse la structure des liens pour des insights SEO
func (l *LinkingMapper) AnalyzeLinkStructure(links []agents.Link) (*agents.LinkAnalysis, error) {
	if len(links) == 0 {
		return &agents.LinkAnalysis{
			LinkEquity:       make(map[string]float64),
			OrphanPages:      []string{},
			HighTrafficPages: []string{},
			Recommendations:  []string{"No links to analyze"},
		}, nil
	}

	linkEquity := l.calculateLinkEquity(links)
	orphanPages := l.findOrphanPages(links)
	highTrafficPages := l.identifyHighTrafficPages(links)
	recommendations := l.generateRecommendations(links, linkEquity, orphanPages)

	return &agents.LinkAnalysis{
		LinkEquity:       linkEquity,
		OrphanPages:      orphanPages,
		HighTrafficPages: highTrafficPages,
		Recommendations:  recommendations,
	}, nil
}

// extractLinksFromPage extrait tous les liens d'une page donnée
func (l *LinkingMapper) extractLinksFromPage(page crawler.PageData, baseDomain string) []agents.Link {
	var links []agents.Link

	// Expression régulière pour les liens
	linkRegex := regexp.MustCompile(`<a[^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := linkRegex.FindAllStringSubmatch(page.Content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			href := strings.TrimSpace(match[1])
			anchorText := strings.TrimSpace(l.cleanHTML(match[2]))
			
			// Skip empty or javascript links, but keep anchors
			if href == "" || strings.HasPrefix(href, "javascript:") {
				continue
			}

			// Déterminer le type de lien avant résolution pour les ancres
			linkType := l.determineLinkType(href, baseDomain)
			
			// Résoudre l'URL complète sauf pour les ancres
			fullURL := href
			if linkType != "anchor" {
				fullURL = l.resolveURL(href, page.URL)
				if fullURL == "" {
					continue
				}
				// Re-déterminer le type après résolution
				linkType = l.determineLinkType(fullURL, baseDomain)
			}
			
			// Détecter les attributs nofollow et noindex
			isNoFollow := strings.Contains(match[0], `rel="nofollow"`) || strings.Contains(match[0], `rel='nofollow'`)
			isNoIndex := strings.Contains(match[0], `rel="noindex"`) || strings.Contains(match[0], `rel='noindex'`)

			link := agents.Link{
				Source:     page.URL,
				Target:     fullURL,
				AnchorText: anchorText,
				Type:       linkType,
				IsNoFollow: isNoFollow,
				IsNoIndex:  isNoIndex,
			}

			// Ne collecter que les liens internes et externes pour les statistiques
			if linkType == "internal" || linkType == "external" {
				links = append(links, link)
			}
		}
	}

	return links
}

// resolveURL résout une URL relative en URL absolue
func (l *LinkingMapper) resolveURL(href, baseURL string) string {
	// URL déjà absolue
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// Parse relative URL
	rel, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Résoudre l'URL relative
	resolved := base.ResolveReference(rel)
	return resolved.String()
}

// determineLinkType détermine si un lien est interne, externe ou ancre
func (l *LinkingMapper) determineLinkType(linkURL, baseDomain string) string {
	if strings.HasPrefix(linkURL, "#") {
		return "anchor"
	}

	linkDomain := l.extractDomain(linkURL)
	if linkDomain == baseDomain {
		return "internal"
	}

	return "external"
}

// extractDomain extrait le domaine d'une URL
func (l *LinkingMapper) extractDomain(urlStr string) string {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return parsed.Host
}

// cleanHTML supprime les balises HTML du texte
func (l *LinkingMapper) cleanHTML(html string) string {
	// Suppression des balises HTML
	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(html, "")
	
	// Normalisation des espaces
	spaceRe := regexp.MustCompile(`\s+`)
	clean = spaceRe.ReplaceAllString(clean, " ")
	
	return strings.TrimSpace(clean)
}

// generateStatistics génère les statistiques de la cartographie des liens
func (l *LinkingMapper) generateStatistics(internalLinks, externalLinks []agents.Link, pageCount int) agents.LinkStatistics {
	totalLinks := len(internalLinks) + len(externalLinks)
	averageLinks := 0.0
	if pageCount > 0 {
		averageLinks = float64(totalLinks) / float64(pageCount)
	}

	return agents.LinkStatistics{
		TotalLinks:    totalLinks,
		InternalCount: len(internalLinks),
		ExternalCount: len(externalLinks),
		AverageLinks:  averageLinks,
	}
}

// calculateLinkEquity calcule l'équité des liens (PageRank simplifié)
func (l *LinkingMapper) calculateLinkEquity(links []agents.Link) map[string]float64 {
	linkEquity := make(map[string]float64)
	linkCounts := make(map[string]int)

	// Compter les liens entrants vers chaque page
	for _, link := range links {
		if link.Type == "internal" && !link.IsNoFollow {
			linkCounts[link.Target]++
		}
	}

	// Calculer l'équité basée sur le nombre de liens entrants
	for target, count := range linkCounts {
		// Équité simple basée sur le logarithme du nombre de liens
		equity := 1.0
		if count > 0 {
			equity = 1.0 + (float64(count) * 0.5)
		}
		linkEquity[target] = equity
	}

	return linkEquity
}

// findOrphanPages trouve les pages qui n'ont aucun lien entrant
func (l *LinkingMapper) findOrphanPages(links []agents.Link) []string {
	linkedPages := make(map[string]bool)
	allSources := make(map[string]bool)

	// Collecter toutes les pages sources et cibles
	for _, link := range links {
		allSources[link.Source] = true
		if link.Type == "internal" {
			linkedPages[link.Target] = true
		}
	}

	// Trouver les pages sources qui ne sont jamais des cibles
	var orphanPages []string
	for source := range allSources {
		if !linkedPages[source] {
			orphanPages = append(orphanPages, source)
		}
	}

	return orphanPages
}

// identifyHighTrafficPages identifie les pages avec beaucoup de liens entrants
func (l *LinkingMapper) identifyHighTrafficPages(links []agents.Link) []string {
	incomingLinks := make(map[string]int)

	// Compter les liens entrants
	for _, link := range links {
		if link.Type == "internal" {
			incomingLinks[link.Target]++
		}
	}

	// Identifier les pages avec plus de 5 liens entrants
	var highTrafficPages []string
	for page, count := range incomingLinks {
		if count >= 5 {
			highTrafficPages = append(highTrafficPages, page)
		}
	}

	// S'assurer que la slice n'est jamais nil
	if highTrafficPages == nil {
		highTrafficPages = []string{}
	}

	return highTrafficPages
}

// generateRecommendations génère des recommandations SEO basées sur l'analyse des liens
func (l *LinkingMapper) generateRecommendations(links []agents.Link, linkEquity map[string]float64, orphanPages []string) []string {
	var recommendations []string

	// Recommandations pour les pages orphelines
	if len(orphanPages) > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("Found %d orphan pages that need internal links for better SEO", len(orphanPages)))
	}

	// Recommandations pour l'équité des liens
	lowEquityCount := 0
	for _, equity := range linkEquity {
		if equity < 1.5 {
			lowEquityCount++
		}
	}

	if lowEquityCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("%d pages have low link equity and could benefit from more internal links", lowEquityCount))
	}

	// Recommandations pour les liens externes
	externalCount := 0
	noFollowCount := 0
	for _, link := range links {
		if link.Type == "external" {
			externalCount++
			if link.IsNoFollow {
				noFollowCount++
			}
		}
	}

	if externalCount > 0 {
		noFollowRatio := float64(noFollowCount) / float64(externalCount) * 100
		if noFollowRatio < 20 {
			recommendations = append(recommendations,
				"Consider adding rel='nofollow' to some external links to preserve link equity")
		}
	}

	// Recommandations pour le texte d'ancrage
	shortAnchorCount := 0
	for _, link := range links {
		if len(link.AnchorText) < 3 {
			shortAnchorCount++
		}
	}

	if shortAnchorCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("%d links have very short anchor text that could be more descriptive", shortAnchorCount))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Link structure appears to be well optimized")
	}

	return recommendations
}