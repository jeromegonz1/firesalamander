package agents

import (
	"context"
)

// Agent représente l'interface commune pour tous les agents spécialisés
type Agent interface {
	Name() string
	Process(ctx context.Context, data interface{}) (*AgentResult, error)
	HealthCheck() error
}

// AgentResult représente le résultat standardisé d'un agent
type AgentResult struct {
	AgentName string                 `json:"agent_name"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data"`
	Errors    []string               `json:"errors,omitempty"`
	Duration  int64                  `json:"duration_ms"`
}

// KeywordAgent - Interface pour l'extraction de mots-clés SEO
type KeywordAgent interface {
	Agent
	ExtractKeywords(content string) (*KeywordResult, error)
	AnalyzeDensity(keywords []string, content string) (*DensityReport, error)
}

// KeywordResult représente les résultats d'extraction de mots-clés
type KeywordResult struct {
	Keywords   []KeywordItem `json:"keywords"`
	TotalCount int           `json:"total_count"`
	Language   string        `json:"language"`
}

// KeywordItem représente un mot-clé avec ses métriques
type KeywordItem struct {
	Term      string  `json:"term"`
	Count     int     `json:"count"`
	Density   float64 `json:"density"`
	Relevance float64 `json:"relevance"`
}

// DensityReport représente l'analyse de densité des mots-clés
type DensityReport struct {
	TotalWords     int                    `json:"total_words"`
	KeywordMetrics map[string]float64     `json:"keyword_metrics"`
	Recommendations []string              `json:"recommendations"`
}

// TechnicalAgent - Interface pour l'audit technique
type TechnicalAgent interface {
	Agent
	AuditPage(page *PageData) (*TechnicalReport, error)
	ValidateStructure(html string) (*StructureResult, error)
}

// PageData représente les données d'une page pour l'audit
type PageData struct {
	URL     string `json:"url"`
	HTML    string `json:"html"`
	Headers map[string]string `json:"headers"`
}

// TechnicalReport représente le rapport d'audit technique
type TechnicalReport struct {
	PageURL      string            `json:"page_url"`
	Performance  PerformanceScore  `json:"performance"`
	Accessibility AccessibilityScore `json:"accessibility"`
	SEO          SEOScore          `json:"seo"`
	Issues       []TechnicalIssue  `json:"issues"`
}

// PerformanceScore représente les métriques de performance
type PerformanceScore struct {
	Score     int   `json:"score"`
	LoadTime  int64 `json:"load_time_ms"`
	Resources int   `json:"resources_count"`
}

// AccessibilityScore représente les métriques d'accessibilité
type AccessibilityScore struct {
	Score  int      `json:"score"`
	Issues []string `json:"issues"`
}

// SEOScore représente les métriques SEO techniques
type SEOScore struct {
	Score           int      `json:"score"`
	MissingElements []string `json:"missing_elements"`
}

// TechnicalIssue représente un problème technique détecté
type TechnicalIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Element     string `json:"element,omitempty"`
}

// StructureResult représente les résultats de validation de structure
type StructureResult struct {
	Valid        bool              `json:"valid"`
	Errors       []StructureError  `json:"errors"`
	Warnings     []StructureError  `json:"warnings"`
	HeadingLevel int               `json:"heading_level"`
}

// StructureError représente une erreur de structure HTML
type StructureError struct {
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Message     string `json:"message"`
	Element     string `json:"element"`
}

// LinkingAgent - Interface pour la cartographie des liens
type LinkingAgent interface {
	Agent
	MapLinks(crawlData *CrawlData) (*LinkMap, error)
	AnalyzeLinkStructure(links []Link) (*LinkAnalysis, error)
}

// LinkMap représente la cartographie des liens
type LinkMap struct {
	InternalLinks []Link           `json:"internal_links"`
	ExternalLinks []Link           `json:"external_links"`
	Statistics    LinkStatistics   `json:"statistics"`
}

// Link représente un lien avec ses propriétés
type Link struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	AnchorText  string `json:"anchor_text"`
	Type        string `json:"type"` // "internal", "external", "anchor"
	IsNoFollow  bool   `json:"is_nofollow"`
	IsNoIndex   bool   `json:"is_noindex"`
}

// LinkStatistics représente les statistiques des liens
type LinkStatistics struct {
	TotalLinks    int     `json:"total_links"`
	InternalCount int     `json:"internal_count"`
	ExternalCount int     `json:"external_count"`
	AverageLinks  float64 `json:"average_links_per_page"`
}

// LinkAnalysis représente l'analyse de la structure des liens
type LinkAnalysis struct {
	LinkEquity      map[string]float64 `json:"link_equity"`
	OrphanPages     []string           `json:"orphan_pages"`
	HighTrafficPages []string           `json:"high_traffic_pages"`
	Recommendations []string           `json:"recommendations"`
}

// BrokenLinksAgent - Interface pour la détection de liens brisés
type BrokenLinksAgent interface {
	Agent
	CheckLinks(urls []string) (*BrokenLinksReport, error)
	ValidateLink(url string) (*LinkStatus, error)
}

// BrokenLinksReport représente le rapport de liens brisés
type BrokenLinksReport struct {
	TotalChecked int          `json:"total_checked"`
	BrokenCount  int          `json:"broken_count"`
	BrokenLinks  []BrokenLink `json:"broken_links"`
	CheckedAt    string       `json:"checked_at"`
}

// BrokenLink représente un lien brisé
type BrokenLink struct {
	URL        string   `json:"url"`
	StatusCode int      `json:"status_code"`
	Error      string   `json:"error,omitempty"`
	FoundOn    []string `json:"found_on"`
}

// LinkStatus représente le statut d'un lien
type LinkStatus struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	IsValid    bool   `json:"is_valid"`
	Error      string `json:"error,omitempty"`
	CheckedAt  string `json:"checked_at"`
}