package constants

// ========================================
// BRAVO-3 SECURITY AGENT CONSTANTS
// ========================================

// Security Agent Identity
const (
	SecurityAgentName = "SECURITY-AGENT"
)

// ========================================
// SECURITY JSON FIELD NAMES
// ========================================

// Security Configuration Fields
const (
	SecurityJSONFieldOWASPTop10      = "owasp_top10"
	SecurityJSONFieldDependencyCheck = "dependency_check"
	SecurityJSONFieldSecretScanning  = "secret_scanning"
	SecurityJSONFieldSQLInjection    = "sql_injection"
	SecurityJSONFieldXSSCheck        = "xss_check"
	SecurityJSONFieldCSRFCheck       = "csrf_check"
	SecurityJSONFieldReportPath      = "report_path"
)

// Security Results Fields
const (
	SecurityJSONFieldTimestamp        = "timestamp"
	SecurityJSONFieldOWASPScore       = "owasp_score"
	SecurityJSONFieldVulnerabilities  = "vulnerabilities"
	SecurityJSONFieldDependencyIssues = "dependency_issues"
	SecurityJSONFieldSecretFindings   = "secret_findings"
	SecurityJSONFieldSecurityHeaders  = "security_headers"
	SecurityJSONFieldOverallRisk      = "overall_risk"
	SecurityJSONFieldPassed           = "passed"
)

// Vulnerability Fields
const (
	SecurityJSONFieldID          = "id"
	SecurityJSONFieldTitle       = "title"
	SecurityJSONFieldSeverity    = "severity"
	SecurityJSONFieldCategory    = "category"
	SecurityJSONFieldDescription = "description"
	SecurityJSONFieldFile        = "file"
	SecurityJSONFieldLine        = "line"
	SecurityJSONFieldCWE         = "cwe"
	SecurityJSONFieldCVSS        = "cvss"
)

// Dependency Issue Fields
const (
	SecurityJSONFieldPackage    = "package"
	SecurityJSONFieldVersion    = "version"
	SecurityJSONFieldCVE        = "cve"
	SecurityJSONFieldFixVersion = "fix_version"
)

// Secret Finding Fields
const (
	SecurityJSONFieldType       = "type"
	SecurityJSONFieldPattern    = "pattern"
	SecurityJSONFieldConfidence = "confidence"
	SecurityJSONFieldEntropy    = "entropy"
)

// Security Headers Fields
const (
	SecurityJSONFieldHSTS                 = "hsts"
	SecurityJSONFieldCSP                  = "csp"
	SecurityJSONFieldXFrameOptions        = "x_frame_options"
	SecurityJSONFieldXContentTypeOptions  = "x_content_type_options"
	SecurityJSONFieldXSSProtection        = "xss_protection"
	SecurityJSONFieldReferrerPolicy       = "referrer_policy"
	SecurityJSONFieldScore                = "score"
)

// ========================================
// SECURITY SEVERITY LEVELS
// ========================================

// Security Severity Values
const (
	SecuritySeverityCritical = "CRITICAL"
	SecuritySeverityHigh     = "HIGH"
	SecuritySeverityMedium   = "MEDIUM"
	SecuritySeverityLow      = "LOW"
)

// ========================================
// SECURITY RISK LEVELS
// ========================================

// Overall Risk Levels
const (
	SecurityRiskLow      = "LOW"
	SecurityRiskMedium   = "MEDIUM"
	SecurityRiskHigh     = "HIGH"
	SecurityRiskCritical = "CRITICAL"
)

// ========================================
// SECURITY CATEGORIES
// ========================================

// Security Analysis Categories
const (
	SecurityCategoryOWASP          = "OWASP"
	SecurityCategoryStaticAnalysis = "Static Analysis"
)

