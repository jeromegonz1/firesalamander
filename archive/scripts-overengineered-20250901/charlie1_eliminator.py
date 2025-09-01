#!/usr/bin/env python3
"""
🤖 CHARLIE-1 SMART ELIMINATOR
Elimination industrielle automatisée des violations hardcoding dans recommendation_engine.go
"""

import re
import os
import json
from pathlib import Path

def load_analysis_results():
    """Charge les résultats de l'analyse CHARLIE-1"""
    try:
        with open('charlie1_analysis.json', 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print("⚠️ charlie1_analysis.json non trouvé. Exécutez d'abord charlie1_detector.py")
        return None

def create_constants_mapping():
    """Crée le mapping complet des constantes pour CHARLIE-1"""
    return {
        # JSON Field Names
        '"id"': 'constants.RecommendationJSONFieldID',
        '"title"': 'constants.RecommendationJSONFieldTitle', 
        '"description"': 'constants.RecommendationJSONFieldDescription',
        '"category"': 'constants.RecommendationJSONFieldCategory',
        '"priority"': 'constants.RecommendationJSONFieldPriority',
        '"impact"': 'constants.RecommendationJSONFieldImpact',
        '"effort"': 'constants.RecommendationJSONFieldEffort',
        '"actions"': 'constants.RecommendationJSONFieldActions',
        '"resources"': 'constants.RecommendationJSONFieldResources',
        '"metrics"': 'constants.RecommendationJSONFieldMetrics',
        '"tags"': 'constants.RecommendationJSONFieldTags',
        '"task"': 'constants.RecommendationJSONFieldTask',
        '"technical"': 'constants.RecommendationJSONFieldTechnical',
        '"estimated_time"': 'constants.RecommendationJSONFieldEstimatedTime',
        
        # Priority Levels
        '"critical"': 'constants.RecommendationPriorityCritical',
        '"high"': 'constants.RecommendationPriorityHigh',
        '"medium"': 'constants.RecommendationPriorityMedium',
        '"low"': 'constants.RecommendationPriorityLow',
        
        # Impact Levels (specific to types)
        'ImpactHigh': 'constants.RecommendationImpactHigh',
        'ImpactMedium': 'constants.RecommendationImpactMedium', 
        'ImpactLow': 'constants.RecommendationImpactLow',
        
        # Effort Levels (specific to types)
        'EffortLow': 'constants.RecommendationEffortLow',
        'EffortMedium': 'constants.RecommendationEffortMedium',
        'EffortHigh': 'constants.RecommendationEffortHigh',
        
        # Priority Types (specific to types)
        'PriorityCritical': 'constants.RecommendationPriorityCritical',
        'PriorityHigh': 'constants.RecommendationPriorityHigh',
        'PriorityMedium': 'constants.RecommendationPriorityMedium',
        'PriorityLow': 'constants.RecommendationPriorityLow',
        
        # Category Names
        '"tags"': 'constants.RecommendationCategoryTags',
        '"performance"': 'constants.RecommendationCategoryPerformance',
        '"security"': 'constants.RecommendationCategorySecurity',
        '"general"': 'constants.RecommendationCategoryGeneral',
        '"technical"': 'constants.RecommendationCategoryTechnical',
        '"content"': 'constants.RecommendationCategoryContent',
        '"mobile"': 'constants.RecommendationCategoryMobile',
        '"structure"': 'constants.RecommendationCategoryStructure',
        '"accessibility"': 'constants.RecommendationCategoryAccessibility',
        
        # Template IDs
        '"missing-title"': 'constants.RecommendationTemplateIDMissingTitle',
        '"missing-meta-desc"': 'constants.RecommendationTemplateIDMissingMetaDesc',
        '"title-length"': 'constants.RecommendationTemplateIDTitleLength',
        '"meta-desc-length"': 'constants.RecommendationTemplateIDMetaDescLength',
        '"missing-h1"': 'constants.RecommendationTemplateIDMissingH1',
        '"multiple-h1"': 'constants.RecommendationTemplateIDMultipleH1',
        '"heading-hierarchy"': 'constants.RecommendationTemplateIDHeadingHierarchy',
        '"missing-alt-tags"': 'constants.RecommendationTemplateIDMissingAltTags',
        '"missing-viewport"': 'constants.RecommendationTemplateIDMissingViewport',
        '"missing-canonical"': 'constants.RecommendationTemplateIDMissingCanonical',
        '"missing-og-tags"': 'constants.RecommendationTemplateIDMissingOGTags',
        '"improve-lcp"': 'constants.RecommendationTemplateIDImproveLCP',
        '"improve-fid"': 'constants.RecommendationTemplateIDImproveFID',
        '"improve-cls"': 'constants.RecommendationTemplateIDImproveCLS',
        '"enable-compression"': 'constants.RecommendationTemplateIDEnableCompression',
        '"configure-caching"': 'constants.RecommendationTemplateIDConfigureCaching',
        '"optimize-images"': 'constants.RecommendationTemplateIDOptimizeImages',
        '"minify-resources"': 'constants.RecommendationTemplateIDMinifyResources',
        '"reduce-page-size"': 'constants.RecommendationTemplateIDReducePageSize',
        '"migrate-https"': 'constants.RecommendationTemplateIDMigrateHTTPS',
        '"fix-mixed-content"': 'constants.RecommendationTemplateIDFixMixedContent',
        '"add-hsts"': 'constants.RecommendationTemplateIDAddHSTS',
        '"make-responsive"': 'constants.RecommendationTemplateIDMakeResponsive',
        '"add-sitemap"': 'constants.RecommendationTemplateIDAddSitemap',
        '"add-robots-txt"': 'constants.RecommendationTemplateIDAddRobotsTxt',
        '"remove-noindex"': 'constants.RecommendationTemplateIDRemoveNoIndex',
        '"fix-duplicate-content"': 'constants.RecommendationTemplateIDFixDuplicateContent',
        '"improve-accessibility"': 'constants.RecommendationTemplateIDImproveAccessibility',
        '"fix-broken-links"': 'constants.RecommendationTemplateIDFixBrokenLinks',
        '"improve-internal-linking"': 'constants.RecommendationTemplateIDImproveInternalLinking',
        '"overall-seo-audit"': 'constants.RecommendationTemplateIDOverallSEOAudit',
        
        # Score Comparisons
        '"poor"': 'constants.RecommendationScorePoor',
        '"good"': 'constants.RecommendationScoreGood',
        '"needs_improvement"': 'constants.RecommendationScoreNeedsImprovement',
        
        # Target Values
        '"≤ 2.5s"': 'constants.RecommendationTargetLCP',
        '"≤ 100ms"': 'constants.RecommendationTargetFID', 
        '"≤ 0.1"': 'constants.RecommendationTargetCLS',
        '"< 2 MB"': 'constants.RecommendationTargetPageSize',
        '"≥ 70%"': 'constants.RecommendationTargetScore',
        '"≥ 3"': 'constants.RecommendationTargetLinks',
        '"≥ 70"': 'constants.RecommendationTargetScore',
        
        # Time Estimates
        '"1-2 heures"': 'constants.RecommendationTimeLow',
        '"4-8 heures"': 'constants.RecommendationTimeMedium',
        '"1-2 jours"': 'constants.RecommendationTimeHigh',
        '"Variable"': 'constants.RecommendationTimeVariable',
        
        # Content Ranges
        '"30-60 caractères"': 'constants.RecommendationTitleRange',
        '"120-160 caractères"': 'constants.RecommendationMetaDescRange',
        
        # Resource Types
        '"documentation"': 'constants.RecommendationResourceTypeDoc',
        '"guide"': 'constants.RecommendationResourceTypeGuide',
        '"tool"': 'constants.RecommendationResourceTypeTool',
        '"example"': 'constants.RecommendationResourceTypeExample',
        
        # Default Messages
        '"Recommandation SEO"': 'constants.RecommendationDefaultTitle',
        '"Amélioration SEO recommandée"': 'constants.RecommendationDefaultDescription',
        
        # Issue Messages
        '"Titre manquant"': 'constants.RecommendationIssueTitleMissing',
        '"Meta description manquante"': 'constants.RecommendationIssueMetaDescMissing',
        '"H1 manquant"': 'constants.RecommendationIssueH1Missing',
        '"Meta viewport manquante"': 'constants.RecommendationIssueViewportMissing',
        
        # Template Titles
        '"Ajouter un titre de page"': 'constants.RecommendationTitleAddPageTitle',
        '"Ajouter une meta description"': 'constants.RecommendationTitleAddMetaDesc',
        '"Migrer vers HTTPS"': 'constants.RecommendationTitleMigrateHTTPS',
        '"Améliorer le Largest Contentful Paint (LCP)"': 'constants.RecommendationTitleImproveLCP',
        
        # Template Descriptions
        '"La page n\'a pas de balise <title>. C\'est un élément fondamental pour le SEO."': 'constants.RecommendationDescMissingTitle',
        '"La page n\'a pas de meta description. Elle influence le taux de clic dans les résultats de recherche."': 'constants.RecommendationDescMissingMetaDesc',
        '"Le site utilise HTTP au lieu de HTTPS. Google favorise les sites sécurisés."': 'constants.RecommendationDescMigrateHTTPS',
        '"Le LCP actuel est de {current_value}, l\'objectif est {target}. Optimisez le chargement du contenu principal."': 'constants.RecommendationDescImproveLCP',
        
        # Action Items
        '"Ajouter une balise <title> descriptive et unique"': 'constants.RecommendationActionAddTitleTag',
        '"Inclure les mots-clés principaux"': 'constants.RecommendationActionIncludeKeywords',
        '"Respecter la longueur optimale (30-60 caractères)"': 'constants.RecommendationActionRespectLength',
        '"Respecter la longueur optimale (120-160 caractères)"': 'constants.RecommendationActionRespectLength',
        '"Ajouter une meta description attrayante"': 'constants.RecommendationActionAddMetaDesc',
        '"Inclure un appel à l\'action"': 'constants.RecommendationActionIncludeCTA',
        '"Obtenir un certificat SSL/TLS"': 'constants.RecommendationActionGetSSLCert',
        '"Configurer le serveur pour HTTPS"': 'constants.RecommendationActionConfigureHTTPS',
        '"Rediriger tout le trafic HTTP vers HTTPS"': 'constants.RecommendationActionRedirectHTTPS',
        '"Mettre à jour les liens internes"': 'constants.RecommendationActionUpdateLinks',
        '"Optimiser les images de l\'above-the-fold"': 'constants.RecommendationActionOptimizeImages',
        '"Améliorer le temps de réponse du serveur"': 'constants.RecommendationActionImproveServer',
        '"Précharger les ressources critiques"': 'constants.RecommendationActionPreloadResources',
        '"Utiliser un CDN"': 'constants.RecommendationActionUseCDN',
        
        # Metrics
        '"Taux de clic"': 'constants.RecommendationMetricCTR',
        '"Position dans les SERP"': 'constants.RecommendationMetricSERPPosition',
        '"Impressions"': 'constants.RecommendationMetricImpressions',
        '"Trust signals"': 'constants.RecommendationMetricTrustSignals',
        '"Ranking boost"': 'constants.RecommendationMetricRankingBoost',
        '"Page Experience"': 'constants.RecommendationMetricPageExperience',
        '"Core Web Vitals"': 'constants.RecommendationMetricCoreWebVitals',
        '"LCP"': 'constants.RecommendationMetricLCP',
        
        # Tags
        '"critique"': 'constants.RecommendationTagCritical',
        '"balises"': 'constants.RecommendationTagTags',
        '"onpage"': 'constants.RecommendationTagOnPage',
        '"meta"': 'constants.RecommendationTagMeta',
        '"ctr"': 'constants.RecommendationTagCTR',
        '"sécurité"': 'constants.RecommendationTagSecurity',
        '"technique"': 'constants.RecommendationTagTechnical',
        '"performance"': 'constants.RecommendationTagPerformance',
        '"core-web-vitals"': 'constants.RecommendationTagCoreWebVitals',
        '"lcp"': 'constants.RecommendationTagLCP',
        
        # Operators
        '"≤"': 'constants.RecommendationOperatorLessEqual',
        '"≥"': 'constants.RecommendationOperatorGreaterEqual', 
        '"<"': 'constants.RecommendationOperatorLess',
        '">"': 'constants.RecommendationOperatorGreater',
        
        # Numeric Constants (only for critical thresholds)
        '20': 'constants.RecommendationMaxRecommendations',
        '90': 'constants.RecommendationRuleMissingTitle',
        '95': 'constants.RecommendationRuleMigrateHTTPS',
        '85': 'constants.RecommendationRuleMissingMetaDesc',
        '75': 'constants.RecommendationRuleImproveLCP',
        '70': 'constants.RecommendationRuleMissingH1',
        '50': 'constants.RecommendationScoreThresholdLow',
        '0.5': 'constants.RecommendationCategoryThresholdLow',
        '0.7': 'constants.RecommendationAccessibilityThreshold',
        '1.0': 'constants.RecommendationAltTextThreshold',
        '3': 'constants.RecommendationMinInternalLinks',
        '2*1024*1024': 'constants.RecommendationMaxPageSizeBytes',
        
        # Priority Weights
        '4': 'constants.RecommendationPriorityWeightCritical', 
        '3': 'constants.RecommendationPriorityWeightHigh',
        '2': 'constants.RecommendationPriorityWeightMedium',
        '1': 'constants.RecommendationPriorityWeightLow',
        '0': 'constants.RecommendationPriorityWeightDefault',
        
        # Placeholders (special handling)
        '"{%s}"': 'constants.RecommendationPlaceholderPattern',
        '"{current_value}"': 'constants.RecommendationPlaceholderCurrentValue',
        '"{target}"': 'constants.RecommendationPlaceholderTarget',
        '"{range}"': 'constants.RecommendationPlaceholderRange',
        '"{category}"': 'constants.RecommendationPlaceholderCategory',
        '"{count}"': 'constants.RecommendationPlaceholderCount',
    }

def eliminate_hardcoding_violations(filepath, constants_mapping):
    """Élimine les violations de hardcoding en utilisant le mapping des constantes"""
    
    print(f"🤖 ÉLIMINATION: {filepath}")
    
    # Lire le fichier
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements = 0
    
    # Effectuer les remplacements
    for hardcoded_value, constant in constants_mapping.items():
        if hardcoded_value in content:
            # Éviter les remplacements dans les commentaires et struct tags
            lines = content.split('\n')
            new_lines = []
            
            for line in lines:
                # Ignorer les commentaires
                if line.strip().startswith('//'):
                    new_lines.append(line)
                    continue
                    
                # Gestion spéciale des struct tags - ne pas remplacer dans `json:"..."`
                if '`json:' in line and hardcoded_value in line:
                    # Ne pas remplacer dans les struct tags JSON
                    new_lines.append(line)
                    continue
                
                # Remplacements normaux
                if hardcoded_value in line:
                    # Vérifier que ce n'est pas dans une struct tag
                    if '`json:' not in line or hardcoded_value not in line[line.find('`json:'):line.find('`', line.find('`json:') + 1) + 1]:
                        new_line = line.replace(hardcoded_value, constant)
                        if new_line != line:
                            replacements += 1
                            print(f"  ✅ Line: {hardcoded_value} → {constant}")
                        new_lines.append(new_line)
                    else:
                        new_lines.append(line)
                else:
                    new_lines.append(line)
            
            content = '\n'.join(new_lines)
    
    # Remplacements spéciaux pour les patterns contextuels
    special_replacements = {
        # Fmt.Sprintf avec placeholders
        'fmt.Sprintf(\"{%s}\", key)': 'fmt.Sprintf(constants.RecommendationPlaceholderPattern, key)',
        
        # Constantes HighQualityScore replacement
        'constants.HighQualityScore': 'constants.RecommendationRuleMissingMetaDesc',
        
        # URL constants docs
        'constants.GoogleTitleLinkDocs': 'constants.GoogleTitleLinkDocs',
        'constants.GoogleSnippetDocs': 'constants.GoogleSnippetDocs',
        'constants.GoogleHTTPSDocs': 'constants.GoogleHTTPSDocs',
        'constants.WebDevLCPDocs': 'constants.WebDevLCPDocs',
    }
    
    for pattern, replacement in special_replacements.items():
        if pattern in content:
            content = content.replace(pattern, replacement)
            replacements += 1
            print(f"  ✅ Special: {pattern} → {replacement}")
    
    # Sauvegarder si des modifications ont été effectuées
    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"✅ ÉLIMINÉ: {replacements} violations dans {filepath}")
        return replacements
    else:
        print(f"ℹ️ Aucune modification nécessaire dans {filepath}")
        return 0

