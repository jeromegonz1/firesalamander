/**
 * Fire Salamander - API Test Agent
 * Tests de contrat, charge et sécurité pour les APIs REST
 */

const axios = require('axios');
const fs = require('fs').promises;
const path = require('path');

class APITestAgent {
    constructor(config = {}) {
        this.config = {
            baseURL: config.baseURL || 'http://localhost:3000',
            timeout: config.timeout || 30000,
            maxRetries: config.maxRetries || 3,
            reportPath: config.reportPath || 'tests/reports/api',
            ...config
        };
        
        this.stats = {
            timestamp: new Date().toISOString(),
            contractTests: [],
            loadTests: [],
            securityTests: [],
            overallScore: 0,
            status: 'unknown'
        };
        
        this.client = axios.create({
            baseURL: this.config.baseURL,
            timeout: this.config.timeout,
            validateStatus: () => true // Accepter tous les status codes
        });
    }

    /**
     * Exécute tous les tests API
     */
    async runFullTestSuite() {
        console.log('🚀 Starting API Test Suite');
        
        try {
            // 1. Tests de contrat (Contract Testing)
            await this.runContractTests();
            
            // 2. Tests de charge (Load Testing)
            await this.runLoadTests();
            
            // 3. Tests de sécurité (Security Testing)
            await this.runSecurityTests();
            
            // Calculer le score global
            this.calculateOverallScore();
            
            // Générer le rapport
            await this.generateReport();
            
            console.log(`✅ API Test Suite completed - Score: ${this.stats.overallScore}/100`);
            
        } catch (error) {
            console.error('❌ API Test Suite failed:', error.message);
            throw error;
        }
        
        return this.stats;
    }

    /**
     * Tests de contrat - Vérification des API endpoints
     */
    async runContractTests() {
        console.log('📋 Running Contract Tests');
        
        const endpoints = [
            { method: 'GET', path: '/', expectedStatus: 200, description: 'Home page' },
            { method: 'GET', path: '/health', expectedStatus: 200, description: 'Health check' },
            { method: 'GET', path: '/debug', expectedStatus: 200, description: 'Debug endpoint' },
            { method: 'GET', path: '/api/nonexistent', expectedStatus: 404, description: 'Non-existent endpoint' }
        ];

        for (const endpoint of endpoints) {
            const testResult = await this.testEndpoint(endpoint);
            this.stats.contractTests.push(testResult);
        }

        // Tests de validation des données
        await this.testDataValidation();
        
        // Tests de pagination
        await this.testPagination();
        
        // Tests d'authentification
        await this.testAuthentication();
    }

    /**
     * Tests de charge - Performance et scalabilité
     */
    async runLoadTests() {
        console.log('⚡ Running Load Tests');
        
        const scenarios = [
            { name: 'Light Load', concurrent: 5, requests: 50 },
            { name: 'Normal Load', concurrent: 10, requests: 100 },
            { name: 'Heavy Load', concurrent: 20, requests: 200 }
        ];

        for (const scenario of scenarios) {
            const result = await this.runLoadScenario(scenario);
            this.stats.loadTests.push(result);
        }
    }

    /**
     * Tests de sécurité - OWASP Top 10
     */
    async runSecurityTests() {
        console.log('🔒 Running Security Tests');
        
        // SQL Injection
        await this.testSQLInjection();
        
        // XSS
        await this.testXSS();
        
        // CSRF
        await this.testCSRF();
        
        // Rate Limiting
        await this.testRateLimiting();
        
        // Headers de sécurité
        await this.testSecurityHeaders();
        
        // Authentication bypass
        await this.testAuthBypass();
    }

    /**
     * Test un endpoint spécifique
     */
    async testEndpoint(endpoint) {
        const startTime = Date.now();
        const testResult = {
            method: endpoint.method,
            path: endpoint.path,
            description: endpoint.description,
            expectedStatus: endpoint.expectedStatus,
            actualStatus: null,
            responseTime: null,
            passed: false,
            errors: [],
            response: null
        };

        try {
            const response = await this.client.request({
                method: endpoint.method,
                url: endpoint.path,
                data: endpoint.data || undefined
            });

            testResult.actualStatus = response.status;
            testResult.responseTime = Date.now() - startTime;
            testResult.response = {
                headers: response.headers,
                data: typeof response.data === 'string' ? response.data.substring(0, 1000) : response.data
            };

            // Vérifier le status code
            if (response.status === endpoint.expectedStatus) {
                testResult.passed = true;
            } else {
                testResult.errors.push(`Expected status ${endpoint.expectedStatus}, got ${response.status}`);
            }

            // Vérifications supplémentaires
            if (endpoint.expectedHeaders) {
                for (const [header, expectedValue] of Object.entries(endpoint.expectedHeaders)) {
                    const actualValue = response.headers[header.toLowerCase()];
                    if (actualValue !== expectedValue) {
                        testResult.errors.push(`Expected header ${header}: ${expectedValue}, got: ${actualValue}`);
                        testResult.passed = false;
                    }
                }
            }

            // Vérifier le temps de réponse
            if (testResult.responseTime > 5000) {
                testResult.errors.push(`Response time too slow: ${testResult.responseTime}ms`);
            }

        } catch (error) {
            testResult.errors.push(`Request failed: ${error.message}`);
            testResult.responseTime = Date.now() - startTime;
        }

        console.log(`${testResult.passed ? '✅' : '❌'} ${endpoint.method} ${endpoint.path} - ${testResult.responseTime}ms`);
        return testResult;
    }

