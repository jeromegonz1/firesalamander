#!/usr/bin/env python3
"""
Integration test for the semantic analysis agent.
"""
import json
import requests
import time
import subprocess
import signal
import os
import sys
from pathlib import Path

def test_semantic_integration():
    """Test the complete semantic analysis pipeline"""
    
    # Sample test data
    test_data = {
        "audit_id": "integration_test_001",
        "crawl_data": {
            "pages": [
                {
                    "url": "https://example.com/",
                    "lang": "fr",
                    "title": "Logiciel de gestion pour cabinets d'avocats",
                    "h1": "Solutions numériques pour professionnels du droit",
                    "h2": ["Gestion de dossiers", "Facturation clients", "Agenda partagé"],
                    "h3": ["Suivi des échéances", "Rapports financiers"],
                    "content": "Notre logiciel de gestion permet aux cabinets d'avocats de gérer efficacement leurs dossiers clients, la facturation et les échéances. Solution cloud sécurisée pour avocats et juristes spécialisés en droit des affaires.",
                    "anchors": [
                        {"text": "Découvrir nos solutions", "href": "/solutions"},
                        {"text": "Contact avocat", "href": "/contact"}
                    ],
                    "depth": 0,
                    "outgoing_links": ["/solutions", "/contact"],
                    "incoming_links": [],
                    "canonical": "https://example.com/"
                },
                {
                    "url": "https://example.com/fonctionnalites",
                    "lang": "fr",
                    "title": "Fonctionnalités - Logiciel juridique pour avocats",
                    "h1": "Fonctionnalités complètes pour votre cabinet",
                    "h2": ["Gestion documentaire", "Comptabilité trust", "Signature électronique"],
                    "h3": ["Archives numériques", "Backup automatique"],
                    "content": "Découvrez toutes les fonctionnalités de notre solution : gestion électronique des documents juridiques, comptabilité spécialisée pour avocats, signature numérique sécurisée, et bien plus pour optimiser votre cabinet d'avocats.",
                    "anchors": [
                        {"text": "Gestion documentaire", "href": "/gestion-documents"},
                        {"text": "Cabinet avocat", "href": "/cabinet"}
                    ],
                    "depth": 1,
                    "outgoing_links": ["/gestion-documents", "/cabinet"],
                    "incoming_links": ["/"],
                    "canonical": "https://example.com/fonctionnalites"
                }
            ],
            "metadata": {
                "total_pages": 2,
                "max_depth_reached": 1,
                "duration_ms": 1500,
                "robots_respected": True,
                "sitemap_found": True
            }
        }
    }
    
    # Test health check
    print("Testing health check...")
    try:
        response = requests.get("http://localhost:8003/health", timeout=5)
        if response.status_code == 200:
            print("✓ Health check passed")
            health_data = response.json()
            print(f"  Service: {health_data.get('service')}")
            print(f"  Version: {health_data.get('version')}")
        else:
            print(f"✗ Health check failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Health check error: {e}")
        return False
    
    # Test semantic analysis
    print("\nTesting semantic analysis...")
    try:
        response = requests.post(
            "http://localhost:8003/analyze",
            json=test_data,
            timeout=30
        )
        
        if response.status_code == 200:
            print("✓ Semantic analysis passed")
            result = response.json()
            
            # Validate result structure
            required_keys = ['audit_id', 'model_version', 'topics', 'suggestions', 'metadata']
            for key in required_keys:
                if key not in result:
                    print(f"✗ Missing key in result: {key}")
                    return False
            
            print(f"  Audit ID: {result['audit_id']}")
            print(f"  Model version: {result['model_version']}")
            print(f"  Topics found: {len(result['topics'])}")
            print(f"  Suggestions: {len(result['suggestions'])}")
            print(f"  Execution time: {result['metadata']['execution_time_ms']}ms")
            
            # Display some results
            if result['topics']:
                print(f"\n  Sample topic: {result['topics'][0]['label']}")
                print(f"  Terms: {', '.join(result['topics'][0]['terms'][:5])}")
            
            if result['suggestions']:
                top_suggestion = result['suggestions'][0]
                print(f"\n  Top suggestion: '{top_suggestion['keyword']}'")
                print(f"  Confidence: {top_suggestion['confidence']:.2f}")
                print(f"  Reason: {top_suggestion['reason']}")
            
            return True
        else:
            print(f"✗ Semantic analysis failed: {response.status_code}")
            print(f"  Response: {response.text}")
            return False
            
    except Exception as e:
        print(f"✗ Semantic analysis error: {e}")
        return False


if __name__ == "__main__":
    print("=== Fire Salamander - Semantic Agent Integration Test ===")
    
    if test_semantic_integration():
        print("\n🎉 All integration tests passed!")
        sys.exit(0)
    else:
        print("\n❌ Integration tests failed!")
        sys.exit(1)