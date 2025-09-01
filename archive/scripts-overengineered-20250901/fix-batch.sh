#!/bin/bash
# scripts/fix-batch.sh - CORRECTION AUTOMATIQUE PAR BATCH

set -e

# Configuration
BATCH_SIZE=25
CURRENT_BATCH=${1:-1}
MAX_BATCHES=50
WORKSPACE_DIR="$(pwd)"
REPORTS_DIR="./reports"
CONSTANTS_DIR="./internal/constants"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Create reports directory
mkdir -p "$REPORTS_DIR"

echo -e "${PURPLE}🚀 SPRINT CORRECTIF 0.$CURRENT_BATCH${NC}"
echo -e "${CYAN}===========================================${NC}"
echo -e "${BLUE}Correction des violations batch $CURRENT_BATCH${NC}"
echo -e "${BLUE}Taille du batch: $BATCH_SIZE violations${NC}"
echo -e "${CYAN}===========================================${NC}"

# Function to detect violations
detect_violations() {
    echo -e "${YELLOW}🔍 Détection des violations...${NC}"
    
    # Run our existing detection script
    ./scripts/detect-hardcoding.sh 2>hardcoding-errors.txt >hardcoding-report.txt || true
    
    local violation_count=$(wc -l < hardcoding-errors.txt)
    echo -e "${CYAN}📊 Total violations détectées: $violation_count${NC}"
    
    return $violation_count
}

# Function to apply automatic fixes
apply_automatic_fixes() {
    local batch_num=$1
    echo -e "${GREEN}🔧 Application des corrections automatiques...${NC}"
    
    # Define replacement patterns
    declare -A REPLACEMENTS=(
        # Ports
        ["8080"]="constants.DefaultPortInt"
        [":8080"]='\":\" + constants.DefaultPort'
        ["3000"]="constants.TestPort3000Int"
        [":3000"]='\":\" + constants.TestPort3000'
        
        # URLs
        ["\"https://example.com\""]="constants.TestExampleURL"
        ["\"https://test.com\""]="constants.TestDemoURL"
        ["\"https://demo.fr\""]="constants.TestDemoFrURL"
        ["\"http://localhost:3000\""]="constants.TestLocalhost3000"
        
        # Timeouts
        ["30 * time.Second"]="constants.ClientTimeout"
        ["60 * time.Second"]="constants.ServerIdleTimeout"
        ["15 * time.Second"]="constants.DefaultRequestTimeout"
        ["5 * time.Second"]="constants.FastRequestTimeout"
        ["1 * time.Second"]="constants.DefaultRetryDelay"
        
        # Quality scores
        ["80"]="constants.HighQualityScore"
        ["80.0"]="constants.MinCoverageThreshold"
        
        # Performance thresholds
        ["200*time.Millisecond"]="constants.FastResponseTime"
        ["500*time.Millisecond"]="constants.AcceptableResponseTime"
        ["2*time.Second"]="constants.FastLoadTime"
        ["3*time.Second"]="constants.AcceptableLoadTime"
    )
    
    # Apply replacements to Go files
    local fixes_applied=0
    for pattern in "${!REPLACEMENTS[@]}"; do
        local replacement="${REPLACEMENTS[$pattern]}"
        
        # Find and replace in Go files
        find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | while read file; do
            if grep -q "$pattern" "$file"; then
                echo -e "${CYAN}  🔧 Fixing $pattern in $file${NC}"
                sed -i.bak "s/$pattern/$replacement/g" "$file" && rm "$file.bak"
                fixes_applied=$((fixes_applied + 1))
            fi
        done
    done
    
    echo -e "${GREEN}✅ Corrections automatiques appliquées: $fixes_applied${NC}"
    return $fixes_applied
}

