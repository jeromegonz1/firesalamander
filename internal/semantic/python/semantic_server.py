#!/usr/bin/env python3
"""
Flask server for the Semantic Analysis Agent.
Provides JSON-RPC endpoints for Fire Salamander orchestrator.
"""
import json
import logging
import sys
from flask import Flask, request, jsonify
from flask_cors import CORS
from pathlib import Path

from semantic_analyzer import SemanticAnalyzer, load_config

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Create Flask app
app = Flask(__name__)
CORS(app)

# Global analyzer instance
analyzer = None


def init_analyzer():
    """Initialize the semantic analyzer"""
    global analyzer
    try:
        # Load config from project root
        config_path = Path(__file__).parent.parent.parent.parent / 'config' / 'semantic.yaml'
        config = load_config(str(config_path))
        analyzer = SemanticAnalyzer(config)
        logger.info("Semantic analyzer initialized successfully")
        return True
    except Exception as e:
        logger.error(f"Failed to initialize analyzer: {e}")
        return False


@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        'status': 'healthy',
        'service': 'semantic-analyzer',
        'version': 'sem-v1.0'
    })


@app.route('/analyze', methods=['POST'])
def analyze_semantic():
    """
    Main analysis endpoint
    Expected input: {
        "audit_id": "string",
        "crawl_data": {...}
    }
    """
    try:
        if not analyzer:
            return jsonify({'error': 'Analyzer not initialized'}), 500
        
        data = request.get_json()
        if not data:
            return jsonify({'error': 'No JSON data provided'}), 400
        
        audit_id = data.get('audit_id')
        crawl_data = data.get('crawl_data')
        
        if not audit_id:
            return jsonify({'error': 'audit_id is required'}), 400
        
        if not crawl_data:
            return jsonify({'error': 'crawl_data is required'}), 400
        
        # Perform analysis
        result = analyzer.analyze(audit_id, crawl_data)
        
        logger.info(f"Analysis completed for audit {audit_id}")
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Analysis failed: {e}")
        return jsonify({'error': str(e)}), 500


@app.route('/topics', methods=['POST'])
def extract_topics():
    """
    Extract topics only (subset of full analysis)
    Expected input: {
        "pages": [...]
    }
    """
    try:
        if not analyzer:
            return jsonify({'error': 'Analyzer not initialized'}), 500
        
        data = request.get_json()
        pages = data.get('pages', [])
        
        if not pages:
            return jsonify({'topics': []})
        
        topics = analyzer._extract_topics(pages)
        return jsonify({'topics': topics})
        
    except Exception as e:
        logger.error(f"Topic extraction failed: {e}")
        return jsonify({'error': str(e)}), 500


@app.route('/keywords', methods=['POST'])
def generate_keywords():
    """
    Generate keyword suggestions only
    Expected input: {
        "pages": [...],
        "limit": 10  # optional
    }
    """
    try:
        if not analyzer:
            return jsonify({'error': 'Analyzer not initialized'}), 500
        
        data = request.get_json()
        pages = data.get('pages', [])
        limit = data.get('limit', 20)
        
        if not pages:
            return jsonify({'suggestions': []})
        
        # Generate candidates
        candidates = analyzer._generate_keyword_candidates(pages)
        
        # Score candidates
        scored_candidates = analyzer._score_candidates(candidates, pages)
        
        # Apply filters
        filtered = analyzer._filter_suggestions(scored_candidates)
        diverse = analyzer._apply_diversity_filter(filtered)
        
        # Limit results
        suggestions = diverse[:limit]
        
        return jsonify({'suggestions': suggestions})
        
    except Exception as e:
        logger.error(f"Keyword generation failed: {e}")
        return jsonify({'error': str(e)}), 500


@app.errorhandler(404)
def not_found(error):
    return jsonify({'error': 'Endpoint not found'}), 404


@app.errorhandler(500)
def internal_error(error):
    return jsonify({'error': 'Internal server error'}), 500


if __name__ == '__main__':
    # Initialize analyzer
    if not init_analyzer():
        logger.error("Failed to initialize analyzer, exiting")
        sys.exit(1)
    
    # Parse command line arguments
    port = 8003  # Default port for semantic agent
    host = '127.0.0.1'
    
    if len(sys.argv) > 1:
        try:
            port = int(sys.argv[1])
        except ValueError:
            logger.error("Invalid port number")
            sys.exit(1)
    
    logger.info(f"Starting semantic analysis server on {host}:{port}")
    app.run(host=host, port=port, debug=False)