    /**
     * Tests de validation des données
     */
    async testDataValidation() {
        const validationTests = [
            {
                method: 'POST',
                path: '/api/sites',
                data: { url: 'invalid-url' },
                expectedStatus: 400,
                description: 'Invalid URL validation'
            },
            {
                method: 'POST',
                path: '/api/sites',
                data: { url: '' },
                expectedStatus: 400,
                description: 'Empty URL validation'
            }
        ];

        for (const test of validationTests) {
            const result = await this.testEndpoint(test);
            this.stats.contractTests.push(result);
        }
    }

    /**
     * Tests de pagination
     */
    async testPagination() {
        const paginationTests = [
            {
                method: 'GET',
                path: '/api/sites?page=1&limit=10',
                expectedStatus: 200,
                description: 'Pagination - first page'
            },
            {
                method: 'GET',
                path: '/api/sites?page=999&limit=10',
                expectedStatus: 200,
                description: 'Pagination - out of range page'
            }
        ];

        for (const test of paginationTests) {
            const result = await this.testEndpoint(test);
            this.stats.contractTests.push(result);
        }
    }

    /**
     * Tests d'authentification
     */
    async testAuthentication() {
        const authTests = [
            {
                method: 'GET',
                path: '/api/admin',
                expectedStatus: 401,
                description: 'Protected endpoint without auth'
            }
        ];

        for (const test of authTests) {
            const result = await this.testEndpoint(test);
            this.stats.contractTests.push(result);
        }
    }

    /**
     * Exécute un scénario de charge
     */
    async runLoadScenario(scenario) {
        console.log(`⚡ Running ${scenario.name} (${scenario.concurrent} concurrent, ${scenario.requests} requests)`);
        
        const results = {
            name: scenario.name,
            concurrent: scenario.concurrent,
            totalRequests: scenario.requests,
            successCount: 0,
            errorCount: 0,
            responseTimes: [],
            errors: [],
            averageResponseTime: 0,
            p95ResponseTime: 0,
            throughput: 0,
            passed: false
        };

        const startTime = Date.now();
        const promises = [];

        // Diviser les requêtes en batches concurrents
        const batchSize = scenario.concurrent;
        const batches = Math.ceil(scenario.requests / batchSize);

        for (let batch = 0; batch < batches; batch++) {
            const batchPromises = [];
            const requestsInThisBatch = Math.min(batchSize, scenario.requests - (batch * batchSize));

            for (let i = 0; i < requestsInThisBatch; i++) {
                batchPromises.push(this.performLoadRequest());
            }

            const batchResults = await Promise.allSettled(batchPromises);
            
            for (const result of batchResults) {
                if (result.status === 'fulfilled') {
                    results.successCount++;
                    results.responseTimes.push(result.value.responseTime);
                } else {
                    results.errorCount++;
                    results.errors.push(result.reason.message);
                }
            }
        }

        const totalTime = Date.now() - startTime;
        
        // Calculer les métriques
        if (results.responseTimes.length > 0) {
            results.averageResponseTime = results.responseTimes.reduce((a, b) => a + b, 0) / results.responseTimes.length;
            results.responseTimes.sort((a, b) => a - b);
            results.p95ResponseTime = results.responseTimes[Math.floor(results.responseTimes.length * 0.95)];
        }
        
        results.throughput = (results.successCount / totalTime) * 1000; // req/sec
        
        // Critères de réussite
        results.passed = results.errorCount === 0 && 
                         results.averageResponseTime < 1000 && 
                         results.p95ResponseTime < 2000;

        console.log(`${results.passed ? '✅' : '❌'} ${scenario.name} - Success: ${results.successCount}/${scenario.requests}, Avg: ${Math.round(results.averageResponseTime)}ms`);
        
        return results;
    }

