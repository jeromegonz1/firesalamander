// DELTA-9 RAMBO CONSTANTS - REPORT MODULE
// ========================================
// "I'm your worst nightmare!"
// "To survive a war, you gotta become war!"
// 
// MISSION: ELIMINATE ALL HARDCODED VALUES IN REPORTS.GO
// STATUS: MAXIMUM CARNAGE AUTHORIZED
// 
// This file contains all constants required for the Fire Salamander
// reporting system. Every hardcoded value has been OBLITERATED and
// replaced with these RAMBO-GRADE constants!

package constants

// ============================================================================
// REPORT FORMATS - The Arsenal of Output Destruction
// ============================================================================
const (
	// Report output formats - Choose your weapon!
	ReportFormatHTML = "html"
	ReportFormatJSON = "json"
	ReportFormatPDF  = "pdf"
	ReportFormatCSV  = "csv"
)

// ============================================================================
// REPORT TYPES - Categories of Digital Warfare
// ============================================================================
const (
	// Report types for different mission objectives
	ReportTypeExecutive  = "executive"
	ReportTypeDetailed   = "detailed"
	ReportTypeTechnical  = "technical"
	ReportTypeComparison = "comparison"
)

// ============================================================================
// JSON FIELD NAMES - The Communication Protocol
// ============================================================================
const (
	// Main configuration fields (unique to reports)
	JSONFieldFormat         = "format"
	JSONFieldIncludeSummary = "include_summary"
	JSONFieldIncludeDetails = "include_details"
	JSONFieldIncludeCharts  = "include_charts"
	JSONFieldIncludeRawData = "include_raw_data"
	JSONFieldCustomSections = "custom_sections"
	JSONFieldBranding      = "branding"

	// Branding fields
	JSONFieldCompanyName = "company_name"
	JSONFieldLogo       = "logo"
	JSONFieldPrimary    = "primary"
	JSONFieldSecondary  = "secondary"
	JSONFieldColors     = "colors"

	// Report structure fields
	JSONFieldID          = "id"
	JSONFieldTitle       = "title"
	JSONFieldGeneratedAt = "generated_at"
	JSONFieldContent     = "content"
	JSONFieldMetadata    = "metadata"
	JSONFieldSize        = "size"

	// Analysis fields (unique to reports)
	JSONFieldURL         = "url"
	JSONFieldDomain      = "domain"
	JSONFieldAnalyzedAt  = "analyzed_at"
	JSONFieldGrade       = "grade"
	JSONFieldTopIssues   = "top_issues"
	JSONFieldTopRecommendations = "top_recommendations"
	JSONFieldKeyMetrics  = "key_metrics"
	JSONFieldPerformance = "performance"
)

// ============================================================================
// STATUS VALUES - Mission Status Indicators
// ============================================================================
// Note: StatusWarning and StatusError are defined in data_integrity_constants.go

// ============================================================================
// COLOR CODES - RAMBO's Paint Job (HEX Values)
// ============================================================================
const (
	// SEPTEO Fire Salamander color palette
	ColorSuccessGreen  = "#28a745"
	ColorWarningYellow = "#ffc107"
	ColorDangerOrange  = "#fd7e14"
	ColorErrorRed      = "#dc3545"
	ColorFireOrange    = "#ff6b35"
	ColorFireYellow    = "#f7931e"
	
	// Additional colors for reports
	ColorPrimary      = "#0066cc"
	ColorSecondary    = "#6c757d"
	ColorInfo         = "#17a2b8"
	ColorLight        = "#f8f9fa"
	ColorDark         = "#343a40"
)

// ============================================================================
// CSS CLASSES - The Style Warriors
// ============================================================================
const (
	// Grade-related classes
	CSSClassGrade            = "grade"
	CSSClassGradeA           = "grade-a"
	CSSClassGradeB           = "grade-b"
	CSSClassGradeC           = "grade-c"
	CSSClassGradeD           = "grade-d"
	CSSClassGradeF           = "grade-f"
	
	// Performance classes
	CSSClassPerformance      = "performance"
	CSSClassPerformanceHigh  = "performance-high"
	CSSClassPerformanceMed   = "performance-medium"
	CSSClassPerformanceLow   = "performance-low"
	
	// Status classes
	CSSClassStatus           = "status"
	CSSClassStatusSuccess    = "status-success"
	CSSClassStatusWarning    = "status-warning"
	CSSClassStatusError      = "status-error"
	CSSClassStatusInfo       = "status-info"
	
	// Report structure classes
	CSSClassReportHeader     = "report-header"
	CSSClassReportBody       = "report-body"
	CSSClassReportFooter     = "report-footer"
	CSSClassReportSection    = "report-section"
	CSSClassReportChart      = "report-chart"
	CSSClassReportTable      = "report-table"
	CSSClassReportSummary    = "report-summary"
	CSSClassReportDetails    = "report-details"
	
	// Branding classes
	CSSClassBranding         = "branding"
	CSSClassLogo            = "logo"
	CSSClassCompanyName     = "company-name"
	
	// Metric classes
	CSSClassMetric          = "metric"
	CSSClassMetricValue     = "metric-value"
	CSSClassMetricLabel     = "metric-label"
	CSSClassScoreDisplay    = "score-display"
)

