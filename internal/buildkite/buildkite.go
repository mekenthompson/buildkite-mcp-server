package buildkite

import (
	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
)

func optionalPaginationParams(r mcp.CallToolRequest) (buildkite.ListOptions, error) {
	page := r.GetInt("page", 1)
	perPage := r.GetInt("perPage", 1)
	return buildkite.ListOptions{
		Page:    page,
		PerPage: perPage,
	}, nil
}

func withPagination() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (min 1)"),
			mcp.Min(1),
		)(tool)

		mcp.WithNumber("perPage",
			mcp.Description("Results per page for pagination (min 1, max 100)"),
			mcp.Min(1),
			mcp.Max(100),
		)(tool)
	}
}
