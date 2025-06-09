package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/buildkite/buildkite-mcp-server/internal/trace"
	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
)

type BuildsClient interface {
	Get(ctx context.Context, org, pipelineSlug, buildNumber string, options *buildkite.BuildGetOptions) (buildkite.Build, *buildkite.Response, error)
	ListByPipeline(ctx context.Context, org, pipelineSlug string, options *buildkite.BuildsListOptions) ([]buildkite.Build, *buildkite.Response, error)
}

// JobSummary represents a summary of jobs grouped by state, with finished jobs classified as passed/failed
type JobSummary struct {
	Total   int            `json:"total"`
	ByState map[string]int `json:"by_state"`
}

// BuildWithSummary represents a build with job summary and optionally full job details
type BuildWithSummary struct {
	buildkite.Build
	JobSummary *JobSummary `json:"job_summary"`
}

func ListBuilds(ctx context.Context, client BuildsClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_builds",
			mcp.WithDescription("List all builds in a pipeline in Buildkite"),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline"),
			),
			mcp.WithString("pipeline_slug",
				mcp.Required(),
				mcp.Description("The slug of the pipeline"),
			),
			mcp.WithString("branch",
				mcp.Description("Filter builds by git branch name"),
			),
			withPagination(),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "List Builds",
				ReadOnlyHint: mcp.ToBoolPtr(true),
			}),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.ListBuilds")
			defer span.End()

			org, err := request.RequireString("org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			pipelineSlug, err := request.RequireString("pipeline_slug")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			branch := request.GetString("branch", "")

			paginationParams, err := optionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			span.SetAttributes(
				attribute.String("org", org),
				attribute.String("pipeline_slug", pipelineSlug),
				attribute.String("branch", branch),
				attribute.Int("page", paginationParams.Page),
				attribute.Int("per_page", paginationParams.PerPage),
			)

			options := &buildkite.BuildsListOptions{
				ExcludeJobs:     true,
				ExcludePipeline: true,
				ListOptions:     paginationParams,
			}
			if branch != "" {
				options.Branch = []string{branch}
			}

			builds, resp, err := client.ListByPipeline(ctx, org, pipelineSlug, options)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get issue: %s", string(body))), nil
			}

			result := PaginatedResult[buildkite.Build]{
				Items: builds,
				Headers: map[string]string{
					"Link": resp.Header.Get("Link"),
				},
			}

			r, err := json.Marshal(&result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal builds: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func GetBuildTestEngineRuns(ctx context.Context, client BuildsClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_build_test_engine_runs",
			mcp.WithDescription("Get test engine runs data for a specific build in Buildkite. This can be used to look up Test Runs."),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline"),
			),
			mcp.WithString("pipeline_slug",
				mcp.Required(),
				mcp.Description("The slug of the pipeline"),
			),
			mcp.WithString("build_number",
				mcp.Required(),
				mcp.Description("The number of the build"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get Build Test Engine Runs",
				ReadOnlyHint: mcp.ToBoolPtr(true),
			}),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.GetBuildTestEngineRuns")
			defer span.End()

			org, err := request.RequireString("org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			pipelineSlug, err := request.RequireString("pipeline_slug")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			buildNumber, err := request.RequireString("build_number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			span.SetAttributes(
				attribute.String("org", org),
				attribute.String("pipeline_slug", pipelineSlug),
				attribute.String("build_number", buildNumber),
			)

			build, resp, err := client.Get(ctx, org, pipelineSlug, buildNumber, &buildkite.BuildGetOptions{
				IncludeTestEngine: true,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get build: %s", string(body))), nil
			}

			// Extract just the test engine runs data
			var testEngineRuns []buildkite.TestEngineRun
			if build.TestEngine != nil {
				testEngineRuns = build.TestEngine.Runs
			}

			r, err := json.Marshal(&testEngineRuns)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal test engine runs: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func GetBuild(ctx context.Context, client BuildsClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_build",
			mcp.WithDescription("Get a build in Buildkite. Includes job summary with counts by state."),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline"),
			),
			mcp.WithString("pipeline_slug",
				mcp.Required(),
				mcp.Description("The slug of the pipeline"),
			),
			mcp.WithString("build_number",
				mcp.Required(),
				mcp.Description("The number of the build"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get Build",
				ReadOnlyHint: mcp.ToBoolPtr(true),
			}),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.GetBuild")
			defer span.End()

			org, err := request.RequireString("org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			pipelineSlug, err := request.RequireString("pipeline_slug")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			buildNumber, err := request.RequireString("build_number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			span.SetAttributes(
				attribute.String("org", org),
				attribute.String("pipeline_slug", pipelineSlug),
				attribute.String("build_number", buildNumber),
			)

			build, resp, err := client.Get(ctx, org, pipelineSlug, buildNumber, &buildkite.BuildGetOptions{
				IncludeTestEngine: true,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get build: %s", string(body))), nil
			}

			// Create job summary
			jobSummary := &JobSummary{
				Total:   len(build.Jobs),
				ByState: make(map[string]int),
			}

			for _, job := range build.Jobs {
				if job.State == "" {
					continue
				}

				jobSummary.ByState[job.State]++
			}

			// Create build with summary - always exclude job details
			buildWithSummary := BuildWithSummary{
				Build:      build,
				JobSummary: jobSummary,
			}
			buildWithSummary.Jobs = nil

			r, err := json.Marshal(&buildWithSummary)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal build: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
