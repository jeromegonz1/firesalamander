package logging

import (
	"encoding/json"
	"time"

	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - AUDIT LOGGER
// TDD + Zero Hardcoding Policy
// ========================================

// auditLogger implémentation du AuditLogger
type auditLogger struct {
	logger *FireSalamanderLogger
}

// NewAuditLogger crée un nouveau logger d'audit
func NewAuditLogger(logger *FireSalamanderLogger) AuditLogger {
	return &auditLogger{
		logger: logger,
	}
}

// ActionStarted log le début d'une action critique
func (a *auditLogger) ActionStarted(action, resource, userID, sessionID string, beforeState interface{}) {
	beforeJSON, _ := json.Marshal(beforeState)
	
	data := map[string]interface{}{
		constants.AuditLogFieldAction:    action,
		constants.AuditLogFieldResource:  resource,
		constants.AuditLogFieldUserID:    userID,
		constants.AuditLogFieldSessionID: sessionID,
		constants.AuditLogFieldBefore:    string(beforeJSON),
		"status":                         "started",
		"start_time":                     time.Now().Format(constants.LogTimestampFormat),
	}
	
	a.logger.Info(constants.LogCategoryAudit, "Action started", data)
}

// ActionCompleted log la fin d'une action critique
func (a *auditLogger) ActionCompleted(action, resource, userID, sessionID string, success bool, afterState interface{}) {
	afterJSON, _ := json.Marshal(afterState)
	
	data := map[string]interface{}{
		constants.AuditLogFieldAction:    action,
		constants.AuditLogFieldResource:  resource,
		constants.AuditLogFieldUserID:    userID,
		constants.AuditLogFieldSessionID: sessionID,
		constants.AuditLogFieldAfter:     string(afterJSON),
		constants.AuditLogFieldSuccess:   success,
		"status":                         "completed",
		"end_time":                       time.Now().Format(constants.LogTimestampFormat),
	}
	
	if success {
		a.logger.Info(constants.LogCategoryAudit, "Action completed successfully", data)
	} else {
		a.logger.Warn(constants.LogCategoryAudit, "Action failed", data)
	}
}

// SecurityEvent log un événement de sécurité
func (a *auditLogger) SecurityEvent(event, description string, data map[string]interface{}) {
	logData := map[string]interface{}{
		"event":       event,
		"description": description,
		"timestamp":   time.Now().Format(constants.LogTimestampFormat),
	}
	
	// Ajouter les données supplémentaires
	if data != nil {
		for k, v := range data {
			logData[k] = v
		}
	}
	
	// Tous les événements de sécurité sont des warnings minimum
	a.logger.Warn(constants.LogCategoryAudit, "Security event", logData)
}