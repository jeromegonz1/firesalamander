/**
 * 🚨 FIRE SALAMANDER ERROR LOGGER
 * Système de logging pour traquer les erreurs critiques comme les boucles infinies
 * 
 * Usage:
 * - Intégrer dans app.js pour monitorer Chart.js
 * - Détecter automatiquement les problèmes de performance
 * - Logger dans /fix/errors.log pour analyse
 */

class FireSalamanderErrorLogger {
    constructor() {
        this.errors = [];
        this.performanceMonitor = {
            heightChecks: 0,
            lastHeight: 0,
            heightGrowthCount: 0,
            startTime: Date.now()
        };
        
        this.init();
    }

    init() {
        console.log('🔍 Fire Salamander Error Logger initialized');
        
        // Monitor DOM mutations pour détecter croissance anormale
        this.observeHeightChanges();
        
        // Monitor JavaScript errors
        this.captureJSErrors();
        
        // Monitor Chart.js specifically
        this.monitorChartJS();
        
        // Performance monitoring
        this.startPerformanceMonitoring();
    }

    /**
     * CRITICAL: Detect infinite height loops
     */
    observeHeightChanges() {
        const observer = new MutationObserver(() => {
            const currentHeight = Math.max(
                document.body.scrollHeight,
                document.documentElement.scrollHeight
            );
            
            // Detect abnormal height growth
            if (currentHeight > this.performanceMonitor.lastHeight + 1000) {
                this.performanceMonitor.heightGrowthCount++;
                
                if (currentHeight > 50000) { // 50k pixels = suspicious
                    this.logCriticalError('INFINITE_HEIGHT_LOOP', {
                        currentHeight,
                        previousHeight: this.performanceMonitor.lastHeight,
                        growthCount: this.performanceMonitor.heightGrowthCount,
                        timeElapsed: Date.now() - this.performanceMonitor.startTime
                    });
                }
            }
            
            this.performanceMonitor.lastHeight = currentHeight;
            this.performanceMonitor.heightChecks++;
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ['style', 'class']
        });
    }

    /**
     * Monitor Chart.js creation for problematic configurations
     */
    monitorChartJS() {
        // Override Chart constructor if available
        if (window.Chart) {
            const originalChart = window.Chart;
            
            window.Chart = (...args) => {
                const [ctx, config] = args;
                
                // Log Chart.js creation
                this.logInfo('CHART_CREATED', {
                    canvasId: ctx.canvas?.id,
                    type: config?.type,
                    maintainAspectRatio: config?.options?.maintainAspectRatio,
                    responsive: config?.options?.responsive
                });
                
                // CRITICAL: Detect dangerous configurations
                if (config?.options?.maintainAspectRatio === false) {
                    this.logCriticalError('DANGEROUS_CHART_CONFIG', {
                        canvasId: ctx.canvas?.id,
                        type: config.type,
                        config: config.options,
                        warning: 'maintainAspectRatio: false can cause infinite loops!'
                    });
                }
                
                return new originalChart(...args);
            };
            
            // Preserve Chart properties
            Object.setPrototypeOf(window.Chart, originalChart);
            Object.keys(originalChart).forEach(key => {
                window.Chart[key] = originalChart[key];
            });
        }
    }

    /**
     * Capture JavaScript errors
     */
    captureJSErrors() {
        window.addEventListener('error', (event) => {
            this.logError('JAVASCRIPT_ERROR', {
                message: event.message,
                filename: event.filename,
                lineno: event.lineno,
                colno: event.colno,
                stack: event.error?.stack
            });
        });

        window.addEventListener('unhandledrejection', (event) => {
            this.logError('UNHANDLED_PROMISE_REJECTION', {
                reason: event.reason,
                promise: event.promise
            });
        });
    }

    /**
     * Monitor performance metrics
     */
    startPerformanceMonitoring() {
        setInterval(() => {
            const metrics = this.getPerformanceMetrics();
            
            // Log si performance dégradée
            if (metrics.memoryUsage > 100 || metrics.height > 20000) {
                this.logWarning('PERFORMANCE_DEGRADATION', metrics);
            }
        }, 5000); // Check every 5 seconds
    }

