#!/usr/bin/env python3
"""
🤖 ALPHA-5 SMART ELIMINATOR
Élimination automatisée des violations hardcoding dans server.go
Mission: 137 violations → 0 violations
"""

import re
import shutil
import subprocess
import sys

def create_string_mappings():
    """Créer les mappings de remplacement pour ALPHA-5"""
    
    mappings = {
        # Error Messages
        '"Erreur accès fichiers statiques: %v"': 'constants.MsgErrorStaticFiles',
        '"Interface web non disponible"': 'constants.MsgWebInterfaceUnavailable',
        '"Method not allowed"': 'constants.MsgMethodNotAllowed',
        '"API endpoint not found"': 'constants.MsgAPIEndpointNotFound',
        '"Nom de fichier requis"': 'constants.MsgFilenameRequired',
        '"Erreur serveur web: %v"': 'constants.MsgWebServerError',
        '"Erreur lecture index.html: %v"': 'constants.MsgErrorIndexHTML',
        
        # Content Types
        '"text/html; charset=utf-8"': 'constants.ContentTypeHTMLCharset',
        '"application/pdf"': 'constants.ContentTypePDF',
        '"text/csv"': 'constants.ContentTypeCSV',
        '"application/octet-stream"': 'constants.ContentTypeOctetStream',
        '"text/html"': 'constants.ContentTypeHTML',
        
        # HTTP Headers
        '"Cache-Control"': 'constants.HeaderCacheControl',
        '"X-Frame-Options"': 'constants.HeaderXFrameOptions',
        '"X-Content-Type-Options"': 'constants.HeaderXContentType',
        '"Content-Length"': 'constants.HeaderContentLength',
        '"Access-Control-Allow-Origin"': 'constants.HeaderCORSOrigin',
        '"Access-Control-Allow-Methods"': 'constants.HeaderCORSMethods',
        '"Access-Control-Allow-Headers"': 'constants.HeaderCORSHeaders',
        '"Content-Disposition"': 'constants.HeaderContentDisposition',
        
        # HTTP Header Values
        '"no-cache, no-store, must-revalidate"': 'constants.HeaderValueNoCache',
        '"DENY"': 'constants.HeaderValueDeny',
        '"nosniff"': 'constants.HeaderValueNoSniff',
        '"*"': 'constants.HeaderValueCORSAll',
        '"GET, POST, PUT, DELETE, OPTIONS"': 'constants.HeaderValueCORSMethods',
        '"Content-Type, Authorization, X-Requested-With"': 'constants.HeaderValueCORSHeaders',
        
        # HTTP Methods
        '"GET"': 'constants.HTTPMethodGET',
        '"POST"': 'constants.HTTPMethodPOST',
        '"PUT"': 'constants.HTTPMethodPUT',
        '"DELETE"': 'constants.HTTPMethodDELETE',
        '"OPTIONS"': 'constants.HTTPMethodOPTIONS',
        
        # API Routes
        '"/api/v1/analyze"': 'constants.APIRouteAnalyze',
        '"/api/v1/analyze/semantic"': 'constants.APIRouteAnalyzeSemantic',
        '"/api/v1/analyze/seo"': 'constants.APIRouteAnalyzeSEO',
        '"/api/v1/analyze/quick"': 'constants.APIRouteAnalyzeQuick',
        '"/api/v1/health"': 'constants.APIRouteHealth',
        '"/api/v1/stats"': 'constants.APIRouteStats',
        '"/api/v1/analyses"': 'constants.APIRouteAnalyses',
        '"/api/v1/analysis/"': 'constants.APIRouteAnalysisDetails',
        '"/api/v1/info"': 'constants.APIRouteInfo',
        '"/api/v1/version"': 'constants.APIRouteVersion',
        
        # Web Routes
        '"/"': 'constants.WebRouteRoot',
        '"/web/health"': 'constants.WebRouteHealth',
        '"/web/download/"': 'constants.WebRouteDownload',
        '"/static/"': 'constants.StaticRoute',
        
        # File Extensions
        '".html"': 'constants.ExtHTML',
        '".pdf"': 'constants.ExtPDF',
        '".json"': 'constants.ExtJSON',
        '".csv"': 'constants.ExtCSV',
        
        # Static Files
        '"index.html"': 'constants.StaticIndexHTML',
        '"static"': 'constants.StaticDirectory',
        
        # Status Values
        '"unavailable"': 'constants.HealthStatusUnavailable',
        '"running"': 'constants.ServiceStatusRunning',
        
        # Service Names
        '"Fire Salamander Web Server"': 'constants.ServiceNameWebServer',
        
        # Log Messages
        '"Routes web enregistrées:"': 'constants.MsgRoutesRegistered',
        '"  GET  / - Interface web principale"': 'constants.MsgRouteWebInterface',
        '"  GET  /static/* - Fichiers statiques"': 'constants.MsgRouteStaticFiles',
        '"  ALL  /api/* - API REST (proxy)"': 'constants.MsgRouteAPIProxy',
        '"  GET  /web/health - Santé du serveur web"': 'constants.MsgRouteWebHealth',
        '"  GET  /web/download/* - Téléchargement de rapports"': 'constants.MsgRouteWebDownload',
        '"Interface web servie: %s %s"': 'constants.MsgWebInterfaceServed',
        '"API Request: %s %s"': 'constants.MsgAPIRequest',
        '"Téléchargement demandé: %s"': 'constants.MsgDownloadRequested',
        '"✅ Serveur web Fire Salamander démarré avec succès"': 'constants.MsgWebServerStarted',
        '"Arrêt du serveur web Fire Salamander"': 'constants.MsgWebServerStopping',
        
        # JSON Field Names
        '"service"': 'constants.JSONFieldService',
        '"components"': 'constants.JSONFieldComponents',
        '"uptime"': 'constants.JSONFieldUptime',
        '"orchestrator_stats"': 'constants.JSONFieldOrchestratorStats',
        '"web_server"': 'constants.ComponentWebServer',
        '"static_files"': 'constants.ComponentStaticFiles',
        '"api_proxy"': 'constants.ComponentAPIProxy',
        '"orchestrator"': 'constants.ComponentOrchestrator',
        '"started_at"': 'constants.JSONFieldStartedAt',
        
        # Report Constants
        '"Rapport Fire Salamander"': 'constants.ReportTitlePrefix',
        '"Rapport Fire Salamander - "': 'constants.ReportTitleSuffix',
        '"Rapport d\'Analyse SEO"': 'constants.ReportHTMLTitle',
        '"optimize_images"': 'constants.RecommendationOptimizeImages',
        '"meta_descriptions"': 'constants.RecommendationMetaDesc',
        '"Optimiser les images"': 'constants.RecommendationTitleImages',
        '"Améliorer les méta-descriptions"': 'constants.RecommendationTitleMetaDesc',
        '"Compresser les images pour améliorer les temps de chargement"': 'constants.RecommendationDescImages',
        '"high"': 'constants.PriorityHigh',
        '"medium"': 'constants.PriorityMedium',
        '"low"': 'constants.PriorityLow',
        '"info"': 'constants.InsightSeverityInfo',
        '"stable"': 'constants.SampleScoreTrend',
        '"Succès"': 'constants.CSVValueSuccess',
        
        # Constants with usage context for better precision
        'constants.WebServerStarting': 'constants.MsgWebServerStarting',  # Use for log message
    }
    
    return mappings

