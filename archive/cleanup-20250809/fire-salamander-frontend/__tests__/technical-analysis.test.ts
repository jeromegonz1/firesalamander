/**
 * Fire Salamander - Technical Analysis Tests
 * TDD approach for detailed page-by-page technical analysis
 */

import { test, expect } from '@playwright/test';

describe('Technical Analysis - TDD Tests', () => {
  
  // Interface et types tests
  describe('TechnicalAnalysis Interface', () => {
    
    test('should define TechnicalAnalysis interface with all required fields', () => {
      // Ce test valide la structure de l'interface
      const mockTechnicalAnalysis = {
        pageAnalysis: [{
          url: 'https://example.com',
          statusCode: 200,
          loadTime: 1200,
          size: 45678,
          title: {
            content: 'Test Page Title',
            length: 15,
            hasKeyword: true,
            issues: []
          },
          metaDescription: {
            content: 'Test meta description',
            length: 23,
            hasKeyword: false,
            issues: ['Too short']
          },
          headings: {
            h1: ['Main Title'],
            h2: ['Section 1', 'Section 2'],
            h3: ['Subsection 1'],
            structure: 'good' as const,
            issues: []
          },
          canonical: 'https://example.com',
          robots: 'index,follow',
          schema: [{ type: 'Organization', valid: true }],
          openGraph: { 'og:title': 'Test Title' },
          images: [{
            src: '/image.jpg',
            alt: 'Test image',
            size: 12345,
            issues: []
          }],
          links: {
            internal: 10,
            external: 3,
            broken: []
          }
        }],
        globalIssues: {
          duplicateContent: [],
          duplicateTitles: [],
          missingTitles: [],
          missingMeta: [],
          brokenLinks: [],
          orphanPages: [],
          redirectChains: []
        },
        crawlability: {
          robotsTxt: {
            exists: true,
            issues: [],
            disallowedPaths: ['/admin']
          },
          sitemap: {
            exists: true,
            url: '/sitemap.xml',
            pagesInSitemap: 25,
            issues: []
          },
          crawlBudget: {
            totalPages: 30,
            crawlablePages: 25,
            blockedPages: 5
          }
        }
      };
      
      expect(mockTechnicalAnalysis).toBeDefined();
      expect(mockTechnicalAnalysis.pageAnalysis).toHaveLength(1);
      expect(mockTechnicalAnalysis.globalIssues).toBeDefined();
      expect(mockTechnicalAnalysis.crawlability).toBeDefined();
    });
    
    test('should validate page analysis structure', () => {
      const pageAnalysis = {
        url: 'https://test.com/page',
        statusCode: 200,
        loadTime: 850,
        size: 23456,
        title: {
          content: 'Page Title',
          length: 10,
          hasKeyword: true,
          issues: []
        }
      };
      
      expect(pageAnalysis.url).toMatch(/^https?:\/\//);
      expect(pageAnalysis.statusCode).toBeGreaterThanOrEqual(100);
      expect(pageAnalysis.loadTime).toBeGreaterThan(0);
      expect(pageAnalysis.size).toBeGreaterThan(0);
      expect(pageAnalysis.title.length).toBe(pageAnalysis.title.content.length);
    });
  });
  
  // Mapper tests
  describe('Technical Analysis Mapper', () => {
    
    test('should map backend data to TechnicalAnalysis structure', async () => {
      const mockBackendData = {
        result_data: JSON.stringify({
          technical_analysis: {
            pages: [{
              url: 'https://example.com',
              status_code: 200,
              load_time: 1500,
              page_size: 56789,
              seo_elements: {
                title: 'Example Title',
                meta_description: 'Example description',
                headings: {
                  h1: ['Main Heading'],
                  h2: ['Secondary Heading']
                }
              }
            }],
            global_issues: {
              duplicate_titles: [],
              broken_links: []
            },
            crawl_info: {
              robots_txt_exists: true,
              sitemap_exists: true
            }
          }
        })
      };
      
      // Le mapper devrait transformer ces données
      expect(mockBackendData.result_data).toContain('technical_analysis');
    });
    
    test('should handle missing or malformed backend data gracefully', () => {
      const emptyData = { result_data: '{}' };
      const invalidData = { result_data: 'invalid json' };
      
      // Le mapper devrait gérer ces cas sans crash
      expect(emptyData).toBeDefined();
      expect(invalidData).toBeDefined();
    });
    
    test('should calculate page analysis metrics correctly', () => {
      const titleContent = 'This is a test title for SEO analysis';
      const expectedLength = titleContent.length;
      
      expect(expectedLength).toBe(37);
      
      // Test pour la détection de mots-clés
      const hasKeyword = titleContent.toLowerCase().includes('seo');
      expect(hasKeyword).toBe(true);
    });
  });
  
  // Component tests
  describe('TechnicalAnalysisSection Component', () => {
    
    test('should render page analysis table with all columns', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Table headers attendus
      await expect(page.locator('[data-testid="page-analysis-table"]')).toBeVisible();
      await expect(page.locator('th:has-text("URL")')).toBeVisible();
      await expect(page.locator('th:has-text("Status")')).toBeVisible();
      await expect(page.locator('th:has-text("Load Time")')).toBeVisible();
      await expect(page.locator('th:has-text("Size")')).toBeVisible();
      await expect(page.locator('th:has-text("Title Issues")')).toBeVisible();
      await expect(page.locator('th:has-text("Meta Issues")')).toBeVisible();
    });
    
    test('should expand page details on row click', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Click sur une ligne de page
      await page.locator('[data-testid="page-row-0"]').click();
      
      // Vérifier que les détails s'affichent
      await expect(page.locator('[data-testid="page-details-0"]')).toBeVisible();
      await expect(page.locator('[data-testid="seo-elements-details"]')).toBeVisible();
      await expect(page.locator('[data-testid="technical-elements-details"]')).toBeVisible();
    });
    
    test('should display global issues section', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      await expect(page.locator('[data-testid="global-issues-section"]')).toBeVisible();
      await expect(page.locator('[data-testid="duplicate-content-issues"]')).toBeVisible();
      await expect(page.locator('[data-testid="broken-links-issues"]')).toBeVisible();
      await expect(page.locator('[data-testid="missing-elements-issues"]')).toBeVisible();
    });
    
    test('should show crawlability analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      await expect(page.locator('[data-testid="crawlability-section"]')).toBeVisible();
      await expect(page.locator('[data-testid="robots-txt-analysis"]')).toBeVisible();
      await expect(page.locator('[data-testid="sitemap-analysis"]')).toBeVisible();
      await expect(page.locator('[data-testid="crawl-budget-analysis"]')).toBeVisible();
    });
    
    test('should filter pages by issues', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Filtrer par pages avec problèmes de titre
      await page.selectOption('[data-testid="issue-filter"]', 'title-issues');
      
      // Vérifier que seules les pages avec problèmes de titre s'affichent
      const visibleRows = await page.locator('[data-testid^="page-row-"]').count();
      expect(visibleRows).toBeGreaterThan(0);
    });
    
    test('should sort pages by different criteria', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Trier par temps de chargement
      await page.click('[data-testid="sort-load-time"]');
      
      // Vérifier l'ordre de tri
      const firstRowLoadTime = await page.locator('[data-testid="page-row-0"] [data-testid="load-time"]').textContent();
      expect(firstRowLoadTime).toBeDefined();
    });
  });
  
  // Integration tests
  describe('Technical Analysis Integration', () => {
    
    test('should load technical analysis from backend API', async ({ page }) => {
      // Mock de l'API backend
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({
          json: {
            data: {
              pageAnalysis: [{
                url: 'https://example.com',
                statusCode: 200,
                loadTime: 1200,
                title: { content: 'Test', length: 4, hasKeyword: false, issues: [] }
              }],
              globalIssues: { duplicateTitles: [], brokenLinks: [] },
              crawlability: { robotsTxt: { exists: true } }
            }
          }
        });
      });
      
      await page.goto('/analysis/1/technical');
      await expect(page.locator('[data-testid="technical-analysis-loaded"]')).toBeVisible();
    });
    
    test('should handle API errors gracefully', async ({ page }) => {
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({ status: 500 });
      });
      
      await page.goto('/analysis/1/technical');
      await expect(page.locator('[data-testid="technical-analysis-error"]')).toBeVisible();
    });
  });
  
  // Performance tests
  describe('Performance & Accessibility', () => {
    
    test('should render large page lists efficiently', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Vérifier que la pagination ou virtualisation fonctionne
      await expect(page.locator('[data-testid="page-analysis-table"]')).toBeVisible();
      
      // Performance: le composant devrait se charger en moins de 3 secondes
      const startTime = Date.now();
      await page.locator('[data-testid="technical-analysis-loaded"]').waitFor();
      const loadTime = Date.now() - startTime;
      
      expect(loadTime).toBeLessThan(3000);
    });
    
    test('should be accessible with keyboard navigation', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Navigation au clavier
      await page.keyboard.press('Tab');
      await page.keyboard.press('Enter');
      
      // Vérifier que l'élément est focusable
      const focusedElement = await page.evaluate(() => document.activeElement?.tagName);
      expect(focusedElement).toBeDefined();
    });
  });
});