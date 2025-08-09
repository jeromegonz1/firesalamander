# 📋 VALIDATION NOTES - Fire Salamander

## 🎯 PROCESSUS DE VALIDATION AUTOMATIQUE

### Convention adoptée : 
À la fin de chaque développement, ouvrir automatiquement le fichier de test HTML dans le navigateur pour validation directe.

### Commande de validation :
```bash
open "file:///Users/jeromegonzalez/claude-code/fire-salamander/test_frontend_complete.html"
```

### Fichiers de test disponibles :
- `test_frontend_complete.html` - Test complet du frontend avec tous les cas d'usage
- `test_expandedissues.html` - Test spécifique pour la fonctionnalité expandedIssues
- Autres fichiers de test à créer selon les besoins

### Workflow de validation :
1. 🔧 Développement terminé
2. 🚀 Serveur démarré (`go run cmd/server/main.go`)
3. 📂 Ouverture automatique du fichier de test
4. ✅ Validation par clicks directs sur les liens
5. 📝 Feedback immédiat

---
**Note ajoutée le :** 2025-08-07
**Par :** Architecte Principal Claude Code