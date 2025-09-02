package constants

// Test URLs for integration testing
const (
	TestURLExample         = "https://example.com"
	TestURLStorageTest     = "https://test-storage.com"
	TestURLIntegrationTest = "https://integration-test.fr"
)

// Integration test data constants
const (
	IntegrationTestDataTask    = "test_task"
	IntegrationTestDataInsight = "test_insight"
	IntegrationTestDataAction  = "test_action"
)

// Integration analysis types
const (
	IntegrationAnalysisTypeContent   = "content"
	IntegrationAnalysisTypeTechnical = "technical"
)

// Integration test categories
const (
	IntegrationTestCategoryPerformance = "performance"
)

// Integration agent names
const (
	IntegrationAgentSEO      = "seo"
	IntegrationAgentSemantic = "semantic"
)

// Analysis statuses
const (
	AnalysisStatusSuccess = "success"
	AnalysisStatusFailed  = "failed"
	AnalysisStatusPending = "pending"
)

// Analysis types
const (
	AnalysisTypeQuick = "quick"
	AnalysisTypeFull  = "full"
)

// Report formats
const (
	ReportFormatHTML = "html"
	ReportFormatJSON = "json"
	ReportFormatCSV  = "csv"
)

// Report types
const (
	ReportTypeExecutive = "executive"
	ReportTypeTechnical = "technical"
)

// SEO Analysis Constants
const (
	// Title length limits (SEO best practices)
	TitleMinLength = 30
	TitleMaxLength = 60
	
	// Tag analyzer title limits (same as above for consistency)
	TagTitleMinLength = 30
	TagTitleMaxLength = 60
	
	// Meta description length limits
	MetaDescMinLength = 120
	MetaDescMaxLength = 160
	TagMetaDescMinLength = 120
	TagMetaDescMaxLength = 160
	
	// Content analysis
	MinContentWords = 300
	OptimalContentWords = 1500
	
	// User agent
	DefaultUserAgent = "FireSalamander-Bot/1.0"
)

// SEO Status Constants
const (
	RealSEOStatusInfo    = "info"
	RealSEOStatusWarn    = "warning"
	RealSEOStatusWarning = "warning"
	RealSEOStatusError   = "error"
)

// SEO Scoring Constants
const (
	MaxTitleScore    = 100
	MaxMetaDescScore = 100
)

// SEO Error Messages
const (
	ErrorTitleEmpty       = "Title tag is empty"
	ErrorTitleMissing     = "Title tag is missing"
	ErrorMetaDescMissing  = "Meta description is missing"
	WarningTitleTooShort  = "Title is too short"
	WarningTitleTooLong   = "Title is too long"
	WarningMetaDescTooShort = "Meta description is too short"
	WarningMetaDescTooLong  = "Meta description is too long"
)

// Additional SEO Error Messages
const (
	ErrorMissingH1               = "Missing H1 tag"
	ErrorMultipleH1              = "Multiple H1 tags found"
	ErrorBrokenHeadingHierarchy  = "Broken heading hierarchy"
	ErrorAllImagesMissingAlt     = "All images are missing alt text"
	WarningMissingAltText        = "Some images are missing alt text"
)

// SEO Scoring Constants (Extended)
const (
	MaxHeadingScore     = 100
	MaxImageScore       = 100
	MaxPerformanceScore = 100
	MaxMobileScore      = 100
	MaxHTTPSScore       = 100
	MaxContentScore     = 100
)

// Performance Constants
const (
	SEOFastLoadTime       = 1000  // 1 second in ms
	SEOAcceptableLoadTime = 3000  // 3 seconds in ms
	SEOSlowLoadTime       = 5000  // 5 seconds in ms
	OptimalPageSize       = 1024000  // 1MB in bytes
	AcceptablePageSize    = 2048000  // 2MB in bytes
	LargePageSize         = 5120000  // 5MB in bytes
)

// SEO Priority Levels
const (
	SEOPriorityCritical = "critical"
	SEOPriorityHigh     = "high"
	SEOPriorityMedium   = "medium"
	SEOPriorityLow      = "low"
)

// SEO Impact Levels
const (
	SEOImpactHigh   = "high"
	SEOImpactMedium = "medium" 
	SEOImpactLow    = "low"
)

// SEO Effort Levels
const (
	EffortQuickWin  = "quick_win"
	EffortModerate  = "moderate"
	EffortComplex   = "complex"
)

// SEO Actions
const (
	ActionAddTitle              = "Add title tag"
	ActionOptimizeTitleLength   = "Optimize title length"
	ActionAddMetaDescription    = "Add meta description"
	ActionAddH1                 = "Add H1 tag"
	ActionAddAltText            = "Add alt text to images"
)

