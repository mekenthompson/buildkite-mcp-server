import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Registers access token-related tools with the MCP server
 */
export function registerAccessTokenTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'access_token',
    'Get details about the current Buildkite API token',
    {},
    async () => {
      try {
        const result = await client.getAccessToken();
        
        return {
          content: [{ 
            type: 'text', 
            text: JSON.stringify(result, null, 2) 
          }]
        };
      } catch (error) {
        return {
          content: [{ 
            type: 'text', 
            text: `Error retrieving access token info: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}