// OWASP Top 10 Categories
const (
	SecurityOWASPInjection              = "A01:2021 - Injection"
	SecurityOWASPBrokenAuth             = "A02:2021 - Broken Authentication"
	SecurityOWASPDataExposure           = "A03:2021 - Sensitive Data Exposure"
	SecurityOWASPXXE                    = "A04:2021 - XML External Entities"
	SecurityOWASPAccessControl          = "A05:2021 - Broken Access Control"
	SecurityOWASPMisconfiguration       = "A06:2021 - Security Misconfiguration"
	SecurityOWASPXSS                    = "A07:2021 - Cross-Site Scripting"
	SecurityOWASPDeserialization        = "A08:2021 - Insecure Deserialization"
	SecurityOWASPVulnerableComponents   = "A09:2021 - Using Components with Known Vulnerabilities"
	SecurityOWASPLogging                = "A10:2021 - Insufficient Logging & Monitoring"
)

// ========================================
// SECURITY SECRET TYPES
// ========================================

// Secret Detection Types
const (
	SecuritySecretTypeAPIKey     = "API_KEY"
	SecuritySecretTypePassword   = "PASSWORD"
	SecuritySecretTypeJWTToken   = "JWT_TOKEN"
	SecuritySecretTypeAWSKey     = "AWS_KEY"
	SecuritySecretTypePrivateKey = "PRIVATE_KEY"
)

// ========================================
// SECURITY TOOLS
// ========================================

// Security Analysis Tools
const (
	SecurityToolGosec        = "gosec"
	SecurityToolNancy        = "nancy"
	SecurityToolGovulncheck  = "govulncheck"
	SecurityToolTruffleHog   = "truffleHog"
	SecurityToolGitleaks     = "gitleaks"
)

// ========================================
// SECURITY TEST TYPES
// ========================================

// OWASP Test Types
const (
	SecurityTestInjection        = "injection"
	SecurityTestBrokenAuth       = "broken_authentication"
	SecurityTestDataExposure     = "sensitive_data_exposure"
	SecurityTestXXE              = "xml_external_entities"
	SecurityTestAccessControl    = "broken_access_control"
	SecurityTestMisconfiguration = "security_misconfiguration"
	SecurityTestXSS              = "cross_site_scripting"
	SecurityTestDeserialization  = "insecure_deserialization"
	SecurityTestVulnComponents   = "vulnerable_components"
	SecurityTestLogging          = "insufficient_logging"
)

// ========================================
// SECURITY FILES AND DIRECTORIES
// ========================================

// Security Report Files
const (
	SecurityFileGosecReport    = "gosec-report.json"
	SecurityFileSecurityReport = "security_report.json"
)

// Directory Patterns to Skip
const (
	SecurityDirVendor      = "vendor/"
	SecurityDirGit         = ".git/"
	SecurityDirNodeModules = "node_modules/"
)

// Default Paths
const (
	SecurityDefaultReportPath = "tests/reports/security"
)

// ========================================
// SECURITY FILE EXTENSIONS
// ========================================

// File Extensions for Scanning
const (
	SecurityExtGo   = ".go"
	SecurityExtYAML = ".yaml"
	SecurityExtYML  = ".yml"
	SecurityExtEnv  = ".env"
)

// ========================================
// SECURITY COMMANDS
// ========================================

// Go Install Commands
const (
	SecurityCmdInstallGosec      = "github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
	SecurityCmdInstallGovulncheck = "golang.org/x/vuln/cmd/govulncheck@latest"
)

// Go Commands
const (
	SecurityCmdGo        = "go"
	SecurityCmdInstall   = "install"
	SecurityCmdGosec     = "gosec"
	SecurityCmdNancy     = "nancy"
	SecurityCmdGovulncheck = "govulncheck"
)

// Command Arguments
const (
	SecurityArgFormat = "-fmt"
	SecurityArgOut    = "-out"
	SecurityArgJSON   = "json"
	SecurityArgSleuth = "sleuth"
	SecurityArgAll    = "./..."
)

// ========================================
// SECURITY SCORING THRESHOLDS
// ========================================

// Base Security Score
const (
	SecurityBaseScore = 100.0
)

// Vulnerability Score Deductions
const (
	SecurityScoreCriticalVuln = 25.0
	SecurityScoreHighVuln     = 15.0
	SecurityScoreMediumVuln   = 8.0
	SecurityScoreLowVuln      = 3.0
)

