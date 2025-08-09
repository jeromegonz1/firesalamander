package constants

// 🚨 FIRE SALAMANDER - REAL SEO ANALYZER Constants
// Zero Hardcoding Policy - All constants for real SEO analysis

// ===== SCORING CONSTANTS =====
const (
	// Maximum scores for each SEO component (totaling 100 points)
	MaxTitleScore       = 20
	MaxMetaDescScore    = 15
	MaxHeadingScore     = 15
	MaxImageScore       = 10
	MaxPerformanceScore = 10
	MaxMobileScore      = 10
	MaxHTTPSScore       = 10
	MaxContentScore     = 10
	MaxSEOScore         = 100

	// Title scoring thresholds
	TitleMinLength     = 30
	TitleMaxLength     = 60
	TitleOptimalLength = 50

	// Meta description scoring thresholds
	MetaDescMinLength     = 120
	MetaDescMaxLength     = 160
	MetaDescOptimalLength = 150

	// Content scoring thresholds
	MinContentWords    = 300
	OptimalContentWords = 800
	
	// SEO Performance thresholds (in milliseconds)
	SEOFastLoadTime       = 1500  // 1.5s
	SEOAcceptableLoadTime = 3000  // 3s
	SEOSlowLoadTime       = 5000  // 5s
	
	// Page size thresholds (in bytes)
	OptimalPageSize    = 500 * 1024     // 500KB
	AcceptablePageSize = 1024 * 1024    // 1MB
	LargePageSize      = 2 * 1024 * 1024 // 2MB
	
	// Request count thresholds
	OptimalRequestCount    = 30
	AcceptableRequestCount = 50
	ExcessiveRequestCount  = 80
)

// ===== ERROR MESSAGES =====
const (
	ErrorTitleMissing               = "Title tag missing"
	ErrorTitleEmpty                 = "Title tag is empty"
	ErrorMetaDescMissing            = "Meta description missing"
	ErrorMetaDescEmpty              = "Meta description is empty"
	ErrorMissingH1                  = "Missing H1 heading"
	ErrorMultipleH1                 = "Multiple H1 tags found"
	ErrorBrokenHeadingHierarchy     = "Broken heading hierarchy detected"
	ErrorNoImages                   = "No images found on page"
	ErrorAllImagesMissingAlt        = "All images missing alt text"
)

// ===== WARNING MESSAGES =====
const (
	WarningTitleTooShort            = "Title too short (less than 30 characters)"
	WarningTitleTooLong             = "Title too long (more than 60 characters)"
	WarningMetaDescTooShort         = "Meta description too short"
	WarningMetaDescTooLong          = "Meta description too long"
	WarningMissingAltText           = "Some images missing alt text"
	WarningPoorPerformance          = "Poor page performance detected"
	WarningLargePageSize            = "Page size too large"
	WarningTooManyRequests          = "Too many HTTP requests"
	WarningNoHTTPS                  = "Page not served over HTTPS"
	WarningNotMobileFriendly        = "Page not mobile-friendly"
)

// ===== SUCCESS MESSAGES =====
const (
	SuccessTitleOptimal             = "Title length is optimal"
	SuccessMetaDescOptimal          = "Meta description length is optimal"
	SuccessGoodHeadingStructure     = "Good heading hierarchy"
	SuccessAllImagesHaveAlt         = "All images have alt text"
	SuccessGoodPerformance          = "Good page performance"
	SuccessHTTPSEnabled             = "Page served over HTTPS"
	SuccessMobileFriendly           = "Page is mobile-friendly"
)

// ===== SEO RECOMMENDATION PRIORITIES =====
const (
	SEOPriorityCritical = "CRITICAL"
	SEOPriorityHigh     = "HIGH"
	SEOPriorityMedium   = "MEDIUM"
	SEOPriorityLow      = "LOW"
)

// ===== SEO RECOMMENDATION IMPACTS =====
const (
	SEOImpactHigh   = "HIGH"
	SEOImpactMedium = "MEDIUM"
	SEOImpactLow    = "LOW"
)

// ===== RECOMMENDATION EFFORTS =====
const (
	EffortQuickWin = "QUICK_WIN"    // < 30 minutes
	EffortModerate = "MODERATE"     // 1-4 hours
	EffortComplex  = "COMPLEX"      // > 1 day
)

