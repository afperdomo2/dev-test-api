import { StyleSheet, Text, TextInput, View } from 'react-native'
import type { ReactNode } from 'react'
import type { TextInputProps } from 'react-native'

import { colors, radius, spacing, typography } from '@/core/theme/tokens'

export type AuthInputProps = {
  label: string
  value: string
  onChangeText: (text: string) => void
  error?: string
  secureTextEntry?: boolean
  keyboardType?: 'default' | 'email-address'
  editable?: boolean
  testID?: string
  textContentType?: TextInputProps['textContentType']
  autoComplete?: TextInputProps['autoComplete']
  returnKeyType?: TextInputProps['returnKeyType']
  onSubmitEditing?: TextInputProps['onSubmitEditing']
  /** Elemento alineado a la derecha dentro del campo (p. ej. el ojo de la contraseña). */
  rightAccessory?: ReactNode
}

export function AuthInput({
  label,
  value,
  onChangeText,
  error,
  secureTextEntry,
  keyboardType = 'default',
  editable = true,
  testID,
  textContentType,
  autoComplete,
  returnKeyType,
  onSubmitEditing,
  rightAccessory,
}: AuthInputProps) {
  return (
    <View style={styles.field}>
      <Text style={styles.label}>{label}</Text>
      <View style={styles.control}>
        <TextInput
          testID={testID}
          style={[
            styles.input,
            rightAccessory ? styles.inputWithAccessory : null,
            error ? styles.inputError : null,
          ]}
          value={value}
          onChangeText={onChangeText}
          secureTextEntry={secureTextEntry}
          keyboardType={keyboardType}
          autoCapitalize="none"
          autoCorrect={false}
          editable={editable}
          textContentType={textContentType}
          autoComplete={autoComplete}
          returnKeyType={returnKeyType}
          onSubmitEditing={onSubmitEditing}
          placeholderTextColor={colors.mutedForeground}
        />
        {rightAccessory ? <View style={styles.accessory}>{rightAccessory}</View> : null}
      </View>
      {error ? <Text style={styles.error}>{error}</Text> : null}
    </View>
  )
}

const styles = StyleSheet.create({
  field: {
    gap: spacing.xs,
  },
  label: {
    fontSize: typography.size.sm,
    fontWeight: typography.weight.medium,
    color: colors.foreground,
  },
  control: {
    justifyContent: 'center',
  },
  input: {
    backgroundColor: colors.background,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    minHeight: 48,
    fontSize: typography.size.md,
    color: colors.foreground,
  },
  inputWithAccessory: {
    paddingRight: spacing['2xl'],
  },
  inputError: {
    borderColor: colors.destructive,
  },
  accessory: {
    position: 'absolute',
    right: 0,
    top: 0,
    bottom: 0,
    justifyContent: 'center',
  },
  error: {
    fontSize: typography.size.xs,
    color: colors.destructive,
  },
})
