#!/bin/bash

# 🔥 Fire Salamander - Lanceur Suite UX/UI
# Script pour lancer facilement tous les tests UX avec standards SEPTEO

set -e

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FIRE_SALAMANDER_URL="http://localhost:8080"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPORTS_DIR="$SCRIPT_DIR/reports"

# Fonctions utilitaires
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${BLUE}"
    echo "🔥 FIRE SALAMANDER - SUITE UX/UI SEPTEO"
    echo "======================================"
    echo -e "${NC}"
}

# Vérifier les prérequis
check_prerequisites() {
    log_info "Vérification des prérequis..."
    
    # Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js n'est pas installé"
        exit 1
    fi
    
    local node_version=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$node_version" -lt 16 ]; then
        log_error "Node.js version 16+ requis (actuel: $(node --version))"
        exit 1
    fi
    
    # NPM packages
    if [ ! -d "$SCRIPT_DIR/node_modules" ]; then
        log_warn "Dépendances NPM manquantes, installation..."
        cd "$SCRIPT_DIR"
        npm install
    fi
    
    # Fire Salamander
    if ! curl -sf "$FIRE_SALAMANDER_URL/api/v1/health" > /dev/null; then
        log_error "Fire Salamander n'est pas accessible sur $FIRE_SALAMANDER_URL"
        log_info "Démarrez Fire Salamander avec: ./fire-salamander --config config.yaml"
        exit 1
    fi
    
    log_success "Prérequis validés"
}

# Créer les dossiers de rapports
setup_reports() {
    log_info "Préparation des dossiers de rapports..."
    
    mkdir -p "$REPORTS_DIR"/{accessibility,lighthouse,design-system,playwright,consolidated}
    mkdir -p "$SCRIPT_DIR/ux/user-flows/recordings"
    mkdir -p "$SCRIPT_DIR/ux/visual-regression/backstop_data"
    
    log_success "Dossiers créés"
}

# Tests visuels de régression  
run_visual_tests() {
    if [ "$SKIP_VISUAL" = "true" ]; then
        log_warn "Tests visuels ignorés (SKIP_VISUAL=true)"
        return 0
    fi
    
    log_info "🎨 Lancement des tests de régression visuelle..."
    
    cd "$SCRIPT_DIR"
    
    # Créer les références si elles n'existent pas
    if [ ! -d "ux/visual-regression/backstop_data/bitmaps_reference" ]; then
        log_info "Création des captures de référence..."
        npm run test:visual:reference || {
            log_warn "Échec création références, continuation..."
        }
    fi
    
    # Lancer les tests visuels
    if npm run test:visual; then
        log_success "Tests visuels réussis"
        return 0
    else
        log_error "Tests visuels échoués"
        return 1
    fi
}

# Tests d'accessibilité
run_accessibility_tests() {
    if [ "$SKIP_ACCESSIBILITY" = "true" ]; then
        log_warn "Tests d'accessibilité ignorés (SKIP_ACCESSIBILITY=true)"
        return 0
    fi
    
    log_info "♿ Lancement des tests d'accessibilité SEPTEO..."
    
    cd "$SCRIPT_DIR"
    
    local exit_code=0
    
    # Tests Axe-core
    if npm run test:axe; then
        log_success "Tests Axe-core réussis"
    else
        log_error "Tests Axe-core échoués"
        exit_code=1
    fi
    
    # Tests Pa11y
    if npm run test:pa11y; then
        log_success "Tests Pa11y réussis"
    else
        log_error "Tests Pa11y échoués"
        exit_code=1
    fi
    
    return $exit_code
}

# Tests de performance Lighthouse
run_performance_tests() {
    if [ "$SKIP_PERFORMANCE" = "true" ]; then
        log_warn "Tests de performance ignorés (SKIP_PERFORMANCE=true)"
        return 0
    fi
    
    log_info "⚡ Lancement des tests de performance Lighthouse..."
    
    cd "$SCRIPT_DIR"
    
    if npm run test:lighthouse; then
        log_success "Tests Lighthouse réussis"
        return 0
    else
        log_error "Tests Lighthouse échoués"
        return 1
    fi
}

# Tests du Design System SEPTEO
run_design_system_tests() {
    if [ "$SKIP_DESIGN_SYSTEM" = "true" ]; then
        log_warn "Tests Design System ignorés (SKIP_DESIGN_SYSTEM=true)"
        return 0
    fi
    
    log_info "🎯 Validation du Design System SEPTEO..."
    
    cd "$SCRIPT_DIR"
    
    if npm run test:design-system; then
        log_success "Design System SEPTEO validé"
        return 0
    else
        log_error "Design System SEPTEO non conforme"
        return 1
    fi
}