// ===== SEO GRADE CONSTANTS =====
const (
	SEOGradeAPlus = "A+"
	SEOGradeA     = "A"
	SEOGradeBPlus = "B+"
	SEOGradeB     = "B"
	SEOGradeCPlus = "C+"
	SEOGradeC     = "C"
	SEOGradeD     = "D"
	SEOGradeF     = "F"
)

// ===== GRADE THRESHOLDS =====
const (
	GradeAThreshold     = 95
	GradeBThreshold     = 90
	GradeBPlusThreshold = 80
	GradeCThreshold     = 70
	GradeCPlusThreshold = 60
	GradeDThreshold     = 50
	// Below 50 = F
)

// ===== HTML PARSING CONSTANTS =====
const (
	HTMLTagTitle       = "title"
	HTMLTagMeta        = "meta"
	HTMLTagH1          = "h1"
	HTMLTagH2          = "h2"
	HTMLTagH3          = "h3"
	HTMLTagH4          = "h4"
	HTMLTagH5          = "h5"
	HTMLTagH6          = "h6"
	HTMLTagImg         = "img"
	HTMLTagLink        = "link"
	HTMLTagA           = "a"
	HTMLTagP           = "p"
	HTMLTagDiv         = "div"
	
	HTMLAttrName       = "name"
	HTMLAttrContent    = "content"
	HTMLAttrProperty   = "property"
	HTMLAttrSrc        = "src"
	HTMLAttrAlt        = "alt"
	HTMLAttrHref       = "href"
	HTMLAttrRel        = "rel"
	HTMLAttrType       = "type"
	HTMLAttrClass      = "class"
	HTMLAttrId         = "id"
)

// ===== META NAMES =====
const (
	MetaNameDescription = "description"
	MetaNameViewport    = "viewport"
	MetaNameRobots      = "robots"
	MetaNameKeywords    = "keywords"
	MetaNameAuthor      = "author"
	MetaPropertyOGTitle = "og:title"
	MetaPropertyOGDesc  = "og:description"
	MetaPropertyOGImage = "og:image"
	MetaPropertyOGType  = "og:type"
	MetaNameTwitterCard = "twitter:card"
	MetaNameTwitterTitle = "twitter:title"
)

// ===== LINK REL VALUES =====
const (
	RelCanonical   = "canonical"
	RelAlternate   = "alternate"
	RelNext        = "next"
	RelPrev        = "prev"
	RelStylesheet  = "stylesheet"
	RelIcon        = "icon"
)

// ===== RECOMMENDATION ACTIONS =====
const (
	ActionAddTitle              = "Add a descriptive title tag"
	ActionOptimizeTitleLength   = "Optimize title length (30-60 characters)"
	ActionAddMetaDescription    = "Add meta description"
	ActionOptimizeMetaLength    = "Optimize meta description length (120-160 characters)"
	ActionAddH1                 = "Add H1 heading"
	ActionFixHeadingHierarchy   = "Fix heading hierarchy structure"
	ActionAddAltText            = "Add alt text to images"
	ActionImprovePerformance    = "Improve page loading performance"
	ActionEnableHTTPS          = "Enable HTTPS protocol"
	ActionOptimizeMobile       = "Optimize for mobile devices"
	ActionReducePageSize       = "Reduce page size and optimize resources"
	ActionReduceRequests       = "Reduce number of HTTP requests"
)

// ===== RECOMMENDATION GUIDES =====
const (
	GuideAddTitle = "Add a <title> tag in the <head> section with 30-60 characters that describe the page content and include main keywords."
	GuideOptimizeTitleLength = "Modify the existing title to be between 30-60 characters. Current title is either too short or too long."
	GuideAddMetaDescription = "Add a <meta name=\"description\" content=\"...\"> tag with 120-160 characters that summarize the page content."
	GuideOptimizeMetaLength = "Adjust meta description length to 120-160 characters for optimal display in search results."
	GuideAddH1 = "Add exactly one <h1> tag that describes the main topic of the page. Use it as the primary heading."
	GuideFixHeadingHierarchy = "Organize headings in proper order: H1 > H2 > H3. Don't skip heading levels."
	GuideAddAltText = "Add descriptive alt attributes to all <img> tags. Describe what the image shows or its purpose."
	GuideImprovePerformance = "Optimize images, enable compression, use CDN, and minimize CSS/JS files to reduce load time."
	GuideEnableHTTPS = "Install SSL certificate and redirect all HTTP traffic to HTTPS for security and SEO benefits."
	GuideOptimizeMobile = "Use responsive design, set viewport meta tag, and test on mobile devices."
)

