#!/usr/bin/env python3
"""
⚔️ DELTA-1 ELIMINATOR
Éliminateur spécialisé pour server.go - 29 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta1_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-1"""
    with open('delta1_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour server.go
    replacements = {
        # HTTP Endpoints
        '"/api/"': 'constants.ServerEndpointAPI',
        
        # Content Types
        '"application/json"': 'constants.ServerContentTypeJSON',
        '"text/html"': 'constants.ServerContentTypeHTML',
        
        # HTTP Headers
        '"Content-Type"': 'constants.ServerHeaderContentType',
        
        # JSON Field Names (dans les structures JSON)
        '"status"': 'constants.ServerJSONFieldStatus',
        '"timestamp"': 'constants.ServerJSONFieldTimestamp',
        '"url"': 'constants.ServerJSONFieldURL',
        '"id"': 'constants.ServerJSONFieldID',
        '"type"': 'constants.ServerJSONFieldType',
        
        # File Extensions
        '".html"': 'constants.ServerExtensionHTML',
        '".json"': 'constants.ServerExtensionJSON',
        
        # Server Config Keys
        '"port"': 'constants.ServerConfigPort',
    }
    
    # Contextes à éviter (où ne PAS remplacer)
    avoid_contexts = [
        r'`json:', # struct tags
        r'//', # commentaires
        r'import', # imports
        r'package', # package declaration
        r'const\s+\w+\s*=', # déclarations de constantes
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
                    if old_value in ['"status"', '"timestamp"', '"url"', '"id"', '"type"']:
                        # Seulement remplacer si c'est dans une structure JSON (avec : ou [])
                        if ':' in line or '[' in line:
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les autres, remplacement direct
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
    
    # Ajouter l'import après les autres imports internes
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
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/server.go'
    
    print("⚔️ DELTA-1 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: server.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-1 SUCCÈS!")
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
            
            with open('delta1_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta1_elimination_report.json")
            
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