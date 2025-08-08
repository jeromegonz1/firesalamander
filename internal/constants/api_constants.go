package constants

// ========================================
// DELTA-2 API CONSTANTS
// Integration API Configuration and HTTP Constants
// ========================================

// ========================================
// API ENDPOINTS
// ========================================

// API V1 Endpoints
const (
	APIEndpointV1Analyze         = "/api/v1/analyze"
	APIEndpointV1AnalyzeSemantic = "/api/v1/analyze/semantic"
	APIEndpointV1AnalyzeSEO      = "/api/v1/analyze/seo"
	APIEndpointV1AnalyzeQuick    = "/api/v1/analyze/quick"
	APIEndpointV1Health          = "/api/v1/health"
	APIEndpointV1Stats           = "/api/v1/stats"
	APIEndpointV1Analyses        = "/api/v1/analyses"
	APIEndpointV1Analysis        = "/api/v1/analysis/"
	APIEndpointV1Info            = "/api/v1/info"
	APIEndpointV1Version         = "/api/v1/version"
	APIEndpointV1Reports         = "api/v1/reports"
	APIEndpointV1Results         = "api/v1/results"
	APIEndpointV1Metrics         = "api/v1/metrics"
	APIEndpointV1Status          = "api/v1/status"
	APIEndpointV1Debug           = "api/v1/debug"
)

// API Root Endpoints
const (
	APIEndpointRoot     = "/api/"
	APIEndpointHealthy  = "/health"
	APIEndpointMetrics  = "/metrics"
	APIEndpointStatus   = "/status"
	APIEndpointDebug    = "/debug"
	APIEndpointAnalyze  = "/analyze"
	APIEndpointResults  = "/results"
	APIEndpointReports  = "/reports"
)

// ========================================
// HTTP METHODS (API Specific)
// ========================================

// HTTP Methods for API
const (
	APIMethodGet     = "GET"
	APIMethodPost    = "POST"
	APIMethodPut     = "PUT"
	APIMethodDelete  = "DELETE"
	APIMethodPatch   = "PATCH"
	APIMethodHead    = "HEAD"
	APIMethodOptions = "OPTIONS"
)

// ========================================
// CONTENT TYPES (API Specific)
// ========================================

// API Content Types
const (
	APIContentTypeJSON      = "application/json"
	APIContentTypeHTML      = "text/html"
	APIContentTypePlain     = "text/plain"
	APIContentTypeXML       = "application/xml"
	APIContentTypeCSV       = "text/csv"
	APIContentTypePDF       = "application/pdf"
	APIContentTypeFormData  = "multipart/form-data"
	APIContentTypeFormURL   = "application/x-www-form-urlencoded"
)

// Charset Specifications for API
const (
	APIContentTypeJSONUTF8  = "application/json; charset=utf-8"
	APIContentTypeHTMLUTF8  = "text/html; charset=utf-8"
	APIContentTypePlainUTF8 = "text/plain; charset=utf-8"
)

// ========================================
// HTTP HEADERS (API Specific)
// ========================================

// Standard API Headers
const (
	APIHeaderContentType   = "Content-Type"
	APIHeaderAuthorization = "Authorization"
	APIHeaderAccept        = "Accept"
	APIHeaderUserAgent     = "User-Agent"
	APIHeaderCacheControl  = "Cache-Control"
	APIHeaderOrigin        = "Origin"
)

// Custom API Headers
const (
	APIHeaderXRequestID     = "X-Request-ID"
	APIHeaderXForwardedFor  = "X-Forwarded-For"
	APIHeaderXRealIP        = "X-Real-IP"
	APIHeaderXAPIKey        = "X-API-Key"
	APIHeaderXVersion       = "X-API-Version"
)

