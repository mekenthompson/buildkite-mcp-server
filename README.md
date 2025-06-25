# buildkite-mcp-server

[![Build status](https://badge.buildkite.com/79fefd75bc7f1898fb35249f7ebd8541a99beef6776e7da1b4.svg?branch=main)](https://buildkite.com/buildkite/buildkite-mcp-server)

> **[Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server exposing Buildkite data (pipelines, builds, jobs, tests) to AI tooling and editors.**

---

## TL;DR Quick-start

```bash
# Requires Docker and a Buildkite API token (see scopes below)
docker run -it --rm -e BUILDKITE_API_TOKEN=BKUA_xxxxx ghcr.io/buildkite/buildkite-mcp-server stdio
```

[![Add to Cursor](https://cursor.com/deeplink/mcp-install-dark.png)](https://cursor.com/install-mcp?name=buildkite&config=eyJjb21tYW5kIjoiZG9ja2VyIHJ1biAtaSAtLXJtIC1lIEJVSUxES0lURV9BUElfVE9LRU4gZ2hjci5pby9idWlsZGtpdGUvYnVpbGRraXRlLW1jcC1zZXJ2ZXIgc3RkaW8iLCJlbnYiOnsiQlVJTERLSVRFX0FQSV9UT0tFTiI6ImJrdWFfeHh4eHh4eHgifX0%3D)

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [API Token Scopes](#api-token-scopes)
- [Installation](#installation)
- [Configuration & Usage](#configuration--usage)
  - [Editors & Tools](#editors--tools)
- [Features](#features)
- [Screenshots](#screenshots)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

---

## Prerequisites

| Requirement | Notes |
|-------------|-------|
| Docker ≥ 20.x | Recommended path – run in an isolated container |
| **OR** Go ≥ 1.22 | Needed only for building/running natively |
| Buildkite API token | Create at https://buildkite.com/user/api-access-tokens |
| Internet access to `ghcr.io` | To pull the pre-built image |

---

## API Token Scopes

### Full functionality

👉 **Quick add:** [Create token with Full functionality](https://buildkite.com/user/api-access-tokens/new?scopes[]=read_clusters&scopes[]=read_pipelines&scopes[]=read_builds&scopes[]=read_build_logs&scopes[]=read_user&scopes[]=read_organizations&scopes[]=read_artifacts&scopes[]=read_suites)

| Scope | Purpose |
|-------|---------|
| `read_clusters` | Access cluster & queue information |
| `read_pipelines` | Pipeline configuration |
| `read_builds` | Builds, jobs & annotations |
| `read_build_logs` | Job log output |
| `read_user` | Current user info |
| `read_organizations` | Organization details |
| `read_artifacts` | Build artifacts & metadata |
| `read_suites` | Buildkite Test Engine data |

### Minimum recommended

👉 **Quick add:** [Create token with Basic functionality](https://buildkite.com/user/api-access-tokens/new?scopes[]=read_builds&scopes[]=read_pipelines&scopes[]=read_user)

| Scope | Purpose |
|-------|---------|
| `read_builds` | Builds, jobs & annotations |
| `read_pipelines` | Pipeline information |
| `read_user` | User identification |

---

## Installation

### 1. Docker (recommended)

```bash
docker pull ghcr.io/buildkite/buildkite-mcp-server
```

Run:

```bash
docker run -it --rm -e BUILDKITE_API_TOKEN=BKUA_xxxxx ghcr.io/buildkite/buildkite-mcp-server stdio
```

### 2. Pre-built binary

Download the latest release from [GitHub Releases](https://github.com/buildkite/buildkite-mcp-server/releases). Binaries are fully-static and require no libc.

### 3. Build from source

```bash
go install github.com/buildkite/buildkite-mcp-server@latest
# or
goreleaser build --snapshot --clean
# or
make build    # uses goreleaser (snapshot)
```

---

## Configuration & Usage

### Editors & Tools

<details>
<summary>Claude Desktop</summary>

```jsonc
{
  "mcpServers": {
    "buildkite": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm", "-e", "BUILDKITE_API_TOKEN",
        "ghcr.io/buildkite/buildkite-mcp-server", "stdio"
      ],
      "env": { "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx" }
    }
  }
}
```

Local binary:

```jsonc
{
  "mcpServers": {
    "buildkite": {
      "command": "buildkite-mcp-server",
      "args": ["stdio"],
      "env": { "BUILDKITE_API_TOKEN": "bkua_xxxxxxxx" }
    }
  }
}
```
</details>

<details>
<summary>Goose</summary>

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

Local:

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
</details>

<details>
<summary>VS Code</summary>

```jsonc
{
  "inputs": [
    {
      "id": "BUILDKITE_API_TOKEN",
      "type": "promptString",
      "description": "Enter your Buildkite Access Token",
      "password": true
    }
  ],
  "servers": {
    "buildkite": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm", "-e", "BUILDKITE_API_TOKEN",
        "ghcr.io/buildkite/buildkite-mcp-server", "stdio"
      ],
      "env": { "BUILDKITE_API_TOKEN": "${input:BUILDKITE_API_TOKEN}" }
    }
  }
}
```
</details>

<details>
<summary>Zed</summary>

```jsonc
// ~/.config/zed/settings.json
{
  "context_servers": {
    "mcp-server-buildkite": {
      "settings": {
        "buildkite_api_token": "your-buildkite-token-here"
      }
    }
  }
}
```
</details>

---

<a name="tools"></a>
<a name="features"></a>
## Tools & Features

| Tool | Description |
|------|-------------|
| `get_cluster` | Detailed cluster info (name, default queue, config) |
| `list_clusters` | List clusters in an organisation |
| `get_cluster_queue` | Details about a specific queue |
| `list_cluster_queues` | List queues in a cluster |
| `get_pipeline` | Detailed pipeline config, steps & stats |
| `list_pipelines` | List all pipelines in an organisation |
| `list_builds` | List builds for a pipeline |
| `get_build` | Detailed build info including jobs |
| `get_build_test_engine_runs` | Test Engine runs for a build |
| `current_user` | Authenticated user info |
| `user_token_organization` | Organisation linked to token |
| `get_jobs` | List jobs for a build |
| `get_job_logs` | Raw log output for a job |
| `list_artifacts` | List artifacts across jobs |
| `get_artifact` | Detailed artifact metadata |
| `list_annotations` | List build annotations |
| `list_test_runs` | Test runs for a suite |
| `get_test_run` | Details of a test run |
| `get_failed_executions` | Failed test executions (with stack traces) |
| `get_test` | Test metadata (for failed executions) |
| `access_token` | Information about the current API token |

---

## Screenshots

![Get Pipeline Tool](docs/images/get_pipeline.png)



---

## Security

To ensure the MCP server is run in a secure environment, we recommend running it in a container.

This image is built from [cgr.dev/chainguard/static](https://images.chainguard.dev/directory/image/static/versions) and runs as an unprivileged user.

---

## Contributing

Development guidelines are in [`DEVELOPMENT.md`](DEVELOPMENT.md).

Run the test suite:

```bash
go test ./...
```

---

## License

MIT © Buildkite

SPDX-License-Identifier: MIT
