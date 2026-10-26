import { MaterialCommunityIcons } from '@expo/vector-icons'
import type { ComponentProps } from 'react'

export type IconName = ComponentProps<typeof MaterialCommunityIcons>['name']

const FALLBACK_ICON: IconName = 'trophy-outline'

/**
 * El backend manda nombres `mdi-*`; MaterialCommunityIcons los usa sin prefijo.
 * Si el nombre no existe en la fuente, se cae al trofeo en vez de romper la pantalla.
 */
export function toMcIcon(name: string): IconName {
  const stripped = name.replace(/^mdi-/, '')
  const glyphMap = MaterialCommunityIcons.glyphMap as Record<string, number> | undefined

  if (glyphMap && !(stripped in glyphMap)) return FALLBACK_ICON

  return stripped as IconName
}
