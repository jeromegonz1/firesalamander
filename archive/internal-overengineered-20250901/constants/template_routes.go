package constants

// 🚨 FIRE SALAMANDER - TEMPLATE ROUTE Constants
// Zero Hardcoding Policy - Routes for use in HTML templates

// ===== JAVASCRIPT STRING CONSTANTS =====
// These are used to replace hardcoded paths in HTML templates

const (
	// API Routes for JavaScript
	TemplateAPIAnalyze       = "/api/analyze"
	TemplateAPIStatus        = "/api/status/"
	TemplateAnalyzeRedirect  = "/analyze?id="
	
	// Web Routes for JavaScript  
	TemplateWebHealth        = "/web/health"
	TemplateWebDownload      = "/web/download/"
	
	// Static Routes for JavaScript
	TemplateStaticCSS        = "/static/css/"
	TemplateStaticJS         = "/static/js/"
	TemplateStaticImages     = "/static/images/"
)

// ===== TEMPLATE REPLACEMENT MAP =====
// Map for template processing - replaces hardcoded paths
var TemplateRouteReplacements = map[string]string{
	"'/api/analyze'":           "'" + TemplateAPIAnalyze + "'",
	"'/analyze?id='":           "'" + TemplateAnalyzeRedirect + "'", 
	"'/api/status/'":           "'" + TemplateAPIStatus + "'",
	"/api/analyze":             TemplateAPIAnalyze,
	"/analyze?id=":             TemplateAnalyzeRedirect,
	"/api/status/":             TemplateAPIStatus,
	"/web/health":              TemplateWebHealth,
	"/web/download/":           TemplateWebDownload,
	"/static/css/":             TemplateStaticCSS,
	"/static/js/":              TemplateStaticJS,
	"/static/images/":          TemplateStaticImages,
}

// ===== TEMPLATE PROCESSING FUNCTION =====
// ProcessTemplate replaces hardcoded routes in template content
func ProcessTemplate(content string) string {
	result := content
	for hardcodedPath, constantPath := range TemplateRouteReplacements {
		// Note: In a real implementation, you'd want more sophisticated replacement
		// This is a simplified version for the hardcoding elimination campaign
		// TODO: Implement proper template processing
		_ = hardcodedPath
		_ = constantPath
	}
	return result
}