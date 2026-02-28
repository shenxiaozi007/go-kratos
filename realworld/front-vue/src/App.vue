<template>
  <div :class="['app-root', `theme-${theme}`]">
    <div class="app-background"></div>

    <main class="app-shell">
      <header class="shell-header">
        <div class="brand">
          <div class="brand-mark">
            <span class="brand-badge"></span>
          </div>
          <div class="brand-text">
            <h1 class="brand-title">Focus Todos</h1>
            <p class="brand-subtitle">以 Apple 风格的优雅界面，整理你今天的待办。</p>
          </div>
        </div>

        <button
          class="icon-button theme-toggle"
          type="button"
          :aria-label="theme === 'light' ? '切换到夜间模式' : '切换到亮色模式'"
          @click="toggleTheme"
        >
          <svg v-if="theme === 'light'" viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="12" cy="12" r="4" />
            <g class="sun-rays">
              <line x1="12" y1="2" x2="12" y2="5" />
              <line x1="12" y1="19" x2="12" y2="22" />
              <line x1="4.22" y1="4.22" x2="6.34" y2="6.34" />
              <line x1="17.66" y1="17.66" x2="19.78" y2="19.78" />
              <line x1="2" y1="12" x2="5" y2="12" />
              <line x1="19" y1="12" x2="22" y2="12" />
              <line x1="4.22" y1="19.78" x2="6.34" y2="17.66" />
              <line x1="17.66" y1="6.34" x2="19.78" y2="4.22" />
            </g>
          </svg>
          <svg v-else viewBox="0 0 24 24" aria-hidden="true">
            <path
              d="M21 13.25A7.75 7.75 0 0 1 10.75 3 7.75 7.75 0 1 0 21 13.25z"
            />
          </svg>
        </button>
      </header>

      <section class="todo-card" aria-label="Todo 列表">
        <section class="todo-input-row">
          <div class="input-with-icon">
            <span class="input-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M11 5h2v6h6v2h-6v6h-2v-6H5v-2h6z" />
              </svg>
            </span>
            <input
              v-model.trim="newTodo"
              class="todo-input"
              type="text"
              name="todo"
              autocomplete="off"
              placeholder="写下你今天最重要的一件事…"
              @keyup.enter="addTodo"
            />
          </div>
          <button
            class="primary-button"
            type="button"
            :disabled="!newTodo || loading"
            @click="addTodo"
          >
            <span class="primary-button-icon" aria-hidden="true">
              <svg v-if="!loading" viewBox="0 0 24 24">
                <path d="M11 5h2v14h-2z" />
                <path d="M5 11h14v2H5z" />
              </svg>
              <svg v-else class="spinner" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none" stroke-dasharray="31.416" stroke-dashoffset="31.416">
                  <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416;0 31.416" repeatCount="indefinite"/>
                  <animate attributeName="stroke-dashoffset" dur="2s" values="0;-15.708;-31.416;-31.416" repeatCount="indefinite"/>
                </circle>
              </svg>
            </span>
            {{ loading ? '处理中...' : '添加' }}
          </button>
        </section>

        <div v-if="error" class="error-message" role="alert">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
          </svg>
          <span>{{ error }}</span>
          <button class="error-close" type="button" @click="error = null" aria-label="关闭错误提示">
            <svg viewBox="0 0 24 24">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </div>

        <section class="todo-meta-row">
          <div class="stats-block" aria-label="完成情况统计">
            <div class="progress-ring">
              <svg viewBox="0 0 60 60">
                <circle class="ring-background" cx="30" cy="30" r="24" />
                <circle
                  class="ring-progress"
                  cx="30"
                  cy="30"
                  r="24"
                  :stroke-dasharray="circumference"
                  :stroke-dashoffset="progressOffset"
                />
              </svg>
              <div class="progress-label">
                <span class="progress-value">{{ completionRate }}%</span>
                <span class="progress-caption">完成度</span>
              </div>
            </div>
            <div class="stats-text">
              <p class="stats-main">
                已完成
                <strong>{{ completedCount }}</strong>
                /
                <span>{{ totalCount }}</span>
              </p>
              <p class="stats-sub">
                还有
                <strong>{{ activeCount }}</strong>
                个待办在进行中
              </p>
            </div>
          </div>

          <div class="filters" role="tablist" aria-label="筛选待办">
            <button
              v-for="option in filterOptions"
              :key="option.value"
              class="filter-pill"
              type="button"
              role="tab"
              :aria-selected="filter === option.value"
              :class="{ 'filter-pill-active': filter === option.value }"
              @click="setFilter(option.value)"
            >
              <span class="filter-dot" :data-type="option.value"></span>
              <span>{{ option.label }}</span>
            </button>
          </div>
        </section>

        <section class="todo-list-wrapper">
          <div v-if="loading && !todos.length" class="loading-state">
            <svg class="loading-spinner" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none" stroke-dasharray="31.416" stroke-dashoffset="31.416">
                <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416;0 31.416" repeatCount="indefinite"/>
                <animate attributeName="stroke-dashoffset" dur="2s" values="0;-15.708;-31.416;-31.416" repeatCount="indefinite"/>
              </circle>
            </svg>
            <p>加载中...</p>
          </div>

          <transition-group
            v-else
            name="todo-list"
            tag="ul"
            class="todo-list"
          >
            <li
              v-for="todo in filteredTodos"
              :key="todo.id"
              class="todo-item"
              :class="{ 'todo-item-completed': todo.completed }"
            >
              <button
                class="icon-button check-button"
                type="button"
                :aria-label="todo.completed ? '标记为未完成' : '标记为已完成'"
                :disabled="loading"
                @click="toggleTodo(todo.id)"
              >
                <span class="checkbox">
                  <span class="checkbox-inner">
                    <svg viewBox="0 0 20 20" aria-hidden="true">
                      <polyline points="4 10.5 8.5 15 16 5" />
                    </svg>
                  </span>
                </span>
              </button>

              <div class="todo-content">
                <p class="todo-text">
                  {{ todo.text }}
                </p>
                <p class="todo-meta">
                  创建于 {{ formatTime(todo.createdAt) }}
                </p>
              </div>

              <button
                class="icon-button delete-button"
                type="button"
                aria-label="删除待办"
                :disabled="loading"
                @click="removeTodo(todo.id)"
              >
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path
                    d="M9 3h6l1 2h4v2H4V5h4l1-2zm-1 6h2v8H8V9zm6 0h2v8h-2V9z"
                  />
                </svg>
              </button>
            </li>
          </transition-group>

          <p v-if="!todos.length" class="empty-state">
            还没有待办。先写下你今天最重要的一件事吧。
          </p>
          <p v-else-if="!filteredTodos.length" class="empty-state">
            当前筛选下没有结果，尝试切换到其它筛选。
          </p>
        </section>
      </section>

      <footer class="shell-footer">
        <span class="footer-text">
          轻盈、专注的 Todo 体验 · 基于 Vue 3
        </span>
      </footer>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { addTodo as apiAddTodo, deleteTodo, getTodoList, updateTodo } from './api/todo'

