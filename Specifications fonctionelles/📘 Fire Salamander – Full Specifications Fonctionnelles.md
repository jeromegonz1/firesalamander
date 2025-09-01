📘 Fire Salamander – Spécifications Fonctionnelles Complètes
Ce document regroupe l’ensemble des spécifications fonctionnelles des agents du projet Fire Salamander. Chaque section correspond à un agent ou un composant du pipeline d’audit SEO et reprend intégralement les spécifications détaillées fournies séparément. Le but est de disposer d’un référentiel unique, versionnable et facile à consulter pour tous les contributeurs du projet.
Sommaire
1.	Agent Crawler
2.	Agent Audit Technique
3.	Agent Analyse Sémantique
4.	Agent Reporting
5.	Agent Orchestrateur
 
🕸 Spécification Fonctionnelle — Agent Crawler
Cette spécification décrit le fonctionnement détaillé de l'agent Crawler de Fire Salamander. Son rôle est d'explorer un domaine web donné, de collecter les données nécessaires à l'audit technique et sémantique, et de produire un index structuré des pages. Le design s'inspire du CDC V4.1 et de la phase 1 de la roadmap.
Objectifs
1.	Explorer intelligemment un site en respectant les règles d’accès (robots.txt) et les contraintes définies par l’utilisateur (profondeur maximale, nombre maximal d’URLs, respect du domaine).
2.	Collecter les données nécessaires aux modules technique et sémantique : URL, langue détectée, balises clés (title, h1–h3), ancres de liens internes, texte principal, meta‐données (canonical, index, nofollow, hreflang), profondeur et maillage interne (liens entrants/sortants).
3.	Garantir la performance et la scalabilité en implémentant un crawl concurrent et un cache pour éviter les re‐fetchs intempestifs.
4.	Fournir un format de sortie standardisé (crawl_index.json) qui servira d’entrée aux autres agents (Audit Technique, Analyse Sémantique).
Entrées
•	Seed URL : domaine à auditer (par ex. https://www.monsite.com).
•	Paramètres d’audit :
•	max_urls : nombre d’URLs maximal (300 par défaut, configurable).
•	max_depth : profondeur maximale de crawl (3 par défaut). La profondeur est calculée à partir de la page d’accueil (niveau 0).
•	respect_robots : booléen déterminant si le crawler tient compte de robots.txt (par défaut true).
•	respect_sitemap : booléen indiquant si le crawler doit d’abord consulter sitemap.xml pour déterminer les pages prioritaires.
•	sampling_strategy : option pour limiter le crawl aux pages « types » (home, catégories, produits phares) quand le site dépasse le nombre maximal d’URLs.
•	language_target : code langue ISO 639 (« fr » pour les analyses sémantiques). Les pages qui ne correspondent pas à cette langue sont collectées mais ne seront pas envoyées au module sémantique.
•	concurrent_requests, request_timeout, retry_attempts, cache_ttl : paramètres de performance (voir § Optimisation).
Comportement et algorithme
1.	Initialisation :
2.	Lire robots.txt (si respect_robots = true) et exclure les chemins interdits.
3.	Récupérer sitemap.xml et ajouter ses URLs à la file de crawl selon l’ordre de priorité (changefreq et priority).
4.	Créer une file de priorité (par défaut BFS) et un ensemble de pages visitées.
5.	Exploration :
6.	Tant que la file n’est pas vide et que max_urls n’est pas atteint :
o	Extraire l’URL actuelle et récupérer la page via HTTP (respect des timeouts et retries).
o	Vérifier la langue de la page (détection par lang ou via un détecteur type langid). Stocker la langue détectée.
o	Analyser le HTML : extraire title, h1–h3, liens internes (<a href>), ancres (texte d’ancre), canonical, attributs rel, meta robots, hreflang, etc.
o	Normaliser les URLs (suppression des fragments, paramètres non pertinents) et s’assurer qu’elles restent dans le domaine.
o	Évaluer la profondeur : si elle dépasse max_depth, ignorer les liens plus profonds.
o	Ajouter les liens valides à la file de crawl (selon la stratégie BFS ou sampling) en évitant les doublons.
o	Construire un enregistrement JSON par page avec tous les champs extraits.
7.	Sampling pour sites volumineux :
8.	Si le nombre de pages explorées dépasse max_urls, appliquer une stratégie de réduction :
o	Inclure systématiquement la page d’accueil.
o	Parcourir les pages listées dans sitemap.xml en priorité.
o	Sélectionner un échantillon de pages par type : catégories principales, produits phares, pages à forte profondeur (pour détecter des failles d’indexation).
9.	Sortie :
10.	Générer un fichier crawl_index.json dans /audits/{audit_id}/ contenant une liste de documents :
{
  "pages": [
    {
      "url": "https://www.monsite.com/",
      "lang": "fr",
      "title": "Camping Nature – Home",
      "h1": "Bienvenue au Camping Nature",
      "h2": ["Nos services", "Nos emplacements"],
      "anchors": [{"text": "Mobil‑home", "href": "/mobil-homes"}],
      "canonical": "https://www.monsite.com/",
      "meta_index": true,
      "depth": 0,
      "outgoing_links": ["/mobil-homes", "/services"],
      "incoming_links": []
    },
    …
  ]
}
Chaque objet page peut contenir des champs supplémentaires (meta description, H3, images, etc.) selon les besoins de l’audit technique.
Optimisation et limitations
•	Performance : pour éviter la surcharge réseau, le nombre de requêtes simultanées est limité (concurrent_requests = 5). Le request_timeout est fixé à 10 secondes et deux retries sont autorisés. Les pages déjà visitées sont stockées dans un cache (TTL = 1 h) pour éviter les re‐fetchs.
•	Respect des règles : le crawler respecte robots.txt et sitemap.xml (sauf configuration contraire). Il ne suit pas les liens externes ni les sous‑domaines qui ne correspondent pas au domaine de base.
•	Langue : seules les pages en français (score de confiance ≥ 0,8) sont envoyées au module sémantique. Les autres sont uniquement analysées techniquement.
•	Limites : pour des sites très volumineux, le crawling complet peut être long. Le sampling permet de réduire l’exploration et de se concentrer sur des pages représentatives.
Interface JSON‑RPC et contrat
•	Requête (exemple) :
{
  "jsonrpc": "2.0",
  "method": "crawl",
  "id": "audit_id",
  "params": {
    "seed_url": "https://www.monsite.com",
    "max_urls": 300,
    "max_depth": 3,
    "respect_robots": true,
    "respect_sitemap": true,
    "language_target": "fr"
  }
}
•	Réponse partielle (streaming) : chaque page explorée peut être envoyée au fur et à mesure (pour l’UI) :
{
  "jsonrpc": "2.0",
  "id": "audit_id",
  "result": {
    "status": "page_found",
    "data": {
      "url": "https://www.monsite.com/mobil-homes",
      "depth": 1,
      "lang": "fr"
    }
  }
}
•	Réponse finale :
{
  "jsonrpc": "2.0",
  "id": "audit_id",
  "result": {
    "status": "complete",
    "data": {
      "crawl_index": "/audits/audit_id/crawl_index.json"
    },
    "metadata": {
      "pages_total": 240,
      "depth_max": 3,
      "duration_ms": 48250
    }
  }
}
Tests et validation (TDD)
1.	Unitaires :
2.	Test de normalisation d’URL (suppression des fragments et paramètres, gestion du trailing slash).
3.	Test de respect de robots.txt (URL non crawlées si disallowed).
4.	Test de detection de langue (FR vs non‑FR).
5.	Test de gestion des depth (pas de liens ajoutés si profondeur > max_depth).
6.	Test de caching (pages visitées ne sont pas réexplorées).
7.	Contrats : valider la structure du crawl_index.json via un schéma JSON et s’assurer que tous les champs obligatoires sont présents.
8.	Intégration : simuler un crawl sur un petit site de test avec sitemap.xml et comparer le nombre de pages trouvées avec les attentes. Vérifier la cohérence des liens entrants/sortants et la profondeur.
9.	Performance : mesurer le temps d’exploration d’un site moyen (ex. 100 pages) et valider que la durée est inférieure au seuil défini (2 s par page p95). Tester le comportement avec sampling.
Le Crawler constitue la première étape du pipeline Fire Salamander et doit être stable, performant et extensible. Ses sorties alimentent directement les modules d’audit technique et sémantique.
 
🧪 Spécification Fonctionnelle — Agent Audit Technique
Cette spécification décrit l’agent Audit Technique de Fire Salamander. Son rôle est d’évaluer l’état technique et SEO d’un site web sur la base des pages collectées par l’agent Crawler. Il repose sur des outils d’audit automatisé (Lighthouse) et des heuristiques personnalisées pour détecter les problèmes SEO courants et analyser la qualité du maillage interne.
Objectifs
1.	Mesurer la performance et la conformité SEO de chaque page grâce à des métriques standardisées : performance, SEO, accessibilité et bonnes pratiques.
2.	Identifier les erreurs et les avertissements techniques : balises manquantes ou mal configurées, structure HTML incorrecte, absence de HTTPS, redirections excessives, images sans attribut alt, duplication de balises meta, titres trop courts ou trop longs, etc.
3.	Analyser le maillage interne : repérer les pages orphelines, mesurer la profondeur moyenne, identifier les ancres pauvres (ex. « cliquez ici »), calculer le ratio de liens internes vs externes, et fournir des recommandations pour équilibrer le “link juice”.
4.	Fournir des recommandations actionnables à partir des résultats d’audit, avec une classification par gravité (critique, élevée, moyenne, faible) pour aider le consultant SEPTEO à prioriser les actions.
5.	Exporter un format de résultats standardisé (tech_result.json) pour être consommé par l’agent Reporting et l’interface utilisateur.
Entrées
•	crawl_index.json : liste des pages collectées par l’agent Crawler, incluant URL, profondeur, langue, balises et maillage interne.
•	Paramètres :
•	device : desktop ou mobile (détermine la configuration Lighthouse).
•	max_concurrent_audits : nombre de pages analysées en parallèle (par défaut 4). Les audits peuvent être gourmands en CPU/GPU.
•	tech_rules.yaml : fichier de configuration regroupant les seuils pour les vérifications (longueur du titre, description, poids des balises, etc.) et la correspondance avec les niveaux de gravité.
Comportement et algorithme
1.	Prise en charge des pages : pour chaque URL en langue cible (mais on analyse aussi les autres pages pour l’aspect technique) :
2.	Lancer Google Lighthouse CLI (via Node.js ou Docker) avec les options : --chrome-flags="--headless" et la configuration device.
3.	Récupérer le rapport JSON Lighthouse. Extraire les scores : performance, accessibility, best-practices, seo.
4.	Convertir les audits (p.ex. link-name, meta-description, canonical) en un format interne en utilisant les seuils définis dans tech_rules.yaml pour déterminer la gravité :
o	Ex. : un titre < 15 caractères est de gravité moyenne ; une absence de balise title est critique.
5.	Collecter les observations sur les redirections HTTP, la compression, le cache et l’utilisation de HTTPS.
6.	Analyse du maillage interne :
7.	Calculer le graph des liens internes à partir de crawl_index.json : sommets = pages, arêtes = liens. Identifier :
o	Les pages orphelines (aucun lien entrant).
o	La profondeur maximale et la profondeur moyenne.
o	Les ancres pauvres : textes génériques ou sans lien avec la page cible.
8.	Fournir des statistiques agrégées (nombre de pages orphelines, distribution des profondeurs) et une liste d’ancres pauvres avec la page source et la page cible.
9.	Classification et recommandation :
10.	Regrouper les résultats en deux catégories : Findings et Warnings. Chaque item a :
o	id : identifiant (ex. missing-title),
o	severity : critical, high, medium, low,
o	message : description lisible par l’utilisateur,
o	evidence : URLs et lignes pertinentes (ex. page où le titre manque).
11.	Classer les items par gravité et fournir un ordre de priorité.
12.	Agrégation des scores :
13.	Calculer les scores moyens pour le site en pondérant chaque page par son importance (poids selon la profondeur et la popularité si les données sont disponibles).
14.	Présenter ces scores de manière synthétique : par exemple, un SEO score global = moyenne des scores Lighthouse seo multipliée par un facteur de pénalité s’il y a des erreurs critiques non résolues.
15.	Sortie :
16.	Générer le fichier tech_result.json dans /audits/{audit_id}/ :
{
  "audit_id": "string",
  "model_version": "tech-v1.0",
  "scores": {
    "performance": 0.78,
    "accessibility": 0.88,
    "best_practices": 0.85,
    "seo": 0.91
  },
  "findings": [
    {"id": "missing-title", "severity": "high", "message": "Titre manquant.", "evidence": ["https://www.monsite.com/a-propos"]},
    {"id": "slow-page", "severity": "medium", "message": "Temps de chargement > 4 s.", "evidence": ["https://www.monsite.com/contact"]}
  ],
  "warnings": [
    {"id": "multiple-h1", "severity": "low", "message": "Plusieurs balises H1.", "evidence": ["https://www.monsite.com/blog"]}
  ],
  "mesh": {
    "orphans": ["https://www.monsite.com/faq"],
    "depth_stats": {"min": 0, "max": 3, "avg": 1.7},
    "weak_anchors": ["cliquez ici"]
  }
}
Tests et validation (TDD)
1.	Unitaires :
2.	Fonction de mapping des audits Lighthouse → findings et warnings selon tech_rules.yaml.
3.	Détection des titres manquants, descriptions trop longues, H1 multiples, images sans alt, liens cassés.
4.	Construction du graph de maillage interne et identification des pages orphelines.
5.	Tests de contrat : validation du tech_result.json contre le schéma JSON officiel. Vérification de la présence des champs obligatoires et du type correct des valeurs.
6.	Intégration : exécuter un audit complet sur un site de test (10 pages) et comparer les résultats avec une référence annotée. Vérifier que les scores Lighthouse sont bien capturés et que les findings sont complets.
7.	Performance : mesurer la durée moyenne d’analyse par page (< 2 s p95) et la consommation CPU. Ajuster max_concurrent_audits si nécessaire.
8.	Robustesse : tester le comportement en cas de page inaccessible (4xx, 5xx). Vérifier que l’erreur est enregistrée mais que l’audit continue.
Contraintes et sécurité
•	Les audits sont exécutés en sandbox (container ou instance isolée) pour éviter les effets de bord sur l’environnement hôte.
•	Les URLs sont normalisées pour éviter des injections via des paramètres malveillants.
•	Les résultats ne doivent contenir aucune donnée personnelle ou sensible ; seules des URLs publiques et des extraits anonymisés sont stockés.
L’agent Audit Technique fournit ainsi un diagnostic complet et priorisé de l’état d’un site, permettant de cibler rapidement les points d’optimisation SEO et de préparer un reporting clair.
 
🧠 Spécification Fonctionnelle — Agent Analyse Sémantique
Cette spécification détaille l’agent Analyse Sémantique de Fire Salamander. Il a pour mission de comprendre les contenus des pages en français, d’en extraire les thèmes et intentions, et de proposer des expressions longue traîne (n‑grammes de 2 mots ou plus) pertinentes pour améliorer la visibilité SEO et l’alignement avec les recherches utilisateurs[1]. Le moteur sémantique adopte une approche hybride combinant règles heuristiques et machine learning, enrichie par les suggestions d’un modèle de langage (IA) optionnel.
Objectifs
1.	Comprendre le contexte métier : identifier les thèmes majeurs d’un site (topics), les intentions des pages (informationnel, commercial, transactionnel, navigationnel) et le type de chaque page (produit, catégorie, article, home).
2.	Extraire des expressions longue traîne (≥ 2 mots) en français, basées sur l’analyse n‑grammes, qui reflètent les recherches des utilisateurs et leur intent[2].
3.	Prioriser les suggestions à l’aide d’un ranker pondéré prenant en compte la pertinence thématique, l’intention cible, la preuve de maillage interne, la lisibilité et le feedback historique.
4.	Permettre l’apprentissage continu grâce au feedback des consultants et aux signaux d’intégration (mots clés intégrés au site), en ajustant automatiquement les pondérations du ranker et, le cas échéant, les hyperparamètres des modèles.
5.	Option IA : proposer, pour des cas premium, des suggestions créatives générées par un modèle de langage (GPT‑3.5 ou Mistral 7B local) afin d’enrichir le vocabulaire au-delà des seules combinaisons détectées, tout en préservant la confidentialité des données.
Entrées
•	semantic_request.json (voir contrat) : identifiant de l’audit, langue (fr), nombre maximal de candidats (max_candidates) et nombre final de suggestions (top_n).
•	crawl_index.json : liste des pages en langue française avec le texte principal, les ancres internes, la profondeur et les métadonnées.
•	config/semantic.yaml : paramètres du moteur (listes de stopwords, patterns POS, poids du ranker, seuils de langue, etc.).
•	banlists/allowlists : listes de mots ou expressions à exclure ou à inclure dans l’analyse (marques, termes trop génériques, jurons, etc.).
Pipeline de traitement
1.	Pré‑traitement linguistique :
2.	Détection de langue : exclure les pages dont la confiance < 0,8.
3.	Nettoyage HTML : extraire le texte principal (suppression du code, navigation, footer). Utiliser des extracteurs basés sur des heuristiques DOM.
4.	Normalisation : passage en minuscules, suppression des accents, lemmatisation avec spaCy fr_core_news_md ou équivalent, suppression des stopwords.
5.	Extraction n‑grammes (2–5) : découper le texte en séquences de mots (bigrammes, trigrammes, etc.) et compter leurs occurrences. Les n‑grammes de plus grande longueur correspondent souvent à des recherches longue traîne, attirant un trafic ciblé[2].
6.	Filtre qualité : éliminer les expressions trop génériques (via banlists), celles avec une densité trop faible ou trop élevée et celles ne respectant pas la syntaxe FR (ex. stopword final).
7.	Compréhension métier et thématique :
8.	Embeddings multilingues (CamemBERT/DistilCamemBERT ou FlauBERT) pour représenter chaque page et chaque n‑gramme dans un espace vectoriel.
9.	Clustering thématique (HDBSCAN ou BERTopic) pour grouper les pages et les n‑grammes par thèmes et détecter des topics (ex. « camping nature », « mobil‑home haut de gamme »).
10.	Intent classifier : modèle (logistic regression ou MLP) entraîné sur un jeu de pages annotées (info/comm/transac/nav) pour prédire l’intention de chaque page et de chaque n‑gramme.
11.	Type de page : détection via règles DOM (présence de prix, de formulaires) et features sémantiques pour distinguer produit, catégorie, article, etc.
12.	Génération de candidates :
13.	Fusion n‑grammes + ancres internes : les ancres fournissent un signal fort de pertinence métier. Les n‑grammes contenus dans les ancres sont priorisés.
14.	Méthodes heuristiques (TextRank, YAKE, TF‑IDF intra-site) pour extraire des keyphrases supplémentaires.
15.	Variantes morpho‑syntaxiques : combinaison de lemmes et d’adjectifs (ex. « camping familial », « mobil‑home familial ») et permutations limitées.
16.	Limiter les candidats : trier par fréquence et pertinence et retenir au maximum max_candidates (500 par défaut).
17.	Scoring et ranking :
18.	Score thématique : similarité cosinus entre le vecteur du n‑gramme et les vecteurs des topics dominants du site.
19.	Score intent : concordance entre l’intention de la page et l’objectif de la suggestion (bonus si le site vise la conversion et que l’intention est transactionnelle).
20.	Preuve maillage : nombre d’ancres internes contenant ou pointant vers l’expression candidate.
21.	Lisibilité et qualité linguistique : heuristiques pour pénaliser les expressions mal formées, avec des stopwords finaux ou un enchaînement incorrect.
22.	Feedback historique : bonus ou malus selon les évaluations précédentes (👍/👎) et la réutilisation effective des suggestions dans les audits suivants.
23.	Score final : combinaison pondérée (les poids sont définis dans config/semantic.yaml). Les candidates sont triées par ce score.
24.	Diversité : appliquer un algorithme de diversification (ex. Maximal Marginal Relevance) pour éviter les suggestions redondantes.
25.	Intégration IA optionnelle :
26.	Prompting : générer des prompts pour un LLM (GPT‑3.5 ou Mistral 7B local) en fournissant les topics, les n‑grammes clés et l’intention cible.
27.	Appel à l’IA : récupérer les suggestions textuelles, les normaliser et les filtrer via les mêmes règles de qualité et de langue.
28.	Fusion : fusionner les suggestions IA avec celles issues du pipeline heuristique ; les classer en fin de liste ou selon un poids IA défini.
29.	Confidentialité : les prompts ne doivent pas contenir d’URL ou d’informations sensibles. Les appels externes doivent respecter la RGPD.
30.	Sortie : création d’un fichier semantic_result.json dans /audits/{audit_id}/ contenant :
{
  "audit_id": "string",
  "model_version": "sem-v1.2",
  "topics": [
    {"id": "t1", "label": "camping nature", "terms": ["camping familial", "mobil-home nature"]},
    …
  ],
  "suggestions": [
    {
      "keyword": "camping familial bord de mer",
      "reason": "topic:t1 + anchor keyword + intent=transactional",
      "confidence": 0.87,
      "evidence": ["https://www.monsite.com/emplacements", "/offres/early-booking"]
    },
    …
  ],
  "metadata": {
    "schema_version": "1.0.0",
    "weights_version": "rank-2025w36",
    "execution_time_ms": 15000,
    "lang": "fr"
  }
}
Tests et validation (TDD)
1.	Unitaires :
2.	Extraction n‑grammes : vérifier que l’on obtient les bigrammes/trigrammes attendus pour un texte donné et que les stopwords sont correctement filtrés.
3.	Détection de langue : tester la fonction de filtrage pour différentes langues.
4.	Scoring : valider les fonctions de calcul du score thématique et du score intent.
5.	Diversification : s’assurer que deux suggestions quasi identiques ne sont pas sélectionnées simultanément.
6.	Contrats : valider que semantic_result.json respecte le schéma JSON (présence des champs topics, suggestions, metadata).
7.	Intégration :
8.	Exécuter le pipeline sur un petit corpus annoté (jeu de 100 sites) et comparer les suggestions avec les expressions préconisées manuellement. Calibrer le ranker afin d’atteindre les critères : Precision@10 ≥ 0.6 et nDCG@10 ≥ 0.7.
9.	Tester l’option IA en mode “premium” : vérifier que le nombre de suggestions ne dépasse pas le plafond, que les prompts sont anonymisés et que le temps d’appel reste acceptable (< 1 s par appel cache compris).
10.	Performance : mesurer le temps total de l’analyse sémantique (< 15 s pour 50 pages) et le nombre de tokens proces¬sés. Mettre en place un cache des embeddings (clé = hash du contenu) pour accélérer les audits récurrents.
11.	Robustesse : tester des pages très courtes ou très longues, des pages sans texte, et des pages multilingues. Vérifier que le moteur sémantique ne plante pas et renvoie une liste vide ou un warning approprié.
Feedback et apprentissage
•	Les consultants peuvent évaluer chaque suggestion (👍/👎) et indiquer une catégorie de refus (too_generic, irrelevant, duplicate). Ce feedback est stocké et utilisé pour ajuster les pondérations du ranker.
•	Les suggestions acceptées et effectivement intégrées dans les sites (vérifiées lors des audits ultérieurs) reçoivent un bonus de confiance. Les suggestions répétitivement rejetées sont pénalisées.
•	Périodiquement (hebdomadaire), un job d’entraînement met à jour le ranker (via régression logistique ou boosting) en utilisant les nouvelles données. La version du modèle est incrémentée et enregistrée dans semantic_result.json.
Contraintes et sécurité
•	RGPD et confidentialité : aucune donnée personnelle n’est envoyée au moteur IA ; les prompts sont anonymisés. Les fichiers d’entrée et de sortie sont conservés 30 jours maximum.
•	Règles de langue FR : l’analyse se limite au français ; des heuristiques de lisibilité pénalisent les expressions mal formées. Les n‑grammes inappropriés (insultes, marques déposées non pertinentes) sont filtrés.
•	Coût et latence : les appels au LLM sont limités (option premium) et mis en cache. La version locale (Mistral 7B) est privilégiée pour réduire les coûts.
L’agent Analyse Sémantique est le cœur de la valeur ajoutée métier de Fire Salamander. En combinant des approches statistiques, des modèles linguistiques et de l’IA générative facultative, il permet de proposer des recommandations adaptées aux sites français et d’améliorer continuellement sa pertinence grâce au feedback humain.
 
📄 Spécification Fonctionnelle — Agent Reporting
L’agent Reporting de Fire Salamander est chargé de transformer les données brutes issues des agents Crawler, Audit Technique et Analyse Sémantique en un document synthétique et compréhensible. Ce rapport est destiné aux consultants SEO de SEPTEO afin de faciliter la lecture des résultats, d’identifier rapidement les priorités et de partager les recommandations avec les clients.
Objectifs
1.	Assembler et synthétiser les sorties des agents (technique, sémantique, maillage) en un rapport cohérent.
2.	Présenter les informations sous plusieurs formes : résumé global (scores et indicateurs clés), détails techniques, suggestions sémantiques, maillage interne, recommandations, suivi historique.
3.	Générer un PDF et une page HTML interactive respectant la charte graphique SEPTEO.
4.	Conserver un historique en versionnant chaque rapport (audit_id et timestamp), accessible depuis l’interface SEPTEO.
Entrées
•	tech_result.json : résultats de l’agent Audit Technique (scores, findings, mesh).
•	semantic_result.json : résultats de l’agent Analyse Sémantique (topics, suggestions, metadata).
•	crawl_index.json : liste des pages et métadonnées (pour référence dans le rapport).
•	audit_request.json : paramètres d’audit (URL, options, mode sémantique).
•	template/ : modèles HTML/CSS destinés à la génération du PDF et de la page interactive. Ces modèles suivent la charte SEPTEO (couleurs, typographie, logo).
Structure du rapport
Le rapport se décompose en sections :
1.	Page de garde :
2.	Titre (« Audit SEO – Fire Salamander »)
3.	Domaine audité et date
4.	Identifiant audit_id
5.	Logo SEPTEO et Fire Salamander
6.	Résumé exécutif :
7.	Score global technique (agrégation des scores Lighthouse)
8.	Nombre de pages analysées et profondeur maximale
9.	Top 3 des erreurs critiques et actions prioritaires
10.	Top 3 des suggestions sémantiques
11.	Sommaire (liens internes vers les sections)
12.	Détails techniques :
13.	Tableau des scores par page (Performance, Accessibilité, Best Practices, SEO)
14.	Liste des findings critiques et élevés (avec explications et URLs)
15.	Analyse du maillage interne : nombre de pages orphelines, distribution de la profondeur, ancres pauvres, suggestions d’amélioration
16.	Graphiques ou diagrammes (ex. donut chart pour la répartition des scores)
17.	Analyse sémantique :
18.	Présentation des topics identifiés (avec keywords représentatifs)
19.	Liste des suggestions top‐N (keyword, raison, confiance, evidence)
20.	Tableau des feedbacks précédents (s’il s’agit d’un audit récurrent)
21.	Mention de l’option IA, si utilisée
22.	Recommandations :
23.	Synthèse des actions prioritaires (technique et contenu)
24.	Checklist d’optimisation SEO (performances, balisage, maillage, contenu)
25.	Liens vers des ressources (guides internes, documentation SEPTEO)
26.	Annexes :
27.	Glossaire des termes techniques
28.	Liste complète des pages auditées (avec statut, profondeur, langue)
29.	Référence de version des modèles utilisés (model_version, schema_version, weights_version)
Génération du PDF
•	L’agent utilise un moteur HTML → PDF (ex. Puppeteer ou wkhtmltopdf) pour convertir un template HTML en PDF.
•	Les templates se trouvent dans templates/report.html et static/report.css et respectent la charte SEPTEO.
•	Les images (logo, icônes) sont intégrées en base64 ou via des liens locaux.
•	La pagination est automatique et une table des matières cliquable est générée.
•	Chaque rapport est enregistré sous /audits/{audit_id}/report.pdf et /audits/{audit_id}/report.html.
•	Le nom de fichier inclut l’identifiant et la date (FSAUDIT-0001_20250901.pdf).
Génération de la page HTML interactive
•	En plus du PDF, une page HTML interactive est générée pour la consultation en ligne.
•	Cette page utilise des composants dynamiques (tableaux filtrables, accordéons) et reprend la structure du PDF.
•	Les données sont chargées via des fichiers JSON (tech_result.json, semantic_result.json) et transformées côté client (utilisation d’Alpine.js ou d’un framework léger).
•	Les utilisateurs peuvent trier, filtrer et télécharger des sections spécifiques (ex. export du tableau des suggestions en CSV).
Interface JSON‑RPC
•	Requête : l’agent Orchestrateur envoie une commande generate_report avec l’identifiant de l’audit.
•	Réponse : l’agent Reporting renvoie un chemin vers les fichiers PDF et HTML une fois la génération terminée, ainsi que des métadonnées (durée, succès/échec).
Tests et validation (TDD)
1.	Unitaires :
2.	Fonction de fusion des données (tech + sémantique) : vérifier que les valeurs agrégées sont correctes.
3.	Génération de la table des matières : tester les ancres et les liens internes.
4.	Formatage des dates, numérotation des audits.
5.	Tests d’intégration :
6.	Générer un rapport complet avec des données de test et vérifier la présence de chaque section et l’absence de contenus manquants.
7.	Vérifier que la pagination du PDF est correcte et que les images se chargent.
8.	Contrôler la conformité du HTML (validation W3C) et l’accessibilité (rôles ARIA, contraste des couleurs).
9.	Tests de performance :
10.	Mesurer la durée de génération du PDF (< 15 s pour un audit moyen). Optimiser en utilisant un moteur de rendu headless.
11.	Tests de régression :
12.	Comparer un rapport généré avec une version de référence pour s’assurer que les changements de template ou de code ne brisent pas la structure (snapshot testing).
Contraintes et sécurité
•	Confidentialité : les rapports ne doivent pas contenir d’informations sensibles ou privées. Les données anonymes (par exemple https://www.monsite.com) sont utilisées dans les captures.
•	Accessibilité : l’interface interactive doit respecter les standards WCAG 2.1 (navigation clavier, contraste). Le PDF doit être consultable sur tous les lecteurs.
•	Stockage : les rapports sont conservés pendant 30 jours (configurable) pour respecter la politique de rétention de SEPTEO.
L’agent Reporting est la dernière étape du pipeline Fire Salamander. Il transforme des résultats techniques et sémantiques bruts en un support clair et actionnable, facilitant la prise de décision pour les consultants SEO.
 
🔧 Spécification Fonctionnelle — Agent Orchestrateur
L’agent Orchestrateur est le chef d’orchestre du pipeline Fire Salamander. Il reçoit les requêtes d’audit, coordonne l’exécution des agents spécialisés (Crawler, Audit Technique, Analyse Sémantique, Reporting), gère les flux JSON‑RPC en streaming et assure la persistance des données et des logs. Ce document détaille son comportement, son interface et ses contraintes.
Objectifs
1.	Centraliser la gestion des audits : recevoir les demandes d’audit, créer des identifiants uniques (audit_id), stocker les paramètres et initialiser l’environnement (répertoires, fichiers de configuration).
2.	Coordonner l’exécution des agents : déclencher le Crawler, puis l’Audit Technique et l’Analyse Sémantique (en parallèle lorsque possible) et enfin le Reporting.
3.	Gérer les flux de données : orchestrer le streaming JSON‑RPC pour permettre un suivi en temps réel dans l’interface utilisateur et relayer les sorties intermédiaires.
4.	Assurer la résilience : suivre l’état de chaque étape, gérer les erreurs, permettre l’annulation ou la reprise d’un audit, et enregistrer les logs pour la traçabilité.
5.	Exposer une API accessible aux clients (interface web ou autres services) pour déclencher, suivre et récupérer les audits.
Entrées
•	audit_request.json : { seed_url, max_urls, max_depth, lang, modes, options } (défini dans le contrat audit_request.schema.json).
•	Webhook/event : optionnel, pour notifier un système externe lorsque l’audit est terminé.
•	Configuration : paramètres globaux (dossiers de travail, secrets, limites de ressources).
Comportement
1. Création d’un audit
1.	Recevoir une requête JSON‑RPC :
{
  "jsonrpc": "2.0",
  "method": "start_audit",
  "id": "client_request_id",
  "params": {
    "seed_url": "https://www.monsite.com",
    "max_urls": 300,
    "max_depth": 3,
    "modes": ["tech", "semantic"],
    "options": { "language": "fr" }
  }
}
1.	Valider les paramètres via audit_request.schema.json.
2.	Générer un audit_id unique (FSAUDIT-0002), créer le dossier /audits/{audit_id}/ et enregistrer un log initial (status = pending).
3.	Retourner un accusé de réception :
{
  "jsonrpc": "2.0",
  "id": "client_request_id",
  "result": {
    "audit_id": "FSAUDIT-0002",
    "status": "created"
  }
}
2. Orchestration des agents
L’orchestrateur s’appuie sur un gestionnaire de tâches (queue) ou un moteur de workflows pour exécuter les étapes :
1.	Crawl : appel à l’agent Crawler avec les paramètres de l’audit. L’orchestrateur écoute les événements (page_found, complete) et met à jour le statut (crawling, crawl_complete).
2.	Audit Technique & Analyse Sémantique : lancer ces deux agents dès que le crawl commence à produire des pages (streaming). Ils peuvent être exécutés en parallèle sur les mêmes entrées (streams indépendants). Le statut passe successivement à tech_running et semantic_running.
3.	Reporting : lorsque le tech_result.json et le semantic_result.json sont disponibles, déclencher la génération du rapport. Le statut passe à reporting.
4.	Finalisation : une fois le PDF et le HTML générés, le statut devient complete. L’orchestrateur enregistre l’horodatage et la durée totale, envoie un webhook si configuré et répercute l’état complet au client via JSON‑RPC.
3. Suivi et streaming
L’orchestrateur diffuse des messages de statut et d’avancement vers l’interface utilisateur via un canal SSE ou websockets :
{
  "jsonrpc": "2.0",
  "id": "FSAUDIT-0002",
  "result": {
    "status": "tech_running",
    "progress": 0.45,
    "message": "37 pages sur 80 analysées techniquement"
  }
}
Des événements peuvent également être enregistrés dans un log pour permettre un diagnostic a posteriori.
4. API publique
L’agent expose les endpoints suivants (via HTTP/REST ou JSON‑RPC) :
Méthode	Description	Entrées	Sorties
POST /audits	Démarrer un nouvel audit	audit_request	audit_id, statut initial
GET /audits/{audit_id}/status	Obtenir l’état courant de l’audit	audit_id	statut (phase, progression, erreurs)
GET /audits/{audit_id}/results	Récupérer les résultats JSON	audit_id	tech_result.json, semantic_result.json
GET /audits/{audit_id}/report	Télécharger le rapport PDF/HTML	audit_id	chemin vers les fichiers
POST /audits/{audit_id}/cancel	Annuler un audit en cours (si possible)	audit_id	statut final
Toutes les requêtes/ réponses peuvent être enveloppées dans un format JSON‑RPC pour l’intégration avec Claude Code.
5. Résilience et gestion des erreurs
•	Suivi de statut : l’audit peut être dans les états created, crawling, tech_running, semantic_running, reporting, complete, cancelled, failed.
•	Reprise : en cas de crash, l’orchestrateur récupère l’état et relance les agents restants (ex. relance de l’audit technique si l’audit sémantique est terminé).
•	Annulation : si l’utilisateur demande l’annulation, l’orchestrateur envoie un signal d’arrêt aux agents en cours et marque l’audit comme cancelled.
•	Timeout : chaque agent a un timeout configurable ; en cas de dépassement, l’orchestrateur marque l’étape comme failed et passe à l’étape suivante (ou arrête l’audit en fonction de la gravité).
•	Logs : tous les événements et erreurs sont enregistrés dans /audits/{audit_id}/audit.log pour diagnostic.
Tests et validation (TDD)
1.	Unitaires :
2.	Génération d’un audit_id unique.
3.	Transition d’états : vérifier que l’orchestrateur passe correctement d’un statut à l’autre et qu’il n’existe pas de transition invalide.
4.	Sérialisation/desérialisation des messages JSON‑RPC.
5.	Tests d’intégration :
6.	Exécuter un audit complet en mode simulé (mocks des agents) et vérifier que la séquence d’appels est correcte et que les fichiers tech_result.json et semantic_result.json sont bien transmis au Reporting.
7.	Tester l’annulation en plein milieu d’un crawl et s’assurer que le statut devient cancelled et que les ressources sont libérées.
8.	Simuler des erreurs (agent qui échoue, timeout) et vérifier que l’orchestrateur gère le cas (retry ou arrêt), met à jour le statut en failed et retourne l’erreur au client.
9.	Performance : vérifier que l’orchestrateur peut gérer plusieurs audits en parallèle (ex. 5 audits simultanés) sans surcharge excessive (moniteur de CPU et mémoire).
Contraintes et sécurité
•	Isolation des agents : chaque agent s’exécute dans un processus ou un container séparé pour éviter qu’une panne n’affecte l’orchestrateur.
•	Confidentialité : l’orchestrateur ne journalise pas de données sensibles (ex. URL exactes dans les logs), sauf autorisation explicite.
•	Rate limiting : limiter le nombre d’audits en parallèle par utilisateur pour éviter les abus.
•	Authentification : si l’API est exposée, utiliser des tokens ou OAuth pour sécuriser les endpoints.
En orchestrant l’ensemble des agents et en fournissant des points d’entrée clairs, l’agent Orchestrateur assure la cohérence et la robustesse du pipeline Fire Salamander.
 
[1] [2] N-gram : importance en SEO - NOIISE
https://www.noiise.com/ressources/seo/ngram/
