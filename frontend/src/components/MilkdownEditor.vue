<script setup lang="ts">
import { watch, ref, onUnmounted, defineComponent, h } from 'vue'
import { Milkdown, MilkdownProvider, useEditor } from '@milkdown/vue'
import { Crepe, CrepeFeature } from '@milkdown/crepe'
import { replaceAll } from '@milkdown/utils'
import { useTheme } from '@/composables/useTheme'
import { Editor } from '@milkdown/kit/core'

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

// ── v-model bridge refs (shared with inner component) ────────
let isInternalUpdate = false
let crepeRef: Crepe | null = null
const editorGetter = ref<(() => Editor | undefined) | null>(null)

// ── Inner editor component (must be child of MilkdownProvider) ──
const MilkdownInner = defineComponent({
  name: 'MilkdownInner',
  setup() {
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

      crepe.on((listener) => {
        listener.markdownUpdated((_ctx, markdown) => {
          if (!isInternalUpdate) {
            emit('update:modelValue', markdown)
          }
        })
      })

      return crepe
    })

    // Expose the editor getter to the outer component
    editorGetter.value = () => editorRef.get()

    return () => h(Milkdown)
  },
})

// ── Watch external modelValue changes -> sync into editor ──
watch(
  () => props.modelValue,
  (newVal) => {
    if (!crepeRef) return
    const getEditor = editorGetter.value
    if (!getEditor) return
    const ed = getEditor()
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
</script>

<template>
  <div class="milkdown-editor-wrapper" :style="{ height: editorHeight }">
    <MilkdownProvider>
      <MilkdownInner />
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