# Tests des parcours utilisateur
run_user_flow_tests() {
    if [ "$SKIP_USER_FLOWS" = "true" ]; then
        log_warn "Tests parcours utilisateur ignorés (SKIP_USER_FLOWS=true)"
        return 0
    fi
    
    log_info "👤 Tests des parcours utilisateur critiques..."
    
    cd "$SCRIPT_DIR"
    
    # Installer les navigateurs Playwright si nécessaire
    if [ ! -d "$HOME/.cache/ms-playwright" ]; then
        log_info "Installation des navigateurs Playwright..."
        npx playwright install
    fi
    
    if npm run test:user-flows; then
        log_success "Parcours utilisateur validés"
        return 0
    else
        log_error "Échec parcours utilisateur"
        return 1
    fi
}

# Génération du rapport consolidé
generate_consolidated_report() {
    log_info "📊 Génération du rapport UX consolidé..."
    
    cd "$SCRIPT_DIR"
    
    if node scripts/generate-ux-report.js; then
        log_success "Rapport consolidé généré"
        
        # Afficher le chemin du rapport
        local report_date=$(date +%Y-%m-%d)
        local html_report="$REPORTS_DIR/consolidated/weekly-ux-report-$report_date.html"
        local exec_report="$REPORTS_DIR/consolidated/executive-ux-summary-$report_date.html"
        
        if [ -f "$html_report" ]; then
            log_info "📄 Rapport détaillé: $html_report"
        fi
        
        if [ -f "$exec_report" ]; then
            log_info "📄 Résumé exécutif: $exec_report"
        fi
        
        return 0
    else
        log_error "Échec génération rapport"
        return 1
    fi
}

# Afficher le résumé des résultats
show_summary() {
    log_info "📈 Résumé des tests UX/UI..."
    
    local json_report="$REPORTS_DIR/ux-consolidated-report.json"
    
    if [ -f "$json_report" ] && command -v jq &> /dev/null; then
        local overall_score=$(jq -r '.septeoCompliance.overall // 0' "$json_report")
        local total_tests=$(jq -r '.summary.totalTests // 0' "$json_report")
        local passed_tests=$(jq -r '.summary.passedTests // 0' "$json_report")
        local recommendations=$(jq -r '.recommendations | length' "$json_report")
        
        echo -e "${BLUE}===========================================${NC}"
        echo -e "${BLUE}🏆 RÉSULTATS FINAUX SEPTEO${NC}"
        echo -e "${BLUE}===========================================${NC}"
        echo -e "📊 Score Global SEPTEO: ${overall_score}%"
        echo -e "🧪 Tests Exécutés: ${passed_tests}/${total_tests}"
        echo -e "💡 Recommandations: ${recommendations}"
        
        if [ "$overall_score" -ge 95 ]; then
            echo -e "${GREEN}🎉 EXCELLENT - Standards SEPTEO respectés!${NC}"
        elif [ "$overall_score" -ge 85 ]; then
            echo -e "${YELLOW}⚠️  BIEN - Quelques améliorations recommandées${NC}"
        else
            echo -e "${RED}❌ À AMÉLIORER - Action requise pour conformité SEPTEO${NC}"
        fi
        
        echo -e "${BLUE}===========================================${NC}"
        
        return 0
    else
        log_warn "Impossible d'afficher le résumé (jq manquant ou rapport absent)"
        return 1
    fi
}

# Lancer le dashboard UX
start_dashboard() {
    log_info "🚀 Démarrage du dashboard UX temps réel..."
    
    cd "$SCRIPT_DIR"
    
    # Vérifier si le dashboard est déjà en cours
    if curl -sf "http://localhost:9002" > /dev/null; then
        log_warn "Dashboard UX déjà en cours sur http://localhost:9002"
        return 0
    fi
    
    # Démarrer le dashboard en arrière-plan
    nohup node scripts/ux-dashboard.js > ux-dashboard.log 2>&1 &
    local dashboard_pid=$!
    
    # Attendre que le dashboard démarre
    sleep 3
    
    if curl -sf "http://localhost:9002" > /dev/null; then
        log_success "Dashboard UX démarré: http://localhost:9002"
        echo "PID: $dashboard_pid"
        echo "$dashboard_pid" > /tmp/fire-salamander-ux-dashboard.pid
        return 0
    else
        log_error "Échec démarrage dashboard UX"
        return 1
    fi
}

