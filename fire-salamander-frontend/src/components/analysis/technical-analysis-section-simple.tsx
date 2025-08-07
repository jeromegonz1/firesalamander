"use client";

import React from 'react';
import { TechnicalAnalysis } from '@/types/technical-analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface TechnicalAnalysisSectionProps {
  technicalData: TechnicalAnalysis;
  analysisId: string;
}

export function TechnicalAnalysisSectionSimple({ 
  technicalData, 
  analysisId 
}: TechnicalAnalysisSectionProps) {

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Analyse Technique</h1>
          <div className="flex items-center space-x-4 mt-2 text-gray-600">
            <span>Analyse #{analysisId}</span>
            <span>{technicalData.pageAnalysis.length} pages analysées</span>
          </div>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Pages analysées</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{technicalData.metrics.totalPages}</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Score de santé</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{technicalData.metrics.healthScore}/100</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Problèmes détectés</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{technicalData.metrics.totalIssues}</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Temps de chargement moyen</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{technicalData.metrics.avgLoadTime}ms</div>
          </CardContent>
        </Card>
      </div>

      {/* Page List */}
      <Card>
        <CardHeader>
          <CardTitle>Pages analysées</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {technicalData.pageAnalysis.slice(0, 10).map((page, index) => (
              <div key={index} className="border p-4 rounded-lg">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">{page.url}</h3>
                    <p className="text-sm text-gray-600">
                      Status: {page.statusCode} | Taille: {Math.round(page.size / 1024)}KB | 
                      Temps: {page.loadTime}ms
                    </p>
                  </div>
                  <div className="text-right">
                    <div className={`px-2 py-1 rounded text-sm ${
                      page.statusCode === 200 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                    }`}>
                      {page.statusCode === 200 ? 'OK' : 'Erreur'}
                    </div>
                  </div>
                </div>
                
                {/* SEO Elements */}
                <div className="mt-3 grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="font-medium">Titre:</span> 
                    <span className={page.title.content ? 'text-green-600' : 'text-red-600'}>
                      {page.title.content ? `"${page.title.content.slice(0, 50)}..."` : 'Manquant'}
                    </span>
                  </div>
                  <div>
                    <span className="font-medium">Meta description:</span> 
                    <span className={page.metaDescription.content ? 'text-green-600' : 'text-red-600'}>
                      {page.metaDescription.content ? `${page.metaDescription.length} caractères` : 'Manquant'}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}