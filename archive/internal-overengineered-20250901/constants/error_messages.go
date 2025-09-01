package constants

// 🚨 FIRE SALAMANDER - ERROR MESSAGE Constants
// Zero Hardcoding Policy - All error/warning/info messages

// ===== ERROR MESSAGES =====
const (
	// General Errors
	ErrorGeneric                = "Error"
	ErrorOccurred              = "An error occurred"
	ErrorUnexpected            = "Unexpected error"
	ErrorInternal              = "Internal error"
	ErrorUnknown               = "Unknown error"
	ErrorProcessing            = "Processing error"
	ErrorTimeout               = "Timeout error"
	ErrorCancelled             = "Operation cancelled"
	
	// Validation Errors
	ErrorInvalid               = "Invalid"
	ErrorInvalidInput          = "Invalid input"
	ErrorInvalidFormat         = "Invalid format"
	ErrorInvalidURL            = "Invalid URL"
	ErrorInvalidParameter      = "Invalid parameter"
	ErrorInvalidConfiguration  = "Invalid configuration"
	ErrorInvalidData           = "Invalid data"
	ErrorInvalidRequest        = "Invalid request"
	ErrorInvalidResponse       = "Invalid response"
	ErrorInvalidAnalysisID     = "Invalid analysis ID"
	
	// Required/Missing Errors
	ErrorRequired              = "Required"
	ErrorFieldRequired         = "Field required"
	ErrorParameterRequired     = "Parameter required"
	ErrorMissing               = "Missing"
	ErrorMissingField          = "Missing field"
	ErrorMissingParameter      = "Missing parameter"
	ErrorMissingConfiguration  = "Missing configuration"
	ErrorNotFound              = "Not found"
	ErrorFileNotFound          = "File not found"
	ErrorResourceNotFound      = "Resource not found"
	
	// Operation Errors
	ErrorFailed                = "Failed"
	ErrorOperationFailed       = "Operation failed"
	ErrorRequestFailed         = "Request failed"
	ErrorConnectionFailed      = "Connection failed"
	ErrorAuthenticationFailed  = "Authentication failed"
	ErrorAuthorizationFailed   = "Authorization failed"
	ErrorValidationFailed      = "Validation failed"
	ErrorParsingFailed         = "Parsing failed"
	ErrorDecodingFailed        = "Decoding failed"
	ErrorEncodingFailed        = "Encoding failed"
	
	// File/IO Errors
	ErrorReadingFile           = "Error reading file"
	ErrorWritingFile           = "Error writing file"
	ErrorCreatingFile          = "Error creating file"
	ErrorDeletingFile          = "Error deleting file"
	ErrorOpeningFile           = "Error opening file"
	ErrorClosingFile           = "Error closing file"
	
	// Network Errors
	ErrorNetwork               = "Network error"
	ErrorConnectionRefused     = "Connection refused"
	ErrorConnectionReset       = "Connection reset"
	ErrorConnectionTimeout     = "Connection timeout"
	ErrorNoConnection          = "No connection"
	ErrorDNSResolution         = "DNS resolution error"
	
	// Database Errors
	ErrorDatabase              = "Database error"
	ErrorDatabaseConnection    = "Database connection error"
	ErrorDatabaseQuery         = "Database query error"
	ErrorDatabaseTransaction   = "Database transaction error"
	
	// HTTP Errors
	ErrorHTTP                  = "HTTP error"
	ErrorHTTPRequest           = "HTTP request error"
	ErrorHTTPResponse          = "HTTP response error"
	ErrorHTTPTimeout           = "HTTP timeout"
)

// ===== WARNING MESSAGES =====
const (
	WarningGeneric             = "Warning"
	WarningDeprecated          = "Deprecated"
	WarningPerformance         = "Performance warning"
	WarningMemory              = "Memory warning"
	WarningCapacity            = "Capacity warning"
	WarningLimit               = "Limit warning"
	WarningThreshold           = "Threshold warning"
	WarningRetry               = "Retry warning"
)

// ===== INFO MESSAGES =====
const (
	InfoGeneric                = "Info"
	InfoStarting               = "Starting"
	InfoStarted                = "Started"
	InfoStopping               = "Stopping"
	InfoStopped                = "Stopped"
	InfoProcessing             = "Processing"
	InfoCompleted              = "Completed"
	InfoReady                  = "Ready"
	InfoInitializing           = "Initializing"
	InfoInitialized            = "Initialized"
	InfoConnecting             = "Connecting"
	InfoConnected              = "Connected"
	InfoDisconnecting          = "Disconnecting"
	InfoDisconnected           = "Disconnected"
	InfoLoading                = "Loading"
	InfoLoaded                 = "Loaded"
	InfoSaving                 = "Saving"
	InfoSaved                  = "Saved"
	InfoCreating               = "Creating"
	InfoCreated                = "Created"
	InfoUpdating               = "Updating"
	InfoUpdated                = "Updated"
	InfoDeleting               = "Deleting"
	InfoDeleted                = "Deleted"
	InfoFound                  = "Found"
	InfoNotFound               = "Not found"
	InfoEmpty                  = "Empty"
	InfoAvailable              = "Available"
	InfoUnavailable            = "Unavailable"
	InfoEnabled                = "Enabled"
	InfoDisabled               = "Disabled"
	InfoActive                 = "Active"
	InfoInactive               = "Inactive"
	InfoPending                = "Pending"
	InfoQueued                 = "Queued"
	InfoScheduled              = "Scheduled"
	InfoExecuting              = "Executing"
	InfoExecuted               = "Executed"
	InfoRetrying               = "Retrying"
	InfoSkipped                = "Skipped"
	InfoIgnored                = "Ignored"
	InfoCancelled              = "Cancelled"
	InfoTimedOut               = "Timed out"
)

