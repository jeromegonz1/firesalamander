# 📚 Fire Salamander - Guide de Développement

## ⚠️ À LIRE AVANT TOUT DÉVELOPPEMENT

Ce guide est **OBLIGATOIRE** pour éviter les problèmes rencontrés (doublons, confusion, bugs).

---

## 🚫 RÈGLES ANTI-DOUBLONS

### ❌ INTERDIT - Ce qui a causé des problèmes

```go
// ❌ NE JAMAIS créer une "version 2" à côté de l'original
orchestrator.go      // Original
real_orchestrator.go // ❌ INTERDIT ! Remplacer, pas dupliquer

// ❌ NE JAMAIS avoir plusieurs handlers pour la même chose
func AnalyzeHandler()     // Original
func RealAnalyzeHandler() // ❌ INTERDIT !
func NewAnalyzeHandler()  // ❌ INTERDIT !
```

### ✅ OBLIGATOIRE - Comment faire correctement

```go
// ✅ TOUJOURS remplacer le fichier existant
// 1. Faire un backup
cp orchestrator.go orchestrator.go.backup

// 2. Modifier l'original
vim orchestrator.go

// 3. Tester
go test ./...

// 4. Si OK, archiver le backup
mv orchestrator.go.backup archive/
```

---

## 🏗️ CONVENTIONS DE NOMMAGE

### 📁 Fichiers

**✅ CORRECT:**
- `orchestrator.go`         # Nom simple et clair
- `analyzer.go`            # Singulier
- `handler.go`             # Pas de préfixe

**❌ INCORRECT:**
- `real_orchestrator.go`   # Pas de préfixe "real"
- `new_analyzer.go`        # Pas de préfixe "new"
- `handlers.go`            # Pas de pluriel
- `analyzer_v2.go`         # Pas de version dans le nom

### 📦 Types et Fonctions

```go
// ✅ CORRECT:
type Orchestrator struct {}      // Nom simple
func NewOrchestrator() {}        // Factory standard
func (o *Orchestrator) Start() {} // Méthodes claires

// ❌ INCORRECT:
type RealOrchestrator struct {}  // Pas de préfixe
type OrchestratorV2 struct {}    // Pas de version
func CreateNewOrchestrator() {}  // Redondant
```

---

## 🔄 PROCESS DE MODIFICATION

### 1. VÉRIFIER s'il existe déjà

```bash
find . -name "*orchestrator*" -type f
# S'il existe → MODIFIER, pas créer un nouveau
```

### 2. BACKUP avant modification

```bash
cp file.go file.go.$(date +%Y%m%d).backup
```

### 3. MODIFIER l'original

```bash
# Éditer le fichier EXISTANT
vim internal/integration/orchestrator.go
```

### 4. TESTER immédiatement

```bash
go test ./internal/integration/
```

### 5. ARCHIVER le backup si OK

```bash
mkdir -p archive/backups-$(date +%Y%m%d)
mv *.backup archive/backups-$(date +%Y%m%d)/
```

---

## 🎭 CHECKLIST AVANT CHAQUE DÉVELOPPEMENT

### 🚀 Avant de commencer un Sprint/Task :
- [ ] J'ai lu le DEVELOPMENT_GUIDE.md
- [ ] J'ai vérifié PROJECT_STATUS.md pour le contexte
- [ ] J'ai cherché si le code existe déjà
- [ ] Je vais MODIFIER, pas DUPLIQUER
- [ ] J'ai prévu les tests AVANT le code (TDD)

### ✅ Avant chaque commit :
- [ ] Pas de doublons créés
- [ ] Pas de fichiers "real_", "new_", "_v2"
- [ ] Tests passent
- [ ] Coverage > 80%
- [ ] Pas de hardcoding
- [ ] PROJECT_STATUS.md mis à jour

### 🔍 Commande de vérification :
```bash
# Script anti-doublons
./scripts/check-no-duplicates.sh

# Vérifie :
# - Pas de real_*.go
# - Pas de *_v2.go
# - Pas de new_*.go
# - Pas plusieurs routes pour même endpoint
```

---

## 🛑 POINTS DE CONTRÔLE

