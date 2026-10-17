import Markdown from 'react-native-markdown-display'
import type { TextStyle } from 'react-native'

import { colors, spacing, typography } from '@/core/theme/tokens'

type MarkdownTextProps = {
  children: string
  style?: TextStyle
}

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
  code_block: {
    backgroundColor: colors.muted,
    color: colors.foreground,
    fontFamily: typography.family.mono,
    fontSize: typography.size.sm,
    padding: spacing.md,
    borderRadius: 8,
    marginVertical: spacing.sm,
  },
  fence: {
    backgroundColor: colors.muted,
    color: colors.foreground,
    fontFamily: typography.family.mono,
    fontSize: typography.size.sm,
    padding: spacing.md,
    borderRadius: 8,
    marginVertical: spacing.sm,
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

export function MarkdownText({ children, style }: MarkdownTextProps) {
  return (
    <Markdown
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
