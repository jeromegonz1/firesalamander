package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"firesalamander/internal/config"
)

// Test de l'orchestrateur principal
func TestOrchestratorBasic(t *testing.T) {
	// Configuration de test
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{
			Port: 8080,
		},
		Crawler: config.CrawlerConfig{
			Workers:   2,
			RateLimit: "1/s",
			UserAgent: "Fire Salamander Test Bot",
		},
		AI: config.AIConfig{
			Enabled:  false, // Désactivé pour les tests
			MockMode: true,
		},
	}

	// Créer l'orchestrateur
	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	if orchestrator == nil {
		t.Fatal("Orchestrateur non créé")
	}

	// Vérifier l'initialisation des composants
	if orchestrator.crawler == nil {
		t.Error("Crawler non initialisé")
	}

	if orchestrator.semanticAnalyzer == nil {
		t.Error("Analyseur sémantique non initialisé")
	}

	if orchestrator.seoAnalyzer == nil {
		t.Error("Analyseur SEO non initialisé")
	}

	t.Log("Orchestrateur créé avec succès")
}

// Test d'analyse complète
func TestFullAnalysis(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{
			Port: 8080,
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

	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	// Démarrer l'orchestrateur
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := orchestrator.Start(ctx); err != nil {
		t.Fatalf("Erreur démarrage orchestrateur: %v", err)
	}
	defer orchestrator.Stop()

	// Lancer une analyse rapide (plus fiable pour les tests)
	options := AnalysisOptions{
		IncludeCrawling:    false, // Désactivé pour éviter les erreurs réseau
		AnalyzePerformance: true,
		UseAIEnrichment:    false,
		Timeout:            15 * time.Second,
	}

	result, err := orchestrator.AnalyzeURL(ctx, "https://example.com", AnalysisTypeQuick, options)
	if err != nil {
		t.Fatalf("Erreur analyse: %v", err)
	}

	// Vérifications de base
	if result == nil {
		t.Fatal("Résultat d'analyse nil")
	}

	if result.URL != "https://example.com" {
		t.Errorf("URL incorrecte: %s", result.URL)
	}

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Score global invalide: %.1f", result.OverallScore)
	}

	if len(result.CategoryScores) == 0 {
		t.Error("Aucun score par catégorie")
	}

	if result.ProcessingTime == 0 {
		t.Error("Temps de traitement invalide")
	}

	t.Logf("Analyse terminée - URL: %s, Score: %.1f, Durée: %v, Statut: %s",
		result.URL, result.OverallScore, result.ProcessingTime, result.Status)

	// Vérifier la présence d'au moins une analyse (sémantique ou SEO)
	if result.SemanticAnalysis == nil && result.SEOAnalysis == nil {
		t.Error("Aucune analyse effectuée")
	}

	// Vérifier les insights cross-modules
	t.Logf("Insights générés: %d", len(result.CrossModuleInsights))
	t.Logf("Actions prioritaires: %d", len(result.PriorityActions))
}

// Test du serveur API
func TestAPIServerBasic(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Test",
			Version: "1.0.0-test",
		},
		Server: config.ServerConfig{
			Port: 8081, // Port différent pour éviter les conflits
		},
		Crawler: config.CrawlerConfig{
			Workers:   1,
			RateLimit: "1/s",
		},
		AI: config.AIConfig{
			Enabled:  false,
			MockMode: true,
		},
	}

	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	// Créer le serveur API
	apiServer := NewAPIServer(orchestrator, cfg)
	if apiServer == nil {
		t.Fatal("Serveur API non créé")
	}

	// Démarrer le serveur (il démarre dans une goroutine)
	if err := apiServer.Start(); err != nil {
		t.Fatalf("Erreur démarrage serveur API: %v", err)
	}

	// Attendre un peu pour que le serveur démarre
	time.Sleep(100 * time.Millisecond)

	// Arrêter le serveur
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Stop(ctx); err != nil {
		t.Errorf("Erreur arrêt serveur API: %v", err)
	}

	t.Log("Serveur API testé avec succès")
}