const THEME_KEY = 'focus-todos-theme'

const theme = ref('light')
const newTodo = ref('')
const todos = ref([])
const filter = ref('all')
const loading = ref(false)
const error = ref(null)

const filterOptions = [
  { value: 'all', label: '全部' },
  { value: 'active', label: '未完成' },
  { value: 'completed', label: '已完成' }
]

const totalCount = computed(() => todos.value.length)
const completedCount = computed(
  () => todos.value.filter((t) => t.completed).length
)
const activeCount = computed(() => totalCount.value - completedCount.value)

const completionRate = computed(() => {
  if (!totalCount.value) return 0
  return Math.round((completedCount.value / totalCount.value) * 100)
})










const filteredTodos = computed(() => {
  if (filter.value === 'active') {
    return todos.value.filter((t) => !t.completed)
  }
  if (filter.value === 'completed') {
    return todos.value.filter((t) => t.completed)
  }
  return todos.value
})

const circumference = 2 * Math.PI * 24
const progressOffset = computed(
  () => circumference - (circumference * completionRate.value) / 100
)

async function addTodo() {
  if (!newTodo.value) return
  const text = newTodo.value.trim()
  if (!text) return

  loading.value = true
  error.value = null

  try {
    const response = await apiAddTodo(text, false)
    // 假设后端返回的数据结构包含 id, value, isCompleted 等字段
    const newItem = {
      id: response.id || response.todo?.id || String(Date.now()),
      text: response.value || response.todo?.value || text,
      completed: response.isCompleted || response.todo?.isCompleted || false,
      createdAt: response.createdAt || response.todo?.createdAt || new Date().toISOString()
    }
    todos.value.unshift(newItem)
    newTodo.value = ''
  } catch (err) {
    error.value = '添加待办失败，请稍后重试'
    console.error('Failed to add todo:', err)
  } finally {
    loading.value = false
  }
}

