#!/usr/bin/env node

/**
 * Fire Salamander - API Test Runner
 * Lance les tests API avec différentes configurations
 */

const APITestAgent = require('./api_test_agent');

async function main() {
    const args = process.argv.slice(2);
    const testType = args.find(arg => arg.startsWith('--type='))?.split('=')[1] || 'all';
    const baseURL = args.find(arg => arg.startsWith('--url='))?.split('=')[1] || 'http://localhost:3000';
    
    console.log('🔥 Fire Salamander - API Test Runner');
    console.log(`Target URL: ${baseURL}`);
    console.log(`Test Type: ${testType}`);
    console.log('==================================');

    const config = {
        baseURL,
        timeout: 30000,
        reportPath: 'tests/reports/api'
    };

    const agent = new APITestAgent(config);

    try {
        let stats;
        
        switch (testType) {
            case 'contract':
                console.log('Running Contract Tests only...');
                await agent.runContractTests();
                stats = agent.stats;
                break;
                
            case 'load':
                console.log('Running Load Tests only...');
                await agent.runLoadTests();
                stats = agent.stats;
                break;
                
            case 'security':
                console.log('Running Security Tests only...');
                await agent.runSecurityTests();
                stats = agent.stats;
                break;
                
            case 'all':
            default:
                console.log('Running Full Test Suite...');
                stats = await agent.runFullTestSuite();
                break;
        }

        // Afficher un résumé
        console.log('\n📊 Test Summary');
        console.log('================');
        console.log(`Overall Score: ${stats.overallScore}/100 (${stats.status})`);
        
        if (stats.contractTests.length > 0) {
            const contractPassed = stats.contractTests.filter(t => t.passed).length;
            console.log(`Contract Tests: ${contractPassed}/${stats.contractTests.length} passed`);
        }
        
        if (stats.loadTests.length > 0) {
            const loadPassed = stats.loadTests.filter(t => t.passed).length;
            console.log(`Load Tests: ${loadPassed}/${stats.loadTests.length} passed`);
        }
        
        if (stats.securityTests.length > 0) {
            const securityPassed = stats.securityTests.filter(t => t.passed).length;
            console.log(`Security Tests: ${securityPassed}/${stats.securityTests.length} passed`);
        }

        // Exit code basé sur le score
        const exitCode = stats.overallScore >= 80 ? 0 : 1;
        console.log(`\n${exitCode === 0 ? '✅' : '❌'} Tests ${exitCode === 0 ? 'PASSED' : 'FAILED'}`);
        
        process.exit(exitCode);
        
    } catch (error) {
        console.error('❌ Test execution failed:', error.message);
        process.exit(1);
    }
}

// Vérifier si le serveur est accessible
async function checkServerHealth(baseURL) {
    try {
        const axios = require('axios');
        const response = await axios.get(`${baseURL}/health`, { timeout: 5000 });
        return response.status === 200;
    } catch (error) {
        return false;
    }
}

// Pre-flight check
async function preFlightCheck(baseURL) {
    console.log('🔍 Performing pre-flight checks...');
    
    const isHealthy = await checkServerHealth(baseURL);
    if (!isHealthy) {
        console.error(`❌ Server at ${baseURL} is not responding. Please start the application first.`);
        console.log('💡 Try running: go run main.go');
        process.exit(1);
    }
    
    console.log('✅ Server is healthy, proceeding with tests...\n');
}

// Main execution
(async () => {
    const baseURL = process.argv.find(arg => arg.startsWith('--url='))?.split('=')[1] || 'http://localhost:3000';
    
    await preFlightCheck(baseURL);
    await main();
})();