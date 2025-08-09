package semantic

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ContentExtractor extracteur de contenu HTML intelligent
type ContentExtractor struct {
	// Configuration d'extraction
	removeElements []string
	preserveElements []string
	
	// Regex pour nettoyage
	whitespaceRegex *regexp.Regexp
	urlRegex        *regexp.Regexp
	emailRegex      *regexp.Regexp
}

// ExtractedContent contenu extrait et structuré
type ExtractedContent struct {
	// Métadonnées de base
	Title           string   `json:"title"`
	MetaDescription string   `json:"meta_description"`
	MetaKeywords    string   `json:"meta_keywords"`
	Language        string   `json:"language"`
	Type            string   `json:"type"` // article, page, product, etc.
	
	// Contenu principal
	CleanText       string   `json:"clean_text"`
	RawHTML         string   `json:"raw_html"`
	
	// Structure
	Headings        []string `json:"headings"`
	HeadingStructure map[string][]string `json:"heading_structure"` // H1, H2, H3, etc.
	
	// Éléments spéciaux
	Links           []Link   `json:"links"`
	Images          []Image  `json:"images"`
	Lists           []string `json:"lists"`
	
	// Métriques
	WordCount       int      `json:"word_count"`
	CharCount       int      `json:"char_count"`
	
	// Analyse de structure
	ContentDensity  float64  `json:"content_density"`
	MainContentArea string   `json:"main_content_area"`
}

// Link lien extrait
type Link struct {
	URL        string `json:"url"`
	Text       string `json:"text"`
	Title      string `json:"title"`
	IsInternal bool   `json:"is_internal"`
	IsExternal bool   `json:"is_external"`
	Anchor     string `json:"anchor"`
}

// Image image extraite
type Image struct {
	URL    string `json:"url"`
	Alt    string `json:"alt"`
	Title  string `json:"title"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// NewContentExtractor crée un nouvel extracteur de contenu
func NewContentExtractor() *ContentExtractor {
	return &ContentExtractor{
		// Éléments à supprimer complètement
		removeElements: []string{
			"script", "style", "nav", "footer", "header", "aside",
			"form", "input", "button", "select", "textarea",
			"iframe", "object", "embed", "canvas", "svg",
		},
		
		// Éléments à préserver le contenu
		preserveElements: []string{
			"p", "div", "span", "h1", "h2", "h3", "h4", "h5", "h6",
			"article", "section", "main", "ul", "ol", "li",
			"blockquote", "pre", "code", "strong", "em", "b", "i",
		},
		
		whitespaceRegex: regexp.MustCompile(`\s+`),
		urlRegex:        regexp.MustCompile(`https?://[^\s]+`),
		emailRegex:      regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	}
}

// Extract extrait le contenu d'une page HTML
func (ce *ContentExtractor) Extract(htmlContent string) (*ExtractedContent, error) {
	log.Printf("Début extraction contenu HTML - Size:%d", len(htmlContent))

	// Parser le HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("erreur parsing HTML: %w", err)
	}

	content := &ExtractedContent{
		RawHTML:         htmlContent,
		HeadingStructure: make(map[string][]string),
		Links:           []Link{},
		Images:          []Image{},
		Lists:           []string{},
	}

	// Extraction des métadonnées
	ce.extractMetadata(doc, content)

	// Identification du contenu principal
	mainContent := ce.findMainContent(doc)

	// Extraction du contenu textuel
	ce.extractTextContent(mainContent, content)

	// Extraction des éléments structurels
	ce.extractStructuralElements(mainContent, content)

	// Nettoyage et finalisation
	ce.cleanAndFinalize(content)

	log.Printf("Extraction terminée - Words:%d Headings:%d Links:%d Images:%d", 
		content.WordCount, len(content.Headings), len(content.Links), len(content.Images))

	return content, nil
}

