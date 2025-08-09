/**
 * Fire Salamander - Security Analysis Mapper Tests (TDD)
 * Lead Tech quality - Tests before mapper implementation
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  SecurityAnalysis,
  SecurityGrade,
  SecuritySeverity,
  HeaderStatus,
  SSLProtocol,
  CipherStrength,
} from '@/types/security-analysis';

// Mock backend data structure pour Security Analysis
interface MockBackendSecurityData {
  id: string;
  url: string;
  analyzed_at: string;
  security_data?: {
    headers?: {
      present: string[];
      missing: string[];
      headers_detail: Array<{
        name: string;
        value: string;
        secure: boolean;
        recommendations?: string[];
      }>;
    };
    ssl?: {
      enabled: boolean;
      protocol: string;
      certificate?: {
        issuer: string;
        subject: string;
        valid_from: string;
        valid_to: string;
        days_remaining: number;
        san: string[];
        chain_valid: boolean;
      };
      cipher_suites?: string[];
      vulnerabilities?: string[];
    };
    cookies?: Array<{
      name: string;
      value: string;
      domain: string;
      path: string;
      secure: boolean;
      http_only: boolean;
      same_site?: string;
    }>;
    vulnerabilities?: Array<{
      type: string;
      severity: string;
      title: string;
      description: string;
      affected_urls: string[];
      cvss_score?: number;
      cwe?: string;
      recommendations: string[];
    }>;
    technologies?: Array<{
      name: string;
      version: string;
      category: string;
      vulnerabilities: number;
    }>;
    permissions?: {
      robots_txt: {
        present: boolean;
        content?: string;
        disallowed_paths?: string[];
      };
      admin_paths?: string[];
      exposed_files?: string[];
    };
    score?: {
      total: number;
      headers: number;
      ssl: number;
      cookies: number;
      vulnerabilities: number;
    };
  };
}

describe('Security Analysis Mapper Tests (TDD)', () => {
  let mockBackendData: MockBackendSecurityData;

  beforeEach(() => {
    mockBackendData = {
      id: 'security-analysis-123',
      url: 'https://example.com',
      analyzed_at: '2024-03-15T10:30:00Z',
      security_data: {
        headers: {
          present: ['Strict-Transport-Security', 'X-Content-Type-Options'],
          missing: ['Content-Security-Policy', 'X-Frame-Options', 'Permissions-Policy'],
          headers_detail: [
            {
              name: 'Strict-Transport-Security',
              value: 'max-age=31536000',
              secure: true,
              recommendations: ['Add includeSubDomains directive', 'Consider HSTS preload'],
            },
            {
              name: 'X-Content-Type-Options',
              value: 'nosniff',
              secure: true,
              recommendations: [],
            },
          ],
        },
        ssl: {
          enabled: true,
          protocol: 'TLSv1.3',
          certificate: {
            issuer: 'Let\'s Encrypt Authority X3',
            subject: 'example.com',
            valid_from: '2024-01-01T00:00:00Z',
            valid_to: '2024-04-01T00:00:00Z',
            days_remaining: 17,
            san: ['example.com', 'www.example.com'],
            chain_valid: true,
          },
          cipher_suites: [
            'TLS_AES_256_GCM_SHA384',
            'TLS_CHACHA20_POLY1305_SHA256',
            'TLS_AES_128_GCM_SHA256',
          ],
          vulnerabilities: [],
        },
        cookies: [
          {
            name: 'session_id',
            value: 'encrypted_value',
            domain: '.example.com',
            path: '/',
            secure: true,
            http_only: true,
            same_site: 'Lax',
          },
          {
            name: 'tracking',
            value: 'ga_value',
            domain: '.example.com',
            path: '/',
            secure: false,
            http_only: false,
            same_site: undefined,
          },
        ],
        vulnerabilities: [
          {
            type: 'Missing Security Header',
            severity: 'high',
            title: 'Content Security Policy Header Missing',
            description: 'The Content-Security-Policy header is not set',
            affected_urls: ['https://example.com/'],
            cvss_score: 7.5,
            cwe: 'CWE-693',
            recommendations: [
              'Implement Content Security Policy',
              'Start with report-only mode',
              'Monitor violations before enforcing',
            ],
          },
          {
            type: 'Insecure Cookie',
            severity: 'medium',
            title: 'Cookie Without Secure Flag',
            description: 'Cookie "tracking" is missing the Secure flag',
            affected_urls: ['https://example.com/'],
            cvss_score: 4.3,
            cwe: 'CWE-614',
            recommendations: [
              'Set Secure flag on all cookies',
              'Review cookie usage',
            ],
          },
        ],
        technologies: [
          {
            name: 'nginx',
            version: '1.18.0',
            category: 'Web Server',
            vulnerabilities: 3,
          },
          {
            name: 'WordPress',
            version: '5.9.0',
            category: 'CMS',
            vulnerabilities: 7,
          },
        ],
        permissions: {
          robots_txt: {
            present: true,
            content: 'User-agent: *\nDisallow: /admin\nDisallow: /backup',
            disallowed_paths: ['/admin', '/backup'],
          },
          admin_paths: ['/admin', '/wp-admin', '/administrator'],
          exposed_files: ['/.git/config', '/backup.sql'],
        },
        score: {
          total: 72,
          headers: 40,
          ssl: 85,
          cookies: 60,
          vulnerabilities: 65,
        },
      },
    };
  });

  describe('mapBackendToSecurityAnalysis function', () => {
    it('should exist and be callable', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      expect(typeof mapBackendToSecurityAnalysis).toBe('function');
    });

    it('should map backend data to SecurityAnalysis structure', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.score).toBeDefined();
      expect(result.headers).toBeDefined();
      expect(result.ssl).toBeDefined();
      expect(result.cookies).toBeDefined();
      expect(result.vulnerabilities).toBeDefined();
    });

    it('should calculate security scores correctly', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.score.overall).toBe(72);
      expect(result.score.grade).toBe(SecurityGrade.B);
      expect(result.headers.score).toBe(40);
      expect(result.ssl.score).toBe(85);
    });

    it('should map security headers correctly', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.headers.securityHeaders).toHaveLength(5); // 2 present + 3 missing
      expect(result.headers.missing).toContain('Content-Security-Policy');
      expect(result.headers.missing).toContain('X-Frame-Options');
      
      const hstsHeader = result.headers.securityHeaders.find(h => h.name === 'Strict-Transport-Security');
      expect(hstsHeader).toBeDefined();
      expect(hstsHeader!.status).toBe(HeaderStatus.PRESENT);
      expect(hstsHeader!.value).toBe('max-age=31536000');
    });

    it('should map SSL/TLS configuration', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.ssl.enabled).toBe(true);
      expect(result.ssl.configuration).toBeDefined();
      expect(result.ssl.configuration!.certificate.daysRemaining).toBe(17);
      expect(result.ssl.configuration!.certificate.isExpiringSoon).toBe(true);
      
      const tls13 = result.ssl.configuration!.protocols.find(p => p.protocol === SSLProtocol.TLS_1_3);
      expect(tls13).toBeDefined();
      expect(tls13!.enabled).toBe(true);
      expect(tls13!.secure).toBe(true);
    });

    it('should analyze cookies security', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.cookies.list).toHaveLength(2);
      expect(result.cookies.totalCookies).toBe(2);
      expect(result.cookies.secureCookies).toBe(1);
      expect(result.cookies.httpOnlyCookies).toBe(1);
      
      const trackingCookie = result.cookies.list.find(c => c.name === 'tracking');
      expect(trackingCookie).toBeDefined();
      expect(trackingCookie!.secure).toBe(false);
      expect(trackingCookie!.issues).toContain('Missing Secure flag');
    });

    it('should categorize vulnerabilities by severity', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.vulnerabilities.high).toHaveLength(1);
      expect(result.vulnerabilities.medium).toHaveLength(1);
      expect(result.vulnerabilities.total).toBe(2);
      
      const cspVuln = result.vulnerabilities.high[0];
      expect(cspVuln.title).toContain('Content Security Policy');
      expect(cspVuln.cvss.score).toBe(7.5);
      expect(cspVuln.cwe).toBe('CWE-693');
    });

    it('should detect technology vulnerabilities', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.infrastructure.technologies).toHaveLength(2);
      
      const wordpress = result.infrastructure.technologies.find(t => t.name === 'WordPress');
      expect(wordpress).toBeDefined();
      expect(wordpress!.vulnerabilities).toBe(7);
      expect(wordpress!.outdated).toBe(true); // Should detect as outdated
    });

    it('should analyze permission issues', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.permissions.robotsTxt.present).toBe(true);
      expect(result.permissions.robotsTxt.sensitive).toContain('/admin');
      expect(result.permissions.robotsTxt.sensitive).toContain('/backup');
      
      expect(result.permissions.adminPaths.found).toContain('/wp-admin');
      expect(result.permissions.backupFiles).toContain('/backup.sql');
    });

    it('should generate security recommendations', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.recommendations.immediate.length).toBeGreaterThan(0);
      
      const cspRecommendation = result.recommendations.immediate.find(r => 
        r.title.includes('Content Security Policy')
      );
      expect(cspRecommendation).toBeDefined();
      expect(cspRecommendation!.severity).toBe(SecuritySeverity.HIGH);
      expect(cspRecommendation!.impact).toBe('high');
    });

    it('should handle empty or invalid data gracefully', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const emptyData = {};
      const result = mapBackendToSecurityAnalysis(emptyData);
      
      expect(result.score.overall).toBe(0);
      expect(result.score.grade).toBe(SecurityGrade.F);
      expect(result.headers.securityHeaders).toHaveLength(0);
      expect(result.vulnerabilities.total).toBe(0);
    });

    it('should set metadata correctly', async () => {
      const { mapBackendToSecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const result = mapBackendToSecurityAnalysis(mockBackendData);
      
      expect(result.metadata.scanDate).toBeTruthy();
      expect(result.metadata.toolVersion).toBe('1.0.0');
      expect(result.metadata.scanDepth).toBe('comprehensive');
    });
  });

  describe('Helper functions', () => {
    it('should have a function to determine security grade', async () => {
      const { getSecurityGrade } = await import('@/lib/mappers/security-mapper');
      
      expect(getSecurityGrade(97)).toBe(SecurityGrade.A_PLUS);
      expect(getSecurityGrade(87)).toBe(SecurityGrade.A);
      expect(getSecurityGrade(72)).toBe(SecurityGrade.B);
      expect(getSecurityGrade(55)).toBe(SecurityGrade.C);
      expect(getSecurityGrade(35)).toBe(SecurityGrade.D);
      expect(getSecurityGrade(15)).toBe(SecurityGrade.F);
    });

    it('should have a function to map header status', async () => {
      const { mapHeaderStatus } = await import('@/lib/mappers/security-mapper');
      
      expect(mapHeaderStatus('present')).toBe(HeaderStatus.PRESENT);
      expect(mapHeaderStatus('missing')).toBe(HeaderStatus.MISSING);
      expect(mapHeaderStatus('misconfigured')).toBe(HeaderStatus.MISCONFIGURED);
      expect(mapHeaderStatus('unknown')).toBe(HeaderStatus.MISSING); // fallback
    });

    it('should have a function to map severity', async () => {
      const { mapSeverity } = await import('@/lib/mappers/security-mapper');
      
      expect(mapSeverity('critical')).toBe(SecuritySeverity.CRITICAL);
      expect(mapSeverity('high')).toBe(SecuritySeverity.HIGH);
      expect(mapSeverity('medium')).toBe(SecuritySeverity.MEDIUM);
      expect(mapSeverity('low')).toBe(SecuritySeverity.LOW);
      expect(mapSeverity('info')).toBe(SecuritySeverity.INFO);
      expect(mapSeverity('unknown')).toBe(SecuritySeverity.INFO); // fallback
    });

    it('should have a function to analyze cookie security', async () => {
      const { analyzeCookieSecurity } = await import('@/lib/mappers/security-mapper');
      
      const cookie = {
        name: 'test',
        secure: false,
        httpOnly: false,
        sameSite: null,
      };
      
      const issues = analyzeCookieSecurity(cookie);
      
      expect(issues).toContain('Missing Secure flag');
      expect(issues).toContain('Missing HttpOnly flag');
      expect(issues).toContain('Missing SameSite attribute');
    });

    it('should have a function to detect sensitive paths', async () => {
      const { detectSensitivePaths } = await import('@/lib/mappers/security-mapper');
      
      const paths = ['/home', '/admin', '/backup', '/api/public', '/.git'];
      const sensitive = detectSensitivePaths(paths);
      
      expect(sensitive).toContain('/admin');
      expect(sensitive).toContain('/backup');
      expect(sensitive).toContain('/.git');
      expect(sensitive).not.toContain('/home');
    });

    it('should have a function to calculate compliance scores', async () => {
      const { calculateComplianceScores } = await import('@/lib/mappers/security-mapper');
      
      const securityData = {
        headers: { score: 60 },
        ssl: { score: 80 },
        cookies: { secureCookies: 3, totalCookies: 5 },
      };
      
      const compliance = calculateComplianceScores(securityData);
      
      expect(compliance.gdpr.score).toBeGreaterThanOrEqual(0);
      expect(compliance.gdpr.score).toBeLessThanOrEqual(100);
      expect(compliance.gdpr.issues.length).toBeGreaterThan(0);
    });

    it('should create empty security analysis for fallback', async () => {
      const { createEmptySecurityAnalysis } = await import('@/lib/mappers/security-mapper');
      
      const emptyAnalysis = createEmptySecurityAnalysis();
      
      expect(emptyAnalysis.score.overall).toBe(0);
      expect(emptyAnalysis.score.grade).toBe(SecurityGrade.F);
      expect(emptyAnalysis.headers.securityHeaders).toHaveLength(0);
      expect(emptyAnalysis.vulnerabilities.total).toBe(0);
      expect(emptyAnalysis.metadata.scanDate).toBeTruthy();
    });
  });

  describe('Security Score Calculations', () => {
    it('should calculate header score based on presence and configuration', async () => {
      const { calculateHeaderScore } = await import('@/lib/mappers/security-mapper');
      
      const headers = {
        present: ['Strict-Transport-Security', 'X-Content-Type-Options'],
        missing: ['Content-Security-Policy', 'X-Frame-Options', 'Permissions-Policy'],
        misconfigured: ['Strict-Transport-Security'], // missing preload
      };
      
      const score = calculateHeaderScore(headers);
      
      expect(score).toBeGreaterThan(0);
      expect(score).toBeLessThan(100);
      expect(score).toBe(40); // 2/5 present, 1 misconfigured
    });

    it('should calculate SSL score based on configuration', async () => {
      const { calculateSSLScore } = await import('@/lib/mappers/security-mapper');
      
      const sslConfig = {
        protocols: ['TLSv1.3', 'TLSv1.2'],
        certificate: { daysRemaining: 17, chainValid: true },
        cipherSuites: ['TLS_AES_256_GCM_SHA384'],
        vulnerabilities: [],
      };
      
      const score = calculateSSLScore(sslConfig);
      
      expect(score).toBeGreaterThan(70); // Good protocols but cert expiring soon
      expect(score).toBeLessThan(100);
    });

    it('should prioritize security recommendations', async () => {
      const { prioritizeRecommendations } = await import('@/lib/mappers/security-mapper');
      
      const issues = [
        { severity: SecuritySeverity.MEDIUM, effort: 'low' },
        { severity: SecuritySeverity.CRITICAL, effort: 'high' },
        { severity: SecuritySeverity.HIGH, effort: 'low' },
        { severity: SecuritySeverity.HIGH, effort: 'medium' },
      ];
      
      const prioritized = prioritizeRecommendations(issues);
      
      // Should prioritize by severity first, then by effort
      expect(prioritized[0].severity).toBe(SecuritySeverity.CRITICAL);
      expect(prioritized[1].severity).toBe(SecuritySeverity.HIGH);
      expect(prioritized[1].effort).toBe('low'); // Lower effort comes first for same severity
    });
  });
});