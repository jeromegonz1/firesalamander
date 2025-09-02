# Scénarios utilisateur BDD - Fire Salamander

## Scénario 1: Audit complet nominal (Happy Path)
**Audit ID**: FS-TEST-001

GIVEN un consultant SEO connecté à Fire Salamander
  AND le site "https://camping-test.fr" est accessible
  AND robots.txt autorise le crawl
WHEN il entre l'URL et clique "Analyser"
THEN le système doit:
  1. Générer audit_id "FS-TEST-001"
  2. Afficher progression: "Crawl en cours... 0/300 pages"
  3. Crawler avec ces étapes visibles:
     - "Analyse robots.txt" ✓
     - "Lecture sitemap.xml" ✓
     - "Crawl: 47 pages trouvées" ✓
  4. Lancer analyses:
     - "Lighthouse: 47/47 pages" ✓
     - "Analyse sémantique FR" ✓
  5. Générer rapport:
     - PDF: FS-TEST-001-report.pdf
     - Taille: < 5MB
     - Temps total: < 5 minutes

AND le rapport doit contenir:
  - Executive summary (1 page)
  - Score global SEO: 72/100
  - Top 5 erreurs critiques
  - 20+ suggestions mots-clés FR
  - Graphique maillage interne

## Scénario 2: Site avec restrictions robots.txt
**Audit ID**: FS-TEST-002

GIVEN un site avec robots.txt restrictif:
```
User-agent: *
Disallow: /admin/
Disallow: /private/
Crawl-delay: 2
```
WHEN le crawler analyse ce site
THEN il doit:
1. Logger: "FS-TEST-002: Restrictions robots.txt détectées"
2. Exclure /admin/* et /private/*
3. Respecter delay 2 secondes entre requêtes
4. Générer rapport avec warning:
   "⚠️ Audit partiel: X pages exclues par robots.txt"
5. Continuer avec les pages autorisées

## Scénario 3: Échec partiel avec fallback
**Audit ID**: FS-TEST-003

GIVEN un audit en cours FS-TEST-003
WHEN Lighthouse timeout sur 5 pages
THEN le système doit:
1. Logger erreurs dans /audits/FS-TEST-003/errors.log
2. Marquer ces pages "technical_analysis": "failed"
3. Continuer les autres pages
4. Générer rapport avec section:
   "Analyse technique incomplète: 42/47 pages analysées"
5. Calculer scores sur les pages réussies uniquement