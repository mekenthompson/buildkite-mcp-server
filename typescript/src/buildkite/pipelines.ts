import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Registers pipeline-related tools with the MCP server
 */
export function registerPipelineTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'list_pipelines',
    'List all pipelines in a buildkite organization',
    {
      org: {
        type: 'string',
        description: 'The organization slug for the owner of the pipeline',
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
        const { org, page, per_page } = args;
        
        const result = await client.listPipelines(org, page, per_page);
        
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
            text: `Error listing pipelines: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );

  server.tool(
    'get_pipeline',
    'Get details of a specific pipeline in Buildkite',
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
      }
    },
    async (args) => {
      try {
        const { org, pipeline_slug } = args;
        
        const result = await client.getPipeline(org, pipeline_slug);
        
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
            text: `Error getting pipeline: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}