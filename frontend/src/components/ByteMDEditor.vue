<script setup lang="ts">
import { watch, onMounted, nextTick } from 'vue'
import { Editor } from '@bytemd/vue-next'
import 'bytemd/dist/index.css'
import 'highlight.js/styles/github-dark.css'
import highlight from '@bytemd/plugin-highlight'
import gfm from '@bytemd/plugin-gfm'

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
function injectCodeBlockHeaders() {
  // ByteMD 预览面板
  const preview = document.querySelector('.bytemd-preview')
  if (!preview) return

  const pres = preview.querySelectorAll('pre')
  pres.forEach((pre) => {
    // 跳过已经处理过的
    if (pre.parentElement?.classList.contains('code-block-wrapper')) return

    const code = pre.querySelector('code')
    const langClass = code?.className.match(/language-(\w+)/)
    const lang = langClass?.[1] ?? ''

    // 创建包裹容器
    const wrapper = document.createElement('div')
    wrapper.className = 'code-block-wrapper'

    // 创建语言标签头
    const header = document.createElement('div')
    header.className = 'code-block-header'
    if (lang) {
      const langLabel = document.createElement('span')
      langLabel.className = 'code-lang'
      langLabel.textContent = lang
      header.appendChild(langLabel)
    }

    // 替换 DOM
    pre.parentNode?.insertBefore(wrapper, pre)
    wrapper.appendChild(header)
    wrapper.appendChild(pre)
  })
}

// 当内容变化时重新注入
watch(() => props.modelValue, () => nextTick(injectCodeBlockHeaders))

onMounted(() => {
  // 首次渲染后注入
  nextTick(injectCodeBlockHeaders)
})
</script>

<template>
  <div class="bytemd-editor-wrapper" :editor-height="editorHeight">
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
}

:deep(.code-block-header) {
  display: flex;
  align-items: center;
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.06);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px 6px 0 0;
}

:deep(.code-lang) {
  font-size: 0.75rem;
  color: #8b949e;
  text-transform: uppercase;
  font-family: 'Fira Code', 'Consolas', monospace;
  user-select: none;
}

:deep(.code-block-wrapper pre) {
  margin: 0;
  border-radius: 0 0 6px 6px;
}
</style>
