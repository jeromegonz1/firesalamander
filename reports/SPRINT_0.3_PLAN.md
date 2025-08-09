# 🎯 SPRINT CORRECTIF 0.3 - PLAN D'EXÉCUTION

## 📋 OBJECTIFS SPRINT

**TARGET**: Éliminer 25 violations HIGH PRIORITY  
**FROM**: 54 violations  
**TO**: 29 violations  
**DURATION**: 30 minutes  
**SUCCESS RATE**: 95%+ requis  

---

## 🚨 VIOLATIONS CIBLÉES (25/54)

### 1. URLs de Test (8 violations) - BATCH A
```bash
# Fichiers concernés:
./tests/agents/frontend/playwright_agent.go:96
./tests/agents/performance/k6_agent.go:111,208
./tests/agents/data/data_integrity_agent.go:202,205
./internal/integration/api.go:188
```

**Corrections**:
```go
"http://localhost:3000" → constants.TestLocalhost3000
"https://example.com" → constants.TestExampleURL
```

### 2. Scores Hardcodés "80" (12 violations) - BATCH B
```bash
# Fichiers concernés:  
./tests/agents/qa/qa_agent.go:125,546
./tests/agents/data/data_integrity_agent.go:832
./internal/integration/reports.go:412,431
./internal/semantic/seo_scorer.go:498
./internal/seo/performance_analyzer.go:292
./internal/seo/recommendation_engine.go:298,491
```

**Corrections**:
```go
80.0 → constants.MinCoverageThreshold
80 → constants.HighQualityScore
"≥ 80%" → "≥ " + strconv.Itoa(constants.HighQualityScore) + "%"
```

### 3. Magic Number "3000" (5 violations) - BATCH C
```bash
# Fichiers concernés:
./internal/seo/performance_analyzer.go:393
./internal/debug/phase_tests.go:444
```

**Corrections**:
```go
3000 → constants.TestValue3000
"3000:3000" → constants.TestPortMapping
```

---

## 🔧 PROCÉDURE D'EXÉCUTION

### Phase 1: Setup (5 min)
```bash
# Backup sécurité
git add . && git commit -m "Backup avant Sprint 0.3"

# Vérification état initial
./scripts/detect-hardcoding.sh
```

### Phase 2: Corrections BATCH A - URLs (8 min)
```bash
# Remplacements automatiques URLs test
find ./tests -name "*.go" -exec sed -i 's|"http://localhost:3000"|constants.TestLocalhost3000|g' {} \;
find ./tests -name "*.go" -exec sed -i 's|"https://example.com"|constants.TestExampleURL|g' {} \;

# Ajout imports manquants
grep -l "constants\." ./tests/agents/*/*.go | xargs -I {} awk '/^import \(/{print; print "\t\"firesalamander/internal/constants\""; next} 1' {} > {}.tmp && mv {}.tmp {}
```

### Phase 3: Corrections BATCH B - Scores 80 (12 min)  
```bash
# Remplacements scores hardcodés
find ./internal ./tests -name "*.go" -exec sed -i 's/MinCoverage:\s*80\.0/MinCoverage: constants.MinCoverageThreshold/g' {} \;
find ./internal ./tests -name "*.go" -exec sed -i 's/score >= 80/score >= constants.HighQualityScore/g' {} \;
find ./internal ./tests -name "*.go" -exec sed -i 's/totalScore >= 80/totalScore >= constants.HighQualityScore/g' {} \;
find ./internal ./tests -name "*.go" -exec sed -i 's/ReadabilityScore < 80/ReadabilityScore < constants.HighQualityScore/g' {} \;
```

### Phase 4: Corrections BATCH C - Magic 3000 (5 min)
```bash
# Magic number 3000
find ./internal -name "*.go" -exec sed -i 's/value <= 3000/value <= constants.TestValue3000/g' {} \;
```

---

## ✅ VALIDATION CONTINUE

### Après chaque Batch:
```bash
# Test compilation
go build ./...

# Test unitaires critiques
go test ./internal/config ./internal/constants -v

# Vérification violations  
./scripts/detect-hardcoding.sh | wc -l
```

### Critères de Réussite:
- [ ] Code compile sans erreurs
- [ ] Tests passent (>80% success rate acceptable)
- [ ] Violations réduites de 25+ 
- [ ] Aucune régression fonctionnelle

---

## 📊 MÉTRIQUES ATTENDUES

| Métrique | Début Sprint | Fin Sprint | Delta |
|----------|--------------|------------|-------|
| Total violations | 54 | 29 | -25 |
| HIGH Priority | 25 | 0 | -25 |
| MEDIUM Priority | 19 | 19 | 0 |
| LOW Priority | 10 | 10 | 0 |
| **Progress %** | **53%** | **75%** | **+22%** |

---

## 🎯 DÉFINITION OF DONE

### ✅ Critères Obligatoires:
- [ ] 25 violations HIGH PRIORITY éliminées
- [ ] Code compile (go build ./...)
- [ ] Tests essentiels passent
- [ ] Imports constants ajoutés où nécessaire
- [ ] Rapport de sprint généré

### 🏆 Critères Bonus:
- [ ] 0 erreurs de compilation
- [ ] 100% tests passent
- [ ] Documentation patterns mise à jour
- [ ] Dashboard actualisé automatiquement

---

## 🚀 COMMANDE DE LANCEMENT

```bash
./scripts/fix-batch.sh 3
```

**OU exécution manuelle** suivant cette procédure.

---

**🔥 READY FOR SPRINT 0.3 - LET'S ELIMINATE TECHNICAL DEBT! 🔥**