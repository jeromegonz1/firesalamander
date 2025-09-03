package page_profiler

// PageRequest représente la requête d'analyse d'une page
type PageRequest struct {
	HTML string `json:"html"`
	URL  string `json:"url"`
}

// PageProfile contient l'analyse complète d'une page
type PageProfile struct {
	URL                string                 `json:"url"`
	MetaTags           map[string]string      `json:"meta_tags"`
	Headings           HeadingStructure       `json:"headings"`
	Images             []ImageInfo            `json:"images"`
	Links              []LinkInfo             `json:"links"`
	SchemaMarkup       SchemaInfo             `json:"schema_markup"`
	ContentStats       ContentStats           `json:"content_stats"`
	CoreWebVitalsHints CoreWebVitalsHints     `json:"core_web_vitals_hints"`
}

// HeadingStructure représente la hiérarchie des titres
type HeadingStructure struct {
	H1 []string `json:"h1"`
	H2 []string `json:"h2"`
	H3 []string `json:"h3"`
	H4 []string `json:"h4"`
	H5 []string `json:"h5"`
	H6 []string `json:"h6"`
}

// ImageInfo contient les informations d'une image
type ImageInfo struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Size   string `json:"size,omitempty"` // "small", "medium", "large"
}

// LinkInfo contient les informations d'un lien
type LinkInfo struct {
	Href   string `json:"href"`
	Text   string `json:"text"`
	Title  string `json:"title,omitempty"`
	Type   string `json:"type"` // "internal", "external", "email", "anchor"
	Target string `json:"target,omitempty"`
}

// SchemaInfo contient les informations Schema.org
type SchemaInfo struct {
	Microdata []MicrodataInfo `json:"microdata"`
	JsonLD    []JsonLDInfo    `json:"json_ld"`
}

// MicrodataInfo représente les microdata Schema.org
type MicrodataInfo struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

// JsonLDInfo représente les JSON-LD Schema.org
type JsonLDInfo struct {
	Type    string                 `json:"@type"`
	Context string                 `json:"@context"`
	Data    map[string]interface{} `json:"data"`
}

// ContentStats contient les statistiques du contenu
type ContentStats struct {
	WordCount      int     `json:"word_count"`
	CharacterCount int     `json:"character_count"`
	ParagraphCount int     `json:"paragraph_count"`
	ListCount      int     `json:"list_count"`
	TextDensity    float64 `json:"text_density"` // Ratio texte/HTML total
}

// CoreWebVitalsHints contient des indices pour les Core Web Vitals
type CoreWebVitalsHints struct {
	LargestContentfulPaint []string `json:"largest_contentful_paint"`
	CumulativeLayoutShift  []string `json:"cumulative_layout_shift"`
	FirstInputDelay        []string `json:"first_input_delay"`
}