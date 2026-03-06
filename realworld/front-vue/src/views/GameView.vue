<template>
  <div class="page page-game">
    <header class="game-header">
      <h1 class="page-title">单人训练场</h1>
      <p class="page-subtitle">WASD 移动，空格开火（占位）。按 <kbd>E</kbd> 结束对局并上报战绩。</p>
    </header>

    <section v-if="result" class="result-panel" role="status" aria-live="polite">
      <div class="panel">
        <h2 class="panel-title">{{ result.success ? '对局已记录' : '上报失败' }}</h2>
        <p>{{ result.message }}</p>
        <button type="button" class="btn btn-primary" @click="backHome">返回首页</button>
      </div>
    </section>

    <section v-else class="game-layout">
      <div class="game-canvas-wrapper">
        <canvas ref="canvas" class="game-canvas" width="640" height="480" role="img" aria-label="坦克大战游戏画面，使用 WASD 移动、空格开火、E 键结束对局"></canvas>
      </div>
      <aside class="game-sidebar" aria-label="游戏状态与操作">
        <div class="panel">
          <h2 class="panel-title">当前状态</h2>
          <p>FPS: {{ fps }}</p>
          <p>坦克位置：({{ tankX }}, {{ tankY }})</p>
        </div>
        <button type="button" class="btn" @click="backHome">
          返回首页
        </button>
      </aside>
    </section>
  </div>
</template>

<script>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { createGameLoop } from '@/game/core/GameLoop'
import { finishSession, getMap, getTankConfig } from '@/api/game'

export default {
  name: 'GameView',
  setup() {
    const routerRef = ref(null)
    const canvas = ref(null)
    const fps = ref(0)
    const tankX = ref(0)
    const tankY = ref(0)
    const result = ref(null)
    const sessionId = ref('')
    const mode = ref(1)
    const mapId = ref(1)
    let loop

    const backHome = () => {
      routerRef.value?.push({ name: 'home' })
    }

    onMounted(async () => {
      const router = useRouter()
      const route = useRoute()
      routerRef.value = router

      const query = route?.query ?? (() => {
        const params = new URLSearchParams(window.location.search)
        return {
          sessionId: params.get('sessionId') ?? params.get('session_id'),
          mode: params.get('mode'),
          mapId: params.get('mapId'),
          difficulty: params.get('difficulty')
        }
      })()
      const sid = query.sessionId ?? query.session_id ?? ''
      const m = parseInt(query.mode, 10) || 1
      const mid = parseInt(query.mapId, 10) || 1
      const difficulty = parseInt(query.difficulty, 10) || 2
      sessionId.value = sid
      mode.value = m
      mapId.value = mid

      if (!sessionId.value) {
        router?.replace({ name: 'home' })
        return
      }
      if (!canvas.value) return

      let tankConfig = null
      let mapData = null
      const [cfgSettled, mapSettled] = await Promise.allSettled([
        getTankConfig(difficulty),
        getMap(mid)
      ])
      if (cfgSettled.status === 'fulfilled' && cfgSettled.value?.config) {
        tankConfig = cfgSettled.value.config
      }
      if (mapSettled.status === 'fulfilled' && mapSettled.value?.map) {
        mapData = mapSettled.value.map
      }

      const { start, stop } = createGameLoop(canvas.value, {
        tankConfig,
        map: mapData,
        onStats(update) {
          fps.value = update.fps
          tankX.value = Math.round(update.tank.x)
          tankY.value = Math.round(update.tank.y)
        },
        onGameEnd(stats) {
          stop()
          finishSession(sessionId.value, {
            mode: mode.value,
            map_id: mapId.value,
            is_win: stats.is_win,
            player1_kills: stats.player1_kills,
            player2_kills: stats.player2_kills,
            enemy_kills: stats.enemy_kills,
            duration_seconds: stats.duration_seconds
          })
            .then(() => {
              result.value = { success: true, message: '战绩已提交到服务器。' }
            })
            .catch((e) => {
              result.value = { success: false, message: e.message || '上报战绩失败' }
            })
        }
      })

      loop = { stop }
      start()
    })

    onBeforeUnmount(() => {
      if (loop) {
        loop.stop()
      }
    })

    return {
      canvas,
      fps,
      tankX,
      tankY,
      result,
      backHome
    }
  }
}
</script>

<style scoped>
.page-game {
  max-width: 1120px;
  margin: 0 auto;
  padding: 24px 16px 32px;
}

.game-header {
  margin-bottom: 16px;
}

.game-layout {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(220px, 1fr);
  gap: 16px;
  align-items: flex-start;
}

.game-canvas-wrapper {
  background: radial-gradient(circle at top, #0f172a, #020617);
  border-radius: 20px;
  padding: 12px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.7);
}

.game-canvas {
  display: block;
  width: 100%;
  border-radius: 16px;
  background: #020617;
}

.game-sidebar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.panel {
  background: rgba(15, 23, 42, 0.95);
  color: #e5e7eb;
  border-radius: 16px;
  padding: 12px 16px;
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.75);
}

.panel-title {
  font-size: 14px;
  margin-bottom: 8px;
}

.btn {
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.5);
  padding: 6px 14px;
  background: rgba(15, 23, 42, 0.9);
  font-size: 14px;
  cursor: pointer;
  color: #e5e7eb;
  transition: all 0.15s ease;
}

.btn:hover {
  background: rgba(30, 64, 175, 0.95);
  border-color: rgba(129, 140, 248, 0.9);
}

.result-panel {
  max-width: 480px;
  margin: 0 auto;
}

.result-panel .panel p {
  margin-bottom: 12px;
}

kbd {
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(15, 23, 42, 0.9);
  font-size: 12px;
}

.btn:focus-visible {
  outline: 2px solid #38bdf8;
  outline-offset: 2px;
}
</style>

