#!/usr/bin/env python3
"""
🔍 DELTA SMART FILTER
Filtrer les vraies violations vs définitions de constantes
"""

import os
import re
import json
from pathlib import Path

def is_constants_file(filepath):
    """Vérifie si c'est un fichier de définition de constantes"""
    return (
        'constants' in filepath.lower() or
        filepath.endswith('_constants.go') or
        'constants.go' in filepath
    )

def is_real_violation(line, value, filepath):
    """Détermine si c'est une vraie violation ou une définition de constante"""
    
    # Si c'est un fichier de constantes, ignorer les définitions
    if is_constants_file(filepath):
        return False
    
    # Ignorer les struct tags JSON
    if '`json:' in line:
        return False
        
    # Ignorer les commentaires
    if line.strip().startswith('//'):
        return False
        
    # Ignorer les imports
    if 'import' in line or 'package' in line:
        return False
    
    # Ignorer les déclarations de constantes dans les fichiers normaux
    if re.match(r'^\s*const\s+\w+\s*=', line.strip()):
        return False
        
    # Ignorer les déclarations de types enum
    if re.match(r'^\s*\w+\s+\w+Type\s*=', line.strip()):
        return False
    
    return True

def scan_real_violations(base_path):
    """Scanne seulement les vraies violations (pas les définitions de constantes)"""
    
    real_violations = []
    file_counts = {}
    
    exclude_patterns = [
        '.git', 'node_modules', 'vendor', '.vscode', '.idea',
        'dist', 'build', 'target', '__pycache__', '.pytest_cache',
        '*.log', '*.tmp', '*.backup', '*.bak',
        'charlie1_backup', 'charlie2_backup', 'bravo3_backup'
    ]
    
    patterns = {
        'string_literals': r'"[A-Za-z][^"]{3,}"',
        'error_messages': r'"[A-Z][^"]*(?:error|failed|unable|invalid|missing|not found)[^"]*"',
        'log_messages': r'"[A-Z][^"]*(?:starting|completed|processing|initializing|validating)[^"]*"',
        'status_values': r'"(?:pending|running|completed|failed|success|error|warning|info)"',
        'http_methods': r'"(?:GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'file_paths': r'"[./][^"]*(?:/[^"]*)*"',
    }
    
    for root, dirs, files in os.walk(base_path):
        # Exclure les répertoires
        dirs[:] = [d for d in dirs if not any(pattern in d for pattern in exclude_patterns)]
        
        for file in files:
            if file.endswith('.go'):
                filepath = os.path.join(root, file)
                relative_path = os.path.relpath(filepath, base_path)
                
                # Exclure les fichiers backup
                if any(pattern in relative_path for pattern in exclude_patterns):
                    continue
                
                file_violations = []
                
                try:
                    with open(filepath, 'r', encoding='utf-8') as f:
                        content = f.read()
                except (UnicodeDecodeError, IOError):
                    continue
                
                line_num = 0
                
                for line in content.split('\n'):
                    line_num += 1
                    line_stripped = line.strip()
                    
                    for category, pattern in patterns.items():
                        matches = re.findall(pattern, line_stripped)
                        for match in matches:
                            if len(match.strip('"')) < 3:
                                continue
                                
                            if is_real_violation(line_stripped, match, relative_path):
                                violation = {
                                    'line': line_num,
                                    'category': category,
                                    'value': match,
                                    'file': relative_path,
                                    'context': line_stripped[:100]
                                }
                                file_violations.append(violation)
                                real_violations.append(violation)
                
                if file_violations:
                    file_counts[relative_path] = len(file_violations)
                    print(f"📁 {relative_path}: {len(file_violations)} VRAIES violations")
    
    return real_violations, file_counts

def generate_real_delta_strategy(file_counts):
    """Génère la stratégie DELTA pour les vraies violations"""
    
    # Trier par nombre de violations
    sorted_files = sorted(file_counts.items(), key=lambda x: x[1], reverse=True)
    
    delta_sprints = []
    sprint_num = 1
    
    for filepath, violation_count in sorted_files:
        if violation_count >= 15:  # Seulement les fichiers avec 15+ vraies violations
            delta_sprints.append({
                'sprint': f'DELTA-{sprint_num}',
                'target': os.path.basename(filepath),
                'filepath': filepath,
                'violations': violation_count,
                'effort': '0.5h' if violation_count < 20 else '1h' if violation_count < 40 else '1.5h'
            })
            sprint_num += 1
            
            if sprint_num > 15:  # Limiter à 15 sprints
                break
    
    return delta_sprints

def main():
    base_path = '/Users/jeromegonzalez/claude-code/fire-salamander'
    
    print("🔍 DELTA SMART FILTER - Identification des VRAIES violations...")
    print("=" * 70)
    
    real_violations, file_counts = scan_real_violations(base_path)
    
    print(f"\n📊 RÉSULTATS APRÈS FILTRAGE INTELLIGENT:")
    print(f"Vraies violations détectées: {len(real_violations)}")
    print(f"Fichiers réellement affectés: {len(file_counts)}")
    
    if len(file_counts) > 0:
        print(f"\n🎯 TOP 10 FICHIERS AVEC VRAIES VIOLATIONS:")
        sorted_files = sorted(file_counts.items(), key=lambda x: x[1], reverse=True)
        for i, (filepath, count) in enumerate(sorted_files[:10], 1):
            print(f"  {i:2d}. {filepath}: {count} violations")
        
        # Générer la vraie stratégie DELTA
        delta_sprints = generate_real_delta_strategy(file_counts)
        
        print(f"\n⚔️ STRATÉGIE DELTA RÉELLE:")
        print(f"Sprints prioritaires: {len(delta_sprints)}")
        total_real_violations = sum(sprint['violations'] for sprint in delta_sprints)
        
        for sprint in delta_sprints:
            print(f"  {sprint['sprint']}: {sprint['target']} ({sprint['violations']} violations - {sprint['effort']})")
        
        print(f"\nTotal violations DELTA réelles: {total_real_violations}")
        
        # Sauvegarder
        results = {
            'real_violations': len(real_violations),
            'files_affected': len(file_counts),
            'top_files': sorted_files,
            'delta_sprints': delta_sprints,
        }
        
        with open('delta_real_analysis.json', 'w') as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        
        print(f"\n✅ Analyse réelle sauvegardée dans delta_real_analysis.json")
    else:
        print(f"\n🎉 AUCUNE VRAIE VIOLATION DÉTECTÉE !")
        print(f"Toutes les 'violations' détectées étaient des définitions légitimes de constantes !")
    
    return len(real_violations)

if __name__ == "__main__":
    main()