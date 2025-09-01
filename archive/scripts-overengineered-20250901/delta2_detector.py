#!/usr/bin/env python3
"""
🔍 DELTA-2 DETECTOR
Détection spécialisée pour api.go - 146 violations
"""

import re
import json
from collections import defaultdict

def analyze_api_file(filepath):
    """Analyse le fichier api.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns spécialisés pour API
    patterns = {
        'http_endpoints': r'\"/(api/[^\"]*|health|status|debug|metrics|analyze|results|reports|[^\"]+)\"',
        'http_methods': r'\"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\"',
        'content_types': r'\"(application/json|text/html|text/plain|application/xml|text/csv|application/pdf)\"',
        'http_headers': r'\"(Content-Type|Authorization|Accept|User-Agent|X-[^\"]*|Cache-Control|Access-Control-[^\"]*|Origin)\"',
        'json_field_names': r'\"(id|name|type|status|data|message|error|timestamp|url|method|path|results|recommendations|metrics|score|title|description|content|value|category|priority|severity|phase|agent|test|check|crawler|seo|semantic|analysis|report)\"',
        'api_error_messages': r'\"[A-Z][^\"]*(?:error|failed|invalid|missing|not found|unauthorized|forbidden|timeout|conflict)[^\"]*\"',
        'api_success_messages': r'\"[A-Z][^\"]*(?:success|completed|processed|analyzed|generated|created|updated|deleted)[^\"]*\"',
        'log_messages': r'\"[A-Z][^\"]*(?:starting|started|processing|completed|analyzing|crawling|generating|validating)[^\"]*\"',
        'status_values': r'\"(pending|running|completed|failed|success|error|warning|info|healthy|degraded|critical)\"',
        'mime_types': r'\"(image/[^\"]*|video/[^\"]*|audio/[^\"]*|font/[^\"]*|application/[^\"]*|text/[^\"]*)\"',
        'file_extensions': r'\"\\.(json|html|css|js|png|jpg|jpeg|gif|svg|ico|xml|txt|csv|pdf|zip|tar|gz)\"',
        'url_patterns': r'\"https?://[^\"]*\"',
        'config_keys': r'\"(host|port|timeout|api_key|secret|token|database|redis|elasticsearch|openai|anthropic)\"',
        'agent_names': r'\"(crawler|seo|semantic|performance|security|qa|data_integrity|frontend|playwright|k6)\"',
        'analysis_types': r'\"(technical|content|performance|security|accessibility|seo|semantic|structural)\"',
        'report_formats': r'\"(html|json|csv|pdf|xml|markdown)\"',
        'test_categories': r'\"(unit|integration|e2e|performance|security|accessibility|seo)\"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires, imports et struct tags
        if (line_stripped.startswith('//') or 
            line_stripped.startswith('import') or 
            line_stripped.startswith('package') or
            '`json:' in line_stripped):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped)
            for match in matches:
                # Nettoyer le match
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
                    'context': line_stripped[:150] + ('...' if len(line_stripped) > 150 else '')
                })
    
    return violations

def categorize_violations(violations):
    """Catégorise les violations par type"""
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    return dict(categories)

def generate_api_constants_mapping(violations):
    """Génère les mappings de constantes pour api.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # HTTP Endpoints
    if 'http_endpoints' in categorized:
        for v in categorized['http_endpoints']:
            value = v['value'].strip('"')
            if value.startswith('/api/'):
                endpoint = value[5:].replace('/', '').replace('-', '').replace('_', '').title()
                const_name = f'constants.APIEndpoint{endpoint}' if endpoint else 'constants.APIEndpointRoot'
            elif value == '/health':
                const_name = 'constants.APIEndpointHealth'
            elif value == '/status':
                const_name = 'constants.APIEndpointStatus'
            elif value == '/debug':
                const_name = 'constants.APIEndpointDebug'
            elif value == '/metrics':
                const_name = 'constants.APIEndpointMetrics'
            elif value == '/analyze':
                const_name = 'constants.APIEndpointAnalyze'
            elif value == '/results':
                const_name = 'constants.APIEndpointResults'
            elif value == '/reports':
                const_name = 'constants.APIEndpointReports'
            else:
                clean_name = value.strip('/').replace('/', '').replace('-', '').replace('_', '').title()
                const_name = f'constants.APIEndpoint{clean_name}'
            
            constants_map[f'"{value}"'] = const_name
    
    # HTTP Methods
    if 'http_methods' in categorized:
        method_map = {
            'GET': 'constants.APIMethodGet',
            'POST': 'constants.APIMethodPost', 
            'PUT': 'constants.APIMethodPut',
            'DELETE': 'constants.APIMethodDelete',
            'PATCH': 'constants.APIMethodPatch',
            'HEAD': 'constants.APIMethodHead',
            'OPTIONS': 'constants.APIMethodOptions'
        }
        for v in categorized['http_methods']:
            if v['value'] in method_map:
                constants_map[f'"{v["value"]}"'] = method_map[v['value']]
    
    # Content Types
    if 'content_types' in categorized:
        content_map = {
            'application/json': 'constants.APIContentTypeJSON',
            'text/html': 'constants.APIContentTypeHTML',
            'text/plain': 'constants.APIContentTypePlain',
            'application/xml': 'constants.APIContentTypeXML',
            'text/csv': 'constants.APIContentTypeCSV',
            'application/pdf': 'constants.APIContentTypePDF'
        }
        for v in categorized['content_types']:
            if v['value'] in content_map:
                constants_map[f'"{v["value"]}"'] = content_map[v['value']]
    
    # HTTP Headers
    if 'http_headers' in categorized:
        header_map = {
            'Content-Type': 'constants.APIHeaderContentType',
            'Authorization': 'constants.APIHeaderAuthorization', 
            'Accept': 'constants.APIHeaderAccept',
            'User-Agent': 'constants.APIHeaderUserAgent',
            'Cache-Control': 'constants.APIHeaderCacheControl',
            'Origin': 'constants.APIHeaderOrigin'
        }
        for v in categorized['http_headers']:
            value = v['value'].strip('"')
            if value in header_map:
                constants_map[f'"{value}"'] = header_map[value]
            elif value.startswith('X-'):
                clean_name = value[2:].replace('-', '').title()
                constants_map[f'"{value}"'] = f'constants.APIHeaderX{clean_name}'
            elif value.startswith('Access-Control-'):
                clean_name = value[15:].replace('-', '').title()
                constants_map[f'"{value}"'] = f'constants.APIHeaderAccessControl{clean_name}'
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_map = {
            'id': 'constants.APIJSONFieldID',
            'name': 'constants.APIJSONFieldName',
            'type': 'constants.APIJSONFieldType',
            'status': 'constants.APIJSONFieldStatus',
            'data': 'constants.APIJSONFieldData',
            'message': 'constants.APIJSONFieldMessage',
            'error': 'constants.APIJSONFieldError',
            'timestamp': 'constants.APIJSONFieldTimestamp',
            'url': 'constants.APIJSONFieldURL',
            'method': 'constants.APIJSONFieldMethod',
            'path': 'constants.APIJSONFieldPath',
            'results': 'constants.APIJSONFieldResults',
            'recommendations': 'constants.APIJSONFieldRecommendations',
            'metrics': 'constants.APIJSONFieldMetrics',
            'score': 'constants.APIJSONFieldScore',
            'title': 'constants.APIJSONFieldTitle',
            'description': 'constants.APIJSONFieldDescription',
            'content': 'constants.APIJSONFieldContent',
            'value': 'constants.APIJSONFieldValue',
            'category': 'constants.APIJSONFieldCategory',
            'priority': 'constants.APIJSONFieldPriority',
            'severity': 'constants.APIJSONFieldSeverity'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_map:
                constants_map[f'"{v["value"]}"'] = json_map[v['value']]
    
    # Status Values
    if 'status_values' in categorized:
        status_map = {
            'pending': 'constants.APIStatusPending',
            'running': 'constants.APIStatusRunning',
            'completed': 'constants.APIStatusCompleted',
            'failed': 'constants.APIStatusFailed',
            'success': 'constants.APIStatusSuccess',
            'error': 'constants.APIStatusError',
            'warning': 'constants.APIStatusWarning',
            'info': 'constants.APIStatusInfo',
            'healthy': 'constants.APIStatusHealthy',
            'degraded': 'constants.APIStatusDegraded',
            'critical': 'constants.APIStatusCritical'
        }
        for v in categorized['status_values']:
            if v['value'] in status_map:
                constants_map[f'"{v["value"]}"'] = status_map[v['value']]
    
    # Agent Names
    if 'agent_names' in categorized:
        agent_map = {
            'crawler': 'constants.APIAgentCrawler',
            'seo': 'constants.APIAgentSEO',
            'semantic': 'constants.APIAgentSemantic',
            'performance': 'constants.APIAgentPerformance',
            'security': 'constants.APIAgentSecurity',
            'qa': 'constants.APIAgentQA',
            'data_integrity': 'constants.APIAgentDataIntegrity',
            'frontend': 'constants.APIAgentFrontend',
            'playwright': 'constants.APIAgentPlaywright',
            'k6': 'constants.APIAgentK6'
        }
        for v in categorized['agent_names']:
            if v['value'] in agent_map:
                constants_map[f'"{v["value"]}"'] = agent_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/api.go'
    
    print("🔍 DELTA-2 DETECTOR - Scanning api.go...")
    
    violations = analyze_api_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_api_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-2:")
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
    
    with open('delta2_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta2_analysis.json")
    return results

if __name__ == "__main__":
    main()