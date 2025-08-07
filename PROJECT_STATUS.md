# Fire Salamander - Project Status

## Version: 2.0 (MVP Restart)
## Architecte Principal: Claude Code
## Date: 2025-01-07

---

## 🏗️ PRINCIPES D'ARCHITECTURE NON-NÉGOCIABLES

### 1. **TDD OBLIGATOIRE** ✅
- Tests d'abord, code ensuite
- RED → GREEN → REFACTOR
- Couverture minimale : 80%
- `go test -cover ./...`

### 2. **NO HARDCODING POLICY** ❌
- Toute valeur dans .env ou configuration
- Aucune chaîne en dur dans le code
- Validation au build : `grep -r "localhost\|8080\|http://" . --exclude-dir=archive`

### 3. **ERROR HANDLING PROFESSIONNEL** 🛡️
- Jamais de `panic()` en production
- Toujours wrapper les erreurs avec contexte
- Format : `fmt.Errorf("operation failed: %w", err)`

### 4. **SOLID PRINCIPLES** 📐
- Single Responsibility : Une fonction = Une responsabilité
- Open/Closed : Extension sans modification
- Interface Segregation : Interfaces spécifiques
- Dependency Inversion : Abstractions, pas de concret

### 5. **CLEAN CODE** 🧹
- Noms explicites (pas d'abréviations)
- Fonctions max 20 lignes
- Commentaires uniquement pour le "pourquoi", pas le "quoi"

---

## 🎯 OBJECTIF MVP (FOCUS STRICT)

### Scope Défini
- ✅ **Analyse SEO basique** (20 pages maximum)
- ✅ **Score simple** (title, meta, h1, images)
- ✅ **Export PDF minimaliste**
- ✅ **Interface web native Go** (html/template)

### SCOPE CREEP INTERDIT ❌
- ❌ Pas de JavaScript frameworks
- ❌ Pas d'ORM complexe
- ❌ Pas d'API REST complète en V1
- ❌ Pas de système d'authentification en V1

---

## 📋 ARCHITECTURE DÉCIDÉE (IMMUTABLE)

```
fire-salamander/
├── cmd/server/           # Point d'entrée uniquement
├── internal/             # Logique métier (non exportée)
│   ├── config/          # Configuration externalisée
│   ├── analyzer/        # Analyse SEO (SOLID)
│   └── crawler/         # Récupération pages (Single Resp.)
├── templates/           # HTML templates Go natifs
├── static/              # CSS minimal, pas de JS lourd
└── tests/               # TDD obligatoire
    ├── unit/           # Tests unitaires
    └── integration/    # Tests d'intégration
```

---

## ✅ ACCOMPLI (Validation Architecte)

- [x] **Archive V1** - Sauvegardée avec documentation post-mortem
- [x] **Nettoyage radical** - Repo propre avec .git préservé
- [x] **Structure MVP** - Séparation claire des responsabilités
- [x] **Configuration externalisée** - .env.example créé, zéro hardcoding
- [x] **Standards qualité** - .gitignore, PROJECT_STATUS.md

---

## 🚧 PROCHAINES ÉTAPES (TDD STRICT)

### Phase 1 : Foundation (Current)
1. [ ] **TDD Config Loader** - Tests RED puis implémentation GREEN
2. [ ] **TDD Basic Server** - HTTP server avec graceful shutdown
3. [ ] **TDD Template Engine** - Rendering HTML simple

### Phase 2 : Core Business
4. [ ] **TDD URL Crawler** - Récupération title/meta uniquement
5. [ ] **TDD SEO Analyzer** - Score basique (0-100)
6. [ ] **TDD Report Generator** - Export PDF minimal

### Phase 3 : MVP Completion
7. [ ] **Integration Tests** - E2E workflow complet
8. [ ] **Performance Tests** - Load testing avec k6
9. [ ] **Security Audit** - Vulnérabilité scanning

---

## 🔧 COMMANDES STANDARDS

```bash
# Tests (OBLIGATOIRE avant commit)
go test ./...
go test -cover ./... -coverprofile=coverage.txt

# Build (zero warnings accepté)
go build -v -o fire-salamander

# Run (avec .env local)
cp .env.example .env
go run main.go

# Linting (installation requise)
golangci-lint run --enable-all

# Security scan
gosec ./...
```

---

## 📝 DÉCISIONS TECHNIQUES DÉFINITIVES

| Composant | Choix | Justification |
|-----------|-------|---------------|
| **Web Server** | `net/http` natif | Simplicité, performance, pas de dépendance |
| **Templates** | `html/template` | Sécurité XSS native, standard Go |
| **Database** | SQLite | Zéro configuration, parfait pour MVP |
| **Config** | `.env` + `os.Getenv()` | Simple, standard, pas de dépendance |
| **Logging** | `slog` (Go 1.21+) | Structured logging natif |
| **Testing** | `testing` standard | TDD natif, pas de framework externe |

---

## 🚨 RED FLAGS (REFUS AUTOMATIQUE)

### Code Review Blockers
- ❌ Hardcoded values (strings, numbers, URLs)
- ❌ `panic()` en production
- ❌ Fonctions > 20 lignes sans justification
- ❌ Tests manquants pour nouvelle feature
- ❌ Noms de variables non explicites (`d`, `tmp`, `data`)

### Architecture Violations
- ❌ Import de packages externes non justifiés
- ❌ Logique métier dans les handlers HTTP
- ❌ SQL queries inline dans le business logic
- ❌ Configuration mélangée avec le code

---

## 📊 MÉTRIQUES QUALITÉ (Monitoring Continu)

```bash
# Coverage minimum
go test -cover ./... | grep "coverage:" | awk '{if($3+0 < 80) exit 1}'

# Complexité cyclomatique (gocyclo)
gocyclo -over 10 .

# Duplication code (dupl)
dupl -t 100 ./...

# Vulnerabilities (gosec)
gosec -quiet ./...
```

---

## 🎯 DÉFINITION OF DONE

### Pour chaque feature :
1. ✅ Tests écrits AVANT le code (TDD)
2. ✅ Coverage ≥ 80%
3. ✅ Zéro hardcoding détecté
4. ✅ Documentation technique à jour
5. ✅ Code review par architecte
6. ✅ Tests d'intégration passent
7. ✅ Performance tests OK
8. ✅ Security scan clean

---

**Architecte Principal :** Claude Code  
**Dernière Révision :** 2025-01-07  
**Prochaine Révision :** Après Phase 1 TDD