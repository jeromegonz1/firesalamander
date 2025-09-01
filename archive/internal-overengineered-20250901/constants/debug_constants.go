package constants

// Debug and Phase Test Constants
// Constantes pour les tests de phase et debugging

// Phase Test Names
const (
	TestConfigurationLoading = "Configuration Loading"
	TestFileStructure       = "File Structure"
	TestGitRepository       = "Git Repository"
	TestHTTPServer          = "HTTP Server"
	TestSEPTEOBranding      = "SEPTEO Branding"
	TestDockerSetup         = "Docker Setup"
	TestDeployScripts       = "Deploy Scripts"
)

// Phase Test Descriptions
const (
	DescConfigurationLoading = "Verify configuration is loaded correctly with all required fields"
	DescFileStructure       = "Verify all required files and directories are present"
	DescGitRepository       = "Verify git repository is properly initialized"
	DescHTTPServer          = "Test server endpoints are responding correctly"
	DescSEPTEOBranding      = "Verify SEPTEO branding is properly integrated"
	DescDockerSetup         = "Verify Docker Compose configuration exists and is valid"
	DescDeployScripts       = "Verify deployment scripts exist and are executable"
)

// Phase Names and Numbers
const (
	Phase1SetupInitial = "Phase 1 - Setup Initial"
	Phase1Number       = "1"
)

// Log Messages for Phase Tests
const (
	LogRunningPhase1Tests      = "🧪 Running Phase 1 Tests - Setup Initial"
	LogConfigurationTest       = "Running configuration test"
	LogFileStructureTest       = "Running file structure test"
	LogGitTest                 = "Running git setup test"
	LogHTTPServerTest          = "Running HTTP server test"
	LogSEPTEOBrandingTest      = "Running SEPTEO branding test"
	LogDockerTest              = "Running Docker setup test"
	LogDeployScriptsTest       = "Running deploy scripts test"
	LogPhase1Completed         = "🧪 Phase 1 Tests completed"
)

// Error Messages
const (
	ErrConfigurationNil        = "Configuration is nil"
	ErrConfigNil               = "config_nil"
	ErrAppNameIncorrect        = "app.name incorrect"
	ErrAppIconIncorrect        = "app.icon incorrect"
	ErrAppPoweredByIncorrect   = "app.powered_by incorrect"
	ErrBrandingColorIncorrect  = "branding.primary_color incorrect"
	ErrServerPortInvalid       = "server.port invalid"
	ErrDockerComposeMissing    = "docker_compose_missing"
	ErrDockerComposeInvalid    = "Docker Compose configuration is incomplete"
	ErrCannotReadDockerCompose = "Cannot read docker-compose.yml"
)

// Success Messages
const (
	MsgConfigurationValid         = "Configuration is valid and complete"
	MsgAllFilesExist             = "All required files and directories exist"
	MsgGitProperlyInitialized    = "Git repository properly initialized"
	MsgAllEndpointsResponding    = "All server endpoints responding correctly"
	MsgSEPTEOBrandingIntegrated  = "SEPTEO branding properly integrated"
	MsgDockerComposeValid        = "Docker Compose configuration is valid"
	MsgAllDeployScriptsReady     = "All deployment scripts exist and are executable"
	MsgEndpointRespondingCorrect = "Endpoint responding correctly"
)

// File and Directory Paths
const (
	RequiredFileGoMod         = "go.mod"
	RequiredFileMainGo        = "main.go"
	RequiredFileReadme        = "README.md"
	RequiredFileGitignore     = ".gitignore"
	RequiredFileDockerCompose = "docker-compose.yml"
	RequiredDirConfig         = "config"
	RequiredDirDeploy         = "deploy"
	RequiredDirTemplates      = "templates"
	RequiredDirInternal       = "internal"
	RequiredDirInternalLogger = "internal/logger"
	RequiredDirInternalDebug  = "internal/debug"
	RequiredDirScripts        = "scripts"
	GitDirectory              = ".git"
)

// Docker Compose Service Names
const (
	DockerServiceApp = "app:"
	DockerServiceDB  = "db:"
)

// Deploy Script Names
const (
	DeployScriptDeploy        = "deploy/deploy.sh"
	DeployScriptSetupInfomaniak = "deploy/setup-infomaniak.sh"
	ScriptsDeploy            = "scripts/deploy.sh"
	ScriptsStart             = "scripts/start.sh"
	ScriptsStop              = "scripts/stop.sh"
	ScriptsRestart           = "scripts/restart.sh"
)

// JSON Field Names for Phase Tests (additional to data_integrity_constants.go)
const (
	JSONFieldPhase       = "phase"  
	JSONFieldTests       = "tests"
	JSONFieldSubTests    = "sub_tests"
	JSONFieldName        = "name"
	JSONFieldTotalTests  = "total_tests"
	JSONFieldPassedTests = "passed_tests"
	JSONFieldFailedTests = "failed_tests"
	JSONFieldDuration    = "duration"
	JSONFieldDetails     = "details"
	// Note: message, error, description already defined in data_integrity_constants.go
)

