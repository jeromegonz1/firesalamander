# Fire Salamander - Spécifications

## Structure
- **[functional/](functional/)** - Spécifications fonctionnelles détaillées des agents
- **[technical/](technical/)** - Architecture technique et contrats API
- **[test-scenarios/](test-scenarios/)** - Scénarios de test et validation

## Agents couverts
1. Agent Crawler - Exploration intelligente des sites
2. Agent Audit Technique - Analyse SEO et performance
3. Agent Analyse Sémantique - Compréhension métier et suggestions
4. Agent Reporting - Génération de rapports
5. Agent Orchestrateur - Coordination du pipeline

## Contrats API
Les schémas JSON-RPC sont définis dans [technical/api-contracts/](technical/api-contracts/).

## Validation
Utilisez `make validate-schemas` pour valider tous les contrats JSON.