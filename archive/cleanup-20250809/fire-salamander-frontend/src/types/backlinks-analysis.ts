/**
 * Fire Salamander - Backlinks Analysis TypeScript Interface
 * Enterprise-grade backlinks analysis with SEO authority metrics
 * Lead Tech Quality - Professional SEO backlinks system
 */

// ===============================
// ENUMS AND CONSTANTS
// ===============================

export enum BacklinkGrade {
  A_PLUS = 'A+',
  A = 'A',
  B = 'B',
  C = 'C',
  D = 'D',
  F = 'F',
}

export enum BacklinkQuality {
  EXCELLENT = 'excellent',
  GOOD = 'good',
  AVERAGE = 'average',
  POOR = 'poor',
  TOXIC = 'toxic',
}

export enum LinkType {
  DOFOLLOW = 'dofollow',
  NOFOLLOW = 'nofollow',
  SPONSORED = 'sponsored',
  UGC = 'ugc',
}

export enum AnchorType {
  EXACT_MATCH = 'exact_match',
  PARTIAL_MATCH = 'partial_match',
  BRANDED = 'branded',
  NAKED_URL = 'naked_url',
  GENERIC = 'generic',
  IMAGE = 'image',
  OTHER = 'other',
}

export enum BacklinkStatus {
  ACTIVE = 'active',
  LOST = 'lost',
  NEW = 'new',
  BROKEN = 'broken',
  REDIRECT = 'redirect',
  TOXIC = 'toxic',
}

export enum DomainCategory {
  NEWS = 'news',
  BLOG = 'blog',
  ECOMMERCE = 'ecommerce',
  CORPORATE = 'corporate',
  GOVERNMENT = 'government',
  EDUCATION = 'education',
  SOCIAL_MEDIA = 'social_media',
  FORUM = 'forum',
  DIRECTORY = 'directory',
  OTHER = 'other',
}

export enum ToxicityReason {
  SPAM = 'spam',
  LOW_QUALITY = 'low_quality',
  UNNATURAL_PATTERN = 'unnatural_pattern',
  SUSPICIOUS_ANCHOR = 'suspicious_anchor',
  PENALIZED_DOMAIN = 'penalized_domain',
  IRRELEVANT_CONTENT = 'irrelevant_content',
  LINK_FARM = 'link_farm',
  PBN = 'pbn',
}

// ===============================
// CORE INTERFACES
// ===============================

export interface BacklinkProfile {
  totalBacklinks: number;
  totalReferringDomains: number;
  totalReferringIPs: number;
  totalReferringSubnets: number;
  newBacklinks30Days: number;
  lostBacklinks30Days: number;
  brokenBacklinks: number;
  dofollowPercentage: number;
  nofollowPercentage: number;
  averageDomainAuthority: number;
  averagePageAuthority: number;
  trustFlow: number;
  citationFlow: number;
  majesticTrustRatio: number;
}

export interface DomainAuthority {
  domainAuthority: number; // DA score 0-100
  pageAuthority: number; // PA score 0-100
  domainRating: number; // Ahrefs DR 0-100
  urlRating: number; // Ahrefs UR 0-100
  trustFlow: number; // Majestic TF 0-100
  citationFlow: number; // Majestic CF 0-100
  spamScore: number; // Moz spam score 0-17
  organicTraffic: number; // Monthly organic traffic
  organicKeywords: number; // Number of ranking keywords
  lastUpdated: string;
}

export interface BacklinkData {
  id: string;
  sourceUrl: string;
  targetUrl: string;
  sourceDomain: string;
  targetDomain: string;
  anchorText: string;
  anchorType: AnchorType;
  linkType: LinkType;
  firstSeen: string;
  lastSeen: string;
  status: BacklinkStatus;
  quality: BacklinkQuality;
  domainAuthority: DomainAuthority;
  contextualRelevance: number; // 0-100
  linkPosition: {
    isMainContent: boolean;
    isNavigation: boolean;
    isFooter: boolean;
    isSidebar: boolean;
    position: number; // Position on page
  };
  suroundingText: string;
  pageTitle: string;
  pageCategory: DomainCategory;
  language: string;
  country: string;
  isRedirect: boolean;
  redirectChain?: string[];
  httpStatus: number;
  toxicityScore: number; // 0-100, higher = more toxic
  toxicityReasons: ToxicityReason[];
}

