# FAQ - Questions fréquentes

## Général

### Q: Quelle est la différence avec Screaming Frog ?
Fire Salamander est optimisé pour le marché français avec analyse sémantique en français via CamemBERT.

### Q: Combien de pages peuvent être analysées ?
300 pages par défaut, configurable jusqu'à 1000.

### Q: Combien de temps prend un audit ?
- Site 50 pages : 2-3 minutes
- Site 300 pages : 5-8 minutes  
- Site 1000 pages : 15-20 minutes

### Q: Fire Salamander fonctionne-t-il sur tous les sites ?
Oui, mais optimisé pour :
- Sites français (détection langue automatique)
- Sites accessibles publiquement
- Sites respectant les standards web

## Technique

### Q: L'audit échoue avec "timeout"
**Solutions :**
- Vérifier la connexion internet
- Réduire max_urls dans `config/crawler.yaml`
- Augmenter les timeouts dans la config

```yaml
# config/crawler.yaml
performance:
  request_timeout: 30s  # Augmenter si sites lents
```

### Q: Pas de suggestions sémantiques
**Causes possibles :**
- Site pas majoritairement en français
- Contenu trop court (< 500 mots)
- Service Python sémantique arrêté

**Solutions :**
```bash
# Vérifier service sémantique
curl http://localhost:5000/health

# Redémarrer si nécessaire
make restart-semantic
```

### Q: Erreur "Lighthouse failed"
```bash
# Réinstaller Lighthouse
npm install -g lighthouse
npm install -g puppeteer

# Vérifier installation
lighthouse --version
```

### Q: "No French content detected"
Le site doit avoir au moins 80% de contenu français. Sites multilingues supportés en mode mixte.

## Erreurs courantes

### 🔴 "connection refused"
```bash
# Vérifier que tous les services sont actifs
make status

# Redémarrer si nécessaire
make restart
```

### 🔴 "Python service unavailable"
```bash
# Vérifier requirements Python
cd internal/semantic/python
pip install -r requirements.txt

# Démarrer manuellement
python semantic_server.py
```

### 🔴 "Config file not found"
```bash
# Vérifier fichiers config
ls config/*.yaml

# Recréer si manquants
make setup-config
```

### 🟡 "Partial results only"
Normal si :
- Site très lent (timeout crawler)
- Service sémantique indisponible
- Restrictions robots.txt strictes

Le rapport partiel reste utilisable.

## Performance

### Q: L'audit est très lent
**Optimisations :**
```yaml
# config/crawler.yaml  
performance:
  concurrent_requests: 10  # Augmenter si serveur puissant
  request_timeout: 15s     # Réduire pour sites rapides
```

### Q: Consommation mémoire élevée
**Solutions :**
- Réduire max_urls à 100-200
- Utiliser DistilCamemBERT au lieu de CamemBERT
- Redémarrer périodiquement les services

```yaml
# config/semantic.yaml
model:
  name: "distilcamembert-base"  # Plus léger
  quantization: true            # Réduit la mémoire
```

## Configuration

### Q: Comment personnaliser les règles SEO ?
Éditer `config/tech_rules.yaml` :
```yaml
title:
  min_length: 30        # Titre minimum
  max_length: 60        # Titre maximum
  missing_severity: high
```

### Q: Comment ajouter des mots vides ?
Éditer `config/stopwords_fr.txt` - un mot par ligne.

### Q: Comment changer les seuils de performance ?
```yaml
# config/lighthouse.yaml
thresholds:
  fcp: 1800    # First Contentful Paint
  lcp: 2500    # Largest Contentful Paint
  fid: 100     # First Input Delay
```

## Rapports

### Q: Comment interpréter le score global ?
- 🟢 **85-100** : Excellent
- 🟡 **70-84** : Bon  
- 🟠 **50-69** : Moyen
- 🔴 **0-49** : Critique

### Q: Que signifient les mots-clés suggérés ?
Mots-clés extraits par IA française avec :
- **Score** : Pertinence commerciale (0-1)
- **Fréquence** : Occurrences dans le contenu
- **Confiance** : Fiabilité de la suggestion

### Q: Comment prioriser les recommandations ?
Ordre de priorité :
1. **Critiques** (impact SEO majeur)
2. **Importantes** (amélioration visible)
3. **Mineures** (optimisations fines)

## Support

### Problème non résolu ?
1. Consulter [common-issues.md](common-issues.md)
2. Vérifier les logs : `make logs`
3. Créer une issue GitHub avec :
   - URL testée
   - Message d'erreur
   - Configuration utilisée
   - Logs pertinents