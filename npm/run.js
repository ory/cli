#!/usr/bin/env node
// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Launcher for the Ory CLI. The actual binary is shipped in per-platform
// packages (see npm/publish.js) that @ory/cli declares as optionalDependencies,
// so that npm only downloads the one matching the current platform.

"use strict"

var spawn = require("child_process").spawn

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

// Termination signals sent to the launcher are forwarded to the binary rather
// than ending the launcher alone: `ory tunnel` and `ory proxy` delete their
// temporary API key on shutdown, which an orphaned binary never gets to do.
var forwardedSignals = ["SIGINT", "SIGTERM", "SIGHUP"]

var child = spawn(binaryPath(), process.argv.slice(2), { stdio: "inherit" })

forwardedSignals.forEach(function (signal) {
  process.on(signal, function () {
    child.kill(signal)
  })
})

child.on("error", function (err) {
  console.error("Unable to start the Ory CLI binary: " + err.message)
  process.exit(1)
})

child.on("exit", function (code, signal) {
  if (signal) {
    // Re-raise with the default disposition restored so the launcher is
    // terminated by the same signal as the binary.
    forwardedSignals.forEach(function (s) {
      process.removeAllListeners(s)
    })
    process.kill(process.pid, signal)
  }
  process.exit(typeof code === "number" ? code : 1)
})
