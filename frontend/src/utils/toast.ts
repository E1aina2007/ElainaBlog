import { createApp } from 'vue'
import ToastComponent from '@/components/Toast.vue'

type ToastType = 'success' | 'error' | 'warning' | 'info'

interface ToastOptions {
  message: string
  type?: ToastType
  duration?: number
}

// Toast容器引用
let toastContainer: HTMLDivElement | null = null

// 获取或创建Toast容器
function getContainer(): HTMLDivElement {
  if (!toastContainer) {
    toastContainer = document.createElement('div')
    toastContainer.id = 'toast-container'
    toastContainer.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      z-index: 10000;
      display: flex;
      flex-direction: column;
      gap: 12px;
      pointer-events: none;
    `
    document.body.appendChild(toastContainer)
  }
  return toastContainer
}

// 显示Toast
function showToast(options: ToastOptions | string): void {
  const container = getContainer()

  // 支持直接传入字符串
  const opts: ToastOptions = typeof options === 'string'
    ? { message: options }
    : options

  // 创建Toast实例容器
  const toastWrapper = document.createElement('div')
  toastWrapper.style.pointerEvents = 'auto'
  container.appendChild(toastWrapper)

  // 创建Vue应用实例
  const app = createApp(ToastComponent, {
    message: opts.message,
    type: opts.type || 'info',
    duration: opts.duration || 3000,
    onClose: () => {
      // 销毁实例并移除DOM
      app.unmount()
      container.removeChild(toastWrapper)

      // 如果容器为空，移除容器
      if (container.children.length === 0) {
        document.body.removeChild(container)
        toastContainer = null
      }
    }
  })

  // 挂载到容器
  app.mount(toastWrapper)
}

// 便捷方法
export const toast = {
  success(message: string, duration?: number) {
    showToast({ message, type: 'success', duration })
  },
  error(message: string, duration?: number) {
    showToast({ message, type: 'error', duration })
  },
  warning(message: string, duration?: number) {
    showToast({ message, type: 'warning', duration })
  },
  info(message: string, duration?: number) {
    showToast({ message, type: 'info', duration })
  }
}

export default toast
