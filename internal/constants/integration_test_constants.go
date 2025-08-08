package constants

// ========================================
// DELTA-4 INTEGRATION TEST CONSTANTS
// Integration Testing Configuration and Validation Constants
// ========================================

// ========================================
// INTEGRATION TEST NAMES
// ========================================

// Core Test Names
const (
	IntegrationTestNameTesting         = "testing"
	IntegrationTestNameIntegration     = "integration"
	IntegrationTestNameAPI             = "api"
	IntegrationTestNameEndpoint        = "endpoint"
	IntegrationTestNameValidation      = "validation"
	IntegrationTestNamePerformance     = "performance"
	IntegrationTestNameSecurity        = "security"
	IntegrationTestNameAccessibility   = "accessibility"
	IntegrationTestNameCompatibility   = "compatibility"
	IntegrationTestNameRegression      = "regression"
	IntegrationTestNameSmoke           = "smoke"
	IntegrationTestNameSanity          = "sanity"
	IntegrationTestNameAcceptance      = "acceptance"
	IntegrationTestNameSystem          = "system"
)

// Specific Test Names
const (
	IntegrationTestNameInsight         = "test_insight"
	IntegrationTestNameTask            = "test_task"
	IntegrationTestNameAction          = "test_action"
	IntegrationTestNameAnalysis        = "test_analysis"
	IntegrationTestNameReport          = "test_report"
	IntegrationTestNameMetrics         = "test_metrics"
	IntegrationTestNameData            = "test_data"
	IntegrationTestNameCrawler         = "test_crawler"
	IntegrationTestNameSEO             = "test_seo"
	IntegrationTestNameSemantic        = "test_semantic"
)

// ========================================
// INTEGRATION TEST PHASES
// ========================================

// BDD Test Phases
const (
	IntegrationTestPhaseGiven    = "given"
	IntegrationTestPhaseWhen     = "when"
	IntegrationTestPhaseThen     = "then"
)

// AAA Test Phases
const (
	IntegrationTestPhaseArrange  = "arrange"
	IntegrationTestPhaseAct      = "act"
	IntegrationTestPhaseAssert   = "assert"
)

// Lifecycle Test Phases
const (
	IntegrationTestPhaseSetup      = "setup"
	IntegrationTestPhaseTeardown   = "teardown"
	IntegrationTestPhaseBefore     = "before"
	IntegrationTestPhaseAfter      = "after"
	IntegrationTestPhasePrepare    = "prepare"
	IntegrationTestPhaseExecute    = "execute"
	IntegrationTestPhaseCleanup    = "cleanup"
	IntegrationTestPhaseInitialize = "initialize"
	IntegrationTestPhaseFinalize   = "finalize"
)

// ========================================
// INTEGRATION ENDPOINTS
// ========================================

// API Endpoints
const (
	IntegrationEndpointAPIRoot     = "/api/"
	IntegrationEndpointAPIAnalyze  = "/api/analyze"
	IntegrationEndpointAPIResults  = "/api/results"
	IntegrationEndpointAPIReports  = "/api/reports"
	IntegrationEndpointAPIMetrics  = "/api/metrics"
	IntegrationEndpointAPIStatus   = "/api/status"
	IntegrationEndpointAPIHealth   = "/api/health"
	IntegrationEndpointAPIDebug    = "/api/debug"
	IntegrationEndpointAPITest     = "/api/test"
)

// System Endpoints
const (
	IntegrationEndpointHealth  = "/health"
	IntegrationEndpointStatus  = "/status"
	IntegrationEndpointDebug   = "/debug"
	IntegrationEndpointMetrics = "/metrics"
	IntegrationEndpointAnalyze = "/analyze"
	IntegrationEndpointResults = "/results"
	IntegrationEndpointReports = "/reports"
	IntegrationEndpointTest    = "/test"
)

// ========================================
// INTEGRATION HTTP METHODS
// ========================================

