#!/usr/bin/env python3
"""
🔍 DELTA-6 DETECTOR
Détection spécialisée pour crawler_test.go - 96 violations
"""

import re
import json
from collections import defaultdict

def analyze_crawler_test_file(filepath):
    """Analyse le fichier crawler_test.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns spécialisés pour CRAWLER TESTING
    patterns = {
        'test_urls': r'"https?://[^"]*"',
        'user_agents': r'"Mozilla/[^"]*|Chrome/[^"]*|Safari/[^"]*|Firefox/[^"]*|Edge/[^"]*|bot/[^"]*|crawler/[^"]*"',
        'http_methods': r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
        'http_status_codes': r'"?(200|201|202|204|301|302|304|400|401|403|404|405|409|422|429|500|501|502|503|504)"?',
        'content_types': r'"(text/html|text/plain|application/json|application/xml|text/css|text/javascript|image/[^"]*)"',
        'html_elements': r'"(html|head|body|title|meta|link|script|style|div|span|p|h1|h2|h3|h4|h5|h6|a|img|ul|ol|li|table|tr|td|th|form|input|button)"',
        'html_attributes': r'"(href|src|alt|title|class|id|name|content|charset|type|rel|media|onclick|onload|style)"',
        'css_selectors': r'"(#[^"]*|\\.[^"]*|\\[[^"]*\\]|[a-z]+:[^"]*)"',
        'xpath_expressions': r'"/[^"]*\\[[^"]*\\]|//[^"]*"',
        'encoding_types': r'"(utf-8|utf-16|iso-8859-1|windows-1252|ascii)"',
        'robots_directives': r'"(User-agent|Disallow|Allow|Crawl-delay|Sitemap|Host)"',
        'sitemap_elements': r'"(urlset|url|loc|lastmod|changefreq|priority)"',
        'crawler_config_keys': r'"(timeout|delay|concurrency|max_depth|max_pages|follow_redirects|respect_robots|user_agent|headers|cookies|proxies)"',
        'test_scenarios': r'"(test_[a-z_]*|should_[a-z_]*|when_[a-z_]*|given_[a-z_]*)"',
        'error_messages': r'"[A-Z][^"]*(?:error|failed|invalid|timeout|not found|forbidden|denied|blocked)[^"]*"',
        'log_messages': r'"[A-Z][^"]*(?:crawling|fetching|parsing|extracting|following|visiting|processing|analyzing)[^"]*"',
        'mime_types': r'"(text/[^"]*|application/[^"]*|image/[^"]*|audio/[^"]*|video/[^"]*|multipart/[^"]*)"',
        'file_extensions': r'"\\.(html|htm|php|asp|jsp|xml|txt|css|js|json|pdf|doc|docx|xls|xlsx|ppt|pptx|zip|tar|gz)"',
        'test_assertions': r'"(assert|expect|should|must|verify|validate|check|ensure|contain|equal|match|include)"',
        'crawler_states': r'"(idle|running|paused|stopped|completed|failed|timeout|blocked|rate_limited)"',
        'link_types': r'"(internal|external|absolute|relative|anchor|mailto|tel|ftp|javascript)"',
        'response_headers': r'"(Content-Type|Content-Length|Last-Modified|ETag|Cache-Control|Set-Cookie|Location|Server|X-[^"]*)"',
        'test_data': r'"(mock_[^"]*|test_[^"]*|sample_[^"]*|dummy_[^"]*|fake_[^"]*)"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires, imports et struct tags
        if (line_stripped.startswith('//') or 
            line_stripped.startswith('import') or 
            line_stripped.startswith('package') or
            '`json:' in line_stripped):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped, re.IGNORECASE)
            for match in matches:
                # Nettoyer le match
                if isinstance(match, tuple):
                    match_value = match[0] if match[0] else match[1] if len(match) > 1 else str(match)
                else:
                    match_value = match
                    
                if len(match_value.strip('"')) < 2:
                    continue
                    
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match_value,
                    'context': line_stripped[:150] + ('...' if len(line_stripped) > 150 else '')
                })
    
    return violations

def categorize_violations(violations):
    """Catégorise les violations par type"""
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    return dict(categories)

def generate_crawler_test_constants_mapping(violations):
    """Génère les mappings de constantes pour crawler_test.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Test URLs
    if 'test_urls' in categorized:
        for v in categorized['test_urls']:
            value = v['value'].strip('"')
            if 'example.com' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerTestURLExample'
            elif 'test.com' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerTestURLTest'
            elif 'localhost' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerTestURLLocalhost'
            elif 'httpbin.org' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerTestURLHTTPBin'
            else:
                # Générer un nom basé sur le domaine
                try:
                    domain = value.split('/')[2].replace('.', '').replace('-', '').title()
                    constants_map[f'"{value}"'] = f'constants.CrawlerTestURL{domain}'
                except:
                    constants_map[f'"{value}"'] = 'constants.CrawlerTestURLGeneric'
    
    # User Agents
    if 'user_agents' in categorized:
        for v in categorized['user_agents']:
            value = v['value'].strip('"')
            if 'Mozilla' in value and 'Chrome' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerUserAgentChrome'
            elif 'Mozilla' in value and 'Firefox' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerUserAgentFirefox'
            elif 'Mozilla' in value and 'Safari' in value:
                constants_map[f'"{value}"'] = 'constants.CrawlerUserAgentSafari'
            elif 'bot' in value.lower() or 'crawler' in value.lower():
                constants_map[f'"{value}"'] = 'constants.CrawlerUserAgentBot'
            else:
                constants_map[f'"{value}"'] = 'constants.CrawlerUserAgentDefault'
    
    # HTTP Methods
    if 'http_methods' in categorized:
        method_map = {
            'GET': 'constants.CrawlerHTTPMethodGet',
            'POST': 'constants.CrawlerHTTPMethodPost',
            'PUT': 'constants.CrawlerHTTPMethodPut',
            'DELETE': 'constants.CrawlerHTTPMethodDelete',
            'PATCH': 'constants.CrawlerHTTPMethodPatch',
            'HEAD': 'constants.CrawlerHTTPMethodHead',
            'OPTIONS': 'constants.CrawlerHTTPMethodOptions'
        }
        for v in categorized['http_methods']:
            if v['value'] in method_map:
                constants_map[f'"{v["value"]}"'] = method_map[v['value']]
    
    # HTTP Status Codes
    if 'http_status_codes' in categorized:
        status_map = {
            '200': 'constants.CrawlerHTTPStatusOK',
            '201': 'constants.CrawlerHTTPStatusCreated',
            '202': 'constants.CrawlerHTTPStatusAccepted',
            '204': 'constants.CrawlerHTTPStatusNoContent',
            '301': 'constants.CrawlerHTTPStatusMovedPermanently',
            '302': 'constants.CrawlerHTTPStatusFound',
            '304': 'constants.CrawlerHTTPStatusNotModified',
            '400': 'constants.CrawlerHTTPStatusBadRequest',
            '401': 'constants.CrawlerHTTPStatusUnauthorized',
            '403': 'constants.CrawlerHTTPStatusForbidden',
            '404': 'constants.CrawlerHTTPStatusNotFound',
            '405': 'constants.CrawlerHTTPStatusMethodNotAllowed',
            '409': 'constants.CrawlerHTTPStatusConflict',
            '422': 'constants.CrawlerHTTPStatusUnprocessableEntity',
            '429': 'constants.CrawlerHTTPStatusTooManyRequests',
            '500': 'constants.CrawlerHTTPStatusInternalServerError',
            '501': 'constants.CrawlerHTTPStatusNotImplemented',
            '502': 'constants.CrawlerHTTPStatusBadGateway',
            '503': 'constants.CrawlerHTTPStatusServiceUnavailable',
            '504': 'constants.CrawlerHTTPStatusGatewayTimeout'
        }
        for v in categorized['http_status_codes']:
            code = v['value'].strip('"')
            if code in status_map:
                constants_map[f'"{code}"'] = status_map[code]
    
    # Content Types
    if 'content_types' in categorized:
        content_map = {
            'text/html': 'constants.CrawlerContentTypeHTML',
            'text/plain': 'constants.CrawlerContentTypePlain',
            'application/json': 'constants.CrawlerContentTypeJSON',
            'application/xml': 'constants.CrawlerContentTypeXML',
            'text/css': 'constants.CrawlerContentTypeCSS',
            'text/javascript': 'constants.CrawlerContentTypeJavaScript',
            'application/javascript': 'constants.CrawlerContentTypeJavaScript'
        }
        for v in categorized['content_types']:
            value = v['value'].strip('"')
            if value in content_map:
                constants_map[f'"{value}"'] = content_map[value]
            elif value.startswith('image/'):
                constants_map[f'"{value}"'] = 'constants.CrawlerContentTypeImage'
    
    # HTML Elements
    if 'html_elements' in categorized:
        element_map = {
            'html': 'constants.CrawlerHTMLElementHTML',
            'head': 'constants.CrawlerHTMLElementHead',
            'body': 'constants.CrawlerHTMLElementBody',
            'title': 'constants.CrawlerHTMLElementTitle',
            'meta': 'constants.CrawlerHTMLElementMeta',
            'link': 'constants.CrawlerHTMLElementLink',
            'script': 'constants.CrawlerHTMLElementScript',
            'style': 'constants.CrawlerHTMLElementStyle',
            'div': 'constants.CrawlerHTMLElementDiv',
            'span': 'constants.CrawlerHTMLElementSpan',
            'p': 'constants.CrawlerHTMLElementP',
            'h1': 'constants.CrawlerHTMLElementH1',
            'h2': 'constants.CrawlerHTMLElementH2',
            'h3': 'constants.CrawlerHTMLElementH3',
            'a': 'constants.CrawlerHTMLElementA',
            'img': 'constants.CrawlerHTMLElementImg',
            'ul': 'constants.CrawlerHTMLElementUL',
            'ol': 'constants.CrawlerHTMLElementOL',
            'li': 'constants.CrawlerHTMLElementLI',
            'table': 'constants.CrawlerHTMLElementTable',
            'tr': 'constants.CrawlerHTMLElementTR',
            'td': 'constants.CrawlerHTMLElementTD',
            'th': 'constants.CrawlerHTMLElementTH',
            'form': 'constants.CrawlerHTMLElementForm',
            'input': 'constants.CrawlerHTMLElementInput',
            'button': 'constants.CrawlerHTMLElementButton'
        }
        for v in categorized['html_elements']:
            if v['value'] in element_map:
                constants_map[f'"{v["value"]}"'] = element_map[v['value']]
    
    # HTML Attributes
    if 'html_attributes' in categorized:
        attr_map = {
            'href': 'constants.CrawlerHTMLAttributeHref',
            'src': 'constants.CrawlerHTMLAttributeSrc',
            'alt': 'constants.CrawlerHTMLAttributeAlt',
            'title': 'constants.CrawlerHTMLAttributeTitle',
            'class': 'constants.CrawlerHTMLAttributeClass',
            'id': 'constants.CrawlerHTMLAttributeID',
            'name': 'constants.CrawlerHTMLAttributeName',
            'content': 'constants.CrawlerHTMLAttributeContent',
            'charset': 'constants.CrawlerHTMLAttributeCharset',
            'type': 'constants.CrawlerHTMLAttributeType',
            'rel': 'constants.CrawlerHTMLAttributeRel',
            'media': 'constants.CrawlerHTMLAttributeMedia',
            'style': 'constants.CrawlerHTMLAttributeStyle'
        }
        for v in categorized['html_attributes']:
            if v['value'] in attr_map:
                constants_map[f'"{v["value"]}"'] = attr_map[v['value']]
    
    # Encoding Types
    if 'encoding_types' in categorized:
        encoding_map = {
            'utf-8': 'constants.CrawlerEncodingUTF8',
            'utf-16': 'constants.CrawlerEncodingUTF16',
            'iso-8859-1': 'constants.CrawlerEncodingISO88591',
            'windows-1252': 'constants.CrawlerEncodingWindows1252',
            'ascii': 'constants.CrawlerEncodingASCII'
        }
        for v in categorized['encoding_types']:
            if v['value'] in encoding_map:
                constants_map[f'"{v["value"]}"'] = encoding_map[v['value']]
    
    # Robots Directives
    if 'robots_directives' in categorized:
        robots_map = {
            'User-agent': 'constants.CrawlerRobotsUserAgent',
            'Disallow': 'constants.CrawlerRobotsDisallow',
            'Allow': 'constants.CrawlerRobotsAllow',
            'Crawl-delay': 'constants.CrawlerRobotsCrawlDelay',
            'Sitemap': 'constants.CrawlerRobotsSitemap',
            'Host': 'constants.CrawlerRobotsHost'
        }
        for v in categorized['robots_directives']:
            if v['value'] in robots_map:
                constants_map[f'"{v["value"]}"'] = robots_map[v['value']]
    
    # Sitemap Elements
    if 'sitemap_elements' in categorized:
        sitemap_map = {
            'urlset': 'constants.CrawlerSitemapURLSet',
            'url': 'constants.CrawlerSitemapURL',
            'loc': 'constants.CrawlerSitemapLoc',
            'lastmod': 'constants.CrawlerSitemapLastmod',
            'changefreq': 'constants.CrawlerSitemapChangefreq',
            'priority': 'constants.CrawlerSitemapPriority'
        }
        for v in categorized['sitemap_elements']:
            if v['value'] in sitemap_map:
                constants_map[f'"{v["value"]}"'] = sitemap_map[v['value']]
    
    # Crawler Config Keys
    if 'crawler_config_keys' in categorized:
        config_map = {
            'timeout': 'constants.CrawlerConfigTimeout',
            'delay': 'constants.CrawlerConfigDelay',
            'concurrency': 'constants.CrawlerConfigConcurrency',
            'max_depth': 'constants.CrawlerConfigMaxDepth',
            'max_pages': 'constants.CrawlerConfigMaxPages',
            'follow_redirects': 'constants.CrawlerConfigFollowRedirects',
            'respect_robots': 'constants.CrawlerConfigRespectRobots',
            'user_agent': 'constants.CrawlerConfigUserAgent',
            'headers': 'constants.CrawlerConfigHeaders',
            'cookies': 'constants.CrawlerConfigCookies',
            'proxies': 'constants.CrawlerConfigProxies'
        }
        for v in categorized['crawler_config_keys']:
            if v['value'] in config_map:
                constants_map[f'"{v["value"]}"'] = config_map[v['value']]
    
    # Crawler States
    if 'crawler_states' in categorized:
        state_map = {
            'idle': 'constants.CrawlerStateIdle',
            'running': 'constants.CrawlerStateRunning',
            'paused': 'constants.CrawlerStatePaused',
            'stopped': 'constants.CrawlerStateStopped',
            'completed': 'constants.CrawlerStateCompleted',
            'failed': 'constants.CrawlerStateFailed',
            'timeout': 'constants.CrawlerStateTimeout',
            'blocked': 'constants.CrawlerStateBlocked',
            'rate_limited': 'constants.CrawlerStateRateLimited'
        }
        for v in categorized['crawler_states']:
            if v['value'] in state_map:
                constants_map[f'"{v["value"]}"'] = state_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/crawler/crawler_test.go'
    
    print("🔍 DELTA-6 DETECTOR - Scanning crawler_test.go...")
    
    violations = analyze_crawler_test_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_crawler_test_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-6:")
    print(f"Total violations détectées: {len(violations)}")
    
    for category, viols in categorized.items():
        print(f"\n🔸 {category.upper()}: {len(viols)} violations")
        for v in viols[:3]:  # Show first 3 of each category
            print(f"  Line {v['line']}: {v['value']}")
        if len(viols) > 3:
            print(f"  ... et {len(viols) - 3} autres")
    
    print(f"\n🏗️ CONSTANTES À CRÉER: {len(constants_map)}")
    print("Preview des mappings:")
    for original, constant in list(constants_map.items())[:10]:
        print(f"  {original} → {constant}")
    
    # Sauvegarder les résultats
    results = {
        'total_violations': len(violations),
        'categories': {k: len(v) for k, v in categorized.items()},
        'violations': violations,
        'constants_mapping': constants_map
    }
    
    with open('delta6_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta6_analysis.json")
    return results

if __name__ == "__main__":
    main()