// SEO Guides
const (
	GuideAddTitle              = "Add a descriptive title tag between 30-60 characters"
	GuideOptimizeTitleLength   = "Adjust title length to be between 30-60 characters"
	GuideAddMetaDescription    = "Add a meta description between 120-160 characters"
	GuideAddH1                 = "Add a single H1 tag with your main keyword"
	GuideAddAltText            = "Add descriptive alt text to all images"
)

// SEO Estimates (in minutes - as strings for struct compatibility)
const (
	EstimateAddTitle      = "15 minutes"
	EstimateOptimizeTitle = "10 minutes"
	EstimateAddMetaDesc   = "20 minutes"
	EstimateAddH1         = "10 minutes"
	EstimateAddAltText    = "30 minutes"
)

// SEO Grade Thresholds
const (
	GradeAThreshold     = 90
	GradeBThreshold     = 80
	GradeBPlusThreshold = 70
	GradeCThreshold     = 60
	GradeCPlusThreshold = 50
	GradeDThreshold     = 40
)

// SEO Grades
const (
	SEOGradeAPlus = "A+"
	SEOGradeA     = "A"
	SEOGradeBPlus = "B+"
	SEOGradeB     = "B"
	SEOGradeC     = "C"
	SEOGradeD     = "D"
	SEOGradeF     = "F"
)

// Keyword Analysis
const (
	MinKeywordLength = 3
	MaxKeywordLength = 50
)

// Performance Analyzer Constants
const (
	ClientTimeout        = 30  // seconds
	OptimalLineLength    = 80  // characters
	AcceptableLoadTime   = 3000 // milliseconds
)

// Recommendation Engine Constants
const (
	RecommendationMaxRecommendations = 10
	RecommendationPlaceholderPattern = "{{%s}}"
)

// Recommendation Template IDs
const (
	RecommendationTemplateIDMissingTitle         = "MISSING_TITLE"
	RecommendationTemplateIDTitleLength          = "TITLE_LENGTH"
	RecommendationTemplateIDMissingMetaDesc      = "MISSING_META_DESC"
	RecommendationTemplateIDMetaDescLength       = "META_DESC_LENGTH"
	RecommendationTemplateIDMissingH1            = "MISSING_H1"
	RecommendationTemplateIDMultipleH1           = "MULTIPLE_H1"
	RecommendationTemplateIDHeadingHierarchy     = "HEADING_HIERARCHY"
	RecommendationTemplateIDMissingAltTags       = "MISSING_ALT_TAGS"
	RecommendationTemplateIDMissingViewport      = "MISSING_VIEWPORT"
	RecommendationTemplateIDMissingCanonical     = "MISSING_CANONICAL"
	RecommendationTemplateIDMissingOGTags        = "MISSING_OG_TAGS"
	RecommendationTemplateIDImproveLCP           = "IMPROVE_LCP"
	RecommendationTemplateIDImproveFID           = "IMPROVE_FID"
	RecommendationTemplateIDImproveCLS           = "IMPROVE_CLS"
	RecommendationTemplateIDEnableCompression    = "ENABLE_COMPRESSION"
	RecommendationTemplateIDConfigureCaching     = "CONFIGURE_CACHING"
	RecommendationTemplateIDOptimizeImages       = "OPTIMIZE_IMAGES"
	RecommendationTemplateIDMinifyResources      = "MINIFY_RESOURCES"
	RecommendationTemplateIDReducePageSize       = "REDUCE_PAGE_SIZE"
	RecommendationTemplateIDMigrateHTTPS         = "MIGRATE_HTTPS"
	RecommendationTemplateIDFixMixedContent      = "FIX_MIXED_CONTENT"
	RecommendationTemplateIDAddHSTS              = "ADD_HSTS"
	RecommendationTemplateIDMakeResponsive       = "MAKE_RESPONSIVE"
	RecommendationTemplateIDAddSitemap           = "ADD_SITEMAP"
	RecommendationTemplateIDAddRobotsTxt         = "ADD_ROBOTS_TXT"
	RecommendationTemplateIDRemoveNoIndex        = "REMOVE_NO_INDEX"
	RecommendationTemplateIDFixDuplicateContent  = "FIX_DUPLICATE_CONTENT"
	RecommendationTemplateIDImproveAccessibility = "IMPROVE_ACCESSIBILITY"
	RecommendationTemplateIDFixBrokenLinks       = "FIX_BROKEN_LINKS"
	RecommendationTemplateIDImproveInternalLinking = "IMPROVE_INTERNAL_LINKING"
)

// Recommendation Issues
const (
	RecommendationIssueTitleMissing = "Page is missing a title tag"
)

// Recommendation Ranges
const (
	RecommendationTitleRange    = "30-60 characters"
	RecommendationMetaDescRange = "120-160 characters"
)

// Recommendation Scores
const (
	RecommendationScorePoor = "poor"
)

// Recommendation Targets
const (
	RecommendationTargetLCP      = "< 2.5s"
	RecommendationTargetFID      = "< 100ms"
	RecommendationTargetCLS      = "< 0.1"
	RecommendationTargetPageSize = "< 2MB"
	RecommendationTargetLinks    = "≥ 10 internal links"
)

