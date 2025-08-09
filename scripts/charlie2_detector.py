#!/usr/bin/env python3
"""
🔍 CHARLIE-2 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans orchestrator.go
"""

import re
import json
from collections import defaultdict

def analyze_orchestrator_file(filepath):
    """Analyse le fichier orchestrator.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection pour ORCHESTRATOR
    patterns = {
        'json_field_names': r'"(status|message|timestamp|phase|duration|error|result|data|config|agents|tasks|logs|metrics|progress|state)"',
        'agent_names': r'"(orchestrator|crawler|seo|qa|security|performance|data-integrity|tag-analyzer)"',
        'phase_names': r'"(initialization|crawling|analysis|scoring|reporting|cleanup|validation|testing|security-scan)"',
        'status_values': r'"(pending|running|completed|failed|cancelled|success|error|warning|info)"',
        'message_templates': r'"([A-Z][^"]*(?:started|completed|failed|error|success|progress|phase|agent)[^"]*)"',
        'config_keys': r'"(max_workers|timeout|retry_attempts|batch_size|concurrent_agents|log_level|debug_mode)"',
        'error_messages': r'"([A-Z][^"]*(?:failed|error|unable|cannot|invalid|missing)[^"]*)"',
        'log_messages': r'"([A-Z][^"]*(?:Starting|Executing|Completed|Processing|Initializing|Validating)[^"]*)"',
        'file_paths': r'"(reports/|logs/|temp/|cache/|config/)[^"]*"',
        'url_patterns': r'"(http://|https://|ws://|wss://)[^"]*"',
        'time_formats': r'"(\d{4}-\d{2}-\d{2}|\d{2}:\d{2}:\d{2}|RFC3339)"',
        'numeric_thresholds': r'\b(10|30|60|100|500|1000|5000|10000)\b',
        'timeout_values': r'\b(30|60|120|300|600|1800|3600)\s*\*\s*time\.(Second|Minute)',
        'goroutine_patterns': r'go\s+func\(\)',
        'channel_operations': r'(make\(chan\s+\w+|<-\s*\w+|\w+\s*<-)',
        'context_patterns': r'context\.(Background|TODO|WithTimeout|WithCancel)\(\)',
        'sync_patterns': r'sync\.(WaitGroup|Mutex|RWMutex)',
        'http_methods': r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'content_types': r'"(application/json|text/html|text/plain|application/xml)"',
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
    """Génère les mappings de constantes pour CHARLIE-2"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_field_map = {
            'status': 'constants.OrchestratorJSONFieldStatus',
            'message': 'constants.OrchestratorJSONFieldMessage',
            'timestamp': 'constants.OrchestratorJSONFieldTimestamp',
            'phase': 'constants.OrchestratorJSONFieldPhase',
            'duration': 'constants.OrchestratorJSONFieldDuration',
            'error': 'constants.OrchestratorJSONFieldError',
            'result': 'constants.OrchestratorJSONFieldResult',
            'data': 'constants.OrchestratorJSONFieldData',
            'config': 'constants.OrchestratorJSONFieldConfig',
            'agents': 'constants.OrchestratorJSONFieldAgents',
            'tasks': 'constants.OrchestratorJSONFieldTasks',
            'logs': 'constants.OrchestratorJSONFieldLogs',
            'metrics': 'constants.OrchestratorJSONFieldMetrics',
            'progress': 'constants.OrchestratorJSONFieldProgress',
            'state': 'constants.OrchestratorJSONFieldState'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_field_map:
                constants_map[f'"{v["value"]}"'] = json_field_map[v['value']]
    
    # Agent Names
    if 'agent_names' in categorized:
        agent_map = {
            'orchestrator': 'constants.OrchestratorAgentNameOrchestrator',
            'crawler': 'constants.OrchestratorAgentNameCrawler',
            'seo': 'constants.OrchestratorAgentNameSEO',
            'qa': 'constants.OrchestratorAgentNameQA',
            'security': 'constants.OrchestratorAgentNameSecurity',
            'performance': 'constants.OrchestratorAgentNamePerformance',
            'data-integrity': 'constants.OrchestratorAgentNameDataIntegrity',
            'tag-analyzer': 'constants.OrchestratorAgentNameTagAnalyzer'
        }
        for v in categorized['agent_names']:
            if v['value'] in agent_map:
                constants_map[f'"{v["value"]}"'] = agent_map[v['value']]
    
    # Phase Names
    if 'phase_names' in categorized:
        phase_map = {
            'initialization': 'constants.OrchestratorPhaseInitialization',
            'crawling': 'constants.OrchestratorPhaseCrawling',
            'analysis': 'constants.OrchestratorPhaseAnalysis',
            'scoring': 'constants.OrchestratorPhaseScoring',
            'reporting': 'constants.OrchestratorPhaseReporting',
            'cleanup': 'constants.OrchestratorPhaseCleanup',
            'validation': 'constants.OrchestratorPhaseValidation',
            'testing': 'constants.OrchestratorPhaseTesting',
            'security-scan': 'constants.OrchestratorPhaseSecurityScan'
        }
        for v in categorized['phase_names']:
            if v['value'] in phase_map:
                constants_map[f'"{v["value"]}"'] = phase_map[v['value']]
    
    # Status Values
    if 'status_values' in categorized:
        status_map = {
            'pending': 'constants.OrchestratorStatusPending',
            'running': 'constants.OrchestratorStatusRunning',
            'completed': 'constants.OrchestratorStatusCompleted',
            'failed': 'constants.OrchestratorStatusFailed',
            'cancelled': 'constants.OrchestratorStatusCancelled',
            'success': 'constants.OrchestratorStatusSuccess',
            'error': 'constants.OrchestratorStatusError',
            'warning': 'constants.OrchestratorStatusWarning',
            'info': 'constants.OrchestratorStatusInfo'
        }
        for v in categorized['status_values']:
            if v['value'] in status_map:
                constants_map[f'"{v["value"]}"'] = status_map[v['value']]
    
    # Config Keys
    if 'config_keys' in categorized:
        config_map = {
            'max_workers': 'constants.OrchestratorConfigMaxWorkers',
            'timeout': 'constants.OrchestratorConfigTimeout',
            'retry_attempts': 'constants.OrchestratorConfigRetryAttempts',
            'batch_size': 'constants.OrchestratorConfigBatchSize',
            'concurrent_agents': 'constants.OrchestratorConfigConcurrentAgents',
            'log_level': 'constants.OrchestratorConfigLogLevel',
            'debug_mode': 'constants.OrchestratorConfigDebugMode'
        }
        for v in categorized['config_keys']:
            if v['value'] in config_map:
                constants_map[f'"{v["value"]}"'] = config_map[v['value']]
    
    # HTTP Methods
    if 'http_methods' in categorized:
        http_map = {
            'GET': 'constants.OrchestratorHTTPMethodGet',
            'POST': 'constants.OrchestratorHTTPMethodPost',
            'PUT': 'constants.OrchestratorHTTPMethodPut',
            'DELETE': 'constants.OrchestratorHTTPMethodDelete',
            'PATCH': 'constants.OrchestratorHTTPMethodPatch',
            'HEAD': 'constants.OrchestratorHTTPMethodHead',
            'OPTIONS': 'constants.OrchestratorHTTPMethodOptions'
        }
        for v in categorized['http_methods']:
            if v['value'] in http_map:
                constants_map[f'"{v["value"]}"'] = http_map[v['value']]
    
    # Content Types
    if 'content_types' in categorized:
        content_map = {
            'application/json': 'constants.OrchestratorContentTypeJSON',
            'text/html': 'constants.OrchestratorContentTypeHTML',
            'text/plain': 'constants.OrchestratorContentTypePlain',
            'application/xml': 'constants.OrchestratorContentTypeXML'
        }
        for v in categorized['content_types']:
            if v['value'] in content_map:
                constants_map[f'"{v["value"]}"'] = content_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/orchestrator.go'
    
    print("🔍 CHARLIE-2 DETECTOR - Scanning orchestrator.go...")
    
    violations = analyze_orchestrator_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION CHARLIE-2:")
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
    
    with open('charlie2_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans charlie2_analysis.json")
    return results

if __name__ == "__main__":
    main()