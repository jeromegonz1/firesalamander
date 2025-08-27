package crawler

import (
	"context"
	"sync"
	"testing"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
)

// ========================================
// TDD RED PHASE - TESTS POUR TERMINATION CONTROLLER
// Tests multi-thread pour éviter les race conditions
// ========================================

// TestTerminationController_NoRaceConditions teste la gestion atomique multi-thread
func TestTerminationController_NoRaceConditions(t *testing.T) {
	// ARRANGE : Configuration pour tests de concurrence
	cfg := &config.CrawlerConfig{
		MaxPages:             100,
		TimeoutSeconds:       30,
		InitialWorkers:       10, // Beaucoup de workers pour stresser le système
		MaxWorkers:           15,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Créer le contrôleur de terminaison
	controller := NewSmartTerminationController(cfg)
	if controller == nil {
		t.Fatal("NewSmartTerminationController devrait retourner un contrôleur, got nil")
	}

	// Test des opérations concurrentes sur 1000 itérations
	numGoroutines := 50
	numIterations := 100
	var wg sync.WaitGroup

	// Compteur pour vérifier la cohérence
	expectedJobs := int32(numGoroutines * numIterations)

	// Lancement des goroutines qui simulent des jobs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Simuler démarrage d'un job
				controller.NotifyJobStarted()
				
				// Petite pause pour augmenter les chances de race condition
				time.Sleep(time.Microsecond)
				
				// Simuler fin d'un job
				controller.NotifyJobCompleted()
			}
		}(i)
	}

	// Attendre que toutes les goroutines finissent
	wg.Wait()

	// ASSERT : Vérifications anti-race condition
	finalJobCount := controller.GetActiveJobsCount()
	if finalJobCount != 0 {
		t.Errorf("Après toutes les opérations, le compteur de jobs actifs devrait être 0, got %d", finalJobCount)
	}

	// Vérifier que le nombre total de jobs traités est correct
	totalJobsProcessed := controller.GetTotalJobsProcessed()
	if totalJobsProcessed != expectedJobs {
		t.Errorf("Nombre total de jobs traités incorrect: expected %d, got %d", expectedJobs, totalJobsProcessed)
	}

	// Vérifier qu'aucune race condition n'a été détectée
	raceConditionsDetected := controller.GetRaceConditionsDetected()
	if raceConditionsDetected > 0 {
		t.Errorf("Des race conditions ont été détectées: %d", raceConditionsDetected)
	}
}

// TestSmartTerminationController_TwoPhaseTermination teste la terminaison en deux phases
func TestSmartTerminationController_TwoPhaseTermination(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             5,
		TimeoutSeconds:       30,
		InitialWorkers:       3,
		MaxWorkers:           5,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	controller := NewSmartTerminationController(cfg)

	// ACT : Test de la terminaison en deux phases
	// Phase 1 : Conditions de pré-terminaison
	controller.NotifyJobStarted()
	controller.NotifyJobStarted()
	controller.NotifyJobStarted()

	phase1Should := controller.ShouldEnterPreTermination()
	if phase1Should {
		t.Error("Ne devrait pas entrer en pré-terminaison avec des jobs actifs")
	}

	// Phase 2 : Terminaison effective
	controller.NotifyJobCompleted()
	controller.NotifyJobCompleted()
	
	phase2Should := controller.ShouldEnterPreTermination()
	if phase2Should {
		t.Error("Ne devrait pas entrer en pré-terminaison tant qu'il y a des jobs")
	}

	// Finir le dernier job
	controller.NotifyJobCompleted()

	finalShould := controller.ShouldTerminate()
	if !finalShould {
		t.Error("Devrait pouvoir terminer quand aucun job n'est actif")
	}

	// ASSERT : Vérifier l'état final
	terminationReason := controller.GetTerminationReason()
	if terminationReason == "" {
		t.Error("Une raison de terminaison devrait être fournie")
	}
}

// TestTerminationController_MultipleConditions teste les conditions multiples
func TestTerminationController_MultipleConditions(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             10,
		TimeoutSeconds:       30,
		InitialWorkers:       2,
		MaxWorkers:           5,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	controller := NewSmartTerminationController(cfg)

	// ACT : Obtenir les conditions de terminaison
	conditions := controller.GetTerminationConditions()

	// ASSERT : Vérifier les conditions obligatoires
	if len(conditions) < 3 {
		t.Errorf("Devrait avoir au moins 3 conditions de terminaison, got %d", len(conditions))
	}

	// Vérifier les types de conditions requis
	conditionTypes := make(map[string]bool)
	for _, condition := range conditions {
		conditionTypes[condition.Type()] = true
	}

	requiredConditions := []string{
		"active_jobs_zero",
		"max_pages_reached", 
		"timeout_reached",
	}

	for _, required := range requiredConditions {
		if !conditionTypes[required] {
			t.Errorf("Condition manquante: %s", required)
		}
	}
}

