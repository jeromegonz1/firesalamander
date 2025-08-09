"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { 
  BarChart3, 
  Search, 
  CheckCircle, 
  Clock, 
  TrendingUp,
  Globe,
  Zap,
  Shield,
  ArrowRight,
  Plus,
  Brain,
  Target,
  Settings,
  AlertCircle
} from "lucide-react";

// Interface pour les types de scan selon Image #1
interface ScanType {
  id: string;
  title: string;
  features: string[];
  estimatedTime: string;
  isSelected: boolean;
  hasAI: boolean;
}

// Composant Scan Type Card selon Image #1
function ScanTypeCard({ scanType, onSelect }: { scanType: ScanType; onSelect: (id: string) => void }) {
  const borderClass = scanType.isSelected ? "border-orange-500 bg-orange-50" : "border-gray-200 hover:border-gray-300";
  
  return (
    <Card 
      className={`cursor-pointer transition-all ${borderClass} relative`}
      onClick={() => onSelect(scanType.id)}
      data-testid={`${scanType.id}-scan`}
    >
      {scanType.isSelected && (
        <div className="absolute top-3 right-3" data-testid={`${scanType.id}-selected-icon`}>
          <div className="w-6 h-6 bg-orange-500 rounded-full flex items-center justify-center">
            <CheckCircle className="h-4 w-4 text-white" />
          </div>
        </div>
      )}
      
      <CardContent className="p-6">
        <h3 className="text-lg font-semibold mb-4">{scanType.title}</h3>
        
        <div className="space-y-2 mb-4">
          {scanType.features.map((feature, index) => (
            <div key={index} className="flex items-center space-x-2 text-sm">
              {scanType.hasAI || !feature.includes('AI-powered') ? (
                <CheckCircle className="h-4 w-4 text-green-500" />
              ) : (
                <div className="h-4 w-4 flex items-center justify-center" data-testid="quick-no-ai">
                  <span className="text-gray-400 line-through">✗</span>
                </div>
              )}
              <span className={scanType.hasAI || !feature.includes('AI-powered') ? '' : 'text-gray-400 line-through'}>
                {feature}
              </span>
            </div>
          ))}
        </div>
        
        <p className="text-sm text-gray-600 font-medium">{scanType.estimatedTime}</p>
      </CardContent>
    </Card>
  );
}

// ScoreCard supprimé - remplacé par les blocs de réassurance

