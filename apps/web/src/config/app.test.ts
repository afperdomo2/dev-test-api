import {
  APP_COLOR_PRESETS,
  DEFAULT_APP_COLOR,
  DEFAULT_APP_NAME,
  resolveAppColor,
  resolveAppName,
} from './app'

describe('config/app', () => {
  describe('resolveAppName', () => {
    it('usa DevTest cuando el valor es ausente o vacío', () => {
      expect(resolveAppName(undefined)).toBe(DEFAULT_APP_NAME)
      expect(resolveAppName('')).toBe(DEFAULT_APP_NAME)
      expect(resolveAppName('   ')).toBe(DEFAULT_APP_NAME)
    })

    it('respeta un nombre personalizado y recorta espacios', () => {
      expect(resolveAppName('  DevTest Pruebas  ')).toBe('DevTest Pruebas')
    })
  })

  describe('resolveAppColor', () => {
    it('usa blue por defecto cuando el valor es ausente o desconocido', () => {
      const blue = APP_COLOR_PRESETS[DEFAULT_APP_COLOR]
      expect(resolveAppColor(undefined)).toBe(blue)
      expect(resolveAppColor('')).toBe(blue)
      expect(resolveAppColor('nope')).toBe(blue)
    })

    it('resuelve presets sin distinguir mayúsculas ni espacios', () => {
      expect(resolveAppColor('red')).toBe(APP_COLOR_PRESETS.red)
      expect(resolveAppColor('RED')).toBe(APP_COLOR_PRESETS.red)
      expect(resolveAppColor('  Amber ')).toBe(APP_COLOR_PRESETS.amber)
    })

    it('devuelve un hex distinto por preset', () => {
      const values = Object.values(APP_COLOR_PRESETS)
      expect(new Set(values).size).toBe(values.length)
    })
  })
})
