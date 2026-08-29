<script setup lang="ts">
import { ref } from 'vue'
import ByteMDEditor from '@/components/ByteMDEditor.vue'

const activeTab = ref('random')

const tabs = [
  { id: 'random', name: '随机数生成器' },
  { id: 'timestamp', name: '时间戳转换' },
  { id: 'color', name: '颜色值转换' },
  { id: 'markdown', name: 'Markdown 预览' },
  { id: 'regex', name: '正则表达式测试' },
]

// ========== 随机数生成器 ==========
const randMin = ref(1)
const randMax = ref(100)
const randCount = ref(1)
const randAllowDup = ref(true)
const randResults = ref<number[]>([])

function generateRandom() {
  const min = Math.ceil(randMin.value)
  const max = Math.floor(randMax.value)
  if (min > max) return
  const range = max - min + 1
  const count = Math.max(1, Math.min(randCount.value, 1000))

  if (!randAllowDup.value && count > range) {
    randResults.value = []
    return
  }

  if (randAllowDup.value) {
    randResults.value = Array.from({ length: count }, () =>
      Math.floor(Math.random() * range) + min
    )
  } else {
    const nums = Array.from({ length: range }, (_, i) => min + i)
    for (let i = nums.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [nums[i], nums[j]] = [nums[j] as number, nums[i] as number]
    }
    randResults.value = nums.slice(0, count)
  }
}

// ========== 时间戳转换 ==========
const tsInput = ref('')
const tsResult = ref('')

function getNow() {
  tsResult.value = Date.now().toString()
}

function unixToDate() {
  const v = tsInput.value.trim()
  if (!v) return
  let ms = Number(v)
  if (v.length <= 10) ms *= 1000
  const d = new Date(ms)
  if (isNaN(d.getTime())) {
    tsResult.value = '无效时间戳'
    return
  }
  tsResult.value = d.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}

function dateToUnix() {
  const v = tsInput.value.trim()
  if (!v) return
  const d = new Date(v)
  if (isNaN(d.getTime())) {
    tsResult.value = '无效日期格式'
    return
  }
  tsResult.value = `${Math.floor(d.getTime() / 1000)} (秒)\n${d.getTime()} (毫秒)`
}

// ========== 颜色值转换 ==========
const hexInput = ref('#7ED7C1')
const rgbResult = ref('')
const hslResult = ref('')
const colorPreview = ref('#7ED7C1')

function hexToRgb(hex: string): [number, number, number] | null {
  const m = hex.replace('#', '')
  if (!/^[0-9a-fA-F]{3}$|^[0-9a-fA-F]{6}$/.test(m)) return null
  const h = m.length === 3
    ? m.split('').map(c => c + c).join('')
    : m
  return [
    parseInt(h.slice(0, 2), 16),
    parseInt(h.slice(2, 4), 16),
    parseInt(h.slice(4, 6), 16),
  ]
}

function rgbToHsl(r: number, g: number, b: number): [number, number, number] {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h = 0, s = 0
  const l = (max + min) / 2
  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r: h = ((g - b) / d + (g < b ? 6 : 0)) / 6; break
      case g: h = ((b - r) / d + 2) / 6; break
      case b: h = ((r - g) / d + 4) / 6; break
    }
  }
  return [Math.round(h * 360), Math.round(s * 100), Math.round(l * 100)]
}

function convertColor() {
  const rgb = hexToRgb(hexInput.value)
  if (!rgb) {
    rgbResult.value = '无效 HEX 值'
    hslResult.value = ''
    colorPreview.value = ''
    return
  }
  const [r, g, b] = rgb
  rgbResult.value = `rgb(${r}, ${g}, ${b})`
  const [h, s, l] = rgbToHsl(r, g, b)
  hslResult.value = `hsl(${h}, ${s}%, ${l}%)`
  colorPreview.value = hexInput.value.startsWith('#') ? hexInput.value : `#${hexInput.value}`
}

