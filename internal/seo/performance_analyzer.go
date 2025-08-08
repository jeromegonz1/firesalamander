package seo

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"firesalamander/internal/constants"
)

// PerformanceAnalyzer analyseur de performances SEO
type PerformanceAnalyzer struct {
	client           *http.Client
	compressionRegex *regexp.Regexp
	cacheRegex       *regexp.Regexp
	imageRegex       *regexp.Regexp
}

// PerformanceMetricsResult résultat des métriques de performance
type PerformanceMetricsResult struct {
	LoadTime         time.Duration     `json:"load_time"`
	PageSize         int64             `json:"page_size"`
	CompressedSize   int64             `json:"compressed_size,omitempty"`
	CompressionRatio float64           `json:"compression_ratio,omitempty"`
	
	// Optimisations
	HasCompression   bool              `json:"has_compression"`
	HasCaching       bool              `json:"has_caching"`
	OptimizedImages  bool              `json:"optimized_images"`
	MinifiedResources bool             `json:"minified_resources"`
	
	// Ressources
	ResourceCounts   ResourceCounts    `json:"resource_counts"`
	LargestResources []Resource        `json:"largest_resources"`
	
	// Core Web Vitals (estimés)
	CoreWebVitals    CoreWebVitals     `json:"core_web_vitals"`
	
	// Headers et optimisations
	HTTPHeaders      HTTPHeaderAnalysis `json:"http_headers"`
	
	Issues           []string          `json:"issues"`
	Recommendations  []string          `json:"recommendations"`
}

// ResourceCounts comptage des ressources
type ResourceCounts struct {
	Images      int `json:"images"`
	Scripts     int `json:"scripts"`
	Stylesheets int `json:"stylesheets"`
	Fonts       int `json:"fonts"`
	Videos      int `json:"videos"`
	Other       int `json:"other"`
}

// Resource ressource web
type Resource struct {
	URL          string        `json:"url"`
	Type         string        `json:"type"`
	Size         int64         `json:"size"`
	LoadTime     time.Duration `json:"load_time"`
	Compressed   bool          `json:"compressed"`
	Cached       bool          `json:"cached"`
}

// CoreWebVitals métriques Core Web Vitals
type CoreWebVitals struct {
	LCP          EstimatedMetric `json:"lcp"` // Largest Contentful Paint
	FID          EstimatedMetric `json:"fid"` // First Input Delay
	CLS          EstimatedMetric `json:"cls"` // Cumulative Layout Shift
	TTFB         time.Duration   `json:"ttfb"` // Time To First Byte
	FCP          EstimatedMetric `json:"fcp"` // First Contentful Paint
	SpeedIndex   EstimatedMetric `json:"speed_index"`
}

// EstimatedMetric métrique estimée
type EstimatedMetric struct {
	Value       float64 `json:"value"`
	Score       string  `json:"score"` // good, needs-improvement, poor
	Threshold   string  `json:"threshold"`
}

// HTTPHeaderAnalysis analyse des headers HTTP
type HTTPHeaderAnalysis struct {
	Compression      CompressionAnalysis `json:"compression"`
	Caching          CachingAnalysis     `json:"caching"`
	Security         SecurityHeaders     `json:"security"`
	Performance      PerformanceHeaders  `json:"performance"`
}

// CompressionAnalysis analyse de la compression
type CompressionAnalysis struct {
	Enabled         bool     `json:"enabled"`
	Algorithm       string   `json:"algorithm"`
	CompressionRate float64  `json:"compression_rate"`
}

// CachingAnalysis analyse du cache
type CachingAnalysis struct {
	HasCacheControl bool          `json:"has_cache_control"`
	CacheControl    string        `json:"cache_control"`
	HasETag         bool          `json:"has_etag"`
	ETag            string        `json:"etag"`
	HasLastModified bool          `json:"has_last_modified"`
	LastModified    string        `json:"last_modified"`
	MaxAge          time.Duration `json:"max_age"`
}

// SecurityHeaders headers de sécurité
type SecurityHeaders struct {
	HasHSTS         bool   `json:"has_hsts"`
	HasCSP          bool   `json:"has_csp"`
	HasXFrame       bool   `json:"has_x_frame"`
	HasXContentType bool   `json:"has_x_content_type"`
}

