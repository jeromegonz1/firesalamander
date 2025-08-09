/**
 * Fire Salamander - Security Analysis Mapper
 * Maps backend security data to SecurityAnalysis interface
 * Lead Tech quality with comprehensive validation and error handling
 */

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
  SSLConfiguration,
  CookieAnalysis,
  ContentSecurityPolicy,
  SecurityVulnerability,
  SecurityBestPractice,
  InfrastructureSecurity,
  PermissionAnalysis,
  SECURITY_HEADERS,
  SECURITY_SCORE_THRESHOLDS,
  OWASP_TOP_10_2021,
} from '@/types/security-analysis';

/**
 * Main mapper function - converts backend data to SecurityAnalysis
 */
export function mapBackendToSecurityAnalysis(backendData: any): SecurityAnalysis {
  try {
    console.log('Mapping backend security data:', backendData);
    
    if (!backendData || typeof backendData !== 'object') {
      console.warn('Invalid backend data provided, using empty analysis');
      return createEmptySecurityAnalysis();
    }

    const securityData = backendData.security_data || {};
    
    // Map all security components
    const headers = mapSecurityHeaders(securityData.headers || {});
    const ssl = mapSSLAnalysis(securityData.ssl || {});
    const cookies = mapCookieAnalysis(securityData.cookies || []);
    const csp = mapContentSecurityPolicy(securityData.csp || {});
    const vulnerabilities = mapVulnerabilities(securityData.vulnerabilities || []);
    const bestPractices = mapBestPractices(securityData);
    const infrastructure = mapInfrastructure(securityData);
    const permissions = mapPermissions(securityData.permissions || {});
    const compliance = mapCompliance(securityData, headers, ssl, cookies);

    // Calculate overall score
    const overallScore = calculateOverallScore({
      headers: headers.score,
      ssl: ssl.score,
      vulnerabilities: calculateVulnerabilityScore(vulnerabilities),
      bestPractices: bestPractices.score,
    });

    // Generate recommendations
    const recommendations = generateRecommendations(
      headers,
      ssl,
      vulnerabilities,
      bestPractices,
      permissions
    );

    // Build the complete analysis
    const analysis: SecurityAnalysis = {
      score: {
        overall: overallScore,
        grade: getSecurityGrade(overallScore),
        trend: determineTrend(overallScore, backendData.previous_score),
        previousScore: backendData.previous_score,
      },
      headers,
      ssl,
      cookies,
      csp,
      vulnerabilities,
      bestPractices,
      infrastructure,
      permissions,
      compliance,
      recommendations,
      metadata: {
        scanDate: backendData.analyzed_at || new Date().toISOString(),
        scanDuration: backendData.scan_duration || 0,
        toolVersion: '1.0.0',
        scanDepth: determineScanDepth(backendData),
        pagesScanned: backendData.pages_scanned || 0,
        errorsEncountered: backendData.errors || 0,
      },
    };

    console.log('Successfully mapped security analysis:', analysis);
    return analysis;

  } catch (error) {
    console.error('Error mapping backend security data:', error);
    return createEmptySecurityAnalysis();
  }
}

/**
 * Maps security headers analysis
 */
function mapSecurityHeaders(headersData: any): SecurityAnalysis['headers'] {
  const securityHeaders: SecurityHeader[] = [];
  const missingHeaders: string[] = [];
  const misconfiguredHeaders: string[] = [];

  // Check all required security headers
  Object.entries(SECURITY_HEADERS).forEach(([headerName, config]) => {
    const headerDetail = headersData.headers_detail?.find(
      (h: any) => h.name === headerName
    );
    const isPresent = headersData.present?.includes(headerName);
    const isMissing = headersData.missing?.includes(headerName);

    if (headerDetail || isPresent) {
      const header: SecurityHeader = {
        name: headerName,
        value: headerDetail?.value || null,
        status: determineHeaderStatus(headerDetail, isPresent),
        recommendation: getHeaderRecommendation(headerName, headerDetail),
        impact: getHeaderImpact(headerName),
        severity: getHeaderSeverity(headerName, headerDetail),
        references: getHeaderReferences(headerName),
      };
      
      securityHeaders.push(header);

      if (header.status === HeaderStatus.MISCONFIGURED) {
        misconfiguredHeaders.push(headerName);
      }
    } else if (isMissing || config.required) {
      // Add missing header
      securityHeaders.push({
        name: headerName,
        value: null,
        status: HeaderStatus.MISSING,
        recommendation: `Implement ${headerName} header with value: ${config.defaultValue}`,
        impact: getHeaderImpact(headerName),
        severity: config.required ? SecuritySeverity.HIGH : SecuritySeverity.MEDIUM,
        references: getHeaderReferences(headerName),
      });
      
      missingHeaders.push(headerName);
    }
  });

  const score = calculateHeaderScore({
    present: headersData.present || [],
    missing: missingHeaders,
    misconfigured: misconfiguredHeaders,
  });

  return {
    securityHeaders,
    score,
    grade: getSecurityGrade(score),
    missing: missingHeaders,
    misconfigured: misconfiguredHeaders,
  };
}

/**
 * Maps SSL/TLS analysis
 */
