import { useCallback, useMemo, useState } from 'react'
import {
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native'
import { useRouter } from 'expo-router'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import type { CreateSessionRequest, SessionDifficulty, Topic } from '@devtest/shared'
import { SESSION_DIFFICULTY_LABELS } from '@devtest/shared'

import { createSession } from '@/core/api/services/sessions.service'
import { getErrorDetail } from '@/core/utils/errors'
import { TopicPicker } from '@/features/sessions/components/TopicPicker'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

type DifficultyOption = {
  value: SessionDifficulty
  icon: keyof typeof MaterialCommunityIcons.glyphMap
}

type QuantityOption = {
  value: number | null
  label: string
}

const DIFFICULTIES: DifficultyOption[] = [
  { value: 'random', icon: 'shuffle-variant' },
  { value: 'beginner', icon: 'signal-cellular-1' },
  { value: 'intermediate', icon: 'signal-cellular-2' },
  { value: 'advanced', icon: 'signal-cellular-3' },
]

const QUANTITIES: QuantityOption[] = [
  { value: 10, label: '10' },
  { value: 25, label: '25' },
  { value: 50, label: '50' },
  { value: null, label: 'Sin límite' },
]

const DEFAULT_QUANTITY = 25

function suggestedName(topics: Topic[]): string {
  if (topics.length === 0) return 'Sesión de estudio'
  const shown = topics.slice(0, 2).map((topic) => topic.name)
  const extra = topics.length > 2 ? ` +${topics.length - 2}` : ''
  return `Estudio · ${shown.join(', ')}${extra}`
}

function isLimitError(err: unknown): boolean {
  return (
    !!err &&
    typeof err === 'object' &&
    'status' in err &&
    (err as { status: number }).status === 422
  )
}

type SectionProps = {
  title: string
  hint?: string
  children: React.ReactNode
}

function Section({ title, hint, children }: SectionProps) {
  return (
    <View style={styles.section}>
      <Text style={styles.sectionTitle}>{title}</Text>
      {hint ? <Text style={styles.sectionHint}>{hint}</Text> : null}
      {children}
    </View>
  )
}

export default function StudySessionScreen() {
  const router = useRouter()
  const [topics, setTopics] = useState<Topic[]>([])
  const [difficulty, setDifficulty] = useState<SessionDifficulty>('random')
  const [quantity, setQuantity] = useState<number | null>(DEFAULT_QUANTITY)
  const [customName, setCustomName] = useState('')
  const [nameDirty, setNameDirty] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [limitReached, setLimitReached] = useState(false)

  const name = nameDirty ? customName : suggestedName(topics)
  const canStart = topics.length > 0 && !submitting

  const topicIDs = useMemo(() => topics.map((topic) => topic.id), [topics])

  const start = useCallback(
    async (unlimited = false) => {
      if (topics.length === 0 || submitting) return
      setSubmitting(true)
      setError(null)
      setLimitReached(false)

      const payload: CreateSessionRequest = {
        name: name.trim() || 'Sesión de estudio',
        mode: 'review',
        difficulty,
        topicIds: topicIDs,
        questionLimit: unlimited || quantity === null ? undefined : quantity,
      }

      try {
        const session = await createSession(payload)
        router.replace(`/(tabs)/sesiones/${session.id}`)
      } catch (err) {
        setError(getErrorDetail(err, 'No se pudo crear la sesión'))
        setLimitReached(isLimitError(err))
      } finally {
        setSubmitting(false)
      }
    },
    [difficulty, name, quantity, router, submitting, topicIDs, topics.length],
  )

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <View style={styles.header}>
        <Pressable
          style={({ pressed }) => [styles.iconButton, pressed && styles.pressed]}
          onPress={() => router.back()}
          hitSlop={8}
          accessibilityRole="button"
          accessibilityLabel="Volver"
        >
          <MaterialCommunityIcons name="chevron-left" size={26} color={colors.foreground} />
        </Pressable>
        <Text style={styles.title}>Nueva sesión</Text>
      </View>

      <ScrollView
        contentContainerStyle={styles.content}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        <Section title="Temas" hint="Solo se muestran temas con preguntas disponibles.">
          <TopicPicker selected={topics} onChange={setTopics} />
        </Section>

        <Section title="Dificultad">
          <View style={styles.difficultyGrid}>
            {DIFFICULTIES.map((option) => {
              const selected = difficulty === option.value
              return (
                <Pressable
                  key={option.value}
                  style={({ pressed }) => [
                    styles.option,
                    selected && styles.optionSelected,
                    pressed && styles.pressed,
                  ]}
                  onPress={() => setDifficulty(option.value)}
                  accessibilityRole="button"
                  accessibilityState={{ selected }}
                  accessibilityLabel={SESSION_DIFFICULTY_LABELS[option.value]}
                >
                  <MaterialCommunityIcons
                    name={option.icon}
                    size={16}
                    color={selected ? colors.primary : colors.mutedForeground}
                  />
                  <Text style={[styles.optionLabel, selected && styles.optionLabelSelected]}>
                    {SESSION_DIFFICULTY_LABELS[option.value]}
                  </Text>
                </Pressable>
              )
            })}
          </View>
        </Section>

        <Section title="Preguntas" hint="Sin límite termina cuando tú finalices la sesión.">
          <View style={styles.quantityRow}>
            {QUANTITIES.map((option) => {
              const selected = quantity === option.value
              return (
                <Pressable
                  key={option.label}
                  style={({ pressed }) => [
                    styles.quantity,
                    selected && styles.optionSelected,
                    pressed && styles.pressed,
                  ]}
                  onPress={() => setQuantity(option.value)}
                  accessibilityRole="button"
                  accessibilityState={{ selected }}
                  accessibilityLabel={
                    option.value === null
                      ? 'Sin límite de preguntas'
                      : `${option.label} preguntas`
                  }
                >
                  <Text style={[styles.optionLabel, selected && styles.optionLabelSelected]}>
                    {option.label}
                  </Text>
                </Pressable>
              )
            })}
          </View>
        </Section>

        <Section title="Nombre">
          <TextInput
            value={name}
            onChangeText={(text) => {
              setNameDirty(true)
              setCustomName(text)
            }}
            placeholder="Nombre de la sesión"
            placeholderTextColor={colors.mutedForeground}
            style={styles.input}
            maxLength={200}
            accessibilityLabel="Nombre de la sesión"
          />
        </Section>

        {error ? (
          <View style={styles.errorBox}>
            <MaterialCommunityIcons
              name="alert-circle-outline"
              size={18}
              color={colors.destructive}
            />
            <Text style={styles.errorText}>{error}</Text>
          </View>
        ) : null}

        <Pressable
          style={({ pressed }) => [
            styles.startButton,
            !canStart && styles.startButtonDisabled,
            pressed && canStart && styles.startButtonPressed,
          ]}
          onPress={() => void start()}
          disabled={!canStart}
          accessibilityRole="button"
          accessibilityState={{ disabled: !canStart, busy: submitting }}
          accessibilityLabel="Empezar sesión"
        >
          {submitting ? (
            <ActivityIndicator color={colors.onPrimary} />
          ) : (
            <Text style={styles.startLabel}>Empezar</Text>
          )}
        </Pressable>

        {limitReached ? (
          <Pressable
            style={({ pressed }) => [styles.unlimitedButton, pressed && styles.pressed]}
            onPress={() => void start(true)}
            disabled={submitting}
            accessibilityRole="button"
            accessibilityLabel="Empezar sin límite de preguntas"
          >
            <Text style={styles.unlimitedLabel}>Empezar sin límite</Text>
          </Pressable>
        ) : null}
      </ScrollView>
    </KeyboardAvoidingView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    paddingHorizontal: spacing.md,
    paddingTop: spacing.lg,
    paddingBottom: spacing.sm,
  },
  iconButton: {
    width: 44,
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    letterSpacing: -0.3,
    color: colors.foreground,
  },
  content: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing['2xl'],
    gap: spacing.lg,
  },
  section: {
    gap: spacing.sm,
  },
  sectionTitle: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  sectionHint: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  difficultyGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.xs,
    flexGrow: 1,
    minWidth: 140,
    minHeight: 44,
    paddingHorizontal: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  optionSelected: {
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}22`,
  },
  optionLabel: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.mutedForeground,
  },
  optionLabelSelected: {
    color: colors.primary,
  },
  quantityRow: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  quantity: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 44,
    paddingHorizontal: spacing.xs,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  input: {
    minHeight: 44,
    paddingHorizontal: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  errorBox: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    padding: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: `${colors.destructive}55`,
    backgroundColor: `${colors.destructive}14`,
  },
  errorText: {
    flex: 1,
    fontSize: typography.size.sm,
    color: colors.destructive,
  },
  startButton: {
    minHeight: 52,
    borderRadius: radius.md,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  startButtonDisabled: {
    opacity: 0.5,
  },
  startButtonPressed: {
    opacity: 0.85,
  },
  startLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.bold,
    color: colors.onPrimary,
  },
  unlimitedButton: {
    minHeight: 44,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.primary,
  },
  unlimitedLabel: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.semibold,
    color: colors.primary,
  },
  pressed: {
    opacity: 0.7,
  },
})
