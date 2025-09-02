# Checklist de Session Fire Salamander

## 📋 Début de session
- [ ] Exécuter `make context` pour charger le contexte
- [ ] Vérifier la branche Git actuelle
- [ ] Charger l'epic en cours depuis `.claude/epics/`
- [ ] Vérifier les tests actuels : `go test ./...`

## 🔧 Pendant la session
- [ ] Suivre TDD : tests d'abord, code ensuite
- [ ] Zero hardcoding : toute config dans YAML
- [ ] Commits atomiques avec messages conventionnels

## 🏁 Fin de session
- [ ] Mettre à jour `.claude/context/current_state.md`
- [ ] Logger dans `.claude/memory/session_$(date +%Y%m%d).md`
- [ ] Exécuter les tests : `make test`
- [ ] Commit avec format : `type(scope): description`
- [ ] Push seulement si tests passent