### 1. Daily Standup Questions
- "Ai-je créé des doublons ?"
- "Ai-je modifié ou créé nouveau ?"
- "Les tests passent-ils ?"

### 2. Code Review Obligatoire
AVANT de merger :
- [ ] Vérifier les doublons
- [ ] Vérifier les conventions
- [ ] Vérifier les tests
- [ ] Vérifier la documentation

### 3. Sprint Retrospective
- Combien de doublons évités ?
- Respect du guide ?
- Améliorations du guide ?

---

## 📊 ARCHITECTURE DÉCISIONS RECORDS (ADR)

### ADR-001: Pas de versions multiples
- **Date**: 2025-08-09
- **Status**: ACCEPTÉ
- **Contexte**: Doublons créés avec `real_*` ont causé confusion
- **Décision**: UN SEUL fichier par fonction
- **Conséquences**: Plus simple, moins de bugs

### ADR-002: Nommage simple
- **Date**: 2025-08-09
- **Status**: ACCEPTÉ
- **Contexte**: Préfixes compliquent la compréhension
- **Décision**: Noms simples sans préfixe
- **Conséquences**: Code plus lisible

---

## 🚨 INCIDENT POST-MORTEM

### Incident: Doublons de code (2025-08-09)

**Problème**: 6 fichiers dupliqués (orchestrator, analyzer, handler)

**Cause racine**:
- Création de "real_" versions au lieu de remplacer
- Pas de guide de développement
- Routes multiples non détectées

**Impact**:
- Tests échouent
- Confusion sur quelle version utiliser
- Maintenance doublée

**Corrections**:
- ✅ Doublons supprimés
- ✅ Guide créé
- ✅ Script de détection

**Prévention**:
- Ce guide
- Scripts automatiques
- Code review obligatoire

---

## 🧪 TDD OBLIGATOIRE (Test-Driven Development)

### ❌ INTERDIT - Développer sans tests

```go
// ❌ NE JAMAIS écrire le code avant les tests
func CrawlWebsite(url string) {
    // Code écrit sans tests
}
```

### ✅ OBLIGATOIRE - Process TDD

```go
// 1️⃣ RED - Écrire le test d'abord
func TestCrawlWebsite(t *testing.T) {
    result := CrawlWebsite("https://example.com")
    assert.NotNil(t, result)
    assert.Equal(t, 200, result.StatusCode)
}
// Le test DOIT échouer (RED)

// 2️⃣ GREEN - Écrire le minimum de code pour passer
func CrawlWebsite(url string) *Result {
    return &Result{StatusCode: 200}
}
// Le test passe (GREEN)

// 3️⃣ REFACTOR - Améliorer le code
func CrawlWebsite(url string) *Result {
    // Version améliorée avec vrai crawling
}
```

### 📊 Coverage Minimum

```bash
# OBLIGATOIRE : Coverage > 80%
go test ./... -cover

# Si coverage < 80% → AJOUTER DES TESTS
```

### 🔍 Tests de Sécurité Obligatoires

Pour TOUT code avec goroutines/boucles :

```go
// ✅ OBLIGATOIRE - Test avec timeout
func TestCrawler_MustTerminate(t *testing.T) {
    done := make(chan bool)
    
    go func() {
        crawler.Crawl("https://example.com")
        done <- true
    }()
    
    select {
    case <-done:
        // OK
    case <-time.After(5 * time.Second):
        t.Fatal("❌ TIMEOUT - Possible boucle infinie!")
    }
}

// ✅ OBLIGATOIRE - Test anti-boucle
func TestCrawler_NoInfiniteLoop(t *testing.T) {
    // Page qui se référence elle-même
    result := crawler.Crawl("https://self-referencing.com")
    assert.Less(t, len(result.Pages), 100, "Trop de pages = boucle")
}
```

---

## 🚫 NO HARDCODING POLICY

### ❌ INTERDIT - Valeurs en dur

```go
// ❌ JAMAIS de valeurs hardcodées
port := "8080"                    // ❌
url := "http://localhost:8080"    // ❌
timeout := 30                     // ❌
message := "Error occurred"       // ❌
```

