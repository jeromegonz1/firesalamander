package constants

// Data Integrity Agent Constants
// Constantes pour l'agent d'intégrité des données

// Database and Connection Constants
const (
	// Database Types
	SQLite3Driver = "sqlite3"
	
	// Default Database Settings
	DefaultDatabasePath = "fire_salamander_dev.db"
	DefaultReportPath   = "tests/reports/data"
	DefaultTimeout      = 30
	
	// Database Tables
	TableCrawlSessions = "crawl_sessions"
	TablePages         = "pages" 
	TableSEOMetrics    = "seo_metrics"
	TableUsers         = "users"
	TableConfigs       = "configs"
	TableAnalyses      = "analyses"
)

// Test Categories
const (
	TestCategorySchemaValidation    = "schema_validation"
	TestCategoryDataConsistency     = "data_consistency"
	TestCategoryReferentialIntegrity = "referential_integrity"
	TestCategoryDataQuality         = "data_quality"
	TestCategoryPerformanceChecks   = "performance_checks"
)

// Status Constants (additional to main constants)
const (
	StatusUnknown   = "unknown"
	StatusPassed    = "passed"
	StatusWarning   = "warning"
	StatusHealthy   = "healthy"
	StatusCritical  = "critical"
	StatusRunning   = "running"
	StatusCompleted = "completed"
)

// Test Names and Descriptions
const (
	TestSchemaValidation    = "Schema Validation"
	TestConstraintsCheck    = "Constraints Check"
	TestDataConsistency     = "Data Consistency"
	TestReferentialCheck    = "Referential Integrity"
	TestDataQuality         = "Data Quality"
	TestPerformanceCheck    = "Performance Check"
	TestIndexValidation     = "Index Validation"
	TestForeignKeyCheck     = "Foreign Key Check"
	TestDataTypeValidation  = "Data Type Validation"
	TestNullConstraintCheck = "Null Constraint Check"
	TestDataConstraints     = "Data Constraints"
	TestTimestampConsistency = "Timestamp Consistency"
	TestNullValues          = "NULL Values"
	TestUniqueConstraints   = "Unique Constraints"
	TestURLQuality          = "URL Quality"
	TestHTTPStatusCodes     = "HTTP Status Codes"
	TestSEOScoreValidity    = "SEO Score Validity"
	TestQueryPerformance    = "Query Performance"
	TestDatabaseSize        = "Database Size"
	TestNumericConsistency  = "Numeric Consistency"
)

// Issue Types
const (
	IssueTypeSchema     = "schema"
	IssueTypeData       = "data" 
	IssueTypeIntegrity  = "integrity"
	IssueTypeQuality    = "quality"
	IssueTypePerformance = "performance"
	IssueTypeConstraint = "constraint"
)

// Severity Levels
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Common Messages and Descriptions
const (
	MsgMissingRequiredTable     = "Missing required table"
	MsgApplicationImpaired      = "Application functionality may be impaired"
	MsgDataInconsistency        = "Data inconsistency detected"
	MsgReferentialIntegrityFail = "Referential integrity violation"
	MsgPerformanceDegradation   = "Performance degradation detected"
	MsgSchemaValidationFailed   = "Schema validation failed"
	MsgConstraintViolation      = "Constraint violation detected"
	MsgDataQualityIssue         = "Data quality issue identified"
	MsgIndexMissing             = "Required index is missing"
	MsgForeignKeyViolation      = "Foreign key constraint violation"
	MsgNullConstraintViolation  = "Null constraint violation"
	MsgDataTypeIncorrect        = "Incorrect data type detected"
	MsgDuplicateRecords         = "Duplicate records found"
	MsgOrphanedRecords          = "Orphaned records detected"
	MsgNullOrEmptyURL           = "NULL or empty URL values found"
	MsgCrawlSessionsNoURL       = "Crawl sessions cannot function without valid URLs"
	MsgDuplicatePageRecords     = "Duplicate page records found"
	MsgInconsistentCrawlResults = "May cause inconsistent crawl results"
	MsgNoDuplicateRecords       = "No duplicate page records found"
	MsgAllTimestampsConsistent  = "All timestamps are consistent"
	MsgPageCountConsistent      = "Page count relationships are consistent"
	MsgAllURLsWellFormed        = "All URLs appear to be well-formed"
	MsgAllStatusCodesValid      = "All status codes are in valid range"
	MsgAllSEOScoresValid        = "All SEO scores are within valid range"
	MsgNoOrphanedRecords        = "No orphaned records found"
)

