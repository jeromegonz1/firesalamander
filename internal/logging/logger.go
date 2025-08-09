package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - LOGGING SYSTEM
// TDD + Zero Hardcoding Policy
// ========================================

// LogLevel représente le niveau de log
type LogLevel int

const (
	TRACE LogLevel = iota - 1
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// String retourne la représentation string du niveau
func (l LogLevel) String() string {
	switch l {
	case TRACE:
		return constants.LogLevelTrace
	case DEBUG:
		return constants.LogLevelDebug
	case INFO:
		return constants.LogLevelInfo
	case WARN:
		return constants.LogLevelWarn
	case ERROR:
		return constants.LogLevelError
	case FATAL:
		return constants.LogLevelFatal
	default:
		return constants.LogLevelInfo
	}
}

// LogEntry représente une entrée de log
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Category   string                 `json:"category"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	File       string                 `json:"file,omitempty"`
	Function   string                 `json:"function,omitempty"`
	Line       int                    `json:"line,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

// LogConfig configuration du système de logging
type LogConfig struct {
	Level           LogLevel
	Format          string // "json" ou "text"
	EnableConsole   bool
	EnableFile      bool
	LogDir          string
	RotationSizeMB  int
	RotationBackups int
	RotationAgeDays int
}

// Logger interface principale
type Logger interface {
	Trace(category, message string, data ...map[string]interface{})
	Debug(category, message string, data ...map[string]interface{})
	Info(category, message string, data ...map[string]interface{})
	Warn(category, message string, data ...map[string]interface{})
	Error(category, message string, err error, data ...map[string]interface{})
	Fatal(category, message string, err error, data ...map[string]interface{})
	
	// Context methods
	WithRequestID(requestID string) Logger
	WithTraceID(traceID string) Logger
	WithData(data map[string]interface{}) Logger
	
	// Specialized loggers
	HTTP() HTTPLogger
	API() APILogger
	Performance() PerformanceLogger
	Audit() AuditLogger
	
	// Management
	SetLevel(level LogLevel)
	Close() error
}

// HTTPLogger interface pour logs HTTP
type HTTPLogger interface {
	RequestReceived(method, url, remoteAddr, userAgent string, requestID string)
	ResponseSent(method, url string, statusCode int, responseTimeMs int64, contentLength int64, requestID string)
	RequestError(method, url string, err error, requestID string)
}

// APILogger interface pour logs API
type APILogger interface {
	RequestReceived(endpoint, method string, requestID string, data map[string]interface{})
	ResponseSent(endpoint, method string, statusCode int, responseTimeMs int64, requestID string)
	ValidationError(endpoint, method string, err error, requestID string)
	ProcessingError(endpoint, method string, err error, requestID string)
}

// PerformanceLogger interface pour logs de performance
type PerformanceLogger interface {
	OperationStarted(operation string, data map[string]interface{})
	OperationCompleted(operation string, durationMs int64, data map[string]interface{})
	MemoryUsage(operation string, memoryBytes int64, goroutines int)
	RequestMetrics(requestsPerSec float64, avgResponseTimeMs int64)
}

// AuditLogger interface pour logs d'audit
type AuditLogger interface {
	ActionStarted(action, resource, userID, sessionID string, beforeState interface{})
	ActionCompleted(action, resource, userID, sessionID string, success bool, afterState interface{})
	SecurityEvent(event, description string, data map[string]interface{})
}

// FireSalamanderLogger implémentation principale
type FireSalamanderLogger struct {
	config    *LogConfig
	writers   map[string]io.Writer
	mu        sync.RWMutex
	requestID string
	traceID   string
	contextData map[string]interface{}
}

// NewLogger crée un nouveau logger Fire Salamander
func NewLogger(config *LogConfig) (Logger, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	// Créer le répertoire de logs si nécessaire
	if config.EnableFile && config.LogDir != "" {
		if err := os.MkdirAll(config.LogDir, constants.LogDirPermissions); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}
	
	logger := &FireSalamanderLogger{
		config:  config,
		writers: make(map[string]io.Writer),
		contextData: make(map[string]interface{}),
	}
	
	// Initialiser les writers
	if err := logger.initializeWriters(); err != nil {
		return nil, fmt.Errorf("failed to initialize writers: %w", err)
	}
	
	return logger, nil
}

// initializeWriters initialise les writers pour chaque type de log
func (l *FireSalamanderLogger) initializeWriters() error {
	// Console writer
	if l.config.EnableConsole {
		l.writers["console"] = os.Stdout
	}
	
	// File writers
	if l.config.EnableFile {
		logFiles := []string{
			constants.LogFileAccess,
			constants.LogFileError,
			constants.LogFileDebug,
			constants.LogFileAudit,
			constants.LogFilePerformance,
			constants.LogFileSystem,
		}
		
		for _, fileName := range logFiles {
			filePath := filepath.Join(l.config.LogDir, fileName)
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.LogFilePermissions)
			if err != nil {
				return fmt.Errorf("failed to open log file %s: %w", fileName, err)
			}
			l.writers[strings.TrimSuffix(fileName, ".log")] = file
		}
	}
	
	return nil
}

