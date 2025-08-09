package seo

import (
	"context"
	"strings"
	"testing"
	"time"

	"firesalamander/internal/constants"
	"golang.org/x/net/html"
)

// Test de l'analyseur de balises SEO
func TestTagAnalyzerBasic(t *testing.T) {
	analyzer := NewTagAnalyzer()

	htmlContent := `
	<!DOCTYPE html>
	<html lang="fr">
	<head>
		<title>Test SEO Page - Guide Complet</title>
		<meta name="description" content="Découvrez notre guide complet pour améliorer votre référencement naturel avec des techniques SEO éprouvées et efficaces.">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="canonical" href="` + constants.TestCanonicalExample + `">
		<meta property="og:title" content="Guide SEO Complet">
		<meta property="og:description" content="Le meilleur guide SEO pour optimiser votre site web">
	</head>
	<body>
		<h1>Guide SEO Complet</h1>
		<h2>Introduction au SEO</h2>
		<p>Le référencement naturel est essentiel.</p>
		<img src="seo-guide.jpg" alt="Guide SEO illustré">
		<h2>Techniques avancées</h2>
		<p>Découvrez les meilleures pratiques.</p>
		<a href="/page-interne">Lien interne</a>
		<a href="` + constants.TestURLExternal + `">Lien externe</a>
	</body>
	</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Erreur parsing HTML: %v", err)
	}

	result, err := analyzer.Analyze(doc, htmlContent)
	if err != nil {
		t.Fatalf("Erreur analyse: %v", err)
	}

	// Tests du titre
	if !result.Title.Present {
		t.Error("Le titre devrait être présent")
	}

	if result.Title.Content != "Test SEO Page - Guide Complet" {
		t.Errorf("Titre incorrect: %s", result.Title.Content)
	}

	if !result.Title.OptimalLength {
		t.Logf("Longueur du titre: %d caractères (considérée non optimale)", result.Title.Length)
	}

	// Tests de la meta description
	if !result.MetaDescription.Present {
		t.Error("La meta description devrait être présente")
	}

	if !result.MetaDescription.OptimalLength {
		t.Error("La longueur de la meta description devrait être optimale")
	}

	if !result.MetaDescription.HasCallToAction {
		t.Logf("Meta description: '%s' - Appel à l'action non détecté", result.MetaDescription.Content)
	}

	// Tests des headings
	if result.Headings.H1Count != 1 {
		t.Errorf("Devrait avoir exactement 1 H1, trouvé: %d", result.Headings.H1Count)
	}

	if !result.Headings.HasHierarchy {
		t.Error("La hiérarchie des headings devrait être correcte")
	}

	// Tests des meta tags
	if !result.MetaTags.HasViewport {
		t.Error("Meta viewport devrait être présent")
	}

	if !result.MetaTags.HasCanonical {
		t.Error("URL canonique devrait être présente")
	}

	if !result.MetaTags.HasOGTags {
		t.Error("Balises Open Graph devraient être présentes")
	}

	// Tests des images
	if result.Images.TotalImages != 1 {
		t.Errorf("Devrait avoir 1 image, trouvé: %d", result.Images.TotalImages)
	}

	if result.Images.AltTextCoverage != 1.0 {
		t.Errorf("Couverture alt text devrait être 100%%, trouvé: %.1f%%", result.Images.AltTextCoverage*100)
	}

	// Tests des liens
	if result.Links.InternalLinks < 1 {
		t.Error("Devrait avoir au moins 1 lien interne")
	}

	if result.Links.ExternalLinks < 1 {
		t.Error("Devrait avoir au moins 1 lien externe")
	}

	t.Logf("Analyse réussie - Titre: %s, Meta: présente, H1: %d, Images: %d",
		result.Title.Content, result.Headings.H1Count, result.Images.TotalImages)
}

// Test de l'analyseur de performance
func TestPerformanceAnalyzerBasic(t *testing.T) {
	analyzer := NewPerformanceAnalyzer()

	// Test avec du contenu HTML simulé
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Test Performance</title>
		<link rel="stylesheet" href="style.css">
		<script src="script.js"></script>
	</head>
	<body>
		<img src="image1.jpg" alt="Image 1">
		<img src="image2.webp" alt="Image 2">
		<img src="image3.png" alt="Image 3">
		<script src="analytics.js"></script>
	</body>
	</html>
	`

	testResult := &PerformanceMetricsResult{
		Issues:          []string{},
		Recommendations: []string{},
	}

	err := analyzer.analyzeHTMLResources(htmlContent, testResult)
	if err != nil {
		t.Fatalf("Erreur analyse HTML: %v", err)
	}

	// Vérifier le comptage des ressources via les métriques
	// (Note: analyzeHTMLResources modifie le résultat passé en paramètre)
	
	t.Logf("Analyse des ressources terminée")
}

