#!/usr/bin/env python3
"""
🔍 ALPHA-4 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans checker.go
"""

import re
import json
from collections import defaultdict

def analyze_checker_file(filepath):
    """Analyse le fichier checker.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection
    patterns = {
        'status_values': r'"(healthy|degraded|ok|error|warning|passed|failed)"',
        'messages': r'"([A-Z][^"]*(?:is|are|configuration|validation|missing|invalid|failed|valid)[^"]*)"',
        'db_types': r'"(sqlite|mysql|firesalamander|localhost)"',
        'paths': r'"(config|deploy|go\.mod|main\.go|\.env)"',
        'content_types': r'"(application/json|text/html|Content-Type)"',
        'environments': r'"(development|production|test)"',
        'error_codes': r'"([a-z_]+_(?:nil|missing|incomplete|invalid|error|failed))"',
        'db_paths': r'"/firesalamander\.db"',
        'hardcoded_booleans': r'\b(true|false)\b(?=\s*,|\s*}|\s*\))',
        'numeric_literals': r'\b(1024|4:)\b',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires et imports
        if line_stripped.startswith('//') or line_stripped.startswith('import') or line_stripped.startswith('package'):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line)
            for match in matches:
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match,
                    'context': line_stripped[:100] + ('...' if len(line_stripped) > 100 else '')
                })
    
    return violations

def categorize_violations(violations):
    """Catégorise les violations par type"""
    
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    
    return dict(categories)

def generate_constants(violations):
    """Génère les constantes pour les violations détectées"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Status Constants
    if 'status_values' in categorized:
        for v in categorized['status_values']:
            value = v['value']
            if value == 'healthy':
                constants_map['"healthy"'] = 'constants.HealthStatusHealthy'
            elif value == 'degraded':
                constants_map['"degraded"'] = 'constants.HealthStatusDegraded'
            elif value == 'ok':
                constants_map['"ok"'] = 'constants.CheckStatusOK'
            elif value == 'error':
                constants_map['"error"'] = 'constants.CheckStatusError'
            elif value == 'warning':
                constants_map['"warning"'] = 'constants.CheckStatusWarning'
            elif value == 'passed':
                constants_map['"passed"'] = 'constants.TestStatusPassed'
            elif value == 'failed':
                constants_map['"failed"'] = 'constants.TestStatusFailed'
    
    # Messages Constants
    if 'messages' in categorized:
        for v in categorized['messages']:
            value = v['value']
            if 'Configuration is nil' in value:
                constants_map[f'"{value}"'] = 'constants.MsgConfigIsNil'
            elif 'Configuration validation failed' in value:
                constants_map[f'"{value}"'] = 'constants.MsgConfigValidationFailed'
            elif 'Configuration is valid' in value:
                constants_map[f'"{value}"'] = 'constants.MsgConfigIsValid'
            elif 'SQLite path is empty' in value:
                constants_map[f'"{value}"'] = 'constants.MsgSQLitePathEmpty'
            elif 'SQLite configuration is valid' in value:
                constants_map[f'"{value}"'] = 'constants.MsgSQLiteConfigValid'
            elif 'MySQL configuration incomplete' in value:
                constants_map[f'"{value}"'] = 'constants.MsgMySQLConfigIncomplete'
            elif 'MySQL configuration is valid' in value:
                constants_map[f'"{value}"'] = 'constants.MsgMySQLConfigValid'
            elif 'Unknown database type' in value:
                constants_map[f'"{value}"'] = 'constants.MsgUnknownDatabaseType'
            elif 'Required files/directories missing' in value:
                constants_map[f'"{value}"'] = 'constants.MsgRequiredFilesMissing'
            elif 'All required files and directories present' in value:
                constants_map[f'"{value}"'] = 'constants.MsgAllFilesPresent'
            elif 'Network configuration valid' in value:
                constants_map[f'"{value}"'] = 'constants.MsgNetworkConfigValid'
            elif 'AI is disabled' in value:
                constants_map[f'"{value}"'] = 'constants.MsgAIDisabled'
            elif 'AI is in mock mode' in value:
                constants_map[f'"{value}"'] = 'constants.MsgAIMockMode'
            elif 'AI configuration is valid' in value:
                constants_map[f'"{value}"'] = 'constants.MsgAIConfigValid'
            elif 'Internal server error' in value:
                constants_map[f'"{value}"'] = 'constants.MsgInternalServerError'
    
    # Database Types
    if 'db_types' in categorized:
        for v in categorized['db_types']:
            value = v['value']
            if value == 'sqlite':
                constants_map['"sqlite"'] = 'constants.DatabaseTypeSQLite'
            elif value == 'mysql':
                constants_map['"mysql"'] = 'constants.DatabaseTypeMySQL'
            elif value == 'firesalamander':
                constants_map['"firesalamander"'] = 'constants.DefaultDatabaseName'
            elif value == 'localhost':
                constants_map['"localhost"'] = 'constants.DefaultHost'
    
    # Paths and Files
    if 'paths' in categorized:
        for v in categorized['paths']:
            value = v['value']
            if value == 'config':
                constants_map['"config"'] = 'constants.ConfigDir'
            elif value == 'deploy':
                constants_map['"deploy"'] = 'constants.DeployDir'
            elif value == 'go.mod':
                constants_map['"go.mod"'] = 'constants.GoModFile'
            elif value == 'main.go':
                constants_map['"main.go"'] = 'constants.MainGoFile'
    
    # Content Types
    if 'content_types' in categorized:
        for v in categorized['content_types']:
            value = v['value']
            if value == 'application/json':
                constants_map['"application/json"'] = 'constants.ContentTypeJSON'
            elif value == 'Content-Type':
                constants_map['"Content-Type"'] = 'constants.HeaderContentType'
    
    # Environments
    if 'environments' in categorized:
        for v in categorized['environments']:
            value = v['value']
            if value == 'development':
                constants_map['"development"'] = 'constants.EnvDevelopment'
    
    # Error Codes
    if 'error_codes' in categorized:
        for v in categorized['error_codes']:
            value = v['value']
            if value == 'config_nil':
                constants_map['"config_nil"'] = 'constants.ErrorCodeConfigNil'
            elif value == 'sqlite_path_missing':
                constants_map['"sqlite_path_missing"'] = 'constants.ErrorCodeSQLitePathMissing'
            elif value == 'mysql_config_incomplete':
                constants_map['"mysql_config_incomplete"'] = 'constants.ErrorCodeMySQLConfigIncomplete'
            elif value == 'unknown_db_type':
                constants_map['"unknown_db_type"'] = 'constants.ErrorCodeUnknownDBType'
    
    # DB Paths
    if 'db_paths' in categorized:
        constants_map['"/firesalamander.db"'] = 'constants.DefaultDatabaseFile'
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/debug/checker.go'
    
    print("🔍 ALPHA-4 DETECTOR - Scanning checker.go...")
    
    violations = analyze_checker_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION ALPHA-4:")
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
    
    with open('alpha4_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans alpha4_analysis.json")
    return results

if __name__ == "__main__":
    main()