// extractMetadata extrait les métadonnées du document
func (ce *ContentExtractor) extractMetadata(doc *html.Node, content *ExtractedContent) {
	ce.walkHTML(doc, func(n *html.Node) {
		switch n.DataAtom {
		case atom.Title:
			if content.Title == "" && n.FirstChild != nil {
				content.Title = strings.TrimSpace(n.FirstChild.Data)
			}
			
		case atom.Meta:
			name := ce.getAttr(n, "name")
			property := ce.getAttr(n, "property")
			contentAttr := ce.getAttr(n, "content")
			
			switch {
			case name == "description":
				content.MetaDescription = contentAttr
			case name == "keywords":
				content.MetaKeywords = contentAttr
			case name == "language" || name == "lang":
				content.Language = contentAttr
			case property == "og:type":
				content.Type = contentAttr
			case property == "og:description" && content.MetaDescription == "":
				content.MetaDescription = contentAttr
			}
			
		case atom.Html:
			if lang := ce.getAttr(n, "lang"); lang != "" && content.Language == "" {
				content.Language = lang
			}
		}
	})

	// Valeurs par défaut
	if content.Language == "" {
		content.Language = "fr" // Par défaut français
	}
	if content.Type == "" {
		content.Type = "page"
	}
}

// findMainContent identifie la zone de contenu principal
func (ce *ContentExtractor) findMainContent(doc *html.Node) *html.Node {
	// Algorithme de détection du contenu principal
	
	// 1. Chercher les éléments sémantiques HTML5
	if main := ce.findNodeByAtom(doc, atom.Main); main != nil {
		return main
	}
	
	if article := ce.findNodeByAtom(doc, atom.Article); article != nil {
		return article
	}
	
	// 2. Chercher par ID/classe commune
	contentSelectors := []string{
		"content", "main-content", "primary", "main",
		"article", "post", "entry", "content-area",
	}
	
	for _, selector := range contentSelectors {
		if node := ce.findNodeByIDOrClass(doc, selector); node != nil {
			return node
		}
	}
	
	// 3. Chercher la div avec le plus de contenu textuel
	bestContent := ce.findContentByDensity(doc)
	if bestContent != nil {
		return bestContent
	}
	
	// 4. Fallback: utiliser le body entier
	return ce.findNodeByAtom(doc, atom.Body)
}

// extractTextContent extrait le contenu textuel
func (ce *ContentExtractor) extractTextContent(node *html.Node, content *ExtractedContent) {
	if node == nil {
		return
	}

	var textParts []string
	
	ce.walkHTML(node, func(n *html.Node) {
		// Ignorer les éléments à supprimer
		if ce.shouldRemoveElement(n) {
			return
		}
		
		// Extraire le texte des nœuds textuels
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textParts = append(textParts, text)
			}
		}
	})
	
	// Assembler le texte complet
	content.CleanText = strings.Join(textParts, " ")
	content.CleanText = ce.cleanText(content.CleanText)
	
	// Calculer les métriques
	content.WordCount = len(strings.Fields(content.CleanText))
	content.CharCount = len(content.CleanText)
}

// extractStructuralElements extrait les éléments structurels
func (ce *ContentExtractor) extractStructuralElements(node *html.Node, content *ExtractedContent) {
	if node == nil {
		return
	}

	ce.walkHTML(node, func(n *html.Node) {
		switch n.DataAtom {
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			heading := ce.extractTextFromNode(n)
			if heading != "" {
				content.Headings = append(content.Headings, heading)
				
				level := n.Data
				if content.HeadingStructure[level] == nil {
					content.HeadingStructure[level] = []string{}
				}
				content.HeadingStructure[level] = append(content.HeadingStructure[level], heading)
			}
			
		case atom.A:
			link := ce.extractLink(n)
			if link.URL != "" {
				content.Links = append(content.Links, link)
			}
			
		case atom.Img:
			image := ce.extractImage(n)
			if image.URL != "" {
				content.Images = append(content.Images, image)
			}
			
		case atom.Ul, atom.Ol:
			listContent := ce.extractList(n)
			if listContent != "" {
				content.Lists = append(content.Lists, listContent)
			}
		}
	})
}

// extractLink extrait les informations d'un lien
func (ce *ContentExtractor) extractLink(n *html.Node) Link {
	link := Link{
		URL:   ce.getAttr(n, "href"),
		Text:  ce.extractTextFromNode(n),
		Title: ce.getAttr(n, "title"),
	}
	
	// Déterminer le type de lien
	if strings.HasPrefix(link.URL, "http") {
		link.IsExternal = true
	} else if strings.HasPrefix(link.URL, "/") || !strings.Contains(link.URL, "://") {
		link.IsInternal = true
	}
	
	// Extraire l'ancre
	if parts := strings.Split(link.URL, "#"); len(parts) > 1 {
		link.Anchor = parts[1]
	}
	
	return link
}

