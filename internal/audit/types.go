package audit

type TechnicalAnalyzer struct {
	Rules TechRules
}

type TechRules struct {
	Title        TitleRules        `yaml:"title"`
	MetaDesc     MetaDescRules     `yaml:"meta_description"`
	Headings     HeadingRules      `yaml:"headings"`
	Images       ImageRules        `yaml:"images"`
	Links        LinkRules         `yaml:"links"`
	Performance  PerformanceRules  `yaml:"performance"`
	MeshAnalysis MeshAnalysisRules `yaml:"mesh_analysis"`
}

type TitleRules struct {
	MinLength        int    `yaml:"min_length"`
	MaxLength        int    `yaml:"max_length"`
	MissingSeverity  string `yaml:"missing_severity"`
	TooShortSeverity string `yaml:"too_short_severity"`
	TooLongSeverity  string `yaml:"too_long_severity"`
}

type MetaDescRules struct {
	MinLength        int    `yaml:"min_length"`
	MaxLength        int    `yaml:"max_length"`
	MissingSeverity  string `yaml:"missing_severity"`
	TooShortSeverity string `yaml:"too_short_severity"`
	TooLongSeverity  string `yaml:"too_long_severity"`
}

type HeadingRules struct {
	H1 H1Rules `yaml:"h1"`
	H2 H2Rules `yaml:"h2"`
	H3 H3Rules `yaml:"h3"`
}

type H1Rules struct {
	Required         bool   `yaml:"required"`
	MultipleSeverity string `yaml:"multiple_severity"`
	MissingSeverity  string `yaml:"missing_severity"`
}

type H2Rules struct {
	MinCount        int    `yaml:"min_count"`
	MissingSeverity string `yaml:"missing_severity"`
}

type H3Rules struct {
	Recommended bool `yaml:"recommended"`
}

type ImageRules struct {
	AltMissingSeverity string `yaml:"alt_missing_severity"`
	OversizedSeverity  string `yaml:"oversized_severity"`
	MaxSizeKB          int    `yaml:"max_size_kb"`
}

type LinkRules struct {
	BrokenSeverity     string   `yaml:"broken_severity"`
	WeakAnchorSeverity string   `yaml:"weak_anchor_severity"`
	WeakAnchors        []string `yaml:"weak_anchors"`
}

type PerformanceRules struct {
	LighthouseThresholds LighthouseThresholds `yaml:"lighthouse_thresholds"`
}

type LighthouseThresholds struct {
	Performance    Threshold `yaml:"performance"`
	Accessibility  Threshold `yaml:"accessibility"`
	BestPractices  Threshold `yaml:"best_practices"`
	SEO            Threshold `yaml:"seo"`
}

type Threshold struct {
	Good              float64 `yaml:"good"`
	NeedsImprovement  float64 `yaml:"needs_improvement"`
}

type MeshAnalysisRules struct {
	OrphanSeverity       string  `yaml:"orphan_severity"`
	MaxDepthWarning      int     `yaml:"max_depth_warning"`
	WeakAnchorThreshold  float64 `yaml:"weak_anchor_threshold"`
}

type TechResult struct {
	AuditID      string     `json:"audit_id"`
	ModelVersion string     `json:"model_version"`
	Scores       Scores     `json:"scores"`
	Findings     []Finding  `json:"findings"`
	Warnings     []Finding  `json:"warnings"`
	Mesh         MeshResult `json:"mesh"`
}

type Scores struct {
	Performance   float64 `json:"performance"`
	Accessibility float64 `json:"accessibility"`
	BestPractices float64 `json:"best_practices"`
	SEO           float64 `json:"seo"`
}

type Finding struct {
	ID       string   `json:"id"`
	Severity string   `json:"severity"`
	Message  string   `json:"message"`
	Evidence []string `json:"evidence"`
}

type MeshResult struct {
	Orphans     []string   `json:"orphans"`
	DepthStats  DepthStats `json:"depth_stats"`
	WeakAnchors []string   `json:"weak_anchors"`
}

type DepthStats struct {
	Min int     `json:"min"`
	Max int     `json:"max"`
	Avg float64 `json:"avg"`
}

type PageAudit struct {
	URL            string
	Title          string
	MetaDesc       string
	H1Count        int
	H2Count        int
	Images         []ImageInfo
	Links          []LinkInfo
	LighthouseData *LighthouseResult
}

type ImageInfo struct {
	Src string
	Alt string
}

type LinkInfo struct {
	Href   string
	Text   string
	Status int
}

type LighthouseResult struct {
	Performance   float64 `json:"performance"`
	Accessibility float64 `json:"accessibility"`
	BestPractices float64 `json:"best-practices"`
	SEO           float64 `json:"seo"`
}