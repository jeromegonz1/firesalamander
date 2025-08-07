/**
 * Fire Salamander - Security Analysis Section Component
 * Lead Tech Architect implementation - Enterprise-grade Security Dashboard
 * Professional OWASP-compliant interface with threat intelligence
 */

'use client';

import { useState, useEffect, useMemo, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  SecurityAnalysis,
  SecurityGrade,
  SecuritySeverity,
  SecurityCategory,
  HeaderStatus,
  SecurityVulnerability,
  SecurityHeader,
  SECURITY_SCORE_THRESHOLDS,
  OWASP_TOP_10_2021,
} from '@/types/security-analysis';
import {
  Shield,
  ShieldAlert,
  ShieldCheck,
  ShieldX,
  Lock,
  Unlock,
  Key,
  Eye,
  EyeOff,
  Globe,
  Server,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Clock,
  Zap,
  Target,
  FileText,
  Code,
  Database,
  Settings,
  Gauge,
  TrendingUp,
  TrendingDown,
  Minus,
  ExternalLink,
  Copy,
  Download,
  RefreshCw,
  Filter,
  Search,
  ChevronDown,
  ChevronRight,
  Info,
  Bug,
  Wrench,
  Calendar,
  MapPin,
  Activity,
  BarChart3,
  PieChart,
} from 'lucide-react';

interface SecurityAnalysisSectionProps {
  securityData: SecurityAnalysis;
  analysisId: string;
}

interface SecurityMetric {
  label: string;
  value: number;
  grade: SecurityGrade;
  trend?: 'up' | 'down' | 'stable';
  icon: typeof Shield;
}

