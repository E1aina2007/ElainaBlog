<script setup lang="ts">
const props = defineProps<{
  currentPage: number
  totalPages: number
  total?: number
}>()

const emit = defineEmits<{
  (e: 'change', page: number): void
}>()

const goToPage = (page: number) => {
  if (page >= 1 && page <= props.totalPages) {
    emit('change', page)
  }
}
</script>

<template>
  <div class="pagination">
    <!-- 上一页 -->
    <button
      class="page-btn"
      :disabled="currentPage === 1"
      @click="goToPage(currentPage - 1)"
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6"></polyline>
      </svg>
    </button>

    <!-- 页码 -->
    <div class="page-numbers">
      <button
        v-for="page in totalPages"
        :key="page"
        class="page-number"
        :class="{ active: page === currentPage }"
        @click="goToPage(page)"
      >
        {{ page }}
      </button>
    </div>

    <!-- 下一页 -->
    <button
      class="page-btn"
      :disabled="currentPage === totalPages"
      @click="goToPage(currentPage + 1)"
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6"></polyline>
      </svg>
    </button>

    <!-- 总数量显示 -->
    <span v-if="props.total && props.total > 0" class="total-info">
      共 {{ props.total }} 条
    </span>
  </div>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 24px 0;
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.page-btn:hover:not(:disabled) {
  color: var(--primary);
  border-color: var(--primary);
  background: var(--bg-secondary);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 6px;
}

.page-number {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  padding: 0 10px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.page-number:hover {
  color: var(--primary);
  border-color: var(--primary);
  background: var(--bg-secondary);
}

.page-number.active {
  color: white;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  border-color: var(--primary);
  box-shadow: var(--shadow-soft);
}

.total-info {
  margin-left: 16px;
  font-size: 13px;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .pagination {
    gap: 4px;
  }

  .page-btn,
  .page-number {
    width: 32px;
    height: 32px;
    font-size: 13px;
  }

  .total-info {
    display: none;
  }
}
</style>