convertColor()

// ========== Markdown 预览 ==========
const mdText = ref('# Hello World\n\n这是一个 **Markdown** 在线预览工具。\n\n- 支持列表\n- 支持 `代码`\n- 支持 [链接](https://example.com)\n\n```js\nconsole.log("Hello!")\n```')

// ========== 正则表达式测试 ==========
function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

const regexPattern = ref('')
const regexFlags = ref('g')
const regexTestStr = ref('')
const regexResult = ref<Array<{ text: string; index: number }>>([])
const regexError = ref('')

function testRegex() {
  regexError.value = ''
  regexResult.value = []
  if (!regexPattern.value) return
  try {
    const re = new RegExp(regexPattern.value, regexFlags.value)
    const str = regexTestStr.value
    if (!regexFlags.value.includes('g')) {
      const m = re.exec(str)
      if (m) regexResult.value = [{ text: m[0], index: m.index }]
      return
    }
    let m: RegExpExecArray | null
    while ((m = re.exec(str)) !== null) {
      regexResult.value.push({ text: m[0], index: m.index })
      if (!m[0]) re.lastIndex++
      if (regexResult.value.length > 999) break
    }
  } catch (e: any) {
    regexError.value = e.message || '正则表达式语法错误'
  }
}

function highlightedText() {
  if (regexResult.value.length === 0) return escapeHtml(regexTestStr.value)
  let html = ''
  let last = 0
  const str = regexTestStr.value
  for (const r of regexResult.value) {
    html += escapeHtml(str.slice(last, r.index))
    html += `<mark>${escapeHtml(r.text)}</mark>`
    last = r.index + r.text.length
  }
  html += escapeHtml(str.slice(last))
  return html
}

function copyText(text: string) {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text).catch(() => fallbackCopy(text))
  } else {
    fallbackCopy(text)
  }
}

function fallbackCopy(text: string) {
  const ta = document.createElement('textarea')
  ta.value = text
  ta.style.cssText = 'position:fixed;left:-9999px;top:-9999px;opacity:0'
  document.body.appendChild(ta)
  ta.select()
  document.execCommand('copy')
  document.body.removeChild(ta)
}
</script>

