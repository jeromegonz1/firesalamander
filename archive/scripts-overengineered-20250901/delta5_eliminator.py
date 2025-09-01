#!/usr/bin/env python3
"""
⚔️ DELTA-5 ELIMINATOR
Éliminateur spécialisé pour security_agent.go - 36 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta5_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-5"""
    with open('delta5_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour security_agent.go
    replacements = {
        # OWASP Categories
        '"injection"': 'constants.SecurityOWASPInjection',
        '"broken_authentication"': 'constants.SecurityOWASPBrokenAuthentication',
        '"sensitive_data"': 'constants.SecurityOWASPSensitiveData',
        '"xml_external_entities"': 'constants.SecurityOWASPXMLExternalEntities',
        '"broken_access_control"': 'constants.SecurityOWASPBrokenAccessControl',
        '"security_misconfiguration"': 'constants.SecurityOWASPSecurityMisconfiguration',
        '"cross_site_scripting"': 'constants.SecurityOWASPCrossSiteScripting',
        '"insecure_deserialization"': 'constants.SecurityOWASPInsecureDeserialization',
        '"components_vulnerabilities"': 'constants.SecurityOWASPComponentsVulnerabilities',
        '"logging_monitoring"': 'constants.SecurityOWASPLoggingMonitoring',
        
        # Security Vulnerabilities
        '"broken_authentication"': 'constants.SecurityVulnerabilityBrokenAuthentication',
        '"sensitive_data_exposure"': 'constants.SecurityVulnerabilitySensitiveDataExposure',
        '"security_misconfiguration"': 'constants.SecurityVulnerabilitySecurityMisconfiguration',
        '"insecure_deserialization"': 'constants.SecurityVulnerabilityInsecureDeserialization',
        
        # Security Levels (seulement dans certains contextes)
        '"error"': 'constants.SecurityLevelError',
        '"warning"': 'constants.SecurityLevelWarning',
        '"info"': 'constants.SecurityLevelInfo',
        '"critical"': 'constants.SecurityLevelCritical',
        '"high"': 'constants.SecurityLevelHigh',
        '"medium"': 'constants.SecurityLevelMedium',
        '"low"': 'constants.SecurityLevelLow',
        '"fatal"': 'constants.SecurityLevelFatal',
        '"debug"': 'constants.SecurityLevelDebug',
        
        # Security Actions
        '"test"': 'constants.SecurityActionTest',
        '"scan"': 'constants.SecurityActionScan',
        '"validate"': 'constants.SecurityActionValidate',
        '"verify"': 'constants.SecurityActionVerify',
        '"audit"': 'constants.SecurityActionAudit',
        '"monitor"': 'constants.SecurityActionMonitor',
        '"log"': 'constants.SecurityActionLog',
        '"alert"': 'constants.SecurityActionAlert',
        '"block"': 'constants.SecurityActionBlock',
        '"allow"': 'constants.SecurityActionAllow',
        '"deny"': 'constants.SecurityActionDeny',
        
        # Security Tests
        '"vulnerability"': 'constants.SecurityTestVulnerability',
        '"penetration"': 'constants.SecurityTestPenetration',
        '"security"': 'constants.SecurityTestSecurity',
        '"compliance"': 'constants.SecurityTestCompliance',
        '"audit"': 'constants.SecurityTestAudit',
        '"assessment"': 'constants.SecurityTestAssessment',
        
        # Security Config Keys
        '"security"': 'constants.SecurityConfigSecurity',
        '"auth"': 'constants.SecurityConfigAuth',
        '"password"': 'constants.SecurityConfigPassword',
        '"PASSWORD"': 'constants.SecurityConfigPassword',
        '"secret"': 'constants.SecurityConfigSecret',
        '"token"': 'constants.SecurityConfigToken',
        '"key"': 'constants.SecurityConfigKey',
        '"ssl"': 'constants.SecurityConfigSSL',
        '"tls"': 'constants.SecurityConfigTLS',
        '"encryption"': 'constants.SecurityConfigEncryption',
        
        # CVE Patterns (garder les CVE spécifiques tels quels)
        # Les CVE sont des identifiants uniques et ne doivent pas être remplacés par des constantes génériques
        
        # Security Categories
        '"authentication"': 'constants.SecurityCategoryAuthentication',
        '"authorization"': 'constants.SecurityCategoryAuthorization',
        '"input_validation"': 'constants.SecurityCategoryInputValidation',
        '"output_encoding"': 'constants.SecurityCategoryOutputEncoding',
        '"session_management"': 'constants.SecurityCategorySessionManagement',
        '"cryptography"': 'constants.SecurityCategoryCryptography',
        '"configuration"': 'constants.SecurityCategoryConfiguration',
        '"logging"': 'constants.SecurityCategoryLogging',
        '"api"': 'constants.SecurityCategoryAPI',
        '"network"': 'constants.SecurityCategoryNetwork',
        '"application"': 'constants.SecurityCategoryApplication',
        '"system"': 'constants.SecurityCategorySystem',
    }
    
    # Contextes à éviter (où ne PAS remplacer)
    avoid_contexts = [
        r'`json:', # struct tags
        r'//', # commentaires
        r'import', # imports
        r'package', # package declaration
        r'const\s+\w+\s*=', # déclarations de constantes
        r'fmt\.', # format strings
        r'log\.', # log statements
        r't\.', # testing methods
        r'CVE-\d{4}-\d+', # CVE identifiers (ne pas remplacer)
    ]
    
    eliminated_count = 0
    changes_made = []
    
    # Appliquer les remplacements ligne par ligne avec contexte
    lines = content.split('\n')
    modified_lines = []
    
    for line_num, line in enumerate(lines, 1):
        original_line = line
        modified_line = line
        
        # Vérifier si on doit éviter cette ligne
        should_avoid = False
        for avoid_pattern in avoid_contexts:
            if re.search(avoid_pattern, line):
                should_avoid = True
                break
        
        # Cas spécial pour les CVE - ne pas remplacer les identifiants CVE
        if 'CVE-' in line:
            should_avoid = True
        
        if not should_avoid:
            # Appliquer les remplacements
            for old_value, new_value in replacements.items():
                if old_value in line:
                    # Vérifications supplémentaires pour éviter les faux positifs
                    
                    # Pour les security levels, s'assurer qu'on est dans un contexte approprié
                    security_levels = ['"error"', '"warning"', '"info"', '"critical"', '"high"', 
                                     '"medium"', '"low"', '"fatal"', '"debug"']
                    
                    if old_value in security_levels:
                        # Seulement remplacer si c'est dans un contexte de niveau/priorité (avec : ou =)
                        if (':' in line or '=' in line or 'Level' in line or 'Priority' in line or 'Severity' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf', 'Sprintf']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les OWASP categories et security actions, contexte de sécurité
                    elif old_value in ['"injection"', '"broken_authentication"', '"sensitive_data"', 
                                     '"xml_external_entities"', '"broken_access_control"',
                                     '"security_misconfiguration"', '"cross_site_scripting"',
                                     '"insecure_deserialization"', '"components_vulnerabilities"',
                                     '"logging_monitoring"', '"test"', '"scan"', '"validate"',
                                     '"verify"', '"audit"', '"monitor"', '"vulnerability"']:
                        # Remplacer dans des contextes de sécurité (listes, affectations, comparaisons)
                        if ('=' in line or ':' in line or '[' in line or 'Type' in line or 'Category' in line or 'Action' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les config keys, seulement dans des contextes de configuration
                    elif old_value in ['"security"', '"auth"', '"password"', '"PASSWORD"', '"secret"',
                                     '"token"', '"key"', '"ssl"', '"tls"', '"encryption"']:
                        # Remplacer dans des contextes de configuration (clés, maps, etc.)
                        if ('Config' in line or 'Key' in line or '[' in line or '=' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les autres, remplacement contextuel
                    else:
                        # Éviter les remplacements dans les strings de format et les logs
                        if not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf', 'Sprintf', '%s', '%d', '%v']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
        
        modified_lines.append(modified_line)
    
    # Reconstituer le contenu
    new_content = '\n'.join(modified_lines)
    
    # Écrire le fichier modifié seulement s'il y a des changements
    if new_content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        
        print(f"✅ {eliminated_count} violations éliminées dans {filepath}")
        
        # Afficher les changements
        print(f"\n📝 CHANGEMENTS APPLIQUÉS:")
        for change in changes_made[:15]:  # Montrer les 15 premiers
            print(f"  Line {change['line']}: {change['old']} → {change['new']}")
        
        if len(changes_made) > 15:
            print(f"  ... et {len(changes_made) - 15} autres changements")
        
        return eliminated_count, changes_made
    else:
        print(f"⚠️ Aucun changement appliqué à {filepath}")
        return 0, []

def add_constants_import(filepath):
    """Ajouter l'import des constantes si nécessaire"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Vérifier si l'import existe déjà
    if 'internal/constants' in content:
        print("📦 Import des constantes déjà présent")
        return
    
    # Ajouter l'import après les autres imports
    lines = content.split('\n')
    import_added = False
    
    for i, line in enumerate(lines):
        # Chercher la fin des imports
        if line.strip() == ')' and not import_added:
            # Chercher le bloc d'import précédent
            for j in range(i-1, -1, -1):
                if 'import' in lines[j] and '(' in lines[j]:
                    # Insérer avant la parenthèse fermante
                    lines.insert(i, '\t"fire-salamander/internal/constants"')
                    import_added = True
                    break
            break
    
    if import_added:
        new_content = '\n'.join(lines)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print("📦 Import des constantes ajouté")

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/tests/agents/security/security_agent.go'
    
    print("⚔️ DELTA-5 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: security_agent.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-5 SUCCÈS!")
            print(f"✅ {eliminated} violations éliminées")
            print(f"💾 Backup disponible: {backup_path}")
            
            # Sauvegarder le rapport
            report = {
                'target_file': filepath,
                'total_detected': total_violations,
                'eliminated_count': eliminated,
                'backup_path': backup_path,
                'changes': changes
            }
            
            with open('delta5_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta5_elimination_report.json")
            
        else:
            print(f"\n⚠️ Aucune violation éliminée")
            print("Les violations détectées peuvent être dans des contextes non remplaçables (CVE, logs, etc.)")
            
    except Exception as e:
        print(f"\n❌ ERREUR durant l'élimination: {e}")
        # Restaurer le backup en cas d'erreur
        shutil.copy2(backup_path, filepath)
        print(f"🔄 Fichier restauré depuis le backup")
        raise
    
    return eliminated

if __name__ == "__main__":
    main()