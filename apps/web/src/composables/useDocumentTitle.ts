import { watchEffect } from 'vue'

/**
 * Mantiene `document.title` sincronizado con el getter.
 *
 * Cada layout declara su propio título (identidad del usuario en el layout
 * autenticado, nombre de la app en el de auth) para que al cambiar de layout o
 * cerrar sesión la pestaña no siga mostrando un usuario anterior.
 */
export function useDocumentTitle(title: () => string): void {
  watchEffect(() => {
    document.title = title()
  })
}
