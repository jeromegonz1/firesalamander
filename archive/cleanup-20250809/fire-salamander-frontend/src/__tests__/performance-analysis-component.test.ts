/**
 * Fire Salamander - Performance Analysis Component Tests (TDD)
 * Lead Tech quality - Component tests before implementation
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import {
  PerformanceAnalysis,
  CoreWebVitals,
  PerformanceGrade,
  DeviceType,
  RecommendationType,
  ImpactLevel,
  PerformanceSortBy,
  PerformanceFilterBy,
} from '@/types/performance-analysis';

// Mock data for testing
const mockPerformanceData: PerformanceAnalysis = {
  pageMetrics: [
    {
      url: 'https://example.com',
      testedAt: '2024-03-15T10:30:00Z',
      mobile: {
        lcp: { value: 2800, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 60 },
        fid: { value: 120, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 70 },
        cls: { value: 0.15, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 65 },
        fcp: { value: 1900, grade: PerformanceGrade.GOOD, percentile: 75 },
        tti: { value: 4200, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 60 },
        speedIndex: { value: 3800, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 65 },
        overallScore: 75,
      },
      desktop: {
        lcp: { value: 1800, grade: PerformanceGrade.EXCELLENT, percentile: 90 },
        fid: { value: 45, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
        cls: { value: 0.08, grade: PerformanceGrade.EXCELLENT, percentile: 92 },
        fcp: { value: 1200, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
        tti: { value: 2800, grade: PerformanceGrade.EXCELLENT, percentile: 88 },
        speedIndex: { value: 2200, grade: PerformanceGrade.EXCELLENT, percentile: 92 },
        overallScore: 90,
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
          images: { count: 15, size: 1500000, unoptimized: 3, oversized: 2, formats: { jpg: 8, png: 5, webp: 2 }, avgSize: 100000, largestSize: 500000, withoutAlt: 1 },
          css: { count: 5, size: 200000, external: 3, inline: 2, blocking: 2, unused: 30, minified: true, avgSize: 40000 },
          js: { count: 8, size: 800000, external: 6, inline: 2, blocking: 1, async: 4, defer: 3, unused: 25, minified: true, avgSize: 100000 },
          fonts: { count: 3, size: 150000, formats: { woff2: 3 }, preloaded: 1, fallbacks: true, avgSize: 50000 },
          total: { count: 31, size: 2650000, requests: 31, transferSize: 1850000, compressionRatio: 30.2 },
        },
        optimization: {
          compression: { enabled: true, algorithm: 'gzip', ratio: 70, supportedTypes: ['text/html'] },
          caching: { browser: {}, cdn: true, serviceWorker: false, etags: true, lastModified: true, staticAssets: { images: '7d', css: '1y', js: '1y', fonts: '1y' } },
          cdn: { enabled: true, provider: 'Cloudflare', endpoints: ['cdn.example.com'], coverage: 85 },
          minification: { css: { enabled: true, savings: 50000, ratio: 25 }, js: { enabled: true, savings: 200000, ratio: 25 }, html: { enabled: true, savings: 5000, ratio: 10 } },
          modernOptimizations: { criticalCSS: true, preloadKey: true, prefetch: false, lazyLoading: true, codeSplitting: true, treeShaking: true, http2: true, http3: false },
        },
        server: { responseTime: 200, statusCode: 200, redirects: 0, redirectChain: [] },
        network: { bandwidth: 'fast-3g', latency: 150, throughput: 1600 },
      },
      scores: { performance: 82, accessibility: 90, bestPractices: 88, seo: 92, pwa: 70 },
    },
  ],
  recommendations: [
    {
      id: 'rec-001',
      type: RecommendationType.IMAGES,
      title: 'Optimize images',
      description: 'Convert images to modern formats',
      impact: ImpactLevel.HIGH,
      pages: [{ url: 'https://example.com', currentValue: 1500000, potentialGain: 750000, priority: 9 }],
      solution: { description: 'Use WebP format', implementation: 'Install imagemin', difficulty: 'medium', estimatedTime: '2-4h', resources: [] },
      estimatedGain: { loadTime: '1.2s', scoreImprovement: 15, bandwidth: '750KB', userExperience: 'Better LCP' },
      metrics: { before: 1500000, after: 750000, improvement: 50 },
      validation: { tool: 'PageSpeed', tested: false, results: '' },
    },
  ],
  summary: {
    totalPages: 1,
    avgPerformanceScore: 82,
    avgLoadTime: { mobile: 2800, desktop: 1800 },
    scoreDistribution: { excellent: 0, good: 1, needsImprovement: 0, poor: 0 },
    coreWebVitals: {
      mobile: { lcp: { avg: 2800, passing: 60 }, fid: { avg: 120, passing: 70 }, cls: { avg: 0.15, passing: 65 } },
      desktop: { lcp: { avg: 1800, passing: 90 }, fid: { avg: 45, passing: 95 }, cls: { avg: 0.08, passing: 92 } },
    },
    opportunities: { highImpact: 1, estimatedGain: { loadTime: 1.2, score: 15, bandwidth: 750000 } },
  },
  config: { testLocation: 'Paris', device: DeviceType.MOBILE, connection: 'Cable', browser: 'Chrome', viewport: { width: 375, height: 667 }, throttling: true },
  metadata: { analysisDate: '2024-03-15T10:30:00Z', analysisId: 'perf-001', testDuration: 45, version: '1.0', tool: 'Fire Salamander' },
};

describe('PerformanceAnalysisSection Component (TDD)', () => {
  beforeEach(() => {
    // Reset DOM before each test
    document.body.innerHTML = '';
  });

  describe('Component Rendering', () => {
    it('should render without crashing', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      expect(() => {
        render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      }).not.toThrow();
    });

    it('should display main header and summary statistics', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      expect(screen.getByText('Analyse Performance')).toBeInTheDocument();
      expect(screen.getByText('1 page analysée')).toBeInTheDocument();
      expect(screen.getByText('82/100')).toBeInTheDocument(); // Performance score
    });

    it('should show Core Web Vitals metrics', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      // Mobile metrics
      expect(screen.getByText('2,800ms')).toBeInTheDocument(); // Mobile LCP
      expect(screen.getByText('120ms')).toBeInTheDocument(); // Mobile FID
      expect(screen.getByText('0.15')).toBeInTheDocument(); // Mobile CLS
      
      // Desktop metrics
      expect(screen.getByText('1,800ms')).toBeInTheDocument(); // Desktop LCP
      expect(screen.getByText('45ms')).toBeInTheDocument(); // Desktop FID
      expect(screen.getByText('0.08')).toBeInTheDocument(); // Desktop CLS
    });

    it('should display tabs for different views', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      expect(screen.getByText('Vue d\'ensemble')).toBeInTheDocument();
      expect(screen.getByText('Core Web Vitals')).toBeInTheDocument();
      expect(screen.getByText('Recommandations')).toBeInTheDocument();
      expect(screen.getByText('Ressources')).toBeInTheDocument();
    });
  });

  describe('Core Web Vitals Tab', () => {
    it('should show mobile and desktop toggle', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      // Click on Core Web Vitals tab
      fireEvent.click(screen.getByText('Core Web Vitals'));
      
      expect(screen.getByText('Mobile')).toBeInTheDocument();
      expect(screen.getByText('Desktop')).toBeInTheDocument();
    });

    it('should switch between mobile and desktop views', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Core Web Vitals'));
      
      // Default should show mobile
      expect(screen.getByText('2,800ms')).toBeInTheDocument();
      
      // Click desktop
      fireEvent.click(screen.getByText('Desktop'));
      
      await waitFor(() => {
        expect(screen.getByText('1,800ms')).toBeInTheDocument();
      });
    });

    it('should show correct performance grades', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Core Web Vitals'));
      
      // Should show "À améliorer" for mobile (NEEDS_IMPROVEMENT)
      expect(screen.getAllByText('À améliorer')).toHaveLength(4); // LCP, FID, CLS, TTI
      
      // Click desktop - should show "Excellent" grades
      fireEvent.click(screen.getByText('Desktop'));
      
      await waitFor(() => {
        expect(screen.getAllByText('Excellent')).toHaveLength(6); // All 6 metrics
      });
    });
  });

  describe('Recommendations Tab', () => {
    it('should display recommendations list', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Recommandations'));
      
      expect(screen.getByText('Optimize images')).toBeInTheDocument();
      expect(screen.getByText('Convert images to modern formats')).toBeInTheDocument();
      expect(screen.getByText('Impact élevé')).toBeInTheDocument();
    });

    it('should show estimated gains', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Recommandations'));
      
      expect(screen.getByText('1.2s')).toBeInTheDocument(); // Load time gain
      expect(screen.getByText('+15 points')).toBeInTheDocument(); // Score improvement
      expect(screen.getByText('750KB')).toBeInTheDocument(); // Bandwidth savings
    });

    it('should filter recommendations by impact', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Recommandations'));
      
      // Should have filter dropdown
      expect(screen.getByText('Tous les impacts')).toBeInTheDocument();
      
      // Filter by high impact
      fireEvent.click(screen.getByText('Tous les impacts'));
      fireEvent.click(screen.getByText('Impact élevé'));
      
      // Should still show the recommendation (it has high impact)
      expect(screen.getByText('Optimize images')).toBeInTheDocument();
    });
  });

  describe('Resources Tab', () => {
    it('should display resource breakdown', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Ressources'));
      
      expect(screen.getByText('Images')).toBeInTheDocument();
      expect(screen.getByText('CSS')).toBeInTheDocument();
      expect(screen.getByText('JavaScript')).toBeInTheDocument();
      expect(screen.getByText('Polices')).toBeInTheDocument();
    });

    it('should show resource sizes and counts', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Ressources'));
      
      expect(screen.getByText('15 fichiers')).toBeInTheDocument(); // Images count
      expect(screen.getByText('1.43 MB')).toBeInTheDocument(); // Images size (1500000 bytes)
      expect(screen.getByText('5 fichiers')).toBeInTheDocument(); // CSS count
      expect(screen.getByText('195.31 KB')).toBeInTheDocument(); // CSS size (200000 bytes)
    });

    it('should show optimization status', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Ressources'));
      
      // Should show optimization indicators
      expect(screen.getByText('Compression activée')).toBeInTheDocument();
      expect(screen.getByText('CDN activé')).toBeInTheDocument();
      expect(screen.getByText('Minification CSS')).toBeInTheDocument();
    });
  });

  describe('Interactions and State Management', () => {
    it('should handle tab switching', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      // Default tab should be "Vue d'ensemble"
      expect(screen.getByText('Score global')).toBeInTheDocument();
      
      // Switch to Core Web Vitals
      fireEvent.click(screen.getByText('Core Web Vitals'));
      expect(screen.getByText('Largest Contentful Paint')).toBeInTheDocument();
      
      // Switch to Recommendations
      fireEvent.click(screen.getByText('Recommandations'));
      expect(screen.getByText('Optimize images')).toBeInTheDocument();
      
      // Switch to Resources
      fireEvent.click(screen.getByText('Ressources'));
      expect(screen.getByText('Répartition des ressources')).toBeInTheDocument();
    });

    it('should handle sorting in recommendations', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      fireEvent.click(screen.getByText('Recommandations'));
      
      // Should have sort dropdown
      expect(screen.getByText('Trier par impact')).toBeInTheDocument();
      
      fireEvent.click(screen.getByText('Trier par impact'));
      fireEvent.click(screen.getByText('Par gain estimé'));
      
      // Should still show recommendations (just one in this case)
      expect(screen.getByText('Optimize images')).toBeInTheDocument();
    });

    it('should handle empty data gracefully', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      const emptyData = {
        ...mockPerformanceData,
        pageMetrics: [],
        recommendations: [],
        summary: { ...mockPerformanceData.summary, totalPages: 0 },
      };
      
      render(<PerformanceAnalysisSection performanceData={emptyData} analysisId="perf-001" />);
      
      expect(screen.getByText('0 page analysée')).toBeInTheDocument();
      expect(screen.getByText('Aucune données disponibles')).toBeInTheDocument();
    });

    it('should show loading state', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      // Test with undefined data to simulate loading
      render(<PerformanceAnalysisSection performanceData={undefined as any} analysisId="perf-001" />);
      
      expect(screen.getByText('Chargement de l\'analyse performance...')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels and roles', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      // Tab navigation should be accessible
      const tabList = screen.getByRole('tablist');
      expect(tabList).toBeInTheDocument();
      
      const tabs = screen.getAllByRole('tab');
      expect(tabs).toHaveLength(4); // 4 tabs
      
      // First tab should be selected
      expect(tabs[0]).toHaveAttribute('aria-selected', 'true');
    });

    it('should support keyboard navigation', async () => {
      const { PerformanceAnalysisSection } = await import('@/components/analysis/performance-analysis-section');
      
      render(<PerformanceAnalysisSection performanceData={mockPerformanceData} analysisId="perf-001" />);
      
      const firstTab = screen.getAllByRole('tab')[0];
      const secondTab = screen.getAllByRole('tab')[1];
      
      // Focus first tab
      firstTab.focus();
      expect(document.activeElement).toBe(firstTab);
      
      // Arrow key navigation
      fireEvent.keyDown(firstTab, { key: 'ArrowRight' });
      expect(document.activeElement).toBe(secondTab);
    });
  });
});