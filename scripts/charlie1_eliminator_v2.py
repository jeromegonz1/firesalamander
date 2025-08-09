#!/usr/bin/env python3
"""
🤖 CHARLIE-1 SMART ELIMINATOR V2
Elimination industrielle ciblée avec contrôle contextuel avancé
"""

import re
import os
import json
from pathlib import Path

def create_selective_constants_mapping():
    """Crée le mapping sélectif pour éviter les erreurs de syntaxe"""
    return {
        # Template IDs (uniquement dans les appels de fonction)
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
        
        # Score Comparisons (string literals)
        '"poor"': 'constants.RecommendationScorePoor',
        '"good"': 'constants.RecommendationScoreGood', 
        '"needs_improvement"': 'constants.RecommendationScoreNeedsImprovement',
        
        # Target Values (string literals)
        '"≤ 2.5s"': 'constants.RecommendationTargetLCP',
        '"≤ 100ms"': 'constants.RecommendationTargetFID',
        '"≤ 0.1"': 'constants.RecommendationTargetCLS',
        '"< 2 MB"': 'constants.RecommendationTargetPageSize',
        '"≥ 70%"': 'constants.RecommendationTargetScore',
        '"≥ 3"': 'constants.RecommendationTargetLinks',
        '"≥ 70"': 'constants.RecommendationTargetScore',
        
        # Time Estimates (string literals)
        '"1-2 heures"': 'constants.RecommendationTimeLow',
        '"4-8 heures"': 'constants.RecommendationTimeMedium',
        '"1-2 jours"': 'constants.RecommendationTimeHigh',
        '"Variable"': 'constants.RecommendationTimeVariable',
        
        # Content Ranges (string literals)
        '"30-60 caractères"': 'constants.RecommendationTitleRange',
        '"120-160 caractères"': 'constants.RecommendationMetaDescRange',
        
        # Resource Types (string literals)
        '"documentation"': 'constants.RecommendationResourceTypeDoc',
        
        # Issue Messages (string literals)
        '"Titre manquant"': 'constants.RecommendationIssueTitleMissing',
        
        # Default Messages (string literals - only in default template)
        '"Recommandation SEO"': 'constants.RecommendationDefaultTitle',
        '"Amélioration SEO recommandée"': 'constants.RecommendationDefaultDescription',
        '"general"': 'constants.RecommendationCategoryGeneral',
        
        # Template Titles (struct literals)
        '"Ajouter un titre de page"': 'constants.RecommendationTitleAddPageTitle',
        '"Ajouter une meta description"': 'constants.RecommendationTitleAddMetaDesc',
        '"Migrer vers HTTPS"': 'constants.RecommendationTitleMigrateHTTPS',
        '"Améliorer le Largest Contentful Paint (LCP)"': 'constants.RecommendationTitleImproveLCP',
        
        # Template Descriptions (struct literals)
        '"La page n\'a pas de balise <title>. C\'est un élément fondamental pour le SEO."': 'constants.RecommendationDescMissingTitle',
        '"La page n\'a pas de meta description. Elle influence le taux de clic dans les résultats de recherche."': 'constants.RecommendationDescMissingMetaDesc',
        '"Le site utilise HTTP au lieu de HTTPS. Google favorise les sites sécurisés."': 'constants.RecommendationDescMigrateHTTPS',
        '"Le LCP actuel est de {current_value}, l\'objectif est {target}. Optimisez le chargement du contenu principal."': 'constants.RecommendationDescImproveLCP',
        
        # Category Names (template structs only, not enum values)
        '"tags"': 'constants.RecommendationCategoryTags',
        '"performance"': 'constants.RecommendationCategoryPerformance', 
        '"security"': 'constants.RecommendationCategorySecurity',
        
        # Action Items (array elements)
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
        
        # Metrics (array elements)
        '"Taux de clic"': 'constants.RecommendationMetricCTR',
        '"Position dans les SERP"': 'constants.RecommendationMetricSERPPosition',
        '"Impressions"': 'constants.RecommendationMetricImpressions',
        '"Trust signals"': 'constants.RecommendationMetricTrustSignals',
        '"Ranking boost"': 'constants.RecommendationMetricRankingBoost',
        '"Page Experience"': 'constants.RecommendationMetricPageExperience',
        '"Core Web Vitals"': 'constants.RecommendationMetricCoreWebVitals',
        '"LCP"': 'constants.RecommendationMetricLCP',
        
        # Tags (array elements)
        '"critique"': 'constants.RecommendationTagCritical',
        '"balises"': 'constants.RecommendationTagTags',
        '"onpage"': 'constants.RecommendationTagOnPage',
        '"meta"': 'constants.RecommendationTagMeta',
        '"ctr"': 'constants.RecommendationTagCTR',
        '"sécurité"': 'constants.RecommendationTagSecurity',
        '"technique"': 'constants.RecommendationTagTechnical',
        '"core-web-vitals"': 'constants.RecommendationTagCoreWebVitals',
        '"lcp"': 'constants.RecommendationTagLCP',
        
        # Numeric Constants (only specific cases)
        '20': 'constants.RecommendationMaxRecommendations',  # Only in > 20 comparison
        '90': 'constants.RecommendationRuleMissingTitle',    # Only in priority rules
        '95': 'constants.RecommendationRuleMigrateHTTPS',    # Only in priority rules  
        '85': 'constants.RecommendationRuleMissingMetaDesc', # Only in priority rules
        '75': 'constants.RecommendationRuleImproveLCP',      # Only in priority rules
        '70': 'constants.RecommendationRuleMissingH1',       # Only in priority rules
        '50': 'constants.RecommendationScoreThresholdLow',   # Only in < 50 comparison
        '0.5': 'constants.RecommendationCategoryThresholdLow', # Only in < 0.5 comparison
        '0.7': 'constants.RecommendationAccessibilityThreshold', # Only in < 0.7 comparison
        '1.0': 'constants.RecommendationAltTextThreshold',   # Only in < 1.0 comparison
        '3': 'constants.RecommendationMinInternalLinks',     # Only in < 3 comparison
        '2*1024*1024': 'constants.RecommendationMaxPageSizeBytes', # Only in size comparison
        
        # Priority Weights (only in switch return statements)
        '4': 'constants.RecommendationPriorityWeightCritical',
        '2': 'constants.RecommendationPriorityWeightMedium', 
        '1': 'constants.RecommendationPriorityWeightLow',
        '0': 'constants.RecommendationPriorityWeightDefault',
        
        # Special replacements
        'constants.HighQualityScore': 'constants.RecommendationRuleMissingMetaDesc',
        '{%s}': 'constants.RecommendationPlaceholderPattern',
        '"technique"': 'constants.RecommendationTagTechnical',
    }

