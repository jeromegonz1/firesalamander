/**
 * Fire Salamander - Technical Analysis Types
 * Professional detailed page-by-page technical analysis (Lead Tech quality)
 */

// Enums pour une meilleure type safety
export enum PageStatus {
  SUCCESS = 200,
  REDIRECT = 300,
  CLIENT_ERROR = 400,
  NOT_FOUND = 404,
  SERVER_ERROR = 500,
}

export enum HeadingStructure {
  GOOD = 'good',
  BAD = 'bad',
}

export enum IssueType {
  TITLE_TOO_SHORT = 'title-too-short',
  TITLE_TOO_LONG = 'title-too-long',
  TITLE_MISSING = 'title-missing',
  TITLE_DUPLICATE = 'title-duplicate',
  META_TOO_SHORT = 'meta-too-short',
  META_TOO_LONG = 'meta-too-long',
  META_MISSING = 'meta-missing',
  META_DUPLICATE = 'meta-duplicate',
  H1_MISSING = 'h1-missing',
  H1_MULTIPLE = 'h1-multiple',
  HEADING_SKIP = 'heading-skip',
  IMAGE_NO_ALT = 'image-no-alt',
  IMAGE_LARGE_SIZE = 'image-large-size',
  LINK_BROKEN = 'link-broken',
  CANONICAL_MISSING = 'canonical-missing',
  ROBOTS_MISSING = 'robots-missing',
}

export enum SimilarityLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
}

// Interfaces principales
export interface PageSEOElement {
  content: string;
  length: number;
  hasKeyword: boolean;
  issues: IssueType[];
  recommendations?: string[];
}

export interface PageHeadings {
  h1: string[];
  h2: string[];
  h3: string[];
  h4?: string[];
  h5?: string[];
  h6?: string[];
  structure: HeadingStructure;
  issues: IssueType[];
  recommendations?: string[];
}

export interface PageImage {
  src: string;
  alt: string;
  title?: string;
  size: number; // en bytes
  dimensions?: {
    width: number;
    height: number;
  };
  format?: string; // jpg, png, webp, etc.
  loading?: 'lazy' | 'eager';
  issues: IssueType[];
  recommendations?: string[];
}

export interface PageLinks {
  internal: number;
  external: number;
  broken: Array<{
    url: string;
    status: number;
    anchorText: string;
    position: 'header' | 'footer' | 'content' | 'navigation';
  }>;
  nofollow: number;
  totalLinks: number;
}

export interface SchemaMarkup {
  type: string; // Organization, Article, Product, etc.
  valid: boolean;
  errors?: string[];
  properties?: Record<string, unknown>;
  url?: string;
}

export interface PageAnalysis {
  url: string;
  statusCode: PageStatus;
  loadTime: number; // en millisecondes
  size: number; // en bytes
  lastCrawled: string; // ISO date string
  depth: number; // profondeur depuis la homepage
  
  // SEO Elements
  title: PageSEOElement;
  metaDescription: PageSEOElement;
  headings: PageHeadings;
  
  // Technical Elements
  canonical: string;
  robots: string;
  metaRobots?: string;
  lang?: string;
  hreflang?: Array<{
    lang: string;
    url: string;
  }>;
  
  // Structured Data
  schema: SchemaMarkup[];
  openGraph: Record<string, string>;
  twitterCard?: Record<string, string>;
  
  // Media & Resources
  images: PageImage[];
  videos?: Array<{
    src: string;
    size: number;
    format: string;
  }>;
  
  // Links
  links: PageLinks;
  
  // Performance
  performance?: {
    fcp: number; // First Contentful Paint
    lcp: number; // Largest Contentful Paint
    cls: number; // Cumulative Layout Shift
    ttfb: number; // Time to First Byte
  };
  
  // Mobile
  mobile?: {
    viewport: boolean;
    responsive: boolean;
    touchTargets: number;
    fontSizes: number[];
  };
}

export interface DuplicateContent {
  pages: string[];
  similarity: number; // pourcentage
  level: SimilarityLevel;
  contentHash?: string;
  affectedElements: Array<'title' | 'meta' | 'h1' | 'content'>;
}

export interface BrokenLink {
  from: string;
  to: string;
  status: number;
  anchorText: string;
  linkType: 'internal' | 'external';
  position: 'header' | 'footer' | 'content' | 'navigation';
  lastChecked: string;
}

export interface RedirectChain {
  chain: string[];
  finalUrl: string;
  totalHops: number;
  totalTime: number;
  statusCodes: number[];
}

export interface GlobalIssues {
  duplicateContent: DuplicateContent[];
  duplicateTitles: Array<{
    title: string;
    pages: string[];
    count: number;
  }>;
  duplicateMeta: Array<{
    meta: string;
    pages: string[];
    count: number;
  }>;
  missingTitles: string[];
  missingMeta: string[];
  brokenLinks: BrokenLink[];
  orphanPages: string[]; // pages sans liens entrants
  redirectChains: RedirectChain[];
  largePages: Array<{
    url: string;
    size: number;
  }>; // pages > 1MB
  slowPages: Array<{
    url: string;
    loadTime: number;
  }>; // pages > 3s
}

