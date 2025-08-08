package constants

// ========================================
// APPLICATION MESSAGES
// ========================================

const (
	// Startup Messages
	AppStarting          = "🔥 Fire Salamander starting"
	AppReady             = "✅ Fire Salamander ready"
	AppShuttingDown      = "🛑 Fire Salamander shutting down"
	AppStopped           = "👋 Fire Salamander stopped"
	
	// Server Messages
	ServerStarting       = "🔥 Démarrage du serveur web Fire Salamander"
	ServerReady          = "✅ Serveur web démarré avec succès"
	ServerStopping       = "🛑 Arrêt du serveur web"
	ServerStopped        = "✅ Serveur web arrêté"
	
	// API Messages
	APIStarting          = "🔌 Démarrage de l'API Fire Salamander"
	APIReady             = "✅ API Fire Salamander démarrée"
	APIStopping          = "🛑 Arrêt de l'API"
	APIStopped           = "✅ API arrêtée"
	
	// Analysis Messages
	AnalysisStarting     = "🚀 Démarrage de l'analyse"
	AnalysisInProgress   = "⏳ Analyse en cours"
	AnalysisComplete     = "✅ Analyse terminée"
	AnalysisFailed       = "❌ Analyse échouée"
	
	// Phase Messages
	PhaseDiscoveryMsg    = "🔍 Découverte du site"
	PhaseSEOAnalysisMsg  = "📊 Analyse SEO technique"
	PhaseAIAnalysisMsg   = "🤖 Analyse sémantique IA"
	PhaseReportGenMsg    = "📋 Génération du rapport"
)

// ========================================
// CLI FLAG DESCRIPTIONS
// ========================================

const (
	ConfigPathDescription = "Chemin vers le fichier de configuration"
	PortDescription       = "Port du serveur web"
	ShowVersionDescription = "Afficher la version"
	WebOnlyDescription    = "Lancer uniquement l'interface web (sans orchestrateur)"
	APIOnlyDescription    = "Lancer uniquement l'API REST (sans interface web)"
	VerboseDescription    = "Activer les logs détaillés"
)

// ========================================
// LOG FORMAT STRINGS
// ========================================

const (
	LogAPIAvailableFormat     = "📡 API disponible sur: http://localhost:%d/api/v1"
	LogInterfaceAvailableFormat = "🔥 Interface Fire Salamander: http://localhost:%d"
	LogAPIIntegratedFormat    = "📡 API REST intégrée: http://localhost:%d/api/v1"
	WebInterfaceAvailableFormat = "📡 Interface web disponible sur: http://localhost:%d"
	APIRestAvailableFormat    = "🔌 API REST disponible sur: http://localhost:%d/api/v1"
	
	// Display Format Strings
	InterfaceWebFormat    = "🌐 Interface Web: http://localhost:%d\n"
	APIRESTFormat         = "📡 API REST: http://localhost:%d/api/v1\n"
	DocInterfaceFormat    = "   - Interface: http://localhost:%d\n"
	DocAPIFormat          = "   - API: http://localhost:%d/api/v1/info\n"
	DocHealthFormat       = "   - Santé: http://localhost:%d/api/v1/health\n"
	
	// Server Status Formats
	ServerStartedFormat   = "🔥 Fire Salamander démarré sur http://%s"
	APIDocAvailableFormat = "API Documentation disponible sur: http://localhost:%d/"
)

// ========================================
// ERROR MESSAGES
// ========================================

const (
	// HTTP Errors
	ErrMethodNotAllowed   = "Method not allowed"
	ErrInvalidJSON        = "Invalid JSON"
	ErrInvalidURL         = "Invalid URL format"
	ErrURLRequired        = "URL is required"
	ErrAnalysisNotFound   = "Analysis not found"
	ErrInvalidAnalysisID  = "Invalid analysis ID"
	ErrRequestFailed      = "Request failed"
	ErrServerError        = "Server error"
	
	// Configuration Errors
	ErrConfigNotFound     = "Configuration file not found"
	ErrConfigInvalid      = "Invalid configuration"
	ErrPortInUse          = "Port already in use"
	
	// Analysis Errors
	ErrAnalysisFailed     = "Analysis failed"
	ErrTimeoutExceeded    = "Timeout exceeded"
	ErrResourceNotFound   = "Resource not found"
	ErrPermissionDenied   = "Permission denied"
	
	// Generic Errors
	ErrInternalError      = "Internal server error"
	ErrBadRequest         = "Bad request"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrNotFound           = "Not found"
)

// ========================================
// SUCCESS MESSAGES
// ========================================

const (
	MsgAnalysisStarted    = "Analysis started successfully"
	MsgAnalysisComplete   = "Analysis completed successfully"
	MsgConfigLoaded       = "Configuration loaded successfully"
	MsgServerStarted      = "Server started successfully"
	MsgResourceCreated    = "Resource created successfully"
	MsgResourceUpdated    = "Resource updated successfully"
	MsgResourceDeleted    = "Resource deleted successfully"
)

// ========================================
// CURL EXAMPLES
// ========================================

const (
	CurlExampleFormat     = "curl -X POST http://localhost:%d/api/v1/analyze/quick \\\n"
	CurlHeaders           = "     -H \"Content-Type: application/json\" \\"
	CurlExampleData       = "     -d '{\"url\": \"https://example.com\"}'"
	APIExampleFormat      = "curl -X POST http://localhost:%d/api/v1/analyze/quick \\\n  -H \"Content-Type: application/json\" \\\n  -d '{\"url\": \"https://example.com\"}'"
)

// ========================================
// UI MESSAGES
// ========================================

const (
	WelcomeMessage        = "Bienvenue dans Fire Salamander"
	AnalyzingMessage      = "Analyse en cours..."
	ResultsReadyMessage   = "Résultats prêts"
	NoResultsMessage      = "Aucun résultat trouvé"
	LoadingMessage        = "Chargement..."
	ErrorMessage          = "Une erreur est survenue"
	TryAgainMessage       = "Veuillez réessayer"
	SuccessMessage        = "Opération réussie"
	UIPlaceholderURL      = "https://your-site.com"
)

// ========================================
// VALIDATION MESSAGES
// ========================================

const (
	ValidURLRequired      = "Une URL valide est requise"
	ValidEmailRequired    = "Une adresse email valide est requise"
	MinLengthRequired     = "Longueur minimale requise"
	MaxLengthExceeded     = "Longueur maximale dépassée"
	NumericValueRequired  = "Valeur numérique requise"
	PositiveValueRequired = "Valeur positive requise"
)