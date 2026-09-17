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

export const shadows = {
  sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
  md: '0 4px 6px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px rgba(0, 0, 0, 0.15)',
} as const