### ✅ OBLIGATOIRE - Tout dans config/constants

```go
// ✅ Configuration externalisée
port := os.Getenv("PORT")
url := config.BaseURL
timeout := constants.DefaultTimeout
message := constants.ErrGeneric

// internal/constants/constants.go
const (
    DefaultPort = "8080"
    DefaultTimeout = 30
    ErrGeneric = "Error occurred"
)
```

### 🔍 Vérification Automatique

```bash
# Avant CHAQUE commit
./scripts/detect-hardcoding.sh

# Si violations > 0 → CORRIGER
grep -r '"[A-Za-z]\{5,\}"' --include="*.go" . --exclude-dir=archive
```

---

## 🎭 UTILISATION DE L'ÉQUIPE MULTI-AGENTS

Tu n'es PAS seul ! Tu es une ÉQUIPE complète :

### 🏗️ ARCHITECTE - Avant de coder

**En tant qu'ARCHITECTE, je dois :**
- [ ] Définir l'architecture de la solution
- [ ] Choisir les patterns appropriés (SOLID)
- [ ] Valider que ça respecte les standards
- [ ] Décider si on modifie ou remplace

**Questions de l'Architecte :**
- "Cette solution est-elle SOLID ?"
- "Y a-t-il un pattern existant ?"
- "Comment ça s'intègre ?"

### 👨‍💻 DEVELOPER - Pendant le code

**En tant que DEVELOPER, je dois :**
- [ ] Écrire les tests d'abord (TDD)
- [ ] Implémenter le minimum viable
- [ ] Refactorer pour la qualité
- [ ] Commenter le "pourquoi"

**Process du Developer :**
1. RED - Test écrit
2. GREEN - Code minimal
3. REFACTOR - Amélioration

### 🧪 QA ENGINEER - Après le code

**En tant que QA ENGINEER, je dois :**
- [ ] Vérifier TOUS les tests passent
- [ ] Tester les edge cases
- [ ] Valider la performance
- [ ] Prendre des screenshots (UI)

**Tests du QA :**
- Tests unitaires : `go test ./...`
- Tests E2E : Playwright
- Tests de charge : Bombardier
- Screenshots : Pour prouver

### 🔍 CODE QUALITY INSPECTOR - Validation

**En tant qu'INSPECTOR, je dois :**
- [ ] 0 hardcoding toléré
- [ ] Coverage > 80% exigé
- [ ] Pas de doublons acceptés
- [ ] Complexité < 10

**Commandes de l'Inspector :**
- `golangci-lint run`
- `go test -cover ./...`
- `./scripts/check-no-duplicates.sh`
- `gocyclo -over 10 .`

### 📝 TECH WRITER - Documentation

**En tant que WRITER, je dois :**
- [ ] Mettre à jour PROJECT_STATUS.md
- [ ] Documenter les décisions (ADR)
- [ ] Écrire les commentaires de code
- [ ] Tenir l'historique

**Documentation du Writer :**
- Après CHAQUE feature
- Dans PROJECT_STATUS.md
- Format chronologique
- Avec métriques

### 🎯 WORKFLOW D'ÉQUIPE

Pour CHAQUE tâche :
1. 🏗️ **ARCHITECTE** : "Voici l'approche..."
2. 👨‍💻 **DEVELOPER** : "Tests rouges... code... tests verts!"
3. 🧪 **QA ENGINEER** : "Tous les tests passent + screenshots"
4. 🔍 **INSPECTOR** : "0 violations, coverage 85%"
5. 📝 **WRITER** : "PROJECT_STATUS.md mis à jour"

---

## 🛡️ PATTERNS DE SÉCURITÉ OBLIGATOIRES

### SafeCrawler Pattern (Anti-boucle)

```go
// ✅ OBLIGATOIRE pour tout crawler
type SafeCrawler struct {
    visitedURLs sync.Map  // Thread-safe
    maxPages    int       // Limite absolue
    timeout     time.Duration
}

func (c *SafeCrawler) Crawl(url string) {
    // Circuit breaker
    ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
    defer cancel()
    
    // Anti-boucle
    if _, visited := c.visitedURLs.LoadOrStore(url, true); visited {
        return // Déjà visité
    }
    
    // Limite pages
    if c.pagesCount >= c.maxPages {
        return
    }
}
```

