/**
 * 游戏 API：对局创建、结束、地图与坦克配置
 * 与后端 Kratos GameService 的 REST 接口对应
 */

const API_BASE_URL = 'http://127.0.0.1:8000'

async function request(url, options = {}) {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    ...options
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`API Error: ${response.status} ${errorText}`)
  }

  return response.json()
}

/**
 * 获取地图列表
 * GET /api/v1/game/maps
 */
export function listMaps() {
  return request('/api/v1/game/maps', { method: 'GET' })
}

/**
 * 获取单张地图详情
 * GET /api/v1/game/maps/{id}
 */
export function getMap(id) {
  return request(`/api/v1/game/maps/${id}`, { method: 'GET' })
}

/**
 * 根据难度获取坦克配置
 * GET /api/v1/game/config/tank?difficulty=1
 * @param {number} difficulty - 1=简单 2=普通 3=困难
 */
export function getTankConfig(difficulty = 2) {
  const params = new URLSearchParams({ difficulty: String(difficulty) })
  return request(`/api/v1/game/config/tank?${params}`, { method: 'GET' })
}

/**
 * 创建对局
 * POST /api/v1/game/session
 * @param {Object} body - { mode, map_id, local_player_count, difficulty }
 * @returns {Promise<{ session_id: string }>}
 */
export function createSession(body) {
  const payload = {
    mode: body.mode,
    map_id: body.map_id,
    local_player_count: body.local_player_count
  }
  if (body.difficulty != null) payload.difficulty = body.difficulty
  return request('/api/v1/game/session', {
    method: 'POST',
    body: JSON.stringify(payload)
  })
}

/**
 * 结束对局并上报战绩
 * PATCH /api/v1/game/session/{session_id}/finish
 * @param {string} sessionId - 对局 ID
 * @param {Object} body - { mode, map_id, is_win, player1_kills, player2_kills, enemy_kills, duration_seconds }
 */
export function finishSession(sessionId, body) {
  return request(`/api/v1/game/session/${encodeURIComponent(sessionId)}/finish`, {
    method: 'PATCH',
    body: JSON.stringify({
      session_id: sessionId,
      mode: body.mode,
      map_id: body.map_id,
      is_win: body.is_win,
      player1_kills: body.player1_kills ?? 0,
      player2_kills: body.player2_kills ?? 0,
      enemy_kills: body.enemy_kills ?? 0,
      duration_seconds: body.duration_seconds ?? 0
    })
  })
}