export interface AnchorTextDistribution {
  anchorText: string;
  type: AnchorType;
  count: number;
  percentage: number;
  averageDA: number;
  isNatural: boolean;
  isOverOptimized: boolean;
  riskScore: number; // 0-100
  topDomains: string[];
}

export interface ReferringDomain {
  domain: string;
  backlinksCount: number;
  firstSeen: string;
  lastSeen: string;
  domainAuthority: DomainAuthority;
  category: DomainCategory;
  language: string;
  country: string;
  isActive: boolean;
  quality: BacklinkQuality;
  linkTypes: {
    dofollow: number;
    nofollow: number;
    sponsored: number;
    ugc: number;
  };
  anchorDiversity: number; // 0-100, higher = more diverse
  contentRelevance: number; // 0-100
  toxicityScore: number;
  topPages: Array<{
    url: string;
    title: string;
    backlinksCount: number;
    pageAuthority: number;
  }>;
}

export interface CompetitorAnalysis {
  competitor: string;
  competitorDomain: string;
  totalBacklinks: number;
  referringDomains: number;
  domainAuthority: DomainAuthority;
  gap: {
    uniqueBacklinks: number;
    uniqueDomains: number;
    opportunities: BacklinkOpportunity[];
  };
  overlap: {
    commonBacklinks: number;
    commonDomains: number;
    sharedDomains: string[];
  };
  strengthComparison: {
    stronger: boolean;
    areas: string[];
    advantages: string[];
    weaknesses: string[];
  };
}

export interface BacklinkOpportunity {
  id: string;
  domain: string;
  url: string;
  title: string;
  domainAuthority: DomainAuthority;
  category: DomainCategory;
  contactInfo: {
    email?: string;
    contactPage?: string;
    socialMedia?: {
      twitter?: string;
      linkedin?: string;
      facebook?: string;
    };
  };
  outreachStatus: 'not_contacted' | 'contacted' | 'responded' | 'acquired' | 'rejected';
  priority: 'high' | 'medium' | 'low';
  estimatedDifficulty: number; // 0-100
  relevanceScore: number; // 0-100
  trafficValue: number; // Estimated monthly traffic value
  competitorLinks: string[]; // Which competitors have links from this domain
  suggestedApproach: string;
  contentOpportunities: string[];
}

export interface LinkBuildingInsights {
  topPerformingContent: Array<{
    url: string;
    title: string;
    backlinksCount: number;
    referringDomains: number;
    socialShares: number;
    contentType: string;
    publishDate: string;
    topics: string[];
  }>;
  contentGaps: Array<{
    topic: string;
    competitorCoverage: number;
    opportunity: number; // 0-100
    suggestedContentType: string;
    estimatedBacklinkPotential: number;
  }>;
  linkBaitOpportunities: Array<{
    type: 'infographic' | 'study' | 'tool' | 'guide' | 'statistics' | 'survey';
    topic: string;
    difficulty: number; // 0-100
    potential: number; // 0-100
    competitorExamples: string[];
    resources: string[];
  }>;
  industryBenchmarks: {
    averageBacklinks: number;
    averageReferringDomains: number;
    averageDomainAuthority: number;
    topPerformers: string[];
    industryTrends: Array<{
      trend: string;
      growth: number; // percentage
      timeframe: string;
    }>;
  };
}

export interface DisavowRecommendations {
  toxicBacklinks: BacklinkData[];
  recommendedActions: Array<{
    backlink: BacklinkData;
    action: 'disavow' | 'contact_removal' | 'monitor' | 'ignore';
    priority: 'high' | 'medium' | 'low';
    reason: string;
    potentialImpact: string;
  }>;
  disavowFile: {
    content: string;
    lastUpdated: string;
    domainsCount: number;
    urlsCount: number;
  };
  riskAssessment: {
    currentRisk: number; // 0-100
    riskAfterDisavow: number; // 0-100
    estimatedRecoveryTime: string;
    confidenceLevel: number; // 0-100
  };
}

