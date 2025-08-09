package constants

// ========================================
// ALPHA-5 WEB SERVER CONSTANTS
// ========================================

// HTTP Methods
const (
	HTTPMethodGET     = "GET"
	HTTPMethodPOST    = "POST"
	HTTPMethodPUT     = "PUT"
	HTTPMethodDELETE  = "DELETE"
	HTTPMethodOPTIONS = "OPTIONS"
)

// Content Types
const (
	ContentTypeHTMLCharset   = "text/html; charset=utf-8"
	ContentTypePDF           = "application/pdf"
	ContentTypeCSV           = "text/csv"
	ContentTypeOctetStream   = "application/octet-stream"
)

// HTTP Headers
const (
	HeaderCacheControl        = "Cache-Control"
	HeaderXFrameOptions       = "X-Frame-Options"
	HeaderXContentType        = "X-Content-Type-Options"
	HeaderContentLength       = "Content-Length"
	HeaderCORSOrigin         = "Access-Control-Allow-Origin"
	HeaderCORSMethods        = "Access-Control-Allow-Methods"
	HeaderCORSHeaders        = "Access-Control-Allow-Headers"
	HeaderContentDisposition = "Content-Disposition"
)

// HTTP Header Values
const (
	HeaderValueNoCache      = "no-cache, no-store, must-revalidate"
	HeaderValueDeny         = "DENY"
	HeaderValueNoSniff      = "nosniff"
	HeaderValueCORSAll      = "*"
	HeaderValueCORSMethods  = "GET, POST, PUT, DELETE, OPTIONS"
	HeaderValueCORSHeaders  = "Content-Type, Authorization, X-Requested-With"
)

// Health Status Values
const (
	HealthStatusUnavailable = "unavailable"
	ServiceStatusRunning    = "running"
)

// API Routes (v1)
const (
	APIRouteAnalyze         = "/api/analyze"
	APIRouteAnalyzeSemantic = "/api/analyze/semantic"
	APIRouteAnalyzeSEO      = "/api/analyze/seo"
	APIRouteAnalyzeQuick    = "/api/analyze/quick"
	APIRouteHealth          = "/api/health"
	APIRouteStats           = "/api/stats"
	APIRouteAnalyses        = "/api/analyses"
	APIRouteAnalysisDetails = "/api/analysis/"
	APIRouteInfo            = "/api/info"
	APIRouteVersion         = "/api/version"
)

// Web Routes
const (
	WebRouteRoot     = "/"
	WebRouteHealth   = "/web/health"
	WebRouteDownload = "/web/download/"
	StaticRoute      = "/static/"
)

// File Extensions
const (
	ExtHTML = ".html"
	ExtPDF  = ".pdf"
	ExtJSON = ".json"
	ExtCSV  = ".csv"
)

// Static Files
const (
	StaticIndexHTML = "index.html"
	StaticDirectory = "static"
)

// ========================================
// WEB SERVER MESSAGES
// ========================================

// Error Messages
const (
	MsgErrorStaticFiles         = "Erreur accès fichiers statiques: %v"
	MsgWebInterfaceUnavailable = "Interface web non disponible"
	MsgMethodNotAllowed        = "Method not allowed"
	MsgAPIEndpointNotFound     = "API endpoint not found"
	MsgFilenameRequired        = "Nom de fichier requis"
	MsgWebServerError          = "Erreur serveur web: %v"
)

// Log Messages
const (
	MsgRoutesRegistered        = "Routes web enregistrées:"
	MsgRouteWebInterface      = "  GET  / - Interface web principale"
	MsgRouteStaticFiles       = "  GET  /static/* - Fichiers statiques"
	MsgRouteAPIProxy          = "  ALL  /api/* - API REST (proxy)"
	MsgRouteWebHealth         = "  GET  /web/health - Santé du serveur web"
	MsgRouteWebDownload       = "  GET  /web/download/* - Téléchargement de rapports"
	MsgWebInterfaceServed     = "Interface web servie: %s %s"
	MsgAPIRequest             = "API Request: %s %s"
	MsgDownloadRequested      = "Téléchargement demandé: %s"
	MsgWebServerStarting      = "Serveur web Fire Salamander en cours de démarrage..."
	MsgWebServerStarted       = "✅ Serveur web Fire Salamander démarré avec succès"
	MsgWebServerStopping      = "Arrêt du serveur web Fire Salamander"
	MsgErrorIndexHTML         = "Erreur lecture index.html: %v"
)

// Service Names
const (
	ServiceNameWebServer = "Fire Salamander Web Server"
)

