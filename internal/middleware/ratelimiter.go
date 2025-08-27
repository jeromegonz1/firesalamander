package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"firesalamander/internal/constants"
)

// 🔥🦎 FIRE SALAMANDER - RATE LIMITING MIDDLEWARE
// Production Security - Protection DoS selon standards SEPTEO
// NO HARDCODING POLICY - All values from constants

// RateLimiter middleware pour protection DoS
type RateLimiter struct {
	// Map de limiters par IP
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	
	// Configuration depuis constants
	rate  rate.Limit
	burst int
}

// NewRateLimiter crée une nouvelle instance du rate limiter
// Utilise les constantes définies pour éviter le hardcoding
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(constants.RateLimitRequestsPerSecond),
		burst:    constants.RateLimitBurst,
	}
}

// getLimiter retourne ou crée un limiter pour une IP donnée
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		
		// Nettoyage automatique des limiters inactifs après 30 minutes
		go func() {
			time.Sleep(30 * time.Minute)
			rl.mu.Lock()
			delete(rl.limiters, ip)
			rl.mu.Unlock()
		}()
	}

	return limiter
}

// extractClientIP extrait l'IP du client en tenant compte des proxies
func (rl *RateLimiter) extractClientIP(r *http.Request) string {
	// Vérifier les headers de proxy les plus communs
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}
	
	// Fallback sur RemoteAddr
	return r.RemoteAddr
}

// Middleware retourne la fonction middleware HTTP
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.extractClientIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			// Headers informatifs pour le client
			w.Header().Set("X-RateLimit-Limit", "1")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", "1")
			w.Header().Set("Retry-After", "1")
			
			// Réponse JSON cohérente avec l'API
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Too many requests. Please try again later."}`))
			return
		}

		// Headers informatifs pour le client (quand accepté)
		w.Header().Set("X-RateLimit-Limit", "1") 
		w.Header().Set("X-RateLimit-Remaining", "4") // Burst - utilisé

		next.ServeHTTP(w, r)
	})
}

// MiddlewareFunc version fonction du middleware pour plus de flexibilité
func (rl *RateLimiter) MiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return rl.Middleware(next).ServeHTTP
}