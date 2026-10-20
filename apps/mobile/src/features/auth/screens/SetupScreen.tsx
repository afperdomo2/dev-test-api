import { useState } from 'react'

import { isRequired, isValidEmail, isValidPassword } from '@devtest/shared'

import { useAuthStore } from '@/core/auth/auth.store'
import { getErrorDetail } from '@/core/utils/errors'
import { AuthButton } from '@/features/auth/components/AuthButton'
import { AuthCard } from '@/features/auth/components/AuthCard'
import { AuthInput } from '@/features/auth/components/AuthInput'
import { AuthLayout } from '@/features/auth/components/AuthLayout'
import { AuthServerError } from '@/features/auth/components/AuthServerError'

export default function SetupScreen() {
  const setup = useAuthStore((s) => s.setup)
  const loading = useAuthStore((s) => s.actionLoading)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [errors, setErrors] = useState<{ email?: string; password?: string; confirm?: string }>({})
  const [serverError, setServerError] = useState<string | null>(null)

  async function submit(): Promise<void> {
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
    setErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) return

    setServerError(null)
    try {
      await setup(email.trim(), password)
    } catch (err: unknown) {
      setServerError(getErrorDetail(err, 'Error en la configuración inicial'))
    }
  }

  return (
    <AuthLayout>
      <AuthCard title="Configuración inicial" subtitle="Crea la cuenta de administrador">
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
        <AuthServerError message={serverError} />
        <AuthButton title="Crear cuenta" onPress={() => void submit()} loading={loading} />
      </AuthCard>
    </AuthLayout>
  )
}