// ============================================================================
// GRADE VALUES - The Judgment System
// ============================================================================
const (
	GradeAPlus  = "A+"
	GradeA      = "A"
	GradeAMinus = "A-"
	GradeBPlus  = "B+"
	GradeB      = "B"
	GradeBMinus = "B-"
	GradeCPlus  = "C+"
	GradeC      = "C"
	GradeCMinus = "C-"
	GradeDPlus  = "D+"
	GradeD      = "D"
	GradeDMinus = "D-"
	GradeF      = "F"
)

// ============================================================================
// DATE/TIME FORMATS - Time is a Weapon
// ============================================================================
const (
	// Go time format constants (Remember: Mon Jan 2 15:04:05 MST 2006 is Unix time 1136239445)
	DateTimeFormat     = "2006-01-02 15:04:05"
	DateOnlyFormat     = "2006-01-02"
	TimeOnlyFormat     = "15:04:05"
	TimestampFormat    = "2006-01-02T15:04:05Z"
	DisplayDateFormat  = "January 2, 2006"
	DisplayTimeFormat  = "3:04 PM"
	ReportIDDateFormat = "20060102_150405"
)

// ============================================================================
// TEMPLATE FUNCTION NAMES - The Most Common Weapons
// ============================================================================
const (
	// Template function identifiers for reports
	TemplateFuncFormat      = "Format"
	TemplateFuncGrade       = "Grade"
	TemplateFuncBranding    = "Branding"
	TemplateFuncScore       = "score"
	TemplateFuncPerformance = "performance"
	TemplateFuncTitle       = "title"
	TemplateFuncURL         = "url"
	TemplateFuncDomain      = "domain"
	TemplateFuncSize        = "size"
	TemplateFuncContent     = "content"
	TemplateFuncMetadata    = "metadata"
)

// ============================================================================
// HTTP STATUS TEMPLATES - Response Warriors
// ============================================================================
const (
	HTTPStatusOK           = 200
	HTTPStatusBadRequest   = 400
	HTTPStatusNotFound     = 404
	HTTPStatusServerError  = 500
)

// ============================================================================
// REPORT SECTIONS - The Battle Zones
// ============================================================================
const (
	SectionExecutiveSummary = "executive_summary"
	SectionTechnicalDetails = "technical_details"
	SectionPerformanceMetrics = "performance_metrics"
	SectionSEOAnalysis      = "seo_analysis"
	SectionRecommendations  = "recommendations"
	SectionRawData         = "raw_data"
	SectionCharts          = "charts"
	SectionAppendix        = "appendix"
)

// ============================================================================
// METRIC TYPES - The Measurement Arsenal
// ============================================================================
const (
	MetricTypeScore        = "score"
	MetricTypePercentage   = "percentage"
	MetricTypeCount        = "count"
	MetricTypeTime         = "time"
	MetricTypeSize         = "size"
	MetricTypeRatio        = "ratio"
)

// ============================================================================
// EMOJIS & ICONS - Visual Warfare Symbols
// ============================================================================
const (
	// Grade emojis
	EmojiGradeA      = "🅰️"
	EmojiGradeB      = "🅱️"
	EmojiGradeC      = "🔤"
	EmojiGradeD      = "📉"
	EmojiGradeF      = "❌"
	
	// Status emojis
	EmojiSuccess     = "✅"
	EmojiWarning     = "⚠️"
	EmojiError       = "🔴"
	EmojiInfo        = "ℹ️"
	EmojiFireSalamander = "🦎"
	EmojiFire        = "🔥"
	
	// Report type emojis
	EmojiExecutive   = "👔"
	EmojiTechnical   = "🔧"
	EmojiDetailed    = "📊"
	EmojiComparison  = "⚖️"
	
	// Analysis emojis
	EmojiChart       = "📈"
	EmojiMetrics     = "📏"
	EmojiRecommendation = "💡"
	EmojiSEO         = "🔍"
	EmojiPerformance = "⚡"
	EmojiSecurity    = "🛡️"
	
	// Rambo warfare emojis
	EmojiRambo       = "🎖️"
	EmojiMission     = "🎯"
	EmojiWeapon      = "🔫"
	EmojiExplosive   = "💥"
	EmojiVictory     = "🏆"
)

// ============================================================================
// RAMBO QUOTES - Motivational Warfare
// ============================================================================
var RamboQuotes = []string{
	"I'm your worst nightmare!",
	"To survive a war, you gotta become war!",
	"When you're pushed, killing's as easy as breathing!",
	"They drew first blood, not me!",
	"This is what we do, who we are. Live for nothing, or die for something!",
	"I could have killed 'em all, I could kill you!",
	"Nothing is over!",
}

// ============================================================================
// REPORT GENERATION SETTINGS - Mission Parameters
// ============================================================================
const (
	// Default values for report generation
	DefaultReportFormat         = ReportFormatHTML
	DefaultReportType          = ReportTypeDetailed
	DefaultIncludeSummary      = true
	DefaultIncludeDetails      = true
	DefaultIncludeCharts       = true
	DefaultIncludeRawData      = false
	MaxReportSizeBytes         = 10485760 // 10MB
	MaxReportGenerationTimeMS  = 30000    // 30 seconds
	ReportCacheExpirationHours = 24
)

// ============================================================================
// MISSION COMPLETE - DELTA-9 RAMBO CONSTANTS DEPLOYED
// All hardcoded values have been OBLITERATED!
// "First blood is drawn. Now it's war!"
// ============================================================================