// ===== TIME ESTIMATES =====
const (
	EstimateAddTitle         = "15 minutes"
	EstimateOptimizeTitle    = "10 minutes"
	EstimateAddMetaDesc      = "20 minutes"
	EstimateOptimizeMeta     = "15 minutes"
	EstimateAddH1            = "10 minutes"
	EstimateFixHeadings      = "30 minutes"
	EstimateAddAltText       = "1-2 hours"
	EstimateImprovePerformance = "4-8 hours"
	EstimateEnableHTTPS      = "2-4 hours"
	EstimateOptimizeMobile   = "1-3 days"
)

// ===== ANALYSIS STATUS =====
const (
	AnalysisStatusPending    = "PENDING"
	AnalysisStatusInProgress = "IN_PROGRESS"
	AnalysisStatusCompleted  = "COMPLETED"
	AnalysisStatusFailed     = "FAILED"
)

// ===== REAL SEO STATUS SEVERITY LEVELS =====
const (
	RealSEOStatusInfo     = "INFO"
	RealSEOStatusWarning  = "WARNING"
	RealSEOStatusError    = "ERROR"
	RealSEOStatusCritical = "CRITICAL"
)

// ===== KEYWORD EXTRACTION =====
const (
	MinKeywordLength     = 3
	MaxKeywordLength     = 20
	MinKeywordFrequency  = 2
	StopWordsThreshold   = 100 // Most common words to ignore
)

// ===== CONTENT QUALITY THRESHOLDS =====
const (
	MinReadabilityScore    = 60  // Flesch Reading Ease
	OptimalReadabilityScore = 70
	MinSentenceLength      = 10
	MaxSentenceLength      = 25
	OptimalParagraphLength = 150 // characters
)

// ===== LINK ANALYSIS =====
const (
	MinInternalLinks     = 3
	OptimalInternalLinks = 10
	MaxExternalLinks     = 10
	MinAnchorTextLength  = 2
	MaxAnchorTextLength  = 60
)

// ===== MOBILE OPTIMIZATION =====
const (
	MinViewportWidth  = 320
	MaxViewportWidth  = 1920
	OptimalTapTarget  = 44 // pixels
	MinTextSize       = 16 // pixels
)

// ===== PERFORMANCE BUDGETS =====
const (
	MaxDOMElements      = 1500
	MaxDOMDepth         = 32
	MaxCSSRules         = 4000
	MaxJavaScriptSize   = 1024 * 1024 // 1MB
	MaxCSSSize          = 512 * 1024  // 512KB
	MaxImageSize        = 2 * 1024 * 1024 // 2MB per image
)

// ===== SEO SPECIFIC HTTP HEADERS =====
const (
	SEOHeaderCSP = "Content-Security-Policy"
	SEOHeaderCanonical = "canonical"
)

// ===== SEO SPECIFIC CONTENT TYPES =====
const (
	SEOContentTypeJS = "application/javascript"
)

// ===== IMAGE FORMATS =====
const (
	ImageFormatJPEG = "jpeg"
	ImageFormatPNG  = "png"
	ImageFormatWEBP = "webp"
	ImageFormatSVG  = "svg"
	ImageFormatGIF  = "gif"
)

// ===== LOG MESSAGES =====
const (
	LogAnalysisStart      = "Starting SEO analysis for URL: %s"
	LogAnalysisComplete   = "SEO analysis completed - Score: %.1f, Grade: %s, Duration: %v"
	LogAnalysisError      = "SEO analysis failed for %s: %v"
	LogTitleAnalysis      = "Title analysis: Present=%t, Length=%d, Score=%d"
	LogMetaAnalysis       = "Meta description analysis: Present=%t, Length=%d, Score=%d"
	LogHeadingAnalysis    = "Heading analysis: H1=%d, H2=%d, H3=%d, Score=%d"
	LogImageAnalysis      = "Image analysis: Total=%d, WithAlt=%d, Coverage=%.2f, Score=%d"
	LogPerformanceAnalysis = "Performance analysis: LoadTime=%.2fs, PageSize=%d bytes, Score=%d"
)