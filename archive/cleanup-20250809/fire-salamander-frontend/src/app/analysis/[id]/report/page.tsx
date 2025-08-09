"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Download, 
  Share2, 
  BarChart3,
  CheckCircle2,
  AlertTriangle,
  Info,
  TrendingUp,
  Globe,
  Zap,
  Brain,
  Target,
  Clock
} from "lucide-react";
import { EnhancedOverviewSection } from "@/components/analysis/enhanced-overview-section";
import { KeywordAnalysisSection } from "@/components/analysis/keyword-analysis-section";
import { mapBackendToOverview } from "@/lib/mappers/overview-mapper";
import { mapBackendToKeywordAnalysis } from "@/lib/mappers/keyword-mapper";

// Composant AI Intelligence Section
function AIIntelligenceSection({ resultData }: { resultData: any }) {
  const topKeywords = [
    { keyword: "hotel cevennes piscine", volume: "2.1K", difficulty: 35, cpc: 1.2 },
    { keyword: "weekend romantique cevennes", volume: "890", difficulty: 28, cpc: 0 },
    { keyword: "hotel spa cevennes", volume: "1.3K", difficulty: 42, cpc: 0 }
  ];
  
  const contentStrategies = [
    { title: "10 Best Hotels in Cévennes with Pools", roi: "High ROI", keywords: 4 },
    { title: "Ultimate Guide to Romantic Weekends", roi: "Medium ROI", keywords: 3 },
    { title: "Cévennes Spa Experience Guide", roi: "High ROI", keywords: 5 }
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Brain className="h-5 w-5 text-purple-600" />
          <span>AI Intelligence</span>
          <Badge className="bg-purple-100 text-purple-700">Powered by GPT-3.5</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top Keyword Opportunities */}
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <Target className="h-4 w-4 text-red-500" />
              <h4 className="font-semibold">Top Keyword Opportunities</h4>
            </div>
            <div className="space-y-3">
              {topKeywords.map((kw, idx) => (
                <div key={idx} className="p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <Target className="h-3 w-3 text-red-500" />
                    <span className="font-medium text-sm">"{kw.keyword}"</span>
                  </div>
                  <div className="text-xs text-gray-600 mt-1">
                    Volume: {kw.volume}, Difficulty: {kw.difficulty}, CPC: €{kw.cpc}
                  </div>
                </div>
              ))}
            </div>
            <div className="mt-4 p-3 bg-red-50 rounded-lg">
              <p className="text-sm text-gray-700">Competitor Gaps: <strong>45 keywords</strong></p>
              <div className="flex space-x-2 mt-2">
                <Button size="sm" variant="outline">Generate More Ideas</Button>
                <Button size="sm" variant="outline">Export</Button>
              </div>
            </div>
          </div>

          {/* Content Strategies */}
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <TrendingUp className="h-4 w-4 text-orange-500" />
              <h4 className="font-semibold">Content Strategies</h4>
            </div>
            <div className="space-y-3">
              {contentStrategies.map((strategy, idx) => (
                <div key={idx} className="p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <TrendingUp className="h-3 w-3 text-orange-500" />
                    <span className="font-medium text-sm">"{strategy.title}"</span>
                  </div>
                  <div className="text-xs text-gray-600 mt-1">
                    {strategy.roi}, targets {strategy.keywords} keywords
                  </div>
                </div>
              ))}
            </div>
            <div className="mt-4 p-3 bg-orange-50 rounded-lg">
              <p className="text-sm text-gray-700">Missing Topics: <strong>12</strong></p>
              <div className="flex space-x-2 mt-2">
                <Button size="sm" className="bg-orange-500 hover:bg-orange-600">Create Content</Button>
                <Button size="sm" variant="outline">View All</Button>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Core Web Vitals
function CoreWebVitalsSection({ resultData }: { resultData: any }) {
  const coreWebVitals = resultData.seo_analysis?.performance_metrics?.core_web_vitals || {};
  
  const getScoreColor = (score: string) => {
    switch (score) {
      case 'good': return 'text-green-600 bg-green-50';
      case 'needs-improvement': return 'text-yellow-600 bg-yellow-50';
      case 'poor': return 'text-red-600 bg-red-50';
      default: return 'text-gray-600 bg-gray-50';
    }
  };
  
  const getScoreStatus = (score: string) => {
    switch (score) {
      case 'good': return 'Good';
      case 'needs-improvement': return 'Needs Improvement';
      case 'poor': return 'Poor';
      default: return 'Unknown';
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Core Web Vitals & Performance</CardTitle>
          <div className="flex space-x-2">
            <Badge variant="outline">Desktop</Badge>
            <Badge variant="secondary">Mobile</Badge>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          {/* Largest Contentful Paint */}
          <div className="text-center">
            <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mb-2 ${getScoreColor(coreWebVitals.lcp?.score || 'good')}`}>
              {getScoreStatus(coreWebVitals.lcp?.score || 'good')}
            </div>
            <div className="text-2xl font-bold text-gray-900">
              {coreWebVitals.lcp?.value || '1.2'}s
            </div>
            <div className="text-sm text-gray-600">Largest Contentful Paint (LCP)</div>
            <div className="text-xs text-gray-500 mt-1">
              Target: &lt; 2.5s
            </div>
          </div>

          {/* First Input Delay */}
          <div className="text-center">
            <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mb-2 ${getScoreColor(coreWebVitals.fid?.score || 'good')}`}>
              {getScoreStatus(coreWebVitals.fid?.score || 'good')}
            </div>
            <div className="text-2xl font-bold text-gray-900">
              {coreWebVitals.fid?.value || '19'}ms
            </div>
            <div className="text-sm text-gray-600">First Input Delay (FID)</div>
            <div className="text-xs text-gray-500 mt-1">
              Target: &lt; 100ms
            </div>
          </div>

          {/* Cumulative Layout Shift */}
          <div className="text-center">
            <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mb-2 ${getScoreColor(coreWebVitals.cls?.score || 'needs-improvement')}`}>
              {getScoreStatus(coreWebVitals.cls?.score || 'needs-improvement')}
            </div>
            <div className="text-2xl font-bold text-gray-900">
              {coreWebVitals.cls?.value || '0.17'}
            </div>
            <div className="text-sm text-gray-600">Cumulative Layout Shift (CLS)</div>
            <div className="text-xs text-gray-500 mt-1">
              Target: &lt; 0.1
            </div>
          </div>
        </div>

        {/* Performance Trends */}
        <div>
          <h4 className="font-semibold mb-4">Performance Trends</h4>
          <div className="h-32 bg-gray-100 rounded-lg flex items-center justify-center">
            <p className="text-gray-500">Performance chart placeholder - LCP, FID, CLS trends over time</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Score Card circulaire (même que dashboard)
function ScoreCard({ title, score, maxScore = 100, color = "blue" }: {
  title: string;
  score: number;
  maxScore?: number;
  color?: "blue" | "green" | "orange" | "red";
}) {
  const percentage = (score / maxScore) * 100;
  const colorClasses = {
    blue: "text-blue-600 border-blue-200 bg-blue-50",
    green: "text-green-600 border-green-200 bg-green-50", 
    orange: "text-orange-600 border-orange-200 bg-orange-50",
    red: "text-red-600 border-red-200 bg-red-50"
  };

  return (
    <Card className={`${colorClasses[color]}`}>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">{title}</p>
            <p className="text-3xl font-bold">
              <span data-testid={title === "Score Global" ? "overall-score" : undefined}>{score}</span>
              <span className="text-lg">/{maxScore}</span>
            </p>
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

// Composant Recommandation
function RecommendationCard({ recommendation }: { recommendation: any }) {
  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case "high": return "bg-red-100 text-red-700 border-red-200";
      case "medium": return "bg-yellow-100 text-yellow-700 border-yellow-200";
      case "low": return "bg-green-100 text-green-700 border-green-200";
      default: return "bg-gray-100 text-gray-700 border-gray-200";
    }
  };

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className={`p-2 rounded-lg ${getPriorityColor(recommendation.priority)}`}>
              {recommendation.priority === "high" && <AlertTriangle className="h-4 w-4" />}
              {recommendation.priority === "medium" && <Info className="h-4 w-4" />}
              {recommendation.priority === "low" && <CheckCircle2 className="h-4 w-4" />}
            </div>
            <div>
              <h4 className="font-semibold">{recommendation.title}</h4>
              <p className="text-sm text-gray-600">{recommendation.category}</p>
            </div>
          </div>
          <Badge className={getPriorityColor(recommendation.priority)}>
            {recommendation.priority}
          </Badge>
        </div>
        
        <p className="text-gray-700 mb-4">{recommendation.description}</p>
        
        {recommendation.action && (
          <div className="p-3 bg-gray-50 rounded-lg">
            <p className="text-sm font-medium text-gray-900 mb-1">Action recommandée:</p>
            <p className="text-sm text-gray-700">{recommendation.action}</p>
          </div>
        )}
        
        <div className="flex items-center justify-between mt-4 text-sm text-gray-500">
          <span>Impact: {recommendation.impact}</span>
          <span>{recommendation.pages} pages affectées</span>
        </div>
      </CardContent>
    </Card>
  );
}

export default function AnalysisReportPage() {
  const params = useParams();
  const [analysis, setAnalysis] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [keywordData, setKeywordData] = useState<any>(null);

  const taskId = params.id as string;

  useEffect(() => {
    const fetchAnalysis = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${taskId}`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
        } else {
          setError("Analyse non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse:", err);
        
        // Données mockées pour la demo
        setAnalysis({
          id: taskId,
          url: "https://www.marina-plage.com/",
          analyzed_at: new Date().toISOString(),
          processing_time: 4500,
          overall_score: 78,
          scores: {
            seo: 82,
            technical: 71,
            performance: 84,
            content: 90
          },
          result_data: JSON.stringify({
            seo_analysis: {
              recommendations: [
                {
                  title: "Optimiser les balises meta description",
                  description: "Plusieurs pages n'ont pas de meta description ou ont des descriptions trop courtes",
                  priority: "high",
                  category: "SEO",
                  impact: "High",
                  pages: 23,
                  action: "Ajouter des meta descriptions uniques de 150-160 caractères pour chaque page"
                },
                {
                  title: "Améliorer la vitesse de chargement",
                  description: "Le temps de chargement initial est supérieur à 3 secondes",
                  priority: "medium",
                  category: "Performance",
                  impact: "Medium",
                  pages: 1,
                  action: "Optimiser les images et utiliser un CDN"
                },
                {
                  title: "Optimiser les images ALT",
                  description: "Des images n'ont pas d'attribut alt ou ont des descriptions peu descriptives",
                  priority: "low",
                  category: "Accessibility",
                  impact: "Low",
                  pages: 12,
                  action: "Ajouter des descriptions alt pertinentes pour toutes les images"
                }
              ]
            }
          })
        });
      } finally {
        setLoading(false);
      }
    };

    fetchAnalysis();
  }, [taskId]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement du rapport...</p>
        </div>
      </div>
    );
  }

  if (error || !analysis) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card className="border-red-200">
          <CardContent className="p-8 text-center">
            <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Rapport non disponible</h2>
            <p className="text-gray-600 mb-4">{error || "Impossible de charger le rapport"}</p>
            <Button onClick={() => window.history.back()}>
              Retour
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const resultData = JSON.parse(analysis.result_data || "{}");
  const recommendations = resultData.seo_analysis?.recommendations || [];
  
  // Extract real backend scores from result_data
  const categoryScores = resultData.category_scores || {};
  const seoScore = Math.round(categoryScores.basics || 0);
  const technicalScore = Math.round(categoryScores.technical || 0);
  const performanceScore = Math.round(categoryScores.performance || 0);
  const overallScore = Math.round((analysis.overall_score || 0) * 100);
  
  // Map to enhanced overview structure
  const overviewData = mapBackendToOverview(analysis);
  
  useEffect(() => {
    if (!analysis) return;
    
    const loadKeywordData = async () => {
      try {
        const data = await mapBackendToKeywordAnalysis(analysis, true);
        setKeywordData(data);
      } catch (error) {
        console.error('Error loading keyword data:', error);
      }
    };
    
    loadKeywordData();
  }, [analysis]);

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Rapport d&apos;Analyse SEO</h1>
          <div className="flex items-center space-x-4 mt-2 text-gray-600">
            <span className="flex items-center space-x-1" data-testid="analyzed-url">
              <Globe className="h-4 w-4" />
              <span>{analysis.url}</span>
            </span>
            <span className="flex items-center space-x-1">
              <Clock className="h-4 w-4" />
              <span>{new Date(analysis.analyzed_at).toLocaleString()}</span>
            </span>
          </div>
        </div>
        
        <div className="flex space-x-2">
          <Button variant="outline">
            <Share2 className="h-4 w-4 mr-2" />
            Partager
          </Button>
          <Button className="bg-orange-500 hover:bg-orange-600">
            <Download className="h-4 w-4 mr-2" />
            Exporter
          </Button>
        </div>
      </div>

      {/* Score Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div data-testid="overall-score-card">
          <ScoreCard 
            title="Score Global" 
            score={overallScore} 
            color={overallScore >= 80 ? "green" : overallScore >= 60 ? "orange" : "red"} 
          />
        </div>
        <ScoreCard title="SEO" score={seoScore} color="blue" />
        <ScoreCard title="Technique" score={technicalScore} color="orange" />
        <ScoreCard title="Performance" score={performanceScore} color="green" />
      </div>

      {/* AI Intelligence Section - Seulement si mode AI Boost (analysis_type === 'full') */}
      {analysis.analysis_type === 'full' && (
        <div className="mb-8">
          <AIIntelligenceSection resultData={resultData} />
        </div>
      )}

      {/* Core Web Vitals & Performance */}
      <div className="mb-8">
        <CoreWebVitalsSection resultData={resultData} />
      </div>

      {/* Tabs pour les différentes sections - 4 onglets */}
      <Tabs defaultValue="overview" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Vue d&apos;ensemble</TabsTrigger>
          <TabsTrigger value="recommendations">Recommandations</TabsTrigger>
          <TabsTrigger value="keywords">Mots-clés</TabsTrigger>
          <TabsTrigger value="seo">SEO</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <EnhancedOverviewSection overview={overviewData} />
        </TabsContent>

        <TabsContent value="recommendations" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-xl font-semibold">Recommandations Actionnables</h3>
            <Badge variant="secondary">
              {recommendations.length} recommandations
            </Badge>
          </div>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {recommendations.map((rec: any, index: number) => (
              <div key={index} data-testid="recommendation-item">
                <RecommendationCard recommendation={rec} />
              </div>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="keywords" className="space-y-6">
          {keywordData ? (
            <KeywordAnalysisSection 
              keywordData={keywordData} 
              analysisId={analysis.id}
            />
          ) : (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500 mx-auto mb-4"></div>
              <p className="text-gray-600">Chargement de l'analyse des mots-clés...</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="seo" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Analyse SEO Technique</CardTitle>
              <CardDescription>
                Détails de l&apos;optimisation pour les moteurs de recherche
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h4 className="font-semibold mb-3">Meta Tags</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                      <span className="text-sm">Title Tag</span>
                      <Badge className="bg-green-100 text-green-700">✓ Présent</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-red-50 rounded-lg">
                      <span className="text-sm">Meta Description</span>
                      <Badge className="bg-red-100 text-red-700">✗ Manquant</Badge>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h4 className="font-semibold mb-3">Structure</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                      <span className="text-sm">H1 Tag</span>
                      <Badge className="bg-green-100 text-green-700">✓ Unique</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                      <span className="text-sm">Structure H2-H6</span>
                      <Badge className="bg-yellow-100 text-yellow-700">⚠ À améliorer</Badge>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

      </Tabs>
    </div>
  );
}