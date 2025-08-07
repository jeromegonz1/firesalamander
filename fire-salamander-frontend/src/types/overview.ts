/**
 * Fire Salamander - Overview Analysis TypeScript Interface
 * Lead Tech quality overview dashboard with comprehensive metrics
 */

export interface OverviewAnalysis {
  // Global Scores
  score: {
    overall: number; // 0-100
    grade: OverviewGrade;
    trend: 'improving' | 'stable' | 'declining';
    previousScore?: number;
    lastUpdated: string;
  };

  // Module Summary
  modules: {
    technical: ModuleSummary;
    performance: ModuleSummary;
    content: ModuleSummary;
    security: ModuleSummary;
    backlinks: ModuleSummary;
  };

  // Key Performance Indicators
  kpis: {
    organicTraffic: {
      current: number;
      change: number;
      trend: 'up' | 'down' | 'stable';
    };
    keywordRankings: {
      total: number;
      top10: number;
      top3: number;
      newRankings: number;
    };
    domainAuthority: {
      current: number;
      change: number;
      trend: 'up' | 'down' | 'stable';
    };
    pageSpeed: {
      mobile: number;
      desktop: number;
      change: number;
    };
  };

  // Critical Issues Summary
  issues: {
    critical: IssueItem[];
    warnings: IssueItem[];
    totalIssues: number;
    resolvedIssues: number;
  };

  // Recommendations Summary
  recommendations: {
    immediate: RecommendationItem[];
    quickWins: RecommendationItem[];
    totalRecommendations: number;
    estimatedImpact: string;
  };

  // Competition Comparison
  competition: {
    position: number;
    totalCompetitors: number;
    aboveAverage: boolean;
    marketShare: number;
    competitorGains: number;
  };

  // Progress Tracking
  progress: {
    completedTasks: number;
    totalTasks: number;
    lastActivity: string;
    nextMilestone: string;
  };

  // Site Information
  siteInfo: {
    url: string;
    title: string;
    description: string;
    lastCrawled: string;
    pagesAnalyzed: number;
    crawlErrors: number;
  };
}

export enum OverviewGrade {
  A_PLUS = 'A+',
  A = 'A',
  B = 'B',
  C = 'C',
  D = 'D',
  F = 'F',
}

export interface ModuleSummary {
  score: number;
  grade: OverviewGrade;
  status: 'completed' | 'in_progress' | 'pending' | 'error';
  lastUpdated: string;
  keyMetrics: Record<string, number | string>;
  criticalIssues: number;
  improvements: number;
}

export interface IssueItem {
  id: string;
  title: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  module: 'technical' | 'performance' | 'content' | 'security' | 'backlinks';
  description: string;
  impact: string;
  effort: 'low' | 'medium' | 'high';
  affectedPages: number;
  estimatedFix: string;
}

export interface RecommendationItem {
  id: string;
  title: string;
  category: 'quick-win' | 'immediate' | 'strategic';
  module: 'technical' | 'performance' | 'content' | 'security' | 'backlinks';
  description: string;
  expectedImpact: string;
  effort: 'low' | 'medium' | 'high';
  timeline: string;
  priority: number;
}

export interface HistoricalData {
  date: string;
  overallScore: number;
  technicalScore: number;
  performanceScore: number;
  contentScore: number;
  securityScore: number;
  backlinksScore: number;
  organicTraffic: number;
  rankings: number;
}

// Type guards
export const isOverviewAnalysis = (obj: any): obj is OverviewAnalysis => {
  return obj && 
    obj.score &&
    typeof obj.score.overall === 'number' &&
    obj.modules &&
    obj.kpis &&
    obj.issues;
};