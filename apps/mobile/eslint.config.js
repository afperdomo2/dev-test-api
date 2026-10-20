// @ts-check
const expoConfig = require('eslint-config-expo/flat')

/**
 * Flat config de ESLint para la app mobile.
 * `eslint-config-expo/flat` ya exporta un array de configuraciones
 * (reglas de Expo + TypeScript + React + React Native).
 *
 * @type {import('eslint').Linter.Config[]}
 */
module.exports = [
  ...expoConfig,
  {
    ignores: ['.expo/*', 'dist/*', 'node_modules/*', 'coverage/*'],
  },
  {
    // Regla nueva del React Compiler (`react-hooks/set-state-in-effect`): marca
    // el patrón estándar de carga inicial dentro de un useEffect, que la app usa
    // en SessionsListPage. Se deja como aviso hasta migrar a un data layer
    // (React Query/SWR); como error bloquearía código que funciona.
    rules: {
      'react-hooks/set-state-in-effect': 'warn',
    },
  },
]
