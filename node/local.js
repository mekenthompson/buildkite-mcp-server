#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

console.log("🧪 Testing buildkite-mcp-server package locally...\n");

// Check if we're in the right directory
const packageJsonPath = path.join(process.cwd(), "package.json");
if (!fs.existsSync(packageJsonPath)) {
  console.error(
    "❌ package.json not found. Run this from the package root directory.",
  );
  process.exit(1);
}

const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, "utf8"));
if (packageJson.name !== "@buildkite/buildkite-mcp-server") {
  console.error(
    "❌ Wrong package. This script is for @buildkite/buildkite-mcp-server",
  );
  process.exit(1);
}

// Test 1: Check all required files exist
console.log("📁 Checking required files...");
const requiredFiles = ["install.js", "index.js", "bin/run.js"];
for (const file of requiredFiles) {
  if (fs.existsSync(file)) {
    console.log(`✅ ${file}`);
  } else {
    console.log(`❌ ${file} - Missing!`);
  }
}

// Test 2: Make bin/run.js executable (Unix-like systems)
if (process.platform !== "win32") {
  try {
    fs.chmodSync("bin/run.js", "755");
    console.log("✅ Made bin/run.js executable");
  } catch (err) {
    console.log("⚠️  Could not make bin/run.js executable:", err.message);
  }
}

// Test 3: Try to run the install script
console.log("\n🔧 Testing install script...");
try {
  require("./install.js");
} catch (err) {
  console.log("❌ Install script failed:", err.message);
  console.log("This might be expected if there are no releases yet.\n");
}

// Test 4: Check if the module can be required
console.log("📦 Testing module import...");
try {
  const buildkite = require("./index.js");
  console.log("✅ Module imported successfully");
  console.log("   Binary path:", buildkite.path);
  console.log("   Required env:", buildkite.requiredEnv);
} catch (err) {
  console.log("❌ Module import failed:", err.message);
}

// Test 5: Create a tarball for testing
console.log("\n📦 Creating test tarball...");
try {
  const { execSync } = require("child_process");
  execSync("npm pack --dry-run", { stdio: "inherit" });
  console.log("✅ Package structure looks good");
} catch (err) {
  console.log("❌ Package creation failed:", err.message);
}

console.log("\n🎯 Local testing complete!");
console.log("\nNext steps:");
console.log(
  "1. If install script failed, make sure the Buildkite repo has releases",
);
console.log(
  "2. Test with: npm pack && npm install -g ./buildkite-buildkite-mcp-server-*.tgz",
);
console.log("3. Then try: buildkite-mcp-server --help");
console.log(
  "4. If everything works, publish with: npm publish --access public",
);

// Explicitly exit the process to prevent hanging
process.exit(0);