// JSON Field Names (Web Server specific)
const (
	JSONFieldService         = "service"
	JSONFieldComponents      = "components"
	JSONFieldUptime          = "uptime"
	JSONFieldOrchestratorStats = "orchestrator_stats"
	JSONFieldWebServer       = "web_server"
	JSONFieldStaticFiles     = "static_files"
	JSONFieldAPIProxy        = "api_proxy"
	JSONFieldOrchestrator    = "orchestrator"
	JSONFieldStartedAt       = "started_at"
)

// Component Names
const (
	ComponentWebServer    = "web_server"
	ComponentStaticFiles  = "static_files"
	ComponentAPIProxy     = "api_proxy"
	ComponentOrchestrator = "orchestrator"
)

// Report Generation Constants
const (
	ReportTitlePrefix       = "Rapport Fire Salamander"
	ReportTitleSuffix       = "Rapport Fire Salamander - "
	ReportHTMLTitle         = "Rapport d'Analyse SEO"
	ReportExecutiveSummary  = "Résumé Exécutif"
	ReportDetailedMetrics   = "Métriques Détaillées"
	ReportPriorityRecs      = "Recommandations Prioritaires"
	ReportGeneratedBy       = "Rapport généré par Fire Salamander v"
)

// Report Metric Names
const (
	MetricSEOTechnical  = "SEO Technique:"
	MetricPerformance   = "Performance:"
	MetricContent       = "Contenu:"
	MetricMobile        = "Mobile:"
)

// Recommendation IDs
const (
	RecommendationOptimizeImages = "optimize_images"
	RecommendationMetaDesc       = "meta_descriptions"
)

// Recommendation Titles
const (
	RecommendationTitleImages   = "Optimiser les images"
	RecommendationTitleMetaDesc = "Améliorer les méta-descriptions"
)

// Recommendation Descriptions
const (
	RecommendationDescImages   = "Compresser les images pour améliorer les temps de chargement"
	RecommendationDescMetaDesc = "Ajouter des méta-descriptions optimisées sur les pages manquantes"
)

// Priority/Impact/Effort Values
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
	ImpactHigh     = "high"
	ImpactMedium   = "medium"
	ImpactLow      = "low"
	EffortHigh     = "high"
	EffortMedium   = "medium"
	EffortLow      = "low"
)

// Insight Types
const (
	InsightTypeContentSEO = "content_seo_alignment"
	InsightSeverityInfo   = "info"
)

// Insight Messages
const (
	InsightTitleContentSEO = "Alignement contenu-SEO détecté"
	InsightDescContentSEO  = "Le titre de la page est cohérent avec les mots-clés identifiés"
)

// CSV Headers and Values
const (
	CSVHeaderURL           = "URL"
	CSVHeaderScoreGlobal   = "Score Global"
	CSVHeaderSEOTechnique  = "SEO Technique"
	CSVHeaderPerformance   = "Performance"
	CSVHeaderContent       = "Contenu"
	CSVHeaderMobile        = "Mobile"
	CSVHeaderStatus        = "Statut"
	CSVHeaderDate          = "Date"
	CSVValueSuccess        = "Succès"
	CSVSummaryTitle        = "Résumé:"
	CSVAverageScore        = "Score moyen"
	CSVSuccessfulAnalyses  = "Analyses réussies"
	CSVFailedAnalyses      = "Analyses échouées"
	CSVLastUpdate          = "Dernière mise à jour"
	CSVMainRecommendations = "Recommandations principales:"
)

// HTML Template Constants
const (
	HTMLDoctype     = "<!DOCTYPE html>"
	HTMLLangFr      = `<html lang="fr">`
	HTMLMetaCharset = `<meta charset="UTF-8">`
	HTMLMetaViewport = `<meta name="viewport" content="width=device-width, initial-scale=1.0">`
)

// CSS Style Values
const (
	CSSFontFamily = "font-family: Arial, sans-serif"
	CSSMaxWidth   = "max-width: 800px"
	CSSMarginAuto = "margin: 0 auto"
	CSSPadding20  = "padding: 20px"
	CSSColorOrange = "color: #ff6b35"
	CSSTextCenter  = "text-align: center"
)

// Time Format Constants
const (
	TimeFormatFrench    = "02/01/2006 à 15:04"
	TimeFormatFrenchDate = "02/01/2006"
	TimeFormatFrenchDateTime = "02/01/2006 15:04"
)

// Analysis Duration and Metadata
const (
	SampleAnalysisDuration = "15.2s"
	SamplePagesAnalyzed    = 1
	SampleIssuesFound      = 3
	SampleScoreTrend       = "stable"
)