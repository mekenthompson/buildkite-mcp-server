import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Registers build-related tools with the MCP server
 */
export function registerBuildTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'list_builds',
    'List all builds in a pipeline in Buildkite',
    {
      org: {
        type: 'string',
        description: 'The organization slug for the owner of the pipeline',
        required: true
      },
      pipeline_slug: {
        type: 'string',
        description: 'The slug of the pipeline',
        required: true
      },
      page: {
        type: 'number',
        description: 'Page number for paginated results',
        default: 1
      },
      per_page: {
        type: 'number',
        description: 'Number of items per page (max 100)',
        default: 30
      }
    },
    async (args) => {
      try {
        const { org, pipeline_slug, page, per_page } = args;
        
        const result = await client.listBuilds(org, pipeline_slug, page, per_page);
        
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
            text: `Error listing builds: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );

  server.tool(
    'get_build',
    'Get details of a specific build in Buildkite',
    {
      org: {
        type: 'string',
        description: 'The organization slug for the owner of the pipeline',
        required: true
      },
      pipeline_slug: {
        type: 'string',
        description: 'The slug of the pipeline',
        required: true
      },
      build_number: {
        type: 'string',
        description: 'The build number to retrieve',
        required: true
      }
    },
    async (args) => {
      try {
        const { org, pipeline_slug, build_number } = args;
        
        const result = await client.getBuild(org, pipeline_slug, build_number);
        
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
            text: `Error getting build: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}