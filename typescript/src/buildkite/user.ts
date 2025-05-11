import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Registers user-related tools with the MCP server
 */
export function registerUserTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'current_user',
    'Get details of the current user in Buildkite',
    {},
    async () => {
      try {
        const result = await client.getCurrentUser();
        
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
            text: `Error retrieving current user: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}