/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '*.css?inline' {
  const css: string
  export default css
}

// @bytemd/vue-next type override: Editor is a runtime value, not just a type
declare module '@bytemd/vue-next' {
  import type { DefineComponent } from 'vue'
  export const Editor: DefineComponent<{
    value?: string
    plugins?: any[]
    sanitize?: (schema: any) => any
    remarkRehype?: Record<string, unknown>
    mode?: string
    previewDebounce?: number
    placeholder?: string
    editorConfig?: Record<string, unknown>
    locale?: Record<string, unknown>
    uploadImages?: (files: File[]) => Promise<{ url: string; alt?: string; title?: string }[]>
  }, {}, any>
  export const Viewer: DefineComponent<{ value?: string; plugins?: any[] }, {}, any>
}
