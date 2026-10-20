/** @type {import('jest').Config} */
module.exports = {
  preset: 'jest-expo',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  testMatch: ['<rootDir>/src/**/*.test.{ts,tsx}'],
  // Jest no conoce los alias de Metro: `@/` (tsconfig paths) y el paquete del
  // workspace, cuyo `main` apunta a TypeScript.
  moduleNameMapper: {
    '^@devtest/shared$': '<rootDir>/../../packages/shared/src/index.ts',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  collectCoverageFrom: ['src/**/*.{ts,tsx}', '!src/**/*.test.{ts,tsx}', '!src/**/__mocks__/**'],
  clearMocks: true,
  // El primer render de cada archivo paga el arranque de React Native + React 19
  // (transform + setup) y supera los 5 s por defecto en máquinas lentas.
  testTimeout: 30000,
}
