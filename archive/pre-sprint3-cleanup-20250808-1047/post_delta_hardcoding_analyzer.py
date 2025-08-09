#!/usr/bin/env python3
"""
FIRE SALAMANDER - Post-DELTA Hardcoding Analyzer
===============================================

Analyseur complet de hardcoding post-corrections DELTA.
Mesure l'impact réel de nos corrections et identifie les violations restantes.

OBJECTIFS:
- Détecter toutes les violations de hardcoding restantes
- Comparer avec l'analyse initiale (4,582 violations)
- Catégoriser par type et sévérité
- Exclure les fichiers de constantes (déjà traités)
- Générer un plan d'action prioritisé

AMÉLIORATIONS:
- Détection plus précise (éviter faux positifs)
- Exclusion des imports, struct tags, commentaires
- Détection contextuelle améliorée
- Scoring de sévérité (Critical, High, Medium, Low)
"""

import os
import re
import json
import ast
from datetime import datetime
from collections import defaultdict
from typing import Dict, List, Tuple, Any, Optional
from dataclasses import dataclass
from enum import Enum

class Severity(Enum):
    CRITICAL = "Critical"
    HIGH = "High"
    MEDIUM = "Medium"
    LOW = "Low"

@dataclass
class Violation:
    file_path: str
    line: int
    category: str
    value: str
    context: str
    severity: Severity
    description: str