async function toggleTodo(id) {
  loading.value = true
  error.value = null

  try {
    const response = await updateTodo(id)
    // 更新本地状态
    const updatedTodo = response.todo || response
    todos.value = todos.value.map((todo) =>
      todo.id === id
        ? {
            ...todo,
            completed: updatedTodo.isCompleted !== undefined
              ? updatedTodo.isCompleted
              : !todo.completed
          }
        : todo
    )
  } catch (err) {
    error.value = '更新待办状态失败，请稍后重试'
    console.error('Failed to update todo:', err)
    // 失败时回滚状态
    todos.value = todos.value.map((todo) =>
      todo.id === id ? { ...todo, completed: !todo.completed } : todo
    )
  } finally {
    loading.value = false
  }
}

async function removeTodo(id) {
  loading.value = true
  error.value = null

  // 乐观更新：先移除 UI 中的项
  const removedTodo = todos.value.find((t) => t.id === id)
  todos.value = todos.value.filter((todo) => todo.id !== id)

  try {
    await deleteTodo(id)
  } catch (err) {
    error.value = '删除待办失败，请稍后重试'
    console.error('Failed to delete todo:', err)
    // 失败时恢复项
    if (removedTodo) {
      todos.value.push(removedTodo)
    }
  } finally {
    loading.value = false
  }
}

function setFilter(value) {
  filter.value = value
}

function formatTime(iso) {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''
  const fromNowMinutes = Math.round(
    (Date.now() - date.getTime()) / 1000 / 60
  )

  if (fromNowMinutes < 1) return '刚刚'
  if (fromNowMinutes < 60) return `${fromNowMinutes} 分钟前`

  const hours = Math.floor(fromNowMinutes / 60)
  if (hours < 24) return `${hours} 小时前`

  const y = date.getFullYear()
  const m = `${date.getMonth() + 1}`.padStart(2, '0')
  const d = `${date.getDate()}`.padStart(2, '0')
  const hh = `${date.getHours()}`.padStart(2, '0')
  const mm = `${date.getMinutes()}`.padStart(2, '0')
  return `${y}/${m}/${d} ${hh}:${mm}`
}

function toggleTheme() {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}

