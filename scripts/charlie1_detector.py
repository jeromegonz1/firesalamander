#!/usr/bin/env python3
"""
🔍 CHARLIE-1 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans recommendation_engine.go
"""

import re
import json
from collections import defaultdict

def analyze_recommendation_engine_file(filepath):
    """Analyse le fichier recommendation_engine.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection pour RECOMMENDATION ENGINE
    patterns = {
        'json_field_names': r'"(id|title|description|category|priority|impact|effort|actions|resources|metrics|tags|task|technical|estimated_time)"',
        'priority_levels': r'"(critical|high|medium|low)"',
        'impact_levels': r'"(high|medium|low)"', 
        'effort_levels': r'"(low|medium|high)"',
        'category_names': r'"(tags|performance|security|general)"',
        'template_ids': r'"(missing-title|missing-meta-desc|title-length|meta-desc-length|missing-h1|multiple-h1|heading-hierarchy|missing-alt-tags|missing-viewport|missing-canonical|missing-og-tags|improve-lcp|improve-fid|improve-cls|enable-compression|configure-caching|optimize-images|minify-resources|reduce-page-size|migrate-https|fix-mixed-content|add-hsts|make-responsive|add-sitemap|add-robots-txt|remove-noindex|fix-duplicate-content|improve-accessibility|fix-broken-links|improve-internal-linking|overall-seo-audit)"',
        'score_comparisons': r'"(poor|good|needs_improvement)"',
        'target_values': r'"(≤ 2\.5s|≤ 100ms|≤ 0\.1|< 2 MB|≥ 70%|≥ \d+%)"',
        'time_estimates': r'"(1-2 heures|4-8 heures|1-2 jours|Variable)"',
        'content_ranges': r'"(30-60 caractères|120-160 caractères)"',
        'error_messages': r'"([A-Z][^"]*(?:manquant|manquante|sans|incorrecte)[^"]*)"',
        'recommendation_titles': r'"(Ajouter|Améliorer|Migrer|Configurer|Optimiser|Corriger|Inclure)[^"]*"',
        'action_descriptions': r'"(Ajouter une|Inclure|Respecter|Obtenir|Configurer|Rediriger|Mettre à jour|Optimiser|Améliorer|Précharger|Utiliser)[^"]*"',
        'resource_types': r'"(documentation)"',
        'metric_names': r'"(Taux de clic|Position dans les SERP|Impressions|Trust signals|Ranking boost|Page Experience|Core Web Vitals|LCP)"',
        'tag_names': r'"(critique|balises|onpage|meta|ctr|sécurité|technique|performance|core-web-vitals|lcp)"',
        'numeric_thresholds': r'\b(20|50|70|85|90|95|0\.5|0\.7|1\.0|2\*1024\*1024|100|2\.5|0\.1)\b',
        'placeholders': r'\{[^}]+\}',
        'operators': r'"(≤|≥|<|>)"',
        'seo_issues': r'"([^"]*(?:Titre|Meta|H1|HTTPS|LCP|performance|sécurité)[^"]*)"',
        'url_docs': r'constants\.(Google|WebDev)[^,\s]+',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires et imports
        if line_stripped.startswith('//') or line_stripped.startswith('import') or line_stripped.startswith('package'):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line, re.IGNORECASE)
            for match in matches:
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match,
                    'context': line_stripped[:120] + ('...' if len(line_stripped) > 120 else '')
                })
    
    return violations

def categorize_violations(violations):
    """Catégorise les violations par type"""
    
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    
    return dict(categories)

def generate_constants_mapping(violations):
    """Génère les mappings de constantes pour CHARLIE-1"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_field_map = {
            'id': 'constants.RecommendationJSONFieldID',
            'title': 'constants.RecommendationJSONFieldTitle',
            'description': 'constants.RecommendationJSONFieldDescription',
            'category': 'constants.RecommendationJSONFieldCategory',
            'priority': 'constants.RecommendationJSONFieldPriority',
            'impact': 'constants.RecommendationJSONFieldImpact',
            'effort': 'constants.RecommendationJSONFieldEffort',
            'actions': 'constants.RecommendationJSONFieldActions',
            'resources': 'constants.RecommendationJSONFieldResources',
            'metrics': 'constants.RecommendationJSONFieldMetrics',
            'tags': 'constants.RecommendationJSONFieldTags',
            'task': 'constants.RecommendationJSONFieldTask',
            'technical': 'constants.RecommendationJSONFieldTechnical',
            'estimated_time': 'constants.RecommendationJSONFieldEstimatedTime'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_field_map:
                constants_map[f'"{v["value"]}"'] = json_field_map[v['value']]
    
    # Priority Levels
    if 'priority_levels' in categorized:
        priority_map = {
            'critical': 'constants.RecommendationPriorityCritical',
            'high': 'constants.RecommendationPriorityHigh',
            'medium': 'constants.RecommendationPriorityMedium',
            'low': 'constants.RecommendationPriorityLow'
        }
        for v in categorized['priority_levels']:
            if v['value'] in priority_map:
                constants_map[f'"{v["value"]}"'] = priority_map[v['value']]
    
    # Impact Levels
    if 'impact_levels' in categorized:
        impact_map = {
            'high': 'constants.RecommendationImpactHigh',
            'medium': 'constants.RecommendationImpactMedium',
            'low': 'constants.RecommendationImpactLow'
        }
        for v in categorized['impact_levels']:
            if v['value'] in impact_map:
                constants_map[f'"{v["value"]}"'] = impact_map[v['value']]
    
    # Effort Levels
    if 'effort_levels' in categorized:
        effort_map = {
            'low': 'constants.RecommendationEffortLow',
            'medium': 'constants.RecommendationEffortMedium',
            'high': 'constants.RecommendationEffortHigh'
        }
        for v in categorized['effort_levels']:
            if v['value'] in effort_map:
                constants_map[f'"{v["value"]}"'] = effort_map[v['value']]
    
    # Category Names
    if 'category_names' in categorized:
        category_map = {
            'tags': 'constants.RecommendationCategoryTags',
            'performance': 'constants.RecommendationCategoryPerformance',
            'security': 'constants.RecommendationCategorySecurity',
            'general': 'constants.RecommendationCategoryGeneral'
        }
        for v in categorized['category_names']:
            if v['value'] in category_map:
                constants_map[f'"{v["value"]}"'] = category_map[v['value']]
    
    # Score Comparisons
    if 'score_comparisons' in categorized:
        score_map = {
            'poor': 'constants.RecommendationScorePoor',
            'good': 'constants.RecommendationScoreGood',
            'needs_improvement': 'constants.RecommendationScoreNeedsImprovement'
        }
        for v in categorized['score_comparisons']:
            if v['value'] in score_map:
                constants_map[f'"{v["value"]}"'] = score_map[v['value']]
    
    # Time Estimates
    if 'time_estimates' in categorized:
        time_map = {
            '1-2 heures': 'constants.RecommendationTimeLow',
            '4-8 heures': 'constants.RecommendationTimeMedium',
            '1-2 jours': 'constants.RecommendationTimeHigh',
            'Variable': 'constants.RecommendationTimeVariable'
        }
        for v in categorized['time_estimates']:
            if v['value'] in time_map:
                constants_map[f'"{v["value"]}"'] = time_map[v['value']]
    
    # Content Ranges
    if 'content_ranges' in categorized:
        range_map = {
            '30-60 caractères': 'constants.RecommendationTitleRange',
            '120-160 caractères': 'constants.RecommendationMetaDescRange'
        }
        for v in categorized['content_ranges']:
            if v['value'] in range_map:
                constants_map[f'"{v["value"]}"'] = range_map[v['value']]
    
    # Target Values
    if 'target_values' in categorized:
        target_map = {
            '≤ 2.5s': 'constants.RecommendationTargetLCP',
            '≤ 100ms': 'constants.RecommendationTargetFID',
            '≤ 0.1': 'constants.RecommendationTargetCLS',
            '< 2 MB': 'constants.RecommendationTargetPageSize',
            '≥ 70%': 'constants.RecommendationTargetScore',
            '≥ 3': 'constants.RecommendationTargetLinks'
        }
        for v in categorized['target_values']:
            if v['value'] in target_map:
                constants_map[f'"{v["value"]}"'] = target_map[v['value']]
    
    # Resource Types
    if 'resource_types' in categorized:
        resource_map = {
            'documentation': 'constants.RecommendationResourceTypeDoc'
        }
        for v in categorized['resource_types']:
            if v['value'] in resource_map:
                constants_map[f'"{v["value"]}"'] = resource_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/recommendation_engine.go'
    
    print("🔍 CHARLIE-1 DETECTOR - Scanning recommendation_engine.go...")
    
    violations = analyze_recommendation_engine_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION CHARLIE-1:")
    print(f"Total violations détectées: {len(violations)}")
    
    for category, viols in categorized.items():
        print(f"\n🔸 {category.upper()}: {len(viols)} violations")
        for v in viols[:3]:  # Show first 3 of each category
            print(f"  Line {v['line']}: {v['value']}")
        if len(viols) > 3:
            print(f"  ... et {len(viols) - 3} autres")
    
    print(f"\n🏗️ CONSTANTES À CRÉER: {len(constants_map)}")
    print("Preview des mappings:")
    for original, constant in list(constants_map.items())[:10]:
        print(f"  {original} → {constant}")
    
    # Sauvegarder les résultats
    results = {
        'total_violations': len(violations),
        'categories': {k: len(v) for k, v in categorized.items()},
        'violations': violations,
        'constants_mapping': constants_map
    }
    
    with open('charlie1_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans charlie1_analysis.json")
    return results

if __name__ == "__main__":
    main()