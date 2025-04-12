package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/buildkite/go-buildkite/v4"
	"github.com/wolfeidau/buildkite-mcp-server/internal/commands"
	"github.com/wolfeidau/buildkite-mcp-server/pkg/applog"
)

const (
	appName = "buildkite-mcp-server"
)

var (
	version = "dev"

	cli struct {
		Stdio    commands.StdioCmd `cmd:"" help:"stdio mcp server."`
		APIToken string            `help:"The Buildkite API token to use." env:"BUILDKITE_API_TOKEN"`
		Debug    bool              `help:"Enable debug mode."`
		Version  kong.VersionFlag
	}
)

func main() {
	ctx := context.Background()

	logger, err := applog.NewTextLogger(
		appName,
		"{{hostname}}_{{username}}_{{timestamp:2006-01-02}}_pid{{pid}}.log",
		&slog.HandlerOptions{AddSource: true},
	)
	if err != nil {
		fmt.Println("failed to create logger:", err)
		os.Exit(1)
	}

	// add logger to context
	ctx = applog.WithLogger(ctx, logger.Logger)

	cmd := kong.Parse(&cli,
		kong.Name(appName),
		kong.Description("A server that proxies requests to the Buildkite API."),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.BindTo(ctx, (*context.Context)(nil)),
	)

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(cli.APIToken))
	if err != nil {
		logger.Error("failed to create buildkite client", "error", err)
		os.Exit(1)
	}

	err = cmd.Run(&commands.Globals{Debug: cli.Debug, Version: version, Client: client})
	cmd.FatalIfErrorf(err)
}
