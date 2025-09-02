package broken

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/constants"
)

// BrokenLinksDetector implémente l'agent de détection de liens brisés
type BrokenLinksDetector struct {
	name       string
	client     *http.Client
	maxWorkers int
}

// NewBrokenLinksDetector crée une nouvelle instance de BrokenLinksDetector
func NewBrokenLinksDetector() *BrokenLinksDetector {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: 10,
		},
	}

	return &BrokenLinksDetector{
		name:       constants.AgentNameBrokenLinks,
		client:     client,
		maxWorkers: 10,
	}
}

// Name retourne le nom de l'agent
func (b *BrokenLinksDetector) Name() string {
	return b.name
}

// Process traite les données d'entrée et détecte les liens brisés
func (b *BrokenLinksDetector) Process(ctx context.Context, data interface{}) (*agents.AgentResult, error) {
	startTime := time.Now()
	
	urls, ok := data.([]string)
	if !ok {
		return &agents.AgentResult{
			AgentName: b.name,
			Status:    constants.StatusFailed,
			Errors:    []string{"invalid input data type, expected []string"},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	report, err := b.CheckLinks(urls)
	if err != nil {
		return &agents.AgentResult{
			AgentName: b.name,
			Status:    constants.StatusFailed,
			Errors:    []string{err.Error()},
			Duration:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	return &agents.AgentResult{
		AgentName: b.name,
		Status:    constants.StatusCompleted,
		Data: map[string]interface{}{
			"broken_links_report": report,
		},
		Duration: time.Since(startTime).Milliseconds(),
	}, nil
}

// HealthCheck vérifie la santé de l'agent
func (b *BrokenLinksDetector) HealthCheck() error {
	// Test simple de validation avec une URL connue
	_, err := b.ValidateLink("https://httpstat.us/200")
	return err
}

// CheckLinks vérifie une liste d'URLs et retourne un rapport des liens brisés
func (b *BrokenLinksDetector) CheckLinks(urls []string) (*agents.BrokenLinksReport, error) {
	if len(urls) == 0 {
		return &agents.BrokenLinksReport{
			TotalChecked: 0,
			BrokenCount:  0,
			BrokenLinks:  []agents.BrokenLink{},
			CheckedAt:    time.Now().Format(time.RFC3339),
		}, nil
	}

	// Canal pour distribuer le travail
	urlChan := make(chan string, len(urls))
	resultChan := make(chan linkCheckResult, len(urls))

	// Démarrer les workers
	var wg sync.WaitGroup
	workers := b.maxWorkers
	if len(urls) < workers {
		workers = len(urls)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go b.worker(&wg, urlChan, resultChan)
	}

	// Envoyer les URLs à vérifier
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Attendre que tous les workers finissent
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collecter les résultats
	var brokenLinks []agents.BrokenLink
	totalChecked := 0

	for result := range resultChan {
		totalChecked++
		if !result.IsValid {
			brokenLinks = append(brokenLinks, agents.BrokenLink{
				URL:        result.URL,
				StatusCode: result.StatusCode,
				Error:      result.Error,
				FoundOn:    []string{}, // À implémenter si nécessaire
			})
		}
	}

	return &agents.BrokenLinksReport{
		TotalChecked: totalChecked,
		BrokenCount:  len(brokenLinks),
		BrokenLinks:  brokenLinks,
		CheckedAt:    time.Now().Format(time.RFC3339),
	}, nil
}

// ValidateLink valide un seul lien
func (b *BrokenLinksDetector) ValidateLink(url string) (*agents.LinkStatus, error) {
	if url == "" {
		return &agents.LinkStatus{
			URL:       url,
			IsValid:   false,
			Error:     "empty URL",
			CheckedAt: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Créer la requête HEAD pour tester l'URL
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return &agents.LinkStatus{
			URL:       url,
			IsValid:   false,
			Error:     fmt.Sprintf("invalid URL: %v", err),
			CheckedAt: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Ajouter des headers pour éviter d'être bloqué
	req.Header.Set("User-Agent", "Fire-Salamander-Link-Checker/1.0")
	req.Header.Set("Accept", "*/*")

	// Effectuer la requête
	resp, err := b.client.Do(req)
	if err != nil {
		return &agents.LinkStatus{
			URL:        url,
			StatusCode: 0,
			IsValid:    false,
			Error:      fmt.Sprintf("request failed: %v", err),
			CheckedAt:  time.Now().Format(time.RFC3339),
		}, nil
	}
	defer resp.Body.Close()

	// Déterminer si le lien est valide basé sur le status code
	isValid := b.isValidStatusCode(resp.StatusCode)
	errorMsg := ""
	if !isValid {
		errorMsg = fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status)
	}

	return &agents.LinkStatus{
		URL:        url,
		StatusCode: resp.StatusCode,
		IsValid:    isValid,
		Error:      errorMsg,
		CheckedAt:  time.Now().Format(time.RFC3339),
	}, nil
}

// linkCheckResult représente le résultat d'une vérification de lien
type linkCheckResult struct {
	URL        string
	StatusCode int
	IsValid    bool
	Error      string
}

// worker traite les URLs depuis un canal
func (b *BrokenLinksDetector) worker(wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- linkCheckResult) {
	defer wg.Done()

	for url := range urlChan {
		status, _ := b.ValidateLink(url) // Ignorer l'erreur car elle est incluse dans LinkStatus
		
		result := linkCheckResult{
			URL:        status.URL,
			StatusCode: status.StatusCode,
			IsValid:    status.IsValid,
			Error:      status.Error,
		}
		
		resultChan <- result
	}
}

// isValidStatusCode détermine si un code de statut HTTP est considéré comme valide
func (b *BrokenLinksDetector) isValidStatusCode(statusCode int) bool {
	// Status codes considérés comme valides
	switch {
	case statusCode >= 200 && statusCode < 300: // Success
		return true
	case statusCode == 301 || statusCode == 302: // Redirections permanentes/temporaires
		return true
	case statusCode == 304: // Not Modified
		return true
	case statusCode == 401: // Unauthorized (peut être normal pour certaines ressources protégées)
		return true
	case statusCode == 403: // Forbidden (ressource existe mais accès interdit)
		return true
	default:
		return false
	}
}

// SetMaxWorkers configure le nombre maximum de workers pour les vérifications parallèles
func (b *BrokenLinksDetector) SetMaxWorkers(workers int) {
	if workers > 0 && workers <= 50 { // Limite raisonnable
		b.maxWorkers = workers
	}
}

// SetTimeout configure le timeout des requêtes HTTP
func (b *BrokenLinksDetector) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		b.client.Timeout = timeout
	}
}

// GetStats retourne des statistiques sur la configuration actuelle
func (b *BrokenLinksDetector) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"max_workers":    b.maxWorkers,
		"timeout":        b.client.Timeout.String(),
		"agent_name":     b.name,
		"transport_type": "http",
	}
}