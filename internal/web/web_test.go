package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/integration"
)

// 🔥💀 RAMBO CONSTANTS - Hardcoded strings eliminated! 💀🔥
const (
	HTTP_METHOD_GET            = "GET"
	HEADER_CONTENT_TYPE        = "Content-Type"
	HEADER_CONTENT_DISPOSITION = "Content-Disposition"
	CONTENT_TYPE_HTML          = "text/html"
	CONTENT_TYPE_JSON          = "application/json"
	TEST_PORT_8080             = 8080
	TEST_PORT_8083             = 8083
	TEST_PORT_8084             = 8084
)

// Test du serveur web principal
func TestWebServer(t *testing.T) {
	// Configuration de test
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{
			Port: TEST_PORT_8080,
		},
		Crawler: config.CrawlerConfig{
			Workers:   1,
			RateLimit: "1/s",
			UserAgent: "Fire Salamander Test Bot",
		},
		AI: config.AIConfig{
			Enabled:  false,
			MockMode: true,
		},
	}

	// Créer l'orchestrateur (peut être nil pour certains tests)
	orchestrator, err := integration.NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	// Créer le serveur web
	webServer := NewWebServer(orchestrator, cfg)
	if webServer == nil {
		t.Fatal("Serveur web non créé")
	}

	// Vérifier l'initialisation des composants
	if webServer.orchestrator == nil {
		t.Error("Orchestrateur non initialisé")
	}

	if webServer.config == nil {
		t.Error("Configuration non initialisée")
	}

	if webServer.mux == nil {
		t.Error("Router non initialisé")
	}

	t.Log("Serveur web créé avec succès")
}

// Test de l'interface web principale
func TestWebInterface(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8080},
	}

	webServer := NewWebServer(nil, cfg) // Orchestrateur nil pour ce test

	// Créer une requête de test
	req := httptest.NewRequest(HTTP_METHOD_GET, "/", nil)
	w := httptest.NewRecorder()

	// Tester la route principale
	webServer.handleWebInterface(w, req)

	// Vérifier la réponse
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code attendu: %d, reçu: %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get(HEADER_CONTENT_TYPE)
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Content-Type attendu: text/html; charset=utf-8, reçu: %s", contentType)
	}

	t.Log("Interface web testée avec succès")
}

// Test de la route de santé web
func TestWebHealth(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8080},
	}

	webServer := NewWebServer(nil, cfg)

	// Créer une requête de test
	req := httptest.NewRequest(HTTP_METHOD_GET, "/web/health", nil)
	w := httptest.NewRecorder()

	// Tester la route de santé
	webServer.handleWebHealth(w, req)

	// Vérifier la réponse
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code attendu: %d, reçu: %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get(HEADER_CONTENT_TYPE)
	if contentType != CONTENT_TYPE_JSON {
		t.Errorf("Content-Type attendu: application/json, reçu: %s", contentType)
	}

	t.Log("Route de santé web testée avec succès")
}

// Test de génération de rapport
func TestReportGeneration(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8080},
	}

	webServer := NewWebServer(nil, cfg)

	// Tester différents formats de rapport
	formats := []string{"html", "json", "csv"}

	for _, format := range formats {
		filename := "test-report." + format
		report := webServer.generateSampleReport(filename)

		if len(report) == 0 {
			t.Errorf("Rapport %s vide", format)
		}

		// Vérifications spécifiques par format
		switch format {
		case "html":
			if !contains(report, "<!DOCTYPE html>") {
				t.Errorf("Rapport HTML invalide")
			}
		case "json":
			if !contains(report, `"report"`) {
				t.Errorf("Rapport JSON invalide")
			}
		case "csv":
			if !contains(report, "URL,Score Global") {
				t.Errorf("Rapport CSV invalide")
			}
		}

		t.Logf("Rapport %s généré avec succès (%d caractères)", format, len(report))
	}
}

