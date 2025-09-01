#!/usr/bin/env python3
"""
🔍 DELTA-4 DETECTOR
Détection spécialisée pour integration_test.go - 108 violations
"""

import re
import json
from collections import defaultdict

def analyze_integration_test_file(filepath):
    """Analyse le fichier integration_test.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns spécialisés pour INTEGRATION TESTING
    patterns = {
        'test_names': r'"(Test[A-Z][^"]*|test_[a-z_]*)"',
        'test_scenarios': r'"(should_[^"]*|when_[^"]*|given_[^"]*|then_[^"]*|scenario_[^"]*|case_[^"]*)"',
        'test_assertions': r'"(expected|actual|assert|verify|validate|check|ensure|confirm|match|equal|contain|include|start|end|greater|less|empty|nil|true|false|success|failure|error|valid|invalid)"',
        'test_phases': r'"(setup|teardown|before|after|arrange|act|assert|given|when|then|prepare|execute|cleanup|initialize|finalize)"',
        'integration_endpoints': r'"/(api/[^"]*|health|status|debug|metrics|analyze|results|reports|test|mock|stub|fake)[^"]*"',
        'http_methods': r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'http_status_codes': r'"(200|201|202|204|301|302|304|400|401|403|404|405|409|422|429|500|501|502|503|504)"',
        'content_types': r'"(application/json|text/html|text/plain|application/xml|text/csv|application/pdf)"',
        'test_data': r'"(mock_[^"]*|test_[^"]*|sample_[^"]*|dummy_[^"]*|fake_[^"]*|example_[^"]*)"',
        'error_messages': r'"[A-Z][^"]*(?:error|failed|invalid|missing|not found|timeout|conflict|denied)[^"]*"',
        'success_messages': r'"[A-Z][^"]*(?:success|completed|passed|valid|correct|found|created|updated|deleted)[^"]*"',
        'log_messages': r'"[A-Z][^"]*(?:starting|started|running|testing|checking|validating|processing|executing)[^"]*"',
        'json_field_names': r'"(id|name|type|status|data|message|error|timestamp|url|method|path|result|response|request|payload|headers|body|test|suite|case|scenario|assertion|expectation|actual|expected|duration|timeout)"',
        'test_config_keys': r'"(timeout|retries|parallel|verbose|debug|mock|stub|fake|cleanup|setup|teardown|database|server|client|host|port)"',
        'agent_names': r'"(crawler|seo|semantic|performance|security|qa|data_integrity|frontend|playwright|k6|integration|unit|e2e)"',
        'analysis_types': r'"(technical|content|performance|security|accessibility|seo|semantic|structural|functional|regression|smoke|sanity)"',
        'test_categories': r'"(unit|integration|e2e|api|ui|performance|security|accessibility|compatibility|regression|smoke|sanity|acceptance|system)"',
        'mock_responses': r'"(mock|stub|fake|test|sample|dummy|example)_response"',
        'test_fixtures': r'"(fixture|testdata|mockdata|sampledata|testcase|scenario|template)"',
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
            matches = re.findall(pattern, line_stripped, re.IGNORECASE)
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

def generate_integration_test_constants_mapping(violations):
    """Génère les mappings de constantes pour integration_test.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Test Names
    if 'test_names' in categorized:
        for v in categorized['test_names']:
            value = v['value'].strip('"')
            if value.startswith('Test'):
                clean_name = value[4:]  # Remove 'Test' prefix
                const_name = f'constants.IntegrationTestName{clean_name}'
            elif value.startswith('test_'):
                clean_name = value[5:].replace('_', '').title()  # Remove 'test_' and convert
                const_name = f'constants.IntegrationTestName{clean_name}'
            else:
                clean_name = value.replace('_', '').title()
                const_name = f'constants.IntegrationTestName{clean_name}'
            
            constants_map[f'"{value}"'] = const_name
    
    # Test Phases
    if 'test_phases' in categorized:
        phase_map = {
            'setup': 'constants.IntegrationTestPhaseSetup',
            'teardown': 'constants.IntegrationTestPhaseTeardown',
            'before': 'constants.IntegrationTestPhaseBefore',
            'after': 'constants.IntegrationTestPhaseAfter',
            'arrange': 'constants.IntegrationTestPhaseArrange',
            'act': 'constants.IntegrationTestPhaseAct',
            'assert': 'constants.IntegrationTestPhaseAssert',
            'given': 'constants.IntegrationTestPhaseGiven',
            'when': 'constants.IntegrationTestPhaseWhen',
            'then': 'constants.IntegrationTestPhaseThen',
            'prepare': 'constants.IntegrationTestPhasePrepare',
            'execute': 'constants.IntegrationTestPhaseExecute',
            'cleanup': 'constants.IntegrationTestPhaseCleanup',
            'initialize': 'constants.IntegrationTestPhaseInitialize',
            'finalize': 'constants.IntegrationTestPhaseFinalize'
        }
        for v in categorized['test_phases']:
            if v['value'] in phase_map:
                constants_map[f'"{v["value"]}"'] = phase_map[v['value']]
    
    # Integration Endpoints
    if 'integration_endpoints' in categorized:
        for v in categorized['integration_endpoints']:
            value = v['value'].strip('"')
            if value.startswith('/api/'):
                endpoint = value[5:].replace('/', '').replace('-', '').replace('_', '').title()
                const_name = f'constants.IntegrationEndpointAPI{endpoint}' if endpoint else 'constants.IntegrationEndpointAPIRoot'
            elif value == '/health':
                const_name = 'constants.IntegrationEndpointHealth'
            elif value == '/status':
                const_name = 'constants.IntegrationEndpointStatus'
            elif value == '/debug':
                const_name = 'constants.IntegrationEndpointDebug'
            elif value == '/metrics':
                const_name = 'constants.IntegrationEndpointMetrics'
            elif value == '/analyze':
                const_name = 'constants.IntegrationEndpointAnalyze'
            elif value == '/results':
                const_name = 'constants.IntegrationEndpointResults'
            elif value == '/reports':
                const_name = 'constants.IntegrationEndpointReports'
            elif value == '/test':
                const_name = 'constants.IntegrationEndpointTest'
            else:
                clean_name = value.strip('/').replace('/', '').replace('-', '').replace('_', '').title()
                const_name = f'constants.IntegrationEndpoint{clean_name}'
            
            constants_map[f'"{value}"'] = const_name
    
    # HTTP Methods
    if 'http_methods' in categorized:
        method_map = {
            'GET': 'constants.IntegrationHTTPMethodGet',
            'POST': 'constants.IntegrationHTTPMethodPost',
            'PUT': 'constants.IntegrationHTTPMethodPut',
            'DELETE': 'constants.IntegrationHTTPMethodDelete',
            'PATCH': 'constants.IntegrationHTTPMethodPatch',
            'HEAD': 'constants.IntegrationHTTPMethodHead',
            'OPTIONS': 'constants.IntegrationHTTPMethodOptions'
        }
        for v in categorized['http_methods']:
            if v['value'] in method_map:
                constants_map[f'"{v["value"]}"'] = method_map[v['value']]
    
    # HTTP Status Codes
    if 'http_status_codes' in categorized:
        status_map = {
            '200': 'constants.IntegrationHTTPStatusOK',
            '201': 'constants.IntegrationHTTPStatusCreated',
            '202': 'constants.IntegrationHTTPStatusAccepted',
            '204': 'constants.IntegrationHTTPStatusNoContent',
            '301': 'constants.IntegrationHTTPStatusMovedPermanently',
            '302': 'constants.IntegrationHTTPStatusFound',
            '304': 'constants.IntegrationHTTPStatusNotModified',
            '400': 'constants.IntegrationHTTPStatusBadRequest',
            '401': 'constants.IntegrationHTTPStatusUnauthorized',
            '403': 'constants.IntegrationHTTPStatusForbidden',
            '404': 'constants.IntegrationHTTPStatusNotFound',
            '405': 'constants.IntegrationHTTPStatusMethodNotAllowed',
            '409': 'constants.IntegrationHTTPStatusConflict',
            '422': 'constants.IntegrationHTTPStatusUnprocessableEntity',
            '429': 'constants.IntegrationHTTPStatusTooManyRequests',
            '500': 'constants.IntegrationHTTPStatusInternalServerError',
            '501': 'constants.IntegrationHTTPStatusNotImplemented',
            '502': 'constants.IntegrationHTTPStatusBadGateway',
            '503': 'constants.IntegrationHTTPStatusServiceUnavailable',
            '504': 'constants.IntegrationHTTPStatusGatewayTimeout'
        }
        for v in categorized['http_status_codes']:
            if v['value'] in status_map:
                constants_map[f'"{v["value"]}"'] = status_map[v['value']]
    
    # Content Types
    if 'content_types' in categorized:
        content_map = {
            'application/json': 'constants.IntegrationContentTypeJSON',
            'text/html': 'constants.IntegrationContentTypeHTML',
            'text/plain': 'constants.IntegrationContentTypePlain',
            'application/xml': 'constants.IntegrationContentTypeXML',
            'text/csv': 'constants.IntegrationContentTypeCSV',
            'application/pdf': 'constants.IntegrationContentTypePDF'
        }
        for v in categorized['content_types']:
            if v['value'] in content_map:
                constants_map[f'"{v["value"]}"'] = content_map[v['value']]
    
    # Test Assertions
    if 'test_assertions' in categorized:
        assertion_map = {
            'expected': 'constants.IntegrationAssertionExpected',
            'actual': 'constants.IntegrationAssertionActual',
            'assert': 'constants.IntegrationAssertionAssert',
            'verify': 'constants.IntegrationAssertionVerify',
            'validate': 'constants.IntegrationAssertionValidate',
            'check': 'constants.IntegrationAssertionCheck',
            'ensure': 'constants.IntegrationAssertionEnsure',
            'confirm': 'constants.IntegrationAssertionConfirm',
            'match': 'constants.IntegrationAssertionMatch',
            'equal': 'constants.IntegrationAssertionEqual',
            'contain': 'constants.IntegrationAssertionContain',
            'include': 'constants.IntegrationAssertionInclude',
            'start': 'constants.IntegrationAssertionStart',
            'end': 'constants.IntegrationAssertionEnd',
            'greater': 'constants.IntegrationAssertionGreater',
            'less': 'constants.IntegrationAssertionLess',
            'empty': 'constants.IntegrationAssertionEmpty',
            'nil': 'constants.IntegrationAssertionNil',
            'true': 'constants.IntegrationAssertionTrue',
            'false': 'constants.IntegrationAssertionFalse',
            'success': 'constants.IntegrationAssertionSuccess',
            'failure': 'constants.IntegrationAssertionFailure',
            'error': 'constants.IntegrationAssertionError',
            'valid': 'constants.IntegrationAssertionValid',
            'invalid': 'constants.IntegrationAssertionInvalid'
        }
        for v in categorized['test_assertions']:
            if v['value'] in assertion_map:
                constants_map[f'"{v["value"]}"'] = assertion_map[v['value']]
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_map = {
            'id': 'constants.IntegrationJSONFieldID',
            'name': 'constants.IntegrationJSONFieldName',
            'type': 'constants.IntegrationJSONFieldType',
            'status': 'constants.IntegrationJSONFieldStatus',
            'data': 'constants.IntegrationJSONFieldData',
            'message': 'constants.IntegrationJSONFieldMessage',
            'error': 'constants.IntegrationJSONFieldError',
            'timestamp': 'constants.IntegrationJSONFieldTimestamp',
            'url': 'constants.IntegrationJSONFieldURL',
            'method': 'constants.IntegrationJSONFieldMethod',
            'path': 'constants.IntegrationJSONFieldPath',
            'result': 'constants.IntegrationJSONFieldResult',
            'response': 'constants.IntegrationJSONFieldResponse',
            'request': 'constants.IntegrationJSONFieldRequest',
            'payload': 'constants.IntegrationJSONFieldPayload',
            'headers': 'constants.IntegrationJSONFieldHeaders',
            'body': 'constants.IntegrationJSONFieldBody',
            'test': 'constants.IntegrationJSONFieldTest',
            'suite': 'constants.IntegrationJSONFieldSuite',
            'case': 'constants.IntegrationJSONFieldCase',
            'scenario': 'constants.IntegrationJSONFieldScenario',
            'assertion': 'constants.IntegrationJSONFieldAssertion',
            'expectation': 'constants.IntegrationJSONFieldExpectation',
            'actual': 'constants.IntegrationJSONFieldActual',
            'expected': 'constants.IntegrationJSONFieldExpected',
            'duration': 'constants.IntegrationJSONFieldDuration',
            'timeout': 'constants.IntegrationJSONFieldTimeout'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_map:
                constants_map[f'"{v["value"]}"'] = json_map[v['value']]
    
    # Test Categories
    if 'test_categories' in categorized:
        category_map = {
            'unit': 'constants.IntegrationTestCategoryUnit',
            'integration': 'constants.IntegrationTestCategoryIntegration',
            'e2e': 'constants.IntegrationTestCategoryE2E',
            'api': 'constants.IntegrationTestCategoryAPI',
            'ui': 'constants.IntegrationTestCategoryUI',
            'performance': 'constants.IntegrationTestCategoryPerformance',
            'security': 'constants.IntegrationTestCategorySecurity',
            'accessibility': 'constants.IntegrationTestCategoryAccessibility',
            'compatibility': 'constants.IntegrationTestCategoryCompatibility',
            'regression': 'constants.IntegrationTestCategoryRegression',
            'smoke': 'constants.IntegrationTestCategorySmoke',
            'sanity': 'constants.IntegrationTestCategorySanity',
            'acceptance': 'constants.IntegrationTestCategoryAcceptance',
            'system': 'constants.IntegrationTestCategorySystem'
        }
        for v in categorized['test_categories']:
            if v['value'] in category_map:
                constants_map[f'"{v["value"]}"'] = category_map[v['value']]
    
    # Agent Names
    if 'agent_names' in categorized:
        agent_map = {
            'crawler': 'constants.IntegrationAgentCrawler',
            'seo': 'constants.IntegrationAgentSEO',
            'semantic': 'constants.IntegrationAgentSemantic',
            'performance': 'constants.IntegrationAgentPerformance',
            'security': 'constants.IntegrationAgentSecurity',
            'qa': 'constants.IntegrationAgentQA',
            'data_integrity': 'constants.IntegrationAgentDataIntegrity',
            'frontend': 'constants.IntegrationAgentFrontend',
            'playwright': 'constants.IntegrationAgentPlaywright',
            'k6': 'constants.IntegrationAgentK6',
            'integration': 'constants.IntegrationAgentIntegration',
            'unit': 'constants.IntegrationAgentUnit',
            'e2e': 'constants.IntegrationAgentE2E'
        }
        for v in categorized['agent_names']:
            if v['value'] in agent_map:
                constants_map[f'"{v["value"]}"'] = agent_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/integration_test.go'
    
    print("🔍 DELTA-4 DETECTOR - Scanning integration_test.go...")
    
    violations = analyze_integration_test_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_integration_test_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-4:")
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
    
    with open('delta4_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta4_analysis.json")
    return results

if __name__ == "__main__":
    main()