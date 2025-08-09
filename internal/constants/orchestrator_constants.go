package constants

// ========================================
// CHARLIE-2 ORCHESTRATOR CONSTANTS
// ========================================

// ========================================
// ORCHESTRATOR JSON FIELD NAMES
// ========================================

// Orchestrator Response Fields
const (
	OrchestratorJSONFieldStatus    = "status"
	OrchestratorJSONFieldMessage   = "message"
	OrchestratorJSONFieldTimestamp = "timestamp"
	OrchestratorJSONFieldPhase     = "phase"
	OrchestratorJSONFieldDuration  = "duration"
	OrchestratorJSONFieldError     = "error"
	OrchestratorJSONFieldResult    = "result"
	OrchestratorJSONFieldData      = "data"
	OrchestratorJSONFieldConfig    = "config"
	OrchestratorJSONFieldAgents    = "agents"
	OrchestratorJSONFieldTasks     = "tasks"
	OrchestratorJSONFieldLogs      = "logs"
	OrchestratorJSONFieldMetrics   = "metrics"
	OrchestratorJSONFieldProgress  = "progress"
	OrchestratorJSONFieldState     = "state"
)

// ========================================
// ORCHESTRATOR AGENT NAMES
// ========================================

// Agent Identifiers
const (
	OrchestratorAgentNameOrchestrator    = "orchestrator"
	OrchestratorAgentNameCrawler         = "crawler"
	OrchestratorAgentNameSEO             = "seo"
	OrchestratorAgentNameQA              = "qa"
	OrchestratorAgentNameSecurity        = "security"
	OrchestratorAgentNamePerformance     = "performance"
	OrchestratorAgentNameDataIntegrity   = "data-integrity"
	OrchestratorAgentNameTagAnalyzer     = "tag-analyzer"
)

// ========================================
// ORCHESTRATOR EXECUTION PHASES
// ========================================

// Phase Names
const (
	OrchestratorPhaseInitialization = "initialization"
	OrchestratorPhaseCrawling       = "crawling"
	OrchestratorPhaseAnalysis       = "analysis"
	OrchestratorPhaseScoring        = "scoring"
	OrchestratorPhaseReporting      = "reporting"
	OrchestratorPhaseCleanup        = "cleanup"
	OrchestratorPhaseValidation     = "validation"
	OrchestratorPhaseTesting        = "testing"
	OrchestratorPhaseSecurityScan   = "security-scan"
)

// ========================================
// ORCHESTRATOR STATUS VALUES
// ========================================

// Task/Agent Status Values
const (
	OrchestratorStatusPending   = "pending"
	OrchestratorStatusRunning   = "running"
	OrchestratorStatusCompleted = "completed"
	OrchestratorStatusFailed    = "failed"
	OrchestratorStatusCancelled = "cancelled"
	OrchestratorStatusSuccess   = "success"
	OrchestratorStatusError     = "error"
	OrchestratorStatusWarning   = "warning"
	OrchestratorStatusInfo      = "info"
)

// ========================================
// ORCHESTRATOR CONFIGURATION KEYS
// ========================================

// Configuration Parameters
const (
	OrchestratorConfigMaxWorkers       = "max_workers"
	OrchestratorConfigTimeout          = "timeout"
	OrchestratorConfigRetryAttempts    = "retry_attempts"
	OrchestratorConfigBatchSize        = "batch_size"
	OrchestratorConfigConcurrentAgents = "concurrent_agents"
	OrchestratorConfigLogLevel         = "log_level"
	OrchestratorConfigDebugMode        = "debug_mode"
)

// ========================================
// ORCHESTRATOR HTTP METHODS
// ========================================

// HTTP Methods for Agent Communication
const (
	OrchestratorHTTPMethodGet     = "GET"
	OrchestratorHTTPMethodPost    = "POST"
	OrchestratorHTTPMethodPut     = "PUT"
	OrchestratorHTTPMethodDelete  = "DELETE"
	OrchestratorHTTPMethodPatch   = "PATCH"
	OrchestratorHTTPMethodHead    = "HEAD"
	OrchestratorHTTPMethodOptions = "OPTIONS"
)

// ========================================
// ORCHESTRATOR CONTENT TYPES
// ========================================

// Content Types for HTTP Communication
const (
	OrchestratorContentTypeJSON  = "application/json"
	OrchestratorContentTypeHTML  = "text/html"
	OrchestratorContentTypePlain = "text/plain"
	OrchestratorContentTypeXML   = "application/xml"
)

// ========================================
// ORCHESTRATOR NUMERIC THRESHOLDS
// ========================================

// Concurrency and Performance Limits
const (
	OrchestratorMaxWorkers             = 10
	OrchestratorDefaultWorkers         = 5
	OrchestratorMaxConcurrentAgents    = 3
	OrchestratorMaxRetryAttempts       = 3
	OrchestratorDefaultBatchSize       = 100
	OrchestratorMaxBatchSize           = 1000
	OrchestratorMaxTimeoutSeconds      = 3600
	OrchestratorDefaultTimeoutSeconds  = 300
	OrchestratorMaxQueueSize           = 5000
	OrchestratorDefaultQueueSize       = 1000
	OrchestratorHealthCheckInterval    = 30  // seconds
	OrchestratorProgressUpdateInterval = 60  // seconds
	OrchestratorMaxLogRetention        = 100 // entries
)

