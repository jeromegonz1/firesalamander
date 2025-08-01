package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jeromegonz1/firesalamander/config"
	"github.com/jeromegonz1/firesalamander/internal/logger"
)

type HealthCheck struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	App       AppInfo                `json:"app"`
	System    SystemInfo             `json:"system"`
	Config    ConfigInfo             `json:"config"`
	Checks    map[string]CheckResult `json:"checks"`
}

type AppInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Icon      string `json:"icon"`
	PoweredBy string `json:"powered_by"`
	Uptime    string `json:"uptime"`
}

type SystemInfo struct {
	Go       string `json:"go_version"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	CPUs     int    `json:"cpus"`
	Memory   string `json:"memory"`
	Debug    bool   `json:"debug_mode"`
}

type ConfigInfo struct {
	Environment string `json:"environment"`
	ServerPort  int    `json:"server_port"`
	Database    string `json:"database_type"`
	AIEnabled   bool   `json:"ai_enabled"`
}

type CheckResult struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var (
	startTime = time.Now()
	log       = logger.New("DEBUG")
)

func NewChecker(cfg *config.Config) *HealthCheck {
	log.Debug("Creating new health checker")
	
	checker := &HealthCheck{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		App: AppInfo{
			Name:      cfg.App.Name,
			Version:   cfg.App.Version,
			Icon:      cfg.App.Icon,
			PoweredBy: cfg.App.PoweredBy,
			Uptime:    time.Since(startTime).String(),
		},
		System: SystemInfo{
			Go:     runtime.Version(),
			OS:     runtime.GOOS,
			Arch:   runtime.GOARCH,
			CPUs:   runtime.NumCPU(),
			Memory: getMemoryUsage(),
			Debug:  logger.IsDebugMode(),
		},
		Config: ConfigInfo{
			Environment: getEnvironment(),
			ServerPort:  cfg.Server.Port,
			Database:    cfg.Database.Type,
			AIEnabled:   cfg.AI.Enabled,
		},
		Checks: make(map[string]CheckResult),
	}
	
	// Exécuter tous les checks
	checker.runAllChecks(cfg)
	
	return checker
}

func (hc *HealthCheck) runAllChecks(cfg *config.Config) {
	log.Debug("Running all health checks")
	
	checks := map[string]func(*config.Config) CheckResult{
		"config":     checkConfig,
		"database":   checkDatabase,
		"filesystem": func(cfg *config.Config) CheckResult { return checkFilesystem() },
		"network":    checkNetwork,
		"ai":         checkAI,
	}
	
	for name, checkFunc := range checks {
		log.Debug("Running check", map[string]interface{}{"check": name})
		result := checkFunc(cfg)
		hc.Checks[name] = result
		
		if result.Status != "ok" {
			log.Warn("Check failed", map[string]interface{}{
				"check":   name,
				"status":  result.Status,
				"message": result.Message,
				"error":   result.Error,
			})
			if hc.Status == "healthy" {
				hc.Status = "degraded"
			}
		} else {
			log.Debug("Check passed", map[string]interface{}{"check": name})
		}
	}
	
	log.Info("Health check completed", map[string]interface{}{
		"status": hc.Status,
		"checks": len(hc.Checks),
	})
}

func checkConfig(cfg *config.Config) CheckResult {
	if cfg == nil {
		return CheckResult{
			Status:  "error",
			Message: "Configuration is nil",
			Error:   "config_nil",
		}
	}
	
	issues := []string{}
	
	if cfg.App.Name == "" {
		issues = append(issues, "app.name is empty")
	}
	if cfg.Server.Port <= 0 {
		issues = append(issues, "server.port is invalid")
	}
	if cfg.Database.Type == "" {
		issues = append(issues, "database.type is empty")
	}
	
	if len(issues) > 0 {
		return CheckResult{
			Status:  "error",
			Message: "Configuration validation failed",
			Data:    issues,
		}
	}
	
	return CheckResult{
		Status:  "ok",
		Message: "Configuration is valid",
		Data: map[string]interface{}{
			"app_name":      cfg.App.Name,
			"server_port":   cfg.Server.Port,
			"database_type": cfg.Database.Type,
		},
	}
}

func checkDatabase(cfg *config.Config) CheckResult {
	switch cfg.Database.Type {
	case "sqlite":
		if cfg.Database.Path == "" {
			return CheckResult{
				Status:  "error",
				Message: "SQLite path is empty",
				Error:   "sqlite_path_missing",
			}
		}
		
		// Vérifier si le répertoire parent existe
		dir := cfg.Database.Path[:len(cfg.Database.Path)-len("/firesalamander.db")]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return CheckResult{
				Status:  "warning",
				Message: "SQLite directory doesn't exist yet",
				Data:    map[string]interface{}{"path": dir},
			}
		}
		
		return CheckResult{
			Status:  "ok",
			Message: "SQLite configuration is valid",
			Data:    map[string]interface{}{"path": cfg.Database.Path},
		}
		
	case "mysql":
		if cfg.Database.Host == "" || cfg.Database.Name == "" {
			return CheckResult{
				Status:  "error",
				Message: "MySQL configuration incomplete",
				Error:   "mysql_config_incomplete",
			}
		}
		
		return CheckResult{
			Status:  "ok",
			Message: "MySQL configuration is valid",
			Data: map[string]interface{}{
				"host": cfg.Database.Host,
				"name": cfg.Database.Name,
			},
		}
		
	default:
		return CheckResult{
			Status:  "error",
			Message: "Unknown database type",
			Error:   "unknown_db_type",
			Data:    map[string]interface{}{"type": cfg.Database.Type},
		}
	}
}

func checkFilesystem() CheckResult {
	dirs := []string{"config", "deploy"}
	files := []string{"go.mod", "main.go"}
	
	missing := []string{}
	
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			missing = append(missing, "directory: "+dir)
		}
	}
	
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missing = append(missing, "file: "+file)
		}
	}
	
	if len(missing) > 0 {
		return CheckResult{
			Status:  "error",
			Message: "Required files/directories missing",
			Data:    missing,
		}
	}
	
	return CheckResult{
		Status:  "ok",
		Message: "All required files and directories present",
	}
}

func checkNetwork(cfg *config.Config) CheckResult {
	// Tester si le port est libre (basique)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	
	return CheckResult{
		Status:  "ok",
		Message: "Network configuration valid",
		Data: map[string]interface{}{
			"port": cfg.Server.Port,
			"addr": addr,
		},
	}
}

func checkAI(cfg *config.Config) CheckResult {
	if !cfg.AI.Enabled {
		return CheckResult{
			Status:  "ok",
			Message: "AI is disabled",
			Data:    map[string]interface{}{"enabled": false},
		}
	}
	
	if cfg.AI.MockMode {
		return CheckResult{
			Status:  "ok",
			Message: "AI is in mock mode",
			Data:    map[string]interface{}{"mock_mode": true},
		}
	}
	
	if cfg.AI.APIKey == "" {
		return CheckResult{
			Status:  "warning",
			Message: "AI is enabled but API key is missing",
			Error:   "api_key_missing",
		}
	}
	
	return CheckResult{
		Status:  "ok",
		Message: "AI configuration is valid",
		Data: map[string]interface{}{
			"enabled":   true,
			"mock_mode": false,
			"api_key":   "***" + cfg.AI.APIKey[len(cfg.AI.APIKey)-4:],
		},
	}
}

func getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%.1f MB", float64(m.Alloc)/1024/1024)
}

func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

// Handler HTTP pour l'endpoint /debug
func DebugHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Debug endpoint called", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		})
		
		// Check if phase tests are requested
		if strings.Contains(r.URL.Query().Get("test"), "phase1") {
			log.Debug("Running Phase 1 tests")
			phaseTests := RunPhase1Tests(cfg)
			
			w.Header().Set("Content-Type", "application/json")
			if phaseTests.Status == "passed" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusExpectationFailed)
			}
			
			if err := json.NewEncoder(w).Encode(phaseTests); err != nil {
				log.Error("Failed to encode phase tests response", map[string]interface{}{
					"error": err.Error(),
				})
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		
		// Standard health check
		checker := NewChecker(cfg)
		
		w.Header().Set("Content-Type", "application/json")
		
		// Status code basé sur l'état
		switch checker.Status {
		case "healthy":
			w.WriteHeader(http.StatusOK)
		case "degraded":
			w.WriteHeader(http.StatusPartialContent)
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		
		if err := json.NewEncoder(w).Encode(checker); err != nil {
			log.Error("Failed to encode debug response", map[string]interface{}{
				"error": err.Error(),
			})
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}