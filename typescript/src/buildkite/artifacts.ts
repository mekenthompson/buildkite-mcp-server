import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { BuildkiteClient } from './client.js';


// Type for Buildkite artifact
interface BuildkiteArtifact {
  id: string;
  job_id: string;
  url: string;
  download_url: string;
  state: string;
  path: string;
  dirname: string;
  filename: string;
  mime_type: string;
  file_size: number;
  sha1sum: string;
  [key: string]: any;
}

/**
 * Registers artifact-related tools with the MCP server
 */
export function registerArtifactTools(server: McpServer, client: BuildkiteClient) {
  server.tool(
    'list_artifacts',
    'List all artifacts for a specific job in Buildkite',
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
        description: 'The ID of the job to list artifacts for',
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
        const { org, pipeline_slug, build_number, job_id, page, per_page } = args;
        
        const result = await client.listArtifacts(org, pipeline_slug, build_number, job_id, page, per_page);
        
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
            text: `Error listing artifacts: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );

  server.tool(
    'get_artifact',
    'Get a specific artifact for a specific job in Buildkite',
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
        description: 'The ID of the job',
        required: true
      },
      artifact_id: {
        type: 'string',
        description: 'The ID of the artifact to retrieve',
        required: true
      }
    },
    async (args) => {
      try {
        const { org, pipeline_slug, build_number, job_id, artifact_id } = args;
        
        // First get the artifact metadata
        const artifact = await client.getArtifact(org, pipeline_slug, build_number, job_id, artifact_id) as BuildkiteArtifact;
        
        // If the artifact has a download URL, fetch its content
        if (artifact && artifact.download_url) {
          const response = await client.downloadArtifact(artifact.download_url);
          let content;
          
          // Try to get content as text, fallback to JSON description if binary
          try {
            content = await response.text();
          } catch (e) {
            content = `Binary artifact available at: ${artifact.download_url}\n\nMetadata: ${JSON.stringify(artifact, null, 2)}`;
          }
          
          return {
            content: [{ 
              type: 'text', 
              text: content
            }]
          };
        } else {
          return {
            content: [{ 
              type: 'text', 
              text: JSON.stringify(artifact, null, 2)
            }]
          };
        }
      } catch (error) {
        return {
          content: [{ 
            type: 'text', 
            text: `Error retrieving artifact: ${error instanceof Error ? error.message : String(error)}` 
          }],
          isError: true
        };
      }
    }
  );
}