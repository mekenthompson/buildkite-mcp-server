package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/buildkite/buildkite-mcp-server/internal/trace"
	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
)

type PipelinesClient interface {
	Get(ctx context.Context, org, pipelineSlug string) (buildkite.Pipeline, *buildkite.Response, error)
	List(ctx context.Context, org string, options *buildkite.PipelineListOptions) ([]buildkite.Pipeline, *buildkite.Response, error)
	Create(ctx context.Context, org string, p buildkite.CreatePipeline) (buildkite.Pipeline, *buildkite.Response, error)
}

func ListPipelines(ctx context.Context, client PipelinesClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_pipelines",
			mcp.WithDescription("List all pipelines in an organization with their basic details, build counts, and current status"),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline"),
			),
			withPagination(),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "List Pipelines",
				ReadOnlyHint: mcp.ToBoolPtr(true),
			}),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.ListPipelines")
			defer span.End()

			org, err := request.RequireString("org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			paginationParams, err := optionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			span.SetAttributes(
				attribute.String("org", org),
				attribute.Int("page", paginationParams.Page),
				attribute.Int("per_page", paginationParams.PerPage),
			)

			pipelines, resp, err := client.List(ctx, org, &buildkite.PipelineListOptions{
				ListOptions: paginationParams,
			})
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

			result := PaginatedResult[buildkite.Pipeline]{
				Items: pipelines,
				Headers: map[string]string{
					"Link": resp.Header.Get("Link"),
				},
			}

			r, err := json.Marshal(&result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal pipelines: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func GetPipeline(ctx context.Context, client PipelinesClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pipeline",
			mcp.WithDescription("Get detailed information about a specific pipeline including its configuration, steps, environment variables, and build statistics"),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline"),
			),
			mcp.WithString("pipeline_slug",
				mcp.Required(),
				mcp.Description("The slug of the pipeline"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get Pipeline",
				ReadOnlyHint: mcp.ToBoolPtr(true),
			}),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.GetPipeline")
			defer span.End()

			org, err := request.RequireString("org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			pipelineSlug, err := request.RequireString("pipeline_slug")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			span.SetAttributes(
				attribute.String("org", org),
				attribute.String("pipeline_slug", pipelineSlug),
			)

			pipeline, resp, err := client.Get(ctx, org, pipelineSlug)
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

			r, err := json.Marshal(&pipeline)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal issue: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

type CreatePipelineArgs struct {
	OrgSlug         string   `json:"org_slug"`
	Name            string   `json:"name"`
	RepositoryURL   string   `json:"repository_url"`
	ClusterID       string   `json:"cluster_id,omitempty"` // Optional, if not provided, the default cluster will be used
	Description     string   `json:"description,omitempty"`
	Configuration   string   `json:"configuration,omitempty"`
	EnvironmentVars []string `json:"environment_variables,omitempty"`
}

func CreatePipeline(ctx context.Context, client PipelinesClient) (tool mcp.Tool, handler mcp.TypedToolHandlerFunc[CreatePipelineArgs]) {
	return mcp.NewTool("create_pipeline",
			mcp.WithDescription("Create a new pipeline in Buildkite using the provided repository URL. The repository URL must be a valid Git repository URL that is accessible to Buildkite."),
			mcp.WithString("org_slug",
				mcp.Required(),
				mcp.Description("The organization slug for the owner of the pipeline. This is used to determine where to create the pipeline."),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("The name of the pipeline"),
			),
			mcp.WithString("repository_url",
				mcp.Required(),
				mcp.Description("The Git repository URL to use for the pipeline"),
			),
			mcp.WithString("description",
				mcp.Description("The description of the pipeline"),
			),
			mcp.WithString("configuration",
				mcp.Description("The pipeline configuration in YAML format. If not provided, the default configuration will be used."),
			),
			mcp.WithObject("environment_variables",
				mcp.Description("A list of environment variables to set for the pipeline builds. Each variable should be a key-value pair in the format 'KEY=VALUE'."),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Create Pipeline",
				ReadOnlyHint: mcp.ToBoolPtr(false),
			}),
		),
		func(ctx context.Context, request mcp.CallToolRequest, args CreatePipelineArgs) (*mcp.CallToolResult, error) {
			ctx, span := trace.Start(ctx, "buildkite.CreatePipeline")
			defer span.End()

			if args.OrgSlug == "" {
				return mcp.NewToolResultError("org_slug is required"), nil
			}
			if args.Name == "" {
				return mcp.NewToolResultError("name is required"), nil
			}
			if args.RepositoryURL == "" {
				return mcp.NewToolResultError("repository_url is required"), nil
			}

			// parse the URL to ensure it's valid
			if _, err := url.Parse(args.RepositoryURL); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid repository URL: %s", err.Error())), nil
			}

			span.SetAttributes(
				attribute.String("name", args.Name),
				attribute.String("repository_url", args.RepositoryURL),
			)

			create := buildkite.CreatePipeline{
				Name:       args.Name,
				Repository: args.RepositoryURL,
			}

			if args.Configuration != "" {
				create.Configuration = args.Configuration
			}

			// if the description is not empty, set it
			if args.Description != "" {
				create.Description = args.Description
			}

			args.EnvironmentVars = args.EnvironmentVars[:0] // Ensure we start with an empty slice

			pipeline, resp, err := client.Create(ctx, args.OrgSlug, create)
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
			r, err := json.Marshal(&pipeline)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal issue: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
