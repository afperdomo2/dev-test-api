import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'

import {
  admin,
  DB_HINT,
  ensureBackendIsUp,
  saveAuthState,
  setupToken,
  trackFailures,
  user,
  warmupVite,
} from './fixtures'

async function login(page: Page, email: string, password: string): Promise<void> {
  await page.getByLabel('Email', { exact: true }).fill(email)
  await page.getByLabel('Contraseña', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Iniciar sesión', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Inicio', exact: true })).toBeVisible()
}

test.describe('camino crítico', () => {
  test.beforeAll(async ({ playwright }) => {
    await warmupVite(playwright, 'http://localhost:3000')
  })

  test('status → setup → login → dashboard → sessions', async ({ page, request }) => {
    await ensureBackendIsUp(request)
    const failures = trackFailures(page)

    // 1) status: sin admin, el guard manda a /setup.
    await page.goto('/')
    await page.waitForURL((url) => url.pathname === '/setup' || url.pathname === '/login')

    if (new URL(page.url()).pathname === '/login') {
      throw new Error(DB_HINT)
    }

    await expect(page).toHaveURL(/\/setup$/)
    await expect(page.getByText('Configuración inicial', { exact: true })).toBeVisible()

    // 2) setup: crear el primer admin.
    await page.getByLabel('Email', { exact: true }).fill(admin.email)
    await page.getByLabel('Contraseña', { exact: true }).fill(admin.password)
    await page.getByLabel('Confirmar contraseña', { exact: true }).fill(admin.password)
    // El campo solo aparece cuando el backend exige SETUP_TOKEN.
    const tokenField = page.getByLabel('Token de instalación', { exact: true })
    if (await tokenField.isVisible().catch(() => false)) {
      await tokenField.fill(setupToken)
    }

    await page.getByRole('button', { name: 'Crear cuenta', exact: true }).click()

    await expect(page).toHaveURL(/\/login$/)
    await expect(page.getByRole('button', { name: 'Iniciar sesión', exact: true })).toBeVisible()

    // 3) login: entra al dashboard y el JWT sobrevive a un reload.
    await login(page, admin.email, admin.password)
    await expect(page).toHaveURL(/\/$/)
    await saveAuthState(page, 'admin')

    await page.reload()
    await expect(page).toHaveURL(/\/$/)
    await expect(page.getByRole('heading', { name: 'Inicio', exact: true })).toBeVisible()

    // 4) dashboard: las cards del admin renderizan con datos reales.
    await expect(page.getByText('Cupo importación hoy', { exact: true })).toBeVisible()
    await expect(page.getByText('Generación por lotes', { exact: true })).toBeVisible()

    // 5) un usuario no-admin: las rutas de estudio excluyen al admin.
    await page.getByRole('link', { name: 'Usuarios', exact: true }).click()
    await expect(page.getByRole('heading', { name: 'Usuarios', exact: true })).toBeVisible()

    await page.getByRole('button', { name: 'Nuevo usuario', exact: true }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
    await dialog.getByLabel('Email', { exact: true }).fill(user.email)
    await dialog.getByLabel('Contraseña', { exact: true }).fill(user.password)
    await dialog.getByLabel('Confirmar contraseña', { exact: true }).fill(user.password)
    await dialog.getByRole('button', { name: 'Crear', exact: true }).click()
    await expect(dialog).toBeHidden()

    await page.getByText('Cerrar sesión', { exact: true }).first().click()
    await expect(page).toHaveURL(/\/login$/)

    // 6) sessions: con un usuario normal, la lista carga contra el backend.
    await login(page, user.email, user.password)
    await page.getByRole('link', { name: 'Sesiones', exact: true }).click()

    await expect(page).toHaveURL(/\/sessions$/)
    await expect(page.getByRole('heading', { name: 'Sesiones', exact: true })).toBeVisible()
    await expect(page.getByText('No hay sesiones aún', { exact: true })).toBeVisible()

    // Aserciones transversales: integración limpia, sin errores silenciosos.
    expect(failures.pageErrors, 'pageerror en el navegador').toEqual([])
    expect(failures.consoleErrors, 'errores de consola').toEqual([])
    expect(failures.apiFailures, 'respuestas /api/** con status >= 400').toEqual([])
  })
})
