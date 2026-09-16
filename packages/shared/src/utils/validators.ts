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