<template>
  <main class="tools-page">
    <div class="container">
      <h1 class="page-title">一些小工具</h1>
      <p class="page-desc"></p>

      <!-- Tab 导航 -->
      <div class="tab-nav">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </button>
      </div>

      <!-- 随机数生成器 -->
      <section v-show="activeTab === 'random'" class="tool-card">
        <h2 class="tool-title">随机数生成器</h2>
        <div class="tool-row">
          <label class="tool-label">
            最小值
            <input v-model.number="randMin" type="number" class="tool-input" />
          </label>
          <label class="tool-label">
            最大值
            <input v-model.number="randMax" type="number" class="tool-input" />
          </label>
          <label class="tool-label">
            数量
            <input v-model.number="randCount" type="number" class="tool-input" min="1" max="1000" />
          </label>
          <label class="tool-label checkbox-label">
            <input v-model="randAllowDup" type="checkbox" />
            允许重复
          </label>
        </div>
        <button class="tool-btn" @click="generateRandom">生成</button>
        <div v-if="randResults.length" class="result-box">
          <div class="result-header">
            <span>结果（{{ randResults.length }} 个）</span>
            <button class="copy-btn" @click="copyText(randResults.join(', '))">复制</button>
          </div>
          <div class="result-content mono">{{ randResults.join(', ') }}</div>
        </div>
      </section>

      <!-- 时间戳转换 -->
      <section v-show="activeTab === 'timestamp'" class="tool-card">
        <h2 class="tool-title">时间戳转换</h2>
        <div class="tool-row">
          <label class="tool-label" style="flex:1">
            输入
            <input v-model="tsInput" type="text" class="tool-input" placeholder="Unix 时间戳（如 1716345600）或日期（如 2024-05-22 12:00:00）" />
          </label>
        </div>
        <div class="btn-group">
          <button class="tool-btn" @click="unixToDate">时间戳 → 日期</button>
          <button class="tool-btn" @click="dateToUnix">日期 → 时间戳</button>
          <button class="tool-btn secondary" @click="getNow">获取当前时间戳</button>
        </div>
        <div v-if="tsResult" class="result-box">
          <div class="result-header">
            <span>结果</span>
            <button class="copy-btn" @click="copyText(tsResult)">复制</button>
          </div>
          <div class="result-content mono">{{ tsResult }}</div>
        </div>
      </section>

      <!-- 颜色值转换 -->
      <section v-show="activeTab === 'color'" class="tool-card">
        <h2 class="tool-title">颜色值转换</h2>
        <div class="tool-row">
          <label class="tool-label">
            HEX
            <div class="color-input-row">
              <input v-model="hexInput" type="text" class="tool-input" placeholder="#7ED7C1" @input="convertColor" />
              <div class="color-swatch" :style="{ background: colorPreview }"></div>
            </div>
          </label>
        </div>
        <button class="tool-btn" @click="convertColor">转换</button>
        <div v-if="rgbResult" class="result-box">
          <div class="result-header">
            <span>RGB</span>
            <button class="copy-btn" @click="copyText(rgbResult)">复制</button>
          </div>
          <div class="result-content mono">{{ rgbResult }}</div>
        </div>
        <div v-if="hslResult" class="result-box" style="margin-top: 8px;">
          <div class="result-header">
            <span>HSL</span>
            <button class="copy-btn" @click="copyText(hslResult)">复制</button>
          </div>
          <div class="result-content mono">{{ hslResult }}</div>
        </div>
      </section>

      <!-- Markdown 预览 -->
      <section v-show="activeTab === 'markdown'" class="tool-card">
        <h2 class="tool-title">Markdown 在线预览</h2>
        <ByteMDEditor
          v-model="mdText"
          editor-height="450px"
          placeholder="输入 Markdown 内容进行预览..."
        />
      </section>

      <!-- 正则表达式测试 -->
      <section v-show="activeTab === 'regex'" class="tool-card">
        <h2 class="tool-title">正则表达式在线测试</h2>
        <details class="help-box">
          <summary class="help-summary">使用说明</summary>
          <div class="help-content">
            <p><strong>基本语法：</strong>在"正则表达式"框中输入模式，"测试字符串"框中输入要匹配的文本，结果会实时显示。</p>
            <p><strong>常用标志（flags）：</strong></p>
            <ul>
              <li><code>g</code> — 全局匹配，查找所有匹配项（默认已开启）</li>
              <li><code>i</code> — 忽略大小写</li>
              <li><code>m</code> — 多行模式，<code>^</code> 和 <code>$</code> 匹配每行的开头和结尾</li>
              <li><code>s</code> — 让 <code>.</code> 也能匹配换行符</li>
            </ul>
            <p><strong>常用语法速查：</strong></p>
            <div class="help-table-wrap">
              <table class="help-table">
                <tr><th>语法</th><th>说明</th><th>示例</th></tr>
                <tr><td><code>.</code></td><td>任意单个字符</td><td><code>a.c</code> 匹配 "abc"</td></tr>
                <tr><td><code>\d</code></td><td>数字 [0-9]</td><td><code>\d+</code> 匹配 "123"</td></tr>
                <tr><td><code>\w</code></td><td>字母数字下划线</td><td><code>\w+</code> 匹配 "hello_1"</td></tr>
                <tr><td><code>\s</code></td><td>空白字符</td><td><code>a\sb</code> 匹配 "a b"</td></tr>
                <tr><td><code>*</code></td><td>0 次或多次</td><td><code>ab*c</code> 匹配 "ac", "abc"</td></tr>
                <tr><td><code>+</code></td><td>1 次或多次</td><td><code>ab+c</code> 匹配 "abc", "abbc"</td></tr>
                <tr><td><code>?</code></td><td>0 次或 1 次</td><td><code>colou?r</code> 匹配 "color", "colour"</td></tr>
                <tr><td><code>{n,m}</code></td><td>n 到 m 次</td><td><code>\d{2,4}</code> 匹配 "12", "1234"</td></tr>
                <tr><td><code>[abc]</code></td><td>字符集合</td><td><code>[aeiou]</code> 匹配元音</td></tr>
                <tr><td><code>[^abc]</code></td><td>非集合中的字符</td><td><code>[^0-9]</code> 匹配非数字</td></tr>
                <tr><td><code>(...)</code></td><td>捕获分组</td><td><code>(\d+)-(\d+)</code></td></tr>
                <tr><td><code>^</code> / <code>$</code></td><td>行首 / 行尾</td><td><code>^Hello</code> 匹配行首的 "Hello"</td></tr>
                <tr><td><code>\b</code></td><td>单词边界</td><td><code>\bcat\b</code> 匹配独立的 "cat"</td></tr>
                <tr><td><code>a|b</code></td><td>或</td><td><code>cat|dog</code> 匹配 "cat" 或 "dog"</td></tr>
              </table>
            </div>
          </div>
        </details>
        <div class="tool-row">
          <label class="tool-label" style="flex:2">
            正则表达式
            <div class="regex-input-row">
              <span class="regex-slash">/</span>
              <input v-model="regexPattern" type="text" class="tool-input regex-input" placeholder="pattern" @input="testRegex" />
              <span class="regex-slash">/</span>
              <input v-model="regexFlags" type="text" class="tool-input regex-flags" placeholder="g" @input="testRegex" />
            </div>
          </label>
        </div>
        <label class="tool-label" style="margin-top: 12px;">
          测试字符串
          <textarea v-model="regexTestStr" class="tool-textarea" placeholder="输入要测试的文本" @input="testRegex" v-tab-indent></textarea>
        </label>
        <div v-if="regexError" class="error-msg">{{ regexError }}</div>
        <div v-if="regexResult.length" class="result-box">
          <div class="result-header">
            <span>匹配结果（{{ regexResult.length }} 个）</span>
          </div>
          <div class="result-content">
            <div v-for="(m, i) in regexResult" :key="i" class="regex-match">
              <span class="match-index">[{{ i }}]</span>
              <span class="match-text">"{{ m.text }}"</span>
              <span class="match-pos">index: {{ m.index }}</span>
            </div>
          </div>
        </div>
        <div v-if="regexTestStr && !regexError" class="result-box highlight-box" style="margin-top: 8px;">
          <div class="result-header"><span>高亮预览</span></div>
          <div class="result-content highlight-text" v-html="highlightedText()"></div>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.tools-page {
  min-height: 100vh;
  padding: 40px 24px;
  background: var(--bg-primary);
}

