package logging

import (
	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - API LOGGER
// TDD + Zero Hardcoding Policy
// ========================================

// apiLogger implémentation du APILogger
type apiLogger struct {
	logger *FireSalamanderLogger
}

// NewAPILogger crée un nouveau logger API
func NewAPILogger(logger *FireSalamanderLogger) APILogger {
	return &apiLogger{
		logger: logger,
	}
}

// RequestReceived log une requête API reçue
func (a *apiLogger) RequestReceived(endpoint, method string, requestID string, data map[string]interface{}) {
	logData := map[string]interface{}{
		"endpoint": endpoint,
		"method":   method,
	}
	
	// Ajouter les données de la requête
	if data != nil {
		for k, v := range data {
			logData[k] = v
		}
	}
	
	contextLogger := a.logger.WithRequestID(requestID)
	contextLogger.Info(constants.LogCategoryAPI, constants.LogMsgAPIRequest, logData)
}

// ResponseSent log une réponse API envoyée
func (a *apiLogger) ResponseSent(endpoint, method string, statusCode int, responseTimeMs int64, requestID string) {
	data := map[string]interface{}{
		"endpoint":                          endpoint,
		"method":                            method,
		constants.HTTPLogFieldStatusCode:    statusCode,
		constants.HTTPLogFieldResponseTime:  responseTimeMs,
	}
	
	contextLogger := a.logger.WithRequestID(requestID)
	
	// Choisir le niveau selon le status code
	if statusCode >= 500 {
		contextLogger.Error(constants.LogCategoryAPI, constants.LogMsgAPIResponse, nil, data)
	} else if statusCode >= 400 {
		contextLogger.Warn(constants.LogCategoryAPI, constants.LogMsgAPIResponse, data)
	} else {
		contextLogger.Info(constants.LogCategoryAPI, constants.LogMsgAPIResponse, data)
	}
}

// ValidationError log une erreur de validation API
func (a *apiLogger) ValidationError(endpoint, method string, err error, requestID string) {
	data := map[string]interface{}{
		"endpoint": endpoint,
		"method":   method,
		"error_type": "validation",
	}
	
	contextLogger := a.logger.WithRequestID(requestID)
	contextLogger.Error(constants.LogCategoryAPI, constants.LogMsgAPIValidation, err, data)
}

// ProcessingError log une erreur de traitement API
func (a *apiLogger) ProcessingError(endpoint, method string, err error, requestID string) {
	data := map[string]interface{}{
		"endpoint": endpoint,
		"method":   method,
		"error_type": "processing",
	}
	
	contextLogger := a.logger.WithRequestID(requestID)
	contextLogger.Error(constants.LogCategoryAPI, constants.LogMsgAPIError, err, data)
}