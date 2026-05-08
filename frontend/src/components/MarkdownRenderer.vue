<script setup lang="ts">
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

interface Props {
  content: string
}

const props = defineProps<Props>()

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight(str: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const result = hljs.highlight(str, { language: lang, ignoreIllegals: true })
        return `<pre class="hljs-code-block"><div class="code-lang">${lang}</div><code class="hljs language-${lang}">${result.value}</code></pre>`
      } catch (_) {}
    }
    const escaped = md.utils.escapeHtml(str)
    return `<pre class="hljs-code-block"><code class="hljs">${escaped}</code></pre>`
  },
})

const renderedContent = computed(() => md.render(props.content || ''))
</script>

<template>
  <div class="markdown-body" v-html="renderedContent" />
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

.markdown-body :deep(code) {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 0.875rem;
  background: var(--bg-secondary);
  padding: 0.125rem 0.375rem;
  border-radius: var(--radius-sm);
  color: var(--primary-dark);
}

.markdown-body :deep(.hljs-code-block) {
  position: relative;
  background: #0d1117;
  padding: 1rem;
  border-radius: var(--radius-md);
  overflow-x: auto;
  margin: 1rem 0;
}

.markdown-body :deep(.hljs-code-block .code-lang) {
  position: absolute;
  top: 0;
  right: 0;
  padding: 2px 10px;
  font-size: 0.75rem;
  color: #8b949e;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 0 var(--radius-md) 0 var(--radius-sm);
  text-transform: uppercase;
  font-family: 'Fira Code', 'Consolas', monospace;
  user-select: none;
}

.markdown-body :deep(.hljs-code-block code) {
  background: transparent;
  padding: 0;
  border-radius: 0;
  color: #e6edf3;
  line-height: 1.6;
  font-size: 0.875rem;
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
