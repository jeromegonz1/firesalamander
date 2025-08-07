/**
 * Fire Salamander - Technical Analysis Mapper Tests
 * TDD tests for backend data mapping to TechnicalAnalysis
 */

import { describe, test, expect, beforeEach } from '@jest/globals';
import {
  TechnicalAnalysis,
  PageAnalysis,
  IssueType,
  PageStatus,
  HeadingStructure,
  SimilarityLevel
} from '@/types/technical-analysis';

describe('Technical Analysis Mapper - TDD Tests', () => {
  
  let mockBackendData: any;
  
  beforeEach(() => {
    mockBackendData = {
      id: 1,
      url: 'https://example.com',
      result_data: JSON.stringify({
        technical_analysis: {
          pages: [
            {
              url: 'https://example.com',
              status_code: 200,
              load_time: 1200,
              page_size: 45678,
              last_crawled: '2025-01-15T10:30:00Z',
              depth: 0,
              seo_elements: {
                title: {
                  content: 'Example Title',
                  length: 13,
                  has_keyword: true,
                  issues: ['title-too-short']
                },
                meta_description: {
                  content: 'This is an example meta description for testing purposes',
                  length: 57,
                  has_keyword: false,
                  issues: []
                },
                headings: {
                  h1: ['Main Heading'],
                  h2: ['Section 1', 'Section 2'],
                  h3: ['Subsection 1'],
                  structure: 'good',
                  issues: []
                }
              },
              technical_elements: {
                canonical: 'https://example.com',
                robots: 'index,follow',
                meta_robots: 'index,follow',
                lang: 'fr',
                hreflang: [
                  { lang: 'en', url: 'https://example.com/en' }
                ]
              },
              schema: [
                {
                  type: 'Organization',
                  valid: true,
                  properties: { name: 'Example Company' }
                }
              ],
              open_graph: {
                'og:title': 'Example Title',
                'og:description': 'Example description'
              },
              images: [
                {
                  src: '/image1.jpg',
                  alt: 'Example image',
                  size: 12345,
                  dimensions: { width: 800, height: 600 },
                  issues: []
                },
                {
                  src: '/image2.jpg',
                  alt: '',
                  size: 67890,
                  issues: ['image-no-alt']
                }
              ],
              links: {
                internal: 15,
                external: 5,
                broken: [
                  {
                    url: 'https://broken.com',
                    status: 404,
                    anchor_text: 'Broken link',
                    position: 'content'
                  }
                ],
                nofollow: 2,
                total: 20
              },
              performance: {
                fcp: 800,
                lcp: 1200,
                cls: 0.05,
                ttfb: 200
              }
            },
            {
              url: 'https://example.com/page2',
              status_code: 404,
              load_time: 0,
              page_size: 0,
              last_crawled: '2025-01-15T10:35:00Z',
              depth: 1,
              seo_elements: {
                title: { content: '', length: 0, has_keyword: false, issues: ['title-missing'] },
                meta_description: { content: '', length: 0, has_keyword: false, issues: ['meta-missing'] },
                headings: { h1: [], h2: [], h3: [], structure: 'bad', issues: ['h1-missing'] }
              }
            }
          ],
          global_issues: {
            duplicate_content: [
              {
                pages: ['https://example.com/page1', 'https://example.com/page2'],
                similarity: 85,
                level: 'high',
                affected_elements: ['title', 'content']
              }
            ],
            duplicate_titles: [
              {
                title: 'Duplicate Title',
                pages: ['https://example.com/dup1', 'https://example.com/dup2'],
                count: 2
              }
            ],
            broken_links: [
              {
                from: 'https://example.com',
                to: 'https://broken-site.com',
                status: 404,
                anchor_text: 'Old link',
                link_type: 'external',
                position: 'content',
                last_checked: '2025-01-15T10:30:00Z'
              }
            ],
            orphan_pages: ['https://example.com/orphan'],
            redirect_chains: [
              {
                chain: ['https://old.com', 'https://redirect.com', 'https://final.com'],
                final_url: 'https://final.com',
                total_hops: 2,
                status_codes: [301, 302]
              }
            ]
          },
          crawl_info: {
            robots_txt: {
              exists: true,
              valid: true,
              url: 'https://example.com/robots.txt',
              user_agents: ['*'],
              disallowed_paths: ['/admin', '/private'],
              sitemap_urls: ['https://example.com/sitemap.xml'],
              issues: []
            },
            sitemap: {
              exists: true,
              url: 'https://example.com/sitemap.xml',
              valid: true,
              format: 'xml',
              pages_in_sitemap: 25,
              last_modified: '2025-01-10T00:00:00Z',
              issues: []
            },
            crawl_budget: {
              total_pages: 30,
              crawlable_pages: 25,
              blocked_pages: 5,
              average_crawl_time: 1.2
            }
          },
          analysis_config: {
            max_depth: 3,
            respect_robots: true,
            user_agent: 'Fire-Salamander-Bot/1.0',
            crawl_delay: 1000,
            max_pages: 100
          },
          metrics: {
            total_pages: 30,
            pages_with_issues: 12,
            avg_load_time: 1500,
            avg_page_size: 35000,
            total_issues: 45,
            health_score: 75
          }
        }
      })
    };
  });
  
  describe('mapBackendToTechnicalAnalysis', () => {
    
    test('should map complete backend data correctly', () => {
      // Cette fonction sera implémentée après les tests
      const mockMapper = (data: any): TechnicalAnalysis => {
        const resultData = JSON.parse(data.result_data);
        const techAnalysis = resultData.technical_analysis;
        
        return {
          pageAnalysis: techAnalysis.pages.map((page: any) => ({
            url: page.url,
            statusCode: page.status_code,
            loadTime: page.load_time,
            size: page.page_size,
            lastCrawled: page.last_crawled,
            depth: page.depth || 0,
            title: {
              content: page.seo_elements?.title?.content || '',
              length: page.seo_elements?.title?.length || 0,
              hasKeyword: page.seo_elements?.title?.has_keyword || false,
              issues: page.seo_elements?.title?.issues || []
            },
            metaDescription: {
              content: page.seo_elements?.meta_description?.content || '',
              length: page.seo_elements?.meta_description?.length || 0,
              hasKeyword: page.seo_elements?.meta_description?.has_keyword || false,
              issues: page.seo_elements?.meta_description?.issues || []
            },
            headings: {
              h1: page.seo_elements?.headings?.h1 || [],
              h2: page.seo_elements?.headings?.h2 || [],
              h3: page.seo_elements?.headings?.h3 || [],
              structure: page.seo_elements?.headings?.structure || 'bad',
              issues: page.seo_elements?.headings?.issues || []
            },
            canonical: page.technical_elements?.canonical || '',
            robots: page.technical_elements?.robots || '',
            schema: page.schema || [],
            openGraph: page.open_graph || {},
            images: page.images || [],
            links: {
              internal: page.links?.internal || 0,
              external: page.links?.external || 0,
              broken: page.links?.broken || []
            }
          })),
          globalIssues: {
            duplicateContent: techAnalysis.global_issues?.duplicate_content || [],
            duplicateTitles: techAnalysis.global_issues?.duplicate_titles || [],
            duplicateMeta: [],
            missingTitles: [],
            missingMeta: [],
            brokenLinks: techAnalysis.global_issues?.broken_links || [],
            orphanPages: techAnalysis.global_issues?.orphan_pages || [],
            redirectChains: techAnalysis.global_issues?.redirect_chains || [],
            largePages: [],
            slowPages: []
          },
          crawlability: {
            robotsTxt: {
              exists: techAnalysis.crawl_info?.robots_txt?.exists || false,
              valid: techAnalysis.crawl_info?.robots_txt?.valid || false,
              userAgents: techAnalysis.crawl_info?.robots_txt?.user_agents || [],
              disallowedPaths: techAnalysis.crawl_info?.robots_txt?.disallowed_paths || [],
              allowedPaths: [],
              sitemapUrls: techAnalysis.crawl_info?.robots_txt?.sitemap_urls || [],
              issues: techAnalysis.crawl_info?.robots_txt?.issues || []
            },
            sitemap: {
              exists: techAnalysis.crawl_info?.sitemap?.exists || false,
              url: techAnalysis.crawl_info?.sitemap?.url || '',
              valid: techAnalysis.crawl_info?.sitemap?.valid || false,
              format: techAnalysis.crawl_info?.sitemap?.format || 'unknown',
              pagesInSitemap: techAnalysis.crawl_info?.sitemap?.pages_in_sitemap || 0,
              images: 0,
              videos: 0,
              issues: techAnalysis.crawl_info?.sitemap?.issues || []
            },
            crawlBudget: {
              totalPages: techAnalysis.crawl_info?.crawl_budget?.total_pages || 0,
              crawlablePages: techAnalysis.crawl_info?.crawl_budget?.crawlable_pages || 0,
              blockedPages: techAnalysis.crawl_info?.crawl_budget?.blocked_pages || 0,
              pagesPerLevel: {},
              averageCrawlTime: techAnalysis.crawl_info?.crawl_budget?.average_crawl_time || 0,
              crawlEfficiency: 0
            }
          },
          metrics: {
            totalPages: techAnalysis.metrics?.total_pages || 0,
            pagesWithIssues: techAnalysis.metrics?.pages_with_issues || 0,
            avgLoadTime: techAnalysis.metrics?.avg_load_time || 0,
            avgPageSize: techAnalysis.metrics?.avg_page_size || 0,
            totalIssues: techAnalysis.metrics?.total_issues || 0,
            issuesByType: {},
            healthScore: techAnalysis.metrics?.health_score || 0
          },
          config: {
            maxDepth: techAnalysis.analysis_config?.max_depth || 3,
            respectRobots: techAnalysis.analysis_config?.respect_robots || true,
            userAgent: techAnalysis.analysis_config?.user_agent || 'Fire-Salamander-Bot/1.0',
            crawlDelay: techAnalysis.analysis_config?.crawl_delay || 1000
          },
          status: {
            analysisDate: new Date().toISOString(),
            crawlDuration: 0,
            crawlStatus: 'completed',
            lastUpdate: new Date().toISOString(),
            version: '1.0'
          }
        };
      };
      
      const result = mockMapper(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.pageAnalysis).toHaveLength(2);
      expect(result.pageAnalysis[0].url).toBe('https://example.com');
      expect(result.pageAnalysis[0].statusCode).toBe(200);
      expect(result.pageAnalysis[0].title.content).toBe('Example Title');
    });
    
    test('should handle page analysis correctly', () => {
      const result = mockBackendData;
      const techData = JSON.parse(result.result_data).technical_analysis;
      
      // Test first page (success case)
      const page1 = techData.pages[0];
      expect(page1.status_code).toBe(PageStatus.SUCCESS);
      expect(page1.seo_elements.title.issues).toContain('title-too-short');
      expect(page1.seo_elements.headings.structure).toBe(HeadingStructure.GOOD);
      
      // Test second page (error case)
      const page2 = techData.pages[1];
      expect(page2.status_code).toBe(PageStatus.NOT_FOUND);
      expect(page2.seo_elements.title.issues).toContain('title-missing');
      expect(page2.seo_elements.headings.structure).toBe(HeadingStructure.BAD);
    });
    
    test('should process global issues correctly', () => {
      const techData = JSON.parse(mockBackendData.result_data).technical_analysis;
      const globalIssues = techData.global_issues;
      
      expect(globalIssues.duplicate_content).toHaveLength(1);
      expect(globalIssues.duplicate_content[0].similarity).toBe(85);
      expect(globalIssues.duplicate_content[0].level).toBe(SimilarityLevel.HIGH);
      
      expect(globalIssues.duplicate_titles).toHaveLength(1);
      expect(globalIssues.duplicate_titles[0].count).toBe(2);
      
      expect(globalIssues.broken_links).toHaveLength(1);
      expect(globalIssues.broken_links[0].status).toBe(404);
    });
    
    test('should handle crawlability data correctly', () => {
      const techData = JSON.parse(mockBackendData.result_data).technical_analysis;
      const crawlInfo = techData.crawl_info;
      
      expect(crawlInfo.robots_txt.exists).toBe(true);
      expect(crawlInfo.robots_txt.disallowed_paths).toContain('/admin');
      expect(crawlInfo.robots_txt.sitemap_urls).toContain('https://example.com/sitemap.xml');
      
      expect(crawlInfo.sitemap.exists).toBe(true);
      expect(crawlInfo.sitemap.format).toBe('xml');
      expect(crawlInfo.sitemap.pages_in_sitemap).toBe(25);
      
      expect(crawlInfo.crawl_budget.total_pages).toBe(30);
      expect(crawlInfo.crawl_budget.crawlable_pages).toBe(25);
    });
  });
  
  describe('Page Analysis Processing', () => {
    
    test('should calculate issue types correctly', () => {
      const page = {
        seo_elements: {
          title: { content: 'Short', length: 5, issues: ['title-too-short'] },
          meta_description: { content: '', length: 0, issues: ['meta-missing'] },
          headings: { h1: [], h2: [], h3: [], issues: ['h1-missing'] }
        },
        images: [
          { src: '/img.jpg', alt: '', issues: ['image-no-alt'] }
        ]
      };
      
      const allIssues = [
        ...page.seo_elements.title.issues,
        ...page.seo_elements.meta_description.issues,
        ...page.seo_elements.headings.issues,
        ...page.images.flatMap(img => img.issues)
      ];
      
      expect(allIssues).toHaveLength(4);
      expect(allIssues).toContain(IssueType.TITLE_TOO_SHORT);
      expect(allIssues).toContain(IssueType.META_MISSING);
      expect(allIssues).toContain(IssueType.H1_MISSING);
      expect(allIssues).toContain(IssueType.IMAGE_NO_ALT);
    });
    
    test('should determine page health status', () => {
      const calculatePageHealth = (issues: string[]) => {
        if (issues.length === 0) return 'good';
        if (issues.length <= 3) return 'warning';
        return 'error';
      };
      
      expect(calculatePageHealth([])).toBe('good');
      expect(calculatePageHealth(['title-too-short'])).toBe('warning');
      expect(calculatePageHealth(['title-missing', 'meta-missing', 'h1-missing', 'image-no-alt'])).toBe('error');
    });
    
    test('should validate URL formats', () => {
      const isValidUrl = (url: string) => {
        try {
          new URL(url);
          return true;
        } catch {
          return false;
        }
      };
      
      expect(isValidUrl('https://example.com')).toBe(true);
      expect(isValidUrl('http://example.com/page')).toBe(true);
      expect(isValidUrl('invalid-url')).toBe(false);
      expect(isValidUrl('')).toBe(false);
    });
  });
  
  describe('Schema and Structured Data Processing', () => {
    
    test('should process schema markup correctly', () => {
      const schemaData = [
        {
          type: 'Organization',
          valid: true,
          properties: { name: 'Company', url: 'https://example.com' }
        },
        {
          type: 'Article',
          valid: false,
          errors: ['Missing required field: headline']
        }
      ];
      
      expect(schemaData).toHaveLength(2);
      expect(schemaData[0].valid).toBe(true);
      expect(schemaData[1].valid).toBe(false);
      expect(schemaData[1].errors).toContain('Missing required field: headline');
    });
    
    test('should process Open Graph data', () => {
      const openGraph = {
        'og:title': 'Page Title',
        'og:description': 'Page Description',
        'og:image': 'https://example.com/image.jpg',
        'og:url': 'https://example.com'
      };
      
      expect(openGraph['og:title']).toBeDefined();
      expect(openGraph['og:image']).toMatch(/^https?:\/\//);
      expect(Object.keys(openGraph)).toHaveLength(4);
    });
  });
  
  describe('Error Handling', () => {
    
    test('should handle missing result_data gracefully', () => {
      const emptyData = { id: 1, url: 'https://example.com', result_data: '{}' };
      
      expect(() => {
        const parsed = JSON.parse(emptyData.result_data);
        const techAnalysis = parsed.technical_analysis || {};
        return techAnalysis;
      }).not.toThrow();
    });
    
    test('should handle malformed JSON gracefully', () => {
      const malformedData = { 
        id: 1, 
        url: 'https://example.com', 
        result_data: 'invalid json{' 
      };
      
      expect(() => {
        try {
          JSON.parse(malformedData.result_data);
        } catch {
          return {}; // Fallback pour données malformées
        }
      }).not.toThrow();
    });
    
    test('should provide default values for missing fields', () => {
      const incompleteData = {
        pages: [
          {
            url: 'https://example.com',
            status_code: 200
            // Champs manquants intentionnellement
          }
        ]
      };
      
      const page = incompleteData.pages[0];
      const processedPage = {
        ...page,
        load_time: page.load_time || 0,
        page_size: page.page_size || 0,
        seo_elements: page.seo_elements || {},
        images: page.images || [],
        links: page.links || { internal: 0, external: 0, broken: [] }
      };
      
      expect(processedPage.load_time).toBe(0);
      expect(processedPage.images).toEqual([]);
      expect(processedPage.links.internal).toBe(0);
    });
  });
  
  describe('Performance Metrics', () => {
    
    test('should process Core Web Vitals correctly', () => {
      const performance = {
        fcp: 800,   // First Contentful Paint
        lcp: 1200,  // Largest Contentful Paint  
        cls: 0.05,  // Cumulative Layout Shift
        ttfb: 200   // Time to First Byte
      };
      
      expect(performance.fcp).toBeLessThan(1800); // Good FCP
      expect(performance.lcp).toBeLessThan(2500); // Good LCP
      expect(performance.cls).toBeLessThan(0.1);  // Good CLS
      expect(performance.ttfb).toBeLessThan(800); // Good TTFB
    });
    
    test('should categorize page speed correctly', () => {
      const categorizeSpeed = (loadTime: number) => {
        if (loadTime < 1000) return 'fast';
        if (loadTime < 3000) return 'moderate';
        return 'slow';
      };
      
      expect(categorizeSpeed(800)).toBe('fast');
      expect(categorizeSpeed(1500)).toBe('moderate');
      expect(categorizeSpeed(4000)).toBe('slow');
    });
  });
});