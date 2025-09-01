#!/usr/bin/env python3
"""
🔍 DELTA-5 DETECTOR
Détection spécialisée pour security_agent.go - 97 violations
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
    
    # Patterns spécialisés pour SECURITY AGENT
    patterns = {
        'security_vulnerabilities': r'"(xss|sql_injection|csrf|clickjacking|directory_traversal|command_injection|ldap_injection|xpath_injection|xxe|ssrf|rfi|lfi|deserialization|buffer_overflow|race_condition|privilege_escalation|information_disclosure|denial_of_service|broken_authentication|session_fixation|insecure_direct_object_references|security_misconfiguration|sensitive_data_exposure|missing_function_level_access_control|using_components_with_known_vulnerabilities|unvalidated_redirects_forwards)"',
        'security_headers': r'"(Content-Security-Policy|X-Frame-Options|X-XSS-Protection|X-Content-Type-Options|Strict-Transport-Security|Referrer-Policy|Feature-Policy|Permissions-Policy|Cross-Origin-Embedder-Policy|Cross-Origin-Opener-Policy|Cross-Origin-Resource-Policy|Expect-CT|Public-Key-Pins|X-Permitted-Cross-Domain-Policies)"',
        'security_protocols': r'"(https|tls|ssl|wss|sftp|ftps|ssh|pgp|gpg|aes|des|rsa|ecdsa|hmac|sha1|sha256|sha512|md5|bcrypt|scrypt|pbkdf2|oauth|saml|jwt|openid)"',
        'security_levels': r'"(low|medium|high|critical|info|warning|error|fatal|trace|debug|none|basic|intermediate|advanced|expert)"',
        'security_actions': r'"(allow|deny|block|permit|reject|ignore|log|alert|warn|monitor|audit|scan|test|validate|verify|authenticate|authorize|encrypt|decrypt|hash|sign|verify_signature|sanitize|escape|filter|whitelist|blacklist|quarantine)"',
        'security_categories': r'"(authentication|authorization|input_validation|output_encoding|error_handling|logging|session_management|cryptography|communication|configuration|business_logic|file_upload|xml|json|api|web_service|database|infrastructure|network|application|system)"',
        'owasp_categories': r'"(A01|A02|A03|A04|A05|A06|A07|A08|A09|A10|injection|broken_authentication|sensitive_data|xml_external_entities|broken_access_control|security_misconfiguration|cross_site_scripting|insecure_deserialization|components_vulnerabilities|logging_monitoring)"',
        'security_tests': r'"(penetration|vulnerability|security|compliance|audit|assessment|scan|analysis|review|inspection|validation|verification)"',
        'cve_patterns': r'"(CVE-[0-9]{4}-[0-9]{4,})"',
        'security_tools': r'"(nmap|nessus|burp|owasp_zap|metasploit|sqlmap|nikto|dirb|gobuster|hydra|john|hashcat|wireshark|tcpdump|ncat|socat|openssl|gpg|ssh_keygen)"',
        'encryption_algorithms': r'"(AES-128|AES-192|AES-256|DES|3DES|RSA-1024|RSA-2048|RSA-4096|ECDSA-256|ECDSA-384|ECDSA-521|SHA-1|SHA-224|SHA-256|SHA-384|SHA-512|MD4|MD5|HMAC-SHA1|HMAC-SHA256|HMAC-SHA512|PBKDF2|BCrypt|SCrypt)"',
        'security_status_codes': r'"(200|201|301|302|400|401|403|404|405|406|409|410|413|414|415|422|429|500|501|502|503|504)"',
        'security_endpoints': r'"/(security|auth|login|logout|register|reset|verify|2fa|mfa|oauth|saml|jwt|api/security|api/auth|api/validate|api/scan|api/audit)[^"]*"',
        'security_config_keys': r'"(security|auth|encryption|ssl|tls|certificate|key|secret|token|password|hash|salt|pepper|nonce|iv|cipher|algorithm|protocol|timeout|max_attempts|lockout|session|cookie|cors|csp|hsts)"',
        'security_error_messages': r'"[A-Z][^"]*(?:unauthorized|forbidden|access denied|invalid credentials|authentication failed|authorization failed|security violation|suspicious activity|attack detected|vulnerability found|compliance violation|audit failed)[^"]*"',
        'security_log_messages': r'"[A-Z][^"]*(?:security scan|vulnerability assessment|penetration test|security audit|compliance check|authentication attempt|authorization request|suspicious request|attack attempt|security event|security alert)[^"]*"',
        'compliance_standards': r'"(PCI-DSS|HIPAA|GDPR|SOX|ISO27001|NIST|OWASP|CIS|SANS|COBIT|ITIL|SOC2|FedRAMP|FISMA)"',
        'security_file_extensions': r'"\\.(key|crt|pem|p12|pfx|jks|keystore|truststore|cer|der|csr|p7b|p7c)"',
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

def generate_security_constants_mapping(violations):
    """Génère les mappings de constantes pour security_agent.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Security Vulnerabilities
    if 'security_vulnerabilities' in categorized:
        vuln_map = {
            'xss': 'constants.SecurityVulnerabilityXSS',
            'sql_injection': 'constants.SecurityVulnerabilitySQLInjection',
            'csrf': 'constants.SecurityVulnerabilityCSRF',
            'clickjacking': 'constants.SecurityVulnerabilityClickjacking',
            'directory_traversal': 'constants.SecurityVulnerabilityDirectoryTraversal',
            'command_injection': 'constants.SecurityVulnerabilityCommandInjection',
            'ldap_injection': 'constants.SecurityVulnerabilityLDAPInjection',
            'xpath_injection': 'constants.SecurityVulnerabilityXPathInjection',
            'xxe': 'constants.SecurityVulnerabilityXXE',
            'ssrf': 'constants.SecurityVulnerabilitySSRF',
            'rfi': 'constants.SecurityVulnerabilityRFI',
            'lfi': 'constants.SecurityVulnerabilityLFI',
            'deserialization': 'constants.SecurityVulnerabilityDeserialization',
            'buffer_overflow': 'constants.SecurityVulnerabilityBufferOverflow',
            'race_condition': 'constants.SecurityVulnerabilityRaceCondition',
            'privilege_escalation': 'constants.SecurityVulnerabilityPrivilegeEscalation',
            'information_disclosure': 'constants.SecurityVulnerabilityInformationDisclosure',
            'denial_of_service': 'constants.SecurityVulnerabilityDenialOfService',
            'broken_authentication': 'constants.SecurityVulnerabilityBrokenAuthentication',
            'session_fixation': 'constants.SecurityVulnerabilitySessionFixation'
        }
        for v in categorized['security_vulnerabilities']:
            if v['value'] in vuln_map:
                constants_map[f'"{v["value"]}"'] = vuln_map[v['value']]
    
    # Security Headers
    if 'security_headers' in categorized:
        header_map = {
            'Content-Security-Policy': 'constants.SecurityHeaderCSP',
            'X-Frame-Options': 'constants.SecurityHeaderXFrameOptions',
            'X-XSS-Protection': 'constants.SecurityHeaderXXSSProtection',
            'X-Content-Type-Options': 'constants.SecurityHeaderXContentTypeOptions',
            'Strict-Transport-Security': 'constants.SecurityHeaderHSTS',
            'Referrer-Policy': 'constants.SecurityHeaderReferrerPolicy',
            'Feature-Policy': 'constants.SecurityHeaderFeaturePolicy',
            'Permissions-Policy': 'constants.SecurityHeaderPermissionsPolicy',
            'Cross-Origin-Embedder-Policy': 'constants.SecurityHeaderCOEP',
            'Cross-Origin-Opener-Policy': 'constants.SecurityHeaderCOOP',
            'Cross-Origin-Resource-Policy': 'constants.SecurityHeaderCORP',
            'Expect-CT': 'constants.SecurityHeaderExpectCT',
            'Public-Key-Pins': 'constants.SecurityHeaderPKP',
            'X-Permitted-Cross-Domain-Policies': 'constants.SecurityHeaderXPermittedCrossDomainPolicies'
        }
        for v in categorized['security_headers']:
            if v['value'] in header_map:
                constants_map[f'"{v["value"]}"'] = header_map[v['value']]
    
    # Security Protocols
    if 'security_protocols' in categorized:
        protocol_map = {
            'https': 'constants.SecurityProtocolHTTPS',
            'tls': 'constants.SecurityProtocolTLS',
            'ssl': 'constants.SecurityProtocolSSL',
            'wss': 'constants.SecurityProtocolWSS',
            'sftp': 'constants.SecurityProtocolSFTP',
            'ftps': 'constants.SecurityProtocolFTPS',
            'ssh': 'constants.SecurityProtocolSSH',
            'pgp': 'constants.SecurityProtocolPGP',
            'gpg': 'constants.SecurityProtocolGPG',
            'oauth': 'constants.SecurityProtocolOAuth',
            'saml': 'constants.SecurityProtocolSAML',
            'jwt': 'constants.SecurityProtocolJWT',
            'openid': 'constants.SecurityProtocolOpenID'
        }
        for v in categorized['security_protocols']:
            if v['value'] in protocol_map:
                constants_map[f'"{v["value"]}"'] = protocol_map[v['value']]
    
    # Security Levels
    if 'security_levels' in categorized:
        level_map = {
            'low': 'constants.SecurityLevelLow',
            'medium': 'constants.SecurityLevelMedium',
            'high': 'constants.SecurityLevelHigh',
            'critical': 'constants.SecurityLevelCritical',
            'info': 'constants.SecurityLevelInfo',
            'warning': 'constants.SecurityLevelWarning',
            'error': 'constants.SecurityLevelError',
            'fatal': 'constants.SecurityLevelFatal',
            'none': 'constants.SecurityLevelNone',
            'basic': 'constants.SecurityLevelBasic',
            'intermediate': 'constants.SecurityLevelIntermediate',
            'advanced': 'constants.SecurityLevelAdvanced',
            'expert': 'constants.SecurityLevelExpert'
        }
        for v in categorized['security_levels']:
            if v['value'] in level_map:
                constants_map[f'"{v["value"]}"'] = level_map[v['value']]
    
    # Security Actions
    if 'security_actions' in categorized:
        action_map = {
            'allow': 'constants.SecurityActionAllow',
            'deny': 'constants.SecurityActionDeny',
            'block': 'constants.SecurityActionBlock',
            'permit': 'constants.SecurityActionPermit',
            'reject': 'constants.SecurityActionReject',
            'ignore': 'constants.SecurityActionIgnore',
            'log': 'constants.SecurityActionLog',
            'alert': 'constants.SecurityActionAlert',
            'warn': 'constants.SecurityActionWarn',
            'monitor': 'constants.SecurityActionMonitor',
            'audit': 'constants.SecurityActionAudit',
            'scan': 'constants.SecurityActionScan',
            'test': 'constants.SecurityActionTest',
            'validate': 'constants.SecurityActionValidate',
            'verify': 'constants.SecurityActionVerify',
            'authenticate': 'constants.SecurityActionAuthenticate',
            'authorize': 'constants.SecurityActionAuthorize',
            'encrypt': 'constants.SecurityActionEncrypt',
            'decrypt': 'constants.SecurityActionDecrypt',
            'hash': 'constants.SecurityActionHash',
            'sign': 'constants.SecurityActionSign',
            'sanitize': 'constants.SecurityActionSanitize',
            'escape': 'constants.SecurityActionEscape',
            'filter': 'constants.SecurityActionFilter',
            'whitelist': 'constants.SecurityActionWhitelist',
            'blacklist': 'constants.SecurityActionBlacklist',
            'quarantine': 'constants.SecurityActionQuarantine'
        }
        for v in categorized['security_actions']:
            if v['value'] in action_map:
                constants_map[f'"{v["value"]}"'] = action_map[v['value']]
    
    # Security Categories
    if 'security_categories' in categorized:
        category_map = {
            'authentication': 'constants.SecurityCategoryAuthentication',
            'authorization': 'constants.SecurityCategoryAuthorization',
            'input_validation': 'constants.SecurityCategoryInputValidation',
            'output_encoding': 'constants.SecurityCategoryOutputEncoding',
            'error_handling': 'constants.SecurityCategoryErrorHandling',
            'logging': 'constants.SecurityCategoryLogging',
            'session_management': 'constants.SecurityCategorySessionManagement',
            'cryptography': 'constants.SecurityCategoryCryptography',
            'communication': 'constants.SecurityCategoryCommunication',
            'configuration': 'constants.SecurityCategoryConfiguration',
            'business_logic': 'constants.SecurityCategoryBusinessLogic',
            'file_upload': 'constants.SecurityCategoryFileUpload',
            'xml': 'constants.SecurityCategoryXML',
            'json': 'constants.SecurityCategoryJSON',
            'api': 'constants.SecurityCategoryAPI',
            'web_service': 'constants.SecurityCategoryWebService',
            'database': 'constants.SecurityCategoryDatabase',
            'infrastructure': 'constants.SecurityCategoryInfrastructure',
            'network': 'constants.SecurityCategoryNetwork',
            'application': 'constants.SecurityCategoryApplication',
            'system': 'constants.SecurityCategorySystem'
        }
        for v in categorized['security_categories']:
            if v['value'] in category_map:
                constants_map[f'"{v["value"]}"'] = category_map[v['value']]
    
    # OWASP Categories
    if 'owasp_categories' in categorized:
        owasp_map = {
            'A01': 'constants.SecurityOWASPA01',
            'A02': 'constants.SecurityOWASPA02',
            'A03': 'constants.SecurityOWASPA03',
            'A04': 'constants.SecurityOWASPA04',
            'A05': 'constants.SecurityOWASPA05',
            'A06': 'constants.SecurityOWASPA06',
            'A07': 'constants.SecurityOWASPA07',
            'A08': 'constants.SecurityOWASPA08',
            'A09': 'constants.SecurityOWASPA09',
            'A10': 'constants.SecurityOWASPA10',
            'injection': 'constants.SecurityOWASPInjection',
            'broken_authentication': 'constants.SecurityOWASPBrokenAuthentication',
            'sensitive_data': 'constants.SecurityOWASPSensitiveData',
            'xml_external_entities': 'constants.SecurityOWASPXMLExternalEntities',
            'broken_access_control': 'constants.SecurityOWASPBrokenAccessControl',
            'security_misconfiguration': 'constants.SecurityOWASPSecurityMisconfiguration',
            'cross_site_scripting': 'constants.SecurityOWASPCrossSiteScripting',
            'insecure_deserialization': 'constants.SecurityOWASPInsecureDeserialization',
            'components_vulnerabilities': 'constants.SecurityOWASPComponentsVulnerabilities',
            'logging_monitoring': 'constants.SecurityOWASPLoggingMonitoring'
        }
        for v in categorized['owasp_categories']:
            if v['value'] in owasp_map:
                constants_map[f'"{v["value"]}"'] = owasp_map[v['value']]
    
    # Compliance Standards
    if 'compliance_standards' in categorized:
        compliance_map = {
            'PCI-DSS': 'constants.SecurityCompliancePCIDSS',
            'HIPAA': 'constants.SecurityComplianceHIPAA',
            'GDPR': 'constants.SecurityComplianceGDPR',
            'SOX': 'constants.SecurityComplianceSOX',
            'ISO27001': 'constants.SecurityComplianceISO27001',
            'NIST': 'constants.SecurityComplianceNIST',
            'OWASP': 'constants.SecurityComplianceOWASP',
            'CIS': 'constants.SecurityComplianceCIS',
            'SANS': 'constants.SecurityComplianceSANS',
            'COBIT': 'constants.SecurityComplianceCOBIT',
            'ITIL': 'constants.SecurityComplianceITIL',
            'SOC2': 'constants.SecurityComplianceSOC2',
            'FedRAMP': 'constants.SecurityComplianceFedRAMP',
            'FISMA': 'constants.SecurityComplianceFISMA'
        }
        for v in categorized['compliance_standards']:
            if v['value'] in compliance_map:
                constants_map[f'"{v["value"]}"'] = compliance_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/tests/agents/security/security_agent.go'
    
    print("🔍 DELTA-5 DETECTOR - Scanning security_agent.go...")
    
    violations = analyze_security_agent_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_security_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-5:")
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
    
    with open('delta5_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta5_analysis.json")
    return results

if __name__ == "__main__":
    main()