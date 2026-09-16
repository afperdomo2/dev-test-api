<script setup lang="ts">
import { computed } from 'vue'

interface TextBlock {
  type: 'text' | 'code'
  lang: string
  text: string
}

interface Props {
  text: string
}

const props = defineProps<Props>()

const codeBlockRegex = /```(\w*)\n([\s\S]*?)```/g

const blocks = computed<Array<TextBlock>>(() => {
  const result: Array<TextBlock> = []
  let lastIndex = 0
  let match: RegExpExecArray | null

  const regex = new RegExp(codeBlockRegex.source, 'g')

  while ((match = regex.exec(props.text)) !== null) {
    if (match.index > lastIndex) {
      result.push({ type: 'text', lang: '', text: props.text.slice(lastIndex, match.index) })
    }
    result.push({ type: 'code', lang: match[1] || 'plaintext', text: match[2] || '' })
    lastIndex = match.index + match[0].length
  }

  if (lastIndex < props.text.length) {
    result.push({ type: 'text', lang: '', text: props.text.slice(lastIndex) })
  }

  if (result.length === 0) {
    result.push({ type: 'text', lang: '', text: props.text })
  }

  return result
})
</script>

<script lang="ts">
function renderInline(text: string): string {
  let html = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-grey-lighten-3 pa-1 rounded text-caption">$1</code>',
  )

  return html
}
</script>

<template>
  <template v-for="(block, i) in blocks" :key="i">
    <div v-if="block.type === 'code'" class="rounded pa-3 bg-grey-lighten-4 overflow-auto mb-2">
      <pre class="text-caption"><code v-highlight="block.lang">{{ block.text }}</code></pre>
    </div>
    <span v-else class="text-body-1" v-html="renderInline(block.text)" />
  </template>
</template>
