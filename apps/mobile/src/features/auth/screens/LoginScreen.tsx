import { useState } from 'react'

import { isRequired, isValidEmail } from '@devtest/shared'

import { useAuthStore } from '@/core/auth/auth.store'
import { getErrorDetail } from '@/core/utils/errors'
import { AuthButton } from '@/features/auth/components/AuthButton'
import { AuthCard } from '@/features/auth/components/AuthCard'
import { AuthInput } from '@/features/auth/components/AuthInput'
import { AuthLayout } from '@/features/auth/components/AuthLayout'
import { AuthServerError } from '@/features/auth/components/AuthServerError'

export default function LoginScreen() {
  const login = useAuthStore((s) => s.login)
  const loading = useAuthStore((s) => s.actionLoading)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({})
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
    }
    setErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) return

    setServerError(null)
    try {
      await login(email.trim(), password)
    } catch (err: unknown) {
      setServerError(getErrorDetail(err, 'Error al iniciar sesión'))
    }
  }

  return (
    <AuthLayout>
      <AuthCard title="Iniciar sesión">
        <AuthInput
          label="Email"
          value={email}
          onChangeText={setEmail}
          error={errors.email}
          keyboardType="email-address"
          editable={!loading}
          testID="login-email"
        />
        <AuthInput
          label="Contraseña"
          value={password}
          onChangeText={setPassword}
          error={errors.password}
          secureTextEntry
          editable={!loading}
          testID="login-password"
        />
        <AuthServerError message={serverError} />
        <AuthButton title="Iniciar sesión" onPress={() => void submit()} loading={loading} />
      </AuthCard>
    </AuthLayout>
  )
}
