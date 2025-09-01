package constants

// ========================================
// BRAVO-1 QA AGENT CONSTANTS
// ========================================

// QA Agent Identity
const (
	QAAgentName = "QA-AGENT"
)

// QA Status Values
const (
	QAStatusPass    = "pass"
	QAStatusFail    = "fail"
	QAStatusWarning = "warning"
	QAStatusSuccess = "success"
	QAStatusError   = "error"
)

// Test Status Values (QA specific - using different names to avoid conflicts)
const (
	QATestStatusPassed  = "passed"
	QATestStatusFailed  = "failed"
	QATestStatusSkipped = "skipped"
)

// Output Formats
const (
	OutputFormatJSON = "json"
	OutputFormatText = "text"
	OutputFormatHTML = "html"
	OutputFormatXML  = "xml"
	OutputFormatCSV  = "csv"
)

// QA Severity Levels (prefixed to avoid conflicts)
const (
	QASeverityHigh     = "high"
	QASeverityMedium   = "medium"
	QASeverityLow      = "low"
	QASeverityCritical = "critical"
	QASeverityWarning  = "warning"
	QASeverityInfo     = "info"
	QASeverityNotice   = "notice"
)

// Confidence Levels
const (
	ConfidenceHigh     = "high"
	ConfidenceMedium   = "medium"
	ConfidenceLow      = "low"
	ConfidenceCertain  = "certain"
	ConfidenceProbable = "probable"
	ConfidencePossible = "possible"
)

// Tool Names
const (
	GoTool         = "go"
	VetTool        = "vet"
	LintTool       = "lint"
	TestTool       = "test"
	BuildTool      = "build"
	GofmtTool      = "gofmt"
	GolintTool     = "golint"
	StaticcheckTool = "staticcheck"
)

// File Extensions (QA specific)
const (
	ExtGo       = ".go"
	ExtTxt      = ".txt"
	ExtMarkdown = ".md"
)

// ========================================
// QA JSON FIELD NAMES
// ========================================

// QA Configuration Fields (prefixed to avoid conflicts)
const (
	QAJSONFieldMinCoverage      = "min_coverage"
	QAJSONFieldEnableVet        = "enable_vet"
	QAJSONFieldEnableLint       = "enable_lint"
	QAJSONFieldEnableSecurity   = "enable_security"
	QAJSONFieldEnableComplexity = "enable_complexity"
	QAJSONFieldOutputFormat     = "output_format"
	QAJSONFieldReportPath       = "report_path"
)

// QA Coverage Fields
const (
	QAJSONFieldTotalCoverage     = "total_coverage"
	QAJSONFieldPackagesCoverage  = "packages_coverage"
	QAJSONFieldThreshold         = "threshold"
	QAJSONFieldPassed            = "passed"
)

// QA Issue Fields (prefixed to avoid conflicts)
const (
	QAJSONFieldFile       = "file"
	QAJSONFieldLine       = "line"
	QAJSONFieldColumn     = "column"
	QAJSONFieldMessage    = "message"
	QAJSONFieldCategory   = "category"
	QAJSONFieldRule       = "rule"
	QAJSONFieldSeverity   = "severity"
	QAJSONFieldConfidence = "confidence"
	QAJSONFieldFunction   = "function"
	QAJSONFieldComplexity = "complexity"
)

// QA Test Fields (prefixed to avoid conflicts)
const (
	QAJSONFieldTotalTests   = "total_tests"
	QAJSONFieldPassedTests  = "passed_tests"
	QAJSONFieldFailedTests  = "failed_tests"
	QAJSONFieldSkippedTests = "skipped_tests"
	QAJSONFieldDuration     = "duration"
)

// QA Result Fields
const (
	QAJSONFieldOverallScore = "overall_score"
	QAJSONFieldStatus       = "status"
	QAJSONFieldTimestamp    = "timestamp"
)

// ========================================
// QA COMMANDS AND TOOLS
// ========================================

// Go Commands
const (
	GoCommandTest     = "go test"
	GoCommandVet      = "go vet"
	GoCommandBuild    = "go build"
	GoCommandMod      = "go mod"
	GoCommandFmt      = "go fmt"
	GoCommandGet      = "go get"
	GoCommandInstall  = "go install"
)

// Command Arguments
const (
	TestArgCover     = "-cover"
	TestArgCoverProf = "-coverprofile"
	TestArgV         = "-v"
	TestArgRace      = "-race"
	TestArgShort     = "-short"
	VetArgAll        = "./..."
	FmtArgDiff       = "-d"
	FmtArgWrite      = "-w"
)

// ========================================
// QA MESSAGES AND TEMPLATES
// ========================================

// Log Messages
const (
	MsgQAAgentStarting        = "🧪 QA Agent starting analysis"
	MsgQAAgentCompleted       = "✅ QA Agent analysis completed"
	MsgRunningUnitTests       = "Running unit tests with coverage"
	MsgRunningVetAnalysis     = "Running go vet analysis"
	MsgRunningLintAnalysis    = "Running lint analysis"
	MsgRunningSecurityCheck   = "Running security analysis"
	MsgRunningComplexityCheck = "Running complexity analysis"
	MsgGeneratingReport       = "Generating QA report"
)

// Error Messages
const (
	MsgUnitTestsFailed      = "Unit tests failed"
	MsgCoverageAnalysisFailed = "Coverage analysis failed"
	MsgVetAnalysisFailed    = "Go vet analysis failed"
	MsgLintAnalysisFailed   = "Lint analysis failed"
	MsgSecurityCheckFailed  = "Security check failed"
	MsgComplexityCheckFailed = "Complexity check failed"
	MsgReportGenerationFailed = "Report generation failed"
)

