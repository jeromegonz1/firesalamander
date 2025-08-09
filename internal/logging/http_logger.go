package logging

import (
	"firesalamander/internal/constants"
)

// ========================================
// FIRE SALAMANDER - HTTP LOGGER
// TDD + Zero Hardcoding Policy
// ========================================

// httpLogger implémentation du HTTPLogger
type httpLogger struct {
	logger *FireSalamanderLogger
}

// NewHTTPLogger crée un nouveau logger HTTP
func NewHTTPLogger(logger *FireSalamanderLogger) HTTPLogger {
	return &httpLogger{
		logger: logger,
	}
}

// RequestReceived log une requête HTTP reçue
func (h *httpLogger) RequestReceived(method, url, remoteAddr, userAgent string, requestID string) {
	data := map[string]interface{}{
		constants.HTTPLogFieldMethod:     method,
		constants.HTTPLogFieldURL:        url,
		constants.HTTPLogFieldRemoteAddr: remoteAddr,
		constants.HTTPLogFieldUserAgent:  userAgent,
		constants.HTTPLogFieldRequestID:  requestID,
	}
	
	contextLogger := h.logger.WithRequestID(requestID)
	contextLogger.Info(constants.LogCategoryHTTP, constants.LogMsgHTTPRequest, data)
}

// ResponseSent log une réponse HTTP envoyée
func (h *httpLogger) ResponseSent(method, url string, statusCode int, responseTimeMs int64, contentLength int64, requestID string) {
	data := map[string]interface{}{
		constants.HTTPLogFieldMethod:        method,
		constants.HTTPLogFieldURL:           url,
		constants.HTTPLogFieldStatusCode:    statusCode,
		constants.HTTPLogFieldResponseTime:  responseTimeMs,
		constants.HTTPLogFieldContentLength: contentLength,
		constants.HTTPLogFieldRequestID:     requestID,
	}
	
	contextLogger := h.logger.WithRequestID(requestID)
	
	// Choisir le niveau selon le status code
	if statusCode >= 500 {
		contextLogger.Error(constants.LogCategoryHTTP, constants.LogMsgHTTPResponse, nil, data)
	} else if statusCode >= 400 {
		contextLogger.Warn(constants.LogCategoryHTTP, constants.LogMsgHTTPResponse, data)
	} else {
		contextLogger.Info(constants.LogCategoryHTTP, constants.LogMsgHTTPResponse, data)
	}
}

// RequestError log une erreur de requête HTTP
func (h *httpLogger) RequestError(method, url string, err error, requestID string) {
	data := map[string]interface{}{
		constants.HTTPLogFieldMethod:    method,
		constants.HTTPLogFieldURL:       url,
		constants.HTTPLogFieldRequestID: requestID,
	}
	
	contextLogger := h.logger.WithRequestID(requestID)
	contextLogger.Error(constants.LogCategoryHTTP, constants.LogMsgHTTPError, err, data)
}