package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	debugMode = false
	logLevels = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
	colors = map[LogLevel]string{
		DEBUG: "\033[36m", // Cyan
		INFO:  "\033[32m", // Green
		WARN:  "\033[33m", // Yellow
		ERROR: "\033[31m", // Red
		FATAL: "\033[35m", // Magenta
	}
	resetColor = "\033[0m"
)

type Logger struct {
	component string
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	File      string                 `json:"file,omitempty"`
	Line      int                    `json:"line,omitempty"`
}

func init() {
	// Active le debug si ENV=development ou DEBUG=true
	env := os.Getenv("ENV")
	debug := os.Getenv("DEBUG")
	debugMode = env == "development" || debug == "true"
	
	if debugMode {
		log.SetFlags(0) // On gère nous-même le formatage
		fmt.Println("🐛 DEBUG MODE ACTIVATED - Fire Salamander Logger")
	}
}

func New(component string) *Logger {
	return &Logger{component: component}
}

func SetDebugMode(enabled bool) {
	debugMode = enabled
	if enabled {
		fmt.Println("🐛 DEBUG MODE ENABLED")
	}
}

func IsDebugMode() bool {
	return debugMode
}

func (l *Logger) log(level LogLevel, message string, data map[string]interface{}) {
	if level == DEBUG && !debugMode {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Level:     logLevels[level],
		Component: l.component,
		Message:   message,
		Data:      data,
	}

	// Ajouter info sur le fichier en mode debug
	if debugMode {
		if pc, file, line, ok := runtime.Caller(2); ok {
			// Extraire juste le nom du fichier
			parts := strings.Split(file, "/")
			entry.File = parts[len(parts)-1]
			entry.Line = line
			
			// Extraire le nom de la fonction
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				funcName := fn.Name()
				if idx := strings.LastIndex(funcName, "."); idx >= 0 {
					funcName = funcName[idx+1:]
				}
				if entry.Data == nil {
					entry.Data = make(map[string]interface{})
				}
				entry.Data["function"] = funcName
			}
		}
	}

	// Format console coloré
	l.logConsole(level, entry)

	// Format JSON pour les logs structurés
	if level >= ERROR {
		l.logJSON(entry)
	}
}

func (l *Logger) logConsole(level LogLevel, entry LogEntry) {
	color := colors[level]
	icon := l.getIcon(level)
	
	// Format base
	output := fmt.Sprintf("%s%s [%s] %s%s %s",
		color, icon, entry.Level, entry.Component, resetColor, entry.Message)
	
	// Ajouter données si présentes
	if len(entry.Data) > 0 {
		dataStr := l.formatData(entry.Data)
		output += fmt.Sprintf(" %s", dataStr)
	}
	
	// Ajouter info debug
	if debugMode && entry.File != "" {
		output += fmt.Sprintf(" \033[90m(%s:%d)\033[0m", entry.File, entry.Line)
	}
	
	fmt.Println(output)
}

func (l *Logger) logJSON(entry LogEntry) {
	jsonData, _ := json.Marshal(entry)
	log.Printf("JSON_LOG: %s", string(jsonData))
}

func (l *Logger) getIcon(level LogLevel) string {
	switch level {
	case DEBUG:
		return "🐛"
	case INFO:
		return "ℹ️ "
	case WARN:
		return "⚠️ "
	case ERROR:
		return "❌"
	case FATAL:
		return "💀"
	default:
		return "📝"
	}
}

func (l *Logger) formatData(data map[string]interface{}) string {
	var parts []string
	for k, v := range data {
		if k == "function" {
			continue // Déjà affiché avec le fichier
		}
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	if len(parts) > 0 {
		return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
	}
	return ""
}

// Méthodes de logging
func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.log(DEBUG, message, d)
}

func (l *Logger) Info(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.log(INFO, message, d)
}

func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.log(WARN, message, d)
}

func (l *Logger) Error(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.log(ERROR, message, d)
}

func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.log(FATAL, message, d)
	os.Exit(1)
}

// Logger global pour usage simple
var globalLogger = New("FIRE-SALAMANDER")

func Debug(message string, data ...map[string]interface{}) {
	globalLogger.Debug(message, data...)
}

func Info(message string, data ...map[string]interface{}) {
	globalLogger.Info(message, data...)
}

func Warn(message string, data ...map[string]interface{}) {
	globalLogger.Warn(message, data...)
}

func Error(message string, data ...map[string]interface{}) {
	globalLogger.Error(message, data...)
}

func Fatal(message string, data ...map[string]interface{}) {
	globalLogger.Fatal(message, data...)
}