#!/bin/bash
# scripts/detect-hardcoding.sh - CODE QUALITY INSPECTOR
# ZERO TOLERANCE pour le hardcoding

set -e

echo "🔍 CODE QUALITY INSPECTOR - Analyse en cours..."
echo "================================================="

VIOLATIONS=0
REPORT_FILE="hardcoding-report.txt"
ERROR_FILE="hardcoding-errors.txt"

# Nettoyer les anciens rapports
rm -f "$REPORT_FILE" "$ERROR_FILE"

echo "📁 Analyse des fichiers Go..."

# 1. golangci-lint - Configuration stricte
echo "🔧 golangci-lint: Analyse statique..."
if ! golangci-lint run --config .golangci.yml --out-format=tab > "$REPORT_FILE" 2>&1; then
    LINT_VIOLATIONS=$(grep -c ":" "$REPORT_FILE" 2>/dev/null || echo "0")
    VIOLATIONS=$((VIOLATIONS + LINT_VIOLATIONS))
    echo "❌ golangci-lint: $LINT_VIOLATIONS violations trouvées"
else
    echo "✅ golangci-lint: Aucune violation"
fi

# 2. Revive - Règles strictes
echo "🔧 revive: Analyse des règles métier..."
if ! revive -config revive.toml ./... > "revive-report.txt" 2>&1; then
    REVIVE_VIOLATIONS=$(grep -c ":" "revive-report.txt" 2>/dev/null || echo "0")
    VIOLATIONS=$((VIOLATIONS + REVIVE_VIOLATIONS))
    echo "❌ revive: $REVIVE_VIOLATIONS violations trouvées"
    cat "revive-report.txt" >> "$ERROR_FILE"
else
    echo "✅ revive: Aucune violation"
fi

# 3. Staticcheck - Analyse avancée
echo "🔧 staticcheck: Analyse statique avancée..."
if ! staticcheck ./... > "staticcheck-report.txt" 2>&1; then
    STATIC_VIOLATIONS=$(grep -c ":" "staticcheck-report.txt" 2>/dev/null || echo "0")
    VIOLATIONS=$((VIOLATIONS + STATIC_VIOLATIONS))
    echo "❌ staticcheck: $STATIC_VIOLATIONS violations trouvées"
    cat "staticcheck-report.txt" >> "$ERROR_FILE"
else
    echo "✅ staticcheck: Aucune violation"
fi

echo "================================================="

# 4. Détection manuelle de patterns critiques
echo "🔍 Détection de patterns hardcodés..."

# Strings hardcodées (> 5 caractères, pas de test, pas de constants/messages)
HARDCODED_STRINGS=$(find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -not -path "./internal/messages/*" -exec grep -Hn '"[A-Za-z][A-Za-z0-9 ]\{4,\}"' {} \; | grep -v "const\|var\|//\|fmt\." | wc -l || echo "0")
if [ "$HARDCODED_STRINGS" -gt 0 ]; then
    echo "❌ Chaînes hardcodées: $HARDCODED_STRINGS trouvées"
    find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -not -path "./internal/messages/*" -exec grep -Hn '"[A-Za-z][A-Za-z0-9 ]\{4,\}"' {} \; | grep -v "const\|var\|//\|fmt\." | head -10 >> "$ERROR_FILE"
    VIOLATIONS=$((VIOLATIONS + HARDCODED_STRINGS))
fi

# URLs hardcodées
HARDCODED_URLS=$(find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -not -path "./internal/messages/*" -exec grep -Hn 'http\(s\)\?://[^"]*' {} \; | wc -l || echo "0")
if [ "$HARDCODED_URLS" -gt 0 ]; then
    echo "❌ URLs hardcodées: $HARDCODED_URLS trouvées"
    find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -not -path "./internal/messages/*" -exec grep -Hn 'http\(s\)\?://[^"]*' {} \; >> "$ERROR_FILE"
    VIOLATIONS=$((VIOLATIONS + HARDCODED_URLS))
fi

