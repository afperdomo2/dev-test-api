/**
 * Mocks centralizados de módulos nativos.
 * Cada test individual solo necesita sobreescribir lo suyo; aquí queda lo común.
 *
 * Nota: `@testing-library/react-native` v14 registra sus matchers
 * automáticamente al importar la librería, así que no hace falta un
 * `extend-expect` aquí.
 */

jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn(async () => null),
  setItemAsync: jest.fn(async () => undefined),
  deleteItemAsync: jest.fn(async () => undefined),
}))

jest.mock('expo-router', () => {
  const router = {
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
    setParams: jest.fn(),
  }
  return {
    useRouter: () => router,
    useLocalSearchParams: jest.fn(() => ({})),
    Link: ({ children }: { children?: unknown }) => children ?? null,
    Stack: { Screen: () => null },
    __router: router,
  }
})

jest.mock('@expo/vector-icons', () => {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const React = require('react')
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { Text } = require('react-native')
  const Icon = ({ name }: { name?: string }) => React.createElement(Text, null, name ?? '')
  return { MaterialCommunityIcons: Icon, Ionicons: Icon, FontAwesome: Icon }
})
