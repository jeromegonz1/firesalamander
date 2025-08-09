"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { TechnicalAnalysisSectionSimple } from "@/components/analysis/technical-analysis-section-simple";
import { mapBackendToTechnicalAnalysis } from "@/lib/mappers/technical-mapper";
import { TechnicalAnalysis } from "@/types/technical-analysis";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle } from "lucide-react";

interface TechnicalIssue {
  type: string;
  severity: 'critical' | 'warning' | 'info';
  message: string;
  solution: string;
  pages_affected: number;
  impact: 'high' | 'medium' | 'low';
}

interface MetaTagStatus {
  tag: string;
  status: 'present' | 'missing' | 'warning';
  count: number;
  recommendation: string;
}

interface CoreWebVital {
  metric: string;
  value: number;
  unit: string;
  status: 'good' | 'needs-improvement' | 'poor';
  threshold_good: number;
  threshold_poor: number;
}

interface SecurityCheck {
  check: string;
  status: 'pass' | 'fail' | 'warning';
  description: string;
  recommendation?: string;
}

interface CrawlabilityCheck {
  check: string;
  status: 'pass' | 'fail' | 'warning';
  details: string;
  impact: string;
}

// Composant Technical Health Score
function TechnicalHealthScore({ score }: { score: number }) {
  const getScoreColor = (score: number) => {
    if (score >= 80) return "text-green-600 bg-green-50 border-green-200";
    if (score >= 60) return "text-orange-600 bg-orange-50 border-orange-200";
    return "text-red-600 bg-red-50 border-red-200";
  };

  const percentage = score;

  return (
    <Card className={`${getScoreColor(score)}`} data-testid="technical-health-score">
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">Score de Santé Technique</p>
            <p className="text-3xl font-bold" data-testid="health-score-value">{score}<span className="text-lg">/100</span></p>
          </div>
          <div className="relative h-16 w-16">
            <svg className="h-16 w-16 transform -rotate-90" viewBox="0 0 100 100">
              <circle
                cx="50"
                cy="50"
                r="40"
                stroke="currentColor"
                strokeWidth="8"
                fill="none"
                className="opacity-20"
              />
              <circle
                cx="50"
                cy="50"
                r="40"
                stroke="currentColor"
                strokeWidth="8"
                fill="none"
                strokeDasharray={`${percentage * 2.51} 251`}
                className="transition-all duration-500"
              />
            </svg>
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-sm font-bold">{Math.round(percentage)}%</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant HTML Structure Analysis
function HTMLStructureAnalysis({ issues }: { issues: TechnicalIssue[] }) {
  return (
    <div className="space-y-4" data-testid="html-validation">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="p-4 bg-green-50 rounded-lg border border-green-200" data-testid="doctype-check">
          <div className="flex items-center space-x-2">
            <CheckCircle2 className="h-5 w-5 text-green-600" />
            <span className="font-medium text-green-900">DOCTYPE</span>
          </div>
          <p className="text-sm text-green-700 mt-1">HTML5 valide détecté</p>
        </div>
        
        <div className="p-4 bg-yellow-50 rounded-lg border border-yellow-200" data-testid="semantic-tags">
          <div className="flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5 text-yellow-600" />
            <span className="font-medium text-yellow-900">Balises Sémantiques</span>
          </div>
          <p className="text-sm text-yellow-700 mt-1">Quelques améliorations possibles</p>
        </div>
        
        <div className="p-4 bg-red-50 rounded-lg border border-red-200">
          <div className="flex items-center space-x-2">
            <X className="h-5 w-5 text-red-600" />
            <span className="font-medium text-red-900">Validation W3C</span>
          </div>
          <p className="text-sm text-red-700 mt-1">3 erreurs détectées</p>
        </div>
      </div>
      
      <div className="space-y-3">
        <h4 className="font-semibold">Problèmes détectés:</h4>
        {issues.slice(0, 3).map((issue, index) => (
          <div key={index} className="p-3 border rounded-lg hover:bg-gray-50 cursor-pointer" data-testid="technical-issue">
            <div className="flex items-start justify-between">
              <div className="flex items-start space-x-3">
                <div className={`p-1 rounded ${
                  issue.severity === 'critical' ? 'bg-red-100 text-red-600' :
                  issue.severity === 'warning' ? 'bg-yellow-100 text-yellow-600' :
                  'bg-blue-100 text-blue-600'
                }`}>
                  {issue.severity === 'critical' && <AlertTriangle className="h-4 w-4" />}
                  {issue.severity === 'warning' && <Info className="h-4 w-4" />}
                  {issue.severity === 'info' && <CheckCircle2 className="h-4 w-4" />}
                </div>
                <div>
                  <p className="font-medium">{issue.type}</p>
                  <p className="text-sm text-gray-600">{issue.message}</p>
                  <p className="text-xs text-gray-500 mt-1">{issue.pages_affected} pages affectées</p>
                </div>
              </div>
              <Badge variant={issue.severity === 'critical' ? 'destructive' : 'secondary'}>
                {issue.severity}
              </Badge>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Composant Meta Tags Audit
function MetaTagsAudit({ metaTags }: { metaTags: MetaTagStatus[] }) {
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'present': return <CheckCircle2 className="h-4 w-4 text-green-600" />;
      case 'missing': return <X className="h-4 w-4 text-red-600" />;
      case 'warning': return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
      default: return <Info className="h-4 w-4 text-gray-600" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'present': return 'bg-green-50 border-green-200';
      case 'missing': return 'bg-red-50 border-red-200';
      case 'warning': return 'bg-yellow-50 border-yellow-200';
      default: return 'bg-gray-50 border-gray-200';
    }
  };

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {metaTags.map((meta, index) => (
          <div key={index} className={`p-4 border rounded-lg ${getStatusColor(meta.status)}`} data-testid={`${meta.tag.toLowerCase().replace(' ', '-')}-tag-status`}>
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center space-x-2">
                {getStatusIcon(meta.status)}
                <span className="font-medium">{meta.tag}</span>
              </div>
              <Badge variant={meta.status === 'present' ? 'default' : 'secondary'}>
                {meta.count} trouvé{meta.count > 1 ? 's' : ''}
              </Badge>
            </div>
            <p className="text-sm text-gray-700">{meta.recommendation}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

// Composant Core Web Vitals
function CoreWebVitalsDisplay({ vitals }: { vitals: CoreWebVital[] }) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'good': return 'bg-green-50 border-green-200 text-green-700';
      case 'needs-improvement': return 'bg-yellow-50 border-yellow-200 text-yellow-700';
      case 'poor': return 'bg-red-50 border-red-200 text-red-700';
      default: return 'bg-gray-50 border-gray-200 text-gray-700';
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'good': return 'Bon';
      case 'needs-improvement': return 'À améliorer';
      case 'poor': return 'Mauvais';
      default: return 'N/A';
    }
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {vitals.map((vital, index) => (
        <div key={index} className={`text-center p-4 rounded-lg border ${getStatusColor(vital.status)}`} data-testid={`${vital.metric.toLowerCase().replace(' ', '-')}-metric`}>
          <div className="text-2xl font-bold">{vital.value}{vital.unit}</div>
          <div className="text-sm text-gray-600 mb-2">{vital.metric}</div>
          <Badge className={`${getStatusColor(vital.status)} border-0`}>
            {getStatusBadge(vital.status)}
          </Badge>
        </div>
      ))}
    </div>
  );
}

// Composant Security Audit
function SecurityAuditDisplay({ securityChecks }: { securityChecks: SecurityCheck[] }) {
  return (
    <div className="space-y-4">
      {securityChecks.map((check, index) => (
        <div key={index} className="flex items-start justify-between p-4 border rounded-lg hover:bg-gray-50" data-testid={`${check.check.toLowerCase().replace(/\s+/g, '-')}-check`}>
          <div className="flex items-start space-x-3">
            <div className={`p-1 rounded ${
              check.status === 'pass' ? 'bg-green-100 text-green-600' :
              check.status === 'fail' ? 'bg-red-100 text-red-600' :
              'bg-yellow-100 text-yellow-600'
            }`}>
              {check.status === 'pass' && <CheckCircle2 className="h-4 w-4" />}
              {check.status === 'fail' && <X className="h-4 w-4" />}
              {check.status === 'warning' && <AlertTriangle className="h-4 w-4" />}
            </div>
            <div>
              <p className="font-medium">{check.check}</p>
              <p className="text-sm text-gray-600">{check.description}</p>
              {check.recommendation && (
                <p className="text-sm text-blue-600 mt-1">{check.recommendation}</p>
              )}
            </div>
          </div>
          <Badge variant={check.status === 'pass' ? 'default' : check.status === 'fail' ? 'destructive' : 'secondary'}>
            {check.status === 'pass' ? '✓ Valide' : check.status === 'fail' ? '✗ Échec' : '⚠ Attention'}
          </Badge>
        </div>
      ))}
    </div>
  );
}

// Composant Mobile Optimization
function MobileOptimizationCheck({ mobileChecks }: { mobileChecks: any[] }) {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="p-4 bg-green-50 rounded-lg border border-green-200" data-testid="viewport-check">
          <div className="flex items-center space-x-2">
            <CheckCircle2 className="h-5 w-5 text-green-600" />
            <span className="font-medium text-green-900">Viewport Meta</span>
          </div>
          <p className="text-sm text-green-700 mt-1">Configuré correctement</p>
        </div>
        
        <div className="p-4 bg-green-50 rounded-lg border border-green-200" data-testid="mobile-friendly-check">
          <div className="flex items-center space-x-2">
            <CheckCircle2 className="h-5 w-5 text-green-600" />
            <span className="font-medium text-green-900">Mobile-Friendly</span>
          </div>
          <p className="text-sm text-green-700 mt-1">Test Google réussi</p>
        </div>
        
        <div className="p-4 bg-yellow-50 rounded-lg border border-yellow-200" data-testid="touch-targets-check">
          <div className="flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5 text-yellow-600" />
            <span className="font-medium text-yellow-900">Touch Targets</span>
          </div>
          <p className="text-sm text-yellow-700 mt-1">Quelques éléments trop petits</p>
        </div>
      </div>
    </div>
  );
}

// Composant Crawlability Analysis
function CrawlabilityAnalysis({ crawlChecks }: { crawlChecks: CrawlabilityCheck[] }) {
  return (
    <div className="space-y-4">
      {crawlChecks.map((check, index) => (
        <div key={index} className="p-4 border rounded-lg hover:bg-gray-50" data-testid={`${check.check.toLowerCase().replace(/\s+/g, '-')}-check`}>
          <div className="flex items-start justify-between">
            <div className="flex items-start space-x-3">
              <div className={`p-1 rounded ${
                check.status === 'pass' ? 'bg-green-100 text-green-600' :
                check.status === 'fail' ? 'bg-red-100 text-red-600' :
                'bg-yellow-100 text-yellow-600'
              }`}>
                {check.status === 'pass' && <CheckCircle2 className="h-4 w-4" />}
                {check.status === 'fail' && <X className="h-4 w-4" />}
                {check.status === 'warning' && <AlertTriangle className="h-4 w-4" />}
              </div>
              <div>
                <p className="font-medium">{check.check}</p>
                <p className="text-sm text-gray-600">{check.details}</p>
                <p className="text-sm text-gray-500 mt-1">Impact: {check.impact}</p>
              </div>
            </div>
            <Badge variant={check.status === 'pass' ? 'default' : check.status === 'fail' ? 'destructive' : 'secondary'}>
              {check.status}
            </Badge>
          </div>
        </div>
      ))}
    </div>
  );
}

// Composant Status Indicator
function StatusIndicator({ status, label }: { status: 'success' | 'warning' | 'error', label: string }) {
  const getStatusClass = () => {
    switch (status) {
      case 'success': return 'bg-green-100 text-green-700 border-green-200';
      case 'warning': return 'bg-yellow-100 text-yellow-700 border-yellow-200'; 
      case 'error': return 'bg-red-100 text-red-700 border-red-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  return (
    <Badge className={`${getStatusClass()} border`} data-testid="status-indicator">
      {label}
    </Badge>
  );
}

// Composant Issue Details Panel
function IssueDetailsPanel({ 
  isOpen, 
  onClose, 
  issue 
}: { 
  isOpen: boolean; 
  onClose: () => void; 
  issue: TechnicalIssue | null; 
}) {
  if (!isOpen || !issue) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" data-testid="issue-details-panel">
      <div className="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">Détails du Problème</h2>
          <Button variant="ghost" size="sm" onClick={onClose} data-testid="close-details">
            <X className="h-4 w-4" />
          </Button>
        </div>
        
        <div className="space-y-6" data-testid="issue-description">
          <div>
            <h3 className="font-semibold mb-2">Type de problème</h3>
            <p className="text-gray-700">{issue.type}</p>
          </div>
          
          <div>
            <h3 className="font-semibold mb-2">Description</h3>
            <p className="text-gray-700">{issue.message}</p>
          </div>
          
          <div data-testid="issue-solution">
            <h3 className="font-semibold mb-2">Solution recommandée</h3>
            <div className="p-4 bg-blue-50 rounded-lg">
              <p className="text-blue-800">{issue.solution}</p>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="p-3 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900">Sévérité</h4>
              <Badge variant={issue.severity === 'critical' ? 'destructive' : 'secondary'}>
                {issue.severity}
              </Badge>
            </div>
            <div className="p-3 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-900">Pages affectées</h4>
              <p className="text-2xl font-bold text-gray-700">{issue.pages_affected}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function AnalysisTechnicalPage() {
  const params = useParams();
  const [analysis, setAnalysis] = useState<any>(null);
  const [technicalData, setTechnicalData] = useState<TechnicalAnalysis | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params.id as string;

  useEffect(() => {
    const fetchTechnicalAnalysis = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}/technical`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
          
          // Map backend data to our TechnicalAnalysis interface
          const mappedData = mapBackendToTechnicalAnalysis(data.data);
          setTechnicalData(mappedData);
        } else {
          setError("Analyse technique non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse technique:", err);
        setError("Erreur de connexion");
        
        // Fallback to mock data for development
        const mockData = createMockTechnicalData();
        setTechnicalData(mockData);
        console.log('Using mock technical data:', mockData);
      } finally {
        setLoading(false);
      }
    };

    fetchTechnicalAnalysis();
  }, [analysisId]);

  const createMockTechnicalData = (): TechnicalAnalysis => {
    return {
      pageAnalysis: [
        {
          url: 'https://example.com',
          statusCode: 200,
          loadTime: 1200,
          size: 45678,
          lastCrawled: new Date().toISOString(),
          depth: 0,
          title: {
            content: 'Example Page Title',
            length: 18,
            hasKeyword: true,
            issues: [],
            recommendations: []
          },
          metaDescription: {
            content: 'This is a good meta description',
            length: 31,
            hasKeyword: false,
            issues: ['meta-too-short'],
            recommendations: ['Rallonger la meta description']
          },
          headings: {
            h1: ['Main Heading'],
            h2: ['Section 1', 'Section 2'],
            h3: ['Subsection'],
            structure: 'good',
            issues: [],
            recommendations: []
          },
          canonical: 'https://example.com',
          robots: 'index,follow',
          schema: [
            { type: 'Organization', valid: true, errors: [] }
          ],
          openGraph: {
            'og:title': 'Example Title',
            'og:description': 'Example Description'
          },
          images: [
            {
              src: '/image1.jpg',
              alt: 'Good image description',
              size: 12345,
              issues: [],
              recommendations: []
            }
          ],
          links: {
            internal: 15,
            external: 5,
            broken: [],
            nofollow: 2,
            totalLinks: 20
          }
        }
      ],
      globalIssues: {
        duplicateContent: [],
        duplicateTitles: [],
        duplicateMeta: [],
        missingTitles: [],
        missingMeta: [],
        brokenLinks: [],
        orphanPages: [],
        redirectChains: [],
        largePages: [],
        slowPages: []
      },
      crawlability: {
        robotsTxt: {
          exists: true,
          valid: true,
          userAgents: ['*'],
          disallowedPaths: ['/admin'],
          allowedPaths: [],
          sitemapUrls: ['https://example.com/sitemap.xml'],
          issues: []
        },
        sitemap: {
          exists: true,
          url: 'https://example.com/sitemap.xml',
          valid: true,
          format: 'xml',
          pagesInSitemap: 25,
          images: 0,
          videos: 0,
          issues: []
        },
        crawlBudget: {
          totalPages: 30,
          crawlablePages: 25,
          blockedPages: 5,
          pagesPerLevel: { 0: 1, 1: 10, 2: 14 },
          averageCrawlTime: 1.2,
          crawlEfficiency: 83.3
        }
      },
      metrics: {
        totalPages: 30,
        pagesWithIssues: 5,
        avgLoadTime: 1500,
        avgPageSize: 35000,
        totalIssues: 12,
        issuesByType: {
          'meta-too-short': 3,
          'image-no-alt': 5,
          'title-missing': 2
        },
        healthScore: 75
      },
      config: {
        maxDepth: 3,
        respectRobots: true,
        userAgent: 'Fire-Salamander-Bot/1.0',
        crawlDelay: 1000
      },
      status: {
        analysisDate: new Date().toISOString(),
        crawlDuration: 5,
        crawlStatus: 'completed',
        lastUpdate: new Date().toISOString(),
        version: '1.0'
      }
    };
  };

  // Legacy mock data pour compatibilité
  const mockTechnicalIssues: TechnicalIssue[] = [
    {
      type: 'Images sans attribut alt',
      severity: 'warning',
      message: '15 images n\'ont pas d\'attribut alt descriptif',
      solution: 'Ajouter des descriptions alt pertinentes pour toutes les images',
      pages_affected: 8,
      impact: 'medium'
    },
    {
      type: 'Meta description manquante',
      severity: 'critical',
      message: '5 pages n\'ont pas de meta description',
      solution: 'Créer des meta descriptions uniques de 150-160 caractères',
      pages_affected: 5,
      impact: 'high'
    },
    {
      type: 'Liens internes cassés',
      severity: 'info',
      message: '3 liens internes retournent une erreur 404',
      solution: 'Vérifier et corriger les URLs des liens cassés',
      pages_affected: 3,
      impact: 'low'
    }
  ];

  const mockMetaTags: MetaTagStatus[] = [
    {
      tag: 'Title Tag',
      status: 'present',
      count: 12,
      recommendation: 'Toutes les pages ont un title tag unique'
    },
    {
      tag: 'Meta Description',
      status: 'missing',
      count: 7,
      recommendation: '5 pages manquent de meta description'
    },
    {
      tag: 'Open Graph Tags',
      status: 'warning',
      count: 8,
      recommendation: 'Certaines balises OG peuvent être améliorées'
    }
  ];

  const mockCoreWebVitals: CoreWebVital[] = [
    {
      metric: 'Largest Contentful Paint',
      value: 1.2,
      unit: 's',
      status: 'good',
      threshold_good: 2.5,
      threshold_poor: 4.0
    },
    {
      metric: 'First Input Delay',
      value: 19,
      unit: 'ms',
      status: 'needs-improvement',
      threshold_good: 100,
      threshold_poor: 300
    },
    {
      metric: 'Cumulative Layout Shift',
      value: 0.17,
      unit: '',
      status: 'poor',
      threshold_good: 0.1,
      threshold_poor: 0.25
    }
  ];

  const mockSecurityChecks: SecurityCheck[] = [
    {
      check: 'HTTPS',
      status: 'pass',
      description: 'Le site utilise HTTPS correctement',
      recommendation: undefined
    },
    {
      check: 'Mixed Content',
      status: 'fail',
      description: 'Contenu HTTP détecté sur pages HTTPS',
      recommendation: 'Migrer toutes les ressources vers HTTPS'
    },
    {
      check: 'Security Headers',
      status: 'warning',
      description: 'Certains en-têtes de sécurité manquent',
      recommendation: 'Ajouter CSP et X-Frame-Options headers'
    }
  ];

  const mockCrawlChecks: CrawlabilityCheck[] = [
    {
      check: 'Robots.txt',
      status: 'pass',
      details: 'Fichier robots.txt présent et bien configuré',
      impact: 'Bon pour l\'indexation'
    },
    {
      check: 'Sitemap XML',
      status: 'pass',
      details: 'Sitemap présent avec 12 URLs',
      impact: 'Facilite la découverte du contenu'
    },
    {
      check: 'Liens internes',
      status: 'warning',
      details: 'Structure de liens internes peut être améliorée',
      impact: 'Distribution du PageRank sous-optimale'
    }
  ];




  console.log('State check:', { loading, error, hasTechnicalData: !!technicalData });

  if (loading) {
    return <div>LOADING STATE</div>;
  }

  if (error && !technicalData) {
    return <div>ERROR STATE: {error}</div>;
  }

  if (!technicalData) {
    return <div>NO DATA STATE</div>;
  }

  console.log('Rendering technical page', { loading, error, technicalData: !!technicalData });

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Debug Info */}
      {process.env.NODE_ENV === 'development' && (
        <div className="mb-4 p-4 bg-gray-100 rounded">
          <p>Technical Data Status: {technicalData ? 'Loaded' : 'Missing'}</p>
          <p>Pages Count: {technicalData?.pageAnalysis?.length || 0}</p>
          <p>Health Score: {technicalData.metrics.healthScore}</p>
          <p>Analysis ID: {analysisId}</p>
        </div>
      )}

      {/* Use the simplified TechnicalAnalysisSection component */}
      <TechnicalAnalysisSectionSimple 
        technicalData={technicalData} 
        analysisId={analysisId} 
      />
    </div>
  );
}
