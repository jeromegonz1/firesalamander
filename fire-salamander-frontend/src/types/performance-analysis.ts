/**
 * Fire Salamander - Performance Analysis Types
 * Professional Core Web Vitals & Performance Analysis (Lead Tech quality)
 * Comparable to GTmetrix, PageSpeed Insights, WebPageTest
 */

// Enums pour une meilleure type safety
export enum PerformanceGrade {
  EXCELLENT = 'A',
  GOOD = 'B',
  NEEDS_IMPROVEMENT = 'C',
  POOR = 'D',
  FAIL = 'F',
}

export enum DeviceType {
  MOBILE = 'mobile',
  DESKTOP = 'desktop',
  TABLET = 'tablet',
}

export enum RecommendationType {
  IMAGES = 'images',
  CSS = 'css',
  JAVASCRIPT = 'javascript',
  FONTS = 'fonts',
  CACHING = 'caching',
  COMPRESSION = 'compression',
  CDN = 'cdn',
  MINIFICATION = 'minification',
  CRITICAL_CSS = 'critical-css',
  PRELOAD = 'preload',
  LAZY_LOADING = 'lazy-loading',
  CODE_SPLITTING = 'code-splitting',
}

export enum ImpactLevel {
  HIGH = 'high',
  MEDIUM = 'medium',
  LOW = 'low',
}

// Core Web Vitals selon Google
export interface CoreWebVitals {
  // Largest Contentful Paint (2.5s good, 4s poor)
  lcp: {
    value: number; // en millisecondes
    grade: PerformanceGrade;
    percentile: number; // 0-100
  };
  
  // First Input Delay (100ms good, 300ms poor)
  fid: {
    value: number; // en millisecondes
    grade: PerformanceGrade;
    percentile: number;
  };
  
  // Cumulative Layout Shift (0.1 good, 0.25 poor)
  cls: {
    value: number; // score CLS
    grade: PerformanceGrade;
    percentile: number;
  };
  
  // First Contentful Paint
  fcp: {
    value: number; // en millisecondes
    grade: PerformanceGrade;
    percentile: number;
  };
  
  // Time to Interactive
  tti: {
    value: number; // en millisecondes
    grade: PerformanceGrade;
    percentile: number;
  };
  
  // Speed Index
  speedIndex: {
    value: number; // en millisecondes
    grade: PerformanceGrade;
    percentile: number;
  };
  
  // Score global (0-100 comme PageSpeed)
  overallScore: number;
}

// Ressources détaillées par type
export interface ResourceMetrics {
  images: {
    count: number;
    size: number; // en bytes
    unoptimized: number; // images non optimisées
    oversized: number; // images trop grandes
    formats: Record<string, number>; // jpg: 5, png: 3, webp: 2, etc.
    avgSize: number;
    largestSize: number;
    withoutAlt: number;
  };
  
  css: {
    count: number;
    size: number;
    external: number;
    inline: number;
    blocking: number; // CSS bloquant le rendu
    unused: number; // CSS inutilisé en %
    minified: boolean;
    avgSize: number;
  };
  
  js: {
    count: number;
    size: number;
    external: number;
    inline: number;
    blocking: number; // JS bloquant
    async: number;
    defer: number;
    unused: number; // JS inutilisé en %
    minified: boolean;
    avgSize: number;
  };
  
  fonts: {
    count: number;
    size: number;
    formats: Record<string, number>; // woff2: 3, woff: 1, ttf: 1
    preloaded: number;
    fallbacks: boolean;
    avgSize: number;
  };
  
  total: {
    count: number;
    size: number;
    requests: number;
    transferSize: number; // taille après compression
    compressionRatio: number; // % de compression
  };
}

// Optimisations techniques détaillées
export interface OptimizationStatus {
  compression: {
    enabled: boolean;
    algorithm: 'gzip' | 'brotli' | 'deflate' | null;
    ratio: number; // % de compression
    supportedTypes: string[];
  };
  
  caching: {
    browser: Record<string, string>; // Cache-Control headers par type
    cdn: boolean;
    serviceWorker: boolean;
    etags: boolean;
    lastModified: boolean;
    staticAssets: {
      images: string; // cache duration
      css: string;
      js: string;
      fonts: string;
    };
  };
  
  cdn: {
    enabled: boolean;
    provider: string | null; // Cloudflare, AWS CloudFront, etc.
    endpoints: string[];
    coverage: number; // % de ressources servies via CDN
  };
  
  minification: {
    css: {
      enabled: boolean;
      savings: number; // bytes économisés
      ratio: number; // % de réduction
    };
    js: {
      enabled: boolean;
      savings: number;
      ratio: number;
    };
    html: {
      enabled: boolean;
      savings: number;
      ratio: number;
    };
  };
  
  // Techniques modernes
  modernOptimizations: {
    criticalCSS: boolean;
    preloadKey: boolean; // preload des ressources critiques
    prefetch: boolean;
    lazyLoading: boolean;
    codeSplitting: boolean;
    treeShaking: boolean;
    http2: boolean;
    http3: boolean;
  };
}

// Détails de performance par page
export interface PagePerformanceMetrics {
  url: string;
  testedAt: string; // ISO date
  
  // Métriques par device
  mobile: CoreWebVitals;
  desktop: CoreWebVitals;
  