// HTTP Methods
const (
	IntegrationHTTPMethodGet     = "GET"
	IntegrationHTTPMethodPost    = "POST"
	IntegrationHTTPMethodPut     = "PUT"
	IntegrationHTTPMethodDelete  = "DELETE"
	IntegrationHTTPMethodPatch   = "PATCH"
	IntegrationHTTPMethodHead    = "HEAD"
	IntegrationHTTPMethodOptions = "OPTIONS"
)

// ========================================
// INTEGRATION HTTP STATUS CODES
// ========================================

// Success Status Codes
const (
	IntegrationHTTPStatusOK        = "200"
	IntegrationHTTPStatusCreated   = "201"
	IntegrationHTTPStatusAccepted  = "202"
	IntegrationHTTPStatusNoContent = "204"
)

// Redirection Status Codes
const (
	IntegrationHTTPStatusMovedPermanently = "301"
	IntegrationHTTPStatusFound           = "302"
	IntegrationHTTPStatusNotModified     = "304"
)

// Client Error Status Codes
const (
	IntegrationHTTPStatusBadRequest          = "400"
	IntegrationHTTPStatusUnauthorized        = "401"
	IntegrationHTTPStatusForbidden           = "403"
	IntegrationHTTPStatusNotFound            = "404"
	IntegrationHTTPStatusMethodNotAllowed    = "405"
	IntegrationHTTPStatusConflict            = "409"
	IntegrationHTTPStatusUnprocessableEntity = "422"
	IntegrationHTTPStatusTooManyRequests     = "429"
)

// Server Error Status Codes
const (
	IntegrationHTTPStatusInternalServerError = "500"
	IntegrationHTTPStatusNotImplemented      = "501"
	IntegrationHTTPStatusBadGateway          = "502"
	IntegrationHTTPStatusServiceUnavailable  = "503"
	IntegrationHTTPStatusGatewayTimeout      = "504"
)

// ========================================
// INTEGRATION CONTENT TYPES
// ========================================

// MIME Content Types
const (
	IntegrationContentTypeJSON  = "application/json"
	IntegrationContentTypeHTML  = "text/html"
	IntegrationContentTypePlain = "text/plain"
	IntegrationContentTypeXML   = "application/xml"
	IntegrationContentTypeCSV   = "text/csv"
	IntegrationContentTypePDF   = "application/pdf"
)

// ========================================
// INTEGRATION TEST ASSERTIONS
// ========================================

// Assertion Keywords
const (
	IntegrationAssertionExpected = "expected"
	IntegrationAssertionActual   = "actual"
	IntegrationAssertionAssert   = "assert"
	IntegrationAssertionVerify   = "verify"
	IntegrationAssertionValidate = "validate"
	IntegrationAssertionCheck    = "check"
	IntegrationAssertionEnsure   = "ensure"
	IntegrationAssertionConfirm  = "confirm"
)

// Comparison Assertions
const (
	IntegrationAssertionMatch    = "match"
	IntegrationAssertionEqual    = "equal"
	IntegrationAssertionContain  = "contain"
	IntegrationAssertionInclude  = "include"
	IntegrationAssertionStart    = "start"
	IntegrationAssertionEnd      = "end"
	IntegrationAssertionGreater  = "greater"
	IntegrationAssertionLess     = "less"
)

// State Assertions
const (
	IntegrationAssertionEmpty   = "empty"
	IntegrationAssertionNil     = "nil"
	IntegrationAssertionTrue    = "true"
	IntegrationAssertionFalse   = "false"
	IntegrationAssertionSuccess = "success"
	IntegrationAssertionFailure = "failure"
	IntegrationAssertionError   = "error"
	IntegrationAssertionValid   = "valid"
	IntegrationAssertionInvalid = "invalid"
)

// ========================================
// INTEGRATION JSON FIELD NAMES
// ========================================

