package constants

// ========================================
// DELTA-1 SERVER CONSTANTS
// Web Server Configuration and HTTP Constants
// ========================================

// ========================================
// HTTP ENDPOINTS
// ========================================

// API Endpoints
const (
	ServerEndpointAPI       = "/api/"
	ServerEndpointHealth    = "/health"
	ServerEndpointMetrics   = "/metrics"
	ServerEndpointStatus    = "/status"
	ServerEndpointDebug     = "/debug"
	ServerEndpointAnalyze   = "/analyze"
	ServerEndpointResults   = "/results"
	ServerEndpointReports   = "/reports"
)

// Static Endpoints
const (
	ServerEndpointStatic    = "/static/"
	ServerEndpointAssets    = "/assets/"
	ServerEndpointUploads   = "/uploads/"
	ServerEndpointDownloads = "/downloads/"
)

// ========================================
// HTTP METHODS
// ========================================

// Standard HTTP Methods
const (
	ServerMethodGet     = "GET"
	ServerMethodPost    = "POST"
	ServerMethodPut     = "PUT"
	ServerMethodDelete  = "DELETE"
	ServerMethodPatch   = "PATCH"
	ServerMethodHead    = "HEAD"
	ServerMethodOptions = "OPTIONS"
)

// ========================================
// CONTENT TYPES
// ========================================

// MIME Types
const (
	ServerContentTypeJSON      = "application/json"
	ServerContentTypeHTML      = "text/html"
	ServerContentTypePlain     = "text/plain"
	ServerContentTypeXML       = "application/xml"
	ServerContentTypeFormData  = "multipart/form-data"
	ServerContentTypeFormURL   = "application/x-www-form-urlencoded"
	ServerContentTypeOctet     = "application/octet-stream"
)

// Charset Specifications
const (
	ServerContentTypeJSONUTF8 = "application/json; charset=utf-8"
	ServerContentTypeHTMLUTF8 = "text/html; charset=utf-8"
	ServerContentTypePlainUTF8 = "text/plain; charset=utf-8"
)

// ========================================
// HTTP HEADERS
// ========================================

// Standard Headers
const (
	ServerHeaderContentType   = "Content-Type"
	ServerHeaderAuthorization = "Authorization"
	ServerHeaderAccept        = "Accept"
	ServerHeaderUserAgent     = "User-Agent"
	ServerHeaderCacheControl  = "Cache-Control"
	ServerHeaderContentLength = "Content-Length"
	ServerHeaderLocation      = "Location"
	ServerHeaderReferer       = "Referer"
	ServerHeaderOrigin        = "Origin"
)

// Custom Headers
const (
	ServerHeaderXRequestID     = "X-Request-ID"
	ServerHeaderXForwardedFor  = "X-Forwarded-For"
	ServerHeaderXRealIP        = "X-Real-IP"
	ServerHeaderXCSRFToken     = "X-CSRF-Token"
	ServerHeaderXAPIKey        = "X-API-Key"
)

