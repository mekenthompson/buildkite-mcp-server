# buildkite-mcp-server

[![Build status](https://badge.buildkite.com/79fefd75bc7f1898fb35249f7ebd8541a99beef6776e7da1b4.svg?branch=main)](https://buildkite.com/buildkite/buildkite-mcp-server)

This is an [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server for [Buildkite](https://buildkite.com). The goal is to provide access to information from buildkite about pipelines, builds and jobs to tools such as [Claude Desktop](https://claude.ai/download), [GitHub Copilot](https://github.com/features/copilot) and other tools, or editors.

[![Add to Cursor](https://cursor.com/deeplink/mcp-install-dark.png)](https://cursor.com/install-mcp?name=buildkite&config=eyJjb21tYW5kIjoiZG9ja2VyIHJ1biAtaSAtLXJtIC1lIEJVSUxES0lURV9BUElfVE9LRU4gZ2hjci5pby9idWlsZGtpdGUvYnVpbGRraXRlLW1jcC1zZXJ2ZXIgc3RkaW8iLCJlbnYiOnsiQlVJTERLSVRFX0FQSV9UT0tFTiI6ImJrdWFfeHh4eHh4eHgifX0%3D)

# Tools

* `get_cluster` - Get detailed information about a specific cluster including its name, description, default queue, and configuration
* `list_clusters` - List all clusters in an organization with their names, descriptions, default queues, and creation details
* `get_cluster_queue` - Get detailed information about a specific queue including its key, description, dispatch status, and hosted agent configuration
* `list_cluster_queues` - List all queues in a cluster with their keys, descriptions, dispatch status, and agent configuration
* `get_pipeline` - Get detailed information about a specific pipeline including its configuration, steps, environment variables, and build statistics
* `list_pipelines` - List all pipelines in an organization with their basic details, build counts, and current status
* `create_pipeline` - Create a new pipeline in Buildkite using the provided repository URL. The repository URL must be a valid Git repository URL that is accessible to Buildkite.
* `list_builds` - List all builds for a pipeline with their status, commit information, and metadata
* `get_build` - Get detailed information about a specific build including its jobs, timing, and execution details
* `get_build_test_engine_runs` - Get test engine runs data for a specific build in Buildkite. This can be used to look up Test Runs.
* `current_user` - Get details about the user account that owns the API token, including name, email, avatar, and account creation date
* `user_token_organization` - Get the organization associated with the user token used for this request
* `get_jobs` - Get all jobs for a specific build including their state, timing, commands, and execution details
* `get_job_logs` - Get the log output and metadata for a specific job, including content, size, and header timestamps
* `list_artifacts` - List all artifacts for a build across all jobs, including file details, paths, sizes, MIME types, and download URLs
* `get_artifact` - Get detailed information about a specific artifact including its metadata, file size, SHA-1 hash, and download URL
* `list_annotations` - List all annotations for a build, including their context, style (success/info/warning/error), rendered HTML content, and creation timestamps
* `list_test_runs` - List all test runs for a test suite in Buildkite Test Engine
* `get_test_run` - Get a specific test run in Buildkite Test Engine
* `get_failed_executions` - Get failed test executions for a specific test run in Buildkite Test Engine. Optionally get the expanded failure details such as full error messages and stack traces.
* `get_test` - Get a specific test in Buildkite Test Engine. This provides additional metadata for failed test executions
* `access_token` - Get information about the current API access token including its scopes and UUID

Example of the `get_pipeline` tool in action.

![Get Pipeline Tool](docs/images/get_pipeline.png)

### Production

To ensure the MCP server is run in a secure environment, we recommend running it in a container.

Pull the pre-built image (recommended):

```bash
docker pull ghcr.io/buildkite/buildkite-mcp-server
```

Or build it yourself using GoReleaser and copy the binary into your path:

```bash
goreleaser build --snapshot --clean
```

# configuration

Create a [Buildkite API Access Token with read access to pipelines].

## Claude Desktop Configuration

Use this configuration if you want to run the server `buildkite-mcp-server` Docker (recommended):

```json
{
    "mcpServers": {
        "buildkite": {
            "command": "docker",
            "args": [
                "run",
                "-i",
                "--rm",
                "-e",
                "BUILDKITE_API_TOKEN",
                "ghcr.io/buildkite/buildkite-mcp-server",
                "stdio"
            ],
            "env": {
                "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx"
            }
        }
    }
}
```

Configuration if you have `buildkite-mcp-server` installed locally.

```json
{
    "mcpServers": {
        "buildkite": {
            "command": "buildkite-mcp-server",
            "args": [
                "stdio"
            ],
            "env": {
                "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx"
            }
        }
    }
}
```

## Goose Configuration

For Docker with [Goose](https://block.github.io/goose/) (recommended):

```yaml
extensions:
  fetch:
    name: Buildkite
    cmd: docker
    args: ["run", "-i", "--rm", "-e", "BUILDKITE_API_TOKEN", "ghcr.io/buildkite/buildkite-mcp-server", "stdio"]
    enabled: true
    envs: { "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx" }
    type: stdio
    timeout: 300
```

Local configuration for Goose:

```yaml
extensions:
  fetch:
    name: Buildkite
    cmd: buildkite-mcp-server
    args: [stdio]
    enabled: true
    envs: { "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx" }
    type: stdio
    timeout: 300
```

## VSCode Configuration

[VSCode](https://code.visualstudio.com/) supports interactive inputs for variables. To get the API token interactively on MCP startup, put the following in `.vscode/mcp.json`

```json
{
    "inputs": [
        {
            "id": "BUILDKITE_API_TOKEN",
            "type": "promptString",
            "description": "Enter your BuildKite Access Token (https://buildkite.com/user/api-access-tokens)",
            "password": true
        }
    ],
    "servers": {
        "buildkite": {
            "command": "docker",
            "args": [
                "run",
                "-i",
                "--rm",
                "-e",
                "BUILDKITE_API_TOKEN",
                "ghcr.io/buildkite/buildkite-mcp-server",
                "stdio"
            ],
            "env": {
                "BUILDKITE_API_TOKEN": "${input:BUILDKITE_API_TOKEN}"
            }
        }
    }
}
```

## Zed

There is a [Zed](https://zed.dev) editor [extension](https://github.com/mcncl/zed-mcp-server-buildkite) available in the [official extension gallery](https://zed.dev/extensions?query=buildkite). During installation it will ask for an API token which will be added to your settings. Or you can manually configure:

```jsonc
// ~/.config/zed/settings.json
{
  "context_servers": {
    "mcp-server-buildkite": {
      "settings": {
        "buildkite_api_token": "your-buildkite-token-here",
      }
    }
  }
}
```

# Security

This container image is built using [cgr.dev/chainguard/static](https://images.chainguard.dev/directory/image/static/versions) base image and is configured to run the MCP server as a non-root user.

# Contributing

Notes on building this project are in the [DEVELOPMENT.md](DEVELOPMENT.md)

## Disclaimer

This project is in the early stages of development and is not yet ready for use.

## License

This project is released under MIT license.


[Buildkite API Access Token with read access to pipelines]: https://buildkite.com/user/api-access-tokens/new?scopes[]=read_pipelines
