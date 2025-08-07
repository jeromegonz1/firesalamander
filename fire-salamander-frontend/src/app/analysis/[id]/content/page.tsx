"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { ContentAnalysisSection } from "@/components/analysis/content-analysis-section";
import { mapBackendToContentAnalysis } from "@/lib/mappers/content-analysis-mapper";
import { 
  ContentAnalysis, 
  ContentType, 
  ContentQuality, 
  ReadabilityLevel,
  ContentIssueType,
  ContentSeverity,
  TopicImportance
} from "@/types/content-analysis";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Activity, FileText } from 'lucide-react';

interface AnalysisData {
  id: string;
  url: string;
  analyzed_at: string;
  content_analysis?: any;
}

export default function AnalysisContentPage() {
  const params = useParams();
  const [contentData, setContentData] = useState<ContentAnalysis | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params.id as string;

  useEffect(() => {
    const fetchContentAnalysis = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}/content`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
          
          // Map backend data to our ContentAnalysis interface
          const mappedData = mapBackendToContentAnalysis(data.data);
          setContentData(mappedData);
        } else {
          setError("Analyse de contenu non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse de contenu:", err);
        setError("Erreur de connexion");
        
        // Fallback to mock data for development
        const mockData = createMockContentData();
        setContentData(mockData);
        console.log('Using mock content data:', mockData);
      } finally {
        setLoading(false);
      }
    };

    fetchContentAnalysis();
  }, [analysisId]);

  const createMockContentData = (): ContentAnalysis => {
    return {
      pages: [
        {
          url: 'https://www.marina-plage.com/',
          title: 'Marina Plage - Restaurant & Bar de plage à Nice',
          contentType: ContentType.HOMEPAGE,
          publishDate: '2024-01-15',
          lastModified: '2024-03-10',
          
          metrics: {
            wordCount: 1200,
            uniqueWords: 680,
            sentenceCount: 85,
            paragraphCount: 18,
            readingTime: 5,
            avgWordsPerSentence: 14.1,
            avgSentencesPerParagraph: 4.7,
          },
          
          quality: {
            overallScore: 68,
            originalityScore: 85,
            topicRelevance: 72,
            keywordIntegration: 65,
            readabilityScore: 74,
            engagementPotential: 70,
            expertiseLevel: 60,
            freshness: 75,
          },
          
          readability: {
            level: ReadabilityLevel.FAIRLY_EASY,
            fleschKincaidGrade: 7.2,
            fleschReadingEase: 74.5,
            gunningFogIndex: 8.1,
            smogIndex: 7.8,
            colemanLiauIndex: 6.9,
            automatedReadabilityIndex: 7.5,
            averageGradeLevel: 7.4,
          },
          
          structure: {
            headings: { h1: 1, h2: 5, h3: 8, h4: 2, h5: 0, h6: 0 },
            lists: { ordered: 2, unordered: 4, totalItems: 25 },
            images: { total: 15, withAlt: 12, withCaption: 8, decorative: 3 },
            links: { internal: 18, external: 6, outbound: 4, anchor: 2 },
            multimedia: { videos: 1, audios: 0, embeds: 2, infographics: 1 },
          },
          
          keywords: {
            primary: [
              {
                keyword: 'restaurant plage Nice',
                density: 1.8,
                frequency: 12,
                positions: [45, 123, 234, 456, 567, 678, 789, 890, 1001, 1100, 1150, 1200],
                inTitle: true,
                inHeadings: true,
                inMeta: true,
                prominence: 92,
              },
              {
                keyword: 'bar plage',
                density: 1.2,
                frequency: 8,
                positions: [67, 134, 289, 445, 623, 778, 945, 1080],
                inTitle: false,
                inHeadings: true,
                inMeta: false,
                prominence: 75,
              },
            ],
            secondary: [
              {
                keyword: 'Marina Plage',
                density: 2.5,
                frequency: 15,
                relevanceScore: 95,
              },
              {
                keyword: 'Côte d\'Azur',
                density: 0.8,
                frequency: 5,
                relevanceScore: 78,
              },
            ],
            lsi: [
              {
                term: 'cuisine méditerranéenne',
                relevance: 88,
                frequency: 4,
              },
              {
                term: 'terrasse vue mer',
                relevance: 82,
                frequency: 3,
              },
            ],
            entities: [
              {
                entity: 'Nice',
                type: 'place',
                confidence: 96,
                mentions: 8,
              },
              {
                entity: 'Méditerranée',
                type: 'place',
                confidence: 89,
                mentions: 3,
              },
            ],
          },
          
          issues: [
            {
              type: ContentIssueType.MISSING_KEYWORDS,
              severity: ContentSeverity.WARNING,
              description: 'Mots-clés importants manquants: "réservation restaurant", "événements privés"',
              affectedElements: ['title', 'headings', 'content'],
              impact: 'Perte potentielle de 15-20% du trafic ciblé',
              howToFix: {
                steps: [
                  'Identifier les mots-clés manquants avec un fort potentiel',
                  'Intégrer naturellement dans le contenu existant',
                  'Optimiser les balises title et meta description',
                  'Créer des sections dédiées pour ces sujets',
                ],
                estimatedTime: '2-3 heures',
                difficulty: 'medium',
                codeExample: `<!-- Exemple d'intégration -->
<h2>Réservation Restaurant Marina Plage</h2>
<p>Réservez votre table dans notre restaurant de plage à Nice...</p>

<!-- Structured data pour événements -->
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Restaurant",
  "name": "Marina Plage",
  "hasMenu": "https://marina-plage.com/menu",
  "acceptsReservations": true
}
</script>`,
                resources: [
                  {
                    title: 'Guide intégration mots-clés',
                    url: 'https://moz.com/learn/seo/on-page-factors',
                    type: 'guide',
                  },
                ],
                beforeAfter: {
                  before: 'Contenu sans mention des réservations et événements privés',
                  after: 'Contenu enrichi avec sections dédiées réservation et événements, mots-clés intégrés naturellement',
                },
                priority: 7,
                potentialGain: {
                  traffic: '+18-25% trafic organique',
                  performance: '+8 points qualité',
                  ranking: 'Amélioration 3-5 positions',
                  engagement: '+12% temps sur page',
                },
              },
            },
            {
              type: ContentIssueType.POOR_STRUCTURE,
              severity: ContentSeverity.INFO,
              description: 'Structure H2/H3 peut être améliorée pour un meilleur balisage sémantique',
              affectedElements: ['headings'],
              impact: 'Impact SEO modéré sur la compréhension du contenu',
              howToFix: {
                steps: [
                  'Restructurer la hiérarchie des titres H2/H3',
                  'Créer une progression logique des sujets',
                  'Ajouter des mots-clés dans les titres',
                ],
                estimatedTime: '1-2 heures',
                difficulty: 'easy',
                priority: 4,
                potentialGain: {
                  traffic: '+5-8% trafic',
                  performance: '+3 points',
                  ranking: 'Amélioration 1-2 positions',
                  engagement: '+5% engagement',
                },
              },
            },
          ],
          
          suggestions: {
            contentLength: {
              current: 1200,
              recommended: 1800,
              reason: 'Analyse concurrentielle: pages similaires font 1800+ mots en moyenne',
            },
            missingTopics: [
              'menu spécialités',
              'événements privés',
              'horaires saisons',
              'réservation en ligne',
            ],
            keywordsToAdd: [
              {
                keyword: 'réservation restaurant Nice',
                searchVolume: 2400,
                difficulty: 45,
                opportunity: 82,
              },
              {
                keyword: 'événements privés plage',
                searchVolume: 890,
                difficulty: 35,
                opportunity: 75,
              },
            ],
            structureImprovements: [
              'Ajouter section FAQ',
              'Créer section témoignages clients',
              'Améliorer le maillage interne',
            ],
            readabilityEnhancements: [
              'Raccourcir certains paragraphes',
              'Ajouter plus de listes à puces',
              'Utiliser plus de sous-titres',
            ],
          },
        },
        // Deuxième page d'exemple
        {
          url: 'https://www.marina-plage.com/menu',
          title: 'Menu Restaurant - Spécialités Méditerranéennes',
          contentType: ContentType.ARTICLE,
          publishDate: '2024-02-01',
          lastModified: '2024-03-05',
          
          metrics: {
            wordCount: 800,
            uniqueWords: 450,
            sentenceCount: 60,
            paragraphCount: 15,
            readingTime: 3,
            avgWordsPerSentence: 13.3,
            avgSentencesPerParagraph: 4.0,
          },
          
          quality: {
            overallScore: 45,
            originalityScore: 55,
            topicRelevance: 85,
            keywordIntegration: 40,
            readabilityScore: 78,
            engagementPotential: 35,
            expertiseLevel: 50,
            freshness: 80,
          },
          
          readability: {
            level: ReadabilityLevel.FAIRLY_EASY,
            fleschKincaidGrade: 6.8,
            fleschReadingEase: 76.2,
            gunningFogIndex: 7.5,
            smogIndex: 7.2,
            colemanLiauIndex: 6.5,
            automatedReadabilityIndex: 7.0,
            averageGradeLevel: 6.8,
          },
          
          structure: {
            headings: { h1: 1, h2: 3, h3: 5, h4: 0, h5: 0, h6: 0 },
            lists: { ordered: 0, unordered: 8, totalItems: 45 },
            images: { total: 8, withAlt: 5, withCaption: 3, decorative: 2 },
            links: { internal: 5, external: 2, outbound: 1, anchor: 1 },
            multimedia: { videos: 0, audios: 0, embeds: 0, infographics: 0 },
          },
          
          keywords: {
            primary: [
              {
                keyword: 'menu restaurant',
                density: 1.5,
                frequency: 6,
                positions: [25, 145, 278, 456, 623, 750],
                inTitle: true,
                inHeadings: true,
                inMeta: true,
                prominence: 88,
              },
            ],
            secondary: [
              {
                keyword: 'spécialités méditerranéennes',
                density: 0.8,
                frequency: 4,
                relevanceScore: 92,
              },
            ],
            lsi: [
              {
                term: 'poissons frais',
                relevance: 85,
                frequency: 3,
              },
            ],
            entities: [
              {
                entity: 'Méditerranée',
                type: 'place',
                confidence: 92,
                mentions: 4,
              },
            ],
          },
          
          issues: [
            {
              type: ContentIssueType.THIN_CONTENT,
              severity: ContentSeverity.CRITICAL,
              description: 'Contenu insuffisant: seulement 800 mots pour une page menu importante',
              affectedElements: ['content'],
              impact: 'Faible performance SEO, manque d\'information pour les utilisateurs',
              howToFix: {
                steps: [
                  'Enrichir les descriptions de plats',
                  'Ajouter des informations sur les ingrédients',
                  'Inclure des histoires sur l\'origine des recettes',
                  'Ajouter section allergènes et informations nutritionnelles',
                  'Créer des recommandations de vin',
                ],
                estimatedTime: '4-6 heures',
                difficulty: 'medium',
                priority: 9,
                potentialGain: {
                  traffic: '+35-50% trafic',
                  performance: '+25 points qualité',
                  ranking: 'Amélioration 5-8 positions',
                  engagement: '+40% temps sur page',
                },
              },
            },
            {
              type: ContentIssueType.MISSING_IMAGES,
              severity: ContentSeverity.WARNING,
              description: 'Manque d\'images pour illustrer les plats du menu',
              affectedElements: ['images'],
              impact: 'Engagement utilisateur réduit, moins d\'attrait visuel',
              howToFix: {
                steps: [
                  'Photographier les plats signature',
                  'Optimiser les images pour le web',
                  'Ajouter des descriptions alt appropriées',
                  'Implémenter lazy loading',
                ],
                estimatedTime: '3-4 heures',
                difficulty: 'easy',
                priority: 6,
                potentialGain: {
                  traffic: '+10-15% trafic',
                  performance: '+5 points',
                  ranking: 'Amélioration 1-3 positions',
                  engagement: '+25% engagement',
                },
              },
            },
          ],
          
          suggestions: {
            contentLength: {
              current: 800,
              recommended: 1500,
              reason: 'Pages menu performantes contiennent 1500-2000 mots avec descriptions détaillées',
            },
            missingTopics: [
              'allergènes et intolérances',
              'origines des produits',
              'accords mets-vins',
              'plats végétariens/vegans',
            ],
            keywordsToAdd: [
              {
                keyword: 'menu méditerranéen Nice',
                searchVolume: 1200,
                difficulty: 40,
                opportunity: 78,
              },
            ],
            structureImprovements: [
              'Organiser par catégories (entrées, plats, desserts)',
              'Ajouter prix et descriptions détaillées',
              'Inclure photos de plats',
            ],
            readabilityEnhancements: [
              'Format liste pour les plats',
              'Descriptions courtes et engageantes',
              'Mise en valeur des spécialités',
            ],
          },
        },
      ],
      
      aiContentAnalysis: {
        topicsCovered: [
          {
            topic: 'Restaurant de plage',
            coverage: 85,
            depth: 70,
            authority: 75,
            searchIntent: 'commercial',
          },
          {
            topic: 'Cuisine méditerranéenne',
            coverage: 75,
            depth: 65,
            authority: 80,
            searchIntent: 'informational',
          },
          {
            topic: 'Localisation Nice',
            coverage: 90,
            depth: 85,
            authority: 90,
            searchIntent: 'navigational',
          },
        ],
        
        topicsMissing: [
          {
            topic: 'Réservation en ligne',
            importance: TopicImportance.HIGH,
            searchVolume: 2400,
            difficulty: 45,
            opportunity: 85,
            relatedKeywords: ['réserver table', 'booking', 'disponibilités'],
          },
          {
            topic: 'Événements privés',
            importance: TopicImportance.HIGH,
            searchVolume: 890,
            difficulty: 35,
            opportunity: 78,
            relatedKeywords: ['soirée privée', 'mariage plage', 'événement corporate'],
          },
          {
            topic: 'Livraison à domicile',
            importance: TopicImportance.MEDIUM,
            searchVolume: 1500,
            difficulty: 55,
            opportunity: 65,
            relatedKeywords: ['delivery', 'commande en ligne', 'uber eats'],
          },
        ],
        
        contentGaps: [
          {
            gap: 'Informations pratiques insuffisantes',
            category: 'depth',
            importance: TopicImportance.HIGH,
            suggestedContent: 'Ajouter horaires détaillés, informations parking, accès transports',
            estimatedLength: 300,
            competitorsCovering: 8,
            potentialTraffic: 1200,
          },
          {
            gap: 'Absence de contenu sur l\'histoire du restaurant',
            category: 'knowledge',
            importance: TopicImportance.MEDIUM,
            suggestedContent: 'Raconter l\'histoire du restaurant, présenter l\'équipe, valeurs',
            estimatedLength: 500,
            competitorsCovering: 5,
            potentialTraffic: 800,
          },
        ],
        
        competitorComparison: {
          avgContentLength: 2200,
          ourAvgLength: 1000,
          lengthRecommendation: 'Augmenter la longueur de contenu de 120% pour égaler la concurrence',
          topicCoverageGap: 35,
          qualityGap: 22,
          topCompetitors: [
            {
              url: 'https://restaurant-concurrent-1.com',
              domain: 'restaurant-concurrent-1.com',
              contentLength: 2800,
              qualityScore: 88,
              topicsCovered: ['restaurant plage', 'réservation', 'événements', 'menu détaillé', 'histoire'],
              strengthsOverUs: [
                'Contenu plus approfondi',
                'Système de réservation intégré',
                'Galerie photo professionnelle',
                'Témoignages clients nombreux',
              ],
            },
            {
              url: 'https://restaurant-concurrent-2.com',
              domain: 'restaurant-concurrent-2.com',
              contentLength: 2100,
              qualityScore: 85,
              topicsCovered: ['cuisine méditerranéenne', 'chef', 'produits locaux', 'wine list'],
              strengthsOverUs: [
                'Focus sur le chef et l\'équipe',
                'Mise en avant des producteurs locaux',
                'Carte des vins détaillée',
              ],
            },
          ],
        },
        
        aiRecommendations: [
          {
            type: 'content_expansion',
            priority: 1,
            impact: 'high',
            description: 'Développer le contenu sur les spécialités et l\'identité culinaire',
            implementation: 'Ajouter 800-1000 mots sur les plats signature, l\'histoire des recettes, et la philosophie culinaire',
            expectedResults: {
              trafficIncrease: '30-45%',
              rankingImprovement: '4-7 positions',
              engagementBoost: '25-35%',
            },
          },
          {
            type: 'restructuring',
            priority: 2,
            impact: 'high',
            description: 'Restructurer l\'information avec des sections claires',
            implementation: 'Organiser en sections: Accueil, Menu, Réservation, Événements, À propos',
            expectedResults: {
              trafficIncrease: '15-25%',
              rankingImprovement: '2-4 positions',
              engagementBoost: '20-30%',
            },
          },
          {
            type: 'keyword_optimization',
            priority: 3,
            impact: 'medium',
            description: 'Optimiser pour les mots-clés manqués à fort potentiel',
            implementation: 'Intégrer naturellement "réservation", "événements privés", "delivery" dans le contenu',
            expectedResults: {
              trafficIncrease: '20-30%',
              rankingImprovement: '3-5 positions',
              engagementBoost: '10-15%',
            },
          },
        ],
        
        contentTone: {
          sentiment: 'positive',
          formality: 45,
          complexity: 35,
          enthusiasm: 75,
          trustworthiness: 68,
          expertise: 55,
          targetAudience: ['Touristes', 'Locaux', 'Familles', 'Couples'],
          toneConsistency: 72,
        },
      },
      
      visualizationData: {
        charts: {
          qualityDistribution: [
            { range: '0-20', count: 0, percentage: 0 },
            { range: '21-40', count: 0, percentage: 0 },
            { range: '41-60', count: 1, percentage: 50 },
            { range: '61-80', count: 1, percentage: 50 },
            { range: '81-100', count: 0, percentage: 0 },
          ],
          
          contentPerformanceTimeline: [
            {
              date: '2024-01-01',
              avgQualityScore: 52,
              totalWords: 1800,
              newContent: 1,
              updatedContent: 0,
            },
            {
              date: '2024-02-01',
              avgQualityScore: 55,
              totalWords: 2000,
              newContent: 1,
              updatedContent: 1,
            },
            {
              date: '2024-03-01',
              avgQualityScore: 57,
              totalWords: 2000,
              newContent: 0,
              updatedContent: 2,
            },
          ],
          
          contentTypeDistribution: [
            {
              type: ContentType.HOMEPAGE,
              count: 1,
              avgQuality: 68,
              avgLength: 1200,
            },
            {
              type: ContentType.ARTICLE,
              count: 1,
              avgQuality: 45,
              avgLength: 800,
            },
          ],
          
          lengthVsPerformance: [
            {
              lengthRange: '0-500',
              avgQualityScore: 0,
              avgEngagement: 0,
              count: 0,
            },
            {
              lengthRange: '501-1000',
              avgQualityScore: 45,
              avgEngagement: 2.1,
              count: 1,
            },
            {
              lengthRange: '1001-1500',
              avgQualityScore: 68,
              avgEngagement: 3.2,
              count: 1,
            },
          ],
        },
        
        heatmaps: {
          keywordDensityMap: [
            {
              page: '/',
              keyword: 'restaurant plage Nice',
              density: 1.8,
              position: { x: 50, y: 120 },
              performance: 85,
            },
            {
              page: '/menu',
              keyword: 'menu restaurant',
              density: 1.5,
              position: { x: 30, y: 80 },
              performance: 75,
            },
          ],
          
          topicCoverageMap: [
            {
              topic: 'Restaurant de plage',
              pages: [
                {
                  url: '/',
                  coverage: 85,
                  depth: 70,
                },
              ],
            },
            {
              topic: 'Cuisine méditerranéenne',
              pages: [
                {
                  url: '/',
                  coverage: 75,
                  depth: 65,
                },
                {
                  url: '/menu',
                  coverage: 90,
                  depth: 80,
                },
              ],
            },
          ],
        },
        
        networkGraphs: {
          internalLinkingGraph: {
            nodes: [
              {
                id: 'homepage',
                url: '/',
                title: 'Marina Plage - Restaurant & Bar de plage à Nice',
                contentType: ContentType.HOMEPAGE,
                qualityScore: 68,
                inLinks: 0,
                outLinks: 5,
                pageRank: 0.8,
              },
              {
                id: 'menu',
                url: '/menu',
                title: 'Menu Restaurant - Spécialités Méditerranéennes',
                contentType: ContentType.ARTICLE,
                qualityScore: 45,
                inLinks: 2,
                outLinks: 2,
                pageRank: 0.4,
              },
            ],
            edges: [
              {
                source: 'homepage',
                target: 'menu',
                anchorText: 'découvrir notre menu',
                strength: 0.8,
                relevance: 0.95,
              },
            ],
          },
          
          semanticGraph: {
            nodes: [
              {
                id: 'restaurant-1',
                term: 'restaurant',
                type: 'keyword',
                importance: 95,
                frequency: 15,
              },
              {
                id: 'plage-1',
                term: 'plage',
                type: 'keyword',
                importance: 90,
                frequency: 12,
              },
              {
                id: 'nice-1',
                term: 'Nice',
                type: 'entity',
                importance: 88,
                frequency: 8,
              },
            ],
            edges: [
              {
                source: 'restaurant-1',
                target: 'plage-1',
                relation: 'related',
                strength: 0.9,
              },
              {
                source: 'plage-1',
                target: 'nice-1',
                relation: 'related',
                strength: 0.85,
              },
            ],
          },
        },
      },
      
      globalMetrics: {
        totalPages: 2,
        totalWords: 2000,
        avgWordsPerPage: 1000,
        avgQualityScore: 57,
        uniqueTopicsCovered: 8,
        contentFreshness: 78,
        duplicateContentPages: 0,
        thinContentPages: 1,
        highQualityPages: 0,
        
        competitorBenchmark: {
          position: 'below',
          gapAnalysis: {
            contentVolume: -55,
            contentQuality: -28,
            topicCoverage: -35,
            keywordCoverage: -42,
          },
        },
      },
      
      globalIssues: [
        {
          id: 'global-thin-content',
          type: ContentIssueType.THIN_CONTENT,
          severity: ContentSeverity.CRITICAL,
          affectedPages: ['/menu'],
          impact: 'Pages avec contenu insuffisant perdent du trafic potentiel',
          description: '1 page sur 2 a un contenu insuffisant (moins de 1000 mots)',
          howToFix: {
            steps: [
              'Identifier toutes les pages avec moins de 1000 mots',
              'Analyser les besoins des utilisateurs pour chaque page',
              'Enrichir le contenu avec des informations pertinentes',
              'Ajouter des éléments visuels et multimédia',
            ],
            estimatedTime: '1-2 jours',
            difficulty: 'medium',
            priority: 9,
            potentialGain: {
              traffic: '+40-60% trafic global',
              performance: '+20 points qualité moyenne',
              ranking: 'Amélioration 3-6 positions moyenne',
              engagement: '+30% engagement moyen',
            },
          },
          priority: 9,
          estimatedROI: {
            trafficGain: 50,
            rankingImprovement: 4,
            implementationCost: 'medium',
          },
        },
        {
          id: 'global-missing-keywords',
          type: ContentIssueType.MISSING_KEYWORDS,
          severity: ContentSeverity.WARNING,
          affectedPages: ['/', '/menu'],
          impact: 'Opportunités de trafic manquées sur des mots-clés importants',
          description: 'Plusieurs mots-clés à fort potentiel ne sont pas exploités',
          howToFix: {
            steps: [
              'Audit complet des mots-clés manqués',
              'Prioriser selon volume de recherche et difficulté',
              'Intégrer dans le contenu existant',
              'Créer nouveau contenu si nécessaire',
            ],
            estimatedTime: '3-5 heures',
            difficulty: 'easy',
            priority: 7,
            potentialGain: {
              traffic: '+25-35% trafic',
              performance: '+10 points',
              ranking: 'Amélioration 2-4 positions',
              engagement: '+15% engagement',
            },
          },
          priority: 7,
          estimatedROI: {
            trafficGain: 30,
            rankingImprovement: 3,
            implementationCost: 'low',
          },
        },
      ],
      
      strategicRecommendations: [
        {
          id: 'content-strategy-expansion',
          category: 'content_strategy',
          title: 'Stratégie d\'expansion du contenu',
          description: 'Développer une stratégie de contenu complète pour couvrir tous les aspects de l\'activité',
          implementation: {
            phases: [
              {
                phase: 'Audit et planification',
                duration: '1 semaine',
                tasks: [
                  'Analyser la concurrence en détail',
                  'Identifier tous les sujets pertinents',
                  'Créer un calendrier éditorial',
                ],
                resources: ['SEO manager', 'Content manager'],
              },
              {
                phase: 'Production de contenu',
                duration: '4 semaines',
                tasks: [
                  'Rédiger du contenu enrichi pour toutes les pages',
                  'Créer des nouvelles sections (réservation, événements)',
                  'Produire du contenu visuel (photos, vidéos)',
                ],
                resources: ['Rédacteur', 'Photographe', 'Web developer'],
              },
              {
                phase: 'Optimisation et suivi',
                duration: '2 semaines',
                tasks: [
                  'Optimiser le SEO on-page',
                  'Mettre en place le tracking',
                  'Analyser les premiers résultats',
                ],
                resources: ['SEO specialist', 'Analytics manager'],
              },
            ],
            totalEffort: '7 semaines',
            priority: 1,
          },
          expectedResults: {
            shortTerm: 'Contenu de qualité et pages optimisées (1-3 mois)',
            mediumTerm: 'Amélioration du trafic organique de 40-60% (3-6 mois)',
            longTerm: 'Position de leader sur les requêtes locales (6+ mois)',
          },
        },
      ],
      
      config: {
        analysisDepth: 'comprehensive',
        includeCompetitorAnalysis: true,
        targetKeywords: ['restaurant plage Nice', 'bar plage', 'cuisine méditerranéenne'],
        competitorDomains: ['restaurant-concurrent-1.com', 'restaurant-concurrent-2.com'],
        contentLanguage: 'fr',
        targetAudience: ['Touristes', 'Locaux', 'Familles'],
      },
      
      metadata: {
        analysisDate: new Date().toISOString(),
        analysisId: analysisId,
        processingTime: 145,
        aiModelUsed: 'GPT-4-Turbo',
        confidenceScore: 88,
        version: '1.0.0',
        tool: 'Fire Salamander Content Analysis',
      },
    };
  };

  console.log('State check:', { loading, error, hasContentData: !!contentData });

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Activity className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-500" />
          <p className="text-gray-600">Chargement de l'analyse de contenu...</p>
        </div>
      </div>
    );
  }

  if (error && !contentData) {
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

  if (!contentData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <FileText className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <h3 className="text-lg font-semibold mb-2">Aucune données de contenu</h3>
          <p className="text-gray-600">L'analyse de contenu n'est pas encore disponible pour cette page.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Debug Info - Development only */}
      {process.env.NODE_ENV === 'development' && (
        <div className="mb-4 p-4 bg-gray-100 rounded">
          <p>Content Data Status: {contentData ? 'Loaded' : 'Missing'}</p>
          <p>Pages Count: {contentData?.pages?.length || 0}</p>
          <p>Average Quality Score: {contentData?.globalMetrics?.avgQualityScore || 0}</p>
          <p>Total Words: {contentData?.globalMetrics?.totalWords || 0}</p>
          <p>Analysis ID: {analysisId}</p>
          <p>AI Recommendations: {contentData?.aiContentAnalysis?.aiRecommendations?.length || 0}</p>
        </div>
      )}

      {/* Use the ContentAnalysisSection component */}
      <ContentAnalysisSection 
        contentData={contentData} 
        analysisId={analysisId} 
      />
    </div>
  );
}