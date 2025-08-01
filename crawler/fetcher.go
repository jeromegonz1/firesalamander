package crawler

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jeromegonz1/firesalamander/internal/logger"
)

var fetcherLog = logger.New("FETCHER")

// Fetcher gère les requêtes HTTP avec optimisations
type Fetcher struct {
	client        *http.Client
	config        *Config
	retryStrategy RetryStrategy
}

// RetryStrategy définit la stratégie de retry
type RetryStrategy struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryStrategy retourne une stratégie de retry par défaut
func DefaultRetryStrategy() RetryStrategy {
	return RetryStrategy{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

// NewFetcher crée un nouveau fetcher HTTP optimisé
func NewFetcher(config *Config) *Fetcher {
	fetcherLog.Debug("Creating new fetcher", map[string]interface{}{
		"timeout":       config.Timeout,
		"retry_attempts": config.RetryAttempts,
		"user_agent":    config.UserAgent,
	})

	// Transport optimisé avec connection pooling
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false, // Nous gérons la compression manuellement
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			// Conserver le User-Agent lors des redirections
			req.Header.Set("User-Agent", config.UserAgent)
			return nil
		},
	}

	return &Fetcher{
		client: client,
		config: config,
		retryStrategy: RetryStrategy{
			MaxAttempts:  config.RetryAttempts,
			InitialDelay: config.RetryDelay,
			MaxDelay:     30 * time.Second,
			Multiplier:   2.0,
		},
	}
}

// Fetch récupère une page avec retry et optimisations
func (f *Fetcher) Fetch(ctx context.Context, targetURL string) (*CrawlResult, error) {
	start := time.Now()
	fetcherLog.Debug("Fetching URL", map[string]interface{}{"url": targetURL})

	var lastErr error
	delay := f.retryStrategy.InitialDelay

	for attempt := 1; attempt <= f.retryStrategy.MaxAttempts; attempt++ {
		result, err := f.doFetch(ctx, targetURL)
		
		if err == nil {
			result.ResponseTime = time.Since(start)
			fetcherLog.Debug("Fetch successful", map[string]interface{}{
				"url":          targetURL,
				"status_code":  result.StatusCode,
				"content_type": result.ContentType,
				"response_time": result.ResponseTime,
				"attempt":      attempt,
			})
			return result, nil
		}

		lastErr = err
		
		// Ne pas retry si c'est une erreur client (4xx)
		if result != nil && result.StatusCode >= 400 && result.StatusCode < 500 {
			fetcherLog.Debug("Client error, not retrying", map[string]interface{}{
				"url":         targetURL,
				"status_code": result.StatusCode,
			})
			result.Error = err
			return result, err
		}

		// Si ce n'est pas la dernière tentative, attendre avant de réessayer
		if attempt < f.retryStrategy.MaxAttempts {
			fetcherLog.Warn("Fetch failed, retrying", map[string]interface{}{
				"url":     targetURL,
				"attempt": attempt,
				"error":   err.Error(),
				"delay":   delay,
			})
			
			select {
			case <-time.After(delay):
				// Augmenter le délai pour la prochaine tentative (backoff exponentiel)
				delay = time.Duration(float64(delay) * f.retryStrategy.Multiplier)
				if delay > f.retryStrategy.MaxDelay {
					delay = f.retryStrategy.MaxDelay
				}
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	fetcherLog.Error("All fetch attempts failed", map[string]interface{}{
		"url":      targetURL,
		"attempts": f.retryStrategy.MaxAttempts,
		"error":    lastErr.Error(),
	})

	return &CrawlResult{
		URL:       targetURL,
		Error:     lastErr,
		CrawledAt: time.Now(),
	}, lastErr
}

// doFetch effectue une requête HTTP unique
func (f *Fetcher) doFetch(ctx context.Context, targetURL string) (*CrawlResult, error) {
	// Créer la requête
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Headers optimisés pour un crawler mobile-first
	req.Header.Set("User-Agent", f.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Effectuer la requête
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Créer le résultat
	result := &CrawlResult{
		URL:         targetURL,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Headers:     make(map[string]string),
		CrawledAt:   time.Now(),
	}

	// Copier les headers importants
	importantHeaders := []string{
		"Content-Type", "Content-Length", "Last-Modified", 
		"ETag", "Cache-Control", "X-Robots-Tag", 
		"Content-Language", "Content-Encoding",
	}
	for _, header := range importantHeaders {
		if value := resp.Header.Get(header); value != "" {
			result.Headers[header] = value
		}
	}

	// Vérifier le status code
	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Lire le body avec gestion de la compression
	body, err := f.readBody(resp)
	if err != nil {
		return result, fmt.Errorf("failed to read body: %w", err)
	}

	result.Body = string(body)

	// Limiter la taille du body stocké (max 10MB)
	maxBodySize := 10 * 1024 * 1024
	if len(result.Body) > maxBodySize {
		fetcherLog.Warn("Body truncated", map[string]interface{}{
			"url":           targetURL,
			"original_size": len(result.Body),
			"max_size":      maxBodySize,
		})
		result.Body = result.Body[:maxBodySize]
	}

	return result, nil
}

// readBody lit le body de la réponse avec décompression si nécessaire
func (f *Fetcher) readBody(resp *http.Response) ([]byte, error) {
	var reader io.Reader = resp.Body

	// Gérer la compression gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// Limiter la lecture à 10MB pour éviter les abus
	limitedReader := io.LimitReader(reader, 10*1024*1024)
	
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// FetchWithMethod effectue une requête avec une méthode HTTP spécifique
func (f *Fetcher) FetchWithMethod(ctx context.Context, method, targetURL string, body io.Reader) (*CrawlResult, error) {
	req, err := http.NewRequestWithContext(ctx, method, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", f.config.UserAgent)
	
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &CrawlResult{
		URL:         targetURL,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		CrawledAt:   time.Now(),
	}

	if resp.StatusCode < 400 {
		bodyBytes, err := f.readBody(resp)
		if err == nil {
			result.Body = string(bodyBytes)
		}
	}

	return result, nil
}

// IsHTML vérifie si le content-type est HTML
func IsHTML(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "text/html") || 
		   strings.Contains(contentType, "application/xhtml+xml")
}

// IsXML vérifie si le content-type est XML
func IsXML(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "text/xml") || 
		   strings.Contains(contentType, "application/xml")
}

// Close ferme les connexions du fetcher
func (f *Fetcher) Close() {
	f.client.CloseIdleConnections()
}