/**
 * 成就 API
 */

import { api } from './api'

/** GET /api/v1/achievements */
export function listAchievements() {
  return api.get('/api/v1/achievements')
}

/** GET /api/v1/achievements/:id */
export function getAchievement(id) {
  return api.get(`/api/v1/achievements/${id}`)
}
