import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { createBuildkiteClient } from './buildkite/client.js';
import { registerPipelineTools } from './buildkite/pipelines.js';
import { registerBuildTools } from './buildkite/builds.js';
import { registerUserTools } from './buildkite/user.js';
import { registerJobTools } from './buildkite/job_logs.js';
import { registerAccessTokenTools } from './buildkite/access_token.js';
import { registerArtifactTools } from './buildkite/artifacts.js';
import { registerOrganizationTools } from './buildkite/organizations.js';

import dotenv from 'dotenv';

// Load environment variables
dotenv.config();

// Configuration
const version = process.env.VERSION || '0.2.0';
const apiToken = process.env.BUILDKITE_API_TOKEN;

if (!apiToken) {
  console.error('BUILDKITE_API_TOKEN environment variable is required');
  process.exit(1);
}


async function main() {
  try {
    console.log(`Starting Buildkite MCP Server (${version})`);
    
    // Create Buildkite client
    const client = createBuildkiteClient(apiToken as string, version);
    
    // Create MCP server
    const server = new McpServer({
      name: 'buildkite-mcp-server',
      version: version
    });
    
    // Register all tools
    registerPipelineTools(server, client);
    registerBuildTools(server, client);
    registerUserTools(server, client);
    registerJobTools(server, client);
    registerAccessTokenTools(server, client);
    registerArtifactTools(server, client);
    registerOrganizationTools(server, client);
    
    // Start server with stdio transport
    const transport = new StdioServerTransport();
    await server.connect(transport);
    
    // Gracefully handle process termination
    process.on('SIGINT', async () => {
      console.log('Shutting down...');
      process.exit(0);
    });
    
  } catch (error) {
    console.error('Error starting server:', error);
    process.exit(1);
  }
}

main();