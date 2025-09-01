package constants

// 🚨 FIRE SALAMANDER - TEST URL Constants
// Zero Hardcoding Policy - All test URLs and network constants

// ===== TEST DOMAINS =====
const (
	TestDomainExample         = "https://example.com"
	TestDomainExampleSecure   = "https://example.com"
	TestDomainTest            = "https://test.com"
	TestDomainLocal           = "http://localhost"
	TestDomainLocalhost       = "localhost"
	TestIP                    = "127.0.0.1"
	TestIPAddress             = "127.0.0.1"
)

// ===== TEST URLs =====
const (
	TestURLExample            = "https://example.com"
	TestURLExampleGuide       = "https://example.com/guide-seo"
	TestURLExampleTest        = "https://example.com/test"
	TestURLIntegrationTest    = "https://integration-test.com"
	TestURLStorageTest        = "https://test-storage.com"
	TestURLBroken             = "http://broken.com"
	TestURLExternal           = "https://external.com"
)

// ===== TEST PORTS =====
const (
	TestPortDefault           = 8080
	TestPortWeb1              = 8083
	TestPortWeb2              = 8084
	TestPortAlternate         = 3000
	TestPortHTTP              = 80
	TestPortHTTPS             = 443
	TestPortDev               = 8000
	TestPortTest              = 9999
)

// ===== TEST SERVER URLs =====
const (
	TestServerLocalhost8083   = "http://localhost:8083"
	TestServerLocalhost8084   = "http://localhost:8084"
	TestServerLocal           = "http://127.0.0.1"
)

// ===== QUERY PARAMETERS =====
const (
	TestQueryURLParam         = "?url=" + TestURLExample
	TestQueryAnalyzeParam     = "/analyze?url=" + TestURLExample
	TestQueryResultsParam     = "/results?url=" + TestURLExample
)

// ===== API ENDPOINT PATHS =====
const (
	TestEndpointHealthPath    = "/web/health"
	TestEndpointAnalyzePath   = "/analyze"
	TestEndpointResultsPath   = "/results" 
	TestEndpointStatusPath    = "/status"
)

// ===== CANONICAL LINKS =====
const (
	TestCanonicalExample      = "https://example.com/guide-seo"
	TestCanonicalTest         = "https://example.com/test"
)

// ===== REQUEST BODIES =====
const (
	TestRequestBodyExample    = `{"url":"` + TestURLExample + `"}`
	TestRequestBodyInvalid    = `{"url":"invalid-url"}`
	TestRequestBodyEmpty      = `{}`
	TestRequestBodyMissing    = `{"other":"value"}`
)