// Test du générateur de rapports
func TestReportGenerator(t *testing.T) {
	// Créer un résultat d'analyse simulé
	result := &UnifiedAnalysisResult{
		TaskID:         "test_task_123",
		URL:            "https://example.com",
		Domain:         "example.com",
		AnalyzedAt:     time.Now(),
		ProcessingTime: 2 * time.Second,
		OverallScore:   75.5,
		Status:         AnalysisStatusSuccess,
		CategoryScores: map[string]float64{
			"content":     80.0,
			"technical":   70.0,
			"performance": 75.0,
			"mobile":      85.0,
		},
		UnifiedMetrics: UnifiedMetrics{
			ContentQualityScore:     80.0,
			TechnicalHealthScore:    70.0,
			SEOReadinessScore:       75.5,
			PerformanceScore:        75.0,
			MobileFriendlinessScore: 85.0,
			UserExperienceScore:     77.5,
		},
		CrossModuleInsights: []CrossModuleInsight{
			{
				Type:        "test_insight",
				Severity:    "info",
				Title:       "Test Insight",
				Description: "Ceci est un insight de test",
				Evidence:    []string{"Evidence 1", "Evidence 2"},
				Modules:     []string{"seo", "semantic"},
				Impact:      "positive",
			},
		},
		PriorityActions: []PriorityAction{
			{
				ID:          "test_action_1",
				Title:       "Action prioritaire de test",
				Description: "Description de l'action de test",
				Priority:    "high",
				Impact:      "high",
				Effort:      "medium",
				Module:      "seo",
				EstimatedTime: "2-4 heures",
			},
		},
	}

	// Créer le générateur de rapports
	generator := NewReportGenerator()
	if generator == nil {
		t.Fatal("Générateur de rapports non créé")
	}

	// Options de rapport
	options := ReportOptions{
		Format:         ReportFormatHTML,
		Type:           ReportTypeExecutive,
		IncludeSummary: true,
		IncludeDetails: true,
		BrandingOptions: BrandingOptions{
			CompanyName: "Fire Salamander Test",
		},
	}

	// Générer le rapport HTML
	report, err := generator.GenerateReport(result, options)
	if err != nil {
		t.Fatalf("Erreur génération rapport HTML: %v", err)
	}

	if report == nil {
		t.Fatal("Rapport non généré")
	}

	if report.Content == "" {
		t.Error("Contenu de rapport vide")
	}

	if report.Size == 0 {
		t.Error("Taille de rapport invalide")
	}

	t.Logf("Rapport HTML généré - Taille: %d octets, Type: %s, Format: %s",
		report.Size, report.Type, report.Format)

	// Générer le rapport JSON
	options.Format = ReportFormatJSON
	reportJSON, err := generator.GenerateReport(result, options)
	if err != nil {
		t.Fatalf("Erreur génération rapport JSON: %v", err)
	}

	if len(reportJSON.Content) == 0 {
		t.Error("Rapport JSON vide")
	}

	t.Logf("Rapport JSON généré - Taille: %d octets", reportJSON.Size)

	// Générer le rapport CSV
	options.Format = ReportFormatCSV
	reportCSV, err := generator.GenerateReport(result, options)
	if err != nil {
		t.Fatalf("Erreur génération rapport CSV: %v", err)
	}

	if len(reportCSV.Content) == 0 {
		t.Error("Rapport CSV vide")
	}

	t.Logf("Rapport CSV généré - Taille: %d octets", reportCSV.Size)
}

