const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function isValidEmail(email: string): boolean {
  return EMAIL_REGEX.test(email)
}

export function isValidPassword(password: string, minLength = 8): boolean {
  return password.length >= minLength
}

export function isRequired(value: string): boolean {
  return value.trim().length > 0
}

export function isMinLength(value: string, min: number): boolean {
  return value.length >= min
}

export function isMaxLength(value: string, max: number): boolean {
  return value.length <= max
}

export interface ValidationRule {
  validate: (value: string) => boolean
  message: string
}

export function requiredRule(msg = 'Este campo es requerido'): ValidationRule {
  return { validate: isRequired, message: msg }
}

export function emailRule(msg = 'Email inválido'): ValidationRule {
  return { validate: isValidEmail, message: msg }
}

export function passwordRule(min = 8, msg?: string): ValidationRule {
  return {
    validate: (v) => isValidPassword(v, min),
    message: msg ?? `Mínimo ${min} caracteres`,
  }
}

export function maxLengthRule(max: number, msg?: string): ValidationRule {
  return {
    validate: (v) => isMaxLength(v, max),
    message: msg ?? `Máximo ${max} caracteres`,
  }
}

export function validateRules(rules: Array<ValidationRule>, value: string): Array<string> {
  return rules.filter((rule) => !rule.validate(value)).map((rule) => rule.message)
}
