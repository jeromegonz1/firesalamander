#!/usr/bin/env python3
"""
Script de diagnostic de l'interface Fire Salamander avec Playwright
"""
import asyncio
from playwright.async_api import async_playwright
import json
import base64

async def test_fire_salamander_interface():
    async with async_playwright() as p:
        # Lancer le navigateur
        browser = await p.chromium.launch(headless=False)
        context = await browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            record_video_dir='videos/'
        )
        
        page = await context.new_page()
        
        # Capturer les erreurs console
        console_messages = []
        page.on('console', lambda msg: console_messages.append({
            'type': msg.type,
            'text': msg.text,
            'location': msg.location
        }))
        
        # Capturer les erreurs réseau
        network_errors = []
        page.on('requestfailed', lambda request: network_errors.append({
            'url': request.url,
            'failure': request.failure
        }))
        
        print("🔍 Ouverture de l'interface Fire Salamander...")
        
        try:
            # 1. Charger la page principale
            response = await page.goto('http://localhost:8080', wait_until='networkidle')
            print(f"✅ Page chargée - Status: {response.status}")
            
            # Attendre un peu pour que tout se charge
            await page.wait_for_timeout(2000)
            
            # 2. Prendre un screenshot
            await page.screenshot(path='fire-salamander-home.png', full_page=True)
            print("📸 Screenshot de la page d'accueil sauvegardé")
            
            # 3. Vérifier les éléments critiques
            critical_elements = {
                'navigation': '.navbar',
                'dashboard': '#dashboard-page',
                'stats-cards': '.stats-grid',
                'charts': '.charts-section',
                'main-content': '.main-content'
            }
            
            print("\n🔍 Vérification des éléments critiques:")
            for name, selector in critical_elements.items():
                element = await page.query_selector(selector)
                if element:
                    is_visible = await element.is_visible()
                    print(f"  ✅ {name}: {'Visible' if is_visible else 'Caché'}")
                else:
                    print(f"  ❌ {name}: Non trouvé!")
            
            # 4. Vérifier le chargement des ressources
            print("\n📦 Vérification des ressources:")
            
            # CSS
            css_loaded = await page.evaluate('''() => {
                const styles = document.querySelector('link[href*="styles.css"]');
                return styles ? styles.sheet !== null : false;
            }''')
            print(f"  {'✅' if css_loaded else '❌'} CSS chargé: {css_loaded}")
            
            # JavaScript
            js_loaded = await page.evaluate('''() => {
                return typeof app !== 'undefined';
            }''')
            print(f"  {'✅' if js_loaded else '❌'} JavaScript app initialisé: {js_loaded}")
            
            # Chart.js
            chartjs_loaded = await page.evaluate('''() => {
                return typeof Chart !== 'undefined';
            }''')
            print(f"  {'✅' if chartjs_loaded else '❌'} Chart.js chargé: {chartjs_loaded}")
            
            # 5. Tester la navigation
            print("\n🧭 Test de navigation:")
            
            # Cliquer sur Analyser
            analyzer_link = await page.query_selector('a[data-page="analyzer"]')
            if analyzer_link:
                await analyzer_link.click()
                await page.wait_for_timeout(1000)
                
                analyzer_visible = await page.is_visible('#analyzer-page')
                print(f"  {'✅' if analyzer_visible else '❌'} Page Analyser accessible")
                
                await page.screenshot(path='fire-salamander-analyzer.png')
                print("  📸 Screenshot de la page Analyser sauvegardé")
            else:
                print("  ❌ Lien Analyser non trouvé")
            
            # 6. Tester l'API
            print("\n🔌 Test de l'API:")
            api_health = await page.evaluate('''async () => {
                try {
                    const response = await fetch('/api/v1/health');
                    const data = await response.json();
                    return { status: response.status, data: data };
                } catch (error) {
                    return { error: error.message };
                }
            }''')
            
            if 'error' in api_health:
                print(f"  ❌ Erreur API: {api_health['error']}")
            else:
                print(f"  ✅ API santé - Status: {api_health['status']}")
                print(f"     Service: {api_health['data']['data']['status']}")
            
            # 7. Analyser les erreurs console
            if console_messages:
                print("\n⚠️  Messages console détectés:")
                for msg in console_messages:
                    if msg['type'] in ['error', 'warning']:
                        print(f"  {msg['type'].upper()}: {msg['text']}")
                        if msg['location']:
                            print(f"     Location: {msg['location']}")
            else:
                print("\n✅ Aucune erreur console détectée")
            
            # 8. Analyser les erreurs réseau
            if network_errors:
                print("\n❌ Erreurs réseau détectées:")
                for error in network_errors:
                    print(f"  URL: {error['url']}")
                    print(f"  Erreur: {error['failure']}")
            else:
                print("\n✅ Aucune erreur réseau détectée")
            
            # 9. Extraire le HTML pour analyse
            html_content = await page.content()
            with open('fire-salamander-page.html', 'w', encoding='utf-8') as f:
                f.write(html_content)
            print("\n📄 HTML de la page sauvegardé pour analyse")
            
            # 10. Inspecter les styles calculés
            print("\n🎨 Inspection des styles:")
            body_styles = await page.evaluate('''() => {
                const body = document.body;
                const styles = window.getComputedStyle(body);
                return {
                    fontFamily: styles.fontFamily,
                    backgroundColor: styles.backgroundColor,
                    color: styles.color
                };
            }''')
            print(f"  Font: {body_styles['fontFamily']}")
            print(f"  Background: {body_styles['backgroundColor']}")
            print(f"  Color: {body_styles['color']}")
            
            # 11. Vérifier si les fichiers statiques sont servis correctement
            print("\n📁 Vérification des fichiers statiques:")
            static_files = [
                '/static/styles.css',
                '/static/app.js',
                '/static/index.html'
            ]
            
            for file in static_files:
                try:
                    response = await page.request.get(f'http://localhost:8080{file}')
                    print(f"  {'✅' if response.ok else '❌'} {file}: {response.status}")
                except Exception as e:
                    print(f"  ❌ {file}: Erreur - {str(e)}")
            
        except Exception as e:
            print(f"\n❌ Erreur lors du test: {str(e)}")
            await page.screenshot(path='fire-salamander-error.png')
            print("📸 Screenshot d'erreur sauvegardé")
        
        finally:
            await browser.close()
            
            # Résumé
            print("\n" + "="*50)
            print("📊 RÉSUMÉ DU DIAGNOSTIC")
            print("="*50)
            print(f"Console messages: {len(console_messages)}")
            print(f"Network errors: {len(network_errors)}")
            print("\n📸 Screenshots sauvegardés:")
            print("  - fire-salamander-home.png")
            print("  - fire-salamander-analyzer.png")
            print("  - fire-salamander-page.html")

if __name__ == "__main__":
    asyncio.run(test_fire_salamander_interface())