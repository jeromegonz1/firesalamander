package constants

// ========================================
// DELTA-6 CRAWLER CONSTANTS
// Web Crawling and Testing Configuration Constants
// ========================================

// ========================================
// CRAWLER TEST URLS
// ========================================

// Standard Test URLs
const (
	CrawlerTestURLExample      = "https://example.com"
	CrawlerTestURLTest         = "https://test.com"
	CrawlerTestURLLocalhost    = "http://localhost"
	CrawlerTestURLHTTPBin      = "https://httpbin.org"
	CrawlerTestURLGeneric      = "https://example.org"
)

// Specific Test URLs
const (
	CrawlerTestURLExamplePage1     = "https://example.com/page1.html"
	CrawlerTestURLExamplePage2     = "https://example.com/page2.html"
	CrawlerTestURLExampleTest      = "https://example.com/test"
	CrawlerTestURLExampleAPI       = "https://example.com/api"
	CrawlerTestURLSitemapSchema    = "http://www.sitemaps.org/schemas/sitemap/0.9"
	CrawlerTestURLRobotsTXT        = "https://example.com/robots.txt"
	CrawlerTestURLSitemapXML       = "https://example.com/sitemap.xml"
)

// Local Test URLs
const (
	CrawlerTestURLLocalhost8080    = "http://localhost:8080"
	CrawlerTestURLLocalhost3000    = "http://localhost:3000"
	CrawlerTestURLLocalhostAPI     = "http://localhost:8080/api"
	CrawlerTestURLLocalhostHealth  = "http://localhost:8080/health"
)

// ========================================
// CRAWLER USER AGENTS
// ========================================

// Browser User Agents
const (
	CrawlerUserAgentChrome  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	CrawlerUserAgentFirefox = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0"
	CrawlerUserAgentSafari  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15"
	CrawlerUserAgentEdge    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59"
)

// Bot User Agents
const (
	CrawlerUserAgentBot        = "FireSalamander/1.0 (+https://github.com/fire-salamander/bot)"
	CrawlerUserAgentGoogleBot  = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	CrawlerUserAgentBingBot    = "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"
	CrawlerUserAgentDefault    = "FireSalamander-Crawler/1.0"
)

// ========================================
// CRAWLER HTTP METHODS
// ========================================

// HTTP Methods for Crawling
const (
	CrawlerHTTPMethodGet     = "GET"
	CrawlerHTTPMethodPost    = "POST"
	CrawlerHTTPMethodPut     = "PUT"
	CrawlerHTTPMethodDelete  = "DELETE"
	CrawlerHTTPMethodPatch   = "PATCH"
	CrawlerHTTPMethodHead    = "HEAD"
	CrawlerHTTPMethodOptions = "OPTIONS"
)

// ========================================
// CRAWLER HTTP STATUS CODES
// ========================================

// Success Status Codes
const (
	CrawlerHTTPStatusOK        = "200"
	CrawlerHTTPStatusCreated   = "201"
	CrawlerHTTPStatusAccepted  = "202"
	CrawlerHTTPStatusNoContent = "204"
)

// Redirection Status Codes
const (
	CrawlerHTTPStatusMovedPermanently = "301"
	CrawlerHTTPStatusFound           = "302"
	CrawlerHTTPStatusNotModified     = "304"
)

// Client Error Status Codes
const (
	CrawlerHTTPStatusBadRequest          = "400"
	CrawlerHTTPStatusUnauthorized        = "401"
	CrawlerHTTPStatusForbidden           = "403"
	CrawlerHTTPStatusNotFound            = "404"
	CrawlerHTTPStatusMethodNotAllowed    = "405"
	CrawlerHTTPStatusConflict            = "409"
	CrawlerHTTPStatusUnprocessableEntity = "422"
	CrawlerHTTPStatusTooManyRequests     = "429"
)

// Server Error Status Codes
const (
	CrawlerHTTPStatusInternalServerError = "500"
	CrawlerHTTPStatusNotImplemented      = "501"
	CrawlerHTTPStatusBadGateway          = "502"
	CrawlerHTTPStatusServiceUnavailable  = "503"
	CrawlerHTTPStatusGatewayTimeout      = "504"
)

// ========================================
// CRAWLER CONTENT TYPES
// ========================================

// Text Content Types
const (
	CrawlerContentTypeHTML       = "text/html"
	CrawlerContentTypePlain      = "text/plain"
	CrawlerContentTypeCSS        = "text/css"
	CrawlerContentTypeJavaScript = "text/javascript"
	CrawlerContentTypeXML        = "text/xml"
)

