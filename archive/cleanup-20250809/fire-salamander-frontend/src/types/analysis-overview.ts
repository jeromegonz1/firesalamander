/**
 * Fire Salamander - Analysis Overview Types
 * Professional SEO Report Data Structure (SEMrush/Ahrefs style)
 */

export interface AnalysisOverview {
  // Metadata
  metadata: {
    url: string;
    analyzedAt: string;
    analysisType: 'Quick Scan' | 'AI Boost Scan' | 'Custom Scan';
    processingTime: string;
    domain: string;
    protocol: 'http' | 'https';
  };

  // Scores détaillés
  scores: {
    global: number;
    seo: {
      score: number;
      trend: 'up' | 'down' | 'stable';
      details: {
        titleOptimization: number;
        metaOptimization: number;
        headingStructure: number;
        keywordUsage: number;
        contentQuality: number;
      };
    };
    technical: {
      score: number;
      trend: 'up' | 'down' | 'stable';
      details: {
        crawlability: number;
        indexability: number;
        siteSpeed: number;
        mobile: number;
        security: number;
      };
    };
    performance: {
      score: number;
      trend: 'up' | 'down' | 'stable';
      coreWebVitals: {
        lcp: { value: number; unit: string; score: 'good' | 'needs-improvement' | 'poor' };
        fid: { value: number; unit: string; score: 'good' | 'needs-improvement' | 'poor' };
        cls: { value: number; unit: string; score: 'good' | 'needs-improvement' | 'poor' };
        ttfb: { value: number; unit: string; score: 'good' | 'needs-improvement' | 'poor' };
        inp: { value: number | null; unit: string; score: 'good' | 'needs-improvement' | 'poor' | 'n/a' };
      };
    };
    accessibility: {
      score: number;
      wcagLevel: 'A' | 'AA' | 'AAA' | 'Fail';
    };
  };

  // Métriques globales
  metrics: {
    totalPages: number;
    pagesAnalyzed: number;
    pagesWithErrors: number;
    
    // Issues breakdown
    issues: {
      total: number;
      critical: number;
      warnings: number;
      notices: number;
      passedChecks: number;
    };
    
    // Performance metrics
    performance: {
      avgLoadTime: number;
      pageSize: number;
      requests: number;
      avgTimeToFirstByte: number;
    };
    
    // Resources
    resources: {
      images: {
        total: number;
        optimized: number;
        missingAlt: number;
        oversized: number;
      };
      scripts: {
        total: number;
        minified: number;
        external: number;
        renderBlocking: number;
      };
      stylesheets: {
        total: number;
        minified: number;
        external: number;
        renderBlocking: number;
      };
      fonts: number;
      videos: number;
    };
    
    // Links analysis
    links: {
      internal: {
        total: number;
        unique: number;
        broken: number;
      };
      external: {
        total: number;
        unique: number;
        broken: number;
        nofollow: number;
      };
    };
    
    // SEO specific
    seo: {
      metaTags: {
        title: boolean;
        description: boolean;
        keywords: boolean;
        ogTags: boolean;
        twitterCards: boolean;
      };
      headings: {
        h1Count: number;
        h2Count: number;
        h3Count: number;
        missingH1: boolean;
        multipleH1: boolean;
      };
      schema: {
        present: boolean;
        types: string[];
      };
    };
  };

  // Top issues for quick overview
  topIssues: Array<{
    id: string;
    title: string;
    category: 'SEO' | 'Technical' | 'Performance' | 'Accessibility' | 'Security';
    severity: 'critical' | 'warning' | 'notice';
    pagesAffected: number;
    description: string;
  }>;

  // Competitor comparison (if available)
  competitorBenchmark?: {
    avgScore: number;
    position: number;
    totalCompetitors: number;
  };
}

// Helper types for UI components
export interface ScoreCardProps {
  title: string;
  score: number;
  maxScore?: number;
  trend?: 'up' | 'down' | 'stable';
  details?: Record<string, number>;
  color?: 'green' | 'yellow' | 'red' | 'blue';
}

export interface MetricCardProps {
  label: string;
  value: number | string;
  unit?: string;
  icon?: React.ComponentType;
  status?: 'good' | 'warning' | 'critical';
  tooltip?: string;
}

export interface IssuesSummaryProps {
  critical: number;
  warnings: number;
  notices: number;
  passed: number;
}

// Enums for consistency
export enum IssueSeverity {
  CRITICAL = 'critical',
  WARNING = 'warning',
  NOTICE = 'notice'
}

export enum ScoreStatus {
  EXCELLENT = 90,
  GOOD = 70,
  AVERAGE = 50,
  POOR = 30
}

// Utility function types
export type ScoreCalculator = (data: any) => number;
export type DataMapper = (backendData: any) => AnalysisOverview;