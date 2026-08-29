<script setup lang="ts">
import { Editor } from '@bytemd/vue-next'
import { useTheme } from '@/composables/useTheme'
import 'bytemd/dist/index.css'
import 'highlight.js/styles/github-dark.css'
import highlight from '@bytemd/plugin-highlight'
import gfm from '@bytemd/plugin-gfm'

const { isDark } = useTheme()

const plugins = [highlight(), gfm()]

withDefaults(
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

/* ── 暗黑模式 ── */
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

.bytemd-dark :deep(.bytemd-toolbar-tab-active),
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
}

.bytemd-dark :deep(.bytemd-dropdown-item:hover) {
  background-color: #161b22;
}

.bytemd-dark :deep(.bytemd-preview .markdown-body) {
  color: #c9d1d9;
}

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

.bytemd-dark :deep(.CodeMirror-selected),
.bytemd-dark :deep(.CodeMirror-focused .CodeMirror-selected),
.bytemd-dark :deep(.CodeMirror-line::selection),
.bytemd-dark :deep(.CodeMirror-line > span::selection) {
  background-color: #264f78;
}

.bytemd-dark :deep(.CodeMirror-activeline-background) {
  background-color: #161b22;
}

.bytemd-dark :deep(.cm-s-default .cm-comment) { color: #8b949e; }
.bytemd-dark :deep(.cm-s-default .cm-keyword) { color: #ff7b72; }
.bytemd-dark :deep(.cm-s-default .cm-string) { color: #a5d6ff; }
.bytemd-dark :deep(.cm-s-default .cm-number) { color: #79c0ff; }
.bytemd-dark :deep(.cm-s-default .cm-def) { color: #d2a8ff; }
.bytemd-dark :deep(.cm-s-default .cm-tag) { color: #7ee787; }
.bytemd-dark :deep(.cm-s-default .cm-attribute) { color: #79c0ff; }

.bytemd-dark :deep(.CodeMirror pre.CodeMirror-placeholder) {
  color: #484f58;
}

/* 暗黑模式预览区代码块：与背景色区分 */
.bytemd-dark :deep(.markdown-body pre) {
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 6px;
}

.bytemd-dark :deep(.markdown-body code) {
  background: #343941;
}
</style>