onMounted(async () => {
  // 加载主题设置
  try {
    const savedTheme = window.localStorage.getItem(THEME_KEY)
    if (savedTheme === 'light' || savedTheme === 'dark') {
      theme.value = savedTheme
    } else if (window.matchMedia) {
      const prefersDark = window.matchMedia(
        '(prefers-color-scheme: dark)'
      ).matches
      theme.value = prefersDark ? 'dark' : 'light'
    }
  } catch (e) {
    // ignore
  }

  // 从 API 加载待办列表
  loading.value = true
  error.value = null

  try {
    const response = await getTodoList()
    // 假设后端返回的数据结构包含 todos 数组或直接是数组
    const todoList = response.todos || response.list || response || []
    todos.value = todoList.map((item) => ({
      id: item.id || String(Date.now() + Math.random()),
      text: item.value || item.text || '',
      completed: Boolean(item.isCompleted !== undefined ? item.isCompleted : item.completed),
      createdAt: item.createdAt || item.created_at || new Date().toISOString()
    }))
  } catch (err) {
    error.value = '加载待办列表失败，请检查网络连接'
    console.error('Failed to load todos:', err)
  } finally {
    loading.value = false
  }
})

watch(theme, (value) => {
  try {
    window.localStorage.setItem(THEME_KEY, value)
  } catch (e) {
    // ignore
  }
})
</script>

<style>
*,
*::before,
*::after {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, system-ui, -system-ui,
    'SF Pro Text', 'SF Pro Icons', 'Helvetica Neue', Helvetica, Arial,
    sans-serif;
  -webkit-font-smoothing: antialiased;
  background-color: #050509;
}

#app {
  height: 100%;
}

.app-root {
  position: relative;
  min-height: 100%;
  color: var(--color-fg-primary);
  background: radial-gradient(circle at top left, #2f5bff 0, #050509 55%);
  background-attachment: fixed;
  overflow: hidden;
}

.app-background {
  pointer-events: none;
  position: fixed;
  inset: 0;
  background:
    radial-gradient(circle at 0% 100%, rgba(255, 105, 180, 0.18), transparent 60%),
    radial-gradient(circle at 100% 0%, rgba(72, 209, 204, 0.18), transparent 60%);
  opacity: 0.9;
}

.app-shell {
  position: relative;
  z-index: 1;
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 20px 24px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.shell-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: conic-gradient(
    from 160deg,
    #ff8a5c,
    #ffd15c,
    #6ce6b7,
    #59a3ff,
    #b28cff,
    #ff8a5c
  );
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px;
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.18),
    0 18px 45px rgba(0, 0, 0, 0.6);
}

.brand-badge {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: inherit;
  background: radial-gradient(
    circle at 20% 0%,
    rgba(255, 255, 255, 0.9),
    transparent 45%
  );
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.18),
    inset 0 10px 20px rgba(255, 255, 255, 0.18),
    inset 0 -14px 22px rgba(0, 0, 0, 0.3);
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-title {
  margin: 0;
  font-size: 24px;
  letter-spacing: -0.02em;
}

.brand-subtitle {
  margin: 0;
  font-size: 13px;
  opacity: 0.78;
}

.icon-button {
  border: none;
  background: transparent;
  padding: 0;
  cursor: pointer;
  color: inherit;
  transition: opacity 0.15s ease;
}

.icon-button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.icon-button svg {
  display: block;
}

.theme-toggle {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.09),
    rgba(255, 255, 255, 0.02)
  );
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.16),
    0 10px 30px rgba(0, 0, 0, 0.55);
  transition:
    transform 0.18s ease,
    box-shadow 0.18s ease,
    background 0.18s ease;
}

.theme-toggle svg {
  width: 20px;
  height: 20px;
  stroke: currentColor;
  fill: none;
  stroke-width: 1.5;
}

.theme-toggle .sun-rays line {
  stroke-linecap: round;
}

.theme-toggle:hover {
  transform: translateY(-1px);
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.22),
    0 16px 40px rgba(0, 0, 0, 0.7);
}

.theme-toggle:active {
  transform: translateY(0);
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.2),
    0 6px 20px rgba(0, 0, 0, 0.6);
}

