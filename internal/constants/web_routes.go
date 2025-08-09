package constants

// 🚨 FIRE SALAMANDER - WEB ROUTES Constants
// Zero Hardcoding Policy - All web routes and paths

// ===== API ROUTES =====
const (
	RouteAPI                 = "/api"
	RouteAPIAnalyze          = "/api/analyze" 
	RouteAPIStatus           = "/api/status"
	RouteAPIResults          = "/api/results"
	RouteAPIHealth           = "/api/health"
	RouteAPIStats            = "/api/stats"
	RouteAPIV1               = "/api/v1"
	RouteAPIV1Health         = "/api/v1/health"
	RouteAPIV1Stats          = "/api/v1/stats"
	RouteAPIV1Analyze        = "/api/v1/analyze"
	RouteAPIV1AnalyzeQuick   = "/api/v1/analyze/quick"
	RouteAPIV1AnalyzeSEO     = "/api/v1/analyze/seo"
	RouteAPIV1AnalyzeFull    = "/api/v1/analyze/full"
)

// ===== WEB ROUTES =====
const (
	RouteWeb                 = "/web"
	RouteWebHealth           = "/web/health"
	RouteWebDownload         = "/web/download"
	RouteWebStatic           = "/web/static"
	RouteWebAssets           = "/web/assets"
	RouteWebUploads          = "/web/uploads"
)

// ===== PAGE ROUTES =====
const (
	RouteHome                = "/"
	RouteAnalyze             = "/analyze"
	RouteResults             = "/results"
	RouteStatus              = "/status"
	RouteReport              = "/report"
	RouteHistory             = "/history"
	// RouteDashboard removed - Fire Salamander only
	RouteAbout               = "/about"
	RouteContact             = "/contact"
	RoutePrivacy             = "/privacy"
	RouteTerms               = "/terms"
)

// ===== STATIC ROUTES =====
const (
	RouteStatic              = "/static"
	RouteAssets              = "/assets"
	RouteUploads             = "/uploads"
	RouteDownloads           = "/downloads"
	RouteImages              = "/images"
	RouteCSS                 = "/css"
	RouteJS                  = "/js"
	RouteFonts               = "/fonts"
)

// ===== ADMIN ROUTES =====
const (
	RouteAdmin               = "/admin"
	// RouteAdminDashboard removed - Fire Salamander only
	RouteAdminUsers          = "/admin/users"
	RouteAdminSettings       = "/admin/settings"
	RouteAdminLogs           = "/admin/logs"
	RouteAdminStats          = "/admin/stats"
)

// ===== AUTHENTICATION ROUTES =====
const (
	RouteAuth                = "/auth"
	RouteLogin               = "/auth/login"
	RouteLogout              = "/auth/logout"
	RouteRegister            = "/auth/register"
	RouteForgotPassword      = "/auth/forgot"
	RouteResetPassword       = "/auth/reset"
)

// ===== RESTRICTED PATHS =====
const (
	RoutePrivate             = "/private"
	RouteInternal            = "/internal"
	RouteSystem              = "/system"
	RouteConfig              = "/config"
)

// ===== PATH PARAMETERS =====
const (
	ParamID                  = "/:id"
	ParamAnalysisID          = "/:analysisId"
	ParamUserID              = "/:userId"
	ParamReportID            = "/:reportId"
)

// ===== QUERY PARAMETERS =====
const (
	QueryParamURL            = "?url="
	QueryParamID             = "?id="
	QueryParamPage           = "?page="
	QueryParamLimit          = "?limit="
	QueryParamFormat         = "?format="
	QueryParamType           = "?type="
)

// ===== JAVASCRIPT ROUTES (for templates) =====
const (
	JSRouteAPIAnalyze        = "'/api/analyze'"
	JSRouteAnalyzeWithID     = "'/analyze?id=' + data.id"
	JSRouteAPIStatus         = "'/api/status/'"
)