// ===== SUCCESS MESSAGES =====
const (
	SuccessGeneric             = "Success"
	SuccessCompleted           = "Successfully completed"
	SuccessCreated             = "Successfully created"
	SuccessUpdated             = "Successfully updated"
	SuccessDeleted             = "Successfully deleted"
	SuccessSaved               = "Successfully saved"
	SuccessLoaded              = "Successfully loaded"
	SuccessConnected           = "Successfully connected"
	SuccessAuthenticated       = "Successfully authenticated"
	SuccessAuthorized          = "Successfully authorized"
	SuccessValidated           = "Successfully validated"
	SuccessProcessed           = "Successfully processed"
	SuccessExecuted            = "Successfully executed"
	SuccessDeployed            = "Successfully deployed"
	SuccessConfigured          = "Successfully configured"
	SuccessInitialized         = "Successfully initialized"
	SuccessRegistered          = "Successfully registered"
	SuccessUnregistered        = "Successfully unregistered"
	SuccessSubscribed          = "Successfully subscribed"
	SuccessUnsubscribed        = "Successfully unsubscribed"
)

// ===== CRAWLER SPECIFIC MESSAGES =====
const (
	CrawlerStarting            = "Crawler starting"
	CrawlerStarted             = "Crawler started"
	CrawlerStopping            = "Crawler stopping"
	CrawlerStopped             = "Crawler stopped"
	CrawlerFetching            = "Fetching page"
	CrawlerFetched             = "Page fetched"
	CrawlerParsing             = "Parsing content"
	CrawlerParsed              = "Content parsed"
	CrawlerError               = "Crawler error"
	CrawlerTimeout             = "Crawler timeout"
	CrawlerRateLimited         = "Rate limited"
	CrawlerRetrying            = "Retrying request"
	CrawlerBlocked             = "Request blocked"
	CrawlerAllowed             = "Request allowed"
	CrawlerQueued              = "Request queued"
	CrawlerCompleted           = "Crawl completed"
	CrawlerFailed              = "Crawl failed"
)

// ===== ANALYSIS SPECIFIC MESSAGES ===== (some already in messages.go)
const (
	AnalysisStarted            = "Analysis started"
	AnalysisProcessing         = "Analysis processing"
	AnalysisCompleted          = "Analysis completed"
	AnalysisError              = "Analysis error"
	AnalysisTimeout            = "Analysis timeout"
	AnalysisCancelled          = "Analysis cancelled"
	AnalysisQueued             = "Analysis queued"
	AnalysisPending            = "Analysis pending"
	AnalysisRunning            = "Analysis running"
	AnalysisFinished           = "Analysis finished"
	AnalysisNotFound           = "Analyse non trouvée"
)

// ===== VALIDATION MESSAGES =====
const (
	ValidationRequired         = "Validation required"
	ValidationInProgress       = "Validation in progress"
	ValidationPassed           = "Validation passed"
	ValidationFailed           = "Validation failed"
	ValidationError            = "Validation error"
	ValidationWarning          = "Validation warning"
	ValidationSkipped          = "Validation skipped"
)

// ===== CONFIGURATION MESSAGES =====
const (
	ConfigLoading              = "Loading configuration"
	ConfigLoaded               = "Configuration loaded"
	ConfigError                = "Configuration error"
	ConfigInvalid              = "Invalid configuration"
	ConfigMissing              = "Missing configuration"
	ConfigUpdated              = "Configuration updated"
	ConfigSaved                = "Configuration saved"
	ConfigApplied              = "Configuration applied"
)

// ===== SERVER MESSAGES ===== (basic ones already in messages.go)
const (
	ServerStarted              = "Server started"
	ServerError                = "Server error"
	ServerListening            = "Server listening"
	ServerShutdown             = "Server shutdown"
	ServerRestarting           = "Server restarting"
	ServerMaintenance          = "Server maintenance"
)

// ===== ORCHESTRATOR MESSAGES =====
const (
	ErrorCreatingOrchestrator  = "Failed to create orchestrator"
	ErrorStartingOrchestrator  = "Failed to start orchestrator"
	ErrorStoppingOrchestrator  = "Error stopping orchestrator"
	ErrorStartingWebServer     = "Failed to start web server"
	ErrorStoppingWebServer     = "Error stopping web server"
)

// ===== FIRE SALAMANDER SPECIFIC =====
const (
	FireSalamanderStarted      = "🔥 Fire Salamander started successfully"
	FireSalamanderStopped      = "🔥 Fire Salamander stopped gracefully"
	ShutdownSignalReceived     = "Shutdown signal received, stopping gracefully..."
)

// ===== JSON KEYS =====
const (
	JSONKeyError               = "error"
	JSONKeyMessage             = "message"
	JSONKeyData                = "data"
	JSONKeyStatus              = "status"
	JSONKeyResult              = "result"
	JSONKeyID                  = "id"
	JSONKeyURL                 = "url"
)

// ===== MESSAGE TEMPLATES =====
const (
	TemplateErrorWithReason    = "%s: %s"
	TemplateFieldRequired      = "%s is required"
	TemplateFieldInvalid       = "%s is invalid"
	TemplateFieldMissing       = "%s is missing"
	TemplateOperationFailed    = "%s operation failed"
	TemplateResourceNotFound   = "%s not found"
	TemplateConnectionFailed   = "Failed to connect to %s"
	TemplateTimeout            = "%s timeout after %s"
	TemplateRetrying           = "Retrying %s (attempt %d/%d)"
	TemplateProgress           = "%s: %d/%d"
	TemplatePercentComplete    = "%s: %.1f%% complete"
)