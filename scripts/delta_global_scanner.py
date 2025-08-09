#!/usr/bin/env python3
"""
🔍 DELTA GLOBAL SCANNER
Scanner industriel pour identifier les 568 violations restantes
"""

import os
import re
import json
from collections import defaultdict
from pathlib import Path

def scan_file_for_violations(filepath):
    """Scanne un fichier pour les violations hardcoding"""
    violations = []
    
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except (UnicodeDecodeError, IOError):
        return violations
    
    line_num = 0
    
    # Patterns universels de détection
    patterns = {
        'string_literals': r'"[A-Za-z][^"]{3,}"',  # Strings de plus de 3 caractères
        'error_messages': r'"[A-Z][^"]*(?:error|failed|unable|invalid|missing|not found)[^"]*"',
        'log_messages': r'"[A-Z][^"]*(?:starting|completed|processing|initializing|validating)[^"]*"',
        'config_keys': r'"[a-z_]+(?:_key|_path|_url|_port|_timeout|_size|_count)"',
        'http_methods': r'"(?:GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'status_values': r'"(?:pending|running|completed|failed|success|error|warning|info)"',
        'file_extensions': r'"\.[a-z]+(?:\.[a-z]+)?"',
        'mime_types': r'"(?:application|text|image|audio|video)/[a-z-]+"',
        'url_schemes': r'"(?:https?|ftp|file|ws)://"',
        'time_formats': r'"(?:\d{4}-\d{2}-\d{2}|\d{2}:\d{2}:\d{2}|RFC3339)"',
        'numeric_strings': r'"\d+(?:\.\d+)?"',
        'database_operations': r'"(?:SELECT|INSERT|UPDATE|DELETE|CREATE|DROP)"',
        'json_field_names': r'"[a-z_]+(?:_id|_name|_type|_status|_data)"',
        'html_tags': r'"</?[a-zA-Z][^">]*>?"',
        'css_selectors': r'"[.#][a-zA-Z][^"]*"',
        'regex_patterns': r'`[^`]{10,}`',  # Regex literals
        'file_paths': r'"[./][^"]*(?:/[^"]*)*"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires et imports
        if line_stripped.startswith('//') or line_stripped.startswith('import') or line_stripped.startswith('package'):
            continue
            
        # Ignorer les struct tags
        if '`json:' in line_stripped:
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped)
            for match in matches:
                # Filtrer les matches trop courts ou trop génériques
                if len(match.strip('"')) < 3:
                    continue
                    
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match,
                    'file': filepath,
                    'context': line_stripped[:100] + ('...' if len(line_stripped) > 100 else '')
                })
    
    return violations

def scan_directory_recursive(base_path, exclude_patterns=None):
    """Scanne récursivement un répertoire pour les violations"""
    
    if exclude_patterns is None:
        exclude_patterns = [
            '.git', 'node_modules', 'vendor', '.vscode', '.idea',
            'dist', 'build', 'target', '__pycache__', '.pytest_cache',
            '*.log', '*.tmp', '*.backup', '*.bak',
            'charlie1_backup', 'charlie2_backup', 'bravo3_backup'
        ]
    
    all_violations = []
    file_counts = defaultdict(int)
    
    for root, dirs, files in os.walk(base_path):
        # Exclure les répertoires
        dirs[:] = [d for d in dirs if not any(pattern in d for pattern in exclude_patterns)]
        
        for file in files:
            if file.endswith('.go'):
                filepath = os.path.join(root, file)
                relative_path = os.path.relpath(filepath, base_path)
                
                # Exclure les fichiers backup et tests générés
                if any(pattern in relative_path for pattern in exclude_patterns):
                    continue
                
                violations = scan_file_for_violations(filepath)
                all_violations.extend(violations)
                
                if violations:
                    file_counts[relative_path] = len(violations)
                    print(f"📁 {relative_path}: {len(violations)} violations")
    
    return all_violations, dict(file_counts)

def analyze_violations(violations):
    """Analyse les violations par catégorie et fichier"""
    
    by_category = defaultdict(list)
    by_file = defaultdict(list)
    
    for violation in violations:
        by_category[violation['category']].append(violation)
        by_file[violation['file']].append(violation)
    
    return dict(by_category), dict(by_file)

def generate_delta_strategy(file_counts):
    """Génère la stratégie d'assaut DELTA basée sur les violations détectées"""
    
    # Trier les fichiers par nombre de violations (descendant)
    sorted_files = sorted(file_counts.items(), key=lambda x: x[1], reverse=True)
    
    # Créer les sprints DELTA
    delta_sprints = []
    sprint_num = 1
    
    for filepath, violation_count in sorted_files[:15]:  # Top 15 files
        if violation_count >= 10:  # Seulement les fichiers avec 10+ violations
            delta_sprints.append({
                'sprint': f'DELTA-{sprint_num}',
                'target': os.path.basename(filepath),
                'filepath': filepath,
                'violations': violation_count,
                'effort': '1h' if violation_count < 30 else '1.5h' if violation_count < 50 else '2h'
            })
            sprint_num += 1
    
    return delta_sprints

def main():
    base_path = '/Users/jeromegonzalez/claude-code/fire-salamander'
    
    print("🔍 DELTA GLOBAL SCANNER - Scanning complet pour violations restantes...")
    print("=" * 70)
    
    # Scanner tous les fichiers Go
    all_violations, file_counts = scan_directory_recursive(base_path)
    
    # Analyser les résultats
    by_category, by_file = analyze_violations(all_violations)
    
    print(f"\n📊 RÉSULTATS SCAN GLOBAL:")
    print(f"Total violations détectées: {len(all_violations)}")
    print(f"Fichiers affectés: {len(file_counts)}")
    
    print(f"\n🎯 TOP 10 FICHIERS LES PLUS CRITIQUES:")
    sorted_files = sorted(file_counts.items(), key=lambda x: x[1], reverse=True)
    for i, (filepath, count) in enumerate(sorted_files[:10], 1):
        relative_path = os.path.relpath(filepath, base_path)
        print(f"  {i:2d}. {relative_path}: {count} violations")
    
    print(f"\n🔸 VIOLATIONS PAR CATÉGORIE:")
    sorted_categories = sorted(by_category.items(), key=lambda x: len(x[1]), reverse=True)
    for category, viols in sorted_categories[:10]:
        print(f"  {category.upper()}: {len(viols)} violations")
    
    # Générer la stratégie DELTA
    delta_sprints = generate_delta_strategy(file_counts)
    
    print(f"\n⚔️ STRATÉGIE D'ASSAUT DELTA:")
    print(f"Sprints planifiés: {len(delta_sprints)}")
    total_delta_violations = sum(sprint['violations'] for sprint in delta_sprints)
    
    for sprint in delta_sprints:
        print(f"  {sprint['sprint']}: {sprint['target']} ({sprint['violations']} violations - {sprint['effort']})")
    
    print(f"\nTotal violations DELTA ciblées: {total_delta_violations}")
    
    # Sauvegarder les résultats
    results = {
        'total_violations': len(all_violations),
        'files_affected': len(file_counts),
        'top_files': sorted_files[:20],
        'categories': {k: len(v) for k, v in by_category.items()},
        'delta_sprints': delta_sprints,
        'violations': all_violations[:100],  # Limiter pour la taille du fichier
    }
    
    with open('delta_global_analysis.json', 'w') as f:
        json.dump(results, f, indent=2, ensure_ascii=False)
    
    print(f"\n✅ Analyse globale sauvegardée dans delta_global_analysis.json")
    
    return results

if __name__ == "__main__":
    main()