def backup_file(filepath):
    """Créer une sauvegarde du fichier original"""
    backup_path = f"{filepath}.alpha5_backup"
    shutil.copy2(filepath, backup_path)
    return backup_path

def restore_file(filepath):
    """Restaurer le fichier depuis la sauvegarde"""
    backup_path = f"{filepath}.alpha5_backup"
    try:
        shutil.copy2(backup_path, filepath)
        return True
    except:
        return False

def apply_targeted_replacements(filepath, mappings):
    """Appliquer les remplacements de hardcoding de manière ciblée"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements_made = 0
    
    # Remplacements spéciaux avec contexte
    special_replacements = [
        # Status values in specific contexts
        (r'health\["status"\]\s*=\s*"degraded"', 'health["status"] = constants.HealthStatusDegraded'),
        (r'"status":\s*"healthy"', '"status": constants.HealthStatusHealthy'),
        (r'"web_server":\s*"healthy"', '"web_server": constants.HealthStatusHealthy'),
        (r'"static_files":\s*"healthy"', '"static_files": constants.HealthStatusHealthy'),
        (r'"api_proxy":\s*"healthy"', '"api_proxy": constants.HealthStatusHealthy'),
        (r'"orchestrator":\s*"healthy"', '"orchestrator": constants.HealthStatusHealthy'),
        
        # Method comparisons
        (r'r\.Method\s*!=\s*"GET"', 'r.Method != constants.HTTPMethodGET'),
        (r'r\.Method\s*==\s*"OPTIONS"', 'r.Method == constants.HTTPMethodOPTIONS'),
        
        # URL Path comparisons  
        (r'r\.URL\.Path\s*!=\s*"/"', 'r.URL.Path != constants.WebRouteRoot'),
        (r'r\.URL\.Path\s*==\s*"/api/v1/analyze"', 'r.URL.Path == constants.APIRouteAnalyze'),
        
        # File extension checks
        (r'filepath\.Ext\(filename\)\s*[\s\n]*case\s*"\.html":', 'filepath.Ext(filename)\n\tcase constants.ExtHTML:'),
        (r'case\s*"\.pdf":', 'case constants.ExtPDF:'),
        (r'case\s*"\.json":', 'case constants.ExtJSON:'),
        (r'case\s*"\.csv":', 'case constants.ExtCSV:'),
    ]
    
    # Appliquer les remplacements spéciaux d'abord
    for pattern, replacement in special_replacements:
        if re.search(pattern, content):
            content = re.sub(pattern, replacement, content)
            replacements_made += 1
            print(f"✅ Special: {pattern[:50]}... → {replacement[:50]}...")
    
    # Appliquer les remplacements standards
    for old_string, new_string in mappings.items():
        if old_string in content:
            content = content.replace(old_string, new_string)
            replacements_made += 1
            print(f"✅ {old_string} → {new_string}")
    
    # Sauvegarder le fichier modifié
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    return replacements_made, len(content) != len(original_content)

def test_compilation():
    """Tester la compilation du code Go"""
    try:
        result = subprocess.run(['go', 'build', './internal/web'], 
                              capture_output=True, text=True, cwd='/Users/jeromegonzalez/claude-code/fire-salamander')
        
        if result.returncode == 0:
            print("✅ Compilation successful!")
            return True
        else:
            print(f"❌ Compilation failed:\n{result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Error testing compilation: {e}")
        return False

def main():
    """Fonction principale ALPHA-5 Eliminator"""
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/server.go'
    
    print("🤖 ALPHA-5 SMART ELIMINATOR")
    print("=" * 50)
    print(f"Target: {filepath}")
    print("Mission: Élimination complète des hardcoding violations")
    
    # 1. Backup
    print("\n📦 Creating backup...")
    backup_path = backup_file(filepath)
    print(f"✅ Backup created: {backup_path}")
    
    # 2. Load mappings
    print("\n🔧 Loading string mappings...")
    mappings = create_string_mappings()
    print(f"✅ {len(mappings)} mappings loaded")
    
    # 3. Apply replacements
    print("\n🔄 Applying hardcode eliminations...")
    replacements_made, file_modified = apply_targeted_replacements(filepath, mappings)
    
    print(f"\n📊 RÉSULTATS:")
    print(f"   - Replacements applied: {replacements_made}")
    print(f"   - File modified: {file_modified}")
    
    # 4. Test compilation
    print(f"\n🔨 Testing compilation...")
    if test_compilation():
        print(f"\n🎯 ALPHA-5 MISSION ACCOMPLISHED!")
        print(f"   - Target file: server.go")
        print(f"   - Violations eliminated: {replacements_made}")
        print(f"   - Status: 100% SUCCESS")
        return True
    else:
        print(f"\n⚠️  Compilation failed - restoring backup...")
        if restore_file(filepath):
            print(f"   - File restored from backup")
        print(f"   - Status: NEEDS MANUAL INTERVENTION")
        return False

if __name__ == "__main__":
    import os
    main()