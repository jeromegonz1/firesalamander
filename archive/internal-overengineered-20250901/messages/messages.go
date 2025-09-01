package messages

import "firesalamander/internal/constants"

// Error messages externalisés
const (
    // HTTP Errors
    ErrMethodNotAllowed    = "Method not allowed"
    ErrInvalidJSON        = "Invalid JSON"
    ErrURLRequired        = "URL is required"
    ErrInvalidURL         = "Invalid URL format"
    ErrAnalysisNotFound   = "Analysis not found"
    ErrAnalysisIncomplete = "Analysis not complete"
    ErrInvalidAnalysisID  = "Invalid analysis ID"
    ErrConnectionFailed   = "Connection failed"
    ErrTimeout            = "Request timeout"
    
    // Server Messages
    ServerStarting        = "Fire Salamander starting"
    ServerStarted         = "Fire Salamander started on"
    ServerStopping        = "Fire Salamander shutting down"
    ServerStopped         = "Fire Salamander stopped"
    APIAvailable          = "API available on"
    WebInterfaceAvailable = "Web interface available on"
    
    // Analysis Messages  
    AnalysisStarted       = "Analysis started"
    AnalysisComplete      = "Analysis complete"
    AnalysisProgress      = "Analysis in progress"
    AnalysisError         = "Analysis error"
    
    // Phase Messages
    PhaseDiscoveryMsg     = "Discovering pages..."
    PhaseSEOAnalysisMsg   = "SEO analysis in progress..."
    PhaseAIAnalysisMsg    = "AI analysis..."
    PhaseReportGenMsg     = "Generating report..."
    
    // Time Estimates
    TimeEstimate30s         = "30 seconds"
    TimeEstimate1m          = "1 minute"  
    TimeEstimate2m          = "2 minutes"
    TimeEstimate5m          = "5 minutes"
    TimeEstimateComplete    = "Complete"
    TimeEstimateCalculating = "Calculating..."
    TimeEstimate2to3min     = "2-3 minutes"
    TimeEstimate1to2min     = "1-2 minutes"
    TimeEstimate45to60s     = "45-60 seconds"
    TimeEstimate30to45s     = "30-45 seconds"
    TimeEstimate15to30s     = "15-30 seconds"
    TimeEstimateFewSeconds  = "Few seconds..."
    
    // UI Messages
    UINewAnalysis         = "New Analysis"
    UIExportPDF          = "Export PDF"
    UICancel             = "Cancel"
    UIAnalyzeButton      = "ANALYZE MY SITE"
    UIAnalyzing          = "ANALYZING..."
    UIPlaceholderURL     = constants.UIPlaceholderURL
    UIURLRequiredHTTPS   = "URL must start with https://"
    UIConnectionError    = "Connection error. Please try again."
    UIPDFGenerated       = "PDF report has been generated and downloaded."
    
    // Log Messages
    LogConfigLoaded       = "Configuration loaded"
    LogConfigError        = "Configuration error"
    LogServerStarted      = "HTTP server started"  
    LogServerError        = "HTTP server error"
    LogAnalysisStarted    = "Analysis started for"
    LogAnalysisComplete   = "Analysis completed for"
    LogAnalysisError      = "Analysis error for"
    LogTemplateError      = "Template error"
    LogDatabaseError      = "Database error"
    
    // Help Messages
    HelpUsage            = "Usage:"
    HelpWebInterface     = "Web Interface:"
    HelpAPIREST          = "API REST:"
    HelpHealthCheck      = "Health Check:"
    HelpExampleCURL      = "Example cURL:"
    HelpPressCtrlC       = "Press Ctrl+C to stop Fire Salamander"
    
    // Mode Messages
    ModeWebOnly          = "Web Interface Only"
    ModeAPIOnly          = "REST API Only"  
    ModeComplete         = "Complete (Web + API + Orchestrator)"
    
    // SEO Messages
    SEOTitleMissing      = "Missing title tags"
    SEOTitleEmpty        = "Empty title tags"  
    SEOTitleTooShort     = "Title tags too short"
    SEOTitleTooLong      = "Title tags too long"
    SEODescMissing       = "Missing meta descriptions"
    SEODescTooShort      = "Meta descriptions too short"
    SEOImageAltMissing   = "Images without alt attributes"
    SEOH1Missing         = "Missing H1 tags"
    SEOH1Multiple        = "Multiple H1 tags"
    SEODuplicateContent  = "Duplicate content detected"
    
    // AI Suggestions
    AIKeywordOptimization = "Keyword optimization"
    AIContentStructure   = "Content structure improvement"
    AITechnicalSEO       = "Technical SEO recommendations"
    AIUserExperience     = "User experience optimization"
)