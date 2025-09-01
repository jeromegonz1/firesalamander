#!/usr/bin/env python3
"""
🔍 DELTA-1 DETECTOR
Détection spécialisée pour server.go - 146 violations
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
    
    # Patterns spécialisés pour SERVER
    patterns = {
        'http_endpoints': r'"/(api/[^"]*|health|metrics|status|debug|[^"]*)"',
        'http_methods': r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'content_types': r'"(application/json|text/html|text/plain|application/xml|multipart/form-data)"',
        'error_messages': r'"([A-Z][^"]*(?:error|failed|unable|invalid|missing|not found|denied)[^"]*)"',
        'log_messages': r'"([A-Z][^"]*(?:starting|started|stopping|stopped|listening|serving|processing)[^"]*)"',
        'server_config': r'"(host|port|timeout|max_connections|buffer_size|read_timeout|write_timeout)"',
        'http_headers': r'"(Content-Type|Authorization|Accept|User-Agent|X-[^"]*|Cache-Control)"',
        'status_codes': r'"(200|201|400|401|403|404|500|502|503)"',
        'mime_extensions': r'"\.(json|html|css|js|png|jpg|jpeg|gif|svg|ico|xml|txt)"',
        'file_paths': r'"(./|/|\.\./)([^"]*/?)*"',
        'template_names': r'"[^"]*\.(?:html|tmpl|tpl)"',
        'cors_origins': r'"https?://[^"]*"',
        'middleware_names': r'"(cors|auth|logging|recovery|compression|rate-limit)"',
        'json_field_names': r'"(id|name|type|status|data|message|error|timestamp|url|method|path)"',
        'environment_vars': r'"[A-Z][A-Z0-9_]*"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires et imports
        if line_stripped.startswith('//') or line_stripped.startswith('import') or line_stripped.startswith('package'):
            continue
            
        # Ignorer les struct tags
        if '`json:' in line_stripped:
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped)
            for match in matches:
                # Nettoyer le match (prendre le premier groupe si c'est un tuple)
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

def categorize_violations(violations):
    """Catégorise les violations par type"""
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    return dict(categories)

def generate_server_constants_mapping(violations):
    """Génère les mappings de constantes pour server.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # HTTP Endpoints
    if 'http_endpoints' in categorized:
        endpoint_map = {}
        for v in categorized['http_endpoints']:
            value = v['value'].strip('"')
            if value.startswith('/api/'):
                const_name = f'constants.ServerEndpointAPI{value[5:].replace("/", "").replace("-", "").title()}'
            elif value == '/health':
                const_name = 'constants.ServerEndpointHealth'
            elif value == '/metrics':
                const_name = 'constants.ServerEndpointMetrics'
            elif value == '/status':
                const_name = 'constants.ServerEndpointStatus'
            elif value == '/debug':
                const_name = 'constants.ServerEndpointDebug'
            else:
                clean_name = value.strip('/').replace('/', '').replace('-', '').title()
                const_name = f'constants.ServerEndpoint{clean_name}'
            
            constants_map[f'"{value}"'] = const_name
    
    # HTTP Methods
    if 'http_methods' in categorized:
        method_map = {
            'GET': 'constants.ServerMethodGet',
            'POST': 'constants.ServerMethodPost',
            'PUT': 'constants.ServerMethodPut',
            'DELETE': 'constants.ServerMethodDelete',
            'PATCH': 'constants.ServerMethodPatch',
            'HEAD': 'constants.ServerMethodHead',
            'OPTIONS': 'constants.ServerMethodOptions'
        }
        for v in categorized['http_methods']:
            if v['value'] in method_map:
                constants_map[f'"{v["value"]}"'] = method_map[v['value']]
    
    # Content Types
    if 'content_types' in categorized:
        content_map = {
            'application/json': 'constants.ServerContentTypeJSON',
            'text/html': 'constants.ServerContentTypeHTML',
            'text/plain': 'constants.ServerContentTypePlain',
            'application/xml': 'constants.ServerContentTypeXML',
            'multipart/form-data': 'constants.ServerContentTypeFormData'
        }
        for v in categorized['content_types']:
            if v['value'] in content_map:
                constants_map[f'"{v["value"]}"'] = content_map[v['value']]
    
    # HTTP Headers
    if 'http_headers' in categorized:
        header_map = {
            'Content-Type': 'constants.ServerHeaderContentType',
            'Authorization': 'constants.ServerHeaderAuthorization',
            'Accept': 'constants.ServerHeaderAccept',
            'User-Agent': 'constants.ServerHeaderUserAgent',
            'Cache-Control': 'constants.ServerHeaderCacheControl'
        }
        for v in categorized['http_headers']:
            value = v['value'].strip('"')
            if value in header_map:
                constants_map[f'"{value}"'] = header_map[value]
            elif value.startswith('X-'):
                clean_name = value[2:].replace('-', '')
                constants_map[f'"{value}"'] = f'constants.ServerHeaderX{clean_name}'
    
    # Server Config Keys
    if 'server_config' in categorized:
        config_map = {
            'host': 'constants.ServerConfigHost',
            'port': 'constants.ServerConfigPort',
            'timeout': 'constants.ServerConfigTimeout',
            'max_connections': 'constants.ServerConfigMaxConnections',
            'buffer_size': 'constants.ServerConfigBufferSize',
            'read_timeout': 'constants.ServerConfigReadTimeout',
            'write_timeout': 'constants.ServerConfigWriteTimeout'
        }
        for v in categorized['server_config']:
            if v['value'] in config_map:
                constants_map[f'"{v["value"]}"'] = config_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/server.go'
    
    print("🔍 DELTA-1 DETECTOR - Scanning server.go...")
    
    violations = analyze_server_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_server_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-1:")
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
    
    with open('delta1_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta1_analysis.json")
    return results

if __name__ == "__main__":
    main()