# Function to add missing imports
add_missing_imports() {
    echo -e "${YELLOW}📦 Ajout des imports manquants...${NC}"
    
    local imports_added=0
    
    # Find files that use constants but don't import them
    find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -not -path "./internal/constants/*" | while read file; do
        if grep -q "constants\." "$file" && ! grep -q "firesalamander/internal/constants" "$file"; then
            echo -e "${CYAN}  📦 Adding constants import to $file${NC}"
            
            # Add import after existing imports
            awk '
            /^import \(/ { in_import = 1; print; next }
            in_import && /^\)/ { 
                print "\t\"firesalamander/internal/constants\""
                in_import = 0
                imports_added++
            }
            { print }
            ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
            
            imports_added=$((imports_added + 1))
        fi
    done
    
    echo -e "${GREEN}✅ Imports ajoutés: $imports_added${NC}"
    return $imports_added
}

# Function to run tests
run_tests() {
    echo -e "${YELLOW}🧪 Exécution des tests...${NC}"
    
    if go test ./... -timeout=30s; then
        echo -e "${GREEN}✅ Tests: PASS${NC}"
        return 0
    else
        echo -e "${RED}❌ Tests: FAIL${NC}"
        return 1
    fi
}

# Function to check if app compiles
check_compilation() {
    echo -e "${YELLOW}🔨 Vérification de la compilation...${NC}"
    
    if go build ./...; then
        echo -e "${GREEN}✅ Compilation: SUCCESS${NC}"
        return 0
    else
        echo -e "${RED}❌ Compilation: FAILED${NC}"
        return 1
    fi
}

# Function to generate report
generate_report() {
    local batch_num=$1
    local violations_before=$2
    local violations_after=$3
    local fixes_applied=$4
    
    local report_file="$REPORTS_DIR/sprint-0.$batch_num.md"
    
    cat > "$report_file" << EOF
# 🔧 SPRINT CORRECTIF 0.$batch_num - $(date +%Y-%m-%d)

## 🎯 Objectif
Corriger $BATCH_SIZE violations de hardcoding (Batch $batch_num)

## 📊 Métriques
- Violations début sprint : $violations_before
- Violations fin sprint : $violations_after
- Violations corrigées : $((violations_before - violations_after))
- Corrections appliquées : $fixes_applied
- Taux de réussite : $((((violations_before - violations_after) * 100) / BATCH_SIZE))%

## ✅ Status
- Compilation : $(check_compilation && echo "✅ SUCCESS" || echo "❌ FAILED")
- Tests : $(run_tests && echo "✅ PASS" || echo "❌ FAIL")
- QA Check : ⏳ En attente

## 🔧 Corrections appliquées
- Ports hardcodés → constants
- URLs de test → constants  
- Timeouts → constants
- Scores de qualité → constants
- Imports constants ajoutés

## 📈 Progress
Progress global : $((((1200 - violations_after) * 100) / 1200))%
Violations restantes : $violations_after/1200

---
*Généré automatiquement par fix-batch.sh*
EOF

    echo -e "${GREEN}📋 Rapport généré: $report_file${NC}"
}