// Composant d'analyse récente
function RecentAnalysisItem({ url, score, date, status }: {
  url: string;
  score: number;
  date: string;
  status: "completed" | "running" | "failed";
}) {
  const statusConfig = {
    completed: { color: "green", icon: CheckCircle, text: "Terminé" },
    running: { color: "blue", icon: Clock, text: "En cours" },
    failed: { color: "red", icon: Clock, text: "Échec" }
  };

  const config = statusConfig[status];
  const StatusIcon = config.icon;

  return (
    <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition-colors">
      <div className="flex items-center space-x-3">
        <div className={`p-2 rounded-lg bg-${config.color}-50`}>
          <StatusIcon className={`h-4 w-4 text-${config.color}-600`} />
        </div>
        <div>
          <p className="font-medium">{url}</p>
          <p className="text-sm text-gray-500">{date}</p>
        </div>
      </div>
      <div className="flex items-center space-x-3">
        <Badge variant={status === "completed" ? "default" : "secondary"}>
          Score: {score}/100
        </Badge>
        <Button variant="ghost" size="sm">
          <ArrowRight className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const router = useRouter();
  const [analyses, setAnalyses] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [url, setUrl] = useState('');
  const [urlError, setUrlError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [selectedScanType, setSelectedScanType] = useState('quick');
  const [options, setOptions] = useState({
    advancedOptions: false,
    competitorAnalysis: false,
    aiInsights: true // Coché par défaut selon Image #1
  });

  // Types de scan réorganisés selon demande utilisateur
  const scanTypes: ScanType[] = [
    {
      id: 'quick',
      title: 'Quick Scan',
      features: [
        'Home page analysis',
        'Basic technical checks',
        'Core Web Vitals',
        'Top 5 issues identified',
        'AI-powered recommendations'
      ],
      estimatedTime: 'Estimated time: 1-2 minutes',
      isSelected: selectedScanType === 'quick',
      hasAI: false
    },
    {
      id: 'comprehensive',
      title: 'IA Boost Scan',
      features: [
        'Full website crawl (up to 500 pages)',
        'Technical SEO analysis',
        'Content quality assessment',
        'Backlink profile analysis',
        'AI-powered recommendations'
      ],
      estimatedTime: 'Estimated time: 4-6 minutes',
      isSelected: selectedScanType === 'comprehensive',
      hasAI: true
    },
    {
      id: 'custom',
      title: 'Custom Scan',
      features: [
        'Select specific pages to scan',
        'Choose analysis modules',
        'Set crawl depth and limits',
        'Configure scan frequency',
        'Optional AI insights'
      ],
      estimatedTime: 'Time varies based on settings',
      isSelected: selectedScanType === 'custom',
      hasAI: true
    }
  ];

  // Récupérer les analyses existantes
  useEffect(() => {
    const fetchAnalyses = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/v1/analyses');
        if (response.ok) {
          const data = await response.json();
          setAnalyses(data.data || []);
        }
      } catch (error) {
        console.error('Erreur lors de la récupération des analyses:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAnalyses();
  }, []);

  // Validation URL
  const validateUrl = (inputUrl: string): boolean => {
    if (!inputUrl.trim()) {
      setUrlError('Please enter a URL');
      return false;
    }

    // Pattern basique pour URL
    const urlPattern = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/;
    if (!urlPattern.test(inputUrl)) {
      setUrlError('Please enter a valid URL');
      return false;
    }

    setUrlError('');
    return true;
  };

  // Soumission du formulaire
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateUrl(url)) {
      return;
    }

    setIsSubmitting(true);
    
    try {
      // Ajouter https:// si pas de protocole
      let formattedUrl = url;
      if (!url.startsWith('http://') && !url.startsWith('https://')) {
        formattedUrl = `https://${url}`;
      }

      const response = await fetch('http://localhost:8080/api/v1/analyze/quick', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: formattedUrl }),
      });

      if (response.ok) {
        const data = await response.json();
        const taskId = data.data.task_id;
        
        // Attendre un peu que l'analyse soit créée en DB
        setTimeout(async () => {
          try {
            // Récupérer l'ID numérique de l'analyse via le task_id
            const analysesResponse = await fetch('http://localhost:8080/api/v1/analyses');
            const analysesData = await analysesResponse.json();
            const analysis = analysesData.data.find((a: any) => a.task_id === taskId);
            
            if (analysis) {
              // Rediriger vers le rapport avec l'ID numérique
              router.push(`/analysis/${analysis.id}/report`);
            } else {
              // Fallback: rediriger vers la progression avec task_id
              router.push(`/analysis/${taskId}/progress`);
            }
          } catch (error) {
            router.push(`/analysis/${taskId}/progress`);
          }
        }, 2000);
      } else {
        setUrlError('Failed to start analysis. Please try again.');
      }
    } catch (error) {
      console.error('Erreur lors du lancement de l\'analyse:', error);
      setUrlError('Network error. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Données mockées pour les stats (si analyses existent)
  const stats = {
    totalAnalyses: analyses.length,
    successfulAnalyses: analyses.filter(a => a.status === 'success').length,
    averageScore: analyses.length > 0 ? Math.round(analyses.reduce((acc, a) => acc + (a.overall_score * 100), 0) / analyses.length) : 0,
    averageTime: "2.3s"
  };

  const recentAnalyses = analyses.slice(0, 3).map(analysis => ({
    url: analysis.url.replace(/^https?:\/\//, ''),
    score: Math.round(analysis.overall_score * 100),
    date: new Date(analysis.created_at).toLocaleDateString(),
    status: analysis.status === 'success' ? 'completed' as const : 'failed' as const
  }));

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement du dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Analysis Form selon Image #1 */}
      <div className="text-center py-12" data-testid="analysis-form">
        <div className="max-w-4xl mx-auto">
          {/* Titre principal selon Image #1 */}
          <div className="mb-8">
            <div className="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <Search className="h-8 w-8 text-orange-600" />
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              New SEO Analysis
            </h1>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto" data-testid="analysis-subtitle">
              Enter a website URL to start a comprehensive SEO audit
            </p>
          </div>

          {/* Formulaire URL */}
          <form onSubmit={handleSubmit} className="mb-8">
            <div className="max-w-2xl mx-auto">
              <div className="relative mb-4">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Globe className="h-5 w-5 text-gray-400" />
                </div>
                <Input
                  id="url-input"
                  type="url"
                  placeholder="Enter website URL (e.g. example.com)"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  className="pl-12 py-4 text-lg border-2 focus:border-orange-500"
                  data-testid="url-input"
                />
              </div>
              
              {urlError && (
                <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg" data-testid="url-validation-error">
                  <div className="flex items-center space-x-2">
                    <AlertCircle className="h-4 w-4 text-red-600" />
                    <span className="text-sm text-red-800">{urlError}</span>
                  </div>
                </div>
              )}
              
              <Button
                type="submit"
                size="lg"
                disabled={isSubmitting}
                className="w-full py-4 text-lg bg-orange-500 hover:bg-orange-600 disabled:opacity-50"
                data-testid="analyze-button"
              >
                {isSubmitting ? (
                  <>
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2" />
                    Starting Analysis...
                  </>
                ) : (
                  <>
                    Analyze Now
                    <ArrowRight className="h-5 w-5 ml-2" />
                  </>
                )}
              </Button>
            </div>
          </form>

          {/* Options selon Image #1 */}
          <div className="max-w-2xl mx-auto mb-8">
            <div className="flex flex-wrap justify-center gap-6 text-sm">
              <label className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  id="advanced-options"
                  checked={options.advancedOptions}
                  onChange={(e) => setOptions({...options, advancedOptions: e.target.checked})}
                  className="rounded border-gray-300 text-orange-500 focus:ring-orange-500"
                  data-testid="advanced-options"
                />
                <span>Advanced Options</span>
              </label>
              
              <label className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  id="competitor-analysis"
                  checked={options.competitorAnalysis}
                  onChange={(e) => setOptions({...options, competitorAnalysis: e.target.checked})}
                  className="rounded border-gray-300 text-orange-500 focus:ring-orange-500"
                  data-testid="competitor-analysis"
                />
                <span>Include Competitor Analysis</span>
              </label>
              
              <label className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  id="ai-insights"
                  checked={options.aiInsights}
                  onChange={(e) => setOptions({...options, aiInsights: e.target.checked})}
                  className="rounded border-gray-300 text-orange-500 focus:ring-orange-500"
                  data-testid="ai-insights"
                />
                <span>Enable AI Insights</span>
              </label>
            </div>
          </div>

          {/* Message de sécurité selon Image #1 */}
          <div className="flex items-center justify-center space-x-2 text-sm text-gray-600" data-testid="security-notice">
            <Shield className="h-4 w-4" data-testid="security-icon" />
            <span>All analyses are private and secure</span>
          </div>
        </div>
      </div>

      {/* Types de scan selon Image #1 */}
      <div className="max-w-6xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {scanTypes.map((scanType) => (
            <ScanTypeCard
              key={scanType.id}
              scanType={scanType}
              onSelect={setSelectedScanType}
            />
          ))}
        </div>
      </div>

      {/* Recent Analyses - Affichées seulement si analyses existent */}
      {analyses.length > 0 && (
        <div className="space-y-8">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Analyses Récentes</CardTitle>
                  <CardDescription>
                    Suivez vos dernières analyses SEO
                  </CardDescription>
                </div>
                <Button variant="outline" size="sm">
                  Voir tout
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentAnalyses.map((analysis, index) => (
                  <RecentAnalysisItem key={index} {...analysis} />
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Bloc de Réassurance selon Image #2 */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6" data-testid="reassurance-blocks">
            <Card className="p-6 text-center hover:shadow-lg transition-shadow border-orange-100">
              <div className="mx-auto w-12 h-12 bg-orange-100 rounded-lg flex items-center justify-center mb-4">
                <Brain className="h-6 w-6 text-orange-600" />
              </div>
              <h3 className="text-lg font-semibold mb-2 text-orange-600">AI-Powered Insights</h3>
              <p className="text-gray-600 text-sm mb-3">
                Leverage advanced AI technology to identify opportunities and generate actionable recommendations.
              </p>
              <div className="space-y-1 text-xs text-left">
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>Keyword opportunities</span>
                </div>
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>Content suggestions</span>
                </div>
              </div>
            </Card>

            <Card className="p-6 text-center hover:shadow-lg transition-shadow border-blue-100">
              <div className="mx-auto w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
                <TrendingUp className="h-6 w-6 text-blue-600" />
              </div>
              <h3 className="text-lg font-semibold mb-2 text-blue-600">Core Web Vitals</h3>
              <p className="text-gray-600 text-sm mb-3">
                Analyze performance metrics that impact your search rankings and user experience.
              </p>
              <div className="space-y-1 text-xs text-left">
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>LCP, FID, CLS measurements</span>
                </div>
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>Mobile & desktop analysis</span>
                </div>
              </div>
            </Card>

            <Card className="p-6 text-center hover:shadow-lg transition-shadow border-purple-100">
              <div className="mx-auto w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
                <Target className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="text-lg font-semibold mb-2 text-purple-600">Backlink Analysis</h3>
              <p className="text-gray-600 text-sm mb-3">
                Discover your backlink profile strength, toxic links, and new link building opportunities.
              </p>
              <div className="space-y-1 text-xs text-left">
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>Domain authority metrics</span>
                </div>
                <div className="flex items-center space-x-2">
                  <CheckCircle className="h-3 w-3 text-green-500" />
                  <span>Competitor link analysis</span>
                </div>
              </div>
            </Card>
          </div>
        </div>
      )}

      {/* Feedback region pour accessibility */}
      <div 
        className="sr-only" 
        data-testid="form-feedback" 
        aria-live="polite" 
        aria-atomic="true"
      >
        {urlError && `Error: ${urlError}`}
        {isSubmitting && "Starting analysis..."}
      </div>
    </div>
  );
}