// Test du gestionnaire de stockage
func TestStorageManager(t *testing.T) {
	// Créer un fichier de base temporaire
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_firesalamander.db")

	// Créer le gestionnaire de stockage
	storage, err := NewStorageManager(dbPath)
	if err != nil {
		t.Fatalf("Erreur création gestionnaire stockage: %v", err)
	}
	defer storage.Close()

	// Vérifier que le fichier DB existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Fichier de base de données non créé")
	}

	// Créer un résultat d'analyse de test
	result := &UnifiedAnalysisResult{
		TaskID:         "test_storage_123",
		URL:            "https://test-storage.com",
		Domain:         "test-storage.com",
		AnalyzedAt:     time.Now(),
		ProcessingTime: 1 * time.Second,
		OverallScore:   85.0,
		Status:         AnalysisStatusSuccess,
	}

	// Sauvegarder l'analyse
	if err := storage.SaveAnalysis(result); err != nil {
		t.Fatalf("Erreur sauvegarde analyse: %v", err)
	}

	// Récupérer l'analyse
	retrieved, err := storage.GetAnalysis("test_storage_123")
	if err != nil {
		t.Fatalf("Erreur récupération analyse: %v", err)
	}

	if retrieved.TaskID != result.TaskID {
		t.Errorf("TaskID incorrect: %s vs %s", retrieved.TaskID, result.TaskID)
	}

	if retrieved.OverallScore != result.OverallScore {
		t.Errorf("Score incorrect: %.1f vs %.1f", retrieved.OverallScore, result.OverallScore)
	}

	// Tester l'historique
	history, err := storage.GetAnalysisHistory("https://test-storage.com", 10)
	if err != nil {
		t.Fatalf("Erreur récupération historique: %v", err)
	}

	if len(history.Analyses) != 1 {
		t.Errorf("Nombre d'analyses incorrect: %d", len(history.Analyses))
	}

	// Tester les statistiques
	stats, err := storage.GetStorageStats()
	if err != nil {
		t.Fatalf("Erreur récupération statistiques: %v", err)
	}

	if stats["total_analyses"].(int) != 1 {
		t.Errorf("Nombre total d'analyses incorrect: %v", stats["total_analyses"])
	}

	t.Logf("Stockage testé avec succès - Stats: %+v", stats)
}

// Test d'intégration complète
func TestCompleteIntegration(t *testing.T) {
	// Configuration complète
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Fire Salamander Integration Test",
			Version: "1.0.0-integration",
		},
		Server: config.ServerConfig{
			Port: 8082,
		},
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
	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("Erreur création orchestrateur: %v", err)
	}

	// 2. Créer le gestionnaire de stockage
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "integration_test.db")
	storage, err := NewStorageManager(dbPath)
	if err != nil {
		t.Fatalf("Erreur création stockage: %v", err)
	}
	defer storage.Close()

	// 3. Créer le serveur API
	apiServer := NewAPIServer(orchestrator, cfg)

	// 4. Créer le générateur de rapports
	reportGenerator := NewReportGenerator()

	// 5. Démarrer l'orchestrateur
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := orchestrator.Start(ctx); err != nil {
		t.Fatalf("Erreur démarrage orchestrateur: %v", err)
	}
	defer orchestrator.Stop()

	// 6. Effectuer une analyse
	options := AnalysisOptions{
		IncludeCrawling:    false,
		AnalyzePerformance: true,
		UseAIEnrichment:    false,
		Timeout:            10 * time.Second,
	}

	result, err := orchestrator.AnalyzeURL(ctx, "https://integration-test.com", AnalysisTypeQuick, options)
	if err != nil {
		t.Fatalf("Erreur analyse intégration: %v", err)
	}

	// 7. Sauvegarder le résultat
	if err := storage.SaveAnalysis(result); err != nil {
		t.Fatalf("Erreur sauvegarde intégration: %v", err)
	}

	// 8. Générer un rapport
	reportOptions := ReportOptions{
		Format:         ReportFormatHTML,
		Type:           ReportTypeExecutive,
		IncludeSummary: true,
		IncludeDetails: true,
	}

	report, err := reportGenerator.GenerateReport(result, reportOptions)
	if err != nil {
		t.Fatalf("Erreur génération rapport intégration: %v", err)
	}

	// 9. Vérifications finales (accepter score 0 pour les erreurs réseau en test)
	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Score global invalide: %.1f", result.OverallScore)
	}

	if len(result.CategoryScores) == 0 {
		t.Error("Scores par catégorie non calculés")
	}

	if report.Size == 0 {
		t.Error("Rapport non généré")
	}

	// 10. Tester le serveur API (démarrage/arrêt)
	if err := apiServer.Start(); err != nil {
		t.Fatalf("Erreur démarrage API intégration: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := apiServer.Stop(ctx); err != nil {
		t.Errorf("Erreur arrêt API intégration: %v", err)
	}

	t.Logf("🎉 Test d'intégration complète réussi!")
	t.Logf("   - Analyse: Score %.1f, Durée %v", result.OverallScore, result.ProcessingTime)
	t.Logf("   - Rapport: %d octets, Format %s", report.Size, report.Format)
	t.Logf("   - Stockage: Analyse sauvegardée avec succès")
	t.Logf("   - API: Serveur démarré et arrêté avec succès")
}