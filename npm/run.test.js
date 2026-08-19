// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

"use strict"

const { test } = require("node:test")
const assert = require("node:assert/strict")
const { spawn } = require("node:child_process")
const fs = require("node:fs")
const os = require("node:os")
const path = require("node:path")

const platformPackage = "@ory/cli-" + process.platform + "-" + process.arch

// The binary stand-in is a Node script: it exits with a requested status,
// echoes its arguments, waits for a signal and records which one arrived, or
// terminates itself with SIGTERM. In wait mode it gives up after ten seconds so
// that a launcher which orphans it cannot hold the test's pipes open forever.
const fakeBinary = `#!/usr/bin/env node
"use strict"
const fs = require("node:fs")
const [mode, arg] = process.argv.slice(2)
switch (mode) {
  case "echo":
    console.log(process.argv.slice(3).join(" "))
    break
  case "exit":
    process.exit(Number(arg))
    break
  case "wait":
    for (const signal of ["SIGINT", "SIGTERM", "SIGHUP"]) {
      process.on(signal, () => {
        fs.writeFileSync(arg, signal)
        process.exit(0)
      })
    }
    fs.writeFileSync(arg, "ready")
    setTimeout(() => process.exit(2), 10000)
    break
  case "self-kill":
    process.kill(process.pid, "SIGTERM")
    break
}
`

// Lays out @ory/cli and the platform package the way npm installs them, so the
// launcher resolves the stand-in through the same require.resolve call it uses
// for the real binary.
function install(t, withBinary) {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "ory-cli-launcher-"))
  t.after(() => fs.rmSync(root, { recursive: true, force: true }))

  const launcherDir = path.join(root, "node_modules", "@ory", "cli", "npm")
  fs.mkdirSync(launcherDir, { recursive: true })
  fs.copyFileSync(
    path.join(__dirname, "run.js"),
    path.join(launcherDir, "run.js"),
  )

  if (withBinary) {
    const binDir = path.join(root, "node_modules", platformPackage, "bin")
    fs.mkdirSync(binDir, { recursive: true })
    fs.writeFileSync(path.join(binDir, "ory"), fakeBinary, { mode: 0o755 })
  }

  return { root, launcher: path.join(launcherDir, "run.js") }
}

function run(launcher, args) {
  const child = spawn(process.execPath, [launcher, ...args], {
    stdio: ["ignore", "pipe", "pipe"],
    env: {
      ...process.env,
      // The stand-in's shebang needs the node running this test on the PATH.
      PATH: path.dirname(process.execPath) + path.delimiter + process.env.PATH,
    },
  })
  let stdout = ""
  let stderr = ""
  child.stdout.on("data", (chunk) => (stdout += chunk))
  child.stderr.on("data", (chunk) => (stderr += chunk))
  const closed = new Promise((resolve) => {
    child.on("close", (code, signal) => resolve({ code, signal }))
  })
  return {
    child,
    done: () => closed.then((result) => ({ ...result, stdout, stderr })),
  }
}

async function waitForFile(file, content) {
  for (let i = 0; i < 200; i++) {
    if (fs.existsSync(file) && fs.readFileSync(file, "utf8") === content) {
      return
    }
    await new Promise((resolve) => setTimeout(resolve, 50))
  }
  throw new Error("timed out waiting for " + file + " to contain " + content)
}

const opts = { skip: process.platform === "win32" && "no shebang support" }

test(
  "passes arguments through and exits with the binary's status",
  opts,
  async (t) => {
    const { launcher } = install(t, true)

    const echoed = await run(launcher, [
      "echo",
      "tunnel",
      "--port",
      "4000",
    ]).done()
    assert.equal(echoed.stdout, "tunnel --port 4000\n")
    assert.equal(echoed.code, 0)

    const failed = await run(launcher, ["exit", "3"]).done()
    assert.equal(failed.code, 3)
  },
)

test("forwards termination signals to the binary", opts, async (t) => {
  for (const signal of ["SIGTERM", "SIGINT", "SIGHUP"]) {
    const { root, launcher } = install(t, true)
    const marker = path.join(root, "signal")
    const proc = run(launcher, ["wait", marker])
    await waitForFile(marker, "ready")

    proc.child.kill(signal)

    const result = await proc.done()
    assert.equal(fs.readFileSync(marker, "utf8"), signal)
    assert.equal(result.code, 0)
    assert.equal(result.signal, null)
  }
})

test(
  "terminates with the signal that terminated the binary",
  opts,
  async (t) => {
    const { launcher } = install(t, true)

    const result = await run(launcher, ["self-kill"]).done()
    assert.equal(result.signal, "SIGTERM")
    assert.equal(result.code, null)
  },
)

test("fails when the platform package is missing", opts, async (t) => {
  const { launcher } = install(t, false)

  const result = await run(launcher, ["version"]).done()
  assert.equal(result.code, 1)
  assert.match(result.stderr, new RegExp(platformPackage + " is missing"))
})
