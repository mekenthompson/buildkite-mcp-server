package buildkite

import (
	"testing"

	"github.com/buildkite/go-buildkite/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func Test_optionalPaginationParams(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]any
		expected  buildkite.ListOptions
		expectErr bool
	}{
		{
			name: "valid pagination parameters",
			args: map[string]any{
				"page":    float64(1),
				"perPage": float64(31),
			},
			expected: buildkite.ListOptions{
				Page:    1,
				PerPage: 31,
			},
			expectErr: false,
		},
		{
			name: "missing pagination parameters should use new defaults (1 per page)",
			args: map[string]any{},
			expected: buildkite.ListOptions{
				Page:    1,
				PerPage: 1,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			req := createMCPRequest(t, tt.args)

			opts, err := optionalPaginationParams(req)
			if tt.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, opts)
			}
		})
	}
}

func createMCPRequest(t *testing.T, args map[string]any) mcp.CallToolRequest {
	t.Helper()
	return mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Arguments: args,
		},
	}
}

func getTextResult(t *testing.T, result *mcp.CallToolResult) mcp.TextContent {
	t.Helper()
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Error("expected text content")
		return mcp.TextContent{}
	}

	return textContent
}
