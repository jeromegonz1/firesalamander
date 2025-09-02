# Fire Salamander - Interface Web MVP

Interface web moderne pour la plateforme d'audit SEO Fire Salamander.

## 🎯 Fonctionnalités

### ✅ Implémentées (MVP)
- **Formulaire d'audit** : Soumission d'audits SEO avec configuration
- **Suivi temps réel** : Progression et statut des agents en direct
- **Visualisation des résultats** : Onglets organisés par type d'analyse
- **Export des rapports** : JSON et HTML
- **Historique local** : Sauvegarde des audits précédents
- **Interface responsive** : Compatible mobile et desktop
- **API REST simulée** : Endpoints pour développement

### 🔄 Statuts des agents affichés
1. **Crawler** - Exploration du site 🕷️
2. **Keyword Extractor** - Analyse des mots-clés 🔍
3. **Technical Auditor** - Audit technique ⚙️
4. **Linking Mapper** - Cartographie des liens 🔗
5. **Broken Links Detector** - Détection liens brisés 🚫

## 🏗️ Architecture

```
web/
├── index.html              # Page principale
├── static/
│   ├── css/
│   │   └── styles.css      # Styles CSS modernes
│   └── js/
│       └── app.js          # Application JavaScript
├── templates/              # Templates (future expansion)
└── README.md              # Cette documentation
```

## 🚀 Démarrage rapide

### Serveur de développement
```bash
# Démarrer le serveur web
cd cmd/webserver
go run main.go

# Ou avec un port spécifique
PORT=3000 go run main.go
```

### Accès
- Interface web : http://localhost:8080
- API : http://localhost:8080/api/v1/*

## 📡 API REST (MVP)

### Démarrer un audit
```http
POST /api/v1/audits
Content-Type: application/json

{
  "siteUrl": "https://example.com",
  "auditType": "complete",
  "maxPages": 10,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Lister les audits
```http
GET /api/v1/audits
```

### Détails d'un audit
```http
GET /api/v1/audits/{auditId}
```

## 🎨 Interface utilisateur

### Sections principales
1. **Header** - Logo et description
2. **Formulaire d'audit** - Configuration et lancement
3. **Progression** - Statut temps réel des agents
4. **Résultats** - Visualisation par onglets
5. **Historique** - Audits précédents
6. **Footer** - Informations légales

### Onglets des résultats
- **Vue d'ensemble** - Résumé et recommandations
- **Mots-clés** - Analyse des termes trouvés
- **Technique** - Scores performance, accessibilité, SEO
- **Liens** - Statistiques et liens brisés

## 🔧 Technologies

### Frontend
- **HTML5** semantic et accessible
- **CSS3** avec variables custom et animations
- **JavaScript Vanilla ES6+** (pas de framework)
- **Design responsive** mobile-first
- **Local Storage** pour l'historique

### Backend (MVP)
- **Go HTTP Server** simple et performant
- **JSON API** REST standardisée
- **CORS** configuré pour développement
- **Middleware** logging et gestion erreurs

## 🎯 Intégration avec les agents

L'interface web communique avec :

### Agents Feature (implémentés)
- `keyword_extractor` - Extraction mots-clés
- `technical_auditor` - Audit technique complet
- `linking_mapper` - Cartographie liens
- `broken_links_detector` - Détection liens brisés

### Orchestrateur V2 (à intégrer)
- Coordination des agents
- Streaming temps réel
- Gestion d'état distribuée

## 📊 Simulation MVP

Pour le MVP, l'interface utilise :
- **Données simulées** réalistes pour démonstration
- **Progression artificielle** avec timers
- **Résultats générés** basés sur de vrais patterns SEO
- **LocalStorage** pour persistance temporaire

## 🔮 Évolution prochaine

### Sprint suivant
- Intégration avec Orchestrateur V2 réel
- WebSocket pour streaming temps réel
- Authentification utilisateurs
- Sauvegarde serveur des audits
- Dashboard analytics

### Fonctionnalités avancées
- Comparaisons d'audits
- Alertes automatiques
- Intégration CI/CD
- API webhooks
- Rapports programmés

## 🧪 Tests

```bash
# Tests serveur web
cd cmd/webserver
go test -v

# Tests interface (manuel)
# Ouvrir http://localhost:8080 et tester :
# 1. Soumission d'audit
# 2. Suivi progression
# 3. Visualisation résultats
# 4. Export JSON/HTML
# 5. Historique local
```

## 📱 Responsive Design

### Breakpoints
- **Desktop** : > 768px (design principal)
- **Tablet** : 481px - 768px (adaptations)
- **Mobile** : < 480px (simplifications)

### Optimisations mobile
- Navigation par onglets verticale
- Cartes d'agents empilées
- Résumé statistiques en 1 colonne
- Formulaires optimisés tactiles

## 🎨 Design System

### Couleurs
- **Primary** : #2563eb (bleu principal)
- **Success** : #059669 (vert succès)  
- **Warning** : #d97706 (orange attention)
- **Error** : #dc2626 (rouge erreur)
- **Background** : #f8fafc (gris clair)

### Typography
- **Système** : -apple-system, Segoe UI, Roboto
- **Tailles** : 0.8rem → 2.5rem avec ratios harmonieux
- **Poids** : 400 (normal), 500 (medium), 600 (semi-bold)

### Espacement
- **Base** : 1rem (16px)
- **Échelle** : 0.25rem, 0.5rem, 1rem, 1.5rem, 2rem
- **Système** : CSS custom properties cohérent

---

🦎 **Fire Salamander MVP** - Interface web moderne pour audits SEO avec agents spécialisés