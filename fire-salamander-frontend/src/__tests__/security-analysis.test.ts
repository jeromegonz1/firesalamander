/**
 * Fire Salamander - Security Analysis Tests (TDD)
 * Lead Tech quality - Tests-first approach for Security Analysis
 * Professional security testing with OWASP standards
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  SecurityAnalysis,
  SecurityGrade,
  SecuritySeverity,
  SecurityCategory,
  HeaderStatus,
  SSLProtocol,
  CipherStrength,
  SecurityHeader,
  SSLCertificate,
  SecurityVulnerability,
  SECURITY_HEADERS,
  OWASP_TOP_10_2021,
  SECURITY_SCORE_THRESHOLDS,
} from '@/types/security-analysis';

describe('Security Analysis Types & Structure', () => {
  describe('SecurityAnalysis Interface', () => {
    let mockSecurityAnalysis: SecurityAnalysis;

    beforeEach(() => {
      mockSecurityAnalysis = {
        score: {
          overall: 78,
          grade: SecurityGrade.B,
          trend: 'improving',
          previousScore: 72,
        },
        
        headers: {
          securityHeaders: [
            {
              name: 'Strict-Transport-Security',
              value: 'max-age=31536000; includeSubDomains',
              status: HeaderStatus.PRESENT,
              recommendation: 'Add preload directive for HSTS preload list inclusion',
              impact: 'Enforces HTTPS connections, preventing protocol downgrade attacks',
              severity: SecuritySeverity.MEDIUM,
              references: [
                'https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security',
                'https://hstspreload.org/',
              ],
            },
            {
              name: 'Content-Security-Policy',
              value: null,
              status: HeaderStatus.MISSING,
              recommendation: 'Implement a Content Security Policy to prevent XSS attacks',
              impact: 'Without CSP, the site is vulnerable to XSS and data injection attacks',
              severity: SecuritySeverity.HIGH,
              references: [
                'https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP',
                'https://csp.withgoogle.com/docs/index.html',
              ],
            },
          ],
          score: 65,
          grade: SecurityGrade.C,
          missing: ['Content-Security-Policy', 'Permissions-Policy'],
          misconfigured: ['X-Frame-Options'],
        },
        
        ssl: {
          enabled: true,
          configuration: {
            protocols: [
              { protocol: SSLProtocol.TLS_1_3, enabled: true, secure: true },
              { protocol: SSLProtocol.TLS_1_2, enabled: true, secure: true },
              { protocol: SSLProtocol.TLS_1_1, enabled: false, secure: false },
            ],
            cipherSuites: [
              {
                name: 'TLS_AES_256_GCM_SHA384',
                strength: CipherStrength.STRONG,
                protocol: 'TLS 1.3',
                keyExchange: 'ECDHE',
                authentication: 'RSA',
                encryption: 'AES256-GCM',
                mac: 'SHA384',
                export: false,
                recommended: true,
              },
            ],
            certificate: {
              issuer: 'Let\'s Encrypt Authority X3',
              subject: 'www.example.com',
              validFrom: '2024-01-01T00:00:00Z',
              validTo: '2024-04-01T00:00:00Z',
              daysRemaining: 25,
              isValid: true,
              isExpired: false,
              isExpiringSoon: true,
              algorithm: 'RSA',
              keySize: 2048,
              san: ['www.example.com', 'example.com'],
              chainValid: true,
              ocspStatus: 'good',
            },
            hsts: {
              enabled: true,
              maxAge: 31536000,
              includeSubdomains: true,
              preload: false,
            },
            forwardSecrecy: true,
            ocspStapling: true,
            sessionResumption: true,
          },
          score: 85,
          grade: SecurityGrade.A,
          issues: ['Certificate expiring soon', 'HSTS preload not enabled'],
          strengths: ['TLS 1.3 enabled', 'Perfect forward secrecy', 'OCSP stapling'],
        },
        
        cookies: {
          list: [
            {
              name: 'sessionId',
              value: 'abc123...',
              domain: '.example.com',
              path: '/',
              expires: null,
              size: 128,
              httpOnly: true,
              secure: true,
              sameSite: 'Lax',
              issues: [],
              recommendations: ['Consider using SameSite=Strict for session cookies'],
            },
            {
              name: 'analytics',
              value: 'ga123...',
              domain: '.example.com',
              path: '/',
              expires: '2025-01-01T00:00:00Z',
              size: 64,
              httpOnly: false,
              secure: false,
              sameSite: null,
              issues: ['Missing Secure flag', 'Missing SameSite attribute'],
              recommendations: [
                'Add Secure flag to ensure cookie is only sent over HTTPS',
                'Set SameSite attribute to prevent CSRF attacks',
              ],
            },
          ],
          totalCookies: 2,
          secureCookies: 1,
          httpOnlyCookies: 1,
          sameSiteCookies: 1,
          issues: ['Some cookies lack Secure flag', 'Some cookies lack SameSite attribute'],
        },
        
        csp: {
          present: false,
          directives: {},
          issues: [
            {
              directive: 'default-src',
              issue: 'No Content Security Policy detected',
              severity: SecuritySeverity.HIGH,
              recommendation: 'Implement CSP starting with: default-src \'self\'',
            },
          ],
          reportUri: null,
          upgradeInsecureRequests: false,
          blockAllMixedContent: false,
        },
        
        vulnerabilities: {
          critical: [],
          high: [
            {
              id: 'VULN-001',
              type: 'Missing Security Headers',
              severity: SecuritySeverity.HIGH,
              category: SecurityCategory.HEADERS,
              title: 'Missing Content Security Policy',
              description: 'The application does not implement a Content Security Policy, making it vulnerable to XSS attacks',
              affectedUrls: ['https://example.com/*'],
              cvss: {
                score: 7.5,
                vector: 'CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:H/A:N',
              },
              cwe: 'CWE-693',
              owasp: ['A05'],
              remediation: {
                summary: 'Implement a Content Security Policy header',
                steps: [
                  'Start with a report-only CSP to monitor violations',
                  'Analyze reports and adjust policy',
                  'Implement enforcing CSP with appropriate directives',
                  'Monitor and refine policy over time',
                ],
                effort: 'medium',
                priority: 8,
              },
              references: [
                { title: 'CSP Guide', url: 'https://csp.withgoogle.com/' },
                { title: 'MDN CSP Docs', url: 'https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP' },
              ],
              exploitability: 'high',
              dateDiscovered: '2024-03-15',
              falsePositive: false,
            },
          ],
          medium: [],
          low: [],
          info: [],
          total: 1,
        },
        
        bestPractices: {
          implemented: [
            {
              id: 'BP-001',
              category: SecurityCategory.SSL_TLS,
              practice: 'Use HTTPS everywhere',
              implemented: true,
              importance: 'critical',
              description: 'All pages are served over HTTPS',
              implementation: 'SSL certificate installed and HTTP redirects to HTTPS',
              benefits: [
                'Encrypts data in transit',
                'Prevents man-in-the-middle attacks',
                'Required for modern web features',
              ],
              references: ['https://developers.google.com/web/fundamentals/security/encrypt-in-transit/why-https'],
            },
          ],
          missing: [
            {
              id: 'BP-002',
              category: SecurityCategory.HEADERS,
              practice: 'Implement security headers',
              implemented: false,
              importance: 'high',
              description: 'Security headers provide defense in depth',
              implementation: 'Configure web server to send security headers',
              benefits: [
                'Prevents clickjacking',
                'Mitigates XSS attacks',
                'Controls resource loading',
              ],
              references: ['https://securityheaders.com/'],
            },
          ],
          score: 70,
        },
        
        infrastructure: {
          webServer: {
            name: 'nginx',
            version: '1.18.0',
            hideVersion: false,
            knownVulnerabilities: 2,
          },
          technologies: [
            {
              name: 'PHP',
              version: '7.4.3',
              category: 'Language',
              outdated: true,
              vulnerabilities: 5,
              latestVersion: '8.3.0',
            },
          ],
          openPorts: [
            {
              port: 443,
              service: 'HTTPS',
              secure: true,
              recommendation: 'Keep open for HTTPS traffic',
            },
            {
              port: 80,
              service: 'HTTP',
              secure: true,
              recommendation: 'Keep open but redirect all traffic to HTTPS',
            },
          ],
          dnssec: false,
          ipv6: true,
          http2: true,
          http3: false,
        },
        
        permissions: {
          robotsTxt: {
            present: true,
            sensitive: ['/admin', '/backup'],
            issues: ['Exposes sensitive paths in robots.txt'],
          },
          sitemapXml: {
            present: true,
            accessible: true,
            sensitive: [],
          },
          adminPaths: {
            found: ['/admin', '/wp-admin'],
            accessible: ['/admin'],
            protected: ['/wp-admin'],
          },
          apiEndpoints: {
            discovered: ['/api/v1/users', '/api/v1/config'],
            authenticated: ['/api/v1/users'],
            rateLimit: false,
          },
          backupFiles: ['backup.sql.gz'],
          configFiles: [],
          sourceCodeExposure: [],
        },
        
        compliance: {
          gdpr: {
            compliant: false,
            issues: ['No cookie consent banner', 'Privacy policy not comprehensive'],
            score: 60,
          },
          pci: {
            level: 4,
            compliant: false,
            issues: ['Weak cryptography', 'Missing security headers'],
          },
          iso27001: {
            controls: 114,
            implemented: 78,
            percentage: 68,
          },
        },
        
        recommendations: {
          immediate: [
            {
              id: 'REC-001',
              title: 'Implement Content Security Policy',
              severity: SecuritySeverity.HIGH,
              effort: 'medium',
              impact: 'high',
              description: 'Add CSP header to prevent XSS attacks',
              implementation: 'Start with report-only mode, then enforce',
            },
          ],
          shortTerm: [],
          longTerm: [],
        },
        
        metadata: {
          scanDate: new Date().toISOString(),
          scanDuration: 180,
          toolVersion: '1.0.0',
          scanDepth: 'comprehensive',
          pagesScanned: 150,
          errorsEncountered: 3,
        },
      };
    });

    it('should have valid security analysis structure', () => {
      expect(mockSecurityAnalysis.score).toBeDefined();
      expect(mockSecurityAnalysis.score.overall).toBeGreaterThanOrEqual(0);
      expect(mockSecurityAnalysis.score.overall).toBeLessThanOrEqual(100);
      expect(Object.values(SecurityGrade)).toContain(mockSecurityAnalysis.score.grade);
    });

    it('should validate security score thresholds', () => {
      const score = mockSecurityAnalysis.score.overall;
      const grade = mockSecurityAnalysis.score.grade;
      const threshold = SECURITY_SCORE_THRESHOLDS[grade];
      
      expect(score).toBeGreaterThanOrEqual(threshold.min);
      expect(score).toBeLessThanOrEqual(threshold.max);
    });

    it('should have valid security headers analysis', () => {
      const { headers } = mockSecurityAnalysis;
      
      expect(headers.securityHeaders).toHaveLength(2);
      expect(headers.missing).toContain('Content-Security-Policy');
      expect(headers.score).toBeGreaterThanOrEqual(0);
      expect(headers.score).toBeLessThanOrEqual(100);
    });

    it('should validate SSL configuration', () => {
      const { ssl } = mockSecurityAnalysis;
      
      expect(ssl.enabled).toBe(true);
      expect(ssl.configuration).toBeDefined();
      expect(ssl.configuration!.protocols).toHaveLength(3);
      expect(ssl.configuration!.certificate.isValid).toBe(true);
      expect(ssl.configuration!.certificate.isExpiringSoon).toBe(true);
      expect(ssl.configuration!.certificate.daysRemaining).toBeLessThan(30);
    });

    it('should properly categorize vulnerabilities', () => {
      const { vulnerabilities } = mockSecurityAnalysis;
      
      expect(vulnerabilities.high).toHaveLength(1);
      expect(vulnerabilities.total).toBe(1);
      expect(vulnerabilities.high[0].cvss.score).toBeGreaterThanOrEqual(7.0);
      expect(vulnerabilities.high[0].owasp).toContain('A05');
    });

    it('should track cookie security attributes', () => {
      const { cookies } = mockSecurityAnalysis;
      
      expect(cookies.totalCookies).toBe(2);
      expect(cookies.secureCookies).toBe(1);
      expect(cookies.httpOnlyCookies).toBe(1);
      expect(cookies.list[1].issues).toContain('Missing Secure flag');
    });

    it('should analyze infrastructure security', () => {
      const { infrastructure } = mockSecurityAnalysis;
      
      expect(infrastructure.webServer.knownVulnerabilities).toBeGreaterThan(0);
      expect(infrastructure.technologies[0].outdated).toBe(true);
      expect(infrastructure.http2).toBe(true);
      expect(infrastructure.dnssec).toBe(false);
    });

    it('should detect permission and access issues', () => {
      const { permissions } = mockSecurityAnalysis;
      
      expect(permissions.robotsTxt.sensitive).toContain('/admin');
      expect(permissions.adminPaths.accessible).toContain('/admin');
      expect(permissions.backupFiles).toHaveLength(1);
      expect(permissions.apiEndpoints.rateLimit).toBe(false);
    });

    it('should evaluate compliance status', () => {
      const { compliance } = mockSecurityAnalysis;
      
      expect(compliance.gdpr.compliant).toBe(false);
      expect(compliance.gdpr.issues).toHaveLength(2);
      expect(compliance.pci.compliant).toBe(false);
      expect(compliance.iso27001.percentage).toBe(68);
    });

    it('should provide prioritized recommendations', () => {
      const { recommendations } = mockSecurityAnalysis;
      
      expect(recommendations.immediate).toHaveLength(1);
      expect(recommendations.immediate[0].severity).toBe(SecuritySeverity.HIGH);
      expect(recommendations.immediate[0].impact).toBe('high');
    });
  });

  describe('Security Header Validation', () => {
    it('should validate required security headers', () => {
      const requiredHeaders = Object.entries(SECURITY_HEADERS)
        .filter(([_, config]) => config.required)
        .map(([name, _]) => name);
      
      expect(requiredHeaders).toContain('Strict-Transport-Security');
      expect(requiredHeaders).toContain('Content-Security-Policy');
      expect(requiredHeaders).toContain('X-Content-Type-Options');
      expect(requiredHeaders).toContain('X-Frame-Options');
    });

    it('should have proper header status values', () => {
      const header: SecurityHeader = {
        name: 'X-Frame-Options',
        value: 'SAMEORIGIN',
        status: HeaderStatus.MISCONFIGURED,
        recommendation: 'Use DENY instead of SAMEORIGIN',
        impact: 'Current setting still allows framing from same origin',
        severity: SecuritySeverity.MEDIUM,
        references: ['https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options'],
      };
      
      expect(Object.values(HeaderStatus)).toContain(header.status);
      expect(header.references).toHaveLength(1);
    });
  });

  describe('SSL/TLS Analysis', () => {
    it('should identify secure protocols', () => {
      const protocols = [
        { protocol: SSLProtocol.TLS_1_3, secure: true },
        { protocol: SSLProtocol.TLS_1_2, secure: true },
        { protocol: SSLProtocol.TLS_1_1, secure: false },
        { protocol: SSLProtocol.TLS_1_0, secure: false },
        { protocol: SSLProtocol.SSL_3_0, secure: false },
      ];
      
      const secureProtocols = protocols.filter(p => p.secure);
      expect(secureProtocols).toHaveLength(2);
      expect(secureProtocols.map(p => p.protocol)).toContain(SSLProtocol.TLS_1_3);
    });

    it('should validate certificate expiration warnings', () => {
      const cert: SSLCertificate = {
        issuer: 'Test CA',
        subject: 'test.com',
        validFrom: '2024-01-01',
        validTo: '2024-04-01',
        daysRemaining: 25,
        isValid: true,
        isExpired: false,
        isExpiringSoon: true,
        algorithm: 'RSA',
        keySize: 2048,
        san: ['test.com', 'www.test.com'],
        chainValid: true,
        ocspStatus: 'good',
      };
      
      expect(cert.isExpiringSoon).toBe(cert.daysRemaining < 30);
      expect(cert.keySize).toBeGreaterThanOrEqual(2048);
    });

    it('should categorize cipher strength', () => {
      const cipherStrengths = Object.values(CipherStrength);
      
      expect(cipherStrengths).toContain(CipherStrength.STRONG);
      expect(cipherStrengths).toContain(CipherStrength.WEAK);
      expect(cipherStrengths).toContain(CipherStrength.INSECURE);
    });
  });

  describe('Vulnerability Analysis', () => {
    it('should properly score vulnerabilities', () => {
      const vuln: SecurityVulnerability = {
        id: 'TEST-001',
        type: 'XSS',
        severity: SecuritySeverity.HIGH,
        category: SecurityCategory.CONTENT,
        title: 'Cross-Site Scripting',
        description: 'Reflected XSS vulnerability',
        affectedUrls: ['/search'],
        cvss: {
          score: 7.5,
          vector: 'CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:H/A:N',
        },
        cwe: 'CWE-79',
        owasp: ['A03'],
        remediation: {
          summary: 'Sanitize user input',
          steps: ['Validate input', 'Encode output'],
          effort: 'low',
          priority: 9,
        },
        references: [],
        exploitability: 'high',
        dateDiscovered: '2024-03-15',
        falsePositive: false,
      };
      
      expect(vuln.cvss.score).toBeGreaterThanOrEqual(7.0); // High severity
      expect(vuln.remediation.priority).toBeGreaterThanOrEqual(7);
      expect(vuln.remediation.priority).toBeLessThanOrEqual(10);
    });

    it('should map to OWASP Top 10', () => {
      const owaspCategories = Object.keys(OWASP_TOP_10_2021);
      
      expect(owaspCategories).toContain('A01');
      expect(owaspCategories).toContain('A03');
      expect(OWASP_TOP_10_2021.A03).toBe('Injection');
      expect(OWASP_TOP_10_2021.A05).toBe('Security Misconfiguration');
    });
  });

  describe('Cookie Security', () => {
    it('should identify insecure cookie configurations', () => {
      const insecureCookie: CookieAnalysis = {
        name: 'session',
        value: 'abc123',
        domain: '.example.com',
        path: '/',
        expires: null,
        size: 32,
        httpOnly: false,
        secure: false,
        sameSite: null,
        issues: [],
        recommendations: [],
      };
      
      // Validation logic
      if (!insecureCookie.secure) {
        insecureCookie.issues.push('Missing Secure flag');
      }
      if (!insecureCookie.httpOnly) {
        insecureCookie.issues.push('Missing HttpOnly flag');
      }
      if (!insecureCookie.sameSite) {
        insecureCookie.issues.push('Missing SameSite attribute');
      }
      
      expect(insecureCookie.issues).toHaveLength(3);
      expect(insecureCookie.issues).toContain('Missing Secure flag');
    });
  });

  describe('Compliance Checks', () => {
    it('should validate GDPR compliance requirements', () => {
      const gdprRequirements = [
        'Cookie consent',
        'Privacy policy',
        'Data retention policy',
        'Right to erasure',
        'Data portability',
      ];
      
      const implemented = ['Privacy policy'];
      const complianceScore = (implemented.length / gdprRequirements.length) * 100;
      
      expect(complianceScore).toBeLessThan(50); // Not compliant
    });

    it('should calculate ISO 27001 compliance percentage', () => {
      const totalControls = 114;
      const implementedControls = 78;
      const percentage = Math.round((implementedControls / totalControls) * 100);
      
      expect(percentage).toBe(68);
    });
  });

  describe('Security Best Practices', () => {
    it('should categorize practice importance', () => {
      const practices = [
        { practice: 'Use HTTPS', importance: 'critical' },
        { practice: 'Security headers', importance: 'high' },
        { practice: 'Rate limiting', importance: 'medium' },
        { practice: 'Security.txt file', importance: 'low' },
      ];
      
      const criticalPractices = practices.filter(p => p.importance === 'critical');
      expect(criticalPractices).toHaveLength(1);
      expect(criticalPractices[0].practice).toBe('Use HTTPS');
    });
  });
});

describe('Security Analysis Helper Functions', () => {
  it('should calculate overall security score', () => {
    const calculateSecurityScore = (components: Record<string, number>): number => {
      const weights = {
        headers: 0.25,
        ssl: 0.25,
        vulnerabilities: 0.30,
        bestPractices: 0.20,
      };
      
      let totalScore = 0;
      for (const [component, score] of Object.entries(components)) {
        totalScore += score * (weights[component as keyof typeof weights] || 0);
      }
      
      return Math.round(totalScore);
    };
    
    const score = calculateSecurityScore({
      headers: 65,
      ssl: 85,
      vulnerabilities: 70,
      bestPractices: 70,
    });
    
    expect(score).toBe(72);
  });

  it('should determine security grade from score', () => {
    const getSecurityGrade = (score: number): SecurityGrade => {
      for (const [grade, threshold] of Object.entries(SECURITY_SCORE_THRESHOLDS)) {
        if (score >= threshold.min && score <= threshold.max) {
          return grade as SecurityGrade;
        }
      }
      return SecurityGrade.F;
    };
    
    expect(getSecurityGrade(96)).toBe(SecurityGrade.A_PLUS);
    expect(getSecurityGrade(87)).toBe(SecurityGrade.A);
    expect(getSecurityGrade(72)).toBe(SecurityGrade.B);
    expect(getSecurityGrade(25)).toBe(SecurityGrade.F);
  });

  it('should prioritize security issues', () => {
    const prioritizeIssues = (issues: SecurityVulnerability[]): SecurityVulnerability[] => {
      const severityOrder = {
        [SecuritySeverity.CRITICAL]: 5,
        [SecuritySeverity.HIGH]: 4,
        [SecuritySeverity.MEDIUM]: 3,
        [SecuritySeverity.LOW]: 2,
        [SecuritySeverity.INFO]: 1,
      };
      
      return issues.sort((a, b) => {
        const severityDiff = severityOrder[b.severity] - severityOrder[a.severity];
        if (severityDiff !== 0) return severityDiff;
        return b.remediation.priority - a.remediation.priority;
      });
    };
    
    const issues: SecurityVulnerability[] = [
      { severity: SecuritySeverity.LOW, remediation: { priority: 5 } } as SecurityVulnerability,
      { severity: SecuritySeverity.CRITICAL, remediation: { priority: 10 } } as SecurityVulnerability,
      { severity: SecuritySeverity.HIGH, remediation: { priority: 8 } } as SecurityVulnerability,
    ];
    
    const prioritized = prioritizeIssues(issues);
    expect(prioritized[0].severity).toBe(SecuritySeverity.CRITICAL);
    expect(prioritized[1].severity).toBe(SecuritySeverity.HIGH);
  });
});