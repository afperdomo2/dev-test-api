import { Platform } from 'react-native'

export const colors = {
  primary: '#22C55E',
  onPrimary: '#0F172A',
  secondary: '#059669',
  onSecondary: '#000000',
  accent: '#D97706',
  onAccent: '#000000',
  background: '#0F172A',
  foreground: '#FFFFFF',
  card: '#192134',
  cardForeground: '#FFFFFF',
  muted: '#10242E',
  mutedForeground: '#94A3B8',
  border: 'rgba(255, 255, 255, 0.08)',
  destructive: '#DC2626',
  onDestructive: '#FFFFFF',
  ring: '#22C55E',
} as const

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  '2xl': 48,
  '3xl': 64,
} as const

export const radius = {
  sm: 8,
  md: 12,
  lg: 16,
} as const

export const typography = {
  family: {
    ui: undefined as string | undefined,
    mono: Platform.select({ ios: 'Menlo', default: 'monospace' }),
  },
  size: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 20,
    xl: 24,
    '2xl': 32,
  } as const,
  weight: {
    regular: '400',
    medium: '500',
    semibold: '600',
    bold: '700',
  } as const,
}

/**
 * Colores de los tokens que devuelve `tokenizeCode` (`@devtest/shared`).
 * La clave es la clase `hljs-*`; si ninguna coincide se usa
 * `colors.foreground`. Contraste verificado sobre `colors.muted` (#10242E).
 */
export const codeTheme: Record<string, string> = {
  'hljs-comment': '#8494A3',
  'hljs-quote': '#8494A3',
  'hljs-doctag': '#8494A3',
  'hljs-meta': '#8494A3',
  'hljs-keyword': '#C792EA',
  'hljs-selector-tag': '#C792EA',
  'hljs-string': '#98C379',
  'hljs-regexp': '#98C379',
  'hljs-addition': '#98C379',
  'hljs-code': '#98C379',
  'hljs-number': '#F78C6C',
  'hljs-title': '#61AFEF',
  'hljs-function': '#61AFEF',
  'hljs-section': '#61AFEF',
  'hljs-link': '#61AFEF',
  'hljs-type': '#E5C07B',
  'hljs-class': '#E5C07B',
  'hljs-attr': '#E5C07B',
  'hljs-attribute': '#E5C07B',
  'hljs-bullet': '#E5C07B',
  'hljs-selector-class': '#E5C07B',
  'hljs-selector-id': '#E5C07B',
  'hljs-built_in': '#56B6C2',
  'hljs-literal': '#56B6C2',
  'hljs-operator': '#56B6C2',
  'hljs-formula': '#56B6C2',
  'hljs-selector-attr': '#56B6C2',
  'hljs-selector-pseudo': '#56B6C2',
  'hljs-variable': '#E06C75',
  'hljs-template-variable': '#E06C75',
  'hljs-property': '#E06C75',
  'hljs-name': '#E06C75',
  'hljs-tag': '#E06C75',
  'hljs-symbol': '#E06C75',
  'hljs-deletion': '#E06C75',
  'hljs-punctuation': colors.mutedForeground,
}

export const shadows = {
  sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
  md: '0 4px 6px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px rgba(0, 0, 0, 0.15)',
} as const