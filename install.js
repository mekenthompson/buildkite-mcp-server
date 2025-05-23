#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");
const os = require("os");
const { createGunzip } = require("zlib");
const tar = require("tar");

// Configuration
const REPO_OWNER = "buildkite";
const REPO_NAME = "buildkite-mcp-server";
const BINARY_NAME = "buildkite-mcp-server";

// Determine platform and architecture
const platform = getPlatform();
const arch = getArch();
const extension = platform === "windows" ? ".exe" : "";

const binaryPath = path.join(__dirname, "bin", BINARY_NAME + extension);
const binaryDir = path.dirname(binaryPath);

// Ensure bin directory exists
if (!fs.existsSync(binaryDir)) {
  fs.mkdirSync(binaryDir, { recursive: true });
}

console.log(`Downloading ${BINARY_NAME} for ${platform}-${arch}...`);

// Get the latest release info
const releaseUrl = `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest`;
const options = {
  headers: {
    "User-Agent": "nodejs-installer",
  },
};

https
  .get(releaseUrl, options, (res) => {
    let data = "";

    res.on("data", (chunk) => {
      data += chunk;
    });

    res.on("end", () => {
      try {
        const releaseInfo = JSON.parse(data);
        const assets = releaseInfo.assets;

        if (!assets || assets.length === 0) {
          console.error(
            "No release assets found. The project may not have published binaries yet.",
          );
          console.log(
            "You can build from source using: goreleaser build --snapshot --clean",
          );
          process.exit(1);
        }

        // Look for the right asset based on platform and architecture
        // - buildkite-mcp-server_darwin_amd64.tar.gz
        // - buildkite-mcp-server_linux_amd64.tar.gz
        // - buildkite-mcp-server_windows_amd64.zip

        const possibleNames = [
          `${BINARY_NAME}_${platform}_${arch}.tar.gz`,
          `${BINARY_NAME}_${platform}_${arch}.zip`,
          `${BINARY_NAME}-${platform}-${arch}.tar.gz`,
          `${BINARY_NAME}-${platform}-${arch}.zip`,
        ];

        let asset = null;
        for (const name of possibleNames) {
          asset = assets.find((a) => a.name === name);
          if (asset) break;
        }

        if (!asset) {
          console.error(`Could not find binary for ${platform}-${arch}`);
          console.log("Available assets:");
          assets.forEach((a) => console.log(`  - ${a.name}`));
          console.log(
            "\nYou can build from source using: goreleaser build --snapshot --clean",
          );
          process.exit(1);
        }

        // Download and extract the binary
        downloadAndExtract(asset.browser_download_url, asset.name, binaryPath);
      } catch (error) {
        console.error("Error parsing release information:", error);
        console.log(
          "The repository may not have releases yet. You can build from source using:",
        );
        console.log("goreleaser build --snapshot --clean");
        process.exit(1);
      }
    });
  })
  .on("error", (err) => {
    console.error("Error fetching release information:", err);
    process.exit(1);
  });

function downloadAndExtract(url, fileName, destPath) {
  const tempFile = path.join(os.tmpdir(), fileName);
  const file = fs.createWriteStream(tempFile);

  console.log(`Downloading ${fileName}...`);

  https
    .get(url, (res) => {
      res.pipe(file);

      file.on("finish", () => {
        file.close(() => {
          try {
            if (fileName.endsWith(".tar.gz")) {
              // Extract tar.gz
              tar.extract({
                file: tempFile,
                cwd: binaryDir,
                sync: true,
              });

              // Find the extracted binary (might be in a subdirectory)
              const extractedBinary = findBinary(
                binaryDir,
                BINARY_NAME + extension,
              );
              if (extractedBinary && extractedBinary !== destPath) {
                fs.renameSync(extractedBinary, destPath);
              }
            } else if (fileName.endsWith(".zip")) {
              // For ZIP files, you'd need a zip extraction library
              // For now, let's assume tar.gz is the primary format
              console.error("ZIP extraction not implemented yet");
              process.exit(1);
            }

            // Make the binary executable
            fs.chmodSync(destPath, "755");

            // Clean up
            fs.unlinkSync(tempFile);

            console.log(`Successfully installed ${BINARY_NAME}`);
            console.log(`Binary location: ${destPath}`);
          } catch (error) {
            console.error("Error extracting binary:", error);
            process.exit(1);
          }
        });
      });
    })
    .on("error", (err) => {
      fs.unlink(tempFile, () => {});
      console.error("Error downloading binary:", err);
      process.exit(1);
    });
}

function findBinary(dir, binaryName) {
  const files = fs.readdirSync(dir, { recursive: true });
  for (const file of files) {
    const fullPath = path.join(dir, file);
    if (
      path.basename(fullPath) === binaryName &&
      fs.statSync(fullPath).isFile()
    ) {
      return fullPath;
    }
  }
  return null;
}

function getPlatform() {
  const platform = os.platform();

  if (platform === "darwin") return "darwin";
  if (platform === "win32") return "windows";
  if (platform === "linux") return "linux";

  throw new Error(`Unsupported platform: ${platform}`);
}

function getArch() {
  const arch = os.arch();

  if (arch === "x64") return "amd64";
  if (arch === "arm64") return "arm64";
  if (arch === "ia32") return "386";

  throw new Error(`Unsupported architecture: ${arch}`);
}
