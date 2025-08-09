#!/usr/bin/env python3
"""
🤖 CHARLIE-2 SMART ELIMINATOR
Elimination industrielle automatisée des violations hardcoding dans orchestrator.go
"""

import re
import os
import json
from pathlib import Path

def create_orchestrator_constants_mapping():
    """Crée le mapping complet des constantes pour CHARLIE-2"""
    return {
        # JSON Field Names
        '"status"': 'constants.OrchestratorJSONFieldStatus',
        '"message"': 'constants.OrchestratorJSONFieldMessage', 
        '"timestamp"': 'constants.OrchestratorJSONFieldTimestamp',
        '"phase"': 'constants.OrchestratorJSONFieldPhase',
        '"duration"': 'constants.OrchestratorJSONFieldDuration',
        '"error"': 'constants.OrchestratorJSONFieldError',
        '"result"': 'constants.OrchestratorJSONFieldResult',
        '"data"': 'constants.OrchestratorJSONFieldData',
        '"config"': 'constants.OrchestratorJSONFieldConfig',
        '"agents"': 'constants.OrchestratorJSONFieldAgents',
        '"tasks"': 'constants.OrchestratorJSONFieldTasks',
        '"logs"': 'constants.OrchestratorJSONFieldLogs',
        '"metrics"': 'constants.OrchestratorJSONFieldMetrics',
        '"progress"': 'constants.OrchestratorJSONFieldProgress',
        '"state"': 'constants.OrchestratorJSONFieldState',
        
        # Agent Names
        '"orchestrator"': 'constants.OrchestratorAgentNameOrchestrator',
        '"crawler"': 'constants.OrchestratorAgentNameCrawler',
        '"seo"': 'constants.OrchestratorAgentNameSEO',
        '"qa"': 'constants.OrchestratorAgentNameQA',
        '"security"': 'constants.OrchestratorAgentNameSecurity',
        '"performance"': 'constants.OrchestratorAgentNamePerformance',
        '"data-integrity"': 'constants.OrchestratorAgentNameDataIntegrity',
        '"tag-analyzer"': 'constants.OrchestratorAgentNameTagAnalyzer',
        
        # Phase Names
        '"initialization"': 'constants.OrchestratorPhaseInitialization',
        '"crawling"': 'constants.OrchestratorPhaseCrawling',
        '"analysis"': 'constants.OrchestratorPhaseAnalysis',
        '"scoring"': 'constants.OrchestratorPhaseScoring',
        '"reporting"': 'constants.OrchestratorPhaseReporting',
        '"cleanup"': 'constants.OrchestratorPhaseCleanup',
        '"validation"': 'constants.OrchestratorPhaseValidation',
        '"testing"': 'constants.OrchestratorPhaseTesting',
        '"security-scan"': 'constants.OrchestratorPhaseSecurityScan',
        
        # Status Values
        '"pending"': 'constants.OrchestratorStatusPending',
        '"running"': 'constants.OrchestratorStatusRunning',
        '"completed"': 'constants.OrchestratorStatusCompleted',
        '"failed"': 'constants.OrchestratorStatusFailed',
        '"cancelled"': 'constants.OrchestratorStatusCancelled',
        '"success"': 'constants.OrchestratorStatusSuccess',
        '"error"': 'constants.OrchestratorStatusError',
        '"warning"': 'constants.OrchestratorStatusWarning',
        '"info"': 'constants.OrchestratorStatusInfo',
        
        # Config Keys
        '"max_workers"': 'constants.OrchestratorConfigMaxWorkers',
        '"timeout"': 'constants.OrchestratorConfigTimeout',
        '"retry_attempts"': 'constants.OrchestratorConfigRetryAttempts',
        '"batch_size"': 'constants.OrchestratorConfigBatchSize',
        '"concurrent_agents"': 'constants.OrchestratorConfigConcurrentAgents',
        '"log_level"': 'constants.OrchestratorConfigLogLevel',
        '"debug_mode"': 'constants.OrchestratorConfigDebugMode',
        
        # HTTP Methods
        '"GET"': 'constants.OrchestratorHTTPMethodGet',
        '"POST"': 'constants.OrchestratorHTTPMethodPost',
        '"PUT"': 'constants.OrchestratorHTTPMethodPut',
        '"DELETE"': 'constants.OrchestratorHTTPMethodDelete',
        '"PATCH"': 'constants.OrchestratorHTTPMethodPatch',
        '"HEAD"': 'constants.OrchestratorHTTPMethodHead',
        '"OPTIONS"': 'constants.OrchestratorHTTPMethodOptions',
        
        # Content Types
        '"application/json"': 'constants.OrchestratorContentTypeJSON',
        '"text/html"': 'constants.OrchestratorContentTypeHTML',
        '"text/plain"': 'constants.OrchestratorContentTypePlain',
        '"application/xml"': 'constants.OrchestratorContentTypeXML',
        
        # Error Types
        '"crawling_error"': 'constants.OrchestratorErrorTypeCrawling',
        '"semantic_error"': 'constants.OrchestratorErrorTypeSemantic',
        '"seo_error"': 'constants.OrchestratorErrorTypeSEO',
        '"validation_error"': 'constants.OrchestratorErrorTypeValidation',
        '"timeout_error"': 'constants.OrchestratorErrorTypeTimeout',
        '"config_error"': 'constants.OrchestratorErrorTypeConfiguration',
        '"network_error"': 'constants.OrchestratorErrorTypeNetwork',
        '"internal_error"': 'constants.OrchestratorErrorTypeInternal',
        
        # File Paths
        '"reports/"': 'constants.OrchestratorPathReports',
        '"logs/"': 'constants.OrchestratorPathLogs',
        '"temp/"': 'constants.OrchestratorPathTemp',
        '"cache/"': 'constants.OrchestratorPathCache',
        '"config/"': 'constants.OrchestratorPathConfig',
        
        # URL Schemes
        '"http://"': 'constants.OrchestratorURLSchemeHTTP',
        '"https://"': 'constants.OrchestratorURLSchemeHTTPS',
        '"ws://"': 'constants.OrchestratorURLSchemeWS',
        '"wss://"': 'constants.OrchestratorURLSchemeWSS',
        
        # Time Formats
        '"RFC3339"': 'constants.OrchestratorTimeFormatRFC3339',
        '"2006-01-02T15:04:05Z"': 'constants.OrchestratorTimeFormatISO8601',
        '"2006-01-02 15:04:05"': 'constants.OrchestratorTimeFormatSimple',
        
        # Numeric Constants (specific contexts only)
        '10': 'constants.OrchestratorMaxWorkers',
        '30': 'constants.OrchestratorHealthCheckInterval',
        '60': 'constants.OrchestratorProgressUpdateInterval',
        '100': 'constants.OrchestratorDefaultBatchSize',
        '500': 'constants.OrchestratorTaskChannelSize',
        '1000': 'constants.OrchestratorMaxBatchSize',
        '5000': 'constants.OrchestratorMaxQueueSize',
        '10000': 'constants.OrchestratorChannelBufferLarge',
        
        # Context Operations
        'context.Background()': 'context.Background()',
        'context.TODO()': 'context.TODO()', 
        'context.WithTimeout': 'context.WithTimeout',
        'context.WithCancel': 'context.WithCancel',
        
        # Sync Patterns
        'sync.WaitGroup': 'sync.WaitGroup',
        'sync.Mutex': 'sync.Mutex',
        'sync.RWMutex': 'sync.RWMutex',
        'WaitGroup': 'sync.WaitGroup',
        'Mutex': 'sync.Mutex',
        'RWMutex': 'sync.RWMutex',
    }

