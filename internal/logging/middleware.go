package logging

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - LOGGING MIDDLEWARE
// TDD + Zero Hardcoding Policy
// ========================================

// responseWriter wrapper pour capturer les informations de réponse
type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int64
}

// NewResponseWriter crée un nouveau wrapper de ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default 200
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += int64(n)
	return n, err
}

// HTTPLoggingMiddleware middleware pour logging des requêtes HTTP
func HTTPLoggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			// Générer ou récupérer le Request ID
			requestID := getOrGenerateRequestID(r)
			
			// Ajouter le Request ID aux headers de réponse
			w.Header().Set(constants.HeaderXRequestID, requestID)
			
			// Wrapper pour capturer la réponse
			wrapped := NewResponseWriter(w)
			
			// Log de la requête entrante
			httpLogger := logger.HTTP()
			httpLogger.RequestReceived(
				r.Method,
				r.URL.String(),
				getClientIP(r),
				r.UserAgent(),
				requestID,
			)
			
			// Traiter la requête
			next.ServeHTTP(wrapped, r)
			
			// Calculer le temps de réponse
			duration := time.Since(startTime)
			responseTimeMs := duration.Nanoseconds() / 1000000
			
			// Log de la réponse
			httpLogger.ResponseSent(
				r.Method,
				r.URL.String(),
				wrapped.statusCode,
				responseTimeMs,
				wrapped.contentLength,
				requestID,
			)
			
			// Log de performance si requête lente
			if responseTimeMs > 1000 { // Plus d'1 seconde
				perfLogger := logger.Performance()
				perfLogger.OperationCompleted(
					"http_request_slow",
					responseTimeMs,
					map[string]interface{}{
						"method": r.Method,
						"url":    r.URL.String(),
						"status": wrapped.statusCode,
					},
				)
			}
		})
	}
}

// APILoggingMiddleware middleware pour logging des requêtes API
func APILoggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ne traiter que les requêtes API
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}
			
			startTime := time.Now()
			requestID := getOrGenerateRequestID(r)
			
			// Wrapper pour capturer la réponse
			wrapped := NewResponseWriter(w)
			
			// Log de la requête API entrante
			apiLogger := logger.API()
			apiLogger.RequestReceived(
				r.URL.Path,
				r.Method,
				requestID,
				map[string]interface{}{
					"content_type":   r.Header.Get("Content-Type"),
					"content_length": r.ContentLength,
					"user_agent":     r.UserAgent(),
					"remote_addr":    getClientIP(r),
				},
			)
			
			// Traiter la requête
			next.ServeHTTP(wrapped, r)
			
			// Calculer le temps de réponse
			duration := time.Since(startTime)
			responseTimeMs := duration.Nanoseconds() / 1000000
			
			// Log de la réponse API
			apiLogger.ResponseSent(
				r.URL.Path,
				r.Method,
				wrapped.statusCode,
				responseTimeMs,
				requestID,
			)
		})
	}
}

// RecoveryMiddleware middleware pour capturer les panics et les logger
func RecoveryMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := getOrGenerateRequestID(r)
					
					logger.Error(
						constants.LogCategorySystem,
						"Panic recovered in HTTP handler",
						err.(error),
						map[string]interface{}{
							"method":     r.Method,
							"url":        r.URL.String(),
							"request_id": requestID,
							"panic":      err,
						},
					)
					
					// Retourner une erreur 500
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// getOrGenerateRequestID récupère ou génère un Request ID
func getOrGenerateRequestID(r *http.Request) string {
	// Vérifier s'il y a déjà un Request ID dans les headers
	if requestID := r.Header.Get(constants.HeaderXRequestID); requestID != "" {
		return requestID
	}
	
	// Générer un nouveau Request ID
	return generateRequestID()
}

// generateRequestID génère un Request ID unique
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return constants.RequestIDPrefix + hex.EncodeToString(bytes)
}

// getClientIP extrait l'adresse IP du client
func getClientIP(r *http.Request) string {
	// Vérifier les headers de proxy
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// Prendre la première IP si plusieurs sont présentes
		if ips := strings.Split(ip, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	
	// Fallback sur RemoteAddr
	if ip := strings.Split(r.RemoteAddr, ":")[0]; ip != "" {
		return ip
	}
	
	return "unknown"
}

// LogRequestContext ajoute les informations de contexte à la requête
type LogRequestContext struct {
	RequestID string
	TraceID   string
	StartTime time.Time
}

// NewLogRequestContext crée un nouveau contexte de requête
func NewLogRequestContext() *LogRequestContext {
	return &LogRequestContext{
		RequestID: generateRequestID(),
		TraceID:   generateTraceID(),
		StartTime: time.Now(),
	}
}

// generateTraceID génère un Trace ID unique
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return constants.TraceIDPrefix + hex.EncodeToString(bytes)
}

// MetricsMiddleware middleware pour collecter des métriques de performance
func MetricsMiddleware(logger Logger) func(http.Handler) http.Handler {
	var (
		requestCount   int64
		totalDuration  int64
		lastMetricsLog time.Time
	)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			next.ServeHTTP(w, r)
			
			// Incrémenter les compteurs
			requestCount++
			duration := time.Since(startTime)
			totalDuration += duration.Nanoseconds() / 1000000
			
			// Logger les métriques toutes les minutes
			if time.Since(lastMetricsLog) > time.Minute {
				perfLogger := logger.Performance()
				avgResponseTime := totalDuration / requestCount
				requestsPerSec := float64(requestCount) / time.Since(lastMetricsLog).Seconds()
				
				perfLogger.RequestMetrics(requestsPerSec, avgResponseTime)
				
				// Reset des compteurs
				requestCount = 0
				totalDuration = 0
				lastMetricsLog = time.Now()
			}
		})
	}
}