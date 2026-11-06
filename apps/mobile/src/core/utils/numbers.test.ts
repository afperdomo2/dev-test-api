import { toPercent, formatDurationMs } from './numbers'

describe('numbers', () => {
  it('convierte 0..1 a porcentaje entero', () => {
    expect(toPercent(0.6667)).toBe(67)
    expect(toPercent(1)).toBe(100)
    expect(toPercent(0)).toBe(0)
  })

  it('formatea duraciones en minutos y horas', () => {
    expect(formatDurationMs(45 * 60 * 1000)).toBe('45 min')
    expect(formatDurationMs(3 * 60 * 60 * 1000 + 20 * 60 * 1000)).toBe('3h 20m')
    expect(formatDurationMs(2 * 60 * 60 * 1000)).toBe('2h')
  })
})
