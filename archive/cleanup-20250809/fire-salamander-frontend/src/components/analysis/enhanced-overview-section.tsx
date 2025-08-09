"use client";

import { AnalysisOverview, ScoreCardProps, MetricCardProps } from '@/types/analysis-overview';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Globe, Shield, Zap, Smartphone, Search, Lock, 
  AlertTriangle, CheckCircle, Info, TrendingUp, TrendingDown,
  Image, Link, FileText, Clock, Server, Activity
} from 'lucide-react';

interface EnhancedOverviewSectionProps {
  overview: AnalysisOverview;
}

// Mini composant pour les scores détaillés
function DetailScore({ label, value, max = 100 }: { label: string; value: number; max?: number }) {
  const percentage = (value / max) * 100;
  const color = percentage >= 70 ? 'bg-green-500' : percentage >= 50 ? 'bg-yellow-500' : 'bg-red-500';
  
  return (
    <div className="flex items-center justify-between py-1">
      <span className="text-sm text-gray-600">{label}</span>
      <div className="flex items-center space-x-2">
        <Progress value={percentage} className="w-20 h-2" />
        <span className="text-sm font-medium w-10 text-right">{value}</span>
      </div>
    </div>
  );
}

// Composant pour les Core Web Vitals
function WebVitalMetric({ 
  label, 
  value, 
  unit, 
  score, 
  testId 
}: { 
  label: string; 
  value: number | null; 
  unit: string; 
  score: string;
  testId: string;
}) {
  const getScoreColor = () => {
    switch (score) {
      case 'good': return 'bg-green-100 text-green-800 border-green-200';
      case 'needs-improvement': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'poor': return 'bg-red-100 text-red-800 border-red-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="text-center p-4 rounded-lg border" data-testid={`${testId}-metric`}>
      <div className={`inline-flex px-2 py-1 rounded-full text-xs font-medium mb-2 ${getScoreColor()}`}
           data-testid={`${testId}-status`}>
        {score === 'good' ? 'Bon' : score === 'needs-improvement' ? 'À améliorer' : score === 'poor' ? 'Mauvais' : 'N/A'}
      </div>
      <div className="text-2xl font-bold" data-testid={`${testId}-value`}>
        {value !== null ? `${value}${unit}` : '-'}
      </div>
      <div className="text-sm text-gray-600">{label}</div>
    </div>
  );
}

// Composant pour afficher les métriques avec icône
function MetricDisplay({ icon: Icon, label, value, status }: {
  icon: React.ComponentType<any>;
  label: string;
  value: string | number;
  status?: 'good' | 'warning' | 'critical';
}) {
  const statusColors = {
    good: 'text-green-600',
    warning: 'text-yellow-600',
    critical: 'text-red-600'
  };

  return (
    <div className="flex items-center space-x-3">
      <Icon className={`h-5 w-5 ${status ? statusColors[status] : 'text-gray-400'}`} />
      <div className="flex-1">
        <p className="text-sm text-gray-600">{label}</p>
        <p className="font-semibold">{value}</p>
      </div>
    </div>
  );
}

export function EnhancedOverviewSection({ overview }: EnhancedOverviewSectionProps) {
  const { scores, metrics, metadata, topIssues } = overview;

  // Déterminer les couleurs des scores
  const getScoreColor = (score: number) => {
    if (score >= 80) return 'green';
    if (score >= 60) return 'orange';
    return 'red';
  };

  return (
    <div className="space-y-6" data-testid="enhanced-overview">
      {/* Header avec metadata */}
      <div className="flex items-center justify-between text-sm text-gray-600">
        <div className="flex items-center space-x-4">
          <span className="flex items-center space-x-1">
            <Globe className="h-4 w-4" />
            <span>{metadata.domain}</span>
          </span>
          <Badge variant="outline">{metadata.analysisType}</Badge>
        </div>
        <div className="flex items-center space-x-4">
          <span>Analysé le {metadata.analyzedAt}</span>
          <span className="flex items-center space-x-1">
            <Clock className="h-4 w-4" />
            <span>{metadata.processingTime}</span>
          </span>
        </div>
      </div>

      {/* Scores principaux */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4" data-testid="score-cards">
        {/* Score Global */}
        <Card className="border-2" data-testid="score-card">
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">Score Global</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div data-testid="global-score" className="text-center flex-1">
                <div className="text-4xl font-bold" data-testid="global-score-value">
                  {scores.global}
                </div>
                <Progress 
                  value={scores.global} 
                  className="mt-2" 
                  data-testid="global-score-indicator"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Score SEO */}
        <Card data-testid="score-card">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center space-x-2">
                <Search className="h-4 w-4" />
                <span>SEO</span>
              </CardTitle>
              <span className="text-2xl font-bold">{scores.seo.score}</span>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-1" data-testid="seo-score-details">
              <DetailScore label="Optimisation titre" value={scores.seo.details.titleOptimization} />
              <DetailScore label="Meta description" value={scores.seo.details.metaOptimization} />
              <DetailScore label="Structure H1-H6" value={scores.seo.details.headingStructure} />
              <DetailScore label="Mots-clés" value={scores.seo.details.keywordUsage} />
              <DetailScore label="Qualité contenu" value={scores.seo.details.contentQuality} />
            </div>
            <div data-testid="seo-score-indicator" className="mt-2">
              <Progress value={scores.seo.score} />
            </div>
          </CardContent>
        </Card>

        {/* Score Technique */}
        <Card data-testid="score-card">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center space-x-2">
                <Shield className="h-4 w-4" />
                <span>Technique</span>
              </CardTitle>
              <span className="text-2xl font-bold">{scores.technical.score}</span>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-1" data-testid="technical-score-details">
              <DetailScore label="Crawlabilité" value={scores.technical.details.crawlability} />
              <DetailScore label="Indexabilité" value={scores.technical.details.indexability} />
              <DetailScore label="Vitesse site" value={scores.technical.details.siteSpeed} />
              <DetailScore label="Mobile" value={scores.technical.details.mobile} />
              <DetailScore label="Sécurité" value={scores.technical.details.security} />
            </div>
            <div data-testid="technical-score-indicator" className="mt-2">
              <Progress value={scores.technical.score} />
            </div>
          </CardContent>
        </Card>

        {/* Score Performance */}
        <Card data-testid="score-card">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg flex items-center space-x-2">
                <Zap className="h-4 w-4" />
                <span>Performance</span>
              </CardTitle>
              <span className="text-2xl font-bold">{scores.performance.score}</span>
            </div>
          </CardHeader>
          <CardContent data-testid="performance-details">
            <div className="grid grid-cols-2 gap-2 text-center">
              <div>
                <div className="text-xs text-gray-600">LCP</div>
                <div className="font-semibold" data-testid="lcp-value">
                  {scores.performance.coreWebVitals.lcp.value}s
                </div>
                <div className={`text-xs ${
                  scores.performance.coreWebVitals.lcp.score === 'good' ? 'text-green-600' : 
                  scores.performance.coreWebVitals.lcp.score === 'needs-improvement' ? 'text-yellow-600' : 'text-red-600'
                }`} data-testid="lcp-status">
                  {scores.performance.coreWebVitals.lcp.score}
                </div>
              </div>
              <div>
                <div className="text-xs text-gray-600">FID</div>
                <div className="font-semibold" data-testid="fid-value">
                  {scores.performance.coreWebVitals.fid.value}ms
                </div>
                <div className={`text-xs ${
                  scores.performance.coreWebVitals.fid.score === 'good' ? 'text-green-600' : 
                  scores.performance.coreWebVitals.fid.score === 'needs-improvement' ? 'text-yellow-600' : 'text-red-600'
                }`} data-testid="fid-status">
                  {scores.performance.coreWebVitals.fid.score}
                </div>
              </div>
              <div>
                <div className="text-xs text-gray-600">CLS</div>
                <div className="font-semibold" data-testid="cls-value">
                  {scores.performance.coreWebVitals.cls.value}
                </div>
                <div className={`text-xs ${
                  scores.performance.coreWebVitals.cls.score === 'good' ? 'text-green-600' : 
                  scores.performance.coreWebVitals.cls.score === 'needs-improvement' ? 'text-yellow-600' : 'text-red-600'
                }`} data-testid="cls-status">
                  {scores.performance.coreWebVitals.cls.score}
                </div>
              </div>
              <div>
                <div className="text-xs text-gray-600">TTFB</div>
                <div className="font-semibold" data-testid="ttfb-value">
                  {scores.performance.coreWebVitals.ttfb.value}ms
                </div>
                <div className={`text-xs ${
                  scores.performance.coreWebVitals.ttfb.score === 'good' ? 'text-green-600' : 
                  scores.performance.coreWebVitals.ttfb.score === 'needs-improvement' ? 'text-yellow-600' : 'text-red-600'
                }`} data-testid="ttfb-status">
                  {scores.performance.coreWebVitals.ttfb.score}
                </div>
              </div>
            </div>
            <div data-testid="performance-score-indicator" className="mt-2">
              <Progress value={scores.performance.score} />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Métriques globales */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6" data-testid="global-metrics">
        {/* Issues Summary */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center space-x-2">
              <AlertTriangle className="h-5 w-5" />
              <span>Résumé des problèmes</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
              <div className="p-3 bg-red-50 rounded-lg">
                <div className="text-2xl font-bold text-red-600" data-testid="critical-issues">
                  {metrics.issues.critical}
                </div>
                <div className="text-sm text-gray-600">Critiques</div>
              </div>
              <div className="p-3 bg-yellow-50 rounded-lg">
                <div className="text-2xl font-bold text-yellow-600" data-testid="warnings">
                  {metrics.issues.warnings}
                </div>
                <div className="text-sm text-gray-600">Avertissements</div>
              </div>
              <div className="p-3 bg-blue-50 rounded-lg">
                <div className="text-2xl font-bold text-blue-600" data-testid="notices">
                  {metrics.issues.notices}
                </div>
                <div className="text-sm text-gray-600">Remarques</div>
              </div>
              <div className="p-3 bg-green-50 rounded-lg">
                <div className="text-2xl font-bold text-green-600" data-testid="passed-checks">
                  {metrics.issues.passedChecks}
                </div>
                <div className="text-sm text-gray-600">Réussis</div>
              </div>
            </div>
            <div className="mt-4 text-center">
              <p className="text-sm text-gray-600">
                Total: <span className="font-semibold" data-testid="total-issues">{metrics.issues.total}</span> problèmes détectés sur <span data-testid="pages-analyzed">{metrics.pagesAnalyzed}</span> page(s)
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Resources & Links */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Ressources & Liens</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <MetricDisplay 
                icon={Image} 
                label="Images" 
                value={`${metrics.resources.images.total} (${metrics.resources.images.missingAlt} sans alt)`}
                status={metrics.resources.images.missingAlt > 0 ? 'warning' : 'good'}
              />
              <MetricDisplay 
                icon={Link} 
                label="Liens internes" 
                value={metrics.links.internal.total}
                data-testid="internal-links"
              />
              <MetricDisplay 
                icon={FileText} 
                label="Scripts" 
                value={metrics.resources.scripts.total}
              />
              <MetricDisplay 
                icon={Link} 
                label="Liens externes" 
                value={metrics.links.external.total}
                data-testid="external-links"
              />
              <MetricDisplay 
                icon={Server} 
                label="Temps de chargement" 
                value={`${metrics.performance.avgLoadTime}ms`}
                status={metrics.performance.avgLoadTime < 1000 ? 'good' : metrics.performance.avgLoadTime < 3000 ? 'warning' : 'critical'}
                data-testid="avg-load-time"
              />
              <MetricDisplay 
                icon={AlertTriangle} 
                label="Liens cassés" 
                value={metrics.links.internal.broken + metrics.links.external.broken}
                status={metrics.links.internal.broken + metrics.links.external.broken > 0 ? 'critical' : 'good'}
                data-testid="broken-links"
              />
            </div>
            <div className="mt-4" data-testid="total-images">
              <p className="text-sm text-gray-600">
                Total images: <span className="font-semibold">{metrics.resources.images.total}</span>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Issues */}
      {topIssues.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Problèmes prioritaires</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topIssues.map((issue) => (
                <div key={issue.id} className="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg">
                  <div className={`p-2 rounded-lg ${
                    issue.severity === 'critical' ? 'bg-red-100' :
                    issue.severity === 'warning' ? 'bg-yellow-100' : 'bg-blue-100'
                  }`}>
                    {issue.severity === 'critical' ? 
                      <AlertTriangle className="h-4 w-4 text-red-600" /> :
                      issue.severity === 'warning' ?
                      <Info className="h-4 w-4 text-yellow-600" /> :
                      <Info className="h-4 w-4 text-blue-600" />
                    }
                  </div>
                  <div className="flex-1">
                    <h4 className="font-semibold">{issue.title}</h4>
                    <p className="text-sm text-gray-600 mt-1">{issue.description}</p>
                    <div className="flex items-center space-x-4 mt-2 text-xs text-gray-500">
                      <span>Catégorie: {issue.category}</span>
                      <span>{issue.pagesAffected} page(s) affectée(s)</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}