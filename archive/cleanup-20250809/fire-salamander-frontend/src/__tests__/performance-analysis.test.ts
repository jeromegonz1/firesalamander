/**
 * Fire Salamander - Performance Analysis Tests (TDD)
 * Lead Tech quality - Tests first approach
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  PerformanceAnalysis,
  CoreWebVitals,
  PerformanceGrade,
  DeviceType,
  RecommendationType,
  ImpactLevel,
  PERFORMANCE_THRESHOLDS,
  PagePerformanceMetrics,
  PerformanceRecommendation,
} from '@/types/performance-analysis';

describe('Performance Analysis Types', () => {
  describe('CoreWebVitals', () => {
    it('should have correct structure for Core Web Vitals', () => {
      const mockCoreWebVitals: CoreWebVitals = {
        lcp: { value: 2000, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
        fid: { value: 50, grade: PerformanceGrade.EXCELLENT, percentile: 90 },
        cls: { value: 0.05, grade: PerformanceGrade.EXCELLENT, percentile: 98 },
        fcp: { value: 1500, grade: PerformanceGrade.GOOD, percentile: 85 },
        tti: { value: 3000, grade: PerformanceGrade.GOOD, percentile: 80 },
        speedIndex: { value: 2800, grade: PerformanceGrade.EXCELLENT, percentile: 92 },
        overallScore: 95,
      };

      expect(mockCoreWebVitals.lcp.value).toBe(2000);
      expect(mockCoreWebVitals.lcp.grade).toBe(PerformanceGrade.EXCELLENT);
      expect(mockCoreWebVitals.overallScore).toBeGreaterThanOrEqual(0);
      expect(mockCoreWebVitals.overallScore).toBeLessThanOrEqual(100);
    });

    it('should validate LCP thresholds', () => {
      const excellentLCP = 2000; // < 2500ms = excellent
      const goodLCP = 3000; // 2500-4000ms = needs improvement
      const poorLCP = 5000; // > 4000ms = poor

      expect(excellentLCP).toBeLessThan(PERFORMANCE_THRESHOLDS.lcp.good);
      expect(goodLCP).toBeGreaterThan(PERFORMANCE_THRESHOLDS.lcp.good);
      expect(goodLCP).toBeLessThan(PERFORMANCE_THRESHOLDS.lcp.poor);
      expect(poorLCP).toBeGreaterThan(PERFORMANCE_THRESHOLDS.lcp.poor);
    });

    it('should validate CLS thresholds', () => {
      const excellentCLS = 0.05; // < 0.1 = excellent
      const goodCLS = 0.15; // 0.1-0.25 = needs improvement
      const poorCLS = 0.35; // > 0.25 = poor

      expect(excellentCLS).toBeLessThan(PERFORMANCE_THRESHOLDS.cls.good);
      expect(goodCLS).toBeGreaterThan(PERFORMANCE_THRESHOLDS.cls.good);
      expect(goodCLS).toBeLessThan(PERFORMANCE_THRESHOLDS.cls.poor);
      expect(poorCLS).toBeGreaterThan(PERFORMANCE_THRESHOLDS.cls.poor);
    });
  });

  describe('PagePerformanceMetrics', () => {
    let mockPageMetrics: PagePerformanceMetrics;

    beforeEach(() => {
      mockPageMetrics = {
        url: 'https://example.com',
        testedAt: new Date().toISOString(),
        mobile: {
          lcp: { value: 2500, grade: PerformanceGrade.GOOD, percentile: 75 },
          fid: { value: 100, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
          cls: { value: 0.1, grade: PerformanceGrade.GOOD, percentile: 80 },
          fcp: { value: 1800, grade: PerformanceGrade.EXCELLENT, percentile: 90 },
          tti: { value: 3800, grade: PerformanceGrade.GOOD, percentile: 70 },
          speedIndex: { value: 3400, grade: PerformanceGrade.GOOD, percentile: 75 },
          overallScore: 85,
        },
        desktop: {
          lcp: { value: 1800, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
          fid: { value: 20, grade: PerformanceGrade.EXCELLENT, percentile: 98 },
          cls: { value: 0.05, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
          fcp: { value: 1200, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
          tti: { value: 2500, grade: PerformanceGrade.EXCELLENT, percentile: 90 },
          speedIndex: { value: 2200, grade: PerformanceGrade.EXCELLENT, percentile: 92 },
          overallScore: 95,
        },
        performance: {
          navigationTiming: {
            dns: 10,
            tcp: 50,
            ssl: 100,
            ttfb: 200,
            domLoaded: 1500,
            pageLoaded: 2500,
            interactive: 3000,
          },
          resources: {
            images: {
              count: 15,
              size: 1500000,
              unoptimized: 3,
              oversized: 2,
              formats: { jpg: 8, png: 5, webp: 2 },
              avgSize: 100000,
              largestSize: 500000,
              withoutAlt: 1,
            },
            css: {
              count: 5,
              size: 200000,
              external: 3,
              inline: 2,
              blocking: 2,
              unused: 30,
              minified: true,
              avgSize: 40000,
            },
            js: {
              count: 8,
              size: 800000,
              external: 6,
              inline: 2,
              blocking: 1,
              async: 4,
              defer: 3,
              unused: 25,
              minified: true,
              avgSize: 100000,
            },
            fonts: {
              count: 3,
              size: 150000,
              formats: { woff2: 3 },
              preloaded: 1,
              fallbacks: true,
              avgSize: 50000,
            },
            total: {
              count: 31,
              size: 2650000,
              requests: 31,
              transferSize: 1850000,
              compressionRatio: 30.2,
            },
          },
          optimization: {
            compression: {
              enabled: true,
              algorithm: 'gzip',
              ratio: 70,
              supportedTypes: ['text/html', 'text/css', 'application/javascript'],
            },
            caching: {
              browser: {
                'text/css': 'public, max-age=31536000',
                'application/javascript': 'public, max-age=31536000',
                'image/*': 'public, max-age=604800',
              },
              cdn: true,
              serviceWorker: false,
              etags: true,
              lastModified: true,
              staticAssets: {
                images: '7 days',
                css: '1 year',
                js: '1 year',
                fonts: '1 year',
              },
            },
            cdn: {
              enabled: true,
              provider: 'Cloudflare',
              endpoints: ['cdn.example.com'],
              coverage: 85,
            },
            minification: {
              css: { enabled: true, savings: 50000, ratio: 25 },
              js: { enabled: true, savings: 200000, ratio: 25 },
              html: { enabled: true, savings: 5000, ratio: 10 },
            },
            modernOptimizations: {
              criticalCSS: true,
              preloadKey: true,
              prefetch: false,
              lazyLoading: true,
              codeSplitting: true,
              treeShaking: true,
              http2: true,
              http3: false,
            },
          },
          server: {
            responseTime: 200,
            statusCode: 200,
            redirects: 0,
            redirectChain: [],
          },
          network: {
            bandwidth: 'fast-3g',
            latency: 150,
            throughput: 1600,
          },
        },
        scores: {
          performance: 85,
          accessibility: 90,
          bestPractices: 88,
          seo: 92,
          pwa: 70,
        },
      };
    });

    it('should have valid page performance metrics structure', () => {
      expect(mockPageMetrics.url).toBe('https://example.com');
      expect(mockPageMetrics.mobile.overallScore).toBeGreaterThanOrEqual(0);
      expect(mockPageMetrics.desktop.overallScore).toBeGreaterThanOrEqual(0);
      expect(mockPageMetrics.performance.resources.total.count).toBeGreaterThan(0);
    });

    it('should have better desktop performance than mobile', () => {
      expect(mockPageMetrics.desktop.overallScore).toBeGreaterThanOrEqual(
        mockPageMetrics.mobile.overallScore
      );
      expect(mockPageMetrics.desktop.lcp.value).toBeLessThanOrEqual(
        mockPageMetrics.mobile.lcp.value
      );
    });

    it('should validate resource metrics', () => {
      const { resources } = mockPageMetrics.performance;
      
      expect(resources.total.count).toBe(
        resources.images.count + 
        resources.css.count + 
        resources.js.count + 
        resources.fonts.count
      );
      
      expect(resources.total.transferSize).toBeLessThan(resources.total.size);
      expect(resources.total.compressionRatio).toBeGreaterThan(0);
      expect(resources.total.compressionRatio).toBeLessThan(100);
    });
  });

  describe('PerformanceRecommendation', () => {
    let mockRecommendation: PerformanceRecommendation;

    beforeEach(() => {
      mockRecommendation = {
        id: 'rec-001',
        type: RecommendationType.IMAGES,
        title: 'Optimize images',
        description: 'Convert images to modern formats and compress them',
        impact: ImpactLevel.HIGH,
        pages: [
          {
            url: 'https://example.com',
            currentValue: 1500000,
            potentialGain: 750000,
            priority: 9,
          },
        ],
        solution: {
          description: 'Convert to WebP format and enable compression',
          implementation: 'Use imagemin or similar tools to optimize images',
          difficulty: 'medium',
          estimatedTime: '2-4 heures',
          resources: ['https://web.dev/serve-images-webp/'],
        },
        estimatedGain: {
          loadTime: 'Réduction de 1.2s',
          scoreImprovement: 15,
          bandwidth: 'Économie de 750KB',
          userExperience: 'Amélioration significative du LCP',
        },
        metrics: {
          before: 1500000,
          after: 750000,
          improvement: 50,
        },
        validation: {
          tool: 'PageSpeed Insights',
          tested: false,
          results: '',
        },
      };
    });

    it('should have valid recommendation structure', () => {
      expect(mockRecommendation.id).toBeTruthy();
      expect(mockRecommendation.impact).toBe(ImpactLevel.HIGH);
      expect(mockRecommendation.pages.length).toBeGreaterThan(0);
      expect(mockRecommendation.metrics.improvement).toBeGreaterThan(0);
    });

    it('should calculate improvement percentage correctly', () => {
      const expectedImprovement = 
        ((mockRecommendation.metrics.before - mockRecommendation.metrics.after) / 
         mockRecommendation.metrics.before) * 100;
      
      expect(mockRecommendation.metrics.improvement).toBe(expectedImprovement);
    });

    it('should have priority between 1 and 10', () => {
      mockRecommendation.pages.forEach(page => {
        expect(page.priority).toBeGreaterThanOrEqual(1);
        expect(page.priority).toBeLessThanOrEqual(10);
      });
    });
  });

  describe('PerformanceAnalysis', () => {
    let mockPerformanceAnalysis: PerformanceAnalysis;

    beforeEach(() => {
      mockPerformanceAnalysis = {
        pageMetrics: [],
        recommendations: [],
        summary: {
          totalPages: 10,
          avgPerformanceScore: 75,
          avgLoadTime: {
            mobile: 3200,
            desktop: 2100,
          },
          scoreDistribution: {
            excellent: 2,
            good: 4,
            needsImprovement: 3,
            poor: 1,
          },
          coreWebVitals: {
            mobile: {
              lcp: { avg: 2800, passing: 60 },
              fid: { avg: 120, passing: 80 },
              cls: { avg: 0.15, passing: 70 },
            },
            desktop: {
              lcp: { avg: 2200, passing: 80 },
              fid: { avg: 50, passing: 95 },
              cls: { avg: 0.08, passing: 90 },
            },
          },
          opportunities: {
            highImpact: 5,
            estimatedGain: {
              loadTime: 1.5,
              score: 20,
              bandwidth: 500000,
            },
          },
        },
        config: {
          testLocation: 'Paris, France',
          device: DeviceType.MOBILE,
          connection: 'Cable',
          browser: 'Chrome',
          viewport: { width: 375, height: 667 },
          throttling: true,
        },
        metadata: {
          analysisDate: new Date().toISOString(),
          analysisId: 'perf-001',
          testDuration: 45,
          version: '1.0',
          tool: 'Fire Salamander Performance',
        },
      };
    });

    it('should have valid performance analysis structure', () => {
      expect(mockPerformanceAnalysis.summary.totalPages).toBeGreaterThan(0);
      expect(mockPerformanceAnalysis.summary.avgPerformanceScore).toBeGreaterThanOrEqual(0);
      expect(mockPerformanceAnalysis.summary.avgPerformanceScore).toBeLessThanOrEqual(100);
      expect(mockPerformanceAnalysis.metadata.analysisId).toBeTruthy();
    });

    it('should have consistent score distribution', () => {
      const { scoreDistribution, totalPages } = mockPerformanceAnalysis.summary;
      const totalDistributed = 
        scoreDistribution.excellent + 
        scoreDistribution.good + 
        scoreDistribution.needsImprovement + 
        scoreDistribution.poor;
      
      expect(totalDistributed).toBe(totalPages);
    });

    it('should have better desktop performance than mobile', () => {
      const { mobile, desktop } = mockPerformanceAnalysis.summary.coreWebVitals;
      
      expect(desktop.lcp.avg).toBeLessThanOrEqual(mobile.lcp.avg);
      expect(desktop.fid.avg).toBeLessThanOrEqual(mobile.fid.avg);
      expect(desktop.cls.avg).toBeLessThanOrEqual(mobile.cls.avg);
      
      expect(desktop.lcp.passing).toBeGreaterThanOrEqual(mobile.lcp.passing);
      expect(desktop.fid.passing).toBeGreaterThanOrEqual(mobile.fid.passing);
      expect(desktop.cls.passing).toBeGreaterThanOrEqual(mobile.cls.passing);
    });

    it('should have valid viewport dimensions', () => {
      const { viewport } = mockPerformanceAnalysis.config;
      expect(viewport.width).toBeGreaterThan(0);
      expect(viewport.height).toBeGreaterThan(0);
    });

    it('should have positive estimated gains', () => {
      const { estimatedGain } = mockPerformanceAnalysis.summary.opportunities;
      expect(estimatedGain.loadTime).toBeGreaterThan(0);
      expect(estimatedGain.score).toBeGreaterThan(0);
      expect(estimatedGain.bandwidth).toBeGreaterThan(0);
    });
  });
});