// PerformanceHeaders headers de performance
type PerformanceHeaders struct {
	HasKeepAlive     bool `json:"has_keep_alive"`
	HasHTTP2         bool `json:"has_http2"`
	HasPreload       bool `json:"has_preload"`
	HasPrefetch      bool `json:"has_prefetch"`
	HasResourceHints bool `json:"has_resource_hints"`
}

// NewPerformanceAnalyzer crée un nouvel analyseur de performance
func NewPerformanceAnalyzer() *PerformanceAnalyzer {
	client := &http.Client{
		Timeout: constants.ClientTimeout,
	}

	return &PerformanceAnalyzer{
		client:           client,
		compressionRegex: regexp.MustCompile(`gzip|deflate|br`),
		cacheRegex:       regexp.MustCompile(`max-age=(\d+)`),
		imageRegex:       regexp.MustCompile(`\.(webp|avif|jpg|jpeg|png|gif|svg)$`),
	}
}

// Analyze effectue l'analyse de performance
func (pa *PerformanceAnalyzer) Analyze(ctx context.Context, targetURL string, htmlContent string) (*PerformanceMetricsResult, error) {
	result := &PerformanceMetricsResult{
		Issues:          []string{},
		Recommendations: []string{},
	}

	// 1. Mesurer le temps de chargement et analyser les headers
	if err := pa.analyzePageLoad(ctx, targetURL, result); err != nil {
		return nil, err
	}

	// 2. Analyser le contenu HTML pour les ressources
	if err := pa.analyzeHTMLResources(htmlContent, result); err != nil {
		return nil, err
	}

	// 3. Estimer les Core Web Vitals
	pa.estimateCoreWebVitals(result)

	// 4. Générer les recommandations
	pa.generatePerformanceRecommendations(result)

	return result, nil
}

// analyzePageLoad analyse le chargement de la page
func (pa *PerformanceAnalyzer) analyzePageLoad(ctx context.Context, targetURL string, result *PerformanceMetricsResult) error {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return err
	}

	// Headers pour optimiser la requête
	req.Header.Set("User-Agent", "Fire Salamander Performance Bot/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	start := time.Now()
	resp, err := pa.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result.LoadTime = time.Since(start)

	// Analyser les headers HTTP
	pa.analyzeHTTPHeaders(resp, result)

	// Mesurer la taille de la page
	buf := make([]byte, 4*1024*1024) // 4MB max
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	result.PageSize = int64(n)

	// Vérifier la compression
	if encoding := resp.Header.Get("Content-Encoding"); encoding != "" {
		result.HasCompression = true
		result.HTTPHeaders.Compression.Enabled = true
		result.HTTPHeaders.Compression.Algorithm = encoding
	}

	return nil
}

// analyzeHTTPHeaders analyse les headers HTTP
func (pa *PerformanceAnalyzer) analyzeHTTPHeaders(resp *http.Response, result *PerformanceMetricsResult) {
	headers := &result.HTTPHeaders

	// Analyse de la compression
	if encoding := resp.Header.Get("Content-Encoding"); encoding != "" {
		headers.Compression.Enabled = true
		headers.Compression.Algorithm = encoding
	}

	// Analyse du cache
	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl != "" {
		headers.Caching.HasCacheControl = true
		headers.Caching.CacheControl = cacheControl
		result.HasCaching = true

		// Extraire max-age
		if matches := pa.cacheRegex.FindStringSubmatch(cacheControl); len(matches) > 1 {
			if maxAge, err := strconv.Atoi(matches[1]); err == nil {
				headers.Caching.MaxAge = time.Duration(maxAge) * time.Second
			}
		}
	}

	// ETag
	if etag := resp.Header.Get("ETag"); etag != "" {
		headers.Caching.HasETag = true
		headers.Caching.ETag = etag
		result.HasCaching = true
	}

	// Last-Modified
	if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
		headers.Caching.HasLastModified = true
		headers.Caching.LastModified = lastMod
	}

	// Headers de sécurité
	headers.Security.HasHSTS = resp.Header.Get("Strict-Transport-Security") != ""
	headers.Security.HasCSP = resp.Header.Get("Content-Security-Policy") != ""
	headers.Security.HasXFrame = resp.Header.Get("X-Frame-Options") != ""
	headers.Security.HasXContentType = resp.Header.Get("X-Content-Type-Options") != ""

	// Headers de performance
	headers.Performance.HasKeepAlive = resp.Header.Get("Connection") == "keep-alive"
	headers.Performance.HasHTTP2 = resp.ProtoMajor == 2
	headers.Performance.HasPreload = strings.Contains(resp.Header.Get("Link"), "preload")
}