// Dependency Issue Score Deductions
const (
	SecurityScoreCriticalDep = 20.0
	SecurityScoreHighDep     = 10.0
	SecurityScoreMediumDep   = 5.0
	SecurityScoreLowDep      = 2.0
)

// Other Security Factors
const (
	SecurityScoreSecretPenalty = 10.0
	SecurityScoreHeadersWeight = 70.0 // Security headers score weight
)

// Risk Level Thresholds
const (
	SecurityThresholdLowRisk      = 90.0
	SecurityThresholdMediumRisk   = 70.0
	SecurityThresholdHighRisk     = 50.0
	SecurityThresholdCriticalRisk = 0.0
)

// ========================================
// SECURITY MESSAGES
// ========================================

// Log Messages
const (
	SecurityMsgStartingScan           = "🔒 Starting OWASP security scan"
	SecurityMsgSecurityScanCompleted  = "Security scan completed"
	SecurityMsgRunningStaticAnalysis  = "Running static code analysis with gosec"
	SecurityMsgCheckingDependencies   = "Checking dependencies for vulnerabilities"
	SecurityMsgScanningSecrets        = "Scanning for hardcoded secrets"
	SecurityMsgRunningDynamicAnalysis = "Running dynamic security analysis"
	SecurityMsgCheckingSecurityHeaders = "Checking HTTP security headers"
)

// Error Messages
const (
	SecurityMsgStaticAnalysisFailed   = "Static analysis failed"
	SecurityMsgDependencyCheckFailed  = "Dependency check failed"
	SecurityMsgSecretScanningFailed   = "Secret scanning failed"
	SecurityMsgDynamicAnalysisFailed  = "Dynamic analysis failed"
	SecurityMsgSecurityHeadersFailed  = "Security headers check failed"
	SecurityMsgReportGenerationFailed = "Failed to generate security report"
)

// Tool Messages
const (
	SecurityMsgGosecNotFound       = "gosec not found, installing..."
	SecurityMsgNancyNotFound       = "nancy not found, using govulncheck instead"
	SecurityMsgFailedInstallGosec  = "failed to install gosec"
	SecurityMsgFailedInstallGovuln = "Failed to install govulncheck"
)

// Analysis Messages
const (
	SecurityMsgGosecFoundIssues       = "gosec found security issues"
	SecurityMsgNancyFoundVulns        = "nancy found vulnerable dependencies"
	SecurityMsgGovulncheckFoundVulns  = "govulncheck found vulnerabilities"
	SecurityMsgFailedParseGosec       = "Failed to parse gosec results"
)

// ========================================
// SECURITY PATTERNS AND EXAMPLES
// ========================================

// Regex Patterns for Secret Detection
const (
	SecurityPatternAPIKey     = `(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?([a-zA-Z0-9]{20,})['"]?`
	SecurityPatternPassword   = `(?i)(password|passwd|pwd)\s*[:=]\s*['"]?([^'"\s]{8,})['"]?`
	SecurityPatternJWTToken   = `eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`
	SecurityPatternAWSKey     = `AKIA[0-9A-Z]{16}`
	SecurityPatternPrivateKey = `-----BEGIN [A-Z ]+PRIVATE KEY-----`
)

// Example Values (for testing/demo)
const (
	SecurityExamplePackage        = "example/vulnerable-package"
	SecurityExampleGitHubPackage  = "github.com/example/vulnerable"
	SecurityExampleVersion        = "v1.0.0"
	SecurityExampleFixVersion     = "v1.0.1"
	SecurityExampleCVE            = "CVE-2023-12345"
	SecurityExampleCVEDemo        = "CVE-2023-99999"
)

// Example Vulnerability Data
const (
	SecurityExampleVulnID          = "SEC-001"
	SecurityExampleVulnTitle       = "Missing Security Headers"
	SecurityExampleVulnDescription = "Application is missing recommended security headers"
	SecurityExampleCWE             = "CWE-16"
	SecurityExampleCVSS            = 3.1
)