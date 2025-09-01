#!/usr/bin/env python3
"""
🤖 CHARLIE-2 SMART ELIMINATOR V2
Elimination industrielle ciblée pour orchestrator.go avec contrôle contextuel
"""

import re
import os
import json

def create_selective_orchestrator_mapping():
    """Crée le mapping sélectif pour orchestrator sans casser les types"""
    return {
        # Agent Names (string literals only)
        '"seo"': 'constants.OrchestratorAgentNameSEO',
        '"crawler"': 'constants.OrchestratorAgentNameCrawler',
        '"performance"': 'constants.OrchestratorAgentNamePerformance',
        '"security"': 'constants.OrchestratorAgentNameSecurity',
        '"qa"': 'constants.OrchestratorAgentNameQA',
        '"semantic"': 'constants.OrchestratorAgentNameSemantic',
        
        # Status Values (string literals only)
        '"pending"': 'constants.OrchestratorStatusPending',
        '"running"': 'constants.OrchestratorStatusRunning',
        '"completed"': 'constants.OrchestratorStatusCompleted',
        '"failed"': 'constants.OrchestratorStatusFailed',
        '"cancelled"': 'constants.OrchestratorStatusCancelled',
        '"success"': 'constants.OrchestratorStatusSuccess',
        '"error"': 'constants.OrchestratorStatusError',
        '"warning"': 'constants.OrchestratorStatusWarning',
        '"info"': 'constants.OrchestratorStatusInfo',
        '"partial"': 'constants.OrchestratorStatusPartial',
        
        # JSON Field Names
        '"status"': 'constants.OrchestratorJSONFieldStatus',
        '"message"': 'constants.OrchestratorJSONFieldMessage',
        '"timestamp"': 'constants.OrchestratorJSONFieldTimestamp',
        '"phase"': 'constants.OrchestratorJSONFieldPhase',
        '"error"': 'constants.OrchestratorJSONFieldError',
        '"result"': 'constants.OrchestratorJSONFieldResult',
        '"data"': 'constants.OrchestratorJSONFieldData',
        
        # Error Types
        '"crawling_error"': 'constants.OrchestratorErrorTypeCrawling',
        '"semantic_error"': 'constants.OrchestratorErrorTypeSemantic',
        '"seo_error"': 'constants.OrchestratorErrorTypeSEO',
        '"validation_error"': 'constants.OrchestratorErrorTypeValidation',
        '"timeout_error"': 'constants.OrchestratorErrorTypeTimeout',
        '"config_error"': 'constants.OrchestratorErrorTypeConfiguration',
        '"network_error"': 'constants.OrchestratorErrorTypeNetwork',
        '"internal_error"': 'constants.OrchestratorErrorTypeInternal',
        
        # Analysis Types
        '"full"': 'constants.OrchestratorAnalysisTypeFull',
        '"semantic"': 'constants.OrchestratorAnalysisTypeSemantic',
        '"quick"': 'constants.OrchestratorAnalysisTypeQuick',
        
        # Insight Types
        '"content_seo_alignment"': 'constants.OrchestratorInsightContentSEOAlignment',
        '"performance_content_mismatch"': 'constants.OrchestratorInsightPerformanceContentMismatch',
        
        # Severity Levels
        '"positive"': 'constants.OrchestratorImpactPositive',
        '"negative"': 'constants.OrchestratorImpactNegative',
        '"neutral"': 'constants.OrchestratorImpactNeutral',
        
        # Categories and scores (string keys)
        '"technical"': 'constants.OrchestratorCategoryTechnical',
        '"content_quality"': 'constants.OrchestratorCategoryContentQuality',
        '"user_experience"': 'constants.OrchestratorCategoryUserExperience',
        '"mobile_friendliness"': 'constants.OrchestratorCategoryMobileFriendliness',
        
        # Common message templates
        '"Variable"': 'constants.OrchestratorTimeVariable',
    }

def get_additional_constants_for_file():
    """Retourne des constantes supplémentaires à ajouter au fichier constants"""
    return '''
// Additional Orchestrator Constants for CHARLIE-2

// Analysis Types
const (
	OrchestratorAnalysisTypeFull     = "full"
	OrchestratorAnalysisTypeSemantic = "semantic"
	OrchestratorAnalysisTypeQuick    = "quick"
)

// Additional Status Values
const (
	OrchestratorStatusPartial = "partial"
)

// Insight Types
const (
	OrchestratorInsightContentSEOAlignment           = "content_seo_alignment"
	OrchestratorInsightPerformanceContentMismatch   = "performance_content_mismatch"
)

// Impact Types
const (
	OrchestratorImpactPositive = "positive"
	OrchestratorImpactNegative = "negative"
	OrchestratorImpactNeutral  = "neutral"
)

// Category Names for Scores
const (
	OrchestratorCategoryTechnical         = "technical"
	OrchestratorCategoryContentQuality    = "content_quality"
	OrchestratorCategoryUserExperience    = "user_experience"
	OrchestratorCategoryMobileFriendliness = "mobile_friendliness"
)

// Agent Names (additional)
const (
	OrchestratorAgentNameSemantic = "semantic"
)

// Time Constants
const (
	OrchestratorTimeVariable = "Variable"
)
'''