// Test Detail Keys
const (
	DetailAppName              = "app_name"
	DetailAppIcon              = "app_icon" 
	DetailPoweredBy            = "powered_by"
	DetailPrimaryColor         = "primary_color"
	DetailServerPort           = "server_port"
	DetailIssuesFound          = "issues_found"
	DetailIssues               = "issues"
	DetailRequiredFiles        = "required_files"
	DetailRequiredDirs         = "required_dirs"
	DetailMissingFiles         = "missing_files"
	DetailMissingDirs          = "missing_dirs"
	DetailMissingFilesList     = "missing_files_list"
	DetailMissingDirsList      = "missing_dirs_list"
	DetailGitInitialized       = "git_initialized"
	DetailGitignoreExists      = "gitignore_exists"
	DetailBaseURL              = "base_url"
	DetailTotalEndpoints       = "total_endpoints"
	DetailPassedEndpoints      = "passed_endpoints"
	DetailURL                  = "url"
	DetailExpectedContentType  = "expected_content_type"
	DetailAccessible           = "accessible"
	DetailStatusCode           = "status_code"
	DetailContentType          = "content_type"
	DetailSepteoLogoIntegrated = "septeo_logo_integrated"
	DetailSepteoOrangeIntegrated = "septeo_orange_integrated"
	DetailSalamanderIconIntegrated = "salamander_icon_integrated"
	DetailHasAppService        = "has_app_service"
	DetailHasDbService         = "has_db_service"
	DetailHasPortMapping       = "has_port_mapping"
	DetailTotalScripts         = "total_scripts"
	DetailMissingScripts       = "missing_scripts"
	DetailNonExecutableScripts = "non_executable_scripts"
	DetailMissingScriptsList   = "missing_scripts_list"
	DetailNonExecutableScriptsList = "non_executable_scripts_list"
)

// Branding Integration Constants
const (
	SepteoLogoPath      = "septeo.svg"
	SepteoOrangeColor   = "#ff6136"
	SalamanderIcon      = "🦎"
)

// HTTP Endpoint Paths
const (
	EndpointHealth = "/health"
	EndpointDebug  = "/debug"
	EndpointHome   = "/"
)

// Content Type Constants
const (
	ContentTypeJSON = "application/json"
	ContentTypeHTML = "text/html"
)

// Test Endpoint Names
const (
	TestEndpointHealth = "Health Endpoint"
	TestEndpointDebug  = "Debug Endpoint"
	TestEndpointHome   = "Home Page"
)

// Error and Status Messages (additional)
const (
	MsgCannotTestServer             = "Cannot test server - configuration is nil"
	MsgEndpointNotAccessible        = "Endpoint not accessible"
	MsgWrongContentType            = "Wrong content type: got %s, expected %s"
	MsgUnexpectedStatusCode        = "Unexpected status code: %d"
	MsgSepteoLogoNotFound          = "SEPTEO logo URL not found in main.go"
	MsgSepteoOrangeNotFound        = "SEPTEO orange color not found"
	MsgSalamanderIconNotFound      = "Fire Salamander icon not found"
	MsgCannotReadMainGo            = "Cannot read main.go file"
	MsgDockerComposeFileMissing    = "docker-compose.yml file missing"
	MsgConfigurationFailed         = "Configuration validation failed: %s"
	MsgMissingFilesAndDirs         = "Missing %d files and %d directories"
	MsgGitSetupIssues              = "Git setup issues: %s"
	MsgOnlyEndpointsResponding     = "Only %d/%d endpoints responding correctly"
	MsgBrandingIssues              = "Branding issues: %s"
	MsgDeployScriptIssues          = "Deploy script issues: %s"
	MsgMissingScriptsCount         = "%d missing scripts"
	MsgNonExecutableScriptsCount   = "%d non-executable scripts"
	MsgGitDirectoryMissing         = ".git directory missing"
	MsgGitignoreMissing            = ".gitignore missing"
)

// ========================================
// ALPHA-4 HEALTH CHECK CONSTANTS
// ========================================

// Health Status Values (ALPHA-4)
const (
	HealthStatusHealthy  = "healthy"
	HealthStatusDegraded = "degraded"
	HealthStatusError    = "error"
)

// Check Status Values (ALPHA-4)
const (
	CheckStatusOK       = "ok"
	CheckStatusError    = "error"
	CheckStatusWarning  = "warning"
)

// Test Status Values (ALPHA-4)
const (
	TestStatusPassed = "passed"
	TestStatusFailed = "failed"
)

// ========================================
// ALPHA-4 DEBUG MESSAGES
// ========================================

