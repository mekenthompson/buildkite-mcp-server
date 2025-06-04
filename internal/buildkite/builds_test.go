package buildkite

import (
	"context"
	"net/http"
	"testing"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/stretchr/testify/require"
)

type MockBuildsClient struct {
	ListByPipelineFunc func(ctx context.Context, org string, pipeline string, opt *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error)
	GetFunc            func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildsListOptions) (buildkite.Build, *buildkite.Response, error)
}

func (m *MockBuildsClient) Get(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildsListOptions) (buildkite.Build, *buildkite.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, org, pipeline, id, opt)
	}
	return buildkite.Build{}, nil, nil
}

func (m *MockBuildsClient) ListByPipeline(ctx context.Context, org string, pipeline string, opt *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error) {
	if m.ListByPipelineFunc != nil {
		return m.ListByPipelineFunc(ctx, org, pipeline, opt)
	}
	return nil, nil, nil
}

var _ BuildsClient = (*MockBuildsClient)(nil)

func TestGetBuildDefault(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildsListOptions) (buildkite.Build, *buildkite.Response, error) {
			// Return build without jobs
			return buildkite.Build{
					ID:        "123",
					Number:    1,
					State:     "running",
					CreatedAt: &buildkite.Timestamp{},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := GetBuild(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test default behavior - jobs excluded by default, summary always included
	request := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
	})
	result, err := handler(ctx, request)
	assert.NoError(err)

	textContent := getTextResult(t, result)
	assert.Equal(`{"id":"123","number":1,"state":"running","blocked":false,"author":{},"created_at":"0001-01-01T00:00:00Z","creator":{"avatar_url":"","created_at":null,"email":"","id":"","name":""},"job_summary":{"total":0,"by_state":{}}}`, textContent.Text)
}

func TestGetBuildWithJobs(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildsListOptions) (buildkite.Build, *buildkite.Response, error) {
			// Create a build with some jobs to test summary functionality
			return buildkite.Build{
					ID:        "123",
					Number:    1,
					State:     "finished",
					CreatedAt: &buildkite.Timestamp{},
					Jobs: []buildkite.Job{
						{ID: "job1", State: "passed"}, // API already coerced
						{ID: "job2", State: "failed"}, // API already coerced
						{ID: "job3", State: "running"},
						{ID: "job4", State: "waiting"},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := GetBuild(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test default behavior - jobs excluded, only summary shown
	request := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
	})
	result, err := handler(ctx, request)
	assert.NoError(err)

	textContent := getTextResult(t, result)
	assert.Contains(textContent.Text, `"job_summary"`)
	assert.Contains(textContent.Text, `"total":4`)
	assert.Contains(textContent.Text, `"by_state":{"failed":1,"passed":1,"running":1,"waiting":1}`)
	assert.NotContains(textContent.Text, `"jobs"`) // Jobs excluded by default

	// Test with jobs explicitly included
	requestIncludeJobs := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"exclude_jobs":  "false",
	})
	resultIncludeJobs, err := handler(ctx, requestIncludeJobs)
	assert.NoError(err)

	textContentIncludeJobs := getTextResult(t, resultIncludeJobs)
	assert.Contains(textContentIncludeJobs.Text, `"job_summary"`)
	assert.Contains(textContentIncludeJobs.Text, `"total":4`)
	assert.Contains(textContentIncludeJobs.Text, `"jobs"`) // Jobs explicitly included
	assert.Contains(textContentIncludeJobs.Text, `"job1"`)
	assert.Contains(textContentIncludeJobs.Text, `"job2"`)
	assert.Contains(textContentIncludeJobs.Text, `"job3"`)
	assert.Contains(textContentIncludeJobs.Text, `"job4"`)
}

