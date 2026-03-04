<template>
  <div class="page page-leaderboard">
    <h1 class="page-title">排行榜</h1>
    <p class="page-subtitle">按胜场排序，展示总场次与胜率。</p>

    <div v-if="loading" class="loading">加载中…</div>
    <p v-else-if="error" class="error-msg">{{ error }}</p>

    <template v-else>
      <div class="panel table-wrap">
        <table class="table">
          <thead>
            <tr>
              <th>排名</th>
              <th>用户</th>
              <th>总场次</th>
              <th>胜场</th>
              <th>胜率</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(e, i) in entries" :key="e.user_id">
              <td class="rank">{{ i + 1 }}</td>
              <td class="username">{{ e.username }}</td>
              <td>{{ e.total_matches }}</td>
              <td>{{ e.total_wins }}</td>
              <td>{{ winRateStr(e.win_rate) }}</td>
            </tr>
          </tbody>
        </table>
        <p v-if="entries.length === 0" class="empty">暂无排行数据</p>
      </div>
    </template>

    <p class="nav-link">
      <router-link to="/">返回首页</router-link>
    </p>
  </div>
</template>

<script>
import { getLeaderboard } from '@/api/stats'

export default {
  name: 'LeaderboardView',
  data() {
    return {
      entries: [],
      loading: true,
      error: ''
    }
  },
  mounted() {
    this.fetchLeaderboard()
  },
  methods: {
    winRateStr(rate) {
      if (rate == null) return '-'
      return `${(rate * 100).toFixed(1)}%`
    },
    async fetchLeaderboard() {
      this.loading = true
      this.error = ''
      try {
        const res = await getLeaderboard(50)
        this.entries = res.entries || []
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
.page-leaderboard {
  max-width: 720px;
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
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.table th {
  font-weight: 600;
  color: #374151;
}

.rank {
  font-weight: 600;
  color: #4f46e5;
  width: 64px;
}

.username {
  font-weight: 500;
}

.empty {
  text-align: center;
  color: #6b7280;
  padding: 24px;
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
