# Definition of Done - Fire Salamander

## ✅ Critères obligatoires pour chaque tâche

### 🧪 Tests & Qualité
- [ ] Tests unitaires TDD (tests écrits AVANT le code)
- [ ] Couverture ≥ 85% pour nouveau code
- [ ] Tous les tests passent (`go test ./...`)
- [ ] Validation JSON Schema si applicable

### 🔧 Code & Configuration  
- [ ] Zéro hardcoding (valeurs dans `config/*.yaml`)
- [ ] Code suit conventions Go/Python
- [ ] Pas de TODO/FIXME en production
- [ ] Logs structurés (JSON)

### 📝 Documentation
- [ ] Mise à jour `.claude/context/current_state.md`
- [ ] ADR créé si décision technique importante
- [ ] Changelog mis à jour si feature/breaking change
- [ ] README mis à jour si nouvelle commande/usage

### 🔄 Process
- [ ] Commit conventionnel (`type(scope): description`)
- [ ] PR review si travail collaboratif
- [ ] CI GitHub Actions passe (vert)
- [ ] Aucune régression introduite

### 🎯 Fonctionnel
- [ ] Spécification respectée
- [ ] Contrats API validés
- [ ] Cas d'erreur gérés
- [ ] Performance dans contraintes (.claude/context/constraints.md)

## 🚫 Bloquants absolus
- Tests qui échouent
- Hardcoding détecté  
- Schémas JSON invalides
- Régression de performance
- Sécurité compromise