#!/usr/bin/env python3
"""
⚔️ DELTA-6 ELIMINATOR
Éliminateur spécialisé pour crawler_test.go - 49 violations ciblées
"""

import re
import json
import shutil
from pathlib import Path

def create_backup(filepath):
    """Créer un backup avant modification"""
    backup_path = f"{filepath}.delta6_backup"
    shutil.copy2(filepath, backup_path)
    print(f"💾 Backup créé: {backup_path}")
    return backup_path

def load_analysis():
    """Charger l'analyse DELTA-6"""
    with open('delta6_analysis.json', 'r') as f:
        return json.load(f)

def execute_elimination(filepath, analysis):
    """Exécuter l'élimination contextuelle des violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
        original_content = content
    
    # Mappings de remplacement précis pour crawler_test.go
    replacements = {
        # Test URLs
        '"http://www.sitemaps.org/schemas/sitemap/0.9"': 'constants.CrawlerTestURLSitemapSchema',
        '"https://example.com/page1.html"': 'constants.CrawlerTestURLExamplePage1',
        '"https://example.com/page2.html"': 'constants.CrawlerTestURLExamplePage2',
        '"https://example.com/test"': 'constants.CrawlerTestURLExampleTest',
        '"https://example.com/page1"': 'constants.CrawlerTestURLExamplePage1',
        '"https://example.com/api"': 'constants.CrawlerTestURLExampleAPI',
        '"https://example.com"': 'constants.CrawlerTestURLExample',
        '"http://localhost:8080"': 'constants.CrawlerTestURLLocalhost8080',
        '"http://localhost"': 'constants.CrawlerTestURLLocalhost',
        
        # HTTP Status Codes
        '"200"': 'constants.CrawlerHTTPStatusOK',
        '"201"': 'constants.CrawlerHTTPStatusCreated',
        '"202"': 'constants.CrawlerHTTPStatusAccepted',
        '"204"': 'constants.CrawlerHTTPStatusNoContent',
        '"301"': 'constants.CrawlerHTTPStatusMovedPermanently',
        '"302"': 'constants.CrawlerHTTPStatusFound',
        '"304"': 'constants.CrawlerHTTPStatusNotModified',
        '"400"': 'constants.CrawlerHTTPStatusBadRequest',
        '"401"': 'constants.CrawlerHTTPStatusUnauthorized',
        '"403"': 'constants.CrawlerHTTPStatusForbidden',
        '"404"': 'constants.CrawlerHTTPStatusNotFound',
        '"405"': 'constants.CrawlerHTTPStatusMethodNotAllowed',
        '"429"': 'constants.CrawlerHTTPStatusTooManyRequests',
        '"500"': 'constants.CrawlerHTTPStatusInternalServerError',
        '"502"': 'constants.CrawlerHTTPStatusBadGateway',
        '"503"': 'constants.CrawlerHTTPStatusServiceUnavailable',
        '"504"': 'constants.CrawlerHTTPStatusGatewayTimeout',
        
        # Content Types
        '"text/html"': 'constants.CrawlerContentTypeHTML',
        '"text/plain"': 'constants.CrawlerContentTypePlain',
        '"application/json"': 'constants.CrawlerContentTypeJSON',
        '"application/xml"': 'constants.CrawlerContentTypeXMLApp',
        '"text/xml"': 'constants.CrawlerContentTypeXML',
        '"text/css"': 'constants.CrawlerContentTypeCSS',
        '"text/javascript"': 'constants.CrawlerContentTypeJavaScript',
        '"application/javascript"': 'constants.CrawlerContentTypeJavaScriptApp',
        
        # Response Headers
        '"Content-Type"': 'constants.CrawlerResponseHeaderContentType',
        '"Content-Length"': 'constants.CrawlerResponseHeaderContentLength',
        '"Last-Modified"': 'constants.CrawlerResponseHeaderLastModified',
        '"ETag"': 'constants.CrawlerResponseHeaderETag',
        '"Cache-Control"': 'constants.CrawlerResponseHeaderCacheControl',
        '"Location"': 'constants.CrawlerResponseHeaderLocation',
        '"Server"': 'constants.CrawlerResponseHeaderServer',
        
        # Encoding Types
        '"UTF-8"': 'constants.CrawlerEncodingUTF8',
        '"UTF-16"': 'constants.CrawlerEncodingUTF16',
        '"ISO-8859-1"': 'constants.CrawlerEncodingISO88591',
        '"ASCII"': 'constants.CrawlerEncodingASCII',
        '"utf-8"': 'constants.CrawlerEncodingUTF8',
        
        # HTTP Methods
        '"GET"': 'constants.CrawlerHTTPMethodGet',
        '"POST"': 'constants.CrawlerHTTPMethodPost',
        '"PUT"': 'constants.CrawlerHTTPMethodPut',
        '"DELETE"': 'constants.CrawlerHTTPMethodDelete',
        '"HEAD"': 'constants.CrawlerHTTPMethodHead',
        '"OPTIONS"': 'constants.CrawlerHTTPMethodOptions',
        
        # HTML Elements (seulement dans contexte parsing/extraction)
        '"html"': 'constants.CrawlerHTMLElementHTML',
        '"head"': 'constants.CrawlerHTMLElementHead',
        '"body"': 'constants.CrawlerHTMLElementBody',
        '"title"': 'constants.CrawlerHTMLElementTitle',
        '"meta"': 'constants.CrawlerHTMLElementMeta',
        '"link"': 'constants.CrawlerHTMLElementLink',
        '"script"': 'constants.CrawlerHTMLElementScript',
        '"a"': 'constants.CrawlerHTMLElementA',
        '"img"': 'constants.CrawlerHTMLElementImg',
        '"div"': 'constants.CrawlerHTMLElementDiv',
        '"p"': 'constants.CrawlerHTMLElementP',
        '"h1"': 'constants.CrawlerHTMLElementH1',
        '"h2"': 'constants.CrawlerHTMLElementH2',
        '"h3"': 'constants.CrawlerHTMLElementH3',
        
        # HTML Attributes (seulement dans contexte parsing)
        '"href"': 'constants.CrawlerHTMLAttributeHref',
        '"src"': 'constants.CrawlerHTMLAttributeSrc',
        '"alt"': 'constants.CrawlerHTMLAttributeAlt',
        '"title"': 'constants.CrawlerHTMLAttributeTitle',
        '"class"': 'constants.CrawlerHTMLAttributeClass',
        '"id"': 'constants.CrawlerHTMLAttributeID',
        '"name"': 'constants.CrawlerHTMLAttributeName',
        '"content"': 'constants.CrawlerHTMLAttributeContent',
        '"type"': 'constants.CrawlerHTMLAttributeType',
        '"rel"': 'constants.CrawlerHTMLAttributeRel',
        
        # Crawler States (dans contexte approprié)
        '"idle"': 'constants.CrawlerStateIdle',
        '"running"': 'constants.CrawlerStateRunning',
        '"paused"': 'constants.CrawlerStatePaused',
        '"stopped"': 'constants.CrawlerStateStopped',
        '"completed"': 'constants.CrawlerStateCompleted',
        '"failed"': 'constants.CrawlerStateFailed',
        '"timeout"': 'constants.CrawlerStateTimeout',
        '"blocked"': 'constants.CrawlerStateBlocked',
        '"rate_limited"': 'constants.CrawlerStateRateLimited',
        
        # Crawler Config Keys
        '"timeout"': 'constants.CrawlerConfigTimeout',
        '"delay"': 'constants.CrawlerConfigDelay',
        '"concurrency"': 'constants.CrawlerConfigConcurrency',
        '"max_depth"': 'constants.CrawlerConfigMaxDepth',
        '"max_pages"': 'constants.CrawlerConfigMaxPages',
        '"user_agent"': 'constants.CrawlerConfigUserAgent',
        '"follow_redirects"': 'constants.CrawlerConfigFollowRedirects',
        '"respect_robots"': 'constants.CrawlerConfigRespectRobots',
        '"headers"': 'constants.CrawlerConfigHeaders',
        '"cookies"': 'constants.CrawlerConfigCookies',
        '"retries"': 'constants.CrawlerConfigRetries',
        
        # Error Messages (messages complets)
        '"Fetch failed: %v"': 'constants.CrawlerErrorFetchFailed + ": %v"',
        '"Fetch failed after retries: %v"': 'constants.CrawlerErrorFetchFailed + " after retries: %v"',
        '"Rate limiter wait failed: %v"': 'constants.CrawlerErrorRateLimited + " wait failed: %v"',
        '"Invalid URL: %s"': 'constants.CrawlerErrorInvalidURL + ": %s"',
        '"Timeout: %v"': 'constants.CrawlerErrorTimeout + ": %v"',
        '"Connection refused"': 'constants.CrawlerErrorConnectionRefused',
        '"Parsing failed"': 'constants.CrawlerErrorParsingFailed',
        '"Content too large"': 'constants.CrawlerErrorContentTooLarge',
        '"Rate limited"': 'constants.CrawlerErrorRateLimited',
        '"Too many requests"': 'constants.CrawlerErrorTooManyRequests',
        '"Access blocked"': 'constants.CrawlerErrorBlocked',
        '"Access forbidden"': 'constants.CrawlerErrorForbidden',
    }
    
    # Contextes à éviter (où ne PAS remplacer)
    avoid_contexts = [
        r'`json:', # struct tags
        r'//', # commentaires
        r'import', # imports
        r'package', # package declaration
        r'const\s+\w+\s*=', # déclarations de constantes
        r'fmt\.', # format strings (éviter les remplacements dans Printf, etc.)
        r'log\.', # log statements
        r't\.', # testing methods
        r'assert\.', # assertion methods
        r'require\.', # require methods
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
                    
                    # Pour les HTML elements et attributes, seulement dans contexte parsing
                    html_elements = ['"html"', '"head"', '"body"', '"title"', '"meta"', '"link"', 
                                   '"script"', '"a"', '"img"', '"div"', '"p"', '"h1"', '"h2"', '"h3"']
                    html_attributes = ['"href"', '"src"', '"alt"', '"title"', '"class"', '"id"', 
                                     '"name"', '"content"', '"type"', '"rel"']
                    
                    if old_value in html_elements or old_value in html_attributes:
                        # Seulement remplacer si on est dans un contexte de parsing HTML
                        if ('Element' in line or 'Attribute' in line or 'Parse' in line or 
                            'Extract' in line or 'Find' in line or 'Select' in line or
                            'querySelector' in line or 'getElementsBy' in line):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les crawler states, seulement dans contexte de state management
                    elif old_value in ['"idle"', '"running"', '"paused"', '"stopped"', '"completed"', 
                                     '"failed"', '"timeout"', '"blocked"', '"rate_limited"']:
                        # Seulement dans contexte de state/status
                        if ('State' in line or 'Status' in line or '==' in line or '!=' in line or
                            'switch' in line or 'case' in line):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les config keys, seulement dans contexte configuration
                    elif old_value in ['"timeout"', '"delay"', '"concurrency"', '"max_depth"', '"max_pages"',
                                     '"user_agent"', '"follow_redirects"', '"respect_robots"', '"headers"',
                                     '"cookies"', '"retries"']:
                        # Seulement dans contexte config/map
                        if ('Config' in line or '[' in line or 'map' in line or '=' in line):
                            modified_line = modified_line.replace(old_value, new_value)
                            eliminated_count += 1
                            changes_made.append({
                                'line': line_num,
                                'old': old_value,
                                'new': new_value,
                                'context': original_line.strip()
                            })
                    
                    # Pour les autres (URLs, status codes, content types, headers, etc.)
                    else:
                        # Éviter les remplacements dans les format strings et tests complexes
                        if not any(avoid in line for avoid in ['%s', '%v', '%d', 'Printf', 'Sprintf', 'Errorf']):
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
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/crawler/crawler_test.go'
    
    print("⚔️ DELTA-6 ELIMINATOR - Attaque en cours...")
    print("=" * 60)
    
    # Charger l'analyse
    analysis = load_analysis()
    total_violations = analysis['total_violations']
    
    print(f"🎯 Cible: crawler_test.go ({total_violations} violations détectées)")
    
    # Créer backup
    backup_path = create_backup(filepath)
    
    try:
        # Ajouter l'import des constantes
        add_constants_import(filepath)
        
        # Exécuter l'élimination
        eliminated, changes = execute_elimination(filepath, analysis)
        
        if eliminated > 0:
            print(f"\n🏆 DELTA-6 SUCCÈS!")
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
            
            with open('delta6_elimination_report.json', 'w') as f:
                json.dump(report, f, indent=2)
            
            print(f"📊 Rapport sauvegardé: delta6_elimination_report.json")
            
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