# Fonction d'aide
show_help() {
    echo "🔥 Fire Salamander - Suite UX/UI SEPTEO"
    echo ""
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo ""
    echo "COMMANDS:"
    echo "  all                 Lancer tous les tests UX (défaut)"
    echo "  visual              Tests de régression visuelle uniquement"
    echo "  accessibility       Tests d'accessibilité uniquement"
    echo "  performance         Tests de performance uniquement"
    echo "  design-system       Validation Design System SEPTEO uniquement"
    echo "  user-flows          Tests parcours utilisateur uniquement"
    echo "  report              Générer le rapport consolidé uniquement"
    echo "  dashboard           Démarrer le dashboard UX temps réel"
    echo "  help                Afficher cette aide"
    echo ""
    echo "OPTIONS:"
    echo "  --skip-visual       Ignorer les tests visuels"
    echo "  --skip-accessibility Ignorer les tests d'accessibilité"
    echo "  --skip-performance  Ignorer les tests de performance"
    echo "  --skip-design-system Ignorer la validation Design System"
    echo "  --skip-user-flows   Ignorer les tests parcours utilisateur"
    echo "  --fast              Mode rapide (ignorer tests longs)"
    echo "  --ci                Mode CI (sans interaction)"
    echo ""
    echo "EXEMPLES:"
    echo "  $0                  # Tous les tests"
    echo "  $0 --fast           # Tests rapides seulement"
    echo "  $0 accessibility    # Tests accessibilité seulement"
    echo "  $0 dashboard        # Dashboard temps réel"
    echo ""
    echo "VARIABLES D'ENVIRONNEMENT:"
    echo "  FIRE_SALAMANDER_URL URL de Fire Salamander (défaut: http://localhost:8080)"
    echo "  SKIP_VISUAL         Ignorer tests visuels (true/false)"
    echo "  SKIP_ACCESSIBILITY  Ignorer tests accessibilité (true/false)"
    echo "  SKIP_PERFORMANCE    Ignorer tests performance (true/false)"
    echo "  SKIP_USER_FLOWS     Ignorer tests parcours (true/false)"
}

# Fonction principale
main() {
    print_header
    
    # Parser les arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-visual)
                SKIP_VISUAL="true"
                shift
                ;;
            --skip-accessibility)
                SKIP_ACCESSIBILITY="true"
                shift
                ;;
            --skip-performance)
                SKIP_PERFORMANCE="true"
                shift
                ;;
            --skip-design-system)
                SKIP_DESIGN_SYSTEM="true"
                shift
                ;;
            --skip-user-flows)
                SKIP_USER_FLOWS="true"
                shift
                ;;
            --fast)
                SKIP_VISUAL="true"
                SKIP_USER_FLOWS="true"
                shift
                ;;
            --ci)
                CI_MODE="true"
                shift
                ;;
            visual|accessibility|performance|design-system|user-flows|report|dashboard|help)
                COMMAND="$1"
                shift
                ;;
            all)
                COMMAND="all"
                shift
                ;;
            *)
                log_error "Option inconnue: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Commande par défaut
    if [ -z "$COMMAND" ]; then
        COMMAND="all"
    fi
    
    # Exécuter la commande
    case $COMMAND in
        help)
            show_help
            exit 0
            ;;
        dashboard)
            check_prerequisites
            start_dashboard
            exit $?
            ;;
        visual)
            check_prerequisites
            setup_reports
            run_visual_tests
            exit $?
            ;;
        accessibility)
            check_prerequisites
            setup_reports
            run_accessibility_tests
            exit $?
            ;;
        performance)
            check_prerequisites
            setup_reports
            run_performance_tests
            exit $?
            ;;
        design-system)
            check_prerequisites
            setup_reports
            run_design_system_tests
            exit $?
            ;;
        user-flows)
            check_prerequisites
            setup_reports
            run_user_flow_tests
            exit $?
            ;;
        report)
            generate_consolidated_report
            show_summary
            exit $?
            ;;
        all)
            # Suite complète
            check_prerequisites
            setup_reports
            
            local overall_exit_code=0
            
            # Lancer tous les tests
            run_visual_tests || overall_exit_code=1
            run_accessibility_tests || overall_exit_code=1
            run_performance_tests || overall_exit_code=1
            run_design_system_tests || overall_exit_code=1
            run_user_flow_tests || overall_exit_code=1
            
            # Générer le rapport consolidé
            generate_consolidated_report || overall_exit_code=1
            
            # Afficher le résumé
            show_summary
            
            if [ $overall_exit_code -eq 0 ]; then
                log_success "🎉 Suite UX/UI terminée avec succès!"
            else
                log_error "💥 Certains tests ont échoué"
            fi
            
            exit $overall_exit_code
            ;;
        *)
            log_error "Commande inconnue: $COMMAND"
            show_help
            exit 1
            ;;
    esac
}

# Point d'entrée
main "$@"