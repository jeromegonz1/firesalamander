"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Search, 
  Zap, 
  Settings, 
  CheckCircle2,
  Globe,
  BarChart3,
  Shield,
  Brain,
  Clock,
  Target
} from "lucide-react";

// Types d'analyse disponibles
const analysisTypes = [
  {
    id: "comprehensive",
    name: "Analyse Complète",
    description: "Audit SEO complet avec toutes les fonctionnalités",
    icon: BarChart3,
    duration: "3-5 min",
    features: ["SEO Technique", "Performance", "Contenu", "Crawling", "IA Insights"],
    color: "blue",
    recommended: true
  },
  {
    id: "quick",
    name: "Analyse Rapide",
    description: "Scan SEO rapide pour un aperçu général",
    icon: Zap,
    duration: "30s",
    features: ["SEO Base", "Performance Core", "Meta Tags"],
    color: "green"
  },
  {
    id: "custom",
    name: "Analyse Personnalisée",
    description: "Configurez votre analyse selon vos besoins",
    icon: Settings,
    duration: "Variable",
    features: ["Options personnalisables"],
    color: "purple"
  }
];

function AnalysisTypeCard({ type, selected, onSelect }: {
  type: typeof analysisTypes[0];
  selected: boolean;
  onSelect: () => void;
}) {
  const Icon = type.icon;
  const colorClasses = {
    blue: selected ? "border-blue-500 bg-blue-50" : "border-gray-200 hover:border-blue-300",
    green: selected ? "border-green-500 bg-green-50" : "border-gray-200 hover:border-green-300",
    purple: selected ? "border-purple-500 bg-purple-50" : "border-gray-200 hover:border-purple-300"
  };

  return (
    <Card 
      className={`cursor-pointer transition-all duration-200 ${colorClasses[type.color as keyof typeof colorClasses]}`}
      onClick={onSelect}
    >
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className={`p-3 rounded-lg ${type.color === 'blue' ? 'bg-blue-100' : type.color === 'green' ? 'bg-green-100' : 'bg-purple-100'}`}>
            <Icon className={`h-6 w-6 ${type.color === 'blue' ? 'text-blue-600' : type.color === 'green' ? 'text-green-600' : 'text-purple-600'}`} />
          </div>
          {type.recommended && (
            <Badge variant="default" className="bg-orange-500">
              Recommandé
            </Badge>
          )}
        </div>
        
        <h3 className="text-lg font-semibold mb-2">{type.name}</h3>
        <p className="text-gray-600 text-sm mb-4">{type.description}</p>
        
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-1 text-sm text-gray-500">
            <Clock className="h-4 w-4" />
            <span>{type.duration}</span>
          </div>
          {selected && (
            <CheckCircle2 className="h-5 w-5 text-green-600" />
          )}
        </div>
        
        <div className="space-y-2">
          {type.features.map((feature, index) => (
            <div key={index} className="flex items-center space-x-2">
              <div className="h-1.5 w-1.5 rounded-full bg-current opacity-60"></div>
              <span className="text-sm text-gray-600">{feature}</span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

export default function NewAnalysisPage() {
  const [selectedType, setSelectedType] = useState("comprehensive");
  const [url, setUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

    setIsLoading(true);
    
    // Simulation d'appel API - à remplacer par le vrai appel
    try {
      const response = await fetch("http://localhost:8080/api/v1/analyze", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          url,
          type: selectedType,
          options: {
            timeout: 60000,
            includeScreenshot: true,
          }
        })
      });
      
      if (response.ok) {
        const result = await response.json();
        // Redirection vers la page de suivi
        window.location.href = `/analysis/${result.data.task_id}/progress`;
      }
    } catch (error) {
      console.error("Erreur lors du lancement de l'analyse:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">
          Nouvelle Analyse SEO
        </h1>
        <p className="text-lg text-gray-600 max-w-2xl mx-auto">
          Analysez les performances SEO de votre site web avec notre technologie d&apos;intelligence artificielle avancée
        </p>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-8">
        {/* URL Input */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Globe className="h-5 w-5" />
              <span>URL à analyser</span>
            </CardTitle>
            <CardDescription>
              Entrez l&apos;URL complète du site web que vous souhaitez analyser
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <Label htmlFor="url">URL du site web</Label>
              <Input
                id="url"
                type="url"
                placeholder="https://example.com"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                required
                className="text-lg"
              />
            </div>
          </CardContent>
        </Card>

        {/* Analysis Type Selection */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Target className="h-5 w-5" />
              <span>Type d&apos;analyse</span>
            </CardTitle>
            <CardDescription>
              Choisissez le type d&apos;analyse qui correspond à vos besoins
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {analysisTypes.map((type) => (
                <AnalysisTypeCard
                  key={type.id}
                  type={type}
                  selected={selectedType === type.id}
                  onSelect={() => setSelectedType(type.id)}
                />
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Advanced Options */}
        {selectedType === "custom" && (
          <Card>
            <CardHeader>
              <CardTitle>Options avancées</CardTitle>
              <CardDescription>
                Personnalisez votre analyse selon vos besoins spécifiques
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Tabs defaultValue="seo" className="w-full">
                <TabsList className="grid w-full grid-cols-4">
                  <TabsTrigger value="seo">SEO</TabsTrigger>
                  <TabsTrigger value="performance">Performance</TabsTrigger>
                  <TabsTrigger value="content">Contenu</TabsTrigger>
                  <TabsTrigger value="ai">IA</TabsTrigger>
                </TabsList>
                
                <TabsContent value="seo" className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="meta-tags" defaultChecked />
                      <label htmlFor="meta-tags">Analyse des meta tags</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="headings" defaultChecked />
                      <label htmlFor="headings">Structure des titres</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="internal-links" defaultChecked />
                      <label htmlFor="internal-links">Liens internes</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="images" defaultChecked />
                      <label htmlFor="images">Optimisation images</label>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="performance" className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="core-vitals" defaultChecked />
                      <label htmlFor="core-vitals">Core Web Vitals</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="lighthouse" defaultChecked />
                      <label htmlFor="lighthouse">Score Lighthouse</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="mobile" defaultChecked />
                      <label htmlFor="mobile">Performance mobile</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="compression" />
                      <label htmlFor="compression">Compression GZIP</label>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="content" className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="keywords" defaultChecked />
                      <label htmlFor="keywords">Analyse mots-clés</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="readability" />
                      <label htmlFor="readability">Lisibilité</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="duplicate" />
                      <label htmlFor="duplicate">Contenu dupliqué</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="semantic" />
                      <label htmlFor="semantic">Analyse sémantique</label>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="ai" className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="ai-insights" defaultChecked />
                      <label htmlFor="ai-insights">Insights IA</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="competitors" />
                      <label htmlFor="competitors">Analyse concurrentielle</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="recommendations" defaultChecked />
                      <label htmlFor="recommendations">Recommandations auto</label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <input type="checkbox" id="trends" />
                      <label htmlFor="trends">Tendances SEO</label>
                    </div>
                  </div>
                </TabsContent>
              </Tabs>
            </CardContent>
          </Card>
        )}

        {/* Submit Button */}
        <div className="flex justify-center">
          <Button 
            type="submit" 
            size="lg" 
            disabled={!url || isLoading}
            className="bg-orange-500 hover:bg-orange-600 px-8 py-3 text-lg"
          >
            {isLoading ? (
              <>
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                Analyse en cours...
              </>
            ) : (
              <>
                <Search className="h-5 w-5 mr-2" />
                Lancer l&apos;Analyse
              </>
            )}
          </Button>
        </div>
      </form>
    </div>
  );
}