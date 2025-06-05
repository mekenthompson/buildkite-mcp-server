package buildkite

import (
	"context"
	"net/http"
	"testing"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJobs(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildGetOptions) (buildkite.Build, *buildkite.Response, error) {
			// Create a build with various job states
			return buildkite.Build{
					ID:        "123",
					Number:    1,
					State:     "finished",
					CreatedAt: &buildkite.Timestamp{},
					Jobs: []buildkite.Job{
						{ID: "job1", State: "passed"},
						{ID: "job2", State: "failed"},
						{ID: "job3", State: "running"},
						{ID: "job4", State: "waiting"},
						{ID: "job5", State: "passed"},
						{ID: "job6", State: "canceled"},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := GetJobs(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test getting all jobs (no filter)
	requestAll := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
	})
	resultAll, err := handler(ctx, requestAll)
	assert.NoError(err)

	textContentAll := getTextResult(t, resultAll)
	assert.Contains(textContentAll.Text, `"job1"`)
	assert.Contains(textContentAll.Text, `"job2"`)
	assert.Contains(textContentAll.Text, `"job3"`)
	assert.Contains(textContentAll.Text, `"job4"`)
	assert.Contains(textContentAll.Text, `"job5"`)
	assert.Contains(textContentAll.Text, `"job6"`)
}

func TestGetJobsWithStateFilter(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildGetOptions) (buildkite.Build, *buildkite.Response, error) {
			// Create a build with various job states
			return buildkite.Build{
					ID:        "123",
					Number:    1,
					State:     "finished",
					CreatedAt: &buildkite.Timestamp{},
					Jobs: []buildkite.Job{
						{ID: "job1", State: "passed"},
						{ID: "job2", State: "failed"},
						{ID: "job3", State: "running"},
						{ID: "job4", State: "waiting"},
						{ID: "job5", State: "passed"},
						{ID: "job6", State: "canceled"},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := GetJobs(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test filtering by "passed" state
	requestPassed := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "passed",
	})
	resultPassed, err := handler(ctx, requestPassed)
	assert.NoError(err)

	textContentPassed := getTextResult(t, resultPassed)
	// Only filtered jobs are returned
	assert.Contains(textContentPassed.Text, `"job1"`)
	assert.Contains(textContentPassed.Text, `"job5"`)
	assert.NotContains(textContentPassed.Text, `"job2"`)
	assert.NotContains(textContentPassed.Text, `"job3"`)
	assert.NotContains(textContentPassed.Text, `"job4"`)
	assert.NotContains(textContentPassed.Text, `"job6"`)

	// Test filtering by "failed" state
	requestFailed := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "failed",
	})
	resultFailed, err := handler(ctx, requestFailed)
	assert.NoError(err)

	textContentFailed := getTextResult(t, resultFailed)
	// Only filtered jobs are returned
	assert.Contains(textContentFailed.Text, `"job2"`)
	assert.NotContains(textContentFailed.Text, `"job1"`)
	assert.NotContains(textContentFailed.Text, `"job3"`)
	assert.NotContains(textContentFailed.Text, `"job4"`)
	assert.NotContains(textContentFailed.Text, `"job5"`)
	assert.NotContains(textContentFailed.Text, `"job6"`)

	// Test filtering by "running" state
	requestRunning := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "running",
	})
	resultRunning, err := handler(ctx, requestRunning)
	assert.NoError(err)

	textContentRunning := getTextResult(t, resultRunning)
	// Only filtered jobs are returned
	assert.Contains(textContentRunning.Text, `"job3"`)
	assert.NotContains(textContentRunning.Text, `"job1"`)
	assert.NotContains(textContentRunning.Text, `"job2"`)
	assert.NotContains(textContentRunning.Text, `"job4"`)
	assert.NotContains(textContentRunning.Text, `"job5"`)
	assert.NotContains(textContentRunning.Text, `"job6"`)
}

func TestGetJobsMissingParameters(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{}

	tool, handler := GetJobs(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test missing org
	requestMissingOrg := createMCPRequest(t, map[string]any{
		"pipeline_slug": "pipeline",
		"build_number":  "1",
	})
	resultMissingOrg, err := handler(ctx, requestMissingOrg)
	assert.NoError(err)
	assert.NotNil(resultMissingOrg)
	assert.Len(resultMissingOrg.Content, 1)
	errorContent, ok := resultMissingOrg.Content[0].(mcp.TextContent)
	assert.True(ok)
	assert.Contains(errorContent.Text, "org")

	// Test missing pipeline_slug
	requestMissingPipeline := createMCPRequest(t, map[string]any{
		"org":          "org",
		"build_number": "1",
	})
	resultMissingPipeline, err := handler(ctx, requestMissingPipeline)
	assert.NoError(err)
	assert.NotNil(resultMissingPipeline)
	assert.Len(resultMissingPipeline.Content, 1)
	errorContent, ok = resultMissingPipeline.Content[0].(mcp.TextContent)
	assert.True(ok)
	assert.Contains(errorContent.Text, "pipeline_slug")

	// Test missing build_number
	requestMissingBuild := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
	})
	resultMissingBuild, err := handler(ctx, requestMissingBuild)
	assert.NoError(err)
	assert.NotNil(resultMissingBuild)
	assert.Len(resultMissingBuild.Content, 1)
	errorContent, ok = resultMissingBuild.Content[0].(mcp.TextContent)
	assert.True(ok)
	assert.Contains(errorContent.Text, "build_number")
}

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