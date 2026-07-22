#!/usr/bin/env node
// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Publishes the Ory CLI to npm for a tagged release. Not part of the npm
// package itself — run by the npm-publish job in .github/workflows/ci.yaml.
//
// For every supported platform this downloads the release archive from GitHub,
// extracts the ory binary into a minimal per-platform package
// (e.g. @ory/cli-linux-x64) and publishes it. It then publishes @ory/cli
// itself, which contains only the npm/run.js launcher plus exact-version
// optionalDependencies on the platform packages, so that npm downloads just
// the binary matching the consumer's platform.
//
// Usage: node npm/publish.js <version> [--dry-run]
//   version    the release tag, with or without the leading "v"
//   --dry-run  build and pack everything, but do not upload to npm

"use strict"

const { execFileSync } = require("child_process")
const fs = require("fs")
const path = require("path")

const platforms = [
  { os: "darwin", cpu: "arm64", assetSuffix: "macOS_arm64.tar.gz" },
  { os: "darwin", cpu: "x64", assetSuffix: "macOS_64bit.tar.gz" },
  { os: "linux", cpu: "arm64", assetSuffix: "linux_arm64.tar.gz" },
  { os: "linux", cpu: "x64", assetSuffix: "linux_64bit.tar.gz" },
  { os: "win32", cpu: "arm64", assetSuffix: "windows_arm64.zip" },
  { os: "win32", cpu: "x64", assetSuffix: "windows_64bit.zip" },
]

const rootDir = path.join(__dirname, "..")
const distDir = path.join(rootDir, "dist", "npm")

function run(cmd, args, opts) {
  console.log("+ " + cmd + " " + args.join(" "))
  return execFileSync(cmd, args, Object.assign({ stdio: "inherit" }, opts))
}

function buildPlatformPackage(platform, version) {
  const pkgName = "@ory/cli-" + platform.os + "-" + platform.cpu
  const asset = "ory_" + version + "-" + platform.assetSuffix
  const url =
    "https://github.com/ory/cli/releases/download/v" + version + "/" + asset
  const pkgDir = path.join(distDir, platform.os + "-" + platform.cpu)
  const binDir = path.join(pkgDir, "bin")
  const binName = platform.os === "win32" ? "ory.exe" : "ory"
  const archive = path.join(distDir, asset)

  fs.mkdirSync(binDir, { recursive: true })
  run("curl", ["-fsSL", "--retry", "3", "-o", archive, url])
  if (asset.endsWith(".zip")) {
    run("unzip", ["-oq", archive, binName, "-d", binDir])
  } else {
    run("tar", ["-xzf", archive, "-C", binDir, binName])
  }
  fs.chmodSync(path.join(binDir, binName), 0o755)
  fs.rmSync(archive)
  fs.copyFileSync(path.join(rootDir, "LICENSE"), path.join(pkgDir, "LICENSE"))
  fs.writeFileSync(
    path.join(pkgDir, "package.json"),
    JSON.stringify(
      {
        name: pkgName,
        version: version,
        description:
          "The Ory CLI binary for " + platform.os + " " + platform.cpu + ".",
        repository: { type: "git", url: "git+https://github.com/ory/cli.git" },
        homepage: "https://ory.com/cli",
        license: "Apache-2.0",
        os: [platform.os],
        cpu: [platform.cpu],
        files: ["bin"],
        preferUnplugged: true,
      },
      null,
      2,
    ) + "\n",
  )
  return { name: pkgName, dir: pkgDir, bin: path.join(binDir, binName) }
}

function smokeTest(pkg, platform, version) {
  if (platform.os !== process.platform || platform.cpu !== process.arch) {
    return
  }
  const out = execFileSync(pkg.bin, ["version"], { encoding: "utf8" })
  console.log(out)
  if (!out.includes(version)) {
    throw new Error(
      "smoke test failed: `ory version` does not mention " + version,
    )
  }
}

function publish(dir, distTag, dryRun) {
  const args = ["publish", "--access", "public", "--tag", distTag]
  if (dryRun) {
    args.push("--dry-run")
  }
  run("npm", args, { cwd: dir })
}

function main() {
  const args = process.argv.slice(2)
  const dryRun = args.includes("--dry-run")
  const version = (args.find((a) => !a.startsWith("--")) || "").replace(
    /^v/,
    "",
  )
  if (!/^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$/.test(version)) {
    console.error("Usage: node npm/publish.js <version> [--dry-run]")
    process.exit(1)
  }

  // Prereleases must not become the version a plain `npm install @ory/cli`
  // resolves to, so publish them under the "next" dist-tag instead of "latest".
  const distTag = version.includes("-") ? "next" : "latest"

  fs.rmSync(distDir, { recursive: true, force: true })
  const built = platforms.map((platform) => {
    const pkg = buildPlatformPackage(platform, version)
    smokeTest(pkg, platform, version)
    return pkg
  })

  // Publish the platform packages before @ory/cli itself so a failure cannot
  // leave @ory/cli pointing at binary packages that do not exist.
  for (const pkg of built) {
    publish(pkg.dir, distTag, dryRun)
  }

  const rootPkgPath = path.join(rootDir, "package.json")
  const rootPkg = JSON.parse(fs.readFileSync(rootPkgPath, "utf8"))
  rootPkg.version = version
  rootPkg.optionalDependencies = {}
  for (const pkg of built) {
    rootPkg.optionalDependencies[pkg.name] = version
  }
  fs.writeFileSync(rootPkgPath, JSON.stringify(rootPkg, null, 2) + "\n")
  publish(rootDir, distTag, dryRun)
}

main()