def should_replace_in_orchestrator_context(line, hardcoded_value, constant):
    """Détermine si un remplacement doit être effectué basé sur le contexte orchestrator"""
    
    # Skip struct tags completely
    if '`json:' in line:
        return False
        
    # Skip comments
    if line.strip().startswith('//'):
        return False
        
    # Skip imports
    if 'import' in line or 'package' in line:
        return False
    
    # Special handling for numeric values - only in specific contexts
    if hardcoded_value in ['10', '30', '60', '100', '500', '1000', '5000', '10000']:
        # Only replace numbers in specific contexts
        if hardcoded_value == '10' and ('workers' in line.lower() or 'maxWorkers' in line):
            return True
        elif hardcoded_value == '30' and ('health' in line.lower() or 'interval' in line.lower()):
            return True
        elif hardcoded_value == '60' and ('progress' in line.lower() or 'update' in line.lower()):
            return True
        elif hardcoded_value == '100' and ('batch' in line.lower() or 'buffer' in line.lower()):
            return True
        elif hardcoded_value == '500' and ('task' in line.lower() or 'queue' in line.lower()):
            return True
        elif hardcoded_value == '1000' and ('max' in line.lower() or 'size' in line.lower()):
            return True
        elif hardcoded_value == '5000' and 'queue' in line.lower():
            return True
        elif hardcoded_value == '10000' and 'buffer' in line.lower():
            return True
        else:
            return False
    
    # Special handling for sync patterns - don't replace type declarations
    if hardcoded_value in ['WaitGroup', 'Mutex', 'RWMutex']:
        # Only replace in variable declarations and usage, not type definitions
        if 'type ' in line or 'struct {' in line:
            return False
        return True
    
    # Special handling for context operations - don't replace unless it's a usage
    if hardcoded_value.startswith('context.'):
        return '(' in line  # Only replace if it's a function call
    
    # Default: allow replacement for string literals and most patterns
    return True

def eliminate_orchestrator_violations(filepath, constants_mapping):
    """Élimine les violations de hardcoding spécifiques à l'orchestrator"""
    
    print(f"🤖 ÉLIMINATION ORCHESTRATOR: {filepath}")
    
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
                if should_replace_in_orchestrator_context(line, hardcoded_value, constant):
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
    print("🤖 CHARLIE-2 SMART ELIMINATOR - Élimination orchestrator...")
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/integration/orchestrator.go'
    
    # Créer une sauvegarde
    backup_path = f"{filepath}.charlie2_backup"
    if not os.path.exists(backup_path):
        with open(filepath, 'r') as original:
            with open(backup_path, 'w') as backup:
                backup.write(original.read())
        print(f"💾 Sauvegarde créée: {backup_path}")
    
    # Obtenir le mapping des constantes
    constants_mapping = create_orchestrator_constants_mapping()
    print(f"📋 {len(constants_mapping)} mappings de constantes chargés")
    
    # Éliminer les violations
    total_eliminated = eliminate_orchestrator_violations(filepath, constants_mapping)
    
    print(f"\n🎯 CHARLIE-2 TERMINÉ:")
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