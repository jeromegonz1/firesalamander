package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/integration"
	"firesalamander/internal/web"
)

var (
	version = "1.0.0"
	banner  = `
🔥 Fire Salamander - Analyseur SEO Avancé
==========================================
Version: %s
Développé par SEPTEO
==========================================
`
)

func main() {
	// Parse command line flags
	var (
		configPath = flag.String("config", constants.DefaultConfigPath, constants.ConfigPathDescription)
		port       = flag.Int("port", constants.DefaultPortInt, constants.PortDescription)
		showVersion = flag.Bool("version", false, constants.ShowVersionDescription)
		webOnly    = flag.Bool("web-only", false, constants.WebOnlyDescription)
		apiOnly    = flag.Bool("api-only", false, constants.APIOnlyDescription)
		verbose    = flag.Bool("verbose", false, constants.VerboseDescription)
	)
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("Fire Salamander v%s\n", version)
		return
	}

	// Show banner
	fmt.Printf(banner, version)

	// Setup logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Load configuration
	cfg, err := loadConfig(*configPath, *port)
	if err != nil {
		log.Fatalf("❌ Erreur chargement configuration: %v", err)
	}

	log.Printf("📋 Configuration chargée depuis: %s", *configPath)
	log.Printf("🔧 Mode: %s", getRunMode(*webOnly, *apiOnly))

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var orchestrator *integration.Orchestrator
	var webServer *web.WebServer
	var apiServer *integration.APIServer

	// Initialize orchestrator (unless web-only mode)
	if !*webOnly {
		log.Printf("🚀 Initialisation de l'orchestrateur Fire Salamander...")
		
		orchestrator, err = integration.NewOrchestrator(cfg)
		if err != nil {
			log.Fatalf("❌ Erreur création orchestrateur: %v", err)
		}

		// Start orchestrator
		if err := orchestrator.Start(ctx); err != nil {
			log.Fatalf("❌ Erreur démarrage orchestrateur: %v", err)
		}
		log.Printf("✅ Orchestrateur démarré avec succès")
	}

	// Start services based on mode
	if *apiOnly {
		// API only mode
		log.Printf("🔌 Démarrage du serveur API uniquement...")
		apiServer = integration.NewAPIServer(orchestrator, cfg)
		
		if err := apiServer.Start(); err != nil {
			log.Fatalf("❌ Erreur démarrage API: %v", err)
		}
		
		log.Printf("✅ Serveur API démarré sur le port %d", cfg.Server.Port)
		log.Printf(constants.LogAPIAvailableFormat, cfg.Server.Port)
		
	} else {
		// Web mode (default) or web-only mode
		log.Printf("🌐 Démarrage du serveur web Fire Salamander...")
		webServer = web.NewWebServer(orchestrator, cfg)
		
		if err := webServer.Start(); err != nil {
			log.Fatalf("❌ Erreur démarrage serveur web: %v", err)
		}
		
		log.Printf("✅ Serveur web démarré sur le port %d", cfg.Server.Port)
		log.Printf(constants.LogInterfaceAvailableFormat, cfg.Server.Port)
		log.Printf(constants.LogAPIIntegratedFormat, cfg.Server.Port)
	}

	// Display startup summary
	displayStartupSummary(cfg, *webOnly, *apiOnly)

	// Wait for shutdown signal
	<-sigChan
	log.Printf("🛑 Signal d'arrêt reçu, fermeture gracieuse...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer shutdownCancel()

	// Shutdown services
	if webServer != nil {
		log.Printf("🌐 Arrêt du serveur web...")
		if err := webServer.Stop(shutdownCtx); err != nil {
			log.Printf("⚠️ Erreur arrêt serveur web: %v", err)
		}
	}

	if apiServer != nil {
		log.Printf("🔌 Arrêt du serveur API...")
		if err := apiServer.Stop(shutdownCtx); err != nil {
			log.Printf("⚠️ Erreur arrêt serveur API: %v", err)
		}
	}

	if orchestrator != nil {
		log.Printf("🚀 Arrêt de l'orchestrateur...")
		if err := orchestrator.Stop(); err != nil {
			log.Printf("⚠️ Erreur arrêt orchestrateur: %v", err)
		}
	}

	// Cancel main context
	cancel()

	log.Printf("✅ Fire Salamander arrêté proprement")
	log.Printf("👋 Merci d'avoir utilisé Fire Salamander!")
}