export interface HistoricalTrends {
  timeline: Array<{
    date: string;
    totalBacklinks: number;
    referringDomains: number;
    domainAuthority: number;
    newBacklinks: number;
    lostBacklinks: number;
    netGain: number;
  }>;
  growthMetrics: {
    backlinksGrowthRate: number; // monthly percentage
    domainsGrowthRate: number; // monthly percentage
    authorityGrowthRate: number; // monthly percentage
    velocityScore: number; // 0-100, natural growth velocity
  };
  seasonalTrends: Array<{
    period: string;
    avgBacklinks: number;
    avgDomains: number;
    pattern: 'increasing' | 'decreasing' | 'stable' | 'seasonal';
  }>;
  milestones: Array<{
    date: string;
    event: string;
    impact: number; // positive or negative change
    description: string;
  }>;
}

// ===============================
// MAIN ANALYSIS INTERFACE
// ===============================

export interface BacklinksAnalysis {
  // Core Metrics
  score: {
    overall: number; // 0-100
    grade: BacklinkGrade;
    trend: 'improving' | 'stable' | 'declining';
    previousScore?: number;
    scoreBreakdown: {
      quantity: number; // Number of backlinks
      quality: number; // Average quality of backlinks  
      diversity: number; // Domain diversity
      authority: number; // Average domain authority
      naturalness: number; // Link naturalness score
      toxicity: number; // Inverse toxicity (100 - avg toxicity)
    };
  };

  // Backlink Profile Overview
  profile: BacklinkProfile;

  // Individual Backlinks
  backlinks: {
    list: BacklinkData[];
    totalCount: number;
    pagination: {
      currentPage: number;
      totalPages: number;
      pageSize: number;
    };
    filtering: {
      availableFilters: {
        quality: BacklinkQuality[];
        status: BacklinkStatus[];
        linkType: LinkType[];
        domainCategories: DomainCategory[];
        countries: string[];
        languages: string[];
      };
    };
  };

  // Anchor Text Analysis
  anchorTexts: {
    distribution: AnchorTextDistribution[];
    naturalness: {
      score: number; // 0-100
      isNatural: boolean;
      overOptimizedAnchors: string[];
      recommendations: string[];
    };
    diversity: {
      uniqueAnchors: number;
      diversityScore: number; // 0-100
      topAnchors: Array<{
        text: string;
        count: number;
        percentage: number;
      }>;
    };
  };

  // Referring Domains
  referringDomains: {
    list: ReferringDomain[];
    totalCount: number;
    qualityDistribution: Record<BacklinkQuality, number>;
    authorityDistribution: {
      high: number; // DA 80+
      medium: number; // DA 50-79
      low: number; // DA 0-49
    };
    categories: Record<DomainCategory, number>;
    geographicDistribution: Record<string, number>;
    topDomains: ReferringDomain[];
  };

  // Competitor Analysis
  competitors: {
    analysis: CompetitorAnalysis[];
    benchmarking: {
      position: number; // Your position among competitors
      totalCompetitors: number;
      aboveAverage: boolean;
      strengths: string[];
      opportunities: string[];
      threats: string[];
    };
    gapAnalysis: {
      totalOpportunities: number;
      highPriorityOpportunities: BacklinkOpportunity[];
      estimatedValue: number;
      timeToAcquire: string;
    };
  };

  // Link Building Opportunities
  opportunities: {
    list: BacklinkOpportunity[];
    totalCount: number;
    priorityDistribution: Record<'high' | 'medium' | 'low', number>;
    estimatedTotalValue: number;
    insights: LinkBuildingInsights;
  };

  // Toxic Links & Disavow
  toxicAnalysis: {
    toxicBacklinks: BacklinkData[];
    toxicDomains: string[];
    overallToxicity: number; // 0-100
    recommendations: DisavowRecommendations;
    riskFactors: Array<{
      factor: string;
      severity: 'high' | 'medium' | 'low';
      affectedLinks: number;
      description: string;
      impact: string;
    }>;
  };

  // Historical Data & Trends
  trends: HistoricalTrends;