// ========================================
// ORCHESTRATOR GOROUTINE PATTERNS
// ========================================

// Goroutine Management
const (
	OrchestratorGoRoutinePoolSize    = 10
	OrchestratorMaxActiveGoroutines  = 50
	OrchestratorGoroutineTimeout     = 30 // minutes
	OrchestratorGoroutineBufferSize  = 100
)

// ========================================
// ORCHESTRATOR CHANNEL OPERATIONS
// ========================================

// Channel Buffer Sizes
const (
	OrchestratorChannelBufferDefault = 100
	OrchestratorChannelBufferLarge   = 1000
	OrchestratorChannelBufferSmall   = 10
	OrchestratorTaskChannelSize      = 500
	OrchestratorResultChannelSize    = 200
	OrchestratorErrorChannelSize     = 50
)

// ========================================
// ORCHESTRATOR ERROR TYPES
// ========================================

// Error Categories
const (
	OrchestratorErrorTypeCrawling     = "crawling_error"
	OrchestratorErrorTypeSemantic     = "semantic_error" 
	OrchestratorErrorTypeSEO          = "seo_error"
	OrchestratorErrorTypeValidation   = "validation_error"
	OrchestratorErrorTypeTimeout      = "timeout_error"
	OrchestratorErrorTypeConfiguration = "config_error"
	OrchestratorErrorTypeNetwork      = "network_error"
	OrchestratorErrorTypeInternal     = "internal_error"
)

// ========================================
// ORCHESTRATOR MESSAGE TEMPLATES
// ========================================

// Standard Messages
const (
	OrchestratorMsgAgentStarted      = "Agent %s started successfully"
	OrchestratorMsgAgentCompleted    = "Agent %s completed with status: %s"
	OrchestratorMsgAgentFailed       = "Agent %s failed with error: %s"
	OrchestratorMsgPhaseStarted      = "Phase %s started"
	OrchestratorMsgPhaseCompleted    = "Phase %s completed in %v"
	OrchestratorMsgTaskQueued        = "Task %s queued for processing"
	OrchestratorMsgTaskProcessing    = "Processing task %s with agent %s"
	OrchestratorMsgProgressUpdate    = "Progress: %d/%d tasks completed (%d%%)"
	OrchestratorMsgSystemReady       = "Orchestrator system ready with %d workers"
	OrchestratorMsgSystemShutdown    = "Orchestrator system shutting down"
)

// Error Messages
const (
	OrchestratorErrAgentNotFound      = "Agent %s not found or not registered"
	OrchestratorErrAgentTimeout       = "Agent %s timed out after %v"
	OrchestratorErrMaxRetriesExceeded = "Maximum retry attempts (%d) exceeded for task %s"
	OrchestratorErrInvalidConfig      = "Invalid configuration: %s"
	OrchestratorErrWorkerPoolFull     = "Worker pool is at maximum capacity (%d)"
	OrchestratorErrTaskQueueFull      = "Task queue is at maximum capacity (%d)"
	OrchestratorErrInvalidPhase       = "Invalid phase transition: %s → %s"
	OrchestratorErrResourceExhausted  = "System resources exhausted: %s"
)

// ========================================
// ORCHESTRATOR LOG MESSAGES
// ========================================

// System Log Messages
const (
	OrchestratorLogStarting        = "Starting orchestrator system"
	OrchestratorLogStopping        = "Stopping orchestrator system"
	OrchestratorLogInitializing    = "Initializing orchestrator components"
	OrchestratorLogExecuting       = "Executing orchestration plan"
	OrchestratorLogCompleted       = "Orchestration completed successfully"
	OrchestratorLogProcessing      = "Processing %d tasks with %d agents"
	OrchestratorLogValidating      = "Validating orchestrator configuration"
	OrchestratorLogCleaning        = "Cleaning up orchestrator resources"
)

// Agent Communication Log Messages
const (
	OrchestratorLogAgentRegistered   = "Agent %s registered successfully"
	OrchestratorLogAgentDeregistered = "Agent %s deregistered"
	OrchestratorLogAgentHealthCheck  = "Health check for agent %s: %s"
	OrchestratorLogAgentCommunication = "Communication with agent %s: %s"
)

// ========================================
// ORCHESTRATOR FILE PATHS
// ========================================

// Default Paths
const (
	OrchestratorPathReports = "reports/"
	OrchestratorPathLogs    = "logs/"
	OrchestratorPathTemp    = "temp/"
	OrchestratorPathCache   = "cache/"
	OrchestratorPathConfig  = "config/"
)