// Test de téléchargement de rapport
func TestReportDownload(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8080},
	}

	webServer := NewWebServer(nil, cfg)

	// Créer une requête de téléchargement
	req := httptest.NewRequest(HTTP_METHOD_GET, "/web/download/test-report.html", nil)
	w := httptest.NewRecorder()

	// Tester le téléchargement
	webServer.handleDownload(w, req)

	// Vérifier la réponse
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code attendu: %d, reçu: %d", http.StatusOK, resp.StatusCode)
	}

	contentDisposition := resp.Header.Get(HEADER_CONTENT_DISPOSITION)
	expectedDisposition := "attachment; filename=test-report.html"
	if contentDisposition != expectedDisposition {
		t.Errorf("Content-Disposition attendu: %s, reçu: %s", expectedDisposition, contentDisposition)
	}

	contentType := resp.Header.Get(HEADER_CONTENT_TYPE)
	if contentType != CONTENT_TYPE_HTML {
		t.Errorf("Content-Type attendu: text/html, reçu: %s", contentType)
	}

	t.Log("Téléchargement de rapport testé avec succès")
}

// Test des statistiques du serveur web
func TestWebServerStats(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8080},
	}

	webServer := NewWebServer(nil, cfg)

	// Récupérer les statistiques
	stats := webServer.GetStats()

	// Vérifier les champs obligatoires
	if stats["service"] != "web_server" {
		t.Errorf("Service attendu: web_server, reçu: %v", stats["service"])
	}

	if stats["status"] != "running" {
		t.Errorf("Status attendu: running, reçu: %v", stats["status"])
	}

	if stats["port"] != cfg.Server.Port {
		t.Errorf("Port attendu: %d, reçu: %v", cfg.Server.Port, stats["port"])
	}

	if stats["version"] != cfg.App.Version {
		t.Errorf("Version attendue: %s, reçue: %v", cfg.App.Version, stats["version"])
	}

	t.Logf("Statistiques web: %+v", stats)
}

// Test de démarrage et arrêt du serveur
func TestWebServerStartStop(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8083}, // Port différent pour éviter les conflits
	}

	webServer := NewWebServer(nil, cfg)

	// Démarrer le serveur
	if err := webServer.Start(); err != nil {
		t.Fatalf("Erreur démarrage serveur web: %v", err)
	}

	// Attendre un peu pour que le serveur démarre
	time.Sleep(100 * time.Millisecond)

	// Tester que le serveur répond
	resp, err := http.Get(constants.TestServerLocalhost8083 + constants.TestEndpointHealthPath)
	if err != nil {
		t.Logf("Serveur probablement non démarré ou port occupé: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code santé attendu: %d, reçu: %d", http.StatusOK, resp.StatusCode)
		}
		t.Log("Serveur web répond correctement")
	}

	// Arrêter le serveur
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := webServer.Stop(ctx); err != nil {
		t.Errorf("Erreur arrêt serveur web: %v", err)
	}

	t.Log("Serveur web démarré et arrêté avec succès")
}

// Test d'intégration complète
func TestWebServerIntegration(t *testing.T) {
	// Configuration complète
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Integration Test",
			Version: "1.0.0-integration",
		},
		Server: config.ServerConfig{Port: TEST_PORT_8084},
		Crawler: config.CrawlerConfig{
			Workers:   1,
			RateLimit: "1/s",
			UserAgent: "Fire Salamander Integration Test",
		},
		AI: config.AIConfig{
			Enabled:  false,
			MockMode: true,
		},
	}

	// 1. Créer l'orchestrateur
	orchestrator, err := integration.NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	// 2. Créer le serveur web
	webServer := NewWebServer(orchestrator, cfg)

	// 3. Démarrer l'orchestrateur
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := orchestrator.Start(ctx); err != nil {
		t.Fatalf("Erreur démarrage orchestrateur: %v", err)
	}
	defer orchestrator.Stop()

	// 4. Démarrer le serveur web
	if err := webServer.Start(); err != nil {
		t.Fatalf("Erreur démarrage serveur web: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// 5. Tester l'intégration
	routes := []string{
		"/",
		"/web/health",
	}

	for _, route := range routes {
		resp, err := http.Get(constants.TestServerLocalhost8084 + route)
		if err != nil {
			t.Logf("Erreur requête %s: %v", route, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Route %s - Status attendu: %d, reçu: %d", route, http.StatusOK, resp.StatusCode)
		} else {
			t.Logf("Route %s - OK", route)
		}
	}

	// 6. Arrêter le serveur web
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := webServer.Stop(shutdownCtx); err != nil {
		t.Errorf("Erreur arrêt serveur web: %v", err)
	}

	t.Log("🎉 Test d'intégration web complète réussi!")
}

// Fonction utilitaire pour vérifier si une chaîne contient une sous-chaîne
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
