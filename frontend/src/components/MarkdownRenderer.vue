<script setup lang="ts">
import { computed, ref, onMounted, watch, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

export interface TocItem {
  level: number
  text: string
  id: string
}

interface Props {
  content: string
}

const props = defineProps<Props>()

const markdownBody = ref<HTMLElement>()

const toc = ref<TocItem[]>([])

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/<[^>]*>/g, '')
    .replace(/[^\w一-鿿\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

// 清洗标题文本（去除 markdown 格式字符），与 extractToc 和 heading_open 渲染器共用同一逻辑
function cleanHeadingText(raw: string): string {
  return raw.replace(/[*_`~[\]()]/g, '').trim()
}

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight(str: string, lang: string) {
    const escapedContent = md.utils.escapeHtml(str)
    const copyBtn = `<button class="code-copy-btn" data-code="${escapedContent}"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg> 复制</button>`
    if (lang && hljs.getLanguage(lang)) {
      try {
        const result = hljs.highlight(str, { language: lang, ignoreIllegals: true })
        return `<pre class="hljs-code-block"><div class="code-block-wrapper"><div class="code-block-header"><span class="code-lang">${lang}</span>${copyBtn}</div><code class="hljs language-${lang}">${result.value}</code></div></pre>`
      } catch {}
    }
    return `<pre class="hljs-code-block"><div class="code-block-wrapper"><div class="code-block-header">${copyBtn}</div><code class="hljs">${escapedContent}</code></div></pre>`
  },
})

// 图片懒加载：为所有 <img> 添加 loading="lazy" decoding="async"
const defaultImageRenderer = md.renderer.rules.image!
md.renderer.rules.image = function (tokens, idx, options, env, self) {
  const token = tokens[idx]
  if (token) {
    token.attrSet('loading', 'lazy')
    token.attrSet('decoding', 'async')
  }
  return defaultImageRenderer(tokens, idx, options, env, self)
}

// 提取 TOC
function extractToc(content: string): TocItem[] {
  const lines = content.split('\n')
  const items: TocItem[] = []
  const slugCount: Record<string, number> = {}
  let inCodeBlock = false

  for (const line of lines) {
    // 检测围栏代码块边界（``` 或 ~~~），允许后面跟随语言标识等 info string
    if (/^\s*`{3,}/.test(line) || /^\s*~{3,}/.test(line)) {
      inCodeBlock = !inCodeBlock
      continue
    }
    // 跳过代码块内的行
    if (inCodeBlock) continue

    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (!match || !match[1] || !match[2]) continue
    const level = match[1].length
    const text = cleanHeadingText(match[2])
    let id = slugify(text)
    if (!id) continue
    if (slugCount[id] !== undefined) {
      slugCount[id]!++
      id = `${id}-${slugCount[id]}`
    } else {
      slugCount[id] = 0
    }
    items.push({ level, text, id })
  }
  return items
}

// 生成 heading ID 映射
function buildHeadingIdMap(content: string): Record<string, string> {
  const items = extractToc(content)
  const map: Record<string, string> = {}
  for (const item of items) {
    map[item.text] = item.id
  }
  return map
}

const headingIdMap = ref<Record<string, string>>({})

const renderedContent = computed(() => md.render(props.content || ''))

// content 变化时同步重建 TOC 与 heading ID 映射（渲染器在 md.render 期间读取，须先于渲染更新）
watch(
  () => props.content,
  (content) => {
    const c = content || ''
    headingIdMap.value = buildHeadingIdMap(c)
    toc.value = extractToc(c)
  },
  { immediate: true },
)

// 自定义 heading 渲染，添加 ID
const defaultHeadingOpen = md.renderer.rules.heading_open || function (tokens, idx, options, _env, self) {
  return self.renderToken(tokens, idx, options)
}

md.renderer.rules.heading_open = function (tokens, idx, options, env, self) {
  const token = tokens[idx]
  if (!token) return defaultHeadingOpen(tokens, idx, options, env, self)
  const nextToken = tokens[idx + 1]
  if (nextToken && nextToken.type === 'inline') {
    // 收集所有可见文本（包括 code_inline），与 extractToc 的文本提取保持一致
    const rawText = nextToken.children
      ?.filter((t) => t.type === 'text' || t.type === 'code_inline')
      .map((t) => t.content)
      .join('') || ''
    const text = cleanHeadingText(rawText)
    const id = headingIdMap.value[text]
    if (id) {
      token.attrSet('id', id)
    }
  }
  return defaultHeadingOpen(tokens, idx, options, env, self)
}

function copyToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard) {
    return navigator.clipboard.writeText(text).then(() => true).catch(() => false)
  }
  try {
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.cssText = 'position:fixed;left:-9999px;top:-9999px;opacity:0'
    document.body.appendChild(ta)
    ta.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(ta)
    return Promise.resolve(ok)
  } catch {
    return Promise.resolve(false)
  }
}

const COPY_SVG = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg> 复制'

function bindCopyButtons() {
  if (!markdownBody.value) return
  const buttons = markdownBody.value.querySelectorAll('.code-copy-btn')
  buttons.forEach((btn) => {
    const el = btn as HTMLElement
    if (el.dataset.bound) return
    el.dataset.bound = '1'
    el.addEventListener('click', () => {
      const code = el.getAttribute('data-code')
      if (!code) return
      copyToClipboard(code).then((ok) => {
        if (ok) {
          el.textContent = '已复制'
          setTimeout(() => { el.innerHTML = COPY_SVG }, 2000)
        }
      })
    })
  })
}

watch(renderedContent, () => nextTick(bindCopyButtons))
onMounted(() => nextTick(bindCopyButtons))

defineExpose({ toc })
</script>

<template>
  <div ref="markdownBody" class="markdown-body" v-html="renderedContent" />
</template>

<style scoped>
.markdown-body :deep(h1) {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 1.5rem 0 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}

.markdown-body :deep(h2) {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 1.25rem 0 0.875rem;
}

.markdown-body :deep(h3) {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 1rem 0 0.75rem;
}

.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0.875rem 0 0.625rem;
}

.markdown-body :deep(p) {
  font-size: 1rem;
  line-height: 1.8;
  color: var(--text-secondary);
  margin: 0.75rem 0;
}

.markdown-body :deep(a) {
  color: var(--primary);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s ease;
}

.markdown-body :deep(a:hover) {
  border-bottom-color: var(--primary);
}

.markdown-body :deep(strong) {
  font-weight: 600;
  color: var(--text-primary);
}

.markdown-body :deep(em) {
  font-style: italic;
  color: var(--text-secondary);
}

.markdown-body :deep(del) {
  text-decoration: line-through;
  color: var(--text-muted);
}

.markdown-body :deep(blockquote) {
  margin: 1rem 0;
  padding: 0.75rem 1rem;
  border-left: 4px solid var(--primary);
  background: var(--bg-secondary);
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
}

.markdown-body :deep(blockquote p) {
  margin: 0;
  color: var(--text-secondary);
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 0.75rem 0;
  padding-left: 1.5rem;
  list-style: revert;
}

.markdown-body :deep(li) {
  margin: 0.375rem 0;
  color: var(--text-secondary);
  line-height: 1.7;
}

/* ── 行内代码 ── */
.markdown-body :deep(code) {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 0.875rem;
  background: var(--bg-secondary);
  padding: 0.15rem 0.4rem;
  border-radius: var(--radius-sm);
  color: var(--color-danger);
}

/* ── 代码块容器（pre 包裹整个代码块） ── */
.markdown-body :deep(.hljs-code-block) {
  background: #0d1117;
  margin: 1rem 0;
  padding: 0;
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* ── 代码块内部 wrapper ── */
.markdown-body :deep(.code-block-wrapper) {
  position: relative;
}

/* ── 代码块语言标签头 ── */
.markdown-body :deep(.code-block-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background: #161b22;
  border-bottom: 1px solid #30363d;
}

.markdown-body :deep(.code-block-header .code-lang) {
  font-size: 0.7rem;
  font-weight: 600;
  color: #58a6ff;
  text-transform: uppercase;
  font-family: 'Fira Code', 'Consolas', monospace;
  user-select: none;
  letter-spacing: 0.5px;
}

.markdown-body :deep(.code-block-header .code-copy-btn) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: transparent;
  color: #8b949e;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-family: 'Fira Code', 'Consolas', monospace;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.markdown-body :deep(.code-block-header .code-copy-btn:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: #e6edf3;
}

/* ── 代码块内容区 ── */
.markdown-body :deep(.hljs-code-block code) {
  display: block;
  background: transparent;
  padding: 1rem;
  margin: 0;
  border-radius: 0;
  color: #e6edf3;
  line-height: 1.6;
  font-size: 0.875rem;
  overflow-x: auto;
}

.markdown-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-md);
  margin: 1rem 0;
}

.markdown-body :deep(hr) {
  border: none;
  height: 1px;
  background: var(--border);
  margin: 1.5rem 0;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
  font-size: 0.9375rem;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  padding: 0.625rem 0.875rem;
  border: 1px solid var(--border);
  text-align: left;
}

.markdown-body :deep(th) {
  background: var(--bg-secondary);
  font-weight: 600;
  color: var(--text-primary);
}

.markdown-body :deep(td) {
  color: var(--text-secondary);
}

.markdown-body :deep(input[type="checkbox"]) {
  margin-right: 0.375rem;
}
</style>

<!-- 日间模式：代码块亮色覆盖 -->
<style>
html:not(.dark) .markdown-body .hljs-code-block {
  background: #f6f8fa;
  border: 1px solid #d0d7de;
}

html:not(.dark) .markdown-body .code-block-header {
  background: #e8eaed;
  border-bottom-color: #d0d7de;
}

html:not(.dark) .markdown-body .code-block-header .code-lang {
  color: #57606a;
}

html:not(.dark) .markdown-body .code-block-header .code-copy-btn {
  color: #57606a;
}

html:not(.dark) .markdown-body .code-block-header .code-copy-btn:hover {
  background: #d0d7de;
  color: #24292e;
}

html:not(.dark) .markdown-body .hljs-code-block code {
  color: #24292e;
}

/* 日间模式语法高亮色 */
html:not(.dark) .markdown-body .hljs-keyword { color: #d73a49; }
html:not(.dark) .markdown-body .hljs-string { color: #032f62; }
html:not(.dark) .markdown-body .hljs-number { color: #005cc5; }
html:not(.dark) .markdown-body .hljs-comment { color: #6a737d; }
html:not(.dark) .markdown-body .hljs-function { color: #6f42c1; }
html:not(.dark) .markdown-body .hljs-title { color: #6f42c1; }
html:not(.dark) .markdown-body .hljs-built_in { color: #005cc5; }
html:not(.dark) .markdown-body .hljs-type { color: #005cc5; }
html:not(.dark) .markdown-body .hljs-literal { color: #005cc5; }
html:not(.dark) .markdown-body .hljs-attr { color: #005cc5; }
html:not(.dark) .markdown-body .hljs-variable { color: #e36209; }
html:not(.dark) .markdown-body .hljs-params { color: #24292e; }
html:not(.dark) .markdown-body .hljs-tag { color: #22863a; }
html:not(.dark) .markdown-body .hljs-name { color: #22863a; }
html:not(.dark) .markdown-body .hljs-selector-class { color: #6f42c1; }
html:not(.dark) .markdown-body .hljs-meta { color: #6a737d; }
html:not(.dark) .markdown-body .hljs-regexp { color: #032f62; }
html:not(.dark) .markdown-body .hljs-subst { color: #24292e; }
</style>
