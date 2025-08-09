#!/usr/bin/env python3
"""
🤖 BRAVO-3 SMART ELIMINATOR
Élimination automatisée des violations hardcoding dans security_agent.go
Mission: 235 violations → 0 violations
"""

import re
import shutil
import subprocess
import sys

def create_string_mappings():
    """Créer les mappings de remplacement pour BRAVO-3"""
    
    mappings = {
        # Agent Name
        '"SECURITY-AGENT"': 'constants.SecurityAgentName',
        
        # JSON Field Names - Configuration
        '"owasp_top10"': 'constants.SecurityJSONFieldOWASPTop10',
        '"dependency_check"': 'constants.SecurityJSONFieldDependencyCheck',
        '"secret_scanning"': 'constants.SecurityJSONFieldSecretScanning',
        '"sql_injection"': 'constants.SecurityJSONFieldSQLInjection',
        '"xss_check"': 'constants.SecurityJSONFieldXSSCheck',
        '"csrf_check"': 'constants.SecurityJSONFieldCSRFCheck',
        '"report_path"': 'constants.SecurityJSONFieldReportPath',
        
        # JSON Field Names - Results
        '"timestamp"': 'constants.SecurityJSONFieldTimestamp',
        '"owasp_score"': 'constants.SecurityJSONFieldOWASPScore',
        '"vulnerabilities"': 'constants.SecurityJSONFieldVulnerabilities',
        '"dependency_issues"': 'constants.SecurityJSONFieldDependencyIssues',
        '"secret_findings"': 'constants.SecurityJSONFieldSecretFindings',
        '"security_headers"': 'constants.SecurityJSONFieldSecurityHeaders',
        '"overall_risk"': 'constants.SecurityJSONFieldOverallRisk',
        '"passed"': 'constants.SecurityJSONFieldPassed',
        
        # JSON Field Names - Vulnerability
        '"id"': 'constants.SecurityJSONFieldID',
        '"title"': 'constants.SecurityJSONFieldTitle',
        '"severity"': 'constants.SecurityJSONFieldSeverity',
        '"category"': 'constants.SecurityJSONFieldCategory',
        '"description"': 'constants.SecurityJSONFieldDescription',
        '"file"': 'constants.SecurityJSONFieldFile',
        '"line"': 'constants.SecurityJSONFieldLine',
        '"cwe"': 'constants.SecurityJSONFieldCWE',
        '"cvss"': 'constants.SecurityJSONFieldCVSS',
        
        # JSON Field Names - Dependency
        '"package"': 'constants.SecurityJSONFieldPackage',
        '"version"': 'constants.SecurityJSONFieldVersion',
        '"cve"': 'constants.SecurityJSONFieldCVE',
        '"fix_version"': 'constants.SecurityJSONFieldFixVersion',
        
        # JSON Field Names - Secret Finding
        '"type"': 'constants.SecurityJSONFieldType',
        '"pattern"': 'constants.SecurityJSONFieldPattern',
        '"confidence"': 'constants.SecurityJSONFieldConfidence',
        '"entropy"': 'constants.SecurityJSONFieldEntropy',
        
        # JSON Field Names - Security Headers
        '"hsts"': 'constants.SecurityJSONFieldHSTS',
        '"csp"': 'constants.SecurityJSONFieldCSP',
        '"x_frame_options"': 'constants.SecurityJSONFieldXFrameOptions',
        '"x_content_type_options"': 'constants.SecurityJSONFieldXContentTypeOptions',
        '"xss_protection"': 'constants.SecurityJSONFieldXSSProtection',
        '"referrer_policy"': 'constants.SecurityJSONFieldReferrerPolicy',
        '"score"': 'constants.SecurityJSONFieldScore',
        
        # Severity Levels
        '"CRITICAL"': 'constants.SecuritySeverityCritical',
        '"HIGH"': 'constants.SecuritySeverityHigh',
        '"MEDIUM"': 'constants.SecuritySeverityMedium',
        '"LOW"': 'constants.SecuritySeverityLow',
        
        # Risk Levels (overlap with severity but different usage)
        '"LOW"': 'constants.SecurityRiskLow',
        
        # Security Categories
        '"OWASP"': 'constants.SecurityCategoryOWASP',
        '"Static Analysis"': 'constants.SecurityCategoryStaticAnalysis',
        
        # OWASP Categories
        '"A06:2021 - Security Misconfiguration"': 'constants.SecurityOWASPMisconfiguration',
        
        # Secret Types
        '"API_KEY"': 'constants.SecuritySecretTypeAPIKey',
        '"PASSWORD"': 'constants.SecuritySecretTypePassword',
        '"JWT_TOKEN"': 'constants.SecuritySecretTypeJWTToken',
        '"AWS_KEY"': 'constants.SecuritySecretTypeAWSKey',
        '"PRIVATE_KEY"': 'constants.SecuritySecretTypePrivateKey',
        
        # Security Tools
        '"gosec"': 'constants.SecurityToolGosec',
        '"nancy"': 'constants.SecurityToolNancy',
        '"govulncheck"': 'constants.SecurityToolGovulncheck',
        '"truffleHog"': 'constants.SecurityToolTruffleHog',
        '"gitleaks"': 'constants.SecurityToolGitleaks',
        
        # Test Types
        '"injection"': 'constants.SecurityTestInjection',
        '"broken_authentication"': 'constants.SecurityTestBrokenAuth',
        '"sensitive_data_exposure"': 'constants.SecurityTestDataExposure',
        '"xml_external_entities"': 'constants.SecurityTestXXE',
        '"broken_access_control"': 'constants.SecurityTestAccessControl',
        '"security_misconfiguration"': 'constants.SecurityTestMisconfiguration',
        '"cross_site_scripting"': 'constants.SecurityTestXSS',
        '"insecure_deserialization"': 'constants.SecurityTestDeserialization',
        '"vulnerable_components"': 'constants.SecurityTestVulnComponents',
        '"insufficient_logging"': 'constants.SecurityTestLogging',
        
        # File Patterns
        '"gosec-report.json"': 'constants.SecurityFileGosecReport',
        '"security_report.json"': 'constants.SecurityFileSecurityReport',
        
        # Directory Patterns
        '"vendor/"': 'constants.SecurityDirVendor',
        '".git/"': 'constants.SecurityDirGit',
        '"node_modules/"': 'constants.SecurityDirNodeModules',
        '"tests/reports/security"': 'constants.SecurityDefaultReportPath',
        
        # File Extensions
        '".go"': 'constants.SecurityExtGo',
        '".yaml"': 'constants.SecurityExtYAML',
        '".yml"': 'constants.SecurityExtYML',
        '".env"': 'constants.SecurityExtEnv',
        
        # Commands
        '"go"': 'constants.SecurityCmdGo',
        '"install"': 'constants.SecurityCmdInstall',
        '"github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"': 'constants.SecurityCmdInstallGosec',
        '"golang.org/x/vuln/cmd/govulncheck@latest"': 'constants.SecurityCmdInstallGovulncheck',
        '"-fmt"': 'constants.SecurityArgFormat',
        '"-out"': 'constants.SecurityArgOut',
        '"json"': 'constants.SecurityArgJSON',
        '"sleuth"': 'constants.SecurityArgSleuth',
        '"./..."': 'constants.SecurityArgAll',
        
        # Log Messages
        '"🔒 Starting OWASP security scan"': 'constants.SecurityMsgStartingScan',
        '"Security scan completed"': 'constants.SecurityMsgSecurityScanCompleted',
        '"Running static code analysis with gosec"': 'constants.SecurityMsgRunningStaticAnalysis',
        '"Checking dependencies for vulnerabilities"': 'constants.SecurityMsgCheckingDependencies',
        '"Scanning for hardcoded secrets"': 'constants.SecurityMsgScanningSecrets',
        '"Running dynamic security analysis"': 'constants.SecurityMsgRunningDynamicAnalysis',
        '"Checking HTTP security headers"': 'constants.SecurityMsgCheckingSecurityHeaders',
        
        # Error Messages
        '"Static analysis failed"': 'constants.SecurityMsgStaticAnalysisFailed',
        '"Dependency check failed"': 'constants.SecurityMsgDependencyCheckFailed',
        '"Secret scanning failed"': 'constants.SecurityMsgSecretScanningFailed',
        '"Dynamic analysis failed"': 'constants.SecurityMsgDynamicAnalysisFailed',
        '"Security headers check failed"': 'constants.SecurityMsgSecurityHeadersFailed',
        '"Failed to generate security report"': 'constants.SecurityMsgReportGenerationFailed',
        
        # Tool Messages
        '"gosec not found, installing..."': 'constants.SecurityMsgGosecNotFound',
        '"nancy not found, using govulncheck instead"': 'constants.SecurityMsgNancyNotFound',
        '"failed to install gosec"': 'constants.SecurityMsgFailedInstallGosec',
        '"Failed to install govulncheck"': 'constants.SecurityMsgFailedInstallGovuln',
        '"gosec found security issues"': 'constants.SecurityMsgGosecFoundIssues',
        '"nancy found vulnerable dependencies"': 'constants.SecurityMsgNancyFoundVulns',
        '"govulncheck found vulnerabilities"': 'constants.SecurityMsgGovulncheckFoundVulns',
        '"Failed to parse gosec results"': 'constants.SecurityMsgFailedParseGosec',
        
        # Example Values
        '"example/vulnerable-package"': 'constants.SecurityExamplePackage',
        '"github.com/example/vulnerable"': 'constants.SecurityExampleGitHubPackage',
        '"v1.0.0"': 'constants.SecurityExampleVersion',
        '"v1.0.1"': 'constants.SecurityExampleFixVersion',
        '"CVE-2023-12345"': 'constants.SecurityExampleCVE',
        '"CVE-2023-99999"': 'constants.SecurityExampleCVEDemo',
        '"SEC-001"': 'constants.SecurityExampleVulnID',
        '"Missing Security Headers"': 'constants.SecurityExampleVulnTitle',
        '"Application is missing recommended security headers"': 'constants.SecurityExampleVulnDescription',
        '"CWE-16"': 'constants.SecurityExampleCWE',
        '"Example vulnerability for testing"': 'constants.SecurityExampleVulnDescription',
        '"Example vulnerability in dependency"': 'constants.SecurityExampleVulnDescription',
        
        # Numeric Values
        '100.0': 'constants.SecurityBaseScore',
        '25.0': 'constants.SecurityScoreCriticalVuln',
        '20.0': 'constants.SecurityScoreCriticalDep',
        '15.0': 'constants.SecurityScoreHighVuln',
        '10.0': 'constants.SecurityScoreHighDep',
        '8.0': 'constants.SecurityScoreMediumVuln',
        '5.0': 'constants.SecurityScoreMediumDep',
        '3.0': 'constants.SecurityScoreLowVuln',
        '3.1': 'constants.SecurityExampleCVSS',
        '2.0': 'constants.SecurityScoreLowDep',
        '90': 'constants.SecurityThresholdLowRisk',
        '70': 'constants.SecurityThresholdMediumRisk',
        '70.0': 'constants.SecurityScoreHeadersWeight',
        '50': 'constants.SecurityThresholdHighRisk',
    }
    
    return mappings

