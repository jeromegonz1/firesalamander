/**
 * Fire Salamander - Performance Mapper Tests (TDD)
 * Lead Tech quality - Tests before implementation
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
} from '@/types/performance-analysis';

// Mock backend data structure
interface MockBackendPerformanceData {
  id: string;
  url: string;
  analyzed_at: string;
  performance_data: {
    mobile: {
      lighthouse: {
        performance: number;
        lcp: number;
        fid: number;
        cls: number;
        fcp: number;
        tti: number;
        speedIndex: number;
      };
      resources: {
        images: number;
        css: number;
        js: number;
        fonts: number;
        totalSize: number;
        requests: number;
      };
    };
    desktop: {
      lighthouse: {
        performance: number;
        lcp: number;
        fid: number;
        cls: number;
        fcp: number;
        tti: number;
        speedIndex: number;
      };
      resources: {
        images: number;
        css: number;
        js: number;
        fonts: number;
        totalSize: number;
        requests: number;
      };
    };
    optimizations: {
      compression: boolean;
      caching: boolean;
      cdn: boolean;
      minification: {
        css: boolean;
        js: boolean;
        html: boolean;
      };
    };
    recommendations: Array<{
      type: string;
      impact: string;
      description: string;
      savings: number;
    }>;
  };
}

describe('Performance Mapper Tests (TDD)', () => {
  let mockBackendData: MockBackendPerformanceData;

  beforeEach(() => {
    mockBackendData = {
      id: 'analysis-123',
      url: 'https://example.com',
      analyzed_at: '2024-03-15T10:30:00Z',
      performance_data: {
        mobile: {
          lighthouse: {
            performance: 75,
            lcp: 2800,
            fid: 120,
            cls: 0.15,
            fcp: 1900,
            tti: 4200,
            speedIndex: 3800,
          },
          resources: {
            images: 1500000,
            css: 200000,
            js: 800000,
            fonts: 150000,
            totalSize: 2650000,
            requests: 31,
          },
        },
        desktop: {
          lighthouse: {
            performance: 90,
            lcp: 1800,
            fid: 45,
            cls: 0.08,
            fcp: 1200,
            tti: 2800,
            speedIndex: 2200,
          },
          resources: {
            images: 1200000,
            css: 180000,
            js: 750000,
            fonts: 120000,
            totalSize: 2250000,
            requests: 28,
          },
        },
        optimizations: {
          compression: true,
          caching: false,
          cdn: true,
          minification: {
            css: true,
            js: false,
            html: true,
          },
        },
        recommendations: [
          {
            type: 'images',
            impact: 'high',
            description: 'Optimize images for better performance',
            savings: 750000,
          },
          {
            type: 'javascript',
            impact: 'medium',
            description: 'Enable JavaScript minification',
            savings: 200000,
          },
        ],
      },
    };
  });

  describe('mapBackendToPerformanceAnalysis function', () => {
    it('should exist and be callable', async () => {
      // This test will fail until we implement the function
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      expect(typeof mapBackendToPerformanceAnalysis).toBe('function');
    });

    it('should map backend data to PerformanceAnalysis structure', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.pageMetrics).toHaveLength(1);
      expect(result.pageMetrics[0].url).toBe('https://example.com');
      expect(result.metadata.analysisId).toBe('analysis-123');
    });

    it('should correctly map Core Web Vitals', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      const mobileCWV = result.pageMetrics[0].mobile;
      
      expect(mobileCWV.lcp.value).toBe(2800);
      expect(mobileCWV.fid.value).toBe(120);
      expect(mobileCWV.cls.value).toBe(0.15);
      expect(mobileCWV.overallScore).toBe(75);
    });

    it('should assign correct performance grades', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      const mobileCWV = result.pageMetrics[0].mobile;
      
      // LCP: 2800ms should be GOOD (< 4000ms but > 2500ms)
      expect(mobileCWV.lcp.grade).toBe(PerformanceGrade.NEEDS_IMPROVEMENT);
      
      // FID: 120ms should be NEEDS_IMPROVEMENT (> 100ms but < 300ms)
      expect(mobileCWV.fid.grade).toBe(PerformanceGrade.NEEDS_IMPROVEMENT);
      
      // CLS: 0.15 should be NEEDS_IMPROVEMENT (> 0.1 but < 0.25)
      expect(mobileCWV.cls.grade).toBe(PerformanceGrade.NEEDS_IMPROVEMENT);
    });

    it('should handle desktop vs mobile differences', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      const { mobile, desktop } = result.pageMetrics[0];
      
      // Desktop should perform better than mobile
      expect(desktop.overallScore).toBeGreaterThan(mobile.overallScore);
      expect(desktop.lcp.value).toBeLessThan(mobile.lcp.value);
      expect(desktop.fid.value).toBeLessThan(mobile.fid.value);
    });

    it('should map resource metrics correctly', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      const resources = result.pageMetrics[0].performance.resources;
      
      expect(resources.images.size).toBe(1500000);
      expect(resources.css.size).toBe(200000);
      expect(resources.js.size).toBe(800000);
      expect(resources.fonts.size).toBe(150000);
      expect(resources.total.size).toBe(2650000);
      expect(resources.total.requests).toBe(31);
    });

    it('should generate recommendations correctly', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      
      expect(result.recommendations).toHaveLength(2);
      expect(result.recommendations[0].type).toBe(RecommendationType.IMAGES);
      expect(result.recommendations[0].impact).toBe(ImpactLevel.HIGH);
      expect(result.recommendations[1].type).toBe(RecommendationType.JAVASCRIPT);
      expect(result.recommendations[1].impact).toBe(ImpactLevel.MEDIUM);
    });

    it('should handle missing or invalid data gracefully', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const invalidData = {};
      const result = mapBackendToPerformanceAnalysis(invalidData);
      
      // Should return a valid structure with fallback data
      expect(result.pageMetrics).toHaveLength(0);
      expect(result.recommendations).toHaveLength(0);
      expect(result.summary.totalPages).toBe(0);
      expect(result.metadata.analysisId).toBeTruthy();
    });

    it('should calculate summary statistics correctly', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      
      expect(result.summary.totalPages).toBe(1);
      expect(result.summary.avgPerformanceScore).toBe(82.5); // (75 + 90) / 2
      expect(result.summary.avgLoadTime.mobile).toBeGreaterThan(0);
      expect(result.summary.avgLoadTime.desktop).toBeGreaterThan(0);
    });

    it('should validate optimization status mapping', async () => {
      const { mapBackendToPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const result = mapBackendToPerformanceAnalysis(mockBackendData);
      const optimization = result.pageMetrics[0].performance.optimization;
      
      expect(optimization.compression.enabled).toBe(true);
      expect(optimization.cdn.enabled).toBe(true);
      expect(optimization.minification.css.enabled).toBe(true);
      expect(optimization.minification.js.enabled).toBe(false);
      expect(optimization.minification.html.enabled).toBe(true);
    });
  });

  describe('Helper functions', () => {
    it('should have a function to calculate performance grade', async () => {
      const { calculatePerformanceGrade } = await import('@/lib/mappers/performance-mapper');
      
      expect(calculatePerformanceGrade(2000, 'lcp')).toBe(PerformanceGrade.EXCELLENT);
      expect(calculatePerformanceGrade(3000, 'lcp')).toBe(PerformanceGrade.NEEDS_IMPROVEMENT);
      expect(calculatePerformanceGrade(5000, 'lcp')).toBe(PerformanceGrade.POOR);
    });

    it('should have a function to create empty performance analysis', async () => {
      const { createEmptyPerformanceAnalysis } = await import('@/lib/mappers/performance-mapper');
      
      const emptyAnalysis = createEmptyPerformanceAnalysis();
      
      expect(emptyAnalysis.pageMetrics).toHaveLength(0);
      expect(emptyAnalysis.recommendations).toHaveLength(0);
      expect(emptyAnalysis.summary.totalPages).toBe(0);
      expect(emptyAnalysis.metadata.analysisId).toBeTruthy();
    });

    it('should have a function to map recommendation types', async () => {
      const { mapRecommendationType } = await import('@/lib/mappers/performance-mapper');
      
      expect(mapRecommendationType('images')).toBe(RecommendationType.IMAGES);
      expect(mapRecommendationType('javascript')).toBe(RecommendationType.JAVASCRIPT);
      expect(mapRecommendationType('css')).toBe(RecommendationType.CSS);
      expect(mapRecommendationType('unknown')).toBe(RecommendationType.IMAGES); // fallback
    });

    it('should have a function to map impact levels', async () => {
      const { mapImpactLevel } = await import('@/lib/mappers/performance-mapper');
      
      expect(mapImpactLevel('high')).toBe(ImpactLevel.HIGH);
      expect(mapImpactLevel('medium')).toBe(ImpactLevel.MEDIUM);
      expect(mapImpactLevel('low')).toBe(ImpactLevel.LOW);
      expect(mapImpactLevel('unknown')).toBe(ImpactLevel.LOW); // fallback
    });

    it('should calculate percentiles correctly', async () => {
      const { calculatePercentile } = await import('@/lib/mappers/performance-mapper');
      
      // LCP of 2000ms should be in high percentile (good performance)
      expect(calculatePercentile(2000, 'lcp')).toBeGreaterThan(80);
      
      // LCP of 5000ms should be in low percentile (poor performance)
      expect(calculatePercentile(5000, 'lcp')).toBeLessThan(20);
    });
  });
});