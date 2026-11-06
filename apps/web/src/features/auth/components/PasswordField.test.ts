import { mount } from '@vue/test-utils'
import PasswordField from './PasswordField.vue'

function mountField(props: Record<string, unknown> = {}) {
  return mount(PasswordField, {
    props: { modelValue: 'secreta123', ...props },
  })
}

describe('PasswordField', () => {
  it('oculta la contraseña por defecto con autofill de login', () => {
    const wrapper = mountField()
    const input = wrapper.find('input')

    expect(input.attributes('type')).toBe('password')
    expect(input.attributes('name')).toBe('password')
    expect(input.attributes('autocomplete')).toBe('current-password')
    expect(input.attributes('required')).toBeDefined()
  })

  it('alterna la visibilidad con un botón accesible', async () => {
    const wrapper = mountField({ dataTest: 'pwd' })
    const toggle = wrapper.find('[data-test="pwd-toggle"]')

    expect(toggle.attributes('aria-label')).toBe('Mostrar contraseña')
    expect(toggle.attributes('type')).toBe('button')

    await toggle.trigger('click')
    expect(wrapper.find('input').attributes('type')).toBe('text')
    expect(wrapper.find('[data-test="pwd-toggle"]').attributes('aria-label')).toBe(
      'Ocultar contraseña',
    )

    await wrapper.find('[data-test="pwd-toggle"]').trigger('click')
    expect(wrapper.find('input').attributes('type')).toBe('password')
    expect(wrapper.find('[data-test="pwd-toggle"]').attributes('aria-label')).toBe(
      'Mostrar contraseña',
    )
  })

  it('vuelve a ocultarse al deshabilitarse el campo (envío en curso)', async () => {
    const wrapper = mountField()

    await wrapper.find('button').trigger('click')
    expect(wrapper.find('input').attributes('type')).toBe('text')

    await wrapper.setProps({ disabled: true })
    expect(wrapper.find('input').attributes('type')).toBe('password')
    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })

  it('emite update:modelValue al escribir', async () => {
    const wrapper = mountField({ modelValue: '' })

    await wrapper.find('input').setValue('nueva')

    expect(wrapper.emitted('update:modelValue')).toEqual([['nueva']])
  })
})
