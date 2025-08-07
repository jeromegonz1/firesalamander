"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { AnalysisLayout } from "@/components/layout/analysis-layout";
import { OverviewSection } from "@/components/analysis/overview-section";
import { mapBackendToOverview } from "@/lib/mappers/overview-mapper";
import { mapBackendToOverviewAnalysis } from "@/lib/mappers/overview-analysis-mapper";
import { OverviewAnalysis, OverviewGrade } from "@/types/overview";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Activity, BarChart3 } from 'lucide-react';

interface AnalysisData {
  id: string;
  url: string;
  analyzed_at: string;
  overview_data?: any;
}

export default function AnalysisOverviewPage() {
  const params = useParams();
  const [overviewData, setOverviewData] = useState<OverviewAnalysis | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params.id as string;

  useEffect(() => {
    const fetchOverviewAnalysis = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}/overview`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
          
          // Map backend data to our OverviewAnalysis interface
          const overviewData = mapBackendToOverview(data.data);
          const mappedData = mapBackendToOverviewAnalysis(overviewData);
          setOverviewData(mappedData);
        } else {
          setError("Vue d'ensemble non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de la vue d'ensemble:", err);
        setError("Erreur de connexion");
        
        // Fallback to mock data for development
        const mockData = createMockOverviewData();
        setOverviewData(mockData);
        console.log('Using mock overview data:', mockData);
      } finally {
        setLoading(false);
      }
    };

    fetchOverviewAnalysis();
  }, [analysisId]);

  const createMockOverviewData = (): OverviewAnalysis => {
    return {
      score: {
        overall: 78,
        grade: OverviewGrade.B,
        trend: 'improving',
        previousScore: 72,
        lastUpdated: new Date().toISOString(),
      },

      modules: {
        technical: {
          score: 85,
          grade: OverviewGrade.A,
          status: 'completed',
          lastUpdated: new Date().toISOString(),
          keyMetrics: {
            metaTags: '92%',
            headingStructure: '88%',
            htmlValidation: '95%',
          },
          criticalIssues: 2,
          improvements: 15,
        },
        performance: {
          score: 72,
          grade: OverviewGrade.B,
          status: 'completed',
          lastUpdated: new Date().toISOString(),
          keyMetrics: {
            lcp: '2.1s',
            fid: '89ms',
            cls: '0.08',
          },
          criticalIssues: 3,
          improvements: 8,
        },
        content: {
          score: 69,
          grade: OverviewGrade.C,
          status: 'completed',
          lastUpdated: new Date().toISOString(),
          keyMetrics: {
            readability: '78%',
            keywordDensity: '2.3%',
            duplicateContent: '5%',
          },
          criticalIssues: 5,
          improvements: 12,
        },
        security: {
          score: 78,
          grade: OverviewGrade.B,
          status: 'completed',
          lastUpdated: new Date().toISOString(),
          keyMetrics: {
            https: '100%',
            headers: '65%',
            vulnerabilities: '3 low',
          },
          criticalIssues: 1,
          improvements: 6,
        },
        backlinks: {
          score: 72,
          grade: OverviewGrade.B,
          status: 'completed',
          lastUpdated: new Date().toISOString(),
          keyMetrics: {
            totalBacklinks: '15,420',
            domainAuthority: '42.8',
            toxicLinks: '3.2%',
          },
          criticalIssues: 4,
          improvements: 23,
        },
      },

      kpis: {
        organicTraffic: {
          current: 45680,
          change: 12.5,
          trend: 'up',
        },
        keywordRankings: {
          total: 2340,
          top10: 189,
          top3: 45,
          newRankings: 23,
        },
        domainAuthority: {
          current: 42,
          change: 3,
          trend: 'up',
        },
        pageSpeed: {
          mobile: 72,
          desktop: 89,
          change: 5,
        },
      },

      issues: {
        critical: [
          {
            id: 'issue-1',
            title: 'Pages sans meta description',
            severity: 'critical',
            module: 'technical',
            description: '45 pages importantes n\'ont pas de meta description',
            impact: 'Réduction significative du CTR dans les SERPs',
            effort: 'low',
            affectedPages: 45,
            estimatedFix: '2-3 heures',
          },
          {
            id: 'issue-2',
            title: 'Images non optimisées',
            severity: 'critical',
            module: 'performance',
            description: 'Les images représentent 78% du poids des pages',
            impact: 'Temps de chargement très lents',
            effort: 'medium',
            affectedPages: 120,
            estimatedFix: '1-2 jours',
          },
          {
            id: 'issue-3',
            title: 'Contenu dupliqué détecté',
            severity: 'critical',
            module: 'content',
            description: '23 pages avec du contenu identique ou très similaire',
            impact: 'Cannibalisation des mots-clés',
            effort: 'high',
            affectedPages: 23,
            estimatedFix: '1 semaine',
          },
        ],
        warnings: [
          {
            id: 'warning-1',
            title: 'Certificat SSL expire bientôt',
            severity: 'high',
            module: 'security',
            description: 'Le certificat SSL expire dans 25 jours',
            impact: 'Perte de confiance et erreurs navigateur',
            effort: 'low',
            affectedPages: 0,
            estimatedFix: '30 minutes',
          },
          {
            id: 'warning-2',
            title: 'Liens toxiques détectés',
            severity: 'medium',
            module: 'backlinks',
            description: '47 backlinks provenant de domaines suspects',
            impact: 'Risque de pénalité algorithmique',
            effort: 'low',
            affectedPages: 0,
            estimatedFix: '1 heure',
          },
        ],
        totalIssues: 87,
        resolvedIssues: 23,
      },

      recommendations: {
        immediate: [
          {
            id: 'rec-1',
            title: 'Ajouter les meta descriptions manquantes',
            category: 'immediate',
            module: 'technical',
            description: 'Rédiger des meta descriptions uniques pour les 45 pages prioritaires',
            expectedImpact: 'Amélioration du CTR de 15-25%',
            effort: 'low',
            timeline: '2-3 jours',
            priority: 1,
          },
          {
            id: 'rec-2',
            title: 'Optimiser les images principales',
            category: 'immediate',
            module: 'performance',
            description: 'Compresser et convertir en WebP les 50 images les plus lourdes',
            expectedImpact: 'Réduction de 40% du temps de chargement',
            effort: 'medium',
            timeline: '1 semaine',
            priority: 2,
          },
        ],
        quickWins: [
          {
            id: 'quick-1',
            title: 'Activer la compression Gzip',
            category: 'quick-win',
            module: 'performance',
            description: 'Configurer la compression sur le serveur web',
            expectedImpact: 'Réduction de 60% de la taille des pages',
            effort: 'low',
            timeline: '1 heure',
            priority: 1,
          },
          {
            id: 'quick-2',
            title: 'Corriger les liens internes brisés',
            category: 'quick-win',
            module: 'technical',
            description: 'Réparer les 12 liens internes retournant des erreurs 404',
            expectedImpact: 'Amélioration de l\'expérience utilisateur',
            effort: 'low',
            timeline: '2 heures',
            priority: 2,
          },
        ],
        totalRecommendations: 28,
        estimatedImpact: 'Amélioration de 15-20 points du score global',
      },

      competition: {
        position: 3,
        totalCompetitors: 8,
        aboveAverage: true,
        marketShare: 12.8,
        competitorGains: 2,
      },

      progress: {
        completedTasks: 67,
        totalTasks: 95,
        lastActivity: 'Il y a 2 heures',
        nextMilestone: 'Score 80+ dans 2 semaines',
      },

      siteInfo: {
        url: 'https://example.com',
        title: 'Example Company - Solutions SEO',
        description: 'Plateforme leader en solutions SEO et marketing digital',
        lastCrawled: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(), // 6 hours ago
        pagesAnalyzed: 247,
        crawlErrors: 3,
      },
    };
  };

  console.log('State check:', { loading, error, hasOverviewData: !!overviewData });

  if (loading) {
    return (
      <AnalysisLayout>
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <Activity className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-500" />
            <p className="text-gray-600">Chargement de la vue d'ensemble...</p>
          </div>
        </div>
      </AnalysisLayout>
    );
  }

  if (error && !overviewData) {
    return (
      <AnalysisLayout>
        <div className="flex items-center justify-center min-h-[400px]">
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2 text-red-600">
                <AlertTriangle className="h-5 w-5" />
                <span>Erreur</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-gray-600 mb-4">{error}</p>
              <Button 
                onClick={() => window.location.reload()} 
                className="w-full"
              >
                Réessayer
              </Button>
            </CardContent>
          </Card>
        </div>
      </AnalysisLayout>
    );
  }

  if (!overviewData) {
    return (
      <AnalysisLayout>
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <BarChart3 className="h-12 w-12 mx-auto mb-4 text-gray-400" />
            <h3 className="text-lg font-semibold mb-2">Aucune données disponibles</h3>
            <p className="text-gray-600">La vue d'ensemble n'est pas encore disponible pour cette analyse.</p>
          </div>
        </div>
      </AnalysisLayout>
    );
  }

  return (
    <AnalysisLayout>
      {/* Debug Info - Development only */}
      {process.env.NODE_ENV === 'development' && (
        <div className="mb-4 p-4 bg-gray-100 rounded">
          <p>Overview Data Status: {overviewData ? 'Loaded' : 'Missing'}</p>
          <p>Overall Score: {overviewData?.score.overall || 0}/100</p>
          <p>Grade: {overviewData?.score.grade || 'N/A'}</p>
          <p>Total Issues: {overviewData?.issues.totalIssues || 0}</p>
          <p>Critical Issues: {overviewData?.issues.critical.length || 0}</p>
          <p>Analysis ID: {analysisId}</p>
          <p>Organic Traffic: {overviewData?.kpis.organicTraffic.current.toLocaleString() || 0}</p>
        </div>
      )}

      <OverviewSection 
        overviewData={overviewData} 
        analysisId={analysisId} 
      />
    </AnalysisLayout>
  );
}