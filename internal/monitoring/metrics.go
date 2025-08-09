package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 🔥🦎 FIRE SALAMANDER - MONITORING TEMPS RÉEL
// NOUVEAU PROCESS V2.0 - SURVEILLANCE ANTI-BOUCLE

// GlobalMetrics - Métriques globales thread-safe
type GlobalMetrics struct {
	mu                    sync.RWMutex
	startTime            time.Time
	totalRequests        int64
	activeAnalyses       int64
	completedAnalyses    int64
	failedAnalyses       int64
	infiniteLoopsDetected int64
	timeoutsOccurred     int64
	maxResponseTimeMs    int64
	totalResponseTimeMs  int64
	urlsProcessed        int64
	uniqueURLs           map[string]int
	lastActivityTime     time.Time
	alerts               []Alert
}

// Alert - Structure d'alerte
type Alert struct {
	Level     string    `json:"level"`     // INFO, WARN, ERROR, CRITICAL
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Component string    `json:"component"`
}

// PerformanceSnapshot - Photo instantanée des performances
type PerformanceSnapshot struct {
	Timestamp           time.Time `json:"timestamp"`
	Goroutines          int       `json:"goroutines"`
	MemoryUsageMB       float64   `json:"memory_usage_mb"`
	CPUUsagePercent     float64   `json:"cpu_usage_percent"`
	ActiveCrawls        int64     `json:"active_crawls"`
	URLsProcessed       int64     `json:"urls_processed"`
	RequestsPerSecond   float64   `json:"requests_per_second"`
	AverageResponseTime float64   `json:"avg_response_time_ms"`
	Alerts              []Alert   `json:"alerts"`
	HealthStatus        string    `json:"health_status"`
}

var globalMetrics = &GlobalMetrics{
	startTime:  time.Now(),
	uniqueURLs: make(map[string]int),
	alerts:     make([]Alert, 0),
}

// IncrementRequests - Incrémenter compteur de requêtes
func IncrementRequests() {
	atomic.AddInt64(&globalMetrics.totalRequests, 1)
	globalMetrics.mu.Lock()
	globalMetrics.lastActivityTime = time.Now()
	globalMetrics.mu.Unlock()
}

// IncrementActiveAnalyses - Incrémenter analyses actives
func IncrementActiveAnalyses() {
	atomic.AddInt64(&globalMetrics.activeAnalyses, 1)
}

// DecrementActiveAnalyses - Décrémenter analyses actives
func DecrementActiveAnalyses() {
	atomic.AddInt64(&globalMetrics.activeAnalyses, -1)
}

// IncrementCompletedAnalyses - Incrémenter analyses complétées
func IncrementCompletedAnalyses() {
	atomic.AddInt64(&globalMetrics.completedAnalyses, 1)
	DecrementActiveAnalyses()
}

// IncrementFailedAnalyses - Incrémenter analyses échouées
func IncrementFailedAnalyses() {
	atomic.AddInt64(&globalMetrics.failedAnalyses, 1)
	DecrementActiveAnalyses()
}

// ReportInfiniteLoop - Signaler une boucle infinie détectée
func ReportInfiniteLoop(url string) {
	atomic.AddInt64(&globalMetrics.infiniteLoopsDetected, 1)
	AddAlert("CRITICAL", fmt.Sprintf("Boucle infinie détectée pour: %s", url), "CRAWLER")
}

// ReportTimeout - Signaler un timeout
func ReportTimeout(component string) {
	atomic.AddInt64(&globalMetrics.timeoutsOccurred, 1)
	AddAlert("WARN", fmt.Sprintf("Timeout dans: %s", component), component)
}

// RecordResponseTime - Enregistrer temps de réponse
func RecordResponseTime(responseTimeMs int64) {
	atomic.AddInt64(&globalMetrics.totalResponseTimeMs, responseTimeMs)
	
	// Mettre à jour le temps max (thread-safe)
	for {
		currentMax := atomic.LoadInt64(&globalMetrics.maxResponseTimeMs)
		if responseTimeMs <= currentMax {
			break
		}
		if atomic.CompareAndSwapInt64(&globalMetrics.maxResponseTimeMs, currentMax, responseTimeMs) {
			break
		}
	}
}

