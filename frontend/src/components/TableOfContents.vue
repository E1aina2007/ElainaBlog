<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

export interface TocItem {
  level: number
  text: string
  id: string
}

interface Props {
  items: TocItem[]
}

const props = defineProps<Props>()

const activeId = ref('')
let observer: IntersectionObserver | null = null
const showPanel = ref(false)

function initObserver() {
  if (observer) {
    observer.disconnect()
  }

  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeId.value = entry.target.id
        }
      }
    },
    { rootMargin: '-20% 0px -60% 0px', threshold: 0 }
  )

  for (const item of props.items) {
    const el = document.getElementById(item.id)
    if (el) observer.observe(el)
  }
}

watch(() => props.items, () => {
  if (props.items.length > 0) {
    setTimeout(initObserver, 100)
  }
}, { immediate: true })

onMounted(() => {
  if (props.items.length > 0) {
    setTimeout(initObserver, 100)
  }
})

onUnmounted(() => {
  observer?.disconnect()
})

function scrollTo(id: string) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    showPanel.value = false
  }
}
</script>

<template>
  <div v-if="items.length > 0">
    <!-- 移动端 TOC 按钮 -->
    <button class="toc-mobile-btn hide-desktop" @click="showPanel = !showPanel">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="3" y1="6" x2="21" y2="6" />
        <line x1="3" y1="12" x2="15" y2="12" />
        <line x1="3" y1="18" x2="18" y2="18" />
      </svg>
    </button>

    <!-- 移动端面板 -->
    <div v-if="showPanel" class="toc-mobile-overlay hide-desktop" @click.self="showPanel = false">
      <div class="toc-mobile-panel">
        <div class="toc-mobile-title">目录</div>
        <div class="toc-mobile-list">
          <div
            v-for="(item, index) in items"
            :key="index"
            class="toc-mobile-item"
            :class="{ active: activeId === item.id }"
            :style="{ paddingLeft: (item.level - 1) * 12 + 12 + 'px' }"
            @click="scrollTo(item.id)"
          >
            {{ item.text }}
          </div>
        </div>
      </div>
    </div>

    <!-- 桌面端 TOC -->
    <div class="toc hide-mobile">
      <div class="toc-title">目录</div>
      <div class="toc-content">
        <div
          v-for="(item, index) in items"
          :key="index"
          class="toc-item"
          :class="{ active: activeId === item.id }"
          :style="{ paddingLeft: (item.level - 1) * 12 + 'px' }"
          @click="scrollTo(item.id)"
        >
          {{ item.text }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toc {
  position: fixed;
  top: 80px;
  right: calc((100vw - 800px) / 2 - 240px);
  width: 200px;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
  padding: 16px 0;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  z-index: 10;
}

@media (max-width: 1280px) {
  .toc {
    right: 20px;
  }
}

.toc-title {
  padding: 0 16px 12px;
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border);
  margin-bottom: 8px;
}

.toc-content {
  padding: 0 8px;
}

.toc-item {
  padding: 6px 8px;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
  border-left: 2px solid transparent;
  transition: all 0.2s ease;
  line-height: 1.4;
}

.toc-item:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.toc-item.active {
  color: var(--primary);
  border-left-color: var(--primary);
  background: var(--primary-lighter);
  font-weight: 500;
}

/* 移动端按钮 */
.toc-mobile-btn {
  position: fixed;
  bottom: 90px;
  right: 20px;
  width: 44px;
  height: 44px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 50%;
  box-shadow: var(--shadow-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 101;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.toc-mobile-btn:hover {
  color: var(--primary);
  border-color: var(--primary);
}

/* 移动端面板 */
.toc-mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 200;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.toc-mobile-panel {
  background: var(--bg-card);
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  width: 100%;
  max-height: 60vh;
  overflow-y: auto;
  padding: 20px;
}

.toc-mobile-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  text-align: center;
}

.toc-mobile-item {
  padding: 10px 12px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.toc-mobile-item:active {
  background: var(--bg-secondary);
  color: var(--primary);
}

.toc-mobile-item.active {
  color: var(--primary);
  font-weight: 500;
  background: var(--primary-lighter);
}
</style>
