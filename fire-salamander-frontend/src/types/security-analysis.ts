/**
 * Fire Salamander - Security Analysis Types
 * Professional Security Analysis with best practices
 * Lead Tech implementation - TDD approach
 */

// Enums pour type safety stricte
export enum SecurityGrade {
  A_PLUS = 'A+',
  A = 'A',
  B = 'B',
  C = 'C',
  D = 'D',
  F = 'F',
}

export enum SecuritySeverity {
  CRITICAL = 'critical',
  HIGH = 'high',
  MEDIUM = 'medium',
  LOW = 'low',
  INFO = 'info',
}

export enum SecurityCategory {
  HEADERS = 'headers',
  SSL_TLS = 'ssl-tls',
  COOKIES = 'cookies',
  CONTENT = 'content',
  INFRASTRUCTURE = 'infrastructure',
  AUTHENTICATION = 'authentication',
  AUTHORIZATION = 'authorization',
  DATA_PROTECTION = 'data-protection',
}

export enum HeaderStatus {
  PRESENT = 'present',
  MISSING = 'missing',
  MISCONFIGURED = 'misconfigured',
  DEPRECATED = 'deprecated',
}

export enum SSLProtocol {
  TLS_1_3 = 'TLS 1.3',
  TLS_1_2 = 'TLS 1.2',
  TLS_1_1 = 'TLS 1.1',
  TLS_1_0 = 'TLS 1.0',
  SSL_3_0 = 'SSL 3.0',
  SSL_2_0 = 'SSL 2.0',
}

export enum CipherStrength {
  STRONG = 'strong',
  MODERATE = 'moderate',
  WEAK = 'weak',
  INSECURE = 'insecure',
}

// ==================== INTERFACES PRINCIPALES ====================

export interface SecurityHeader {
  name: string;
  value: string | null;
  status: HeaderStatus;
  recommendation: string;
  impact: string;
  severity: SecuritySeverity;
  references: string[];
}

export interface SSLCertificate {
  issuer: string;
  subject: string;
  validFrom: string;
  validTo: string;
  daysRemaining: number;
  isValid: boolean;
  isExpired: boolean;
  isExpiringSoon: boolean; // < 30 jours
  algorithm: string;
  keySize: number;
  san: string[]; // Subject Alternative Names
  chainValid: boolean;
  ocspStatus: 'good' | 'revoked' | 'unknown';
}

export interface SSLConfiguration {
  protocols: {
    protocol: SSLProtocol;
    enabled: boolean;
    secure: boolean;
  }[];
  cipherSuites: {
    name: string;
    strength: CipherStrength;
    protocol: string;
    keyExchange: string;
    authentication: string;
    encryption: string;
    mac: string;
    export: boolean;
    recommended: boolean;
  }[];
  certificate: SSLCertificate;
  hsts: {
    enabled: boolean;
    maxAge: number;
    includeSubdomains: boolean;
    preload: boolean;
  };
  forwardSecrecy: boolean;
  ocspStapling: boolean;
  sessionResumption: boolean;
}

export interface CookieAnalysis {
  name: string;
  value: string;
  domain: string;
  path: string;
  expires: string | null;
  size: number;
  httpOnly: boolean;
  secure: boolean;
  sameSite: 'Strict' | 'Lax' | 'None' | null;
  issues: string[];
  recommendations: string[];
}

export interface ContentSecurityPolicy {
  present: boolean;
  directives: Record<string, string[]>;
  issues: {
    directive: string;
    issue: string;
    severity: SecuritySeverity;
    recommendation: string;
  }[];
  reportUri: string | null;
  upgradeInsecureRequests: boolean;
  blockAllMixedContent: boolean;
}

export interface SecurityVulnerability {
  id: string;
  type: string;
  severity: SecuritySeverity;
  category: SecurityCategory;
  title: string;
  description: string;
  affectedUrls: string[];
  cvss: {
    score: number;
    vector: string;
  };
  cwe: string; // Common Weakness Enumeration
  owasp: string[]; // OWASP Top 10 categories
  remediation: {
    summary: string;
    steps: string[];
    effort: 'low' | 'medium' | 'high';
    priority: number; // 1-10
  };
  references: {
    title: string;
    url: string;
  }[];
  exploitability: 'high' | 'medium' | 'low' | 'theoretical';
  dateDiscovered: string;
  falsePositive: boolean;
}

export interface SecurityBestPractice {
  id: string;
  category: SecurityCategory;
  practice: string;
  implemented: boolean;
  importance: 'critical' | 'high' | 'medium' | 'low';
  description: string;
  implementation: string;
  benefits: string[];
  references: string[];
}

export interface InfrastructureSecurity {
  webServer: {
    name: string;
    version: string;
    hideVersion: boolean;
    knownVulnerabilities: number;
  };
  technologies: {
    name: string;
    version: string;
    category: string;
    outdated: boolean;
    vulnerabilities: number;
    latestVersion: string;
  }[];
  openPorts: {
    port: number;
    service: string;
    secure: boolean;
    recommendation: string;
  }[];
  dnssec: boolean;
  ipv6: boolean;
  http2: boolean;
  http3: boolean;
}

