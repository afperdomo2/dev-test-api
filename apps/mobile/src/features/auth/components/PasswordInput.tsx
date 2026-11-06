import { useState } from 'react'
import { Pressable, StyleSheet } from 'react-native'
import { MaterialCommunityIcons } from '@expo/vector-icons'

import { colors, spacing } from '@/core/theme/tokens'
import { AuthInput } from '@/features/auth/components/AuthInput'
import type { AuthInputProps } from '@/features/auth/components/AuthInput'

export type PasswordInputProps = Omit<
  AuthInputProps,
  'secureTextEntry' | 'rightAccessory' | 'keyboardType'
>

/**
 * Campo de contraseña con visibilidad controlada.
 *
 * - Arranca siempre oculto (`secureTextEntry`), sin persistir el estado.
 * - El botón del ojo es un control real: etiqueta accesible dinámica y
 *   área táctil de 44×44 (iOS) / 48×48 con `hitSlop` (Android).
 * - Al bloquearse el campo (p. ej. mientras se envía el formulario) vuelve
 *   a ocultarse para no dejar la contraseña en texto plano.
 */
export function PasswordInput({ editable = true, testID, ...props }: PasswordInputProps) {
  const [visible, setVisible] = useState(false)
  const [wasEditable, setWasEditable] = useState(editable)

  // Ajuste de estado en render (patrón recomendado por React): al desbloquear
  // el campo tras un envío, la contraseña vuelve a estar oculta.
  if (editable !== wasEditable) {
    setWasEditable(editable)
    if (!editable) setVisible(false)
  }

  const toggle = (
    <Pressable
      onPress={() => setVisible((current) => !current)}
      disabled={!editable}
      accessibilityRole="button"
      accessibilityLabel={visible ? 'Ocultar contraseña' : 'Mostrar contraseña'}
      accessibilityHint={
        visible
          ? 'Vuelve a ocultar el texto de la contraseña'
          : 'Muestra la contraseña en texto plano'
      }
      accessibilityState={{ disabled: !editable }}
      hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
      style={({ pressed }) => [styles.toggle, pressed ? styles.togglePressed : null]}
      testID={testID ? `${testID}-toggle` : undefined}
    >
      <MaterialCommunityIcons
        name={visible ? 'eye-off-outline' : 'eye-outline'}
        size={22}
        color={colors.mutedForeground}
        accessible={false}
      />
    </Pressable>
  )

  return (
    <AuthInput
      {...props}
      testID={testID}
      editable={editable}
      secureTextEntry={!visible}
      rightAccessory={toggle}
    />
  )
}

const styles = StyleSheet.create({
  toggle: {
    width: 48,
    height: 48,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: spacing.xs,
  },
  togglePressed: {
    opacity: 0.7,
  },
})