function mapSSLAnalysis(sslData: any): SecurityAnalysis['ssl'] {
  if (!sslData.enabled) {
    return {
      enabled: false,
      configuration: null,
      score: 0,
      grade: SecurityGrade.F,
      issues: ['SSL/TLS not enabled'],
      strengths: [],
    };
  }

  const configuration = mapSSLConfiguration(sslData);
  const { score, issues, strengths } = evaluateSSLConfiguration(configuration);

  return {
    enabled: true,
    configuration,
    score,
    grade: getSecurityGrade(score),
    issues,
    strengths,
  };
}

/**
 * Maps SSL configuration details
 */
function mapSSLConfiguration(sslData: any): SSLConfiguration {
  const protocols = mapSSLProtocols(sslData.protocol);
  const cipherSuites = mapCipherSuites(sslData.cipher_suites || []);
  const certificate = mapSSLCertificate(sslData.certificate || {});

  return {
    protocols,
    cipherSuites,
    certificate,
    hsts: {
      enabled: sslData.hsts?.enabled || false,
      maxAge: sslData.hsts?.max_age || 0,
      includeSubdomains: sslData.hsts?.include_subdomains || false,
      preload: sslData.hsts?.preload || false,
    },
    forwardSecrecy: sslData.forward_secrecy || false,
    ocspStapling: sslData.ocsp_stapling || false,
    sessionResumption: sslData.session_resumption || false,
  };
}

/**
 * Maps SSL protocols
 */
function mapSSLProtocols(protocol: string): SSLConfiguration['protocols'] {
  const protocols = [
    { protocol: SSLProtocol.TLS_1_3, enabled: false, secure: true },
    { protocol: SSLProtocol.TLS_1_2, enabled: false, secure: true },
    { protocol: SSLProtocol.TLS_1_1, enabled: false, secure: false },
    { protocol: SSLProtocol.TLS_1_0, enabled: false, secure: false },
    { protocol: SSLProtocol.SSL_3_0, enabled: false, secure: false },
    { protocol: SSLProtocol.SSL_2_0, enabled: false, secure: false },
  ];

  // Enable the detected protocol
  if (protocol) {
    const protocolEntry = protocols.find(p => 
      p.protocol.toLowerCase().replace(/\s/g, '') === protocol.toLowerCase().replace(/\s/g, '')
    );
    if (protocolEntry) {
      protocolEntry.enabled = true;
    }
  }

  return protocols;
}

/**
 * Maps cipher suites
 */
function mapCipherSuites(cipherSuitesData: string[]): SSLConfiguration['cipherSuites'] {
  return cipherSuitesData.map(suite => ({
    name: suite,
    strength: determineCipherStrength(suite),
    protocol: extractProtocolFromCipher(suite),
    keyExchange: extractKeyExchange(suite),
    authentication: extractAuthentication(suite),
    encryption: extractEncryption(suite),
    mac: extractMAC(suite),
    export: suite.includes('EXPORT'),
    recommended: isRecommendedCipher(suite),
  }));
}

/**
 * Maps SSL certificate
 */
function mapSSLCertificate(certData: any): SSLCertificate {
  const daysRemaining = certData.days_remaining || 0;
  const isExpiringSoon = daysRemaining < 30 && daysRemaining > 0;
  const isExpired = daysRemaining <= 0;

  return {
    issuer: certData.issuer || 'Unknown',
    subject: certData.subject || 'Unknown',
    validFrom: certData.valid_from || '',
    validTo: certData.valid_to || '',
    daysRemaining,
    isValid: certData.chain_valid && !isExpired,
    isExpired,
    isExpiringSoon,
    algorithm: certData.algorithm || 'RSA',
    keySize: certData.key_size || 2048,
    san: certData.san || [],
    chainValid: certData.chain_valid || false,
    ocspStatus: certData.ocsp_status || 'unknown',
  };
}

/**
 * Maps cookie analysis
 */
function mapCookieAnalysis(cookiesData: any[]): SecurityAnalysis['cookies'] {
  const cookieList: CookieAnalysis[] = cookiesData.map(cookie => {
    const issues = analyzeCookieSecurity(cookie);
    const recommendations = generateCookieRecommendations(cookie, issues);

    return {
      name: cookie.name,
      value: cookie.value,
      domain: cookie.domain,
      path: cookie.path,
      expires: cookie.expires || null,
      size: cookie.value?.length || 0,
      httpOnly: cookie.http_only || false,
      secure: cookie.secure || false,
      sameSite: cookie.same_site || null,
      issues,
      recommendations,
    };
  });

  const totalCookies = cookieList.length;
  const secureCookies = cookieList.filter(c => c.secure).length;
  const httpOnlyCookies = cookieList.filter(c => c.httpOnly).length;
  const sameSiteCookies = cookieList.filter(c => c.sameSite).length;

  const issues: string[] = [];
  if (secureCookies < totalCookies) {
    issues.push('Some cookies lack Secure flag');
  }
  if (httpOnlyCookies < totalCookies) {
    issues.push('Some cookies lack HttpOnly flag');
  }
  if (sameSiteCookies < totalCookies) {
    issues.push('Some cookies lack SameSite attribute');
  }

  return {
    list: cookieList,
    totalCookies,
    secureCookies,
    httpOnlyCookies,
    sameSiteCookies,
    issues,
  };
}

