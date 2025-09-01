"""
Tests for the semantic analysis Flask server.
"""
import pytest
import json
from unittest.mock import patch, Mock
from semantic_server import app


@pytest.fixture
def client():
    """Create test client"""
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client


@pytest.fixture
def mock_analyzer():
    """Mock analyzer for testing"""
    mock = Mock()
    mock.analyze.return_value = {
        'audit_id': 'test_123',
        'model_version': 'sem-v1.0',
        'topics': [
            {
                'id': 'topic_0',
                'label': 'Logiciel juridique',
                'terms': ['logiciel', 'avocat', 'cabinet']
            }
        ],
        'suggestions': [
            {
                'keyword': 'logiciel avocat',
                'reason': 'Très pertinent pour le thématique',
                'confidence': 0.85,
                'evidence': ['https://example.com (titre)']
            }
        ],
        'metadata': {
            'schema_version': '1.0',
            'weights_version': '1.0',
            'execution_time_ms': 250,
            'lang': 'fr'
        }
    }
    return mock


class TestSemanticServer:
    """Test the Flask API endpoints"""
    
    def test_health_check(self, client):
        """Test health check endpoint"""
        response = client.get('/health')
        assert response.status_code == 200
        
        data = json.loads(response.data)
        assert data['status'] == 'healthy'
        assert data['service'] == 'semantic-analyzer'
        assert data['version'] == 'sem-v1.0'
    
    @patch('semantic_server.analyzer')
    def test_analyze_endpoint(self, mock_analyzer_global, client, mock_analyzer):
        """Test main analysis endpoint"""
        mock_analyzer_global.return_value = mock_analyzer
        mock_analyzer_global.analyze = mock_analyzer.analyze
        
        # Set the global analyzer
        import semantic_server
        semantic_server.analyzer = mock_analyzer
        
        test_data = {
            'audit_id': 'test_123',
            'crawl_data': {
                'pages': [
                    {
                        'url': 'https://example.com',
                        'lang': 'fr',
                        'title': 'Logiciel pour avocats',
                        'content': 'Notre solution de gestion...'
                    }
                ]
            }
        }
        
        response = client.post('/analyze', json=test_data)
        assert response.status_code == 200
        
        data = json.loads(response.data)
        assert data['audit_id'] == 'test_123'
        assert 'suggestions' in data
        assert 'topics' in data
        assert 'metadata' in data
    
    def test_analyze_missing_data(self, client):
        """Test analysis with missing required data"""
        # Missing audit_id
        response = client.post('/analyze', json={'crawl_data': {}})
        assert response.status_code == 400
        
        # Missing crawl_data
        response = client.post('/analyze', json={'audit_id': 'test'})
        assert response.status_code == 400
        
        # No JSON data
        response = client.post('/analyze')
        assert response.status_code == 400
    
    def test_404_handler(self, client):
        """Test 404 error handling"""
        response = client.get('/nonexistent')
        assert response.status_code == 404
        
        data = json.loads(response.data)
        assert 'error' in data