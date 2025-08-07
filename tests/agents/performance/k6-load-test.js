import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

/**
 * Fire Salamander - Performance Test Agent (k6)
 * Tests de charge et performance
 */

// Métriques custom
const errorRate = Rate('errors');
const responseTime = Trend('response_time');
const requestCount = Counter('requests_total');

// Configuration des scénarios de test
export const options = {
  scenarios: {
    // Test de montée en charge graduelle
    ramp_up: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '30s', target: 5 },   // Montée à 5 utilisateurs
        { duration: '1m', target: 10 },   // Montée à 10 utilisateurs
        { duration: '1m', target: 20 },   // Montée à 20 utilisateurs
        { duration: '30s', target: 0 },   // Redescente
      ],
    },
    
    // Test de charge constante
    constant_load: {
      executor: 'constant-vus',
      vus: 10,
      duration: '2m',
      startTime: '3m', // Démarre après ramp_up
    },
    
    // Test de spike
    spike_test: {
      executor: 'constant-vus',
      vus: 50,
      duration: '30s',
      startTime: '5m', // Démarre après constant_load
    }
  },
  
  // Seuils de performance
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% des requêtes < 2s
    http_req_failed: ['rate<0.1'],     // Moins de 10% d'erreurs
    errors: ['rate<0.05'],             // Moins de 5% d'erreurs custom
  },
};

// URL de base (configurable via ENV)
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  // Test de la page d'accueil
  const homeResponse = http.get(`${BASE_URL}/`);
  check(homeResponse, {
    'home page status is 200': (r) => r.status === 200,
    'home page contains Fire Salamander': (r) => r.body.includes('Fire Salamander'),
    'home page contains SEPTEO': (r) => r.body.includes('SEPTEO'),
    'home page response time < 1s': (r) => r.timings.duration < 1000,
  }) || errorRate.add(1);
  
  responseTime.add(homeResponse.timings.duration);
  requestCount.add(1);
  
  sleep(1);
  
  // Test de l'endpoint health
  const healthResponse = http.get(`${BASE_URL}/health`);
  check(healthResponse, {
    'health endpoint status is 200': (r) => r.status === 200,
    'health endpoint returns JSON': (r) => r.headers['Content-Type'].includes('application/json'),
    'health endpoint contains healthy': (r) => r.body.includes('healthy'),
    'health endpoint response time < 500ms': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);
  
  responseTime.add(healthResponse.timings.duration);
  requestCount.add(1);
  
  sleep(1);
  
  // Test de l'endpoint debug
  const debugResponse = http.get(`${BASE_URL}/debug`);
  check(debugResponse, {
    'debug endpoint status is 200': (r) => r.status === 200,
    'debug endpoint returns JSON': (r) => r.headers['Content-Type'].includes('application/json'),
    'debug endpoint response time < 2s': (r) => r.timings.duration < 2000,
  }) || errorRate.add(1);
  
  responseTime.add(debugResponse.timings.duration);
  requestCount.add(1);
  
  sleep(2);
}

// Setup : exécuté une fois au début
export function setup() {
  console.log('🔥 Fire Salamander Performance Test Starting');
  console.log(`Target: ${BASE_URL}`);
  
  // Vérifier que le serveur est accessible
  const response = http.get(`${BASE_URL}/health`);
  if (response.status !== 200) {
    console.error(`❌ Server not accessible: ${response.status}`);
    return null;
  }
  
  console.log('✅ Server is accessible, starting load tests...');
  return { baseUrl: BASE_URL };
}

// Teardown : exécuté une fois à la fin
export function teardown(data) {
  console.log('🏁 Fire Salamander Performance Test Completed');
  
  // Ici on pourrait envoyer les résultats à un système de monitoring
  // ou générer un rapport personnalisé
}