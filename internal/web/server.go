package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"firesalamander/internal/api"
	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/integration"
)

//go:embed static/*
var staticFiles embed.FS

// WebServer serveur web pour l'interface Fire Salamander
type WebServer struct {
	orchestrator *integration.Orchestrator
	apiServer    *integration.APIServer
	config       *config.Config
	server       *http.Server
	mux          *http.ServeMux
}

// NewWebServer crée un nouveau serveur web
func NewWebServer(orchestrator *integration.Orchestrator, cfg *config.Config) *WebServer {
	mux := http.NewServeMux()
	
	webServer := &WebServer{
		orchestrator: orchestrator,
		config:       cfg,
		mux:          mux,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  constants.ServerReadTimeout,
			WriteTimeout: constants.ServerWriteTimeout,
			IdleTimeout:  constants.ServerIdleTimeout,
		},
	}

	// Créer le serveur API intégré
	webServer.apiServer = integration.NewAPIServer(orchestrator, cfg)

	// 🔥🦎 SPRINT 5: Initialize Real Orchestrator for API
	log.Printf("🔥🦎 Initializing Real Fire Salamander API...")
	api.InitOrchestrator()
	log.Printf("✅ Real Fire Salamander API ready!")

	// Enregistrer les routes
	webServer.registerRoutes()

	return webServer
}

// registerRoutes enregistre toutes les routes
func (ws *WebServer) registerRoutes() {
	// Servir les fichiers statiques embarqués
	staticFS, err := fs.Sub(staticFiles, constants.StaticDirectory)
	if err != nil {
		log.Fatalf(constants.MsgErrorStaticFiles, err)
	}
	
	fileServer := http.FileServer(http.FS(staticFS))
	
	// Route pour les assets statiques
	ws.mux.Handle(constants.StaticRoute, http.StripPrefix(constants.StaticRoute, fileServer))
	
	// Route racine - servir l'interface web
	ws.mux.HandleFunc(constants.WebRouteRoot, ws.handleWebInterface)
	
	// Routes API - proxy vers le serveur API
	ws.mux.HandleFunc(constants.ServerEndpointAPI, ws.handleAPI)
	
	// Route de santé spécifique au serveur web
	ws.mux.HandleFunc(constants.WebRouteHealth, ws.handleWebHealth)
	
	// Route pour télécharger des rapports
	ws.mux.HandleFunc(constants.WebRouteDownload, ws.handleDownload)
	
	log.Printf(constants.MsgRoutesRegistered)
	log.Printf(constants.MsgRouteWebInterface)
	log.Printf(constants.MsgRouteStaticFiles)
	log.Printf(constants.MsgRouteAPIProxy)
	log.Printf(constants.MsgRouteWebHealth)
	log.Printf(constants.MsgRouteWebDownload)
	
	// 🔥🦎 SPRINT 5: Log Real API Routes
	log.Printf("🔥🦎 REAL Fire Salamander API Routes:")
	log.Printf("  POST /api/real/analyze - Start real SEO analysis")
	log.Printf("  GET  /api/real/status/{id} - Get real-time analysis status") 
	log.Printf("  GET  /api/real/results/{id} - Get real analysis results")
}

// handleWebInterface sert l'interface web principale
func (ws *WebServer) handleWebInterface(w http.ResponseWriter, r *http.Request) {
	// Si ce n'est pas la racine, vérifier si c'est un fichier statique
	if r.URL.Path != constants.WebRouteRoot {
		// Essayer de servir comme fichier statique
		staticFS, _ := fs.Sub(staticFiles, constants.StaticDirectory)
		if _, err := fs.Stat(staticFS, r.URL.Path[1:]); err == nil {
			http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
			return
		}
		
		// Si pas trouvé, rediriger vers la page d'accueil (SPA routing)
		if r.URL.Path != constants.WebRouteRoot {
			http.Redirect(w, r, constants.WebRouteRoot, http.StatusFound)
			return
		}
	}

	// Servir index.html pour toutes les routes SPA
	staticFS, _ := fs.Sub(staticFiles, constants.StaticDirectory)
	indexContent, err := fs.ReadFile(staticFS, constants.StaticIndexHTML)
	if err != nil {
		http.Error(w, constants.MsgWebInterfaceUnavailable, http.StatusInternalServerError)
		log.Printf(constants.MsgErrorIndexHTML, err)
		return
	}

	// Définir les headers appropriés
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeHTMLCharset)
	w.Header().Set(constants.HeaderCacheControl, constants.HeaderValueNoCache)
	w.Header().Set(constants.HeaderXFrameOptions, constants.HeaderValueDeny)
	w.Header().Set(constants.HeaderXContentType, constants.HeaderValueNoSniff)
	w.Header().Set(constants.HeaderContentLength, fmt.Sprintf("%d", len(indexContent)))
	
	// Écrire le contenu
	w.WriteHeader(http.StatusOK)
	w.Write(indexContent)
	
	log.Printf(constants.MsgWebInterfaceServed, r.Method, r.URL.Path)
}

