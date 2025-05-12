#!/usr/bin/env node

import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
import fs from 'fs';
import dotenv from 'dotenv';

// Load environment variables from .env file if present
dotenv.config();

// Get the current script directory
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const packageJsonPath = resolve(__dirname, '../package.json');
const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));

// Buildkite API client
class BuildkiteClient {
  constructor(apiToken, baseUrl = 'https://api.buildkite.com/v2') {
    this.apiToken = apiToken;
    this.baseUrl = baseUrl;
    this.userAgent = `buildkite-mcp-server/${packageJson.version}`;
  }

  async request(path, options = {}) {
    const url = `${this.baseUrl}${path}`;
    const headers = {
      'Authorization': `Bearer ${this.apiToken}`,
      'User-Agent': this.userAgent,
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const { default: fetch } = await import('node-fetch');
    
    try {
      const response = await fetch(url, { ...options, headers });
      
      if (!response.ok) {
        const text = await response.text();
        throw new Error(`Buildkite API error: ${response.status} ${response.statusText} - ${text}`);
      }
      
      return response.json();
    } catch (error) {
      console.error('API request failed:', error.message);
      throw error;
    }
  }

  // API methods
  async getCurrentUser() {
    return this.request('/user');
  }

  async getOrganizations() {
    return this.request('/organizations');
  }

  async getPipelines(org) {
    return this.request(`/organizations/${org}/pipelines`);
  }

  async getPipeline(org, pipeline) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}`);
  }

  async getBuilds(org, pipeline) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}/builds`);
  }
  
  async getBuild(org, pipeline, buildNumber) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}/builds/${buildNumber}`);
  }
  
  async getJobLogs(org, pipeline, buildNumber, jobId) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}/builds/${buildNumber}/jobs/${jobId}/log`);
  }
  
  async getArtifacts(org, pipeline, buildNumber, jobId) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}/builds/${buildNumber}/jobs/${jobId}/artifacts`);
  }
  
  async getArtifact(org, pipeline, buildNumber, jobId, artifactId) {
    return this.request(`/organizations/${org}/pipelines/${pipeline}/builds/${buildNumber}/jobs/${jobId}/artifacts/${artifactId}`);
  }
  
  async getAccessToken() {
    return this.request('/access-token');
  }
}

async function main() {
  console.error(`Buildkite MCP Server v${packageJson.version}`);
  
  // Check for API token
  const apiToken = process.env.BUILDKITE_API_TOKEN;
  if (!apiToken) {
    console.error('Error: BUILDKITE_API_TOKEN environment variable is required');
    process.exit(1);
  }
  
  // Create API client
  const client = new BuildkiteClient(apiToken);
  
  // Create server
  const server = new McpServer({
    name: 'buildkite-mcp-server',
    version: packageJson.version
  });
  
  // Register tools
  
  // User tools
  server.tool(
    'current_user',
    'Get details of the current user in Buildkite',
    {},
    async () => {
      try {
        const result = await client.getCurrentUser();
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Organization tools
  server.tool(
    'user_token_organization',
    'Get the organization associated with the user token used for this request',
    {},
    async () => {
      try {
        const result = await client.getOrganizations();
        const org = result[0] || null;
        return {
          content: [{ type: 'text', text: JSON.stringify(org, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Pipeline tools
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
        const result = await client.getPipelines(args.org);
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
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
        const result = await client.getPipeline(args.org, args.pipeline_slug);
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Build tools
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
        const result = await client.getBuilds(args.org, args.pipeline_slug);
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
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
        const result = await client.getBuild(args.org, args.pipeline_slug, args.build_number);
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Job logs tool
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
        const result = await client.getJobLogs(args.org, args.pipeline_slug, args.build_number, args.job_id);
        
        // Handle different log response formats
        let logText = '';
        if (typeof result === 'string') {
          logText = result;
        } else if (result && result.content) {
          logText = result.content;
        } else {
          logText = JSON.stringify(result, null, 2);
        }
        
        return {
          content: [{ type: 'text', text: logText }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Artifact tools
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
      }
    },
    async (args) => {
      try {
        const result = await client.getArtifacts(args.org, args.pipeline_slug, args.build_number, args.job_id);
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
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
        const result = await client.getArtifact(
          args.org, 
          args.pipeline_slug, 
          args.build_number, 
          args.job_id, 
          args.artifact_id
        );
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Access token tool
  server.tool(
    'access_token',
    'Get details about the current Buildkite API token',
    {},
    async () => {
      try {
        const result = await client.getAccessToken();
        return {
          content: [{ type: 'text', text: JSON.stringify(result, null, 2) }]
        };
      } catch (error) {
        return {
          content: [{ type: 'text', text: `Error: ${error.message}` }],
          isError: true
        };
      }
    }
  );
  
  // Add a prompt for user_token_organization
  server.prompt(
    'user_token_organization_prompt',
    "When asked for detail of a user's pipelines, start by looking up the user's token organization",
    {},
    async () => {
      try {
        const orgs = await client.getOrganizations();
        const org = orgs[0];
        
        return {
          messages: [{
            role: 'assistant',
            content: {
              type: 'text',
              text: `I'll help you find information about your Buildkite pipelines. First, I'll look up your organization.\n\nI found that your user token is associated with the "${org.name}" organization (slug: ${org.slug}).`
            }
          }]
        };
      } catch (error) {
        return {
          messages: [{
            role: 'assistant',
            content: {
              type: 'text',
              text: `I'll help you find information about your Buildkite pipelines, but I'm having trouble retrieving your organization information. Could you please specify which organization you're interested in?`
            }
          }]
        };
      }
    }
  );
  
  // Start server with stdio transport
  console.error('Starting server...');
  const transport = new StdioServerTransport();
  await server.connect(transport);
  
  // Handle process termination
  process.on('SIGINT', async () => {
    console.error('Shutting down...');
    process.exit(0);
  });
}

main().catch(error => {
  console.error('Error:', error);
  process.exit(1);
});