// General Messages (ALPHA-4)
const (
	MsgCheckFailed                = "Check failed"
	MsgCheckPassed               = "Check passed"
	MsgHealthCheckCompleted      = "Health check completed"
	MsgDebugEndpointCalled       = "Debug endpoint called"
	MsgRunningPhase1Tests        = "Running Phase 1 tests"
	MsgCreatingHealthChecker     = "Creating new health checker"
	MsgRunningAllHealthChecks    = "Running all health checks"
	MsgRunningCheck              = "Running check"
	MsgDebugResponseEncodeFailed = "Failed to encode debug response"
	MsgPhaseTestsEncodeFailed    = "Failed to encode phase tests response"
	MsgInternalServerError       = "Internal server error"
)

// Configuration Messages (ALPHA-4)
const (
	MsgConfigIsNil             = "Configuration is nil"
	MsgConfigValidationFailed  = "Configuration validation failed"
	MsgConfigIsValid          = "Configuration is valid"
	MsgAppNameEmpty           = "app.name is empty"
	MsgServerPortInvalid      = "server.port is invalid"
	MsgDBPathEmpty            = "database.path is empty"
)

// Database Messages (ALPHA-4)
const (
	MsgSQLitePathEmpty         = "SQLite path is empty"
	MsgSQLiteConfigValid      = "SQLite configuration is valid"
	MsgSQLiteDirNotExist      = "SQLite directory doesn't exist yet"
	MsgMySQLConfigIncomplete  = "MySQL configuration incomplete"
	MsgMySQLConfigValid       = "MySQL configuration is valid"
	MsgUnknownDatabaseType    = "Unknown database type"
)

// Filesystem Messages (ALPHA-4)
const (
	MsgRequiredFilesMissing = "Required files/directories missing"
	MsgAllFilesPresent     = "All required files and directories present"
)

// Network Messages (ALPHA-4)
const (
	MsgNetworkConfigValid = "Network configuration valid"
)

// AI Messages (ALPHA-4)
const (
	MsgAIDisabled       = "AI is disabled"
	MsgAIMockMode      = "AI is in mock mode"
	MsgAIConfigValid   = "AI configuration is valid"
)

// ========================================
// ALPHA-4 ERROR CODES
// ========================================

const (
	ErrorCodeConfigNil               = "config_nil"
	ErrorCodeSQLitePathMissing      = "sqlite_path_missing"
	ErrorCodeMySQLConfigIncomplete  = "mysql_config_incomplete"
	ErrorCodeUnknownDBType          = "unknown_db_type"
)

// ========================================
// ALPHA-4 DATABASE CONSTANTS
// ========================================

const (
	DatabaseTypeSQLite    = "sqlite"
	DatabaseTypeMySQL     = "mysql"
	DefaultDatabaseName   = "firesalamander"
	DefaultDatabaseFile   = "/firesalamander.db"
)

// ========================================
// ALPHA-4 DIRECTORY & FILE PATHS
// ========================================

const (
	DeployDir    = "deploy"
	GoModFile    = "go.mod"
	MainGoFile   = "main.go"
)

// ========================================
// ALPHA-4 HTTP CONSTANTS
// ========================================

const (
	HeaderContentType  = "Content-Type"
)

// ========================================
// ALPHA-4 ENVIRONMENT CONSTANTS
// ========================================

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// ========================================
// ALPHA-4 DEBUG FIELD NAMES
// ========================================

const (
	DebugFieldCheck      = "check"
	DebugFieldStatus     = "status"
	DebugFieldMessage    = "message"
	DebugFieldError      = "error"
	DebugFieldMethod     = "method"
	DebugFieldPath       = "path"
	DebugFieldRemote     = "remote"
	DebugFieldChecks     = "checks"
	DebugFieldEnabled    = "enabled"
	DebugFieldMockMode   = "mock_mode"
	DebugFieldAPIKey     = "api_key"
	DebugFieldHost       = "host"
	DebugFieldName       = "name"
	DebugFieldType       = "type"
	DebugFieldPort       = "port"
	DebugFieldAddr       = "addr"
	DebugFieldAppName    = "app_name"
	DebugFieldServerPort = "server_port"
	DebugFieldDBType     = "database_type"
)

// ========================================
// ALPHA-4 DEBUG PREFIXES & SUFFIXES
// ========================================

const (
	DebugDirPrefix     = "directory: "
	DebugFilePrefix    = "file: "
	DebugMaskSuffix    = "***"
	DebugQueryParam    = "test"
	DebugPhase1Value   = "phase1"
	DebugEnvVariable   = "ENV"
)

// ========================================
// ALPHA-4 NUMERIC DEBUG CONSTANTS
// ========================================

const (
	MemoryDivisor1024   = 1024
	MemoryDivisor1024x2 = 1024 * 1024
	APIKeySuffixLength  = 4
)