// handleAPI proxy les requêtes API vers le serveur API intégré
func (ws *WebServer) handleAPI(w http.ResponseWriter, r *http.Request) {
	// Ajouter les headers CORS
	w.Header().Set(constants.HeaderCORSOrigin, constants.HeaderValueCORSAll)
	w.Header().Set(constants.HeaderCORSMethods, constants.HeaderValueCORSMethods)
	w.Header().Set(constants.HeaderCORSHeaders, constants.HeaderValueCORSHeaders)
	
	// Gérer les requêtes OPTIONS pour CORS
	if r.Method == constants.HTTPMethodOPTIONS {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Log de la requête API
	log.Printf(constants.MsgAPIRequest, r.Method, r.URL.Path)
	
	// Proxy vers les handlers de l'API server
	switch {
	case r.URL.Path == constants.APIRouteAnalyze:
		ws.apiServer.HandleAnalyze(w, r)
	case r.URL.Path == constants.APIRouteAnalyzeSemantic:
		ws.apiServer.HandleSemanticAnalysis(w, r)
	case r.URL.Path == constants.APIRouteAnalyzeSEO:
		ws.apiServer.HandleSEOAnalysis(w, r)
	case r.URL.Path == constants.APIRouteAnalyzeQuick:
		ws.apiServer.HandleQuickAnalysis(w, r)
	case r.URL.Path == constants.APIRouteHealth:
		ws.apiServer.HandleHealth(w, r)
	case r.URL.Path == constants.APIRouteStats:
		ws.apiServer.HandleStats(w, r)
	case r.URL.Path == constants.APIRouteAnalyses:
		ws.apiServer.HandleAnalyses(w, r)
	case strings.HasPrefix(r.URL.Path, constants.APIRouteAnalysisDetails):
		ws.apiServer.HandleAnalysisDetails(w, r)
	case r.URL.Path == constants.APIRouteInfo:
		ws.apiServer.HandleInfo(w, r)
	case r.URL.Path == constants.APIRouteVersion:
		ws.apiServer.HandleVersion(w, r)
	
	// 🔥🦎 SPRINT 5: REAL API ROUTES - Fire Salamander Integration  
	// Routes principales (utilisées par le frontend)
	case r.URL.Path == "/api/analyze":
		api.AnalyzeHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/status/"):
		api.StatusHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/results/"):
		api.ResultsHandler(w, r)
		
	// Routes explicites /api/real/* (pour compatibilité)
	case r.URL.Path == "/api/real/analyze":
		api.AnalyzeHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/real/status/"):
		api.StatusHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/real/results/"):
		api.ResultsHandler(w, r)
	
	default:
		// Route API non trouvée
		http.Error(w, constants.MsgAPIEndpointNotFound, http.StatusNotFound)
	}
}

// handleWebHealth retourne l'état de santé du serveur web
func (ws *WebServer) handleWebHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.HTTPMethodGET {
		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(constants.ServerHeaderContentType, constants.ServerContentTypeJSON)
	
	health := map[string]interface{}{
		constants.ServerJSONFieldStatus:     "healthy",
		"service":    "Fire Salamander Web Server",
		"version":    constants.AppVersion,
		constants.ServerJSONFieldTimestamp:  time.Now().Format(time.RFC3339),
		"uptime":     time.Since(time.Now()), // Placeholder - devrait être le vrai uptime
		"components": map[string]string{
			"web_server":    "healthy",
			"static_files":  "healthy",
			"api_proxy":     "healthy",
			"orchestrator":  "healthy",
		},
	}

	// Vérifier la santé de l'orchestrateur
	if ws.orchestrator != nil {
		stats := ws.orchestrator.GetStats()
		if stats != nil {
			health["orchestrator_stats"] = stats
		}
	} else {
		health["components"].(map[string]string)["orchestrator"] = "unavailable"
		health[constants.ServerJSONFieldStatus] = "degraded"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		constants.ServerJSONFieldStatus: "%s",
		"service": "%s",
		"version": "%s",
		constants.ServerJSONFieldTimestamp: "%s",
		"components": %v
	}`,
		health[constants.ServerJSONFieldStatus],
		health["service"],
		health["version"],
		health[constants.ServerJSONFieldTimestamp],
		`{"web_server":"healthy","static_files":"healthy","api_proxy":"healthy","orchestrator":"healthy"}`,
	)
}

// handleDownload gère le téléchargement de rapports
func (ws *WebServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != constants.HTTPMethodGET {
		http.Error(w, constants.MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Extraire le nom du fichier de l'URL
	filename := filepath.Base(r.URL.Path)
	if filename == "." || filename == "/" {
		http.Error(w, constants.MsgFilenameRequired, http.StatusBadRequest)
		return
	}

	log.Printf(constants.MsgDownloadRequested, filename)

	// Pour l'instant, retourner un rapport d'exemple
	// Dans une vraie implémentation, récupérer le rapport depuis le stockage
	sampleReport := ws.generateSampleReport(filename)
	
	// Définir les headers de téléchargement
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	// Déterminer le type de contenu basé sur l'extension
	ext := filepath.Ext(filename)
	switch ext {
	case constants.ServerExtensionHTML:
		w.Header().Set(constants.ServerHeaderContentType, constants.ServerContentTypeHTML)
	case ".pdf":
		w.Header().Set(constants.ServerHeaderContentType, "application/pdf")
	case constants.ServerExtensionJSON:
		w.Header().Set(constants.ServerHeaderContentType, constants.ServerContentTypeJSON)
	case ".csv":
		w.Header().Set(constants.ServerHeaderContentType, "text/csv")
	default:
		w.Header().Set(constants.ServerHeaderContentType, "application/octet-stream")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sampleReport))
}

// generateSampleReport génère un rapport d'exemple
func (ws *WebServer) generateSampleReport(filename string) string {
	ext := filepath.Ext(filename)
	
	switch ext {
	case constants.ServerExtensionHTML:
		return ws.generateHTMLReport()
	case constants.ServerExtensionJSON:
		return ws.generateJSONReport()
	case ".csv":
		return ws.generateCSVReport()
	default:
		return "Rapport Fire Salamander - " + filename
	}
}

// generateHTMLReport génère un rapport HTML
func (ws *WebServer) generateHTMLReport() string {
	return `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Rapport Fire Salamander</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; border-bottom: 2px solid #ff6b35; padding-bottom: 20px; margin-bottom: 30px; }
        .logo { font-size: 2rem; margin-bottom: 10px; }
        .score { font-size: 3rem; font-weight: bold; color: #ff6b35; text-align: center; margin: 20px 0; }
        .section { margin: 30px 0; padding: 20px; background: #f8f9fa; border-radius: 8px; }
        .metric { display: flex; justify-content: space-between; margin: 10px 0; }
        .recommendation { background: white; padding: 15px; margin: 10px 0; border-left: 4px solid #ff6b35; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🔥 Fire Salamander</div>
        <h1>Rapport d'Analyse SEO</h1>
        <p>Généré le ` + time.Now().Format("02/01/2006 à 15:04") + `</p>
    </div>
    
    <div class="section">
        <h2>Résumé Exécutif</h2>
        <div class="score">87/100</div>
        <p>Votre site présente de bonnes performances SEO avec quelques axes d'amélioration identifiés.</p>
    </div>
    
    <div class="section">
        <h2>Métriques Détaillées</h2>
        <div class="metric"><span>SEO Technique:</span><span>85/100</span></div>
        <div class="metric"><span>Performance:</span><span>72/100</span></div>
        <div class="metric"><span>Contenu:</span><span>90/100</span></div>
        <div class="metric"><span>Mobile:</span><span>88/100</span></div>
    </div>
    
    <div class="section">
        <h2>Recommandations Prioritaires</h2>
        <div class="recommendation">
            <h3>Optimiser les images</h3>
            <p>Compresser les images pour améliorer les temps de chargement.</p>
        </div>
        <div class="recommendation">
            <h3>Améliorer les méta-descriptions</h3>
            <p>Ajouter des méta-descriptions optimisées sur les pages manquantes.</p>
        </div>
    </div>
    
    <footer style="text-align: center; margin-top: 50px; color: #666;">
        <p>Rapport généré par Fire Salamander v` + ws.config.App.Version + `</p>
    </footer>
</body>
</html>`
}

// generateJSONReport génère un rapport JSON
func (ws *WebServer) generateJSONReport() string {
	return `{
  "report": {
    "title": "Rapport Fire Salamander",
    "generated_at": "` + time.Now().Format(time.RFC3339) + `",
    "version": "` + ws.config.App.Version + `",
    "overall_score": ` + fmt.Sprintf("%d", constants.SampleOverallScore) + `,
    constants.ServerJSONFieldURL: "` + constants.TestExampleURL + `",
    "categories": {
      "technical": ` + fmt.Sprintf("%d", constants.SampleTechnicalScore) + `,
      "performance": ` + fmt.Sprintf("%d", constants.SamplePerformanceScore) + `,
      "content": ` + fmt.Sprintf("%d", constants.SampleContentScore) + `,
      "mobile": ` + fmt.Sprintf("%d", constants.SampleMobileScore) + `
    },
    "recommendations": [
      {
        constants.ServerJSONFieldID: "optimize_images",
        "title": "Optimiser les images",
        "description": "Compresser les images pour améliorer les temps de chargement",
        "priority": "high",
        "impact": "high",
        "effort": "medium"
      },
      {
        constants.ServerJSONFieldID: "meta_descriptions",
        "title": "Améliorer les méta-descriptions",
        "description": "Ajouter des méta-descriptions optimisées sur les pages manquantes",
        "priority": "medium",
        "impact": "medium",
        "effort": "low"
      }
    ],
    "insights": [
      {
        constants.ServerJSONFieldType: "content_seo_alignment",
        "severity": "info",
        "title": "Alignement contenu-SEO détecté",
        "description": "Le titre de la page est cohérent avec les mots-clés identifiés"
      }
    ],
    "metadata": {
      "analysis_duration": "15.2s",
      "pages_analyzed": 1,
      "issues_found": 3,
      "score_trend": "stable"
    }
  }
}`
}

// generateCSVReport génère un rapport CSV
func (ws *WebServer) generateCSVReport() string {
	return `URL,Score Global,SEO Technique,Performance,Contenu,Mobile,Statut,Date
` + constants.TestExampleURL + `,87,85,72,90,88,Succès,` + time.Now().Format("02/01/2006") + `
` + constants.TestDemoURL + `,72,78,65,85,70,Succès,` + time.Now().AddDate(0, 0, -1).Format("02/01/2006") + `
` + constants.TestDemoFrURL + `,91,92,88,95,89,Succès,` + time.Now().AddDate(0, 0, -2).Format("02/01/2006") + `

Résumé:
Score moyen,83.3
Analyses réussies,3
Analyses échouées,0
Dernière mise à jour,` + time.Now().Format("02/01/2006 15:04") + `

Recommandations principales:
1. Optimiser les images pour améliorer les performances
2. Améliorer les méta-descriptions
3. Optimiser la structure des URLs`
}

// Start démarre le serveur web
func (ws *WebServer) Start() error {
	log.Printf(constants.MsgWebServerStarting)
	log.Printf(constants.WebInterfaceAvailableFormat, ws.config.Server.Port)
	log.Printf(constants.APIRestAvailableFormat, ws.config.Server.Port)
	
	// Démarrer le serveur dans une goroutine
	go func() {
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(constants.MsgWebServerError, err)
		}
	}()
	
	log.Printf(constants.MsgWebServerStarted)
	
	return nil
}

// Stop arrête le serveur web
func (ws *WebServer) Stop(ctx context.Context) error {
	log.Printf(constants.MsgWebServerStopping)
	return ws.server.Shutdown(ctx)
}

// GetStats retourne les statistiques du serveur web
func (ws *WebServer) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"service":    "web_server",
		constants.ServerJSONFieldStatus:     "running",
		constants.ServerConfigPort:       ws.config.Server.Port,
		"version":    constants.AppVersion,
		"started_at": time.Now().Format(time.RFC3339),
	}
	
	// Ajouter les stats de l'orchestrateur si disponible
	if ws.orchestrator != nil {
		if orchestratorStats := ws.orchestrator.GetStats(); orchestratorStats != nil {
			stats["orchestrator"] = orchestratorStats
		}
	}
	
	return stats
}