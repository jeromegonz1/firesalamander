package v2

import (
	"testing"
	"firesalamander/internal/constants"
	"github.com/stretchr/testify/assert"
)

// TestProgressCalculationFix - Test to verify that progress calculations are correct (0-100%)
func TestProgressCalculationFix(t *testing.T) {
	o := &orchestratorV2{}
	
	tests := []struct {
		step     string
		expected float64
		name     string
	}{
		{constants.PipelineStepInitializing, constants.ProgressInitialized, "initializing step"},
		{constants.PipelineStepCrawling, constants.ProgressCrawlingComplete, "crawling step"},
		{constants.PipelineStepAnalyzing, constants.ProgressAnalysisComplete, "analyzing step"},
		{constants.PipelineStepReporting, constants.ProgressReportingComplete, "reporting step"},
		{constants.PipelineStepCompleted, constants.ProgressFullComplete, "completed step"},
		{"unknown", 0.0, "unknown step"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := o.calculateProgressFromStep(tt.step)
			assert.Equal(t, tt.expected, progress, "Progress should be correct for step %s", tt.step)
			// Verify progress is never above 100% or below 0%
			assert.True(t, progress >= 0.0 && progress <= 100.0, "Progress should be between 0 and 100, got %.1f", progress)
		})
	}
}

// TestProgressValues - Test that all progress constants are reasonable percentages
func TestProgressValues(t *testing.T) {
	progressValues := []struct {
		name  string
		value float64
	}{
		{"ProgressInitialized", constants.ProgressInitialized},
		{"ProgressCrawlingComplete", constants.ProgressCrawlingComplete},
		{"ProgressAnalysisComplete", constants.ProgressAnalysisComplete},
		{"ProgressReportingComplete", constants.ProgressReportingComplete},
		{"ProgressFullComplete", constants.ProgressFullComplete},
	}

	for _, pv := range progressValues {
		t.Run(pv.name, func(t *testing.T) {
			// All progress values should be between 0 and 100
			assert.True(t, pv.value >= 0.0 && pv.value <= 100.0, 
				"%s should be between 0 and 100, got %.1f", pv.name, pv.value)
			
			// Make sure we don't have crazy values like 500%, 3000%, 10000%
			assert.True(t, pv.value < 500.0, 
				"%s should not be above 500%%, got %.1f", pv.name, pv.value)
		})
	}
}

// TestProgressNeverMultipliedBy100 - Ensure progress values are not accidentally multiplied by 100
func TestProgressNeverMultipliedBy100(t *testing.T) {
	o := &orchestratorV2{}
	
	// Test each step to make sure we don't get values like 500%, 3000%, etc.
	steps := []string{
		constants.PipelineStepInitializing,
		constants.PipelineStepCrawling,
		constants.PipelineStepAnalyzing,
		constants.PipelineStepReporting,
		constants.PipelineStepCompleted,
	}

	for _, step := range steps {
		progress := o.calculateProgressFromStep(step)
		
		// These would be the bad values if someone accidentally multiplied by 100
		badValues := []float64{500.0, 3000.0, 8000.0, 9500.0, 10000.0}
		
		for _, badValue := range badValues {
			assert.NotEqual(t, badValue, progress, 
				"Progress for step %s should not be %.1f (this suggests multiplication by 100)", step, badValue)
		}
	}
}