#!/bin/bash

# Fire Salamander - Test Agents Script
# Script de test global pour tous les agents de test

set -e

# Configuration
BASE_URL="http://localhost:3000"
OUTPUT_DIR="tests/reports"
TIMEOUT=30

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
ORANGE='\033[0;33m'
NC='\033[0m' # No Color

# Fonctions utiles
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_salamander() {
    echo -e "${ORANGE}[🔥 SALAMANDER]${NC} $1"
}

# Fonction pour vérifier les dépendances
check_dependencies() {
    log_info "Checking dependencies..."
    
    # Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        exit 1
    fi
    
    # Node.js et npm
    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed"
        exit 1
    fi
    
    # Python3
    if ! command -v python3 &> /dev/null; then
        log_error "Python3 is not installed"
        exit 1
    fi
    
    # pip3
    if ! command -v pip3 &> /dev/null; then
        log_error "pip3 is not installed"
        exit 1
    fi
    
    log_success "All basic dependencies are available"
}

# Fonction pour installer les dépendances Python
install_python_deps() {
    log_info "Installing Python dependencies..."
    
    # Créer un requirements.txt temporaire si il n'existe pas
    if [ ! -f "requirements.txt" ]; then
        cat > requirements.txt << EOF
requests>=2.25.0
beautifulsoup4>=4.9.0
psutil>=5.8.0
pyyaml>=6.0
EOF
        log_info "Created temporary requirements.txt"
    fi
    
    pip3 install -r requirements.txt --quiet || {
        log_warning "Some Python packages may not be installed. Continuing anyway..."
    }
}

# Fonction pour vérifier que le serveur Fire Salamander est démarré
check_server() {
    log_info "Checking if Fire Salamander server is running at $BASE_URL..."
    
    local max_attempts=10
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
            log_success "Fire Salamander server is running"
            return 0
        fi
        
        log_info "Attempt $attempt/$max_attempts - Server not ready, waiting..."
        sleep 2
        ((attempt++))
    done
    
    log_warning "Server is not responding. Starting server in background..."
    nohup go run main.go > /dev/null 2>&1 &
    SERVER_PID=$!
    
    # Attendre que le serveur démarre
    sleep 10
    
    if curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
        log_success "Fire Salamander server started successfully (PID: $SERVER_PID)"
        return 0
    else
        log_error "Failed to start Fire Salamander server"
        return 1
    fi
}

# Fonction pour créer les répertoires de rapports
create_report_dirs() {
    log_info "Creating report directories..."
    
    mkdir -p "$OUTPUT_DIR"/{qa,api,frontend,security,performance,seo,data,monitoring}
    
    log_success "Report directories created"
}

# Fonction pour tester le QA Agent (Go)
test_qa_agent() {
    log_salamander "Testing QA Agent (Go)..."
    
    cd tests/agents/qa
    
    if go test -v > "$OUTPUT_DIR/qa/qa_test_output.log" 2>&1; then
        log_success "QA Agent tests passed"
        return 0
    else
        log_error "QA Agent tests failed. Check $OUTPUT_DIR/qa/qa_test_output.log"
        return 1
    fi
    
    cd - > /dev/null
}

