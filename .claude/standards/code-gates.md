# Quality Gates - Fire Salamander

## Gates obligatoires avant merge

### Gate 1 : Compilation ✅
```bash
go build ./...
```

### Gate 2 : Tests ✅
```bash
go test ./... -cover
# Coverage > 85%
```

### Gate 3 : Standards ✅
```bash
# Pas de hardcoding
grep -r "localhost\|8080\|300" internal/
```

### Gate 4 : Architecture ✅
- Interfaces définies
- SOLID respecté
- Pas de dépendances circulaires

### Gate 5 : Documentation ✅
- README à jour
- Commentaires sur fonctions publiques
- CHANGELOG mis à jour

**Si un gate échoue = STOP, corriger avant de continuer**