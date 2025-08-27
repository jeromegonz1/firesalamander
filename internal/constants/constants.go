package constants

import "time"

// ========================================
// APPLICATION CONSTANTS
// ========================================

const (
	AppName        = "Fire Salamander"
	AppIcon        = "🦎"
	PoweredBy      = "SEPTEO"
	PrimaryColor   = "#ff6136"
	AppVersion     = "2.0.0"
	VendorName     = "SEPTEO"
	AppDescription = "SEO Analysis Tool for SEPTEO"
)

// ========================================
// NETWORK & SERVER CONSTANTS
// ========================================

const (
	DefaultPort         = "8080"
	DefaultPortInt      = 8080
	DefaultHost         = "localhost"
	DefaultProtocol     = "http"
	DefaultSecureProtocol = "https"
	
	// Test Ports
	TestPort3000        = "3000"
	TestPort3000Int     = 3000
)

// ========================================
// LIMITS & THRESHOLDS
// ========================================

const (
	// Page Limits
	MaxPages           = 20
	MaxCrawlDepth      = 3
	DefaultMaxPages    = 20
	DefaultMinPages    = 10
	
	// Rate Limiting - Production Security
	RateLimitRequestsPerSecond = 1    // 1 request per second
	RateLimitBurst            = 5    // Burst of 5 requests allowed
	
	// Site Type Limits
	EcommerceSiteMaxPages = 700
	EcommerceSiteMinPages = 100
	TestSiteMinPages   = 20
	TestSiteMaxPages   = 50
	CorporateSiteMinPages = 30
	CorporateSiteMaxPages = 110
	TechnicalSiteMinPages = 15
	TechnicalSiteMaxPages = 80
	BlogSiteMinPages   = 25
	BlogSiteMaxPages   = 200
	
	// Quality Scores
	HighQualityScore      = 80
	MinCoverageThreshold  = 80.0
	QualityThreshold80    = 80
	ExcellentScore90      = 90
	GoodScore85          = 85
	AcceptableScore70    = 70
	
	// Performance Limits
	MaxPageSize2MB       = 2 * 1024 * 1024
	MaxPageSize1MB       = 1024 * 1024
	MaxPageSize500KB     = 500 * 1024
	OptimalLineLength    = 80
	MaxLineLength        = 120
	
	// Request Limits
	MaxRetries           = 3
	MaxRedirects         = 10
	MaxIdleConns         = 100
	MaxIdleConnsPerHost  = 10
)

// ========================================
// TIMEOUTS & DURATIONS
// ========================================

const (
	// Server Timeouts
	ServerReadTimeout     = 30 * time.Second
	ServerWriteTimeout    = 30 * time.Second
	ServerIdleTimeout     = 60 * time.Second
	ShutdownTimeout       = 30 * time.Second
	
	// Client Timeouts
	ClientTimeout         = 30 * time.Second
	MaxHTTPTimeout        = 30 * time.Second
	ClientKeepAlive       = 30 * time.Second
	ClientTLSTimeout      = 10 * time.Second
	ClientIdleTimeout     = 90 * time.Second
	ClientExpectTimeout   = 1 * time.Second
	
	// Request Timeouts
	DefaultRequestTimeout = 15 * time.Second
	FastRequestTimeout    = 5 * time.Second
	LongRequestTimeout    = 60 * time.Second
	DefaultRetryDelay     = 1 * time.Second
	
	// Cache Durations
	RobotsCacheDuration   = 7 * 24 * time.Hour
	DefaultCacheExpiry    = 24 * time.Hour
	CacheCleanupInterval  = 5 * time.Minute
	RobotsCleanupInterval = 1 * time.Hour
	
	// Analysis Timeouts
	DefaultSimulationSpeed = 500 * time.Millisecond
	FastSimulationSpeed   = 100 * time.Millisecond
	SlowSimulationSpeed   = 1000 * time.Millisecond
	
	// Performance Thresholds
	FastLoadTime          = 2 * time.Second
	AcceptableLoadTime    = 3 * time.Second
	SlowLoadTime          = 5 * time.Second
	FastResponseTime      = 200 * time.Millisecond
	AcceptableResponseTime = 500 * time.Millisecond
	SlowResponseTime      = 1 * time.Second
	
	// Test Performance Values
	TestMinResponseTime   = 45 * time.Millisecond
	TestMaxResponseTime   = 890 * time.Millisecond
	TestAvgResponseTime   = 145 * time.Millisecond
	TestP50ResponseTime   = 132 * time.Millisecond
	TestP90ResponseTime   = 178 * time.Millisecond
	TestP95ResponseTime   = 189 * time.Millisecond
	TestP99ResponseTime   = 195 * time.Millisecond
	
	// Test Durations
	TestDuration2Min      = 2 * time.Minute
	TestDuration5Min      = 5 * time.Minute
	TestRampUpTime       = 30 * time.Second
	TestRampDownTime     = 30 * time.Second
)

// ========================================
// PATHS & FILES
// ========================================

