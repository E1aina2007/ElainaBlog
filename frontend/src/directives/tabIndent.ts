import type { Directive } from 'vue'

function handleKeydown(e: KeyboardEvent) {
  if (e.key !== 'Tab') return
  e.preventDefault()

  const el = e.target as HTMLTextAreaElement
  const start = el.selectionStart
  const end = el.selectionEnd
  const text = el.value

  if (e.shiftKey) {
    const lineStart = text.lastIndexOf('\n', start - 1) + 1
    if (text.substring(lineStart, lineStart + 2) === '  ') {
      el.value = text.substring(0, lineStart) + text.substring(lineStart + 2)
      el.selectionStart = Math.max(start - 2, lineStart)
      el.selectionEnd = Math.max(end - 2, lineStart)
    }
  } else {
    el.value = text.substring(0, start) + '  ' + text.substring(end)
    el.selectionStart = el.selectionEnd = start + 2
  }

  el.dispatchEvent(new Event('input', { bubbles: true }))
}

export const vTabIndent: Directive = {
  mounted(el) {
    el.addEventListener('keydown', handleKeydown)
  },
  unmounted(el) {
    el.removeEventListener('keydown', handleKeydown)
  },
}
