# Checklist Pull Request - Fire Salamander

## 📋 Avant de créer la PR

### Code & Tests
- [ ] Tous les tests passent (`make test`)
- [ ] Couverture ≥ 85% pour nouveau code
- [ ] Aucun hardcoding (config dans YAML)
- [ ] Schémas JSON validés (`make validate-schemas`)
- [ ] Code formaté (`go fmt`, `black` pour Python)

### Documentation
- [ ] `.claude/context/current_state.md` mis à jour
- [ ] ADR créé si décision technique importante
- [ ] README mis à jour si nouvelle fonctionnalité
- [ ] Epic mis à jour avec avancement

### Process
- [ ] Branche nommée selon conventions
- [ ] Commits conventionnels (`type(scope): description`)
- [ ] PR liée à issue/epic correspondant

## 🔍 Template Description PR

```markdown
## 🎯 Objectif
Résumé en 1-2 phrases de ce que fait cette PR.

## 🔧 Changements
- [ ] Fonctionnalité A ajoutée
- [ ] Bug B corrigé  
- [ ] Configuration C mise à jour

## 🧪 Tests
- [ ] Tests unitaires: X nouveaux tests
- [ ] Coverage: Y% → Z%
- [ ] Tests d'intégration si applicable

## 📋 Definition of Done
- [ ] Tous les critères DoD respectés
- [ ] CI passe (tous jobs verts)
- [ ] Documentation mise à jour

## 🔗 Références
- Epic: `.claude/epics/epic-XXX-agent.md`
- Issue: #123
- ADR: `CDC/decisions/ADR-XXX.md` (si applicable)
```

## ✅ Validation automatique

Le CI GitHub Actions vérifie automatiquement :
- Tests Go et Python
- Validation schémas JSON
- Linting et formatage
- Build réussi

## 🚫 Critères de rejet

- Tests qui échouent
- Hardcoding détecté
- Pas de tests pour nouveau code
- Documentation manquante
- Commits non conventionnels