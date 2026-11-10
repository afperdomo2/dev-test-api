import { IMPORT_PROMPT_KINDS, buildImportPrompt, getImportCsvTemplate } from './questionImport'
import type { Topic } from '@devtest/shared'

const topic: Topic = {
  id: 't1',
  slug: 'go',
  name: 'Go',
  category: 'lenguajes',
  isSystem: true,
  createdAt: '2026-01-01T00:00:00Z',
}

describe('questionImport prompt', () => {
  it('ofrece "Mixta" junto a los 3 tipos concretos', () => {
    expect(IMPORT_PROMPT_KINDS.map((k) => k.value)).toEqual([
      'mixed',
      'single_choice',
      'multiple_choice',
      'true_false',
    ])
  })

  it('genera 20 preguntas por defecto', () => {
    const prompt = buildImportPrompt([topic])
    expect(prompt).toContain('EXACTAMENTE 20 preguntas')
    expect(prompt).toContain('total 21 filas')
  })

  it('el prompt mixto reparte entre los 3 tipos y define sus reglas', () => {
    const prompt = buildImportPrompt([topic], 20, { kind: 'mixed' })
    expect(prompt).toContain('MEZCLANDO los 3 tipos')
    expect(prompt).toContain('type=single_choice: EXACTAMENTE 4 opciones')
    expect(prompt).toContain('type=multiple_choice: EXACTAMENTE 4 opciones')
    expect(prompt).toContain('type=true_false: EXACTAMENTE 2 opciones')
    expect(prompt).toContain('NO uses code_completion')
  })

  it('el prompt de un tipo concreto no menciona la mezcla', () => {
    const prompt = buildImportPrompt([topic], 5, { kind: 'true_false' })
    expect(prompt).toContain('tipo true_false')
    expect(prompt).not.toContain('MEZCLANDO')
  })

  it('la plantilla genérica incluye un ejemplo de cada tipo', () => {
    const csv = getImportCsvTemplate('mixed')
    expect(csv).toContain('single_choice,')
    expect(csv).toContain('multiple_choice,')
    expect(csv).toContain('true_false,')
  })
})
