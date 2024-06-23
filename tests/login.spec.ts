import { test, expect } from "@playwright/test"
import { ChildProcessWithoutNullStreams, spawn } from "child_process"
import { randomBytes } from "crypto"
import { unlink } from "fs/promises"
import readline from "node:readline/promises"

function generateRandomFileName(extension: string): string {
  const randomString = randomBytes(16).toString("hex")
  return `${randomString}${extension}`
}

test.describe("should be able to login with the CLI", () => {
  const config = generateRandomFileName(".cli-config.json")
  let url: string = ""
  let child: ChildProcessWithoutNullStreams
  let rl: readline.Interface

  test.beforeAll(async () => {
    child = spawn("./cli", ["auth"], {
      env: {
        HOME: "/dev/null",
        ORY_CLOUD_ORYAPIS_URL: "https://oryapis:8080",
        ORY_CLOUD_CONSOLE_URL: "https://console.ory:8080",
        ORY_CLOUD_CONFIG_PATH: config,
      },
      stdio: "pipe",
      cwd: process.cwd(),
      detached: false,
    })
    child.on("error", (error) => {
      test.fail(true, "Error running the CLI command")
    })

    rl = readline.createInterface(child.stderr)
    await expect(async () => {
      const line = await rl[Symbol.asyncIterator]().next()
      expect(line.done).toBeFalsy()
      const match = line.value.match(
        new RegExp("https://project.console.ory.*"),
      )
      expect(match).toBeTruthy()
      url = match[0]
    }).toPass()
  })

  test.afterAll(async () => {
    child.kill()
    await unlink(config)
  })

  test("with email and password", async ({ page }) => {
    await page.goto(url)
    const emailInput = await page.locator(
      `[data-testid="node/input/identifier"] input`,
    )
    await emailInput.fill("<email>")
    const passwordInput = await page.locator(
      `[data-testid="node/input/password"] input`,
    )
    await passwordInput.fill("<password>")
    const submit = page.locator(
      '[type="submit"][name="method"][value="password"]',
    )
    await submit.click()

    const allow = await page.getByRole("button", { name: "Allow" })
    await allow.click()

    let success = false
    for await (const line of rl) {
      if (line.includes("Successfully logged into Ory Network")) {
        success = true
      }
    }
    expect(success).toBeTruthy()
  })
})
