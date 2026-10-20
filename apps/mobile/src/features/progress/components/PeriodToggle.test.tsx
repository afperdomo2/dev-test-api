import { fireEvent, render, screen } from '@testing-library/react-native'

import { PeriodToggle } from './PeriodToggle'

describe('PeriodToggle', () => {
  it('marca el periodo activo y emite el elegido', async () => {
    const onChange = jest.fn()

    await render(<PeriodToggle value="week" onChange={onChange} />)

    expect(screen.getByLabelText('Semana').props.accessibilityState).toEqual({ selected: true })
    expect(screen.getByLabelText('Mes').props.accessibilityState).toEqual({ selected: false })

    await fireEvent.press(screen.getByLabelText('Mes'))

    expect(onChange).toHaveBeenCalledWith('month')
  })
})