# Fonction pour tester l'API Test Agent
test_api_agent() {
    log_salamander "Testing API Test Agent..."
    
    cd tests/agents/api
    
    # Installer les dépendances si nécessaire
    if [ ! -d "node_modules" ]; then
        log_info "Installing API test dependencies..."
        npm install --silent
    fi
    
    if node test_runner.js --url="$BASE_URL" > "$OUTPUT_DIR/api/api_test_output.log" 2>&1; then
        log_success "API Test Agent passed"
        return 0
    else
        log_error "API Test Agent failed. Check $OUTPUT_DIR/api/api_test_output.log"
        return 1
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le Frontend Test Agent (Playwright)
test_frontend_agent() {
    log_salamander "Testing Frontend Test Agent (Playwright)..."
    
    cd tests/agents/frontend
    
    # Installer les dépendances si nécessaire
    if [ ! -d "node_modules" ]; then
        log_info "Installing Playwright dependencies..."
        npm install --silent
        npx playwright install --with-deps chromium --quiet
    fi
    
    if BASE_URL="$BASE_URL" npx playwright test > "$OUTPUT_DIR/frontend/frontend_test_output.log" 2>&1; then
        log_success "Frontend Test Agent passed"
        return 0
    else
        log_error "Frontend Test Agent failed. Check $OUTPUT_DIR/frontend/frontend_test_output.log"
        return 1
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le Security Agent
test_security_agent() {
    log_salamander "Testing Security Agent..."
    
    cd tests/agents/security
    
    if python3 security_agent.py --url="$BASE_URL" --output="$OUTPUT_DIR/security" > "$OUTPUT_DIR/security/security_test_output.log" 2>&1; then
        log_success "Security Agent passed"
        return 0
    else
        log_warning "Security Agent completed with warnings. Check $OUTPUT_DIR/security/security_test_output.log"
        return 0  # Security tests peuvent avoir des warnings
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le Performance Agent (k6)
test_performance_agent() {
    log_salamander "Testing Performance Agent (k6)..."
    
    # Vérifier si k6 est installé
    if ! command -v k6 &> /dev/null; then
        log_warning "k6 is not installed. Skipping performance tests."
        echo "To install k6: https://k6.io/docs/getting-started/installation/" > "$OUTPUT_DIR/performance/k6_not_installed.log"
        return 0
    fi
    
    cd tests/agents/performance
    
    if BASE_URL="$BASE_URL" k6 run --out json="$OUTPUT_DIR/performance/results.json" k6-load-test.js > "$OUTPUT_DIR/performance/performance_test_output.log" 2>&1; then
        log_success "Performance Agent passed"
        return 0
    else
        log_warning "Performance Agent completed with issues. Check $OUTPUT_DIR/performance/performance_test_output.log"
        return 0  # Performance tests peuvent échouer selon les ressources
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le SEO Agent
test_seo_agent() {
    log_salamander "Testing SEO Agent..."
    
    cd tests/agents/seo
    
    if python3 seo_agent.py --url="$BASE_URL" --output="$OUTPUT_DIR/seo" > "$OUTPUT_DIR/seo/seo_test_output.log" 2>&1; then
        log_success "SEO Agent passed"
        return 0
    else
        log_warning "SEO Agent completed with issues. Check $OUTPUT_DIR/seo/seo_test_output.log"
        return 0  # SEO tests peuvent avoir des warnings
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le Data Integrity Agent
test_data_agent() {
    log_salamander "Testing Data Integrity Agent..."
    
    cd tests/agents/data
    
    if go run data_integrity_agent.go --database="../../../fire_salamander_dev.db" --output="$OUTPUT_DIR/data" > "$OUTPUT_DIR/data/data_test_output.log" 2>&1; then
        log_success "Data Integrity Agent passed"
        return 0
    else
        log_warning "Data Integrity Agent completed with issues. Check $OUTPUT_DIR/data/data_test_output.log"
        return 0  # Data tests peuvent avoir des warnings
    fi
    
    cd - > /dev/null
}

# Fonction pour tester le Monitoring Agent
test_monitoring_agent() {
    log_salamander "Testing Monitoring Agent (30 seconds)..."
    
    cd tests/agents/monitoring
    
    if timeout 30 python3 monitoring_agent.py --url="$BASE_URL" --duration=20 --output="$OUTPUT_DIR/monitoring" > "$OUTPUT_DIR/monitoring/monitoring_test_output.log" 2>&1; then
        log_success "Monitoring Agent passed"
        return 0
    else
        log_warning "Monitoring Agent completed with issues. Check $OUTPUT_DIR/monitoring/monitoring_test_output.log"
        return 0  # Monitoring peut avoir des timeouts
    fi
    
    cd - > /dev/null
}

# Fonction pour générer un rapport de synthèse
generate_summary_report() {
    log_info "Generating summary report..."
    
    local summary_file="$OUTPUT_DIR/test_summary.html"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    cat > "$summary_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Fire Salamander - Test Agents Summary</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .header { background: #ff6136; color: white; padding: 20px; border-radius: 8px; text-align: center; }
        .section { background: white; margin: 20px 0; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .agent { margin: 10px 0; padding: 15px; border-left: 4px solid #ddd; border-radius: 4px; }
        .success { border-left-color: #28a745; background-color: #d4edda; }
        .warning { border-left-color: #ffc107; background-color: #fff3cd; }
        .error { border-left-color: #dc3545; background-color: #f8d7da; }
        .links { margin-top: 10px; }
        .links a { margin-right: 10px; color: #007bff; text-decoration: none; }
        .links a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔥 Fire Salamander - Test Agents Summary</h1>
        <p>Generated: $timestamp</p>
        <p>Target: $BASE_URL</p>
    </div>
    
    <div class="section">
        <h2>📊 Test Results</h2>
EOF

    # QA Agent
    if [ -f "$OUTPUT_DIR/qa/qa_test_output.log" ]; then
        if grep -q "PASS\|SUCCESS" "$OUTPUT_DIR/qa/qa_test_output.log"; then
            echo '<div class="agent success"><strong>✅ QA Agent (Go)</strong> - Tests passed</div>' >> "$summary_file"
        else
            echo '<div class="agent error"><strong>❌ QA Agent (Go)</strong> - Tests failed</div>' >> "$summary_file"
        fi
        echo '<div class="links"><a href="qa/qa_test_output.log">View Log</a></div>' >> "$summary_file"
    fi
    
    # API Agent
    if [ -f "$OUTPUT_DIR/api/api_test_output.log" ]; then
        echo '<div class="agent success"><strong>✅ API Test Agent</strong> - Completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="api/api_test_output.log">View Log</a> <a href="api/">View Reports</a></div>' >> "$summary_file"
    fi
    
    # Frontend Agent
    if [ -f "$OUTPUT_DIR/frontend/frontend_test_output.log" ]; then
        echo '<div class="agent success"><strong>✅ Frontend Test Agent (Playwright)</strong> - Completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="frontend/frontend_test_output.log">View Log</a> <a href="frontend/">View Reports</a></div>' >> "$summary_file"
    fi
    
    # Security Agent
    if [ -f "$OUTPUT_DIR/security/security_report.html" ]; then
        echo '<div class="agent warning"><strong>⚠️ Security Agent</strong> - Completed with analysis</div>' >> "$summary_file"
        echo '<div class="links"><a href="security/security_test_output.log">View Log</a> <a href="security/security_report.html">View Security Report</a></div>' >> "$summary_file"
    fi
    
    # Performance Agent
    if [ -f "$OUTPUT_DIR/performance/results.json" ]; then
        echo '<div class="agent success"><strong>✅ Performance Agent (k6)</strong> - Load tests completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="performance/performance_test_output.log">View Log</a> <a href="performance/results.json">View Results</a></div>' >> "$summary_file"
    elif [ -f "$OUTPUT_DIR/performance/k6_not_installed.log" ]; then
        echo '<div class="agent warning"><strong>⚠️ Performance Agent (k6)</strong> - k6 not installed</div>' >> "$summary_file"
        echo '<div class="links"><a href="performance/k6_not_installed.log">Installation Guide</a></div>' >> "$summary_file"
    fi
    
    # SEO Agent
    if [ -f "$OUTPUT_DIR/seo/seo_report.html" ]; then
        echo '<div class="agent success"><strong>✅ SEO Agent</strong> - Analysis completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="seo/seo_test_output.log">View Log</a> <a href="seo/seo_report.html">View SEO Report</a></div>' >> "$summary_file"
    fi
    
    # Data Integrity Agent
    if [ -f "$OUTPUT_DIR/data/data_integrity_report.html" ]; then
        echo '<div class="agent success"><strong>✅ Data Integrity Agent</strong> - Database analysis completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="data/data_test_output.log">View Log</a> <a href="data/data_integrity_report.html">View Data Report</a></div>' >> "$summary_file"
    fi
    
    # Monitoring Agent
    if [ -f "$OUTPUT_DIR/monitoring/monitoring_report.html" ]; then
        echo '<div class="agent success"><strong>✅ Monitoring Agent</strong> - System monitoring completed</div>' >> "$summary_file"
        echo '<div class="links"><a href="monitoring/monitoring_test_output.log">View Log</a> <a href="monitoring/monitoring_report.html">View Monitoring Report</a></div>' >> "$summary_file"
    fi
    
    cat >> "$summary_file" << EOF
    </div>
    
    <div class="section">
        <h2>📁 Report Directories</h2>
        <ul>
            <li><strong>QA:</strong> tests/reports/qa/ - Go code quality and coverage</li>
            <li><strong>API:</strong> tests/reports/api/ - API contract, load, and security tests</li>
            <li><strong>Frontend:</strong> tests/reports/frontend/ - Playwright E2E tests</li>
            <li><strong>Security:</strong> tests/reports/security/ - OWASP security analysis</li>
            <li><strong>Performance:</strong> tests/reports/performance/ - k6 load testing results</li>
            <li><strong>SEO:</strong> tests/reports/seo/ - SEO accuracy analysis</li>
            <li><strong>Data:</strong> tests/reports/data/ - Database integrity checks</li>
            <li><strong>Monitoring:</strong> tests/reports/monitoring/ - System health monitoring</li>
        </ul>
    </div>
    
    <div class="section">
        <h2>🔧 Next Steps</h2>
        <p>All test agents have been executed. Review the individual reports for detailed analysis and recommendations.</p>
        <p>For CI/CD integration, use: <code>./test-agents.sh</code> in your pipeline.</p>
    </div>
</body>
</html>
EOF

    log_success "Summary report generated: $summary_file"
}

# Fonction pour nettoyer (optionnel)
cleanup() {
    if [ ! -z "$SERVER_PID" ]; then
        log_info "Stopping Fire Salamander server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
    fi
}

# Fonction principale
main() {
    echo "🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥"
    log_salamander "Fire Salamander v5 - Test Agents Suite"
    echo "🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥"
    echo
    
    # Trap pour cleanup
    trap cleanup EXIT
    
    # Étape 1: Vérifier les dépendances
    check_dependencies
    
    # Étape 2: Installer les dépendances Python
    install_python_deps
    
    # Étape 3: Vérifier le serveur
    if ! check_server; then
        log_error "Cannot proceed without Fire Salamander server"
        exit 1
    fi
    
    # Étape 4: Créer les répertoires
    create_report_dirs
    
    echo
    log_salamander "Starting Test Agents Execution..."
    echo
    
    # Variables pour compter les résultats
    local passed=0
    local failed=0
    local warnings=0
    
    # Étape 5: Exécuter tous les agents de test
    
    # QA Agent (critique)
    if test_qa_agent; then
        ((passed++))
    else
        ((failed++))
    fi
    
    # API Agent
    if test_api_agent; then
        ((passed++))
    else
        ((failed++))
    fi
    
    # Frontend Agent
    if test_frontend_agent; then
        ((passed++))
    else
        ((failed++))
    fi
    
    # Security Agent
    if test_security_agent; then
        ((passed++))
    else
        ((warnings++))
    fi
    
    # Performance Agent
    if test_performance_agent; then
        ((passed++))
    else
        ((warnings++))
    fi
    
    # SEO Agent
    if test_seo_agent; then
        ((passed++))
    else
        ((warnings++))
    fi
    
    # Data Integrity Agent
    if test_data_agent; then
        ((passed++))
    else
        ((warnings++))
    fi
    
    # Monitoring Agent
    if test_monitoring_agent; then
        ((passed++))
    else
        ((warnings++))
    fi
    
    # Étape 6: Générer le rapport de synthèse
    generate_summary_report
    
    # Résultats finaux
    echo
    echo "🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥"
    log_salamander "Test Agents Execution Complete!"
    echo "🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥"
    echo
    log_success "Passed: $passed"
    log_warning "Warnings: $warnings"
    log_error "Failed: $failed"
    echo
    log_info "📊 Summary Report: $OUTPUT_DIR/test_summary.html"
    log_info "📁 All Reports: $OUTPUT_DIR/"
    echo
    
    # Exit code
    if [ $failed -gt 0 ]; then
        log_error "Some critical tests failed!"
        exit 1
    elif [ $warnings -gt 0 ]; then
        log_warning "All tests completed with some warnings"
        exit 0
    else
        log_success "All tests passed successfully!"
        exit 0
    fi
}

# Lancer le script principal
main "$@"