// extractImage extrait les informations d'une image
func (ce *ContentExtractor) extractImage(n *html.Node) Image {
	return Image{
		URL:    ce.getAttr(n, "src"),
		Alt:    ce.getAttr(n, "alt"),
		Title:  ce.getAttr(n, "title"),
		Width:  ce.getAttr(n, "width"),
		Height: ce.getAttr(n, "height"),
	}
}

// extractList extrait le contenu d'une liste
func (ce *ContentExtractor) extractList(n *html.Node) string {
	var items []string
	
	ce.walkHTML(n, func(child *html.Node) {
		if child.DataAtom == atom.Li {
			item := ce.extractTextFromNode(child)
			if item != "" {
				items = append(items, item)
			}
		}
	})
	
	return strings.Join(items, "; ")
}

// cleanAndFinalize nettoie et finalise le contenu
func (ce *ContentExtractor) cleanAndFinalize(content *ExtractedContent) {
	// Calculer la densité de contenu
	if content.CharCount > 0 {
		content.ContentDensity = float64(content.WordCount) / float64(content.CharCount) * 100
	}
	
	// Identifier la zone de contenu principal
	content.MainContentArea = ce.identifyMainContentArea(content)
}

// identifyMainContentArea identifie la zone principale du contenu
func (ce *ContentExtractor) identifyMainContentArea(content *ExtractedContent) string {
	// Heuristique simple basée sur la structure
	if len(content.HeadingStructure["h1"]) > 0 {
		return "article"
	} else if len(content.Headings) > 3 {
		return "documentation"
	} else if len(content.Links) > content.WordCount/10 {
		return "navigation"
	} else {
		return "content"
	}
}

// cleanText nettoie le texte extrait
func (ce *ContentExtractor) cleanText(text string) string {
	// Supprimer les URLs
	text = ce.urlRegex.ReplaceAllString(text, "")
	
	// Supprimer les emails
	text = ce.emailRegex.ReplaceAllString(text, "")
	
	// Normaliser les espaces
	text = ce.whitespaceRegex.ReplaceAllString(text, " ")
	
	// Supprimer les caractères spéciaux en début/fin
	text = strings.TrimSpace(text)
	
	return text
}

// Fonctions utilitaires

func (ce *ContentExtractor) walkHTML(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ce.walkHTML(c, fn)
	}
}

func (ce *ContentExtractor) getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (ce *ContentExtractor) extractTextFromNode(n *html.Node) string {
	var texts []string
	
	ce.walkHTML(n, func(child *html.Node) {
		if child.Type == html.TextNode {
			text := strings.TrimSpace(child.Data)
			if text != "" {
				texts = append(texts, text)
			}
		}
	})
	
	return strings.Join(texts, " ")
}

func (ce *ContentExtractor) shouldRemoveElement(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	
	for _, elem := range ce.removeElements {
		if n.Data == elem {
			return true
		}
	}
	
	return false
}

func (ce *ContentExtractor) findNodeByAtom(doc *html.Node, atom atom.Atom) *html.Node {
	var result *html.Node
	
	ce.walkHTML(doc, func(n *html.Node) {
		if result == nil && n.DataAtom == atom {
			result = n
		}
	})
	
	return result
}

func (ce *ContentExtractor) findNodeByIDOrClass(doc *html.Node, selector string) *html.Node {
	var result *html.Node
	
	ce.walkHTML(doc, func(n *html.Node) {
		if result != nil {
			return
		}
		
		id := ce.getAttr(n, "id")
		class := ce.getAttr(n, "class")
		
		if strings.Contains(id, selector) || strings.Contains(class, selector) {
			result = n
		}
	})
	
	return result
}

func (ce *ContentExtractor) findContentByDensity(doc *html.Node) *html.Node {
	var bestNode *html.Node
	bestScore := 0
	
	ce.walkHTML(doc, func(n *html.Node) {
		if n.DataAtom == atom.Div {
			textLength := len(ce.extractTextFromNode(n))
			if textLength > bestScore {
				bestScore = textLength
				bestNode = n
			}
		}
	})
	
	return bestNode
}