// AddURLProcessed - Ajouter une URL traitée
func AddURLProcessed(url string) {
	atomic.AddInt64(&globalMetrics.urlsProcessed, 1)
	
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	
	globalMetrics.uniqueURLs[url]++
	
	// Détecter URLs trop souvent crawlées
	if globalMetrics.uniqueURLs[url] > 3 {
		AddAlert("ERROR", fmt.Sprintf("URL crawlée %d fois: %s", globalMetrics.uniqueURLs[url], url), "ANTI-LOOP")
	}
}

// AddAlert - Ajouter une alerte
func AddAlert(level, message, component string) {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()
	
	alert := Alert{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Component: component,
	}
	
	globalMetrics.alerts = append(globalMetrics.alerts, alert)
	
	// Garder seulement les 50 dernières alertes
	if len(globalMetrics.alerts) > 50 {
		globalMetrics.alerts = globalMetrics.alerts[1:]
	}
	
	// Log critique
	if level == "CRITICAL" {
		fmt.Printf("🚨 CRITICAL ALERT: %s [%s]\n", message, component)
	}
}

// GetCurrentSnapshot - Obtenir photo instantanée
func GetCurrentSnapshot() PerformanceSnapshot {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	totalRequests := atomic.LoadInt64(&globalMetrics.totalRequests)
	totalResponseTime := atomic.LoadInt64(&globalMetrics.totalResponseTimeMs)
	
	var avgResponseTime float64
	if totalRequests > 0 {
		avgResponseTime = float64(totalResponseTime) / float64(totalRequests)
	}
	
	uptime := time.Since(globalMetrics.startTime)
	requestsPerSecond := float64(totalRequests) / uptime.Seconds()
	
	globalMetrics.mu.RLock()
	alertsCopy := make([]Alert, len(globalMetrics.alerts))
	copy(alertsCopy, globalMetrics.alerts)
	globalMetrics.mu.RUnlock()
	
	// Déterminer le statut de santé
	healthStatus := determineHealthStatus()
	
	return PerformanceSnapshot{
		Timestamp:           time.Now(),
		Goroutines:          runtime.NumGoroutine(),
		MemoryUsageMB:       float64(memStats.Alloc) / 1024 / 1024,
		ActiveCrawls:        atomic.LoadInt64(&globalMetrics.activeAnalyses),
		URLsProcessed:       atomic.LoadInt64(&globalMetrics.urlsProcessed),
		RequestsPerSecond:   requestsPerSecond,
		AverageResponseTime: avgResponseTime,
		Alerts:              alertsCopy,
		HealthStatus:        healthStatus,
	}
}

// determineHealthStatus - Déterminer le statut de santé
func determineHealthStatus() string {
	// Vérifications critiques
	if atomic.LoadInt64(&globalMetrics.infiniteLoopsDetected) > 0 {
		return "CRITICAL"
	}
	
	if atomic.LoadInt64(&globalMetrics.activeAnalyses) > 10 {
		return "DEGRADED"
	}
	
	if runtime.NumGoroutine() > 100 {
		return "WARNING"
	}
	
	// Vérifier ratio succès/échec
	completed := atomic.LoadInt64(&globalMetrics.completedAnalyses)
	failed := atomic.LoadInt64(&globalMetrics.failedAnalyses)
	
	if completed+failed > 0 {
		successRate := float64(completed) / float64(completed+failed)
		if successRate < 0.5 {
			return "DEGRADED"
		}
	}
	
	return "HEALTHY"
}

