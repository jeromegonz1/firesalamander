package api

import (
	"math/rand"
	"net/url"
	"time"
)

// SimulateAnalysis - Simuler une analyse SEO en arrière-plan
func SimulateAnalysis(analysisID string) {
	_, exists := Store.Get(analysisID)
	if !exists {
		return
	}

	// Phase 1: Découverte des pages (0-30%)
	simulatePhase(analysisID, "Découverte des pages...", 0, 30, 500*time.Millisecond, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.Status = "analyzing"
			a.PagesFound = int(float64(progress) * 1.5) // Simule la découverte progressive
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 2: Analyse SEO (30-70%)
	simulatePhase(analysisID, "Analyse SEO en cours...", 30, 70, 800*time.Millisecond, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.PagesAnalyzed = int(float64(a.PagesFound) * float64(progress-30) / 40.0) // Analyse progressive
			a.IssuesFound = int(float64(progress-30) * 0.1) // Accumule les problèmes
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 3: Analyse IA (70-95%)
	simulatePhase(analysisID, "Analyse IA...", 70, 95, 600*time.Millisecond, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 4: Génération rapport (95-100%)
	simulatePhase(analysisID, "Génération rapport...", 95, 100, 300*time.Millisecond, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			if progress == 100 {
				a.Status = "complete"
				a.EstimatedTime = "Terminé"
				// Générer les résultats finaux
				a.Results = GenerateTestResults(a.URL)
			} else {
				a.EstimatedTime = "Quelques secondes..."
			}
		})
	})
}

// simulatePhase - Simuler une phase de progression
func simulatePhase(analysisID, description string, startProgress, endProgress int, interval time.Duration, callback func(int)) {
	for progress := startProgress; progress <= endProgress; progress += rand.Intn(3) + 1 {
		if progress > endProgress {
			progress = endProgress
		}
		
		callback(progress)
		
		// Variation aléatoire du timing pour rendre plus réaliste
		sleepTime := interval + time.Duration(rand.Intn(200))*time.Millisecond
		time.Sleep(sleepTime)
		
		// Vérifier que l'analyse existe toujours
		if _, exists := Store.Get(analysisID); !exists {
			return
		}
	}
}

// calculateRemainingTime - Calculer le temps restant estimé
func calculateRemainingTime(progress int) string {
	switch {
	case progress < 20:
		return "2-3 minutes"
	case progress < 40:
		return "1-2 minutes"
	case progress < 60:
		return "45-60 secondes"
	case progress < 80:
		return "30-45 secondes"
	case progress < 95:
		return "15-30 secondes"
	case progress < 100:
		return "Quelques secondes..."
	default:
		return "Terminé"
	}
}

// SimulateRealisticPages - Simuler un nombre réaliste de pages selon le domaine
func simulateRealisticPages(analysisURL string) int {
	parsedURL, err := url.Parse(analysisURL)
	if err != nil {
		return 20 // Valeur par défaut
	}

	domain := parsedURL.Host
	
	// Simuler différents types de sites
	switch {
	case contains(domain, []string{"github.com", "gitlab.com"}):
		return rand.Intn(100) + 50 // Sites techniques : 50-150 pages
	case contains(domain, []string{"blog", "news", "journal"}):
		return rand.Intn(200) + 100 // Blogs/News : 100-300 pages  
	case contains(domain, []string{"shop", "store", "commerce"}):
		return rand.Intn(500) + 200 // E-commerce : 200-700 pages
	case domain == "example.com":
		return rand.Intn(30) + 20 // Site de test : 20-50 pages
	default:
		return rand.Intn(80) + 30 // Site corporate : 30-110 pages
	}
}

// contains - Vérifier si un domaine contient certains mots-clés
func contains(domain string, keywords []string) bool {
	for _, keyword := range keywords {
		if len(domain) > 0 && len(keyword) > 0 {
			for i := 0; i <= len(domain)-len(keyword); i++ {
				if domain[i:i+len(keyword)] == keyword {
					return true
				}
			}
		}
	}
	return false
}

// GenerateRealisticIssues - Générer des problèmes réalistes selon le type de site
func generateRealisticIssues(domain string, pagesCount int) []ResultIssue {
	issues := []ResultIssue{
		{
			Title:       "Balises title manquantes",
			Count:       max(1, pagesCount/10),
			Description: "Certaines pages n'ont pas de balise title ou celle-ci est vide.",
			Solution:    "Ajoutez une balise title unique et descriptive pour chaque page.",
		},
		{
			Title:       "Images sans attribut alt",
			Count:       max(2, pagesCount/8),
			Description: "Des images n'ont pas d'attribut alt pour l'accessibilité.",
			Solution:    "Ajoutez des attributs alt descriptifs à toutes vos images.",
		},
	}

	// Ajouter des problèmes spécifiques selon le type de site
	if contains(domain, []string{"blog", "news"}) {
		issues = append(issues, ResultIssue{
			Title:       "Dates de publication manquantes",
			Count:       max(1, pagesCount/15),
			Description: "Articles sans dates structurées pour les moteurs de recherche.",
			Solution:    "Ajoutez des métadonnées de date avec schema.org.",
		})
	}

	return issues
}

// max - Fonction utilitaire max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}