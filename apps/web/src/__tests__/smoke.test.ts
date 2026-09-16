describe('vitest setup', () => {
  it('runs with globals', () => {
    expect(1 + 1).toBe(2)
  })

  it('has jsdom', () => {
    expect(typeof window).toBe('object')
    expect(typeof document).toBe('object')
  })
})
