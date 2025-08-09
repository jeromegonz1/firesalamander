/**
 * Fire Salamander - Technical Analysis Component Tests
 * TDD tests for TechnicalAnalysisSection component
 */

import { test, expect } from '@playwright/test';

describe('TechnicalAnalysisSection Component - TDD Tests', () => {
  
  // Mock data setup for consistent testing
  const mockTechnicalData = {
    pageAnalysis: [
      {
        url: 'https://example.com',
        statusCode: 200,
        loadTime: 1200,
        size: 45678,
        lastCrawled: '2025-01-15T10:30:00Z',
        depth: 0,
        title: {
          content: 'Example Page Title',
          length: 18,
          hasKeyword: true,
          issues: [],
          recommendations: []
        },
        metaDescription: {
          content: 'This is a good meta description',
          length: 31,
          hasKeyword: false,
          issues: ['meta-too-short'],
          recommendations: ['Rallonger la meta description']
        },
        headings: {
          h1: ['Main Heading'],
          h2: ['Section 1', 'Section 2'],
          h3: ['Subsection'],
          structure: 'good',
          issues: [],
          recommendations: []
        },
        canonical: 'https://example.com',
        robots: 'index,follow',
        schema: [
          { type: 'Organization', valid: true, errors: [] }
        ],
        openGraph: {
          'og:title': 'Example Title',
          'og:description': 'Example Description'
        },
        images: [
          {
            src: '/image1.jpg',
            alt: 'Good image description',
            size: 12345,
            issues: [],
            recommendations: []
          },
          {
            src: '/image2.jpg',
            alt: '',
            size: 1200000, // Large image
            issues: ['image-no-alt', 'image-large-size'],
            recommendations: ['Ajouter un texte alternatif', 'Optimiser la taille']
          }
        ],
        links: {
          internal: 15,
          external: 5,
          broken: [
            {
              url: 'https://broken.com',
              status: 404,
              anchorText: 'Broken link',
              position: 'content'
            }
          ],
          nofollow: 2,
          totalLinks: 20
        },
        performance: {
          fcp: 800,
          lcp: 1200,
          cls: 0.05,
          ttfb: 200
        },
        issueCount: 3,
        issueTypes: ['meta-too-short', 'image-no-alt', 'image-large-size'],
        overallHealth: 'warning'
      },
      {
        url: 'https://example.com/page2',
        statusCode: 404,
        loadTime: 0,
        size: 0,
        title: {
          content: '',
          length: 0,
          hasKeyword: false,
          issues: ['title-missing'],
          recommendations: ['Ajouter un titre']
        },
        issueCount: 5,
        overallHealth: 'error'
      }
    ],
    globalIssues: {
      duplicateContent: [
        {
          pages: ['https://example.com/page1', 'https://example.com/page2'],
          similarity: 85,
          level: 'high',
          affectedElements: ['title', 'content']
        }
      ],
      duplicateTitles: [
        {
          title: 'Duplicate Title',
          pages: ['https://example.com/dup1', 'https://example.com/dup2'],
          count: 2
        }
      ],
      brokenLinks: [
        {
          from: 'https://example.com',
          to: 'https://broken-site.com',
          status: 404,
          anchorText: 'Old link',
          linkType: 'external',
          position: 'content',
          lastChecked: '2025-01-15T10:30:00Z'
        }
      ],
      orphanPages: ['https://example.com/orphan'],
      redirectChains: [
        {
          chain: ['https://old.com', 'https://redirect.com', 'https://final.com'],
          finalUrl: 'https://final.com',
          totalHops: 2,
          statusCodes: [301, 302]
        }
      ]
    },
    crawlability: {
      robotsTxt: {
        exists: true,
        valid: true,
        userAgents: ['*'],
        disallowedPaths: ['/admin', '/private'],
        sitemapUrls: ['https://example.com/sitemap.xml'],
        issues: []
      },
      sitemap: {
        exists: true,
        url: 'https://example.com/sitemap.xml',
        valid: true,
        format: 'xml',
        pagesInSitemap: 25,
        images: 10,
        videos: 2,
        issues: []
      },
      crawlBudget: {
        totalPages: 30,
        crawlablePages: 25,
        blockedPages: 5,
        pagesPerLevel: { 0: 1, 1: 10, 2: 14, 3: 5 },
        averageCrawlTime: 1.2,
        crawlEfficiency: 83.3
      }
    },
    metrics: {
      totalPages: 30,
      pagesWithIssues: 12,
      avgLoadTime: 1500,
      avgPageSize: 35000,
      totalIssues: 45,
      issuesByType: {
        'title-missing': 5,
        'meta-too-short': 8,
        'image-no-alt': 12,
        'link-broken': 3
      },
      healthScore: 75
    }
  };

  test.beforeEach(async ({ page }) => {
    // Mock the API responses for consistent testing
    await page.route('**/api/v1/analysis/*/technical', async route => {
      await route.fulfill({
        json: { data: mockTechnicalData }
      });
    });
  });

  describe('Component Rendering', () => {
    
    test('should render main technical analysis sections', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Verify main sections are present
      await expect(page.locator('[data-testid="technical-analysis-section"]')).toBeVisible();
      await expect(page.locator('[data-testid="metrics-overview"]')).toBeVisible();
      await expect(page.locator('[data-testid="page-analysis-section"]')).toBeVisible();
      await expect(page.locator('[data-testid="global-issues-section"]')).toBeVisible();
      await expect(page.locator('[data-testid="crawlability-section"]')).toBeVisible();
    });
    
    test('should display correct metrics overview', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Check metrics cards
      await expect(page.locator('[data-testid="total-pages-metric"]')).toContainText('30');
      await expect(page.locator('[data-testid="pages-with-issues-metric"]')).toContainText('12');
      await expect(page.locator('[data-testid="health-score-metric"]')).toContainText('75');
      await expect(page.locator('[data-testid="total-issues-metric"]')).toContainText('45');
    });
    
    test('should render health score gauge correctly', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const healthGauge = page.locator('[data-testid="health-score-gauge"]');
      await expect(healthGauge).toBeVisible();
      
      // Check gauge color based on score (75 = warning/orange)
      await expect(healthGauge.locator('[data-testid="gauge-fill"]')).toHaveClass(/warning|orange/);
      await expect(healthGauge.locator('[data-testid="gauge-score"]')).toContainText('75');
    });
  });

  describe('Page Analysis Table', () => {
    
    test('should render page analysis table with all columns', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const table = page.locator('[data-testid="page-analysis-table"]');
      await expect(table).toBeVisible();
      
      // Check table headers
      await expect(table.locator('th[data-testid="url-header"]')).toContainText('URL');
      await expect(table.locator('th[data-testid="status-header"]')).toContainText('Status');
      await expect(table.locator('th[data-testid="load-time-header"]')).toContainText('Load Time');
      await expect(table.locator('th[data-testid="size-header"]')).toContainText('Size');
      await expect(table.locator('th[data-testid="issues-header"]')).toContainText('Issues');
      await expect(table.locator('th[data-testid="health-header"]')).toContainText('Health');
    });
    
    test('should display page data correctly in table rows', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // First page (success)
      const firstRow = page.locator('[data-testid="page-row-0"]');
      await expect(firstRow.locator('[data-testid="page-url"]')).toContainText('https://example.com');
      await expect(firstRow.locator('[data-testid="page-status"]')).toContainText('200');
      await expect(firstRow.locator('[data-testid="page-load-time"]')).toContainText('1.2s');
      await expect(firstRow.locator('[data-testid="page-size"]')).toContainText('44.6 KB');
      await expect(firstRow.locator('[data-testid="page-issues"]')).toContainText('3');
      await expect(firstRow.locator('[data-testid="page-health"]')).toHaveClass(/warning/);
      
      // Second page (error)
      const secondRow = page.locator('[data-testid="page-row-1"]');
      await expect(secondRow.locator('[data-testid="page-status"]')).toContainText('404');
      await expect(secondRow.locator('[data-testid="page-health"]')).toHaveClass(/error/);
    });
    
    test('should expand page details on row click', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Click on first page row
      await page.locator('[data-testid="page-row-0"]').click();
      
      // Verify details panel appears
      const detailsPanel = page.locator('[data-testid="page-details-0"]');
      await expect(detailsPanel).toBeVisible();
      
      // Check sections in details
      await expect(detailsPanel.locator('[data-testid="seo-elements-section"]')).toBeVisible();
      await expect(detailsPanel.locator('[data-testid="technical-elements-section"]')).toBeVisible();
      await expect(detailsPanel.locator('[data-testid="images-section"]')).toBeVisible();
      await expect(detailsPanel.locator('[data-testid="links-section"]')).toBeVisible();
      await expect(detailsPanel.locator('[data-testid="performance-section"]')).toBeVisible();
    });
    
    test('should show SEO elements details correctly', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      await page.locator('[data-testid="page-row-0"]').click();
      
      const seoSection = page.locator('[data-testid="seo-elements-section"]');
      
      // Title details
      await expect(seoSection.locator('[data-testid="title-content"]')).toContainText('Example Page Title');
      await expect(seoSection.locator('[data-testid="title-length"]')).toContainText('18');
      
      // Meta description details
      await expect(seoSection.locator('[data-testid="meta-content"]')).toContainText('This is a good meta description');
      await expect(seoSection.locator('[data-testid="meta-issues"]')).toContainText('meta-too-short');
      
      // Headings
      await expect(seoSection.locator('[data-testid="h1-list"]')).toContainText('Main Heading');
      await expect(seoSection.locator('[data-testid="h2-list"]')).toContainText('Section 1');
      await expect(seoSection.locator('[data-testid="heading-structure"]')).toContainText('good');
    });
    
    test('should show images analysis correctly', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      await page.locator('[data-testid="page-row-0"]').click();
      
      const imagesSection = page.locator('[data-testid="images-section"]');
      await expect(imagesSection).toBeVisible();
      
      // First image (good)
      const image1 = imagesSection.locator('[data-testid="image-0"]');
      await expect(image1.locator('[data-testid="image-src"]')).toContainText('/image1.jpg');
      await expect(image1.locator('[data-testid="image-alt"]')).toContainText('Good image description');
      await expect(image1.locator('[data-testid="image-size"]')).toContainText('12.1 KB');
      await expect(image1.locator('[data-testid="image-issues"]')).toHaveCount(0);
      
      // Second image (with issues)
      const image2 = imagesSection.locator('[data-testid="image-1"]');
      await expect(image2.locator('[data-testid="image-issues"]')).toContainText('image-no-alt');
      await expect(image2.locator('[data-testid="image-issues"]')).toContainText('image-large-size');
    });
  });

  describe('Filtering and Sorting', () => {
    
    test('should filter pages by status code', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Apply 404 filter
      await page.selectOption('[data-testid="status-filter"]', '404');
      
      // Should only show 404 pages
      await expect(page.locator('[data-testid^="page-row-"]')).toHaveCount(1);
      await expect(page.locator('[data-testid="page-row-0"] [data-testid="page-status"]')).toContainText('404');
    });
    
    test('should filter pages by issue type', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Apply title issues filter
      await page.selectOption('[data-testid="issue-filter"]', 'title-issues');
      
      // Should show only pages with title issues
      const visibleRows = page.locator('[data-testid^="page-row-"]');
      const count = await visibleRows.count();
      expect(count).toBeGreaterThan(0);
      
      // Verify all shown pages have title issues
      for (let i = 0; i < count; i++) {
        const row = visibleRows.nth(i);
        await row.click();
        await expect(page.locator('[data-testid="title-issues"]')).toBeVisible();
        await row.click(); // Close details
      }
    });
    
    test('should sort pages by load time', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Sort by load time descending
      await page.click('[data-testid="sort-load-time"]');
      
      // Verify sort order
      const firstRowTime = await page.locator('[data-testid="page-row-0"] [data-testid="page-load-time"]').textContent();
      const secondRowTime = await page.locator('[data-testid="page-row-1"] [data-testid="page-load-time"]').textContent();
      
      // Parse times and verify descending order
      const time1 = parseFloat(firstRowTime?.replace('s', '') || '0');
      const time2 = parseFloat(secondRowTime?.replace('s', '') || '0');
      expect(time1).toBeGreaterThanOrEqual(time2);
    });
    
    test('should search pages by URL', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Search for specific URL
      await page.fill('[data-testid="page-search"]', 'page2');
      
      // Should filter to matching URLs
      await expect(page.locator('[data-testid^="page-row-"]')).toHaveCount(1);
      await expect(page.locator('[data-testid="page-row-0"] [data-testid="page-url"]')).toContainText('page2');
    });
  });

  describe('Global Issues Section', () => {
    
    test('should display duplicate content issues', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const duplicateSection = page.locator('[data-testid="duplicate-content-section"]');
      await expect(duplicateSection).toBeVisible();
      
      // Check duplicate content item
      const duplicateItem = duplicateSection.locator('[data-testid="duplicate-item-0"]');
      await expect(duplicateItem.locator('[data-testid="similarity"]')).toContainText('85%');
      await expect(duplicateItem.locator('[data-testid="similarity-level"]')).toContainText('high');
      await expect(duplicateItem.locator('[data-testid="affected-pages"]')).toContainText('2 pages');
    });
    
    test('should display broken links analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const brokenLinksSection = page.locator('[data-testid="broken-links-section"]');
      await expect(brokenLinksSection).toBeVisible();
      
      // Check broken link item
      const brokenLink = brokenLinksSection.locator('[data-testid="broken-link-0"]');
      await expect(brokenLink.locator('[data-testid="broken-url"]')).toContainText('https://broken-site.com');
      await expect(brokenLink.locator('[data-testid="broken-status"]')).toContainText('404');
      await expect(brokenLink.locator('[data-testid="broken-from"]')).toContainText('https://example.com');
    });
    
    test('should display redirect chains analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const redirectSection = page.locator('[data-testid="redirect-chains-section"]');
      await expect(redirectSection).toBeVisible();
      
      // Check redirect chain item
      const redirectChain = redirectSection.locator('[data-testid="redirect-chain-0"]');
      await expect(redirectChain.locator('[data-testid="chain-hops"]')).toContainText('2 hops');
      await expect(redirectChain.locator('[data-testid="final-url"]')).toContainText('https://final.com');
      
      // Check chain visualization
      const chainSteps = redirectChain.locator('[data-testid="chain-step"]');
      await expect(chainSteps).toHaveCount(3); // 3 URLs in chain
    });
    
    test('should show orphan pages', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const orphanSection = page.locator('[data-testid="orphan-pages-section"]');
      await expect(orphanSection).toBeVisible();
      await expect(orphanSection.locator('[data-testid="orphan-page-0"]')).toContainText('https://example.com/orphan');
    });
  });

  describe('Crawlability Section', () => {
    
    test('should display robots.txt analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const robotsSection = page.locator('[data-testid="robots-txt-section"]');
      await expect(robotsSection).toBeVisible();
      
      // Check robots.txt status
      await expect(robotsSection.locator('[data-testid="robots-exists"]')).toContainText('✓ Exists');
      await expect(robotsSection.locator('[data-testid="robots-valid"]')).toContainText('✓ Valid');
      
      // Check disallowed paths
      await expect(robotsSection.locator('[data-testid="disallowed-paths"]')).toContainText('/admin');
      await expect(robotsSection.locator('[data-testid="disallowed-paths"]')).toContainText('/private');
      
      // Check sitemap URLs
      await expect(robotsSection.locator('[data-testid="sitemap-urls"]')).toContainText('sitemap.xml');
    });
    
    test('should display sitemap analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const sitemapSection = page.locator('[data-testid="sitemap-section"]');
      await expect(sitemapSection).toBeVisible();
      
      // Check sitemap metrics
      await expect(sitemapSection.locator('[data-testid="sitemap-exists"]')).toContainText('✓ Exists');
      await expect(sitemapSection.locator('[data-testid="sitemap-format"]')).toContainText('XML');
      await expect(sitemapSection.locator('[data-testid="pages-in-sitemap"]')).toContainText('25');
      await expect(sitemapSection.locator('[data-testid="images-in-sitemap"]')).toContainText('10');
    });
    
    test('should display crawl budget analysis', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const budgetSection = page.locator('[data-testid="crawl-budget-section"]');
      await expect(budgetSection).toBeVisible();
      
      // Check crawl budget metrics
      await expect(budgetSection.locator('[data-testid="total-pages"]')).toContainText('30');
      await expect(budgetSection.locator('[data-testid="crawlable-pages"]')).toContainText('25');
      await expect(budgetSection.locator('[data-testid="blocked-pages"]')).toContainText('5');
      await expect(budgetSection.locator('[data-testid="crawl-efficiency"]')).toContainText('83.3%');
      
      // Check depth distribution chart
      await expect(budgetSection.locator('[data-testid="depth-distribution"]')).toBeVisible();
    });
  });

  describe('Export and Actions', () => {
    
    test('should show export options', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      const exportButton = page.locator('[data-testid="export-technical"]');
      await expect(exportButton).toBeVisible();
      
      await exportButton.click();
      
      // Check export options
      await expect(page.locator('[data-testid="export-pdf"]')).toBeVisible();
      await expect(page.locator('[data-testid="export-excel"]')).toBeVisible();
      await expect(page.locator('[data-testid="export-csv"]')).toBeVisible();
    });
    
    test('should allow bulk actions on selected pages', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Select multiple pages
      await page.check('[data-testid="page-checkbox-0"]');
      await page.check('[data-testid="page-checkbox-1"]');
      
      // Check bulk actions appear
      await expect(page.locator('[data-testid="bulk-actions"]')).toBeVisible();
      await expect(page.locator('[data-testid="reanalyze-selected"]')).toBeVisible();
      await expect(page.locator('[data-testid="export-selected"]')).toBeVisible();
    });
  });

  describe('Real-time Updates', () => {
    
    test('should refresh data when analysis is updated', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Mock updated data
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({
          json: { 
            data: { 
              ...mockTechnicalData,
              metrics: {
                ...mockTechnicalData.metrics,
                healthScore: 85 // Updated score
              }
            }
          }
        });
      });
      
      // Click refresh
      await page.click('[data-testid="refresh-analysis"]');
      
      // Check updated data
      await expect(page.locator('[data-testid="health-score-metric"]')).toContainText('85');
    });
    
    test('should show analysis progress for ongoing crawls', async ({ page }) => {
      // Mock ongoing analysis
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({
          json: { 
            data: { 
              ...mockTechnicalData,
              status: {
                ...mockTechnicalData.status,
                crawlStatus: 'in-progress'
              }
            }
          }
        });
      });
      
      await page.goto('/analysis/1/technical');
      
      // Check progress indicator
      await expect(page.locator('[data-testid="analysis-progress"]')).toBeVisible();
      await expect(page.locator('[data-testid="crawl-status"]')).toContainText('in-progress');
    });
  });

  describe('Error Handling', () => {
    
    test('should handle API errors gracefully', async ({ page }) => {
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({ status: 500 });
      });
      
      await page.goto('/analysis/1/technical');
      
      await expect(page.locator('[data-testid="technical-analysis-error"]')).toBeVisible();
      await expect(page.locator('[data-testid="error-message"]')).toContainText('Unable to load technical analysis');
      await expect(page.locator('[data-testid="retry-button"]')).toBeVisible();
    });
    
    test('should show empty state for no pages', async ({ page }) => {
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({
          json: { 
            data: { 
              ...mockTechnicalData,
              pageAnalysis: []
            }
          }
        });
      });
      
      await page.goto('/analysis/1/technical');
      
      await expect(page.locator('[data-testid="no-pages-message"]')).toBeVisible();
      await expect(page.locator('[data-testid="start-analysis-button"]')).toBeVisible();
    });
  });

  describe('Performance and Accessibility', () => {
    
    test('should load large datasets efficiently', async ({ page }) => {
      // Mock large dataset
      const largePageList = Array.from({ length: 1000 }, (_, i) => ({
        url: `https://example.com/page${i}`,
        statusCode: 200,
        loadTime: Math.random() * 3000,
        size: Math.random() * 100000,
        issueCount: Math.floor(Math.random() * 10),
        overallHealth: i % 3 === 0 ? 'good' : i % 3 === 1 ? 'warning' : 'error'
      }));
      
      await page.route('**/api/v1/analysis/*/technical', async route => {
        await route.fulfill({
          json: { 
            data: { 
              ...mockTechnicalData,
              pageAnalysis: largePageList
            }
          }
        });
      });
      
      const startTime = Date.now();
      await page.goto('/analysis/1/technical');
      
      // Should load within reasonable time
      await expect(page.locator('[data-testid="page-analysis-table"]')).toBeVisible();
      const loadTime = Date.now() - startTime;
      expect(loadTime).toBeLessThan(5000);
      
      // Should implement virtualization or pagination
      const visibleRows = await page.locator('[data-testid^="page-row-"]').count();
      expect(visibleRows).toBeLessThanOrEqual(100); // Should not render all 1000 rows
    });
    
    test('should be keyboard accessible', async ({ page }) => {
      await page.goto('/analysis/1/technical');
      
      // Test keyboard navigation through table
      await page.keyboard.press('Tab');
      await page.keyboard.press('Enter');
      
      // Should open page details
      await expect(page.locator('[data-testid="page-details-0"]')).toBeVisible();
      
      // Test filter navigation
      await page.keyboard.press('Tab');
      await page.keyboard.press('ArrowDown');
      
      const activeElement = await page.evaluate(() => document.activeElement?.getAttribute('data-testid'));
      expect(activeElement).toMatch(/filter|search/);
    });
  });
});