// SQL Query Constants
const (
	QueryGetTableList    = "SELECT name FROM sqlite_master WHERE type='table'"
	QueryTableExists     = "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	QueryRowCount        = "SELECT COUNT(*) FROM "
	QueryConstraintCheck = "PRAGMA foreign_key_check"
	QueryIndexList       = "PRAGMA index_list(%s)"
	QueryTableInfo       = "PRAGMA table_info(%s)"
	QueryIntegrityCheck  = "PRAGMA integrity_check"
)

// Report and File Constants  
const (
	ReportFileExtension     = ".json"
	ReportDateTimeFormat    = "2006-01-02_15-04-05"
	BackupFileExtension     = ".backup"
	TempFilePrefix          = "temp_"
	LogFileExtension        = ".log"
	ConfigFileExtension     = ".yaml"
)

// JSON Field Names (API Contract Constants)
const (
	JSONFieldTimestamp      = "timestamp"
	JSONFieldDatabase       = "database"
	JSONFieldTestResults    = "test_results"
	JSONFieldIssues         = "issues"
	JSONFieldOverallScore   = "overall_score"
	JSONFieldStatus         = "status"
	JSONFieldTest           = "test"
	JSONFieldDescription    = "description"
	JSONFieldValue          = "value"
	JSONFieldExpected       = "expected"
	JSONFieldSeverity       = "severity"
	JSONFieldType           = "type"
	JSONFieldTable          = "table"
	JSONFieldColumn         = "column"
	JSONFieldIssue          = "issue"
	JSONFieldImpact         = "impact"
	JSONFieldCount          = "count"
	JSONFieldTimeout        = "timeout"
	JSONFieldDatabasePath   = "database_path"
	JSONFieldReportPath     = "report_path"
	JSONFieldTestCategories = "test_categories"
)

// Performance Thresholds  
const (
	MaxQueryExecutionTime     = 5000  // milliseconds
	MaxTableScanTime          = 10000 // milliseconds
	MaxIndexLookupTime        = 1000  // milliseconds
	MinAcceptablePerformance  = 80    // percentage
	MaxAcceptableRecordCount  = 1000000
	// HighQualityScore already defined in main constants
)

// Validation Rules
const (
	MinTableNameLength = 2
	MaxTableNameLength = 64
	MinColumnNameLength = 1
	MaxColumnNameLength = 64
	MaxStringLength     = 255
	MaxTextLength       = 65535
)

// SQL Keywords and Reserved Words
const (
	SQLKeywordCHECK = "CHECK"
)

// Error and Impact Messages (additional to main constants)
const (
	MsgIncorrectSessionDurationCalculations = "Incorrect session duration calculations"
	MsgDataInconsistencyAndPotentialErrors  = "Data inconsistency and potential application errors"
)

// Query Performance Test Names
const (
	QueryNameSimpleCount = "Simple Count"
	QueryNameComplexJoin = "Complex Join"
)

// Quality Status Labels (scoring system)
const (
	StatusExcellent  = "excellent"
	StatusAcceptable = "acceptable"
)

// HTML Template CSS Classes
const (
	HTMLClassHeader    = "header"
	HTMLClassSection   = "section"
	HTMLClassScore     = "score"
	HTMLClassTestResult = "test-result"
	HTMLClassPass      = "pass"
	HTMLClassFail      = "fail"
	HTMLClassWarning   = "warning"
	HTMLClassError     = "error"
	HTMLClassSeverityHigh   = "severity-high"
	HTMLClassSeverityMedium = "severity-medium"
	HTMLClassSeverityLow    = "severity-low"
	HTMLClassStatusExcellent        = "status-excellent"
	HTMLClassStatusGood             = "status-good"
	HTMLClassStatusAcceptable       = "status-acceptable"
	HTMLClassStatusNeedsImprovement = "status-needs_improvement"
	HTMLClassStatusPoor             = "status-poor"
)