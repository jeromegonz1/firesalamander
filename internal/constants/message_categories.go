// Package constants provides centralized message management with category-based organization
// following clean architecture principles for separation of concerns and future extensibility.
//
// Architecture Design:
// - Separation of concerns: Messages grouped by functional category and audience
// - Type safety: Strong typing with custom types for message categories
// - Maintainability: Single source of truth with consistent naming conventions
// - Extensibility: Template-ready structure for dynamic content and i18n
// - Future-proof: Ready for internationalization with key-based message resolution
package constants

import (
	"fmt"
)

// MessageCategory represents the functional category of a message for architectural separation
type MessageCategory string

// MessageAudience defines the target audience for proper message routing
type MessageAudience string

// MessageSeverity indicates the importance level for consistent UX
type MessageSeverity string

// Message categories for architectural separation of concerns
const (
	CategoryError   MessageCategory = "error"   // Error conditions and failures
	CategorySuccess MessageCategory = "success" // Successful operations and completion
	CategoryInfo    MessageCategory = "info"    // Informational messages
	CategoryWarning MessageCategory = "warning" // Warning conditions
	CategoryUI      MessageCategory = "ui"      // User interface elements
	CategoryLog     MessageCategory = "log"     // System logging messages
	CategoryHelp    MessageCategory = "help"    // Help and documentation
	CategoryTime    MessageCategory = "time"    // Time estimates and progress
	CategoryPhase   MessageCategory = "phase"   // Analysis phase indicators
	CategoryMode    MessageCategory = "mode"    // Operational mode descriptions
	CategorySEO     MessageCategory = "seo"     // SEO analysis results
	CategoryAI      MessageCategory = "ai"      // AI-driven recommendations
)

// Message audiences for proper routing and formatting
const (
	AudienceUser     MessageAudience = "user"     // End users via UI
	AudienceInternal MessageAudience = "internal" // System/developer logs
	AudienceAPI      MessageAudience = "api"      // API consumers
)

// Message severity levels for consistent UX (using Message prefix to avoid conflicts)
const (
	MessageSeverityLow    MessageSeverity = "low"    // Minor issues or info
	MessageSeverityMedium MessageSeverity = "medium" // Important notifications
	MessageSeverityHigh   MessageSeverity = "high"   // Critical issues requiring attention
)

// MessageMetadata contains architectural metadata for proper message handling
type MessageMetadata struct {
	Category     MessageCategory `json:"category"`
	Audience     MessageAudience `json:"audience"`
	Severity     MessageSeverity `json:"severity"`
	I18nReady    bool           `json:"i18n_ready"`    // Prepared for internationalization
	HasTemplate  bool           `json:"has_template"`  // Contains template placeholders
	TechnicalLevel string       `json:"technical_level"` // user-friendly, semi-technical, technical
}

// ArchitecturalMessage represents a message with its metadata for clean architecture
type ArchitecturalMessage struct {
	Key      string          `json:"key"`      // Unique identifier for i18n
	Content  string          `json:"content"`  // Message content or template
	Metadata MessageMetadata `json:"metadata"` // Architectural metadata
}

