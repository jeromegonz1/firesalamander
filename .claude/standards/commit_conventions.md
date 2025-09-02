# Conventions de Commit - Fire Salamander

## Format Conventional Commits

```
type(scope): description

[body optionnel]

[footer optionnel]
```

## Types autorisés

- **feat**: Nouvelle fonctionnalité
- **fix**: Correction de bug
- **docs**: Documentation uniquement
- **test**: Ajout/modification tests
- **refactor**: Refactoring sans changement fonctionnel
- **chore**: Maintenance (build, config, etc.)
- **perf**: Amélioration performance

## Scopes recommandés

- **crawler**: Agent exploration web
- **technical**: Agent audit technique
- **semantic**: Agent analyse sémantique  
- **report**: Agent génération rapports
- **orchestrator**: Agent coordination
- **config**: Configuration et YAML
- **ci**: CI/CD et GitHub Actions
- **docs**: Documentation

## Exemples

✅ **Bons commits:**
```
feat(crawler): implement robots.txt parser with caching
fix(semantic): resolve French tokenization edge cases
test(orchestrator): add integration tests for audit pipeline
docs(readme): update installation instructions
refactor(report): extract template helpers to utils
```

❌ **Commits à éviter:**
```
fix stuff
update code
wip
test
```

## Règles

1. **Description**: Impératif présent, max 50 caractères
2. **Body**: Optionnel, explique le "pourquoi"
3. **Footer**: Références issues/breaking changes
4. **Langue**: Français pour description, anglais accepté