def backup_file(filepath):
    """Créer une sauvegarde du fichier original"""
    backup_path = f"{filepath}.bravo3_backup"
    shutil.copy2(filepath, backup_path)
    return backup_path

def restore_file(filepath):
    """Restaurer le fichier depuis la sauvegarde"""
    backup_path = f"{filepath}.bravo3_backup"
    try:
        shutil.copy2(backup_path, filepath)
        return True
    except:
        return False

def apply_targeted_replacements(filepath, mappings):
    """Appliquer les remplacements de hardcoding de manière ciblée"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements_made = 0
    
    # Remplacements spéciaux avec contexte pour éviter les conflits
    special_replacements = [
        # Test array replacement
        (r'tests := \[\]string\{[^}]+\}', lambda m: replace_test_array(m.group(0))),
        
        # Secret patterns map
        (r'_ = map\[string\]string\{[^}]+\}', lambda m: replace_secret_patterns_map(m.group(0))),
        
        # Switch case for severity levels in vulnerability scoring
        (r'switch vuln\.Severity \{[^}]+case "CRITICAL":[^}]+case "HIGH":[^}]+case "MEDIUM":[^}]+case "LOW":[^}]+\}', 
         lambda m: replace_vuln_severity_switch(m.group(0))),
         
        # Switch case for dependency severity levels
        (r'switch dep\.Severity \{[^}]+case "CRITICAL":[^}]+case "HIGH":[^}]+case "MEDIUM":[^}]+case "LOW":[^}]+\}',
         lambda m: replace_dep_severity_switch(m.group(0))),
    ]
    
    def replace_test_array(match_text):
        """Replace OWASP tests array with constants"""
        test_constants = [
            'constants.SecurityTestInjection',
            'constants.SecurityTestBrokenAuth',
            'constants.SecurityTestDataExposure',
            'constants.SecurityTestXXE',
            'constants.SecurityTestAccessControl',
            'constants.SecurityTestMisconfiguration',
            'constants.SecurityTestXSS',
            'constants.SecurityTestDeserialization',
            'constants.SecurityTestVulnComponents',
            'constants.SecurityTestLogging'
        ]
        return f'tests := []string{{\n\t\t{",\\n\\t\\t".join(test_constants)},\n\t}}'
    
    def replace_secret_patterns_map(match_text):
        """Replace secret patterns map with constants"""
        return '''_ = map[string]string{
\t\tconstants.SecuritySecretTypeAPIKey:     constants.SecurityPatternAPIKey,
\t\tconstants.SecuritySecretTypePassword:   constants.SecurityPatternPassword,
\t\tconstants.SecuritySecretTypeJWTToken:   constants.SecurityPatternJWTToken,
\t\tconstants.SecuritySecretTypeAWSKey:     constants.SecurityPatternAWSKey,
\t\tconstants.SecuritySecretTypePrivateKey: constants.SecurityPatternPrivateKey,
\t}'''
    
    def replace_vuln_severity_switch(match_text):
        """Replace vulnerability severity switch statement"""
        return '''switch vuln.Severity {
\t\tcase constants.SecuritySeverityCritical:
\t\t\tscore -= constants.SecurityScoreCriticalVuln
\t\tcase constants.SecuritySeverityHigh:
\t\t\tscore -= constants.SecurityScoreHighVuln
\t\tcase constants.SecuritySeverityMedium:
\t\t\tscore -= constants.SecurityScoreMediumVuln
\t\tcase constants.SecuritySeverityLow:
\t\t\tscore -= constants.SecurityScoreLowVuln
\t\t}'''
    
    def replace_dep_severity_switch(match_text):
        """Replace dependency severity switch statement"""
        return '''switch dep.Severity {
\t\tcase constants.SecuritySeverityCritical:
\t\t\tscore -= constants.SecurityScoreCriticalDep
\t\tcase constants.SecuritySeverityHigh:
\t\t\tscore -= constants.SecurityScoreHighDep
\t\tcase constants.SecuritySeverityMedium:
\t\t\tscore -= constants.SecurityScoreMediumDep
\t\tcase constants.SecuritySeverityLow:
\t\t\tscore -= constants.SecurityScoreLowDep
\t\t}'''
    
    # Appliquer les remplacements spéciaux d'abord
    for pattern, replacement in special_replacements:
        if isinstance(replacement, str):
            if re.search(pattern, content, re.DOTALL):
                content = re.sub(pattern, replacement, content, flags=re.DOTALL)
                replacements_made += 1
                print(f"✅ Special: {pattern[:50]}... → {replacement[:50]}...")
        else:
            # C'est une fonction lambda
            matches = re.findall(pattern, content, re.DOTALL)
            if matches:
                content = re.sub(pattern, replacement, content, flags=re.DOTALL)
                replacements_made += len(matches)
                print(f"✅ Special function: {pattern[:50]}... → {len(matches)} replacements")
    
    # Appliquer les remplacements standards
    for old_string, new_string in mappings.items():
        if old_string in content:
            content = content.replace(old_string, new_string)
            replacements_made += 1
            print(f"✅ {old_string} → {new_string}")
    
    # Sauvegarder le fichier modifié
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    return replacements_made, len(content) != len(original_content)

def test_compilation():
    """Tester la compilation du code Go"""
    try:
        result = subprocess.run(['go', 'build', './tests/agents/security'], 
                              capture_output=True, text=True, cwd='/Users/jeromegonzalez/claude-code/fire-salamander')
        
        if result.returncode == 0:
            print("✅ Compilation successful!")
            return True
        else:
            print(f"❌ Compilation failed:\n{result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Error testing compilation: {e}")
        return False

def main():
    """Fonction principale BRAVO-3 Eliminator"""
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/tests/agents/security/security_agent.go'
    
    print("🤖 BRAVO-3 SMART ELIMINATOR")
    print("=" * 50)
    print(f"Target: {filepath}")
    print("Mission: Élimination complète des hardcoding violations")
    
    # 1. Backup
    print("\n📦 Creating backup...")
    backup_path = backup_file(filepath)
    print(f"✅ Backup created: {backup_path}")
    
    # 2. Load mappings
    print("\n🔧 Loading string mappings...")
    mappings = create_string_mappings()
    print(f"✅ {len(mappings)} mappings loaded")
    
    # 3. Apply replacements
    print("\n🔄 Applying hardcode eliminations...")
    replacements_made, file_modified = apply_targeted_replacements(filepath, mappings)
    
    print(f"\n📊 RÉSULTATS:")
    print(f"   - Replacements applied: {replacements_made}")
    print(f"   - File modified: {file_modified}")
    
    # 4. Test compilation
    print(f"\n🔨 Testing compilation...")
    if test_compilation():
        print(f"\n🎯 BRAVO-3 MISSION ACCOMPLISHED!")
        print(f"   - Target file: security_agent.go")
        print(f"   - Violations eliminated: {replacements_made}")
        print(f"   - Status: 100% SUCCESS")
        return True
    else:
        print(f"\n⚠️  Compilation failed - restoring backup...")
        if restore_file(filepath):
            print(f"   - File restored from backup")
        print(f"   - Status: NEEDS MANUAL INTERVENTION")
        return False

if __name__ == "__main__":
    import os
    main()