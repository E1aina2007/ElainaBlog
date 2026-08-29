<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory, toggleCategoryTop, type Category } from '@/api/category'
import toast from '@/utils/toast'

const categories = ref<Category[]>([])
const loading = ref(false)
const editingId = ref<number | null>(null)
const newCategoryName = ref('')
const editName = ref('')

const fetchCategories = async () => {
  loading.value = true
  try {
    categories.value = await getCategoryList()
  } catch (error) {
    console.error('获取分类列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  if (!newCategoryName.value.trim()) {
    toast.warning('请输入分类名称')
    return
  }
  try {
    await createCategory(newCategoryName.value.trim())
    newCategoryName.value = ''
    fetchCategories()
    toast.success('创建成功')
  } catch (error) {
    console.error('创建失败:', error)
    toast.error('创建失败')
  }
}

const startEdit = (category: Category) => {
  editingId.value = category.id
  editName.value = category.name
}

const handleUpdate = async () => {
  if (!editName.value.trim() || !editingId.value) return
  try {
    await updateCategory(editingId.value, editName.value.trim())
    editingId.value = null
    fetchCategories()
    toast.success('更新成功')
  } catch (error) {
    console.error('更新失败:', error)
    toast.error('更新失败')
  }
}

const handleDelete = async (category: Category) => {
  if (!confirm(`确定要删除分类 "${category.name}" 吗？`)) return
  try {
    await deleteCategory(category.id)
    fetchCategories()
    toast.success('删除成功')
  } catch (error) {
    console.error('删除失败:', error)
    toast.error('删除失败，可能该分类下还有文章')
  }
}

const cancelEdit = () => {
  editingId.value = null
  editName.value = ''
}

const handleToggleTop = async (category: Category) => {
  try {
    await toggleCategoryTop(category.id, !category.is_top)
    fetchCategories()
    toast.success(category.is_top ? '已取消置顶' : '已置顶')
  } catch {
    toast.error('操作失败')
  }
}

onMounted(fetchCategories)
</script>

<template>
  <div class="category-list">
    <div class="page-header">
      <h2>分类管理</h2>
    </div>

    <div class="add-section">
      <div class="add-form">
        <input
          v-model="newCategoryName"
          type="text"
          placeholder="输入新分类名称..."
          @keyup.enter="handleCreate"
        />
        <button class="btn-primary" @click="handleCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          添加分类
        </button>
      </div>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>分类名称</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="category in categories" :key="category.id">
            <td>{{ category.id }}</td>
            <td>
              <div v-if="editingId === category.id" class="edit-form">
                <input
                  v-model="editName"
                  type="text"
                  @keyup.enter="handleUpdate"
                  @keyup.esc="cancelEdit"
                  v-focus
                />
                <button class="btn-icon success" @click="handleUpdate">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </button>
                <button class="btn-icon" @click="cancelEdit">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>
              <span v-else class="category-name">
                {{ category.name }}
                <span v-if="category.is_top" class="top-badge">置顶</span>
              </span>
            </td>
            <td>{{ category.created_at ? new Date(category.created_at).toLocaleDateString('zh-CN') : '-' }}</td>
            <td class="actions">
              <button
                v-if="editingId !== category.id"
                class="action-btn"
                :class="{ pinned: category.is_top }"
                @click="handleToggleTop(category)"
                :title="category.is_top ? '取消置顶' : '置顶'"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" :fill="category.is_top ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                  <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"></path>
                </svg>
              </button>
              <button
                v-if="editingId !== category.id"
                class="action-btn edit"
                @click="startEdit(category)"
                title="编辑"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                </svg>
              </button>
              <button class="action-btn delete" @click="handleDelete(category)" title="删除">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </button>
            </td>
          </tr>
          <tr v-if="categories.length === 0">
            <td colspan="4" class="empty-cell">
              <div class="empty-state">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                </svg>
                <p>暂无分类，请添加</p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.category-list {
  max-width: 800px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  color: var(--text-primary);
}

.add-section {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-card);
}

.add-form {
  display: flex;
  gap: 12px;
}

.add-form input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid var(--input-border);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.add-form input:focus {
  border-color: var(--color-indigo);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--color-indigo);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background: var(--color-indigo-hover);
}

.table-container {
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: var(--shadow-card);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 16px;
  text-align: left;
  font-size: 14px;
}

.data-table th {
  background: var(--toolbar-bg);
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--input-border);
}

.data-table td {
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border);
}

.data-table tr:hover td {
  background: var(--toolbar-bg);
}

.category-name {
  font-weight: 500;
  color: var(--text-primary);
}

.edit-form {
  display: flex;
  align-items: center;
  gap: 8px;
}

.edit-form input {
  padding: 8px 12px;
  border: 1px solid var(--color-indigo);
  border-radius: 6px;
  font-size: 14px;
  outline: none;
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: var(--border);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: var(--input-border);
}

.btn-icon.success {
  background: color-mix(in srgb, var(--color-success) 10%, transparent);
  color: var(--color-success);
}

.btn-icon.success:hover {
  background: color-mix(in srgb, var(--color-success) 15%, transparent);
}

.actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: var(--border);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--input-border);
  color: var(--text-primary);
}

.action-btn.edit:hover {
  background: color-mix(in srgb, var(--color-indigo) 10%, transparent);
  color: var(--color-indigo);
}

.action-btn.delete:hover {
  background: color-mix(in srgb, var(--color-danger) 10%, transparent);
  color: var(--color-danger);
}

.action-btn.pinned {
  background: color-mix(in srgb, var(--color-warning) 15%, transparent);
  color: var(--color-warning);
}

.action-btn.pinned:hover {
  background: color-mix(in srgb, var(--color-warning) 25%, transparent);
}

.top-badge {
  display: inline-block;
  margin-left: 6px;
  padding: 1px 6px;
  font-size: 11px;
  font-weight: 500;
  color: var(--color-warning);
  background: color-mix(in srgb, var(--color-warning) 10%, transparent);
  border-radius: 4px;
  vertical-align: middle;
}

.empty-cell {
  padding: 60px 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--text-muted);
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}
</style>