// Standard Response Fields
const (
	IntegrationJSONFieldID        = "id"
	IntegrationJSONFieldName      = "name"
	IntegrationJSONFieldType      = "type"
	IntegrationJSONFieldStatus    = "status"
	IntegrationJSONFieldData      = "data"
	IntegrationJSONFieldMessage   = "message"
	IntegrationJSONFieldError     = "error"
	IntegrationJSONFieldTimestamp = "timestamp"
	IntegrationJSONFieldURL       = "url"
	IntegrationJSONFieldMethod    = "method"
	IntegrationJSONFieldPath      = "path"
)

// Test Response Fields
const (
	IntegrationJSONFieldResult      = "result"
	IntegrationJSONFieldResponse    = "response"
	IntegrationJSONFieldRequest     = "request"
	IntegrationJSONFieldPayload     = "payload"
	IntegrationJSONFieldHeaders     = "headers"
	IntegrationJSONFieldBody        = "body"
	IntegrationJSONFieldTest        = "test"
	IntegrationJSONFieldSuite       = "suite"
	IntegrationJSONFieldCase        = "case"
	IntegrationJSONFieldScenario    = "scenario"
	IntegrationJSONFieldAssertion   = "assertion"
	IntegrationJSONFieldExpectation = "expectation"
	IntegrationJSONFieldActual      = "actual"
	IntegrationJSONFieldExpected    = "expected"
	IntegrationJSONFieldDuration    = "duration"
	IntegrationJSONFieldTimeout     = "timeout"
)

// ========================================
// INTEGRATION TEST CATEGORIES
// ========================================

// Test Types
const (
	IntegrationTestCategoryUnit          = "unit"
	IntegrationTestCategoryIntegration   = "integration"
	IntegrationTestCategoryE2E           = "e2e"
	IntegrationTestCategoryAPI           = "api"
	IntegrationTestCategoryUI            = "ui"
	IntegrationTestCategoryPerformance   = "performance"
	IntegrationTestCategorySecurity      = "security"
	IntegrationTestCategoryAccessibility = "accessibility"
	IntegrationTestCategoryCompatibility = "compatibility"
	IntegrationTestCategoryRegression    = "regression"
	IntegrationTestCategorySmoke         = "smoke"
	IntegrationTestCategorySanity        = "sanity"
	IntegrationTestCategoryAcceptance    = "acceptance"
	IntegrationTestCategorySystem        = "system"
)

// ========================================
// INTEGRATION AGENT NAMES
// ========================================

// Analysis Agents
const (
	IntegrationAgentCrawler       = "crawler"
	IntegrationAgentSEO          = "seo"
	IntegrationAgentSemantic     = "semantic"
	IntegrationAgentPerformance  = "performance"
	IntegrationAgentSecurity     = "security"
	IntegrationAgentQA           = "qa"
	IntegrationAgentDataIntegrity = "data_integrity"
	IntegrationAgentFrontend     = "frontend"
	IntegrationAgentPlaywright   = "playwright"
	IntegrationAgentK6           = "k6"
)

// Test Agents
const (
	IntegrationAgentIntegration = "integration"
	IntegrationAgentUnit        = "unit"
	IntegrationAgentE2E         = "e2e"
)

// ========================================
// INTEGRATION ANALYSIS TYPES
// ========================================

// Analysis Categories
const (
	IntegrationAnalysisTypeTechnical     = "technical"
	IntegrationAnalysisTypeContent       = "content"
	IntegrationAnalysisTypePerformance   = "performance"
	IntegrationAnalysisTypeSecurity      = "security"
	IntegrationAnalysisTypeAccessibility = "accessibility"
	IntegrationAnalysisTypeSEO           = "seo"
	IntegrationAnalysisTypeSemantic      = "semantic"
	IntegrationAnalysisTypeStructural    = "structural"
	IntegrationAnalysisTypeFunctional    = "functional"
	IntegrationAnalysisTypeRegression    = "regression"
	IntegrationAnalysisTypeSmoke         = "smoke"
	IntegrationAnalysisTypeSanity        = "sanity"
)

