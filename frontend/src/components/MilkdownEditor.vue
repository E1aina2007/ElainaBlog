<script setup lang="ts">
import { watch, ref, onUnmounted } from 'vue'
import { Milkdown, MilkdownProvider, useEditor } from '@milkdown/vue'
import { Crepe, CrepeFeature } from '@milkdown/crepe'
import { replaceAll } from '@milkdown/utils'
import { useTheme } from '@/composables/useTheme'

// Common base styles (layout, toolbar, ProseMirror, etc.)
import '@milkdown/crepe/theme/common/style.css'
// Theme CSS as raw strings for dynamic swapping
import nordLightCSS from '@milkdown/crepe/theme/nord.css?inline'
import nordDarkCSS from '@milkdown/crepe/theme/nord-dark.css?inline'

const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    onUpload?: (file: File) => Promise<string>
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

const { isDark } = useTheme()

// ── v-model bridge ──────────────────────────────────────────
let isInternalUpdate = false

// Keep reference to Crepe for getMarkdown()
let crepeRef: Crepe | null = null

const editorRef = useEditor((root) => {
  const featureConfigs: Record<string, unknown> = {}

  if (props.placeholder) {
    featureConfigs[CrepeFeature.Placeholder] = {
      text: props.placeholder,
      mode: 'doc' as const,
    }
  }

  if (props.onUpload) {
    featureConfigs[CrepeFeature.ImageBlock] = {
      onUpload: props.onUpload,
    }
  }

  const crepe = new Crepe({
    root,
    defaultValue: props.modelValue,
    featureConfigs,
  })

  crepeRef = crepe

  // Listen for user edits -> emit v-model
  crepe.on((listener) => {
    listener.markdownUpdated((_ctx, markdown) => {
      if (!isInternalUpdate) {
        emit('update:modelValue', markdown)
      }
    })
  })

  return crepe
})

// Watch external modelValue changes -> sync into editor
watch(
  () => props.modelValue,
  (newVal) => {
    if (!crepeRef) return
    const ed = editorRef.get()
    if (!ed) return
    const currentMd = crepeRef.getMarkdown()
    if (currentMd === newVal) return

    isInternalUpdate = true
    ed.action(replaceAll(newVal))
    setTimeout(() => {
      isInternalUpdate = false
    }, 0)
  }
)

// ── Theme CSS swapping ──────────────────────────────────────
let themeStyleEl: HTMLStyleElement | null = null

function applyTheme(dark: boolean) {
  if (themeStyleEl) {
    themeStyleEl.remove()
    themeStyleEl = null
  }
  themeStyleEl = document.createElement('style')
  themeStyleEl.setAttribute('data-milkdown-theme', '')
  themeStyleEl.textContent = dark ? nordDarkCSS : nordLightCSS
  document.head.appendChild(themeStyleEl)
}

watch(isDark, (dark) => applyTheme(dark), { immediate: true })

onUnmounted(() => {
  themeStyleEl?.remove()
})
</script>

<template>
  <div class="milkdown-editor-wrapper" :style="{ height: editorHeight }">
    <MilkdownProvider>
      <Milkdown />
    </MilkdownProvider>
  </div>
</template>

<style scoped>
.milkdown-editor-wrapper {
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}
</style>