    /**
     * Effectue une requête de charge
     */
    async performLoadRequest() {
        const startTime = Date.now();
        
        try {
            const response = await this.client.get('/health');
            const responseTime = Date.now() - startTime;
            
            if (response.status === 200) {
                return { responseTime, success: true };
            } else {
                throw new Error(`HTTP ${response.status}`);
            }
        } catch (error) {
            throw new Error(`Load request failed: ${error.message}`);
        }
    }

    /**
     * Tests d'injection SQL
     */
    async testSQLInjection() {
        const injectionPayloads = [
            "'; DROP TABLE users; --",
            "' OR '1'='1",
            "'; UNION SELECT * FROM users; --",
            "admin'--",
            "admin' OR 1=1--"
        ];

        for (const payload of injectionPayloads) {
            const result = await this.testEndpoint({
                method: 'GET',
                path: `/api/search?q=${encodeURIComponent(payload)}`,
                expectedStatus: 400,
                description: `SQL Injection test: ${payload.substring(0, 20)}...`
            });
            
            // Vérifier que l'injection n'a pas réussi
            if (result.actualStatus === 200 && result.response?.data?.includes('users')) {
                result.passed = false;
                result.errors.push('Possible SQL injection vulnerability');
            }
            
            this.stats.securityTests.push(result);
        }
    }

    /**
     * Tests XSS
     */
    async testXSS() {
        const xssPayloads = [
            "<script>alert('XSS')</script>",
            "<img src=x onerror=alert('XSS')>",
            "javascript:alert('XSS')",
            "<svg onload=alert('XSS')>"
        ];

        for (const payload of xssPayloads) {
            const result = await this.testEndpoint({
                method: 'POST',
                path: '/api/feedback',
                data: { message: payload },
                expectedStatus: 400,
                description: `XSS test: ${payload.substring(0, 20)}...`
            });
            
            this.stats.securityTests.push(result);
        }
    }

    /**
     * Tests CSRF
     */
    async testCSRF() {
        // Test sans token CSRF
        const result = await this.testEndpoint({
            method: 'POST',
            path: '/api/admin/delete',
            data: { id: 1 },
            expectedStatus: 403,
            description: 'CSRF protection test'
        });
        
        this.stats.securityTests.push(result);
    }

    /**
     * Tests de rate limiting
     */
    async testRateLimiting() {
        console.log('🔄 Testing rate limiting');
        
        const promises = [];
        for (let i = 0; i < 50; i++) {
            promises.push(this.client.get('/api/test'));
        }

        const results = await Promise.allSettled(promises);
        const rateLimitedCount = results.filter(r => 
            r.status === 'fulfilled' && r.value.status === 429
        ).length;

        const testResult = {
            method: 'GET',
            path: '/api/test',
            description: 'Rate limiting test (50 rapid requests)',
            passed: rateLimitedCount > 0,
            rateLimitedRequests: rateLimitedCount,
            totalRequests: 50,
            errors: rateLimitedCount === 0 ? ['No rate limiting detected'] : []
        };

        this.stats.securityTests.push(testResult);
    }

    /**
     * Tests des headers de sécurité
     */
    async testSecurityHeaders() {
        const response = await this.client.get('/');
        const headers = response.headers;
        
        const requiredHeaders = {
            'x-content-type-options': 'nosniff',
            'x-frame-options': ['DENY', 'SAMEORIGIN'],
            'x-xss-protection': '1; mode=block',
            'strict-transport-security': (value) => value && value.includes('max-age')
        };

        const testResult = {
            method: 'GET',
            path: '/',
            description: 'Security headers test',
            passed: true,
            errors: [],
            headers: headers
        };

        for (const [header, expected] of Object.entries(requiredHeaders)) {
            const actual = headers[header];
            
            if (typeof expected === 'function') {
                if (!expected(actual)) {
                    testResult.errors.push(`Missing or invalid security header: ${header}`);
                    testResult.passed = false;
                }
            } else if (Array.isArray(expected)) {
                if (!expected.includes(actual)) {
                    testResult.errors.push(`Invalid ${header} header: expected one of ${expected.join(', ')}, got ${actual}`);
                    testResult.passed = false;
                }
            } else if (actual !== expected) {
                testResult.errors.push(`Missing security header: ${header}`);
                testResult.passed = false;
            }
        }

        this.stats.securityTests.push(testResult);
    }

