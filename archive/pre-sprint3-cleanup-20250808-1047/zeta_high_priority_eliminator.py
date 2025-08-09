#!/usr/bin/env python3
"""
🔥 FIRE SALAMANDER - ZETA HIGH PRIORITY ELIMINATOR
Mission: Éliminer les 4 violations HIGH priorité identifiées dans post_delta_analysis.json
Phase: ZETA - Finalisation optimisations production
"""

import os
import re
import json
import sys
from datetime import datetime
from pathlib import Path

class ZetaHighPriorityEliminator:
    """Élimine les 4 violations HIGH priorité de configuration serveur hardcodée"""
    
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.violations_file = self.project_root / "post_delta_analysis.json"
        self.constants_dir = self.project_root / "internal" / "constants"
        self.server_constants_file = self.constants_dir / "server_constants.go"
        
        # Les 4 violations HIGH identifiées
        self.high_violations = [
            {
                "file": "cmd/fire-salamander/main.go",
                "line": 169,
                "pattern": 'Host: "localhost"',
                "replacement": 'Host: constants.ServerDefaultHost',
                "description": "Configuration serveur hardcodée - main.go"
            },
            {
                "file": "internal/config/config.go", 
                "line": 79,
                "pattern": 'cfg.Server.Host = "localhost"',
                "replacement": 'cfg.Server.Host = constants.ServerDefaultHost',
                "description": "Configuration serveur hardcodée - config.go"
            },
            {
                "file": "internal/seo/analyzer.go",
                "line": 449,
                "pattern": '"localhost"',
                "replacement": 'constants.ServerDefaultHost',
                "description": "Configuration serveur hardcodée - analyzer.go"
            },
            {
                "file": "internal/debug/phase_tests.go",
                "line": 279,
                "pattern": 'constants.HTTPPrefix+"localhost:%d"',
                "replacement": 'constants.HTTPPrefix+constants.ServerDefaultHost+":%d"',
                "description": "Configuration serveur hardcodée - phase_tests.go"
            }
        ]
        
        print(f"🎯 ZETA HIGH PRIORITY ELIMINATOR - Démarrage")
        print(f"📁 Répertoire projet: {self.project_root}")
        print(f"🔍 {len(self.high_violations)} violations HIGH à corriger")
    
    def verify_constants_available(self) -> bool:
        """Vérifie que les constantes nécessaires sont disponibles"""
        print(f"\n🔧 Vérification des constantes dans {self.server_constants_file}")
        
        if not self.server_constants_file.exists():
            print(f"❌ Fichier de constantes serveur non trouvé: {self.server_constants_file}")
            return False
            
        with open(self.server_constants_file, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # Vérifier que ServerDefaultHost existe
        if "ServerDefaultHost" not in content:
            print("❌ Constante ServerDefaultHost non trouvée")
            return False
            
        print("✅ Constantes serveur disponibles")
        return True
    
    def fix_violation(self, violation: dict) -> bool:
        """Corrige une violation spécifique"""
        file_path = self.project_root / violation["file"]
        
        if not file_path.exists():
            print(f"❌ Fichier non trouvé: {file_path}")
            return False
            
        print(f"🔧 Correction de {violation['file']}...")
        
        # Lire le contenu du fichier
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # Effectuer le remplacement
        if violation["pattern"] in content:
            # Vérifier que l'import constants est présent
            if not self.ensure_constants_import(file_path, content):
                return False
                
            # Effectuer le remplacement
            new_content = content.replace(violation["pattern"], violation["replacement"])
            
            # Sauvegarder
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
                
            print(f"✅ {violation['description']} - Corrigé")
            return True
        else:
            print(f"⚠️  Pattern non trouvé dans {violation['file']}: {violation['pattern']}")
            return False
    
    def ensure_constants_import(self, file_path: Path, content: str) -> bool:
        """S'assure que l'import constants est présent"""
        # Vérifier si l'import existe déjà
        if 'import (' in content and '"fire-salamander/internal/constants"' in content:
            return True
        elif 'import "fire-salamander/internal/constants"' in content:
            return True
            
        # Ajouter l'import
        lines = content.split('\n')
        import_added = False
        
        for i, line in enumerate(lines):
            # Chercher une section d'imports existante
            if line.strip().startswith('import ('):
                # Ajouter dans la section d'imports existante
                j = i + 1
                while j < len(lines) and not lines[j].strip().startswith(')'):
                    j += 1
                if j < len(lines):
                    lines.insert(j, '\t"fire-salamander/internal/constants"')
                    import_added = True
                    break
            elif line.strip().startswith('import "') and not import_added:
                # Remplacer l'import simple par un import multiple
                lines[i] = 'import ('
                lines.insert(i + 1, '\t' + line.strip().replace('import ', ''))
                lines.insert(i + 2, '\t"fire-salamander/internal/constants"')
                lines.insert(i + 3, ')')
                import_added = True
                break
        
        if not import_added:
            # Ajouter après la déclaration package
            for i, line in enumerate(lines):
                if line.strip().startswith('package '):
                    lines.insert(i + 2, 'import "fire-salamander/internal/constants"')
                    import_added = True
                    break
        
        if import_added:
            new_content = '\n'.join(lines)
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            print(f"📦 Import constants ajouté à {file_path.name}")
            
        return import_added
    
    def test_compilation(self) -> bool:
        """Teste que le code compile après les modifications"""
        print(f"\n🔨 Test de compilation...")
        
        try:
            os.chdir(self.project_root)
            result = os.system("go build ./...")
            
            if result == 0:
                print("✅ Compilation réussie")
                return True
            else:
                print("❌ Erreur de compilation")
                return False
        except Exception as e:
            print(f"❌ Erreur lors du test de compilation: {e}")
            return False
    
    def generate_report(self, results: dict) -> None:
        """Génère un rapport des corrections effectuées"""
        report = {
            "mission": "ZETA HIGH PRIORITY ELIMINATOR",
            "timestamp": datetime.now().isoformat(),
            "summary": {
                "total_violations": len(self.high_violations),
                "fixed_violations": sum(1 for r in results.values() if r['success']),
                "failed_violations": sum(1 for r in results.values() if not r['success'])
            },
            "details": []
        }
        
        for key, result in results.items():
            violation = result['violation']
            success = result['success']
            report["details"].append({
                "file": violation["file"],
                "line": violation["line"],
                "description": violation["description"],
                "status": "FIXED" if success else "FAILED",
                "pattern": violation["pattern"],
                "replacement": violation["replacement"]
            })
        
        report_file = self.project_root / "zeta_high_priority_report.json"
        with open(report_file, 'w', encoding='utf-8') as f:
            json.dump(report, f, indent=2, ensure_ascii=False)
            
        print(f"\n📋 Rapport généré: {report_file}")
        
        # Affichage résumé
        print(f"\n🎯 RÉSUMÉ MISSION ZETA:")
        print(f"✅ Violations corrigées: {report['summary']['fixed_violations']}")
        print(f"❌ Échecs: {report['summary']['failed_violations']}")
        print(f"📊 Taux de succès: {(report['summary']['fixed_violations'] / report['summary']['total_violations']) * 100:.1f}%")
    
    def run(self) -> bool:
        """Lance l'élimination des violations HIGH"""
        print(f"\n🚀 Démarrage mission ZETA HIGH PRIORITY ELIMINATOR")
        
        # Vérification des prérequis
        if not self.verify_constants_available():
            print("❌ Prérequis non satisfaits - Arrêt")
            return False
        
        # Correction des violations
        results = {}
        for i, violation in enumerate(self.high_violations):
            results[f"violation_{i}_{violation['file'].split('/')[-1]}"] = {
                'violation': violation,
                'success': self.fix_violation(violation)
            }
        
        # Test de compilation
        compilation_success = self.test_compilation()
        
        # Génération du rapport
        self.generate_report(results)
        
        # Résultat final
        all_fixed = all(r['success'] for r in results.values())
        if all_fixed and compilation_success:
            print(f"\n🎉 MISSION ZETA ACCOMPLIE - Toutes les violations HIGH éliminées avec succès!")
            return True
        else:
            print(f"\n⚠️  MISSION ZETA PARTIELLEMENT ACCOMPLIE - Vérifier le rapport")
            return False

def main():
    """Point d'entrée principal"""
    if len(sys.argv) != 2:
        print("Usage: python3 zeta_high_priority_eliminator.py <project_root>")
        sys.exit(1)
        
    project_root = sys.argv[1]
    if not os.path.exists(project_root):
        print(f"❌ Répertoire projet non trouvé: {project_root}")
        sys.exit(1)
        
    eliminator = ZetaHighPriorityEliminator(project_root)
    success = eliminator.run()
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()