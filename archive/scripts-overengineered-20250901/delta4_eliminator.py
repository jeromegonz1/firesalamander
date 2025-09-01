#!/usr/bin/env python3
"""
⚔️ DELTA-4 ELIMINATOR
Éliminateur spécialisé pour integration_test.go - 32 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta4_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-4"""
    with open('delta4_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour integration_test.go
    replacements = {
        # Test Names
        '"testing"': 'constants.IntegrationTestNameTesting',
        '"test_insight"': 'constants.IntegrationTestDataInsight',
        '"test_task_123"': 'constants.IntegrationTestDataTask + "_123"',
        '"test_action_1"': 'constants.IntegrationTestDataAction + "_1"',
        '"test_analysis"': 'constants.IntegrationTestDataAnalysis',
        '"test_report"': 'constants.IntegrationTestDataReport',
        '"test_metrics"': 'constants.IntegrationTestDataMetrics',
        
        # Analysis Types
        '"content"': 'constants.IntegrationAnalysisTypeContent',
        '"technical"': 'constants.IntegrationAnalysisTypeTechnical',
        '"performance"': 'constants.IntegrationAnalysisTypePerformance',
        '"security"': 'constants.IntegrationAnalysisTypeSecurity',
        '"accessibility"': 'constants.IntegrationAnalysisTypeAccessibility',
        '"seo"': 'constants.IntegrationAnalysisTypeSEO',
        '"semantic"': 'constants.IntegrationAnalysisTypeSemantic',
        '"structural"': 'constants.IntegrationAnalysisTypeStructural',
        '"functional"': 'constants.IntegrationAnalysisTypeFunctional',
        '"regression"': 'constants.IntegrationAnalysisTypeRegression',
        '"smoke"': 'constants.IntegrationAnalysisTypeSmoke',
        '"sanity"': 'constants.IntegrationAnalysisTypeSanity',
        
        # Agent Names (seulement dans contexte approprié)
        '"crawler"': 'constants.IntegrationAgentCrawler',
        '"seo"': 'constants.IntegrationAgentSEO',
        '"semantic"': 'constants.IntegrationAgentSemantic',
        '"performance"': 'constants.IntegrationAgentPerformance',
        '"security"': 'constants.IntegrationAgentSecurity',
        '"qa"': 'constants.IntegrationAgentQA',
        '"data_integrity"': 'constants.IntegrationAgentDataIntegrity',
        '"frontend"': 'constants.IntegrationAgentFrontend',
        '"playwright"': 'constants.IntegrationAgentPlaywright',
        '"k6"': 'constants.IntegrationAgentK6',
        '"integration"': 'constants.IntegrationAgentIntegration',
        '"unit"': 'constants.IntegrationAgentUnit',
        '"e2e"': 'constants.IntegrationAgentE2E',
        
        # Test Categories
        '"unit"': 'constants.IntegrationTestCategoryUnit',
        '"integration"': 'constants.IntegrationTestCategoryIntegration',
        '"e2e"': 'constants.IntegrationTestCategoryE2E',
        '"api"': 'constants.IntegrationTestCategoryAPI',
        '"ui"': 'constants.IntegrationTestCategoryUI',
        '"performance"': 'constants.IntegrationTestCategoryPerformance',
        '"security"': 'constants.IntegrationTestCategorySecurity',
        '"accessibility"': 'constants.IntegrationTestCategoryAccessibility',
        '"compatibility"': 'constants.IntegrationTestCategoryCompatibility',
        '"regression"': 'constants.IntegrationTestCategoryRegression',
        '"smoke"': 'constants.IntegrationTestCategorySmoke',
        '"sanity"': 'constants.IntegrationTestCategorySanity',
        '"acceptance"': 'constants.IntegrationTestCategoryAcceptance',
        '"system"': 'constants.IntegrationTestCategorySystem',
        
        # Error Messages (seulement les messages complets, pas les formats)
        '"URL incorrecte: %s"': 'constants.IntegrationErrorURLInvalid + ": %s"',
        '"Score global invalide: %.1f"': 'constants.IntegrationErrorScoreInvalid + ": %.1f"',
        '"Temps de traitement invalide"': 'constants.IntegrationErrorProcessingTimeInvalid',
        '"Taille de rapport invalide"': 'constants.IntegrationErrorReportSizeInvalid',
        '"Format de données invalide"': 'constants.IntegrationErrorDataFormatInvalid',
        
        # JSON Field Names (seulement dans contexte JSON)
        '"id"': 'constants.IntegrationJSONFieldID',
        '"name"': 'constants.IntegrationJSONFieldName',
        '"type"': 'constants.IntegrationJSONFieldType',
        '"status"': 'constants.IntegrationJSONFieldStatus',
        '"data"': 'constants.IntegrationJSONFieldData',
        '"message"': 'constants.IntegrationJSONFieldMessage',
        '"error"': 'constants.IntegrationJSONFieldError',
        '"timestamp"': 'constants.IntegrationJSONFieldTimestamp',
        '"url"': 'constants.IntegrationJSONFieldURL',
        '"method"': 'constants.IntegrationJSONFieldMethod',
        '"path"': 'constants.IntegrationJSONFieldPath',
        '"result"': 'constants.IntegrationJSONFieldResult',
        '"response"': 'constants.IntegrationJSONFieldResponse',
        '"request"': 'constants.IntegrationJSONFieldRequest',
        '"duration"': 'constants.IntegrationJSONFieldDuration',
        '"timeout"': 'constants.IntegrationJSONFieldTimeout',
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
        r'assert\.', # assertion methods
        r'require\.', # require methods
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
        
        if not should_avoid:
            # Appliquer les remplacements
            for old_value, new_value in replacements.items():
                if old_value in line:
                    # Vérifications supplémentaires pour éviter les faux positifs
                    
                    # Pour les JSON field names, s'assurer qu'on est dans un contexte JSON
                    json_fields = ['"id"', '"name"', '"type"', '"status"', '"data"', '"message"', 
                                 '"error"', '"timestamp"', '"url"', '"method"', '"path"', '"result"',
                                 '"response"', '"request"', '"duration"', '"timeout"']
                    
                    if old_value in json_fields:
                        # Seulement remplacer si c'est dans une structure JSON (avec : ou [])
                        if (':' in line or '[' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf', 'Sprintf', 't.', 'assert.', 'require.']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les agent names et analysis types, vérifier le contexte
                    elif old_value in ['"crawler"', '"seo"', '"semantic"', '"performance"', '"security"', 
                                     '"qa"', '"data_integrity"', '"frontend"', '"playwright"', '"k6"',
                                     '"integration"', '"unit"', '"e2e"', '"content"', '"technical"',
                                     '"accessibility"', '"structural"', '"functional"', '"regression"',
                                     '"smoke"', '"sanity"']:
                        # Seulement remplacer dans des contextes appropriés (affectation, comparaison, etc.)
                        if ('=' in line or ':' in line or 'assert' in line or 'expect' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf', 'Sprintf']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les error messages et test data, remplacement direct
                    else:
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
        for change in changes_made[:10]:  # Montrer les 10 premiers
            print(f"  Line {change['line']}: {change['old']} → {change['new']}")
        
        if len(changes_made) > 10:
            print(f"  ... et {len(changes_made) - 10} autres changements")
        
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
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/integration_test.go'
    
    print("⚔️ DELTA-4 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: integration_test.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-4 SUCCÈS!")
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
            
            with open('delta4_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta4_elimination_report.json")
            
        else:
            print(f"\n⚠️ Aucune violation éliminée")
            print("Les violations détectées peuvent être dans des contextes non remplaçables")
            
    except Exception as e:
        print(f"\n❌ ERREUR durant l'élimination: {e}")
        # Restaurer le backup en cas d'erreur
        shutil.copy2(backup_path, filepath)
        print(f"🔄 Fichier restauré depuis le backup")
        raise
    
    return eliminated

if __name__ == "__main__":
    main()