#!/usr/bin/env python3
"""
Fire Salamander - Monitoring Agent
Agent de surveillance système et monitoring de santé
"""

import json
import os
import time
import requests
import psutil
import subprocess
from datetime import datetime, timedelta
from typing import Dict, List, Optional
import logging
import threading
import signal
import sys

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class MonitoringAgent:
    """Agent de monitoring pour Fire Salamander"""
    
    def __init__(self, config: Optional[Dict] = None):
        self.config = config or self.default_config()
        self.stats = {
            'timestamp': datetime.now().isoformat(),
            'target_url': self.config['target_url'],
            'monitoring_results': {},
            'alerts': [],
            'metrics': {},
            'overall_status': 'unknown',
            'uptime_percentage': 0.0
        }
        self.monitoring_active = False
        self.monitoring_thread = None
        
    def default_config(self) -> Dict:
        """Configuration par défaut"""
        return {
            'target_url': 'http://localhost:3000',
            'check_interval': 30,  # secondes
            'timeout': 10,
            'report_path': 'tests/reports/monitoring',
            'duration': 300,  # 5 minutes de monitoring par défaut
            'thresholds': {
                'response_time_warning': 2.0,  # secondes
                'response_time_critical': 5.0,
                'cpu_warning': 80.0,  # pourcentage
                'cpu_critical': 95.0,
                'memory_warning': 80.0,
                'memory_critical': 95.0,
                'disk_warning': 80.0,
                'disk_critical': 95.0
            },
            'endpoints_to_monitor': [
                {'path': '/', 'name': 'Homepage'},
                {'path': '/health', 'name': 'Health Check'},
                {'path': '/debug', 'name': 'Debug Info'}
            ]
        }
    
    def run_monitoring_session(self) -> Dict:
        """Lance une session de monitoring complète"""
        logger.info("🔍 Starting monitoring session")
        
        try:
            # 1. Test initial de connectivité
            if not self._test_initial_connectivity():
                self.stats['overall_status'] = 'critical'
                self.stats['alerts'].append({
                    'level': 'critical',
                    'message': 'Initial connectivity test failed',
                    'timestamp': datetime.now().isoformat()
                })
                return self.stats
            
            # 2. Monitoring continu
            self._start_continuous_monitoring()
            
            # 3. Tests de charge légers
            self._run_light_load_tests()
            
            # 4. Vérification des ressources système
            self._check_system_resources()
            
            # 5. Génération du rapport
            self._generate_report()
            
            logger.info(f"✅ Monitoring session completed - Status: {self.stats['overall_status']}")
            
        except KeyboardInterrupt:
            logger.info("⏹️ Monitoring interrupted by user")
            self.monitoring_active = False
        except Exception as e:
            logger.error(f"❌ Monitoring session failed: {str(e)}")
            self.stats['overall_status'] = 'error'
            self.stats['error'] = str(e)
        
        return self.stats
    
    def _test_initial_connectivity(self) -> bool:
        """Test initial de connectivité"""
        logger.info("🔌 Testing initial connectivity")
        
        connectivity_tests = []
        
        for endpoint in self.config['endpoints_to_monitor']:
            url = f"{self.config['target_url']}{endpoint['path']}"
            
            try:
                start_time = time.time()
                response = requests.get(url, timeout=self.config['timeout'])
                response_time = time.time() - start_time
                
                test_result = {
                    'endpoint': endpoint['name'],
                    'url': url,
                    'status_code': response.status_code,
                    'response_time': response_time,
                    'status': 'pass' if response.status_code == 200 else 'fail',
                    'timestamp': datetime.now().isoformat()
                }
                
                if response.status_code != 200:
                    test_result['status'] = 'fail'
                    self.stats['alerts'].append({
                        'level': 'warning',
                        'endpoint': endpoint['name'],
                        'message': f'HTTP {response.status_code} response',
                        'timestamp': datetime.now().isoformat()
                    })
                elif response_time > self.config['thresholds']['response_time_warning']:
                    test_result['status'] = 'warning'
                    self.stats['alerts'].append({
                        'level': 'warning',
                        'endpoint': endpoint['name'],
                        'message': f'Slow response time: {response_time:.2f}s',
                        'timestamp': datetime.now().isoformat()
                    })
                
                connectivity_tests.append(test_result)
                
            except requests.exceptions.RequestException as e:
                connectivity_tests.append({
                    'endpoint': endpoint['name'],
                    'url': url,
                    'status': 'fail',
                    'error': str(e),
                    'timestamp': datetime.now().isoformat()
                })
                
                self.stats['alerts'].append({
                    'level': 'critical',
                    'endpoint': endpoint['name'],
                    'message': f'Connection failed: {str(e)}',
                    'timestamp': datetime.now().isoformat()
                })
        
        self.stats['monitoring_results']['connectivity_tests'] = connectivity_tests
        
        # Retourner True si au moins un endpoint répond
        return any(test['status'] == 'pass' for test in connectivity_tests)
    
    def _start_continuous_monitoring(self):
        """Démarre le monitoring continu"""
        logger.info(f"🔄 Starting continuous monitoring for {self.config['duration']} seconds")
        
        self.monitoring_active = True
        start_time = time.time()
        end_time = start_time + self.config['duration']
        
        monitoring_data = {
            'checks': [],
            'uptime_checks': 0,
            'successful_checks': 0,
            'failed_checks': 0,
            'average_response_time': 0.0,
            'min_response_time': float('inf'),
            'max_response_time': 0.0
        }
        
        response_times = []
        
        while time.time() < end_time and self.monitoring_active:
            check_start = time.time()
            
            # Vérifier chaque endpoint
            for endpoint in self.config['endpoints_to_monitor']:
                url = f"{self.config['target_url']}{endpoint['path']}"
                
                try:
                    response = requests.get(url, timeout=self.config['timeout'])
                    response_time = time.time() - check_start
                    
                    check_result = {
                        'timestamp': datetime.now().isoformat(),
                        'endpoint': endpoint['name'],
                        'status_code': response.status_code,
                        'response_time': response_time,
                        'status': 'success' if response.status_code == 200 else 'fail'
                    }
                    
                    monitoring_data['checks'].append(check_result)
                    monitoring_data['uptime_checks'] += 1
                    
                    if response.status_code == 200:
                        monitoring_data['successful_checks'] += 1
                        response_times.append(response_time)
                        monitoring_data['min_response_time'] = min(monitoring_data['min_response_time'], response_time)
                        monitoring_data['max_response_time'] = max(monitoring_data['max_response_time'], response_time)
                        
                        # Vérifier les seuils
                        if response_time > self.config['thresholds']['response_time_critical']:
                            self.stats['alerts'].append({
                                'level': 'critical',
                                'endpoint': endpoint['name'],
                                'message': f'Critical response time: {response_time:.2f}s',
                                'timestamp': datetime.now().isoformat()
                            })
                        elif response_time > self.config['thresholds']['response_time_warning']:
                            self.stats['alerts'].append({
                                'level': 'warning',
                                'endpoint': endpoint['name'],
                                'message': f'High response time: {response_time:.2f}s',
                                'timestamp': datetime.now().isoformat()
                            })
                    else:
                        monitoring_data['failed_checks'] += 1
                        self.stats['alerts'].append({
                            'level': 'error',
                            'endpoint': endpoint['name'],
                            'message': f'HTTP {response.status_code} error',
                            'timestamp': datetime.now().isoformat()
                        })
                        
                except requests.exceptions.RequestException as e:
                    monitoring_data['checks'].append({
                        'timestamp': datetime.now().isoformat(),
                        'endpoint': endpoint['name'],
                        'status': 'error',
                        'error': str(e)
                    })
                    
                    monitoring_data['uptime_checks'] += 1
                    monitoring_data['failed_checks'] += 1
                    
                    self.stats['alerts'].append({
                        'level': 'critical',
                        'endpoint': endpoint['name'],
                        'message': f'Connection error: {str(e)}',
                        'timestamp': datetime.now().isoformat()
                    })
            
            # Attendre l'intervalle suivant
            time.sleep(self.config['check_interval'])
        
        # Calculer les métriques finales
        if response_times:
            monitoring_data['average_response_time'] = sum(response_times) / len(response_times)
        else:
            monitoring_data['min_response_time'] = 0.0
        
        if monitoring_data['uptime_checks'] > 0:
            self.stats['uptime_percentage'] = (monitoring_data['successful_checks'] / monitoring_data['uptime_checks']) * 100
        
        self.stats['monitoring_results']['continuous_monitoring'] = monitoring_data
        logger.info(f"📊 Uptime: {self.stats['uptime_percentage']:.1f}%")
    
    def _run_light_load_tests(self):
        """Lance des tests de charge légers"""
        logger.info("⚡ Running light load tests")
        
        load_tests = []
        
        # Test 1: Requêtes simultanées (5 threads)
        concurrent_results = self._test_concurrent_requests(5, 10)
        load_tests.append({
            'test': 'Concurrent Requests (5 threads)',
            'results': concurrent_results,
            'status': 'pass' if concurrent_results['success_rate'] > 90 else 'warning'
        })
        
        # Test 2: Burst test (10 requêtes rapides)
        burst_results = self._test_burst_requests(10)
        load_tests.append({
            'test': 'Burst Test (10 requests)',
            'results': burst_results,
            'status': 'pass' if burst_results['success_rate'] > 90 else 'warning'
        })
        
        self.stats['monitoring_results']['load_tests'] = load_tests
    
    def _test_concurrent_requests(self, num_threads: int, requests_per_thread: int) -> Dict:
        """Test de requêtes simultanées"""
        results = []
        threads = []
        
        def make_requests():
            thread_results = []
            for _ in range(requests_per_thread):
                try:
                    start_time = time.time()
                    response = requests.get(self.config['target_url'], timeout=self.config['timeout'])
                    response_time = time.time() - start_time
                    
                    thread_results.append({
                        'status_code': response.status_code,
                        'response_time': response_time,
                        'success': response.status_code == 200
                    })
                except Exception as e:
                    thread_results.append({
                        'error': str(e),
                        'success': False
                    })
            
            results.extend(thread_results)
        
        # Lancer les threads
        for _ in range(num_threads):
            thread = threading.Thread(target=make_requests)
            threads.append(thread)
            thread.start()
        
        # Attendre que tous se terminent
        for thread in threads:
            thread.join()
        
        # Analyser les résultats
        total_requests = len(results)
        successful_requests = sum(1 for r in results if r.get('success', False))
        success_rate = (successful_requests / total_requests) * 100 if total_requests > 0 else 0
        
        response_times = [r['response_time'] for r in results if 'response_time' in r]
        avg_response_time = sum(response_times) / len(response_times) if response_times else 0
        
        return {
            'total_requests': total_requests,
            'successful_requests': successful_requests,
            'success_rate': success_rate,
            'average_response_time': avg_response_time,
            'max_response_time': max(response_times) if response_times else 0,
            'min_response_time': min(response_times) if response_times else 0
        }
    
    def _test_burst_requests(self, num_requests: int) -> Dict:
        """Test de rafale de requêtes"""
        results = []
        
        start_time = time.time()
        for _ in range(num_requests):
            try:
                req_start = time.time()
                response = requests.get(self.config['target_url'], timeout=self.config['timeout'])
                response_time = time.time() - req_start
                
                results.append({
                    'status_code': response.status_code,
                    'response_time': response_time,
                    'success': response.status_code == 200
                })
            except Exception as e:
                results.append({
                    'error': str(e),
                    'success': False
                })
        
        total_time = time.time() - start_time
        
        successful_requests = sum(1 for r in results if r.get('success', False))
        success_rate = (successful_requests / num_requests) * 100
        
        response_times = [r['response_time'] for r in results if 'response_time' in r]
        avg_response_time = sum(response_times) / len(response_times) if response_times else 0
        
        return {
            'total_requests': num_requests,
            'successful_requests': successful_requests,
            'success_rate': success_rate,
            'total_time': total_time,
            'requests_per_second': num_requests / total_time if total_time > 0 else 0,
            'average_response_time': avg_response_time
        }
    
    def _check_system_resources(self):
        """Vérifie les ressources système"""
        logger.info("💻 Checking system resources")
        
        system_metrics = {}
        
        try:
            # CPU Usage
            cpu_percent = psutil.cpu_percent(interval=1)
            system_metrics['cpu'] = {
                'usage_percent': cpu_percent,
                'status': self._get_resource_status('cpu', cpu_percent)
            }
            
            if cpu_percent > self.config['thresholds']['cpu_warning']:
                level = 'critical' if cpu_percent > self.config['thresholds']['cpu_critical'] else 'warning'
                self.stats['alerts'].append({
                    'level': level,
                    'message': f'High CPU usage: {cpu_percent:.1f}%',
                    'timestamp': datetime.now().isoformat()
                })
            
            # Memory Usage
            memory = psutil.virtual_memory()
            system_metrics['memory'] = {
                'total_gb': round(memory.total / (1024**3), 2),
                'used_gb': round(memory.used / (1024**3), 2),
                'usage_percent': memory.percent,
                'status': self._get_resource_status('memory', memory.percent)
            }
            
            if memory.percent > self.config['thresholds']['memory_warning']:
                level = 'critical' if memory.percent > self.config['thresholds']['memory_critical'] else 'warning'
                self.stats['alerts'].append({
                    'level': level,
                    'message': f'High memory usage: {memory.percent:.1f}%',
                    'timestamp': datetime.now().isoformat()
                })
            
            # Disk Usage
            disk = psutil.disk_usage('/')
            disk_percent = (disk.used / disk.total) * 100
            system_metrics['disk'] = {
                'total_gb': round(disk.total / (1024**3), 2),
                'used_gb': round(disk.used / (1024**3), 2),
                'free_gb': round(disk.free / (1024**3), 2),
                'usage_percent': disk_percent,
                'status': self._get_resource_status('disk', disk_percent)
            }
            
            if disk_percent > self.config['thresholds']['disk_warning']:
                level = 'critical' if disk_percent > self.config['thresholds']['disk_critical'] else 'warning'
                self.stats['alerts'].append({
                    'level': level,
                    'message': f'High disk usage: {disk_percent:.1f}%',
                    'timestamp': datetime.now().isoformat()
                })
            
            # Network I/O
            network = psutil.net_io_counters()
            system_metrics['network'] = {
                'bytes_sent': network.bytes_sent,
                'bytes_recv': network.bytes_recv,
                'packets_sent': network.packets_sent,
                'packets_recv': network.packets_recv
            }
            
            # Process Count
            process_count = len(psutil.pids())
            system_metrics['processes'] = {
                'count': process_count,
                'status': 'normal' if process_count < 500 else 'warning'
            }
            
        except Exception as e:
            logger.warning(f"Could not gather all system metrics: {e}")
            system_metrics['error'] = str(e)
        
        self.stats['metrics']['system'] = system_metrics
    
    def _get_resource_status(self, resource_type: str, usage_percent: float) -> str:
        """Détermine le statut d'une ressource"""
        warning_threshold = self.config['thresholds'][f'{resource_type}_warning']
        critical_threshold = self.config['thresholds'][f'{resource_type}_critical']
        
        if usage_percent >= critical_threshold:
            return 'critical'
        elif usage_percent >= warning_threshold:
            return 'warning'
        else:
            return 'normal'
    
    def _determine_overall_status(self):
        """Détermine le statut global du système"""
        # Compter les alertes par niveau
        critical_alerts = sum(1 for alert in self.stats['alerts'] if alert['level'] == 'critical')
        warning_alerts = sum(1 for alert in self.stats['alerts'] if alert['level'] == 'warning')
        
        # Vérifier l'uptime
        uptime = self.stats.get('uptime_percentage', 0)
        
        if critical_alerts > 0 or uptime < 50:
            self.stats['overall_status'] = 'critical'
        elif warning_alerts > 5 or uptime < 90:
            self.stats['overall_status'] = 'warning'
        elif uptime >= 99:
            self.stats['overall_status'] = 'excellent'
        elif uptime >= 95:
            self.stats['overall_status'] = 'good'
        else:
            self.stats['overall_status'] = 'acceptable'
    
    def _generate_report(self):
        """Génère le rapport de monitoring"""
        try:
            os.makedirs(self.config['report_path'], exist_ok=True)
            
            # Déterminer le statut global
            self._determine_overall_status()
            
            # Rapport JSON
            json_report_path = os.path.join(self.config['report_path'], 'monitoring_report.json')
            with open(json_report_path, 'w') as f:
                json.dump(self.stats, f, indent=2)
            
            # Rapport HTML
            html_report = self._generate_html_report()
            html_report_path = os.path.join(self.config['report_path'], 'monitoring_report.html')
            with open(html_report_path, 'w') as f:
                f.write(html_report)
            
            logger.info(f"📊 Monitoring report generated: {json_report_path}")
            
        except Exception as e:
            logger.error(f"Failed to generate monitoring report: {str(e)}")
    
    def _generate_html_report(self) -> str:
        """Génère un rapport HTML"""
        
        # Compter les alertes
        critical_alerts = sum(1 for alert in self.stats['alerts'] if alert['level'] == 'critical')
        warning_alerts = sum(1 for alert in self.stats['alerts'] if alert['level'] == 'warning')
        error_alerts = sum(1 for alert in self.stats['alerts'] if alert['level'] == 'error')
        
        html = f"""
<!DOCTYPE html>
<html>
<head>
    <title>Fire Salamander - Monitoring Report</title>
    <meta charset="UTF-8">
    <style>
        body {{ font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }}
        .header {{ background: #ff6136; color: white; padding: 20px; border-radius: 8px; }}
        .status {{ font-size: 1.5em; font-weight: bold; }}
        .status-excellent {{ color: #28a745; }}
        .status-good {{ color: #17a2b8; }}
        .status-acceptable {{ color: #ffc107; }}
        .status-warning {{ color: #fd7e14; }}
        .status-critical {{ color: #dc3545; }}
        .section {{ background: white; margin: 20px 0; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
        .metric {{ display: inline-block; margin: 10px; padding: 15px; background: #f8f9fa; border-radius: 5px; min-width: 120px; text-align: center; }}
        .alert {{ margin: 10px 0; padding: 10px; border-left: 4px solid #ddd; border-radius: 4px; }}
        .alert-critical {{ border-left-color: #dc3545; background-color: #f8d7da; }}
        .alert-warning {{ border-left-color: #ffc107; background-color: #fff3cd; }}
        .alert-error {{ border-left-color: #fd7e14; background-color: #f8d7da; }}
        table {{ width: 100%; border-collapse: collapse; margin: 10px 0; }}
        th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
        th {{ background-color: #f2f2f2; }}
        .resource-normal {{ color: #28a745; }}
        .resource-warning {{ color: #ffc107; }}
        .resource-critical {{ color: #dc3545; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>📊 Fire Salamander - Monitoring Report</h1>
        <div class="status status-{self.stats['overall_status']}">
            Status: {self.stats['overall_status'].upper()}
        </div>
        <p>Target: {self.stats['target_url']}</p>
        <p>Generated: {self.stats['timestamp']}</p>
    </div>
    
    <div class="section">
        <h2>📈 Overview</h2>
        <div class="metric">
            <strong>{self.stats['uptime_percentage']:.1f}%</strong><br>
            Uptime
        </div>
        <div class="metric">
            <strong>{len(self.stats['alerts'])}</strong><br>
            Total Alerts
        </div>
        <div class="metric">
            <strong>{critical_alerts}</strong><br>
            Critical
        </div>
        <div class="metric">
            <strong>{warning_alerts}</strong><br>
            Warnings
        </div>
        <div class="metric">
            <strong>{error_alerts}</strong><br>
            Errors
        </div>
    </div>
        """
        
        # Section des alertes
        if self.stats['alerts']:
            html += '<div class="section"><h2>🚨 Alerts</h2>'
            for alert in self.stats['alerts'][-20:]:  # Dernières 20 alertes
                alert_class = f"alert-{alert['level']}"
                endpoint_info = f" [{alert['endpoint']}]" if 'endpoint' in alert else ""
                html += f"""
                <div class="alert {alert_class}">
                    <strong>[{alert['level'].upper()}]{endpoint_info}</strong> {alert['message']}<br>
                    <small>{alert['timestamp']}</small>
                </div>
                """
            html += '</div>'
        
        # Section des métriques système
        if 'system' in self.stats['metrics']:
            system = self.stats['metrics']['system']
            html += '<div class="section"><h2>💻 System Resources</h2>'
            html += '<table><tr><th>Resource</th><th>Current</th><th>Status</th></tr>'
            
            if 'cpu' in system:
                status_class = f"resource-{system['cpu']['status']}"
                html += f'<tr><td>CPU Usage</td><td>{system["cpu"]["usage_percent"]:.1f}%</td><td class="{status_class}">{system["cpu"]["status"].upper()}</td></tr>'
            
            if 'memory' in system:
                status_class = f"resource-{system['memory']['status']}"
                html += f'<tr><td>Memory Usage</td><td>{system["memory"]["usage_percent"]:.1f}% ({system["memory"]["used_gb"]} GB / {system["memory"]["total_gb"]} GB)</td><td class="{status_class}">{system["memory"]["status"].upper()}</td></tr>'
            
            if 'disk' in system:
                status_class = f"resource-{system['disk']['status']}"
                html += f'<tr><td>Disk Usage</td><td>{system["disk"]["usage_percent"]:.1f}% ({system["disk"]["used_gb"]} GB / {system["disk"]["total_gb"]} GB)</td><td class="{status_class}">{system["disk"]["status"].upper()}</td></tr>'
            
            if 'processes' in system:
                html += f'<tr><td>Processes</td><td>{system["processes"]["count"]}</td><td>{system["processes"]["status"].upper()}</td></tr>'
            
            html += '</table></div>'
        
        # Section des tests de connectivité
        if 'connectivity_tests' in self.stats['monitoring_results']:
            html += '<div class="section"><h2>🔌 Connectivity Tests</h2>'
            html += '<table><tr><th>Endpoint</th><th>Status</th><th>Response Time</th><th>Status Code</th></tr>'
            
            for test in self.stats['monitoring_results']['connectivity_tests']:
                status_color = '#28a745' if test['status'] == 'pass' else '#dc3545' if test['status'] == 'fail' else '#ffc107'
                response_time = f"{test.get('response_time', 0):.3f}s" if 'response_time' in test else 'N/A'
                status_code = test.get('status_code', 'N/A')
                
                html += f'<tr><td>{test["endpoint"]}</td><td style="color: {status_color}">{test["status"].upper()}</td><td>{response_time}</td><td>{status_code}</td></tr>'
            
            html += '</table></div>'
        
        # Section du monitoring continu
        if 'continuous_monitoring' in self.stats['monitoring_results']:
            monitoring = self.stats['monitoring_results']['continuous_monitoring']
            html += '<div class="section"><h2>🔄 Continuous Monitoring</h2>'
            html += f"""
            <div class="metric">
                <strong>{monitoring['uptime_checks']}</strong><br>
                Total Checks
            </div>
            <div class="metric">
                <strong>{monitoring['successful_checks']}</strong><br>
                Successful
            </div>
            <div class="metric">
                <strong>{monitoring['failed_checks']}</strong><br>
                Failed
            </div>
            <div class="metric">
                <strong>{monitoring.get('average_response_time', 0):.3f}s</strong><br>
                Avg Response
            </div>
            <div class="metric">
                <strong>{monitoring.get('max_response_time', 0):.3f}s</strong><br>
                Max Response
            </div>
            """
            html += '</div>'
        
        # Section des tests de charge
        if 'load_tests' in self.stats['monitoring_results']:
            html += '<div class="section"><h2>⚡ Load Tests</h2>'
            
            for test in self.stats['monitoring_results']['load_tests']:
                results = test['results']
                status_color = '#28a745' if test['status'] == 'pass' else '#ffc107'
                
                html += f"""
                <h3 style="color: {status_color}">{test['test']} - {test['status'].upper()}</h3>
                <p><strong>Success Rate:</strong> {results.get('success_rate', 0):.1f}%</p>
                <p><strong>Average Response Time:</strong> {results.get('average_response_time', 0):.3f}s</p>
                <p><strong>Total Requests:</strong> {results.get('total_requests', 0)}</p>
                """
                
                if 'requests_per_second' in results:
                    html += f"<p><strong>Requests/Second:</strong> {results['requests_per_second']:.1f}</p>"
        
        html += '</div>'
        
        html += """
</body>
</html>
        """
        
        return html
    
    def stop_monitoring(self):
        """Arrête le monitoring"""
        self.monitoring_active = False
        if self.monitoring_thread and self.monitoring_thread.is_alive():
            self.monitoring_thread.join()