// Message registry following architectural patterns - organized by category for maintainability
var MessageRegistry = map[string]ArchitecturalMessage{
	
	// === ERROR MESSAGES === (CategoryError)
	"ERR_METHOD_NOT_ALLOWED": {
		Key:     "ERR_METHOD_NOT_ALLOWED",
		Content: "Method not allowed",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"ERR_INVALID_JSON": {
		Key:     "ERR_INVALID_JSON",
		Content: "Invalid JSON format",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "technical",
		},
	},
	"ERR_URL_REQUIRED": {
		Key:     "ERR_URL_REQUIRED",
		Content: "URL is required for analysis",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"ERR_INVALID_URL": {
		Key:     "ERR_INVALID_URL",
		Content: "Invalid URL format provided",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"ERR_ANALYSIS_NOT_FOUND": {
		Key:     "ERR_ANALYSIS_NOT_FOUND",
		Content: "Analysis not found",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"ERR_ANALYSIS_INCOMPLETE": {
		Key:     "ERR_ANALYSIS_INCOMPLETE",
		Content: "Analysis is not yet complete",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityHigh,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"ERR_CONNECTION_FAILED": {
		Key:     "ERR_CONNECTION_FAILED",
		Content: "Connection failed - please check your network",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityHigh,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"ERR_REQUEST_TIMEOUT": {
		Key:     "ERR_REQUEST_TIMEOUT",
		Content: "Request timeout - please try again",
		Metadata: MessageMetadata{
			Category:       CategoryError,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},

	// === SUCCESS MESSAGES === (CategorySuccess)
	"SUCCESS_SERVER_STARTING": {
		Key:     "SUCCESS_SERVER_STARTING",
		Content: "Fire Salamander is starting up",
		Metadata: MessageMetadata{
			Category:       CategorySuccess,
			Audience:       AudienceInternal,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"SUCCESS_SERVER_STARTED": {
		Key:     "SUCCESS_SERVER_STARTED",
		Content: "Fire Salamander started successfully on {port}",
		Metadata: MessageMetadata{
			Category:       CategorySuccess,
			Audience:       AudienceInternal,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    true,
			TechnicalLevel: "user-friendly",
		},
	},
	"SUCCESS_ANALYSIS_STARTED": {
		Key:     "SUCCESS_ANALYSIS_STARTED",
		Content: "Analysis has been initiated successfully",
		Metadata: MessageMetadata{
			Category:       CategorySuccess,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"SUCCESS_ANALYSIS_COMPLETE": {
		Key:     "SUCCESS_ANALYSIS_COMPLETE",
		Content: "Analysis completed successfully",
		Metadata: MessageMetadata{
			Category:       CategorySuccess,
			Audience:       AudienceUser,
			Severity:       MessageSeverityHigh,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},

	// === INFO MESSAGES === (CategoryInfo)
	"INFO_SERVER_STOPPING": {
		Key:     "INFO_SERVER_STOPPING",
		Content: "Fire Salamander is shutting down gracefully",
		Metadata: MessageMetadata{
			Category:       CategoryInfo,
			Audience:       AudienceInternal,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"INFO_SERVER_STOPPED": {
		Key:     "INFO_SERVER_STOPPED",
		Content: "Fire Salamander has been stopped",
		Metadata: MessageMetadata{
			Category:       CategoryInfo,
			Audience:       AudienceInternal,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},

	// === UI MESSAGES === (CategoryUI)
	"UI_NEW_ANALYSIS": {
		Key:     "UI_NEW_ANALYSIS",
		Content: "New Analysis",
		Metadata: MessageMetadata{
			Category:       CategoryUI,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"UI_EXPORT_PDF": {
		Key:     "UI_EXPORT_PDF",
		Content: "Export PDF Report",
		Metadata: MessageMetadata{
			Category:       CategoryUI,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"UI_CANCEL": {
		Key:     "UI_CANCEL",
		Content: "Cancel",
		Metadata: MessageMetadata{
			Category:       CategoryUI,
			Audience:       AudienceUser,
			Severity:       MessageSeverityLow,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"UI_ANALYZE_BUTTON": {
		Key:     "UI_ANALYZE_BUTTON",
		Content: "ANALYZE MY WEBSITE",
		Metadata: MessageMetadata{
			Category:       CategoryUI,
			Audience:       AudienceUser,
			Severity:       MessageSeverityLow,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"UI_ANALYZING": {
		Key:     "UI_ANALYZING",
		Content: "ANALYZING YOUR WEBSITE...",
		Metadata: MessageMetadata{
			Category:       CategoryUI,
			Audience:       AudienceUser,
			Severity:       MessageSeverityHigh,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},

	// === PHASE MESSAGES === (CategoryPhase)
	"PHASE_DISCOVERY": {
		Key:     "PHASE_DISCOVERY",
		Content: "Discovering and mapping your website pages...",
		Metadata: MessageMetadata{
			Category:       CategoryPhase,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"PHASE_SEO_ANALYSIS": {
		Key:     "PHASE_SEO_ANALYSIS",
		Content: "Performing comprehensive SEO analysis...",
		Metadata: MessageMetadata{
			Category:       CategoryPhase,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"PHASE_AI_ANALYSIS": {
		Key:     "PHASE_AI_ANALYSIS",
		Content: "Running AI-powered content analysis...",
		Metadata: MessageMetadata{
			Category:       CategoryPhase,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "semi-technical",
		},
	},
	"PHASE_REPORT_GENERATION": {
		Key:     "PHASE_REPORT_GENERATION",
		Content: "Generating comprehensive analysis report...",
		Metadata: MessageMetadata{
			Category:       CategoryPhase,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},

	// === TIME ESTIMATES === (CategoryTime)
	"TIME_30_SECONDS": {
		Key:     "TIME_30_SECONDS",
		Content: "30 seconds remaining",
		Metadata: MessageMetadata{
			Category:       CategoryTime,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"TIME_1_MINUTE": {
		Key:     "TIME_1_MINUTE",
		Content: "1 minute remaining",
		Metadata: MessageMetadata{
			Category:       CategoryTime,
			Audience:       AudienceUser,
			Severity:       MessageSeverityMedium,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"TIME_CALCULATING": {
		Key:     "TIME_CALCULATING",
		Content: "Calculating time estimate...",
		Metadata: MessageMetadata{
			Category:       CategoryTime,
			Audience:       AudienceUser,
			Severity:       MessageSeverityLow,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
	"TIME_COMPLETE": {
		Key:     "TIME_COMPLETE",
		Content: "Analysis Complete",
		Metadata: MessageMetadata{
			Category:       CategoryTime,
			Audience:       AudienceUser,
			Severity:       MessageSeverityHigh,
			I18nReady:      true,
			HasTemplate:    false,
			TechnicalLevel: "user-friendly",
		},
	},
}

// === ARCHITECTURAL MESSAGE MANAGER ===

// MessageManager provides centralized message management following clean architecture principles
type MessageManager struct {
	registry map[string]ArchitecturalMessage
	locale   string // Future i18n support
}

// NewMessageManager creates a new message manager instance
func NewMessageManager() *MessageManager {
	return &MessageManager{
		registry: MessageRegistry,
		locale:   "en", // Default locale
	}
}

// GetMessage retrieves a message by key with architectural metadata
func (mm *MessageManager) GetMessage(key string) (ArchitecturalMessage, error) {
	msg, exists := mm.registry[key]
	if !exists {
		return ArchitecturalMessage{}, fmt.Errorf("message key '%s' not found in registry", key)
	}
	return msg, nil
}

// GetMessageContent returns just the content string for backward compatibility
func (mm *MessageManager) GetMessageContent(key string) string {
	msg, err := mm.GetMessage(key)
	if err != nil {
		return fmt.Sprintf("[MISSING_MESSAGE: %s]", key)
	}
	return msg.Content
}

// GetMessagesByCategory returns all messages for a specific category
func (mm *MessageManager) GetMessagesByCategory(category MessageCategory) map[string]ArchitecturalMessage {
	result := make(map[string]ArchitecturalMessage)
	for key, msg := range mm.registry {
		if msg.Metadata.Category == category {
			result[key] = msg
		}
	}
	return result
}

// GetMessagesByAudience returns all messages for a specific audience
func (mm *MessageManager) GetMessagesByAudience(audience MessageAudience) map[string]ArchitecturalMessage {
	result := make(map[string]ArchitecturalMessage)
	for key, msg := range mm.registry {
		if msg.Metadata.Audience == audience {
			result[key] = msg
		}
	}
	return result
}

// FormatMessage applies template parameters to a message (future implementation)
func (mm *MessageManager) FormatMessage(key string, params map[string]string) (string, error) {
	msg, err := mm.GetMessage(key)
	if err != nil {
		return "", err
	}
	
	content := msg.Content
	if msg.Metadata.HasTemplate {
		// Template implementation would go here
		// For now, simple placeholder replacement
		for paramKey, paramValue := range params {
			placeholder := fmt.Sprintf("{%s}", paramKey)
			content = fmt.Sprintf(content, paramValue) // Basic implementation
			_ = placeholder // Avoid unused variable warning
		}
	}
	
	return content, nil
}

// ValidateMessageRegistry performs architectural validation of the message registry
func (mm *MessageManager) ValidateMessageRegistry() []string {
	var errors []string
	
	for key, msg := range mm.registry {
		// Validate key consistency
		if msg.Key != key {
			errors = append(errors, fmt.Sprintf("Key mismatch for %s: registry key differs from message key", key))
		}
		
		// Validate required fields
		if msg.Content == "" {
			errors = append(errors, fmt.Sprintf("Empty content for message key: %s", key))
		}
		
		// Validate metadata consistency
		if msg.Metadata.Category == "" {
			errors = append(errors, fmt.Sprintf("Missing category for message key: %s", key))
		}
		
		if msg.Metadata.Audience == "" {
			errors = append(errors, fmt.Sprintf("Missing audience for message key: %s", key))
		}
	}
	
	return errors
}

// === MIGRATION HELPERS FOR EXISTING CODE ===

// Legacy constants for backward compatibility during migration
// These will be deprecated once the eliminator script completes the refactoring

// GetLegacyMessage provides backward compatibility during architectural migration
func GetLegacyMessage(oldConstant string) string {
	// Map old constant names to new message keys
	legacyMap := map[string]string{
		"ErrMethodNotAllowed":    "ERR_METHOD_NOT_ALLOWED",
		"ErrInvalidJSON":        "ERR_INVALID_JSON", 
		"ErrURLRequired":        "ERR_URL_REQUIRED",
		"ErrInvalidURL":         "ERR_INVALID_URL",
		"ErrAnalysisNotFound":   "ERR_ANALYSIS_NOT_FOUND",
		"ErrAnalysisIncomplete": "ERR_ANALYSIS_INCOMPLETE",
		"ErrConnectionFailed":   "ERR_CONNECTION_FAILED",
		"ErrTimeout":            "ERR_REQUEST_TIMEOUT",
		"ServerStarting":        "SUCCESS_SERVER_STARTING",
		"ServerStarted":         "SUCCESS_SERVER_STARTED",
		"ServerStopping":        "INFO_SERVER_STOPPING",
		"ServerStopped":         "INFO_SERVER_STOPPED",
		"AnalysisStarted":       "SUCCESS_ANALYSIS_STARTED",
		"AnalysisComplete":      "SUCCESS_ANALYSIS_COMPLETE",
		"UINewAnalysis":         "UI_NEW_ANALYSIS",
		"UIExportPDF":           "UI_EXPORT_PDF",
		"UICancel":              "UI_CANCEL",
		"UIAnalyzeButton":       "UI_ANALYZE_BUTTON",
		"UIAnalyzing":           "UI_ANALYZING",
		"PhaseDiscoveryMsg":     "PHASE_DISCOVERY",
		"PhaseSEOAnalysisMsg":   "PHASE_SEO_ANALYSIS",
		"PhaseAIAnalysisMsg":    "PHASE_AI_ANALYSIS",
		"PhaseReportGenMsg":     "PHASE_REPORT_GENERATION",
		"TimeEstimate30s":       "TIME_30_SECONDS",
		"TimeEstimate1m":        "TIME_1_MINUTE",
		"TimeEstimateCalculating": "TIME_CALCULATING",
		"TimeEstimateComplete":  "TIME_COMPLETE",
	}
	
	if newKey, exists := legacyMap[oldConstant]; exists {
		mm := NewMessageManager()
		return mm.GetMessageContent(newKey)
	}
	
	return fmt.Sprintf("[UNMAPPED_LEGACY: %s]", oldConstant)
}

// === ARCHITECTURAL CONSTANTS FOR FUTURE EXTENSIBILITY ===

// I18n readiness markers for future internationalization
const (
	// Default locale configuration
	DefaultLocale = "en"
	
	// Template placeholder patterns
	TemplateStartDelim = "{"
	TemplateEndDelim   = "}"
	
	// Message validation patterns
	MaxMessageLength = 500
	MinMessageLength = 1
)

// Future extension points for clean architecture
type MessageExtension interface {
	FormatMessage(msg ArchitecturalMessage, params map[string]interface{}) string
	ValidateMessage(msg ArchitecturalMessage) error
	TranslateMessage(msg ArchitecturalMessage, locale string) string
}