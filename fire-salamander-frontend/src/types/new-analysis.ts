/**
 * Fire Salamander - New Analysis TypeScript Interfaces
 * Lead Tech quality form system for creating new SEO analyses
 */

export interface NewAnalysisForm {
  // Basic Information
  url: string;
  name?: string;
  description?: string;
  
  // Analysis Configuration
  analysisType: AnalysisType;
  modules: AnalysisModule[];
  depth: AnalysisDepth;
  
  // Advanced Options
  crawlSettings: CrawlSettings;
  schedulingSettings: SchedulingSettings;
  notificationSettings: NotificationSettings;
  
  // Competition
  competitors?: string[];
  
  // Keywords (optional)
  targetKeywords?: string[];
  
  // Validation
  isValid?: boolean;
  validationErrors?: ValidationError[];
}

export enum AnalysisType {
  QUICK = 'quick',          // 15 minutes - basic analysis
  STANDARD = 'standard',    // 1 hour - comprehensive analysis
  DEEP = 'deep',           // 4 hours - in-depth analysis
  CONTINUOUS = 'continuous' // ongoing monitoring
}

export enum AnalysisModule {
  TECHNICAL = 'technical',
  PERFORMANCE = 'performance',
  CONTENT = 'content',
  SECURITY = 'security',
  BACKLINKS = 'backlinks',
  KEYWORDS = 'keywords',
  COMPETITION = 'competition',
}

export enum AnalysisDepth {
  SURFACE = 'surface',      // Homepage + key pages (5-10 pages)
  STANDARD = 'standard',    // Full site crawl (up to 1000 pages)
  DEEP = 'deep',           // Complete crawl (unlimited)
}

export interface CrawlSettings {
  maxPages: number;
  crawlDepth: number;
  respectRobots: boolean;
  includeSubdomains: boolean;
  excludePatterns: string[];
  includePatterns: string[];
  userAgent: string;
  crawlDelay: number; // seconds
  timeout: number;    // seconds
}

export interface SchedulingSettings {
  immediate: boolean;
  scheduledFor?: string; // ISO date
  recurring: boolean;
  frequency?: 'daily' | 'weekly' | 'monthly';
  timezone: string;
}

export interface NotificationSettings {
  email: boolean;
  emailAddress?: string;
  webhook: boolean;
  webhookUrl?: string;
  onCompletion: boolean;
  onIssues: boolean;
  onThresholds: boolean;
}

export interface ValidationError {
  field: string;
  message: string;
  severity: 'error' | 'warning';
}

export interface AnalysisEstimation {
  estimatedDuration: string;
  estimatedPages: number;
  estimatedCost: number;
  modulesIncluded: AnalysisModule[];
  features: string[];
}

export interface AnalysisPreset {
  id: string;
  name: string;
  description: string;
  type: AnalysisType;
  modules: AnalysisModule[];
  depth: AnalysisDepth;
  crawlSettings: Partial<CrawlSettings>;
  icon: string;
  popular?: boolean;
  recommended?: boolean;
}

// Default presets
export const ANALYSIS_PRESETS: AnalysisPreset[] = [
  {
    id: 'quick-check',
    name: 'Vérification Rapide',
    description: 'Analyse express des points critiques (15 min)',
    type: AnalysisType.QUICK,
    modules: [AnalysisModule.TECHNICAL, AnalysisModule.PERFORMANCE],
    depth: AnalysisDepth.SURFACE,
    crawlSettings: { maxPages: 10, crawlDepth: 2 },
    icon: 'Zap',
    popular: true,
  },
  {
    id: 'full-seo-audit',
    name: 'Audit SEO Complet',
    description: 'Analyse complète de tous les aspects SEO (1-2h)',
    type: AnalysisType.STANDARD,
    modules: [
      AnalysisModule.TECHNICAL,
      AnalysisModule.PERFORMANCE,
      AnalysisModule.CONTENT,
      AnalysisModule.SECURITY,
      AnalysisModule.BACKLINKS,
    ],
    depth: AnalysisDepth.STANDARD,
    crawlSettings: { maxPages: 1000, crawlDepth: 5 },
    icon: 'Search',
    recommended: true,
  },
  {
    id: 'security-focus',
    name: 'Audit Sécurité',
    description: 'Focus sur la sécurité et les vulnérabilités (30 min)',
    type: AnalysisType.STANDARD,
    modules: [AnalysisModule.SECURITY, AnalysisModule.TECHNICAL],
    depth: AnalysisDepth.STANDARD,
    crawlSettings: { maxPages: 500, crawlDepth: 3 },
    icon: 'Shield',
  },
  {
    id: 'performance-deep',
    name: 'Performance Avancée',
    description: 'Analyse poussée des performances (45 min)',
    type: AnalysisType.STANDARD,
    modules: [AnalysisModule.PERFORMANCE, AnalysisModule.TECHNICAL],
    depth: AnalysisDepth.DEEP,
    crawlSettings: { maxPages: 2000, crawlDepth: 4 },
    icon: 'Activity',
  },
  {
    id: 'competitor-analysis',
    name: 'Analyse Concurrentielle',
    description: 'Comparaison avec vos concurrents (2h)',
    type: AnalysisType.DEEP,
    modules: [
      AnalysisModule.BACKLINKS,
      AnalysisModule.KEYWORDS,
      AnalysisModule.CONTENT,
      AnalysisModule.COMPETITION,
    ],
    depth: AnalysisDepth.STANDARD,
    crawlSettings: { maxPages: 1000, crawlDepth: 3 },
    icon: 'Users',
  },
];

