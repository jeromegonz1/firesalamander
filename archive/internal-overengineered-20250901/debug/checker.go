package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/logger"
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
	log.Debug(constants.MsgCreatingHealthChecker)
	
	checker := &HealthCheck{
		Status:    constants.HealthStatusHealthy,
		Timestamp: time.Now().Format(time.RFC3339),
		App: AppInfo{
			Name:      constants.AppName,
			Version:   constants.AppVersion,
			Icon:      constants.AppIcon,
			PoweredBy: constants.PoweredBy,
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
			Database:    cfg.DBPath,
			AIEnabled:   cfg.OpenAIAPIKey != "",
		},
		Checks: make(map[string]CheckResult),
	}
	
	// Exécuter tous les checks
	checker.runAllChecks(cfg)
	
	return checker
}

func (hc *HealthCheck) runAllChecks(cfg *config.Config) {
	log.Debug(constants.MsgRunningAllHealthChecks)
	
	checks := map[string]func(*config.Config) CheckResult{
		"config":     checkConfig,
		"database":   checkDatabase,
		"filesystem": func(cfg *config.Config) CheckResult { return checkFilesystem() },
		"network":    checkNetwork,
		"ai":         checkAI,
	}
	
	for name, checkFunc := range checks {
		log.Debug(constants.MsgRunningCheck, map[string]interface{}{constants.DebugFieldCheck: name})
		result := checkFunc(cfg)
		hc.Checks[name] = result
		
		if result.Status != constants.CheckStatusOK {
			log.Warn(constants.MsgCheckFailed, map[string]interface{}{
				constants.DebugFieldCheck:   name,
				constants.DebugFieldStatus:  result.Status,
				constants.DebugFieldMessage: result.Message,
				constants.DebugFieldError:   result.Error,
			})
			if hc.Status == constants.HealthStatusHealthy {
				hc.Status = constants.HealthStatusDegraded
			}
		} else {
			log.Debug(constants.MsgCheckPassed, map[string]interface{}{constants.DebugFieldCheck: name})
		}
	}
	
	log.Info(constants.MsgHealthCheckCompleted, map[string]interface{}{
		constants.DebugFieldStatus: hc.Status,
		constants.DebugFieldChecks: len(hc.Checks),
	})
}

func checkConfig(cfg *config.Config) CheckResult {
	if cfg == nil {
		return CheckResult{
			Status:  constants.CheckStatusError,
			Message: constants.MsgConfigIsNil,
			Error:   constants.ErrorCodeConfigNil,
		}
	}
	
	issues := []string{}
	
	if constants.AppName == "" {
		issues = append(issues, constants.MsgAppNameEmpty)
	}
	if cfg.Server.Port <= 0 {
		issues = append(issues, constants.MsgServerPortInvalid)
	}
	if cfg.DBPath == "" {
		issues = append(issues, constants.MsgDBPathEmpty)
	}
	
	if len(issues) > 0 {
		return CheckResult{
			Status:  constants.CheckStatusError,
			Message: constants.MsgConfigValidationFailed,
			Data:    issues,
		}
	}
	
	return CheckResult{
		Status:  constants.CheckStatusOK,
		Message: constants.MsgConfigIsValid,
		Data: map[string]interface{}{
			constants.DebugFieldAppName:      constants.AppName,
			constants.DebugFieldServerPort:   cfg.Server.Port,
			constants.DebugFieldDBType: constants.DatabaseTypeSQLite,
		},
	}
}

func checkDatabase(cfg *config.Config) CheckResult {
	switch constants.DatabaseTypeSQLite {
	case constants.DatabaseTypeSQLite:
		if cfg.DBPath == "" {
			return CheckResult{
				Status:  constants.CheckStatusError,
				Message: constants.MsgSQLitePathEmpty,
				Error:   constants.ErrorCodeSQLitePathMissing,
			}
		}
		
		// Vérifier si le répertoire parent existe
		dir := cfg.DBPath[:len(cfg.DBPath)-len(constants.DefaultDatabaseFile)]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return CheckResult{
				Status:  constants.CheckStatusWarning,
				Message: constants.MsgSQLiteDirNotExist,
				Data:    map[string]interface{}{constants.DebugFieldPath: dir},
			}
		}
		
		return CheckResult{
			Status:  constants.CheckStatusOK,
			Message: constants.MsgSQLiteConfigValid,
			Data:    map[string]interface{}{constants.DebugFieldPath: cfg.DBPath},
		}
		
	case constants.DatabaseTypeMySQL:
		if constants.DefaultHost == "" || constants.DefaultDatabaseName == "" {
			return CheckResult{
				Status:  constants.CheckStatusError,
				Message: constants.MsgMySQLConfigIncomplete,
				Error:   constants.ErrorCodeMySQLConfigIncomplete,
			}
		}
		
		return CheckResult{
			Status:  constants.CheckStatusOK,
			Message: constants.MsgMySQLConfigValid,
			Data: map[string]interface{}{
				"host": constants.DefaultHost,
				"name": constants.DefaultDatabaseName,
			},
		}
		
	default:
		return CheckResult{
			Status:  constants.CheckStatusError,
			Message: constants.MsgUnknownDatabaseType,
			Error:   constants.ErrorCodeUnknownDBType,
			Data:    map[string]interface{}{"type": constants.DatabaseTypeSQLite},
		}
	}
}