.todo-card {
  position: relative;
  border-radius: 28px;
  padding: 20px 20px 18px;
  background: radial-gradient(
      circle at 0 0,
      rgba(255, 255, 255, 0.18),
      transparent 55%
    ),
    linear-gradient(
      145deg,
      rgba(15, 15, 25, 0.9),
      rgba(15, 15, 28, 0.96)
    );
  box-shadow:
    0 18px 45px rgba(0, 0, 0, 0.65),
    0 0 0 1px rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(30px);
  -webkit-backdrop-filter: blur(30px);
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.todo-input-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.input-with-icon {
  flex: 1 1 220px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 999px;
  background: rgba(7, 7, 16, 0.94);
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.09),
    0 14px 30px rgba(0, 0, 0, 0.7);
}

.input-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent-primary);
}

.input-icon svg {
  width: 18px;
  height: 18px;
  stroke: currentColor;
  stroke-width: 1.5;
  fill: none;
}

.todo-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: var(--color-fg-primary);
  font-size: 15px;
}

.todo-input::placeholder {
  color: var(--color-fg-muted);
}

.primary-button {
  position: relative;
  flex-shrink: 0;
  border-radius: 999px;
  border: none;
  padding: 10px 18px 10px 14px;
  font-size: 14px;
  letter-spacing: 0.02em;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #050509;
  cursor: pointer;
  background: linear-gradient(135deg, #ffcf71, #ff9a62);
  box-shadow:
    0 12px 30px rgba(255, 149, 0, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.6);
  transition:
    transform 0.16s ease,
    box-shadow 0.16s ease,
    filter 0.16s ease;
}

.primary-button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.primary-button-icon svg {
  width: 16px;
  height: 16px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.6;
}

.primary-button-icon .spinner {
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.primary-button:hover:not(:disabled) {
  transform: translateY(-1px);
  filter: brightness(1.03);
  box-shadow:
    0 18px 40px rgba(255, 149, 0, 0.7),
    0 0 0 1px rgba(255, 255, 255, 0.7);
}

.primary-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow:
    0 8px 24px rgba(255, 149, 0, 0.55),
    0 0 0 1px rgba(255, 255, 255, 0.7);
}

.primary-button:disabled {
  opacity: 0.55;
  cursor: default;
  box-shadow:
    0 8px 18px rgba(0, 0, 0, 0.45),
    0 0 0 1px rgba(255, 255, 255, 0.35);
}

.todo-meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: center;
  justify-content: space-between;
}

.stats-block {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.progress-ring {
  position: relative;
  width: 60px;
  height: 60px;
}

.progress-ring svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.ring-background {
  fill: none;
  stroke: rgba(255, 255, 255, 0.12);
  stroke-width: 4;
}

.ring-progress {
  fill: none;
  stroke: var(--color-accent-primary);
  stroke-width: 4;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.4s ease;
}

.progress-label {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1px;
}

.progress-value {
  font-size: 14px;
  font-weight: 600;
}

.progress-caption {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  opacity: 0.75;
}

.stats-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
}

.stats-main {
  margin: 0;
  font-weight: 500;
}

.stats-main strong {
  color: var(--color-accent-primary);
}

.stats-sub {
  margin: 0;
  opacity: 0.8;
}

.filters {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px;
  border-radius: 999px;
  background: rgba(7, 7, 16, 0.9);
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.08),
    0 10px 24px rgba(0, 0, 0, 0.7);
}

.filter-pill {
  border-radius: 999px;
  border: none;
  padding: 5px 10px;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--color-fg-muted);
  background: transparent;
  cursor: pointer;
  transition:
    color 0.15s ease,
    background 0.15s ease,
    box-shadow 0.15s ease,
    transform 0.15s ease;
}

.filter-pill .filter-dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.38);
}

.filter-pill .filter-dot[data-type='active'] {
  background: #ff9a62;
}

.filter-pill .filter-dot[data-type='completed'] {
  background: #6ce6b7;
}