    /**
     * Tests de contournement d'authentification
     */
    async testAuthBypass() {
        const bypassAttempts = [
            { path: '/api/admin/../public', description: 'Path traversal auth bypass' },
            { path: '/api/admin%2f..%2fpublic', description: 'URL encoded path traversal' },
            { path: '/api/admin/', description: 'Trailing slash bypass' }
        ];

        for (const attempt of bypassAttempts) {
            const result = await this.testEndpoint({
                method: 'GET',
                path: attempt.path,
                expectedStatus: 401,
                description: attempt.description
            });
            
            this.stats.securityTests.push(result);
        }
    }

    /**
     * Calcule le score global
     */
    calculateOverallScore() {
        let score = 100;
        
        // Score des tests de contrat (40%)
        const contractPassed = this.stats.contractTests.filter(t => t.passed).length;
        const contractTotal = this.stats.contractTests.length;
        const contractScore = contractTotal > 0 ? (contractPassed / contractTotal) * 40 : 0;
        
        // Score des tests de charge (30%)
        const loadPassed = this.stats.loadTests.filter(t => t.passed).length;
        const loadTotal = this.stats.loadTests.length;
        const loadScore = loadTotal > 0 ? (loadPassed / loadTotal) * 30 : 0;
        
        // Score des tests de sécurité (30%)
        const securityPassed = this.stats.securityTests.filter(t => t.passed).length;
        const securityTotal = this.stats.securityTests.length;
        const securityScore = securityTotal > 0 ? (securityPassed / securityTotal) * 30 : 0;
        
        this.stats.overallScore = Math.round(contractScore + loadScore + securityScore);
        
        // Déterminer le statut
        if (this.stats.overallScore >= 90) {
            this.stats.status = 'excellent';
        } else if (this.stats.overallScore >= 80) {
            this.stats.status = 'good';
        } else if (this.stats.overallScore >= 70) {
            this.stats.status = 'acceptable';
        } else {
            this.stats.status = 'needs_improvement';
        }
    }

    /**
     * Génère le rapport de test
     */
    async generateReport() {
        try {
            await fs.mkdir(this.config.reportPath, { recursive: true });
            
            const reportFile = path.join(this.config.reportPath, 'api_test_report.json');
            const htmlReportFile = path.join(this.config.reportPath, 'api_test_report.html');
            
            // Rapport JSON
            await fs.writeFile(reportFile, JSON.stringify(this.stats, null, 2));
            
            // Rapport HTML
            const htmlReport = this.generateHTMLReport();
            await fs.writeFile(htmlReportFile, htmlReport);
            
            console.log(`📊 Report generated: ${reportFile}`);
            
        } catch (error) {
            console.error('Failed to generate report:', error);
        }
    }

    /**
     * Génère un rapport HTML
     */
    generateHTMLReport() {
        return `
<!DOCTYPE html>
<html>
<head>
    <title>Fire Salamander - API Test Report</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { background: #ff6136; color: white; padding: 20px; border-radius: 8px; }
        .score { font-size: 2em; font-weight: bold; }
        .section { margin: 20px 0; }
        .test-result { margin: 10px 0; padding: 10px; border-left: 4px solid #ddd; }
        .passed { border-left-color: #28a745; }
        .failed { border-left-color: #dc3545; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🕷️ Fire Salamander - API Test Report</h1>
        <div class="score">Score: ${this.stats.overallScore}/100 (${this.stats.status})</div>
        <p>Generated: ${this.stats.timestamp}</p>
    </div>
    
    <div class="section">
        <h2>📋 Contract Tests</h2>
        ${this.stats.contractTests.map(test => `
            <div class="test-result ${test.passed ? 'passed' : 'failed'}">
                <strong>${test.method} ${test.path}</strong> - ${test.description}<br>
                Status: ${test.actualStatus} | Response Time: ${test.responseTime}ms
                ${test.errors.length > 0 ? `<br>Errors: ${test.errors.join(', ')}` : ''}
            </div>
        `).join('')}
    </div>
    
    <div class="section">
        <h2>⚡ Load Tests</h2>
        ${this.stats.loadTests.map(test => `
            <div class="test-result ${test.passed ? 'passed' : 'failed'}">
                <strong>${test.name}</strong><br>
                Success: ${test.successCount}/${test.totalRequests} | 
                Avg Response: ${Math.round(test.averageResponseTime)}ms | 
                P95: ${Math.round(test.p95ResponseTime)}ms
            </div>
        `).join('')}
    </div>
    
    <div class="section">
        <h2>🔒 Security Tests</h2>
        ${this.stats.securityTests.map(test => `
            <div class="test-result ${test.passed ? 'passed' : 'failed'}">
                <strong>${test.description}</strong>
                ${test.errors.length > 0 ? `<br>Issues: ${test.errors.join(', ')}` : ''}
            </div>
        `).join('')}
    </div>
</body>
</html>
        `;
    }
}

module.exports = APITestAgent;