// CORS Headers
const (
	ServerHeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	ServerHeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	ServerHeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	ServerHeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	ServerHeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	ServerHeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// ========================================
// HTTP STATUS CODES
// ========================================

// Success Codes
const (
	ServerStatusOK           = 200
	ServerStatusCreated      = 201
	ServerStatusAccepted     = 202
	ServerStatusNoContent    = 204
)

// Redirection Codes
const (
	ServerStatusMovedPermanently = 301
	ServerStatusFound           = 302
	ServerStatusNotModified     = 304
)

// Client Error Codes
const (
	ServerStatusBadRequest          = 400
	ServerStatusUnauthorized        = 401
	ServerStatusForbidden           = 403
	ServerStatusNotFound            = 404
	ServerStatusMethodNotAllowed    = 405
	ServerStatusNotAcceptable       = 406
	ServerStatusConflict            = 409
	ServerStatusUnprocessableEntity = 422
	ServerStatusTooManyRequests     = 429
)

// Server Error Codes
const (
	ServerStatusInternalServerError = 500
	ServerStatusNotImplemented      = 501
	ServerStatusBadGateway          = 502
	ServerStatusServiceUnavailable  = 503
	ServerStatusGatewayTimeout      = 504
)

// ========================================
// SERVER CONFIGURATION KEYS
// ========================================

// Basic Configuration
const (
	ServerConfigHost         = "host"
	ServerConfigPort         = "port"
	ServerConfigTimeout      = "timeout"
	ServerConfigReadTimeout  = "read_timeout"
	ServerConfigWriteTimeout = "write_timeout"
	ServerConfigIdleTimeout  = "idle_timeout"
)

// Advanced Configuration
const (
	ServerConfigMaxConnections    = "max_connections"
	ServerConfigBufferSize        = "buffer_size"
	ServerConfigMaxRequestSize    = "max_request_size"
	ServerConfigMaxResponseSize   = "max_response_size"
	ServerConfigKeepAlive         = "keep_alive"
	ServerConfigCompression       = "compression"
)

// SSL/TLS Configuration
const (
	ServerConfigTLSEnabled    = "tls_enabled"
	ServerConfigTLSCertFile   = "tls_cert_file"
	ServerConfigTLSKeyFile    = "tls_key_file"
	ServerConfigTLSMinVersion = "tls_min_version"
	ServerConfigTLSMaxVersion = "tls_max_version"
)

// ========================================
// FILE EXTENSIONS AND TEMPLATES
// ========================================

// Template Extensions
const (
	ServerExtensionHTML = ".html"
	ServerExtensionJSON = ".json"
	ServerExtensionXML  = ".xml"
	ServerExtensionTXT  = ".txt"
	ServerExtensionCSS  = ".css"
	ServerExtensionJS   = ".js"
)

// Static File Extensions
const (
	ServerExtensionPNG  = ".png"
	ServerExtensionJPG  = ".jpg"
	ServerExtensionJPEG = ".jpeg"
	ServerExtensionGIF  = ".gif"
	ServerExtensionSVG  = ".svg"
	ServerExtensionICO  = ".ico"
)

// Template Names
const (
	ServerTemplateIndex   = "index.html"
	ServerTemplateError   = "error.html"
	ServerTemplate404     = "404.html"
	ServerTemplate500     = "500.html"
	ServerTemplateHealth  = "health.html"
	ServerTemplateMetrics = "metrics.html"
)

// ========================================
// JSON FIELD NAMES (for API responses)
// ========================================

// Standard Response Fields
const (
	ServerJSONFieldStatus    = "status"
	ServerJSONFieldMessage   = "message"
	ServerJSONFieldData      = "data"
	ServerJSONFieldError     = "error"
	ServerJSONFieldTimestamp = "timestamp"
	ServerJSONFieldID        = "id"
	ServerJSONFieldName      = "name"
	ServerJSONFieldType      = "type"
	ServerJSONFieldURL       = "url"
	ServerJSONFieldMethod    = "method"
	ServerJSONFieldPath      = "path"
)

// Request/Response Metadata
const (
	ServerJSONFieldRequestID   = "request_id"
	ServerJSONFieldUserAgent   = "user_agent"
	ServerJSONFieldRemoteAddr  = "remote_addr"
	ServerJSONFieldContentType = "content_type"
	ServerJSONFieldSize        = "size"
	ServerJSONFieldDuration    = "duration"
)

// ========================================
// MIDDLEWARE NAMES
// ========================================

// Core Middleware
const (
	ServerMiddlewareCORS        = "cors"
	ServerMiddlewareAuth        = "auth"
	ServerMiddlewareLogging     = "logging"
	ServerMiddlewareRecovery    = "recovery"
	ServerMiddlewareCompression = "compression"
	ServerMiddlewareRateLimit   = "rate-limit"
)

// Security Middleware
const (
	ServerMiddlewareCSRF     = "csrf"
	ServerMiddlewareSecure   = "secure"
	ServerMiddlewareJWT      = "jwt"
	ServerMiddlewareAPIKey   = "api-key"
	ServerMiddlewareThrottle = "throttle"
)

// ========================================
// ERROR MESSAGES
// ========================================

// Generic Errors
const (
	ServerErrorInternal       = "Internal server error"
	ServerErrorNotFound       = "Resource not found"
	ServerErrorBadRequest     = "Bad request"
	ServerErrorUnauthorized   = "Unauthorized access"
	ServerErrorForbidden      = "Access forbidden"
	ServerErrorMethodNotAllowed = "Method not allowed"
	ServerErrorTimeout        = "Request timeout"
)

// Validation Errors
const (
	ServerErrorInvalidJSON    = "Invalid JSON format"
	ServerErrorMissingField   = "Required field missing"
	ServerErrorInvalidFormat  = "Invalid data format"
	ServerErrorOutOfRange     = "Value out of range"
	ServerErrorTooLarge       = "Request too large"
)

// ========================================
// LOG MESSAGES
// ========================================

// Server Lifecycle
const (
	ServerLogStarting      = "Server starting on %s:%d"
	ServerLogStarted       = "Server started successfully"
	ServerLogStopping      = "Server stopping..."
	ServerLogStopped       = "Server stopped"
	ServerLogListening     = "Server listening on %s"
	ServerLogServing       = "Serving request %s %s"
	ServerLogProcessing    = "Processing request %s"
	ServerLogShuttingDown  = "Server shutting down gracefully"
)

// Request Processing
const (
	ServerLogRequestReceived = "Request received: %s %s from %s"
	ServerLogRequestComplete = "Request completed: %s %s (%d) in %v"
	ServerLogRequestFailed   = "Request failed: %s %s - %s"
	ServerLogMiddleware      = "Middleware %s executed"
	ServerLogHandlerCalled   = "Handler %s called"
)

// ========================================
// ENVIRONMENT VARIABLES
// ========================================

// Server Configuration Environment Variables
const (
	ServerEnvHost         = "SERVER_HOST"
	ServerEnvPort         = "SERVER_PORT"
	ServerEnvTimeout      = "SERVER_TIMEOUT"
	ServerEnvDebug        = "SERVER_DEBUG"
	ServerEnvEnvironment  = "SERVER_ENVIRONMENT"
	ServerEnvLogLevel     = "SERVER_LOG_LEVEL"
)

// Security Environment Variables
const (
	ServerEnvSecretKey    = "SERVER_SECRET_KEY"
	ServerEnvJWTSecret    = "JWT_SECRET"
	ServerEnvAPIKey       = "API_KEY"
	ServerEnvCSRFKey      = "CSRF_KEY"
	ServerEnvTLSCert      = "TLS_CERT_FILE"
	ServerEnvTLSKey       = "TLS_KEY_FILE"
)

// ========================================
// DEFAULT VALUES
// ========================================

// Server Defaults
const (
	ServerDefaultHost    = "localhost"
	ServerDefaultPort    = 8080
	ServerDefaultTimeout = 30 // seconds
)

// Buffer and Size Defaults
const (
	ServerDefaultBufferSize     = 4096
	ServerDefaultMaxRequestSize = 1024 * 1024 // 1MB
	ServerDefaultMaxConnections = 1000
)