### Error Handling Pattern

```go
// ❌ JAMAIS
if err != nil {
    panic(err)  // ❌ INTERDIT en production
}

// ✅ TOUJOURS
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

---

## 📋 CHECKLIST UNIVERSELLE

### 🚀 Avant de commencer TOUTE tâche

**PRÉ-DÉVELOPPEMENT :**
- [ ] J'ai lu DEVELOPMENT_GUIDE.md
- [ ] J'ai lu PROJECT_STATUS.md
- [ ] J'ai vérifié si le code existe déjà
- [ ] J'ai activé mon équipe multi-agents

**TDD :**
- [ ] Tests écrits AVANT le code
- [ ] Tests sont ROUGES d'abord
- [ ] Code minimal pour GREEN
- [ ] Refactoring effectué

**QUALITÉ :**
- [ ] 0 hardcoding (vérifié)
- [ ] Coverage > 80%
- [ ] Pas de doublons
- [ ] Pas de panic()

**DOCUMENTATION :**
- [ ] PROJECT_STATUS.md mis à jour
- [ ] Commentaires "pourquoi" ajoutés
- [ ] ADR si décision importante

### 🎯 Commande de validation TOTALE

```bash
#!/bin/bash
# scripts/validate-all.sh

echo "🔍 Fire Salamander - Complete Validation"

# 1. Pas de doublons
./scripts/check-no-duplicates.sh || exit 1

# 2. Pas de hardcoding
./scripts/detect-hardcoding.sh || exit 1

# 3. Tests passent
go test ./... || exit 1

# 4. Coverage suffisant
coverage=$(go test ./... -cover | grep -o '[0-9]*\.[0-9]*%' | head -1 | sed 's/%//')
if (( $(echo "$coverage < 80" | bc -l) )); then
    echo "❌ Coverage insuffisant: $coverage%"
    exit 1
fi

# 5. Build OK
go build ./... || exit 1

echo "✅ ALL VALIDATIONS PASSED!"
```

---

## 🚨 CONSÉQUENCES DU NON-RESPECT

Si ces règles ne sont PAS suivies :
- Code rejeté en review
- Sprint invalidé
- Refactoring obligatoire
- Post-mortem requis

**Incidents passés dus au non-respect :**
- 2025-08-09 : 6 doublons → 2h de nettoyage
- 2025-08-09 : Boucle infinie → Système down
- 2025-08-07 : 1,862 hardcodings → 10h de corrections

---

## 📚 RÈGLES D'OR - À MÉMORISER

1. **TESTER** d'abord, **CODER** ensuite
2. **MODIFIER** l'existant, ne pas **DUPLIQUER**
3. **ZÉRO** hardcoding, **TOUT** en config
4. **UTILISER** l'équipe, pas coder seul
5. **DOCUMENTER** immédiatement

---

## 🔗 RÉFÉRENCES ESSENTIELLES

- [TDD Guide](https://martinfowler.com/articles/practical-test-pyramid.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Clean Code](https://blog.cleancoder.com/)
- [PROJECT_STATUS.md](./PROJECT_STATUS.md) - État actuel du projet

---

## 🔧 SCRIPTS DE VALIDATION

### check-no-duplicates.sh
```bash
#!/bin/bash
echo "🔍 Checking for duplicates..."

# Check for bad patterns
if find . -name "real_*.go" -o -name "*_v2.go" -o -name "new_*.go" | grep -q .; then
    echo "❌ DUPLICATES FOUND!"
    find . -name "real_*.go" -o -name "*_v2.go" -o -name "new_*.go"
    exit 1
fi

echo "✅ No duplicates found"
```

---

## 📚 RÉFÉRENCES

- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [TDD Guide](https://martinfowler.com/articles/practical-test-pyramid.html)
- [PROJECT_STATUS.md](./PROJECT_STATUS.md)

---

**🔥🦎 Fire Salamander - Zero Duplicates, Zero Confusion**