// Default form values
export const getDefaultFormValues = (): NewAnalysisForm => ({
  url: '',
  name: '',
  description: '',
  analysisType: AnalysisType.STANDARD,
  modules: [
    AnalysisModule.TECHNICAL,
    AnalysisModule.PERFORMANCE,
    AnalysisModule.CONTENT,
    AnalysisModule.SECURITY,
    AnalysisModule.BACKLINKS,
  ],
  depth: AnalysisDepth.STANDARD,
  crawlSettings: {
    maxPages: 1000,
    crawlDepth: 5,
    respectRobots: true,
    includeSubdomains: false,
    excludePatterns: [],
    includePatterns: [],
    userAgent: 'Fire Salamander SEO Crawler',
    crawlDelay: 1,
    timeout: 30,
  },
  schedulingSettings: {
    immediate: true,
    recurring: false,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  },
  notificationSettings: {
    email: true,
    webhook: false,
    onCompletion: true,
    onIssues: true,
    onThresholds: false,
  },
  competitors: [],
  targetKeywords: [],
});

// Form validation
export const validateAnalysisForm = (form: NewAnalysisForm): ValidationError[] => {
  const errors: ValidationError[] = [];

  // URL validation
  if (!form.url) {
    errors.push({
      field: 'url',
      message: 'L\'URL est obligatoire',
      severity: 'error',
    });
  } else {
    try {
      new URL(form.url);
      if (!form.url.startsWith('http://') && !form.url.startsWith('https://')) {
        errors.push({
          field: 'url',
          message: 'L\'URL doit commencer par http:// ou https://',
          severity: 'error',
        });
      }
    } catch {
      errors.push({
        field: 'url',
        message: 'Format d\'URL invalide',
        severity: 'error',
      });
    }
  }

  // Modules validation
  if (form.modules.length === 0) {
    errors.push({
      field: 'modules',
      message: 'Au moins un module doit être sélectionné',
      severity: 'error',
    });
  }

  // Email validation
  if (form.notificationSettings.email && form.notificationSettings.emailAddress) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(form.notificationSettings.emailAddress)) {
      errors.push({
        field: 'notificationSettings.emailAddress',
        message: 'Adresse email invalide',
        severity: 'error',
      });
    }
  }

  // Webhook validation
  if (form.notificationSettings.webhook && form.notificationSettings.webhookUrl) {
    try {
      new URL(form.notificationSettings.webhookUrl);
    } catch {
      errors.push({
        field: 'notificationSettings.webhookUrl',
        message: 'URL webhook invalide',
        severity: 'error',
      });
    }
  }

  // Crawl settings validation
  if (form.crawlSettings.maxPages < 1) {
    errors.push({
      field: 'crawlSettings.maxPages',
      message: 'Le nombre de pages doit être supérieur à 0',
      severity: 'error',
    });
  }

  if (form.crawlSettings.crawlDepth < 1) {
    errors.push({
      field: 'crawlSettings.crawlDepth',
      message: 'La profondeur doit être supérieure à 0',
      severity: 'error',
    });
  }

  // Warnings
  if (form.crawlSettings.maxPages > 10000) {
    errors.push({
      field: 'crawlSettings.maxPages',
      message: 'Un grand nombre de pages peut considérablement ralentir l\'analyse',
      severity: 'warning',
    });
  }

  return errors;
};