export function SecurityAnalysisSection({ 
  securityData, 
  analysisId 
}: SecurityAnalysisSectionProps) {
  // État pour la gestion de l'interface enterprise
  const [activeTab, setActiveTab] = useState<'dashboard' | 'headers' | 'ssl' | 'vulnerabilities' | 'compliance' | 'recommendations'>('dashboard');
  const [expandedVulns, setExpandedVulns] = useState<Set<string>>(new Set());
  const [selectedSeverities, setSelectedSeverities] = useState<Set<SecuritySeverity>>(new Set(Object.values(SecuritySeverity)));
  const [selectedCategories, setSelectedCategories] = useState<Set<SecurityCategory>>(new Set(Object.values(SecurityCategory)));
  const [searchQuery, setSearchQuery] = useState('');
  const [showOnlyActionable, setShowOnlyActionable] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState(false);

  // Métriques dashboard calculées avec useMemo pour les performances
  const dashboardMetrics = useMemo((): SecurityMetric[] => [
    {
      label: 'Score Global',
      value: securityData.score.overall,
      grade: securityData.score.grade,
      trend: securityData.score.trend === 'improving' ? 'up' : 
             securityData.score.trend === 'declining' ? 'down' : 'stable',
      icon: Shield,
    },
    {
      label: 'Headers Sécurité',
      value: securityData.headers.score,
      grade: securityData.headers.grade,
      icon: Globe,
    },
    {
      label: 'SSL/TLS',
      value: securityData.ssl.score,
      grade: securityData.ssl.grade,
      icon: Lock,
    },
    {
      label: 'Conformité',
      value: Math.round((securityData.compliance.gdpr.score + 
                        (securityData.compliance.pci.compliant ? 100 : 50) +
                        securityData.compliance.iso27001.percentage) / 3),
      grade: getGradeFromScore(Math.round((securityData.compliance.gdpr.score + 
                        (securityData.compliance.pci.compliant ? 100 : 50) +
                        securityData.compliance.iso27001.percentage) / 3)),
      icon: FileText,
    },
  ], [securityData]);

  // Vulnérabilités filtrées avec performance optimisée
  const filteredVulnerabilities = useMemo(() => {
    const allVulns = [
      ...securityData.vulnerabilities.critical,
      ...securityData.vulnerabilities.high,
      ...securityData.vulnerabilities.medium,
      ...securityData.vulnerabilities.low,
      ...securityData.vulnerabilities.info,
    ];

    return allVulns.filter(vuln => {
      // Filtre par sévérité
      if (!selectedSeverities.has(vuln.severity)) return false;
      
      // Filtre par catégorie
      if (!selectedCategories.has(vuln.category)) return false;
      
      // Filtre par recherche
      if (searchQuery && !vuln.title.toLowerCase().includes(searchQuery.toLowerCase()) &&
          !vuln.description.toLowerCase().includes(searchQuery.toLowerCase())) {
        return false;
      }
      
      // Filtre actionnable uniquement
      if (showOnlyActionable && vuln.remediation.effort === 'high') return false;
      
      return true;
    }).sort((a, b) => {
      // Tri par priorité : sévérité puis facilité de correction
      const severityOrder = {
        [SecuritySeverity.CRITICAL]: 5,
        [SecuritySeverity.HIGH]: 4,
        [SecuritySeverity.MEDIUM]: 3,
        [SecuritySeverity.LOW]: 2,
        [SecuritySeverity.INFO]: 1,
      };
      
      const severityDiff = severityOrder[b.severity] - severityOrder[a.severity];
      if (severityDiff !== 0) return severityDiff;
      
      return b.remediation.priority - a.remediation.priority;
    });
  }, [securityData.vulnerabilities, selectedSeverities, selectedCategories, searchQuery, showOnlyActionable]);

  // Handlers optimisés avec useCallback
  const toggleVulnerabilityExpansion = useCallback((id: string) => {
    setExpandedVulns(prev => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  }, []);

  const handleSeverityFilter = useCallback((severity: SecuritySeverity) => {
    setSelectedSeverities(prev => {
      const newSet = new Set(prev);
      if (newSet.has(severity)) {
        newSet.delete(severity);
      } else {
        newSet.add(severity);
      }
      return newSet;
    });
  }, []);

  const exportSecurityReport = useCallback(() => {
    const reportData = {
      analysis: securityData,
      timestamp: new Date().toISOString(),
      analysisId,
    };
    
    const blob = new Blob([JSON.stringify(reportData, null, 2)], {
      type: 'application/json'
    });
    
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `security-report-${analysisId}-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [securityData, analysisId]);

  // Fonctions utilitaires pour l'affichage
  const getGradeColor = (grade: SecurityGrade): string => {
    const colors = {
      [SecurityGrade.A_PLUS]: 'text-green-600 bg-green-50 border-green-200',
      [SecurityGrade.A]: 'text-green-600 bg-green-50 border-green-200',
      [SecurityGrade.B]: 'text-blue-600 bg-blue-50 border-blue-200',
      [SecurityGrade.C]: 'text-yellow-600 bg-yellow-50 border-yellow-200',
      [SecurityGrade.D]: 'text-orange-600 bg-orange-50 border-orange-200',
      [SecurityGrade.F]: 'text-red-600 bg-red-50 border-red-200',
    };
    return colors[grade];
  };

  const getSeverityIcon = (severity: SecuritySeverity) => {
    switch (severity) {
      case SecuritySeverity.CRITICAL: return <ShieldX className="h-4 w-4 text-red-600" />;
      case SecuritySeverity.HIGH: return <ShieldAlert className="h-4 w-4 text-red-500" />;
      case SecuritySeverity.MEDIUM: return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      case SecuritySeverity.LOW: return <Info className="h-4 w-4 text-blue-500" />;
      case SecuritySeverity.INFO: return <CheckCircle className="h-4 w-4 text-gray-500" />;
    }
  };

  const getSeverityBadgeVariant = (severity: SecuritySeverity) => {
    switch (severity) {
      case SecuritySeverity.CRITICAL: return 'destructive';
      case SecuritySeverity.HIGH: return 'destructive';
      case SecuritySeverity.MEDIUM: return 'default';
      case SecuritySeverity.LOW: return 'secondary';
      case SecuritySeverity.INFO: return 'outline';
    }
  };

  const getHeaderStatusIcon = (status: HeaderStatus) => {
    switch (status) {
      case HeaderStatus.PRESENT: return <CheckCircle className="h-4 w-4 text-green-500" />;
      case HeaderStatus.MISSING: return <XCircle className="h-4 w-4 text-red-500" />;
      case HeaderStatus.MISCONFIGURED: return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      case HeaderStatus.DEPRECATED: return <Clock className="h-4 w-4 text-gray-500" />;
    }
  };

  const getTrendIcon = (trend?: 'up' | 'down' | 'stable') => {
    switch (trend) {
      case 'up': return <TrendingUp className="h-4 w-4 text-green-500" />;
      case 'down': return <TrendingDown className="h-4 w-4 text-red-500" />;
      default: return <Minus className="h-4 w-4 text-gray-500" />;
    }
  };

  function getGradeFromScore(score: number): SecurityGrade {
    for (const [grade, threshold] of Object.entries(SECURITY_SCORE_THRESHOLDS)) {
      if (score >= threshold.min && score <= threshold.max) {
        return grade as SecurityGrade;
      }
    }
    return SecurityGrade.F;
  }

  // Auto-refresh logic
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (autoRefresh) {
      interval = setInterval(() => {
        console.log('Auto-refreshing security data...');
        // Trigger refresh logic here
      }, 30000); // 30 seconds
    }
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [autoRefresh]);

  return (
    <div className="space-y-6">
      {/* Header avec contrôles enterprise */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Shield className="h-6 w-6 text-blue-600" />
            Analyse de Sécurité
            <Badge className={`ml-2 ${getGradeColor(securityData.score.grade)}`}>
              {securityData.score.grade}
            </Badge>
          </h2>
          <p className="text-sm text-gray-600 mt-1">
            Dernière analyse: {new Date(securityData.metadata.scanDate).toLocaleDateString('fr-FR')} 
            • {securityData.metadata.pagesScanned} pages scannées
            • {securityData.vulnerabilities.total} vulnérabilités détectées
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setAutoRefresh(!autoRefresh)}
            className={autoRefresh ? 'bg-green-50 border-green-200' : ''}
          >
            <RefreshCw className={`h-4 w-4 mr-1 ${autoRefresh ? 'animate-spin text-green-600' : ''}`} />
            Auto-refresh
          </Button>
          <Button variant="outline" size="sm" onClick={exportSecurityReport}>
            <Download className="h-4 w-4 mr-1" />
            Export
          </Button>
        </div>
      </div>

      {/* Métriques dashboard en cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {dashboardMetrics.map((metric, index) => (
          <Card key={index} className="relative overflow-hidden">
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <metric.icon className="h-8 w-8 text-blue-500 opacity-80" />
                {metric.trend && getTrendIcon(metric.trend)}
              </div>
              <div className="space-y-1">
                <p className="text-sm font-medium text-gray-600">{metric.label}</p>
                <div className="flex items-center gap-2">
                  <span className="text-2xl font-bold">{metric.value}</span>
                  <Badge className={`${getGradeColor(metric.grade)} border text-xs`}>
                    {metric.grade}
                  </Badge>
                </div>
                <Progress value={metric.value} className="h-1.5" />
              </div>
            </CardContent>
            {/* Barre colorée selon le grade */}
            <div className={`absolute bottom-0 left-0 right-0 h-1 ${metric.grade === SecurityGrade.F ? 'bg-red-500' : 
                             metric.grade === SecurityGrade.D ? 'bg-orange-500' :
                             metric.grade === SecurityGrade.C ? 'bg-yellow-500' :
                             metric.grade === SecurityGrade.B ? 'bg-blue-500' :
                             'bg-green-500'}`} />
          </Card>
        ))}
      </div>

      {/* Alertes critiques en haut */}
      {securityData.vulnerabilities.critical.length > 0 && (
        <Alert className="border-red-200 bg-red-50">
          <ShieldAlert className="h-4 w-4 text-red-600" />
          <AlertDescription className="text-red-800">
            <strong>Attention critique :</strong> {securityData.vulnerabilities.critical.length} vulnérabilité(s) critique(s) 
            détectée(s) nécessitant une action immédiate.
            <Button variant="link" className="text-red-800 p-0 h-auto font-semibold ml-2"
                    onClick={() => setActiveTab('vulnerabilities')}>
              Voir les détails →
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Navigation par onglets enterprise */}
      <Tabs value={activeTab} onValueChange={(value: any) => setActiveTab(value)} className="w-full">
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="dashboard" className="flex items-center gap-1">
            <BarChart3 className="h-4 w-4" />
            <span className="hidden sm:inline">Vue d'ensemble</span>
          </TabsTrigger>
          <TabsTrigger value="headers" className="flex items-center gap-1">
            <Globe className="h-4 w-4" />
            <span className="hidden sm:inline">Headers</span>
            {securityData.headers.missing.length > 0 && (
              <Badge variant="destructive" className="ml-1 h-5 w-5 p-0 text-xs">
                {securityData.headers.missing.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="ssl" className="flex items-center gap-1">
            <Lock className="h-4 w-4" />
            <span className="hidden sm:inline">SSL/TLS</span>
            {!securityData.ssl.enabled && (
              <ShieldX className="h-3 w-3 text-red-500 ml-1" />
            )}
          </TabsTrigger>
          <TabsTrigger value="vulnerabilities" className="flex items-center gap-1">
            <Bug className="h-4 w-4" />
            <span className="hidden sm:inline">Vulnérabilités</span>
            {securityData.vulnerabilities.total > 0 && (
              <Badge variant="destructive" className="ml-1 h-5 w-5 p-0 text-xs">
                {securityData.vulnerabilities.total}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="compliance" className="flex items-center gap-1">
            <FileText className="h-4 w-4" />
            <span className="hidden sm:inline">Conformité</span>
          </TabsTrigger>
          <TabsTrigger value="recommendations" className="flex items-center gap-1">
            <Target className="h-4 w-4" />
            <span className="hidden sm:inline">Actions</span>
            {securityData.recommendations.immediate.length > 0 && (
              <Badge className="ml-1 h-5 w-5 p-0 text-xs bg-orange-500">
                {securityData.recommendations.immediate.length}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        {/* Dashboard Overview */}
        <TabsContent value="dashboard" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            
            {/* Score trend et détails */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Gauge className="h-5 w-5 text-blue-500" />
                  Évolution du Score de Sécurité
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-3xl font-bold">{securityData.score.overall}/100</span>
                    <div className="text-right">
                      <Badge className={`${getGradeColor(securityData.score.grade)} border`}>
                        Grade {securityData.score.grade}
                      </Badge>
                      {securityData.score.previousScore && (
                        <p className="text-sm text-gray-600 mt-1">
                          Précédent: {securityData.score.previousScore}
                          {getTrendIcon(securityData.score.trend)}
                        </p>
                      )}
                    </div>
                  </div>
                  <Progress value={securityData.score.overall} className="h-3" />
                  
                  <div className="grid grid-cols-2 gap-4 pt-4 border-t">
                    <div className="text-center">
                      <p className="text-2xl font-bold text-green-600">
                        {securityData.bestPractices.implemented.length}
                      </p>
                      <p className="text-sm text-gray-600">Bonnes pratiques</p>
                    </div>
                    <div className="text-center">
                      <p className="text-2xl font-bold text-red-600">
                        {securityData.vulnerabilities.total}
                      </p>
                      <p className="text-sm text-gray-600">Vulnérabilités</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Distribution des vulnérabilités */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="h-5 w-5 text-red-500" />
                  Répartition des Vulnérabilités
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {[
                    { severity: SecuritySeverity.CRITICAL, count: securityData.vulnerabilities.critical.length, color: 'bg-red-500' },
                    { severity: SecuritySeverity.HIGH, count: securityData.vulnerabilities.high.length, color: 'bg-orange-500' },
                    { severity: SecuritySeverity.MEDIUM, count: securityData.vulnerabilities.medium.length, color: 'bg-yellow-500' },
                    { severity: SecuritySeverity.LOW, count: securityData.vulnerabilities.low.length, color: 'bg-blue-500' },
                    { severity: SecuritySeverity.INFO, count: securityData.vulnerabilities.info.length, color: 'bg-gray-500' },
                  ].map(({ severity, count, color }) => (
                    <div key={severity} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <div className={`w-3 h-3 rounded ${color}`} />
                        <span className="capitalize text-sm">{severity}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className="text-xs">{count}</Badge>
                        <div className="w-16 bg-gray-200 rounded-full h-2">
                          <div 
                            className={`${color} h-2 rounded-full`}
                            style={{ 
                              width: securityData.vulnerabilities.total > 0 
                                ? `${(count / securityData.vulnerabilities.total) * 100}%` 
                                : '0%' 
                            }}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
                
                {securityData.vulnerabilities.total === 0 && (
                  <div className="text-center py-4 text-green-600">
                    <CheckCircle className="h-12 w-12 mx-auto mb-2" />
                    <p className="font-medium">Aucune vulnérabilité détectée</p>
                  </div>
                )}
              </CardContent>
            </Card>

          </div>

          {/* Recommandations immédiates en dashboard */}
          {securityData.recommendations.immediate.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="h-5 w-5 text-orange-500" />
                  Actions Immédiates Recommandées
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-3">
                  {securityData.recommendations.immediate.slice(0, 3).map((rec, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        {getSeverityIcon(rec.severity)}
                        <div>
                          <p className="font-medium">{rec.title}</p>
                          <p className="text-sm text-gray-600">{rec.description}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className="text-xs">
                          Effort: {rec.effort}
                        </Badge>
                        <Badge className="text-xs bg-green-500">
                          Impact: {rec.impact}
                        </Badge>
                      </div>
                    </div>
                  ))}
                  {securityData.recommendations.immediate.length > 3 && (
                    <Button variant="outline" onClick={() => setActiveTab('recommendations')}>
                      Voir toutes les recommandations ({securityData.recommendations.immediate.length})
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Security Headers Analysis */}
        <TabsContent value="headers" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Globe className="h-5 w-5 text-blue-500" />
                  Headers de Sécurité
                  <Badge className={`${getGradeColor(securityData.headers.grade)} border`}>
                    {securityData.headers.grade} ({securityData.headers.score}/100)
                  </Badge>
                </div>
                <div className="text-sm text-gray-600">
                  {securityData.headers.securityHeaders.filter(h => h.status === HeaderStatus.PRESENT).length}/
                  {securityData.headers.securityHeaders.length} configurés
                </div>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {securityData.headers.securityHeaders.map((header, index) => (
                  <div key={index} className="border rounded-lg p-4">
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex items-center gap-2">
                        {getHeaderStatusIcon(header.status)}
                        <h4 className="font-mono text-sm font-medium">{header.name}</h4>
                        <Badge variant={getSeverityBadgeVariant(header.severity)}>
                          {header.severity}
                        </Badge>
                      </div>
                      {header.value && (
                        <Button variant="ghost" size="sm" onClick={() => navigator.clipboard.writeText(header.value!)}>
                          <Copy className="h-3 w-3" />
                        </Button>
                      )}
                    </div>
                    
                    {header.value ? (
                      <div className="bg-gray-50 rounded p-2 mb-2">
                        <code className="text-xs break-all">{header.value}</code>
                      </div>
                    ) : (
                      <div className="bg-red-50 rounded p-2 mb-2">
                        <span className="text-xs text-red-600">Header manquant</span>
                      </div>
                    )}
                    
                    <div className="space-y-2">
                      <p className="text-sm text-gray-700">{header.impact}</p>
                      
                      {header.recommendation && (
                        <Alert>
                          <Wrench className="h-4 w-4" />
                          <AlertDescription>
                            <strong>Recommandation:</strong> {header.recommendation}
                          </AlertDescription>
                        </Alert>
                      )}
                      
                      {header.references.length > 0 && (
                        <div className="flex gap-2">
                          {header.references.map((ref, refIndex) => (
                            <Button
                              key={refIndex}
                              variant="link"
                              size="sm"
                              className="h-auto p-0 text-xs"
                              onClick={() => window.open(ref, '_blank')}
                            >
                              <ExternalLink className="h-3 w-3 mr-1" />
                              Documentation
                            </Button>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* SSL/TLS Analysis */}
        <TabsContent value="ssl" className="space-y-4">
          {securityData.ssl.enabled ? (
            <div className="grid gap-4">
              {/* SSL Overview */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Lock className="h-5 w-5 text-green-500" />
                    Configuration SSL/TLS
                    <Badge className={`${getGradeColor(securityData.ssl.grade)} border`}>
                      {securityData.ssl.grade} ({securityData.ssl.score}/100)
                    </Badge>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    
                    {/* Certificate Info */}
                    {securityData.ssl.configuration && (
                      <div>
                        <h4 className="font-medium mb-3 flex items-center gap-2">
                          <Key className="h-4 w-4" />
                          Certificat SSL
                        </h4>
                        <div className="space-y-2 text-sm">
                          <div className="flex justify-between">
                            <span className="text-gray-600">Émetteur:</span>
                            <span className="font-mono">{securityData.ssl.configuration.certificate.issuer}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-600">Sujet:</span>
                            <span className="font-mono">{securityData.ssl.configuration.certificate.subject}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-600">Expire dans:</span>
                            <span className={securityData.ssl.configuration.certificate.isExpiringSoon ? 'text-orange-600 font-medium' : ''}>
                              {securityData.ssl.configuration.certificate.daysRemaining} jours
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-600">Taille clé:</span>
                            <span>{securityData.ssl.configuration.certificate.keySize} bits</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-gray-600">Chaîne valide:</span>
                            <span className={securityData.ssl.configuration.certificate.chainValid ? 'text-green-600' : 'text-red-600'}>
                              {securityData.ssl.configuration.certificate.chainValid ? '✓ Oui' : '✗ Non'}
                            </span>
                          </div>
                        </div>
                        
                        {securityData.ssl.configuration.certificate.san.length > 0 && (
                          <div className="mt-3">
                            <span className="text-sm text-gray-600">Domaines alternatifs:</span>
                            <div className="mt-1 flex flex-wrap gap-1">
                              {securityData.ssl.configuration.certificate.san.map((domain, index) => (
                                <Badge key={index} variant="outline" className="text-xs">
                                  {domain}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        )}
                      </div>
                    )}

                    {/* Protocol & Ciphers */}
                    {securityData.ssl.configuration && (
                      <div>
                        <h4 className="font-medium mb-3 flex items-center gap-2">
                          <Shield className="h-4 w-4" />
                          Protocoles & Chiffrement
                        </h4>
                        <div className="space-y-3">
                          <div>
                            <span className="text-sm text-gray-600 block mb-1">Protocoles supportés:</span>
                            <div className="space-y-1">
                              {securityData.ssl.configuration.protocols.map((proto, index) => (
                                <div key={index} className="flex items-center justify-between text-sm">
                                  <span className="font-mono">{proto.protocol}</span>
                                  <div className="flex items-center gap-1">
                                    {proto.enabled ? (
                                      <CheckCircle className="h-3 w-3 text-green-500" />
                                    ) : (
                                      <XCircle className="h-3 w-3 text-gray-400" />
                                    )}
                                    <Badge variant={proto.secure && proto.enabled ? "default" : proto.enabled ? "destructive" : "outline"} className="text-xs">
                                      {proto.enabled ? (proto.secure ? "Sécurisé" : "Obsolète") : "Désactivé"}
                                    </Badge>
                                  </div>
                                </div>
                              ))}
                            </div>
                          </div>

                          <div>
                            <span className="text-sm text-gray-600 block mb-1">Chiffrement recommandé:</span>
                            <div className="space-y-1">
                              {securityData.ssl.configuration.cipherSuites.slice(0, 3).map((cipher, index) => (
                                <div key={index} className="flex items-center justify-between text-xs">
                                  <span className="font-mono truncate">{cipher.name}</span>
                                  <Badge variant={cipher.recommended ? "default" : "outline"} className="text-xs">
                                    {cipher.strength}
                                  </Badge>
                                </div>
                              ))}
                            </div>
                          </div>
                        </div>
                      </div>
                    )}

                  </div>

                  {/* SSL Issues & Strengths */}
                  <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                    {securityData.ssl.issues.length > 0 && (
                      <Alert className="border-orange-200 bg-orange-50">
                        <AlertTriangle className="h-4 w-4 text-orange-600" />
                        <AlertDescription>
                          <strong>Points d'amélioration:</strong>
                          <ul className="mt-1 text-sm list-disc list-inside">
                            {securityData.ssl.issues.map((issue, index) => (
                              <li key={index}>{issue}</li>
                            ))}
                          </ul>
                        </AlertDescription>
                      </Alert>
                    )}
                    
                    {securityData.ssl.strengths.length > 0 && (
                      <Alert className="border-green-200 bg-green-50">
                        <CheckCircle className="h-4 w-4 text-green-600" />
                        <AlertDescription>
                          <strong>Points forts:</strong>
                          <ul className="mt-1 text-sm list-disc list-inside">
                            {securityData.ssl.strengths.map((strength, index) => (
                              <li key={index}>{strength}</li>
                            ))}
                          </ul>
                        </AlertDescription>
                      </Alert>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <Card>
              <CardContent className="text-center py-8">
                <Unlock className="h-16 w-16 text-red-500 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-red-600 mb-2">SSL/TLS non configuré</h3>
                <p className="text-gray-600 mb-4">
                  Ce site n'utilise pas SSL/TLS, ce qui représente un risque critique pour la sécurité.
                </p>
                <Alert className="border-red-200 bg-red-50">
                  <ShieldAlert className="h-4 w-4 text-red-600" />
                  <AlertDescription className="text-red-800">
                    <strong>Action critique requise:</strong> Implémentez SSL/TLS immédiatement pour sécuriser 
                    les communications et respecter les standards de sécurité web modernes.
                  </AlertDescription>
                </Alert>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Vulnerabilities Analysis */}
        <TabsContent value="vulnerabilities" className="space-y-4">
          {/* Filtres de vulnérabilités */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
                <div className="flex items-center gap-2">
                  <Filter className="h-4 w-4 text-gray-500" />
                  <span className="text-sm font-medium">Filtres:</span>
                </div>
                
                <div className="flex flex-wrap gap-2">
                  {Object.values(SecuritySeverity).map(severity => (
                    <Button
                      key={severity}
                      variant={selectedSeverities.has(severity) ? "default" : "outline"}
                      size="sm"
                      onClick={() => handleSeverityFilter(severity)}
                      className="text-xs capitalize"
                    >
                      {severity}
                      <Badge variant="outline" className="ml-1 h-4 w-4 p-0 text-xs">
                        {severity === SecuritySeverity.CRITICAL ? securityData.vulnerabilities.critical.length :
                         severity === SecuritySeverity.HIGH ? securityData.vulnerabilities.high.length :
                         severity === SecuritySeverity.MEDIUM ? securityData.vulnerabilities.medium.length :
                         severity === SecuritySeverity.LOW ? securityData.vulnerabilities.low.length :
                         securityData.vulnerabilities.info.length}
                      </Badge>
                    </Button>
                  ))}
                </div>

                <div className="flex items-center gap-2">
                  <div className="relative">
                    <Search className="h-4 w-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Rechercher..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-9 pr-4 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>
                  
                  <Button
                    variant={showOnlyActionable ? "default" : "outline"}
                    size="sm"
                    onClick={() => setShowOnlyActionable(!showOnlyActionable)}
                    className="text-xs"
                  >
                    <Target className="h-3 w-3 mr-1" />
                    Actionables
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Liste des vulnérabilités */}
          <div className="space-y-3">
            {filteredVulnerabilities.length === 0 ? (
              <Card>
                <CardContent className="text-center py-8">
                  <CheckCircle className="h-16 w-16 text-green-500 mx-auto mb-4" />
                  <h3 className="text-lg font-semibold text-green-600 mb-2">
                    {securityData.vulnerabilities.total === 0 
                      ? "Aucune vulnérabilité détectée" 
                      : "Aucune vulnérabilité ne correspond aux filtres"}
                  </h3>
                  <p className="text-gray-600">
                    {securityData.vulnerabilities.total === 0 
                      ? "Félicitations ! Votre site respecte les bonnes pratiques de sécurité." 
                      : "Essayez de modifier vos critères de filtrage."}
                  </p>
                </CardContent>
              </Card>
            ) : (
              filteredVulnerabilities.map((vuln, index) => (
                <Card key={vuln.id} className="border-l-4" style={{
                  borderLeftColor: vuln.severity === SecuritySeverity.CRITICAL ? '#dc2626' :
                                  vuln.severity === SecuritySeverity.HIGH ? '#ea580c' :
                                  vuln.severity === SecuritySeverity.MEDIUM ? '#ca8a04' :
                                  vuln.severity === SecuritySeverity.LOW ? '#2563eb' : '#6b7280'
                }}>
                  <CardContent className="p-4">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-start gap-3 flex-1">
                        {getSeverityIcon(vuln.severity)}
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <h4 className="font-semibold">{vuln.title}</h4>
                            <Badge variant={getSeverityBadgeVariant(vuln.severity)}>
                              {vuln.severity}
                            </Badge>
                            {vuln.cvss.score > 0 && (
                              <Badge variant="outline" className="text-xs">
                                CVSS {vuln.cvss.score}
                              </Badge>
                            )}
                            {vuln.owasp.length > 0 && (
                              <Badge variant="outline" className="text-xs">
                                OWASP {vuln.owasp.join(', ')}
                              </Badge>
                            )}
                          </div>
                          <p className="text-sm text-gray-600 mb-2">{vuln.description}</p>
                          
                          {vuln.affectedUrls.length > 0 && (
                            <div className="mb-2">
                              <span className="text-xs text-gray-500">Pages affectées: </span>
                              <div className="flex flex-wrap gap-1 mt-1">
                                {vuln.affectedUrls.slice(0, 3).map((url, urlIndex) => (
                                  <Badge key={urlIndex} variant="outline" className="text-xs">
                                    {url}
                                  </Badge>
                                ))}
                                {vuln.affectedUrls.length > 3 && (
                                  <Badge variant="outline" className="text-xs">
                                    +{vuln.affectedUrls.length - 3} autres
                                  </Badge>
                                )}
                              </div>
                            </div>
                          )}

                          <div className="flex items-center gap-4 text-xs text-gray-500">
                            <span>Priorité: {vuln.remediation.priority}/10</span>
                            <span>Effort: {vuln.remediation.effort}</span>
                            <span>Exploitabilité: {vuln.exploitability}</span>
                            {vuln.cwe && <span>CWE: {vuln.cwe}</span>}
                          </div>
                        </div>
                      </div>

                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => toggleVulnerabilityExpansion(vuln.id)}
                      >
                        {expandedVulns.has(vuln.id) ? 
                          <ChevronDown className="h-4 w-4" /> : 
                          <ChevronRight className="h-4 w-4" />
                        }
                      </Button>
                    </div>

                    {expandedVulns.has(vuln.id) && (
                      <div className="border-t pt-4 space-y-4">
                        {/* Remediation Steps */}
                        <div>
                          <h5 className="font-medium text-sm mb-2 flex items-center gap-2">
                            <Wrench className="h-4 w-4 text-blue-500" />
                            Plan de correction
                          </h5>
                          <div className="bg-blue-50 border border-blue-200 rounded p-3">
                            <p className="text-sm font-medium text-blue-800 mb-2">
                              {vuln.remediation.summary}
                            </p>
                            {vuln.remediation.steps.length > 0 && (
                              <ol className="text-xs text-blue-700 list-decimal list-inside space-y-1">
                                {vuln.remediation.steps.map((step, stepIndex) => (
                                  <li key={stepIndex}>{step}</li>
                                ))}
                              </ol>
                            )}
                          </div>
                        </div>

                        {/* References */}
                        {vuln.references.length > 0 && (
                          <div>
                            <h5 className="font-medium text-sm mb-2">Références</h5>
                            <div className="flex flex-wrap gap-2">
                              {vuln.references.map((ref, refIndex) => (
                                <Button
                                  key={refIndex}
                                  variant="outline"
                                  size="sm"
                                  onClick={() => window.open(ref.url, '_blank')}
                                  className="text-xs"
                                >
                                  <ExternalLink className="h-3 w-3 mr-1" />
                                  {ref.title}
                                </Button>
                              ))}
                            </div>
                          </div>
                        )}

                        {/* CVSS Details */}
                        {vuln.cvss.vector && (
                          <div>
                            <h5 className="font-medium text-sm mb-1">Détails CVSS</h5>
                            <code className="text-xs bg-gray-100 p-1 rounded break-all">
                              {vuln.cvss.vector}
                            </code>
                          </div>
                        )}
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))
            )}
          </div>

          {/* Statistiques des vulnérabilités */}
          <Card>
            <CardContent className="p-4">
              <div className="flex justify-between items-center text-sm">
                <span>
                  Affichage de {filteredVulnerabilities.length} sur {securityData.vulnerabilities.total} vulnérabilités
                </span>
                <div className="flex items-center gap-2">
                  <span className="text-gray-500">Priorité moyenne:</span>
                  <Badge variant="outline">
                    {filteredVulnerabilities.length > 0 
                      ? Math.round(filteredVulnerabilities.reduce((sum, v) => sum + v.remediation.priority, 0) / filteredVulnerabilities.length)
                      : 0}/10
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Compliance Analysis */}
        <TabsContent value="compliance" className="space-y-4">
          <div className="grid gap-4">
            
            {/* GDPR Compliance */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5 text-blue-500" />
                  Conformité GDPR
                  <Badge variant={securityData.compliance.gdpr.compliant ? "default" : "destructive"}>
                    {securityData.compliance.gdpr.compliant ? "Conforme" : "Non conforme"}
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-2xl font-bold">
                    {securityData.compliance.gdpr.score}/100
                  </span>
                  <Progress value={securityData.compliance.gdpr.score} className="w-48" />
                </div>
                
                {securityData.compliance.gdpr.issues.length > 0 && (
                  <div>
                    <h5 className="font-medium text-sm mb-2 text-red-600">Points non conformes:</h5>
                    <ul className="text-sm space-y-1">
                      {securityData.compliance.gdpr.issues.map((issue, index) => (
                        <li key={index} className="flex items-center gap-2">
                          <XCircle className="h-3 w-3 text-red-500" />
                          {issue}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* PCI DSS Compliance */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5 text-green-500" />
                  Conformité PCI-DSS
                  <Badge variant="outline">
                    Niveau {securityData.compliance.pci.level}
                  </Badge>
                  <Badge variant={securityData.compliance.pci.compliant ? "default" : "destructive"}>
                    {securityData.compliance.pci.compliant ? "Conforme" : "Non conforme"}
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                {securityData.compliance.pci.issues.length > 0 ? (
                  <div>
                    <h5 className="font-medium text-sm mb-2 text-red-600">Exigences non respectées:</h5>
                    <ul className="text-sm space-y-1">
                      {securityData.compliance.pci.issues.map((issue, index) => (
                        <li key={index} className="flex items-center gap-2">
                          <XCircle className="h-3 w-3 text-red-500" />
                          {issue}
                        </li>
                      ))}
                    </ul>
                  </div>
                ) : (
                  <div className="text-center py-4">
                    <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-2" />
                    <p className="font-medium text-green-600">Conformité PCI-DSS respectée</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* ISO 27001 Compliance */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="h-5 w-5 text-purple-500" />
                  Conformité ISO 27001
                  <Badge variant="outline">
                    {securityData.compliance.iso27001.percentage}%
                  </Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-gray-600">Contrôles implémentés</span>
                    <span className="font-medium">
                      {securityData.compliance.iso27001.implemented}/{securityData.compliance.iso27001.controls}
                    </span>
                  </div>
                  <Progress value={securityData.compliance.iso27001.percentage} className="w-full" />
                  
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t">
                    <div className="text-center">
                      <p className="text-2xl font-bold text-green-600">
                        {securityData.compliance.iso27001.implemented}
                      </p>
                      <p className="text-sm text-gray-600">Implémentés</p>
                    </div>
                    <div className="text-center">
                      <p className="text-2xl font-bold text-red-600">
                        {securityData.compliance.iso27001.controls - securityData.compliance.iso27001.implemented}
                      </p>
                      <p className="text-sm text-gray-600">Manquants</p>
                    </div>
                    <div className="text-center">
                      <p className="text-2xl font-bold text-blue-600">
                        {securityData.compliance.iso27001.percentage}%
                      </p>
                      <p className="text-sm text-gray-600">Conformité</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

          </div>
        </TabsContent>

        {/* Recommendations */}
        <TabsContent value="recommendations" className="space-y-4">
          <div className="space-y-4">
            
            {/* Actions immédiates */}
            {securityData.recommendations.immediate.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Zap className="h-5 w-5 text-red-500" />
                    Actions Immédiates
                    <Badge variant="destructive">
                      {securityData.recommendations.immediate.length}
                    </Badge>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {securityData.recommendations.immediate.map((rec, index) => (
                      <div key={index} className="border border-red-200 rounded-lg p-4 bg-red-50">
                        <div className="flex items-start justify-between mb-3">
                          <div className="flex items-start gap-3">
                            {getSeverityIcon(rec.severity)}
                            <div>
                              <h4 className="font-semibold text-red-800">{rec.title}</h4>
                              <p className="text-sm text-red-700 mt-1">{rec.description}</p>
                            </div>
                          </div>
                          <div className="flex gap-2">
                            <Badge className="text-xs bg-orange-100 text-orange-800">
                              Effort: {rec.effort}
                            </Badge>
                            <Badge className="text-xs bg-green-100 text-green-800">
                              Impact: {rec.impact}
                            </Badge>
                          </div>
                        </div>
                        
                        <Alert className="border-red-300 bg-red-100">
                          <Wrench className="h-4 w-4 text-red-600" />
                          <AlertDescription className="text-red-800">
                            <strong>Implémentation:</strong> {rec.implementation}
                          </AlertDescription>
                        </Alert>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Summary des actions */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5 text-blue-500" />
                  Résumé des Actions de Sécurité
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-3xl font-bold text-red-600">
                      {securityData.recommendations.immediate.length}
                    </p>
                    <p className="text-sm text-gray-600">Actions immédiates</p>
                    <p className="text-xs text-gray-500 mt-1">À faire aujourd'hui</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-3xl font-bold text-orange-600">
                      {securityData.recommendations.shortTerm.length}
                    </p>
                    <p className="text-sm text-gray-600">Court terme</p>
                    <p className="text-xs text-gray-500 mt-1">Cette semaine</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-3xl font-bold text-blue-600">
                      {securityData.recommendations.longTerm.length}
                    </p>
                    <p className="text-sm text-gray-600">Long terme</p>
                    <p className="text-xs text-gray-500 mt-1">Ce mois</p>
                  </div>
                </div>
                
                {securityData.recommendations.immediate.length === 0 && 
                 securityData.recommendations.shortTerm.length === 0 && 
                 securityData.recommendations.longTerm.length === 0 && (
                  <div className="text-center py-8">
                    <CheckCircle className="h-16 w-16 text-green-500 mx-auto mb-4" />
                    <h3 className="text-lg font-semibold text-green-600 mb-2">
                      Aucune action requise
                    </h3>
                    <p className="text-gray-600">
                      Votre configuration de sécurité respecte les meilleures pratiques.
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>

          </div>
        </TabsContent>

      </Tabs>

      {/* Footer avec métadonnées */}
      <Card className="bg-gray-50">
        <CardContent className="p-4">
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4 text-sm text-gray-600">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-1">
                <Calendar className="h-4 w-4" />
                <span>Analysé le {new Date(securityData.metadata.scanDate).toLocaleDateString('fr-FR')}</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-4 w-4" />
                <span>Durée: {securityData.metadata.scanDuration}s</span>
              </div>
              <div className="flex items-center gap-1">
                <MapPin className="h-4 w-4" />
                <span>Profondeur: {securityData.metadata.scanDepth}</span>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <span>Fire Salamander Security v{securityData.metadata.toolVersion}</span>
              <Badge variant="outline" className="text-xs">
                ID: {analysisId.slice(0, 8)}
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}