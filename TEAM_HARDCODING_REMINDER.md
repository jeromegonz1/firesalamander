# 📢 RAPPEL IMPORTANT À TOUTE L'ÉQUIPE FIRE SALAMANDER

## ❌ POURQUOI LE HARDCODING EST NOTRE ENNEMI #1

### 1. 🔧 **DIFFICILE À MAINTENIR**
- Si une valeur change, il faut modifier le code source → **risque d'erreur**
- Exemple concret : Changer le timeout de 30s à 60s nécessite :
  - Trouver TOUTES les occurrences dans le code
  - Recompiler
  - Redéployer
  - Risquer d'en oublier une → **BUG EN PRODUCTION**

### 2. 🌍 **PAS PORTABLE**
- Une valeur adaptée à un environnement (dev, staging, prod) peut ne pas l'être à un autre
- Exemples critiques :
  - `localhost:8080` en dev → `api.septeo.com` en prod
  - Timeout 5s en dev → 30s en prod (réseau plus lent)
  - 10 pages max en test → 1000 pages en prod

### 3. 🔐 **SÉCURITÉ COMPROMISE**
- Les secrets (clés API, tokens) codés en dur peuvent être exposés dans le repo Git
- **RAPPEL** : Notre code est versionné, TOUT l'historique est accessible
- Un secret commité = secret compromis À VIE
- Exemple : `APIKey = "sk-1234567890"` → **FUITE DE DONNÉES**

### 4. 🚫 **MOINS FLEXIBLE**
- Impossible de personnaliser le comportement sans changer le code
- Le client veut 50 pages au lieu de 20 ? → Recompilation nécessaire
- Besoin d'ajuster les seuils SEO ? → Modification du code
- Test A/B impossible sans déployer 2 versions

---

## 📊 IMPACT RÉEL SUR FIRE SALAMANDER

### Violations actuelles : **257** 😱

Cela signifie :
- **257 points de maintenance** potentiels
- **257 risques de bugs** lors des changements
- **257 blocages** pour la configuration client
- **257 violations** de nos standards SEPTEO

### Exemples concrets dans notre code :

```go
// ❌ MAUVAIS
simulatePhase(analysisID, "Découverte des pages...", 0, 30, 500*time.Millisecond, ...)

// ✅ BON
simulatePhase(analysisID, cfg.Messages.Discovery, cfg.Phases.Discovery.Start, cfg.Phases.Discovery.End, cfg.Phases.Discovery.Interval, ...)
```

```go
// ❌ MAUVAIS
return "1-2 minutes"

// ✅ BON
return cfg.Messages.EstimatedTime.OneToTwo
```

```html
<!-- ❌ MAUVAIS -->
<script src="https://cdn.tailwindcss.com"></script>

<!-- ✅ BON -->
<script src="{{.Config.CDN.TailwindURL}}"></script>
```

---

## 🎯 RÈGLES D'OR ANTI-HARDCODING

### 1. **LA RÈGLE DES 3 QUESTIONS**
Avant d'écrire une valeur, demandez-vous :
- Cette valeur pourrait-elle changer ? → **CONFIG**
- Cette valeur est-elle différente selon l'environnement ? → **ENV**
- Cette valeur est-elle sensible ? → **SECRET**

### 2. **TYPES DE VALEURS À EXTERNALISER**

| Type | Exemples | Solution |
|------|----------|----------|
| **URLs** | API endpoints, CDNs | `.env` ou `config.yaml` |
| **Timeouts** | 30s, 1m, 500ms | `config.yaml` |
| **Limites** | Max pages, taille fichiers | `config.yaml` |
| **Messages** | Erreurs, UI, logs | `messages.yaml` |
| **Chemins** | Templates, assets | `.env` |
| **Ports** | 8080, 3000 | `.env` |
| **Secrets** | API keys, tokens | `.env` (jamais committé) |

### 3. **HIÉRARCHIE DE CONFIGURATION**

```
1. Variables d'environnement (.env) - PRIORITÉ MAX
   └── Surcharge tout le reste
   
2. Fichiers de configuration (config.yaml)
   └── Valeurs par défaut de l'application
   
3. Valeurs par défaut dans le code
   └── UNIQUEMENT si vraiment nécessaire
```

---

## 🛠️ OUTILS À VOTRE DISPOSITION

### 1. **Test Anti-Hardcoding**
```bash
# À lancer AVANT chaque commit
go test ./internal/qa -run TestNoHardcoding

# Si le test échoue → CORRIGEZ avant de continuer
```

### 2. **Configuration Centralisée**
```go
// Utilisez TOUJOURS
cfg := config.Load()
cfg.Simulation.Phases.Discovery.Duration

// Au lieu de
500 * time.Millisecond
```

### 3. **Messages Externalisés**
```go
// Utilisez
cfg.Messages.Errors.InvalidURL

// Au lieu de
"Invalid URL format"
```

---

## 📝 CHECKLIST AVANT CHAQUE COMMIT

- [ ] J'ai lancé `go test ./internal/qa -run TestNoHardcoding`
- [ ] Aucune URL n'est hardcodée
- [ ] Aucun timeout/durée n'est hardcodé
- [ ] Aucun message n'est hardcodé
- [ ] Aucun port/chemin n'est hardcodé
- [ ] Aucune limite/seuil n'est hardcodé
- [ ] Les valeurs sont dans `.env` ou `config.yaml`

---

## 🎖️ ENGAGEMENT D'ÉQUIPE

En tant que membre de l'équipe Fire Salamander, je m'engage à :

1. **NE JAMAIS** hardcoder de valeurs
2. **TOUJOURS** externaliser dans la configuration
3. **VÉRIFIER** avec le test anti-hardcoding
4. **REFUSER** les PR avec du hardcoding
5. **ÉDUQUER** mes collègues sur ces pratiques

---

## 💡 RAPPEL FINAL

**Le hardcoding n'est pas une question de "rapidité" mais de DETTE TECHNIQUE.**

Chaque valeur hardcodée aujourd'hui = 
- 🐛 Un bug potentiel demain
- 💰 Des heures de maintenance
- 😤 Un client frustré
- 📉 Une réputation ternie

**ENSEMBLE, construisons un Fire Salamander MAINTENABLE, FLEXIBLE et PROFESSIONNEL !**

---

*Document créé par Claude Code, Architecte Principal*
*Date : 2025-08-07*
*À lire et appliquer par TOUS les membres de l'équipe*