def signal_handler(signum, frame):
    """Gestionnaire de signal pour arrêt propre"""
    logger.info("Received signal, stopping monitoring...")
    sys.exit(0)

def main():
    """Point d'entrée principal"""
    import argparse
    
    # Gestionnaire de signal pour Ctrl+C
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    parser = argparse.ArgumentParser(description='Fire Salamander Monitoring Agent')
    parser.add_argument('--url', default='http://localhost:3000', help='Target URL')
    parser.add_argument('--duration', type=int, default=300, help='Monitoring duration in seconds')
    parser.add_argument('--interval', type=int, default=30, help='Check interval in seconds')
    parser.add_argument('--output', default='tests/reports/monitoring', help='Output directory')
    args = parser.parse_args()
    
    config = {
        'target_url': args.url,
        'duration': args.duration,
        'check_interval': args.interval,
        'report_path': args.output
    }
    
    agent = MonitoringAgent(config)
    
    try:
        results = agent.run_monitoring_session()
        
        print(f"\n📊 Monitoring Results:")
        print(f"Status: {results['overall_status']}")
        print(f"Uptime: {results['uptime_percentage']:.1f}%")
        print(f"Alerts: {len(results['alerts'])}")
        print(f"Report: {args.output}/monitoring_report.html")
        
        # Exit code basé sur le statut
        exit_code = 0
        if results['overall_status'] in ['critical', 'warning']:
            exit_code = 1
        
        sys.exit(exit_code)
        
    except KeyboardInterrupt:
        logger.info("⏹️ Monitoring stopped by user")
        agent.stop_monitoring()
        sys.exit(0)

if __name__ == '__main__':
    main()