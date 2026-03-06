<template>
  <div class="page page-history">
    <h1 class="page-title">历史战绩</h1>
    <p class="page-subtitle">最近对局记录，分页展示。</p>

    <div v-if="loading" class="loading" role="status" aria-live="polite" aria-busy="true">加载中…</div>
    <p v-else-if="error" class="error-msg" role="alert">{{ error }}</p>

    <template v-else>
      <div class="panel table-wrap">
        <table class="table">
          <thead>
            <tr>
              <th scope="col">对局 ID</th>
              <th scope="col">模式</th>
              <th scope="col">地图</th>
              <th scope="col">胜负</th>
              <th scope="col">P1击杀</th>
              <th scope="col">P2击杀</th>
              <th scope="col">敌方击杀</th>
              <th scope="col">时长(秒)</th>
              <th scope="col">时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in records" :key="r.session_id">
              <td class="session-id">{{ r.session_id }}</td>
              <td>{{ modeLabel(r.mode) }}</td>
              <td>{{ r.map_id }}</td>
              <td :class="r.is_win ? 'win' : 'lose'">{{ r.is_win ? '胜' : '负' }}</td>
              <td>{{ r.player1_kills }}</td>
              <td>{{ r.player2_kills }}</td>
              <td>{{ r.enemy_kills }}</td>
              <td>{{ r.duration_seconds }}</td>
              <td>{{ formatTime(r.created_at) }}</td>
            </tr>
          </tbody>
        </table>
        <p v-if="records.length === 0" class="empty">暂无战绩</p>
      </div>
      <nav class="pagination" aria-label="历史战绩分页">
        <button type="button" class="btn" :disabled="page <= 1" @click="loadPage(page - 1)" :aria-label="`上一页，当前第 ${page} 页`">上一页</button>
        <span class="page-info" aria-live="polite">第 {{ page }} 页，共 {{ total }} 条</span>
        <button type="button" class="btn" :disabled="page * pageSize >= total" @click="loadPage(page + 1)" :aria-label="`下一页，当前第 ${page} 页`">下一页</button>
      </nav>
    </template>

    <p class="nav-link">
      <router-link to="/">返回首页</router-link>
    </p>
  </div>
</template>

<script>
import { listMatches } from '@/api/stats'

const MODE_MAP = { 0: '未知', 1: '单人', 2: '双人' }

export default {
  name: 'HistoryView',
  data() {
    return {
      records: [],
      page: 1,
      pageSize: 20,
      total: 0,
      loading: true,
      error: ''
    }
  },
  mounted() {
    this.loadPage(1)
  },
  methods: {
    modeLabel(mode) {
      return MODE_MAP[mode] ?? '未知'
    },
    formatTime(ts) {
      if (!ts) return '-'
      const d = new Date(ts * 1000)
      return d.toLocaleString('zh-CN')
    },
    async loadPage(p) {
      this.loading = true
      this.error = ''
      try {
        const res = await listMatches(p, this.pageSize)
        this.records = res.records || []
        this.page = res.page ?? p
        this.pageSize = res.page_size ?? this.pageSize
        this.total = res.total ?? 0
      } catch (e) {
        this.error = e.message || '加载失败'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.page-history {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
}

.page-subtitle {
  color: #606266;
  margin-bottom: 16px;
}

.loading,
.error-msg {
  margin: 16px 0;
}

.error-msg {
  color: #dc2626;
}

.panel {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.08);
  margin-bottom: 16px;
}

.table-wrap {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.table th,
.table td {
  padding: 8px 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.table th {
  font-weight: 600;
  color: #374151;
}

.session-id {
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.win {
  color: #059669;
  font-weight: 500;
}

.lose {
  color: #dc2626;
}

.empty {
  text-align: center;
  color: #6b7280;
  padding: 24px;
}

.pagination {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.page-info {
  color: #6b7280;
  font-size: 14px;
}

.btn {
  border-radius: 8px;
  border: 1px solid #d1d5db;
  padding: 6px 14px;
  background: #f9fafb;
  cursor: pointer;
  font-size: 14px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn:hover:not(:disabled) {
  background: #f3f4f6;
}

.btn:focus-visible {
  outline: 2px solid #4f46e5;
  outline-offset: 2px;
}

.nav-link a:focus-visible {
  outline: 2px solid #4f46e5;
  outline-offset: 2px;
  border-radius: 4px;
}

.nav-link {
  margin-top: 16px;
}

.nav-link a {
  color: #4f46e5;
  text-decoration: none;
}

.nav-link a:hover {
  text-decoration: underline;
}
</style>