// ========================================
// INTEGRATION TEST DATA
// ========================================

// Mock Data Prefixes
const (
	IntegrationTestDataMock    = "mock_"
	IntegrationTestDataTest    = "test_"
	IntegrationTestDataSample  = "sample_"
	IntegrationTestDataDummy   = "dummy_"
	IntegrationTestDataFake    = "fake_"
	IntegrationTestDataExample = "example_"
)

// Test Data Types
const (
	IntegrationTestDataTask     = "test_task"
	IntegrationTestDataInsight  = "test_insight"
	IntegrationTestDataAction   = "test_action"
	IntegrationTestDataAnalysis = "test_analysis"
	IntegrationTestDataReport   = "test_report"
	IntegrationTestDataMetrics  = "test_metrics"
	IntegrationTestDataCrawler  = "test_crawler"
	IntegrationTestDataSEO      = "test_seo"
	IntegrationTestDataSemantic = "test_semantic"
)

// ========================================
// INTEGRATION TEST CONFIG KEYS
// ========================================

// Test Configuration
const (
	IntegrationConfigTimeout   = "timeout"
	IntegrationConfigRetries   = "retries"
	IntegrationConfigParallel  = "parallel"
	IntegrationConfigVerbose   = "verbose"
	IntegrationConfigDebug     = "debug"
)

// Mock Configuration
const (
	IntegrationConfigMock     = "mock"
	IntegrationConfigStub     = "stub"
	IntegrationConfigFake     = "fake"
	IntegrationConfigCleanup  = "cleanup"
	IntegrationConfigSetup    = "setup"
	IntegrationConfigTeardown = "teardown"
)

// Service Configuration
const (
	IntegrationConfigDatabase = "database"
	IntegrationConfigServer   = "server"
	IntegrationConfigClient   = "client"
	IntegrationConfigHost     = "host"
	IntegrationConfigPort     = "port"
)

// ========================================
// INTEGRATION ERROR MESSAGES
// ========================================

// Validation Errors
const (
	IntegrationErrorURLInvalid               = "URL incorrecte"
	IntegrationErrorScoreInvalid             = "Score global invalide"
	IntegrationErrorProcessingTimeInvalid    = "Temps de traitement invalide"
	IntegrationErrorReportSizeInvalid        = "Taille de rapport invalide"
	IntegrationErrorDataFormatInvalid        = "Format de données invalide"
	IntegrationErrorRequestInvalid           = "Requête invalide"
	IntegrationErrorResponseInvalid          = "Réponse invalide"
	IntegrationErrorPayloadInvalid           = "Payload invalide"
)

// Test Execution Errors
const (
	IntegrationErrorTestFailed        = "Test échoué"
	IntegrationErrorTestTimeout       = "Timeout du test"
	IntegrationErrorTestSetupFailed   = "Échec de la configuration du test"
	IntegrationErrorTestCleanupFailed = "Échec du nettoyage du test"
	IntegrationErrorAssertionFailed   = "Assertion échouée"
	IntegrationErrorMockSetupFailed   = "Échec de la configuration du mock"
)

// Service Errors
const (
	IntegrationErrorServiceUnavailable = "Service indisponible"
	IntegrationErrorServiceTimeout     = "Timeout du service"
	IntegrationErrorServiceConnection  = "Erreur de connexion au service"
	IntegrationErrorDatabaseConnection = "Erreur de connexion à la base de données"
	IntegrationErrorAPIConnection      = "Erreur de connexion à l'API"
)

// ========================================
// INTEGRATION SUCCESS MESSAGES
// ========================================

