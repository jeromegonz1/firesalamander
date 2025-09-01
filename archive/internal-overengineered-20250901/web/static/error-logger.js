/**
 * 🔥 FIRE SALAMANDER - ERROR LOGGER & MONITORING
 * Système de monitoring pour détecter les problèmes de boucles infinies et autres erreurs critiques
 */

console.log('🔍 Fire Salamander Error Logger Starting...');

class FireSalamanderLogger {
    constructor() {
        this.errors = [];
        this.heightHistory = [];
        this.maxHeightThreshold = 50000; // 50k pixels threshold
        this.heightCheckInterval = null;
        this.startTime = Date.now();
        
        this.initialize();
    }
    
    initialize() {
        // Monitor DOM height changes
        this.startHeightMonitoring();
        
        // Monitor JavaScript errors
        this.setupErrorHandling();
        
        // Monitor Chart.js dangerous configurations
        this.monitorChartJsConfigs();
        
        // Monitor memory usage if available
        this.monitorMemoryUsage();
        
        console.log('✅ Fire Salamander Error Logger initialized');
    }
    
    startHeightMonitoring() {
        this.heightCheckInterval = setInterval(() => {
            const currentHeight = Math.max(
                document.body.scrollHeight,
                document.body.offsetHeight,
                document.documentElement.clientHeight,
                document.documentElement.scrollHeight,
                document.documentElement.offsetHeight
            );
            
            this.heightHistory.push({
                timestamp: Date.now(),
                height: currentHeight
            });
            
            // Keep only last 50 measurements
            if (this.heightHistory.length > 50) {
                this.heightHistory.shift();
            }
            
            // Check for dangerous height
            if (currentHeight > this.maxHeightThreshold) {
                this.logCriticalError('INFINITE_HEIGHT_LOOP', {
                    currentHeight,
                    threshold: this.maxHeightThreshold,
                    heightHistory: this.heightHistory.slice(-10),
                    timestamp: new Date().toISOString()
                });
            }
            
            // Check for rapid growth
            if (this.heightHistory.length >= 3) {
                const recent = this.heightHistory.slice(-3);
                const growth = recent[2].height - recent[0].height;
                const timeSpan = recent[2].timestamp - recent[0].timestamp;
                const growthRate = growth / (timeSpan / 1000); // pixels per second
                
                if (growthRate > 1000) { // More than 1000px/s growth
                    this.logWarning('RAPID_HEIGHT_GROWTH', {
                        growthRate: Math.round(growthRate),
                        recent,
                        timestamp: new Date().toISOString()
                    });
                }
            }
            
        }, 1000); // Check every second
    }
    
    setupErrorHandling() {
        // Catch unhandled JavaScript errors
        window.addEventListener('error', (event) => {
            this.logError('JAVASCRIPT_ERROR', {
                message: event.message,
                filename: event.filename,
                lineno: event.lineno,
                colno: event.colno,
                stack: event.error ? event.error.stack : null,
                timestamp: new Date().toISOString()
            });
        });
        
        // Catch unhandled promise rejections
        window.addEventListener('unhandledrejection', (event) => {
            // Ignore empty rejections which are common in normal operation
            if (!event.reason || 
                (typeof event.reason === 'object' && Object.keys(event.reason).length === 0) ||
                event.reason === '') {
                // Prevent default error logging for empty/harmless rejections  
                event.preventDefault();
                return;
            }
            
            this.logWarning('UNHANDLED_PROMISE_REJECTION', {
                reason: event.reason,
                timestamp: new Date().toISOString()
            });
        });
    }
    
    monitorChartJsConfigs() {
        // Intercept Chart.js creation to check for dangerous configs
        const originalConsoleWarn = console.warn;
        console.warn = (...args) => {
            if (args.some(arg => typeof arg === 'string' && arg.includes('maintainAspectRatio'))) {
                this.logWarning('DANGEROUS_CHARTJS_CONFIG', {
                    warning: args.join(' '),
                    timestamp: new Date().toISOString()
                });
            }
            originalConsoleWarn.apply(console, args);
        };
        
        // Monitor for Chart.js infinite loop patterns
        const originalConsoleError = console.error;
        console.error = (...args) => {
            const errorMessage = args.join(' ').toLowerCase();
            if (errorMessage.includes('canvas') && errorMessage.includes('resize')) {
                this.logCriticalError('CHARTJS_RESIZE_LOOP', {
                    error: args.join(' '),
                    timestamp: new Date().toISOString()
                });
            }
            originalConsoleError.apply(console, args);
        };
    }
    
