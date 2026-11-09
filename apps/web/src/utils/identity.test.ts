import { buildDocumentTitle, roleLabel, shortRoleLabel, userInitial } from './identity'

describe('roleLabel', () => {
  it.each([
    [true, 'Administrador'],
    [false, 'Usuario'],
  ])('isAdmin=%s → %s', (isAdmin, expected) => {
    expect(roleLabel(isAdmin)).toBe(expected)
  })
})

describe('shortRoleLabel', () => {
  it('usa Admin para administradores', () => {
    expect(shortRoleLabel(true)).toBe('Admin')
  })

  it('usa Usuario para cuentas normales', () => {
    expect(shortRoleLabel(false)).toBe('Usuario')
  })
})

describe('buildDocumentTitle', () => {
  it('pone el rol primero y el nombre de la app al final', () => {
    expect(buildDocumentTitle({ email: 'admin@empresa.com', isAdmin: true }, 'DevTest')).toBe(
      'Admin · admin@empresa.com — DevTest',
    )
  })

  it('usa Usuario para cuentas normales', () => {
    expect(buildDocumentTitle({ email: 'juan@empresa.com', isAdmin: false }, 'DevTest')).toBe(
      'Usuario · juan@empresa.com — DevTest',
    )
  })

  it('cae al nombre de la app cuando no hay usuario', () => {
    expect(buildDocumentTitle(null, 'DevTest')).toBe('DevTest')
  })

  it('cae al nombre de la app cuando el email está vacío', () => {
    expect(buildDocumentTitle({ email: '', isAdmin: false }, 'DevTest')).toBe('DevTest')
  })
})

describe('userInitial', () => {
  it('devuelve la primera letra en mayúscula', () => {
    expect(userInitial('juan@empresa.com')).toBe('J')
  })

  it('ignora espacios alrededor del email', () => {
    expect(userInitial('  admin@empresa.com')).toBe('A')
  })

  it('devuelve ? sin email', () => {
    expect(userInitial(undefined)).toBe('?')
    expect(userInitial('')).toBe('?')
    expect(userInitial('   ')).toBe('?')
  })
})