// loadConfig charge la configuration depuis le fichier ou utilise les valeurs par défaut
func loadConfig(configPath string, port int) (*config.Config, error) {
	// Try to load from file
	cfg, err := config.Load()
	if err != nil {
		// If file doesn't exist, create default config
		if os.IsNotExist(err) {
			log.Printf("📝 Fichier de configuration non trouvé, utilisation de la configuration par défaut")
			cfg = &config.Config{Server: config.ServerConfig{Port: 8080, Host: constants.ServerDefaultHost}}
		} else {
			return nil, fmt.Errorf("erreur lecture configuration: %w", err)
		}
	}

	// Override port if specified
	if port != constants.DefaultPortInt {
		cfg.Server.Port = port
	}

	// Set version
	cfg.App.Version = version

	return cfg, nil
}

// getRunMode retourne le mode d'exécution
func getRunMode(webOnly, apiOnly bool) string {
	if webOnly {
		return constants.RunModeWebOnly
	}
	if apiOnly {
		return constants.RunModeAPIOnly
	}
	return constants.RunModeComplete
}

// displayStartupSummary affiche un résumé du démarrage
func displayStartupSummary(cfg *config.Config, webOnly, apiOnly bool) {
	fmt.Println()
	fmt.Println("🎯 FIRE SALAMANDER DÉMARRÉ AVEC SUCCÈS")
	fmt.Println("=====================================")
	
	if !webOnly {
		fmt.Printf("🚀 Orchestrateur: Actif (%d workers)\n", cfg.Crawler.Workers)
		fmt.Printf("🔍 Analyses disponibles: Sémantique, SEO, Complète, Rapide\n")
	}
	
	if !apiOnly {
		fmt.Printf(constants.InterfaceWebFormat, cfg.Server.Port)
		fmt.Printf("   - Dashboard de monitoring\n")
		fmt.Printf("   - Outil d'analyse interactif\n")  
		fmt.Printf("   - Historique et rapports\n")
	}
	
	fmt.Printf(constants.APIRESTFormat, cfg.Server.Port)
	fmt.Printf("   - POST /analyze (analyse complète)\n")
	fmt.Printf("   - POST /analyze/quick (analyse rapide)\n")
	fmt.Printf("   - POST /analyze/seo (analyse SEO)\n")
	fmt.Printf("   - POST /analyze/semantic (analyse sémantique)\n")
	fmt.Printf("   - GET  /health (santé du service)\n")
	fmt.Printf("   - GET  /stats (statistiques)\n")
	
	fmt.Println()
	fmt.Println("📚 Documentation:")
	fmt.Printf(constants.DocInterfaceFormat, cfg.Server.Port)
	fmt.Printf(constants.DocAPIFormat, cfg.Server.Port)
	fmt.Printf(constants.DocHealthFormat, cfg.Server.Port)
	
	fmt.Println()
	fmt.Println("🔥 Prêt à analyser vos sites web!")
	fmt.Println("=====================================")
	
	if !webOnly && !apiOnly {
		fmt.Println()
		fmt.Println("💡 Exemple d'utilisation:")
		fmt.Printf(constants.CurlExampleFormat, cfg.Server.Port)
		fmt.Println(constants.CurlHeaders)
		fmt.Println(constants.CurlExampleData)
		fmt.Println()
	}
	
	fmt.Println("Appuyez sur Ctrl+C pour arrêter Fire Salamander")
	fmt.Println()
}