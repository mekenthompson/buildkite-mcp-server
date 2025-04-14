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

type OrganizationsClient interface {
	List(ctx context.Context, options *buildkite.OrganizationListOptions) ([]buildkite.Organization, *buildkite.Response, error)
}

func UserTokenOrganization(ctx context.Context, client OrganizationsClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("user_token_organization",
			mcp.WithDescription("Get the organization associated with the user token used for this request"),
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Ctx(ctx).Debug().Msg("Getting current user token organization")

			orgs, resp, err := client.List(ctx, &buildkite.OrganizationListOptions{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if resp.StatusCode != 200 {
				return mcp.NewToolResultError("failed to get current user organizations"), nil
			}

			if len(orgs) == 0 {
				return mcp.NewToolResultError("no organization found for the current user token"), nil
			}

			r, err := json.Marshal(&orgs[0])
			if err != nil {
				return nil, fmt.Errorf("failed to marshal user organizations: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
