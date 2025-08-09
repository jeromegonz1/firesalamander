package logging

import (
	"runtime"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - PERFORMANCE LOGGER
// TDD + Zero Hardcoding Policy
// ========================================

// performanceLogger implémentation du PerformanceLogger
type performanceLogger struct {
	logger *FireSalamanderLogger
	operations map[string]time.Time // Track des opérations en cours
}

// NewPerformanceLogger crée un nouveau logger de performance
func NewPerformanceLogger(logger *FireSalamanderLogger) PerformanceLogger {
	return &performanceLogger{
		logger: logger,
		operations: make(map[string]time.Time),
	}
}

// OperationStarted log le début d'une opération
func (p *performanceLogger) OperationStarted(operation string, data map[string]interface{}) {
	logData := map[string]interface{}{
		constants.PerfLogFieldOperation: operation,
		"status": "started",
	}
	
	// Ajouter les données supplémentaires
	if data != nil {
		for k, v := range data {
			logData[k] = v
		}
	}
	
	// Enregistrer le temps de début
	p.operations[operation] = time.Now()
	
	p.logger.Debug(constants.LogCategoryPerformance, "Operation started", logData)
}

// OperationCompleted log la fin d'une opération
func (p *performanceLogger) OperationCompleted(operation string, durationMs int64, data map[string]interface{}) {
	logData := map[string]interface{}{
		constants.PerfLogFieldOperation: operation,
		constants.PerfLogFieldDurationMs: durationMs,
		"status": "completed",
	}
	
	// Ajouter les données supplémentaires
	if data != nil {
		for k, v := range data {
			logData[k] = v
		}
	}
	
	// Nettoyer le tracking
	delete(p.operations, operation)
	
	// Choisir le niveau selon la durée
	if durationMs > 5000 { // Plus de 5 secondes = warning
		p.logger.Warn(constants.LogCategoryPerformance, "Operation completed (slow)", logData)
	} else if durationMs > 1000 { // Plus d'1 seconde = info
		p.logger.Info(constants.LogCategoryPerformance, "Operation completed", logData)
	} else { // Rapide = debug
		p.logger.Debug(constants.LogCategoryPerformance, "Operation completed", logData)
	}
}

// MemoryUsage log l'utilisation mémoire
func (p *performanceLogger) MemoryUsage(operation string, memoryBytes int64, goroutines int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	data := map[string]interface{}{
		constants.PerfLogFieldOperation:      operation,
		constants.PerfLogFieldMemoryUsage:   memoryBytes,
		constants.PerfLogFieldGoroutines:    goroutines,
		"heap_alloc_bytes":                  m.HeapAlloc,
		"heap_sys_bytes":                    m.HeapSys,
		"heap_inuse_bytes":                  m.HeapInuse,
		"stack_inuse_bytes":                 m.StackInuse,
		"next_gc_bytes":                     m.NextGC,
		"num_gc":                            m.NumGC,
	}
	
	// Warning si utilisation mémoire élevée
	if memoryBytes > 100*1024*1024 { // Plus de 100MB
		p.logger.Warn(constants.LogCategoryPerformance, "High memory usage detected", data)
	} else {
		p.logger.Debug(constants.LogCategoryPerformance, "Memory usage", data)
	}
}

// RequestMetrics log les métriques de requêtes
func (p *performanceLogger) RequestMetrics(requestsPerSec float64, avgResponseTimeMs int64) {
	data := map[string]interface{}{
		constants.PerfLogFieldRequestsPerSec: requestsPerSec,
		"avg_response_time_ms":               avgResponseTimeMs,
		"timestamp":                          time.Now().Format(constants.LogTimestampFormat),
	}
	
	// Warning si performance dégradée
	if requestsPerSec < 1.0 || avgResponseTimeMs > 2000 {
		p.logger.Warn(constants.LogCategoryPerformance, "Performance degradation detected", data)
	} else {
		p.logger.Info(constants.LogCategoryPerformance, "Request metrics", data)
	}
}