// Recommendation Thresholds
const (
	RecommendationAltTextThreshold      = 0.8  // 80% of images should have alt text
	RecommendationMaxPageSizeBytes      = 2048000 // 2MB
	RecommendationAccessibilityThreshold = 80
	RecommendationMinInternalLinks      = 10
	RecommendationRuleMissingMetaDesc   = 80
)

// Recommendation Defaults
const (
	RecommendationDefaultTitle       = "Improve SEO"
	RecommendationDefaultDescription = "SEO improvement recommendation"
	RecommendationCategoryGeneral    = "general"
)

// Recommendation Tags
const (
	RecommendationTagTechnical = "technical"
)

// Recommendation Time Estimates
const (
	RecommendationTimeLow      = "low"
	RecommendationTimeMedium   = "medium"
	RecommendationTimeHigh     = "high"
	RecommendationTimeVariable = "variable"
)

// Tag Analyzer Constants - Regex Patterns
const (
	TagRegexTitlePattern    = `<title[^>]*>(.*?)</title>`
	TagRegexMetaDescPattern = `<meta\s+name\s*=\s*["']description["']\s+content\s*=\s*["']([^"']*)["'][^>]*>`
	TagRegexURLPattern      = `https?://[^\s<>"']+`
	TagRegexImageExtPattern = `\.(jpg|jpeg|png|gif|webp|svg)$`
)

// Tag Analyzer Messages
const (
	TagMsgTitleMissing       = "Title tag is missing"
	TagMsgTitleTooShort      = "Title is too short"
	TagMsgTitleTooLong       = "Title is too long"
	TagMsgTitleDuplicates    = "Duplicate title tags found"
	TagMsgMetaDescMissing    = "Meta description is missing"
	TagMsgMetaDescTooShort   = "Meta description is too short"
	TagMsgMetaDescTooLong    = "Meta description is too long"
	TagMsgNoH1               = "Missing H1 tag"
	TagMsgMultipleH1         = "Multiple H1 tags found"
	TagMsgBadHierarchy       = "Heading hierarchy is not logical"
	TagMsgViewportMissing    = "Viewport meta tag is missing"
)

// Tag Analyzer Recommendations
const (
	TagRecommendAddTitle           = "Add a descriptive title tag"
	TagRecommendExtendTitle        = "Extend the title to at least 30 characters"
	TagRecommendShortenTitle       = "Shorten the title to maximum 60 characters"
	TagRecommendAvoidDuplicates    = "Remove duplicate title tags"
	TagRecommendAddMetaDesc        = "Add a meta description"
	TagRecommendExtendMetaDesc     = "Extend meta description to at least 120 characters"
	TagRecommendShortenMetaDesc    = "Shorten meta description to maximum 160 characters"
	TagRecommendAddCTA             = "Add a call-to-action in meta description"
	TagRecommendAddH1              = "Add a single H1 tag"
	TagRecommendSingleH1           = "Use only one H1 tag per page"
	TagRecommendRespectHierarchy   = "Fix heading hierarchy (H1->H2->H3)"
	TagRecommendAddCanonical       = "Add canonical URL"
	TagRecommendAddViewport        = "Add viewport meta tag for mobile"
)

// Tag Analyzer Values
const (
	TagMinWordLength      = 5
	TagMetaNameDescription = "description"
	TagMetaNameRobots     = "robots"
	TagMetaNameViewport   = "viewport"
	TagValueCanonical     = "canonical"
	TagPrefixOG           = "og:"
)

// Technical Auditor Constants
const (
	DefaultRequestTimeout = 30000  // 30 seconds in milliseconds
	FastRequestTimeout    = 5000   // 5 seconds in milliseconds
	HTTPPrefix           = "http://"
	HTTPSPrefix          = "https://"
	HTTPStatusOK         = 200
	HTTPStatusBadRequest = 400
)

// Test Constants for cmd/server
const (
	TestQueryURLParam      = "url"
	TestQueryAnalyzeParam  = "analyze"
	TestQueryResultsParam  = "results"
)

// Test Constants for SEO tests
const (
	MaxSEOScore             = 100
	TestCanonicalExample    = "https://example.com/"
	TestCanonicalTest       = "https://test.com/"
	TestURLExternal         = "https://external.com"
	TestURLExampleTest      = "https://example-test.com"
	TestURLBroken           = "https://broken-link.invalid"
)

// Configuration Defaults (to replace hardcoding)
const (
	// Crawler Configuration
	DefaultMaxURLs = 300
	DefaultServerPort = 8080
	DefaultSemanticServiceURL = "http://localhost:5000"
	
	// Performance Thresholds
	FIDGoodThreshold = 100   // ms
	FIDNeedsImprovementThreshold = 300 // ms
)

// Performance test values (from performance_analyzer.go)
const (
	TestValue3000 = 3000
)