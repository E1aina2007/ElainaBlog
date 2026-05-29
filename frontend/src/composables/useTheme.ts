import { useDark, useToggle } from '@vueuse/core'
import { type Ref } from 'vue'

let isDark: Ref<boolean>
let toggleDark: (value?: boolean) => void

export function useTheme() {
  if (!isDark) {
    isDark = useDark({
      selector: 'html',
      attribute: 'class',
      valueDark: 'dark',
      valueLight: '',
      storageKey: 'elaina-theme-dark',
    })
    toggleDark = useToggle(isDark)
  }
  return { isDark, toggleDark }
}
