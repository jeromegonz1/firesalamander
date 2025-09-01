#!/usr/bin/env python3
"""
Fire Salamander - SEO Accuracy Agent
Validation SEO automatisée et analyse de conformité
"""

import json
import os
import time
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin, urlparse
from datetime import datetime
from typing import Dict, List, Optional, Tuple
import logging
import re

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class SEOAgent:
    """Agent de validation SEO pour Fire Salamander"""
    
    def __init__(self, config: Optional[Dict] = None):
        self.config = config or self.default_config()
        self.stats = {
            'timestamp': datetime.now().isoformat(),
            'target_url': self.config['target_url'],
            'seo_results': {},
            'issues': [],
            'recommendations': [],
            'overall_score': 0,
            'status': 'unknown'
        }
        
    def default_config(self) -> Dict:
        """Configuration par défaut"""
        return {
            'target_url': 'http://localhost:3000',
            'timeout': 30,
            'report_path': 'tests/reports/seo',
            'user_agent': 'Mozilla/5.0 (SEO Agent) Fire Salamander/1.0',
            'max_pages': 10,  # Limite pour éviter trop de requêtes
            'seo_rules': {
                'title_min_length': 30,
                'title_max_length': 60,
                'meta_description_min_length': 120,
                'meta_description_max_length': 160,
                'h1_max_count': 1,
                'image_alt_required': True,
                'internal_links_min': 2
            }
        }
    
    def run_full_seo_audit(self) -> Dict:
        """Lance un audit SEO complet"""
        logger.info("🔍 Starting SEO accuracy audit")
        
        try:
            # 1. Analyser la page d'accueil
            self._analyze_homepage()
            
            # 2. Découvrir les pages du site
            discovered_pages = self._discover_pages()
            
            # 3. Analyser chaque page découverte
            for page_url in discovered_pages[:self.config['max_pages']]:
                self._analyze_page(page_url)
            
            # 4. Tests structurels globaux
            self._test_site_structure()
            
            # 5. Tests techniques SEO
            self._test_technical_seo()
            
            # 6. Calculer le score global
            self._calculate_seo_score()
            
            # 7. Générer les recommandations
            self._generate_recommendations()
            
            # 8. Générer le rapport
            self._generate_report()
            
            logger.info(f"✅ SEO audit completed - Score: {self.stats['overall_score']}/100")
            
        except Exception as e:
            logger.error(f"❌ SEO audit failed: {str(e)}")
            self.stats['status'] = 'error'
            self.stats['error'] = str(e)
        
        return self.stats
    
    def _analyze_homepage(self):
        """Analyse la page d'accueil"""
        logger.info("🏠 Analyzing homepage SEO")
        
        try:
            response = requests.get(
                self.config['target_url'], 
                timeout=self.config['timeout'],
                headers={'User-Agent': self.config['user_agent']}
            )
            
            if response.status_code != 200:
                self.stats['issues'].append({
                    'type': 'critical',
                    'page': 'Homepage',
                    'issue': f'Homepage not accessible (HTTP {response.status_code})',
                    'impact': 'Critical - Search engines cannot index the site'
                })
                return
            
            soup = BeautifulSoup(response.text, 'html.parser')
            homepage_analysis = self._analyze_page_seo(soup, self.config['target_url'])
            
            self.stats['seo_results']['homepage'] = homepage_analysis
            
        except Exception as e:
            self.stats['issues'].append({
                'type': 'error',
                'page': 'Homepage',
                'issue': f'Cannot analyze homepage: {str(e)}',
                'impact': 'Critical'
            })
    
    def _discover_pages(self) -> List[str]:
        """Découvre les pages du site via les liens internes"""
        logger.info("🕷️ Discovering site pages")
        
        discovered_pages = set()
        
        try:
            response = requests.get(
                self.config['target_url'],
                timeout=self.config['timeout'],
                headers={'User-Agent': self.config['user_agent']}
            )
            
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                
                # Trouver tous les liens internes
                for link in soup.find_all('a', href=True):
                    href = link['href']
                    full_url = urljoin(self.config['target_url'], href)
                    
                    # Vérifier si c'est un lien interne
                    parsed_base = urlparse(self.config['target_url'])
                    parsed_link = urlparse(full_url)
                    
                    if (parsed_link.netloc == parsed_base.netloc and 
                        not href.startswith('#') and 
                        not href.startswith('mailto:') and
                        not href.startswith('tel:')):
                        discovered_pages.add(full_url)
                
                # Ajouter des pages communes si elles existent
                common_pages = ['/about', '/contact', '/services', '/blog', '/sitemap.xml', '/robots.txt']
                for page in common_pages:
                    test_url = urljoin(self.config['target_url'], page)
                    try:
                        test_response = requests.head(test_url, timeout=5)
                        if test_response.status_code == 200:
                            discovered_pages.add(test_url)
                    except:
                        pass
            
        except Exception as e:
            logger.warning(f"Page discovery failed: {str(e)}")
        
        discovered_list = list(discovered_pages)
        logger.info(f"📄 Discovered {len(discovered_list)} pages")
        
        return discovered_list
    
    def _analyze_page(self, url: str):
        """Analyse une page spécifique"""
        try:
            response = requests.get(
                url,
                timeout=self.config['timeout'],
                headers={'User-Agent': self.config['user_agent']}
            )
            
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                page_analysis = self._analyze_page_seo(soup, url)
                
                # Utiliser l'URL comme clé (nettoyer pour éviter les caractères spéciaux)
                page_key = url.replace(self.config['target_url'], '').strip('/') or 'root'
                page_key = re.sub(r'[^a-zA-Z0-9_-]', '_', page_key)
                
                self.stats['seo_results'][f'page_{page_key}'] = page_analysis
                
        except Exception as e:
            logger.warning(f"Cannot analyze page {url}: {str(e)}")
    
    def _analyze_page_seo(self, soup: BeautifulSoup, url: str) -> Dict:
        """Analyse SEO détaillée d'une page"""
        analysis = {
            'url': url,
            'title': {},
            'meta_description': {},
            'headings': {},
            'images': {},
            'links': {},
            'content': {},
            'issues': [],
            'score': 0
        }
        
        # 1. Analyse du title
        title_tag = soup.find('title')
        if title_tag:
            title_text = title_tag.get_text().strip()
            analysis['title'] = {
                'text': title_text,
                'length': len(title_text),
                'status': 'good'
            }
            
            # Vérifier la longueur
            if len(title_text) < self.config['seo_rules']['title_min_length']:
                analysis['title']['status'] = 'warning'
                analysis['issues'].append({
                    'type': 'warning',
                    'element': 'title',
                    'issue': f'Title too short ({len(title_text)} chars, minimum {self.config["seo_rules"]["title_min_length"]})'
                })
            elif len(title_text) > self.config['seo_rules']['title_max_length']:
                analysis['title']['status'] = 'warning'
                analysis['issues'].append({
                    'type': 'warning',
                    'element': 'title',
                    'issue': f'Title too long ({len(title_text)} chars, maximum {self.config["seo_rules"]["title_max_length"]})'
                })
        else:
            analysis['title'] = {'status': 'critical', 'text': '', 'length': 0}
            analysis['issues'].append({
                'type': 'critical',
                'element': 'title',
                'issue': 'Missing title tag'
            })
        
        # 2. Analyse de la meta description
        meta_desc = soup.find('meta', attrs={'name': 'description'})
        if meta_desc and meta_desc.get('content'):
            desc_text = meta_desc['content'].strip()
            analysis['meta_description'] = {
                'text': desc_text,
                'length': len(desc_text),
                'status': 'good'
            }
            
            # Vérifier la longueur
            if len(desc_text) < self.config['seo_rules']['meta_description_min_length']:
                analysis['meta_description']['status'] = 'warning'
                analysis['issues'].append({
                    'type': 'warning',
                    'element': 'meta_description',
                    'issue': f'Meta description too short ({len(desc_text)} chars)'
                })
            elif len(desc_text) > self.config['seo_rules']['meta_description_max_length']:
                analysis['meta_description']['status'] = 'warning'
                analysis['issues'].append({
                    'type': 'warning',
                    'element': 'meta_description',
                    'issue': f'Meta description too long ({len(desc_text)} chars)'
                })
        else:
            analysis['meta_description'] = {'status': 'critical', 'text': '', 'length': 0}
            analysis['issues'].append({
                'type': 'critical',
                'element': 'meta_description',
                'issue': 'Missing meta description'
            })
        
        # 3. Analyse des headings
        headings_analysis = {'h1': 0, 'h2': 0, 'h3': 0, 'h4': 0, 'h5': 0, 'h6': 0}
        for i in range(1, 7):
            headings = soup.find_all(f'h{i}')
            headings_analysis[f'h{i}'] = len(headings)
        
        analysis['headings'] = headings_analysis
        
        # Vérifier H1
        if headings_analysis['h1'] == 0:
            analysis['issues'].append({
                'type': 'critical',
                'element': 'h1',
                'issue': 'Missing H1 tag'
            })
        elif headings_analysis['h1'] > self.config['seo_rules']['h1_max_count']:
            analysis['issues'].append({
                'type': 'warning',
                'element': 'h1',
                'issue': f'Multiple H1 tags found ({headings_analysis["h1"]})'
            })
        
        # 4. Analyse des images
        images = soup.find_all('img')
        images_without_alt = [img for img in images if not img.get('alt')]
        
        analysis['images'] = {
            'total': len(images),
            'without_alt': len(images_without_alt),
            'status': 'good' if len(images_without_alt) == 0 else 'warning'
        }
        
        if images_without_alt:
            analysis['issues'].append({
                'type': 'warning',
                'element': 'images',
                'issue': f'{len(images_without_alt)} images without alt text'
            })
        
        # 5. Analyse des liens
        internal_links = []
        external_links = []
        
        for link in soup.find_all('a', href=True):
            href = link['href']
            if href.startswith('http'):
                parsed_link = urlparse(href)
                parsed_base = urlparse(url)
                if parsed_link.netloc == parsed_base.netloc:
                    internal_links.append(href)
                else:
                    external_links.append(href)
            elif not href.startswith('#') and not href.startswith('mailto:'):
                internal_links.append(href)
        
        analysis['links'] = {
            'internal': len(internal_links),
            'external': len(external_links),
            'status': 'good'
        }
        
        if len(internal_links) < self.config['seo_rules']['internal_links_min']:
            analysis['links']['status'] = 'warning'
            analysis['issues'].append({
                'type': 'warning',
                'element': 'links',
                'issue': f'Few internal links ({len(internal_links)})'
            })
        
        # 6. Analyse du contenu
        text_content = soup.get_text()
        word_count = len(text_content.split())
        
        analysis['content'] = {
            'word_count': word_count,
            'character_count': len(text_content),
            'status': 'good' if word_count > 200 else 'warning'
        }
        
        if word_count < 200:
            analysis['issues'].append({
                'type': 'warning',
                'element': 'content',
                'issue': f'Low content word count ({word_count} words)'
            })
        
        # Calculer le score de la page
        analysis['score'] = self._calculate_page_score(analysis)
        
        return analysis
    
    def _calculate_page_score(self, analysis: Dict) -> int:
        """Calcule le score SEO d'une page"""
        score = 100
        
        # Pénalités par type d'issue
        for issue in analysis['issues']:
            if issue['type'] == 'critical':
                score -= 20
            elif issue['type'] == 'warning':
                score -= 10
            elif issue['type'] == 'info':
                score -= 5
        
        return max(0, score)
    
    def _test_site_structure(self):
        """Tests de structure globale du site"""
        logger.info("🏗️ Testing site structure")
        
        structure_tests = []
        
        # Test 1: Robots.txt
        try:
            robots_url = urljoin(self.config['target_url'], '/robots.txt')
            response = requests.get(robots_url, timeout=10)
            
            if response.status_code == 200:
                structure_tests.append({
                    'test': 'robots.txt',
                    'status': 'pass',
                    'description': 'robots.txt is accessible'
                })
                
                # Analyser le contenu basique
                if 'User-agent:' in response.text:
                    structure_tests.append({
                        'test': 'robots.txt format',
                        'status': 'pass',
                        'description': 'robots.txt has proper format'
                    })
                else:
                    structure_tests.append({
                        'test': 'robots.txt format',
                        'status': 'warning',
                        'description': 'robots.txt format may be incorrect'
                    })
            else:
                structure_tests.append({
                    'test': 'robots.txt',
                    'status': 'warning',
                    'description': 'robots.txt not found (recommended for SEO)'
                })
        except Exception as e:
            structure_tests.append({
                'test': 'robots.txt',
                'status': 'error',
                'description': f'Cannot test robots.txt: {str(e)}'
            })
        
        # Test 2: Sitemap.xml
        try:
            sitemap_url = urljoin(self.config['target_url'], '/sitemap.xml')
            response = requests.get(sitemap_url, timeout=10)
            
            if response.status_code == 200:
                structure_tests.append({
                    'test': 'sitemap.xml',
                    'status': 'pass',
                    'description': 'XML sitemap is accessible'
                })
                
                # Vérifier que c'est du XML valide
                if '<?xml' in response.text and '<urlset' in response.text:
                    structure_tests.append({
                        'test': 'sitemap.xml format',
                        'status': 'pass',
                        'description': 'XML sitemap has proper format'
                    })
                else:
                    structure_tests.append({
                        'test': 'sitemap.xml format',
                        'status': 'warning',
                        'description': 'XML sitemap format may be incorrect'
                    })
            else:
                structure_tests.append({
                    'test': 'sitemap.xml',
                    'status': 'warning',
                    'description': 'XML sitemap not found (recommended for SEO)'
                })
        except Exception as e:
            structure_tests.append({
                'test': 'sitemap.xml',
                'status': 'error',
                'description': f'Cannot test sitemap.xml: {str(e)}'
            })
        
        self.stats['seo_results']['site_structure'] = structure_tests
    
    def _test_technical_seo(self):
        """Tests techniques SEO"""
        logger.info("⚙️ Testing technical SEO")
        
        technical_tests = []
        
        try:
            response = requests.get(
                self.config['target_url'],
                timeout=self.config['timeout'],
                headers={'User-Agent': self.config['user_agent']}
            )
            
            # Test 1: Response time
            if response.elapsed.total_seconds() < 3:
                technical_tests.append({
                    'test': 'Page Load Speed',
                    'status': 'pass',
                    'description': f'Good response time: {response.elapsed.total_seconds():.2f}s'
                })
            else:
                technical_tests.append({
                    'test': 'Page Load Speed',
                    'status': 'warning',
                    'description': f'Slow response time: {response.elapsed.total_seconds():.2f}s'
                })
            
            # Test 2: Content-Type
            content_type = response.headers.get('Content-Type', '')
            if 'text/html' in content_type:
                technical_tests.append({
                    'test': 'Content-Type',
                    'status': 'pass',
                    'description': f'Proper HTML content type: {content_type}'
                })
            else:
                technical_tests.append({
                    'test': 'Content-Type',
                    'status': 'warning',
                    'description': f'Unusual content type: {content_type}'
                })
            
            # Test 3: Charset
            if 'charset=utf-8' in content_type.lower():
                technical_tests.append({
                    'test': 'Character Encoding',
                    'status': 'pass',
                    'description': 'UTF-8 encoding specified'
                })
            else:
                technical_tests.append({
                    'test': 'Character Encoding',
                    'status': 'warning',
                    'description': 'UTF-8 encoding not explicitly specified'
                })
            
            # Test 4: Gzip compression
            encoding = response.headers.get('Content-Encoding', '')
            if 'gzip' in encoding:
                technical_tests.append({
                    'test': 'Compression',
                    'status': 'pass',
                    'description': 'Gzip compression enabled'
                })
            else:
                technical_tests.append({
                    'test': 'Compression',
                    'status': 'warning',
                    'description': 'Gzip compression not detected'
                })
            
            # Test 5: Mobile viewport
            soup = BeautifulSoup(response.text, 'html.parser')
            viewport_meta = soup.find('meta', attrs={'name': 'viewport'})
            
            if viewport_meta:
                technical_tests.append({
                    'test': 'Mobile Viewport',
                    'status': 'pass',
                    'description': f'Viewport meta tag present: {viewport_meta.get("content", "")}'
                })
            else:
                technical_tests.append({
                    'test': 'Mobile Viewport',
                    'status': 'warning',
                    'description': 'Missing viewport meta tag for mobile responsiveness'
                })
            
        except Exception as e:
            technical_tests.append({
                'test': 'Technical SEO',
                'status': 'error',
                'description': f'Cannot perform technical tests: {str(e)}'
            })
        
        self.stats['seo_results']['technical_seo'] = technical_tests
    
    def _calculate_seo_score(self):
        """Calcule le score SEO global"""
        total_score = 0
        page_count = 0
        
        # Moyenne des scores des pages
        for key, results in self.stats['seo_results'].items():
            if key.startswith('page_') or key == 'homepage':
                if isinstance(results, dict) and 'score' in results:
                    total_score += results['score']
                    page_count += 1
        
        if page_count > 0:
            average_page_score = total_score / page_count
        else:
            average_page_score = 0
        
        # Ajustements pour les tests globaux
        adjustment = 0
        
        # Bonus/malus pour structure du site
        if 'site_structure' in self.stats['seo_results']:
            structure_tests = self.stats['seo_results']['site_structure']
            for test in structure_tests:
                if test['status'] == 'pass':
                    adjustment += 2
                elif test['status'] == 'warning':
                    adjustment -= 1
        
        # Bonus/malus pour SEO technique
        if 'technical_seo' in self.stats['seo_results']:
            technical_tests = self.stats['seo_results']['technical_seo']
            for test in technical_tests:
                if test['status'] == 'pass':
                    adjustment += 1
                elif test['status'] == 'warning':
                    adjustment -= 2
        
        final_score = max(0, min(100, int(average_page_score + adjustment)))
        self.stats['overall_score'] = final_score
        
        # Déterminer le statut
        if final_score >= 90:
            self.stats['status'] = 'excellent'
        elif final_score >= 80:
            self.stats['status'] = 'good'
        elif final_score >= 70:
            self.stats['status'] = 'acceptable'
        elif final_score >= 60:
            self.stats['status'] = 'needs_improvement'
        else:
            self.stats['status'] = 'poor'
    
    def _generate_recommendations(self):
        """Génère des recommandations SEO"""
        recommendations = []
        
        # Analyser tous les problèmes trouvés
        all_issues = []
        for key, results in self.stats['seo_results'].items():
            if isinstance(results, dict) and 'issues' in results:
                all_issues.extend(results['issues'])
            elif isinstance(results, list):
                # Pour site_structure et technical_seo
                for test in results:
                    if test['status'] in ['warning', 'fail']:
                        all_issues.append({
                            'type': test['status'],
                            'element': test['test'],
                            'issue': test['description']
                        })
        
        # Générer des recommandations basées sur les problèmes
        issue_types = {}
        for issue in all_issues:
            issue_type = issue.get('element', 'unknown')
            if issue_type not in issue_types:
                issue_types[issue_type] = []
            issue_types[issue_type].append(issue)
        
        # Recommandations par type de problème
        if 'title' in issue_types:
            recommendations.append({
                'priority': 'high',
                'category': 'Content',
                'recommendation': 'Optimize page titles to be between 30-60 characters and include target keywords',
                'impact': 'High - Titles are crucial for search rankings'
            })
        
        if 'meta_description' in issue_types:
            recommendations.append({
                'priority': 'high',
                'category': 'Content',
                'recommendation': 'Add compelling meta descriptions (120-160 characters) for all pages',
                'impact': 'High - Improves click-through rates from search results'
            })
        
        if 'h1' in issue_types:
            recommendations.append({
                'priority': 'high',
                'category': 'Content',
                'recommendation': 'Ensure each page has exactly one H1 tag with target keywords',
                'impact': 'High - H1 tags help search engines understand page structure'
            })
        
        if 'images' in issue_types:
            recommendations.append({
                'priority': 'medium',
                'category': 'Accessibility',
                'recommendation': 'Add descriptive alt text to all images for better accessibility and SEO',
                'impact': 'Medium - Helps with image search and accessibility'
            })
        
        if 'robots.txt' in issue_types:
            recommendations.append({
                'priority': 'medium',
                'category': 'Technical',
                'recommendation': 'Create a robots.txt file to guide search engine crawlers',
                'impact': 'Medium - Helps search engines crawl your site efficiently'
            })
        
        if 'sitemap.xml' in issue_types:
            recommendations.append({
                'priority': 'medium',
                'category': 'Technical',
                'recommendation': 'Generate and maintain an XML sitemap for better indexation',
                'impact': 'Medium - Helps search engines discover all your pages'
            })
        
        # Recommandations générales si le score est bas
        if self.stats['overall_score'] < 70:
            recommendations.append({
                'priority': 'high',
                'category': 'Strategy',
                'recommendation': 'Conduct comprehensive SEO audit and create optimization roadmap',
                'impact': 'Critical - Multiple SEO issues detected requiring immediate attention'
            })
        
        self.stats['recommendations'] = recommendations
    
    def _generate_report(self):
        """Génère le rapport SEO"""
        try:
            os.makedirs(self.config['report_path'], exist_ok=True)
            
            # Rapport JSON
            json_report_path = os.path.join(self.config['report_path'], 'seo_report.json')
            with open(json_report_path, 'w') as f:
                json.dump(self.stats, f, indent=2)
            
            # Rapport HTML
            html_report = self._generate_html_report()
            html_report_path = os.path.join(self.config['report_path'], 'seo_report.html')
            with open(html_report_path, 'w') as f:
                f.write(html_report)
            
            logger.info(f"📊 SEO report generated: {json_report_path}")
            
        except Exception as e:
            logger.error(f"Failed to generate SEO report: {str(e)}")
    
    def _generate_html_report(self) -> str:
        """Génère un rapport HTML"""
        
        # Compter les pages analysées
        page_count = sum(1 for key in self.stats['seo_results'].keys() 
                        if key.startswith('page_') or key == 'homepage')
        
        # Compter les problèmes par priorité
        high_priority = sum(1 for rec in self.stats['recommendations'] if rec['priority'] == 'high')
        medium_priority = sum(1 for rec in self.stats['recommendations'] if rec['priority'] == 'medium')
        low_priority = sum(1 for rec in self.stats['recommendations'] if rec['priority'] == 'low')
        
        html = f"""
<!DOCTYPE html>
<html>
<head>
    <title>Fire Salamander - SEO Report</title>
    <meta charset="UTF-8">
    <style>
        body {{ font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }}
        .header {{ background: #ff6136; color: white; padding: 20px; border-radius: 8px; }}
        .score {{ font-size: 2em; font-weight: bold; }}
        .status-excellent {{ color: #28a745; }}
        .status-good {{ color: #17a2b8; }}
        .status-acceptable {{ color: #ffc107; }}
        .status-needs_improvement {{ color: #fd7e14; }}
        .status-poor {{ color: #dc3545; }}
        .section {{ background: white; margin: 20px 0; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
        .page-analysis {{ margin: 10px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }}
        .issue {{ margin: 5px 0; padding: 8px; border-left: 4px solid #ddd; }}
        .critical {{ border-left-color: #dc3545; background-color: #f8d7da; }}
        .warning {{ border-left-color: #ffc107; background-color: #fff3cd; }}
        .info {{ border-left-color: #17a2b8; background-color: #d1ecf1; }}
        .recommendation {{ margin: 10px 0; padding: 15px; border-radius: 5px; }}
        .high-priority {{ background-color: #f8d7da; border-left: 4px solid #dc3545; }}
        .medium-priority {{ background-color: #fff3cd; border-left: 4px solid #ffc107; }}
        .low-priority {{ background-color: #d4edda; border-left: 4px solid #28a745; }}
        table {{ width: 100%; border-collapse: collapse; margin: 10px 0; }}
        th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
        th {{ background-color: #f2f2f2; }}
        .metric {{ display: inline-block; margin: 10px; padding: 10px; background: #f8f9fa; border-radius: 5px; min-width: 100px; text-align: center; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>🔍 Fire Salamander - SEO Report</h1>
        <div class="score status-{self.stats['status']}">
            Score: {self.stats['overall_score']}/100 ({self.stats['status']})
        </div>
        <p>Target: {self.stats['target_url']}</p>
        <p>Generated: {self.stats['timestamp']}</p>
    </div>
    
    <div class="section">
        <h2>📊 Summary</h2>
        <div class="metric">
            <strong>{page_count}</strong><br>
            Pages Analyzed
        </div>
        <div class="metric">
            <strong>{len(self.stats['recommendations'])}</strong><br>
            Recommendations
        </div>
        <div class="metric">
            <strong>{high_priority}</strong><br>
            High Priority
        </div>
        <div class="metric">
            <strong>{medium_priority}</strong><br>
            Medium Priority
        </div>
    </div>
        """
        
        # Ajouter les recommandations
        if self.stats['recommendations']:
            html += '<div class="section"><h2>💡 Recommendations</h2>'
            for rec in self.stats['recommendations']:
                priority_class = f"{rec['priority']}-priority"
                html += f"""
                <div class="recommendation {priority_class}">
                    <strong>[{rec['priority'].upper()}] {rec['category']}</strong><br>
                    {rec['recommendation']}<br>
                    <em>Impact: {rec['impact']}</em>
                </div>
                """
            html += '</div>'
        
        # Ajouter l'analyse détaillée des pages
        for key, results in self.stats['seo_results'].items():
            if key.startswith('page_') or key == 'homepage':
                page_title = 'Homepage' if key == 'homepage' else key.replace('page_', '').replace('_', ' ').title()
                html += f'<div class="section"><h2>📄 {page_title}</h2>'
                html += f'<div class="page-analysis">'
                html += f'<p><strong>URL:</strong> {results.get("url", "N/A")}</p>'
                html += f'<p><strong>Page Score:</strong> {results.get("score", 0)}/100</p>'
                
                # Détails SEO
                if 'title' in results:
                    title_info = results['title']
                    html += f'<p><strong>Title:</strong> "{title_info.get("text", "")}" ({title_info.get("length", 0)} chars)</p>'
                
                if 'meta_description' in results:
                    desc_info = results['meta_description']
                    html += f'<p><strong>Meta Description:</strong> "{desc_info.get("text", "")}" ({desc_info.get("length", 0)} chars)</p>'
                
                if 'headings' in results:
                    headings = results['headings']
                    html += f'<p><strong>Headings:</strong> H1:{headings.get("h1", 0)}, H2:{headings.get("h2", 0)}, H3:{headings.get("h3", 0)}</p>'
                
                # Issues de la page
                for issue in results.get('issues', []):
                    issue_class = issue.get('type', 'info')
                    html += f'<div class="issue {issue_class}"><strong>{issue.get("element", "")}</strong>: {issue.get("issue", "")}</div>'
                
                html += '</div></div>'
        
        # Ajouter les tests techniques
        for key, results in self.stats['seo_results'].items():
            if key in ['site_structure', 'technical_seo']:
                section_title = key.replace('_', ' ').title()
                html += f'<div class="section"><h2>{section_title}</h2>'
                
                if isinstance(results, list):
                    for test in results:
                        status_class = test.get('status', 'info')
                        html += f"""
                        <div class="issue {status_class}">
                            <strong>{test.get('test', 'Unknown Test')}</strong><br>
                            {test.get('description', 'No description')}
                        </div>
                        """
                
                html += '</div>'
        
        html += """
</body>
</html>
        """
        
        return html

def main():
    """Point d'entrée principal"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Fire Salamander SEO Agent')
    parser.add_argument('--url', default='http://localhost:3000', help='Target URL')
    parser.add_argument('--output', default='tests/reports/seo', help='Output directory')
    parser.add_argument('--max-pages', type=int, default=10, help='Maximum pages to analyze')
    args = parser.parse_args()
    
    config = {
        'target_url': args.url,
        'report_path': args.output,
        'max_pages': args.max_pages
    }
    
    agent = SEOAgent(config)
    results = agent.run_full_seo_audit()
    
    print(f"\n🔍 SEO Audit Results:")
    print(f"Score: {results['overall_score']}/100 ({results['status']})")
    print(f"Recommendations: {len(results['recommendations'])}")
    print(f"Report: {args.output}/seo_report.html")
    
    # Exit code basé sur le score
    exit_code = 0 if results['overall_score'] >= 70 else 1
    exit(exit_code)

if __name__ == '__main__':
    main()