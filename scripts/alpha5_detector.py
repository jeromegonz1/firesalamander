#!/usr/bin/env python3
"""
🔍 ALPHA-5 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans server.go
"""

import re
import json
from collections import defaultdict

def analyze_server_file(filepath):
    """Analyse le fichier server.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection
    patterns = {
        'error_messages': r'"([A-Z][^"]*(?:erreur|error|non disponible|not found|not allowed|requis|required)[^"]*)"',
        'content_types': r'"(text/html|application/json|application/pdf|text/csv|application/octet-stream)"',
        'http_headers': r'"(Content-Type|Cache-Control|X-Frame-Options|X-Content-Type-Options|Content-Length|Access-Control-[^"]+|Content-Disposition)"',
        'status_values': r'"(healthy|degraded|unavailable|running)"',
        'http_methods': r'"(GET|POST|PUT|DELETE|OPTIONS)"',
        'api_routes': r'"/api/v1/[^"]*"',
        'web_routes': r'"/web/[^"]*"',
        'static_routes': r'"/static[^"]*"',
        'file_extensions': r'"\.(html|pdf|json|csv)"',
        'service_names': r'"([^"]*Fire Salamander[^"]*)"',
        'log_messages': r'"([^"]*(?:enregistrées|servie|demandé|démarré)[^"]*)"',
        'html_templates': r'<!DOCTYPE html>|<html[^>]*>|<head[^>]*>|<body[^>]*>',
        'css_styles': r'font-family:|font-size:|color:|background:|border:|margin:|padding:',
        'json_keys': r'"(status|service|version|timestamp|components|title|description|priority|impact|effort)"',
        'hardcoded_values': r'\b(87|85|72|90|88|15\.2s|"info"|"high"|"medium"|"low"|"stable"|"Succès")\b',
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
    """Génère les mappings de constantes pour ALPHA-5"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Error Messages
    if 'error_messages' in categorized:
        for v in categorized['error_messages']:
            value = v['value']
            if 'Erreur accès fichiers statiques' in value:
                constants_map[f'"{value}"'] = 'constants.MsgErrorStaticFiles'
            elif 'Interface web non disponible' in value:
                constants_map[f'"{value}"'] = 'constants.MsgWebInterfaceUnavailable'
            elif 'Method not allowed' in value:
                constants_map[f'"{value}"'] = 'constants.MsgMethodNotAllowed'
            elif 'API endpoint not found' in value:
                constants_map[f'"{value}"'] = 'constants.MsgAPIEndpointNotFound'
            elif 'Nom de fichier requis' in value:
                constants_map[f'"{value}"'] = 'constants.MsgFilenameRequired'
    
    # Content Types
    if 'content_types' in categorized:
        content_type_map = {
            'text/html': 'constants.ContentTypeHTML',
            'application/json': 'constants.ContentTypeJSON',
            'application/pdf': 'constants.ContentTypePDF',
            'text/csv': 'constants.ContentTypeCSV',
            'application/octet-stream': 'constants.ContentTypeOctetStream'
        }
        for v in categorized['content_types']:
            if v['value'] in content_type_map:
                constants_map[f'"{v["value"]}"'] = content_type_map[v['value']]
    
    # HTTP Headers
    if 'http_headers' in categorized:
        header_map = {
            'Content-Type': 'constants.HeaderContentType',
            'Cache-Control': 'constants.HeaderCacheControl',
            'X-Frame-Options': 'constants.HeaderXFrameOptions',
            'X-Content-Type-Options': 'constants.HeaderXContentType',
            'Content-Length': 'constants.HeaderContentLength',
            'Access-Control-Allow-Origin': 'constants.HeaderCORSOrigin',
            'Access-Control-Allow-Methods': 'constants.HeaderCORSMethods',
            'Access-Control-Allow-Headers': 'constants.HeaderCORSHeaders',
            'Content-Disposition': 'constants.HeaderContentDisposition'
        }
        for v in categorized['http_headers']:
            if v['value'] in header_map:
                constants_map[f'"{v["value"]}"'] = header_map[v['value']]
    
    # Status Values
    if 'status_values' in categorized:
        status_map = {
            'healthy': 'constants.HealthStatusHealthy',
            'degraded': 'constants.HealthStatusDegraded',
            'unavailable': 'constants.HealthStatusUnavailable',
            'running': 'constants.ServiceStatusRunning'
        }
        for v in categorized['status_values']:
            if v['value'] in status_map:
                constants_map[f'"{v["value"]}"'] = status_map[v['value']]
    
    # HTTP Methods
    if 'http_methods' in categorized:
        method_map = {
            'GET': 'constants.HTTPMethodGET',
            'POST': 'constants.HTTPMethodPOST',
            'PUT': 'constants.HTTPMethodPUT',
            'DELETE': 'constants.HTTPMethodDELETE',
            'OPTIONS': 'constants.HTTPMethodOPTIONS'
        }
        for v in categorized['http_methods']:
            if v['value'] in method_map:
                constants_map[f'"{v["value"]}"'] = method_map[v['value']]
    
    # API Routes
    if 'api_routes' in categorized:
        route_map = {
            '/api/v1/analyze': 'constants.APIRouteAnalyze',
            '/api/v1/analyze/semantic': 'constants.APIRouteAnalyzeSemantic',
            '/api/v1/analyze/seo': 'constants.APIRouteAnalyzeSEO',
            '/api/v1/analyze/quick': 'constants.APIRouteAnalyzeQuick',
            '/api/v1/health': 'constants.APIRouteHealth',
            '/api/v1/stats': 'constants.APIRouteStats',
            '/api/v1/analyses': 'constants.APIRouteAnalyses',
            '/api/v1/analysis/': 'constants.APIRouteAnalysisDetails',
            '/api/v1/info': 'constants.APIRouteInfo',
            '/api/v1/version': 'constants.APIRouteVersion'
        }
        for v in categorized['api_routes']:
            if v['value'] in route_map:
                constants_map[f'"{v["value"]}"'] = route_map[v['value']]
    
    # Web Routes
    if 'web_routes' in categorized:
        web_route_map = {
            '/web/health': 'constants.WebRouteHealth',
            '/web/download/': 'constants.WebRouteDownload'
        }
        for v in categorized['web_routes']:
            if v['value'] in web_route_map:
                constants_map[f'"{v["value"]}"'] = web_route_map[v['value']]
    
    # Static Routes
    if 'static_routes' in categorized:
        constants_map['"/static/"'] = 'constants.StaticRoute'
    
    # File Extensions
    if 'file_extensions' in categorized:
        ext_map = {
            '.html': 'constants.ExtHTML',
            '.pdf': 'constants.ExtPDF', 
            '.json': 'constants.ExtJSON',
            '.csv': 'constants.ExtCSV'
        }
        for v in categorized['file_extensions']:
            if v['value'] in ext_map:
                constants_map[f'"{v["value"]}"'] = ext_map[v['value']]
    
    # Service Names
    if 'service_names' in categorized:
        for v in categorized['service_names']:
            if 'Fire Salamander Web Server' in v['value']:
                constants_map[f'"{v["value"]}"'] = 'constants.ServiceNameWebServer'
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/server.go'
    
    print("🔍 ALPHA-5 DETECTOR - Scanning server.go...")
    
    violations = analyze_server_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION ALPHA-5:")
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
    
    with open('alpha5_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans alpha5_analysis.json")
    return results

if __name__ == "__main__":
    main()