import { render, screen } from '@testing-library/react-native'

import { HighlightedCode } from './HighlightedCode'

describe('HighlightedCode', () => {
  it('pinta el bloque con colores por token', async () => {
    await render(<HighlightedCode code="func main() {}" language="go" />)

    // `func` es keyword en Go: violeta del tema (#C792EA).
    expect(screen.getByText('func')).toHaveStyle({ color: '#C792EA' })
    expect(screen.getByText('func')).toBeOnTheScreen()
  })

  it('mantiene el código completo como texto', async () => {
    await render(<HighlightedCode code="const a = 1" language="typescript" />)

    expect(screen.getByText('const')).toBeOnTheScreen()
    expect(screen.getByText('1')).toBeOnTheScreen()
  })

  it('pinta plano cuando no hay lenguaje', async () => {
    await render(<HighlightedCode code="texto sin lenguaje" />)

    expect(screen.getByText('texto sin lenguaje')).toBeOnTheScreen()
  })
})
