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
	// Should always have pagination metadata (default page size 25)
	assert.Contains(textContentAll.Text, `"page":1`)
	assert.Contains(textContentAll.Text, `"per_page":25`)
	assert.Contains(textContentAll.Text, `"total":6`)
	assert.Contains(textContentAll.Text, `"has_next":false`)
	assert.Contains(textContentAll.Text, `"has_prev":false`)
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
	// Should always have pagination metadata (default page size 25)
	assert.Contains(textContentPassed.Text, `"page":1`)
	assert.Contains(textContentPassed.Text, `"per_page":25`)
	assert.Contains(textContentPassed.Text, `"total":2`)
	assert.Contains(textContentPassed.Text, `"has_next":false`)
	assert.Contains(textContentPassed.Text, `"has_prev":false`)

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
	// Should always have pagination metadata (default page size 25)
	assert.Contains(textContentFailed.Text, `"page":1`)
	assert.Contains(textContentFailed.Text, `"per_page":25`)
	assert.Contains(textContentFailed.Text, `"total":1`)
	assert.Contains(textContentFailed.Text, `"has_next":false`)
	assert.Contains(textContentFailed.Text, `"has_prev":false`)

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
	// Should always have pagination metadata (default page size 25)
	assert.Contains(textContentRunning.Text, `"page":1`)
	assert.Contains(textContentRunning.Text, `"per_page":25`)
	assert.Contains(textContentRunning.Text, `"total":1`)
	assert.Contains(textContentRunning.Text, `"has_next":false`)
	assert.Contains(textContentRunning.Text, `"has_prev":false`)
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

func TestGetJobsPagination(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildGetOptions) (buildkite.Build, *buildkite.Response, error) {
			// Create a build with 6 jobs to test pagination
			return buildkite.Build{
					ID:        "123",
					Number:    1,
					State:     "finished",
					CreatedAt: &buildkite.Timestamp{},
					Jobs: []buildkite.Job{
						{ID: "job1", State: "passed", Agent: buildkite.Agent{ID: "agent1", Name: "test-agent-1"}},
						{ID: "job2", State: "failed", Agent: buildkite.Agent{ID: "agent2", Name: "test-agent-2"}},
						{ID: "job3", State: "running", Agent: buildkite.Agent{ID: "agent3", Name: "test-agent-3"}},
						{ID: "job4", State: "waiting"},
						{ID: "job5", State: "passed", Agent: buildkite.Agent{ID: "agent5", Name: "test-agent-5"}},
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

	// Test first page with page size of 2
	requestFirstPage := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"page":          float64(1),
		"perPage":       float64(2),
	})
	resultFirstPage, err := handler(ctx, requestFirstPage)
	assert.NoError(err)

	textContentFirstPage := getTextResult(t, resultFirstPage)
	// Should contain first 2 jobs
	assert.Contains(textContentFirstPage.Text, `"job1"`)
	assert.Contains(textContentFirstPage.Text, `"job2"`)
	assert.NotContains(textContentFirstPage.Text, `"job3"`)
	assert.NotContains(textContentFirstPage.Text, `"job4"`)
	// Should have pagination metadata
	assert.Contains(textContentFirstPage.Text, `"page":1`)
	assert.Contains(textContentFirstPage.Text, `"per_page":2`)
	assert.Contains(textContentFirstPage.Text, `"total":6`)
	assert.Contains(textContentFirstPage.Text, `"has_next":true`)
	assert.Contains(textContentFirstPage.Text, `"has_prev":false`)

	// Test second page with page size of 2
	requestSecondPage := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"page":          float64(2),
		"perPage":       float64(2),
	})
	resultSecondPage, err := handler(ctx, requestSecondPage)
	assert.NoError(err)

	textContentSecondPage := getTextResult(t, resultSecondPage)
	// Should contain next 2 jobs
	assert.NotContains(textContentSecondPage.Text, `"job1"`)
	assert.NotContains(textContentSecondPage.Text, `"job2"`)
	assert.Contains(textContentSecondPage.Text, `"job3"`)
	assert.Contains(textContentSecondPage.Text, `"job4"`)
	// Should have pagination metadata
	assert.Contains(textContentSecondPage.Text, `"page":2`)
	assert.Contains(textContentSecondPage.Text, `"per_page":2`)
	assert.Contains(textContentSecondPage.Text, `"total":6`)
	assert.Contains(textContentSecondPage.Text, `"has_next":true`)
	assert.Contains(textContentSecondPage.Text, `"has_prev":true`)

	// Test last page
	requestLastPage := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"page":          float64(3),
		"perPage":       float64(2),
	})
	resultLastPage, err := handler(ctx, requestLastPage)
	assert.NoError(err)

	textContentLastPage := getTextResult(t, resultLastPage)
	// Should contain last 2 jobs
	assert.Contains(textContentLastPage.Text, `"job5"`)
	assert.Contains(textContentLastPage.Text, `"job6"`)
	// Should have pagination metadata
	assert.Contains(textContentLastPage.Text, `"page":3`)
	assert.Contains(textContentLastPage.Text, `"per_page":2`)
	assert.Contains(textContentLastPage.Text, `"total":6`)
	assert.Contains(textContentLastPage.Text, `"has_next":false`)
	assert.Contains(textContentLastPage.Text, `"has_prev":true`)

	// Test page beyond available data
	requestBeyond := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"page":          float64(5),
		"perPage":       float64(2),
	})
	resultBeyond, err := handler(ctx, requestBeyond)
	assert.NoError(err)

	textContentBeyond := getTextResult(t, resultBeyond)
	// Should contain empty items array
	assert.Contains(textContentBeyond.Text, `"items":[]`)
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