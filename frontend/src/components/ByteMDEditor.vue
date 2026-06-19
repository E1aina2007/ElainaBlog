<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { Editor } from '@bytemd/vue-next'
import { useTheme } from '@/composables/useTheme'
import 'bytemd/dist/index.css'
import 'highlight.js/styles/github-dark.css'
import highlight from '@bytemd/plugin-highlight'
import gfm from '@bytemd/plugin-gfm'

const { isDark } = useTheme()

const plugins = [highlight(), gfm()]

const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    uploadImages?: (files: File[]) => Promise<{ url: string }[]>
    editorHeight?: string
  }>(),
  {
    placeholder: '',
    editorHeight: '500px',
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

function handleChange(value: string) {
  emit('update:modelValue', value)
}

// ── 为预览面板代码块注入语言标签 ──────────────────────────
let observer: MutationObserver | null = null

function injectCodeBlockHeaders(previewEl: Element) {
  const pres = previewEl.querySelectorAll('pre')
  pres.forEach((pre) => {
    if (pre.parentElement?.classList.contains('code-block-wrapper')) return

    const code = pre.querySelector('code')
    const langClass = code?.className.match(/language-(\w+)/)
    const lang = langClass?.[1] ?? ''

    const wrapper = document.createElement('div')
    wrapper.className = 'code-block-wrapper'

    const header = document.createElement('div')
    header.className = 'code-block-header'
    if (lang) {
      const langLabel = document.createElement('span')
      langLabel.className = 'code-lang'
      langLabel.textContent = lang
      header.appendChild(langLabel)
    }

    pre.parentNode?.insertBefore(wrapper, pre)
    wrapper.appendChild(header)
    wrapper.appendChild(pre)
  })
}

function setupObserver() {
  const preview = document.querySelector('.bytemd-preview')
  if (!preview) return

  // 首次立即处理
  injectCodeBlockHeaders(preview)

  // 监听子节点变化
  observer = new MutationObserver(() => {
    injectCodeBlockHeaders(preview)
  })
  observer.observe(preview, { childList: true, subtree: true })
}

onMounted(() => {
  // ByteMD 渲染是异步的，延迟启动观察者
  setTimeout(setupObserver, 300)
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<template>
  <div class="bytemd-editor-wrapper" :class="{ 'bytemd-dark': isDark }" :editor-height="editorHeight">
    <Editor
      :value="modelValue"
      :plugins="plugins"
      mode="split"
      :placeholder="placeholder"
      :upload-images="uploadImages"
      @change="handleChange"
    />
  </div>
</template>

<style scoped>
.bytemd-editor-wrapper {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.bytemd-editor-wrapper :deep(.bytemd) {
  height: v-bind(editorHeight) !important;
}

/* ── 预览面板代码块语言标签样式 ── */
:deep(.code-block-wrapper) {
  position: relative;
  margin: 1rem 0;
  border-radius: 6px;
  overflow: hidden;
}

:deep(.code-block-header) {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background: #161b22;
  border-bottom: 1px solid #30363d;
}

:deep(.code-lang) {
  font-size: 0.7rem;
  font-weight: 600;
  color: #58a6ff;
  text-transform: uppercase;
  font-family: 'Fira Code', 'Consolas', monospace;
  user-select: none;
  letter-spacing: 0.5px;
}

:deep(.code-block-wrapper pre) {
  margin: 0;
  border-radius: 0;
}

/* ── 暗黑模式 ── */
:deep(.bytemd-dark.bytemd-editor-wrapper) .bytemd,
.bytemd-editor-wrapper:deep(.bytemd-dark) .bytemd,
.bytemd-dark :deep(.bytemd) {
  color: #c9d1d9;
  border-color: #30363d;
  background-color: #0d1117;
}

.bytemd-dark :deep(.bytemd-toolbar) {
  background-color: #161b22;
  border-bottom-color: #30363d;
}

.bytemd-dark :deep(.bytemd-toolbar-icon:hover) {
  background-color: #30363d;
}

.bytemd-dark :deep(.bytemd-toolbar-tab-active) {
  color: #58a6ff;
}

.bytemd-dark :deep(.bytemd-toolbar-icon-active) {
  color: #58a6ff;
}

.bytemd-dark :deep(.bytemd-split .bytemd-preview) {
  border-left-color: #30363d;
}

.bytemd-dark :deep(.bytemd-sidebar) {
  border-left-color: #30363d;
  background-color: #0d1117;
}

.bytemd-dark :deep(.bytemd-status) {
  border-top-color: #30363d;
  color: #8b949e;
}

.bytemd-dark :deep(.bytemd-dropdown-title) {
  border-bottom-color: #30363d;
  color: #c9d1d9;
}

.bytemd-dark :deep(.bytemd-dropdown-item:hover) {
  background-color: #161b22;
}

.bytemd-dark :deep(.bytemd-preview .markdown-body) {
  color: #c9d1d9;
}

.bytemd-dark :deep(.bytemd-help) {
  color: #8b949e;
}

.bytemd-dark :deep(.bytemd-toc-active) {
  color: #58a6ff;
  background-color: #161b22;
}

/* CodeMirror editor dark mode */
.bytemd-dark :deep(.CodeMirror) {
  color: #c9d1d9;
  background-color: #0d1117;
}

.bytemd-dark :deep(.CodeMirror-gutters) {
  background-color: #0d1117;
  border-right-color: #30363d;
}

.bytemd-dark :deep(.CodeMirror-cursor) {
  border-left-color: #c9d1d9;
}

.bytemd-dark :deep(.CodeMirror-selected) {
  background-color: #264f78;
}

.bytemd-dark :deep(.CodeMirror-focused .CodeMirror-selected) {
  background-color: #264f78;
}

.bytemd-dark :deep(.CodeMirror-activeline-background) {
  background-color: #161b22;
}

.bytemd-dark :deep(.CodeMirror-line::selection),
.bytemd-dark :deep(.CodeMirror-line > span::selection),
.bytemd-dark :deep(.CodeMirror-line > span > span::selection) {
  background-color: #264f78;
}

.bytemd-dark :deep(.cm-s-default .cm-comment) {
  color: #8b949e;
}

.bytemd-dark :deep(.cm-s-default .cm-keyword) {
  color: #ff7b72;
}

.bytemd-dark :deep(.cm-s-default .cm-string) {
  color: #a5d6ff;
}

.bytemd-dark :deep(.cm-s-default .cm-number) {
  color: #79c0ff;
}

.bytemd-dark :deep(.cm-s-default .cm-def) {
  color: #d2a8ff;
}

.bytemd-dark :deep(.cm-s-default .cm-variable-2) {
  color: #ffa657;
}

.bytemd-dark :deep(.cm-s-default .cm-tag) {
  color: #7ee787;
}

.bytemd-dark :deep(.cm-s-default .cm-attribute) {
  color: #79c0ff;
}

.bytemd-dark :deep(.cm-s-default .cm-header) {
  color: #79c0ff;
}

.bytemd-dark :deep(.cm-s-default .cm-quote) {
  color: #7ee787;
}

.bytemd-dark :deep(.CodeMirror pre.CodeMirror-placeholder) {
  color: #484f58;
}
</style>
