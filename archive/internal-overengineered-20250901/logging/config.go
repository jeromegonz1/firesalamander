package logging

import (
	"os"
	"strconv"
	"strings"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - LOGGING CONFIG
// TDD + Zero Hardcoding Policy
// ========================================

// LoadConfigFromEnv charge la configuration de logging depuis les variables d'environnement
func LoadConfigFromEnv() *LogConfig {
	config := &LogConfig{
		Level:           parseLogLevel(getEnvOrDefault(constants.EnvLogLevel, constants.DefaultLogLevel)),
		Format:          getEnvOrDefault(constants.EnvLogFormat, constants.DefaultLogFormat),
		EnableConsole:   getEnvBoolOrDefault("ENABLE_CONSOLE_LOGGING", constants.DefaultEnableConsole),
		EnableFile:      getEnvBoolOrDefault("ENABLE_FILE_LOGGING", constants.DefaultEnableFile),
		LogDir:          getEnvOrDefault(constants.EnvLogDir, constants.LogDefaultDir),
		RotationSizeMB:  getEnvIntOrDefault(constants.EnvLogRotationSize, constants.LogRotationMaxSizeMB),
		RotationBackups: getEnvIntOrDefault(constants.EnvLogRotationBackups, constants.LogRotationMaxBackups),
		RotationAgeDays: getEnvIntOrDefault(constants.EnvLogRotationAge, constants.LogRotationMaxAgeDays),
	}
	
	return config
}

// parseLogLevel parse le niveau de log depuis une string
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case constants.LogLevelTrace:
		return TRACE
	case constants.LogLevelDebug:
		return DEBUG
	case constants.LogLevelInfo:
		return INFO
	case constants.LogLevelWarn:
		return WARN
	case constants.LogLevelError:
		return ERROR
	case constants.LogLevelFatal:
		return FATAL
	default:
		return INFO // Défaut sécurisé
	}
}

// getEnvOrDefault récupère une variable d'environnement ou retourne la valeur par défaut
func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolOrDefault récupère une variable d'environnement booléenne
func getEnvBoolOrDefault(envVar string, defaultValue bool) bool {
	if value := os.Getenv(envVar); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvIntOrDefault récupère une variable d'environnement entière
func getEnvIntOrDefault(envVar string, defaultValue int) int {
	if value := os.Getenv(envVar); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}

// DefaultConfig retourne une configuration par défaut pour le développement
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:           DEBUG,
		Format:          "json",
		EnableConsole:   true,
		EnableFile:      true,
		LogDir:          constants.LogDefaultDir,
		RotationSizeMB:  constants.LogRotationMaxSizeMB,
		RotationBackups: constants.LogRotationMaxBackups,
		RotationAgeDays: constants.LogRotationMaxAgeDays,
	}
}

// ProductionConfig retourne une configuration pour la production
func ProductionConfig() *LogConfig {
	return &LogConfig{
		Level:           INFO,
		Format:          "json",
		EnableConsole:   false,
		EnableFile:      true,
		LogDir:          "/var/log/fire-salamander",
		RotationSizeMB:  constants.LogRotationMaxSizeMB,
		RotationBackups: constants.LogRotationMaxBackups,
		RotationAgeDays: constants.LogRotationMaxAgeDays,
	}
}