/**
 * Maps Content Security Policy
 */
function mapContentSecurityPolicy(cspData: any): ContentSecurityPolicy {
  if (!cspData || !cspData.present) {
    return {
      present: false,
      directives: {},
      issues: [{
        directive: 'default-src',
        issue: 'No Content Security Policy detected',
        severity: SecuritySeverity.HIGH,
        recommendation: "Implement CSP starting with: default-src 'self'",
      }],
      reportUri: null,
      upgradeInsecureRequests: false,
      blockAllMixedContent: false,
    };
  }

  return {
    present: true,
    directives: cspData.directives || {},
    issues: mapCSPIssues(cspData.directives || {}),
    reportUri: cspData.report_uri || null,
    upgradeInsecureRequests: cspData.upgrade_insecure_requests || false,
    blockAllMixedContent: cspData.block_all_mixed_content || false,
  };
}

/**
 * Maps vulnerabilities
 */
function mapVulnerabilities(vulnsData: any[]): SecurityAnalysis['vulnerabilities'] {
  const categorized = {
    critical: [] as SecurityVulnerability[],
    high: [] as SecurityVulnerability[],
    medium: [] as SecurityVulnerability[],
    low: [] as SecurityVulnerability[],
    info: [] as SecurityVulnerability[],
  };

  vulnsData.forEach((vuln, index) => {
    const severity = mapSeverity(vuln.severity);
    const mappedVuln: SecurityVulnerability = {
      id: vuln.id || `VULN-${index + 1}`,
      type: vuln.type || 'Unknown',
      severity,
      category: mapVulnerabilityCategory(vuln.type),
      title: vuln.title || vuln.type,
      description: vuln.description || '',
      affectedUrls: vuln.affected_urls || [],
      cvss: {
        score: vuln.cvss_score || 0,
        vector: vuln.cvss_vector || '',
      },
      cwe: vuln.cwe || '',
      owasp: mapToOWASP(vuln.type, vuln.cwe),
      remediation: {
        summary: vuln.recommendations?.[0] || 'Apply security best practices',
        steps: vuln.recommendations || [],
        effort: determineRemediationEffort(vuln),
        priority: calculateVulnerabilityPriority(severity, vuln),
      },
      references: vuln.references?.map((ref: any) => ({
        title: ref.title || ref,
        url: ref.url || ref,
      })) || [],
      exploitability: determineExploitability(vuln),
      dateDiscovered: vuln.date_discovered || new Date().toISOString(),
      falsePositive: vuln.false_positive || false,
    };

    categorized[severity].push(mappedVuln);
  });

  const total = vulnsData.length;

  return { ...categorized, total };
}

/**
 * Maps best practices
 */
function mapBestPractices(securityData: any): SecurityAnalysis['bestPractices'] {
  const implemented: SecurityBestPractice[] = [];
  const missing: SecurityBestPractice[] = [];

  // Check HTTPS
  const httpsEnabled = securityData.ssl?.enabled || false;
  const httpsPractice: SecurityBestPractice = {
    id: 'BP-001',
    category: SecurityCategory.SSL_TLS,
    practice: 'Use HTTPS everywhere',
    implemented: httpsEnabled,
    importance: 'critical',
    description: httpsEnabled ? 'All pages are served over HTTPS' : 'Site not using HTTPS',
    implementation: 'SSL certificate installed and HTTP redirects to HTTPS',
    benefits: [
      'Encrypts data in transit',
      'Prevents man-in-the-middle attacks',
      'Required for modern web features',
    ],
    references: ['https://developers.google.com/web/fundamentals/security/encrypt-in-transit/why-https'],
  };

  if (httpsEnabled) {
    implemented.push(httpsPractice);
  } else {
    missing.push(httpsPractice);
  }

  // Check security headers
  const hasSecurityHeaders = (securityData.headers?.present?.length || 0) >= 3;
  const headersPractice: SecurityBestPractice = {
    id: 'BP-002',
    category: SecurityCategory.HEADERS,
    practice: 'Implement security headers',
    implemented: hasSecurityHeaders,
    importance: 'high',
    description: 'Security headers provide defense in depth',
    implementation: 'Configure web server to send security headers',
    benefits: [
      'Prevents clickjacking',
      'Mitigates XSS attacks',
      'Controls resource loading',
    ],
    references: ['https://securityheaders.com/'],
  };

  if (hasSecurityHeaders) {
    implemented.push(headersPractice);
  } else {
    missing.push(headersPractice);
  }

  const score = (implemented.length / (implemented.length + missing.length)) * 100;

  return { implemented, missing, score };
}

/**
 * Maps infrastructure security
 */