.container {
  max-width: 900px;
  margin: 0 auto;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-desc {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  margin-bottom: 24px;
}

/* Tab 导航 */
.tab-nav {
  display: flex;
  gap: 4px;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 6px;
  box-shadow: var(--shadow-card);
  margin-bottom: 24px;
  overflow-x: auto;
}

.tab-btn {
  flex: 1;
  min-width: max-content;
  padding: 10px 16px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tab-btn:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.tab-btn.active {
  color: white;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  box-shadow: var(--shadow-soft);
}

/* 工具卡片 */
.tool-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-card);
  margin-bottom: 24px;
}

.tool-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}

.tool-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  align-items: flex-end;
  margin-bottom: 16px;
}

.tool-label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.checkbox-label {
  flex-direction: row;
  align-items: center;
  gap: 8px;
  padding-bottom: 8px;
}

.tool-input {
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  outline: none;
  transition: border-color 0.2s;
  font-family: 'Fira Code', 'Consolas', monospace;
}

.tool-input:focus {
  border-color: var(--primary);
}

.tool-textarea {
  width: 100%;
  min-height: 100px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  outline: none;
  resize: vertical;
  font-family: 'Fira Code', 'Consolas', monospace;
  line-height: 1.6;
}

.tool-textarea:focus {
  border-color: var(--primary);
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: var(--shadow-soft);
}

