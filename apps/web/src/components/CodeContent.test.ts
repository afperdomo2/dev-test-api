/* eslint-disable no-useless-escape */
import { mount } from '@vue/test-utils'
import CodeContent from './CodeContent.vue'

vi.mock('@/directives/highlight', () => ({
  vHighlight: {},
}))

describe('CodeContent', () => {
  it('renders plain text when no code blocks', () => {
    const wrapper = mount(CodeContent, {
      props: { text: 'Hello world' },
    })
    expect(wrapper.text()).toContain('Hello world')
  })

  it('renders code block', () => {
    const text = 'before\n```javascript\nconsole.log(\"hi\")\n```\nafter'
    const wrapper = mount(CodeContent, {
      props: { text },
    })
    expect(wrapper.text()).toContain('before')
    expect(wrapper.text()).toContain('console.log')
    expect(wrapper.text()).toContain('after')
  })

  it('renders multiple blocks', () => {
    const text = 'a ```go\nfmt.Println()\n``` b ```python\nprint()\n``` c'
    const wrapper = mount(CodeContent, {
      props: { text },
    })
    expect(wrapper.html()).toContain('fmt.Println')
    expect(wrapper.html()).toContain('print()')
  })

  it('handles bold and inline code', () => {
    const wrapper = mount(CodeContent, {
      props: { text: '**bold** and `inline`' },
    })
    expect(wrapper.html()).toContain('<strong>bold</strong>')
    expect(wrapper.html()).toContain('<code')
    expect(wrapper.html()).toContain('inline')
  })

  it('handles text with no markers as single block', () => {
    const wrapper = mount(CodeContent, {
      props: { text: 'just text' },
    })
    expect(wrapper.text()).toContain('just text')
  })
})
