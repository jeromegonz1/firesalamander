#!/usr/bin/env python3
"""
🔍 BRAVO-1 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans qa_agent.go
"""

import re
import json
from collections import defaultdict

def analyze_qa_agent_file(filepath):
    """Analyse le fichier qa_agent.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection
    patterns = {
        'json_field_names': r'"(min_coverage|enable_vet|enable_lint|enable_security|enable_complexity|output_format|report_path|total_coverage|packages_coverage|threshold|passed|file|line|column|message|category|rule|severity|confidence|function|complexity|total_tests|passed_tests|failed_tests|skipped_tests|duration|overall_score|status)"',
        'status_values': r'"(pass|fail|warning|passed|failed|skipped|success|error|critical|info)"',
        'output_formats': r'"(json|text|html|xml|csv)"',
        'severity_levels': r'"(high|medium|low|critical|warning|info|notice)"',
        'confidence_levels': r'"(high|medium|low|certain|probable|possible)"',
        'tool_names': r'"(QA-AGENT|go|vet|lint|test|build|gofmt|golint|staticcheck)"',
        'command_strings': r'"(go\s+\w+[^"]*)"',
        'file_patterns': r'"\*\.(go|test\.go|mod|sum)"',
        'log_messages': r'"([A-Z][^"]*(?:starting|completed|failed|analyzing|found|executing)[^"]*)"',
        'error_messages': r'"([A-Z][^"]*(?:error|failed|unable|cannot|missing)[^"]*)"',
        'coverage_messages': r'"([^"]*(?:coverage|test|analysis)[^"]*)"',
        'report_templates': r'"([^"]*(?:Report|Analysis|Summary|Results)[^"]*)"',
        'hardcoded_paths': r'"/[^"]*"',
        'hardcoded_extensions': r'"\.(go|json|html|xml|txt|md)"',
        'numeric_literals': r'\b(80|85|90|95|100|0\.8|0\.85|0\.9|0\.95)\b',
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
    """Génère les mappings de constantes pour BRAVO-1"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_field_map = {
            'min_coverage': 'constants.JSONFieldMinCoverage',
            'enable_vet': 'constants.JSONFieldEnableVet',
            'enable_lint': 'constants.JSONFieldEnableLint',
            'enable_security': 'constants.JSONFieldEnableSecurity',
            'enable_complexity': 'constants.JSONFieldEnableComplexity',
            'output_format': 'constants.JSONFieldOutputFormat',
            'report_path': 'constants.JSONFieldReportPath',
            'total_coverage': 'constants.JSONFieldTotalCoverage',
            'packages_coverage': 'constants.JSONFieldPackagesCoverage',
            'threshold': 'constants.JSONFieldThreshold',
            'passed': 'constants.JSONFieldPassed',
            'file': 'constants.JSONFieldFile',
            'line': 'constants.JSONFieldLine',
            'column': 'constants.JSONFieldColumn',
            'message': 'constants.JSONFieldMessage',
            'category': 'constants.JSONFieldCategory',
            'rule': 'constants.JSONFieldRule',
            'severity': 'constants.JSONFieldSeverity',
            'confidence': 'constants.JSONFieldConfidence',
            'function': 'constants.JSONFieldFunction',
            'complexity': 'constants.JSONFieldComplexity',
            'total_tests': 'constants.JSONFieldTotalTests',
            'passed_tests': 'constants.JSONFieldPassedTests',
            'failed_tests': 'constants.JSONFieldFailedTests',
            'skipped_tests': 'constants.JSONFieldSkippedTests',
            'duration': 'constants.JSONFieldDuration',
            'overall_score': 'constants.JSONFieldOverallScore',
            'status': 'constants.JSONFieldStatus'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_field_map:
                constants_map[f'"{v["value"]}"'] = json_field_map[v['value']]
    
    # Status Values
    if 'status_values' in categorized:
        status_map = {
            'pass': 'constants.QAStatusPass',
            'fail': 'constants.QAStatusFail',
            'warning': 'constants.QAStatusWarning',
            'passed': 'constants.TestStatusPassed',
            'failed': 'constants.TestStatusFailed',
            'skipped': 'constants.TestStatusSkipped',
            'success': 'constants.QAStatusSuccess',
            'error': 'constants.QAStatusError',
            'critical': 'constants.SeverityCritical',
            'info': 'constants.SeverityInfo'
        }
        for v in categorized['status_values']:
            if v['value'] in status_map:
                constants_map[f'"{v["value"]}"'] = status_map[v['value']]
    
    # Output Formats
    if 'output_formats' in categorized:
        format_map = {
            'json': 'constants.OutputFormatJSON',
            'text': 'constants.OutputFormatText',
            'html': 'constants.OutputFormatHTML',
            'xml': 'constants.OutputFormatXML',
            'csv': 'constants.OutputFormatCSV'
        }
        for v in categorized['output_formats']:
            if v['value'] in format_map:
                constants_map[f'"{v["value"]}"'] = format_map[v['value']]
    
    # Severity Levels
    if 'severity_levels' in categorized:
        severity_map = {
            'high': 'constants.SeverityHigh',
            'medium': 'constants.SeverityMedium',
            'low': 'constants.SeverityLow',
            'critical': 'constants.SeverityCritical',
            'warning': 'constants.SeverityWarning',
            'info': 'constants.SeverityInfo',
            'notice': 'constants.SeverityNotice'
        }
        for v in categorized['severity_levels']:
            if v['value'] in severity_map:
                constants_map[f'"{v["value"]}"'] = severity_map[v['value']]
    
    # Confidence Levels
    if 'confidence_levels' in categorized:
        confidence_map = {
            'high': 'constants.ConfidenceHigh',
            'medium': 'constants.ConfidenceMedium',
            'low': 'constants.ConfidenceLow',
            'certain': 'constants.ConfidenceCertain',
            'probable': 'constants.ConfidenceProbable',
            'possible': 'constants.ConfidencePossible'
        }
        for v in categorized['confidence_levels']:
            if v['value'] in confidence_map:
                constants_map[f'"{v["value"]}"'] = confidence_map[v['value']]
    
    # Tool Names
    if 'tool_names' in categorized:
        tool_map = {
            'QA-AGENT': 'constants.QAAgentName',
            'go': 'constants.GoTool',
            'vet': 'constants.VetTool',
            'lint': 'constants.LintTool',
            'test': 'constants.TestTool',
            'build': 'constants.BuildTool',
            'gofmt': 'constants.GofmtTool',
            'golint': 'constants.GolintTool',
            'staticcheck': 'constants.StaticcheckTool'
        }
        for v in categorized['tool_names']:
            if v['value'] in tool_map:
                constants_map[f'"{v["value"]}"'] = tool_map[v['value']]
    
    # File Extensions
    if 'hardcoded_extensions' in categorized:
        ext_map = {
            '.go': 'constants.ExtGo',
            '.json': 'constants.ExtJSON',
            '.html': 'constants.ExtHTML',
            '.xml': 'constants.ExtXML',
            '.txt': 'constants.ExtTxt',
            '.md': 'constants.ExtMarkdown'
        }
        for v in categorized['hardcoded_extensions']:
            if v['value'] in ext_map:
                constants_map[f'"{v["value"]}"'] = ext_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/tests/agents/qa/qa_agent.go'
    
    print("🔍 BRAVO-1 DETECTOR - Scanning qa_agent.go...")
    
    violations = analyze_qa_agent_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION BRAVO-1:")
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
    
    with open('bravo1_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans bravo1_analysis.json")
    return results

if __name__ == "__main__":
    main()