def should_replace_in_context(line, hardcoded_value, constant):
    """Détermine si un remplacement doit être effectué basé sur le contexte"""
    
    # Skip struct tags completely
    if '`json:' in line:
        return False
        
    # Skip comments
    if line.strip().startswith('//'):
        return False
        
    # Skip const/type/var declarations
    if re.match(r'^\s*(const|type|var)\s*\(?\s*$', line.strip()):
        return False
        
    # Skip const declarations with assignments (like Priority = "critical")
    if re.match(r'^\s*\w+\s+(Priority|Impact|Effort)\s*=', line.strip()):
        return False
        
    # Special handling for numeric values
    if hardcoded_value in ['20', '90', '95', '85', '75', '70', '50']:
        # Only replace in specific contexts
        if hardcoded_value == '20' and '> 20' in line and 'len(recommendations)' in line:
            return True
        elif hardcoded_value in ['90', '95', '85', '75', '70'] and 'priorityRules[' in line:
            return True
        elif hardcoded_value == '50' and 'OverallScore <' in line:
            return True
        else:
            return False
    
    if hardcoded_value in ['0.5', '0.7', '1.0', '3']:
        # Only replace in specific comparison contexts
        if hardcoded_value == '0.5' and 'score <' in line:
            return True
        elif hardcoded_value == '0.7' and 'Score <' in line:
            return True
        elif hardcoded_value == '1.0' and 'AltTextCoverage <' in line:
            return True
        elif hardcoded_value == '3' and 'InternalLinks <' in line:
            return True
        else:
            return False
    
    if hardcoded_value in ['4', '2', '1', '0']:
        # Only replace in return statements of weight functions
        if 'return ' in line and ('getPriorityWeight' in line or 'getImpactWeight' in line):
            return True
        else:
            return False
            
    if hardcoded_value == '2*1024*1024':
        # Only in PageSize comparison
        if 'PageSize >' in line:
            return True
        else:
            return False
    
    # Default: allow replacement for string literals
    return True

def eliminate_hardcoding_violations_selective(filepath, constants_mapping):
    """Élimine sélectivement les violations avec contrôle contextuel"""
    
    print(f"🤖 ÉLIMINATION SÉLECTIVE: {filepath}")
    
    # Lire le fichier
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements = 0
    
    # Split into lines for context-aware processing
    lines = content.split('\n')
    new_lines = []
    
    for line_num, line in enumerate(lines):
        new_line = line
        
        for hardcoded_value, constant in constants_mapping.items():
            if hardcoded_value in line:
                if should_replace_in_context(line, hardcoded_value, constant):
                    old_line = new_line
                    new_line = new_line.replace(hardcoded_value, constant)
                    if new_line != old_line:
                        replacements += 1
                        print(f"  ✅ Line {line_num + 1}: {hardcoded_value} → {constant}")
        
        new_lines.append(new_line)
    
    content = '\n'.join(new_lines)
    
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
    print("🤖 CHARLIE-1 SMART ELIMINATOR V2 - Élimination sélective intelligente...")
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/recommendation_engine.go'
    
    # Créer une sauvegarde V2
    backup_path = f"{filepath}.charlie1_v2_backup"
    if not os.path.exists(backup_path):
        with open(filepath, 'r') as original:
            with open(backup_path, 'w') as backup:
                backup.write(original.read())
        print(f"💾 Sauvegarde V2 créée: {backup_path}")
    
    # Obtenir le mapping sélectif des constantes
    constants_mapping = create_selective_constants_mapping()
    print(f"📋 {len(constants_mapping)} mappings sélectifs chargés")
    
    # Éliminer les violations sélectivement
    total_eliminated = eliminate_hardcoding_violations_selective(filepath, constants_mapping)
    
    print(f"\n🎯 CHARLIE-1 V2 TERMINÉ:")
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