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