<script setup lang="ts">
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

/* Match editor height */
.bytemd-editor-wrapper :deep(.bytemd) {
  height: v-bind(editorHeight) !important;
}
</style>
