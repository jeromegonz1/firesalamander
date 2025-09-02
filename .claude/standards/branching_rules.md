# Règles de Branching - Fire Salamander

## Structure des branches

```
main
├── feature/agent-crawler-v2
├── feature/semantic-ml-integration  
├── fix/orchestrator-timeout
└── docs/api-documentation
```

## Conventions de nommage

### Format général
`type/description-kebab-case`

### Types de branches
- **feature/**: Nouvelles fonctionnalités
- **fix/**: Corrections de bugs
- **docs/**: Documentation uniquement
- **refactor/**: Refactoring sans changement fonctionnel
- **test/**: Tests ou amélioration coverage
- **chore/**: Maintenance, build, config

### Exemples
```
feature/agent-crawler-sitemap-support
feature/semantic-camembert-integration
fix/report-template-encoding-issue
docs/api-contracts-documentation
refactor/orchestrator-error-handling
test/crawler-edge-cases-coverage
```

## Workflow Git

1. **Création branche**:
   ```bash
   git checkout -b feature/agent-crawler-advanced
   ```

2. **Développement**: Commits atomiques fréquents

3. **Tests**: `make test` avant push

4. **Push**: 
   ```bash
   git push -u origin feature/agent-crawler-advanced
   ```

5. **Pull Request**: Vers `main` avec template

6. **Merge**: Squash and merge après review

## Règles de protection

- `main`: Branch protégée
- PR review obligatoire pour `main`
- CI doit passer (tous tests verts)
- Pas de force push sur `main`

## Naming tips

- Max 50 caractères
- Descriptif et spécifique
- Éviter acronymes obscurs
- Utiliser scope agent si applicable