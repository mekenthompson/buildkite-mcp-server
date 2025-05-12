# Buildkite MCP Server (TypeScript)

A TypeScript implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server for Buildkit using TypeScript and the official SDK. For the Go implementation, see [buildkite-mcp](https://github.com/buildkite/buildkite-mcp).

## Overview

This TypeScript implementation provides all the functionality of the original Go server, enabling LLM integrations to access Buildkite data through a standardized MCP interface. 

## Features

- Full MCP protocol implementation using the official TypeScript SDK
- All the same tools as the Go implementation
- TypeScript typing for better IDE support and type safety
- Modern ES module structure

## Installation & Usage

Install globally or use directly with npx:

```bash
npm install -g @buildkite/buildkite-mcp   # global install (optional)
npx @buildkite/buildkite-mcp              # or run directly
```

Or clone and build from source:

```bash
git clone https://github.com/buildkite/buildkite-mcp-server.git
cd buildkite-mcp-server/typescript
npm install && npm run build
```


## Usage

### Using NPX

The easiest way to use this MCP server is via NPX:

```bash
# Set your Buildkite API token
export BUILDKITE_API_TOKEN=your_token_here

# Install globally (optional)
npm install -g @buildkite/buildkite-mcp

# Or run directly with npx
npx @buildkite/buildkite-mcp
```

### As a Library

You can also use the TypeScript source as a library in your project:

```ts
import { createBuildkiteClient } from '@buildkite/buildkite-mcp';
// ... use the SDK programmatically
```

### As a stdio MCP server

```bash
# Set your Buildkite API token
export BUILDKITE_API_TOKEN=your_token_here

# Run the server
npm run stdio
```

## Tools

| Tool Name | Description |
|-----------|-------------|
| `get_pipeline` | Get details of a specific pipeline in Buildkite |
| `list_pipelines` | List all pipelines in a buildkite organization |
| `list_builds` | List all builds in a pipeline in Buildkite |
| `get_build` | Get details of a specific build |
| `get_job_logs` | Get logs for a specific job in Buildkite |
| `list_artifacts` | List all artifacts for a specific job in Buildkite |
| `get_artifact` | Get a specific artifact for a specific job in Buildkite |
| `current_user` | Get details of the current user in Buildkite |
| `user_token_organization` | Get the organization associated with the user token |

## Configuration

The server reads configuration from environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `BUILDKITE_API_TOKEN` | Your Buildkite API token | Yes |
| `BUILDKITE_API_BASE_URL` | Custom API URL (defaults to https://api.buildkite.com/v2) | No |

## Development

```bash
# Run in development mode with hot reload
npm run dev

# Type-check the project
npm run tsc -- --noEmit

# Build the project
npm run build
```

## License

This project is released under MIT license.