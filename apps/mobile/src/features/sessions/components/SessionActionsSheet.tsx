import { forwardRef, useCallback, useMemo } from 'react'
import { Alert, Pressable, StyleSheet, Text, View } from 'react-native'
import BottomSheet, { BottomSheetBackdrop, type BottomSheetBackdropProps } from '@gorhom/bottom-sheet'
import type { Session } from '@devtest/shared'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type SessionActionsSheetProps = {
  session: Session | null
  onContinue: (id: string) => void
  onFinish: (id: string) => void
  onReview: (id: string) => void
  onDelete: (id: string) => void
}

export const SessionActionsSheet = forwardRef<BottomSheet, SessionActionsSheetProps>(
  ({ session, onContinue, onFinish, onReview, onDelete }, ref) => {
    const snapPoints = useMemo(() => ['25%', '40%'], [])

    const renderBackdrop = useCallback(
      (props: BottomSheetBackdropProps) => (
        <BottomSheetBackdrop {...props} disappearsOnIndex={-1} appearsOnIndex={0} />
      ),
      [],
    )

    const handleDelete = useCallback(() => {
      if (!session) return
      if (session.answerCount > 0) {
        Alert.alert(
          'No se puede eliminar',
          'Esta sesion tiene respuestas. Finalizala en su lugar.',
          [{ text: 'OK' }],
        )
        return
      }
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
    }, [session, onDelete])

    if (!session) return null

    const isInProgress = session.status === 'in_progress'
    const isCompleted = session.status === 'completed'

    return (
      <BottomSheet
        ref={ref}
        index={-1}
        snapPoints={snapPoints}
        enablePanDownToClose
        backdropComponent={renderBackdrop}
        backgroundStyle={styles.background}
        handleIndicatorStyle={styles.handle}
        handleStyle={styles.handleContainer}
      >
        <View style={styles.content}>
          <View style={styles.header}>
            <Text style={styles.title} numberOfLines={1}>
              {session.name}
            </Text>
            <Text style={styles.subtitle}>{session.status === 'in_progress' ? 'En progreso' : session.status === 'completed' ? 'Completada' : 'Cancelada'}</Text>
          </View>

          <View style={styles.actions}>
            {isInProgress && (
              <>
                <ActionRow
                  icon="▶"
                  label="Continuar sesion"
                  onPress={() => {
                    ;(ref as React.RefObject<BottomSheet>).current?.close()
                    onContinue(session.id)
                  }}
                />
                <ActionRow
                  icon="⏹"
                  label="Finalizar sesion"
                  onPress={() => {
                    ;(ref as React.RefObject<BottomSheet>).current?.close()
                    onFinish(session.id)
                  }}
                />
              </>
            )}

            {(isCompleted || session.status === 'cancelled') && (
              <ActionRow
                icon="📋"
                label="Ver resumen"
                onPress={() => {
                  ;(ref as React.RefObject<BottomSheet>).current?.close()
                  onReview(session.id)
                }}
              />
            )}

            <View style={styles.separator} />

            <ActionRow
              icon="🗑"
              label="Eliminar"
              destructive
              onPress={() => {
                ;(ref as React.RefObject<BottomSheet>).current?.close()
                handleDelete()
              }}
            />
          </View>
        </View>
      </BottomSheet>
    )
  },
)

SessionActionsSheet.displayName = 'SessionActionsSheet'

function ActionRow({
  icon,
  label,
  destructive,
  onPress,
}: {
  icon: string
  label: string
  destructive?: boolean
  onPress: () => void
}) {
  return (
    <Pressable
      style={({ pressed }) => [styles.actionRow, pressed && styles.actionRowPressed]}
      onPress={onPress}
    >
      <Text style={styles.actionIcon}>{icon}</Text>
      <Text style={[styles.actionLabel, destructive && styles.actionLabelDestructive]}>
        {label}
      </Text>
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
    flex: 1,
    paddingHorizontal: spacing.lg,
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
  },
  actionRowPressed: {
    backgroundColor: `${colors.foreground}08`,
  },
  actionIcon: {
    fontSize: typography.size.lg,
    width: 28,
    textAlign: 'center',
  },
  actionLabel: {
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  actionLabelDestructive: {
    color: colors.destructive,
  },
})