// Application Content Types
const (
	CrawlerContentTypeJSON            = "application/json"
	CrawlerContentTypeXMLApp          = "application/xml"
	CrawlerContentTypeJavaScriptApp   = "application/javascript"
	CrawlerContentTypePDF             = "application/pdf"
	CrawlerContentTypeOctetStream     = "application/octet-stream"
)

// Image Content Types
const (
	CrawlerContentTypeImage     = "image/"
	CrawlerContentTypeImagePNG  = "image/png"
	CrawlerContentTypeImageJPEG = "image/jpeg"
	CrawlerContentTypeImageGIF  = "image/gif"
	CrawlerContentTypeImageSVG  = "image/svg+xml"
	CrawlerContentTypeImageICO  = "image/x-icon"
)

// ========================================
// CRAWLER HTML ELEMENTS
// ========================================

// Document Structure Elements
const (
	CrawlerHTMLElementHTML   = "html"
	CrawlerHTMLElementHead   = "head"
	CrawlerHTMLElementBody   = "body"
	CrawlerHTMLElementTitle  = "title"
	CrawlerHTMLElementMeta   = "meta"
)

// Link and Media Elements
const (
	CrawlerHTMLElementLink   = "link"
	CrawlerHTMLElementScript = "script"
	CrawlerHTMLElementStyle  = "style"
	CrawlerHTMLElementA      = "a"
	CrawlerHTMLElementImg    = "img"
)

// Content Structure Elements
const (
	CrawlerHTMLElementDiv    = "div"
	CrawlerHTMLElementSpan   = "span"
	CrawlerHTMLElementP      = "p"
	CrawlerHTMLElementH1     = "h1"
	CrawlerHTMLElementH2     = "h2"
	CrawlerHTMLElementH3     = "h3"
	CrawlerHTMLElementH4     = "h4"
	CrawlerHTMLElementH5     = "h5"
	CrawlerHTMLElementH6     = "h6"
)

// List Elements
const (
	CrawlerHTMLElementUL     = "ul"
	CrawlerHTMLElementOL     = "ol"
	CrawlerHTMLElementLI     = "li"
)

// Table Elements
const (
	CrawlerHTMLElementTable  = "table"
	CrawlerHTMLElementTR     = "tr"
	CrawlerHTMLElementTD     = "td"
	CrawlerHTMLElementTH     = "th"
)

// Form Elements
const (
	CrawlerHTMLElementForm   = "form"
	CrawlerHTMLElementInput  = "input"
	CrawlerHTMLElementButton = "button"
	CrawlerHTMLElementSelect = "select"
	CrawlerHTMLElementOption = "option"
	CrawlerHTMLElementTextarea = "textarea"
)

// ========================================
// CRAWLER HTML ATTRIBUTES
// ========================================

// Link Attributes
const (
	CrawlerHTMLAttributeHref = "href"
	CrawlerHTMLAttributeSrc  = "src"
	CrawlerHTMLAttributeAlt  = "alt"
	CrawlerHTMLAttributeRel  = "rel"
)

// Identification Attributes
const (
	CrawlerHTMLAttributeID    = "id"
	CrawlerHTMLAttributeClass = "class"
	CrawlerHTMLAttributeName  = "name"
	CrawlerHTMLAttributeTitle = "title"
)

// Content Attributes
const (
	CrawlerHTMLAttributeContent = "content"
	CrawlerHTMLAttributeType    = "type"
	CrawlerHTMLAttributeCharset = "charset"
	CrawlerHTMLAttributeMedia   = "media"
	CrawlerHTMLAttributeStyle   = "style"
)

// Event Attributes
const (
	CrawlerHTMLAttributeOnclick = "onclick"
	CrawlerHTMLAttributeOnload  = "onload"
	CrawlerHTMLAttributeOnchange = "onchange"
	CrawlerHTMLAttributeOnsubmit = "onsubmit"
)

// ========================================
// CRAWLER CSS SELECTORS
// ========================================

// ID and Class Selectors
const (
	CrawlerCSSSelectorID         = "#"
	CrawlerCSSSelectorClass      = "."
	CrawlerCSSSelectorAttribute  = "["
	CrawlerCSSSelectorPseudo     = ":"
)

