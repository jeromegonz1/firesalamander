# 🔥 Fire Salamander SEO Analyzer

**Fire Salamander Agency Edition** - Outil SEO interne pour SEPTEO

![Fire Salamander](https://img.shields.io/badge/Fire%20Salamander-🦎-orange)
![SEPTEO](https://img.shields.io/badge/Powered%20by-SEPTEO-ff6136)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8)
![Status](https://img.shields.io/badge/Status-In%20Development-yellow)

## 🎯 Description

Fire Salamander est un outil d'analyse SEO avancé développé pour l'agence SEPTEO. Il combine analyse sémantique locale et intelligence artificielle pour fournir des insights SEO précis tout en maîtrisant les coûts.

### ✨ Caractéristiques

- 🕷️ **Crawling Intelligent** - Exploration respectueuse avec cache optimisé
- 🧠 **Analyse Hybride** - Sémantique locale + enrichissement IA sélectif  
- 📊 **Audit SEO Complet** - Technical SEO, contenu, performance
- 📈 **Rapports Visuels** - Templates SEPTEO avec insights IA
- ⚡ **Performance** - Architecture Go optimisée, <40€/mois
- 🎨 **Brand SEPTEO** - Design system et identité visuelle intégrés

## 🚀 Installation

### Développement Local (Mac)

```bash
# Cloner le repository
git clone https://github.com/jeromegonz1/firesalamander.git
cd firesalamander

# Démarrer avec Docker Compose
docker-compose up

# Ou démarrer directement
go run main.go
```

### Production (Infomaniak)

```bash
# Configuration initiale du serveur
./deploy/setup-infomaniak.sh

# Déploiement
export DEPLOY_HOST=your-server.com
export DEPLOY_USER=your-user  
./deploy/deploy.sh
```

## 📋 Configuration

### Variables d'environnement

```bash
# Développement
ENV=development

# Production
ENV=production
DB_NAME=firesalamander
DB_USER=firesalamander
DB_PASS=your-password
OPENAI_API_KEY=your-openai-key
```

### Fichiers de configuration

- `config/config.dev.yaml` - Configuration développement
- `config/config.prod.yaml` - Configuration production

## 🏗️ Architecture

```
firesalamander/
├── main.go                 # Point d'entrée
├── config/                 # Configuration multi-env
├── crawler/               # Module crawling
├── semantic/              # Analyse sémantique hybride
├── seo/                   # Audit SEO technique
├── reports/               # Génération rapports
├── api/                   # API REST
├── web/                   # Interface HTMX
└── deploy/                # Scripts déploiement
```

## 🛠️ Stack Technique

- **Backend**: Go 1.22+
- **Frontend**: HTML/HTMX (sans build)
- **Database**: SQLite (dev) / MySQL (prod)
- **Cache**: In-memory (bigcache)
- **Queue**: Embedded (asynq)
- **AI**: OpenAI GPT-3.5 (usage optimisé)

## 📊 Fonctionnalités

### Phase 1 ✅ - Setup Initial
- [x] Structure projet Go
- [x] Configuration multi-environnements
- [x] Serveur HTTP avec branding SEPTEO
- [x] Docker Compose développement
- [x] Scripts déploiement Infomaniak

### Phase 2 🚧 - Module Crawler
- [ ] Fetcher HTTP optimisé
- [ ] Parser robots.txt et sitemaps
- [ ] Cache intelligent
- [ ] Rate limiting configurable

### Phase 3 📋 - Analyse Sémantique
- [ ] Extraction n-grammes français
- [ ] Scoring TF-IDF local
- [ ] ML local avec apprentissage
- [ ] Interface OpenAI sélective

### Phase 4 📋 - Analyse SEO
- [ ] Technical SEO (title, meta, headings)
- [ ] Analyse contenu et mots-clés
- [ ] Performance (Core Web Vitals)
- [ ] Scoring intelligent

### Phase 5 📋 - Rapports
- [ ] Templates HTML avec design SEPTEO
- [ ] Sections modulaires
- [ ] Export PDF (optionnel)
- [ ] Storage et historique

### Phase 6 📋 - API & Interface
- [ ] API REST endpoints
- [ ] Interface HTMX
- [ ] Dashboard sites
- [ ] Lancement analyses

### Phase 7 📋 - Déploiement
- [ ] Monitoring et logs
- [ ] Backup automatique
- [ ] CI/CD GitHub
- [ ] Production Infomaniak

## 🌐 Endpoints

```bash
GET  /                     # Dashboard principal
GET  /health              # Health check
POST /api/sites           # Ajouter site
GET  /api/sites           # Liste sites  
POST /api/sites/:id/analyze # Lancer analyse
GET  /api/reports/:id     # Récupérer rapport
```

## 💰 Coûts

- **Serveur**: Infomaniak (déjà payé)
- **OpenAI**: ~2€/mois (usage optimisé)
- **Total**: <40€/mois

## 🔧 Commandes Utiles

```bash
# Développement
go run main.go
docker-compose up

# Tests
go test ./...

# Build
go build -o firesalamander

# Déploiement
./deploy/deploy.sh
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/jeromegonz1/firesalamander/issues)
- **Contact**: Équipe SEPTEO

## 📄 Licence

Propriété SEPTEO - Usage interne uniquement

---

**🦎 Fire Salamander** - Propulsé par **SEPTEO**