// Success Messages
const (
	MsgAllTestsPassed       = "All unit tests passed"
	MsgCoverageThresholdMet = "Coverage threshold met"
	MsgNoVetIssues         = "No vet issues found"
	MsgNoLintIssues        = "No lint issues found"
	MsgNoSecurityIssues    = "No security issues found"
	MsgComplexityAcceptable = "Code complexity is acceptable"
)

// Warning Messages
const (
	MsgCoverageThresholdNotMet = "Coverage threshold not met"
	MsgVetIssuesFound          = "Go vet issues found"
	MsgLintIssuesFound         = "Lint issues found"
	MsgSecurityIssuesFound     = "Security issues found"
	MsgHighComplexityFound     = "High complexity functions found"
)

// ========================================
// QA REPORT TEMPLATES
// ========================================

// Report Sections
const (
	ReportSectionSummary     = "Executive Summary"
	ReportSectionCoverage    = "Test Coverage Analysis"
	ReportSectionTests       = "Unit Tests Results"
	ReportSectionVet         = "Go Vet Analysis"
	ReportSectionLint        = "Code Style Analysis"
	ReportSectionSecurity    = "Security Analysis"
	ReportSectionComplexity  = "Code Complexity Analysis"
	ReportSectionRecommend   = "Recommendations"
)

// Report Templates
const (
	ReportTitleQA           = "Fire Salamander QA Report"
	ReportSubtitleCoverage  = "Code Coverage Analysis"
	ReportSubtitleTests     = "Test Results"
	ReportSubtitleIssues    = "Code Quality Issues"
	ReportSubtitleSummary   = "Quality Summary"
)

// ========================================
// QA THRESHOLDS AND LIMITS
// ========================================

// Coverage Thresholds
const (
	CoverageThresholdHigh     = 90.0
	CoverageThresholdMedium   = 80.0
	CoverageThresholdLow      = 70.0
	CoverageThresholdMinimum  = 60.0
)

// Complexity Thresholds
const (
	ComplexityThresholdHigh   = 15
	ComplexityThresholdMedium = 10
	ComplexityThresholdLow    = 5
)

// Quality Score Thresholds
const (
	QualityScoreExcellent = 95.0
	QualityScoreGood      = 85.0
	QualityScoreAverage   = 75.0
	QualityScorePoor      = 60.0
	
	// Note: Score thresholds are defined in constants.go to avoid duplication
)

// ========================================
// QA FILE PATTERNS
// ========================================

// Test File Patterns
const (
	TestFilePattern     = "*_test.go"
	BenchFilePattern    = "*_bench.go"
	ExampleFilePattern  = "*_example.go"
)

// Source File Patterns
const (
	GoFilePattern       = "*.go"
	GoModFilePattern    = "go.mod"
	GoSumFilePattern    = "go.sum"
)

// Report File Patterns
const (
	CoverageProfileFile = "coverage.out"
	CoverageHTMLFile   = "coverage.html"
	QAReportJSONFile   = "qa_report.json"
	QAReportHTMLFile   = "qa_report.html"
	QAReportTextFile   = "qa_report.txt"
)

// ========================================
// QA ANALYSIS CATEGORIES
// ========================================

// Vet Categories
const (
	VetCategoryAsmdecl      = "asmdecl"
	VetCategoryAssign       = "assign"
	VetCategoryAtomic       = "atomic"
	VetCategoryBools        = "bools"
	VetCategoryBuildtag     = "buildtag"
	VetCategoryCgocall      = "cgocall"
	VetCategoryComposites   = "composites"
	VetCategoryCopylocks    = "copylocks"
	VetCategoryHttpresponse = "httpresponse"
	VetCategoryLoopclosure  = "loopclosure"
	VetCategoryLostcancel   = "lostcancel"
	VetCategoryNilfunc      = "nilfunc"
	VetCategoryPrintf       = "printf"
	VetCategoryShift        = "shift"
	VetCategoryStructtag    = "structtag"
	VetCategoryTests        = "tests"
	VetCategoryUnreachable  = "unreachable"
	VetCategoryUnsafeptr    = "unsafeptr"
	VetCategoryUnusedresult = "unusedresult"
)

// Lint Rules
const (
	LintRuleExportedFunc       = "exported-func"
	LintRuleExportedType       = "exported-type" 
	LintRuleExportedVar        = "exported-var"
	LintRulePackageComment     = "package-comment"
	LintRuleVarNaming          = "var-naming"
	LintRuleFuncNaming         = "func-naming"
	LintRuleTypeNaming         = "type-naming"
	LintRuleConstNaming        = "const-naming"
	LintRuleInterfaceNaming    = "interface-naming"
	LintRuleReceiverNaming     = "receiver-naming"
)

// ========================================
// QA TIME FORMATS AND DURATIONS
// ========================================

// Time Formats
const (
	QATimeFormatISO     = "2006-01-02T15:04:05Z07:00"
	QATimeFormatSimple  = "2006-01-02 15:04:05"
	QATimeFormatReport  = "January 2, 2006 at 15:04"
)

// Default Durations
const (
	DefaultTestTimeout    = "10m"
	DefaultVetTimeout     = "5m"
	DefaultLintTimeout    = "5m"
	DefaultBuildTimeout   = "5m"
)

// ========================================
// QA PATHS AND DIRECTORIES
// ========================================

// Default Paths
const (
	DefaultQAReportsDir    = "tests/reports/qa"
	DefaultCoverageDir     = "coverage"
	DefaultTestDataDir     = "testdata"
	DefaultBenchmarksDir   = "benchmarks"
)

// Config Paths
const (
	QAConfigFile        = ".qa-config.json"
	GolangCIConfigFile  = ".golangci.yml" 
	StaticcheckFile     = "staticcheck.conf"
)