export interface RobotsTxtAnalysis {
  exists: boolean;
  url?: string;
  valid: boolean;
  userAgents: string[];
  disallowedPaths: string[];
  allowedPaths: string[];
  sitemapUrls: string[];
  crawlDelay?: number;
  issues: string[];
  recommendations?: string[];
}

export interface SitemapAnalysis {
  exists: boolean;
  url: string;
  valid: boolean;
  format: 'xml' | 'txt' | 'unknown';
  pagesInSitemap: number;
  lastModified?: string;
  images: number;
  videos: number;
  issues: string[];
  recommendations?: string[];
  indexSitemaps?: string[]; // pour les sitemap index
}

export interface CrawlBudget {
  totalPages: number;
  crawlablePages: number;
  blockedPages: number;
  indexedPages?: number;
  pagesPerLevel: Record<number, number>; // depth -> count
  averageCrawlTime: number;
  crawlEfficiency: number; // pourcentage
}

export interface TechnicalAnalysis {
  // Analyse par page
  pageAnalysis: PageAnalysis[];
  
  // Problèmes globaux
  globalIssues: GlobalIssues;
  
  // Crawlability
  crawlability: {
    robotsTxt: RobotsTxtAnalysis;
    sitemap: SitemapAnalysis;
    crawlBudget: CrawlBudget;
  };
  
  // Métriques globales
  metrics: {
    totalPages: number;
    pagesWithIssues: number;
    avgLoadTime: number;
    avgPageSize: number;
    totalIssues: number;
    issuesByType: Record<IssueType, number>;
    healthScore: number; // 0-100
  };
  
  // Configuration d'analyse
  config: {
    maxDepth: number;
    respectRobots: boolean;
    userAgent: string;
    crawlDelay: number;
    maxPages?: number;
    includedPaths?: string[];
    excludedPaths?: string[];
  };
  
  // Statut
  status: {
    analysisDate: string;
    crawlDuration: number; // en minutes
    crawlStatus: 'completed' | 'partial' | 'failed';
    lastUpdate: string;
    version: string;
  };
}

// Types utilitaires pour les composants UI
export interface PageAnalysisTableRow extends PageAnalysis {
  issueCount: number;
  issueTypes: IssueType[];
  overallHealth: 'good' | 'warning' | 'error';
}

export interface TechnicalIssueCard {
  type: IssueType;
  severity: 'low' | 'medium' | 'high' | 'critical';
  title: string;
  description: string;
  affectedPages: string[];
  recommendation: string;
  impact: string;
  effort: 'low' | 'medium' | 'high';
}

export interface CrawlabilityMetric {
  name: string;
  value: number | string | boolean;
  status: 'success' | 'warning' | 'error';
  description: string;
  recommendation?: string;
}

// Enums pour le filtrage et tri
export enum PageSortBy {
  URL = 'url',
  STATUS_CODE = 'statusCode',
  LOAD_TIME = 'loadTime',
  SIZE = 'size',
  ISSUES = 'issues',
  LAST_CRAWLED = 'lastCrawled',
}

export enum PageFilterBy {
  ALL = 'all',
  ERRORS_ONLY = 'errors',
  WARNINGS_ONLY = 'warnings',
  TITLE_ISSUES = 'title-issues',
  META_ISSUES = 'meta-issues',
  HEADING_ISSUES = 'heading-issues',
  IMAGE_ISSUES = 'image-issues',
  LINK_ISSUES = 'link-issues',
  SLOW_PAGES = 'slow-pages',
  LARGE_PAGES = 'large-pages',
}

export enum GlobalIssueType {
  DUPLICATE_CONTENT = 'duplicate-content',
  BROKEN_LINKS = 'broken-links',
  REDIRECT_CHAINS = 'redirect-chains',
  ORPHAN_PAGES = 'orphan-pages',
  MISSING_ELEMENTS = 'missing-elements',
}

// Types pour les requêtes API
export interface TechnicalAnalysisRequest {
  analysisId: number;
  includePerformance?: boolean;
  includeMobile?: boolean;
  maxPages?: number;
  respectRobots?: boolean;
}

export interface PageAnalysisRequest {
  url: string;
  includeImages?: boolean;
  includeLinks?: boolean;
  includeSchema?: boolean;
}

// Types pour l'export
export interface TechnicalAnalysisExport {
  format: 'pdf' | 'excel' | 'csv' | 'json';
  sections: Array<'pages' | 'issues' | 'crawlability' | 'metrics'>;
  filters?: {
    issueTypes?: IssueType[];
    severityLevel?: ('low' | 'medium' | 'high' | 'critical')[];
    pageStatus?: PageStatus[];
  };
}

// Utility types
export type PageURL = string;
export type ISODateString = string;
export type HTTPStatusCode = number;
export type ByteSize = number;
export type MillisecondsTime = number;