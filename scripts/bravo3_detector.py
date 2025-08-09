#!/usr/bin/env python3
"""
🔍 BRAVO-3 HARDCODE DETECTOR
Détection industrielle des violations hardcoding dans security_agent.go
"""

import re
import json
from collections import defaultdict

def analyze_security_agent_file(filepath):
    """Analyse le fichier security_agent.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns de détection pour SECURITY AGENT
    patterns = {
        'agent_name': r'"(SECURITY-AGENT)"',
        'json_field_names': r'"(owasp_top10|dependency_check|secret_scanning|sql_injection|xss_check|csrf_check|report_path|timestamp|owasp_score|vulnerabilities|dependency_issues|secret_findings|security_headers|overall_risk|passed|id|title|severity|category|description|file|line|cwe|cvss|package|version|cve|fix_version|type|pattern|confidence|entropy|hsts|csp|x_frame_options|x_content_type_options|xss_protection|referrer_policy|score)"',
        'severity_levels': r'"(CRITICAL|HIGH|MEDIUM|LOW)"',
        'risk_levels': r'"(LOW|MEDIUM|HIGH|CRITICAL)"',
        'security_categories': r'"(OWASP|Static Analysis)"',
        'owasp_categories': r'"(A\d{2}:\d{4} - [^"]+)"',
        'vulnerability_types': r'"(API_KEY|PASSWORD|JWT_TOKEN|AWS_KEY|PRIVATE_KEY)"',
        'file_extensions': r'\.(go|yaml|yml|env)"',
        'security_tools': r'"(gosec|nancy|govulncheck|truffleHog|gitleaks)"',
        'security_commands': r'"(go install|github\.com/[^"]+)"',
        'test_types': r'"(injection|broken_authentication|sensitive_data_exposure|xml_external_entities|broken_access_control|security_misconfiguration|cross_site_scripting|insecure_deserialization|vulnerable_components|insufficient_logging)"',
        'cve_patterns': r'"(CVE-\d{4}-\d{4,})"',
        'file_patterns': r'"(gosec-report\.json|security_report\.json)"',
        'directory_patterns': r'"(vendor/|\.git/|node_modules/|tests/reports/security)"',
        'regex_patterns': r'`[^`]+`',
        'log_messages': r'"([🔒🛡️⚠️][^"]*|[A-Z][^"]*(?:security|vulnerability|analysis|scan)[^"]*)"',
        'error_messages': r'"([A-Z][^"]*(?:failed|error|not found)[^"]*)"',
        'success_messages': r'"([A-Z][^"]*(?:completed|found|installed)[^"]*)"',
        'package_examples': r'"(example/[^"]+|github\.com/example/[^"]+)"',
        'version_examples': r'"(v\d+\.\d+\.\d+)"',
        'numeric_scores': r'\b(100\.0|90|70|50|25\.0|20\.0|15\.0|10\.0|8\.0|5\.0|3\.1|3\.0|2\.0)\b',
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
    """Génère les mappings de constantes pour BRAVO-3"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Agent Name
    if 'agent_name' in categorized:
        constants_map['"SECURITY-AGENT"'] = 'constants.SecurityAgentName'
    
    # JSON Field Names
    if 'json_field_names' in categorized:
        json_field_map = {
            'owasp_top10': 'constants.SecurityJSONFieldOWASPTop10',
            'dependency_check': 'constants.SecurityJSONFieldDependencyCheck',
            'secret_scanning': 'constants.SecurityJSONFieldSecretScanning',
            'sql_injection': 'constants.SecurityJSONFieldSQLInjection',
            'xss_check': 'constants.SecurityJSONFieldXSSCheck',
            'csrf_check': 'constants.SecurityJSONFieldCSRFCheck',
            'report_path': 'constants.SecurityJSONFieldReportPath',
            'timestamp': 'constants.SecurityJSONFieldTimestamp',
            'owasp_score': 'constants.SecurityJSONFieldOWASPScore',
            'vulnerabilities': 'constants.SecurityJSONFieldVulnerabilities',
            'dependency_issues': 'constants.SecurityJSONFieldDependencyIssues',
            'secret_findings': 'constants.SecurityJSONFieldSecretFindings',
            'security_headers': 'constants.SecurityJSONFieldSecurityHeaders',
            'overall_risk': 'constants.SecurityJSONFieldOverallRisk',
            'passed': 'constants.SecurityJSONFieldPassed',
            'id': 'constants.SecurityJSONFieldID',
            'title': 'constants.SecurityJSONFieldTitle',
            'severity': 'constants.SecurityJSONFieldSeverity',
            'category': 'constants.SecurityJSONFieldCategory',
            'description': 'constants.SecurityJSONFieldDescription',
            'file': 'constants.SecurityJSONFieldFile',
            'line': 'constants.SecurityJSONFieldLine',
            'cwe': 'constants.SecurityJSONFieldCWE',
            'cvss': 'constants.SecurityJSONFieldCVSS',
            'package': 'constants.SecurityJSONFieldPackage',
            'version': 'constants.SecurityJSONFieldVersion',
            'cve': 'constants.SecurityJSONFieldCVE',
            'fix_version': 'constants.SecurityJSONFieldFixVersion',
            'type': 'constants.SecurityJSONFieldType',
            'pattern': 'constants.SecurityJSONFieldPattern',
            'confidence': 'constants.SecurityJSONFieldConfidence',
            'entropy': 'constants.SecurityJSONFieldEntropy',
            'hsts': 'constants.SecurityJSONFieldHSTS',
            'csp': 'constants.SecurityJSONFieldCSP',
            'x_frame_options': 'constants.SecurityJSONFieldXFrameOptions',
            'x_content_type_options': 'constants.SecurityJSONFieldXContentTypeOptions',
            'xss_protection': 'constants.SecurityJSONFieldXSSProtection',
            'referrer_policy': 'constants.SecurityJSONFieldReferrerPolicy',
            'score': 'constants.SecurityJSONFieldScore'
        }
        for v in categorized['json_field_names']:
            if v['value'] in json_field_map:
                constants_map[f'"{v["value"]}"'] = json_field_map[v['value']]
    
    # Severity Levels
    if 'severity_levels' in categorized:
        severity_map = {
            'CRITICAL': 'constants.SecuritySeverityCritical',
            'HIGH': 'constants.SecuritySeverityHigh',
            'MEDIUM': 'constants.SecuritySeverityMedium',
            'LOW': 'constants.SecuritySeverityLow'
        }
        for v in categorized['severity_levels']:
            if v['value'] in severity_map:
                constants_map[f'"{v["value"]}"'] = severity_map[v['value']]
    
    # Risk Levels (same as severity but for overall risk)
    if 'risk_levels' in categorized:
        risk_map = {
            'LOW': 'constants.SecurityRiskLow',
            'MEDIUM': 'constants.SecurityRiskMedium',
            'HIGH': 'constants.SecurityRiskHigh',
            'CRITICAL': 'constants.SecurityRiskCritical'
        }
        for v in categorized['risk_levels']:
            if v['value'] in risk_map:
                constants_map[f'"{v["value"]}"'] = risk_map[v['value']]
    
    # Security Categories
    if 'security_categories' in categorized:
        category_map = {
            'OWASP': 'constants.SecurityCategoryOWASP',
            'Static Analysis': 'constants.SecurityCategoryStaticAnalysis'
        }
        for v in categorized['security_categories']:
            if v['value'] in category_map:
                constants_map[f'"{v["value"]}"'] = category_map[v['value']]
    
    # Vulnerability Types
    if 'vulnerability_types' in categorized:
        vuln_type_map = {
            'API_KEY': 'constants.SecuritySecretTypeAPIKey',
            'PASSWORD': 'constants.SecuritySecretTypePassword',
            'JWT_TOKEN': 'constants.SecuritySecretTypeJWTToken',
            'AWS_KEY': 'constants.SecuritySecretTypeAWSKey',
            'PRIVATE_KEY': 'constants.SecuritySecretTypePrivateKey'
        }
        for v in categorized['vulnerability_types']:
            if v['value'] in vuln_type_map:
                constants_map[f'"{v["value"]}"'] = vuln_type_map[v['value']]
    
    # Security Tools
    if 'security_tools' in categorized:
        tool_map = {
            'gosec': 'constants.SecurityToolGosec',
            'nancy': 'constants.SecurityToolNancy',
            'govulncheck': 'constants.SecurityToolGovulncheck',
            'truffleHog': 'constants.SecurityToolTruffleHog',
            'gitleaks': 'constants.SecurityToolGitleaks'
        }
        for v in categorized['security_tools']:
            if v['value'] in tool_map:
                constants_map[f'"{v["value"]}"'] = tool_map[v['value']]
    
    # Test Types
    if 'test_types' in categorized:
        test_map = {
            'injection': 'constants.SecurityTestInjection',
            'broken_authentication': 'constants.SecurityTestBrokenAuth',
            'sensitive_data_exposure': 'constants.SecurityTestDataExposure',
            'xml_external_entities': 'constants.SecurityTestXXE',
            'broken_access_control': 'constants.SecurityTestAccessControl',
            'security_misconfiguration': 'constants.SecurityTestMisconfiguration',
            'cross_site_scripting': 'constants.SecurityTestXSS',
            'insecure_deserialization': 'constants.SecurityTestDeserialization',
            'vulnerable_components': 'constants.SecurityTestVulnComponents',
            'insufficient_logging': 'constants.SecurityTestLogging'
        }
        for v in categorized['test_types']:
            if v['value'] in test_map:
                constants_map[f'"{v["value"]}"'] = test_map[v['value']]
    
    # File Patterns
    if 'file_patterns' in categorized:
        file_map = {
            'gosec-report.json': 'constants.SecurityFileGosecReport',
            'security_report.json': 'constants.SecurityFileSecurityReport'
        }
        for v in categorized['file_patterns']:
            if v['value'] in file_map:
                constants_map[f'"{v["value"]}"'] = file_map[v['value']]
    
    # Directory Patterns
    if 'directory_patterns' in categorized:
        dir_map = {
            'vendor/': 'constants.SecurityDirVendor',
            '.git/': 'constants.SecurityDirGit',
            'node_modules/': 'constants.SecurityDirNodeModules',
            'tests/reports/security': 'constants.SecurityDefaultReportPath'
        }
        for v in categorized['directory_patterns']:
            if v['value'] in dir_map:
                constants_map[f'"{v["value"]}"'] = dir_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/tests/agents/security/security_agent.go'
    
    print("🔍 BRAVO-3 DETECTOR - Scanning security_agent.go...")
    
    violations = analyze_security_agent_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION BRAVO-3:")
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
    
    with open('bravo3_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans bravo3_analysis.json")
    return results

if __name__ == "__main__":
    main()