<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  getFriendLinkList,
  createFriendLink,
  updateFriendLink,
  deleteFriendLink,
  type FriendLink,
} from '@/api/friendlink'
import { getFaviconUrl } from '@/utils/favicon'
import toast from '@/utils/toast'

const links = ref<FriendLink[]>([])
const loading = ref(false)
const showForm = ref(false)
const editingId = ref<number | null>(null)

const form = ref({
  name: '',
  url: '',
  avatar: '',
  description: '',
  sort_order: 0,
})

const resetForm = () => {
  form.value = { name: '', url: '', avatar: '', description: '', sort_order: 0 }
  editingId.value = null
}

const fetchLinks = async () => {
  loading.value = true
  try {
    links.value = (await getFriendLinkList()) ?? []
  } catch {
    toast.error('获取友链列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  resetForm()
  showForm.value = true
}

const handleEdit = (link: FriendLink) => {
  editingId.value = link.id
  form.value = {
    name: link.name,
    url: link.url,
    avatar: link.avatar || '',
    description: link.description || '',
    sort_order: link.sort_order || 0,
  }
  showForm.value = true
}

const handleSubmit = async () => {
  if (!form.value.name.trim() || !form.value.url.trim()) {
    toast.warning('站点名称和链接不能为空')
    return
  }

  try {
    if (editingId.value) {
      await updateFriendLink({ id: editingId.value, ...form.value })
      toast.success('更新成功')
    } else {
      await createFriendLink(form.value)
      toast.success('创建成功')
    }
    showForm.value = false
    resetForm()
    await fetchLinks()
  } catch {
    toast.error('操作失败')
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('确定删除这条友链？')) return
  try {
    await deleteFriendLink(id)
    toast.success('删除成功')
    await fetchLinks()
  } catch {
    toast.error('删除失败')
  }
}

const handleCancel = () => {
  showForm.value = false
  resetForm()
}

onMounted(fetchLinks)
</script>

<template>
  <div class="friendlink-page">
    <div class="page-header">
      <h2 class="page-title">友情链接管理</h2>
      <button class="btn-primary" @click="handleAdd">+ 添加友链</button>
    </div>

    <!-- 友链列表 -->
    <div class="table-wrapper" v-if="!showForm">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="links.length === 0" class="empty-state">暂无友情链接</div>
      <table v-else class="data-table">
        <thead>
          <tr>
            <th>站点名称</th>
            <th>链接</th>
            <th>描述</th>
            <th>排序</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="link in links" :key="link.id">
            <td class="col-name">
              <div class="link-info">
                <img
                  :src="link.avatar || getFaviconUrl(link.url)"
                  class="link-avatar"
                  alt=""
                  @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <span>{{ link.name }}</span>
              </div>
            </td>
            <td class="col-url">
              <a :href="link.url.startsWith('http') ? link.url : 'https://' + link.url" target="_blank" rel="noopener noreferrer" class="link-url">
                {{ link.url }}
              </a>
            </td>
            <td class="col-desc">{{ link.description || '-' }}</td>
            <td class="col-sort">{{ link.sort_order }}</td>
            <td class="col-actions">
              <button class="btn-action edit" @click="handleEdit(link)">编辑</button>
              <button class="btn-action delete" @click="handleDelete(link.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 编辑表单 -->
    <div v-if="showForm" class="form-card">
      <h3 class="form-title">{{ editingId ? '编辑友链' : '添加友链' }}</h3>
      <div class="form-group">
        <label class="form-label">站点名称 <span class="required">*</span></label>
        <input v-model="form.name" type="text" class="form-input" placeholder="请输入站点名称" />
      </div>
      <div class="form-group">
        <label class="form-label">站点链接 <span class="required">*</span></label>
        <input v-model="form.url" type="text" class="form-input" placeholder="https://example.com" />
      </div>
      <div class="form-group">
        <label class="form-label">头像/Logo URL</label>
        <div class="avatar-input-row">
          <input v-model="form.avatar" type="text" class="form-input" placeholder="留空则自动获取网站图标" />
          <img
            v-if="form.url"
            :src="form.avatar || getFaviconUrl(form.url)"
            class="avatar-preview"
            alt=""
            @error="($event.target as HTMLImageElement).style.display = 'none'"
          />
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">站点描述</label>
        <textarea v-model="form.description" class="form-input" rows="3" placeholder="简短描述"></textarea>
      </div>
      <div class="form-group">
        <label class="form-label">排序权重</label>
        <input v-model.number="form.sort_order" type="number" class="form-input" placeholder="越大越靠前" />
      </div>
      <div class="form-actions">
        <button class="btn-outline" @click="handleCancel">取消</button>
        <button class="btn-primary" @click="handleSubmit">
          {{ editingId ? '保存修改' : '创建' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.friendlink-page {
  width: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.btn-primary {
  padding: 10px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background: var(--primary-dark);
  transform: translateY(-1px);
}

.btn-outline {
  padding: 10px 20px;
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-outline:hover {
  color: var(--text-primary);
  border-color: var(--primary-light);
}

/* 表格 */
.table-wrapper {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
  font-size: 0.875rem;
}

.data-table th {
  background: var(--bg-secondary);
  font-weight: 600;
  color: var(--text-secondary);
}

.col-name {
  font-weight: 500;
  color: var(--text-primary);
}

.link-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.link-avatar {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  object-fit: cover;
}

.link-url {
  color: var(--primary-dark);
  text-decoration: none;
  word-break: break-all;
}

.link-url:hover {
  text-decoration: underline;
}

.col-desc {
  color: var(--text-secondary);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.col-sort {
  color: var(--text-muted);
  text-align: center;
}

.col-actions {
  white-space: nowrap;
}

.btn-action {
  padding: 6px 14px;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s;
  margin-right: 8px;
}

.btn-action.edit {
  background: rgba(126, 215, 193, 0.15);
  color: var(--primary-dark);
}

.btn-action.edit:hover {
  background: rgba(126, 215, 193, 0.3);
}

.btn-action.delete {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-danger);
}

.btn-action.delete:hover {
  background: rgba(239, 68, 68, 0.2);
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 48px 0;
  color: var(--text-muted);
  font-size: 0.875rem;
}

/* 表单 */
.form-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: var(--shadow-soft);
}

.form-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.required {
  color: var(--color-danger);
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  outline: none;
  transition: border-color 0.2s;
  resize: vertical;
}

.form-input:focus {
  border-color: var(--primary);
}

.avatar-input-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-input-row .form-input {
  flex: 1;
}

.avatar-preview {
  width: 36px;
  height: 36px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
  border: 1px solid var(--border);
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .data-table th:nth-child(3),
  .data-table td:nth-child(3),
  .data-table th:nth-child(4),
  .data-table td:nth-child(4) {
    display: none;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