  // Détails techniques
  performance: {
    // Timing détaillé (Navigation Timing API)
    navigationTiming: {
      dns: number; // DNS lookup
      tcp: number; // TCP connection
      ssl: number; // SSL handshake
      ttfb: number; // Time to First Byte
      domLoaded: number; // DOM content loaded
      pageLoaded: number; // Page fully loaded
      interactive: number; // Time to Interactive
    };
    
    // Ressources
    resources: ResourceMetrics;
    
    // Optimisations
    optimization: OptimizationStatus;
    
    // Métriques serveur
    server: {
      responseTime: number;
      statusCode: number;
      redirects: number;
      redirectChain: string[];
    };
    
    // Métriques réseau
    network: {
      bandwidth: 'slow-3g' | 'fast-3g' | '4g' | 'broadband';
      latency: number;
      throughput: number;
    };
  };
  
  // Scores et grades
  scores: {
    performance: number; // 0-100
    accessibility: number;
    bestPractices: number;
    seo: number;
    pwa: number;
  };
}

// Recommandations avec impact quantifié
export interface PerformanceRecommendation {
  id: string;
  type: RecommendationType;
  title: string;
  description: string;
  impact: ImpactLevel;
  
  // Pages affectées
  pages: Array<{
    url: string;
    currentValue: number;
    potentialGain: number;
    priority: number; // 1-10
  }>;
  
  // Solutions techniques
  solution: {
    description: string;
    implementation: string;
    difficulty: 'easy' | 'medium' | 'hard';
    estimatedTime: string; // "2-4 heures"
    resources: string[]; // liens utiles
  };
  
  // Impact estimé
  estimatedGain: {
    loadTime: string; // "Réduction de 1.2s"
    scoreImprovement: number; // +15 points
    bandwidth: string; // "Économie de 245KB"
    userExperience: string; // Impact UX
  };
  
  // Mesures
  metrics: {
    before: number;
    after: number;
    improvement: number; // %
  };
  
  // Validation
  validation: {
    tool: string; // GTmetrix, PageSpeed, WebPageTest
    tested: boolean;
    results: string;
  };
}

// Analyse comparative
export interface CompetitorBenchmark {
  competitors: Array<{
    url: string;
    scores: {
      performance: number;
      lcp: number;
      fid: number;
      cls: number;
    };
    loadTime: number;
    pageSize: number;
  }>;
  
  ranking: {
    performance: number; // position dans le benchmark
    opportunities: string[];
  };
}

// Interface principale - Performance Analysis
export interface PerformanceAnalysis {
  // Métriques par page
  pageMetrics: PagePerformanceMetrics[];
  
  // Recommandations prioritaires
  recommendations: PerformanceRecommendation[];
  
  // Vue d'ensemble
  summary: {
    totalPages: number;
    avgPerformanceScore: number;
    avgLoadTime: {
      mobile: number;
      desktop: number;
    };
    
    // Distribution des scores
    scoreDistribution: {
      excellent: number; // pages avec score > 90
      good: number; // 75-90
      needsImprovement: number; // 50-74
      poor: number; // < 50
    };
    
    // Core Web Vitals globaux
    coreWebVitals: {
      mobile: {
        lcp: { avg: number, passing: number }; // % de pages qui passent
        fid: { avg: number, passing: number };
        cls: { avg: number, passing: number };
      };
      desktop: {
        lcp: { avg: number, passing: number };
        fid: { avg: number, passing: number };
        cls: { avg: number, passing: number };
      };
    };
    
    // Opportunités d'amélioration
    opportunities: {
      highImpact: number; // nombre de recommandations high impact
      estimatedGain: {
        loadTime: number; // secondes économisées
        score: number; // points de score
        bandwidth: number; // bytes économisés
      };
    };
  };
  
  // Benchmark concurrentiel
  benchmark?: CompetitorBenchmark;
  
  // Configuration de test
  config: {
    testLocation: string; // "Paris, France"
    device: DeviceType;
    connection: string; // "Cable"
    browser: string; // "Chrome"
    viewport: {
      width: number;
      height: number;
    };
    throttling: boolean;
  };
  
  // Métadonnées
  metadata: {
    analysisDate: string;
    analysisId: string;
    testDuration: number; // en secondes
    version: string;
    tool: string; // "Fire Salamander Performance"
  };
}

// Types pour les tableaux et filtres
export interface PagePerformanceTableRow extends PagePerformanceMetrics {
  id: string;
  selected?: boolean;
}

export enum PerformanceSortBy {
  SCORE = 'score',
  LOAD_TIME = 'loadTime',
  LCP = 'lcp',
  FID = 'fid',
  CLS = 'cls',
  PAGE_SIZE = 'pageSize',
  URL = 'url',
}

export enum PerformanceFilterBy {
  ALL = 'all',
  MOBILE = 'mobile',
  DESKTOP = 'desktop',
  POOR = 'poor',
  NEEDS_IMPROVEMENT = 'needs-improvement',
  GOOD = 'good',
  EXCELLENT = 'excellent',
}

// Utilitaires pour les calculs
export interface PerformanceThresholds {
  lcp: { good: number; poor: number };
  fid: { good: number; poor: number };
  cls: { good: number; poor: number };
  fcp: { good: number; poor: number };
  tti: { good: number; poor: number };
  speedIndex: { good: number; poor: number };
}

export const PERFORMANCE_THRESHOLDS: PerformanceThresholds = {
  lcp: { good: 2500, poor: 4000 },
  fid: { good: 100, poor: 300 },
  cls: { good: 0.1, poor: 0.25 },
  fcp: { good: 1800, poor: 3000 },
  tti: { good: 3800, poor: 7300 },
  speedIndex: { good: 3400, poor: 5800 },
};