// writeLog écrit une entrée de log
func (l *FireSalamanderLogger) writeLog(level LogLevel, category, message string, err error, data map[string]interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	// Vérifier le niveau
	if level < l.config.Level {
		return
	}
	
	// Créer l'entrée de log
	entry := l.createLogEntry(level, category, message, err, data)
	
	// Formater l'entrée
	var output string
	if l.config.Format == constants.DefaultLogFormat {
		output = l.formatJSON(entry)
	} else {
		output = l.formatText(entry)
	}
	
	// Écrire vers les destinations appropriées
	l.writeToDestinations(category, level, output)
}

// createLogEntry crée une entrée de log
func (l *FireSalamanderLogger) createLogEntry(level LogLevel, category, message string, err error, data map[string]interface{}) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now().Format(constants.LogTimestampFormat),
		Level:     level.String(),
		Category:  category,
		Message:   message,
		Data:      make(map[string]interface{}),
		TraceID:   l.traceID,
		RequestID: l.requestID,
	}
	
	// Ajouter les données contextuelles
	for k, v := range l.contextData {
		entry.Data[k] = v
	}
	
	// Ajouter les données passées en paramètre
	for k, v := range data {
		entry.Data[k] = v
	}
	
	// Ajouter l'erreur si présente
	if err != nil {
		entry.Error = err.Error()
	}
	
	// Informations de caller pour DEBUG et ERROR
	if level <= DEBUG || level >= ERROR {
		if pc, file, line, ok := runtime.Caller(3); ok {
			entry.File = filepath.Base(file)
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}
	
	return entry
}

// formatJSON formate en JSON
func (l *FireSalamanderLogger) formatJSON(entry *LogEntry) string {
	data, _ := json.Marshal(entry.Data)
	return fmt.Sprintf(constants.LogFormatJSON,
		entry.Timestamp,
		entry.Level,
		entry.Category,
		entry.Message,
		string(data),
		entry.TraceID,
		entry.RequestID,
	) + "\n"
}

// formatText formate en texte
func (l *FireSalamanderLogger) formatText(entry *LogEntry) string {
	return fmt.Sprintf(constants.LogFormatText,
		entry.Timestamp,
		entry.Level,
		entry.Category,
		entry.Message,
		entry.RequestID,
	) + "\n"
}

// writeToDestinations écrit vers les bonnes destinations
func (l *FireSalamanderLogger) writeToDestinations(category string, level LogLevel, output string) {
	// Console pour tous
	if writer, ok := l.writers["console"]; ok {
		writer.Write([]byte(output))
	}
	
	// Fichiers spécialisés
	if l.config.EnableFile {
		switch {
		case strings.Contains(category, constants.LogCategoryHTTP):
			if writer, ok := l.writers["access"]; ok {
				writer.Write([]byte(output))
			}
		case level >= ERROR:
			if writer, ok := l.writers["error"]; ok {
				writer.Write([]byte(output))
			}
		case level == DEBUG || level == TRACE:
			if writer, ok := l.writers["debug"]; ok {
				writer.Write([]byte(output))
			}
		case strings.Contains(category, constants.LogCategoryAudit):
			if writer, ok := l.writers["audit"]; ok {
				writer.Write([]byte(output))
			}
		case strings.Contains(category, constants.LogCategoryPerformance):
			if writer, ok := l.writers["performance"]; ok {
				writer.Write([]byte(output))
			}
		default:
			if writer, ok := l.writers["system"]; ok {
				writer.Write([]byte(output))
			}
		}
	}
}

// Implémentation des méthodes Logger
func (l *FireSalamanderLogger) Trace(category, message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(TRACE, category, message, nil, dataMap)
}

func (l *FireSalamanderLogger) Debug(category, message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(DEBUG, category, message, nil, dataMap)
}

func (l *FireSalamanderLogger) Info(category, message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(INFO, category, message, nil, dataMap)
}

func (l *FireSalamanderLogger) Warn(category, message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(WARN, category, message, nil, dataMap)
}

func (l *FireSalamanderLogger) Error(category, message string, err error, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(ERROR, category, message, err, dataMap)
}

func (l *FireSalamanderLogger) Fatal(category, message string, err error, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.writeLog(FATAL, category, message, err, dataMap)
	os.Exit(1)
}

// Méthodes de contexte
func (l *FireSalamanderLogger) WithRequestID(requestID string) Logger {
	newLogger := *l
	newLogger.requestID = requestID
	return &newLogger
}

func (l *FireSalamanderLogger) WithTraceID(traceID string) Logger {
	newLogger := *l
	newLogger.traceID = traceID
	return &newLogger
}

func (l *FireSalamanderLogger) WithData(data map[string]interface{}) Logger {
	newLogger := *l
	newLogger.contextData = make(map[string]interface{})
	for k, v := range l.contextData {
		newLogger.contextData[k] = v
	}
	for k, v := range data {
		newLogger.contextData[k] = v
	}
	return &newLogger
}

// Méthodes de gestion
func (l *FireSalamanderLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

func (l *FireSalamanderLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	for name, writer := range l.writers {
		if name != "console" {
			if closer, ok := writer.(io.Closer); ok {
				closer.Close()
			}
		}
	}
	return nil
}

// Loggers spécialisés - implémentés dans des fichiers séparés
func (l *FireSalamanderLogger) HTTP() HTTPLogger {
	return NewHTTPLogger(l)
}

func (l *FireSalamanderLogger) API() APILogger {
	return NewAPILogger(l)
}

func (l *FireSalamanderLogger) Performance() PerformanceLogger {
	return NewPerformanceLogger(l)
}

func (l *FireSalamanderLogger) Audit() AuditLogger {
	return NewAuditLogger(l)
}