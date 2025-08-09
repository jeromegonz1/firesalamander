#!/usr/bin/env python3
"""
⚔️ DELTA-2 ELIMINATOR
Éliminateur spécialisé pour api.go - 60 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta2_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-2"""
    with open('delta2_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour api.go
    replacements = {
        # HTTP Endpoints API V1
        '"api/v1/analyze"': 'constants.APIEndpointV1Analyze',
        '"api/v1/analyze/semantic"': 'constants.APIEndpointV1AnalyzeSemantic',
        '"api/v1/analyze/seo"': 'constants.APIEndpointV1AnalyzeSEO',
        '"api/v1/analyze/quick"': 'constants.APIEndpointV1AnalyzeQuick',
        '"api/v1/health"': 'constants.APIEndpointV1Health',
        '"api/v1/stats"': 'constants.APIEndpointV1Stats',
        '"api/v1/analyses"': 'constants.APIEndpointV1Analyses',
        '"api/v1/analysis/"': 'constants.APIEndpointV1Analysis',
        '"api/v1/info"': 'constants.APIEndpointV1Info',
        '"api/v1/version"': 'constants.APIEndpointV1Version',
        '"api/v1/reports"': 'constants.APIEndpointV1Reports',
        '"api/v1/results"': 'constants.APIEndpointV1Results',
        '"api/v1/metrics"': 'constants.APIEndpointV1Metrics',
        '"api/v1/status"': 'constants.APIEndpointV1Status',
        '"api/v1/debug"': 'constants.APIEndpointV1Debug',
        
        # HTTP Methods
        '"GET"': 'constants.APIMethodGet',
        '"POST"': 'constants.APIMethodPost',
        '"PUT"': 'constants.APIMethodPut',
        '"DELETE"': 'constants.APIMethodDelete',
        '"PATCH"': 'constants.APIMethodPatch',
        '"HEAD"': 'constants.APIMethodHead',
        '"OPTIONS"': 'constants.APIMethodOptions',
        
        # Content Types
        '"application/json"': 'constants.APIContentTypeJSON',
        '"text/html"': 'constants.APIContentTypeHTML',
        '"text/plain"': 'constants.APIContentTypePlain',
        '"application/xml"': 'constants.APIContentTypeXML',
        '"text/csv"': 'constants.APIContentTypeCSV',
        '"application/pdf"': 'constants.APIContentTypePDF',
        
        # HTTP Headers
        '"Content-Type"': 'constants.APIHeaderContentType',
        '"Authorization"': 'constants.APIHeaderAuthorization',
        '"Accept"': 'constants.APIHeaderAccept',
        '"User-Agent"': 'constants.APIHeaderUserAgent',
        '"Cache-Control"': 'constants.APIHeaderCacheControl',
        '"Origin"': 'constants.APIHeaderOrigin',
        '"Access-Control-Allow-Origin"': 'constants.APIHeaderAccessControlAllowOrigin',
        '"Access-Control-Allow-Methods"': 'constants.APIHeaderAccessControlAllowMethods',
        '"Access-Control-Allow-Headers"': 'constants.APIHeaderAccessControlAllowHeaders',
        '"Access-Control-Expose-Headers"': 'constants.APIHeaderAccessControlExposeHeaders',
        '"Access-Control-Max-Age"': 'constants.APIHeaderAccessControlMaxAge',
        '"Access-Control-Allow-Credentials"': 'constants.APIHeaderAccessControlAllowCredentials',
        
        # JSON Field Names (seulement dans contexte JSON)
        '"id"': 'constants.APIJSONFieldID',
        '"name"': 'constants.APIJSONFieldName',
        '"type"': 'constants.APIJSONFieldType',
        '"status"': 'constants.APIJSONFieldStatus',
        '"data"': 'constants.APIJSONFieldData',
        '"message"': 'constants.APIJSONFieldMessage',
        '"error"': 'constants.APIJSONFieldError',
        '"timestamp"': 'constants.APIJSONFieldTimestamp',
        '"url"': 'constants.APIJSONFieldURL',
        '"method"': 'constants.APIJSONFieldMethod',
        '"path"': 'constants.APIJSONFieldPath',
        '"results"': 'constants.APIJSONFieldResults',
        '"recommendations"': 'constants.APIJSONFieldRecommendations',
        '"metrics"': 'constants.APIJSONFieldMetrics',
        '"score"': 'constants.APIJSONFieldScore',
        '"title"': 'constants.APIJSONFieldTitle',
        '"description"': 'constants.APIJSONFieldDescription',
        '"content"': 'constants.APIJSONFieldContent',
        '"value"': 'constants.APIJSONFieldValue',
        '"category"': 'constants.APIJSONFieldCategory',
        '"priority"': 'constants.APIJSONFieldPriority',
        '"severity"': 'constants.APIJSONFieldSeverity',
        
        # Status Values
        '"pending"': 'constants.APIStatusPending',
        '"running"': 'constants.APIStatusRunning',
        '"completed"': 'constants.APIStatusCompleted',
        '"failed"': 'constants.APIStatusFailed',
        '"success"': 'constants.APIStatusSuccess',
        '"error"': 'constants.APIStatusError',
        '"warning"': 'constants.APIStatusWarning',
        '"info"': 'constants.APIStatusInfo',
        '"healthy"': 'constants.APIStatusHealthy',
        '"degraded"': 'constants.APIStatusDegraded',
        '"critical"': 'constants.APIStatusCritical',
        
        # Agent Names
        '"crawler"': 'constants.APIAgentCrawler',
        '"seo"': 'constants.APIAgentSEO',
        '"semantic"': 'constants.APIAgentSemantic',
        '"performance"': 'constants.APIAgentPerformance',
        '"security"': 'constants.APIAgentSecurity',
        '"qa"': 'constants.APIAgentQA',
        '"data_integrity"': 'constants.APIAgentDataIntegrity',
        '"frontend"': 'constants.APIAgentFrontend',
        '"playwright"': 'constants.APIAgentPlaywright',
        '"k6"': 'constants.APIAgentK6',
        
        # Error Messages
        '"Requête JSON invalide: "': 'constants.APIErrorInvalidJSON + ": "',
        '"Requête JSON invalide:"': 'constants.APIErrorInvalidJSON + ":"',
        '"Requête JSON invalide"': 'constants.APIErrorInvalidJSON',
    }
    
    # Contextes à éviter (où ne PAS remplacer)
    avoid_contexts = [
        r'`json:', # struct tags
        r'//', # commentaires
        r'import', # imports
        r'package', # package declaration
        r'const\s+\w+\s*=', # déclarations de constantes
        r'fmt\.', # format strings
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
                                 '"error"', '"timestamp"', '"url"', '"method"', '"path"', '"results"',
                                 '"recommendations"', '"metrics"', '"score"', '"title"', '"description"',
                                 '"content"', '"value"', '"category"', '"priority"', '"severity"']
                    
                    if old_value in json_fields:
                        # Seulement remplacer si c'est dans une structure JSON (avec : ou [])
                        if (':' in line or '[' in line) and not 'fmt.' in line:
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les status values et agent names, même vérification
                    elif old_value in ['"pending"', '"running"', '"completed"', '"failed"', '"success"', 
                                     '"error"', '"warning"', '"info"', '"healthy"', '"degraded"', '"critical"',
                                     '"crawler"', '"seo"', '"semantic"', '"performance"', '"security"', '"qa"',
                                     '"data_integrity"', '"frontend"', '"playwright"', '"k6"']:
                        if (':' in line or '=' in line) and not 'fmt.' in line:
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les autres (endpoints, methods, headers, etc.), remplacement direct
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
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/api.go'
    
    print("⚔️ DELTA-2 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: api.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-2 SUCCÈS!")
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
            
            with open('delta2_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta2_elimination_report.json")
            
        else:
            print(f"\n⚠️ Aucune violation éliminée")
            print("Vérifier les patterns de remplacement")
            
    except Exception as e:
        print(f"\n❌ ERREUR durant l'élimination: {e}")
        # Restaurer le backup en cas d'erreur
        shutil.copy2(backup_path, filepath)
        print(f"🔄 Fichier restauré depuis le backup")
        raise
    
    return eliminated

if __name__ == "__main__":
    main()