.tool-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-hover);
}

.tool-btn.secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  box-shadow: none;
  border: 1px solid var(--border);
}

.tool-btn.secondary:hover {
  border-color: var(--primary);
  color: var(--primary);
  transform: none;
  box-shadow: none;
}

.btn-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

/* 结果区域 */
.result-box {
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  overflow: hidden;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 14px;
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid var(--border);
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.copy-btn {
  padding: 2px 10px;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  color: var(--primary);
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.result-content {
  padding: 12px 14px;
  font-size: 0.875rem;
  color: var(--text-primary);
  line-height: 1.6;
  word-break: break-all;
}

.result-content.mono {
  font-family: 'Fira Code', 'Consolas', monospace;
}

/* 颜色输入 */
.color-input-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.color-swatch {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  border: 2px solid var(--border);
  flex-shrink: 0;
}

/* Markdown 编辑器圆角适配 */
.tool-card :deep(.bytemd-editor-wrapper) {
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* 使用说明 */
.help-box {
  margin-bottom: 20px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.help-summary {
  padding: 10px 16px;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--primary);
  background: var(--bg-secondary);
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.help-summary:hover {
  background: var(--primary-lighter);
}

.help-content {
  padding: 16px;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  line-height: 1.7;
}

.help-content p {
  margin: 0 0 8px;
}

.help-content ul {
  padding-left: 1.25rem;
  margin: 0 0 12px;
}

.help-content li {
  margin: 2px 0;
}

.help-content code {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 0.8rem;
  background: var(--bg-secondary);
  padding: 1px 5px;
  border-radius: 3px;
  color: var(--primary-dark);
}

.help-table-wrap {
  overflow-x: auto;
  margin-top: 8px;
}

.help-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
}

.help-table th,
.help-table td {
  padding: 6px 10px;
  border: 1px solid var(--border);
  text-align: left;
}

.help-table th {
  background: var(--bg-secondary);
  font-weight: 600;
  color: var(--text-primary);
}

.help-table td code {
  font-size: 0.78rem;
}

/* 正则输入 */
.regex-input-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.regex-slash {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-muted);
  font-family: 'Fira Code', 'Consolas', monospace;
}

.regex-input {
  flex: 1;
}

.regex-flags {
  width: 50px;
}

.regex-match {
  display: flex;
  gap: 12px;
  padding: 4px 0;
  font-size: 0.8125rem;
  font-family: 'Fira Code', 'Consolas', monospace;
}

.match-index {
  color: var(--text-muted);
  min-width: 30px;
}

.match-text {
  color: var(--primary-dark);
  font-weight: 500;
}

.match-pos {
  color: var(--text-secondary);
}

.error-msg {
  padding: 10px 14px;
  background: rgba(239, 68, 68, 0.08);
  color: var(--color-danger);
  border-radius: var(--radius-sm);
  font-size: 0.8125rem;
  margin-top: 12px;
}

.highlight-box .highlight-text {
  font-family: 'Fira Code', 'Consolas', monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

.highlight-box .highlight-text :deep(mark) {
  background: rgba(126, 215, 193, 0.3);
  color: var(--text-primary);
  padding: 1px 2px;
  border-radius: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .tools-page {
    padding: 24px 16px;
  }

  .tool-card {
    padding: 20px;
  }

  .tool-row {
    flex-direction: column;
  }

  .tab-nav {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }
}
</style>