// checkSystemHealth - Vérifications de santé automatiques
func checkSystemHealth() []string {
	alerts := make([]string, 0)
	
	// 1. Vérifier les boucles infinies
	if loops := atomic.LoadInt64(&globalMetrics.infiniteLoopsDetected); loops > 0 {
		alerts = append(alerts, fmt.Sprintf("🔄 %d boucles infinies détectées", loops))
	}
	
	// 2. Vérifier les fuites mémoire
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryMB := float64(memStats.Alloc) / 1024 / 1024
	
	if memoryMB > 500 { // Plus de 500MB
		alerts = append(alerts, fmt.Sprintf("💾 Fuite mémoire possible: %.1fMB", memoryMB))
	}
	
	// 3. Vérifier les deadlocks
	if goroutines := runtime.NumGoroutine(); goroutines > 50 {
		alerts = append(alerts, fmt.Sprintf("🔒 Possible deadlock: %d goroutines", goroutines))
	}
	
	// 4. Vérifier activité récente
	globalMetrics.mu.RLock()
	lastActivity := globalMetrics.lastActivityTime
	globalMetrics.mu.RUnlock()
	
	if time.Since(lastActivity) > 5*time.Minute {
		alerts = append(alerts, "⏰ Pas d'activité récente (possible freeze)")
	}
	
	return alerts
}

// MetricsHandler - Endpoint HTTP pour les métriques
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	snapshot := GetCurrentSnapshot()
	
	// Ajouter vérifications système
	systemAlerts := checkSystemHealth()
	
	response := map[string]interface{}{
		"performance":        snapshot,
		"system_alerts":      systemAlerts,
		"uptime_seconds":     int(time.Since(globalMetrics.startTime).Seconds()),
		"total_requests":     atomic.LoadInt64(&globalMetrics.totalRequests),
		"active_analyses":    atomic.LoadInt64(&globalMetrics.activeAnalyses), // 🔥🦎 SPRINT 6: Fix null value
		"completed_analyses": atomic.LoadInt64(&globalMetrics.completedAnalyses),
		"failed_analyses":    atomic.LoadInt64(&globalMetrics.failedAnalyses),
		"infinite_loops":     atomic.LoadInt64(&globalMetrics.infiniteLoopsDetected),
		"timeouts":           atomic.LoadInt64(&globalMetrics.timeoutsOccurred),
		"max_response_time_ms": atomic.LoadInt64(&globalMetrics.maxResponseTimeMs),
		"fire_salamander":    "v2.0-safety",
	}
	
	json.NewEncoder(w).Encode(response)
}

// HealthHandler - Endpoint de santé simple
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	healthStatus := determineHealthStatus()
	
	statusCode := http.StatusOK
	if healthStatus == "CRITICAL" {
		statusCode = http.StatusServiceUnavailable
	} else if healthStatus == "DEGRADED" || healthStatus == "WARNING" {
		statusCode = http.StatusPartialContent
	}
	
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             healthStatus,
		"timestamp":          time.Now(),
		"uptime_seconds":     int(time.Since(globalMetrics.startTime).Seconds()),
		"active_analyses":    atomic.LoadInt64(&globalMetrics.activeAnalyses),
		"infinite_loops":     atomic.LoadInt64(&globalMetrics.infiniteLoopsDetected),
		"goroutines":         runtime.NumGoroutine(),
	})
}

// ResetMetrics - Reset des métriques (pour tests)
func ResetMetrics() {
	atomic.StoreInt64(&globalMetrics.totalRequests, 0)
	atomic.StoreInt64(&globalMetrics.activeAnalyses, 0)
	atomic.StoreInt64(&globalMetrics.completedAnalyses, 0)
	atomic.StoreInt64(&globalMetrics.failedAnalyses, 0)
	atomic.StoreInt64(&globalMetrics.infiniteLoopsDetected, 0)
	atomic.StoreInt64(&globalMetrics.timeoutsOccurred, 0)
	atomic.StoreInt64(&globalMetrics.maxResponseTimeMs, 0)
	atomic.StoreInt64(&globalMetrics.totalResponseTimeMs, 0)
	atomic.StoreInt64(&globalMetrics.urlsProcessed, 0)
	
	globalMetrics.mu.Lock()
	globalMetrics.startTime = time.Now()
	globalMetrics.uniqueURLs = make(map[string]int)
	globalMetrics.alerts = make([]Alert, 0)
	globalMetrics.lastActivityTime = time.Now()
	globalMetrics.mu.Unlock()
}