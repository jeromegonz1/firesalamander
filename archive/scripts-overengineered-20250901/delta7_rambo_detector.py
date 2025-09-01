#!/usr/bin/env python3
"""
🔍 DELTA-7 RAMBO DETECTOR
Mode Rambo - Détection ultra-rapide pour seo_test.go - 94 violations
"""

import re
import json
from collections import defaultdict

def analyze_seo_test_file(filepath):
    """Mode Rambo - Analyse ultra-rapide du fichier seo_test.go"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns Rambo pour SEO TESTING - Mode Killer
    patterns = {
        'seo_metrics': r'"(title_length|meta_description_length|h1_count|h2_count|img_without_alt|keyword_density|internal_links|external_links|word_count|readability_score|page_speed|mobile_friendly|schema_markup|canonical_url|sitemap_indexed|robots_indexed)"',
        'html_tags': r'"(title|meta|h1|h2|h3|h4|h5|h6|img|a|p|div|span|article|section|header|footer|nav|main|aside)"',
        'meta_names': r'"(description|keywords|author|viewport|robots|canonical|og:title|og:description|og:image|og:url|og:type|twitter:card|twitter:title|twitter:description|twitter:image)"',
        'seo_attributes': r'"(alt|title|href|src|content|name|property|rel|hreflang|lang|role|aria-label|aria-describedby)"',
        'seo_values': r'"(index|noindex|follow|nofollow|noimageindex|noarchive|nosnippet|max-snippet|max-image-preview|max-video-preview)"',
        'page_types': r'"(website|article|product|blog|homepage|category|search|contact|about|privacy|terms)"',
        'test_urls': r'"https?://[^"]*"',
        'content_types': r'"(text/html|application/json|application/xml|text/css|text/javascript|image/[^"]*)"',
        'http_status_codes': r'"?(200|201|404|500|301|302)"?',
        'seo_tools': r'"(lighthouse|pagespeed|screaming_frog|semrush|ahrefs|moz|yoast|rankmath)"',
        'schema_types': r'"(Article|BlogPosting|Product|Organization|Person|WebPage|WebSite|BreadcrumbList|FAQPage|HowTo|Recipe|Review|Event|LocalBusiness)"',
        'performance_metrics': r'"(lcp|fid|cls|fcp|ttfb|tti|tbt|speed_index|first_paint|largest_paint)"',
        'seo_errors': r'"[A-Z][^"]*(?:missing|empty|too long|too short|duplicate|invalid|not found|error)[^"]*"',
        'seo_recommendations': r'"[A-Z][^"]*(?:add|optimize|improve|reduce|increase|fix|update|include)[^"]*"',
        'test_scenarios': r'"(should_[a-z_]*|test_[a-z_]*|when_[a-z_]*|given_[a-z_]*|verify_[a-z_]*)"',
        'seo_categories': r'"(technical|content|performance|mobile|security|social|local|international)"',
        'audit_levels': r'"(critical|high|medium|low|info|warning|error|success)"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Skip comments, imports, struct tags - Mode Rambo
        if (line_stripped.startswith('//') or 
            line_stripped.startswith('import') or 
            line_stripped.startswith('package') or
            '`json:' in line_stripped):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped, re.IGNORECASE)
            for match in matches:
                if isinstance(match, tuple):
                    match_value = match[0] if match[0] else match[1] if len(match) > 1 else str(match)
                else:
                    match_value = match
                    
                if len(match_value.strip('"')) < 2:
                    continue
                    
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match_value,
                    'context': line_stripped[:120] + ('...' if len(line_stripped) > 120 else '')
                })
    
    return violations

def generate_seo_constants_mapping(violations):
    """Génération Rambo des mappings SEO constants"""
    
    constants_map = {}
    categorized = defaultdict(list)
    for violation in violations:
        categorized[violation['category']].append(violation)
    
    # SEO Metrics - Mode Killer
    if 'seo_metrics' in categorized:
        metrics_map = {
            'title_length': 'constants.SEOMetricTitleLength',
            'meta_description_length': 'constants.SEOMetricMetaDescriptionLength',
            'h1_count': 'constants.SEOMetricH1Count',
            'h2_count': 'constants.SEOMetricH2Count',
            'img_without_alt': 'constants.SEOMetricImagesWithoutAlt',
            'keyword_density': 'constants.SEOMetricKeywordDensity',
            'internal_links': 'constants.SEOMetricInternalLinks',
            'external_links': 'constants.SEOMetricExternalLinks',
            'word_count': 'constants.SEOMetricWordCount',
            'readability_score': 'constants.SEOMetricReadabilityScore',
            'page_speed': 'constants.SEOMetricPageSpeed',
            'mobile_friendly': 'constants.SEOMetricMobileFriendly',
            'schema_markup': 'constants.SEOMetricSchemaMarkup',
            'canonical_url': 'constants.SEOMetricCanonicalURL',
            'sitemap_indexed': 'constants.SEOMetricSitemapIndexed',
            'robots_indexed': 'constants.SEOMetricRobotsIndexed'
        }
        for v in categorized['seo_metrics']:
            if v['value'] in metrics_map:
                constants_map[f'"{v["value"]}"'] = metrics_map[v['value']]
    
    # HTML Tags - Mode Rambo
    if 'html_tags' in categorized:
        tags_map = {
            'title': 'constants.SEOHTMLTagTitle',
            'meta': 'constants.SEOHTMLTagMeta',
            'h1': 'constants.SEOHTMLTagH1',
            'h2': 'constants.SEOHTMLTagH2',
            'h3': 'constants.SEOHTMLTagH3',
            'h4': 'constants.SEOHTMLTagH4',
            'h5': 'constants.SEOHTMLTagH5',
            'h6': 'constants.SEOHTMLTagH6',
            'img': 'constants.SEOHTMLTagImg',
            'a': 'constants.SEOHTMLTagA',
            'p': 'constants.SEOHTMLTagP',
            'div': 'constants.SEOHTMLTagDiv',
            'span': 'constants.SEOHTMLTagSpan',
            'article': 'constants.SEOHTMLTagArticle',
            'section': 'constants.SEOHTMLTagSection',
            'header': 'constants.SEOHTMLTagHeader',
            'footer': 'constants.SEOHTMLTagFooter',
            'nav': 'constants.SEOHTMLTagNav',
            'main': 'constants.SEOHTMLTagMain',
            'aside': 'constants.SEOHTMLTagAside'
        }
        for v in categorized['html_tags']:
            if v['value'] in tags_map:
                constants_map[f'"{v["value"]}"'] = tags_map[v['value']]
    
    # Meta Names - Mode Killer
    if 'meta_names' in categorized:
        meta_map = {
            'description': 'constants.SEOMetaNameDescription',
            'keywords': 'constants.SEOMetaNameKeywords',
            'author': 'constants.SEOMetaNameAuthor',
            'viewport': 'constants.SEOMetaNameViewport',
            'robots': 'constants.SEOMetaNameRobots',
            'canonical': 'constants.SEOMetaNameCanonical',
            'og:title': 'constants.SEOMetaNameOGTitle',
            'og:description': 'constants.SEOMetaNameOGDescription',
            'og:image': 'constants.SEOMetaNameOGImage',
            'og:url': 'constants.SEOMetaNameOGURL',
            'og:type': 'constants.SEOMetaNameOGType',
            'twitter:card': 'constants.SEOMetaNameTwitterCard',
            'twitter:title': 'constants.SEOMetaNameTwitterTitle',
            'twitter:description': 'constants.SEOMetaNameTwitterDescription',
            'twitter:image': 'constants.SEOMetaNameTwitterImage'
        }
        for v in categorized['meta_names']:
            if v['value'] in meta_map:
                constants_map[f'"{v["value"]}"'] = meta_map[v['value']]
    
    # SEO Values - Mode Rambo
    if 'seo_values' in categorized:
        values_map = {
            'index': 'constants.SEOValueIndex',
            'noindex': 'constants.SEOValueNoIndex',
            'follow': 'constants.SEOValueFollow',
            'nofollow': 'constants.SEOValueNoFollow',
            'noimageindex': 'constants.SEOValueNoImageIndex',
            'noarchive': 'constants.SEOValueNoArchive',
            'nosnippet': 'constants.SEOValueNoSnippet',
            'max-snippet': 'constants.SEOValueMaxSnippet',
            'max-image-preview': 'constants.SEOValueMaxImagePreview',
            'max-video-preview': 'constants.SEOValueMaxVideoPreview'
        }
        for v in categorized['seo_values']:
            if v['value'] in values_map:
                constants_map[f'"{v["value"]}"'] = values_map[v['value']]
    
    # Schema Types - Mode Killer
    if 'schema_types' in categorized:
        schema_map = {
            'Article': 'constants.SEOSchemaTypeArticle',
            'BlogPosting': 'constants.SEOSchemaTypeBlogPosting',
            'Product': 'constants.SEOSchemaTypeProduct',
            'Organization': 'constants.SEOSchemaTypeOrganization',
            'Person': 'constants.SEOSchemaTypePerson',
            'WebPage': 'constants.SEOSchemaTypeWebPage',
            'WebSite': 'constants.SEOSchemaTypeWebSite',
            'BreadcrumbList': 'constants.SEOSchemaTypeBreadcrumbList',
            'FAQPage': 'constants.SEOSchemaTypeFAQPage',
            'HowTo': 'constants.SEOSchemaTypeHowTo',
            'Recipe': 'constants.SEOSchemaTypeRecipe',
            'Review': 'constants.SEOSchemaTypeReview',
            'Event': 'constants.SEOSchemaTypeEvent',
            'LocalBusiness': 'constants.SEOSchemaTypeLocalBusiness'
        }
        for v in categorized['schema_types']:
            if v['value'] in schema_map:
                constants_map[f'"{v["value"]}"'] = schema_map[v['value']]
    
    # Performance Metrics - Mode Rambo
    if 'performance_metrics' in categorized:
        perf_map = {
            'lcp': 'constants.SEOPerformanceMetricLCP',
            'fid': 'constants.SEOPerformanceMetricFID',
            'cls': 'constants.SEOPerformanceMetricCLS',
            'fcp': 'constants.SEOPerformanceMetricFCP',
            'ttfb': 'constants.SEOPerformanceMetricTTFB',
            'tti': 'constants.SEOPerformanceMetricTTI',
            'tbt': 'constants.SEOPerformanceMetricTBT',
            'speed_index': 'constants.SEOPerformanceMetricSpeedIndex',
            'first_paint': 'constants.SEOPerformanceMetricFirstPaint',
            'largest_paint': 'constants.SEOPerformanceMetricLargestPaint'
        }
        for v in categorized['performance_metrics']:
            if v['value'] in perf_map:
                constants_map[f'"{v["value"]}"'] = perf_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/seo_test.go'
    
    print("🔥 DELTA-7 RAMBO DETECTOR - Mode Killer Activé!")
    print("=" * 70)
    print("🎯 FIRST BLOOD - Scanning seo_test.go...")
    
    violations = analyze_seo_test_file(filepath)
    constants_map = generate_seo_constants_mapping(violations)
    
    print(f"\n💀 RAMBO RÉSULTATS:")
    print(f"🎯 Violations exterminées: {len(violations)}")
    print(f"🏗️ Constantes à déployer: {len(constants_map)}")
    
    # Show categories breakdown - Mode Rambo
    categorized = defaultdict(list)
    for violation in violations:
        categorized[violation['category']].append(violation)
    
    print(f"\n⚔️ RAMBO BREAKDOWN:")
    for category, viols in categorized.items():
        print(f"🔸 {category.upper()}: {len(viols)} kills")
        for v in viols[:2]:  # Show first 2 of each category - Mode Rambo
            print(f"  💥 Line {v['line']}: {v['value']}")
        if len(viols) > 2:
            print(f"  🎯 ... et {len(viols) - 2} autres cibles")
    
    print(f"\n🏗️ RAMBO CONSTANTS PREVIEW:")
    for original, constant in list(constants_map.items())[:8]:
        print(f"  💀 {original} → {constant}")
    
    # Save results - Mode Rambo
    results = {
        'total_violations': len(violations),
        'categories': {k: len(v) for k, v in categorized.items()},
        'violations': violations,
        'constants_mapping': constants_map
    }
    
    with open('delta7_rambo_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n🎯 RAMBO DATA EXTRACTED: delta7_rambo_analysis.json")
    print("💀 READY FOR ELIMINATION - NOTHING LEFT STANDING!")
    return results

if __name__ == "__main__":
    main()