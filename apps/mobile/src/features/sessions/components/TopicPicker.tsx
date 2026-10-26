import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  ActivityIndicator,
  FlatList,
  Modal,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'
import type { Topic } from '@devtest/shared'

import { listTopics } from '@/core/api/services/topics.service'
import { getErrorDetail } from '@/core/utils/errors'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'

const MAX_TOPICS = 8
const SEARCH_DEBOUNCE_MS = 400
const RESULTS_PER_PAGE = 50

type TopicPickerProps = {
  selected: Topic[]
  onChange: (topics: Topic[]) => void
  max?: number
}

export function TopicPicker({ selected, onChange, max = MAX_TOPICS }: TopicPickerProps) {
  const [open, setOpen] = useState(false)
  const [draft, setDraft] = useState<Topic[]>(selected)
  const [search, setSearch] = useState('')
  const [results, setResults] = useState<Topic[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const requestRef = useRef(0)
  const mountedRef = useRef(true)

  useEffect(() => {
    mountedRef.current = true
    return () => {
      mountedRef.current = false
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
  }, [])

  // La carga se dispara desde handlers (abrir / escribir), no desde un effect:
  // evita renders en cascada y peticiones duplicadas al montar.
  const fetchTopics = useCallback(async (term: string) => {
    const requestID = requestRef.current + 1
    requestRef.current = requestID
    setLoading(true)
    setError(null)
    try {
      const res = await listTopics({
        search: term || undefined,
        hasQuestions: true,
        perPage: RESULTS_PER_PAGE,
      })
      if (mountedRef.current && requestRef.current === requestID) {
        setResults(res.data ?? [])
      }
    } catch (err) {
      if (mountedRef.current && requestRef.current === requestID) {
        setError(getErrorDetail(err, 'No se pudieron cargar los temas'))
      }
    } finally {
      if (mountedRef.current && requestRef.current === requestID) {
        setLoading(false)
      }
    }
  }, [])

  const draftIDs = useMemo(() => new Set(draft.map((topic) => topic.id)), [draft])
  const atMax = draft.length >= max

  const openPicker = useCallback(() => {
    setDraft(selected)
    setSearch('')
    setOpen(true)
    void fetchTopics('')
  }, [fetchTopics, selected])

  const handleSearchChange = useCallback(
    (text: string) => {
      setSearch(text)
      if (debounceRef.current) clearTimeout(debounceRef.current)
      debounceRef.current = setTimeout(() => {
        void fetchTopics(text.trim())
      }, SEARCH_DEBOUNCE_MS)
    },
    [fetchTopics],
  )

  const closePicker = useCallback(() => setOpen(false), [])

  const confirm = useCallback(() => {
    onChange(draft)
    setOpen(false)
  }, [draft, onChange])

  const toggle = useCallback(
    (topic: Topic) => {
      setDraft((prev) => {
        if (prev.some((item) => item.id === topic.id)) {
          return prev.filter((item) => item.id !== topic.id)
        }
        if (prev.length >= max) return prev
        return [...prev, topic]
      })
    },
    [max],
  )

  const removeSelected = useCallback(
    (topicID: string) => {
      onChange(selected.filter((topic) => topic.id !== topicID))
    },
    [onChange, selected],
  )

  return (
    <View style={styles.container}>
      {selected.length > 0 ? (
        <View style={styles.chips}>
          {selected.map((topic) => (
            <Pressable
              key={topic.id}
              style={({ pressed }) => [styles.chip, pressed && styles.pressed]}
              onPress={() => removeSelected(topic.id)}
              hitSlop={8}
              accessibilityRole="button"
              accessibilityLabel={`Quitar ${topic.name}`}
            >
              <Text style={styles.chipText} numberOfLines={1}>
                {topic.name}
              </Text>
              <MaterialCommunityIcons name="close" size={14} color={colors.foreground} />
            </Pressable>
          ))}
        </View>
      ) : (
        <Text style={styles.hint}>Elige al menos un tema para empezar.</Text>
      )}

      <Pressable
        style={({ pressed }) => [styles.addButton, pressed && styles.pressed]}
        onPress={openPicker}
        accessibilityRole="button"
        accessibilityLabel="Añadir temas"
      >
        <MaterialCommunityIcons name="plus" size={18} color={colors.primary} />
        <Text style={styles.addLabel}>
          {selected.length > 0 ? 'Cambiar temas' : 'Añadir temas'}
        </Text>
        <Text style={styles.counter}>
          {selected.length}/{max}
        </Text>
      </Pressable>

      <Modal
        visible={open}
        animationType="slide"
        onRequestClose={closePicker}
        statusBarTranslucent
      >
        <View style={styles.modal}>
          <View style={styles.modalHeader}>
            <Text style={styles.modalTitle}>Elegir temas</Text>
            <Pressable
              style={({ pressed }) => [styles.iconButton, pressed && styles.pressed]}
              onPress={closePicker}
              hitSlop={8}
              accessibilityRole="button"
              accessibilityLabel="Cerrar selector de temas"
            >
              <MaterialCommunityIcons name="close" size={22} color={colors.foreground} />
            </Pressable>
          </View>

          <View style={styles.searchRow}>
            <MaterialCommunityIcons name="magnify" size={18} color={colors.mutedForeground} />
            <TextInput
              value={search}
              onChangeText={handleSearchChange}
              placeholder="Buscar tema"
              placeholderTextColor={colors.mutedForeground}
              style={styles.searchInput}
              autoCorrect={false}
              autoCapitalize="none"
              accessibilityLabel="Buscar tema"
            />
            {search.length > 0 && (
              <Pressable
                onPress={() => handleSearchChange('')}
                hitSlop={12}
                accessibilityRole="button"
                accessibilityLabel="Limpiar búsqueda"
              >
                <MaterialCommunityIcons
                  name="close-circle"
                  size={18}
                  color={colors.mutedForeground}
                />
              </Pressable>
            )}
          </View>

          {atMax && (
            <Text style={styles.maxHint}>
              Máximo {max} temas. Quita uno para añadir otro.
            </Text>
          )}

          {loading ? (
            <View style={styles.center}>
              <ActivityIndicator color={colors.primary} />
            </View>
          ) : error ? (
            <View style={styles.center}>
              <Text style={styles.errorText}>{error}</Text>
            </View>
          ) : (
            <FlatList
              data={results}
              keyExtractor={(topic) => topic.id}
              keyboardShouldPersistTaps="handled"
              contentContainerStyle={styles.list}
              ListEmptyComponent={
                <Text style={styles.empty}>No hay temas con preguntas para esa búsqueda.</Text>
              }
              renderItem={({ item }) => {
                const checked = draftIDs.has(item.id)
                const disabled = !checked && atMax
                return (
                  <Pressable
                    style={({ pressed }) => [
                      styles.row,
                      pressed && styles.pressed,
                      disabled && styles.rowDisabled,
                    ]}
                    onPress={() => toggle(item)}
                    disabled={disabled}
                    accessibilityRole="checkbox"
                    accessibilityState={{ checked, disabled }}
                    accessibilityLabel={`${item.name}, ${item.category}`}
                  >
                    <MaterialCommunityIcons
                      name={checked ? 'checkbox-marked' : 'checkbox-blank-outline'}
                      size={22}
                      color={checked ? colors.primary : colors.mutedForeground}
                    />
                    <View style={styles.rowBody}>
                      <Text style={styles.rowTitle} numberOfLines={1}>
                        {item.name}
                      </Text>
                      <Text style={styles.rowSubtitle} numberOfLines={1}>
                        {item.category}
                      </Text>
                    </View>
                  </Pressable>
                )
              }}
            />
          )}

          <View style={styles.modalFooter}>
            <Text style={styles.counter}>
              {draft.length}/{max} temas
            </Text>
            <Pressable
              style={({ pressed }) => [styles.confirmButton, pressed && styles.confirmPressed]}
              onPress={confirm}
              accessibilityRole="button"
              accessibilityLabel="Confirmar temas"
            >
              <Text style={styles.confirmLabel}>Listo</Text>
            </Pressable>
          </View>
        </View>
      </Modal>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    gap: spacing.sm,
  },
  chips: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.xs,
  },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    maxWidth: '100%',
    paddingHorizontal: spacing.sm + 4,
    paddingVertical: spacing.xs + 2,
    borderRadius: radius.sm,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  chipText: {
    flexShrink: 1,
    fontSize: typography.size.sm,
    color: colors.foreground,
  },
  hint: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  addButton: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    minHeight: 44,
    paddingHorizontal: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    borderStyle: 'dashed',
  },
  addLabel: {
    flex: 1,
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.primary,
  },
  counter: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  pressed: {
    opacity: 0.7,
  },
  modal: {
    flex: 1,
    backgroundColor: colors.background,
    paddingTop: spacing['2xl'],
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.lg,
    gap: spacing.md,
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  modalTitle: {
    fontSize: typography.size.xl,
    fontWeight: typography.weight.bold,
    color: colors.foreground,
  },
  iconButton: {
    width: 44,
    height: 44,
    alignItems: 'center',
    justifyContent: 'center',
  },
  searchRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    minHeight: 44,
    paddingHorizontal: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.card,
  },
  searchInput: {
    flex: 1,
    minHeight: 44,
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  maxHint: {
    fontSize: typography.size.xs,
    color: colors.accent,
  },
  center: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
  },
  errorText: {
    fontSize: typography.size.sm,
    color: colors.destructive,
    textAlign: 'center',
  },
  empty: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
    textAlign: 'center',
    paddingVertical: spacing.xl,
  },
  list: {
    paddingBottom: spacing.md,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    minHeight: 56,
    paddingVertical: spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  rowDisabled: {
    opacity: 0.4,
  },
  rowBody: {
    flex: 1,
    gap: 2,
  },
  rowTitle: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.medium,
    color: colors.foreground,
  },
  rowSubtitle: {
    fontSize: typography.size.xs,
    color: colors.mutedForeground,
  },
  modalFooter: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: spacing.md,
  },
  confirmButton: {
    minHeight: 44,
    minWidth: 120,
    paddingHorizontal: spacing.xl,
    borderRadius: radius.md,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  confirmPressed: {
    opacity: 0.85,
  },
  confirmLabel: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.onPrimary,
  },
})
