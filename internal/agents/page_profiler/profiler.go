package page_profiler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"firesalamander/internal/agents"
	"golang.org/x/net/html"
)

// PageProfiler analyse la structure HTML des pages
type PageProfiler struct {
	name string
}

// NewPageProfiler crée une nouvelle instance de PageProfiler
func NewPageProfiler() *PageProfiler {
	return &PageProfiler{
		name: "page_profiler",
	}
}

// Name retourne le nom de l'agent
func (p *PageProfiler) Name() string {
	return p.name
}

// Process analyse une page HTML et retourne un profil détaillé
func (p *PageProfiler) Process(ctx context.Context, input interface{}) (*agents.AgentResult, error) {
	request, ok := input.(PageRequest)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected PageRequest")
	}

	if request.HTML == "" {
		// Retourner profil vide mais valide
		return &agents.AgentResult{
			AgentName: p.name,
			Status:    "completed",
			Data: map[string]interface{}{
				"meta_tags":             map[string]string{},
				"headings":              HeadingStructure{},
				"images":                []ImageInfo{},
				"links":                 []LinkInfo{},
				"schema_markup":         SchemaInfo{},
				"content_stats":         ContentStats{WordCount: 0},
				"core_web_vitals_hints": CoreWebVitalsHints{},
			},
		}, nil
	}

	profile := p.analyzePage(request.HTML, request.URL)

	return &agents.AgentResult{
		AgentName: p.name,
		Status:    "completed",
		Data: map[string]interface{}{
			"meta_tags":             profile.MetaTags,
			"headings":              profile.Headings,
			"images":                profile.Images,
			"links":                 profile.Links,
			"schema_markup":         profile.SchemaMarkup,
			"content_stats":         profile.ContentStats,
			"core_web_vitals_hints": profile.CoreWebVitalsHints,
		},
	}, nil
}

// HealthCheck vérifie que l'agent peut fonctionner
func (p *PageProfiler) HealthCheck() error {
	// Vérifier que les dépendances sont disponibles
	testHTML := `<html><head><title>Test</title></head><body><p>Test</p></body></html>`
	_, err := html.Parse(strings.NewReader(testHTML))
	return err
}

// analyzePage effectue l'analyse complète d'une page HTML
func (p *PageProfiler) analyzePage(htmlContent, pageURL string) PageProfile {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		// Retourner profil vide en cas d'erreur de parsing
		return PageProfile{
			URL:                pageURL,
			MetaTags:           make(map[string]string),
			Headings:           HeadingStructure{},
			Images:             []ImageInfo{},
			Links:              []LinkInfo{},
			SchemaMarkup:       SchemaInfo{},
			ContentStats:       ContentStats{},
			CoreWebVitalsHints: CoreWebVitalsHints{},
		}
	}

	profile := PageProfile{
		URL:                pageURL,
		MetaTags:           p.extractMetaTags(doc),
		Headings:           p.extractHeadings(doc),
		Images:             p.extractImages(doc),
		Links:              p.extractLinks(doc, pageURL),
		SchemaMarkup:       p.extractSchemaMarkup(doc, htmlContent),
		ContentStats:       p.calculateContentStats(doc),
		CoreWebVitalsHints: p.generateCoreWebVitalsHints(doc, htmlContent),
	}

	return profile
}

// extractMetaTags extrait tous les meta tags
func (p *PageProfiler) extractMetaTags(doc *html.Node) map[string]string {
	metaTags := make(map[string]string)
	
	var findMeta func(*html.Node)
	findMeta = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" {
				if n.FirstChild != nil {
					metaTags["title"] = n.FirstChild.Data
				}
			} else if n.Data == "meta" {
				name := ""
				property := ""
				content := ""
				
				for _, attr := range n.Attr {
					switch attr.Key {
					case "name":
						name = attr.Val
					case "property":
						property = attr.Val
					case "content":
						content = attr.Val
					}
				}
				
				if name != "" {
					metaTags[name] = content
				} else if property != "" {
					metaTags[property] = content
				}
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMeta(c)
		}
	}
	
	findMeta(doc)
	return metaTags
}