def should_replace_in_orchestrator(line, hardcoded_value, constant):
    """Détermine si un remplacement doit être effectué"""
    
    # Skip struct tags
    if '`json:' in line:
        return False
        
    # Skip comments
    if line.strip().startswith('//'):
        return False
        
    # Skip imports and package
    if 'import' in line or 'package' in line:
        return False
    
    # Skip type declarations - ne pas remplacer dans les const/type blocks
    if re.match(r'^\s*(const|type|var)\s*\(?\s*$', line.strip()):
        return False
    
    # Skip const assignments like TaskStatus = "pending"
    if re.match(r'^\s*\w+\s+(TaskStatus|AnalysisStatus|AnalysisType)\s*=', line.strip()):
        return False
    
    # Only replace string literals, not bare identifiers in type definitions
    return True

def eliminate_orchestrator_violations_v2(filepath, constants_mapping):
    """Élimine sélectivement les violations orchestrator"""
    
    print(f"🤖 ÉLIMINATION ORCHESTRATOR V2: {filepath}")
    
    # Lire le fichier
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements = 0
    
    # Split into lines for context-aware processing
    lines = content.split('\n')
    new_lines = []
    
    for line_num, line in enumerate(lines):
        new_line = line
        
        for hardcoded_value, constant in constants_mapping.items():
            if hardcoded_value in line:
                if should_replace_in_orchestrator(line, hardcoded_value, constant):
                    old_line = new_line
                    new_line = new_line.replace(hardcoded_value, constant)
                    if new_line != old_line:
                        replacements += 1
                        print(f"  ✅ Line {line_num + 1}: {hardcoded_value} → {constant}")
        
        new_lines.append(new_line)
    
    content = '\n'.join(new_lines)
    
    # Sauvegarder si des modifications ont été effectuées
    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"✅ ÉLIMINÉ: {replacements} violations dans {filepath}")
        return replacements
    else:
        print(f"ℹ️ Aucune modification nécessaire dans {filepath}")
        return 0

def main():
    print("🤖 CHARLIE-2 SMART ELIMINATOR V2 - Élimination sélective orchestrator...")
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/orchestrator.go'
    constants_file = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/constants/orchestrator_constants.go'
    
    # Créer une sauvegarde V2
    backup_path = f"{filepath}.charlie2_v2_backup"
    if not os.path.exists(backup_path):
        with open(filepath, 'r') as original:
            with open(backup_path, 'w') as backup:
                backup.write(original.read())
        print(f"💾 Sauvegarde V2 créée: {backup_path}")
    
    # Ajouter les constantes manquantes au fichier constants
    with open(constants_file, 'a') as f:
        f.write(get_additional_constants_for_file())
    print(f"📋 Constantes supplémentaires ajoutées à {constants_file}")
    
    # Obtenir le mapping sélectif
    constants_mapping = create_selective_orchestrator_mapping()
    print(f"📋 {len(constants_mapping)} mappings sélectifs chargés")
    
    # Éliminer les violations sélectivement
    total_eliminated = eliminate_orchestrator_violations_v2(filepath, constants_mapping)
    
    print(f"\n🎯 CHARLIE-2 V2 TERMINÉ:")
    print(f"✅ Total violations éliminées: {total_eliminated}")
    print(f"📁 Fichier traité: {filepath}")
    print(f"💾 Sauvegarde: {backup_path}")
    
    # Tester la compilation
    print("\n🔨 Test de compilation...")
    import subprocess
    try:
        result = subprocess.run(['go', 'build', './internal/integration/...'], 
                              capture_output=True, text=True, cwd='/Users/jeromegonzalez/claude-code/fire-salamander')
        if result.returncode == 0:
            print("✅ Compilation réussie!")
        else:
            print("⚠️ Erreurs de compilation:")
            print(result.stderr)
    except Exception as e:
        print(f"⚠️ Impossible de tester la compilation: {e}")
    
    return total_eliminated

if __name__ == "__main__":
    main()