export interface PermissionAnalysis {
  robotsTxt: {
    present: boolean;
    sensitive: string[]; // Sensitive paths exposed
    issues: string[];
  };
  sitemapXml: {
    present: boolean;
    accessible: boolean;
    sensitive: string[];
  };
  adminPaths: {
    found: string[];
    accessible: string[];
    protected: string[];
  };
  apiEndpoints: {
    discovered: string[];
    authenticated: string[];
    rateLimit: boolean;
  };
  backupFiles: string[];
  configFiles: string[];
  sourceCodeExposure: string[];
}

// ==================== INTERFACE PRINCIPALE ====================

export interface SecurityAnalysis {
  // Score global de sécurité
  score: {
    overall: number; // 0-100
    grade: SecurityGrade;
    trend: 'improving' | 'stable' | 'declining';
    previousScore?: number;
  };

  // Headers de sécurité
  headers: {
    securityHeaders: SecurityHeader[];
    score: number;
    grade: SecurityGrade;
    missing: string[];
    misconfigured: string[];
  };

  // SSL/TLS
  ssl: {
    enabled: boolean;
    configuration: SSLConfiguration | null;
    score: number;
    grade: SecurityGrade;
    issues: string[];
    strengths: string[];
  };

  // Cookies
  cookies: {
    list: CookieAnalysis[];
    totalCookies: number;
    secureCookies: number;
    httpOnlyCookies: number;
    sameSiteCookies: number;
    issues: string[];
  };

  // Content Security Policy
  csp: ContentSecurityPolicy;

  // Vulnérabilités détectées
  vulnerabilities: {
    critical: SecurityVulnerability[];
    high: SecurityVulnerability[];
    medium: SecurityVulnerability[];
    low: SecurityVulnerability[];
    info: SecurityVulnerability[];
    total: number;
  };

  // Best practices
  bestPractices: {
    implemented: SecurityBestPractice[];
    missing: SecurityBestPractice[];
    score: number;
  };

  // Infrastructure
  infrastructure: InfrastructureSecurity;

  // Permissions et accès
  permissions: PermissionAnalysis;

  // Métriques de conformité
  compliance: {
    gdpr: {
      compliant: boolean;
      issues: string[];
      score: number;
    };
    pci: {
      level: number;
      compliant: boolean;
      issues: string[];
    };
    iso27001: {
      controls: number;
      implemented: number;
      percentage: number;
    };
  };

  // Recommandations prioritaires
  recommendations: {
    immediate: {
      id: string;
      title: string;
      severity: SecuritySeverity;
      effort: 'low' | 'medium' | 'high';
      impact: 'low' | 'medium' | 'high';
      description: string;
      implementation: string;
    }[];
    shortTerm: any[];
    longTerm: any[];
  };

  // Métadonnées
  metadata: {
    scanDate: string;
    scanDuration: number; // secondes
    toolVersion: string;
    scanDepth: 'basic' | 'standard' | 'comprehensive';
    pagesScanned: number;
    errorsEncountered: number;
  };
}

// ==================== TYPES UTILITAIRES ====================

export interface SecurityIssueAction {
  issue: string;
  action: string;
  priority: number;
  effort: 'low' | 'medium' | 'high';
  codeExample?: string;
  documentation: string[];
}

export enum SecurityScanType {
  HEADERS = 'headers',
  SSL = 'ssl',
  VULNERABILITIES = 'vulnerabilities',
  PERMISSIONS = 'permissions',
  INFRASTRUCTURE = 'infrastructure',
  FULL = 'full',
}

// ==================== CONSTANTES ====================

export const SECURITY_HEADERS = {
  'Strict-Transport-Security': {
    required: true,
    defaultValue: 'max-age=31536000; includeSubDomains; preload',
  },
  'X-Content-Type-Options': {
    required: true,
    defaultValue: 'nosniff',
  },
  'X-Frame-Options': {
    required: true,
    defaultValue: 'DENY',
  },
  'X-XSS-Protection': {
    required: false, // Deprecated but still checked
    defaultValue: '1; mode=block',
  },
  'Content-Security-Policy': {
    required: true,
    defaultValue: "default-src 'self'",
  },
  'Referrer-Policy': {
    required: true,
    defaultValue: 'strict-origin-when-cross-origin',
  },
  'Permissions-Policy': {
    required: true,
    defaultValue: 'geolocation=(), microphone=(), camera=()',
  },
};

export const OWASP_TOP_10_2021 = {
  A01: 'Broken Access Control',
  A02: 'Cryptographic Failures',
  A03: 'Injection',
  A04: 'Insecure Design',
  A05: 'Security Misconfiguration',
  A06: 'Vulnerable and Outdated Components',
  A07: 'Identification and Authentication Failures',
  A08: 'Software and Data Integrity Failures',
  A09: 'Security Logging and Monitoring Failures',
  A10: 'Server-Side Request Forgery',
};

export const SECURITY_SCORE_THRESHOLDS = {
  [SecurityGrade.A_PLUS]: { min: 95, max: 100 },
  [SecurityGrade.A]: { min: 85, max: 94 },
  [SecurityGrade.B]: { min: 70, max: 84 },
  [SecurityGrade.C]: { min: 50, max: 69 },
  [SecurityGrade.D]: { min: 30, max: 49 },
  [SecurityGrade.F]: { min: 0, max: 29 },
};