// extractHeadings extrait la hiérarchie des titres H1-H6
func (p *PageProfiler) extractHeadings(doc *html.Node) HeadingStructure {
	headings := HeadingStructure{
		H1: []string{},
		H2: []string{},
		H3: []string{},
		H4: []string{},
		H5: []string{},
		H6: []string{},
	}
	
	var findHeadings func(*html.Node)
	findHeadings = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				if text := p.extractText(n); text != "" {
					headings.H1 = append(headings.H1, text)
				}
			case "h2":
				if text := p.extractText(n); text != "" {
					headings.H2 = append(headings.H2, text)
				}
			case "h3":
				if text := p.extractText(n); text != "" {
					headings.H3 = append(headings.H3, text)
				}
			case "h4":
				if text := p.extractText(n); text != "" {
					headings.H4 = append(headings.H4, text)
				}
			case "h5":
				if text := p.extractText(n); text != "" {
					headings.H5 = append(headings.H5, text)
				}
			case "h6":
				if text := p.extractText(n); text != "" {
					headings.H6 = append(headings.H6, text)
				}
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findHeadings(c)
		}
	}
	
	findHeadings(doc)
	return headings
}

// extractImages extrait toutes les informations des images
func (p *PageProfiler) extractImages(doc *html.Node) []ImageInfo {
	var images []ImageInfo
	
	var findImages func(*html.Node)
	findImages = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			img := ImageInfo{}
			
			for _, attr := range n.Attr {
				switch attr.Key {
				case "src":
					img.Src = attr.Val
					img.Format = p.getImageFormat(attr.Val)
				case "alt":
					img.Alt = attr.Val
				case "width":
					fmt.Sscanf(attr.Val, "%d", &img.Width)
				case "height":
					fmt.Sscanf(attr.Val, "%d", &img.Height)
				}
			}
			
			img.Size = p.categorizeImageSize(img.Width, img.Height)
			images = append(images, img)
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findImages(c)
		}
	}
	
	findImages(doc)
	return images
}

// extractLinks extrait tous les liens avec classification
func (p *PageProfiler) extractLinks(doc *html.Node, baseURL string) []LinkInfo {
	var links []LinkInfo
	
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			link := LinkInfo{}
			
			for _, attr := range n.Attr {
				switch attr.Key {
				case "href":
					link.Href = attr.Val
				case "title":
					link.Title = attr.Val
				case "target":
					link.Target = attr.Val
				}
			}
			
			if link.Href != "" {
				link.Text = p.extractText(n)
				link.Type = p.classifyLink(link.Href, baseURL)
				links = append(links, link)
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	
	findLinks(doc)
	return links
}

// extractSchemaMarkup extrait les données Schema.org
func (p *PageProfiler) extractSchemaMarkup(doc *html.Node, htmlContent string) SchemaInfo {
	schema := SchemaInfo{
		Microdata: []MicrodataInfo{},
		JsonLD:    []JsonLDInfo{},
	}
	
	// Extraire microdata
	var findMicrodata func(*html.Node)
	findMicrodata = func(n *html.Node) {
		if n.Type == html.ElementNode {
			itemType := ""
			hasItemScope := false
			
			for _, attr := range n.Attr {
				if attr.Key == "itemtype" {
					itemType = attr.Val
				} else if attr.Key == "itemscope" {
					hasItemScope = true
				}
			}
			
			if itemType != "" && hasItemScope {
				properties := make(map[string]string)
				p.extractItemProps(n, properties)
				
				schema.Microdata = append(schema.Microdata, MicrodataInfo{
					Type:       itemType,
					Properties: properties,
				})
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMicrodata(c)
		}
	}
	
	findMicrodata(doc)
	
	// Extraire JSON-LD
	schema.JsonLD = p.extractJsonLD(htmlContent)
	
	return schema
}

// calculateContentStats calcule les statistiques du contenu
func (p *PageProfiler) calculateContentStats(doc *html.Node) ContentStats {
	allText := p.extractAllText(doc)
	htmlContent := p.nodeToString(doc)
	
	words := strings.Fields(allText)
	paragraphs := p.countElements(doc, "p")
	lists := p.countElements(doc, "ul") + p.countElements(doc, "ol")
	
	var textDensity float64
	if len(htmlContent) > 0 {
		textDensity = float64(len(allText)) / float64(len(htmlContent))
	}
	
	return ContentStats{
		WordCount:      len(words),
		CharacterCount: len(allText),
		ParagraphCount: paragraphs,
		ListCount:      lists,
		TextDensity:    textDensity,
	}
}

// generateCoreWebVitalsHints génère des indices pour les Core Web Vitals
func (p *PageProfiler) generateCoreWebVitalsHints(doc *html.Node, htmlContent string) CoreWebVitalsHints {
	hints := CoreWebVitalsHints{
		LargestContentfulPaint: []string{},
		CumulativeLayoutShift:  []string{},
		FirstInputDelay:        []string{},
	}
	
	// LCP: Détecter grandes images
	var checkImages func(*html.Node)
	checkImages = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			width, height := 0, 0
			for _, attr := range n.Attr {
				switch attr.Key {
				case "width":
					fmt.Sscanf(attr.Val, "%d", &width)
				case "height":
					fmt.Sscanf(attr.Val, "%d", &height)
				}
			}
			
			if width > 1000 || height > 800 {
				hints.LargestContentfulPaint = append(hints.LargestContentfulPaint, "large image detected")
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			checkImages(c)
		}
	}
	
	checkImages(doc)
	
	// CLS: Détecter styles inline
	if strings.Contains(htmlContent, "style=") {
		hints.CumulativeLayoutShift = append(hints.CumulativeLayoutShift, "inline styles detected")
	}
	
	// FID: Détecter scripts
	if strings.Contains(htmlContent, "<script") {
		if strings.Contains(htmlContent, "async") {
			hints.FirstInputDelay = append(hints.FirstInputDelay, "async scripts found")
		} else {
			hints.FirstInputDelay = append(hints.FirstInputDelay, "blocking scripts detected")
		}
	}
	
	return hints
}

// Méthodes utilitaires
func (p *PageProfiler) extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	
	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result.WriteString(p.extractText(c))
	}
	return strings.TrimSpace(result.String())
}