const (
	// Directories
	TemplatesDir          = "templates"
	StaticDir            = "static"
	DataDir              = "data"
	ConfigDir            = "config"
	ScriptsDir           = "scripts"
	
	// Configuration Files
	DefaultConfigFile     = "config.yaml"
	DefaultConfigPath     = "config.yaml"
	DefaultEnvFile        = ".env"
	ExampleConfigFile     = "config.example.yaml"
	
	// Template Files
	HomeTemplate         = "home.html"
	AnalyzingTemplate    = "analyzing.html"
	ResultsTemplate      = "results.html"
	BaseTemplate         = "base.html"
	
	// Log Files
	DefaultLogFile       = "fire-salamander.log"
	ErrorLogFile         = "errors.log"
	AccessLogFile        = "access.log"
)

// ========================================
// URLS & ENDPOINTS
// ========================================

const (
	// Protocol Prefixes
	HTTPSPrefix          = "https://"
	HTTPPrefix           = "http://"
	HTTPSPrefixLength    = 8
	HTTPPrefixLength     = 7
	
	// Test URLs
	TestExampleURL       = "https://example.com"
	TestDemoURL          = "https://test.com"
	TestDemoFrURL        = "https://demo.fr"
	TestLocalhost3000    = "http://localhost:3000"
	TestDefaultDomain    = "example.com"
	
	// API Paths
	APIBasePath          = "/api"
	APIVersion           = "v1"
	APIHealthPath        = "/health"
	APIAnalyzePath       = "/analyze"
	APIStatusPath        = "/status"
	APIResultsPath       = "/results"
	APIInfoPath          = "/info"
	
	// External APIs
	OpenAIAPIURL         = "https://api.openai.com/v1/chat/completions"
	
	// AI Configuration
	AIRequestTimeout = 30 * time.Second
	DefaultCacheTTL  = 1 * time.Hour
	DefaultMaxTokens = 1000
	
	// Documentation URLs
	GoogleTitleLinkDocs  = "https://developers.google.com/search/docs/appearance/title-link"
	GoogleSnippetDocs    = "https://developers.google.com/search/docs/appearance/snippet"
	GoogleHTTPSDocs      = "https://developers.google.com/search/docs/advanced/security/https"
	WebDevLCPDocs        = "https://web.dev/lcp/"
)

// ========================================
// STATUS & STATES
// ========================================

const (
	// HTTP Status Messages
	StatusOK             = "OK"
	StatusError          = "ERROR"
	StatusProcessing     = "PROCESSING"
	StatusComplete       = "COMPLETE"
	StatusFailed         = "FAILED"
	StatusPending        = "PENDING"
	StatusStarted        = "STARTED"
	StatusAnalyzing      = "ANALYZING"
	
	// Analysis Types
	AnalysisTypeFull     = "full"
	AnalysisTypeQuick    = "quick"
	AnalysisTypeSEO      = "seo"
	AnalysisTypeSemantic = "semantic"
	
	// Analysis Phases
	PhaseDiscovery       = "discovery"
	PhaseSEOAnalysis     = "seo_analysis"
	PhaseAIAnalysis      = "ai_analysis"
	PhaseReportGen       = "report_generation"
	
	// Run Modes
	RunModeWebOnly       = "Interface Web Uniquement"
	RunModeAPIOnly       = "API REST Uniquement"
	RunModeComplete      = "Complet (Web + API + Orchestrateur)"
)

// ========================================
// USER AGENTS & IDENTIFIERS
// ========================================

const (
	SEOBotUserAgent      = "Fire Salamander SEO Bot/1.0 (+https://firesalamander.dev)"
	TestUserAgent        = "FireSalamander-Test/2.0"
	DefaultUserAgent     = "Fire Salamander/2.0"
)

// ========================================
// PROGRESS & SIMULATION VALUES
// ========================================

const (
	// Progress Thresholds
	ProgressThreshold20  = 20
	ProgressThreshold40  = 40
	ProgressThreshold60  = 60
	ProgressThreshold80  = 80
	ProgressThreshold95  = 95
	ProgressThreshold100 = 100
	
	// Progress Control
	DefaultProgressStart = 0
	DefaultProgressEnd   = 100
	ProgressRandomStep   = 3
	
	// Phase Progress Ranges
	PhaseDiscoveryStart  = 0
	PhaseDiscoveryEnd    = 25
	PhaseSEOStart        = 25
	PhaseSEOEnd          = 60
	PhaseAIStart         = 60
	PhaseAIEnd           = 90
	PhaseReportStart     = 90
	PhaseReportEnd       = 100
	
	// Phase Speeds
	PhaseSEOSpeed        = 800 * time.Millisecond
	PhaseAISpeed         = 600 * time.Millisecond
	PhaseReportSpeed     = 400 * time.Millisecond
	
	// Simulation Factors
	PageDiscoveryFactor     = 0.8   // Factor for page discovery simulation
	AnalysisProgressRatio   = 50.0  // Ratio for analysis progress
	IssueAccumulationRate   = 0.3   // Rate of issue accumulation
	TimingVariation         = 200 * time.Millisecond // Random timing variation
)

// ========================================
// SAMPLE DATA & TEST VALUES
// ========================================

const (
	// Sample Report Scores
	SampleOverallScore   = 87
	SampleTechnicalScore = 85
	SamplePerformanceScore = 72
	SampleContentScore   = 90
	SampleMobileScore    = 88
	
	// Test Magic Numbers
	TestValue3000        = 3000
	TestMagicNumber42    = 42
	TestMagicNumber100   = 100
	TestMagicNumber256   = 256
	TestMagicNumber1024  = 1024
)