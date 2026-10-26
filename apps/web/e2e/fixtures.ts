import { mkdir, readdir } from 'node:fs/promises'
import { join, relative } from 'node:path'
import type { APIRequestContext, Page, Playwright } from '@playwright/test'

/** URL directa del backend (Playwright comprueba /health antes de empezar). */
export const API_URL = process.env.E2E_API_URL ?? 'http://localhost:8080'

export type Credentials = {
  email: string
  password: string
}

/** El setup crea el primer admin; la DB debe estar vacía (OQ-1). */
export const admin: Credentials = {
  email: 'e2e-admin@example.com',
  password: 'E2eAdmin123!',
}

/** Usuario no-admin: las rutas de estudio (`/sessions`) excluyen al admin. */
export const user: Credentials = {
  email: 'e2e-user@example.com',
  password: 'E2eUser123!',
}

/**
 * Token de instalación. Debe coincidir con el `SETUP_TOKEN` con el que arranca
 * el backend de E2E (por defecto: `e2e-setup-token`).
 */
export const setupToken = process.env.E2E_SETUP_TOKEN ?? 'e2e-setup-token'

export const DB_HINT =
  'La DB de E2E ya está configurada (existe un admin). Recrea una DB limpia:\n' +
  '  1) DROP DATABASE devtest_e2e; CREATE DATABASE devtest_e2e;\n' +
  '  2) DB_NAME=devtest_e2e make migrate-up\n' +
  '  3) DB_NAME=devtest_e2e make run'

/**
 * Falla con un mensaje accionable si el backend no está escuchando, en vez de
 * dejar un timeout críptico en el primer goto.
 */
export async function ensureBackendIsUp(request: APIRequestContext): Promise<void> {
  try {
    const response = await request.get(`${API_URL}/health`, { timeout: 5_000 })
    if (!response.ok()) throw new Error(`HTTP ${response.status()}`)
  } catch {
    throw new Error(
      `El backend no responde en ${API_URL}/health. Levántalo con "DB_NAME=devtest_e2e make run".`,
    )
  }
}

type FailureTracker = {
  consoleErrors: Array<string>
  pageErrors: Array<string>
  apiFailures: Array<string>
}

/** Registra errores de consola, pageerror y respuestas /api/** >= 400. */
export function trackFailures(page: Page): FailureTracker {
  const tracker: FailureTracker = {
    consoleErrors: [],
    pageErrors: [],
    apiFailures: [],
  }

  page.on('console', (message) => {
    if (message.type() === 'error') tracker.consoleErrors.push(message.text())
  })
  page.on('pageerror', (error) => tracker.pageErrors.push(error.message))
  page.on('response', (response) => {
    if (!response.url().includes('/api/')) return
    if (response.status() >= 400) {
      tracker.apiFailures.push(
        `${response.status()} ${response.request().method()} ${response.url()}`,
      )
    }
  })

  return tracker
}

export async function saveAuthState(page: Page, name: string): Promise<void> {
  const directory = join(process.cwd(), 'e2e', '.auth')
  await mkdir(directory, { recursive: true })
  await page.context().storageState({ path: join(directory, `${name}.json`) })
}

async function sourceFiles(directory: string): Promise<Array<string>> {
  const entries = await readdir(directory, { withFileTypes: true })
  const files: Array<string> = []

  for (const entry of entries) {
    const path = join(directory, entry.name)
    if (entry.isDirectory()) {
      if (entry.name === '__tests__') continue
      files.push(...(await sourceFiles(path)))
      continue
    }
    if (!/\.(vue|ts)$/.test(entry.name)) continue
    if (/\.(test|spec)\.ts$/.test(entry.name)) continue
    files.push(path)
  }

  return files
}

/**
 * Calienta el grafo de módulos de Vite antes de abrir el navegador.
 *
 * Sin esto, la primera visita a una ruta lazy hace que Vite descubra
 * componentes Vuetify nuevos (`vuetify/components/VDataTable`, …),
 * reoptimice y recargue la página en mitad de la navegación: el click se
 * pierde y el test ve una página en blanco.
 */
export async function warmupVite(playwright: Playwright, baseURL: string): Promise<void> {
  const context = await playwright.request.newContext({ baseURL })

  try {
    for (const file of await sourceFiles(join(process.cwd(), 'src'))) {
      const path = relative(process.cwd(), file).split('\\').join('/')
      await context.get(`/${path}`)
    }
  } finally {
    await context.dispose()
  }

  // La reoptimización de Vite es asíncrona (tiene debounce interno).
  await new Promise((resolve) => setTimeout(resolve, 3_000))
}
