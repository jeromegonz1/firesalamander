package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"firesalamander/internal/config"
	"golang.org/x/net/html"
)

type Crawler struct {
	Config     config.CrawlerConfig
	Visited    map[string]bool
	Queue      []CrawlTask
	Results    []PageData
	mutex      sync.RWMutex
	client     *http.Client
}

func NewCrawler(config config.CrawlerConfig) *Crawler {
	return &Crawler{
		Config:  config,
		Visited: make(map[string]bool),
		Queue:   make([]CrawlTask, 0),
		Results: make([]PageData, 0),
		client: &http.Client{
			Timeout: config.Performance.RequestTimeout,
		},
	}
}

func (c *Crawler) ShouldCrawlURL(urlStr string) bool {
	// Check extensions
	for _, ext := range c.Config.Exclusions.Extensions {
		if strings.HasSuffix(urlStr, ext) {
			return false
		}
	}

	// Check patterns
	for _, pattern := range c.Config.Exclusions.Patterns {
		if strings.Contains(urlStr, pattern) {
			return false
		}
	}

	return true
}

func (c *Crawler) RespectDepthLimit(depth int) bool {
	return depth <= c.Config.Limits.MaxDepth
}

func NormalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove fragment
	u.Fragment = ""

	// Remove UTM parameters
	q := u.Query()
	utmParams := []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content"}
	for _, param := range utmParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()

	// Remove trailing slash for non-root paths
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}

func ExtractContent(pageURL, htmlContent string, depth int) (*PageData, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	page := &PageData{
		URL:           pageURL,
		Depth:         depth,
		H2:            make([]string, 0),
		H3:            make([]string, 0),
		Anchors:       make([]Anchor, 0),
		OutgoingLinks: make([]string, 0),
		IncomingLinks: make([]string, 0),
		MetaIndex:     true, // default
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "html":
				// Extract language
				for _, attr := range n.Attr {
					if attr.Key == "lang" {
						page.Lang = attr.Val
					}
				}
			case "title":
				page.Title = getTextContent(n)
			case "link":
				// Extract canonical
				var rel, href string
				for _, attr := range n.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if rel == "canonical" {
					page.Canonical = href
				}
			case "meta":
				// Check robots meta
				var name, content string
				for _, attr := range n.Attr {
					if attr.Key == "name" {
						name = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if name == "robots" && strings.Contains(content, "noindex") {
					page.MetaIndex = false
				}
			case "h1":
				if page.H1 == "" { // Only first H1
					page.H1 = getTextContent(n)
				}
			case "h2":
				page.H2 = append(page.H2, getTextContent(n))
			case "h3":
				page.H3 = append(page.H3, getTextContent(n))
			case "a":
				var href, text string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				text = getTextContent(n)
				if href != "" && text != "" {
					page.Anchors = append(page.Anchors, Anchor{
						Text: strings.TrimSpace(text),
						Href: href,
					})
					// Add to outgoing links if internal
					if isInternalLink(href, pageURL) {
						page.OutgoingLinks = append(page.OutgoingLinks, href)
					}
				}
			case "p", "div", "article", "section":
				// Extract main content
				content := getTextContent(n)
				if len(content) > 50 { // Meaningful content
					page.Content += content + " "
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	// Detect language if not set
	if page.Lang == "" {
		page.Lang = DetectLanguage(page.Content)
	}

	// Clean content
	page.Content = strings.TrimSpace(page.Content)

	return page, nil
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getTextContent(c)
	}
	return strings.TrimSpace(result)
}

func isInternalLink(href, baseURL string) bool {
	if strings.HasPrefix(href, "/") {
		return true
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	link, err := url.Parse(href)
	if err != nil {
		return false
	}

	return link.Host == base.Host
}

func DetectLanguage(text string) string {
	if text == "" {
		return "unknown"
	}

	// Simple French detection
	frenchWords := []string{"le", "la", "les", "un", "une", "des", "et", "à", "de", "du", "pour", "avec", "sur", "dans", "nous", "sommes"}
	englishWords := []string{"the", "and", "or", "to", "of", "in", "for", "with", "on", "at", "by", "we", "are", "an", "a"}

	text = strings.ToLower(text)
	words := strings.Fields(text)
	frenchCount := 0
	englishCount := 0

	for _, word := range words {
		for _, frWord := range frenchWords {
			if word == frWord {
				frenchCount++
			}
		}
		for _, enWord := range englishWords {
			if word == enWord {
				englishCount++
			}
		}
	}

	if frenchCount > englishCount {
		return "fr"
	} else if englishCount > frenchCount {
		return "en"
	}

	return "unknown"
}