// CORS Headers for API
const (
	APIHeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	APIHeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	APIHeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	APIHeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	APIHeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	APIHeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// ========================================
// JSON FIELD NAMES (API Responses)
// ========================================

// Standard API Response Fields
const (
	APIJSONFieldID          = "id"
	APIJSONFieldName        = "name"
	APIJSONFieldType        = "type"
	APIJSONFieldStatus      = "status"
	APIJSONFieldData        = "data"
	APIJSONFieldMessage     = "message"
	APIJSONFieldError       = "error"
	APIJSONFieldTimestamp   = "timestamp"
	APIJSONFieldURL         = "url"
	APIJSONFieldMethod      = "method"
	APIJSONFieldPath        = "path"
)

// Analysis Response Fields
const (
	APIJSONFieldResults         = "results"
	APIJSONFieldRecommendations = "recommendations"
	APIJSONFieldMetrics         = "metrics"
	APIJSONFieldScore           = "score"
	APIJSONFieldTitle           = "title"
	APIJSONFieldDescription     = "description"
	APIJSONFieldContent         = "content"
	APIJSONFieldValue           = "value"
	APIJSONFieldCategory        = "category"
	APIJSONFieldPriority        = "priority"
	APIJSONFieldSeverity        = "severity"
)

// Agent Response Fields
const (
	APIJSONFieldPhase      = "phase"
	APIJSONFieldAgent      = "agent"
	APIJSONFieldTest       = "test"
	APIJSONFieldCheck      = "check"
	APIJSONFieldCrawler    = "crawler"
	APIJSONFieldSEO        = "seo"
	APIJSONFieldSemantic   = "semantic"
	APIJSONFieldAnalysis   = "analysis"
	APIJSONFieldReport     = "report"
)

// Request Metadata Fields
const (
	APIJSONFieldRequestID   = "request_id"
	APIJSONFieldUserAgent   = "user_agent"
	APIJSONFieldRemoteAddr  = "remote_addr"
	APIJSONFieldContentType = "content_type"
	APIJSONFieldSize        = "size"
	APIJSONFieldDuration    = "duration"
)

// ========================================
// API STATUS VALUES
// ========================================

// Analysis Status Values
const (
	APIStatusPending   = "pending"
	APIStatusRunning   = "running"
	APIStatusCompleted = "completed"
	APIStatusFailed    = "failed"
	APIStatusSuccess   = "success"
	APIStatusError     = "error"
	APIStatusWarning   = "warning"
	APIStatusInfo      = "info"
)

// Health Status Values
const (
	APIStatusHealthy   = "healthy"
	APIStatusDegraded  = "degraded"
	APIStatusCritical  = "critical"
	APIStatusUnknown   = "unknown"
)

// Agent Status Values
const (
	APIStatusActive    = "active"
	APIStatusInactive  = "inactive"
	APIStatusDisabled  = "disabled"
)

// ========================================
// API AGENT NAMES
// ========================================

// Analysis Agents
const (
	APIAgentCrawler       = "crawler"
	APIAgentSEO          = "seo"
	APIAgentSemantic     = "semantic"
	APIAgentPerformance  = "performance"
	APIAgentSecurity     = "security"
	APIAgentQA           = "qa"
	APIAgentDataIntegrity = "data_integrity"
	APIAgentFrontend     = "frontend"
	APIAgentPlaywright   = "playwright"
	APIAgentK6           = "k6"
)

// Analysis Types
const (
	APIAnalysisTypeTechnical     = "technical"
	APIAnalysisTypeContent       = "content"
	APIAnalysisTypePerformance   = "performance"
	APIAnalysisTypeSecurity      = "security"
	APIAnalysisTypeAccessibility = "accessibility"
	APIAnalysisTypeSEO           = "seo"
	APIAnalysisTypeSemantic      = "semantic"
	APIAnalysisTypeStructural    = "structural"
)

// ========================================
// API ERROR MESSAGES
// ========================================

// Request Validation Errors
const (
	APIErrorInvalidJSON    = "Requête JSON invalide"
	APIErrorMissingField   = "Champ requis manquant"
	APIErrorInvalidFormat  = "Format de données invalide"
	APIErrorInvalidURL     = "URL invalide"
	APIErrorURLRequired    = "URL requise pour l'analyse"
)

// Analysis Errors
const (
	APIErrorAnalysisFailed    = "L'analyse a échoué"
	APIErrorAgentNotFound     = "Agent d'analyse non trouvé"
	APIErrorAnalysisTimeout   = "Timeout de l'analyse"
	APIErrorAnalysisConflict  = "Conflit lors de l'analyse"
	APIErrorInsufficientData  = "Données insuffisantes pour l'analyse"
)

// Server Errors
const (
	APIErrorInternalServer = "Erreur interne du serveur"
	APIErrorServiceUnavailable = "Service temporairement indisponible"
	APIErrorTooManyRequests    = "Trop de requêtes"
	APIErrorUnauthorized       = "Non autorisé"
	APIErrorForbidden          = "Accès interdit"
	APIErrorNotFound           = "Ressource non trouvée"
)

// ========================================
// API SUCCESS MESSAGES
// ========================================

// Analysis Success Messages
const (
	APISuccessAnalysisStarted   = "Analyse démarrée avec succès"
	APISuccessAnalysisCompleted = "Analyse terminée avec succès"
	APISuccessDataProcessed     = "Données traitées avec succès"
	APISuccessReportGenerated   = "Rapport généré avec succès"
)

// Generic Success Messages
const (
	APISuccessRequestProcessed = "Requête traitée avec succès"
	APISuccessDataRetrieved    = "Données récupérées avec succès"
	APISuccessOperationComplete = "Opération terminée avec succès"
)

// ========================================
// API LOG MESSAGES
// ========================================

// Request Processing Logs
const (
	APILogRequestReceived    = "Requête API reçue: %s %s"
	APILogRequestProcessing  = "Traitement de la requête: %s"
	APILogRequestCompleted   = "Requête terminée: %s %s (%d) en %v"
	APILogRequestFailed      = "Requête échouée: %s %s - %s"
)

// Analysis Processing Logs
const (
	APILogAnalysisStarted    = "Analyse démarrée: %s pour %s"
	APILogAnalysisProgress   = "Progression analyse: %s - %s"
	APILogAnalysisCompleted  = "Analyse terminée: %s en %v"
	APILogAnalysisFailed     = "Analyse échouée: %s - %s"
)

// Agent Processing Logs
const (
	APILogAgentStarted      = "Agent démarré: %s"
	APILogAgentProcessing   = "Agent en cours: %s - %s"
	APILogAgentCompleted    = "Agent terminé: %s en %v"
	APILogAgentFailed       = "Agent échoué: %s - %s"
)

// ========================================
// API CONFIGURATION KEYS
// ========================================

// API Server Configuration
const (
	APIConfigHost         = "api_host"
	APIConfigPort         = "api_port"
	APIConfigTimeout      = "api_timeout"
	APIConfigMaxRequests  = "api_max_requests"
	APIConfigRateLimit    = "api_rate_limit"
	APIConfigCORSEnabled  = "api_cors_enabled"
	APIConfigCORSOrigins  = "api_cors_origins"
)

// API Authentication Configuration
const (
	APIConfigAuthEnabled  = "api_auth_enabled"
	APIConfigAPIKey       = "api_key"
	APIConfigJWTSecret    = "jwt_secret"
	APIConfigTokenExpiry  = "token_expiry"
)

// Analysis Configuration
const (
	APIConfigAnalysisTimeout     = "analysis_timeout"
	APIConfigMaxConcurrentAnalysis = "max_concurrent_analysis"
	APIConfigAnalysisRetries     = "analysis_retries"
	APIConfigCacheEnabled        = "analysis_cache_enabled"
	APIConfigCacheTTL            = "analysis_cache_ttl"
)

// ========================================
// API MIME TYPES
// ========================================

// Document Types
const (
	APIMimeTypeJSON = "application/json"
	APIMimeTypeXML  = "application/xml"
	APIMimeTypeHTML = "text/html"
	APIMimeTypeTXT  = "text/plain"
	APIMimeTypeCSV  = "text/csv"
	APIMimeTypePDF  = "application/pdf"
)

// Image Types
const (
	APIMimeTypePNG  = "image/png"
	APIMimeTypeJPG  = "image/jpeg"
	APIMimeTypeGIF  = "image/gif"
	APIMimeTypeSVG  = "image/svg+xml"
	APIMimeTypeICO  = "image/x-icon"
	APIMimeTypeWEBP = "image/webp"
)

// ========================================
// API REPORT FORMATS
// ========================================

// Report Output Formats
const (
	APIReportFormatHTML     = "html"
	APIReportFormatJSON     = "json"
	APIReportFormatCSV      = "csv"
	APIReportFormatPDF      = "pdf"
	APIReportFormatXML      = "xml"
	APIReportFormatMarkdown = "markdown"
)

// ========================================
// API TEST CATEGORIES
// ========================================

// Testing Categories
const (
	APITestCategoryUnit          = "unit"
	APITestCategoryIntegration   = "integration"
	APITestCategoryE2E           = "e2e"
	APITestCategoryPerformance   = "performance"
	APITestCategorySecurity      = "security"
	APITestCategoryAccessibility = "accessibility"
	APITestCategorySEO           = "seo"
)

// ========================================
// API VERSION CONSTANTS
// ========================================

// API Version Information
const (
	APIVersionV1     = "v1"
	APIVersionHeader = "X-API-Version"
	APIVersionPath   = "/api/v1"
)

// ========================================
// API DEFAULT VALUES
// ========================================

// Default API Configuration
const (
	APIDefaultTimeout     = 30 // seconds
	APIDefaultMaxRequests = 1000
	APIDefaultRateLimit   = 100 // requests per minute
	APIDefaultPort        = 8080
	APIDefaultHost        = "localhost"
)

// Default Analysis Configuration
const (
	APIDefaultAnalysisTimeout     = 300 // seconds
	APIDefaultMaxConcurrentAnalysis = 5
	APIDefaultAnalysisRetries     = 3
	APIDefaultCacheTTL           = 3600 // seconds
)