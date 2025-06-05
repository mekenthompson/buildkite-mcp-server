package buildkite

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

type MockTestExecutionsClient struct {
	GetFailedExecutionsFunc func(ctx context.Context, org, slug, runID string, opt *buildkite.FailedExecutionsOptions) ([]buildkite.FailedExecution, *buildkite.Response, error)
}

func (m *MockTestExecutionsClient) GetFailedExecutions(ctx context.Context, org, slug, runID string, opt *buildkite.FailedExecutionsOptions) ([]buildkite.FailedExecution, *buildkite.Response, error) {
	if m.GetFailedExecutionsFunc != nil {
		return m.GetFailedExecutionsFunc(ctx, org, slug, runID, opt)
	}
	return nil, nil, nil
}

var _ TestExecutionsClient = (*MockTestExecutionsClient)(nil)

func TestGetFailedExecutions(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	failedExecutions := []buildkite.FailedExecution{
		{
			ExecutionID:   "exec-1",
			RunID:         "run-123",
			TestID:        "test-456",
			TestName:      "Test Case 1",
			FailureReason: "Assertion failed",
			Duration:      1.5,
		},
		{
			ExecutionID:   "exec-2",
			RunID:         "run-123",
			TestID:        "test-789",
			TestName:      "Test Case 2",
			FailureReason: "Timeout",
			Duration:      30.0,
		},
	}

	mockClient := &MockTestExecutionsClient{
		GetFailedExecutionsFunc: func(ctx context.Context, org, slug, runID string, opt *buildkite.FailedExecutionsOptions) ([]buildkite.FailedExecution, *buildkite.Response, error) {
			return failedExecutions, &buildkite.Response{
				Response: &http.Response{
					StatusCode: http.StatusOK,
				},
			}, nil
		},
	}

	tool, handler := GetFailedTestExecutions(ctx, mockClient)

	// Test tool properties
	assert.Equal("get_failed_executions", tool.Name)
	assert.Equal("Get failed test executions for a specific test run in Buildkite Test Engine", tool.Description)
	assert.True(*tool.Annotations.ReadOnlyHint)

	// Test successful request
	request := createMCPRequest(t, map[string]any{
		"org":                      "org",
		"test_suite_slug":          "suite1",
		"run_id":                   "run1",
		"include_failure_expanded": true,
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.NotNil(result)

	// Check the result contains failed execution data
	textContent := result.Content[0].(mcp.TextContent)
	assert.Contains(textContent.Text, "exec-1")
	assert.Contains(textContent.Text, "exec-2")
	assert.Contains(textContent.Text, "Test Case 1")
	assert.Contains(textContent.Text, "Assertion failed")
	assert.Contains(textContent.Text, "Timeout")
}

func TestGetFailedExecutionsMissingOrg(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	mockClient := &MockTestExecutionsClient{}

	_, handler := GetFailedTestExecutions(ctx, mockClient)

	request := createMCPRequest(t, map[string]any{
		"test_suite_slug": "suite1",
		"run_id":          "run1",
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.True(result.IsError)
	assert.Contains(result.Content[0].(mcp.TextContent).Text, "org")
}

func TestGetFailedExecutionsMissingTestSuiteSlug(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	mockClient := &MockTestExecutionsClient{}

	_, handler := GetFailedTestExecutions(ctx, mockClient)

	request := createMCPRequest(t, map[string]any{
		"org":    "org",
		"run_id": "run1",
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.True(result.IsError)
	assert.Contains(result.Content[0].(mcp.TextContent).Text, "test_suite_slug")
}

func TestGetFailedExecutionsMissingRunID(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	mockClient := &MockTestExecutionsClient{}

	_, handler := GetFailedTestExecutions(ctx, mockClient)

	request := createMCPRequest(t, map[string]any{
		"org":             "org",
		"test_suite_slug": "suite1",
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.True(result.IsError)
	assert.Contains(result.Content[0].(mcp.TextContent).Text, "run_id")
}

func TestGetFailedExecutionsWithError(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	mockClient := &MockTestExecutionsClient{
		GetFailedExecutionsFunc: func(ctx context.Context, org, slug, runID string, opt *buildkite.FailedExecutionsOptions) ([]buildkite.FailedExecution, *buildkite.Response, error) {
			return []buildkite.FailedExecution{}, &buildkite.Response{}, fmt.Errorf("API error")
		},
	}

	_, handler := GetFailedTestExecutions(ctx, mockClient)

	request := createMCPRequest(t, map[string]any{
		"org":             "org",
		"test_suite_slug": "suite1",
		"run_id":          "run1",
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.True(result.IsError)
	assert.Contains(result.Content[0].(mcp.TextContent).Text, "API error")
}

func TestGetFailedExecutionsHTTPError(t *testing.T) {
	assert := require.New(t)

	ctx := context.Background()
	mockClient := &MockTestExecutionsClient{
		GetFailedExecutionsFunc: func(ctx context.Context, org, slug, runID string, opt *buildkite.FailedExecutionsOptions) ([]buildkite.FailedExecution, *buildkite.Response, error) {
			return []buildkite.FailedExecution{}, &buildkite.Response{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("Failed executions not found")),
				},
			}, nil
		},
	}

	_, handler := GetFailedTestExecutions(ctx, mockClient)

	request := createMCPRequest(t, map[string]any{
		"org":             "org",
		"test_suite_slug": "suite1",
		"run_id":          "run1",
	})

	result, err := handler(ctx, request)
	assert.NoError(err)
	assert.True(result.IsError)
	assert.Contains(result.Content[0].(mcp.TextContent).Text, "Failed executions not found")
}