// Common CSS Selectors
const (
	CrawlerCSSSelectorLinks      = "a[href]"
	CrawlerCSSSelectorImages     = "img[src]"
	CrawlerCSSSelectorScripts    = "script[src]"
	CrawlerCSSSelectorStyles     = "link[rel='stylesheet']"
	CrawlerCSSSelectorMeta       = "meta[name]"
	CrawlerCSSSelectorTitle      = "title"
)

// ========================================
// CRAWLER XPATH EXPRESSIONS
// ========================================

// Common XPath Expressions
const (
	CrawlerXPathLinks       = "//a[@href]"
	CrawlerXPathImages      = "//img[@src]"
	CrawlerXPathScripts     = "//script[@src]"
	CrawlerXPathTitle       = "//title"
	CrawlerXPathMeta        = "//meta[@name]"
	CrawlerXPathHeadings    = "//h1 | //h2 | //h3 | //h4 | //h5 | //h6"
)

// ========================================
// CRAWLER ENCODING TYPES
// ========================================

// Character Encodings
const (
	CrawlerEncodingUTF8          = "UTF-8"
	CrawlerEncodingUTF16         = "UTF-16"
	CrawlerEncodingISO88591      = "ISO-8859-1"
	CrawlerEncodingWindows1252   = "Windows-1252"
	CrawlerEncodingASCII         = "ASCII"
)

// ========================================
// CRAWLER ROBOTS DIRECTIVES
// ========================================

// Robots.txt Directives
const (
	CrawlerRobotsUserAgent   = "User-agent"
	CrawlerRobotsDisallow    = "Disallow"
	CrawlerRobotsAllow       = "Allow"
	CrawlerRobotsCrawlDelay  = "Crawl-delay"
	CrawlerRobotsSitemap     = "Sitemap"
	CrawlerRobotsHost        = "Host"
)

// ========================================
// CRAWLER SITEMAP ELEMENTS
// ========================================

// XML Sitemap Elements
const (
	CrawlerSitemapURLSet     = "urlset"
	CrawlerSitemapURL        = "url"
	CrawlerSitemapLoc        = "loc"
	CrawlerSitemapLastmod    = "lastmod"
	CrawlerSitemapChangefreq = "changefreq"
	CrawlerSitemapPriority   = "priority"
)

// ========================================
// CRAWLER CONFIG KEYS
// ========================================

// Basic Configuration
const (
	CrawlerConfigTimeout         = "timeout"
	CrawlerConfigDelay           = "delay"
	CrawlerConfigConcurrency     = "concurrency"
	CrawlerConfigMaxDepth        = "max_depth"
	CrawlerConfigMaxPages        = "max_pages"
	CrawlerConfigUserAgent       = "user_agent"
)

// Behavior Configuration
const (
	CrawlerConfigFollowRedirects = "follow_redirects"
	CrawlerConfigRespectRobots   = "respect_robots"
	CrawlerConfigIgnoreSSL       = "ignore_ssl"
	CrawlerConfigAllowDuplicates = "allow_duplicates"
)

// Network Configuration
const (
	CrawlerConfigHeaders    = "headers"
	CrawlerConfigCookies    = "cookies"
	CrawlerConfigProxies    = "proxies"
	CrawlerConfigRetries    = "retries"
	CrawlerConfigBackoff    = "backoff"
)

// Filtering Configuration
const (
	CrawlerConfigIncludePatterns = "include_patterns"
	CrawlerConfigExcludePatterns = "exclude_patterns"
	CrawlerConfigAllowedDomains  = "allowed_domains"
	CrawlerConfigBlockedDomains  = "blocked_domains"
)

// ========================================
// CRAWLER STATES
// ========================================

// Crawler Execution States
const (
	CrawlerStateIdle        = "idle"
	CrawlerStateRunning     = "running"
	CrawlerStatePaused      = "paused"
	CrawlerStateStopped     = "stopped"
	CrawlerStateCompleted   = "completed"
	CrawlerStateFailed      = "failed"
	CrawlerStateTimeout     = "timeout"
	CrawlerStateBlocked     = "blocked"
	CrawlerStateRateLimited = "rate_limited"
)

// ========================================
// CRAWLER LINK TYPES
// ========================================

// Link Classification
const (
	CrawlerLinkTypeInternal   = "internal"
	CrawlerLinkTypeExternal   = "external"
	CrawlerLinkTypeAbsolute   = "absolute"
	CrawlerLinkTypeRelative   = "relative"
	CrawlerLinkTypeAnchor     = "anchor"
	CrawlerLinkTypeMailto     = "mailto"
	CrawlerLinkTypeTel        = "tel"
	CrawlerLinkTypeFTP        = "ftp"
	CrawlerLinkTypeJavaScript = "javascript"
)