  // Advanced Analytics
  analytics: {
    linkVelocity: {
      current: number; // links per month
      natural: number; // expected natural velocity
      isNatural: boolean;
      warning: string | null;
    };
    linkPatterns: {
      patternAnalysis: Array<{
        pattern: string;
        confidence: number; // 0-100
        risk: 'high' | 'medium' | 'low';
        description: string;
      }>;
      footprintDetection: Array<{
        footprint: string;
        occurrences: number;
        risk: number; // 0-100
        affectedLinks: string[];
      }>;
    };
    contentAnalysis: {
      topLinkingPages: Array<{
        url: string;
        title: string;
        linksCount: number;
        contentType: string;
        topics: string[];
      }>;
      linkingContext: Array<{
        context: string;
        frequency: number;
        naturalness: number; // 0-100
      }>;
    };
  };

  // Recommendations & Action Items
  recommendations: {
    immediate: Array<{
      id: string;
      title: string;
      priority: 'high' | 'medium' | 'low';
      category: 'acquisition' | 'cleanup' | 'optimization' | 'monitoring';
      description: string;
      expectedImpact: string;
      effort: 'low' | 'medium' | 'high';
      timeline: string;
      resources: string[];
    }>;
    strategic: Array<{
      id: string;
      title: string;
      goal: string;
      tactics: string[];
      kpis: string[];
      timeline: string;
      budget: string;
    }>;
  };

  // Metadata
  metadata: {
    scanDate: string;
    scanDuration: number; // seconds
    toolVersion: string;
    dataFreshness: number; // hours since last update
    coverage: {
      crawledDomains: number;
      totalFoundBacklinks: number;
      verifiedBacklinks: number;
      verificationRate: number; // percentage
    };
    limitations: string[];
    dataSources: string[];
  };
}

// ===============================
// UTILITY INTERFACES
// ===============================

export interface BacklinksFilters {
  quality?: BacklinkQuality[];
  status?: BacklinkStatus[];
  linkType?: LinkType[];
  domainAuthority?: {
    min: number;
    max: number;
  };
  dateRange?: {
    from: string;
    to: string;
  };
  domains?: string[];
  anchorTypes?: AnchorType[];
  countries?: string[];
  languages?: string[];
  categories?: DomainCategory[];
  toxicityThreshold?: number;
}

export interface BacklinksSorting {
  field: keyof BacklinkData | 'domainAuthority' | 'pageAuthority' | 'firstSeen' | 'lastSeen';
  direction: 'asc' | 'desc';
}

export interface BacklinksExport {
  format: 'csv' | 'xlsx' | 'pdf' | 'json';
  filters?: BacklinksFilters;
  fields: string[];
  includeMetadata: boolean;
}

// ===============================
// API RESPONSE INTERFACES
// ===============================

export interface BacklinksAPIResponse {
  success: boolean;
  data: BacklinksAnalysis;
  pagination?: {
    total: number;
    page: number;
    limit: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
  errors?: string[];
  warnings?: string[];
  metadata: {
    requestId: string;
    timestamp: string;
    processingTime: number;
    rateLimit: {
      remaining: number;
      reset: string;
    };
  };
}

export interface BacklinksBulkResponse {
  results: Array<{
    domain: string;
    analysis: BacklinksAnalysis | null;
    error?: string;
  }>;
  summary: {
    total: number;
    successful: number;
    failed: number;
    processingTime: number;
  };
}

// Type guards for runtime type checking
export const isBacklinkData = (obj: any): obj is BacklinkData => {
  return obj && 
    typeof obj.id === 'string' &&
    typeof obj.sourceUrl === 'string' &&
    typeof obj.targetUrl === 'string' &&
    Object.values(BacklinkStatus).includes(obj.status);
};

export const isBacklinksAnalysis = (obj: any): obj is BacklinksAnalysis => {
  return obj &&
    obj.score &&
    typeof obj.score.overall === 'number' &&
    Object.values(BacklinkGrade).includes(obj.score.grade) &&
    obj.profile &&
    obj.backlinks &&
    Array.isArray(obj.backlinks.list);
};