// Report File Names
const (
	OrchestratorFileOrchestrationReport = "orchestration_report.json"
	OrchestratorFileAgentMetrics        = "agent_metrics.json"
	OrchestratorFileSystemStatus        = "system_status.json"
	OrchestratorFileErrorLog            = "orchestrator_errors.log"
	OrchestratorFilePerformanceLog      = "performance_metrics.log"
)

// ========================================
// ORCHESTRATOR TIME FORMATS
// ========================================

// Time Format Constants
const (
	OrchestratorTimeFormatRFC3339    = "RFC3339"
	OrchestratorTimeFormatISO8601    = "2006-01-02T15:04:05Z"
	OrchestratorTimeFormatSimple     = "2006-01-02 15:04:05"
	OrchestratorTimeFormatLogEntry   = "15:04:05.000"
)

// ========================================
// ORCHESTRATOR URL PATTERNS
// ========================================

// URL Schemes for Agent Communication
const (
	OrchestratorURLSchemeHTTP  = "http://"
	OrchestratorURLSchemeHTTPS = "https://"
	OrchestratorURLSchemeWS    = "ws://"
	OrchestratorURLSchemeWSS   = "wss://"
)

// Default Endpoints
const (
	OrchestratorEndpointHealth   = "/health"
	OrchestratorEndpointStatus   = "/status"
	OrchestratorEndpointMetrics  = "/metrics"
	OrchestratorEndpointAgents   = "/agents"
	OrchestratorEndpointTasks    = "/tasks"
)

// ========================================
// ORCHESTRATOR CONTEXT OPERATIONS
// ========================================

// Context Types for Different Operations
const (
	OrchestratorContextBackground = "Background"
	OrchestratorContextTODO       = "TODO"
	OrchestratorContextTimeout    = "WithTimeout"
	OrchestratorContextCancel     = "WithCancel"
)

// ========================================
// ORCHESTRATOR SYNCHRONIZATION
// ========================================

// Sync Primitive Names
const (
	OrchestratorSyncWaitGroup = "WaitGroup"
	OrchestratorSyncMutex     = "Mutex"
	OrchestratorSyncRWMutex   = "RWMutex"
)

// ========================================
// ORCHESTRATOR OPERATION MODES
// ========================================

// Execution Modes
const (
	OrchestratorModeSequential = "sequential"
	OrchestratorModeParallel   = "parallel"
	OrchestratorModePipeline   = "pipeline"
	OrchestratorModeAdaptive   = "adaptive"
)

// Quality Modes
const (
	OrchestratorQualityFast     = "fast"
	OrchestratorQualityBalanced = "balanced"
	OrchestratorQualityThorough = "thorough"
)
// Additional Orchestrator Constants for CHARLIE-2

// Analysis Types
const (
	OrchestratorAnalysisTypeFull     = "full"
	OrchestratorAnalysisTypeSemantic = "semantic"
	OrchestratorAnalysisTypeQuick    = "quick"
)

// Additional Status Values
const (
	OrchestratorStatusPartial = "partial"
)

// Insight Types
const (
	OrchestratorInsightContentSEOAlignment           = "content_seo_alignment"
	OrchestratorInsightPerformanceContentMismatch   = "performance_content_mismatch"
)

// Impact Types
const (
	OrchestratorImpactPositive = "positive"
	OrchestratorImpactNegative = "negative"
	OrchestratorImpactNeutral  = "neutral"
)

// Category Names for Scores
const (
	OrchestratorCategoryTechnical         = "technical"
	OrchestratorCategoryContentQuality    = "content_quality"
	OrchestratorCategoryUserExperience    = "user_experience"
	OrchestratorCategoryMobileFriendliness = "mobile_friendliness"
)

// Agent Names (additional)
const (
	OrchestratorAgentNameSemantic = "semantic"
)

// Time Constants
const (
	OrchestratorTimeVariable = "Variable"
)

// ========================================
// SPRINT 5 - REAL ORCHESTRATOR CONSTANTS
// ========================================

// Real Orchestrator Status Values
const (
	OrchestratorStatusStarting    = "starting"
	OrchestratorStatusCrawling    = "crawling" 
	OrchestratorStatusAnalyzing   = "analyzing"
	OrchestratorStatusAggregating = "aggregating"
	OrchestratorStatusComplete    = "complete"
)

// Real Orchestrator Configuration
const (
	OrchestratorMaxPages        = 20
	OrchestratorInitialWorkers  = 5
	OrchestratorAnalysisTimeout = 120000000000 // 2 minutes in nanoseconds
)

// Test Constants
const (
	MaxAnalysisWaitTime        = 30000000000  // 30 seconds in nanoseconds
	AnalysisPollingInterval    = 500000000    // 500ms in nanoseconds
	TestServerDelay            = 100000000    // 100ms in nanoseconds
	MaxUpdateChecks            = 20
	UpdateCheckInterval        = 250000000    // 250ms in nanoseconds
	ErrorAnalysisTimeout       = 10000000000  // 10 seconds in nanoseconds
	ErrorCheckInterval         = 500000000    // 500ms in nanoseconds
	MaxAcceptableAnalysisTime  = 120000000000 // 2 minutes in nanoseconds
)
