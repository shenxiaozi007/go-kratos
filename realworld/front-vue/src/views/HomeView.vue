<template>
  <div class="page page-home">
    <h1 class="page-title">坦克大战</h1>
    <p class="page-subtitle">选择模式、难度与地图，开始一局战斗。</p>

    <section class="panel">
      <h2 class="panel-title">模式</h2>
      <div class="button-group">
        <button type="button" class="btn" :class="{ 'btn-primary': mode === 'single' }" @click="mode = 'single'">
          单人模式
        </button>
        <button type="button" class="btn" disabled>
          本地双人（待实现）
        </button>
      </div>
    </section>

    <section class="panel">
      <h2 class="panel-title">难度</h2>
      <div class="button-group">
        <button type="button" class="btn" :class="{ 'btn-primary': difficulty === 1 }" @click="difficulty = 1">
          简单
        </button>
        <button type="button" class="btn" :class="{ 'btn-primary': difficulty === 2 }" @click="difficulty = 2">
          普通
        </button>
        <button type="button" class="btn" :class="{ 'btn-primary': difficulty === 3 }" @click="difficulty = 3">
          困难
        </button>
      </div>
    </section>

    <section class="panel">
      <h2 class="panel-title">地图</h2>
      <div v-if="mapsLoading" class="loading-inline">加载地图中…</div>
      <div v-else class="button-group">
        <button
          v-for="m in maps"
          :key="m.id"
          type="button"
          class="btn"
          :class="{ 'btn-primary': mapId === m.id }"
          @click="mapId = m.id"
        >
          {{ m.name }}
        </button>
      </div>
      <p v-if="maps.length === 0 && !mapsLoading" class="hint">暂无地图，请先启动后端并执行建表与种子。</p>
    </section>

    <section class="panel actions">
      <button type="button" class="btn btn-primary btn-lg" :disabled="loading" @click="goGame">
        {{ loading ? '创建对局中…' : '进入游戏' }}
      </button>
      <p v-if="error" class="error-msg">{{ error }}</p>
    </section>
  </div>
</template>

<script>
import { createSession, listMaps } from '@/api/game'

export default {
  name: 'HomeView',
  data() {
    return {
      mode: 'single',
      difficulty: 2,
      mapId: 1,
      maps: [],
      mapsLoading: false,
      loading: false,
      error: ''
    }
  },
  mounted() {
    this.loadMaps()
  },
  methods: {
    async loadMaps() {
      this.mapsLoading = true
      try {
        const res = await listMaps()
        this.maps = res.maps || []
        if (this.maps.length > 0 && !this.maps.find(m => m.id === this.mapId)) {
          this.mapId = this.maps[0].id
        }
      } catch (_) {
        this.maps = []
      } finally {
        this.mapsLoading = false
      }
    },
    async goGame() {
      if (this.loading) return
      this.error = ''
      this.loading = true
      try {
        const mode = this.mode === 'single' ? 1 : 2
        const localPlayerCount = this.mode === 'single' ? 1 : 2
        const res = await createSession({
          mode,
          map_id: this.mapId,
          local_player_count: localPlayerCount,
          difficulty: this.difficulty
        })
        const sessionId = res.session_id
        if (!sessionId) {
          this.error = '未获取到对局 ID'
          return
        }
        this.$router.push({
          name: 'game',
          query: {
            sessionId,
            mode: String(mode),
            mapId: String(this.mapId),
            difficulty: String(this.difficulty)
          }
        })
      } catch (e) {
        this.error = e.message || '创建对局失败'
        console.error(e)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 16px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
}

.page-subtitle {
  color: #606266;
  margin-bottom: 24px;
}

.panel {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 16px 20px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.08);
  margin-bottom: 16px;
}

.panel-title {
  font-size: 16px;
  margin-bottom: 12px;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.btn {
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.5);
  padding: 6px 14px;
  background: rgba(248, 250, 252, 0.8);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn:hover:not(:disabled) {
  background: #fff;
  border-color: #6366f1;
  color: #111827;
}

.btn-primary {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  border-color: transparent;
}

.btn-primary:hover:not(:disabled) {
  filter: brightness(1.05);
}

.btn-lg {
  padding: 10px 24px;
  font-size: 16px;
}

.actions {
  text-align: right;
}

.error-msg {
  color: #dc2626;
  font-size: 14px;
  margin-top: 8px;
}

.loading-inline,
.hint {
  color: #6b7280;
  font-size: 14px;
}
</style>

