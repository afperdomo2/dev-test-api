import { isValidEmail, isValidPassword, isRequired, isMaxLength } from '@devtest/shared'

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

export function confirmPasswordRule(
  password: string,
  msg = 'Las contraseñas no coinciden',
): ValidationRule {
  return { validate: (v) => v === password, message: msg }
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