function mapInfrastructure(securityData: any): InfrastructureSecurity {
  const technologies = (securityData.technologies || []).map((tech: any) => ({
    name: tech.name,
    version: tech.version,
    category: tech.category,
    outdated: isOutdatedVersion(tech),
    vulnerabilities: tech.vulnerabilities || 0,
    latestVersion: tech.latest_version || 'Unknown',
  }));

  return {
    webServer: {
      name: securityData.web_server?.name || 'Unknown',
      version: securityData.web_server?.version || 'Unknown',
      hideVersion: securityData.web_server?.hide_version || false,
      knownVulnerabilities: securityData.web_server?.vulnerabilities || 0,
    },
    technologies,
    openPorts: (securityData.open_ports || []).map((port: any) => ({
      port: port.port || port,
      service: port.service || 'Unknown',
      secure: isSecurePort(port.port || port),
      recommendation: getPortRecommendation(port.port || port),
    })),
    dnssec: securityData.dnssec || false,
    ipv6: securityData.ipv6 || false,
    http2: securityData.http2 || false,
    http3: securityData.http3 || false,
  };
}

/**
 * Maps permissions analysis
 */
function mapPermissions(permissionsData: any): PermissionAnalysis {
  const robotsTxt = permissionsData.robots_txt || {};
  const sensitive = detectSensitivePaths(robotsTxt.disallowed_paths || []);

  return {
    robotsTxt: {
      present: robotsTxt.present || false,
      sensitive,
      issues: sensitive.length > 0 ? ['Exposes sensitive paths in robots.txt'] : [],
    },
    sitemapXml: {
      present: permissionsData.sitemap_xml?.present || false,
      accessible: permissionsData.sitemap_xml?.accessible || false,
      sensitive: detectSensitivePaths(permissionsData.sitemap_xml?.urls || []),
    },
    adminPaths: {
      found: permissionsData.admin_paths || [],
      accessible: permissionsData.admin_paths?.filter((p: string) => 
        !permissionsData.protected_paths?.includes(p)
      ) || [],
      protected: permissionsData.protected_paths || [],
    },
    apiEndpoints: {
      discovered: permissionsData.api_endpoints || [],
      authenticated: permissionsData.authenticated_endpoints || [],
      rateLimit: permissionsData.rate_limiting || false,
    },
    backupFiles: permissionsData.exposed_files?.filter((f: string) => 
      f.includes('backup') || f.endsWith('.sql') || f.endsWith('.zip')
    ) || [],
    configFiles: permissionsData.exposed_files?.filter((f: string) => 
      f.includes('config') || f.endsWith('.env') || f.endsWith('.ini')
    ) || [],
    sourceCodeExposure: permissionsData.exposed_files?.filter((f: string) => 
      f.includes('.git') || f.endsWith('.py') || f.endsWith('.php')
    ) || [],
  };
}

/**
 * Maps compliance status
 */
function mapCompliance(
  securityData: any,
  headers: SecurityAnalysis['headers'],
  ssl: SecurityAnalysis['ssl'],
  cookies: SecurityAnalysis['cookies']
): SecurityAnalysis['compliance'] {
  const gdprIssues: string[] = [];
  if (cookies.totalCookies > 0 && !securityData.cookie_consent) {
    gdprIssues.push('No cookie consent banner');
  }
  if (!securityData.privacy_policy_comprehensive) {
    gdprIssues.push('Privacy policy not comprehensive');
  }

  const gdprScore = calculateGDPRScore(gdprIssues, cookies);

  const pciIssues: string[] = [];
  if (!ssl.enabled) {
    pciIssues.push('Not using HTTPS');
  }
  if (headers.missing.length > 3) {
    pciIssues.push('Missing security headers');
  }

  return {
    gdpr: {
      compliant: gdprIssues.length === 0,
      issues: gdprIssues,
      score: gdprScore,
    },
    pci: {
      level: 4,
      compliant: pciIssues.length === 0,
      issues: pciIssues,
    },
    iso27001: {
      controls: 114,
      implemented: Math.round((headers.score + ssl.score) / 2 * 1.14), // Rough estimate
      percentage: Math.round((headers.score + ssl.score) / 2),
    },
  };
}

// ==================== HELPER FUNCTIONS ====================

/**
 * Creates an empty SecurityAnalysis structure
 */
export function createEmptySecurityAnalysis(): SecurityAnalysis {
  return {
    score: {
      overall: 0,
      grade: SecurityGrade.F,
      trend: 'stable',
    },
    headers: {
      securityHeaders: [],
      score: 0,
      grade: SecurityGrade.F,
      missing: [],
      misconfigured: [],
    },
    ssl: {
      enabled: false,
      configuration: null,
      score: 0,
      grade: SecurityGrade.F,
      issues: ['SSL/TLS not enabled'],
      strengths: [],
    },
    cookies: {
      list: [],
      totalCookies: 0,
      secureCookies: 0,
      httpOnlyCookies: 0,
      sameSiteCookies: 0,
      issues: [],
    },
    csp: {
      present: false,
      directives: {},
      issues: [],
      reportUri: null,
      upgradeInsecureRequests: false,
      blockAllMixedContent: false,
    },
    vulnerabilities: {
      critical: [],
      high: [],
      medium: [],
      low: [],
      info: [],
      total: 0,
    },
    bestPractices: {
      implemented: [],
      missing: [],
      score: 0,
    },
    infrastructure: {
      webServer: {
        name: 'Unknown',
        version: 'Unknown',
        hideVersion: false,
        knownVulnerabilities: 0,
      },
      technologies: [],
      openPorts: [],
      dnssec: false,
      ipv6: false,
      http2: false,
      http3: false,
    },
    permissions: {
      robotsTxt: { present: false, sensitive: [], issues: [] },
      sitemapXml: { present: false, accessible: false, sensitive: [] },
      adminPaths: { found: [], accessible: [], protected: [] },
      apiEndpoints: { discovered: [], authenticated: [], rateLimit: false },
      backupFiles: [],
      configFiles: [],
      sourceCodeExposure: [],
    },
    compliance: {
      gdpr: { compliant: false, issues: [], score: 0 },
      pci: { level: 4, compliant: false, issues: [] },
      iso27001: { controls: 114, implemented: 0, percentage: 0 },
    },
    recommendations: {
      immediate: [],
      shortTerm: [],
      longTerm: [],
    },
    metadata: {
      scanDate: new Date().toISOString(),
      scanDuration: 0,
      toolVersion: '1.0.0',
      scanDepth: 'basic',
      pagesScanned: 0,
      errorsEncountered: 0,
    },
  };
}