// analyzeHTMLResources analyse les ressources dans le HTML
func (pa *PerformanceAnalyzer) analyzeHTMLResources(htmlContent string, result *PerformanceMetricsResult) error {
	// Compter les différents types de ressources
	result.ResourceCounts.Images = strings.Count(htmlContent, "<img")
	result.ResourceCounts.Scripts = strings.Count(htmlContent, "<script")
	result.ResourceCounts.Stylesheets = strings.Count(htmlContent, `<link rel="stylesheet"`)
	result.ResourceCounts.Fonts = strings.Count(htmlContent, ".woff") + strings.Count(htmlContent, ".ttf")

	// Vérifier l'optimisation des images
	webpCount := strings.Count(strings.ToLower(htmlContent), ".webp")
	avifCount := strings.Count(strings.ToLower(htmlContent), ".avif")
	
	if result.ResourceCounts.Images > 0 {
		optimizedRatio := float64(webpCount+avifCount) / float64(result.ResourceCounts.Images)
		result.OptimizedImages = optimizedRatio > 0.5
	}

	// Vérifier la minification (heuristique simple)
	result.MinifiedResources = pa.checkMinification(htmlContent)

	return nil
}

// checkMinification vérifie si les ressources sont minifiées
func (pa *PerformanceAnalyzer) checkMinification(htmlContent string) bool {
	// Heuristique simple : vérifier les espaces et retours à la ligne
	lines := strings.Split(htmlContent, "\n")
	shortLines := 0
	
	for _, line := range lines {
		if len(strings.TrimSpace(line)) < constants.OptimalLineLength {
			shortLines++
		}
	}

	// Si plus de 70% des lignes sont courtes, probablement pas minifié
	return float64(shortLines)/float64(len(lines)) < 0.7
}

// estimateCoreWebVitals estime les Core Web Vitals
func (pa *PerformanceAnalyzer) estimateCoreWebVitals(result *PerformanceMetricsResult) {
	// TTFB (Time To First Byte) - basé sur le load time
	result.CoreWebVitals.TTFB = result.LoadTime

	// LCP (Largest Contentful Paint) - estimation basée sur la taille et les ressources
	lcpValue := float64(result.LoadTime.Milliseconds())
	if result.ResourceCounts.Images > 5 {
		lcpValue *= 1.2
	}
	if result.PageSize > 1024*1024 { // > 1MB
		lcpValue *= 1.3
	}

	result.CoreWebVitals.LCP = EstimatedMetric{
		Value:     lcpValue,
		Score:     pa.scoreLCP(lcpValue),
		Threshold: "LCP ≤ 2.5s (good), ≤ 4.0s (needs improvement), > 4.0s (poor)",
	}

	// FID (First Input Delay) - estimation basée sur les scripts
	fidValue := float64(result.ResourceCounts.Scripts * 10) // 10ms par script
	result.CoreWebVitals.FID = EstimatedMetric{
		Value:     fidValue,
		Score:     pa.scoreFID(fidValue),
		Threshold: "FID ≤ 100ms (good), ≤ 300ms (needs improvement), > 300ms (poor)",
	}

	// CLS (Cumulative Layout Shift) - estimation basée sur les ressources
	clsValue := 0.0
	if result.ResourceCounts.Images > 0 && !result.OptimizedImages {
		clsValue += 0.1
	}
	if !result.MinifiedResources {
		clsValue += 0.05
	}

	result.CoreWebVitals.CLS = EstimatedMetric{
		Value:     clsValue,
		Score:     pa.scoreCLS(clsValue),
		Threshold: "CLS ≤ 0.1 (good), ≤ 0.25 (needs improvement), > 0.25 (poor)",
	}

	// FCP (First Contentful Paint)
	fcpValue := lcpValue * 0.6 // Estimation
	result.CoreWebVitals.FCP = EstimatedMetric{
		Value:     fcpValue,
		Score:     pa.scoreFCP(fcpValue),
		Threshold: "FCP ≤ 1.8s (good), ≤ 3.0s (needs improvement), > 3.0s (poor)",
	}

	// Speed Index
	speedIndexValue := lcpValue * 0.8
	result.CoreWebVitals.SpeedIndex = EstimatedMetric{
		Value:     speedIndexValue,
		Score:     pa.scoreSpeedIndex(speedIndexValue),
		Threshold: "SI ≤ 3.4s (good), ≤ 5.8s (needs improvement), > 5.8s (poor)",
	}
}

