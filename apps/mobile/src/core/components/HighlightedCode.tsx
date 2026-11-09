import { memo, useMemo } from 'react'
import { StyleSheet, Text } from 'react-native'
import type { StyleProp, TextStyle } from 'react-native'
import { tokenizeCode } from '@devtest/shared'

import { resolveCodeColor } from '@/core/markdown/code'
import { colors, spacing, typography } from '@/core/theme/tokens'

type HighlightedCodeProps = {
  code: string
  /** Lenguaje del fence (`ts`, `python`, …). Sin lenguaje se pinta plano. */
  language?: string
  style?: StyleProp<TextStyle>
}

/**
 * Bloque de código con resaltado de sintaxis. Usa el tokenizador compartido
 * (`@devtest/shared`) para que web y mobile pinten igual; el color sale de
 * `codeTheme` en los design tokens.
 */
export const HighlightedCode = memo(function HighlightedCode({
  code,
  language,
  style,
}: HighlightedCodeProps) {
  const tokens = useMemo(() => tokenizeCode(code, language), [code, language])

  return (
    <Text style={[styles.fence, style]} selectable>
      {tokens.map((token, index) => {
        const color = resolveCodeColor(token.classes)

        return (
          <Text key={index} style={color ? { color } : undefined}>
            {token.text}
          </Text>
        )
      })}
    </Text>
  )
})

const styles = StyleSheet.create({
  fence: {
    backgroundColor: colors.muted,
    color: colors.foreground,
    fontFamily: typography.family.mono,
    fontSize: typography.size.sm,
    // Explícito para no heredar el `lineHeight` del cuerpo del markdown, que
    // hacía que el bloque ocupara mucho más alto.
    lineHeight: 18,
    padding: spacing.md,
    borderRadius: 8,
    marginVertical: spacing.sm,
  },
})