/**
 * Determines security grade from score
 */
export function getSecurityGrade(score: number): SecurityGrade {
  for (const [grade, threshold] of Object.entries(SECURITY_SCORE_THRESHOLDS)) {
    if (score >= threshold.min && score <= threshold.max) {
      return grade as SecurityGrade;
    }
  }
  return SecurityGrade.F;
}

/**
 * Maps header status
 */
export function mapHeaderStatus(status: string): HeaderStatus {
  const statusMap: Record<string, HeaderStatus> = {
    'present': HeaderStatus.PRESENT,
    'missing': HeaderStatus.MISSING,
    'misconfigured': HeaderStatus.MISCONFIGURED,
    'deprecated': HeaderStatus.DEPRECATED,
  };
  return statusMap[status.toLowerCase()] || HeaderStatus.MISSING;
}

/**
 * Maps security severity
 */
export function mapSeverity(severity: string): SecuritySeverity {
  const severityMap: Record<string, SecuritySeverity> = {
    'critical': SecuritySeverity.CRITICAL,
    'high': SecuritySeverity.HIGH,
    'medium': SecuritySeverity.MEDIUM,
    'low': SecuritySeverity.LOW,
    'info': SecuritySeverity.INFO,
    'informational': SecuritySeverity.INFO,
  };
  return severityMap[severity.toLowerCase()] || SecuritySeverity.INFO;
}

/**
 * Analyzes cookie security issues
 */
export function analyzeCookieSecurity(cookie: any): string[] {
  const issues: string[] = [];
  
  if (!cookie.secure) {
    issues.push('Missing Secure flag');
  }
  if (!cookie.http_only && cookie.name.toLowerCase().includes('session')) {
    issues.push('Missing HttpOnly flag');
  }
  if (!cookie.same_site) {
    issues.push('Missing SameSite attribute');
  }
  
  return issues;
}

/**
 * Detects sensitive paths
 */
export function detectSensitivePaths(paths: string[]): string[] {
  const sensitivePatterns = [
    '/admin',
    '/backup',
    '/.git',
    '/config',
    '/api/private',
    '/wp-admin',
    '/phpmyadmin',
    '/.env',
    '/database',
  ];
  
  return paths.filter(path => 
    sensitivePatterns.some(pattern => 
      path.toLowerCase().includes(pattern.toLowerCase())
    )
  );
}

/**
 * Calculates compliance scores
 */
export function calculateComplianceScores(securityData: any): any {
  const gdprIssues: string[] = [];
  
  if (!securityData.cookies?.secureCookies) {
    gdprIssues.push('Insecure cookies');
  }
  if (!securityData.headers?.score || securityData.headers.score < 60) {
    gdprIssues.push('Insufficient security headers');
  }
  
  const gdprScore = Math.max(0, 100 - (gdprIssues.length * 20));
  
  return {
    gdpr: {
      compliant: gdprIssues.length === 0,
      issues: gdprIssues,
      score: gdprScore,
    },
  };
}

/**
 * Calculates header score
 */
export function calculateHeaderScore(headers: any): number {
  const totalHeaders = Object.keys(SECURITY_HEADERS).filter(h => 
    SECURITY_HEADERS[h as keyof typeof SECURITY_HEADERS].required
  ).length;
  
  const presentHeaders = headers.present?.length || 0;
  const misconfiguredHeaders = headers.misconfigured?.length || 0;
  
  const score = Math.round(
    ((presentHeaders - misconfiguredHeaders * 0.5) / totalHeaders) * 100
  );
  
  return Math.max(0, Math.min(100, score));
}

/**
 * Calculates SSL score
 */
export function calculateSSLScore(sslConfig: any): number {
  let score = 100;
  
  // Protocol penalties
  if (!sslConfig.protocols?.includes('TLSv1.3')) score -= 10;
  if (sslConfig.protocols?.includes('TLSv1.1')) score -= 20;
  if (sslConfig.protocols?.includes('TLSv1.0')) score -= 30;
  
  // Certificate penalties
  if (sslConfig.certificate?.daysRemaining < 30) score -= 15;
  if (!sslConfig.certificate?.chainValid) score -= 20;
  
  // Missing features
  if (!sslConfig.forwardSecrecy) score -= 10;
  if (!sslConfig.ocspStapling) score -= 5;
  
  return Math.max(0, score);
}