.filter-pill-active {
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-fg-primary);
  box-shadow:
    0 7px 16px rgba(0, 0, 0, 0.7),
    0 0 0 1px rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.todo-list-wrapper {
  margin-top: 4px;
  max-height: min(420px, 60vh);
  overflow: auto;
  padding-right: 4px;
}

.todo-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.todo-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 16px;
  background: rgba(5, 5, 14, 0.94);
  box-shadow:
    0 10px 26px rgba(0, 0, 0, 0.75),
    0 0 0 1px rgba(255, 255, 255, 0.05);
  transition:
    transform 0.13s ease,
    box-shadow 0.13s ease,
    background 0.13s ease,
    opacity 0.13s ease;
}

.todo-item:hover {
  transform: translateY(-1px);
  box-shadow:
    0 16px 36px rgba(0, 0, 0, 0.8),
    0 0 0 1px rgba(255, 255, 255, 0.09);
}

.todo-item-completed {
  opacity: 0.7;
  background: radial-gradient(
      circle at 0 0,
      rgba(108, 230, 183, 0.16),
      transparent 55%
    ),
    rgba(5, 12, 12, 0.96);
}

.check-button {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox {
  width: 22px;
  height: 22px;
  border-radius: 999px;
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.24),
    inset 0 0 0 1px rgba(0, 0, 0, 0.6),
    0 8px 20px rgba(0, 0, 0, 0.6);
  background: radial-gradient(
    circle at 30% 0,
    rgba(255, 255, 255, 0.7),
    transparent 55%
  );
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox-inner {
  width: 100%;
  height: 100%;
  border-radius: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at 20% 0, #6ce6b7, #1b6c5b);
  opacity: 0;
  transform: scale(0.4);
  transition:
    opacity 0.16s ease,
    transform 0.16s ease;
}

.checkbox-inner svg {
  width: 14px;
  height: 14px;
  stroke: #04100b;
  stroke-width: 2;
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.todo-item-completed .checkbox-inner {
  opacity: 1;
  transform: scale(1);
}

.todo-content {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.todo-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.4;
  word-break: break-word;
}

.todo-item-completed .todo-text {
  text-decoration: line-through;
  text-decoration-thickness: 1px;
  text-decoration-color: rgba(255, 255, 255, 0.5);
  opacity: 0.88;
}

.todo-meta {
  margin: 0;
  font-size: 11px;
  opacity: 0.7;
}

.delete-button {
  width: 32px;
  height: 32px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.02);
  color: var(--color-fg-muted);
  transition:
    background 0.15s ease,
    color 0.15s ease,
    transform 0.12s ease;
}

.delete-button svg {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.delete-button:hover {
  background: rgba(255, 80, 80, 0.16);
  color: #ff7a7a;
  transform: translateY(-1px);
}

.delete-button:active {
  transform: translateY(0);
}

.empty-state {
  margin: 10px 0 0;
  font-size: 13px;
  color: var(--color-fg-muted);
}

.error-message {
  margin-top: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  background: rgba(255, 59, 48, 0.12);
  border: 1px solid rgba(255, 59, 48, 0.25);
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #ff3b30;
  animation: slideDown 0.2s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-message svg:first-child {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  fill: currentColor;
}

.error-message span {
  flex: 1;
}

.error-close {
  width: 24px;
  height: 24px;
  padding: 0;
  border: none;
  background: transparent;
  color: currentColor;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: background-color 0.15s ease;
  flex-shrink: 0;
}

.error-close:hover {
  background: rgba(255, 59, 48, 0.15);
}

.error-close svg {
  width: 18px;
  height: 18px;
  fill: currentColor;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 16px;
  color: var(--color-fg-muted);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  color: var(--color-accent-primary);
}

.loading-state p {
  margin: 0;
  font-size: 14px;
  letter-spacing: 0.02em;
}

.shell-footer {
  margin-top: auto;
  padding-top: 18px;
  display: flex;
  justify-content: center;
}

.footer-text {
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  opacity: 0.72;
}

.todo-list-enter-active,
.todo-list-leave-active,
.todo-list-move {
  transition:
    opacity 0.18s ease,
    transform 0.18s ease;
}

.todo-list-enter-from,
.todo-list-leave-to {
  opacity: 0;
  transform: translateY(4px) scale(0.99);
}

.todo-list-leave-active {
  position: absolute;
}

.app-root.theme-light {
  --color-fg-primary: #0a0c1a;
  --color-fg-muted: #5d6475;
  --color-accent-primary: #336dff;
  background: radial-gradient(circle at top left, #e4eeff 0, #f7f7fb 55%);
}

.app-root.theme-light .app-background {
  background:
    radial-gradient(circle at 0% 100%, rgba(255, 194, 222, 0.4), transparent 60%),
    radial-gradient(circle at 100% 0%, rgba(187, 213, 255, 0.6), transparent 60%);
  opacity: 0.85;
}

.app-root.theme-light .todo-card {
  background: radial-gradient(
      circle at 0 0,
      rgba(255, 255, 255, 0.9),
      rgba(255, 255, 255, 0.5) 50%,
      rgba(255, 255, 255, 0.9)
    );
  box-shadow:
    0 20px 40px rgba(111, 134, 194, 0.36),
    0 0 0 1px rgba(255, 255, 255, 0.9);
}

.app-root.theme-light .input-with-icon {
  background: rgba(246, 247, 255, 0.9);
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.9),
    0 14px 30px rgba(155, 169, 214, 0.5);
}

.app-root.theme-light .filters {
  background: rgba(246, 247, 255, 0.9);
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.9),
    0 10px 24px rgba(155, 169, 214, 0.45);
}

.app-root.theme-light .filter-pill-active {
  background: rgba(255, 255, 255, 0.9);
  box-shadow:
    0 7px 20px rgba(155, 169, 214, 0.7),
    0 0 0 1px rgba(51, 109, 255, 0.52);
}

.app-root.theme-light .todo-item {
  background: rgba(250, 250, 255, 0.95);
  box-shadow:
    0 10px 26px rgba(155, 169, 214, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.9);
}

.app-root.theme-light .todo-item-completed {
  background: radial-gradient(
      circle at 0 0,
      rgba(108, 230, 183, 0.12),
      transparent 55%
    ),
    rgba(245, 252, 248, 0.96);
}

.app-root.theme-light .theme-toggle {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.95),
    rgba(244, 248, 255, 0.92)
  );
  box-shadow:
    0 12px 30px rgba(155, 169, 214, 0.7),
    0 0 0 1px rgba(255, 255, 255, 1);
}

.app-root.theme-light body {
  background-color: #f5f5fb;
}

.app-root.theme-dark {
  --color-fg-primary: #f6f7ff;
  --color-fg-muted: #9399b4;
  --color-accent-primary: #6ce6b7;
}

.app-root.theme-dark .todo-card {
  color: var(--color-fg-primary);
}

.app-root.theme-dark .primary-button {
  color: #050509;
}

.app-root.theme-dark .todo-input::placeholder {
  color: var(--color-fg-muted);
}

@media (max-width: 720px) {
  .app-shell {
    padding: 18px 14px 18px;
  }

  .shell-header {
    align-items: flex-start;
  }

  .brand-title {
    font-size: 21px;
  }

  .brand-subtitle {
    font-size: 12px;
  }

  .todo-card {
    border-radius: 20px;
    padding: 16px 14px 14px;
  }

  .todo-meta-row {
    align-items: flex-start;
  }

  .todo-list-wrapper {
    max-height: none;
    max-height: 60vh;
  }
}

@media (max-width: 480px) {
  .todo-input-row {
    align-items: stretch;
  }

  .primary-button {
    width: 100%;
    justify-content: center;
  }

  .filters {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
