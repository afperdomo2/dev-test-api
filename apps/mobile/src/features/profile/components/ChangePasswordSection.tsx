import { useState } from 'react'
import { StyleSheet, Text, View } from 'react-native'

import { isMaxLength, isRequired, isValidPassword } from '@devtest/shared'

import { changePassword } from '@/core/api/services/auth.service'
import { colors, radius, spacing, typography } from '@/core/theme/tokens'
import { getErrorDetail } from '@/core/utils/errors'
import { AuthButton } from '@/features/auth/components/AuthButton'
import { PasswordInput } from '@/features/auth/components/PasswordInput'

type FieldErrors = {
  currentPassword?: string
  newPassword?: string
  confirmPassword?: string
}

export function ChangePasswordSection() {
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [errors, setErrors] = useState<FieldErrors>({})
  const [serverError, setServerError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  function validate(): boolean {
    const next: FieldErrors = {}

    if (!isRequired(currentPassword)) {
      next.currentPassword = 'La contraseña actual es obligatoria'
    } else if (!isMaxLength(currentPassword, 72)) {
      next.currentPassword = 'Máximo 72 caracteres'
    }

    if (!isRequired(newPassword)) {
      next.newPassword = 'La nueva contraseña es obligatoria'
    } else if (!isValidPassword(newPassword, 8)) {
      next.newPassword = 'Mínimo 8 caracteres'
    } else if (!isMaxLength(newPassword, 72)) {
      next.newPassword = 'Máximo 72 caracteres'
    } else if (newPassword === currentPassword) {
      next.newPassword = 'La nueva contraseña debe ser diferente a la actual'
    }

    if (!isRequired(confirmPassword)) {
      next.confirmPassword = 'Confirma la nueva contraseña'
    } else if (confirmPassword !== newPassword) {
      next.confirmPassword = 'Las contraseñas no coinciden'
    }

    setErrors(next)
    setServerError(null)
    setSuccess(null)

    return Object.keys(next).length === 0
  }

  function reset() {
    setCurrentPassword('')
    setNewPassword('')
    setConfirmPassword('')
    setErrors({})
  }

  async function submit(): Promise<void> {
    if (!validate()) return

    setLoading(true)
    try {
      await changePassword({ currentPassword, newPassword })
      reset()
      setSuccess('Contraseña actualizada correctamente')
    } catch (err: unknown) {
      setServerError(getErrorDetail(err, 'Error al actualizar la contraseña'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <View style={styles.card}>
      <Text style={styles.title}>Cambiar contraseña</Text>
      <Text style={styles.subtitle}>Mínimo 8 caracteres. Tu sesión se mantiene activa.</Text>

      <View style={styles.form}>
        <PasswordInput
          label="Contraseña actual"
          value={currentPassword}
          onChangeText={setCurrentPassword}
          error={errors.currentPassword}
          textContentType="password"
          autoComplete="current-password"
          returnKeyType="next"
          editable={!loading}
          testID="change-password-current"
        />

        <PasswordInput
          label="Nueva contraseña"
          value={newPassword}
          onChangeText={setNewPassword}
          error={errors.newPassword}
          textContentType="newPassword"
          autoComplete="new-password"
          returnKeyType="next"
          editable={!loading}
          testID="change-password-new"
        />

        <PasswordInput
          label="Confirmar nueva contraseña"
          value={confirmPassword}
          onChangeText={setConfirmPassword}
          error={errors.confirmPassword}
          textContentType="newPassword"
          autoComplete="new-password"
          returnKeyType="done"
          onSubmitEditing={() => void submit()}
          editable={!loading}
          testID="change-password-confirm"
        />

        {serverError ? (
          <Text style={styles.error} accessibilityLiveRegion="polite">
            {serverError}
          </Text>
        ) : null}

        {success ? (
          <Text style={styles.success} accessibilityLiveRegion="polite">
            {success}
          </Text>
        ) : null}

        <AuthButton title="Actualizar contraseña" onPress={() => void submit()} loading={loading} />
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.lg,
    gap: spacing.xs,
  },
  title: {
    fontSize: typography.size.md,
    fontWeight: typography.weight.semibold,
    color: colors.foreground,
  },
  subtitle: {
    fontSize: typography.size.sm,
    color: colors.mutedForeground,
  },
  form: {
    marginTop: spacing.sm,
    gap: spacing.md,
  },
  error: {
    fontSize: typography.size.sm,
    color: colors.destructive,
  },
  success: {
    fontSize: typography.size.sm,
    color: colors.primary,
  },
})
