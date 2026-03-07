import { defineConfig, devices } from '@playwright/test'

// 前端地址：可通过环境变量覆盖，例如 PLAYWRIGHT_BASE_URL=http://localhost:5174
const baseURL = process.env.PLAYWRIGHT_BASE_URL ?? 'http://192.168.31.115:5174'

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  timeout: 60000,
  expect: { timeout: 15000 },
})
