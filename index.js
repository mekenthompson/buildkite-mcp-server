const path = require("path");
const os = require("os");

// Export the path to the binary for programmatic use
const binName =
  os.platform() === "win32"
    ? "buildkite-mcp-server.exe"
    : "buildkite-mcp-server";
const binaryPath = path.join(__dirname, "bin", binName);

module.exports = {
  // Path to the binary
  path: binaryPath,

  // Binary name
  name: "buildkite-mcp-server",

  // Our default command is stdio
  stdio: [binaryPath, "stdio"],

  // We'll look for the API token in the environment
  requiredEnv: ["BUILDKITE_API_TOKEN"],
};
