import * as expoRouter from 'expo-router'

type RouterMock = {
  push: jest.Mock
  replace: jest.Mock
  back: jest.Mock
  setParams: jest.Mock
}

/** Router mockeado globalmente en `jest.setup.ts`. */
export const routerMock = (expoRouter as unknown as { __router: RouterMock }).__router

export function resetRouterMock(): void {
  routerMock.push.mockReset()
  routerMock.replace.mockReset()
  routerMock.back.mockReset()
  routerMock.setParams.mockReset()
}