/**
 * Prioritizes security recommendations
 */
export function prioritizeRecommendations(issues: any[]): any[] {
  const severityOrder = {
    [SecuritySeverity.CRITICAL]: 5,
    [SecuritySeverity.HIGH]: 4,
    [SecuritySeverity.MEDIUM]: 3,
    [SecuritySeverity.LOW]: 2,
    [SecuritySeverity.INFO]: 1,
  };
  
  const effortOrder = {
    'low': 1,
    'medium': 2,
    'high': 3,
  };
  
  return issues.sort((a, b) => {
    // First by severity
    const severityDiff = severityOrder[b.severity] - severityOrder[a.severity];
    if (severityDiff !== 0) return severityDiff;
    
    // Then by effort (lower effort first)
    return effortOrder[a.effort] - effortOrder[b.effort];
  });
}

// ==================== PRIVATE HELPER FUNCTIONS ====================

function calculateOverallScore(scores: Record<string, number>): number {
  const weights = {
    headers: 0.25,
    ssl: 0.25,
    vulnerabilities: 0.30,
    bestPractices: 0.20,
  };
  
  let totalScore = 0;
  for (const [component, score] of Object.entries(scores)) {
    totalScore += score * (weights[component as keyof typeof weights] || 0);
  }
  
  return Math.round(totalScore);
}

function calculateVulnerabilityScore(vulnerabilities: SecurityAnalysis['vulnerabilities']): number {
  let score = 100;
  
  // Deduct points based on vulnerability severity
  score -= vulnerabilities.critical.length * 20;
  score -= vulnerabilities.high.length * 10;
  score -= vulnerabilities.medium.length * 5;
  score -= vulnerabilities.low.length * 2;
  
  return Math.max(0, score);
}

function determineHeaderStatus(headerDetail: any, isPresent: boolean): HeaderStatus {
  if (!headerDetail && !isPresent) return HeaderStatus.MISSING;
  if (headerDetail?.recommendations?.length > 0) return HeaderStatus.MISCONFIGURED;
  return HeaderStatus.PRESENT;
}

function getHeaderRecommendation(headerName: string, headerDetail: any): string {
  if (headerDetail?.recommendations?.length > 0) {
    return headerDetail.recommendations.join('. ');
  }
  
  const config = SECURITY_HEADERS[headerName as keyof typeof SECURITY_HEADERS];
  if (!config) return `Configure ${headerName} header`;
  
  return `Set ${headerName} to: ${config.defaultValue}`;
}

function getHeaderImpact(headerName: string): string {
  const impacts: Record<string, string> = {
    'Strict-Transport-Security': 'Enforces HTTPS connections, preventing protocol downgrade attacks',
    'Content-Security-Policy': 'Prevents XSS and data injection attacks',
    'X-Frame-Options': 'Prevents clickjacking attacks',
    'X-Content-Type-Options': 'Prevents MIME type sniffing',
    'Referrer-Policy': 'Controls information sent in Referer header',
    'Permissions-Policy': 'Controls browser features and APIs',
  };
  
  return impacts[headerName] || `Improves security by implementing ${headerName}`;
}

function getHeaderSeverity(headerName: string, headerDetail: any): SecuritySeverity {
  const config = SECURITY_HEADERS[headerName as keyof typeof SECURITY_HEADERS];
  
  if (!config?.required) return SecuritySeverity.MEDIUM;
  if (headerDetail?.status === HeaderStatus.MISSING) return SecuritySeverity.HIGH;
  if (headerDetail?.status === HeaderStatus.MISCONFIGURED) return SecuritySeverity.MEDIUM;
  
  return SecuritySeverity.INFO;
}

function getHeaderReferences(headerName: string): string[] {
  const baseUrl = 'https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/';
  return [baseUrl + headerName];
}

function determineCipherStrength(cipher: string): CipherStrength {
  if (cipher.includes('AES_256_GCM') || cipher.includes('CHACHA20')) {
    return CipherStrength.STRONG;
  }
  if (cipher.includes('AES_128_GCM')) {
    return CipherStrength.MODERATE;
  }
  if (cipher.includes('3DES') || cipher.includes('RC4')) {
    return CipherStrength.WEAK;
  }
  if (cipher.includes('NULL') || cipher.includes('EXPORT')) {
    return CipherStrength.INSECURE;
  }
  return CipherStrength.MODERATE;
}

function isRecommendedCipher(cipher: string): boolean {
  const recommended = [
    'TLS_AES_256_GCM_SHA384',
    'TLS_CHACHA20_POLY1305_SHA256',
    'TLS_AES_128_GCM_SHA256',
    'ECDHE-RSA-AES256-GCM-SHA384',
    'ECDHE-RSA-AES128-GCM-SHA256',
  ];
  return recommended.some(rec => cipher.includes(rec));
}

function extractProtocolFromCipher(cipher: string): string {
  if (cipher.startsWith('TLS_')) return 'TLS 1.3';
  if (cipher.includes('ECDHE')) return 'TLS 1.2';
  return 'Unknown';
}

