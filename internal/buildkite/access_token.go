package buildkite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type AccessTokenClient interface {
	Get(ctx context.Context) (buildkite.AccessToken, *buildkite.Response, error)
}

func AccessToken(ctx context.Context, client AccessTokenClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("access_token",
			mcp.WithDescription("Get the details for the API access token that was used to authenticate the request"),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Ctx(ctx).Debug().Msg("Getting access token")

			token, resp, err := client.Get(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if resp.StatusCode != 200 {
				return mcp.NewToolResultError("failed to get access token"), nil
			}

			r, err := json.Marshal(&token)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal token: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
