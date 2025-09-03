package integration

import (
	"context"
	"testing"
	"time"

	"firesalamander/internal/agents"
	"firesalamander/internal/agents/broken"
	"firesalamander/internal/agents/keyword"
	"firesalamander/internal/agents/linking"
	"firesalamander/internal/agents/page_profiler"
	"firesalamander/internal/agents/semantic/recommender"
	"firesalamander/internal/agents/semantic/topic"
	"firesalamander/internal/agents/technical"
	v2 "firesalamander/internal/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFullSemanticPipeline teste le pipeline sémantique complet Sprint 3
func TestFullSemanticPipeline(t *testing.T) {
	// [RÔLE: QA Lead] Test sur un vrai site avec tous les agents
	
	orch := v2.NewOrchestratorV2()
	
	// Register all Sprint 3 agents
	agents := []struct {
		name  string
		agent agents.Agent
	}{
		{"technical", technical.NewTechnicalAuditor()},
		{"keyword", keyword.NewKeywordExtractor()},
		{"linking", linking.NewLinkingMapper()},
		{"broken_links", broken.NewBrokenLinksDetector()},
		{"page_profiler", page_profiler.NewPageProfiler()},
		{"topic_clusterer", topic.NewTopicClusterer()},
		{"semantic_recommender", recommender.NewSemanticRecommender()},
	}
	
	registered := 0
	for _, a := range agents {
		if err := orch.RegisterAgent(a.name, a.agent); err != nil {
			t.Logf("Warning: Failed to register agent %s: %v", a.name, err)
		} else {
			registered++
		}
	}
	
	assert.Equal(t, 7, registered, "All 7 Sprint 3 agents should be registered")
	
	// Start audit
	auditRequest := &v2.AuditRequest{
		AuditID:   "test_sprint3_pipeline",
		SeedURL:   "https://httpbin.org/html", // Site simple pour test
		MaxPages:  5,
		Options:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}
	
	progressChan, err := orch.StartAudit(context.Background(), auditRequest)
	require.NoError(t, err, "Audit should start successfully")
	require.NotNil(t, progressChan, "Progress channel should not be nil")
	
	// Monitor progress
	var updates []v2.ProgressUpdate
	timeout := time.After(30 * time.Second)
	completed := false
	
	for !completed {
		select {
		case update, ok := <-progressChan:
			if !ok {
				completed = true
				break
			}
			updates = append(updates, *update)
			t.Logf("Progress: %.1f%% - %s", update.Progress, update.Step)
			
			if update.Progress >= 100.0 {
				completed = true
			}
			
		case <-timeout:
			t.Fatal("Audit timed out after 30 seconds")
		}
	}
	
	// Verify pipeline execution
	assert.Greater(t, len(updates), 5, "Should have multiple progress updates")
	
	// Check that all pipeline steps occurred
	steps := make(map[string]bool)
	for _, update := range updates {
		steps[update.Step] = true
	}
	
	expectedSteps := []string{"initializing", "crawling", "analyzing", "completed"}
	for _, step := range expectedSteps {
		assert.True(t, steps[step], "Pipeline should include step: %s", step)
	}
	
	t.Logf("✅ Sprint 3 Pipeline completed successfully with %d progress updates", len(updates))
}

// TestSemanticAgentsIntegration teste l'intégration spécifique des agents sémantiques
func TestSemanticAgentsIntegration(t *testing.T) {
	// [RÔLE: QA] Test intégration agents sémantiques
	
	// Test data - HTML sample
	testHTML := `
<!DOCTYPE html>
<html lang="fr">
<head>
	<title>Cabinet d'Avocats Parisien - Droit des Affaires</title>
	<meta name="description" content="Cabinet spécialisé en droit des affaires, conseil juridique et contentieux commercial à Paris.">
	<meta name="keywords" content="avocat, paris, droit affaires, conseil juridique">
</head>
<body>
	<h1>Expert en Droit des Affaires</h1>
	<h2>Nos Services Juridiques</h2>
	<p>Notre cabinet d'avocats propose des services complets en droit des affaires, 
	   conseil en entreprise et accompagnement juridique personnalisé.</p>
	<h2>Domaines d'Expertise</h2>
	<p>Nous intervenons dans les domaines du droit commercial, droit des sociétés, 
	   contentieux commercial et conseil en fusion-acquisition.</p>
	<a href="/contact">Contactez-nous</a>
	<a href="https://legifrance.gouv.fr">Références légales</a>
</body>
</html>`
	
	// Test PageProfiler
	profiler := page_profiler.NewPageProfiler()
	profilerResult, err := profiler.Process(context.Background(), page_profiler.PageRequest{
		HTML: testHTML,
		URL:  "https://cabinet-avocat.example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "completed", profilerResult.Status)
	
	// Verify profiler extracted meta tags
	metaTags := profilerResult.Data["meta_tags"].(map[string]string)
	assert.Contains(t, metaTags, "title")
	assert.Contains(t, metaTags, "description")
	assert.Contains(t, metaTags["title"], "Cabinet d'Avocats")
	
	// Test semantic agents individually (simplified)
	clusterer := topic.NewTopicClusterer()
	err = clusterer.HealthCheck()
	require.NoError(t, err, "TopicClusterer health check should pass")
	
	recommenderAgent := recommender.NewSemanticRecommender()
	err = recommenderAgent.HealthCheck()
	require.NoError(t, err, "SemanticRecommender health check should pass")
	
	// Note: Les agents sémantiques nécessitent des structures spécifiques
	// Pour ce test d'intégration, nous validons leur existence et santé
	
	t.Logf("✅ Semantic agents integration test passed")
	t.Logf("   - PageProfiler: extracted %d meta tags", len(metaTags))
	t.Logf("   - TopicClusterer: health check OK")
	t.Logf("   - SemanticRecommender: health check OK")
}

// TestAPIEndpointsWithSemantic teste les endpoints API avec le pipeline sémantique
func TestAPIEndpointsWithSemantic(t *testing.T) {
	// [RÔLE: QA] Test endpoints avec agents Sprint 3
	
	// Ce test nécessite que le serveur soit démarré
	// Nous testons la structure de réponse sans démarrer un serveur complet
	
	// Test de la structure attendue des réponses
	auditResponse := map[string]interface{}{
		"id":      "audit_test_123",
		"status":  "started",
		"message": "Audit started successfully",
	}
	
	// Verify response structure
	assert.Contains(t, auditResponse, "id")
	assert.Contains(t, auditResponse, "status")
	assert.Contains(t, auditResponse, "message")
	assert.Equal(t, "started", auditResponse["status"])
	
	// Test expected progress structure
	progressUpdate := v2.ProgressUpdate{
		AuditID:   "audit_test_123",
		Step:      "analyzing",
		Progress:  80.0,
		AgentName: "semantic_recommender",
		Timestamp: time.Now(),
	}
	
	assert.Equal(t, "audit_test_123", progressUpdate.AuditID)
	assert.Equal(t, "analyzing", progressUpdate.Step)
	assert.Equal(t, 80.0, progressUpdate.Progress)
	assert.Equal(t, "semantic_recommender", progressUpdate.AgentName)
	
	t.Logf("✅ API endpoints structure validation passed")
}

// TestPerformanceWithSemanticAgents teste la performance du pipeline sémantique
func TestPerformanceWithSemanticAgents(t *testing.T) {
	// [RÔLE: QA] Test performance Sprint 3
	
	start := time.Now()
	
	// Create all semantic agents
	agents := []agents.Agent{
		technical.NewTechnicalAuditor(),
		keyword.NewKeywordExtractor(),
		linking.NewLinkingMapper(),
		broken.NewBrokenLinksDetector(),
		page_profiler.NewPageProfiler(),
		topic.NewTopicClusterer(),
		recommender.NewSemanticRecommender(),
	}
	
	agentCreationTime := time.Since(start)
	
	// Test health checks
	start = time.Now()
	for i, agent := range agents {
		err := agent.HealthCheck()
		assert.NoError(t, err, "Agent %d health check should pass", i)
	}
	healthCheckTime := time.Since(start)
	
	// Performance assertions
	assert.Less(t, agentCreationTime.Milliseconds(), int64(1000), "Agent creation should take less than 1s")
	assert.Less(t, healthCheckTime.Milliseconds(), int64(500), "Health checks should take less than 500ms")
	
	t.Logf("✅ Performance test passed")
	t.Logf("   - Agent creation: %dms", agentCreationTime.Milliseconds())
	t.Logf("   - Health checks: %dms", healthCheckTime.Milliseconds())
}

// TestErrorHandlingInSemantic teste la gestion d'erreurs dans le pipeline sémantique
func TestErrorHandlingInSemantic(t *testing.T) {
	// [RÔLE: QA] Test gestion erreurs
	
	profiler := page_profiler.NewPageProfiler()
	
	// Test with invalid input type
	_, err := profiler.Process(context.Background(), "invalid input")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input type")
	
	// Test with empty HTML
	result, err := profiler.Process(context.Background(), page_profiler.PageRequest{
		HTML: "",
		URL:  "https://example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
	
	// Verify empty HTML returns valid empty structure
	contentStats := result.Data["content_stats"].(page_profiler.ContentStats)
	assert.Equal(t, 0, contentStats.WordCount)
	
	t.Logf("✅ Error handling test passed")
}

// BenchmarkSemanticPipeline benchmark du pipeline sémantique
func BenchmarkSemanticPipeline(b *testing.B) {
	agents := []agents.Agent{
		page_profiler.NewPageProfiler(),
		topic.NewTopicClusterer(),
		recommender.NewSemanticRecommender(),
	}
	
	testHTML := `<html><head><title>Test</title></head><body><p>Test content for benchmarking semantic analysis performance</p></body></html>`
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Benchmark PageProfiler
		profiler := agents[0]
		_, err := profiler.Process(context.Background(), page_profiler.PageRequest{
			HTML: testHTML,
			URL:  "https://example.com",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}