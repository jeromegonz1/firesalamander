package constants

// ========================================
// FIRE SALAMANDER - LOGGING CONSTANTS
// TDD + Zero Hardcoding Policy
// ========================================

// ========================================
// LOG LEVELS
// ========================================

// Log Levels suivant les standards syslog
const (
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarn    = "WARN"
	LogLevelError   = "ERROR"
	LogLevelFatal   = "FATAL"
	LogLevelTrace   = "TRACE"
)

// Log Level Numbers (pour comparaison)
const (
	LogLevelDebugNum   = 0
	LogLevelInfoNum    = 1
	LogLevelWarnNum    = 2
	LogLevelErrorNum   = 3
	LogLevelFatalNum   = 4
	LogLevelTraceNum   = -1
)

// ========================================
// LOG CATEGORIES
// ========================================

// Catégories de logs Fire Salamander
const (
	LogCategoryHTTP        = "HTTP"
	LogCategoryAPI         = "API"
	LogCategoryCrawler     = "CRAWLER"
	LogCategorySEO         = "SEO"
	LogCategoryOrchestrator = "ORCHESTRATOR"
	LogCategoryDatabase    = "DATABASE"
	LogCategoryAuth        = "AUTH"
	LogCategoryPerformance = "PERFORMANCE"
	LogCategorySystem      = "SYSTEM"
	LogCategoryDebug       = "DEBUG"
	LogCategoryAudit       = "AUDIT"
)

// ========================================
// LOG FORMATS
// ========================================

// Format de timestamp standard
const (
	LogTimestampFormat     = "2006-01-02T15:04:05.000Z"
	LogTimestampFormatFile = "2006-01-02"
)

// Templates de log JSON
const (
	LogFormatJSON = `{"timestamp":"%s","level":"%s","category":"%s","message":"%s","data":%s,"trace_id":"%s","request_id":"%s"}`
	LogFormatText = "%s [%s] [%s] %s - %s"
)

// ========================================
// LOG FILES
// ========================================

// Noms de fichiers de log
const (
	LogFileAccess      = "access.log"
	LogFileError       = "error.log"
	LogFileDebug       = "debug.log"
	LogFileAudit       = "audit.log"
	LogFilePerformance = "performance.log"
	LogFileSystem      = "system.log"
)

// Répertoires de logs
const (
	LogDefaultDir     = "./logs"
	LogDirPermissions = 0755
	LogFilePermissions = 0644
)

// ========================================
// HTTP ACCESS LOG FIELDS
// ========================================

// Champs obligatoires pour access logs
const (
	HTTPLogFieldMethod        = "method"
	HTTPLogFieldURL           = "url"
	HTTPLogFieldStatusCode    = "status_code"
	HTTPLogFieldResponseTime  = "response_time_ms"
	HTTPLogFieldUserAgent     = "user_agent"
	HTTPLogFieldRemoteAddr    = "remote_addr"
	HTTPLogFieldReferer       = "referer"
	HTTPLogFieldContentLength = "content_length"
	HTTPLogFieldRequestID     = "request_id"
)

// ========================================
// ERROR LOG FIELDS
// ========================================

// Champs pour error logs
const (
	ErrorLogFieldErrorType    = "error_type"
	ErrorLogFieldErrorCode    = "error_code"
	ErrorLogFieldStackTrace   = "stack_trace"
	ErrorLogFieldFunction     = "function"
	ErrorLogFieldFile         = "file"
	ErrorLogFieldLine         = "line"
	ErrorLogFieldComponent    = "component"
)

// ========================================
// PERFORMANCE LOG FIELDS
// ========================================

// Métriques de performance
const (
	PerfLogFieldOperation      = "operation"
	PerfLogFieldDurationMs     = "duration_ms"
	PerfLogFieldMemoryUsage    = "memory_usage_bytes"
	PerfLogFieldGoroutines     = "goroutines_count"
	PerfLogFieldCPUPercent     = "cpu_percent"
	PerfLogFieldRequestsPerSec = "requests_per_second"
)

// ========================================
// AUDIT LOG FIELDS
// ========================================

// Champs pour audit logs (actions critiques)
const (
	AuditLogFieldAction     = "action"
	AuditLogFieldResource   = "resource"
	AuditLogFieldUserID     = "user_id"
	AuditLogFieldSessionID  = "session_id"
	AuditLogFieldBefore     = "before_state"
	AuditLogFieldAfter      = "after_state"
	AuditLogFieldSuccess    = "success"
)

// ========================================
// LOG ROTATION
// ========================================

// Configuration rotation des logs
const (
	LogRotationMaxSizeMB   = 100
	LogRotationMaxBackups  = 10
	LogRotationMaxAgeDays  = 30
	LogRotationCompress    = true
)

// ========================================
// LOG MESSAGES STANDARDISÉS
// ========================================

// Messages de démarrage/arrêt
const (
	LogMsgServerStarting    = "Fire Salamander server starting"
	LogMsgServerStarted     = "Fire Salamander server started successfully"
	LogMsgServerStopping    = "Fire Salamander server stopping"
	LogMsgServerStopped     = "Fire Salamander server stopped"
)

// Messages HTTP
const (
	LogMsgHTTPRequest       = "HTTP request received"
	LogMsgHTTPResponse      = "HTTP response sent"
	LogMsgHTTPError         = "HTTP request error"
)

// Messages API
const (
	LogMsgAPIRequest        = "API request received"
	LogMsgAPIResponse       = "API response sent"
	LogMsgAPIError          = "API request failed"
	LogMsgAPIValidation     = "API validation error"
)

// Messages Crawler
const (
	LogMsgCrawlerStarted    = "Crawler started"
	LogMsgCrawlerStopped    = "Crawler stopped"
	LogMsgCrawlerPageFound  = "Page found by crawler"
	LogMsgCrawlerPageError  = "Crawler page error"
)

// Messages SEO Analysis
const (
	LogMsgSEOAnalysisStarted   = "SEO analysis started"
	LogMsgSEOAnalysisCompleted = "SEO analysis completed"
	LogMsgSEOAnalysisError     = "SEO analysis error"
)

// ========================================
// CONFIGURATION ENVIRONMENT VARIABLES
// ========================================

// Variables d'environnement pour logging
const (
	EnvLogLevel          = "LOG_LEVEL"
	EnvLogFormat         = "LOG_FORMAT"
	EnvLogDir            = "LOG_DIR"
	EnvLogRotationSize   = "LOG_ROTATION_SIZE_MB"
	EnvLogRotationBackups = "LOG_ROTATION_BACKUPS"
	EnvLogRotationAge    = "LOG_ROTATION_AGE_DAYS"
)

// Valeurs par défaut
const (
	DefaultLogLevel        = LogLevelInfo
	DefaultLogFormat       = "json"
	DefaultEnableConsole   = true
	DefaultEnableFile      = true
)

// ========================================
// TRACE ET REQUEST IDS
// ========================================

// Headers pour tracing
const (
	HeaderXRequestID  = "X-Request-ID"
	HeaderXTraceID    = "X-Trace-ID"
	HeaderXSpanID     = "X-Span-ID"
)

// Préfixes pour génération d'IDs
const (
	RequestIDPrefix = "req-"
	TraceIDPrefix   = "trace-"
	SpanIDPrefix    = "span-"
)