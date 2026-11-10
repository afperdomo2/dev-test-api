import type { User } from '@/types/user.types'

export const ROLE_LABELS = {
  admin: 'Administrador',
  user: 'Usuario',
} as const

/** Etiqueta corta para el título de pestaña, donde el espacio es limitado. */
export const ROLE_SHORT_LABELS = {
  admin: 'Admin',
  user: 'Usuario',
} as const

export function roleLabel(isAdmin: boolean): string {
  return isAdmin ? ROLE_LABELS.admin : ROLE_LABELS.user
}

export function shortRoleLabel(isAdmin: boolean): string {
  return isAdmin ? ROLE_SHORT_LABELS.admin : ROLE_SHORT_LABELS.user
}

/**
 * Título de la pestaña con la identidad del usuario logueado.
 *
 * El rol va primero a propósito: los navegadores truncan el título por la
 * derecha, así que un email largo nunca oculta si la pestaña es de un admin o
 * de un usuario normal.
 */
export function buildDocumentTitle(
  user: Pick<User, 'email' | 'isAdmin'> | null,
  appName: string,
): string {
  if (!user?.email) return appName
  return `${shortRoleLabel(user.isAdmin)} · ${user.email} — ${appName}`
}

/** Inicial para el avatar; `?` mientras la sesión no tiene usuario. */
export function userInitial(email?: string): string {
  const first = email?.trim().charAt(0)
  return first ? first.toUpperCase() : '?'
}
