import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Registers organization-related tools with the MCP server
 */
export function registerOrganizationTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'user_token_organization',
    'Get the organization associated with the user token used for this request',
    {},
    async () => {
      try {
        const result = await client.getUserTokenOrganization();
        
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
            text: `Error retrieving organization for user token: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );

  server.prompt(
    'user_token_organization_prompt',
    'When asked for detail of a user\'s pipelines, start by looking up the user\'s token organization',
    {},
    async () => {
      try {
        const org = await client.getUserTokenOrganization();
        
        return {
          messages: [
            {
              role: 'assistant',
              content: {
                type: 'text',
                text: `I'll help you find information about your Buildkite pipelines. First, I'll look up your organization.\n\nI found that your user token is associated with the "${org.name}" organization (slug: ${org.slug}).`
              }
            }
          ]
        };
      } catch (error) {
        return {
          messages: [
            {
              role: 'assistant',
              content: {
                type: 'text',
                text: `I'll help you find information about your Buildkite pipelines, but I'm having trouble retrieving your organization information. Could you please specify which organization you're interested in?`
              }
            }
          ]
        };
      }
    }
  );
}