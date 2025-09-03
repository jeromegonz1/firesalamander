package agents

// Common types shared between agents to avoid import cycles

// CrawlData represents crawl result data that can be passed between agents
type CrawlData struct {
	Pages    []PageInfo `json:"pages"`
	Metadata CrawlMeta  `json:"metadata"`
}

// PageInfo represents information about a crawled page
type PageInfo struct {
	URL           string            `json:"url"`
	Lang          string            `json:"lang"`
	Title         string            `json:"title"`
	H1            string            `json:"h1"`
	H2            []string          `json:"h2"`
	H3            []string          `json:"h3"`
	Anchors       []AnchorInfo      `json:"anchors"`
	Canonical     string            `json:"canonical"`
	MetaIndex     bool              `json:"meta_index"`
	Depth         int               `json:"depth"`
	OutgoingLinks []string          `json:"outgoing_links"`
	IncomingLinks []string          `json:"incoming_links"`
	Content       string            `json:"content"`
}

// AnchorInfo represents an anchor link
type AnchorInfo struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

// CrawlMeta represents metadata about the crawl operation
type CrawlMeta struct {
	TotalPages      int  `json:"total_pages"`
	MaxDepthReached int  `json:"max_depth_reached"`
	DurationMs      int  `json:"duration_ms"`
	RobotsRespected bool `json:"robots_respected"`
	SitemapFound    bool `json:"sitemap_found"`
}