// Test de l'auditeur technique
func TestTechnicalAuditorBasic(t *testing.T) {
	auditor := NewTechnicalAuditor()

	htmlContent := `
	<!DOCTYPE html>
	<html lang="fr">
	<head>
		<title>Test Audit Technique</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="robots" content="index, follow">
		<link rel="canonical" href="` + constants.TestCanonicalTest + `">
	</head>
	<body>
		<h1>Test Audit</h1>
		<p>Contenu de test pour l'audit technique.</p>
		<img src="test.jpg" alt="Image de test">
	</body>
	</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Erreur parsing HTML: %v", err)
	}

	ctx := context.Background()
	result, err := auditor.Audit(ctx, constants.TestURLExampleTest, doc, htmlContent)
	if err != nil {
		t.Fatalf("Erreur audit: %v", err)
	}

	// Tests d'audit mobile
	if !result.Mobile.HasViewport {
		t.Error("Meta viewport devrait être détecté")
	}

	if !result.Mobile.IsResponsive {
		t.Error("Le site devrait être détecté comme responsive")
	}

	// Tests de structure
	if !result.Structure.ValidHTML {
		t.Error("HTML devrait être valide (DOCTYPE présent)")
	}

	// Tests d'accessibilité
	if result.Accessibility.AltTagCoverage != 1.0 {
		t.Errorf("Couverture alt text devrait être 100%%, trouvé: %.1f%%", result.Accessibility.AltTagCoverage*100)
	}

	// Tests d'indexabilité
	if result.Indexability.BlockedByRobots {
		t.Error("La page ne devrait pas être bloquée par robots")
	}

	if !result.Indexability.HasCanonical {
		t.Error("URL canonique devrait être détectée")
	}

	t.Logf("Audit technique réussi - Mobile: %.1f%%, Structure: %.1f%%, Accessibilité: %.1f%%",
		result.Mobile.MobileScore*100, result.Structure.StructureScore*100, result.Accessibility.Score*100)
}

// Test du moteur de recommandations
func TestRecommendationEngineBasic(t *testing.T) {
	engine := NewRecommendationEngine()

	// Créer une analyse simulée avec des problèmes
	analysis := &SEOAnalysisResult{
		URL:          constants.TestURLExampleTest,
		Domain:       "example.com",
		Protocol:     "https",
		StatusCode:   constants.HTTPStatusOK,
		OverallScore: 45.5, // Score faible
		CategoryScores: map[string]float64{
			"tags":        0.3, // 30% - faible
			"performance": 0.6, // 60% - correct
			"technical":   0.4, // 40% - faible
			"basics":      0.8, // 80% - bon
		},
		TagAnalysis: TagAnalysisResult{
			Title: TitleAnalysis{
				Present:       false, // Problème critique
				OptimalLength: false,
			},
			MetaDescription: MetaDescAnalysis{
				Present:          true,
				OptimalLength:    false, // Problème
				HasCallToAction:  false,
			},
			Headings: HeadingAnalysis{
				H1Count:      0, // Problème
				HasHierarchy: false,
			},
			MetaTags: MetaTagsAnalysis{
				HasViewport:  false, // Problème critique
				HasCanonical: false,
				HasOGTags:    false,
			},
			Images: ImageAnalysis{
				TotalImages:     3,
				ImagesWithAlt:   1,
				AltTextCoverage: 0.33, // Problème
			},
		},
		PerformanceMetrics: PerformanceMetricsResult{
			LoadTime:        5 * time.Second, // Problème
			HasCompression:  false,           // Problème
			HasCaching:      false,           // Problème
			OptimizedImages: false,           // Problème
			CoreWebVitals: CoreWebVitals{
				LCP: EstimatedMetric{
					Value: 4500, // Problème
					Score: "poor",
				},
				FID: EstimatedMetric{
					Value: constants.HTTPStatusOK, // OK
					Score: "needs-improvement",
				},
				CLS: EstimatedMetric{
					Value: 0.3, // Problème
					Score: "poor",
				},
			},
		},
		TechnicalAudit: TechnicalAuditResult{
			Security: SecurityAudit{
				HasHTTPS:     true, // OK
				MixedContent: false,
			},
			Mobile: MobileAudit{
				HasViewport:   false, // Problème
				IsResponsive:  false, // Problème
				MobileScore:   0.2,
			},
			Structure: StructureAudit{
				HasSitemap:   false, // Problème
				HasRobotsTxt: false, // Problème
				ValidHTML:    true,
			},
			Accessibility: AccessibilityAudit{
				Score:          0.4, // Problème
				AltTagCoverage: 0.33,
			},
			Indexability: IndexabilityAudit{
				HasNoIndex:     false,
				HasCanonical:   false, // Problème
				BlockedByRobots: false,
			},
			Crawlability: CrawlabilityAudit{
				InternalLinks: 1, // Problème
				BrokenLinks:   []string{constants.TestURLBroken},
			},
		},
	}

	recommendations := engine.GenerateRecommendations(analysis)

	if len(recommendations) == 0 {
		t.Fatal("Aucune recommandation générée")
	}

	// Vérifier qu'on a des recommandations critiques
	hasCritical := false
	hasTitle := false
	hasViewport := false

	for _, rec := range recommendations {
		if rec.Priority == PriorityCritical {
			hasCritical = true
		}
		if rec.ID == "missing-title" {
			hasTitle = true
		}
		if rec.ID == "missing-viewport" {
			hasViewport = true
		}
	}

	if !hasCritical {
		t.Error("Devrait avoir au moins une recommandation critique")
	}

	if !hasTitle {
		t.Error("Devrait avoir une recommandation pour le titre manquant")
	}

	if !hasViewport {
		t.Error("Devrait avoir une recommandation pour le viewport manquant")
	}

	// Vérifier que les recommandations sont triées par priorité
	if len(recommendations) > 1 {
		firstPrio := engine.getPriorityWeight(recommendations[0].Priority)
		secondPrio := engine.getPriorityWeight(recommendations[1].Priority)
		
		if firstPrio < secondPrio {
			t.Error("Les recommandations devraient être triées par priorité décroissante")
		}
	}

	t.Logf("Recommandations générées: %d", len(recommendations))
	
	// Afficher les 3 premières recommandations pour debug
	for i, rec := range recommendations {
		if i >= 3 {
			break
		}
		t.Logf("Rec %d: %s (Priorité: %s, Impact: %s)", i+1, rec.Title, rec.Priority, rec.Impact)
	}
}

// Test d'intégration complète du module SEO
func TestSEOAnalyzerIntegration(t *testing.T) {
	// Note: Ce test nécessiterait un serveur HTTP mock pour être complet
	// Pour l'instant, on teste seulement la création de l'analyseur
	
	analyzer := NewSEOAnalyzer()
	
	if analyzer == nil {
		t.Fatal("L'analyseur SEO n'a pas pu être créé")
	}

	if analyzer.tagAnalyzer == nil {
		t.Error("L'analyseur de balises n'est pas initialisé")
	}

	if analyzer.performanceAnalyzer == nil {
		t.Error("L'analyseur de performance n'est pas initialisé")
	}

	if analyzer.technicalAuditor == nil {
		t.Error("L'auditeur technique n'est pas initialisé")
	}

	if analyzer.recommendationEngine == nil {
		t.Error("Le moteur de recommandations n'est pas initialisé")
	}

	t.Log("Intégration SEO - Tous les composants sont initialisés correctement")
}