// Test Success Messages
const (
	IntegrationSuccessTestPassed       = "Test réussi"
	IntegrationSuccessTestCompleted    = "Test terminé avec succès"
	IntegrationSuccessSuiteCompleted   = "Suite de tests terminée avec succès"
	IntegrationSuccessAssertionPassed  = "Assertion réussie"
	IntegrationSuccessValidationPassed = "Validation réussie"
	IntegrationSuccessSetupCompleted   = "Configuration terminée avec succès"
	IntegrationSuccessCleanupCompleted = "Nettoyage terminé avec succès"
)

// Analysis Success Messages
const (
	IntegrationSuccessAnalysisCompleted   = "Analyse terminée avec succès"
	IntegrationSuccessReportGenerated     = "Rapport généré avec succès"
	IntegrationSuccessMetricsCalculated   = "Métriques calculées avec succès"
	IntegrationSuccessDataProcessed       = "Données traitées avec succès"
	IntegrationSuccessValidationCompleted = "Validation terminée avec succès"
)

// ========================================
// INTEGRATION LOG MESSAGES
// ========================================

// Test Execution Logs
const (
	IntegrationLogTestStarted    = "Test démarré: %s"
	IntegrationLogTestRunning    = "Exécution du test: %s"
	IntegrationLogTestCompleted  = "Test terminé: %s en %v"
	IntegrationLogTestFailed     = "Test échoué: %s - %s"
	IntegrationLogSuiteStarted   = "Suite démarrée: %s"
	IntegrationLogSuiteCompleted = "Suite terminée: %s - %d/%d tests réussis"
)

// Setup and Cleanup Logs
const (
	IntegrationLogSetupStarted    = "Configuration démarrée"
	IntegrationLogSetupCompleted  = "Configuration terminée"
	IntegrationLogCleanupStarted  = "Nettoyage démarré"
	IntegrationLogCleanupCompleted = "Nettoyage terminé"
	IntegrationLogMockSetup       = "Configuration du mock: %s"
	IntegrationLogMockCleanup     = "Nettoyage du mock: %s"
)

// Analysis Logs
const (
	IntegrationLogAnalysisStarted   = "Analyse démarrée: %s"
	IntegrationLogAnalysisProgress  = "Progression de l'analyse: %s - %s"
	IntegrationLogAnalysisCompleted = "Analyse terminée: %s en %v"
	IntegrationLogAnalysisFailed    = "Analyse échouée: %s - %s"
)

// ========================================
// INTEGRATION MOCK RESPONSES
// ========================================

// Mock Response Types
const (
	IntegrationMockResponseSuccess = "mock_response_success"
	IntegrationMockResponseError   = "mock_response_error"
	IntegrationMockResponseTimeout = "mock_response_timeout"
	IntegrationMockResponseEmpty   = "mock_response_empty"
	IntegrationMockResponseLarge   = "mock_response_large"
)

// ========================================
// INTEGRATION TEST FIXTURES
// ========================================

// Test Data Fixtures
const (
	IntegrationFixtureTestData  = "fixture_testdata"
	IntegrationFixtureMockData  = "fixture_mockdata"
	IntegrationFixtureSampleData = "fixture_sampledata"
	IntegrationFixtureTestCase  = "fixture_testcase"
	IntegrationFixtureScenario  = "fixture_scenario"
	IntegrationFixtureTemplate  = "fixture_template"
)

// ========================================
// INTEGRATION DEFAULT VALUES
// ========================================

// Default Test Configuration
const (
	IntegrationDefaultTimeout      = 30 // seconds
	IntegrationDefaultRetries      = 3
	IntegrationDefaultParallelJobs = 4
)

// Default Test Timeouts
const (
	IntegrationDefaultUnitTestTimeout        = 5  // seconds
	IntegrationDefaultIntegrationTestTimeout = 30 // seconds
	IntegrationDefaultE2ETestTimeout         = 60 // seconds
	IntegrationDefaultPerformanceTestTimeout = 300 // seconds
)

// Default Test Sizes
const (
	IntegrationDefaultSmallTestSize  = 100
	IntegrationDefaultMediumTestSize = 1000
	IntegrationDefaultLargeTestSize  = 10000
)