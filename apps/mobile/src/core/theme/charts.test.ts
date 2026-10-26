import { heatLevel } from './charts'

describe('heatLevel', () => {
  it('sin actividad devuelve el nivel 0', () => {
    expect(heatLevel(0, 10)).toBe(0)
    expect(heatLevel(3, 0)).toBe(0)
  })

  it('reparte la intensidad en tres tramos sobre el máximo', () => {
    expect(heatLevel(1, 9)).toBe(1)
    expect(heatLevel(5, 9)).toBe(2)
    expect(heatLevel(9, 9)).toBe(3)
  })
})
