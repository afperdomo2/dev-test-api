const TOKEN_KEY = 'auth_token'
const THEME_KEY = 'app_theme'

export type Theme = 'light' | 'dark'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

/**
 * Tema persistido por dispositivo. Devuelve `null` si no hay valor guardado o
 * si el valor es inválido (p. ej. un `localStorage` manipulado a mano).
 */
export function getTheme(): Theme | null {
  const value = localStorage.getItem(THEME_KEY)
  return value === 'light' || value === 'dark' ? value : null
}

export function setTheme(theme: Theme): void {
  localStorage.setItem(THEME_KEY, theme)
}