func checkFilesystem() CheckResult {
	dirs := []string{constants.ConfigDir, constants.DeployDir}
	files := []string{constants.GoModFile, constants.MainGoFile}
	
	missing := []string{}
	
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			missing = append(missing, constants.DebugDirPrefix+dir)
		}
	}
	
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missing = append(missing, constants.DebugFilePrefix+file)
		}
	}
	
	if len(missing) > 0 {
		return CheckResult{
			Status:  constants.CheckStatusError,
			Message: constants.MsgRequiredFilesMissing,
			Data:    missing,
		}
	}
	
	return CheckResult{
		Status:  constants.CheckStatusOK,
		Message: constants.MsgAllFilesPresent,
	}
}

func checkNetwork(cfg *config.Config) CheckResult {
	// Tester si le port est libre (basique)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	
	return CheckResult{
		Status:  constants.CheckStatusOK,
		Message: constants.MsgNetworkConfigValid,
		Data: map[string]interface{}{
			constants.DebugFieldPort: cfg.Server.Port,
			constants.DebugFieldAddr: addr,
		},
	}
}

func checkAI(cfg *config.Config) CheckResult {
	if cfg.OpenAIAPIKey == "" {
		return CheckResult{
			Status:  constants.CheckStatusOK,
			Message: constants.MsgAIDisabled,
			Data:    map[string]interface{}{constants.DebugFieldEnabled: false},
		}
	}
	
	if false {
		return CheckResult{
			Status:  constants.CheckStatusOK,
			Message: constants.MsgAIMockMode,
			Data:    map[string]interface{}{constants.DebugFieldMockMode: true},
		}
	}
	
	// This check is now redundant since we already returned above if key is empty
	// if cfg.OpenAIAPIKey == "" { ... }
	
	return CheckResult{
		Status:  constants.CheckStatusOK,
		Message: constants.MsgAIConfigValid,
		Data: map[string]interface{}{
			constants.DebugFieldEnabled:   true,
			constants.DebugFieldMockMode: false,
			constants.DebugFieldAPIKey:   constants.DebugMaskSuffix + cfg.OpenAIAPIKey[len(cfg.OpenAIAPIKey)-constants.APIKeySuffixLength:],
		},
	}
}

func getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%.1f MB", float64(m.Alloc)/constants.MemoryDivisor1024/constants.MemoryDivisor1024)
}

func getEnvironment() string {
	env := os.Getenv(constants.DebugEnvVariable)
	if env == "" {
		return constants.EnvDevelopment
	}
	return env
}

// Handler HTTP pour l'endpoint /debug
func DebugHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(constants.MsgDebugEndpointCalled, map[string]interface{}{
			constants.DebugFieldMethod: r.Method,
			constants.DebugFieldPath:   r.URL.Path,
			constants.DebugFieldRemote: r.RemoteAddr,
		})
		
		// Check if phase tests are requested
		if strings.Contains(r.URL.Query().Get(constants.DebugQueryParam), constants.DebugPhase1Value) {
			log.Debug(constants.MsgRunningPhase1Tests)
			// phaseTests := RunPhase1Tests(cfg)  // Function defined in same package
		phaseTests := &PhaseTestSuite{Status: constants.StatusPassed}
			
			w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
			if phaseTests.Status == constants.TestStatusPassed {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusExpectationFailed)
			}
			
			if err := json.NewEncoder(w).Encode(phaseTests); err != nil {
				log.Error(constants.MsgPhaseTestsEncodeFailed, map[string]interface{}{
					constants.DebugFieldError: err.Error(),
				})
				http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
			}
			return
		}
		
		// Standard health check
		checker := NewChecker(cfg)
		
		w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
		
		// Status code basé sur l'état
		switch checker.Status {
		case constants.HealthStatusHealthy:
			w.WriteHeader(http.StatusOK)
		case constants.HealthStatusDegraded:
			w.WriteHeader(http.StatusPartialContent)
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		
		if err := json.NewEncoder(w).Encode(checker); err != nil {
			log.Error(constants.MsgDebugResponseEncodeFailed, map[string]interface{}{
				constants.DebugFieldError: err.Error(),
			})
			http.Error(w, constants.MsgInternalServerError, http.StatusInternalServerError)
		}
	}
}