// Fonctions de scoring pour Core Web Vitals

func (pa *PerformanceAnalyzer) scoreLCP(value float64) string {
	if value <= 2500 {
		return "good"
	} else if value <= 4000 {
		return "needs-improvement"
	}
	return "poor"
}

func (pa *PerformanceAnalyzer) scoreFID(value float64) string {
	if value <= 100 {
		return "good"
	} else if value <= 300 {
		return "needs-improvement"
	}
	return "poor"
}

func (pa *PerformanceAnalyzer) scoreCLS(value float64) string {
	if value <= 0.1 {
		return "good"
	} else if value <= 0.25 {
		return "needs-improvement"
	}
	return "poor"
}

func (pa *PerformanceAnalyzer) scoreFCP(value float64) string {
	if value <= 1800 {
		return "good"
	} else if value <= constants.TestValue3000 {
		return "needs-improvement"
	}
	return "poor"
}

func (pa *PerformanceAnalyzer) scoreSpeedIndex(value float64) string {
	if value <= 3400 {
		return "good"
	} else if value <= 5800 {
		return "needs-improvement"
	}
	return "poor"
}

// generatePerformanceRecommendations génère les recommandations de performance
func (pa *PerformanceAnalyzer) generatePerformanceRecommendations(result *PerformanceMetricsResult) {
	// Temps de chargement
	if result.LoadTime > constants.AcceptableLoadTime {
		result.Issues = append(result.Issues, "Temps de chargement élevé")
		result.Recommendations = append(result.Recommendations, "Optimiser le temps de réponse du serveur")
	}

	// Taille de la page
	if result.PageSize > 2*1024*1024 { // > 2MB
		result.Issues = append(result.Issues, "Page trop volumineuse")
		result.Recommendations = append(result.Recommendations, "Réduire la taille de la page")
	}

	// Compression
	if !result.HasCompression {
		result.Issues = append(result.Issues, "Compression désactivée")
		result.Recommendations = append(result.Recommendations, "Activer la compression GZIP/Brotli")
	}

	// Cache
	if !result.HasCaching {
		result.Issues = append(result.Issues, "Cache non configuré")
		result.Recommendations = append(result.Recommendations, "Configurer les headers de cache")
	}

	// Images
	if !result.OptimizedImages && result.ResourceCounts.Images > 0 {
		result.Issues = append(result.Issues, "Images non optimisées")
		result.Recommendations = append(result.Recommendations, "Utiliser des formats d'images modernes (WebP, AVIF)")
	}

	// Minification
	if !result.MinifiedResources {
		result.Recommendations = append(result.Recommendations, "Minifier les ressources CSS et JavaScript")
	}

	// Core Web Vitals
	if result.CoreWebVitals.LCP.Score == "poor" {
		result.Issues = append(result.Issues, "LCP élevé")
		result.Recommendations = append(result.Recommendations, "Optimiser le Largest Contentful Paint")
	}

	if result.CoreWebVitals.FID.Score == "poor" {
		result.Issues = append(result.Issues, "FID élevé")
		result.Recommendations = append(result.Recommendations, "Réduire le JavaScript et optimiser l'interactivité")
	}

	if result.CoreWebVitals.CLS.Score == "poor" {
		result.Issues = append(result.Issues, "CLS élevé")
		result.Recommendations = append(result.Recommendations, "Stabiliser la mise en page et les décalages")
	}

	// Sécurité
	if !result.HTTPHeaders.Security.HasHSTS {
		result.Recommendations = append(result.Recommendations, "Ajouter le header HSTS pour la sécurité")
	}

	// HTTP/2
	if !result.HTTPHeaders.Performance.HasHTTP2 {
		result.Recommendations = append(result.Recommendations, "Migrer vers HTTP/2 pour de meilleures performances")
	}
}