import { fireEvent, render, screen } from '@testing-library/react-native'

import ProgressHubScreen from './ProgressHubScreen'
import { resetRouterMock, routerMock } from '@/__tests__/router'

describe('ProgressHubScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    resetRouterMock()
  })

  it('muestra la fila Temas', async () => {
    await render(<ProgressHubScreen />)

    expect(screen.getByText('Temas')).toBeOnTheScreen()
    expect(screen.getByLabelText('Temas')).toBeOnTheScreen()
    expect(screen.getByText('Mi semana')).toBeOnTheScreen()
    expect(screen.getByText('Logros')).toBeOnTheScreen()
  })

  it('al pulsar Temas navega a la pantalla de temas', async () => {
    await render(<ProgressHubScreen />)

    await fireEvent.press(screen.getByLabelText('Temas'))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/temas')
  })

  it('al pulsar Mi semana navega al informe', async () => {
    await render(<ProgressHubScreen />)

    await fireEvent.press(screen.getByLabelText('Mi semana'))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/semana')
  })

  it('al pulsar Logros navega a la pantalla de logros', async () => {
    await render(<ProgressHubScreen />)

    await fireEvent.press(screen.getByLabelText('Logros'))

    expect(routerMock.push).toHaveBeenCalledWith('/(tabs)/progreso/logros')
  })
})