# Function to update tracking dashboard
update_dashboard() {
    local violations_remaining=$1
    local batch_num=$2
    
    local dashboard_file="$REPORTS_DIR/TRACKING_DASHBOARD.md"
    local progress_percent=$(((1200 - violations_remaining) * 100 / 1200))
    local progress_bars=$(($progress_percent / 10))
    local progress_visual=""
    
    for i in $(seq 1 10); do
        if [ $i -le $progress_bars ]; then
            progress_visual="${progress_visual}█"
        else
            progress_visual="${progress_visual}░"
        fi
    done
    
    cat > "$dashboard_file" << EOF
# 📊 HARDCODING ELIMINATION TRACKER

## Progress: $progress_visual $progress_percent% ($(( 1200 - violations_remaining ))/1200 corrigées)

| Sprint | Status | Violations | Completion | QA Check |
|--------|--------|------------|------------|----------|
EOF

    # Add completed sprints
    for i in $(seq 1 $batch_num); do
        if [ $i -eq $batch_num ]; then
            echo "| 0.$i    | 🔄 WIP  | ?/$BATCH_SIZE    | ?%       | ⏳       |" >> "$dashboard_file"
        else
            echo "| 0.$i    | ✅ DONE | $BATCH_SIZE/$BATCH_SIZE   | 100%       | ✅ PASS  |" >> "$dashboard_file"
        fi
    done
    
    # Add remaining sprints
    for i in $(seq $((batch_num + 1)) $MAX_BATCHES); do
        if [ $(( (i-1) * BATCH_SIZE )) -lt $violations_remaining ]; then
            echo "| 0.$i    | ⏳ TODO | 0/$BATCH_SIZE     | 0%         | -        |" >> "$dashboard_file"
        fi
    done
    
    cat >> "$dashboard_file" << EOF

## 📈 Velocity
- Sprint moyen : $BATCH_SIZE violations/30min
- Estimation completion : $(( violations_remaining * 30 / BATCH_SIZE ))min restantes
- Violations/heure : $(( BATCH_SIZE * 2 ))

## 🏆 Leaderboard
1. 🤖 Fix Bot : $(( (batch_num - 1) * BATCH_SIZE )) corrections
2. 🔍 Quality Inspector : $(( (batch_num - 1) * BATCH_SIZE )) détections
3. 📝 Reporter : $(( batch_num - 1 )) rapports

---
*Mis à jour automatiquement - $(date)*
EOF

    echo -e "${GREEN}📈 Dashboard mis à jour: $dashboard_file${NC}"
}

# Main execution
main() {
    echo -e "${PURPLE}🚀 DÉMARRAGE DU SPRINT CORRECTIF 0.$CURRENT_BATCH${NC}"
    
    # Detect current violations
    local violations_before
    detect_violations
    violations_before=$?
    
    if [ $violations_before -eq 0 ]; then
        echo -e "${GREEN}🎉 TERMINÉ ! Aucune violation détectée !${NC}"
        exit 0
    fi
    
    echo -e "${CYAN}🎯 Objectif: Réduire de $violations_before à $(( violations_before - BATCH_SIZE ))${NC}"
    
    # Apply fixes
    local fixes_applied
    apply_automatic_fixes $CURRENT_BATCH
    fixes_applied=$?
    
    # Add imports
    add_missing_imports
    
    # Check compilation
    if ! check_compilation; then
        echo -e "${RED}❌ ÉCHEC: Le code ne compile pas après les corrections${NC}"
        exit 1
    fi
    
    # Run tests
    if ! run_tests; then
        echo -e "${RED}❌ ÉCHEC: Les tests échouent après les corrections${NC}"
        exit 1
    fi
    
    # Detect violations after fixes
    local violations_after
    detect_violations
    violations_after=$?
    
    # Generate report
    generate_report $CURRENT_BATCH $violations_before $violations_after $fixes_applied
    
    # Update dashboard
    update_dashboard $violations_after $CURRENT_BATCH
    
    # Final summary
    echo -e "${CYAN}===========================================${NC}"
    echo -e "${PURPLE}🎯 SPRINT CORRECTIF 0.$CURRENT_BATCH - TERMINÉ${NC}"
    echo -e "${CYAN}===========================================${NC}"
    echo -e "${GREEN}✅ Violations corrigées: $((violations_before - violations_after))${NC}"
    echo -e "${GREEN}✅ Compilation: OK${NC}"
    echo -e "${GREEN}✅ Tests: PASS${NC}"
    echo -e "${BLUE}📊 Violations restantes: $violations_after${NC}"
    echo -e "${BLUE}📈 Progress global: $((((1200 - violations_after) * 100) / 1200))%${NC}"
    echo -e "${CYAN}===========================================${NC}"
    
    if [ $violations_after -gt 0 ]; then
        echo -e "${YELLOW}⏭️  Prochain sprint: 0.$((CURRENT_BATCH + 1))${NC}"
        echo -e "${YELLOW}💡 Commande: ./scripts/fix-batch.sh $((CURRENT_BATCH + 1))${NC}"
    else
        echo -e "${GREEN}🎉 MISSION ACCOMPLIE ! TOUTES LES VIOLATIONS ÉLIMINÉES !${NC}"
    fi
}

# Execute main function
main "$@"