class PostDeltaHardcodingAnalyzer:
    def __init__(self, root_path: str):
        self.root_path = root_path
        self.violations = []
        self.file_stats = defaultdict(int)
        self.category_stats = defaultdict(int)
        self.severity_stats = defaultdict(int)
        
        # Fichiers à exclure (déjà corrigés ou non pertinents)
        self.excluded_files = {
            'internal/constants/',  # Fichiers de constantes déjà créés
            'node_modules/',
            'vendor/',
            '.git/',
            'test_',
            '_test.go',
            '.bak',
            'backup',
            'archive/',
            'backups/',
        }
        
        # Extensions à inclure
        self.included_extensions = {'.go'}
        
        # Patterns de hardcoding améliorés avec scoring
        self.patterns = {
            # CRITICAL - Endpoints et URLs sensibles
            'api_endpoints': {
                'patterns': [
                    r'"(/api/[^"]*)"',
                    r'"(/v\d+/[^"]*)"',
                    r'"(https?://[^"]*)"',
                    r'`(/api/[^`]*)`',
                    r'HandleFunc\(["\']([^"\']*)["\']',
                ],
                'severity': Severity.CRITICAL,
                'description': 'Endpoints API ou URLs hardcodées'
            },
            
            # HIGH - Configuration et sécurité
            'database_config': {
                'patterns': [
                    r'"(postgres://[^"]*)"',
                    r'"(mysql://[^"]*)"',
                    r'"(mongodb://[^"]*)"',
                    r'"(redis://[^"]*)"',
                    r'DSN.*["\']([^"\']*)["\']',
                ],
                'severity': Severity.HIGH,
                'description': 'Configuration base de données hardcodée'
            },
            
            'security_keys': {
                'patterns': [
                    r'"(sk-[a-zA-Z0-9]{20,})"',
                    r'"([A-Za-z0-9]{32,})"',  # Clés potentielles
                    r'token.*["\']([^"\']{20,})["\']',
                ],
                'severity': Severity.HIGH,
                'description': 'Clés de sécurité potentiellement hardcodées'
            },
            
            'server_config': {
                'patterns': [
                    r':[0-9]{4,5}["\']',  # Ports
                    r'"(localhost[^"]*)"',
                    r'"(127\.0\.0\.1[^"]*)"',
                    r'"(0\.0\.0\.0[^"]*)"',
                ],
                'severity': Severity.HIGH,
                'description': 'Configuration serveur hardcodée'
            },
            
            # MEDIUM - Headers HTTP et content types
            'http_headers': {
                'patterns': [
                    r'Header\(\)\.Set\(["\']([A-Za-z-]+)["\']',
                    r'Header\(\)\.Add\(["\']([A-Za-z-]+)["\']',
                    r'Header\(\)\.Get\(["\']([A-Za-z-]+)["\']',
                ],
                'severity': Severity.MEDIUM,
                'description': 'Headers HTTP hardcodés'
            },
            
            'content_types': {
                'patterns': [
                    r'"(application/[^"]*)"',
                    r'"(text/[^"]*)"',
                    r'"(image/[^"]*)"',
                    r'"(multipart/[^"]*)"',
                ],
                'severity': Severity.MEDIUM,
                'description': 'Content-Types hardcodés'
            },
            
            'http_methods': {
                'patterns': [
                    r'Method.*["\']([A-Z]+)["\']',
                    r'\.([GET|POST|PUT|DELETE|PATCH]+)\(',
                ],
                'severity': Severity.MEDIUM,
                'description': 'Méthodes HTTP hardcodées'
            },
            
            # MEDIUM - Messages et interface
            'error_messages': {
                'patterns': [
                    r'fmt\.Errorf\(["\']([^"\']{20,})["\']',
                    r'errors\.New\(["\']([^"\']{10,})["\']',
                    r'panic\(["\']([^"\']*)["\']',
                ],
                'severity': Severity.MEDIUM,
                'description': 'Messages d\'erreur hardcodés'
            },
            
            'log_messages': {
                'patterns': [
                    r'log\..*\(["\']([^"\']{15,})["\']',
                    r'logger\..*\(["\']([^"\']{15,})["\']',
                    r'Printf\(["\']([^"\']{15,})["\']',
                ],
                'severity': Severity.MEDIUM,
                'description': 'Messages de log hardcodés'
            },
            
            # LOW - JSON et données structurées
            'json_fields': {
                'patterns': [
                    r'["\']([a-zA-Z_][a-zA-Z0-9_]*)["\']:\s*["\']',
                    r'json:"([^"]*)"',
                    r'bson:"([^"]*)"',
                ],
                'severity': Severity.LOW,
                'description': 'Champs JSON hardcodés'
            },
            
            'template_names': {
                'patterns': [
                    r'\.ParseFiles\(["\']([^"\']*\.html)["\']',
                    r'\.ParseGlob\(["\']([^"\']*)["\']',
                    r'template.*["\']([^"\']*\.html)["\']',
                ],
                'severity': Severity.LOW,
                'description': 'Noms de templates hardcodés'
            },
            
            'file_extensions': {
                'patterns': [
                    r'["\'](\.[a-zA-Z0-9]{2,4})["\']',
                    r'filepath\.Ext.*["\'](\.[^"\']*)["\']',
                ],
                'severity': Severity.LOW,
                'description': 'Extensions de fichiers hardcodées'
            }
        }
        
        # Exclusions pour éviter les faux positifs
        self.exclusions = [
            # Imports
            r'^import\s+',
            r'^\s*["\'].*["\']$',  # Lignes avec seulement une string
            
            # Commentaires
            r'^\s*//',
            r'^\s*/\*',
            r'\*/\s*$',
            
            # Struct tags
            r'`[^`]*`',
            
            # Constantes déjà définies
            r'const\s+\w+\s*=',
            
            # Tests
            r'func\s+Test',
            r'func\s+Benchmark',
            
            # Packages standards
            r'"(fmt|os|io|log|net|http|json|time|strings|strconv)"',
            
            # Valeurs techniques courantes
            r'"(utf-8|UTF-8|ascii|ASCII)"',
            r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
            
            # Protocoles standards
            r'"(http|https|ftp|smtp|tcp|udp)"',
        ]

    def should_exclude_file(self, file_path: str) -> bool:
        """Détermine si un fichier doit être exclu de l'analyse."""
        # Vérifier les exclusions
        for exclusion in self.excluded_files:
            if exclusion in file_path:
                return True
        
        # Vérifier l'extension
        _, ext = os.path.splitext(file_path)
        return ext not in self.included_extensions

    def should_exclude_line(self, line: str) -> bool:
        """Détermine si une ligne doit être exclue de l'analyse."""
        for exclusion_pattern in self.exclusions:
            if re.search(exclusion_pattern, line.strip(), re.IGNORECASE):
                return True
        return False

    def extract_context(self, lines: List[str], line_idx: int, context_lines: int = 2) -> str:
        """Extrait le contexte autour d'une ligne."""
        start = max(0, line_idx - context_lines)
        end = min(len(lines), line_idx + context_lines + 1)
        
        context = []
        for i in range(start, end):
            prefix = ">>> " if i == line_idx else "    "
            context.append(f"{prefix}{lines[i].strip()}")
        
        return "\\n".join(context)

    def analyze_file(self, file_path: str) -> List[Violation]:
        """Analyse un fichier pour détecter les violations de hardcoding."""
        violations = []
        
        try:
            with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                lines = f.readlines()
            
            for line_idx, line in enumerate(lines):
                line_num = line_idx + 1
                
                # Exclure certaines lignes
                if self.should_exclude_line(line):
                    continue
                
                # Analyser chaque catégorie de patterns
                for category, config in self.patterns.items():
                    for pattern in config['patterns']:
                        matches = re.finditer(pattern, line, re.IGNORECASE)
                        
                        for match in matches:
                            value = match.group(1) if match.groups() else match.group(0)
                            
                            # Éviter les valeurs trop courtes ou communes
                            if len(value) < 2:
                                continue
                            
                            # Éviter certaines valeurs communes
                            common_values = {'', 'ok', 'id', 'url', 'api', 'get', 'post'}
                            if value.lower() in common_values:
                                continue
                            
                            violation = Violation(
                                file_path=os.path.relpath(file_path, self.root_path),
                                line=line_num,
                                category=category,
                                value=value,
                                context=self.extract_context(lines, line_idx),
                                severity=config['severity'],
                                description=config['description']
                            )
                            
                            violations.append(violation)
                            
        except Exception as e:
            print(f"Erreur lors de l'analyse de {file_path}: {e}")
        
        return violations

    def scan_directory(self) -> None:
        """Scan récursif du répertoire."""
        print("🔍 Début du scan post-DELTA...")
        
        total_files = 0
        analyzed_files = 0
        
        for root, dirs, files in os.walk(self.root_path):
            # Exclure certains dossiers
            dirs[:] = [d for d in dirs if not any(exc in os.path.join(root, d) for exc in self.excluded_files)]
            
            for file in files:
                file_path = os.path.join(root, file)
                total_files += 1
                
                if self.should_exclude_file(file_path):
                    continue
                
                analyzed_files += 1
                if analyzed_files % 10 == 0:
                    print(f"   Analysé {analyzed_files} fichiers...")
                
                file_violations = self.analyze_file(file_path)
                self.violations.extend(file_violations)
                
                # Mise à jour des statistiques
                rel_path = os.path.relpath(file_path, self.root_path)
                self.file_stats[rel_path] = len(file_violations)
                
                for violation in file_violations:
                    self.category_stats[violation.category] += 1
                    self.severity_stats[violation.severity.value] += 1
        
        print(f"✅ Scan terminé: {analyzed_files}/{total_files} fichiers analysés")
        print(f"📊 Total violations détectées: {len(self.violations)}")

    def generate_comparison_metrics(self) -> Dict[str, Any]:
        """Génère les métriques de comparaison avec l'analyse initiale."""
        initial_violations = 4582  # Référence initiale
        current_violations = len(self.violations)
        
        reduction = initial_violations - current_violations
        reduction_percent = (reduction / initial_violations) * 100 if initial_violations > 0 else 0
        
        return {
            "initial_violations": initial_violations,
            "current_violations": current_violations,
            "violations_eliminated": reduction,
            "reduction_percentage": round(reduction_percent, 2),
            "remaining_work": current_violations,
            "completion_percentage": round(reduction_percent, 2),
            "progress_status": "EXCELLENT" if reduction_percent > 80 else 
                            "GOOD" if reduction_percent > 60 else
                            "IN_PROGRESS" if reduction_percent > 40 else
                            "NEEDS_WORK"
        }

    def generate_priority_plan(self) -> Dict[str, Any]:
        """Génère un plan d'action prioritisé."""
        # Top 10 fichiers avec le plus de violations
        top_files = sorted(
            [(path, count) for path, count in self.file_stats.items() if count > 0],
            key=lambda x: x[1],
            reverse=True
        )[:10]
        
        # Violations par sévérité
        critical_violations = [v for v in self.violations if v.severity == Severity.CRITICAL]
        high_violations = [v for v in self.violations if v.severity == Severity.HIGH]
        
        # Catégories prioritaires
        priority_categories = sorted(
            self.category_stats.items(),
            key=lambda x: x[1],
            reverse=True
        )[:5]
        
        return {
            "immediate_actions": {
                "critical_violations": len(critical_violations),
                "high_priority_violations": len(high_violations),
                "top_files_to_fix": top_files[:5],
                "priority_categories": priority_categories
            },
            "recommended_phases": [
                {
                    "phase": "PHASE 1 - CRITIQUE",
                    "focus": "Éliminer toutes les violations CRITICAL",
                    "target_violations": len(critical_violations),
                    "estimated_effort": "2-4 heures"
                },
                {
                    "phase": "PHASE 2 - HAUTE PRIORITÉ", 
                    "focus": "Traiter les violations HIGH",
                    "target_violations": len(high_violations),
                    "estimated_effort": "4-6 heures"
                },
                {
                    "phase": "PHASE 3 - NETTOYAGE",
                    "focus": "Optimiser les violations MEDIUM/LOW",
                    "target_violations": len([v for v in self.violations if v.severity in [Severity.MEDIUM, Severity.LOW]]),
                    "estimated_effort": "6-8 heures"
                }
            ]
        }

    def generate_report(self) -> Dict[str, Any]:
        """Génère le rapport complet d'analyse."""
        comparison_metrics = self.generate_comparison_metrics()
        priority_plan = self.generate_priority_plan()
        
        # Transformer les violations pour JSON
        violations_data = []
        for v in self.violations:
            violations_data.append({
                "file": v.file_path,
                "line": v.line,
                "category": v.category,
                "value": v.value,
                "context": v.context,
                "severity": v.severity.value,
                "description": v.description
            })
        
        report = {
            "analysis_info": {
                "timestamp": datetime.now().isoformat(),
                "analyzer_version": "POST_DELTA_v1.0",
                "mission": "Évaluation impact corrections DELTA",
                "scan_path": self.root_path
            },
            
            "summary": {
                "total_violations": len(self.violations),
                "files_affected": len([f for f, c in self.file_stats.items() if c > 0]),
                "categories_found": len(self.category_stats),
                "severity_breakdown": dict(self.severity_stats)
            },
            
            "comparison_with_initial": comparison_metrics,
            "priority_action_plan": priority_plan,
            
            "detailed_stats": {
                "violations_by_category": dict(self.category_stats),
                "violations_by_file": dict(self.file_stats),
                "top_violating_files": sorted(
                    [(path, count) for path, count in self.file_stats.items() if count > 0],
                    key=lambda x: x[1],
                    reverse=True
                )[:20]
            },
            
            "violations": violations_data,
            
            "recommendations": {
                "next_steps": [
                    "Corriger immédiatement les violations CRITICAL",
                    "Créer des constantes pour les violations HIGH",
                    "Planifier la refactorisation des violations MEDIUM",
                    "Documenter les violations LOW acceptables"
                ],
                "tools_needed": [
                    "Scripts d'automatisation pour les remplacements",
                    "Linters configurés pour éviter les régressions",
                    "Tests pour valider les corrections"
                ]
            }
        }
        
        return report

