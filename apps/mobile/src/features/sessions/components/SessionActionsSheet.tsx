import { forwardRef, useCallback } from 'react'
import { Alert, Pressable, StyleSheet, Text, View } from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import {
  BottomSheetBackdrop,
  BottomSheetModal,
  BottomSheetView,
  type BottomSheetBackdropProps,
} from '@gorhom/bottom-sheet'
import type { Session } from '@devtest/shared'
import { SESSION_MODE_LABELS, SESSION_STATUS_LABELS } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const DELETE_WINDOW_MS = 24 * 60 * 60 * 1000

type SessionActionsSheetProps = {
  session: Session | null
  onContinue: (id: string) => void
  onFinish: (id: string) => void
  onReview: (id: string) => void
  onDelete: (id: string) => void
}

export const SessionActionsSheet = forwardRef<BottomSheetModal, SessionActionsSheetProps>(
  ({ session, onContinue, onFinish, onReview, onDelete }, ref) => {
    const dismiss = useCallback(() => {
      ;(ref as React.RefObject<BottomSheetModal>).current?.dismiss()
    }, [ref])

    const renderBackdrop = useCallback(
      (props: BottomSheetBackdropProps) => (
        <BottomSheetBackdrop {...props} disappearsOnIndex={-1} appearsOnIndex={0} />
      ),
      [],
    )

    const handleFinish = useCallback(() => {
      if (!session) return
      dismiss()
      Alert.alert(
        'Finalizar sesion',
        'Al finalizar se cerrara la sesion y podras ver tu resumen de respuestas.',
        [
          { text: 'Cancelar', style: 'cancel' },
          { text: 'Finalizar', onPress: () => onFinish(session.id) },
        ],
      )
    }, [session, dismiss, onFinish])

    const handleDelete = useCallback(() => {
      if (!session) return
      dismiss()
      Alert.alert(
        'Eliminar sesion',
        'Esta accion no se puede deshacer. La sesion se borrara permanentemente.',
        [
          { text: 'Cancelar', style: 'cancel' },
          {
            text: 'Eliminar',
            style: 'destructive',
            onPress: () => onDelete(session.id),
          },
        ],
      )
    }, [session, dismiss, onDelete])

    const isInProgress = session?.status === 'in_progress'
    const hasResult = session?.status === 'completed' || session?.status === 'cancelled'

    const deleteBlockedReason = getDeleteBlockedReason(session)

    return (
      <BottomSheetModal
        ref={ref}
        enablePanDownToClose
        backdropComponent={renderBackdrop}
        backgroundStyle={styles.background}
        handleIndicatorStyle={styles.handle}
        handleStyle={styles.handleContainer}
      >
        {session && (
          <BottomSheetView style={styles.content}>
            <View style={styles.header}>
              <Text style={styles.title} numberOfLines={1}>
                {session.name}
              </Text>
              <Text style={styles.subtitle}>
                {SESSION_STATUS_LABELS[session.status]} · {SESSION_MODE_LABELS[session.mode]}
              </Text>
            </View>

            <View style={styles.actions}>
              {isInProgress && (
                <>
                  <ActionRow
                    icon="play"
                    label="Continuar sesion"
                    onPress={() => {
                      dismiss()
                      onContinue(session.id)
                    }}
                  />
                  <ActionRow
                    icon="stop"
                    label="Finalizar sesion"
                    onPress={handleFinish}
                  />
                </>
              )}

              {hasResult && (
                <ActionRow
                  icon="file-document-outline"
                  label="Ver resumen"
                  onPress={() => {
                    dismiss()
                    onReview(session.id)
                  }}
                />
              )}

              <View style={styles.separator} />

              <ActionRow
                icon="delete-outline"
                label="Eliminar sesion"
                destructive
                disabled={deleteBlockedReason !== null}
                hint={deleteBlockedReason ?? undefined}
                onPress={handleDelete}
              />
            </View>
          </BottomSheetView>
        )}
      </BottomSheetModal>
    )
  },
)

SessionActionsSheet.displayName = 'SessionActionsSheet'

function getDeleteBlockedReason(session: Session | null): string | null {
  if (!session) return null
  if (session.answerCount > 0) return 'No disponible: la sesion tiene respuestas'
  const createdAt = new Date(session.createdAt).getTime()
  if (Number.isNaN(createdAt) || Date.now() - createdAt > DELETE_WINDOW_MS) {
    return 'Solo dentro de las primeras 24 horas'
  }
  return null
}

function ActionRow({
  icon,
  label,
  hint,
  destructive,
  disabled,
  onPress,
}: {
  icon: keyof typeof MaterialCommunityIcons.glyphMap
  label: string
  hint?: string
  destructive?: boolean
  disabled?: boolean
  onPress: () => void
}) {
  const iconColor = disabled
    ? colors.mutedForeground
    : destructive
      ? colors.destructive
      : colors.foreground

  return (
    <Pressable
      style={({ pressed }) => [
        styles.actionRow,
        pressed && !disabled && styles.actionRowPressed,
        disabled && styles.actionRowDisabled,
      ]}
      onPress={onPress}
      disabled={disabled}
      accessibilityRole="button"
      accessibilityLabel={hint ? `${label}. ${hint}` : label}
      accessibilityState={{ disabled: Boolean(disabled) }}
    >
      <MaterialCommunityIcons name={icon} size={20} color={iconColor} />
      <View style={styles.actionTexts}>
        <Text
          style={[
            styles.actionLabel,
            destructive && !disabled && styles.actionLabelDestructive,
            disabled && styles.actionLabelDisabled,
          ]}
        >
          {label}
        </Text>
        {hint && <Text style={styles.actionHint}>{hint}</Text>}
      </View>
    </Pressable>
  )
}

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.card,
  },
  handle: {
    backgroundColor: colors.mutedForeground,
    width: 36,
  },
  handleContainer: {
    top: 0,
  },
  content: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xl,
  },
  header: {
    paddingBottom: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    gap: spacing.xs,
  },
  title: {
    fontSize: typography.size.lg,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  subtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  actions: {
    paddingTop: spacing.sm,
  },
  separator: {
    height: 1,
    backgroundColor: colors.border,
    marginVertical: spacing.sm,
  },
  actionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.sm,
    borderRadius: radius.sm,
    gap: spacing.md,
    minHeight: 48,
  },
  actionRowPressed: {
    backgroundColor: `${colors.foreground}08`,
  },
  actionRowDisabled: {
    opacity: 0.5,
  },
  actionTexts: {
    flex: 1,
    gap: 2,
  },
  actionLabel: {
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  actionLabelDestructive: {
    color: colors.destructive,
  },
  actionLabelDisabled: {
    color: colors.mutedForeground,
  },
  actionHint: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
})