# Ports et nombres magiques critiques
MAGIC_PORTS=$(find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -exec grep -Hn '\b\(8080\|3000\|5432\|6379\|27017\|443\|80\)\b' {} \; | wc -l || echo "0")
if [ "$MAGIC_PORTS" -gt 0 ]; then
    echo "❌ Ports hardcodés: $MAGIC_PORTS trouvés"
    find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -exec grep -Hn '\b\(8080\|3000\|5432\|6379\|27017\|443\|80\)\b' {} \; >> "$ERROR_FILE"
    VIOLATIONS=$((VIOLATIONS + MAGIC_PORTS))
fi

# Durées hardcodées
HARDCODED_DURATIONS=$(find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -exec grep -Hn '[0-9]\+\s*\*\s*time\.\(Second\|Minute\|Hour\|Millisecond\)' {} \; | wc -l || echo "0")
if [ "$HARDCODED_DURATIONS" -gt 0 ]; then
    echo "❌ Durées hardcodées: $HARDCODED_DURATIONS trouvées"
    find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" -not -path "./internal/constants/*" -exec grep -Hn '[0-9]\+\s*\*\s*time\.\(Second\|Minute\|Hour\|Millisecond\)' {} \; >> "$ERROR_FILE"
    VIOLATIONS=$((VIOLATIONS + HARDCODED_DURATIONS))
fi

# 5. Vérification des templates HTML
echo "🌐 Vérification des templates HTML..."

TEMPLATE_VIOLATIONS=0
if [ -d "templates" ]; then
    # URLs CDN hardcodées
    CDN_URLS=$(find templates -name "*.html" -exec grep -Hn 'https://[^"]*cdn[^"]*' {} \; | wc -l || echo "0")
    if [ "$CDN_URLS" -gt 0 ]; then
        echo "❌ URLs CDN hardcodées: $CDN_URLS trouvées"
        TEMPLATE_VIOLATIONS=$((TEMPLATE_VIOLATIONS + CDN_URLS))
    fi

    # Timeouts JS hardcodés
    JS_TIMEOUTS=$(find templates -name "*.html" -exec grep -Hn '[0-9]\+\s*)\s*;\s*//.*\(timeout\|delay\|interval\)' {} \; | wc -l || echo "0")
    if [ "$JS_TIMEOUTS" -gt 0 ]; then
        echo "❌ Timeouts JS hardcodés: $JS_TIMEOUTS trouvés"
        TEMPLATE_VIOLATIONS=$((TEMPLATE_VIOLATIONS + JS_TIMEOUTS))
    fi
fi

VIOLATIONS=$((VIOLATIONS + TEMPLATE_VIOLATIONS))

echo "================================================="
echo "📊 RAPPORT FINAL CODE QUALITY INSPECTOR"
echo "================================================="
echo "🔍 Total violations détectées: $VIOLATIONS"
echo "📁 Fichiers analysés: $(find . -name "*.go" -not -path "./vendor/*" | wc -l) Go + $(find templates -name "*.html" 2>/dev/null | wc -l || echo "0") HTML"
echo ""

if [ $VIOLATIONS -eq 0 ]; then
    echo "🎉 STATUS: ✅ CONFORMITÉ TOTALE"
    echo "👏 Félicitations! Aucune violation détectée."
else
    echo "🚨 STATUS: ❌ NON CONFORME"
    echo "⚠️  Action requise: Corriger TOUTES les violations"
    echo ""
    echo "📋 Détail des violations dans:"
    [ -f "$ERROR_FILE" ] && echo "   - $ERROR_FILE"
    [ -f "$REPORT_FILE" ] && echo "   - $REPORT_FILE"
    [ -f "revive-report.txt" ] && echo "   - revive-report.txt"
    [ -f "staticcheck-report.txt" ] && echo "   - staticcheck-report.txt"
fi

echo "================================================="

# Nettoyer les fichiers temporaires
rm -f "revive-report.txt" "staticcheck-report.txt"

exit $VIOLATIONS