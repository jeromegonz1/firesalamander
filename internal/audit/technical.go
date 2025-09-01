package audit

import (
	"fmt"
	"firesalamander/internal/crawler"
)

func NewTechnicalAnalyzer(rules TechRules) *TechnicalAnalyzer {
	return &TechnicalAnalyzer{
		Rules: rules,
	}
}

func (ta *TechnicalAnalyzer) ValidateTitle(title, url string) []Finding {
	var findings []Finding

	if title == "" {
		findings = append(findings, Finding{
			ID:       "missing-title",
			Severity: ta.Rules.Title.MissingSeverity,
			Message:  "Titre manquant",
			Evidence: []string{url},
		})
		return findings
	}

	titleLen := len(title)
	if titleLen < ta.Rules.Title.MinLength {
		findings = append(findings, Finding{
			ID:       "title-too-short",
			Severity: ta.Rules.Title.TooShortSeverity,
			Message:  fmt.Sprintf("Titre trop court (%d caractères, minimum %d)", titleLen, ta.Rules.Title.MinLength),
			Evidence: []string{url},
		})
	}

	if titleLen > ta.Rules.Title.MaxLength {
		findings = append(findings, Finding{
			ID:       "title-too-long",
			Severity: ta.Rules.Title.TooLongSeverity,
			Message:  fmt.Sprintf("Titre trop long (%d caractères, maximum %d)", titleLen, ta.Rules.Title.MaxLength),
			Evidence: []string{url},
		})
	}

	return findings
}

func (ta *TechnicalAnalyzer) ValidateHeadings(h1Count, h2Count int, url string) []Finding {
	var findings []Finding

	// Check H1
	if ta.Rules.Headings.H1.Required && h1Count == 0 {
		findings = append(findings, Finding{
			ID:       "missing-h1",
			Severity: ta.Rules.Headings.H1.MissingSeverity,
			Message:  "Balise H1 manquante",
			Evidence: []string{url},
		})
	}

	if h1Count > 1 {
		findings = append(findings, Finding{
			ID:       "multiple-h1",
			Severity: ta.Rules.Headings.H1.MultipleSeverity,
			Message:  fmt.Sprintf("Plusieurs balises H1 (%d trouvées)", h1Count),
			Evidence: []string{url},
		})
	}

	// Check H2
	if h2Count < ta.Rules.Headings.H2.MinCount {
		findings = append(findings, Finding{
			ID:       "missing-h2",
			Severity: ta.Rules.Headings.H2.MissingSeverity,
			Message:  "Balises H2 insuffisantes pour structurer le contenu",
			Evidence: []string{url},
		})
	}

	return findings
}

func (ta *TechnicalAnalyzer) AnalyzeMesh(pages []crawler.PageData) MeshResult {
	// Build incoming links map
	incomingLinks := make(map[string][]string)
	allURLs := make(map[string]bool)
	
	// Initialize maps
	for _, page := range pages {
		allURLs[page.URL] = true
		incomingLinks[page.URL] = make([]string, 0)
	}

	// Build incoming links
	for _, page := range pages {
		for _, outgoing := range page.OutgoingLinks {
			if allURLs[outgoing] {
				incomingLinks[outgoing] = append(incomingLinks[outgoing], page.URL)
			}
		}
	}

	// Find orphans
	var orphans []string
	for url := range allURLs {
		if len(incomingLinks[url]) == 0 && url != pages[0].URL { // Exclude home page
			orphans = append(orphans, url)
		}
	}

	// Calculate depth stats (simplified)
	var depths []int
	for range pages {
		depths = append(depths, 0) // TODO: calculate actual depth
	}

	depthStats := DepthStats{
		Min: 0,
		Max: len(pages) - 1, // Simplified
		Avg: float64(len(pages)) / 2.0,
	}

	// Find weak anchors
	var weakAnchors []string
	anchorCounts := make(map[string]int)
	
	for _, page := range pages {
		// TODO: Extract anchors from PageData - for now just check URL
		_ = page.URL // Use page to avoid compiler error
	}

	for anchor, count := range anchorCounts {
		if count > 1 {
			for _, weak := range ta.Rules.Links.WeakAnchors {
				if anchor == weak {
					weakAnchors = append(weakAnchors, anchor)
				}
			}
		}
	}

	return MeshResult{
		Orphans:     orphans,
		DepthStats:  depthStats,
		WeakAnchors: weakAnchors,
	}
}

func (ta *TechnicalAnalyzer) ComputeScores(results []LighthouseResult) Scores {
	if len(results) == 0 {
		return Scores{}
	}

	var performance, accessibility, bestPractices, seo float64

	for _, result := range results {
		performance += result.Performance
		accessibility += result.Accessibility
		bestPractices += result.BestPractices
		seo += result.SEO
	}

	count := float64(len(results))
	return Scores{
		Performance:   performance / count,
		Accessibility: accessibility / count,
		BestPractices: bestPractices / count,
		SEO:          seo / count,
	}
}

func (ta *TechnicalAnalyzer) Analyze(crawlResult *crawler.CrawlResult, auditID string) (*TechResult, error) {
	var allFindings []Finding
	var allWarnings []Finding
	var lighthouseResults []LighthouseResult

	// Analyze each page
	for _, page := range crawlResult.Pages {
		// Validate title
		titleFindings := ta.ValidateTitle(page.Title, page.URL)
		allFindings = append(allFindings, titleFindings...)

		// Validate headings
		h1Count := 0
		if page.H1 != "" {
			h1Count = 1
		}
		headingFindings := ta.ValidateHeadings(h1Count, len(page.H2), page.URL)
		allFindings = append(allFindings, headingFindings...)

		// Mock Lighthouse for now
		lighthouseResults = append(lighthouseResults, LighthouseResult{
			Performance:   0.85,
			Accessibility: 0.90,
			BestPractices: 0.88,
			SEO:          0.92,
		})
	}

	// Analyze mesh
	mesh := ta.AnalyzeMesh(crawlResult.Pages)

	// Compute scores
	scores := ta.ComputeScores(lighthouseResults)

	return &TechResult{
		AuditID:      auditID,
		ModelVersion: "tech-v1.0",
		Scores:       scores,
		Findings:     allFindings,
		Warnings:     allWarnings,
		Mesh:         mesh,
	}, nil
}