    monitorMemoryUsage() {
        if ('memory' in performance) {
            setInterval(() => {
                const memory = performance.memory;
                const memoryUsage = {
                    used: Math.round(memory.usedJSHeapSize / 1048576), // MB
                    total: Math.round(memory.totalJSHeapSize / 1048576), // MB
                    limit: Math.round(memory.jsHeapSizeLimit / 1048576), // MB
                    timestamp: Date.now()
                };
                
                // Check for memory leaks (>500MB or >80% of limit)
                if (memoryUsage.used > 500 || memoryUsage.used > (memoryUsage.limit * 0.8)) {
                    this.logWarning('HIGH_MEMORY_USAGE', memoryUsage);
                }
                
            }, 10000); // Check every 10 seconds
        }
    }
    
    logInfo(type, data) {
        const entry = {
            level: 'INFO',
            type,
            data,
            timestamp: new Date().toISOString(),
            sessionId: this.getSessionId()
        };
        
        this.errors.push(entry);
        console.log(`ℹ️  [${type}]`, data);
        this.saveToStorage();
    }
    
    logWarning(type, data) {
        const entry = {
            level: 'WARNING',
            type,
            data,
            timestamp: new Date().toISOString(),
            sessionId: this.getSessionId()
        };
        
        this.errors.push(entry);
        console.warn(`⚠️  [${type}]`, data);
        this.saveToStorage();
        
        // Show user notification for warnings
        this.showNotification('warning', `Warning: ${type}`, data);
    }
    
    logError(type, data) {
        const entry = {
            level: 'ERROR',
            type,
            data,
            timestamp: new Date().toISOString(),
            sessionId: this.getSessionId()
        };
        
        this.errors.push(entry);
        console.error(`❌ [${type}]`, data);
        this.saveToStorage();
        
        // Show user notification for errors
        this.showNotification('error', `Error: ${type}`, data);
    }
    
    logCriticalError(type, data) {
        const entry = {
            level: 'CRITICAL',
            type,
            data,
            timestamp: new Date().toISOString(),
            sessionId: this.getSessionId()
        };
        
        this.errors.push(entry);
        console.error(`🚨 [CRITICAL] [${type}]`, data);
        this.saveToStorage();
        
        // Show prominent notification for critical errors
        this.showNotification('critical', `CRITICAL: ${type}`, data);
        
        // Try to recover from infinite loops
        if (type === 'INFINITE_HEIGHT_LOOP') {
            this.attemptInfiniteLoopRecovery();
        }
    }
    
    attemptInfiniteLoopRecovery() {
        console.log('🩹 Attempting to recover from infinite loop...');
        
        // Stop all Chart.js instances
        if (window.Chart && window.Chart.instances) {
            Object.values(window.Chart.instances).forEach(chart => {
                try {
                    chart.destroy();
                    console.log('🔧 Destroyed Chart.js instance:', chart.id);
                } catch (e) {
                    console.error('Failed to destroy chart:', e);
                }
            });
        }
        
        // Reset body styles that might cause infinite growth
        document.body.style.height = 'auto';
        document.body.style.maxHeight = '100vh';
        document.documentElement.style.height = 'auto';
        document.documentElement.style.maxHeight = '100vh';
        
        // Find and fix problematic chart containers
        const chartContainers = document.querySelectorAll('.chart-container, canvas');
        chartContainers.forEach(container => {
            container.style.height = '300px';
            container.style.maxHeight = '300px';
            container.style.overflow = 'hidden';
        });
        
        this.logInfo('RECOVERY_ATTEMPTED', {
            chartsDestroyed: window.Chart && window.Chart.instances ? Object.keys(window.Chart.instances).length : 0,
            containersFixed: chartContainers.length,
            timestamp: new Date().toISOString()
        });
    }
    
