package commands

import (
	"context"

	"github.com/buildkite/buildkite-mcp-server/internal/buildkite"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type StdioCmd struct{}

func (c *StdioCmd) Run(ctx context.Context, globals *Globals) error {
	s := server.NewMCPServer(
		"github-mcp-server",
		globals.Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging())

	// add the logger to the context
	ctx = globals.Logger.WithContext(ctx)

	log.Ctx(ctx).Info().Str("version", globals.Version).Msg("Starting Buildkite MCP server")

	s.AddTool(buildkite.GetPipeline(ctx, globals.Client.Pipelines))
	s.AddTool(buildkite.ListPipelines(ctx, globals.Client.Pipelines))
	s.AddTool(buildkite.ListBuilds(ctx, globals.Client.Builds))
	s.AddTool(buildkite.GetBuild(ctx, globals.Client.Builds))
	s.AddTool(buildkite.CurrentUser(ctx, globals.Client.User))
	s.AddTool(buildkite.GetJobLogs(ctx, globals.Client))
	s.AddTool(buildkite.AccessToken(ctx, globals.Client.AccessTokens))

	return server.ServeStdio(s)
}
