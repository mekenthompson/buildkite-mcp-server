import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


/**
 * Process and format job logs for better readability
 */
function processJobLogs(logs: any): string {
  // If logs are already in string format, return them directly
  if (typeof logs === 'string') {
    return logs;
  }

  // If logs are in an object with content field, extract that
  if (logs && typeof logs === 'object' && logs.content) {
    return logs.content;
  }

  // Otherwise, stringify the logs object
  return JSON.stringify(logs, null, 2);
}

/**
 * Registers job log-related tools with the MCP server
 */
export function registerJobTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'get_job_logs',
    'Get logs for a specific job in Buildkite',
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
        description: 'The build number containing the job',
        required: true
      },
      job_id: {
        type: 'string',
        description: 'The ID of the job to get logs for',
        required: true
      }
    },
    async (args) => {
      try {
        const { org, pipeline_slug, build_number, job_id } = args;
        
        const result = await client.getJobLogs(org, pipeline_slug, build_number, job_id);
        
        // Process logs to a readable format
        const processedLogs = processJobLogs(result);
        
        return {
          content: [{ 
            type: 'text', 
            text: processedLogs
          }]
        };
      } catch (error) {
        return {
          content: [{ 
            type: 'text', 
            text: `Error retrieving job logs: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}