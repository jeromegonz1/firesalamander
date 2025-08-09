#!/usr/bin/env python3
"""
⚔️ DELTA-3 ELIMINATOR
Éliminateur spécialisé pour ngram_analyzer.go - 6 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta3_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-3"""
    with open('delta3_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour ngram_analyzer.go
    replacements = {
        # Analysis Types (N-Gram)
        '"unigram"': 'constants.SemanticAnalysisUnigram',
        '"bigram"': 'constants.SemanticAnalysisBigram',
        '"trigram"': 'constants.SemanticAnalysisTrigram',
        '"4-gram"': 'constants.SemanticAnalysisFourGram',
        '"5-gram"': 'constants.SemanticAnalysisFiveGram',
        '"n-gram"': 'constants.SemanticAnalysisNGram',
        
        # NLP Constants
        '"unicode"': 'constants.SemanticNLPUnicode',
        '"ascii"': 'constants.SemanticNLPASCII',
        '"utf8"': 'constants.SemanticNLPUTF8',
        '"stopwords"': 'constants.SemanticNLPStopwords',
        '"punctuation"': 'constants.SemanticNLPPunctuation',
        '"whitespace"': 'constants.SemanticNLPWhitespace',
        
        # Data Types
        '"online"': 'constants.SemanticDataTypeOnline',
        '"offline"': 'constants.SemanticDataTypeOffline',
        '"batch"': 'constants.SemanticDataTypeBatch',
        '"stream"': 'constants.SemanticDataTypeStream',
        '"real_time"': 'constants.SemanticDataTypeRealTime',
        
        # Semantic Field Names (seulement dans contexte approprié)
        '"text"': 'constants.SemanticFieldText',
        '"content"': 'constants.SemanticFieldContent',
        '"token"': 'constants.SemanticFieldToken',
        '"word"': 'constants.SemanticFieldWord',
        '"phrase"': 'constants.SemanticFieldPhrase',
        '"sentence"': 'constants.SemanticFieldSentence',
        '"document"': 'constants.SemanticFieldDocument',
        '"frequency"': 'constants.SemanticFieldFrequency',
        '"score"': 'constants.SemanticFieldScore',
        '"weight"': 'constants.SemanticFieldWeight',
        '"similarity"': 'constants.SemanticFieldSimilarity',
        '"distance"': 'constants.SemanticFieldDistance',
        '"threshold"': 'constants.SemanticFieldThreshold',
        
        # Analysis Methods
        '"tf-idf"': 'constants.SemanticAnalysisTFIDF',
        '"cosine"': 'constants.SemanticAnalysisCosine',
        '"euclidean"': 'constants.SemanticAnalysisEuclidean',
        '"jaccard"': 'constants.SemanticAnalysisJaccard',
        '"levenshtein"': 'constants.SemanticAnalysisLevenshtein',
        '"stemming"': 'constants.SemanticAnalysisStemming',
        '"lemmatization"': 'constants.SemanticAnalysisLemmatization',
        '"tokenization"': 'constants.SemanticAnalysisTokenization',
        '"normalization"': 'constants.SemanticAnalysisNormalization',
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
                    
                    # Pour les semantic field names, s'assurer qu'on est dans un contexte approprié
                    semantic_fields = ['"text"', '"content"', '"token"', '"word"', '"phrase"', 
                                     '"sentence"', '"document"', '"frequency"', '"score"', '"weight"',
                                     '"similarity"', '"distance"', '"threshold"']
                    
                    if old_value in semantic_fields:
                        # Seulement remplacer si c'est dans une structure de données (avec : ou =)
                        if (':' in line or '=' in line or '[' in line) and not any(avoid in line for avoid in ['fmt.', 'log.', 'Printf', 'Sprintf']):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les analysis types et data types, remplacement direct
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
        for change in changes_made:
            print(f"  Line {change['line']}: {change['old']} → {change['new']}")
        
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
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/semantic/ngram_analyzer.go'
    
    print("⚔️ DELTA-3 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: ngram_analyzer.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-3 SUCCÈS!")
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
            
            with open('delta3_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta3_elimination_report.json")
            
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