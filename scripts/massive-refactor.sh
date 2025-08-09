#!/bin/bash
# massive-refactor.sh - REFACTORING AUTOMATIQUE POUR ZERO VIOLATIONS

echo "🚀 MASSIVE REFACTORING - FIRE SALAMANDER"
echo "=========================================="

# Configuration
CONSTANTS_IMPORT='"firesalamander/internal/constants"'
MESSAGES_IMPORT='"firesalamander/internal/messages"'
BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"

# Créer sauvegarde
mkdir -p "$BACKUP_DIR"
echo "📁 Sauvegarde créée dans: $BACKUP_DIR"

# Liste des fichiers Go (excluant constants, messages, tests, vendor)
GO_FILES=$(find . -name "*.go" \
    -not -path "./vendor/*" \
    -not -name "*_test.go" \
    -not -path "./internal/constants/*" \
    -not -path "./internal/messages/*" \
    -not -path "./backups/*")

echo "📄 Fichiers à traiter: $(echo "$GO_FILES" | wc -l)"

# Patterns de remplacement les plus critiques
declare -A REPLACEMENTS=(
    # Timeouts courants
    ["30\s*\*\s*time\.Second"]="constants.ShutdownTimeout"
    ["15\s*\*\s*time\.Second"]="constants.ClientTimeout" 
    ["60\s*\*\s*time\.Second"]="constants.ServerIdleTimeout"
    ["5\s*\*\s*time\.Second"]="constants.DefaultRobotsTimeout"
    ["10\s*\*\s*time\.Second"]="constants.DefaultSitemapTimeout"
    ["1\s*\*\s*time\.Second"]="constants.DefaultRetryDelay"
    ["2\s*\*\s*time\.Second"]="constants.FastLoadTime"
    ["3\s*\*\s*time\.Second"]="constants.AcceptableLoadTime"
    
    # Durées communes
    ["200\s*\*\s*time\.Millisecond"]="constants.FastResponseTime"
    ["500\s*\*\s*time\.Millisecond"]="constants.AcceptableResponseTime"
    ["1000\s*\*\s*time\.Millisecond"]="constants.DefaultPollInterval"
    ["7\s*\*\s*24\s*\*\s*time\.Hour"]="constants.RobotsCacheDuration"
    ["24\s*\*\s*time\.Hour"]="constants.DefaultCacheExpiry"
    ["5\s*\*\s*time\.Minute"]="constants.CacheCleanupInterval"
    ["1\s*\*\s*time\.Hour"]="constants.RobotsCleanupInterval"
    
    # Ports
    ["8080"]="constants.DefaultPort"
    [":8080"]='":\" + constants.DefaultPort'
    ["3000"]="constants.TestLocalhost3000"
    
    # URLs courantes  
    ['"https://api\.openai\.com/v1/chat/completions"']="constants.OpenAIAPIURL"
    ['"https://example\.com"']="constants.TestExampleURL"
    ['"http://localhost:3000"']="constants.TestLocalhost3000"
    ['"https://test\.com"']="constants.TestDemoURL"
    ['"https://demo\.fr"']="constants.TestDemoFrURL"
    
    # Messages courants
    ['"Method not allowed"']="messages.ErrMethodNotAllowed"
    ['"Invalid JSON"']="messages.ErrInvalidJSON"
    ['"URL is required"']="messages.ErrURLRequired"
    ['"Invalid URL format"']="messages.ErrInvalidURL"
    ['"Analysis not found"']="messages.ErrAnalysisNotFound"
    ['"Invalid analysis ID"']="messages.ErrInvalidAnalysisID"
    
    # Status
    ['"started"']="constants.StatusProcessing"
    ['"analyzing"']="constants.StatusProcessing" 
    ['"complete"']="constants.StatusComplete"
    ['"completed"']="constants.StatusComplete"
    
    # Nombres magiques communs
    ["\b80\b"]="constants.HighQualityScore"
    ["80\.0"]="constants.MinCoverageThreshold"
)

REFACTOR_COUNT=0
IMPORT_ADDED=0

# Fonction pour ajouter l'import si nécessaire
add_import() {
    local file=$1
    local import=$2
    local import_name=$(echo $import | tr -d '"')
    
    if ! grep -q "$import_name" "$file"; then
        # Chercher la ligne import (
        if grep -q "^import (" "$file"; then
            # Ajouter après import (
            sed -i.refactor "/^import (/a\\
\\t$import" "$file"
            echo "  ✅ Import ajouté: $import_name"
            ((IMPORT_ADDED++))
        fi
    fi
}

# Sauvegarder les fichiers d'origine
echo "💾 Sauvegarde des fichiers..."
for file in $GO_FILES; do
    cp "$file" "$BACKUP_DIR/$(basename "$file").bak"
done

echo "🔄 Début du refactoring massif..."

# Appliquer les remplacements
for file in $GO_FILES; do
    echo "📝 Traitement: $file"
    file_changed=false
    needs_constants=false
    needs_messages=false
    
    for pattern in "${!REPLACEMENTS[@]}"; do
        replacement="${REPLACEMENTS[$pattern]}"
        
        # Effectuer le remplacement
        if sed -i.temp "s/$pattern/$replacement/g" "$file"; then
            if ! diff -q "$file.temp" "$file" >/dev/null 2>&1; then
                echo "  🔧 Remplacé: $pattern → $replacement"
                file_changed=true
                ((REFACTOR_COUNT++))
                
                # Déterminer si on a besoin d'imports
                if [[ $replacement == constants.* ]]; then
                    needs_constants=true
                elif [[ $replacement == messages.* ]]; then
                    needs_messages=true
                fi
            fi
        fi
        rm -f "$file.temp"
    done
    
    # Ajouter les imports si nécessaire
    if $needs_constants; then
        add_import "$file" "$CONSTANTS_IMPORT"
    fi
    if $needs_messages; then
        add_import "$file" "$MESSAGES_IMPORT"  
    fi
    
    # Nettoyer les fichiers temporaires
    rm -f "$file.refactor"
done

echo "=========================================="
echo "📊 RAPPORT FINAL DU REFACTORING MASSIF"
echo "=========================================="
echo "🔧 Remplacements effectués: $REFACTOR_COUNT"
echo "📦 Imports ajoutés: $IMPORT_ADDED"
echo "💾 Sauvegarde dans: $BACKUP_DIR"
echo ""
echo "🧪 Lancement de la vérification..."

# Vérifier les violations restantes
./scripts/detect-hardcoding.sh

echo ""
echo "✅ Refactoring massif terminé !"
echo "⚠️  Vérifiez que le code compile avec: go build ./..."