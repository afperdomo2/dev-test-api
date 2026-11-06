import { useState } from 'react'
import { StyleSheet, Text } from 'react-native'

import { isRequired, isValidEmail, isValidPassword } from '@devtest/shared'

import { useAuthStore } from '@/core/auth/auth.store'
import { colors, spacing, typography } from '@/core/theme/tokens'
import { getErrorDetail } from '@/core/utils/errors'
import { AuthButton } from '@/features/auth/components/AuthButton'
import { AuthCard } from '@/features/auth/components/AuthCard'
import { AuthInput } from '@/features/auth/components/AuthInput'
import { AuthLayout } from '@/features/auth/components/AuthLayout'
import { AuthServerError } from '@/features/auth/components/AuthServerError'

export default function SetupScreen() {
  const setup = useAuthStore((s) => s.setup)
  const loading = useAuthStore((s) => s.actionLoading)
  const setupRequiresToken = useAuthStore((s) => s.setupRequiresToken)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [setupToken, setSetupToken] = useState('')
  const [errors, setErrors] = useState<{
    email?: string
    password?: string
    confirm?: string
    token?: string
  }>({})
  const [serverError, setServerError] = useState<string | null>(null)

  async function submit(): Promise<void> {
    if (!setupRequiresToken) return

    const nextErrors: typeof errors = {}
    if (!isRequired(email)) {
      nextErrors.email = 'El email es obligatorio'
    } else if (!isValidEmail(email.trim())) {
      nextErrors.email = 'Email no válido'
    }
    if (!isRequired(password)) {
      nextErrors.password = 'La contraseña es obligatoria'
    } else if (!isValidPassword(password)) {
      nextErrors.password = 'Mínimo 8 caracteres'
    }
    if (!isRequired(confirmPassword)) {
      nextErrors.confirm = 'Confirma la contraseña'
    } else if (confirmPassword !== password) {
      nextErrors.confirm = 'Las contraseñas no coinciden'
    }
    if (!isRequired(setupToken)) {
      nextErrors.token = 'El token de instalación es obligatorio'
    }
    setErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) return

    setServerError(null)
    try {
      await setup(email.trim(), password, setupToken)
    } catch (err: unknown) {
      setServerError(getErrorDetail(err, 'Error en la configuración inicial'))
    }
  }

  return (
    <AuthLayout>
      <AuthCard title="Configuración inicial" subtitle="Crea la cuenta de administrador">
        {!setupRequiresToken ? (
          <Text style={styles.disabled} testID="setup-disabled">
            La configuración inicial está deshabilitada en este servidor. Contacta al
            administrador.
          </Text>
        ) : null}

        <AuthInput
          label="Email"
          value={email}
          onChangeText={setEmail}
          error={errors.email}
          keyboardType="email-address"
          editable={!loading}
          testID="setup-email"
        />
        <AuthInput
          label="Contraseña"
          value={password}
          onChangeText={setPassword}
          error={errors.password}
          secureTextEntry
          editable={!loading}
          testID="setup-password"
        />
        <AuthInput
          label="Confirmar contraseña"
          value={confirmPassword}
          onChangeText={setConfirmPassword}
          error={errors.confirm}
          secureTextEntry
          editable={!loading}
          testID="setup-confirm"
        />
        {setupRequiresToken ? (
          <AuthInput
            label="Token de instalación"
            value={setupToken}
            onChangeText={setSetupToken}
            error={errors.token}
            secureTextEntry
            editable={!loading}
            testID="setup-token"
          />
        ) : null}
        <AuthServerError message={serverError} />
        <AuthButton
          title="Crear cuenta"
          onPress={() => void submit()}
          loading={loading}
          disabled={!setupRequiresToken}
        />
      </AuthCard>
    </AuthLayout>
  )
}

const styles = StyleSheet.create({
  disabled: {
    fontSize: typography.size.sm,
    color: colors.foreground,
    textAlign: 'center',
    marginBottom: spacing.md,
    opacity: 0.8,
  },
})
