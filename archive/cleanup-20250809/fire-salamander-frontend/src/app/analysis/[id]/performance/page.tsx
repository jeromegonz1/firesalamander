"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { PerformanceAnalysisSection } from "@/components/analysis/performance-analysis-section";
import { mapBackendToPerformanceAnalysis } from "@/lib/mappers/performance-mapper";
import { PerformanceAnalysis, DeviceType, PerformanceGrade, RecommendationType, ImpactLevel } from "@/types/performance-analysis";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Activity, Gauge } from 'lucide-react';

interface AnalysisData {
  id: string;
  url: string;
  analyzed_at: string;
  performance_data?: any;
}

export default function AnalysisPerformancePage() {
  const params = useParams();
  const [performanceData, setPerformanceData] = useState<PerformanceAnalysis | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params.id as string;

  useEffect(() => {
    const fetchPerformanceAnalysis = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}/performance`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
          
          // Map backend data to our PerformanceAnalysis interface
          const mappedData = mapBackendToPerformanceAnalysis(data.data);
          setPerformanceData(mappedData);
        } else {
          setError("Analyse performance non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse performance:", err);
        setError("Erreur de connexion");
        
        // Fallback to mock data for development
        const mockData = createMockPerformanceData();
        setPerformanceData(mockData);
        console.log('Using mock performance data:', mockData);
      } finally {
        setLoading(false);
      }
    };

    fetchPerformanceAnalysis();
  }, [analysisId]);

  const createMockPerformanceData = (): PerformanceAnalysis => {
    return {
      pageMetrics: [
        {
          url: 'https://www.marina-plage.com/',
          testedAt: new Date().toISOString(),
          mobile: {
            lcp: { value: 2800, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 65 },
            fid: { value: 120, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 60 },
            cls: { value: 0.15, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 55 },
            fcp: { value: 1900, grade: PerformanceGrade.GOOD, percentile: 75 },
            tti: { value: 4200, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 58 },
            speedIndex: { value: 3800, grade: PerformanceGrade.NEEDS_IMPROVEMENT, percentile: 62 },
            overallScore: 72,
          },
          desktop: {
            lcp: { value: 1800, grade: PerformanceGrade.EXCELLENT, percentile: 92 },
            fid: { value: 45, grade: PerformanceGrade.EXCELLENT, percentile: 98 },
            cls: { value: 0.08, grade: PerformanceGrade.EXCELLENT, percentile: 95 },
            fcp: { value: 1200, grade: PerformanceGrade.EXCELLENT, percentile: 94 },
            tti: { value: 2800, grade: PerformanceGrade.EXCELLENT, percentile: 88 },
            speedIndex: { value: 2200, grade: PerformanceGrade.EXCELLENT, percentile: 90 },
            overallScore: 94,
          },
          performance: {
            navigationTiming: {
              dns: 25,
              tcp: 78,
              ssl: 145,
              ttfb: 320,
              domLoaded: 1800,
              pageLoaded: 2800,
              interactive: 3200,
            },
            resources: {
              images: {
                count: 18,
                size: 2856000, // 2.86MB
                unoptimized: 12,
                oversized: 5,
                formats: { jpg: 12, png: 4, webp: 2 },
                avgSize: 158667,
                largestSize: 650000,
                withoutAlt: 3,
              },
              css: {
                count: 6,
                size: 245000,
                external: 4,
                inline: 2,
                blocking: 3,
                unused: 35,
                minified: false,
                avgSize: 40833,
              },
              js: {
                count: 12,
                size: 1200000,
                external: 8,
                inline: 4,
                blocking: 2,
                async: 5,
                defer: 5,
                unused: 40,
                minified: false,
                avgSize: 100000,
              },
              fonts: {
                count: 4,
                size: 180000,
                formats: { woff2: 2, woff: 1, ttf: 1 },
                preloaded: 0,
                fallbacks: false,
                avgSize: 45000,
              },
              total: {
                count: 40,
                size: 4481000, // 4.48MB
                requests: 40,
                transferSize: 3136700, // 30% compression
                compressionRatio: 30.0,
              },
            },
            optimization: {
              compression: {
                enabled: false,
                algorithm: null,
                ratio: 0,
                supportedTypes: [],
              },
              caching: {
                browser: {
                  'text/css': 'max-age=3600',
                  'application/javascript': 'max-age=3600',
                  'image/*': 'max-age=86400',
                },
                cdn: false,
                serviceWorker: false,
                etags: false,
                lastModified: true,
                staticAssets: {
                  images: '1 day',
                  css: '1 hour',
                  js: '1 hour',
                  fonts: '1 week',
                },
              },
              cdn: {
                enabled: false,
                provider: null,
                endpoints: [],
                coverage: 0,
              },
              minification: {
                css: { enabled: false, savings: 0, ratio: 0 },
                js: { enabled: false, savings: 0, ratio: 0 },
                html: { enabled: false, savings: 0, ratio: 0 },
              },
              modernOptimizations: {
                criticalCSS: false,
                preloadKey: false,
                prefetch: false,
                lazyLoading: false,
                codeSplitting: false,
                treeShaking: false,
                http2: false,
                http3: false,
              },
            },
            server: {
              responseTime: 320,
              statusCode: 200,
              redirects: 1,
              redirectChain: ['https://marina-plage.com', 'https://www.marina-plage.com'],
            },
            network: {
              bandwidth: 'fast-3g',
              latency: 150,
              throughput: 1.6,
            },
          },
          scores: {
            performance: 83,
            accessibility: 87,
            bestPractices: 79,
            seo: 91,
            pwa: 42,
          },
        },
      ],
      recommendations: [
        {
          id: 'optimize-images-001',
          type: RecommendationType.IMAGES,
          title: 'Optimiser les images',
          description: '12 images peuvent être optimisées pour réduire leur taille de 65%',
          impact: ImpactLevel.HIGH,
          pages: [{
            url: 'https://www.marina-plage.com/',
            currentValue: 2856000,
            potentialGain: 1856000,
            priority: 9,
          }],
          solution: {
            description: 'Convertir les images au format WebP et appliquer une compression moderne',
            implementation: 'Utiliser des outils comme imagemin, squoosh.app ou un service CDN avec optimisation automatique',
            difficulty: 'medium',
            estimatedTime: '3-5 heures',
            resources: [
              'https://web.dev/serve-images-webp/',
              'https://squoosh.app/',
              'https://imagemin.com/',
            ],
          },
          estimatedGain: {
            loadTime: 'Réduction de 2.1s',
            scoreImprovement: 18,
            bandwidth: 'Économie de 1.86MB',
            userExperience: 'Amélioration significative du LCP et de la perception de vitesse',
          },
          metrics: {
            before: 2856000,
            after: 1000000,
            improvement: 65.0,
          },
          validation: {
            tool: 'PageSpeed Insights',
            tested: false,
            results: '',
          },
        },
        {
          id: 'enable-compression-002',
          type: RecommendationType.COMPRESSION,
          title: 'Activer la compression serveur',
          description: 'La compression Gzip/Brotli n\'est pas activée sur le serveur',
          impact: ImpactLevel.HIGH,
          pages: [{
            url: 'https://www.marina-plage.com/',
            currentValue: 4481000,
            potentialGain: 3136700,
            priority: 8,
          }],
          solution: {
            description: 'Activer la compression Gzip ou Brotli sur le serveur web',
            implementation: 'Configurer nginx, Apache ou le CDN pour compresser automatiquement les ressources texte',
            difficulty: 'easy',
            estimatedTime: '30 minutes - 1 heure',
            resources: [
              'https://web.dev/reduce-network-payloads-using-text-compression/',
              'https://nginx.org/en/docs/http/ngx_http_gzip_module.html',
            ],
          },
          estimatedGain: {
            loadTime: 'Réduction de 1.8s',
            scoreImprovement: 15,
            bandwidth: 'Économie de 3.14MB',
            userExperience: 'Amélioration du temps de téléchargement pour toutes les ressources',
          },
          metrics: {
            before: 4481000,
            after: 1344300,
            improvement: 70.0,
          },
          validation: {
            tool: 'GTmetrix',
            tested: false,
            results: '',
          },
        },
        {
          id: 'minify-resources-003',
          type: RecommendationType.MINIFICATION,
          title: 'Minifier CSS et JavaScript',
          description: 'Les fichiers CSS et JS ne sont pas minifiés, économies potentielles de 445KB',
          impact: ImpactLevel.MEDIUM,
          pages: [{
            url: 'https://www.marina-plage.com/',
            currentValue: 1445000,
            potentialGain: 445000,
            priority: 6,
          }],
          solution: {
            description: 'Minifier les fichiers CSS et JavaScript en supprimant les espaces, commentaires et caractères inutiles',
            implementation: 'Intégrer des outils de minification dans le processus de build (UglifyJS, Terser, cssnano)',
            difficulty: 'medium',
            estimatedTime: '2-3 heures',
            resources: [
              'https://web.dev/reduce-network-payloads-using-text-compression/',
              'https://terser.org/',
              'https://cssnano.co/',
            ],
          },
          estimatedGain: {
            loadTime: 'Réduction de 0.8s',
            scoreImprovement: 8,
            bandwidth: 'Économie de 445KB',
            userExperience: 'Amélioration du FCP et du temps d\'interactivité',
          },
          metrics: {
            before: 1445000,
            after: 1000000,
            improvement: 30.8,
          },
          validation: {
            tool: 'WebPageTest',
            tested: false,
            results: '',
          },
        },
        {
          id: 'lazy-loading-004',
          type: RecommendationType.LAZY_LOADING,
          title: 'Implémenter le lazy loading',
          description: 'Charger les images hors écran uniquement lors du scroll',
          impact: ImpactLevel.MEDIUM,
          pages: [{
            url: 'https://www.marina-plage.com/',
            currentValue: 18,
            potentialGain: 12,
            priority: 5,
          }],
          solution: {
            description: 'Ajouter l\'attribut loading="lazy" aux images et implémenter l\'Intersection Observer API',
            implementation: 'Utiliser la propriété CSS content-visibility et les attributs HTML natifs',
            difficulty: 'easy',
            estimatedTime: '1-2 heures',
            resources: [
              'https://web.dev/browser-level-image-lazy-loading/',
              'https://web.dev/lazy-loading-images/',
            ],
          },
          estimatedGain: {
            loadTime: 'Réduction de 1.2s',
            scoreImprovement: 12,
            bandwidth: 'Économie initiale de 1.5MB',
            userExperience: 'Amélioration du LCP pour le contenu above-the-fold',
          },
          metrics: {
            before: 18,
            after: 6,
            improvement: 66.7,
          },
          validation: {
            tool: 'Lighthouse',
            tested: false,
            results: '',
          },
        },
      ],
      summary: {
        totalPages: 1,
        avgPerformanceScore: 83,
        avgLoadTime: {
          mobile: 2800,
          desktop: 1800,
        },
        scoreDistribution: {
          excellent: 0,
          good: 1,
          needsImprovement: 0,
          poor: 0,
        },
        coreWebVitals: {
          mobile: {
            lcp: { avg: 2800, passing: 65 },
            fid: { avg: 120, passing: 60 },
            cls: { avg: 0.15, passing: 55 },
          },
          desktop: {
            lcp: { avg: 1800, passing: 92 },
            fid: { avg: 45, passing: 98 },
            cls: { avg: 0.08, passing: 95 },
          },
        },
        opportunities: {
          highImpact: 2,
          estimatedGain: {
            loadTime: 4.9,
            score: 33,
            bandwidth: 5000000,
          },
        },
      },
      config: {
        testLocation: 'Paris, France',
        device: DeviceType.MOBILE,
        connection: 'Fast 3G',
        browser: 'Chrome',
        viewport: { width: 375, height: 667 },
        throttling: true,
      },
      metadata: {
        analysisDate: new Date().toISOString(),
        analysisId: analysisId,
        testDuration: 58,
        version: '1.0.0',
        tool: 'Fire Salamander Performance',
      },
    };
  };

  console.log('State check:', { loading, error, hasPerformanceData: !!performanceData });

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Activity className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-500" />
          <p className="text-gray-600">Chargement de l'analyse performance...</p>
        </div>
      </div>
    );
  }

  if (error && !performanceData) {
    return (
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
    );
  }

  if (!performanceData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Gauge className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <h3 className="text-lg font-semibold mb-2">Aucune données performance</h3>
          <p className="text-gray-600">L'analyse performance n'est pas encore disponible pour cette page.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Debug Info - Development only */}
      {process.env.NODE_ENV === 'development' && (
        <div className="mb-4 p-4 bg-gray-100 rounded">
          <p>Performance Data Status: {performanceData ? 'Loaded' : 'Missing'}</p>
          <p>Pages Count: {performanceData?.pageMetrics?.length || 0}</p>
          <p>Performance Score: {performanceData?.summary?.avgPerformanceScore || 0}</p>
          <p>Analysis ID: {analysisId}</p>
          <p>Recommendations: {performanceData?.recommendations?.length || 0}</p>
        </div>
      )}

      {/* Use the PerformanceAnalysisSection component */}
      <PerformanceAnalysisSection 
        performanceData={performanceData} 
        analysisId={analysisId} 
      />
    </div>
  );
}