func TestGetBuildWithJobStateFilter(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	client := &MockBuildsClient{
		GetFunc: func(ctx context.Context, org string, pipeline string, id string, opt *buildkite.BuildsListOptions) (buildkite.Build, *buildkite.Response, error) {
			// Create a build with various job states (API already coerced finished jobs)
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

	tool, handler := GetBuild(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test filtering by "passed" pseudo state (default: include jobs when filtering)
	requestPassed := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "passed",
	})
	resultPassed, err := handler(ctx, requestPassed)
	assert.NoError(err)

	textContentPassed := getTextResult(t, resultPassed)
	assert.Contains(textContentPassed.Text, `"job_summary"`)
	assert.Contains(textContentPassed.Text, `"total":6`) // Summary shows all jobs, not filtered
	assert.Contains(textContentPassed.Text, `"by_state":{"canceled":1,"failed":1,"passed":2,"running":1,"waiting":1}`)
	assert.Contains(textContentPassed.Text, `"jobs"`) // Jobs included by default when filtering
	// Only filtered jobs are returned
	assert.Contains(textContentPassed.Text, `"job1"`)
	assert.Contains(textContentPassed.Text, `"job5"`)
	assert.NotContains(textContentPassed.Text, `"job2"`)
	assert.NotContains(textContentPassed.Text, `"job3"`)
	assert.NotContains(textContentPassed.Text, `"job4"`)
	assert.NotContains(textContentPassed.Text, `"job6"`)

	// Test filtering by "failed" pseudo state
	requestFailed := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "failed",
	})
	resultFailed, err := handler(ctx, requestFailed)
	assert.NoError(err)

	textContentFailed := getTextResult(t, resultFailed)
	assert.Contains(textContentFailed.Text, `"job_summary"`)
	assert.Contains(textContentFailed.Text, `"total":6`) // Summary shows all jobs, not filtered
	assert.Contains(textContentFailed.Text, `"by_state":{"canceled":1,"failed":1,"passed":2,"running":1,"waiting":1}`)
	assert.Contains(textContentFailed.Text, `"jobs"`)
	// Only filtered jobs are returned
	assert.Contains(textContentFailed.Text, `"job2"`)
	assert.NotContains(textContentFailed.Text, `"job1"`)
	assert.NotContains(textContentFailed.Text, `"job3"`)
	assert.NotContains(textContentFailed.Text, `"job4"`)
	assert.NotContains(textContentFailed.Text, `"job5"`)
	assert.NotContains(textContentFailed.Text, `"job6"`)

	// Test filtering by actual state "running"
	requestRunning := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "running",
	})
	resultRunning, err := handler(ctx, requestRunning)
	assert.NoError(err)

	textContentRunning := getTextResult(t, resultRunning)
	assert.Contains(textContentRunning.Text, `"job_summary"`)
	assert.Contains(textContentRunning.Text, `"total":6`) // Summary shows all jobs, not filtered
	assert.Contains(textContentRunning.Text, `"by_state":{"canceled":1,"failed":1,"passed":2,"running":1,"waiting":1}`)
	assert.Contains(textContentRunning.Text, `"jobs"`)
	// Only filtered jobs are returned
	assert.Contains(textContentRunning.Text, `"job3"`)
	assert.NotContains(textContentRunning.Text, `"job1"`)
	assert.NotContains(textContentRunning.Text, `"job2"`)
	assert.NotContains(textContentRunning.Text, `"job4"`)
	assert.NotContains(textContentRunning.Text, `"job5"`)
	assert.NotContains(textContentRunning.Text, `"job6"`)

	// Test filtering with jobs explicitly excluded
	requestPassedExcluded := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"build_number":  "1",
		"job_state":     "passed",
		"exclude_jobs":  "true",
	})
	resultPassedExcluded, err := handler(ctx, requestPassedExcluded)
	assert.NoError(err)

	textContentPassedExcluded := getTextResult(t, resultPassedExcluded)
	assert.Contains(textContentPassedExcluded.Text, `"job_summary"`)
	assert.Contains(textContentPassedExcluded.Text, `"total":6`) // Summary shows all jobs, not filtered
	assert.Contains(textContentPassedExcluded.Text, `"by_state":{"canceled":1,"failed":1,"passed":2,"running":1,"waiting":1}`)
	assert.NotContains(textContentPassedExcluded.Text, `"jobs"`) // Jobs explicitly excluded
}

func TestListBuilds(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	var capturedOptions *buildkite.BuildsListOptions
	client := &MockBuildsClient{
		ListByPipelineFunc: func(ctx context.Context, org string, pipeline string, opt *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error) {
			capturedOptions = opt
			return []buildkite.Build{
					{
						ID:        "123",
						Number:    1,
						State:     "running",
						CreatedAt: &buildkite.Timestamp{},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := ListBuilds(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	request := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
	})
	result, err := handler(ctx, request)
	assert.NoError(err)

	textContent := getTextResult(t, result)

	assert.Equal(`{"headers":{"Link":""},"items":[{"id":"123","number":1,"state":"running","blocked":false,"author":{},"created_at":"0001-01-01T00:00:00Z","creator":{"avatar_url":"","created_at":null,"email":"","id":"","name":""}}]}`, textContent.Text)

	// Verify default pagination parameters - ensure they are set to 1 per page
	assert.NotNil(capturedOptions)
	assert.Equal(1, capturedOptions.Page)
	assert.Equal(1, capturedOptions.PerPage)
	assert.Nil(capturedOptions.Branch) // Branch should be nil when not specified
}

func TestListBuildsWithCustomPagination(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	var capturedOptions *buildkite.BuildsListOptions
	client := &MockBuildsClient{
		ListByPipelineFunc: func(ctx context.Context, org string, pipeline string, opt *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error) {
			capturedOptions = opt
			return []buildkite.Build{
					{
						ID:        "123",
						Number:    1,
						State:     "running",
						CreatedAt: &buildkite.Timestamp{},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := ListBuilds(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test with custom pagination parameters
	request := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"page":          float64(3),
		"perPage":       float64(50),
	})
	_, err := handler(ctx, request)
	assert.NoError(err)

	// Verify custom pagination parameters were used
	assert.NotNil(capturedOptions)
	assert.Equal(3, capturedOptions.Page)
	assert.Equal(50, capturedOptions.PerPage)
	assert.Nil(capturedOptions.Branch) // Branch should be nil when not specified
}

func TestListBuildsWithBranchFilter(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	var capturedOptions *buildkite.BuildsListOptions
	client := &MockBuildsClient{
		ListByPipelineFunc: func(ctx context.Context, org string, pipeline string, opt *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error) {
			capturedOptions = opt
			return []buildkite.Build{
					{
						ID:        "123",
						Number:    1,
						State:     "running",
						CreatedAt: &buildkite.Timestamp{},
					},
				}, &buildkite.Response{
					Response: &http.Response{
						StatusCode: 200,
					},
				}, nil
		},
	}

	tool, handler := ListBuilds(ctx, client)
	assert.NotNil(tool)
	assert.NotNil(handler)

	// Test with branch filter
	request := createMCPRequest(t, map[string]any{
		"org":           "org",
		"pipeline_slug": "pipeline",
		"branch":        "main",
	})
	_, err := handler(ctx, request)
	assert.NoError(err)

	// Verify branch filter was applied
	assert.NotNil(capturedOptions)
	assert.Equal([]string{"main"}, capturedOptions.Branch)
	assert.Equal(1, capturedOptions.Page)
	assert.Equal(1, capturedOptions.PerPage)
}
