/**
 * 用户资料 API
 */

import { api } from './api'

/** GET /api/v1/user/profile */
export function getProfile() {
  return api.get('/api/v1/user/profile')
}

/** PUT /api/v1/user/profile */
export function updateProfile(body) {
  return api.put('/api/v1/user/profile', body)
}