    getPerformanceMetrics() {
        const height = Math.max(
            document.body.scrollHeight,
            document.documentElement.scrollHeight
        );
        
        return {
            timestamp: new Date().toISOString(),
            height,
            heightChecks: this.performanceMonitor.heightChecks,
            heightGrowthCount: this.performanceMonitor.heightGrowthCount,
            memoryUsage: performance.memory?.usedJSHeapSize / 1024 / 1024 || 0,
            chartCount: document.querySelectorAll('canvas').length,
            uptime: Date.now() - this.performanceMonitor.startTime
        };
    }

    /**
     * Logging methods
     */
    logCriticalError(type, data) {
        const error = {
            level: 'CRITICAL',
            type,
            timestamp: new Date().toISOString(),
            data,
            userAgent: navigator.userAgent,
            url: window.location.href,
            performance: this.getPerformanceMetrics()
        };
        
        this.errors.push(error);
        console.error('🚨 CRITICAL ERROR:', error);
        
        // Save to localStorage for persistence
        this.saveToStorage(error);
        
        // Try to send to server (if available)
        this.sendToServer(error);
    }

    logError(type, data) {
        const error = {
            level: 'ERROR',
            type,
            timestamp: new Date().toISOString(),
            data
        };
        
        this.errors.push(error);
        console.error('❌ ERROR:', error);
        this.saveToStorage(error);
    }

    logWarning(type, data) {
        const warning = {
            level: 'WARNING',
            type,
            timestamp: new Date().toISOString(),
            data
        };
        
        this.errors.push(warning);
        console.warn('⚠️ WARNING:', warning);
        this.saveToStorage(warning);
    }

    logInfo(type, data) {
        const info = {
            level: 'INFO',
            type,
            timestamp: new Date().toISOString(),
            data
        };
        
        console.log('ℹ️ INFO:', info);
    }

    /**
     * Persistence
     */
    saveToStorage(error) {
        try {
            const stored = JSON.parse(localStorage.getItem('fire-salamander-errors') || '[]');
            stored.push(error);
            
            // Keep only last 100 errors
            if (stored.length > 100) {
                stored.splice(0, stored.length - 100);
            }
            
            localStorage.setItem('fire-salamander-errors', JSON.stringify(stored));
        } catch (e) {
            console.error('Failed to save error to localStorage:', e);
        }
    }

    sendToServer(error) {
        // Try to send critical errors to server
        if (window.fetch) {
            fetch('/api/v1/errors', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(error)
            }).catch(() => {
                // Silent fail - server might not have this endpoint
            });
        }
    }

    /**
     * Public API
     */
    getStoredErrors() {
        try {
            return JSON.parse(localStorage.getItem('fire-salamander-errors') || '[]');
        } catch (e) {
            return [];
        }
    }

    clearStoredErrors() {
        localStorage.removeItem('fire-salamander-errors');
        this.errors = [];
        console.log('🧹 Error log cleared');
    }

    exportErrors() {
        const errors = this.getStoredErrors();
        const blob = new Blob([JSON.stringify(errors, null, 2)], {
            type: 'application/json'
        });
        
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `fire-salamander-errors-${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    // Debug: Force height monitoring test
    simulateInfiniteLoop() {
        console.warn('🧪 SIMULATING infinite height loop for testing...');
        document.body.style.height = '100000px';
        setTimeout(() => {
            document.body.style.height = '200000px';
        }, 1000);
    }
}

// Auto-initialize if in browser environment
if (typeof window !== 'undefined') {
    window.FireSalamanderErrorLogger = FireSalamanderErrorLogger;
    
    // Auto-start logging
    document.addEventListener('DOMContentLoaded', () => {
        window.fireSalamanderLogger = new FireSalamanderErrorLogger();
        
        // Expose global functions for debugging
        window.getFireSalamanderErrors = () => window.fireSalamanderLogger.getStoredErrors();
        window.clearFireSalamanderErrors = () => window.fireSalamanderLogger.clearStoredErrors();
        window.exportFireSalamanderErrors = () => window.fireSalamanderLogger.exportErrors();
        
        console.log('🔥 Fire Salamander Error Logger ready! Available commands:');
        console.log('- getFireSalamanderErrors() - View logged errors');
        console.log('- clearFireSalamanderErrors() - Clear error log');
        console.log('- exportFireSalamanderErrors() - Export errors as JSON');
    });
}