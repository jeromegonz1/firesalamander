package api

import (
	"math/rand"
	"net/url"
	"time"

	"firesalamander/internal/constants"
	"firesalamander/internal/messages"
)

// SimulateAnalysis - Simuler une analyse SEO en arrière-plan
func SimulateAnalysis(analysisID string) {
	_, exists := Store.Get(analysisID)
	if !exists {
		return
	}

	// Phase 1: Découverte des pages
	simulatePhase(analysisID, messages.PhaseDiscoveryMsg, constants.DefaultProgressStart, 30, constants.DefaultSimulationSpeed, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.Status = constants.StatusProcessing
			a.PagesFound = int(float64(progress) * constants.PageDiscoveryFactor) // Facteur de découverte
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 2: Analyse SEO
	simulatePhase(analysisID, messages.PhaseSEOAnalysisMsg, constants.PhaseSEOStart, constants.PhaseSEOEnd, constants.PhaseSEOSpeed, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.PagesAnalyzed = int(float64(a.PagesFound) * float64(progress-constants.PhaseSEOStart) / constants.AnalysisProgressRatio) // Ratio d'analyse
			a.IssuesFound = int(float64(progress-constants.PhaseSEOStart) * constants.IssueAccumulationRate) // Taux de problèmes
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 3: Analyse IA
	simulatePhase(analysisID, messages.PhaseAIAnalysisMsg, constants.PhaseAIStart, constants.PhaseAIEnd, constants.PhaseAISpeed, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			a.EstimatedTime = calculateRemainingTime(progress)
		})
	})

	// Phase 4: Génération rapport
	simulatePhase(analysisID, messages.PhaseReportGenMsg, constants.PhaseReportStart, constants.PhaseReportEnd, constants.PhaseReportSpeed, func(progress int) {
		Store.Update(analysisID, func(a *AnalysisState) {
			a.Progress = progress
			if progress == constants.DefaultProgressEnd {
				a.Status = constants.StatusComplete
				a.EstimatedTime = messages.TimeEstimateComplete
				// Générer les résultats finaux
				// TODO: Implement test result generation
				a.Results = &ResultsResponse{
					Score: 75,
					PagesCount: a.PagesFound,
					Issues: []ResultIssue{},
					Warnings: []ResultWarning{},
					Analysis: AnalysisResult{
						Score: 75,
						PagesAnalyzed: a.PagesAnalyzed,
					},
				}
			} else {
				a.EstimatedTime = messages.TimeEstimateCalculating
			}
		})
	})
}

// simulatePhase - Simuler une phase de progression
func simulatePhase(analysisID, description string, startProgress, endProgress int, interval time.Duration, callback func(int)) {
	for progress := startProgress; progress <= endProgress; progress += rand.Intn(constants.ProgressRandomStep) + 1 {
		if progress > endProgress {
			progress = endProgress
		}
		
		callback(progress)
		
		// Variation aléatoire du timing pour rendre plus réaliste
		sleepTime := interval + time.Duration(rand.Intn(int(constants.TimingVariation/time.Millisecond)))*time.Millisecond
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
	case progress < constants.ProgressThreshold20:
		return messages.TimeEstimate2to3min
	case progress < constants.ProgressThreshold40:
		return messages.TimeEstimate1to2min
	case progress < constants.ProgressThreshold60:
		return messages.TimeEstimate45to60s
	case progress < constants.ProgressThreshold80:
		return messages.TimeEstimate30to45s
	case progress < constants.ProgressThreshold95:
		return messages.TimeEstimate15to30s
	case progress < constants.ProgressThreshold100:
		return messages.TimeEstimateFewSeconds
	default:
		return messages.TimeEstimateComplete
	}
}

// SimulateRealisticPages - Simuler un nombre réaliste de pages selon le domaine
func simulateRealisticPages(analysisURL string) int {
	parsedURL, err := url.Parse(analysisURL)
	if err != nil {
		return constants.DefaultMinPages
	}

	domain := parsedURL.Host
	
	// Simuler différents types de sites
	switch {
	case contains(domain, []string{"github.com", "gitlab.com"}):
		return rand.Intn(constants.TechnicalSiteMaxPages-constants.TechnicalSiteMinPages) + constants.TechnicalSiteMinPages
	case contains(domain, []string{"blog", "news", "journal"}):
		return rand.Intn(constants.BlogSiteMaxPages-constants.BlogSiteMinPages) + constants.BlogSiteMinPages  
	case contains(domain, []string{"shop", "store", "commerce"}):
		return rand.Intn(constants.EcommerceSiteMaxPages-constants.EcommerceSiteMinPages) + constants.EcommerceSiteMinPages
	case domain == constants.TestDefaultDomain:
		return rand.Intn(constants.TestSiteMaxPages-constants.TestSiteMinPages) + constants.TestSiteMinPages
	default:
		return rand.Intn(constants.CorporateSiteMaxPages-constants.CorporateSiteMinPages) + constants.CorporateSiteMinPages
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