def main():
    """Point d'entrée principal."""
    print("🔥 FIRE SALAMANDER - Post-DELTA Hardcoding Analyzer")
    print("=" * 55)
    
    # Chemin du projet
    root_path = "/Users/jeromegonzalez/claude-code/fire-salamander"
    
    if not os.path.exists(root_path):
        print(f"❌ Chemin non trouvé: {root_path}")
        return
    
    # Créer l'analyseur
    analyzer = PostDeltaHardcodingAnalyzer(root_path)
    
    # Scanner le projet
    analyzer.scan_directory()
    
    # Générer le rapport
    print("📋 Génération du rapport d'analyse...")
    report = analyzer.generate_report()
    
    # Sauvegarder le rapport
    output_file = os.path.join(root_path, "post_delta_analysis.json")
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(report, f, indent=2, ensure_ascii=False)
    
    print(f"✅ Rapport sauvegardé: {output_file}")
    
    # Afficher le résumé
    print("\\n📊 RÉSUMÉ DE L'ANALYSE POST-DELTA")
    print("=" * 40)
    print(f"Violations totales: {report['summary']['total_violations']}")
    print(f"Fichiers affectés: {report['summary']['files_affected']}")
    
    comparison = report['comparison_with_initial']
    print(f"\\n🎯 COMPARAISON AVEC L'ANALYSE INITIALE")
    print(f"Violations initiales: {comparison['initial_violations']}")
    print(f"Violations actuelles: {comparison['current_violations']}")
    print(f"Violations éliminées: {comparison['violations_eliminated']}")
    print(f"Réduction: {comparison['reduction_percentage']}%")
    print(f"Statut: {comparison['progress_status']}")
    
    print(f"\\n⚠️  RÉPARTITION PAR SÉVÉRITÉ")
    for severity, count in report['summary']['severity_breakdown'].items():
        print(f"{severity}: {count}")
    
    print(f"\\n🎯 TOP 5 FICHIERS À CORRIGER")
    for file_path, count in report['detailed_stats']['top_violating_files'][:5]:
        print(f"{file_path}: {count} violations")
    
    print(f"\\n📋 PROCHAINES ACTIONS RECOMMANDÉES:")
    for action in report['recommendations']['next_steps']:
        print(f"• {action}")

if __name__ == "__main__":
    main()