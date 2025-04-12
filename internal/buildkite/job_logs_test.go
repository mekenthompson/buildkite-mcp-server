package buildkite

import (
	"context"
	"testing"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJobLogs(t *testing.T) {
	// Test the tool definition
	t.Run("ToolDefinition", func(t *testing.T) {
		tool, _ := GetJobLogs(context.Background(), nil)

		assert.Equal(t, "get_job_logs", tool.Name)
		assert.Contains(t, tool.Description, "Get the logs of a job")
	})

	t.Run("MissingParameters", func(t *testing.T) {
		assert := require.New(t)
		_, handler := GetJobLogs(context.Background(), &buildkite.Client{})

		// Test missing org parameter
		req := createMCPRequest(t, map[string]any{
			"pipeline_slug": "test-pipeline",
			"build_number":  "123",
			"job_uuid":      "job-123",
		})
		result, err := handler(context.Background(), req)
		assert.NoError(err)
		assert.NotNil(result)
		assert.NotEmpty(result.Content)

		// Test missing pipeline_slug parameter
		req = createMCPRequest(t, map[string]any{
			"org":          "test-org",
			"build_number": "123",
			"job_uuid":     "job-123",
		})
		result, err = handler(context.Background(), req)
		assert.NoError(err)
		assert.NotNil(result)
		assert.NotEmpty(result.Content)

		// Test missing build_number parameter
		req = createMCPRequest(t, map[string]any{
			"org":           "test-org",
			"pipeline_slug": "test-pipeline",
			"job_uuid":      "job-123",
		})
		result, err = handler(context.Background(), req)
		assert.NoError(err)
		assert.NotNil(result)
		assert.NotEmpty(result.Content)

		// Test missing job_uuid parameter
		req = createMCPRequest(t, map[string]any{
			"org":           "test-org",
			"pipeline_slug": "test-pipeline",
			"build_number":  "123",
		})
		result, err = handler(context.Background(), req)
		assert.NoError(err)
		assert.NotNil(result)
		assert.NotEmpty(result.Content)
	})
}
