#!/usr/bin/env node
// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Launcher for the Ory CLI. The actual binary is shipped in per-platform
// packages (see npm/publish.js) that @ory/cli declares as optionalDependencies,
// so that npm only downloads the one matching the current platform.

"use strict"

var spawnSync = require("child_process").spawnSync

var packages = {
  "darwin arm64": "@ory/cli-darwin-arm64",
  "darwin x64": "@ory/cli-darwin-x64",
  "linux arm64": "@ory/cli-linux-arm64",
  "linux x64": "@ory/cli-linux-x64",
  "win32 arm64": "@ory/cli-win32-arm64",
  "win32 x64": "@ory/cli-win32-x64",
}

function binaryPath() {
  var platformKey = process.platform + " " + process.arch
  var pkg = packages[platformKey]
  if (!pkg) {
    console.error(
      "@ory/cli does not ship a prebuilt Ory CLI binary for " +
        platformKey +
        ".\nSee https://www.ory.com/docs/guides/cli/installation for other installation options.",
    )
    process.exit(1)
  }
  var bin = process.platform === "win32" ? "bin/ory.exe" : "bin/ory"
  try {
    return require.resolve(pkg + "/" + bin)
  } catch (err) {
    console.error(
      "The Ory CLI binary package " +
        pkg +
        " is missing. It is an optional dependency of @ory/cli, so make sure\n" +
        "optional dependencies are not disabled (e.g. via --omit=optional or\n" +
        "--no-optional) and reinstall.",
    )
    process.exit(1)
  }
}

var result = spawnSync(binaryPath(), process.argv.slice(2), {
  stdio: "inherit",
})
if (result.error) {
  throw result.error
}
if (result.signal) {
  process.kill(process.pid, result.signal)
}
process.exit(typeof result.status === "number" ? result.status : 1)
