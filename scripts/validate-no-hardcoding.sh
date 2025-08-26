#!/bin/bash

# ========================================
# FIRE SALAMANDER - VALIDATION ANTI-HARDCODING
# Script de contrôle automatique pour conformité SEPTEO
# ========================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Couleurs pour output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
VIOLATIONS_FOUND=0
FILES_CHECKED=0

echo -e "${BLUE}🔥 FIRE SALAMANDER - VALIDATION ANTI-HARDCODING${NC}"
echo -e "${BLUE}=================================================${NC}"
echo ""

# Fonction pour afficher une violation
report_violation() {
    local file="$1"
    local line="$2"
    local content="$3"
    local reason="$4"
    
    echo -e "${RED}❌ VIOLATION DÉTECTÉE${NC}"
    echo -e "   📄 Fichier: ${file}"
    echo -e "   📍 Ligne: ${line}"
    echo -e "   💬 Contenu: ${content}"
    echo -e "   🔍 Raison: ${reason}"
    echo ""
    
    VIOLATIONS_FOUND=$((VIOLATIONS_FOUND + 1))
}

# Fonction pour vérifier un fichier
check_file() {
    local file="$1"
    FILES_CHECKED=$((FILES_CHECKED + 1))
    
    echo -e "${BLUE}🔍 Vérification: ${file}${NC}"
    
    # Patterns de hardcoding à détecter
    local patterns=(
        # Timeouts hardcodés
        ":[[:space:]]*[0-9]+[[:space:]]*\\*[[:space:]]*time\\.Second"
        ":[[:space:]]*[0-9]+[[:space:]]*\\*[[:space:]]*time\\.Millisecond"
        ":[[:space:]]*[0-9]+[[:space:]]*\\*[[:space:]]*time\\.Minute"
        
        # Tailles hardcodées
        "[[:space:]][0-9]+[[:space:]]*\\*[[:space:]]*1024[[:space:]]*\\*[[:space:]]*1024"
        "maxBodySize[[:space:]]*:=[[:space:]]*[0-9]+"
        "maxURLs[[:space:]]*:[[:space:]]*[0-9]+"
        
        # Ports hardcodés
        "localhost:[0-9]+"
        ":[0-9]{4,5}[^0-9]"
        
        # Limites hardcodées
        "len\\([^)]+\\)[[:space:]]*>=[[:space:]]*[0-9]+"
        "len\\([^)]+\\)[[:space:]]*<[[:space:]]*[0-9]+"
        
        # User agents hardcodés
        "User-Agent.*[^constants\\.]"
        
        # Seuils de performance hardcodés
        "[0-9]+ms"
        "threshold.*[0-9]+"
    )
    
    local descriptions=(
        "Timeout hardcodé avec time.Second"
        "Timeout hardcodé avec time.Millisecond"  
        "Timeout hardcodé avec time.Minute"
        "Taille mémoire hardcodée (MB)"
        "Taille body hardcodée"
        "Limite URLs hardcodée"
        "Port localhost hardcodé"
        "Port réseau hardcodé"
        "Limite de longueur hardcodée (>=)"
        "Limite de longueur hardcodée (<)"
        "User-Agent hardcodé"
        "Seuil en millisecondes hardcodé"
        "Seuil/threshold hardcodé"
    )
    
    # Vérifier chaque pattern
    for i in "${!patterns[@]}"; do
        local pattern="${patterns[$i]}"
        local description="${descriptions[$i]}"
        
        # Rechercher le pattern dans le fichier
        while IFS= read -r line_content; do
            local line_number=$(echo "$line_content" | cut -d: -f1)
            local line_text=$(echo "$line_content" | cut -d: -f2-)
            
            # Vérifier si c'est dans un commentaire (ignorer)
            if [[ "$line_text" =~ ^[[:space:]]*// ]]; then
                continue
            fi
            
            # Exclure les lignes avec constants.
            if [[ "$line_text" =~ constants\. ]]; then
                continue
            fi
            
            # Exclure les définitions de constantes elles-mêmes
            if [[ "$line_text" =~ ^[[:space:]]*const[[:space:]] ]] || [[ "$line_text" =~ ^[[:space:]]*[A-Z][A-Za-z0-9_]*[[:space:]]*= ]]; then
                continue
            fi
            
            report_violation "$file" "$line_number" "$(echo "$line_text" | xargs)" "$description"
            
        done < <(grep -n -E "$pattern" "$file" 2>/dev/null || true)
    done
    
    # Vérifications spécifiques par type de fichier
    if [[ "$file" == *_test.go ]]; then
        check_test_file_hardcoding "$file"
    fi
    
    if [[ "$file" == *crawler*.go ]]; then
        check_crawler_specific_hardcoding "$file"
    fi
}

# Vérifications spécifiques aux fichiers de test
check_test_file_hardcoding() {
    local file="$1"
    
    # Vérifier les timeouts de test hardcodés
    while IFS= read -r line_content; do
        local line_number=$(echo "$line_content" | cut -d: -f1)
        local line_text=$(echo "$line_content" | cut -d: -f2-)
        
        if [[ ! "$line_text" =~ constants\. ]]; then
            report_violation "$file" "$line_number" "$(echo "$line_text" | xargs)" "Timeout de test hardcodé"
        fi
    done < <(grep -n -E "[0-9]+[[:space:]]*\\*[[:space:]]*time\\.(Second|Millisecond)" "$file" 2>/dev/null || true)
}

# Vérifications spécifiques aux crawlers
check_crawler_specific_hardcoding() {
    local file="$1"
    
    # Vérifier les URLs de test hardcodées
    while IFS= read -r line_content; do
        local line_number=$(echo "$line_content" | cut -d: -f1)
        local line_text=$(echo "$line_content" | cut -d: -f2-)
        
        if [[ "$line_text" =~ \"https?://[^\"]+\" ]] && [[ ! "$line_text" =~ constants\. ]]; then
            report_violation "$file" "$line_number" "$(echo "$line_text" | xargs)" "URL hardcodée dans le code"
        fi
    done < <(grep -n -E "\"https?://[^\"]+\"" "$file" 2>/dev/null || true)
    
    # Vérifier les retry attempts hardcodés
    while IFS= read -r line_content; do
        local line_number=$(echo "$line_content" | cut -d: -f1)
        local line_text=$(echo "$line_content" | cut -d: -f2-)
        
        if [[ "$line_text" =~ [Rr]etry.*[0-9]+ ]] && [[ ! "$line_text" =~ constants\. ]]; then
            report_violation "$file" "$line_number" "$(echo "$line_text" | xargs)" "Nombre de retry hardcodé"
        fi
    done < <(grep -n -E "[Rr]etry.*[0-9]+" "$file" 2>/dev/null || true)
}

# Fonction pour vérifier l'utilisation correcte des constantes
check_constants_usage() {
    echo -e "${BLUE}🔧 Vérification de l'utilisation des constantes...${NC}"
    
    # Vérifier que les fichiers importent bien le package constants
    local crawler_files=($(find "$PROJECT_ROOT/internal/crawler" -name "*.go" -not -name "*_test.go"))
    
    for file in "${crawler_files[@]}"; do
        if ! grep -q "firesalamander/internal/constants" "$file"; then
            echo -e "${YELLOW}⚠️ Warning: ${file} n'importe pas le package constants${NC}"
        fi
    done
}

# Fonction principale
main() {
    echo -e "${BLUE}Démarrage de la validation anti-hardcoding...${NC}"
    echo ""
    
    # Vérifier que nous sommes dans le bon répertoire
    if [[ ! -d "$PROJECT_ROOT/internal/crawler" ]]; then
        echo -e "${RED}❌ Erreur: Répertoire internal/crawler non trouvé${NC}"
        exit 1
    fi
    
    # Trouver tous les fichiers Go dans le crawler
    local files=($(find "$PROJECT_ROOT/internal/crawler" -name "*.go"))
    
    if [[ ${#files[@]} -eq 0 ]]; then
        echo -e "${RED}❌ Erreur: Aucun fichier Go trouvé dans internal/crawler${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}📁 ${#files[@]} fichiers Go trouvés à vérifier${NC}"
    echo ""
    
    # Vérifier chaque fichier
    for file in "${files[@]}"; do
        check_file "$file"
    done
    
    # Vérifications supplémentaires
    check_constants_usage
    
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}RÉSULTATS DE LA VALIDATION${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "📊 Fichiers vérifiés: ${FILES_CHECKED}"
    
    if [[ $VIOLATIONS_FOUND -eq 0 ]]; then
        echo -e "${GREEN}✅ SUCCÈS: Aucune violation de hardcoding détectée!${NC}"
        echo -e "${GREEN}🔥 Le code respecte les standards SEPTEO${NC}"
        exit 0
    else
        echo -e "${RED}❌ ÉCHEC: ${VIOLATIONS_FOUND} violation(s) de hardcoding détectée(s)${NC}"
        echo -e "${RED}🚨 Le code ne respecte pas les standards SEPTEO${NC}"
        echo ""
        echo -e "${YELLOW}📋 Actions recommandées:${NC}"
        echo -e "${YELLOW}1. Déplacer les valeurs hardcodées vers internal/constants/crawler_constants.go${NC}"
        echo -e "${YELLOW}2. Utiliser les constantes dans le code: constants.NomConstante${NC}"
        echo -e "${YELLOW}3. Externaliser la configuration via environment variables${NC}"
        echo -e "${YELLOW}4. Re-exécuter ce script pour validation${NC}"
        exit 1
    fi
}

# Vérifier les paramètres de ligne de commande
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -h, --help     Afficher cette aide"
            echo "  -v, --verbose  Mode verbose"
            echo ""
            echo "Ce script vérifie l'absence de hardcoding dans le module crawler."
            echo "Conformité aux standards SEPTEO - Zéro tolérance hardcoding."
            exit 0
            ;;
        -v|--verbose)
            set -x
            shift
            ;;
        *)
            echo "Option inconnue: $1"
            exit 1
            ;;
    esac
done

# Exécution
main "$@"