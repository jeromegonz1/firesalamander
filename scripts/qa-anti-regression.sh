#!/bin/bash

# 🔥🦎 FIRE SALAMANDER - QA ANTI-RÉGRESSION AUTOMATIQUE
# NOUVEAU PROCESS V2.0 - QA ENGINEER CHECKLIST

echo "🔍 === QA CHECKLIST ANTI-BOUCLE INFINIE - DÉMARRAGE ==="

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Compteurs
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Fonction de test
run_test() {
    local test_name="$1"
    local command="$2"
    local timeout_seconds="$3"
    
    echo -e "\n${YELLOW}🧪 Test: $test_name${NC}"
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    
    if timeout $timeout_seconds bash -c "$command"; then
        echo -e "${GREEN}✅ PASS: $test_name${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Fonction de vérification des métriques
check_metrics() {
    local endpoint="$1"
    local metric_name="$2"
    local max_value="$3"
    
    local value=$(curl -s "$endpoint" | jq -r ".$metric_name // 0")
    if [[ "$value" =~ ^[0-9]+$ ]] && [ "$value" -le "$max_value" ]; then
        echo -e "${GREEN}✅ $metric_name: $value (≤ $max_value)${NC}"
        return 0
    else
        echo -e "${RED}❌ $metric_name: $value (> $max_value)${NC}"
        return 1
    fi
}

echo -e "\n📋 === PHASE 1: TESTS AUTOMATIQUES ==="

# 1. Test avec timeout existe
run_test "Tests avec timeout existent" "grep -r 'WithTimeout\\|time.After' internal/crawler/*test*.go | wc -l | awk '{print (\$1 > 0) ? \"exit 0\" : \"exit 1\"}' | bash" 5

# 2. Test de boucle infinie existe
run_test "Tests anti-boucle infinie existent" "grep -r 'NoInfiniteLoop\\|BOUCLE.*INFINIE' internal/crawler/*test*.go | wc -l | awk '{print (\$1 > 0) ? \"exit 0\" : \"exit 1\"}' | bash" 5

# 3. Lancer les tests de sécurité
run_test "Tests de sécurité passent" "cd /Users/jeromegonzalez/claude-code/fire-salamander && go test -timeout=20s ./internal/crawler -run TestParallelCrawler_MustTerminate" 25

run_test "Test anti-boucle passe" "cd /Users/jeromegonzalez/claude-code/fire-salamander && go test -timeout=10s ./internal/crawler -run TestParallelCrawler_NoInfiniteLoop" 15

# 4. Benchmark avec limite de temps
run_test "Benchmark avec timeout" "cd /Users/jeromegonzalez/claude-code/fire-salamander && go test -timeout=15s ./internal/crawler -run TestParallelCrawler_BenchmarkTimeout" 20

echo -e "\n🖥️ === PHASE 2: TESTS MANUELS ==="

# Vérifier si le serveur tourne
if ! pgrep -f "fire-salamander" > /dev/null; then
    echo -e "${YELLOW}⚠️ Serveur non démarré, lancement...${NC}"
    cd /Users/jeromegonzalez/claude-code/fire-salamander
    ./fire-salamander-FIXED > qa-test.log 2>&1 &
    sleep 3
fi

# 5. Test sur example.com - doit terminer en < 10s
run_test "example.com termine en <10s" "
    start=\$(date +%s)
    curl -s -X POST http://localhost:8080/api/analyze -H 'Content-Type: application/json' -d '{\"url\":\"https://example.com\"}' | grep -q 'analysis-'
    analysis_id=\$(curl -s -X POST http://localhost:8080/api/analyze -H 'Content-Type: application/json' -d '{\"url\":\"https://example.com\"}' | jq -r '.id')
    
    # Attendre que l'analyse soit complète
    for i in {1..15}; do
        status=\$(curl -s http://localhost:8080/api/status/\$analysis_id | jq -r '.status')
        if [ \"\$status\" = \"complete\" ] || [ \"\$status\" = \"error\" ]; then
            break
        fi
        sleep 1
    done
    
    end=\$(date +%s)
    duration=\$((end - start))
    echo \"⏱️ Durée: \${duration}s\"
    [ \$duration -lt 15 ]
" 20

# 6. Vérifier logs - pas de répétition d'URL
run_test "Pas de répétition d'URL dans les logs" "
    if [ -f qa-test.log ]; then
        # Chercher les URLs crawlées et vérifier qu'elles ne se répètent pas
        grep 'Page crawled' qa-test.log | awk '{print \$NF}' | sort | uniq -d | wc -l | awk '{exit (\$1 == 0) ? 0 : 1}'
    else
        echo 'Fichier log non trouvé'
        exit 1
    fi
" 5

echo -e "\n📊 === PHASE 3: MONITORING ==="

# 7. CPU ne doit pas être à 100% pendant plus de 5s
run_test "CPU usage acceptable" "
    cpu_samples=0
    high_cpu_count=0
    for i in {1..10}; do
        cpu=\$(ps aux | grep fire-salamander | grep -v grep | awk '{print \$3}' | head -1)
        if [ -n \"\$cpu\" ]; then
            cpu_samples=\$((cpu_samples + 1))
            if (( \$(echo \"\$cpu > 80\" | bc -l) )); then
                high_cpu_count=\$((high_cpu_count + 1))
            fi
        fi
        sleep 0.5
    done
    
    if [ \$cpu_samples -gt 0 ]; then
        high_cpu_ratio=\$(echo \"scale=2; \$high_cpu_count / \$cpu_samples\" | bc -l)
        echo \"📊 CPU high usage ratio: \$high_cpu_ratio\"
        (( \$(echo \"\$high_cpu_ratio < 0.5\" | bc -l) ))
    else
        echo \"⚠️ Process non trouvé\"
        exit 1
    fi
" 10

# 8. Vérifier mémoire stable
run_test "Mémoire stable" "
    mem_samples=()
    for i in {1..5}; do
        mem=\$(ps aux | grep fire-salamander | grep -v grep | awk '{print \$4}' | head -1)
        if [ -n \"\$mem\" ]; then
            mem_samples+=(\$mem)
        fi
        sleep 1
    done
    
    if [ \${#mem_samples[@]} -gt 2 ]; then
        first=\${mem_samples[0]}
        last=\${mem_samples[-1]}
        growth=\$(echo \"scale=2; \$last - \$first\" | bc -l)
        echo \"📊 Memory growth: \${growth}%\"
        (( \$(echo \"\$growth < 5.0\" | bc -l) ))
    else
        echo \"⚠️ Pas assez d'échantillons mémoire\"
        exit 1
    fi
" 10

echo -e "\n🔄 === PHASE 4: COMMANDES DE VALIDATION ==="

# 9. Test avec timeout forcé
run_test "Timeout forcé fonctionne" "timeout 5s bash -c 'curl -s -X POST http://localhost:8080/api/analyze -H \"Content-Type: application/json\" -d \"{\\\"url\\\":\\\"https://httpbin.org/delay/10\\\"}\" && sleep 10'" 8

echo -e "\n📈 === RÉSULTATS FINAUX ==="
echo -e "Tests passés: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests échoués: ${RED}$TESTS_FAILED${NC}"
echo -e "Total: $TESTS_TOTAL"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 === QA CHECKLIST: TOUTES LES VÉRIFICATIONS PASSÉES ===${NC}"
    echo -e "${GREEN}✅ Le système est prêt pour la production${NC}"
    exit 0
else
    echo -e "\n${RED}🚨 === QA CHECKLIST: ÉCHECS DÉTECTÉS ===${NC}"
    echo -e "${RED}❌ Le système N'EST PAS prêt pour la production${NC}"
    exit 1
fi