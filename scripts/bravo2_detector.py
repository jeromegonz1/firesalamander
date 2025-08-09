#!/usr/bin/env python3
"""
🔍 BRAVO-2 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans tag_analyzer.go
"""

import re
import json
from collections import defaultdict

def analyze_tag_analyzer_file(filepath):
    """Analyse le fichier tag_analyzer.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection pour TAG ANALYZER
    patterns = {
        'json_field_names': r'"(title|meta_description|headings|meta_tags|images|links|microdata|present|content|length|optimal_length|has_keywords|duplicate_words|issues|recommendations|has_call_to_action|h1_count|h1_content|heading_structure|has_hierarchy|missing_levels|has_robots|robots_content|has_canonical|canonical_url|has_og_tags|og_tags|has_twitter_card|twitter_card|has_viewport|viewport_content|total_images|images_with_alt|alt_text_coverage|missing_alt_images|optimized_formats|lazy_loading|internal_links|external_links|nofollow_links|broken_links|anchor_optimization|has_json_ld|json_ld_types|has_microdata|microdata_types|has_rdfa|property|name)"',
        'meta_names': r'"(description|robots|viewport)"',
        'heading_levels': r'"(h1|h2|h3|h4|h5|h6)"',
        'html_attributes': r'"(src|alt|loading|href|rel|type|property|content|name)"',
        'html_values': r'"(lazy|canonical|nofollow|application/ld\+json|og:|twitter:|itemscope|typeof|vocab)"',
        'error_messages': r'"([A-Z][^"]*(?:manquant|manquante|trop|incorrect|sans)[^"]*)"',
        'recommendations': r'"(Ajouter|Étendre|Raccourcir|Éviter|Respecter|Implémenter|Optimiser|Améliorer)[^"]*"',
        'call_to_action_words': r'"(découvrir|contacter|commander|acheter|télécharger|essayer|commencer|cliquer|discover|contact|order|buy|download|try|start|click)"',
        'bad_anchor_texts': r'"(cliquez ici|click here|lire la suite|read more|ici|here)"',
        'regex_patterns': r'`[^`]+`',
        'length_thresholds': r'\b(30|60|120|160|200|300)\b',
        'file_extensions': r'\\.(webp|avif|jpg|jpeg|png|gif|svg)',
        'url_protocols': r'https?://',
        'magic_numbers': r'\b(0\.7|1\.0|3)\b',
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
    """Génère les mappings de constantes pour BRAVO-2"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_field_map = {
            'title': 'constants.TagJSONFieldTitle',
            'meta_description': 'constants.TagJSONFieldMetaDescription',
            'headings': 'constants.TagJSONFieldHeadings',
            'meta_tags': 'constants.TagJSONFieldMetaTags',
            'images': 'constants.TagJSONFieldImages',
            'links': 'constants.TagJSONFieldLinks',
            'microdata': 'constants.TagJSONFieldMicrodata',
            'present': 'constants.TagJSONFieldPresent',
            'content': 'constants.TagJSONFieldContent',
            'length': 'constants.TagJSONFieldLength',
            'optimal_length': 'constants.TagJSONFieldOptimalLength',
            'has_keywords': 'constants.TagJSONFieldHasKeywords',
            'duplicate_words': 'constants.TagJSONFieldDuplicateWords',
            'issues': 'constants.TagJSONFieldIssues',
            'recommendations': 'constants.TagJSONFieldRecommendations',
            'has_call_to_action': 'constants.TagJSONFieldHasCallToAction',
            'h1_count': 'constants.TagJSONFieldH1Count',
            'h1_content': 'constants.TagJSONFieldH1Content',
            'heading_structure': 'constants.TagJSONFieldHeadingStructure',
            'has_hierarchy': 'constants.TagJSONFieldHasHierarchy',
            'missing_levels': 'constants.TagJSONFieldMissingLevels',
            'has_robots': 'constants.TagJSONFieldHasRobots',
            'robots_content': 'constants.TagJSONFieldRobotsContent',
            'has_canonical': 'constants.TagJSONFieldHasCanonical',
            'canonical_url': 'constants.TagJSONFieldCanonicalURL',
            'has_og_tags': 'constants.TagJSONFieldHasOGTags',
            'og_tags': 'constants.TagJSONFieldOGTags',
            'has_twitter_card': 'constants.TagJSONFieldHasTwitterCard',
            'twitter_card': 'constants.TagJSONFieldTwitterCard',
            'has_viewport': 'constants.TagJSONFieldHasViewport',
            'viewport_content': 'constants.TagJSONFieldViewportContent',
            'total_images': 'constants.TagJSONFieldTotalImages',
            'images_with_alt': 'constants.TagJSONFieldImagesWithAlt',
            'alt_text_coverage': 'constants.TagJSONFieldAltTextCoverage',
            'missing_alt_images': 'constants.TagJSONFieldMissingAltImages',
            'optimized_formats': 'constants.TagJSONFieldOptimizedFormats',
            'lazy_loading': 'constants.TagJSONFieldLazyLoading',
            'internal_links': 'constants.TagJSONFieldInternalLinks',
            'external_links': 'constants.TagJSONFieldExternalLinks',
            'nofollow_links': 'constants.TagJSONFieldNoFollowLinks',
            'broken_links': 'constants.TagJSONFieldBrokenLinks',
            'anchor_optimization': 'constants.TagJSONFieldAnchorOptimization',
            'has_json_ld': 'constants.TagJSONFieldHasJSONLD',
            'json_ld_types': 'constants.TagJSONFieldJSONLDTypes',
            'has_microdata': 'constants.TagJSONFieldHasMicrodata',
            'microdata_types': 'constants.TagJSONFieldMicrodataTypes',
            'has_rdfa': 'constants.TagJSONFieldHasRDFa',
            'property': 'constants.TagJSONFieldProperty',
            'name': 'constants.TagJSONFieldName'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_field_map:
                constants_map[f'"{v["value"]}"'] = json_field_map[v['value']]
    
    # Meta Names
    if 'meta_names' in categorized:
        meta_name_map = {
            'description': 'constants.TagMetaNameDescription',
            'robots': 'constants.TagMetaNameRobots',
            'viewport': 'constants.TagMetaNameViewport'
        }
        for v in categorized['meta_names']:
            if v['value'] in meta_name_map:
                constants_map[f'"{v["value"]}"'] = meta_name_map[v['value']]
    
    # HTML Attributes
    if 'html_attributes' in categorized:
        attr_map = {
            'src': 'constants.TagAttrSrc',
            'alt': 'constants.TagAttrAlt',
            'loading': 'constants.TagAttrLoading',
            'href': 'constants.TagAttrHref',
            'rel': 'constants.TagAttrRel',
            'type': 'constants.TagAttrType',
            'property': 'constants.TagAttrProperty',
            'content': 'constants.TagAttrContent',
            'name': 'constants.TagAttrName'
        }
        for v in categorized['html_attributes']:
            if v['value'] in attr_map:
                constants_map[f'"{v["value"]}"'] = attr_map[v['value']]
    
    # HTML Values
    if 'html_values' in categorized:
        value_map = {
            'lazy': 'constants.TagValueLazy',
            'canonical': 'constants.TagValueCanonical',
            'nofollow': 'constants.TagValueNoFollow',
            'application/ld+json': 'constants.TagValueJSONLD',
            'og:': 'constants.TagPrefixOG',
            'twitter:': 'constants.TagPrefixTwitter',
            'itemscope': 'constants.TagValueItemScope',
            'typeof': 'constants.TagValueTypeOf',
            'vocab': 'constants.TagValueVocab'
        }
        for v in categorized['html_values']:
            if v['value'] in value_map:
                constants_map[f'"{v["value"]}"'] = value_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/tag_analyzer.go'
    
    print("🔍 BRAVO-2 DETECTOR - Scanning tag_analyzer.go...")
    
    violations = analyze_tag_analyzer_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION BRAVO-2:")
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
    
    with open('bravo2_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans bravo2_analysis.json")
    return results

if __name__ == "__main__":
    main()