import {
  isValidEmail,
  isValidPassword,
  isRequired,
  isMinLength,
  isMaxLength,
  requiredRule,
  emailRule,
  passwordRule,
  confirmPasswordRule,
  maxLengthRule,
  validateRules,
} from './validators'

describe('validators', () => {
  describe('isValidEmail', () => {
    it('returns true for valid emails', () => {
      expect(isValidEmail('user@example.com')).toBe(true)
      expect(isValidEmail('a.b+tag@sub.domain.co')).toBe(true)
    })

    it('returns false for invalid emails', () => {
      expect(isValidEmail('')).toBe(false)
      expect(isValidEmail('plainaddress')).toBe(false)
      expect(isValidEmail('@missing.com')).toBe(false)
      expect(isValidEmail('user@')).toBe(false)
      expect(isValidEmail('user @example.com')).toBe(false)
    })
  })

  describe('isValidPassword', () => {
    it('validates default min 8', () => {
      expect(isValidPassword('12345678')).toBe(true)
      expect(isValidPassword('1234567')).toBe(false)
      expect(isValidPassword('')).toBe(false)
    })

    it('respects custom minLength', () => {
      expect(isValidPassword('abc', 3)).toBe(true)
      expect(isValidPassword('ab', 3)).toBe(false)
    })
  })

  describe('isRequired', () => {
    it('true when non-empty trimmed', () => {
      expect(isRequired('a')).toBe(true)
      expect(isRequired('  x  ')).toBe(true)
    })

    it('false when empty or whitespace', () => {
      expect(isRequired('')).toBe(false)
      expect(isRequired('   ')).toBe(false)
      expect(isRequired('\t\n')).toBe(false)
    })
  })

  describe('isMinLength / isMaxLength', () => {
    it('checks boundaries', () => {
      expect(isMinLength('abc', 3)).toBe(true)
      expect(isMinLength('ab', 3)).toBe(false)
      expect(isMaxLength('abc', 3)).toBe(true)
      expect(isMaxLength('abcd', 3)).toBe(false)
    })
  })

  describe('rule factories', () => {
    it('requiredRule uses default message', () => {
      const rule = requiredRule()
      expect(rule.message).toBe('Este campo es requerido')
      expect(rule.validate('')).toBe(false)
      expect(rule.validate('hi')).toBe(true)
    })

    it('requiredRule custom message', () => {
      const rule = requiredRule('custom')
      expect(rule.message).toBe('custom')
    })

    it('emailRule validates', () => {
      const rule = emailRule()
      expect(rule.validate('a@b.com')).toBe(true)
      expect(rule.validate('bad')).toBe(false)
      expect(rule.message).toBe('Email inválido')
    })

    it('passwordRule default 8', () => {
      const rule = passwordRule()
      expect(rule.validate('12345678')).toBe(true)
      expect(rule.validate('1234567')).toBe(false)
      expect(rule.message).toBe('Mínimo 8 caracteres')
    })

    it('passwordRule custom min and msg', () => {
      const rule = passwordRule(4, 'short')
      expect(rule.validate('abcd')).toBe(true)
      expect(rule.validate('abc')).toBe(false)
      expect(rule.message).toBe('short')
    })

    it('confirmPasswordRule matches', () => {
      const rule = confirmPasswordRule('secret')
      expect(rule.validate('secret')).toBe(true)
      expect(rule.validate('other')).toBe(false)
      expect(rule.message).toBe('Las contraseñas no coinciden')
    })

    it('maxLengthRule', () => {
      const rule = maxLengthRule(5)
      expect(rule.validate('12345')).toBe(true)
      expect(rule.validate('123456')).toBe(false)
      expect(rule.message).toBe('Máximo 5 caracteres')
    })

    it('maxLengthRule custom message', () => {
      const rule = maxLengthRule(3, 'too long')
      expect(rule.message).toBe('too long')
    })
  })

  describe('validateRules', () => {
    it('returns empty when all pass', () => {
      const rules = [requiredRule(), emailRule()]
      expect(validateRules(rules, 'a@b.com')).toEqual([])
    })

    it('returns failing messages', () => {
      const rules = [requiredRule('req'), emailRule('bad email')]
      expect(validateRules(rules, '')).toEqual(['req', 'bad email'])
      expect(validateRules(rules, 'not-an-email')).toEqual(['bad email'])
    })

    it('handles empty rules', () => {
      expect(validateRules([], 'any')).toEqual([])
    })
  })
})
