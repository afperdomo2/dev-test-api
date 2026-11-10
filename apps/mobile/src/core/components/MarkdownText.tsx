import Markdown from 'react-native-markdown-display'
import type { ASTNode, RenderRules } from 'react-native-markdown-display'
import type { TextStyle } from 'react-native'

import { HighlightedCode } from '@/core/components/HighlightedCode'
import { stripTrailingNewline } from '@/core/markdown/code'
import { colors, spacing, typography } from '@/core/theme/tokens'

type MarkdownTextProps = {
  children: string
  style?: TextStyle
}

/** `sourceInfo` (el lenguaje del fence) existe en runtime pero no en los tipos. */
type CodeNode = ASTNode & { sourceInfo?: string }

const markdownStyles = {
  body: {
    color: colors.foreground,
    fontSize: typography.size.md,
    lineHeight: 22,
  },
  paragraph: {
    marginTop: 0,
    marginBottom: 0,
  },
  strong: {
    fontWeight: '700' as TextStyle['fontWeight'],
  },
  em: {
    fontStyle: 'italic' as TextStyle['fontStyle'],
  },
  code_inline: {
    backgroundColor: colors.muted,
    color: colors.primary,
    fontFamily: typography.family.mono,
    fontSize: typography.size.sm,
    paddingHorizontal: 4,
    paddingVertical: 2,
    borderRadius: 4,
  },
  bullet_list: {
    marginVertical: spacing.xs,
  },
  ordered_list: {
    marginVertical: spacing.xs,
  },
  list_item: {
    marginTop: spacing.xs,
  },
  link: {
    color: colors.primary,
  },
}

/**
 * Bloques de código con resaltado de sintaxis. El lenguaje viene del fence
 * (```` ```ts ````) y lo resuelve el tokenizador de `@devtest/shared`.
 *
 * No recibe `inheritedStyles`: markdown-display cuela el `fontSize`/`lineHeight`
 * del cuerpo y el bloque quedaba más grande que antes. El look del bloque vive
 * entero en `HighlightedCode`.
 */
function renderCodeBlock(node: ASTNode) {
  return (
    <HighlightedCode
      key={node.key}
      code={stripTrailingNewline(node.content)}
      language={(node as CodeNode).sourceInfo}
    />
  )
}

const markdownRules: RenderRules = {
  fence: renderCodeBlock,
  code_block: renderCodeBlock,
}

export function MarkdownText({ children, style }: MarkdownTextProps) {
  return (
    <Markdown
      rules={markdownRules}
      style={
        {
          ...markdownStyles,
          body: { ...markdownStyles.body, ...style },
        } as never
      }
    >
      {children}
    </Markdown>
  )
}