def main():
    print("🤖 CHARLIE-1 SMART ELIMINATOR - Démarrage de l'élimination industrielle...")
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/recommendation_engine.go'
    
    # Créer une sauvegarde
    backup_path = f"{filepath}.charlie1_backup"
    if not os.path.exists(backup_path):
        with open(filepath, 'r') as original:
            with open(backup_path, 'w') as backup:
                backup.write(original.read())
        print(f"💾 Sauvegarde créée: {backup_path}")
    
    # Obtenir le mapping des constantes
    constants_mapping = create_constants_mapping()
    print(f"📋 {len(constants_mapping)} mappings de constantes chargés")
    
    # Éliminer les violations
    total_eliminated = eliminate_hardcoding_violations(filepath, constants_mapping)
    
    print(f"\n🎯 CHARLIE-1 TERMINÉ:")
    print(f"✅ Total violations éliminées: {total_eliminated}")
    print(f"📁 Fichier traité: {filepath}")
    print(f"💾 Sauvegarde: {backup_path}")
    
    # Tester la compilation
    print("\n🔨 Test de compilation...")
    import subprocess
    try:
        result = subprocess.run(['go', 'build', './internal/seo/...'], 
                              capture_output=True, text=True, cwd='/Users/jeromegonzalez/claude-code/fire-salamander')
        if result.returncode == 0:
            print("✅ Compilation réussie!")
        else:
            print("⚠️ Erreurs de compilation:")
            print(result.stderr)
    except Exception as e:
        print(f"⚠️ Impossible de tester la compilation: {e}")
    
    return total_eliminated

if __name__ == "__main__":
    main()