/**
 * 战绩与排行榜 API
 * GET /api/v1/stats/matches, GET /api/v1/stats/leaderboard
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
 * 分页查询历史战绩
 * GET /api/v1/stats/matches?page=1&page_size=20
 */
export function listMatches(page = 1, pageSize = 20) {
  const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) })
  return request(`/api/v1/stats/matches?${params}`)
}

/**
 * 获取排行榜
 * GET /api/v1/stats/leaderboard?limit=50
 */
export function getLeaderboard(limit = 50) {
  return request(`/api/v1/stats/leaderboard?limit=${limit}`)
}