// TestTerminationController_GracefulShutdown teste l'arrêt propre
func TestTerminationController_GracefulShutdown(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             20,
		TimeoutSeconds:       30,
		InitialWorkers:       5,
		MaxWorkers:           10,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	controller := NewSmartTerminationController(cfg)

	// ACT : Simuler des jobs qui doivent se terminer proprement
	ctx, cancel := context.WithCancel(context.Background())
	
	// Démarrer des jobs
	for i := 0; i < 3; i++ {
		controller.NotifyJobStarted()
	}

	// Demander un arrêt gracieux
	shutdownStarted := time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond) // Simuler du travail
		for i := 0; i < 3; i++ {
			controller.NotifyJobCompleted()
		}
	}()

	// Attendre la terminaison ou le timeout
	shutdownComplete := controller.WaitForGracefulShutdown(ctx, 5*time.Second)
	shutdownDuration := time.Since(shutdownStarted)

	// Annuler le contexte
	cancel()

	// ASSERT
	if !shutdownComplete {
		t.Error("Le shutdown gracieux devrait se terminer avec succès")
	}

	if shutdownDuration > 2*time.Second {
		t.Errorf("Le shutdown gracieux ne devrait pas prendre plus de 2s, got %v", shutdownDuration)
	}

	finalJobCount := controller.GetActiveJobsCount()
	if finalJobCount != 0 {
		t.Errorf("Après le shutdown gracieux, aucun job ne devrait être actif, got %d", finalJobCount)
	}
}

// TestTerminationController_MemoryLeaks teste qu'il n'y a pas de fuites mémoire
func TestTerminationController_MemoryLeaks(t *testing.T) {
	// ARRANGE
	cfg := &config.CrawlerConfig{
		MaxPages:             1000,
		TimeoutSeconds:       60,
		InitialWorkers:       10,
		MaxWorkers:           20,
		MinWorkers:           1,
		UserAgent:            constants.ParallelCrawlerUserAgent,
		FastThresholdMs:      constants.DefaultFastThresholdMs,
		SlowThresholdMs:      constants.DefaultSlowThresholdMs,
		ErrorThresholdPercent: constants.DefaultErrorThresholdPercent,
		AdaptIntervalSeconds: constants.DefaultAdaptIntervalSeconds,
	}

	// ACT : Créer et détruire plusieurs contrôleurs
	for i := 0; i < 100; i++ {
		controller := NewSmartTerminationController(cfg)
		
		// Simuler de l'activité
		for j := 0; j < 10; j++ {
			controller.NotifyJobStarted()
			controller.NotifyJobCompleted()
		}
		
		// Nettoyer
		controller.Cleanup()
		controller = nil
	}

	// ASSERT : Forcer le garbage collector pour détecter les fuites
	// Note: Ce test est plus pour la documentation, il faudrait des outils
	// spécifiques pour vraiment détecter les fuites mémoire
	t.Log("Test de fuites mémoire complété - utilisez des outils de profiling pour validation")
}

// ========================================
// INTERFACES POUR LE TERMINATION CONTROLLER (À IMPLÉMENTER)
// ========================================

// ISmartTerminationController interface étendue pour la terminaison intelligente
type ISmartTerminationController interface {
	ITerminationController
	
	// Nouvelles méthodes pour les tests
	GetActiveJobsCount() int32
	GetTotalJobsProcessed() int32
	GetRaceConditionsDetected() int32
	GetTerminationReason() string
	ShouldEnterPreTermination() bool
	WaitForGracefulShutdown(ctx context.Context, timeout time.Duration) bool
	Cleanup()
}

// ========================================
// TYPES À IMPLÉMENTER (STUBS POUR TDD)
// ========================================

// SmartTerminationController implémentation intelligente (À IMPLÉMENTER)
type SmartTerminationController struct {
	// À implémenter selon l'architecture
}

// ========================================
// FONCTIONS À IMPLÉMENTER (STUBS POUR TESTS ROUGES)
// ========================================

// NewSmartTerminationController crée un contrôleur de terminaison intelligent
func NewSmartTerminationController(cfg *config.CrawlerConfig) ISmartTerminationController {
	// TODO: Implémenter selon l'architecture reçue
	// Cette fonction doit retourner nil pour que les tests échouent (RED phase)
	return nil
}