// ========================================
// CRAWLER RESPONSE HEADERS
// ========================================

// Standard HTTP Headers
const (
	CrawlerResponseHeaderContentType   = "Content-Type"
	CrawlerResponseHeaderContentLength = "Content-Length"
	CrawlerResponseHeaderLastModified  = "Last-Modified"
	CrawlerResponseHeaderETag          = "ETag"
	CrawlerResponseHeaderCacheControl  = "Cache-Control"
	CrawlerResponseHeaderLocation      = "Location"
	CrawlerResponseHeaderServer        = "Server"
)

// Cookie Headers
const (
	CrawlerResponseHeaderSetCookie = "Set-Cookie"
	CrawlerResponseHeaderCookie    = "Cookie"
)

// Custom Headers
const (
	CrawlerResponseHeaderXRobotsTag     = "X-Robots-Tag"
	CrawlerResponseHeaderXFrameOptions  = "X-Frame-Options"
	CrawlerResponseHeaderXPoweredBy     = "X-Powered-By"
	CrawlerResponseHeaderXRequestID     = "X-Request-ID"
)

// ========================================
// CRAWLER TEST DATA
// ========================================

// Mock Data Prefixes
const (
	CrawlerTestDataMock    = "mock_"
	CrawlerTestDataTest    = "test_"
	CrawlerTestDataSample  = "sample_"
	CrawlerTestDataDummy   = "dummy_"
	CrawlerTestDataFake    = "fake_"
)

// Test Data Types
const (
	CrawlerTestDataHTML     = "test_html"
	CrawlerTestDataXML      = "test_xml"
	CrawlerTestDataJSON     = "test_json"
	CrawlerTestDataCSS      = "test_css"
	CrawlerTestDataJS       = "test_js"
	CrawlerTestDataRobots   = "test_robots"
	CrawlerTestDataSitemap  = "test_sitemap"
)

// ========================================
// CRAWLER TEST SCENARIOS
// ========================================

// Test Scenario Names
const (
	CrawlerTestScenarioBasic         = "test_basic_crawl"
	CrawlerTestScenarioDeepCrawl     = "test_deep_crawl"
	CrawlerTestScenarioRedirects     = "test_redirects"
	CrawlerTestScenarioRobots        = "test_robots_txt"
	CrawlerTestScenarioSitemap       = "test_sitemap_xml"
	CrawlerTestScenarioRateLimit     = "test_rate_limit"
	CrawlerTestScenarioTimeout       = "test_timeout"
	CrawlerTestScenarioErrorHandling = "test_error_handling"
)

// BDD Test Scenarios
const (
	CrawlerTestScenarioShouldCrawl    = "should_crawl_successfully"
	CrawlerTestScenarioShouldFollow   = "should_follow_links"
	CrawlerTestScenarioShouldRespect  = "should_respect_robots"
	CrawlerTestScenarioShouldHandle   = "should_handle_errors"
	CrawlerTestScenarioWhenBlocked    = "when_blocked_by_robots"
	CrawlerTestScenarioWhenTimeout    = "when_request_timeout"
	CrawlerTestScenarioGivenURL       = "given_valid_url"
)

// ========================================
// CRAWLER TEST ASSERTIONS
// ========================================

// Test Assertion Keywords
const (
	CrawlerTestAssertAssert   = "assert"
	CrawlerTestAssertExpect   = "expect"
	CrawlerTestAssertShould   = "should"
	CrawlerTestAssertMust     = "must"
	CrawlerTestAssertVerify   = "verify"
	CrawlerTestAssertValidate = "validate"
	CrawlerTestAssertCheck    = "check"
	CrawlerTestAssertEnsure   = "ensure"
	CrawlerTestAssertContain  = "contain"
	CrawlerTestAssertEqual    = "equal"
	CrawlerTestAssertMatch    = "match"
	CrawlerTestAssertInclude  = "include"
)

// ========================================
// CRAWLER FILE EXTENSIONS
// ========================================

// Web File Extensions
const (
	CrawlerFileExtensionHTML = ".html"
	CrawlerFileExtensionHTM  = ".htm"
	CrawlerFileExtensionPHP  = ".php"
	CrawlerFileExtensionASP  = ".asp"
	CrawlerFileExtensionJSP  = ".jsp"
	CrawlerFileExtensionXML  = ".xml"
	CrawlerFileExtensionTXT  = ".txt"
)

