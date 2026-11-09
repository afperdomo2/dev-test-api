import { render, screen } from '@testing-library/react-native'

import { typography } from '@/core/theme/tokens'

import { MarkdownText } from './MarkdownText'

describe('MarkdownText', () => {
  it('resalta los bloques de código según el lenguaje del fence', async () => {
    const markdown = ['Usa **defer** para cerrar recursos:', '', '```go', 'func main() {}', '```'].join(
      '\n',
    )

    await render(<MarkdownText>{markdown}</MarkdownText>)

    // `func` es keyword en Go: violeta del tema (#C792EA).
    expect(screen.getByText('func')).toHaveStyle({ color: '#C792EA' })
    expect(screen.getByText('defer')).toBeOnTheScreen()
  })

  it('mantiene el tamaño del bloque de código (no hereda el del cuerpo)', async () => {
    await render(<MarkdownText>{['```go', 'func main() {}', '```'].join('\n')}</MarkdownText>)

    expect(screen.getByText(/func main/)).toHaveStyle({
      fontSize: typography.size.sm,
      lineHeight: 18,
    })
  })

  it('pinta plano un fence sin lenguaje', async () => {
    await render(<MarkdownText>{['```', 'sin lenguaje', '```'].join('\n')}</MarkdownText>)

    expect(screen.getByText('sin lenguaje')).toBeOnTheScreen()
  })
})
