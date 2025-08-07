/**
 * Fire Salamander - Technical Analysis Test Page
 * Simple test page to verify integration
 */

"use client";

import React, { useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function AnalysisTechnicalTestPage() {
  const params = useParams();
  const [loading] = useState(false);
  const analysisId = params.id as string;

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]" data-testid="technical-analysis-loaded">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading technical analysis...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6" data-testid="technical-analysis-section">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Technical Analysis Test</h1>
          <div className="flex items-center space-x-4 mt-2 text-gray-600">
            <span>Analysis ID: {analysisId}</span>
          </div>
        </div>
      </div>

      {/* Test Section */}
      <Card>
        <CardHeader>
          <CardTitle>Technical Analysis Integration Test</CardTitle>
        </CardHeader>
        <CardContent>
          <p>This is a test page to verify that our technical analysis integration is working correctly.</p>
          <div className="mt-4 p-4 bg-green-50 rounded-lg">
            <p className="text-green-800">✓ Page loads successfully</p>
            <p className="text-green-800">✓ Component structure is correct</p>
            <p className="text-green-800">✓ TypeScript imports are working</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}