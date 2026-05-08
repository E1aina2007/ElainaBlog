<script setup lang="ts">
import { ref, computed } from 'vue'
import { uploadAvatar } from '@/api/upload'
import toast from '@/utils/toast'

const props = defineProps<{
  modelValue: string
  size?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const size = computed(() => props.size || 120)
const isDragging = ref(false)
const isUploading = ref(false)
const previewUrl = ref('')

// 当前显示的头像
const displayAvatar = computed(() => {
  return previewUrl.value || props.modelValue || ''
})

// 选择文件
function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files && input.files[0]) {
    processFile(input.files[0])
  }
}

// 拖拽
function handleDrop(e: DragEvent) {
  isDragging.value = false
  if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
    processFile(e.dataTransfer.files[0])
  }
}

function handleDragOver() {
  isDragging.value = true
}

function handleDragLeave() {
  isDragging.value = false
}

// 处理文件：裁剪为圆形后上传
async function processFile(file: File) {
  if (!file.type.startsWith('image/')) {
    toast.warning('请选择图片文件')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    toast.warning('图片大小不能超过5MB')
    return
  }

  isUploading.value = true

  try {
    // 读取图片并裁剪为圆形
    const croppedBlob = await cropToCircle(file)
    const croppedFile = new File([croppedBlob], file.name, { type: 'image/png' })

    // 上传到头像专用目录
    const result = await uploadAvatar(croppedFile)
    previewUrl.value = result.url
    emit('update:modelValue', result.url)
    toast.success('头像上传成功')
  } catch (err: any) {
    toast.error(err?.message || '上传失败')
  } finally {
    isUploading.value = false
  }
}

// 将图片裁剪为圆形
function cropToCircle(file: File): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    const reader = new FileReader()

    reader.onload = (e) => {
      img.src = e.target?.result as string
    }
    reader.onerror = reject
    reader.readAsDataURL(file)

    img.onload = () => {
      const outputSize = 256
      const canvas = document.createElement('canvas')
      canvas.width = outputSize
      canvas.height = outputSize
      const ctx = canvas.getContext('2d')!

      // 绘制圆形裁剪路径
      ctx.beginPath()
      ctx.arc(outputSize / 2, outputSize / 2, outputSize / 2, 0, Math.PI * 2)
      ctx.closePath()
      ctx.clip()

      // 计算居中裁剪（取正方形中心区域）
      const srcSize = Math.min(img.width, img.height)
      const sx = (img.width - srcSize) / 2
      const sy = (img.height - srcSize) / 2

      ctx.drawImage(img, sx, sy, srcSize, srcSize, 0, 0, outputSize, outputSize)

      canvas.toBlob((blob) => {
        if (blob) resolve(blob)
        else reject(new Error('裁剪失败'))
      }, 'image/png')
    }
    img.onerror = () => reject(new Error('图片加载失败'))
  })
}
</script>

<template>
  <div class="avatar-upload" :style="{ width: `${size}px` }">
    <div
      class="avatar-drop-zone"
      :class="{ dragging: isDragging, uploading: isUploading }"
      @drop.prevent="handleDrop"
      @dragover.prevent="handleDragOver"
      @dragleave="handleDragLeave"
      @click="($refs.fileInput as HTMLInputElement)?.click()"
    >
      <img v-if="displayAvatar" :src="displayAvatar" class="avatar-preview" alt="头像" />
      <div v-else class="avatar-placeholder">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
          <circle cx="12" cy="7" r="4" />
        </svg>
        <span>上传头像</span>
      </div>
      <div v-if="isUploading" class="avatar-overlay">
        <div class="spinner"></div>
      </div>
      <div v-else class="avatar-hover-overlay">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
          <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z" />
          <circle cx="12" cy="13" r="4" />
        </svg>
        <span>更换头像</span>
      </div>
    </div>
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      class="file-input"
      @change="handleFileChange"
    />
  </div>
</template>

<style scoped>
.avatar-upload {
  display: inline-block;
}

.avatar-drop-zone {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  border: 3px dashed var(--border, #ddd);
  transition: border-color 0.2s, transform 0.2s;
  background: var(--bg-secondary, #f5f5f5);
}

.avatar-drop-zone:hover {
  border-color: var(--primary, #7ed7c1);
  transform: scale(1.02);
}

.avatar-drop-zone.dragging {
  border-color: var(--primary, #7ed7c1);
  background: rgba(126, 215, 193, 0.1);
}

.avatar-drop-zone.uploading {
  pointer-events: none;
}

.avatar-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: var(--text-muted, #999);
  font-size: 0.75rem;
}

.avatar-hover-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: white;
  font-size: 0.75rem;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-drop-zone:hover .avatar-hover-overlay {
  opacity: 1;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.file-input {
  display: none;
}
</style>