func (p *PageProfiler) extractAllText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	
	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result.WriteString(p.extractAllText(c))
	}
	return result.String()
}

func (p *PageProfiler) getImageFormat(src string) string {
	ext := strings.ToLower(filepath.Ext(src))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".webp":
		return "webp"
	case ".svg":
		return "svg"
	default:
		return "unknown"
	}
}

func (p *PageProfiler) categorizeImageSize(width, height int) string {
	pixels := width * height
	if pixels > 500000 {
		return "large"
	} else if pixels > 50000 {
		return "medium"
	} else if pixels > 0 {
		return "small"
	}
	return "unknown"
}

func (p *PageProfiler) classifyLink(href, baseURL string) string {
	if strings.HasPrefix(href, "#") {
		return "anchor"
	}
	if strings.HasPrefix(href, "mailto:") {
		return "email"
	}
	
	parsedURL, err := url.Parse(href)
	if err != nil {
		return "unknown"
	}
	
	if parsedURL.IsAbs() {
		baseParsed, err := url.Parse(baseURL)
		if err != nil {
			return "external"
		}
		
		if parsedURL.Host == baseParsed.Host {
			return "internal"
		}
		return "external"
	}
	
	return "internal"
}

func (p *PageProfiler) countElements(doc *html.Node, tagName string) int {
	count := 0
	var countNodes func(*html.Node)
	countNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tagName {
			count++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			countNodes(c)
		}
	}
	countNodes(doc)
	return count
}

func (p *PageProfiler) extractItemProps(n *html.Node, props map[string]string) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "itemprop" {
				text := p.extractText(n)
				if text == "" && len(n.Attr) > 0 {
					// Si pas de texte, chercher dans les attributs
					for _, a := range n.Attr {
						if a.Key == "content" {
							text = a.Val
							break
						}
					}
				}
				props[attr.Val] = text
			}
		}
	}
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.extractItemProps(c, props)
	}
}

func (p *PageProfiler) extractJsonLD(htmlContent string) []JsonLDInfo {
	var jsonLDs []JsonLDInfo
	
	// Regex pour extraire JSON-LD - plus flexible avec les espaces
	re := regexp.MustCompile(`(?s)<script[^>]*type\s*=\s*["']application/ld\+json["'][^>]*>(.*?)</script>`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			jsonContent := strings.TrimSpace(match[1])
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(jsonContent), &data); err == nil {
				jsonLD := JsonLDInfo{
					Data: data,
				}
				
				if t, ok := data["@type"].(string); ok {
					jsonLD.Type = t
				}
				if c, ok := data["@context"].(string); ok {
					jsonLD.Context = c
				}
				
				jsonLDs = append(jsonLDs, jsonLD)
			}
		}
	}
	
	return jsonLDs
}

func (p *PageProfiler) nodeToString(n *html.Node) string {
	var buf strings.Builder
	html.Render(&buf, n)
	return buf.String()
}