function extractKeyExchange(cipher: string): string {
  if (cipher.includes('ECDHE')) return 'ECDHE';
  if (cipher.includes('DHE')) return 'DHE';
  if (cipher.includes('RSA')) return 'RSA';
  return 'Unknown';
}

function extractAuthentication(cipher: string): string {
  if (cipher.includes('RSA')) return 'RSA';
  if (cipher.includes('ECDSA')) return 'ECDSA';
  if (cipher.includes('PSK')) return 'PSK';
  return 'Unknown';
}

function extractEncryption(cipher: string): string {
  if (cipher.includes('AES256-GCM')) return 'AES256-GCM';
  if (cipher.includes('AES128-GCM')) return 'AES128-GCM';
  if (cipher.includes('CHACHA20')) return 'CHACHA20-POLY1305';
  return 'Unknown';
}

function extractMAC(cipher: string): string {
  if (cipher.includes('SHA384')) return 'SHA384';
  if (cipher.includes('SHA256')) return 'SHA256';
  if (cipher.includes('SHA')) return 'SHA';
  return 'AEAD';
}

function evaluateSSLConfiguration(config: SSLConfiguration): {
  score: number;
  issues: string[];
  strengths: string[];
} {
  let score = 100;
  const issues: string[] = [];
  const strengths: string[] = [];

  // Check protocols
  const hasInsecureProtocols = config.protocols.some(p => !p.secure && p.enabled);
  if (hasInsecureProtocols) {
    score -= 30;
    issues.push('Insecure protocols enabled');
  }

  const hasTLS13 = config.protocols.find(p => p.protocol === SSLProtocol.TLS_1_3)?.enabled;
  if (hasTLS13) {
    strengths.push('TLS 1.3 enabled');
  } else {
    score -= 10;
    issues.push('TLS 1.3 not enabled');
  }

  // Check certificate
  if (config.certificate.isExpiringSoon) {
    score -= 15;
    issues.push('Certificate expiring soon');
  }
  if (!config.certificate.chainValid) {
    score -= 20;
    issues.push('Invalid certificate chain');
  }

  // Check HSTS
  if (!config.hsts.enabled) {
    score -= 10;
    issues.push('HSTS not enabled');
  } else if (!config.hsts.preload) {
    score -= 5;
    issues.push('HSTS preload not enabled');
  }

  // Check features
  if (config.forwardSecrecy) strengths.push('Perfect forward secrecy');
  else { score -= 10; issues.push('No forward secrecy'); }

  if (config.ocspStapling) strengths.push('OCSP stapling');
  else { score -= 5; issues.push('No OCSP stapling'); }

  return {
    score: Math.max(0, score),
    issues,
    strengths,
  };
}

function generateCookieRecommendations(cookie: any, issues: string[]): string[] {
  const recommendations: string[] = [];
  
  if (issues.includes('Missing Secure flag')) {
    recommendations.push('Add Secure flag to ensure cookie is only sent over HTTPS');
  }
  if (issues.includes('Missing HttpOnly flag')) {
    recommendations.push('Add HttpOnly flag to prevent JavaScript access');
  }
  if (issues.includes('Missing SameSite attribute')) {
    recommendations.push('Set SameSite attribute to prevent CSRF attacks');
  }
  
  return recommendations;
}

function mapCSPIssues(directives: Record<string, any>): ContentSecurityPolicy['issues'] {
  const issues: ContentSecurityPolicy['issues'] = [];
  
  if (!directives['default-src']) {
    issues.push({
      directive: 'default-src',
      issue: 'Missing default-src directive',
      severity: SecuritySeverity.HIGH,
      recommendation: "Add default-src 'self' as fallback policy",
    });
  }
  
  if (directives['script-src']?.includes("'unsafe-inline'")) {
    issues.push({
      directive: 'script-src',
      issue: "Using 'unsafe-inline' for scripts",
      severity: SecuritySeverity.HIGH,
      recommendation: "Remove 'unsafe-inline' and use nonces or hashes",
    });
  }
  
  return issues;
}

function mapVulnerabilityCategory(type: string): SecurityCategory {
  const typeMap: Record<string, SecurityCategory> = {
    'Missing Security Header': SecurityCategory.HEADERS,
    'SSL': SecurityCategory.SSL_TLS,
    'Cookie': SecurityCategory.COOKIES,
    'XSS': SecurityCategory.CONTENT,
    'SQL Injection': SecurityCategory.DATA_PROTECTION,
  };
  
  return typeMap[type] || SecurityCategory.CONTENT;
}

function mapToOWASP(type: string, cwe: string): string[] {
  const owaspMap: Record<string, string[]> = {
    'XSS': ['A03'],
    'SQL Injection': ['A03'],
    'Missing Security Header': ['A05'],
    'Insecure Cookie': ['A05'],
    'Weak Cryptography': ['A02'],
  };
  
  return owaspMap[type] || ['A05']; // Default to Security Misconfiguration
}

function determineRemediationEffort(vuln: any): 'low' | 'medium' | 'high' {
  if (vuln.type.includes('Header')) return 'low';
  if (vuln.type.includes('Cookie')) return 'low';
  if (vuln.type.includes('SSL')) return 'medium';
  return 'high';
}

