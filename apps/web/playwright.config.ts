import { defineConfig } from '@playwright/test'

/**
 * E2E del camino crítico (B5.1).
 *
 * Requisitos externos (no los levanta Playwright):
 *   1. DB dedicada y vacía: DB_NAME=devtest_e2e make migrate-up
 *   2. Backend: make run (con DB_NAME=devtest_e2e)
 *
 * El frontend lo levanta `webServer`, salvo que `PLAYWRIGHT_BASE_URL` apunte a
 * un stack ya levantado (por ejemplo `make e2e-docker`).
 */
const baseURL = process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:3000'
const useExternalServer = Boolean(process.env.PLAYWRIGHT_BASE_URL)

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false, // la suite comparte estado en la DB
  workers: 1,
  retries: 0,
  // Vite en frío (primer arranque + optimización de Vuetify) puede tardar;
  // 30 s (el default) no alcanza para cargar el grafo de módulos.
  timeout: 120_000,
  // El primer render de un chunk lazy (dashboard) puede coincidir con la
  // optimización de dependencias de Vite y recargar la página.
  expect: { timeout: 30_000 },
  reporter: [['list'], ['html', { open: 'never' }]],
  outputDir: 'test-results',
  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  // `make e2e-docker` ya levanta el stack: no arrancar Vite en ese caso.
  ...(useExternalServer
    ? {}
    : {
        webServer: {
          command: 'pnpm dev',
          url: 'http://localhost:3000',
          reuseExistingServer: true,
          timeout: 120_000,
        },
      }),
})