    showNotification(level, title, data) {
        const notification = document.createElement('div');
        notification.className = `error-logger-notification error-logger-${level}`;
        notification.innerHTML = `
            <div class="notification-header">
                <span class="notification-icon">${this.getNotificationIcon(level)}</span>
                <span class="notification-title">${title}</span>
                <button class="notification-close" onclick="this.parentElement.parentElement.remove()">×</button>
            </div>
            <div class="notification-body">
                <small>${JSON.stringify(data, null, 2)}</small>
            </div>
        `;
        
        // Add styles if not present
        if (!document.getElementById('error-logger-styles')) {
            const styles = document.createElement('style');
            styles.id = 'error-logger-styles';
            styles.textContent = `
                .error-logger-notification {
                    position: fixed;
                    top: 20px;
                    right: 20px;
                    min-width: 300px;
                    max-width: 500px;
                    border-radius: 8px;
                    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                    z-index: 10000;
                    font-family: monospace;
                    font-size: 12px;
                    animation: slideIn 0.3s ease-out;
                }
                .error-logger-warning { background: #fff3cd; border-left: 4px solid #ffc107; }
                .error-logger-error { background: #f8d7da; border-left: 4px solid #dc3545; }
                .error-logger-critical { background: #d4a5a5; border-left: 4px solid #8b0000; animation: pulse 1s infinite; }
                .notification-header {
                    padding: 10px;
                    border-bottom: 1px solid rgba(0,0,0,0.1);
                    display: flex;
                    align-items: center;
                    gap: 8px;
                }
                .notification-close {
                    margin-left: auto;
                    background: none;
                    border: none;
                    font-size: 18px;
                    cursor: pointer;
                }
                .notification-body {
                    padding: 10px;
                    max-height: 200px;
                    overflow-y: auto;
                }
                @keyframes slideIn {
                    from { transform: translateX(100%); opacity: 0; }
                    to { transform: translateX(0); opacity: 1; }
                }
                @keyframes pulse {
                    0%, 100% { opacity: 1; }
                    50% { opacity: 0.7; }
                }
            `;
            document.head.appendChild(styles);
        }
        
        document.body.appendChild(notification);
        
        // Auto-remove after delay (except critical)
        if (level !== 'critical') {
            setTimeout(() => {
                if (notification.parentElement) {
                    notification.remove();
                }
            }, level === 'warning' ? 5000 : 8000);
        }
    }
    
    getNotificationIcon(level) {
        const icons = {
            warning: '⚠️',
            error: '❌',
            critical: '🚨'
        };
        return icons[level] || 'ℹ️';
    }
    
    saveToStorage() {
        try {
            // Keep only last 100 errors
            if (this.errors.length > 100) {
                this.errors = this.errors.slice(-100);
            }
            
            localStorage.setItem('fire-salamander-errors', JSON.stringify({
                errors: this.errors,
                lastUpdate: Date.now(),
                sessionId: this.getSessionId()
            }));
        } catch (e) {
            console.error('Failed to save errors to localStorage:', e);
        }
    }
    
    getSessionId() {
        if (!this.sessionId) {
            this.sessionId = `fs-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
        }
        return this.sessionId;
    }
    
    exportErrors() {
        return {
            errors: this.errors,
            heightHistory: this.heightHistory,
            sessionInfo: {
                sessionId: this.getSessionId(),
                startTime: this.startTime,
                userAgent: navigator.userAgent,
                url: window.location.href,
                timestamp: new Date().toISOString()
            }
        };
    }
    
    clearErrors() {
        this.errors = [];
        this.heightHistory = [];
        localStorage.removeItem('fire-salamander-errors');
        console.log('🧹 Error logger cleared');
    }
    
    getStats() {
        const levels = this.errors.reduce((acc, error) => {
            acc[error.level] = (acc[error.level] || 0) + 1;
            return acc;
        }, {});
        
        return {
            totalErrors: this.errors.length,
            byLevel: levels,
            sessionDuration: Date.now() - this.startTime,
            currentHeight: document.body.scrollHeight,
            maxHeightReached: Math.max(...this.heightHistory.map(h => h.height), 0),
            memoryUsage: 'memory' in performance ? performance.memory : null
        };
    }
}

// Initialize global logger
if (typeof window !== 'undefined') {
    window.fireSalamanderLogger = new FireSalamanderLogger();
    
    // Expose logger functions globally
    window.logFireSalamanderError = (type, data) => window.fireSalamanderLogger.logError(type, data);
    window.logFireSalamanderWarning = (type, data) => window.fireSalamanderLogger.logWarning(type, data);
    window.exportFireSalamanderErrors = () => window.fireSalamanderLogger.exportErrors();
    window.clearFireSalamanderErrors = () => window.fireSalamanderLogger.clearErrors();
    window.getFireSalamanderStats = () => window.fireSalamanderLogger.getStats();
    
    console.log('✅ Fire Salamander Error Logger ready');
    console.log('📊 Available commands: exportFireSalamanderErrors(), clearFireSalamanderErrors(), getFireSalamanderStats()');
}