// Style and Script Extensions
const (
	CrawlerFileExtensionCSS  = ".css"
	CrawlerFileExtensionJS   = ".js"
	CrawlerFileExtensionJSON = ".json"
)

// Document Extensions
const (
	CrawlerFileExtensionPDF  = ".pdf"
	CrawlerFileExtensionDOC  = ".doc"
	CrawlerFileExtensionDOCX = ".docx"
	CrawlerFileExtensionXLS  = ".xls"
	CrawlerFileExtensionXLSX = ".xlsx"
	CrawlerFileExtensionPPT  = ".ppt"
	CrawlerFileExtensionPPTX = ".pptx"
)

// Archive Extensions
const (
	CrawlerFileExtensionZIP = ".zip"
	CrawlerFileExtensionTAR = ".tar"
	CrawlerFileExtensionGZ  = ".gz"
	CrawlerFileExtensionRAR = ".rar"
)

// ========================================
// CRAWLER ERROR MESSAGES
// ========================================

// Network Errors
const (
	CrawlerErrorFetchFailed       = "Échec de la récupération"
	CrawlerErrorTimeout           = "Timeout de la requête"
	CrawlerErrorConnectionRefused = "Connexion refusée"
	CrawlerErrorDNSLookupFailed   = "Échec de la résolution DNS"
	CrawlerErrorSSLError          = "Erreur SSL/TLS"
)

// Content Errors
const (
	CrawlerErrorInvalidURL       = "URL invalide"
	CrawlerErrorInvalidHTML      = "HTML invalide"
	CrawlerErrorParsingFailed    = "Échec du parsing"
	CrawlerErrorEncodingError    = "Erreur d'encodage"
	CrawlerErrorContentTooLarge  = "Contenu trop volumineux"
)

// Rate Limiting Errors
const (
	CrawlerErrorRateLimited     = "Limitation de débit appliquée"
	CrawlerErrorTooManyRequests = "Trop de requêtes"
	CrawlerErrorBlocked         = "Accès bloqué par robots.txt"
	CrawlerErrorForbidden       = "Accès interdit"
)

// ========================================
// CRAWLER LOG MESSAGES
// ========================================

// Crawling Activity Logs
const (
	CrawlerLogCrawlingStarted   = "Crawling démarré: %s"
	CrawlerLogFetchingURL       = "Récupération de l'URL: %s"
	CrawlerLogParsingHTML       = "Analyse HTML: %s"
	CrawlerLogExtractingLinks   = "Extraction des liens: %s"
	CrawlerLogFollowingLink     = "Suivi du lien: %s"
	CrawlerLogVisitingPage      = "Visite de la page: %s"
	CrawlerLogProcessingPage    = "Traitement de la page: %s"
	CrawlerLogAnalyzingContent  = "Analyse du contenu: %s"
	CrawlerLogCrawlingCompleted = "Crawling terminé: %s"
)

// Rate Limiting Logs
const (
	CrawlerLogRateLimitApplied = "Limitation appliquée: %s"
	CrawlerLogDelayApplied     = "Délai appliqué: %v"
	CrawlerLogWaitingForRate   = "Attente pour limitation: %v"
)

// ========================================
// CRAWLER DEFAULT VALUES
// ========================================

// Default Configuration Values
const (
	CrawlerDefaultTimeout       = 30    // seconds
	CrawlerDefaultDelay         = 1     // second
	CrawlerDefaultConcurrency   = 10    // concurrent requests
	CrawlerDefaultMaxDepth      = 3     // levels
	CrawlerDefaultMaxPages      = 1000  // pages
	CrawlerDefaultRetries       = 3     // retry attempts
	CrawlerDefaultBackoff       = 2     // exponential backoff factor
)

// Default Limits
const (
	CrawlerDefaultMaxFileSize   = 10485760 // 10MB
	CrawlerDefaultMaxRedirects  = 10       // redirect hops
	CrawlerDefaultMaxLinks      = 10000    // links per page
	CrawlerDefaultBufferSize    = 4096     // bytes
)

// Default Timeouts (in seconds)
const (
	CrawlerDefaultConnectTimeout = 10
	CrawlerDefaultReadTimeout    = 30
	CrawlerDefaultWriteTimeout   = 10
	CrawlerDefaultIdleTimeout    = 90
)