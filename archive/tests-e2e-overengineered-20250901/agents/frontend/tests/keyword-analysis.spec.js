/**
 * Fire Salamander - Keyword Analysis Section Tests
 * TDD Suite for Professional Keyword Analysis (SEMrush/Ahrefs style)
 * 
 * Requirements: Complete keyword analysis with AI suggestions
 */

const { test, expect } = require('@playwright/test');

test.describe('Keyword Analysis Section - TDD Suite', () => {
  const reportUrl = '/analysis/26/report';

  test.beforeEach(async ({ page }) => {
    await page.goto(reportUrl);
    await page.locator('[data-testid="tab-keywords"]').click();
    await expect(page.locator('[data-testid="keyword-analysis-section"]')).toBeVisible();
  });

  test.describe('Found Keywords Display', () => {
    test('should display keywords found on page with metrics', async ({ page }) => {
      const keywordsList = page.locator('[data-testid="found-keywords-list"]');
      await expect(keywordsList).toBeVisible();
      
      // First keyword item
      const firstKeyword = keywordsList.locator('[data-testid="keyword-item"]').first();
      await expect(firstKeyword).toBeVisible();
      
      // Check all metrics are displayed
      await expect(firstKeyword.locator('[data-testid="keyword-text"]')).toBeVisible();
      await expect(firstKeyword.locator('[data-testid="keyword-density"]')).toBeVisible();
      await expect(firstKeyword.locator('[data-testid="keyword-occurrences"]')).toBeVisible();
      await expect(firstKeyword.locator('[data-testid="keyword-prominence"]')).toBeVisible();
      
      // Check locations badges
      const locations = firstKeyword.locator('[data-testid="keyword-locations"] .badge');
      expect(await locations.count()).toBeGreaterThan(0);
    });

    test('should highlight keywords by location (title, h1, meta, content)', async ({ page }) => {
      const keywordItem = page.locator('[data-testid="keyword-item"]').first();
      const locationBadges = keywordItem.locator('[data-testid="location-badge"]');
      
      // Verify location types
      const badgeTexts = await locationBadges.allTextContents();
      const validLocations = ['title', 'h1', 'meta', 'content'];
      
      for (const badge of badgeTexts) {
        expect(validLocations).toContain(badge.toLowerCase());
      }
    });

    test('should sort keywords by density/prominence', async ({ page }) => {
      const sortButton = page.locator('[data-testid="sort-keywords"]');
      await sortButton.click();
      
      const sortOptions = page.locator('[data-testid="sort-option"]');
      await sortOptions.filter({ hasText: 'Densité' }).click();
      
      // Verify sorting changed
      const densities = await page.locator('[data-testid="keyword-density"]').allTextContents();
      const values = densities.map(d => parseFloat(d.replace('%', '')));
      
      // Check descending order
      for (let i = 1; i < values.length; i++) {
        expect(values[i - 1]).toBeGreaterThanOrEqual(values[i]);
      }
    });
  });

  test.describe('N-grams Analysis', () => {
    test('should display bigrams and trigrams analysis', async ({ page }) => {
      const ngramsSection = page.locator('[data-testid="ngrams-analysis"]');
      await expect(ngramsSection).toBeVisible();
      
      // Bigrams tab
      await ngramsSection.locator('[data-testid="bigrams-tab"]').click();
      const bigramsList = ngramsSection.locator('[data-testid="bigrams-list"]');
      await expect(bigramsList).toBeVisible();
      
      const firstBigram = bigramsList.locator('[data-testid="ngram-item"]').first();
      await expect(firstBigram.locator('[data-testid="ngram-phrase"]')).toBeVisible();
      await expect(firstBigram.locator('[data-testid="ngram-count"]')).toBeVisible();
      
      // Trigrams tab
      await ngramsSection.locator('[data-testid="trigrams-tab"]').click();
      const trigramsList = ngramsSection.locator('[data-testid="trigrams-list"]');
      await expect(trigramsList).toBeVisible();
    });

    test('should show ngrams visualization chart', async ({ page }) => {
      const ngramsChart = page.locator('[data-testid="ngrams-chart"]');
      await expect(ngramsChart).toBeVisible();
      
      // Check chart has data
      const chartBars = ngramsChart.locator('[data-testid="chart-bar"]');
      expect(await chartBars.count()).toBeGreaterThan(0);
    });
  });

  test.describe('AI Suggestions with ChatGPT', () => {
    test('should display AI-powered keyword suggestions', async ({ page }) => {
      const aiSection = page.locator('[data-testid="ai-suggestions"]');
      await expect(aiSection).toBeVisible();
      
      // Check AI badge
      await expect(aiSection.locator('[data-testid="ai-badge"]')).toContainText('GPT-3.5');
      
      // Check suggestions list
      const suggestionsList = aiSection.locator('[data-testid="ai-suggestions-list"]');
      await expect(suggestionsList).toBeVisible();
      
      const firstSuggestion = suggestionsList.locator('[data-testid="ai-suggestion-item"]').first();
      await expect(firstSuggestion).toBeVisible();
    });

    test('should show complete AI suggestion details', async ({ page }) => {
      const suggestion = page.locator('[data-testid="ai-suggestion-item"]').first();
      
      // All required fields
      await expect(suggestion.locator('[data-testid="suggestion-keyword"]')).toBeVisible();
      await expect(suggestion.locator('[data-testid="suggestion-volume"]')).toBeVisible();
      await expect(suggestion.locator('[data-testid="suggestion-difficulty"]')).toBeVisible();
      await expect(suggestion.locator('[data-testid="suggestion-cpc"]')).toBeVisible();
      await expect(suggestion.locator('[data-testid="suggestion-intent"]')).toBeVisible();
      await expect(suggestion.locator('[data-testid="suggestion-reason"]')).toBeVisible();
    });

    test('should expand to show content ideas', async ({ page }) => {
      const suggestion = page.locator('[data-testid="ai-suggestion-item"]').first();
      
      // Click to expand
      await suggestion.locator('[data-testid="expand-suggestion"]').click();
      
      // Check content ideas are visible
      const contentIdeas = suggestion.locator('[data-testid="content-ideas-list"]');
      await expect(contentIdeas).toBeVisible();
      
      const ideas = contentIdeas.locator('[data-testid="content-idea"]');
      expect(await ideas.count()).toBeGreaterThan(0);
    });

    test('should allow refreshing AI suggestions', async ({ page }) => {
      const refreshButton = page.locator('[data-testid="refresh-ai-suggestions"]');
      await expect(refreshButton).toBeVisible();
      
      await refreshButton.click();
      
      // Check loading state
      await expect(page.locator('[data-testid="ai-loading"]')).toBeVisible();
      
      // Wait for new suggestions
      await expect(page.locator('[data-testid="ai-loading"]')).not.toBeVisible({ timeout: 10000 });
    });
  });

  test.describe('Semantic Analysis', () => {
    test('should display main topics and entities', async ({ page }) => {
      const semanticSection = page.locator('[data-testid="semantic-analysis"]');
      await expect(semanticSection).toBeVisible();
      
      // Main topics
      const topicsList = semanticSection.locator('[data-testid="main-topics"]');
      await expect(topicsList).toBeVisible();
      const topics = topicsList.locator('[data-testid="topic-tag"]');
      expect(await topics.count()).toBeGreaterThan(0);
      
      // Entities
      const entitiesList = semanticSection.locator('[data-testid="entities-list"]');
      await expect(entitiesList).toBeVisible();
      const entities = entitiesList.locator('[data-testid="entity-item"]');
      expect(await entities.count()).toBeGreaterThan(0);
    });

    test('should show readability metrics', async ({ page }) => {
      const readabilityCard = page.locator('[data-testid="readability-metrics"]');
      await expect(readabilityCard).toBeVisible();
      
      // Check all metrics
      await expect(readabilityCard.locator('[data-testid="readability-score"]')).toBeVisible();
      await expect(readabilityCard.locator('[data-testid="readability-level"]')).toBeVisible();
      await expect(readabilityCard.locator('[data-testid="avg-sentence-length"]')).toBeVisible();
      await expect(readabilityCard.locator('[data-testid="avg-word-length"]')).toBeVisible();
      
      // Visual indicator
      await expect(readabilityCard.locator('[data-testid="readability-gauge"]')).toBeVisible();
    });

    test('should display sentiment analysis', async ({ page }) => {
      const sentimentBadge = page.locator('[data-testid="sentiment-badge"]');
      await expect(sentimentBadge).toBeVisible();
      
      // Check valid sentiment values
      const sentimentText = await sentimentBadge.textContent();
      expect(['Positive', 'Neutral', 'Negative', 'Mixed']).toContain(sentimentText);
    });
  });

  test.describe('Competitor Gap Analysis', () => {
    test('should display competitor keyword gaps when available', async ({ page }) => {
      const competitorSection = page.locator('[data-testid="competitor-gaps"]');
      
      // This section may not always be visible (depends on analysis type)
      const isVisible = await competitorSection.isVisible();
      
      if (isVisible) {
        const gapsList = competitorSection.locator('[data-testid="keyword-gaps-list"]');
        await expect(gapsList).toBeVisible();
        
        const firstGap = gapsList.locator('[data-testid="gap-item"]').first();
        await expect(firstGap.locator('[data-testid="gap-keyword"]')).toBeVisible();
        await expect(firstGap.locator('[data-testid="competitor-rank"]')).toBeVisible();
        await expect(firstGap.locator('[data-testid="gap-volume"]')).toBeVisible();
        await expect(firstGap.locator('[data-testid="gap-opportunity"]')).toBeVisible();
      }
    });

    test('should highlight high opportunity gaps', async ({ page }) => {
      const highOpportunities = page.locator('[data-testid="gap-opportunity"][data-level="high"]');
      const count = await highOpportunities.count();
      
      if (count > 0) {
        // Check styling
        const firstHigh = highOpportunities.first();
        await expect(firstHigh).toHaveClass(/high-opportunity/);
      }
    });
  });

  test.describe('Export and Actions', () => {
    test('should allow exporting keyword data', async ({ page }) => {
      const exportButton = page.locator('[data-testid="export-keywords"]');
      await expect(exportButton).toBeVisible();
      
      await exportButton.click();
      
      // Check export options
      const exportMenu = page.locator('[data-testid="export-menu"]');
      await expect(exportMenu).toBeVisible();
      await expect(exportMenu.locator('[data-testid="export-csv"]')).toBeVisible();
      await expect(exportMenu.locator('[data-testid="export-pdf"]')).toBeVisible();
    });

    test('should allow filtering keywords', async ({ page }) => {
      const filterInput = page.locator('[data-testid="filter-keywords"]');
      await expect(filterInput).toBeVisible();
      
      // Type filter
      await filterInput.fill('septeo');
      
      // Check filtered results
      const visibleKeywords = page.locator('[data-testid="keyword-item"]:visible');
      const keywords = await visibleKeywords.allTextContents();
      
      for (const keyword of keywords) {
        expect(keyword.toLowerCase()).toContain('septeo');
      }
    });
  });

  test.describe('Mobile Responsiveness', () => {
    test('should adapt layout on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      
      const keywordSection = page.locator('[data-testid="keyword-analysis-section"]');
      await expect(keywordSection).toBeVisible();
      
      // Check cards stack vertically
      const cards = keywordSection.locator('[data-testid="analysis-card"]');
      const firstCard = await cards.first().boundingBox();
      const secondCard = await cards.nth(1).boundingBox();
      
      if (firstCard && secondCard) {
        expect(secondCard.y).toBeGreaterThan(firstCard.y + firstCard.height);
      }
    });
  });
});