function calculateVulnerabilityPriority(severity: SecuritySeverity, vuln: any): number {
  const basePriority = {
    [SecuritySeverity.CRITICAL]: 10,
    [SecuritySeverity.HIGH]: 8,
    [SecuritySeverity.MEDIUM]: 5,
    [SecuritySeverity.LOW]: 3,
    [SecuritySeverity.INFO]: 1,
  };
  
  return basePriority[severity] || 1;
}

function determineExploitability(vuln: any): 'high' | 'medium' | 'low' | 'theoretical' {
  if (vuln.cvss_score >= 8) return 'high';
  if (vuln.cvss_score >= 5) return 'medium';
  if (vuln.cvss_score >= 2) return 'low';
  return 'theoretical';
}

function isOutdatedVersion(tech: any): boolean {
  if (!tech.version || !tech.latest_version) return false;
  
  // Simple version comparison (would need more sophisticated logic in production)
  const current = tech.version.split('.').map(Number);
  const latest = tech.latest_version.split('.').map(Number);
  
  for (let i = 0; i < Math.max(current.length, latest.length); i++) {
    const c = current[i] || 0;
    const l = latest[i] || 0;
    if (c < l) return true;
    if (c > l) return false;
  }
  
  return false;
}

function isSecurePort(port: number): boolean {
  const securePorts = [443, 22, 3389]; // HTTPS, SSH, RDP
  return securePorts.includes(port);
}

function getPortRecommendation(port: number): string {
  const recommendations: Record<number, string> = {
    80: 'Keep open but redirect all traffic to HTTPS',
    443: 'Keep open for HTTPS traffic',
    22: 'Restrict to specific IP addresses',
    3306: 'Close to public, use SSH tunneling',
    5432: 'Close to public, use SSH tunneling',
    21: 'Replace FTP with SFTP',
  };
  
  return recommendations[port] || 'Review necessity and restrict access';
}

function calculateGDPRScore(issues: string[], cookies: SecurityAnalysis['cookies']): number {
  let score = 100;
  
  score -= issues.length * 20;
  
  if (cookies.totalCookies > 0) {
    const securePercentage = (cookies.secureCookies / cookies.totalCookies) * 100;
    if (securePercentage < 100) score -= 10;
  }
  
  return Math.max(0, score);
}

function determineTrend(
  currentScore: number,
  previousScore?: number
): 'improving' | 'stable' | 'declining' {
  if (!previousScore) return 'stable';
  
  const diff = currentScore - previousScore;
  if (diff > 5) return 'improving';
  if (diff < -5) return 'declining';
  return 'stable';
}

function determineScanDepth(backendData: any): 'basic' | 'standard' | 'comprehensive' {
  const pagesScanned = backendData.pages_scanned || 0;
  if (pagesScanned > 100) return 'comprehensive';
  if (pagesScanned > 20) return 'standard';
  return 'basic';
}

function generateRecommendations(
  headers: SecurityAnalysis['headers'],
  ssl: SecurityAnalysis['ssl'],
  vulnerabilities: SecurityAnalysis['vulnerabilities'],
  bestPractices: SecurityAnalysis['bestPractices'],
  permissions: PermissionAnalysis
): SecurityAnalysis['recommendations'] {
  const immediate: SecurityAnalysis['recommendations']['immediate'] = [];
  
  // Add critical header recommendations
  if (headers.missing.includes('Content-Security-Policy')) {
    immediate.push({
      id: 'REC-001',
      title: 'Implement Content Security Policy',
      severity: SecuritySeverity.HIGH,
      effort: 'medium',
      impact: 'high',
      description: 'Add CSP header to prevent XSS attacks',
      implementation: 'Start with report-only mode, then enforce',
    });
  }
  
  // Add SSL recommendations
  if (ssl.configuration?.certificate.isExpiringSoon) {
    immediate.push({
      id: 'REC-002',
      title: 'Renew SSL Certificate',
      severity: SecuritySeverity.HIGH,
      effort: 'low',
      impact: 'high',
      description: 'Certificate expires in less than 30 days',
      implementation: 'Renew certificate or enable auto-renewal',
    });
  }
  
  // Add vulnerability recommendations
  vulnerabilities.critical.forEach((vuln, index) => {
    immediate.push({
      id: `REC-VULN-${index + 1}`,
      title: vuln.title,
      severity: vuln.severity,
      effort: vuln.remediation.effort,
      impact: 'high',
      description: vuln.description,
      implementation: vuln.remediation.summary,
    });
  });
  
  // Sort by priority
  immediate.sort((a, b) => {
    const severityOrder = {
      [SecuritySeverity.CRITICAL]: 5,
      [SecuritySeverity.HIGH]: 4,
      [SecuritySeverity.MEDIUM]: 3,
      [SecuritySeverity.LOW]: 2,
      [SecuritySeverity.INFO]: 1,
    };
    return severityOrder[b.severity] - severityOrder[a.severity];
  });
  
  return {
    immediate: immediate.slice(0, 5), // Top 5 immediate actions
    shortTerm: [],
    longTerm: [],
  };
}