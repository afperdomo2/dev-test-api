import type { Tip } from '@devtest/shared'

import { clearPendingTip, getPendingTip, setPendingTip } from './dailyTipStore'

const tip: Tip = {
  id: 'tip1',
  userId: 'u1',
  topicId: 't1',
  topic: 'Go',
  content: 'Usa `defer` para cerrar recursos.',
  source: 'manual',
  createdAt: '2026-10-19T00:00:00Z',
  updatedAt: '2026-10-19T00:00:00Z',
}

describe('dailyTipStore', () => {
  beforeEach(() => {
    clearPendingTip()
  })

  it('guarda y devuelve el tip de su sesión', () => {
    setPendingTip('s1', tip)

    expect(getPendingTip('s1')).toEqual(tip)
  })

  it('no devuelve el tip de otra sesión', () => {
    setPendingTip('s1', tip)

    expect(getPendingTip('s2')).toBeNull()
  })

  it('get no consume: se puede leer dos veces', () => {
    setPendingTip('s1', tip)

    expect(getPendingTip('s1')).toEqual(tip)
    expect(getPendingTip('s1')).toEqual(tip)
  })

  it('sin tip no guarda nada', () => {
    setPendingTip('s1', undefined)

    expect(getPendingTip('s1')).toBeNull()
  })

  it('clear lo borra', () => {
    setPendingTip('s1', tip)
    clearPendingTip()

    expect(getPendingTip('s1')).toBeNull()
  })
})
